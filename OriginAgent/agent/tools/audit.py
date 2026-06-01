"""Safe audit events for generic tool calls."""

from __future__ import annotations

import hashlib
import json
import os
import threading
import uuid
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Literal, Protocol

from loguru import logger

from OriginAgent.utils.helpers import ensure_dir

ToolCallStatus = Literal[
    "success",
    "validation_error",
    "policy_denied",
    "error",
    "interrupted",
]
ToolAuditMode = Literal["off", "minimal", "security"]


@dataclass(frozen=True)
class ToolAuditConfig:
    mode: ToolAuditMode = "minimal"
    security_tools: tuple[str, ...] = (
        "exec",
        "message",
        "web_fetch",
        "content_read",
        "cron",
        "spawn",
        "originagent_device_*",
        "mcp_*",
    )
    security_on_policy_denial: bool = True

    @classmethod
    def from_config(cls, value: object | None) -> "ToolAuditConfig":
        if value is None:
            return cls()
        default = cls()
        return cls(
            mode=getattr(value, "mode", default.mode),
            security_tools=tuple(getattr(value, "security_tools", default.security_tools)),
            security_on_policy_denial=bool(
                getattr(value, "security_on_policy_denial", default.security_on_policy_denial)
            ),
        )


@dataclass(frozen=True)
class ToolCallAuditEvent:
    tool_name: str
    status: ToolCallStatus
    duration_ms: int
    read_only: bool
    exclusive: bool
    error_kind: str | None = None
    actor_id_hash: str | None = None
    session_key_hash: str | None = None
    subagent_task_id: str | None = None
    parent_session_key_hash: str | None = None
    origin_channel: str | None = None
    origin_chat_id_hash: str | None = None
    target_kind: str | None = None
    target_hash: str | None = None
    policy_rule: str | None = None
    result_size: int | None = None
    created_at: str = field(
        default_factory=lambda: datetime.now(timezone.utc).isoformat()
    )
    event_id: str = field(default_factory=lambda: f"tool_audit_{uuid.uuid4().hex[:12]}")
    prev_hash: str | None = None
    event_hash: str | None = None

    def to_dict(self) -> dict[str, object]:
        return asdict(self)


class ToolAuditSink(Protocol):
    def record(self, event: ToolCallAuditEvent) -> None:
        ...


class InMemoryToolAuditSink:
    def __init__(self) -> None:
        self.events: list[ToolCallAuditEvent] = []

    def record(self, event: ToolCallAuditEvent) -> None:
        self.events.append(event)


class JsonlToolAuditSink:
    """Best-effort file sink for generic tool audit events."""

    def __init__(self, workspace: Path):
        self.workspace = Path(workspace)
        self.path = self.workspace / "memory" / "audit" / "tool_calls.jsonl"
        self._last_hash: str | None = self._load_last_hash()
        self._lock = threading.Lock()

    def _load_last_hash(self) -> str | None:
        try:
            if not self.path.exists():
                return None
            last = ""
            with self.path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    if line.strip():
                        last = line
            if not last:
                return None
            data = json.loads(last)
            value = data.get("event_hash")
            return value if isinstance(value, str) else None
        except Exception:
            return None

    @staticmethod
    def _hash_event(data: dict[str, object]) -> str:
        payload = json.dumps(data, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
        return hashlib.sha256(payload.encode("utf-8")).hexdigest()

    def record(self, event: ToolCallAuditEvent) -> None:
        try:
            with self._lock:
                ensure_dir(self.path.parent)
                event_data = event.to_dict()
                event_data["prev_hash"] = self._last_hash
                event_data["event_hash"] = None
                event_data["event_hash"] = self._hash_event(event_data)
                with self.path.open("a", encoding="utf-8") as handle:
                    handle.write(json.dumps(event_data, ensure_ascii=False, sort_keys=True))
                    handle.write("\n")
                    handle.flush()
                    os.fsync(handle.fileno())
                self._last_hash = str(event_data["event_hash"])
        except Exception as exc:
            logger.debug("Tool audit write failed: {}", exc)
