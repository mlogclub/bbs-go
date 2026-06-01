"""Durable records for subagent task, lifecycle, and tool activity."""

from __future__ import annotations

import json
import os
import threading
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Literal

from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.utils.helpers import ensure_dir, truncate_text

SubagentTaskStatus = Literal[
    "spawned",
    "running",
    "completed",
    "failed",
    "interrupted",
    "cancelled",
]
SubagentLifecycleState = Literal[
    "spawned",
    "running",
    "awaiting_tools",
    "tools_completed",
    "final_response",
    "completed",
    "failed",
    "interrupted",
    "cancelled",
]
SubagentToolStatus = Literal["success", "denied", "failed", "interrupted"]

_SUMMARY_MAX_CHARS = 240
_RESULT_MAX_CHARS = 400
_PROFILE_MAX_CHARS = 160
_TOOL_LIST_MAX = 24
_RECENT_TASK_LIMIT = 10


def _utcnow() -> str:
    return datetime.now(timezone.utc).isoformat()


def summarize_text(value: Any, max_chars: int = _SUMMARY_MAX_CHARS) -> str:
    text = redact_memory_text(str(value or "").strip())
    if not text:
        return ""
    return truncate_text(text, max_chars).replace("\n... (truncated)", " ...")


def summarize_payload(value: Any, max_chars: int = _RESULT_MAX_CHARS) -> str:
    if value is None:
        return ""
    if isinstance(value, str):
        return summarize_text(value, max_chars=max_chars)
    try:
        rendered = json.dumps(value, ensure_ascii=False, sort_keys=True)
    except Exception:
        rendered = str(value)
    return summarize_text(rendered, max_chars=max_chars)


@dataclass(frozen=True)
class SubagentTaskRecord:
    subagent_id: str
    root_subagent_id: str | None
    parent_subagent_id: str | None
    subagent_depth: int
    parent_session_key: str | None
    origin_channel: str | None
    origin_chat_id: str | None
    origin_message_id: str | None
    task_label: str
    task_summary: str
    delegated_profile_summary: str
    allowed_tools_summary: list[str]
    provider_summary: str
    isolation_mode: str
    terminal_status: SubagentTaskStatus
    stop_reason: str | None = None
    failure_summary: str | None = None
    started_at: str | None = None
    ended_at: str | None = None
    created_at: str = field(default_factory=_utcnow)

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class SubagentLifecycleRecord:
    subagent_id: str
    root_subagent_id: str | None
    parent_subagent_id: str | None
    subagent_depth: int
    parent_session_key: str | None
    state: SubagentLifecycleState
    detail: str | None = None
    created_at: str = field(default_factory=_utcnow)

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class SubagentToolRecord:
    subagent_id: str
    root_subagent_id: str | None
    parent_subagent_id: str | None
    subagent_depth: int
    parent_session_key: str | None
    tool_name: str
    status: SubagentToolStatus
    started_at: str
    ended_at: str
    duration_ms: int
    argument_summary: str
    result_summary: str | None = None
    policy_rule: str | None = None
    error_kind: str | None = None
    created_at: str = field(default_factory=_utcnow)

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


class JsonlSubagentRecordStore:
    """Append-only JSONL store for subagent execution records."""

    def __init__(self, workspace: Path):
        root = Path(workspace) / "memory" / "subagents"
        self._tasks_path = root / "tasks.jsonl"
        self._lifecycle_path = root / "lifecycle.jsonl"
        self._tools_path = root / "tools.jsonl"
        self._lock = threading.Lock()

    def append_task(self, record: SubagentTaskRecord) -> None:
        self._append(self._tasks_path, record.to_dict())

    def append_lifecycle(self, record: SubagentLifecycleRecord) -> None:
        self._append(self._lifecycle_path, record.to_dict())

    def append_tool(self, record: SubagentToolRecord) -> None:
        self._append(self._tools_path, record.to_dict())

    def recent_task_summary(self, limit: int = _RECENT_TASK_LIMIT) -> dict[str, Any]:
        records = list(self._iter_jsonl(self._tasks_path))
        if not records:
            return {
                "subagent_task_total": 0,
                "subagent_recent_task_count": 0,
                "subagent_terminal_status_counts": {},
                "subagent_recent_tasks": [],
                "subagent_last_task_at": None,
            }
        latest_by_task: dict[str, dict[str, Any]] = {}
        ordered_ids: list[str] = []
        for item in records:
            subagent_id = str(item.get("subagent_id") or "").strip()
            if not subagent_id:
                continue
            if subagent_id not in latest_by_task:
                ordered_ids.append(subagent_id)
            latest_by_task[subagent_id] = item
        canonical = [latest_by_task[subagent_id] for subagent_id in ordered_ids if subagent_id in latest_by_task]
        status_counts: dict[str, int] = {}
        for item in canonical:
            status = str(item.get("terminal_status") or "unknown")
            status_counts[status] = status_counts.get(status, 0) + 1
        recent = canonical[-max(1, limit):]
        return {
            "subagent_task_total": len(canonical),
            "subagent_recent_task_count": len(recent),
            "subagent_terminal_status_counts": status_counts,
            "subagent_recent_tasks": [
                {
                    "subagent_id": item.get("subagent_id"),
                    "root_subagent_id": item.get("root_subagent_id"),
                    "parent_subagent_id": item.get("parent_subagent_id"),
                    "subagent_depth": item.get("subagent_depth"),
                    "task_label": item.get("task_label"),
                    "terminal_status": item.get("terminal_status"),
                    "stop_reason": item.get("stop_reason"),
                    "provider_summary": item.get("provider_summary"),
                    "isolation_mode": item.get("isolation_mode"),
                    "ended_at": item.get("ended_at"),
                }
                for item in reversed(recent)
            ],
            "subagent_last_task_at": recent[-1].get("ended_at") or recent[-1].get("created_at"),
        }

    @staticmethod
    def delegated_profile_summary(
        *,
        capability_snapshot: Any,
        allow_web: bool,
        allowed_tool_names: set[str] | frozenset[str],
    ) -> str:
        parts = [
            f"read={'on' if getattr(capability_snapshot, 'can_read_files', False) else 'off'}",
            f"write={'on' if getattr(capability_snapshot, 'can_write_files', False) else 'off'}",
            f"exec={'on' if getattr(capability_snapshot, 'can_exec', False) else 'off'}",
            f"cron={'on' if getattr(capability_snapshot, 'can_create_cron', False) else 'off'}",
            f"spawn={'on' if getattr(capability_snapshot, 'can_spawn', False) else 'off'}",
            f"web={'on' if allow_web else 'off'}",
            f"tools={min(len(allowed_tool_names), _TOOL_LIST_MAX)}",
        ]
        return summarize_text(", ".join(parts), max_chars=_PROFILE_MAX_CHARS)

    @staticmethod
    def allowed_tools_summary(allowed_tool_names: set[str] | frozenset[str]) -> list[str]:
        return sorted(str(name) for name in allowed_tool_names)[:_TOOL_LIST_MAX]

    def _append(self, path: Path, payload: dict[str, Any]) -> None:
        with self._lock:
            ensure_dir(path.parent)
            with path.open("a", encoding="utf-8") as handle:
                handle.write(json.dumps(payload, ensure_ascii=False, sort_keys=True))
                handle.write("\n")
                handle.flush()
                os.fsync(handle.fileno())

    @staticmethod
    def _iter_jsonl(path: Path) -> list[dict[str, Any]]:
        if not path.exists():
            return []
        items: list[dict[str, Any]] = []
        try:
            with path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(data, dict):
                        items.append(data)
        except Exception:
            return []
        return items
