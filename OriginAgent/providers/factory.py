"""Create LLM providers from config."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from OriginAgent.config.schema import Config, InlineFallbackConfig, ModelPresetConfig
from OriginAgent.providers.base import GenerationSettings, LLMProvider
from OriginAgent.providers.registry import find_by_name


@dataclass(frozen=True)
class ProviderSnapshot:
    provider: LLMProvider
    model: str
    context_window_tokens: int
    signature: tuple[object, ...]


def _provider_config_for(
    config: Config,
    *,
    model: str,
    provider_name: str,
):
    if provider_name and provider_name != "auto":
        spec = find_by_name(provider_name)
        if spec:
            return getattr(config.providers, spec.name, None), spec.name, spec
        return None, None, None
    matched_name = config.get_provider_name(model)
    spec = find_by_name(matched_name) if matched_name else None
    return config.get_provider(model), matched_name, spec


def _fill_preset_defaults(config: Config, preset: ModelPresetConfig) -> ModelPresetConfig:
    defaults = config.agents.defaults
    return ModelPresetConfig(
        model=preset.model,
        provider=preset.provider or "auto",
        max_tokens=preset.max_tokens if preset.max_tokens is not None else defaults.max_tokens,
        context_window_tokens=(
            preset.context_window_tokens
            if preset.context_window_tokens is not None
            else defaults.context_window_tokens
        ),
        temperature=preset.temperature if preset.temperature is not None else defaults.temperature,
        reasoning_effort=(
            preset.reasoning_effort
            if preset.reasoning_effort is not None
            else defaults.reasoning_effort
        ),
        fallback_models=preset.fallback_models,
    )


def _api_base_for(config: Config, model: str, p: object, spec: object) -> str | None:
    if p and getattr(p, "api_base", None):
        return p.api_base
    if spec and getattr(spec, "default_api_base", ""):
        return spec.default_api_base
    return config.get_api_base(model)


def _make_plain_provider(config: Config, preset: ModelPresetConfig) -> LLMProvider:
    """Create the LLM provider implied by config."""
    preset = _fill_preset_defaults(config, preset)
    model = preset.model
    p, provider_name, spec = _provider_config_for(
        config,
        model=model,
        provider_name=preset.provider,
    )
    backend = spec.backend if spec else "openai_compat"

    if backend == "azure_openai":
        if not p or not p.api_key or not p.api_base:
            raise ValueError("Azure OpenAI requires api_key and api_base in config.")
    elif backend == "openai_compat" and not model.startswith("bedrock/"):
        needs_key = not (p and p.api_key)
        exempt = spec and (spec.is_oauth or spec.is_local or spec.is_direct)
        if needs_key and not exempt:
            raise ValueError(f"No API key configured for provider '{provider_name}'.")

    if backend == "openai_codex":
        from OriginAgent.providers.openai_codex_provider import OpenAICodexProvider

        provider = OpenAICodexProvider(default_model=model)
    elif backend == "azure_openai":
        from OriginAgent.providers.azure_openai_provider import AzureOpenAIProvider

        provider = AzureOpenAIProvider(
            api_key=p.api_key,
            api_base=p.api_base,
            default_model=model,
        )
    elif backend == "github_copilot":
        from OriginAgent.providers.github_copilot_provider import GitHubCopilotProvider

        provider = GitHubCopilotProvider(default_model=model)
    elif backend == "anthropic":
        from OriginAgent.providers.anthropic_provider import AnthropicProvider

        provider = AnthropicProvider(
            api_key=p.api_key if p else None,
            api_base=_api_base_for(config, model, p, spec),
            default_model=model,
            extra_headers=p.extra_headers if p else None,
        )
    elif backend == "bedrock":
        from OriginAgent.providers.bedrock_provider import BedrockProvider

        provider = BedrockProvider(
            api_key=p.api_key if p else None,
            api_base=p.api_base if p else None,
            default_model=model,
            region=getattr(p, "region", None) if p else None,
            profile=getattr(p, "profile", None) if p else None,
            extra_body=p.extra_body if p else None,
        )
    else:
        from OriginAgent.providers.openai_compat_provider import OpenAICompatProvider

        provider = OpenAICompatProvider(
            api_key=p.api_key if p else None,
            api_base=_api_base_for(config, model, p, spec),
            default_model=model,
            extra_headers=p.extra_headers if p else None,
            spec=spec,
            extra_body=p.extra_body if p else None,
        )

    provider.generation = GenerationSettings(
        temperature=preset.temperature,
        max_tokens=preset.max_tokens,
        reasoning_effort=preset.reasoning_effort,
    )
    return provider


def make_plain_provider(config: Config, preset: ModelPresetConfig) -> LLMProvider:
    """Create one provider from a preset without wrapping fallback models."""
    return _make_plain_provider(config, preset)


def _resolve_fallback_presets(config: Config, primary: ModelPresetConfig) -> list[ModelPresetConfig]:
    presets: list[ModelPresetConfig] = []
    for fallback in primary.fallback_models:
        if isinstance(fallback, str):
            presets.append(_fill_preset_defaults(config, config.resolve_preset(fallback)))
        elif isinstance(fallback, InlineFallbackConfig):
            presets.append(
                _fill_preset_defaults(
                    config,
                    ModelPresetConfig(
                        model=fallback.model,
                        provider=fallback.provider,
                        max_tokens=fallback.max_tokens,
                        context_window_tokens=fallback.context_window_tokens,
                        temperature=fallback.temperature,
                        reasoning_effort=fallback.reasoning_effort,
                    ),
                )
            )
    return [preset for preset in presets if preset.model != primary.model or preset.provider != primary.provider]


def make_provider(config: Config, preset: ModelPresetConfig | None = None) -> LLMProvider:
    """Create the LLM provider implied by config and optional model preset."""
    primary = _fill_preset_defaults(config, preset or config.resolve_preset())
    provider = _make_plain_provider(config, primary)
    fallback_presets = _resolve_fallback_presets(config, primary)
    if not fallback_presets:
        return provider

    from OriginAgent.providers.fallback_provider import FallbackProvider

    return FallbackProvider(
        primary=provider,
        fallback_presets=fallback_presets,
        provider_factory=lambda fallback: _make_plain_provider(config, fallback),
    )


def provider_signature(config: Config, preset: ModelPresetConfig | None = None) -> tuple[object, ...]:
    """Return the config fields that affect the primary LLM provider."""
    resolved = _fill_preset_defaults(config, preset or config.resolve_preset())
    model = resolved.model
    p, provider_name, _ = _provider_config_for(
        config,
        model=model,
        provider_name=resolved.provider,
    )
    return (
        model,
        resolved.provider,
        provider_name,
        p.api_key if p else None,
        _api_base_for(config, model, p, find_by_name(provider_name) if provider_name else None),
        p.extra_headers if p else None,
        p.extra_body if p else None,
        getattr(p, "region", None) if p else None,
        getattr(p, "profile", None) if p else None,
        resolved.max_tokens,
        resolved.temperature,
        resolved.reasoning_effort,
        resolved.context_window_tokens,
        tuple(item if isinstance(item, str) else item.model_dump_json() for item in resolved.fallback_models),
    )


def build_provider_snapshot(config: Config, preset_name: str | None = None) -> ProviderSnapshot:
    preset = _fill_preset_defaults(config, config.resolve_preset(preset_name))
    return ProviderSnapshot(
        provider=make_provider(config, preset),
        model=preset.model,
        context_window_tokens=preset.context_window_tokens or config.agents.defaults.context_window_tokens,
        signature=provider_signature(config, preset),
    )


def load_provider_snapshot(
    config_path: Path | None = None,
    *,
    preset_name: str | None = None,
) -> ProviderSnapshot:
    from OriginAgent.config.loader import load_config, resolve_config_env_vars

    return build_provider_snapshot(
        resolve_config_env_vars(load_config(config_path)),
        preset_name=preset_name,
    )
