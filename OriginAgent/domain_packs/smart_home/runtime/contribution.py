"""Runtime contribution for the smart_home domain pack."""

from __future__ import annotations

from OriginAgent.agent.domain_packs import DomainRuntimeContribution

from OriginAgent.domain_packs.smart_home.runtime.device_factory import build_device_action_executor


def build_runtime_contribution(context) -> DomainRuntimeContribution:
    tools_config = getattr(context.config, "device", None)
    executor = context.overrides.get("device_action_executor")
    if executor is None and tools_config is not None:
        executor = build_device_action_executor(
            workspace=context.workspace,
            config=tools_config,
        )
    return DomainRuntimeContribution(
        tool_context={
            "device_action_executor": executor,
            "device_registry": context.overrides.get("device_registry"),
        }
    )
