"""Factory for productized device action execution."""

from __future__ import annotations

from pathlib import Path

from loguru import logger

from OriginAgent.agent.action_runtime import SafeActionExecutor
from OriginAgent.agent.audit import AuditLogger
from OriginAgent.agent.confirmation import ConfirmationManager
from OriginAgent.agent.facts import FactStore
from OriginAgent.config.schema import DeviceToolsConfig

from .action_safety import SmartHomeActionSafetyGate
from .confirmation_prompts import SmartHomeConfirmationPromptBuilder
from .device_actions import DeviceActionSchemaRegistry, TypedActionPlanner
from .device_backends import DeviceActionExecutor
from .device_integrations import RealLightingBackend
from .devices import DEVICE_SCOPE_REDACTOR
from .facts import SMART_HOME_FACT_STORE_CONFIG
from .permissions import PermissionResolver
from .presence import PresenceStore


class _NoopLightingClient:
    def set_power(self, device_id: str, power: str):
        return {"ok": True}

    def set_brightness(self, device_id: str, brightness: int):
        return {"ok": True}

    def set_color_temperature(self, device_id: str, temperature: str):
        return {"ok": True}


def build_device_action_executor(
    *,
    workspace: Path,
    config: DeviceToolsConfig,
    audit_logger: AuditLogger | None = None,
    presence_store: PresenceStore | None = None,
    fact_store: FactStore | None = None,
    permission_resolver: PermissionResolver | None = None,
) -> DeviceActionExecutor | None:
    if not config.enabled:
        return None
    if not config.lighting_enabled:
        return None
    if config.backend == "none":
        return None
    if config.mode == "real":
        logger.warning("Device gateway real mode is not enabled in this release; device tools disabled.")
        return None
    if config.backend != "fake":
        logger.warning("Unsupported device backend '{}'; device tools disabled", config.backend)
        return None

    audit = audit_logger or AuditLogger(workspace, scope_redactor=DEVICE_SCOPE_REDACTOR)
    presence = presence_store or PresenceStore(workspace)
    facts = fact_store or FactStore(workspace, config=SMART_HOME_FACT_STORE_CONFIG)
    permissions = permission_resolver or PermissionResolver()
    backend = RealLightingBackend(_NoopLightingClient(), real_mode=False)
    safe_executor = SafeActionExecutor(
        gate=SmartHomeActionSafetyGate(presence, facts),
        confirmation_manager=ConfirmationManager(
            workspace,
            audit_logger=audit,
            prompt_builder=SmartHomeConfirmationPromptBuilder(),
        ),
        backend=backend,
        permission_resolver=permissions,
        audit_logger=audit,
        scope_redactor=DEVICE_SCOPE_REDACTOR,
    )
    return DeviceActionExecutor(
        TypedActionPlanner(DeviceActionSchemaRegistry()),
        safe_executor,
        audit_logger=audit,
    )
