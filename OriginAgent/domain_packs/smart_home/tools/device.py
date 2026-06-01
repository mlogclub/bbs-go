"""Typed device tools exposed to model providers."""

from __future__ import annotations

import asyncio
from contextvars import ContextVar, Token
from typing import Any

from OriginAgent.agent.tools.base import Tool, tool_parameters
from OriginAgent.agent.tools.schema import IntegerSchema, StringSchema, tool_parameters_schema

from OriginAgent.domain_packs.smart_home.runtime.device_actions import TypedDeviceAction
from OriginAgent.domain_packs.smart_home.runtime.device_backends import DeviceActionExecutor
from OriginAgent.domain_packs.smart_home.runtime.devices import DeviceRegistry
from OriginAgent.domain_packs.smart_home.tools.device_messages import device_human_message


class _LightingToolBase(Tool):
    canonical_action_id = ""
    action_type = ""

    def __init__(
        self,
        executor: DeviceActionExecutor,
        *,
        device_registry: DeviceRegistry | None = None,
        real_mode: bool = False,
    ):
        self._executor = executor
        self._device_registry = device_registry
        self._real_mode = real_mode
        self._actor_id_var: ContextVar[str] = ContextVar(
            f"{self.name}_actor_id",
            default="",
        )
        self._trigger_var: ContextVar[str] = ContextVar(
            f"{self.name}_trigger",
            default="",
        )

    @classmethod
    def enabled(cls, ctx: Any) -> bool:
        device_config = getattr(getattr(ctx, "config", None), "device", None)
        return (
            getattr(ctx, "device_action_executor", None) is not None
            and bool(getattr(device_config, "enabled", False))
            and bool(getattr(device_config, "lighting_enabled", False))
        )

    @classmethod
    def create(cls, ctx: Any) -> Tool:
        device_config = getattr(getattr(ctx, "config", None), "device", None)
        real_mode = bool(getattr(device_config, "mode", None) == "real")
        return cls(
            ctx.device_action_executor,
            device_registry=getattr(ctx, "device_registry", None),
            real_mode=real_mode,
        )

    @property
    def description(self) -> str:
        return "Submit a typed low-risk lighting action through the OriginAgent device gateway."

    def set_context(self, actor_id: str, trigger: str) -> tuple[Token[str], Token[str]]:
        actor_token = self._actor_id_var.set(actor_id)
        trigger_token = self._trigger_var.set(trigger)
        return actor_token, trigger_token

    def reset_context(self, tokens: tuple[Token[str], Token[str]]) -> None:
        actor_token, trigger_token = tokens
        self._actor_id_var.reset(actor_token)
        self._trigger_var.reset(trigger_token)

    def _actor_id(self) -> str:
        actor_id = self._actor_id_var.get().strip()
        if not actor_id:
            raise PermissionError("device tool has no actor context")
        return actor_id

    def _trigger(self) -> str:
        trigger = self._trigger_var.get().strip()
        if not trigger:
            raise PermissionError("device tool has no trigger context")
        return trigger

    def _submit(
        self,
        *,
        device_id: str,
        room: str | None,
        parameters: dict[str, Any],
    ) -> dict[str, Any]:
        if self._real_mode:
            if self._device_registry is None:
                raise PermissionError("real-mode device tools require a device registry")
            record = self._device_registry.resolve(
                actor_id=self._actor_id(),
                domain="lighting",
                room=room,
                device_ref=device_id,
            )
            device_id = record.device_id
            room = record.room or room
        action = TypedDeviceAction(
            action_type=self.action_type,
            domain="lighting",
            device_id=device_id,
            room=room,
            parameters=parameters,
            requested_by=self._actor_id(),
            trigger=self._trigger(),
        )
        result = self._executor.submit_typed(action)
        return {
            "status": _tool_status(result.status),
            "execution_status": result.status,
            "action_id": result.action_id,
            "confirmation_id": result.confirmation_id,
            "backend_called": result.backend_called,
            "permission_status": result.permission_status,
            "human_message": device_human_message(result.status, result.reason),
        }

    async def _submit_async(
        self,
        *,
        device_id: str,
        room: str | None,
        parameters: dict[str, Any],
    ) -> dict[str, Any]:
        return await asyncio.to_thread(
            self._submit,
            device_id=device_id,
            room=room,
            parameters=parameters,
        )


def _tool_status(execution_status: str) -> str:
    if execution_status in {"executed", "dry_run"}:
        return "success"
    if execution_status in {"pending_confirmation", "needs_confirmation"}:
        return "pending_confirmation"
    if execution_status in {"denied", "deny"}:
        return "denied"
    return "failed"


@tool_parameters(
    tool_parameters_schema(
        device_id=StringSchema("Lighting device id", min_length=1, max_length=128),
        room=StringSchema("Optional room name", nullable=True),
        power=StringSchema("Power state", enum=["on", "off"]),
        required=["device_id", "power"],
        additional_properties=False,
    )
)
class LightingSetPowerTool(_LightingToolBase):
    canonical_action_id = "originagent.device.lighting.set_power"
    action_type = "set_light_power"
    name = "originagent_device_lighting_set_power"

    async def execute(
        self,
        device_id: str,
        power: str,
        room: str | None = None,
        **kwargs: Any,
    ) -> dict[str, Any]:
        return await self._submit_async(
            device_id=device_id,
            room=room,
            parameters={"power": power},
        )


@tool_parameters(
    tool_parameters_schema(
        device_id=StringSchema("Lighting device id", min_length=1, max_length=128),
        room=StringSchema("Optional room name", nullable=True),
        brightness=IntegerSchema(
            description="Brightness 0-100",
            minimum=0,
            maximum=100,
        ),
        required=["device_id", "brightness"],
        additional_properties=False,
    )
)
class LightingSetBrightnessTool(_LightingToolBase):
    canonical_action_id = "originagent.device.lighting.set_brightness"
    action_type = "set_light_brightness"
    name = "originagent_device_lighting_set_brightness"

    async def execute(
        self,
        device_id: str,
        brightness: int,
        room: str | None = None,
        **kwargs: Any,
    ) -> dict[str, Any]:
        return await self._submit_async(
            device_id=device_id,
            room=room,
            parameters={"brightness": brightness},
        )


@tool_parameters(
    tool_parameters_schema(
        device_id=StringSchema("Lighting device id", min_length=1, max_length=128),
        room=StringSchema("Optional room name", nullable=True),
        temperature=StringSchema(
            "Color temperature",
            enum=["warm", "neutral", "cool"],
        ),
        required=["device_id", "temperature"],
        additional_properties=False,
    )
)
class LightingSetColorTemperatureTool(_LightingToolBase):
    canonical_action_id = "originagent.device.lighting.set_color_temperature"
    action_type = "set_light_color_temperature"
    name = "originagent_device_lighting_set_color_temperature"

    async def execute(
        self,
        device_id: str,
        temperature: str,
        room: str | None = None,
        **kwargs: Any,
    ) -> dict[str, Any]:
        return await self._submit_async(
            device_id=device_id,
            room=room,
            parameters={"temperature": temperature},
        )


@tool_parameters(
    tool_parameters_schema(
        device_ref=StringSchema("Lighting device reference", min_length=1, max_length=128),
        room=StringSchema("Optional room name", nullable=True),
        power=StringSchema("Power state", enum=["on", "off"]),
        required=["device_ref", "power"],
        additional_properties=False,
    )
)
class RealModeLightingSetPowerTool(LightingSetPowerTool):
    async def execute(
        self,
        device_ref: str,
        power: str,
        room: str | None = None,
        **kwargs: Any,
    ) -> dict[str, Any]:
        return await self._submit_async(device_id=device_ref, room=room, parameters={"power": power})


@tool_parameters(
    tool_parameters_schema(
        device_ref=StringSchema("Lighting device reference", min_length=1, max_length=128),
        room=StringSchema("Optional room name", nullable=True),
        brightness=IntegerSchema(
            description="Brightness 0-100",
            minimum=0,
            maximum=100,
        ),
        required=["device_ref", "brightness"],
        additional_properties=False,
    )
)
class RealModeLightingSetBrightnessTool(LightingSetBrightnessTool):
    async def execute(
        self,
        device_ref: str,
        brightness: int,
        room: str | None = None,
        **kwargs: Any,
    ) -> dict[str, Any]:
        return await self._submit_async(
            device_id=device_ref,
            room=room,
            parameters={"brightness": brightness},
        )


@tool_parameters(
    tool_parameters_schema(
        device_ref=StringSchema("Lighting device reference", min_length=1, max_length=128),
        room=StringSchema("Optional room name", nullable=True),
        temperature=StringSchema(
            "Color temperature",
            enum=["warm", "neutral", "cool"],
        ),
        required=["device_ref", "temperature"],
        additional_properties=False,
    )
)
class RealModeLightingSetColorTemperatureTool(LightingSetColorTemperatureTool):
    async def execute(
        self,
        device_ref: str,
        temperature: str,
        room: str | None = None,
        **kwargs: Any,
    ) -> dict[str, Any]:
        return await self._submit_async(
            device_id=device_ref,
            room=room,
            parameters={"temperature": temperature},
        )


def lighting_tools(
    executor: DeviceActionExecutor,
    *,
    device_registry: DeviceRegistry | None = None,
    real_mode: bool = False,
) -> list[Tool]:
    if real_mode:
        return [
            RealModeLightingSetPowerTool(executor, device_registry=device_registry, real_mode=True),
            RealModeLightingSetBrightnessTool(executor, device_registry=device_registry, real_mode=True),
            RealModeLightingSetColorTemperatureTool(executor, device_registry=device_registry, real_mode=True),
        ]
    return [
        LightingSetPowerTool(executor, device_registry=device_registry, real_mode=real_mode),
        LightingSetBrightnessTool(executor, device_registry=device_registry, real_mode=real_mode),
        LightingSetColorTemperatureTool(executor, device_registry=device_registry, real_mode=real_mode),
    ]
