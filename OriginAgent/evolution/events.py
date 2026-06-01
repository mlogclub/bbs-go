"""Evolution lifecycle events, not MessageBus events (see OriginAgent.bus.events)."""

from __future__ import annotations

import uuid
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from enum import Enum
from typing import Any, Mapping

EVENT_SCHEMA_VERSION = "originagent.evolution.event.v1"
SOURCE_EVENT_STREAMS = frozenset({"evolution", "skill_lifecycle", "domain_pack"})


class EventValidationError(ValueError):
    """Raised when an evolution event is invalid."""


class EventType(str, Enum):
    MODULE_PROPOSED = "module_proposed"
    MODULE_MANIFEST_VALIDATED = "module_manifest_validated"
    MODULE_INSTALLED_STAGING = "module_installed_staging"
    MODULE_FAILED = "module_failed"
    MODULE_PERMISSION_CHECKED = "module_permission_checked"
    MODULE_VERIFIED = "module_verified"
    MODULE_ACTIVATED = "module_activated"
    MODULE_DEPRECATED = "module_deprecated"
    STATE_BRANCH_CREATED = "state_branch_created"
    STATE_BRANCH_MERGED = "state_branch_merged"
    STATE_BRANCH_DISCARDED = "state_branch_discarded"
    MODULE_ROLLBACK_STARTED = "module_rollback_started"
    MODULE_ROLLBACK_SUCCEEDED = "module_rollback_succeeded"
    MODULE_ROLLBACK_FAILED = "module_rollback_failed"
    DIRTY_ROLLBACK = "dirty_rollback"
    TEARDOWN_STARTED = "teardown_started"
    TEARDOWN_SUCCEEDED = "teardown_succeeded"
    TEARDOWN_FAILED = "teardown_failed"
    MODULE_FORCE_CLEAN_REQUESTED = "module_force_clean_requested"
    MODULE_FORCE_CLEAN_SUCCEEDED = "module_force_clean_succeeded"
    EXTERNAL_SIDE_EFFECT_ABANDONED = "external_side_effect_abandoned"
    TELEMETRY_RECORDED = "telemetry_recorded"
    MEMORY_VAULT_EXPORTED = "memory_vault_exported"
    MEMORY_VAULT_IMPORTED = "memory_vault_imported"
    SCORE_COMMITTED = "score_committed"
    SCORE_REVEALED = "score_revealed"
    UNMAPPED = "unmapped"


PHASE_1A_EVENT_TYPES = frozenset(
    {
        EventType.MODULE_PROPOSED,
        EventType.MODULE_MANIFEST_VALIDATED,
        EventType.MODULE_INSTALLED_STAGING,
        EventType.MODULE_FAILED,
    }
)

_CREATEABLE_EVENT_TYPES = frozenset(EventType)

_SKILL_ACTION_MAP = {
    "verify": EventType.MODULE_VERIFIED,
    "activate": EventType.MODULE_ACTIVATED,
    "reject": EventType.MODULE_FAILED,
    "deprecate": EventType.MODULE_DEPRECATED,
    "always_on": None,
    "always_off": None,
}

_DOMAIN_PACK_ACTION_MAP = {
    "install": EventType.MODULE_INSTALLED_STAGING,
    "upgrade": None,
    "uninstall": None,
    "enable": None,
    "disable": None,
    "activate": EventType.MODULE_ACTIVATED,
    "deactivate": EventType.MODULE_DEPRECATED,
    "eval": EventType.MODULE_VERIFIED,
    "move_to_domain": None,
}


@dataclass(frozen=True)
class EvolutionEvent:
    """One append-only local evolution event."""

    event_id: str
    event_type: str
    created_at: str
    schema_version: str = EVENT_SCHEMA_VERSION
    actor: str = "user"
    actor_public_key: str = ""
    module_id: str = ""
    module_version: str = ""
    module_type: str = ""
    source_event_stream: str = "evolution"
    source_event_id: str = ""
    artifact_digest: str = ""
    state_branch_id: str = ""
    capability_snapshot: dict[str, Any] = field(default_factory=dict)
    result: dict[str, Any] = field(default_factory=dict)
    previous_event_hash: str | None = None
    event_hash: str = ""
    signature: str = ""
    signature_scheme: str = "ed25519"

    def __post_init__(self) -> None:
        if self.schema_version != EVENT_SCHEMA_VERSION:
            raise EventValidationError(f"unsupported schema_version: {self.schema_version}")
        if not self.event_id:
            raise EventValidationError("event_id is required")
        if self.event_type not in {event.value for event in _CREATEABLE_EVENT_TYPES}:
            raise EventValidationError(f"unsupported event_type: {self.event_type}")
        if not self.created_at:
            raise EventValidationError("created_at is required")
        if self.source_event_stream not in SOURCE_EVENT_STREAMS:
            raise EventValidationError(f"unsupported source_event_stream: {self.source_event_stream}")

    @classmethod
    def new(
        cls,
        event_type: EventType | str,
        *,
        actor: str = "user",
        actor_public_key: str = "",
        module_id: str = "",
        module_version: str = "",
        module_type: str = "",
        source_event_stream: str = "evolution",
        source_event_id: str = "",
        artifact_digest: str = "",
        state_branch_id: str = "",
        capability_snapshot: dict[str, Any] | None = None,
        result: dict[str, Any] | None = None,
    ) -> "EvolutionEvent":
        """Create a local evolution event with an ISO 8601 UTC timestamp."""

        normalized = _event_type_value(event_type)
        return cls(
            event_id=f"evolution_{uuid.uuid4().hex}",
            event_type=normalized,
            created_at=datetime.now(timezone.utc).isoformat(),
            actor=actor,
            actor_public_key=actor_public_key,
            module_id=module_id,
            module_version=module_version,
            module_type=module_type,
            source_event_stream=source_event_stream,
            source_event_id=source_event_id,
            artifact_digest=artifact_digest,
            state_branch_id=state_branch_id,
            capability_snapshot=dict(capability_snapshot or {}),
            result=dict(result or {}),
        )

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


def normalize_event_type(source_event_stream: str, raw_event: Mapping[str, Any]) -> EventType | None:
    """Normalize an evolution or legacy lifecycle event type."""

    source = str(source_event_stream or "").strip()
    if source == "evolution":
        try:
            return EventType(str(raw_event.get("event_type") or ""))
        except ValueError:
            return None
    if source == "skill_lifecycle":
        action = str(raw_event.get("action") or "").strip().lower()
        return _SKILL_ACTION_MAP.get(action)
    if source == "domain_pack":
        action = str(raw_event.get("action") or "").strip().lower()
        return _DOMAIN_PACK_ACTION_MAP.get(action)
    return None


def _event_type_value(event_type: EventType | str) -> str:
    if isinstance(event_type, EventType):
        return event_type.value
    value = str(event_type or "")
    try:
        normalized = EventType(value)
    except ValueError as exc:
        raise EventValidationError(f"unsupported event_type: {value}") from exc
    return normalized.value
