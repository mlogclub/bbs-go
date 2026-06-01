"""Optional real-device integrations for low-risk validation pilots."""

from __future__ import annotations

from typing import Any, Protocol

from OriginAgent.agent.action_runtime import ActionIntent
from .permissions import infer_device_domain

_LIGHTING_DOMAIN = "lighting"
_LIGHTING_ACTIONS = {
    "set_light_power",
    "set_light_brightness",
    "set_light_color_temperature",
}
_LIGHT_TEMPERATURES = {"warm", "neutral", "cool"}
_LIGHT_POWER_STATES = {"on", "off"}


class LightingDeviceClient(Protocol):
    def set_power(self, device_id: str, power: str) -> dict[str, Any]:
        ...

    def set_brightness(self, device_id: str, brightness: int) -> dict[str, Any]:
        ...

    def set_color_temperature(self, device_id: str, temperature: str) -> dict[str, Any]:
        ...


class RealLightingBackend:
    def __init__(
        self,
        client: LightingDeviceClient,
        *,
        real_mode: bool = False,
    ):
        self.client = client
        self.real_mode = real_mode

    def execute(self, intent: ActionIntent) -> dict[str, Any]:
        payload = _validate_lighting_intent(intent)
        device_id = payload["device_id"]
        operation = _operation_for_action(intent.action)

        if not self.real_mode:
            return _backend_result(operation=operation, dry_run=True)

        if intent.action == "set_light_power":
            self.client.set_power(device_id, _validated_power(payload))
        elif intent.action == "set_light_brightness":
            self.client.set_brightness(device_id, _validated_brightness(payload))
        elif intent.action == "set_light_color_temperature":
            self.client.set_color_temperature(device_id, _validated_temperature(payload))
        else:
            raise ValueError(f"unsupported lighting action: {intent.action}")

        return _backend_result(operation=operation, dry_run=False)


def _validate_lighting_intent(intent: ActionIntent) -> dict[str, Any]:
    if intent.action not in _LIGHTING_ACTIONS:
        raise ValueError(f"unsupported lighting action: {intent.action}")

    inferred_domain = infer_device_domain(intent.scope, intent.action)
    if inferred_domain != _LIGHTING_DOMAIN:
        raise ValueError(f"scope is not a lighting scope: {intent.scope}")

    scope_parts = [part for part in str(intent.scope or "").casefold().split(".") if part]
    if _LIGHTING_DOMAIN not in scope_parts:
        raise ValueError("scope must contain lighting domain")

    payload = intent.payload if isinstance(intent.payload, dict) else {}
    payload_action_type = payload.get("action_type")
    payload_domain = payload.get("domain")
    payload_device_id = payload.get("device_id")

    if payload_action_type != intent.action:
        raise ValueError("payload action_type does not match intent action")
    if payload_domain != _LIGHTING_DOMAIN:
        raise ValueError("payload domain does not match lighting domain")
    if not isinstance(payload_device_id, str) or not payload_device_id.strip():
        raise ValueError("payload device_id is required")
    if payload_device_id.casefold() not in scope_parts:
        raise ValueError("payload device_id does not match scope")

    _validate_action_parameters(intent.action, payload)
    return payload


def _validate_action_parameters(action: str, payload: dict[str, Any]) -> None:
    if action == "set_light_power":
        _validated_power(payload)
        return
    if action == "set_light_brightness":
        _validated_brightness(payload)
        return
    if action == "set_light_color_temperature":
        _validated_temperature(payload)
        return
    raise ValueError(f"unsupported lighting action: {action}")


def _validated_power(payload: dict[str, Any]) -> str:
    power = payload.get("power")
    if not isinstance(power, str):
        raise ValueError("power is required")
    normalized = power.strip().casefold()
    if normalized not in _LIGHT_POWER_STATES:
        raise ValueError("power must be on or off")
    return normalized


def _validated_brightness(payload: dict[str, Any]) -> int:
    value = payload.get("brightness")
    if isinstance(value, bool):
        raise ValueError("brightness must be an integer between 0 and 100")
    try:
        brightness = int(str(value).strip())
    except (TypeError, ValueError):
        raise ValueError("brightness must be an integer between 0 and 100") from None
    if brightness < 0 or brightness > 100:
        raise ValueError("brightness must be an integer between 0 and 100")
    return brightness


def _validated_temperature(payload: dict[str, Any]) -> str:
    temperature = payload.get("temperature")
    if not isinstance(temperature, str):
        raise ValueError("temperature is required")
    normalized = temperature.strip().casefold()
    if normalized not in _LIGHT_TEMPERATURES:
        raise ValueError("temperature must be warm, neutral, or cool")
    return normalized


def _operation_for_action(action: str) -> str:
    return {
        "set_light_power": "set_power",
        "set_light_brightness": "set_brightness",
        "set_light_color_temperature": "set_color_temperature",
    }[action]


def _backend_result(*, operation: str, dry_run: bool) -> dict[str, Any]:
    return {
        "backend": "real_lighting",
        "accepted": True,
        "dry_run": dry_run,
        "device_id_present": True,
        "operation": operation,
    }
