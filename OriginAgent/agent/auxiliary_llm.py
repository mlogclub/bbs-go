"""Task-aware auxiliary LLM routing for background jobs.

This module keeps OriginAgent's existing provider abstractions and adds a small
fallback router for background work such as consolidation and Dream.  It does
not introduce a second provider stack: candidates are still ordinary
``LLMProvider`` instances built by the existing provider factory.
"""

from __future__ import annotations

import asyncio
import time
from collections.abc import Awaitable, Callable
from dataclasses import dataclass
from typing import Any

from loguru import logger

from OriginAgent.config.schema import (
    AuxiliaryConfig,
    AuxiliaryTaskConfig,
    FallbackCandidate,
    InlineFallbackConfig,
    ModelPresetConfig,
)
from OriginAgent.providers.base import GenerationSettings, LLMProvider, LLMResponse


_OPENROUTER_AUX_MODEL = "google/gemini-2.5-flash"

_AUX_DEFAULT_MODELS: dict[str, str] = {
    "deepseek": "deepseek-chat",
    "openai": "gpt-4o-mini",
    "gemini": "gemini-2.0-flash",
    "zhipu": "glm-4-flash",
    "dashscope": "qwen-turbo",
    "moonshot": "moonshot-v1-8k",
    "minimax": "MiniMax-Text-01",
    "mistral": "mistral-small-latest",
    "stepfun": "step-1-8k",
    "groq": "llama-3.1-8b-instant",
}

_PAYMENT_TOKENS = (
    "402",
    "insufficient_quota",
    "quota_exhausted",
    "quota exceeded",
    "quota_exceeded",
    "billing_hard_limit",
    "billing hard limit",
    "insufficient_balance",
    "insufficient balance",
    "credit balance too low",
    "payment required",
    "payment_required",
    "out of credits",
    "out of quota",
)

_NON_FALLBACK_STATUSES = frozenset({400, 401, 403, 404, 422})
_NON_FALLBACK_TOKENS = (
    "authentication",
    "auth",
    "permission",
    "invalid_request",
    "context_length",
    "content_filter",
    "refusal",
)


ProviderFactory = Callable[[ModelPresetConfig], LLMProvider]


@dataclass(slots=True)
class _Candidate:
    key: str
    provider: LLMProvider
    model: str
    source: str
    max_tokens: int | None = None
    temperature: float | None = None
    reasoning_effort: str | None = None


class AuxiliaryLLMRouter:
    """Route background LLM calls through task-aware fallback candidates."""

    def __init__(
        self,
        *,
        primary_provider: LLMProvider,
        primary_model: str,
        auxiliary_config: AuxiliaryConfig | None = None,
        config: Any | None = None,
        provider_factory: ProviderFactory | None = None,
        primary_provider_name: str | None = None,
        clock: Callable[[], float] = time.monotonic,
    ):
        self.primary_provider = primary_provider
        self.primary_model = primary_model
        self.primary_provider_name = primary_provider_name or "primary"
        self.config = config
        self.auxiliary_config = auxiliary_config or AuxiliaryConfig()
        self._provider_factory = provider_factory
        self._clock = clock
        self._unhealthy_until: dict[str, float] = {}

    def set_primary(
        self,
        provider: LLMProvider,
        model: str,
        provider_name: str | None = None,
    ) -> None:
        changed = provider is not self.primary_provider or model != self.primary_model
        self.primary_provider = provider
        self.primary_model = model
        if provider_name:
            self.primary_provider_name = provider_name
        if changed:
            self._unhealthy_until.pop("primary", None)

    def task_provider(self, task: str) -> "AuxiliaryTaskProvider":
        return AuxiliaryTaskProvider(self, task)

    def _task_config(self, task: str | None) -> AuxiliaryTaskConfig:
        if not task:
            return AuxiliaryTaskConfig()
        return self.auxiliary_config.tasks.get(task, AuxiliaryTaskConfig())

    def _candidate_model(self, task_cfg: AuxiliaryTaskConfig, requested: str | None) -> str:
        return task_cfg.model_override or requested or self.primary_model

    def _candidate_key(self, source: str, provider_name: str, model: str) -> str:
        if source == "primary":
            return "primary"
        return f"{provider_name}:{model}"

    def _is_unhealthy(self, key: str) -> bool:
        until = self._unhealthy_until.get(key)
        if until is None:
            return False
        if until <= self._clock():
            self._unhealthy_until.pop(key, None)
            return False
        return True

    def _mark_unhealthy(self, key: str, ttl: float, reason: str) -> None:
        if ttl <= 0:
            return
        self._unhealthy_until[key] = self._clock() + ttl
        logger.warning(
            "Auxiliary LLM candidate '{}' marked unhealthy for {}s ({})",
            key,
            int(ttl),
            reason,
        )

    def _provider_factory_or_none(self) -> ProviderFactory | None:
        if self._provider_factory:
            return self._provider_factory
        if self.config is None:
            return None

        def _factory(preset: ModelPresetConfig) -> LLMProvider:
            from OriginAgent.providers.factory import make_plain_provider

            return make_plain_provider(self.config, preset)

        return _factory

    def _preset_from_fallback(self, fallback: FallbackCandidate) -> ModelPresetConfig | None:
        if isinstance(fallback, str):
            if self.config is None:
                return None
            return self.config.resolve_preset(fallback)
        if isinstance(fallback, InlineFallbackConfig):
            return ModelPresetConfig(
                model=fallback.model,
                provider=fallback.provider,
                max_tokens=fallback.max_tokens,
                context_window_tokens=fallback.context_window_tokens,
                temperature=fallback.temperature,
                reasoning_effort=fallback.reasoning_effort,
            )
        return None

    def _configured_fallbacks(self, task_cfg: AuxiliaryTaskConfig) -> list[ModelPresetConfig]:
        raw = list(task_cfg.fallback_models)
        if not raw and self.config is not None:
            raw = list(getattr(self.config.agents.defaults, "fallback_models", []))

        presets: list[ModelPresetConfig] = []
        for fallback in raw:
            try:
                preset = self._preset_from_fallback(fallback)
            except Exception:
                logger.exception("Auxiliary LLM: failed to resolve fallback preset")
                continue
            if preset is not None:
                presets.append(preset)
        return presets

    def _default_api_key_fallbacks(self) -> list[ModelPresetConfig]:
        if self.config is None:
            return []
        try:
            from OriginAgent.providers.registry import PROVIDERS
        except Exception:
            return []

        presets: list[ModelPresetConfig] = []
        primary_name = (self.primary_provider_name or "").replace("-", "_")
        for spec in PROVIDERS:
            if spec.name == primary_name:
                continue
            provider_cfg = getattr(self.config.providers, spec.name, None)
            if not provider_cfg or not getattr(provider_cfg, "api_key", None):
                continue
            if spec.name == "openrouter":
                presets.append(ModelPresetConfig(
                    model=_OPENROUTER_AUX_MODEL,
                    provider="openrouter",
                ))
                continue
            model = _AUX_DEFAULT_MODELS.get(spec.name)
            if model:
                presets.append(ModelPresetConfig(model=model, provider=spec.name))
        return presets

    def _candidates(
        self,
        task: str | None,
        requested_model: str | None,
    ) -> list[_Candidate]:
        task_cfg = self._task_config(task)
        primary_model = self._candidate_model(task_cfg, requested_model)
        candidates = [
            _Candidate(
                key="primary",
                provider=self.primary_provider,
                model=primary_model,
                source="primary",
            )
        ]
        if not self.auxiliary_config.enabled:
            return candidates

        factory = self._provider_factory_or_none()
        if factory is None:
            return candidates

        seen = {"primary"}
        fallback_presets = [
            *self._configured_fallbacks(task_cfg),
            *self._default_api_key_fallbacks(),
        ]
        for preset in fallback_presets:
            provider_name = preset.provider or "auto"
            key = self._candidate_key("fallback", provider_name, preset.model)
            if key in seen:
                continue
            seen.add(key)
            try:
                provider = factory(preset)
            except Exception as exc:
                logger.warning(
                    "Auxiliary LLM: failed to create fallback '{}' ({}): {}",
                    preset.model,
                    provider_name,
                    exc,
                )
                continue
            candidates.append(_Candidate(
                key=key,
                provider=provider,
                model=preset.model,
                source=provider_name,
                max_tokens=preset.max_tokens,
                temperature=preset.temperature,
                reasoning_effort=preset.reasoning_effort,
            ))
        return candidates

    @staticmethod
    def _is_payment_error(response: LLMResponse) -> bool:
        if response.error_status_code == 402:
            return True
        haystack = " ".join(
            str(value or "").lower()
            for value in (
                response.error_type,
                response.error_code,
                response.error_kind,
                response.content,
            )
        )
        return any(token in haystack for token in _PAYMENT_TOKENS)

    @staticmethod
    def _is_non_fallback_error(response: LLMResponse) -> bool:
        if response.error_status_code in _NON_FALLBACK_STATUSES:
            return True
        haystack = " ".join(
            str(value or "").lower()
            for value in (
                response.error_type,
                response.error_code,
                response.error_kind,
                response.content,
            )
        )
        return any(token in haystack for token in _NON_FALLBACK_TOKENS)

    @staticmethod
    def _is_transient_fallback_error(response: LLMResponse) -> bool:
        if response.error_status_code and response.error_status_code >= 500:
            return True
        return LLMProvider._is_transient_response(response)

    def _fallback_reason(self, response: LLMResponse) -> str | None:
        if self._is_payment_error(response):
            return "payment"
        if self._is_non_fallback_error(response):
            return None
        if self._is_transient_fallback_error(response):
            return "transient"
        return None

    def _cooldown_for(self, response: LLMResponse, reason: str) -> float:
        if reason == "payment":
            return float(self.auxiliary_config.payment_cooldown_s)
        retry_after = LLMProvider._extract_retry_after_from_response(response)
        if retry_after is not None:
            return max(retry_after, 0.0)
        return float(self.auxiliary_config.transient_cooldown_s)

    async def _invoke_candidate(
        self,
        candidate: _Candidate,
        *,
        messages: list[dict[str, Any]],
        tools: list[dict[str, Any]] | None,
        max_tokens: int | None,
        temperature: float | None,
        reasoning_effort: str | None,
        tool_choice: str | dict[str, Any] | None,
        retry_mode: str,
        on_retry_wait: Callable[[str], Awaitable[None]] | None,
        timeout_s: float | None,
        stream: bool,
        on_content_delta: Callable[[str], Awaitable[None]] | None,
        on_thinking_delta: Callable[[str], Awaitable[None]] | None,
    ) -> LLMResponse:
        kwargs: dict[str, Any] = {
            "messages": messages,
            "tools": tools,
            "model": candidate.model,
            "retry_mode": retry_mode,
            "on_retry_wait": on_retry_wait,
        }
        effective_max_tokens = (
            candidate.max_tokens
            if candidate.source != "primary" and candidate.max_tokens is not None
            else max_tokens
        )
        effective_temperature = (
            candidate.temperature
            if candidate.source != "primary" and candidate.temperature is not None
            else temperature
        )
        effective_reasoning_effort = (
            candidate.reasoning_effort
            if candidate.source != "primary" and candidate.reasoning_effort is not None
            else reasoning_effort
        )
        if effective_max_tokens is not None:
            kwargs["max_tokens"] = effective_max_tokens
        if effective_temperature is not None:
            kwargs["temperature"] = effective_temperature
        if effective_reasoning_effort is not None:
            kwargs["reasoning_effort"] = effective_reasoning_effort
        if tool_choice is not None:
            kwargs["tool_choice"] = tool_choice

        if stream:
            kwargs["on_content_delta"] = on_content_delta
            kwargs["on_thinking_delta"] = on_thinking_delta
            coro = candidate.provider.chat_stream_with_retry(**kwargs)
        else:
            coro = candidate.provider.chat_with_retry(**kwargs)

        if timeout_s is None or timeout_s <= 0 or stream:
            return await coro
        try:
            return await asyncio.wait_for(coro, timeout=timeout_s)
        except asyncio.TimeoutError:
            return LLMResponse(
                content=f"Error calling LLM: timed out after {timeout_s:g}s",
                finish_reason="error",
                error_kind="timeout",
            )

    async def call_llm(
        self,
        task: str,
        *,
        messages: list[dict[str, Any]],
        tools: list[dict[str, Any]] | None = None,
        model: str | None = None,
        max_tokens: int | None = None,
        temperature: float | None = None,
        reasoning_effort: str | None = None,
        tool_choice: str | dict[str, Any] | None = None,
        retry_mode: str = "standard",
        on_retry_wait: Callable[[str], Awaitable[None]] | None = None,
        timeout: float | None = None,
        stream: bool = False,
        on_content_delta: Callable[[str], Awaitable[None]] | None = None,
        on_thinking_delta: Callable[[str], Awaitable[None]] | None = None,
    ) -> LLMResponse:
        task_cfg = self._task_config(task)
        timeout_s = timeout if timeout is not None else task_cfg.timeout_s
        last_response: LLMResponse | None = None
        candidates = self._candidates(task, model)
        for candidate in candidates:
            if self._is_unhealthy(candidate.key):
                logger.info(
                    "Auxiliary LLM task '{}' skipping unhealthy candidate '{}'",
                    task,
                    candidate.key,
                )
                continue
            try:
                response = await self._invoke_candidate(
                    candidate,
                    messages=messages,
                    tools=tools,
                    max_tokens=max_tokens,
                    temperature=temperature,
                    reasoning_effort=reasoning_effort,
                    tool_choice=tool_choice,
                    retry_mode=retry_mode,
                    on_retry_wait=on_retry_wait,
                    timeout_s=timeout_s,
                    stream=stream,
                    on_content_delta=on_content_delta,
                    on_thinking_delta=on_thinking_delta,
                )
            except asyncio.CancelledError:
                raise
            except Exception as exc:
                response = LLMResponse(
                    content=f"Error calling LLM: {exc}",
                    finish_reason="error",
                )

            if response.finish_reason != "error":
                return response

            last_response = response
            reason = self._fallback_reason(response)
            if reason is None:
                return response
            self._mark_unhealthy(
                candidate.key,
                self._cooldown_for(response, reason),
                reason,
            )
            logger.warning(
                "Auxiliary LLM task '{}' failed on '{}' ({}), trying fallback",
                task,
                candidate.key,
                reason,
            )

        return last_response or LLMResponse(
            content=f"No LLM provider configured for auxiliary task '{task}'",
            finish_reason="error",
            error_kind="configuration",
            error_should_retry=False,
        )


class AuxiliaryTaskProvider(LLMProvider):
    """LLMProvider facade that routes all calls as one auxiliary task."""

    def __init__(self, router: AuxiliaryLLMRouter, task: str):
        super().__init__()
        self.router = router
        self.task = task
        primary_generation = getattr(router.primary_provider, "generation", None)
        self.generation = GenerationSettings(
            temperature=getattr(primary_generation, "temperature", 0.7),
            max_tokens=getattr(primary_generation, "max_tokens", 4096),
            reasoning_effort=getattr(primary_generation, "reasoning_effort", None),
        )

    @property
    def supports_progress_deltas(self) -> bool:
        return bool(getattr(self.router.primary_provider, "supports_progress_deltas", False))

    def get_default_model(self) -> str:
        return self.router.primary_model

    async def chat(
        self,
        messages: list[dict[str, Any]],
        tools: list[dict[str, Any]] | None = None,
        model: str | None = None,
        max_tokens: int = 4096,
        temperature: float = 0.7,
        reasoning_effort: str | None = None,
        tool_choice: str | dict[str, Any] | None = None,
    ) -> LLMResponse:
        return await self.router.call_llm(
            self.task,
            messages=messages,
            tools=tools,
            model=model,
            max_tokens=max_tokens,
            temperature=temperature,
            reasoning_effort=reasoning_effort,
            tool_choice=tool_choice,
        )

    async def chat_with_retry(
        self,
        messages: list[dict[str, Any]],
        tools: list[dict[str, Any]] | None = None,
        model: str | None = None,
        max_tokens: object = LLMProvider._SENTINEL,
        temperature: object = LLMProvider._SENTINEL,
        reasoning_effort: object = LLMProvider._SENTINEL,
        tool_choice: str | dict[str, Any] | None = None,
        retry_mode: str = "standard",
        on_retry_wait: Callable[[str], Awaitable[None]] | None = None,
    ) -> LLMResponse:
        if max_tokens is self._SENTINEL or max_tokens is None:
            max_tokens = self.generation.max_tokens
        if temperature is self._SENTINEL or temperature is None:
            temperature = self.generation.temperature
        if reasoning_effort is self._SENTINEL:
            reasoning_effort = self.generation.reasoning_effort
        return await self.router.call_llm(
            self.task,
            messages=messages,
            tools=tools,
            model=model,
            max_tokens=max_tokens,
            temperature=temperature,
            reasoning_effort=reasoning_effort,
            tool_choice=tool_choice,
            retry_mode=retry_mode,
            on_retry_wait=on_retry_wait,
        )

    async def chat_stream_with_retry(
        self,
        messages: list[dict[str, Any]],
        tools: list[dict[str, Any]] | None = None,
        model: str | None = None,
        max_tokens: object = LLMProvider._SENTINEL,
        temperature: object = LLMProvider._SENTINEL,
        reasoning_effort: object = LLMProvider._SENTINEL,
        tool_choice: str | dict[str, Any] | None = None,
        on_content_delta: Callable[[str], Awaitable[None]] | None = None,
        on_thinking_delta: Callable[[str], Awaitable[None]] | None = None,
        retry_mode: str = "standard",
        on_retry_wait: Callable[[str], Awaitable[None]] | None = None,
    ) -> LLMResponse:
        if max_tokens is self._SENTINEL or max_tokens is None:
            max_tokens = self.generation.max_tokens
        if temperature is self._SENTINEL or temperature is None:
            temperature = self.generation.temperature
        if reasoning_effort is self._SENTINEL:
            reasoning_effort = self.generation.reasoning_effort
        return await self.router.call_llm(
            self.task,
            messages=messages,
            tools=tools,
            model=model,
            max_tokens=max_tokens,
            temperature=temperature,
            reasoning_effort=reasoning_effort,
            tool_choice=tool_choice,
            retry_mode=retry_mode,
            on_retry_wait=on_retry_wait,
            stream=True,
            on_content_delta=on_content_delta,
            on_thinking_delta=on_thinking_delta,
        )


async def call_llm(
    task: str,
    *,
    router: AuxiliaryLLMRouter | None = None,
    provider: LLMProvider | None = None,
    messages: list[dict[str, Any]],
    tools: list[dict[str, Any]] | None = None,
    model: str | None = None,
    max_tokens: int | None = None,
    temperature: float | None = None,
    reasoning_effort: str | None = None,
    tool_choice: str | dict[str, Any] | None = None,
    retry_mode: str = "standard",
    on_retry_wait: Callable[[str], Awaitable[None]] | None = None,
    timeout: float | None = None,
) -> LLMResponse:
    """Call a task-scoped auxiliary LLM route, or the direct provider fallback."""

    if router is not None:
        return await router.call_llm(
            task,
            messages=messages,
            tools=tools,
            model=model,
            max_tokens=max_tokens,
            temperature=temperature,
            reasoning_effort=reasoning_effort,
            tool_choice=tool_choice,
            retry_mode=retry_mode,
            on_retry_wait=on_retry_wait,
            timeout=timeout,
        )
    if provider is None:
        return LLMResponse(
            content=f"No LLM provider configured for auxiliary task '{task}'",
            finish_reason="error",
            error_kind="configuration",
            error_should_retry=False,
        )
    kwargs: dict[str, Any] = {
        "messages": messages,
        "tools": tools,
        "model": model,
        "retry_mode": retry_mode,
        "on_retry_wait": on_retry_wait,
    }
    if max_tokens is not None:
        kwargs["max_tokens"] = max_tokens
    if temperature is not None:
        kwargs["temperature"] = temperature
    if reasoning_effort is not None:
        kwargs["reasoning_effort"] = reasoning_effort
    if tool_choice is not None:
        kwargs["tool_choice"] = tool_choice
    coro = provider.chat_with_retry(**kwargs)
    if timeout is None or timeout <= 0:
        return await coro
    try:
        return await asyncio.wait_for(coro, timeout=timeout)
    except asyncio.TimeoutError:
        return LLMResponse(
            content=f"Error calling LLM: timed out after {timeout:g}s",
            finish_reason="error",
            error_kind="timeout",
        )
