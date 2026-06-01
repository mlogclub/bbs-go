"""Typed device action schema for low-risk Phase 6 actions."""

from __future__ import annotations

import re
from dataclasses import dataclass, field, replace
from typing import Any

from OriginAgent.agent.action_runtime import ActionIntent

ALLOWED_DEVICE_DOMAINS = {"lighting", "media", "climate"}
DISABLED_DEVICE_DOMAINS = {
    "lock",
    "security",
    "camera",
    "gas",
    "presence",
    "appliance",
}

_IDENTIFIER_RE = re.compile(r"^[A-Za-z0-9_-]+$")


class DeviceActionSchemaError(ValueError):
    code = "schema_error"

    def __init__(self, message: str, *, code: str | None = None):
        super().__init__(message)
        self.code = code or self.code


class UnsupportedDomainError(DeviceActionSchemaError):
    code = "unsupported_domain"


class UnsupportedActionError(DeviceActionSchemaError):
    code = "unsupported_action"


class InvalidActionParameterError(DeviceActionSchemaError):
    code = "invalid_parameter"


@dataclass
class TypedDeviceAction:
    action_type: str
    device_id: str
    domain: str
    room: str | None = None
    parameters: dict[str, Any] = field(default_factory=dict)
    requested_by: str | None = None
    trigger: str = "user_initiated"
    idempotency_key: str | None = None


class DeviceActionSchemaRegistry:
    _ACTION_DOMAINS = {
        "set_light_power": "lighting",
        "set_light_brightness": "lighting",
        "set_light_color_temperature": "lighting",
        "set_media_volume": "media",
        "play_media": "media",
        "stop_media": "media",
        "set_temperature": "climate",
        "set_hvac_mode": "climate",
    }
    _ENUMS = {
        "power": {"on", "off"},
        "temperature": {"warm", "neutral", "cool"},
        "source": {"ambient", "playlist", "radio"},
        "mode": {"off", "heat", "cool", "auto", "fan"},
    }

    def validate(self, action: TypedDeviceAction) -> TypedDeviceAction:
        normalized = self._normalize_action(action)
        expected_domain = self._ACTION_DOMAINS.get(normalized.action_type)
        if expected_domain is None:
            raise UnsupportedActionError(
                f"unsupported action_type {normalized.action_type}",
                code="unsupported_action",
            )
        if normalized.domain != expected_domain:
            raise UnsupportedActionError(
                f"action_type {normalized.action_type} is not supported for domain {normalized.domain}",
                code="domain_action_mismatch",
            )
        self._validate_parameters(normalized)
        return normalized

    def infer_risk(self, action: TypedDeviceAction) -> str:
        validated = self.validate(action)
        if validated.domain == "climate":
            return "medium"
        if (
            validated.action_type == "set_media_volume"
            and validated.parameters.get("volume", 0) > 70
        ):
            return "medium"
        return "low"

    def infer_scope(self, action: TypedDeviceAction) -> str:
        validated = self.validate(action)
        if validated.room:
            return f"home.{validated.room}.{validated.domain}.{validated.device_id}"
        return f"home.{validated.domain}.{validated.device_id}"

    def requires_presence_empty(self, action: TypedDeviceAction) -> bool:
        self.validate(action)
        return False

    def _normalize_action(self, action: TypedDeviceAction) -> TypedDeviceAction:
        action_type = _required_identifier(action.action_type, "action_type")
        device_id = _required_identifier(action.device_id, "device_id")
        domain = _required_identifier(action.domain, "domain")
        room = _optional_identifier(action.room, "room")
        if domain in DISABLED_DEVICE_DOMAINS:
            raise UnsupportedDomainError(
                f"domain {domain} is disabled in Phase 6",
                code="disabled_domain",
            )
        if domain not in ALLOWED_DEVICE_DOMAINS:
            raise UnsupportedDomainError(
                f"unknown domain {domain}",
                code="unsupported_domain",
            )
        if not isinstance(action.parameters, dict):
            raise InvalidActionParameterError(
                "parameters must be a dict",
                code="invalid_parameters",
            )
        return replace(
            action,
            action_type=action_type,
            device_id=device_id,
            domain=domain,
            room=room,
            parameters=dict(action.parameters),
            trigger=str(action.trigger or "user_initiated").strip().lower(),
            requested_by=_optional_string(action.requested_by),
            idempotency_key=_optional_string(action.idempotency_key),
        )

    def _validate_parameters(self, action: TypedDeviceAction) -> None:
        parameters = action.parameters
        if action.action_type == "set_light_power":
            _require_exact_keys(parameters, {"power"})
            _require_enum(parameters, "power", self._ENUMS["power"])
            return
        if action.action_type == "set_light_brightness":
            _require_exact_keys(parameters, {"brightness"})
            _require_int_range(parameters, "brightness", 0, 100)
            return
        if action.action_type == "set_light_color_temperature":
            _require_exact_keys(parameters, {"temperature"})
            _require_enum(parameters, "temperature", self._ENUMS["temperature"])
            return
        if action.action_type == "set_media_volume":
            _require_exact_keys(parameters, {"volume"})
            _require_int_range(parameters, "volume", 0, 100)
            return
        if action.action_type == "play_media":
            _require_exact_keys(parameters, {"source"})
            _require_enum(parameters, "source", self._ENUMS["source"])
            return
        if action.action_type == "stop_media":
            _require_exact_keys(parameters, set())
            return
        if action.action_type == "set_temperature":
            _require_exact_keys(parameters, {"temperature_c"})
            _require_float_range(parameters, "temperature_c", 16, 30)
            return
        if action.action_type == "set_hvac_mode":
            _require_exact_keys(parameters, {"mode"})
            _require_enum(parameters, "mode", self._ENUMS["mode"])


class TypedActionPlanner:
    def __init__(self, registry: DeviceActionSchemaRegistry):
        self.registry = registry

    def to_intent(self, action: TypedDeviceAction) -> ActionIntent:
        validated = self.registry.validate(action)
        payload = {
            "device_id": validated.device_id,
            "domain": validated.domain,
            "action_type": validated.action_type,
            **validated.parameters,
        }
        return ActionIntent(
            action=validated.action_type,
            scope=self.registry.infer_scope(validated),
            trigger=validated.trigger,
            risk=self.registry.infer_risk(validated),
            requested_by=validated.requested_by,
            requires_presence_empty=self.registry.requires_presence_empty(validated),
            payload=payload,
            idempotency_key=validated.idempotency_key,
        )


def _required_identifier(value: Any, field_name: str) -> str:
    normalized = str(value or "").strip().lower()
    if not normalized:
        raise InvalidActionParameterError(f"{field_name} is required", code="invalid_identifier")
    if not _IDENTIFIER_RE.match(normalized):
        raise InvalidActionParameterError(
            f"{field_name} must contain only letters, numbers, underscores, or hyphens",
            code="invalid_identifier",
        )
    return normalized


def _optional_identifier(value: Any, field_name: str) -> str | None:
    if value is None:
        return None
    normalized = str(value).strip().lower()
    if not normalized:
        return None
    if not _IDENTIFIER_RE.match(normalized):
        raise InvalidActionParameterError(
            f"{field_name} must contain only letters, numbers, underscores, or hyphens",
            code="invalid_identifier",
        )
    return normalized


def _optional_string(value: Any) -> str | None:
    if value is None:
        return None
    normalized = str(value).strip()
    return normalized or None


def _require_exact_keys(parameters: dict[str, Any], expected: set[str]) -> None:
    actual = set(parameters.keys())
    if actual != expected:
        raise InvalidActionParameterError(
            f"expected parameters {sorted(expected)}, got {sorted(actual)}",
            code="invalid_parameter_keys",
        )


def _require_enum(parameters: dict[str, Any], key: str, allowed: set[str]) -> None:
    value = parameters.get(key)
    if not isinstance(value, str) or value.strip().lower() not in allowed:
        raise InvalidActionParameterError(
            f"{key} must be one of {sorted(allowed)}",
            code="invalid_parameter",
        )
    parameters[key] = value.strip().lower()


def _require_int_range(parameters: dict[str, Any], key: str, minimum: int, maximum: int) -> None:
    value = parameters.get(key)
    if isinstance(value, bool) or not isinstance(value, int):
        raise InvalidActionParameterError(
            f"{key} must be an int",
            code="invalid_parameter",
        )
    if value < minimum or value > maximum:
        raise InvalidActionParameterError(
            f"{key} must be between {minimum} and {maximum}",
            code="invalid_parameter",
        )


def _require_float_range(
    parameters: dict[str, Any],
    key: str,
    minimum: float,
    maximum: float,
) -> None:
    value = parameters.get(key)
    if isinstance(value, bool) or not isinstance(value, int | float):
        raise InvalidActionParameterError(
            f"{key} must be numeric",
            code="invalid_parameter",
        )
    if value < minimum or value > maximum:
        raise InvalidActionParameterError(
            f"{key} must be between {minimum} and {maximum}",
            code="invalid_parameter",
        )
    parameters[key] = float(value)
