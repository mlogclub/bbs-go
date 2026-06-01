"""Shared path protection policy for runtime state and workspace tools."""

from __future__ import annotations

from pathlib import Path

from OriginAgent.config.paths import get_cron_dir, get_data_dir, get_media_dir
from OriginAgent.security.policy import PolicyDeniedError

_NOTE = (
    " (this is a hard policy boundary, not a transient failure; "
    "do not retry with shell tricks or alternative tools)"
)


def _resolve_loose(path: Path) -> Path:
    return path.expanduser().resolve(strict=False)


def _is_under(path: Path, root: Path) -> bool:
    try:
        path.relative_to(root)
        return True
    except ValueError:
        return False


class ProtectedPathPolicy:
    """Central read/list/write policy for OriginAgent runtime state paths."""

    def __init__(self, workspace: Path | None = None):
        self.workspace = _resolve_loose(Path(workspace)) if workspace is not None else None
        try:
            data_dir = _resolve_loose(get_data_dir())
        except Exception:
            data_dir = _resolve_loose(Path.cwd() / ".originagent")
        try:
            cron_dir = _resolve_loose(get_cron_dir())
        except Exception:
            cron_dir = data_dir / "cron"
        try:
            media_dir = _resolve_loose(get_media_dir())
        except Exception:
            media_dir = data_dir / "media"

        workspace_roots: list[Path] = []
        if self.workspace is not None:
            workspace_roots.extend([
                self.workspace / "memory" / "audit",
                self.workspace / "memory" / "action",
                self.workspace / "memory" / "security",
                self.workspace / "memory" / "history.jsonl",
                self.workspace / "memory" / ".dream_cursor",
                self.workspace / ".dream_cursor",
                self.workspace / ".cursor",
                self.workspace / "pending_confirmations.json",
                self.workspace / "confirmations",
                self.workspace / "cron",
                self.workspace / "device",
            ])

        self._protected_write_paths = tuple(
            _resolve_loose(path)
            for path in (
                *workspace_roots,
                data_dir / "memory" / "audit",
                data_dir / "memory" / "action",
                data_dir / "memory" / "security",
                data_dir / "memory" / "history.jsonl",
                data_dir / "memory" / ".dream_cursor",
                data_dir / "pending_confirmations.json",
                data_dir / "confirmations",
                cron_dir,
                data_dir / "device",
            )
        )
        self._protected_read_paths = tuple(
            _resolve_loose(path)
            for path in (
                data_dir / "pending_confirmations.json",
                data_dir / "confirmations",
                cron_dir,
                data_dir / "device",
                data_dir / "memory" / "audit",
                data_dir / "memory" / "action",
                data_dir / "memory" / "security",
                data_dir / "memory" / "history.jsonl",
                *(workspace_roots if self.workspace is not None else ()),
            )
        )
        self.read_only_roots = (media_dir,)

    def assert_can_read(self, path: Path) -> None:
        resolved = _resolve_loose(Path(path))
        if self._matches(resolved, self._protected_read_paths):
            raise PolicyDeniedError(
                f"Path {path} is protected runtime state and cannot be read by generic file tools{_NOTE}",
                code="protected_path_read",
                boundary="filesystem",
                policy_rule="protected_path_read",
            )

    def assert_can_list(self, path: Path) -> None:
        self.assert_can_read(path)

    def assert_can_write(self, path: Path) -> None:
        resolved = _resolve_loose(Path(path))
        if self._matches(resolved, self._protected_write_paths):
            raise PolicyDeniedError(
                f"Path {path} is protected runtime state and cannot be modified{_NOTE}",
                code="protected_path_write",
                boundary="filesystem",
                policy_rule="protected_path_write",
            )
        if self._matches(resolved, self.read_only_roots):
            raise PolicyDeniedError(
                f"Path {path} is in a read-only tool root and cannot be modified{_NOTE}",
                code="read_only_root_write",
                boundary="filesystem",
                policy_rule="read_only_root_write",
            )

    @staticmethod
    def _matches(path: Path, roots: tuple[Path, ...]) -> bool:
        return any(path == root or _is_under(path, root) for root in roots)
