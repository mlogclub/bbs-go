"""Helpers for runtime model preset selection."""

from __future__ import annotations

from collections.abc import Callable
from typing import Any

from OriginAgent.config.schema import ModelPresetConfig
from OriginAgent.providers.base import LLMProvider
from OriginAgent.providers.factory import ProviderSnapshot, build_provider_snapshot

PresetSnapshotLoader = Callable[[str], ProviderSnapshot]


def default_selection_signature(signature: tuple[object, ...] | None) -> tuple[object, ...] | None:
    return signature[:2] if signature else None


def configured_model_presets(config: Any) -> dict[str, ModelPresetConfig]:
    default_preset = config.resolve_default_preset()
    presets = {"default": default_preset}
    for name, preset in (config.model_presets or {}).items():
        presets[name] = _fill_from_default(preset, default_preset)
    return presets


def _fill_from_default(preset: ModelPresetConfig, default: ModelPresetConfig) -> ModelPresetConfig:
    return ModelPresetConfig(
        model=preset.model,
        provider=preset.provider,
        max_tokens=preset.max_tokens if preset.max_tokens is not None else default.max_tokens,
        context_window_tokens=preset.context_window_tokens
        if preset.context_window_tokens is not None
        else default.context_window_tokens,
        temperature=preset.temperature if preset.temperature is not None else default.temperature,
        reasoning_effort=preset.reasoning_effort
        if preset.reasoning_effort is not None
        else default.reasoning_effort,
        fallback_models=preset.fallback_models,
    )


def make_preset_snapshot_loader(
    config: Any,
    provider_snapshot_loader: Callable[..., ProviderSnapshot] | None,
) -> PresetSnapshotLoader:
    if provider_snapshot_loader is not None:
        return lambda name: provider_snapshot_loader(preset_name=name)
    return lambda name: build_provider_snapshot(config, preset_name=name)


def build_static_preset_snapshot(
    provider: LLMProvider,
    name: str,
    preset: ModelPresetConfig,
) -> ProviderSnapshot:
    provider.generation = preset.to_generation_settings()
    return ProviderSnapshot(
        provider=provider,
        model=preset.model,
        context_window_tokens=preset.context_window_tokens or 0,
        signature=("model_preset", name, preset.model_dump_json()),
    )


def build_runtime_preset_snapshot(
    *,
    name: str,
    presets: dict[str, ModelPresetConfig],
    provider: LLMProvider,
    loader: PresetSnapshotLoader | None,
) -> ProviderSnapshot:
    if loader is not None:
        return loader(name)
    return build_static_preset_snapshot(provider, name, presets[name])


def normalize_preset_name(name: str | None, presets: dict[str, ModelPresetConfig]) -> str:
    if not isinstance(name, str) or not name.strip():
        raise ValueError("model_preset must be a non-empty string")
    name = name.strip()
    if name not in presets:
        raise KeyError(f"model_preset {name!r} not found. Available: {', '.join(presets) or '(none)'}")
    return name
