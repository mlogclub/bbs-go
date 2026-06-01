"""Provider wrapper that transparently fails over to fallback models."""

from __future__ import annotations

import time
from collections.abc import Awaitable, Callable
from typing import Any

from loguru import logger

from OriginAgent.providers.base import LLMProvider, LLMResponse

_PRIMARY_FAILURE_THRESHOLD = 3
_PRIMARY_COOLDOWN_S = 60
_MISSING = object()
_FALLBACK_ERROR_KINDS = frozenset({"timeout", "connection", "server_error", "rate_limit", "overloaded"})
_NON_FALLBACK_ERROR_KINDS = frozenset({
    "authentication",
    "auth",
    "permission",
    "content_filter",
    "refusal",
    "context_length",
    "invalid_request",
})
_FALLBACK_ERROR_TOKENS = (
    "rate limit",
    "too_many_requests",
    "overloaded",
    "server error",
    "temporarily unavailable",
    "timeout",
    "timed out",
    "connection",
    "insufficient_quota",
    "quota",
    "billing_hard_limit",
    "insufficient_balance",
    "out of credits",
)


class FallbackProvider(LLMProvider):
    """Wrap a primary provider and fail over before any fallback content streams."""

    def __init__(
        self,
        primary: LLMProvider,
        fallback_presets: list[Any],
        provider_factory: Callable[[Any], LLMProvider],
    ):
        self._primary = primary
        self._fallback_presets = list(fallback_presets)
        self._provider_factory = provider_factory
        self._primary_failures = 0
        self._primary_tripped_at: float | None = None

    @property
    def generation(self):
        return self._primary.generation

    @generation.setter
    def generation(self, value):
        self._primary.generation = value

    @property
    def supports_progress_deltas(self) -> bool:
        return bool(getattr(self._primary, "supports_progress_deltas", False))

    def get_default_model(self) -> str:
        return self._primary.get_default_model()

    def _primary_available(self) -> bool:
        if self._primary_tripped_at is None:
            return True
        return time.monotonic() - self._primary_tripped_at >= _PRIMARY_COOLDOWN_S

    async def chat(self, **kwargs: Any) -> LLMResponse:
        if not self._fallback_presets:
            return await self._primary.chat(**kwargs)
        return await self._try_with_fallback(lambda p, kw: p.chat(**kw), kwargs, None)

    async def chat_stream(self, **kwargs: Any) -> LLMResponse:
        if not self._fallback_presets:
            return await self._primary.chat_stream(**kwargs)
        has_streamed = [False]
        original_delta = kwargs.get("on_content_delta")

        async def _tracking_delta(text: str) -> None:
            if text:
                has_streamed[0] = True
            if original_delta:
                await original_delta(text)

        kwargs["on_content_delta"] = _tracking_delta
        return await self._try_with_fallback(lambda p, kw: p.chat_stream(**kw), kwargs, has_streamed)

    async def _try_with_fallback(
        self,
        call: Callable[[LLMProvider, dict[str, Any]], Awaitable[LLMResponse]],
        kwargs: dict[str, Any],
        has_streamed: list[bool] | None,
    ) -> LLMResponse:
        primary_model = kwargs.get("model") or self._primary.get_default_model()
        last_response: LLMResponse | None = None

        if self._primary_available():
            response = await call(self._primary, kwargs)
            if response.finish_reason != "error":
                self._primary_failures = 0
                self._primary_tripped_at = None
                return response
            if has_streamed is not None and has_streamed[0]:
                return response
            if not self._should_fallback(response):
                return response
            last_response = response
            self._primary_failures += 1
            if self._primary_failures >= _PRIMARY_FAILURE_THRESHOLD:
                self._primary_tripped_at = time.monotonic()
                logger.warning("Primary model '{}' circuit opened", primary_model)

        for fallback in self._fallback_presets:
            if has_streamed is not None and has_streamed[0]:
                break
            try:
                fallback_provider = self._provider_factory(fallback)
            except Exception as exc:
                logger.warning("Failed to create fallback provider '{}': {}", fallback.model, exc)
                continue

            original_values = {
                name: kwargs.get(name, _MISSING)
                for name in ("model", "max_tokens", "temperature", "reasoning_effort")
            }
            kwargs["model"] = fallback.model
            kwargs["max_tokens"] = fallback.max_tokens
            kwargs["temperature"] = fallback.temperature
            if fallback.reasoning_effort is None:
                kwargs.pop("reasoning_effort", None)
            else:
                kwargs["reasoning_effort"] = fallback.reasoning_effort
            try:
                response = await call(fallback_provider, kwargs)
            finally:
                for name, value in original_values.items():
                    if value is _MISSING:
                        kwargs.pop(name, None)
                    else:
                        kwargs[name] = value
            if response.finish_reason != "error":
                logger.info("Fallback '{}' succeeded after primary '{}'", fallback.model, primary_model)
                return response
            last_response = response

        return last_response or LLMResponse(
            content=f"Primary model '{primary_model}' circuit open and no fallbacks available",
            finish_reason="error",
        )

    @staticmethod
    def _should_fallback(response: LLMResponse) -> bool:
        if response.error_should_retry is False:
            return False
        status = response.error_status_code
        kind = (response.error_kind or "").lower()
        error_type = (response.error_type or "").lower()
        code = (response.error_code or "").lower()
        text = (response.content or "").lower()
        if status in {400, 401, 403, 404, 422}:
            return False
        if kind in _NON_FALLBACK_ERROR_KINDS:
            return False
        if any(token in value for value in (kind, error_type, code) for token in _NON_FALLBACK_ERROR_KINDS):
            return False
        if response.error_should_retry is True:
            return True
        if status is not None and (status in {408, 409, 429} or 500 <= status <= 599):
            return True
        if kind in _FALLBACK_ERROR_KINDS:
            return True
        return any(token in value for value in (kind, error_type, code, text) for token in _FALLBACK_ERROR_TOKENS)
