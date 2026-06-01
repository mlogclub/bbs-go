"""Generic permission request/decision model."""

from __future__ import annotations

from dataclasses import InitVar, dataclass, field
from typing import Any

VALID_PERMISSION_ACTIONS = {
    "execute_action",
    "confirm_action",
    "create_rule",
    "manage_presence",
    "manage_facts",
}
VALID_PERMISSION_DECISIONS = {"allow", "deny", "ask_admin"}
VALID_RISKS = {"low", "medium", "high"}
VALID_TRIGGERS = {"user_initiated", "scheduled", "system", "subagent"}


@dataclass
class PermissionRequest:
    actor_id: str | None
    action: str
    scope: str
    risk: str
    trigger: str
    permission: str
    attributes: dict[str, str] = field(default_factory=dict)
    device_domain: InitVar[str | None] = None

    def __post_init__(self, device_domain: str | None) -> None:
        self.actor_id = _normalize_actor_id(self.actor_id)
        self.action = str(self.action or "").strip()
        self.scope = str(self.scope or "").strip()
        self.risk = str(self.risk or "low").strip().lower()
        self.trigger = str(self.trigger or "").strip().lower()
        self.permission = _normalize_permission(self.permission)
        if self.risk not in VALID_RISKS:
            raise ValueError(f"invalid risk: {self.risk!r}")
        if self.trigger not in VALID_TRIGGERS:
            raise ValueError(f"invalid trigger: {self.trigger!r}")
        self.attributes = _normalize_attributes(self.attributes)
        if device_domain is not None:
            self.attributes.setdefault("device_domain", _normalize_attribute_value(device_domain))

    @property
    def device_domain(self) -> str:
        return self.attributes.get("device_domain", "general")

    def attribute(self, name: str, default: str | None = None) -> str | None:
        normalized_name = str(name or "").strip().lower()
        if not normalized_name:
            return default
        return self.attributes.get(normalized_name, default)


@dataclass
class PermissionDecision:
    decision: str
    reason: str
    actor_role: str = "unknown"

    def __post_init__(self) -> None:
        self.decision = _normalize_decision(self.decision)
        self.actor_role = _normalize_attribute_value(self.actor_role) or "unknown"
        self.reason = str(self.reason or "")


class PermissionResolver:
    """Default resolver for domain-neutral actions.

    Domain packs should provide stricter resolvers for domain-specific action
    semantics.
    """

    def evaluate(self, request: PermissionRequest) -> PermissionDecision:
        if request.permission in {"manage_presence", "manage_facts", "create_rule"}:
            return PermissionDecision(
                decision="ask_admin",
                reason=f"{request.permission} requires explicit administrator approval",
            )
        if request.permission in {"confirm_action", "execute_action"}:
            if request.risk == "high":
                return PermissionDecision(
                    decision="ask_admin",
                    reason="high-risk actions require explicit administrator approval",
                )
            return PermissionDecision(
                decision="allow",
                reason="default resolver allows low/medium risk actions",
            )
        return PermissionDecision(
            decision="deny",
            reason=f"unsupported permission: {request.permission}",
        )


def _normalize_actor_id(actor_id: str | None) -> str | None:
    if actor_id is None:
        return None
    normalized = str(actor_id).strip()
    return normalized or None


def _normalize_permission(permission: str | None) -> str:
    normalized = str(permission or "").strip().lower()
    if normalized not in VALID_PERMISSION_ACTIONS:
        raise ValueError(f"invalid permission action: {permission!r}")
    return normalized


def _normalize_decision(decision: str | None) -> str:
    normalized = str(decision or "").strip().lower()
    if normalized not in VALID_PERMISSION_DECISIONS:
        raise ValueError(f"invalid permission decision: {decision!r}")
    return normalized


def _normalize_attribute_value(value: Any) -> str:
    return str(value or "").strip().lower()


def _normalize_attributes(attributes: dict[str, Any] | None) -> dict[str, str]:
    normalized: dict[str, str] = {}
    if not isinstance(attributes, dict):
        return normalized
    for key, value in attributes.items():
        key_text = str(key or "").strip().lower()
        value_text = str(value or "").strip().lower()
        if key_text and value_text:
            normalized[key_text] = value_text
    return normalized
