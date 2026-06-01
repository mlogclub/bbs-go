"""Shell execution tool."""

import asyncio
import os
import re
import signal
import shutil
import sys
from contextlib import suppress
from pathlib import Path
from typing import Any, Literal

from loguru import logger

from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.tools.limits import ToolLimits
from OriginAgent.agent.tools.sandbox import wrap_command
from OriginAgent.agent.tools.schema import IntegerSchema, StringSchema, tool_parameters_schema
from OriginAgent.config.paths import get_media_dir
from OriginAgent.security.paths import ProtectedPathPolicy
from OriginAgent.security.policy import PolicyDeniedError

_IS_WINDOWS = sys.platform == "win32"


# Policy note appended to recoverable workspace-boundary guard errors.
_WORKSPACE_BOUNDARY_NOTE = (
    "\n\nNote: this is a hard policy boundary, not a transient failure. "
    "Do NOT retry with shell tricks (symlinks, base64 piping, alternative "
    "tools, working_dir overrides). If the user genuinely needs this "
    "resource, tell them you cannot reach it under the current "
    "restrict_to_workspace policy and ask how to proceed."
)
_UNSAFE_EXEC_MARKER = (
    "[unsafe-exec profile=local_dev sandbox=none] "
    "local_dev unsafe exec does not provide sandbox isolation."
)


class ExecTool(Tool):
    """Tool to execute shell commands."""

    def __init__(
        self,
        timeout: int = 60,
        working_dir: str | None = None,
        deny_patterns: list[str] | None = None,
        allow_patterns: list[str] | None = None,
        restrict_to_workspace: bool = False,
        sandbox: str = "",
        path_append: str = "",
        allowed_env_keys: list[str] | None = None,
        limits: ToolLimits | None = None,
        protected_policy: ProtectedPathPolicy | None = None,
        security_profile: Literal["secure", "local_dev", "disabled"] = "secure",
        allow_unsafe_exec: bool = False,
        shell_syntax_policy: Literal["restricted", "shell"] = "restricted",
    ):
        self._limits = limits or ToolLimits()
        self.timeout = timeout
        self.working_dir = working_dir
        self.sandbox = sandbox
        self.deny_patterns = (deny_patterns or []) + [
            r"\brm\s+-[rf]{1,2}\b",          # rm -r, rm -rf, rm -fr
            r"\bdel\s+/[fq]\b",              # del /f, del /q
            r"\brmdir\s+/s\b",               # rmdir /s
            r"(?:^|[;&|]\s*)format\b",       # format (as standalone command only)
            r"\b(mkfs|diskpart)\b",          # disk operations
            r"\bdd\s+if=",                   # dd
            r">\s*/dev/sd",                  # write to disk
            r"\b(shutdown|reboot|poweroff)\b",  # system power
            r":\(\)\s*\{.*\};\s*:",          # fork bomb
            # Block writes to OriginAgent internal state files (#2989).
            # history.jsonl / .dream_cursor are managed by append_history();
            # direct writes corrupt the cursor format and crash /dream.
            r">>?\s*\S*(?:history\.jsonl|\.dream_cursor)",            # > / >> redirect
            r"\btee\b[^|;&<>]*(?:history\.jsonl|\.dream_cursor)",     # tee / tee -a
            r"\b(?:cp|mv)\b(?:\s+[^\s|;&<>]+)+\s+\S*(?:history\.jsonl|\.dream_cursor)",  # cp/mv target
            r"\bdd\b[^|;&<>]*\bof=\S*(?:history\.jsonl|\.dream_cursor)",  # dd of=
            r"\bsed\s+-i[^|;&<>]*(?:history\.jsonl|\.dream_cursor)",  # sed -i
        ]
        self.allow_patterns = allow_patterns or []
        self.restrict_to_workspace = restrict_to_workspace
        self.path_append = path_append
        self.allowed_env_keys = allowed_env_keys or []
        self.security_profile = security_profile
        self.allow_unsafe_exec = allow_unsafe_exec
        self.shell_syntax_policy = (
            shell_syntax_policy
            if security_profile == "local_dev" and allow_unsafe_exec
            else "restricted"
        )
        workspace_for_policy = Path(working_dir) if working_dir else None
        self._protected_policy = protected_policy or ProtectedPathPolicy(workspace_for_policy)

    @property
    def name(self) -> str:
        return "exec"

    _MAX_TIMEOUT = 600
    _MAX_OUTPUT = 10_000

    # Kernel device files safe as stdio redirect targets (#3599).
    _BENIGN_DEVICE_PATHS: frozenset[str] = frozenset({
        "/dev/null",
        "/dev/zero",
        "/dev/full",
        "/dev/random",
        "/dev/urandom",
        "/dev/stdin",
        "/dev/stdout",
        "/dev/stderr",
        "/dev/tty",
    })

    @property
    def parameters(self) -> dict[str, Any]:
        max_timeout = self._limits.exec_max_timeout_seconds
        return tool_parameters_schema(
            command=StringSchema("The shell command to execute"),
            working_dir=StringSchema("Optional working directory for the command"),
            timeout=IntegerSchema(
                60,
                description=(
                    "Timeout in seconds. Increase for long-running commands "
                    f"like compilation or installation (default 60, max {max_timeout})."
                ),
                minimum=1,
                maximum=max_timeout,
            ),
            required=["command"],
        )

    @property
    def description(self) -> str:
        max_output = self._limits.exec_max_output_chars
        return (
            "Execute a shell command and return its output. "
            "Prefer read_file/write_file/edit_file over cat/echo/sed, "
            "and grep/glob over shell find/grep. "
            "Use -y or --yes flags to avoid interactive prompts. "
            f"Output is truncated at {max_output:,} chars; timeout defaults to 60s."
        )

    @property
    def exclusive(self) -> bool:
        return True

    async def execute(
        self, command: str, working_dir: str | None = None,
        timeout: int | None = None, **kwargs: Any,
    ) -> str:
        cwd = working_dir or self.working_dir or os.getcwd()

        # Prevent an LLM-supplied working_dir from escaping the configured
        # workspace when restrict_to_workspace is enabled (#2826). Without
        # this, a caller can pass working_dir="/etc" and then all absolute
        # paths under /etc would pass the _guard_command check that anchors
        # on cwd.
        if self.working_dir:
            try:
                requested = Path(cwd).expanduser().resolve()
                workspace_root = Path(self.working_dir).expanduser().resolve()
            except Exception:
                return (
                    "Error: working_dir could not be resolved"
                    + _WORKSPACE_BOUNDARY_NOTE
                )
            if requested != workspace_root and workspace_root not in requested.parents:
                return (
                    "Error: working_dir is outside the configured workspace"
                    + _WORKSPACE_BOUNDARY_NOTE
                )

        try:
            unsafe_exec = self._assert_sandbox_available()
            guard_error = self._guard_command(command, cwd)
            if guard_error:
                return guard_error

            if self.sandbox and not unsafe_exec:
                if _IS_WINDOWS:
                    raise PolicyDeniedError(
                        f"sandbox '{self.sandbox}' is not supported on Windows; command was not executed",
                        code="sandbox_unsupported",
                        boundary="exec",
                        policy_rule="windows_sandbox_unsupported",
                    )
                workspace = self.working_dir or cwd
                command = wrap_command(self.sandbox, command, workspace, cwd)
                cwd = str(Path(workspace).resolve())
        except PolicyDeniedError as exc:
            return f"Error: {exc}"

        effective_timeout = min(timeout or self.timeout, self._limits.exec_max_timeout_seconds)
        env = self._build_env()

        if self.path_append:
            if _IS_WINDOWS:
                env["PATH"] = env.get("PATH", "") + os.pathsep + self.path_append
            else:
                env["ORIGINAGENT_PATH_APPEND"] = self.path_append
                command = f'export PATH="$PATH{os.pathsep}$ORIGINAGENT_PATH_APPEND"; {command}'

        try:
            process = await self._spawn(command, cwd, env)

            try:
                stdout, stderr = await asyncio.wait_for(
                    process.communicate(),
                    timeout=effective_timeout,
                )
            except asyncio.TimeoutError:
                await self._kill_process(process)
                return f"Error: Command timed out after {effective_timeout} seconds"
            except asyncio.CancelledError:
                await self._kill_process(process)
                raise

            output_parts = []

            if stdout:
                output_parts.append(stdout.decode("utf-8", errors="replace"))

            if stderr:
                stderr_text = stderr.decode("utf-8", errors="replace")
                if stderr_text.strip():
                    output_parts.append(f"STDERR:\n{stderr_text}")

            output_parts.append(f"\nExit code: {process.returncode}")

            result = "\n".join(output_parts) if output_parts else "(no output)"
            if unsafe_exec:
                result = f"{_UNSAFE_EXEC_MARKER}\n{result}"

            max_len = self._limits.exec_max_output_chars
            if len(result) > max_len:
                half = max_len // 2
                result = (
                    result[:half]
                    + f"\n\n... ({len(result) - max_len:,} chars truncated) ...\n\n"
                    + result[-half:]
                )

            return result

        except Exception as e:
            return f"Error executing command: {str(e)}"

    @staticmethod
    async def _spawn(
        command: str, cwd: str, env: dict[str, str],
    ) -> asyncio.subprocess.Process:
        """Launch *command* in a platform-appropriate shell."""
        if _IS_WINDOWS:
            # create_subprocess_exec re-quotes args via list2cmdline, which
            # breaks commands containing paths with spaces (e.g. "D:\Program
            # Files\python.exe" "script.py"). create_subprocess_shell passes
            # the raw command string to COMSPEC without re-quoting.
            return await asyncio.create_subprocess_shell(
                command,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
                cwd=cwd,
                env=env,
            )
        bash = shutil.which("bash") or "/bin/bash"
        return await asyncio.create_subprocess_exec(
            bash, "-l", "-c", command,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd=cwd,
            env=env,
            start_new_session=True,
        )

    @staticmethod
    async def _kill_process(process: asyncio.subprocess.Process) -> None:
        """Kill a subprocess and reap it to prevent zombies."""
        if _IS_WINDOWS:
            process.kill()
        else:
            try:
                os.killpg(os.getpgid(process.pid), signal.SIGKILL)
            except ProcessLookupError:
                pass
            except Exception:
                process.kill()
        try:
            with suppress(asyncio.TimeoutError):
                await asyncio.wait_for(process.wait(), timeout=5.0)
        finally:
            if not _IS_WINDOWS:
                try:
                    os.waitpid(process.pid, os.WNOHANG)
                except (ProcessLookupError, ChildProcessError) as e:
                    logger.debug("Process already reaped or not found: {}", e)

    def _build_env(self) -> dict[str, str]:
        """Build a minimal environment for subprocess execution.

        On Unix, only HOME/LANG/TERM are passed; ``bash -l`` sources the
        user's profile which sets PATH and other essentials.

        On Windows, ``cmd.exe`` has no login-profile mechanism, so a curated
        set of system variables (including PATH) is forwarded.  API keys and
        other secrets are still excluded.
        """
        if _IS_WINDOWS:
            sr = os.environ.get("SYSTEMROOT", r"C:\Windows")
            env = {
                "SYSTEMROOT": sr,
                "COMSPEC": os.environ.get("COMSPEC", f"{sr}\\system32\\cmd.exe"),
                "USERPROFILE": os.environ.get("USERPROFILE", ""),
                "HOMEDRIVE": os.environ.get("HOMEDRIVE", "C:"),
                "HOMEPATH": os.environ.get("HOMEPATH", "\\"),
                "TEMP": os.environ.get("TEMP", f"{sr}\\Temp"),
                "TMP": os.environ.get("TMP", f"{sr}\\Temp"),
                "PATHEXT": os.environ.get("PATHEXT", ".COM;.EXE;.BAT;.CMD"),
                "PATH": os.environ.get("PATH", f"{sr}\\system32;{sr}"),
                "APPDATA": os.environ.get("APPDATA", ""),
                "LOCALAPPDATA": os.environ.get("LOCALAPPDATA", ""),
                "ProgramData": os.environ.get("ProgramData", ""),
                "ProgramFiles": os.environ.get("ProgramFiles", ""),
                "ProgramFiles(x86)": os.environ.get("ProgramFiles(x86)", ""),
                "ProgramW6432": os.environ.get("ProgramW6432", ""),
            }
            for key in self.allowed_env_keys:
                val = os.environ.get(key)
                if val is not None:
                    env[key] = val
            return env
        home = os.environ.get("HOME", "/tmp")
        env = {
            "HOME": home,
            "LANG": os.environ.get("LANG", "C.UTF-8"),
            "TERM": os.environ.get("TERM", "dumb"),
        }
        for key in self.allowed_env_keys:
            val = os.environ.get(key)
            if val is not None:
                env[key] = val
        return env

    def _guard_command(self, command: str, cwd: str) -> str | None:
        """Best-effort safety guard for potentially destructive commands."""
        cmd = command.strip()
        lower = cmd.lower()

        explicitly_allowed = bool(self.allow_patterns) and any(
            re.search(p, lower) for p in self.allow_patterns
        )
        for pattern in self.deny_patterns:
            if re.search(pattern, lower):
                return "Error: Command blocked by deny pattern filter"

        if self.shell_syntax_policy == "restricted" and self._has_shell_control_syntax(cmd):
            return "Error: Command blocked by shell syntax policy"

        if self.allow_patterns and not explicitly_allowed:
            return "Error: Command blocked by allowlist filter (not in allowlist)"

        from OriginAgent.security.network import contains_internal_url
        if contains_internal_url(cmd):
            # The runner turns this marker into a non-retryable security hint.
            return "Error: Command blocked by safety guard (internal/private URL detected)"

        for raw in self._extract_absolute_paths(cmd):
            try:
                expanded = os.path.expandvars(raw.strip())
                if self._is_benign_device_path(expanded):
                    continue
                p = Path(expanded).expanduser().resolve()
            except Exception:
                continue
            if self._is_benign_device_path(str(p)):
                continue
            try:
                self._protected_policy.assert_can_write(Path(str(p)))
            except (PermissionError, TypeError, ValueError) as exc:
                if isinstance(exc, PermissionError):
                    return f"Error: {exc}"

        if self.restrict_to_workspace:
            if "..\\" in cmd or "../" in cmd:
                return (
                    "Error: Command blocked by safety guard (path traversal detected)"
                    + _WORKSPACE_BOUNDARY_NOTE
                )

            cwd_path = Path(cwd).resolve()

            for raw in self._extract_absolute_paths(cmd):
                try:
                    expanded = os.path.expandvars(raw.strip())
                    # Match against the un-resolved path first.  On Linux,
                    # /dev/stderr is a symlink to /proc/self/fd/2 and
                    # ``Path.resolve()`` would mask the device-file intent.
                    if self._is_benign_device_path(expanded):
                        continue
                    p = Path(expanded).expanduser().resolve()
                except Exception:
                    continue

                if self._is_benign_device_path(str(p)):
                    continue
                media_path = get_media_dir().resolve()
                if (p.is_absolute()
                    and cwd_path not in p.parents
                    and p != cwd_path
                    and media_path not in p.parents
                    and p != media_path
                ):
                    return (
                        "Error: Command blocked by safety guard (path outside working dir)"
                        + _WORKSPACE_BOUNDARY_NOTE
                    )

        return None

    def _assert_sandbox_available(self) -> bool:
        if self.security_profile == "disabled":
            raise PolicyDeniedError(
                "exec profile is disabled; command was not executed",
                code="exec_profile_disabled",
                boundary="exec",
                policy_rule="exec_profile_disabled",
            )
        if self.security_profile not in {"secure", "local_dev"}:
            raise PolicyDeniedError(
                f"exec profile '{self.security_profile}' is not supported; command was not executed",
                code="exec_profile_invalid",
                boundary="exec",
                policy_rule="exec_profile_invalid",
            )
        sandbox = (self.sandbox or "").strip()
        if _IS_WINDOWS and sandbox and sandbox != "none":
            if self.security_profile == "local_dev" and self.allow_unsafe_exec:
                return True
            raise PolicyDeniedError(
                f"sandbox '{sandbox}' is not supported on Windows; command was not executed",
                code="sandbox_unsupported",
                boundary="exec",
                policy_rule="windows_sandbox_unsupported",
            )
        if self.working_dir and not self.restrict_to_workspace:
            if (
                self.security_profile == "local_dev"
                and self.allow_unsafe_exec
                and sandbox == "bwrap"
                and shutil.which("bwrap") is not None
            ):
                return False
            if (
                self.security_profile == "local_dev"
                and self.allow_unsafe_exec
                and (not sandbox or sandbox == "none" or sandbox == "bwrap")
            ):
                return True
            raise PolicyDeniedError(
                "workspace exec requires restrict_to_workspace with a supported sandbox; command was not executed",
                code="sandbox_required",
                boundary="exec",
                policy_rule="sandbox_required",
            )
        if not self.restrict_to_workspace:
            return False
        if not sandbox or sandbox == "none":
            if self.security_profile == "local_dev" and self.allow_unsafe_exec:
                return True
            raise PolicyDeniedError(
                "restrict_to_workspace requires a supported sandbox; command was not executed",
                code="sandbox_required",
                boundary="exec",
                policy_rule="sandbox_required",
            )
        if sandbox != "bwrap":
            raise PolicyDeniedError(
                f"sandbox '{sandbox}' is not an approved workspace sandbox",
                code="sandbox_unapproved",
                boundary="exec",
                policy_rule="sandbox_required",
            )
        if shutil.which("bwrap") is None:
            if self.security_profile == "local_dev" and self.allow_unsafe_exec:
                return True
            raise PolicyDeniedError(
                "sandbox 'bwrap' is configured but not available; command was not executed",
                code="sandbox_missing",
                boundary="exec",
                policy_rule="sandbox_backend_missing",
            )
        return False

    @classmethod
    def _is_benign_device_path(cls, path: str) -> bool:
        """Return True for kernel device files that should never be workspace-blocked."""
        if path in cls._BENIGN_DEVICE_PATHS:
            return True
        return path.startswith("/dev/fd/")

    @staticmethod
    def _extract_absolute_paths(command: str) -> list[str]:
        # Windows: match drive-root paths like `C:\` as well as `C:\path\to\file`
        # NOTE: `*` is required so `C:\` (nothing after the slash) is still extracted.
        win_paths = re.findall(r"[A-Za-z]:\\[^\s\"'|><;]*", command)
        posix_paths = re.findall(r"(?:^|[\s|>'\"])(/[^\s\"'>;|<]+)", command) # POSIX: /absolute only
        home_paths = re.findall(r"(?:^|[\s>'\"])(~[^\s\"'>;|<]*)", command) # POSIX/Windows home shortcut: ~
        return win_paths + posix_paths + home_paths

    @staticmethod
    def _has_shell_control_syntax(command: str) -> bool:
        """Detect shell control syntax outside single-quoted literals.

        Double-quoted semicolons and pipes are ordinary argument text, but
        command substitutions still execute inside double quotes.
        """
        in_single = False
        in_double = False
        escaped = False
        i = 0
        while i < len(command):
            ch = command[i]
            if escaped:
                escaped = False
                i += 1
                continue
            if ch == "\\":
                escaped = True
                i += 1
                continue
            if ch == "'" and not in_double:
                in_single = not in_single
                i += 1
                continue
            if ch == '"' and not in_single:
                in_double = not in_double
                i += 1
                continue
            if in_single:
                i += 1
                continue
            if ch == "`":
                return True
            if ch == "$" and i + 1 < len(command) and command[i + 1] == "(":
                return True
            if not in_double and ch in {";", "|", "&", "\n", "\r"}:
                return True
            i += 1
        return False
