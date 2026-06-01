"""Read-only trial runner for governed workflow proposals."""

from __future__ import annotations

import fnmatch
import tempfile
from dataclasses import asdict, dataclass, field
from pathlib import Path
from typing import Any

from OriginAgent.agent.evolution_sandbox import SandboxEvaluator
from OriginAgent.agent.evolution_trial_logs import EvolutionTrialLogStore
from OriginAgent.utils.helpers import ensure_dir, truncate_text

_DEFAULT_READ_ONLY_TOOLS = {"read_file", "glob", "grep"}
_MAX_STEP_OUTPUT_CHARS = 2000


@dataclass(frozen=True)
class TrialRunResult:
    """Serializable result for a read-only workflow trial run."""

    status: str
    mode: str = "trial"
    read_only: bool = True
    isolated_workspace: bool = True
    gate_status: str = ""
    step_results: list[dict[str, Any]] = field(default_factory=list)
    log_id: str = ""
    summary: dict[str, Any] = field(default_factory=dict)
    policy: dict[str, Any] = field(default_factory=dict)

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class TrialRunner:
    """Run a read-only workflow trial inside an isolated temporary workspace."""

    def __init__(self, workspace: Path, config: Any | None = None) -> None:
        self.workspace = Path(workspace)
        self.config = config
        self.sandbox = SandboxEvaluator(self.workspace, config)
        self.logs = EvolutionTrialLogStore(self.workspace)

    def run_workflow_payload(
        self,
        payload: dict[str, Any],
        *,
        fixtures: dict[str, str] | None = None,
    ) -> dict[str, Any]:
        """Evaluate and replay a workflow payload with read-only fixture data."""

        gate = self.sandbox.evaluate_trial_workflow_payload(payload)
        policy = gate.get("policy") if isinstance(gate.get("policy"), dict) else {}
        if str(gate.get("status") or "") != "passed":
            log = self._append_log(
                payload,
                status=str(gate.get("status") or "blocked"),
                step_logs=_gate_step_logs(gate),
                summary={
                    "gate_status": str(gate.get("status") or ""),
                    "steps_checked": _safe_int(_mapping(gate.get("replay_summary")).get("steps_checked"), 0),
                    "blocked_steps": _safe_int(_mapping(gate.get("replay_summary")).get("blocked_steps"), 0),
                    "executed_steps": 0,
                },
                metadata={"trial_gate": gate},
            )
            return TrialRunResult(
                status=str(gate.get("status") or "blocked"),
                gate_status=str(gate.get("status") or ""),
                step_results=_gate_step_logs(gate),
                log_id=str(log.get("trial_id") or ""),
                summary=dict(log.get("summary") or {}),
                policy=policy,
            ).to_json()

        steps = payload.get("steps") if isinstance(payload.get("steps"), list) else []
        with self._trial_temp_directory() as tmp:
            root = Path(tmp).resolve()
            self._seed_fixtures(root, fixtures or {})
            step_results = [
                self._run_step(step, index=index, root=root)
                for index, step in enumerate(steps)
            ]

        status = _overall_status(step_results)
        summary = {
            "gate_status": "passed",
            "steps_checked": len([step for step in step_results if step.get("status") != "invalid"]),
            "blocked_steps": len([step for step in step_results if step.get("status") == "blocked"]),
            "failed_steps": len([step for step in step_results if step.get("status") == "failed"]),
            "executed_steps": len([step for step in step_results if step.get("executed") is True]),
        }
        log = self._append_log(
            payload,
            status=status,
            step_logs=step_results,
            summary=summary,
            metadata={"trial_gate": gate},
        )
        return TrialRunResult(
            status=status,
            gate_status="passed",
            step_results=step_results,
            log_id=str(log.get("trial_id") or ""),
            summary=dict(log.get("summary") or summary),
            policy=policy,
        ).to_json()

    def _run_step(self, step: Any, *, index: int, root: Path) -> dict[str, Any]:
        if not isinstance(step, dict):
            return _step_result(
                index=index,
                step={},
                status="invalid",
                output="",
                executed=False,
                issues=[{"code": "trial_step_not_mapping", "severity": "reject", "message": "Workflow step must be a mapping."}],
            )
        tool = str(step.get("tool") or "").strip()
        allowed_tools = self._allowed_tools()
        if not tool:
            return _step_result(index=index, step=step, status="skipped", output="", executed=False)
        if tool not in allowed_tools:
            return _step_result(
                index=index,
                step=step,
                status="blocked",
                output=f"[blocked by trial runner] `{tool}` is not an allowed read-only trial tool.",
                executed=False,
                issues=[{"code": "trial_tool_not_read_only", "severity": "pending", "message": f"Tool `{tool}` is not allowed in trial."}],
            )
        try:
            if tool == "read_file":
                output = self._read_file(root, step)
            elif tool == "glob":
                output = self._glob(root, step)
            elif tool == "grep":
                output = self._grep(root, step)
            else:
                output = f"[blocked by trial runner] `{tool}` has no read-only trial implementation."
                return _step_result(
                    index=index,
                    step=step,
                    status="blocked",
                    output=output,
                    executed=False,
                    issues=[{"code": "trial_tool_unimplemented", "severity": "pending", "message": output}],
                )
        except Exception as exc:
            return _step_result(
                index=index,
                step=step,
                status="failed",
                output=f"{type(exc).__name__}: {exc}",
                executed=True,
                issues=[{"code": "trial_step_failed", "severity": "reject", "message": truncate_text(str(exc), 300)}],
            )
        return _step_result(index=index, step=step, status="passed", output=output, executed=True)

    def _read_file(self, root: Path, step: dict[str, Any]) -> str:
        path = _step_path(step)
        target = _resolve_inside(root, path)
        if not target.is_file():
            raise FileNotFoundError(path)
        return target.read_text(encoding="utf-8", errors="replace")

    def _glob(self, root: Path, step: dict[str, Any]) -> str:
        pattern = str(step.get("pattern") or step.get("path") or "*").strip() or "*"
        if _unsafe_path(pattern):
            raise ValueError("glob pattern leaves trial workspace")
        matches = [
            str(path.relative_to(root)).replace("\\", "/")
            for path in root.rglob("*")
            if path.is_file() and fnmatch.fnmatch(str(path.relative_to(root)).replace("\\", "/"), pattern)
        ]
        return "\n".join(sorted(matches))

    def _grep(self, root: Path, step: dict[str, Any]) -> str:
        pattern = str(step.get("pattern") or "").strip()
        if not pattern:
            raise ValueError("grep pattern is required")
        path_value = str(step.get("path") or step.get("file") or ".").strip() or "."
        target = _resolve_inside(root, path_value)
        files = [target] if target.is_file() else [path for path in target.rglob("*") if path.is_file()]
        lines: list[str] = []
        needle = pattern.casefold()
        for path in sorted(files):
            rel = str(path.relative_to(root)).replace("\\", "/")
            content = path.read_text(encoding="utf-8", errors="replace")
            for number, line in enumerate(content.splitlines(), start=1):
                if needle in line.casefold():
                    lines.append(f"{rel}:{number}:{line}")
        return "\n".join(lines)

    def _seed_fixtures(self, root: Path, fixtures: dict[str, str]) -> None:
        for raw_path, content in fixtures.items():
            target = _resolve_inside(root, str(raw_path))
            ensure_dir(target.parent)
            target.write_text(str(content), encoding="utf-8")

    def _trial_temp_directory(self) -> tempfile.TemporaryDirectory[str]:
        trial_config = getattr(self.config, "trial", None)
        configured = str(getattr(trial_config, "temp_dir", "") or "").strip()
        if configured:
            base = ensure_dir(Path(configured))
            return tempfile.TemporaryDirectory(prefix="originagent_trial_run_", dir=str(base))
        return tempfile.TemporaryDirectory(prefix="originagent_trial_run_")

    def _allowed_tools(self) -> set[str]:
        sandbox_config = getattr(self.config, "sandbox", None)
        raw_tools = getattr(sandbox_config, "read_only_tools", None) or []
        configured = {str(item).strip() for item in raw_tools if str(item).strip()}
        return (configured or set(_DEFAULT_READ_ONLY_TOOLS)) & set(_DEFAULT_READ_ONLY_TOOLS)

    def _append_log(
        self,
        payload: dict[str, Any],
        *,
        status: str,
        step_logs: list[dict[str, Any]],
        summary: dict[str, Any],
        metadata: dict[str, Any],
    ) -> dict[str, Any]:
        evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
        trial_config = getattr(self.config, "trial", None)
        return self.logs.append_log(
            opportunity_id=str(evolution.get("opportunity_id") or ""),
            proposal_id=str(payload.get("review_proposal_id") or ""),
            artifact_type="workflow",
            artifact_name=str(payload.get("workflow_name") or payload.get("subject_id") or ""),
            artifact_path=str(payload.get("subject_path") or ""),
            status=status,
            step_logs=step_logs,
            summary=summary,
            metadata=metadata,
            max_step_output_chars=_safe_int(
                getattr(trial_config, "max_step_output_chars", _MAX_STEP_OUTPUT_CHARS)
                if trial_config is not None
                else _MAX_STEP_OUTPUT_CHARS,
                _MAX_STEP_OUTPUT_CHARS,
            ),
        )


def _gate_step_logs(gate: dict[str, Any]) -> list[dict[str, Any]]:
    step_results = gate.get("step_results") if isinstance(gate.get("step_results"), list) else []
    logs: list[dict[str, Any]] = []
    for step in step_results:
        if not isinstance(step, dict):
            continue
        logs.append({
            "index": _safe_int(step.get("index"), 0),
            "title": str(step.get("title") or ""),
            "tool": str(step.get("tool") or ""),
            "status": str(step.get("status") or "unknown"),
            "output": "",
            "executed": False,
            "read_only": True,
            "isolated_workspace": True,
            "issues": step.get("issues") if isinstance(step.get("issues"), list) else [],
            "metadata": {"executed": False, "source": "trial_gate"},
        })
    return logs


def _step_result(
    *,
    index: int,
    step: dict[str, Any],
    status: str,
    output: str,
    executed: bool,
    issues: list[dict[str, Any]] | None = None,
) -> dict[str, Any]:
    return {
        "index": index + 1,
        "title": str(step.get("title") or "")[:120],
        "tool": str(step.get("tool") or "").strip(),
        "status": status,
        "output": output,
        "executed": executed,
        "read_only": True,
        "isolated_workspace": True,
        "issues": issues or [],
    }


def _overall_status(step_results: list[dict[str, Any]]) -> str:
    if any(str(step.get("status") or "") == "failed" for step in step_results):
        return "failed"
    if any(str(step.get("status") or "") == "blocked" for step in step_results):
        return "blocked"
    return "passed"


def _step_path(step: dict[str, Any]) -> str:
    return str(step.get("path") or step.get("file") or step.get("file_path") or "").strip()


def _resolve_inside(root: Path, value: str) -> Path:
    if not value:
        raise ValueError("path is required")
    target = (root / value).resolve()
    target.relative_to(root)
    return target


def _unsafe_path(value: str) -> bool:
    try:
        Path(value).resolve().relative_to(Path(".").resolve())
    except Exception:
        pass
    return value.startswith("/") or "\\" in value and ":" in value or ".." in Path(value).parts


def _mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def _safe_int(value: Any, default: int) -> int:
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default
