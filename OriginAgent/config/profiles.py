"""Runtime profile presets for common OriginAgent operating modes."""

from __future__ import annotations

from OriginAgent.config.schema import Config, DeviceToolsConfig, ExecToolConfig, RuntimeProfile, ToolAuditConfig


def build_runtime_profile_defaults(profile: RuntimeProfile) -> Config:
    """Return a config object containing only conservative profile defaults."""

    config = Config()
    config.runtime.profile = profile
    if profile == "default":
        return config
    if profile in {"safe", "household_safe"}:
        config.tools.audit = ToolAuditConfig(mode="minimal")
        config.tools.exec = ExecToolConfig(profile="secure", allow_unsafe_exec=False)
        config.tools.device = DeviceToolsConfig(enabled=False, mode="dry_run")
        return config
    if profile == "local_dev":
        config.tools.audit = ToolAuditConfig(mode="minimal")
        config.tools.exec = ExecToolConfig(profile="local_dev", allow_unsafe_exec=False)
        config.tools.device = DeviceToolsConfig(enabled=False, mode="dry_run")
        return config
    if profile == "automation":
        config.tools.audit = ToolAuditConfig(mode="minimal")
        config.tools.exec = ExecToolConfig(profile="secure", allow_unsafe_exec=False)
        config.tools.device = DeviceToolsConfig(enabled=False, mode="dry_run")
        return config
    return config


def apply_runtime_profile(config: Config) -> Config:
    """Apply conservative profile defaults without overriding explicit differences.

    This v1 helper intentionally treats the schema defaults as the only values
    eligible for replacement. It keeps user-specified values intact without
    requiring raw config field-presence tracking in the loader.
    """

    profile = config.runtime.profile
    if profile == "default":
        return config
    defaults = Config()
    profile_defaults = build_runtime_profile_defaults(profile)
    updated = config.model_copy(deep=True)
    if config.tools.audit == defaults.tools.audit:
        updated.tools.audit = profile_defaults.tools.audit
    if config.tools.exec == defaults.tools.exec:
        updated.tools.exec = profile_defaults.tools.exec
    if config.tools.device == defaults.tools.device:
        updated.tools.device = profile_defaults.tools.device
    return updated
