"""Append-only audit logs for action explainability.

Audit logs are diagnostic records. They do not authorize actions and should
never contain raw payloads, prompts, source excerpts, or private evidence.
"""

from __future__ import annotations

import json
import os
import re
import uuid
import hashlib
from contextlib import suppress
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from filelock import FileLock

from OriginAgent.agent.action_privacy import FORBIDDEN_METADATA_KEYS
from OriginAgent.utils.helpers import ensure_dir, truncate_text

AUDIT_REASON_MAX_CHARS = 2000
VALID_AUDIT_EVENT_TYPES = {
    "action_decision",
    "confirmation_event",
    "permission_decision",
}
ACTION_DECISIONS_FILE = "action_decisions.jsonl"
CONFIRMATION_EVENTS_FILE = "confirmation_events.jsonl"
PERMISSION_DECISIONS_FILE = "permission_decisions.jsonl"
AUDIT_FILE_BY_TYPE = {
    "action_decision": ACTION_DECISIONS_FILE,
    "confirmation_event": CONFIRMATION_EVENTS_FILE,
    "permission_decision": PERMISSION_DECISIONS_FILE,
}
AUDIT_FORBIDDEN_KEYS = {
    *FORBIDDEN_METADATA_KEYS,
    "access_token",
    "api_key",
    "apikey",
    "authorization",
    "bearer",
    "display_name",
    "password",
    "private_key",
    "prompt",
    "raw_payload",
    "raw_backend_result",
    "raw_evidence",
    "source_excerpt",
    "history",
    "secret",
    "token",
}

_PRIVATE_KEY_RE = re.compile(
    r"-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----.*?-----END [A-Z0-9 ]*PRIVATE KEY-----",
    re.DOTALL,
)
_BEARER_TOKEN_RE = re.compile(r"(?i)\bBearer\s+[A-Za-z0-9._~+/=-]{8,}")
_OPENAI_KEY_RE = re.compile(r"\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b")
_GITHUB_TOKEN_RE = re.compile(r"\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{20,}\b")
_SECRET_ASSIGNMENT_RE = re.compile(
    r"(?i)\b(api[_-]?key|token|secret|password)\b(\s*[:=]\s*)([\"']?)"
    r"[^\"'\s,;]{8,}([\"']?)"
)
_RAW_EVIDENCE_RE = re.compile(
    r"(?i)\b(source_excerpt|prompt|raw_payload|raw_backend_result|raw_evidence|history)"
    r"\b\s*[:=]\s*(\{[^}]*\}|\"[^\"]*\"|'[^']*'|[^,;\n]+)"
)
_EMAIL_RE = re.compile(r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b")
_CHINA_ID_RE = re.compile(
    r"(?<!\d)\d{6}(?:18|19|20)\d{2}(?:0[1-9]|1[0-2])"
    r"(?:0[1-9]|[12]\d|3[01])\d{3}[\dXx](?!\d)"
)
_LONG_NUMBER_RE = re.compile(r"(?<!\d)\d{16,}(?!\d)")
ScopeRedactor = Callable[[str | None], str | None]


@dataclass
class AuditEvent:
    event_id: str
    event_type: str
    created_at: str
    actor_id: str | None
    action_id: str | None
    confirmation_id: str | None
    scope: str | None
    action: str | None
    risk: str | None
    trigger: str | None
    decision: str
    reason: str
    metadata: dict[str, str] = field(default_factory=dict)
    prev_hash: str | None = None
    event_hash: str | None = None

    def __post_init__(self) -> None:
        if self.event_type not in VALID_AUDIT_EVENT_TYPES:
            raise ValueError(f"invalid audit event type: {self.event_type!r}")
        self.event_id = _sanitize_optional_string(self.event_id) or f"audit_{uuid.uuid4().hex[:12]}"
        self.created_at = _sanitize_optional_string(self.created_at) or _format_datetime(None)
        self.actor_id = _sanitize_optional_string(self.actor_id)
        self.action_id = _sanitize_optional_string(self.action_id)
        self.confirmation_id = _sanitize_optional_string(self.confirmation_id)
        self.scope = _sanitize_scope(_sanitize_optional_string(self.scope))
        self.action = _sanitize_optional_string(self.action)
        self.risk = _sanitize_optional_string(self.risk)
        self.trigger = _sanitize_optional_string(self.trigger)
        self.decision = _sanitize_optional_string(self.decision) or "unknown"
        self.reason = _sanitize_optional_string(self.reason) or ""
        self.metadata = sanitize_audit_metadata(self.metadata)

    @classmethod
    def from_dict(cls, raw: dict[str, Any]) -> "AuditEvent":
        return cls(
            event_id=str(raw.get("event_id", "")),
            event_type=str(raw.get("event_type", "")),
            created_at=str(raw.get("created_at", "")),
            actor_id=_optional_raw_string(raw.get("actor_id")),
            action_id=_optional_raw_string(raw.get("action_id")),
            confirmation_id=_optional_raw_string(raw.get("confirmation_id")),
            scope=_optional_raw_string(raw.get("scope")),
            action=_optional_raw_string(raw.get("action")),
            risk=_optional_raw_string(raw.get("risk")),
            trigger=_optional_raw_string(raw.get("trigger")),
            decision=str(raw.get("decision", "")),
            reason=str(raw.get("reason", "")),
            metadata=raw.get("metadata") if isinstance(raw.get("metadata"), dict) else {},
            prev_hash=_optional_raw_string(raw.get("prev_hash")),
            event_hash=_optional_raw_string(raw.get("event_hash")),
        )

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


class AuditLogger:
    def __init__(
        self,
        workspace: Path,
        *,
        lock_factory: Callable[[], FileLock] | None = None,
        redactor: Callable[[str], str] | None = None,
        scope_redactor: ScopeRedactor | None = None,
    ):
        self.workspace = Path(workspace)
        self.memory_dir = ensure_dir(self.workspace / "memory")
        self.audit_dir = self.memory_dir / "audit"
        self._lock_file = self.memory_dir / ".lock"
        self._lock_factory = lock_factory
        self._redactor = redactor
        self._scope_redactor = scope_redactor or _default_scope_redactor
        self._last_hash_by_file: dict[str, str | None] = {}

    def log_action_decision(
        self,
        *,
        action_id: str,
        decision: str,
        reason: str,
        actor_id: str | None = None,
        confirmation_id: str | None = None,
        action: str | None = None,
        scope: str | None = None,
        risk: str | None = None,
        trigger: str | None = None,
        metadata: dict[str, Any] | None = None,
        created_at: datetime | str | None = None,
    ) -> AuditEvent:
        return self._log(
            AuditEvent(
                event_id=_new_event_id(),
                event_type="action_decision",
                created_at=_format_datetime(created_at),
                actor_id=actor_id,
                action_id=action_id,
                confirmation_id=confirmation_id,
                scope=self._scope_redactor(scope),
                action=action,
                risk=risk,
                trigger=trigger,
                decision=decision,
                reason=reason,
                metadata=sanitize_audit_metadata(metadata or {}, redactor=self._redactor),
            )
        )

    def log_confirmation_event(
        self,
        *,
        confirmation_id: str,
        decision: str,
        reason: str,
        actor_id: str | None = None,
        action_id: str | None = None,
        action: str | None = None,
        scope: str | None = None,
        risk: str | None = None,
        trigger: str | None = None,
        metadata: dict[str, Any] | None = None,
        created_at: datetime | str | None = None,
    ) -> AuditEvent:
        return self._log(
            AuditEvent(
                event_id=_new_event_id(),
                event_type="confirmation_event",
                created_at=_format_datetime(created_at),
                actor_id=actor_id,
                action_id=action_id,
                confirmation_id=confirmation_id,
                scope=self._scope_redactor(scope),
                action=action,
                risk=risk,
                trigger=trigger,
                decision=decision,
                reason=reason,
                metadata=sanitize_audit_metadata(metadata or {}, redactor=self._redactor),
            )
        )

    def log_permission_decision(
        self,
        *,
        action_id: str,
        actor_id: str | None,
        action: str,
        scope: str,
        risk: str,
        trigger: str,
        permission: str,
        device_domain: str | None = None,
        attributes: dict[str, Any] | None = None,
        decision: str,
        reason: str,
        actor_role: str,
        confirmation_id: str | None = None,
        metadata: dict[str, Any] | None = None,
        created_at: datetime | str | None = None,
    ) -> AuditEvent:
        merged_metadata = {
            "permission": permission,
            "actor_role": actor_role,
            **(metadata or {}),
        }
        if device_domain is not None:
            merged_metadata["device_domain"] = device_domain
        for key, value in (attributes or {}).items():
            normalized_key = _sanitize_optional_string(key)
            if normalized_key is None or value is None:
                continue
            merged_metadata[normalized_key] = str(value)
        return self._log(
            AuditEvent(
                event_id=_new_event_id(),
                event_type="permission_decision",
                created_at=_format_datetime(created_at),
                actor_id=actor_id,
                action_id=action_id,
                confirmation_id=confirmation_id,
                scope=self._scope_redactor(scope),
                action=action,
                risk=risk,
                trigger=trigger,
                decision=decision,
                reason=reason,
                metadata=sanitize_audit_metadata(merged_metadata, redactor=self._redactor),
            )
        )

    def find_by_action_id(self, action_id: str) -> list[AuditEvent]:
        target = _sanitize_optional_string(action_id)
        if not target:
            return []
        events = [
            event
            for event in self._read_all_events()
            if event.action_id == target
        ]
        related_confirmations = {
            event.confirmation_id
            for event in events
            if event.confirmation_id
        }
        if related_confirmations:
            events.extend(
                event
                for event in self._read_all_events()
                if event.action_id != target
                and event.confirmation_id in related_confirmations
            )
        return _sort_events(_dedupe_events(events))

    def find_by_confirmation_id(self, confirmation_id: str) -> list[AuditEvent]:
        target = _sanitize_optional_string(confirmation_id)
        if not target:
            return []
        return _sort_events([
            event
            for event in self._read_all_events()
            if event.confirmation_id == target
        ])

    def explain_action(self, action_id: str) -> list[AuditEvent]:
        return self.find_by_action_id(action_id)

    def _log(self, event: AuditEvent) -> AuditEvent:
        filename = AUDIT_FILE_BY_TYPE[event.event_type]
        path = self.audit_dir / filename
        with self._locked():
            if not self.audit_dir.exists():
                self.audit_dir.mkdir(parents=True, exist_ok=True)
                _fsync_parent(self.audit_dir)
            event.prev_hash = self._last_hash(filename, path)
            event.event_hash = None
            event.event_hash = _hash_audit_event(event.to_dict())
            line = json.dumps(event.to_dict(), ensure_ascii=False, sort_keys=True) + "\n"
            with path.open("a", encoding="utf-8") as handle:
                handle.write(line)
                handle.flush()
                os.fsync(handle.fileno())
            _fsync_parent(path)
            self._last_hash_by_file[filename] = event.event_hash
        return event

    def _last_hash(self, filename: str, path: Path) -> str | None:
        if filename in self._last_hash_by_file:
            return self._last_hash_by_file[filename]
        value = _read_last_event_hash(path)
        self._last_hash_by_file[filename] = value
        return value

    def _read_all_events(self) -> list[AuditEvent]:
        events: list[AuditEvent] = []
        for filename in AUDIT_FILE_BY_TYPE.values():
            events.extend(self._read_file_events(self.audit_dir / filename))
        return events

    def _read_file_events(self, path: Path) -> list[AuditEvent]:
        try:
            lines = path.read_text(encoding="utf-8").splitlines()
        except FileNotFoundError:
            return []
        except OSError:
            return []
        events: list[AuditEvent] = []
        for line in lines:
            if not line.strip():
                continue
            with suppress(json.JSONDecodeError, TypeError, ValueError):
                raw = json.loads(line)
                if isinstance(raw, dict):
                    events.append(AuditEvent.from_dict(raw))
        return events

    def _locked(self) -> FileLock:
        if self._lock_factory is not None:
            return self._lock_factory()
        return FileLock(str(self._lock_file))


def sanitize_audit_metadata(
    metadata: dict[str, Any],
    *,
    redactor: Callable[[str], str] | None = None,
) -> dict[str, str]:
    sanitized: dict[str, str] = {}
    if not isinstance(metadata, dict):
        return sanitized
    for key, value in metadata.items():
        if not isinstance(key, str):
            continue
        normalized_key = key.strip()
        if not normalized_key:
            continue
        if normalized_key.casefold() in AUDIT_FORBIDDEN_KEYS:
            continue
        if value is None:
            continue
        sanitized[normalized_key] = _sanitize_text(
            _stringify_value(value),
            redactor=redactor,
        )
    return sanitized


def _hash_audit_event(data: dict[str, Any]) -> str:
    payload = dict(data)
    payload["event_hash"] = None
    serialized = json.dumps(payload, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(serialized.encode("utf-8")).hexdigest()


def _read_last_event_hash(path: Path) -> str | None:
    try:
        last = ""
        with path.open("r", encoding="utf-8") as handle:
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


def _sanitize_optional_string(
    value: Any,
    *,
    redactor: Callable[[str], str] | None = None,
) -> str | None:
    if value is None:
        return None
    sanitized = _sanitize_text(str(value), redactor=redactor).strip()
    return sanitized or None


def _sanitize_text(
    value: str,
    *,
    redactor: Callable[[str], str] | None = None,
) -> str:
    cleaned = str(value or "")
    cleaned = _PRIVATE_KEY_RE.sub("[REDACTED_PRIVATE_KEY]", cleaned)
    cleaned = _BEARER_TOKEN_RE.sub("[REDACTED_BEARER_TOKEN]", cleaned)
    cleaned = _OPENAI_KEY_RE.sub("[REDACTED_SECRET]", cleaned)
    cleaned = _GITHUB_TOKEN_RE.sub("[REDACTED_SECRET]", cleaned)
    cleaned = _SECRET_ASSIGNMENT_RE.sub(
        lambda match: (
            f"{match.group(1)}{match.group(2)}"
            f"{match.group(3)}[REDACTED_SECRET]{match.group(4)}"
        ),
        cleaned,
    )
    cleaned = _RAW_EVIDENCE_RE.sub("[REDACTED_EVIDENCE]", cleaned)
    cleaned = _EMAIL_RE.sub("[REDACTED_EMAIL]", cleaned)
    cleaned = _CHINA_ID_RE.sub("[REDACTED_ID]", cleaned)
    cleaned = _LONG_NUMBER_RE.sub("[REDACTED_LONG_NUMBER]", cleaned)
    if redactor is not None:
        cleaned = redactor(cleaned)
    return truncate_text(cleaned, AUDIT_REASON_MAX_CHARS)[:AUDIT_REASON_MAX_CHARS]


def _sanitize_scope(scope: str | None) -> str | None:
    if scope is None:
        return None
    normalized = str(scope).strip().lower()
    if not normalized:
        return None
    normalized = re.sub(r"\s+", ".", normalized)
    normalized = re.sub(r"\.+", ".", normalized).strip(".")
    return normalized or None


def _default_scope_redactor(scope: str | None) -> str | None:
    normalized = _sanitize_scope(scope)
    if normalized is None:
        return None
    parts = [part for part in normalized.split(".") if part]
    if parts and parts[0] == "home" and len(parts) >= 3:
        return ".".join([*parts[:-1], "<device>"])
    return normalized


def _stringify_value(value: Any) -> str:
    if isinstance(value, str):
        return value
    if isinstance(value, bool | int | float):
        return str(value)
    if isinstance(value, list):
        return json.dumps(
            [_stringify_list_item(item) for item in value if item is not None],
            ensure_ascii=False,
            sort_keys=True,
            separators=(",", ":"),
        )
    if isinstance(value, dict):
        return "[REDACTED_COMPLEX_VALUE]"
    return str(value)


def _stringify_list_item(value: Any) -> str:
    if isinstance(value, str):
        return _sanitize_text(value)
    if isinstance(value, bool | int | float):
        return str(value)
    return "[REDACTED_COMPLEX_VALUE]"


def _optional_raw_string(value: Any) -> str | None:
    if value is None:
        return None
    return str(value)


def _new_event_id() -> str:
    return f"audit_{uuid.uuid4().hex[:12]}"


def _format_datetime(value: datetime | str | None) -> str:
    if isinstance(value, str):
        return _sanitize_text(value)
    if value is None:
        value = datetime.now(timezone.utc)
    if value.tzinfo is None:
        value = value.replace(tzinfo=timezone.utc)
    return value.astimezone(timezone.utc).isoformat()


def _sort_events(events: list[AuditEvent]) -> list[AuditEvent]:
    return sorted(events, key=lambda event: event.created_at)


def _dedupe_events(events: list[AuditEvent]) -> list[AuditEvent]:
    by_id: dict[str, AuditEvent] = {}
    for event in events:
        by_id[event.event_id] = event
    return list(by_id.values())


def _fsync_parent(path: Path) -> None:
    if os.name == "nt":
        return
    with suppress(OSError):
        fd = os.open(str(path.parent), os.O_RDONLY)
        try:
            os.fsync(fd)
        finally:
            os.close(fd)
