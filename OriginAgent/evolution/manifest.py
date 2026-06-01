"""Evolution module manifest validation."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any

MODULE_SCHEMA_VERSION = "originagent.evolution.module.v1"
ALLOWED_MODULE_TYPES = frozenset({"skill", "workflow", "domain_pack", "tool"})


class ManifestValidationError(ValueError):
    """Raised when an evolution module manifest is invalid."""


@dataclass(frozen=True)
class EvolutionManifest:
    """Validated local evolution module manifest."""

    schema_version: str
    module_id: str
    module_type: str
    version: str
    target_originagent: dict[str, Any] = field(default_factory=dict)
    target_module_api: str = ""
    permissions: dict[str, Any] = field(default_factory=dict)
    runtime_requirements: dict[str, Any] = field(default_factory=dict)
    context_budget: dict[str, Any] = field(default_factory=dict)
    external_endpoints: tuple[str, ...] = ()
    external_side_effects: dict[str, Any] = field(default_factory=dict)
    tests: dict[str, Any] = field(default_factory=dict)


def validate_manifest(raw: dict[str, Any]) -> EvolutionManifest:
    """Validate an already parsed evolution module manifest."""

    if not isinstance(raw, dict):
        raise ManifestValidationError("manifest must be a mapping")

    schema_version = _required_string(raw, "schema_version")
    if schema_version != MODULE_SCHEMA_VERSION:
        raise ManifestValidationError(f"unsupported schema_version: {schema_version}")

    module_id = _required_string(raw, "module_id")
    module_type = _required_string(raw, "module_type")
    if module_type not in ALLOWED_MODULE_TYPES:
        allowed = ", ".join(sorted(ALLOWED_MODULE_TYPES))
        raise ManifestValidationError(f"module_type must be one of: {allowed}")

    version = _required_string(raw, "version")

    return EvolutionManifest(
        schema_version=schema_version,
        module_id=module_id,
        module_type=module_type,
        version=version,
        target_originagent=_optional_mapping(raw, "target_originagent"),
        target_module_api=_optional_string(raw, "target_module_api"),
        permissions=_optional_mapping(raw, "permissions"),
        runtime_requirements=_optional_mapping(raw, "runtime_requirements"),
        context_budget=_optional_mapping(raw, "context_budget"),
        external_endpoints=_optional_string_tuple(raw, "external_endpoints"),
        external_side_effects=_optional_mapping(raw, "external_side_effects"),
        tests=_optional_mapping(raw, "tests"),
    )


def _required_string(raw: dict[str, Any], key: str) -> str:
    value = raw.get(key)
    if not isinstance(value, str) or not value.strip():
        raise ManifestValidationError(f"missing required field: {key}")
    return value.strip()


def _optional_string(raw: dict[str, Any], key: str) -> str:
    value = raw.get(key)
    if value is None:
        return ""
    if not isinstance(value, str):
        raise ManifestValidationError(f"{key} must be a string")
    return value.strip()


def _optional_mapping(raw: dict[str, Any], key: str) -> dict[str, Any]:
    value = raw.get(key)
    if value is None:
        return {}
    if not isinstance(value, dict):
        raise ManifestValidationError(f"{key} must be a mapping")
    return dict(value)


def _optional_string_tuple(raw: dict[str, Any], key: str) -> tuple[str, ...]:
    value = raw.get(key)
    if value is None:
        return ()
    if not isinstance(value, list):
        raise ManifestValidationError(f"{key} must be a list")
    endpoints: list[str] = []
    for item in value:
        if not isinstance(item, str) or not item.strip():
            raise ManifestValidationError(f"{key} entries must be non-empty strings")
        endpoints.append(item.strip())
    return tuple(endpoints)
