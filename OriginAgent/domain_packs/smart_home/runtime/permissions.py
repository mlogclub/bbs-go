"""Smart-home permission model."""

from __future__ import annotations

from dataclasses import dataclass

from OriginAgent.agent.permissions import PermissionDecision, PermissionRequest

VALID_HOUSEHOLD_ROLES = {
    "admin",
    "resident",
    "guest",
    "child",
    "elder",
    "service_person",
    "unknown",
}
VALID_DEVICE_DOMAINS = {
    "lighting",
    "climate",
    "media",
    "lock",
    "security",
    "camera",
    "gas",
    "appliance",
    "presence",
    "general",
}
HIGH_RISK_DEVICE_DOMAINS = {"lock", "security", "camera", "gas", "presence"}


@dataclass
class HouseholdActor:
    actor_id: str
    role: str
    display_name: str | None = None

    def __post_init__(self) -> None:
        self.actor_id = _normalize_actor_id(self.actor_id) or "unknown"
        self.role = _normalize_role(self.role)
        if self.display_name is not None:
            self.display_name = self.display_name.strip() or None


class PermissionResolver:
    def __init__(self, actors: dict[str, HouseholdActor] | None = None):
        self.actors: dict[str, HouseholdActor] = {}
        for actor_id, actor in (actors or {}).items():
            normalized_id = _normalize_actor_id(actor_id)
            if normalized_id:
                self.actors[normalized_id] = actor

    def resolve_actor(self, actor_id: str | None) -> HouseholdActor:
        normalized_id = _normalize_actor_id(actor_id)
        if not normalized_id:
            return HouseholdActor("unknown", "unknown")
        return self.actors.get(normalized_id, HouseholdActor(normalized_id, "unknown"))

    def evaluate(self, request: PermissionRequest) -> PermissionDecision:
        actor = self.resolve_actor(request.actor_id)
        role = actor.role

        if request.permission in {"manage_presence", "manage_facts", "create_rule"}:
            if role == "admin":
                return _allow(role, f"{role} may {request.permission}")
            return _ask_admin_or_deny(role, request.permission)

        if request.permission == "confirm_action":
            return self._evaluate_confirm(request, role)

        if request.permission == "execute_action":
            return self._evaluate_execute(request, role)

        return _deny(role, f"unsupported permission: {request.permission}")

    def _evaluate_confirm(
        self,
        request: PermissionRequest,
        role: str,
    ) -> PermissionDecision:
        if role == "admin":
            return _allow(role, "admin may confirm action")
        if _requires_admin(request):
            return PermissionDecision(
                decision="ask_admin",
                reason="action requires administrator confirmation",
                actor_role=role,
            )
        if role in {"resident", "elder"} and request.risk in {"low", "medium"}:
            return _allow(role, f"{role} may confirm low/medium action")
        return _deny(role, f"{role} cannot confirm action")

    def _evaluate_execute(
        self,
        request: PermissionRequest,
        role: str,
    ) -> PermissionDecision:
        if role == "admin":
            return _allow(role, "admin may execute action")
        if _requires_admin(request):
            return PermissionDecision(
                decision="ask_admin",
                reason="action requires administrator permission",
                actor_role=role,
            )
        if role in {"resident", "elder"} and request.risk in {"low", "medium"}:
            return _allow(role, f"{role} may execute low/medium action")
        if (
            role in {"guest", "child", "service_person", "unknown"}
            and request.risk == "low"
            and request.trigger == "user_initiated"
        ):
            return _allow(role, f"{role} may execute low-risk user action")
        return _deny(role, f"{role} cannot execute action")


def infer_device_domain(scope: str, action: str) -> str:
    text = f"{scope or ''} {action or ''}".strip().casefold()
    if not text:
        return "general"
    if _contains_any(text, ("unlock", "door", ".lock", " lock", "lock_")):
        return "lock"
    if _contains_any(text, ("security", "alarm", "arm_", " arm", "disarm")):
        return "security"
    if "camera" in text:
        return "camera"
    if _contains_any(text, ("gas", "stove")):
        return "gas"
    if "presence" in text:
        return "presence"
    if _contains_any(text, ("light", "lighting")):
        return "lighting"
    if _contains_any(text, ("climate", "thermostat", ".ac", " ac", "hvac")):
        return "climate"
    if _contains_any(text, ("media", "music", "speaker")):
        return "media"
    if "appliance" in text:
        return "appliance"
    return "general"


def _contains_any(text: str, needles: tuple[str, ...]) -> bool:
    return any(needle in text for needle in needles)


def _requires_admin(request: PermissionRequest) -> bool:
    return request.risk == "high" or request.device_domain in HIGH_RISK_DEVICE_DOMAINS


def _allow(role: str, reason: str) -> PermissionDecision:
    return PermissionDecision(decision="allow", reason=reason, actor_role=role)


def _deny(role: str, reason: str) -> PermissionDecision:
    return PermissionDecision(decision="deny", reason=reason, actor_role=role)


def _ask_admin_or_deny(role: str, permission: str) -> PermissionDecision:
    if role in {"resident", "elder", "guest", "child", "service_person", "unknown"}:
        return PermissionDecision(
            decision="ask_admin",
            reason=f"{permission} requires administrator permission",
            actor_role=role,
        )
    return _deny(role, f"{role} cannot {permission}")


def _normalize_actor_id(actor_id: str | None) -> str | None:
    if actor_id is None:
        return None
    normalized = str(actor_id).strip()
    return normalized or None


def _normalize_role(role: str | None) -> str:
    normalized = str(role or "unknown").strip().lower()
    if normalized not in VALID_HOUSEHOLD_ROLES:
        raise ValueError(f"invalid household role: {role!r}")
    return normalized


def _normalize_device_domain(device_domain: str | None) -> str:
    normalized = str(device_domain or "general").strip().lower()
    if normalized not in VALID_DEVICE_DOMAINS:
        raise ValueError(f"invalid device domain: {device_domain!r}")
    return normalized
