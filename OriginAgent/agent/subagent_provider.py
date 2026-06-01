"""Delegated provider selection for subagents."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Callable

from OriginAgent.providers.base import LLMProvider
from OriginAgent.providers.factory import ProviderSnapshot


@dataclass(frozen=True)
class SubagentProviderSelection:
    provider: LLMProvider
    model: str
    provider_summary: str
    selection_mode: str = "inherit_parent"
    provider_alias: str | None = None


class SubagentProviderSelector:
    """Small phase-2 provider selector with safe inherit-parent fallback."""

    def __init__(
        self,
        *,
        provider: LLMProvider,
        model: str,
        preset_snapshot_loader: Callable[[str], ProviderSnapshot] | None = None,
        delegated_preset: str | None = None,
    ) -> None:
        self._provider = provider
        self._model = model
        self._preset_snapshot_loader = preset_snapshot_loader
        self._delegated_preset = delegated_preset

    def set_runtime(self, provider: LLMProvider, model: str) -> None:
        self._provider = provider
        self._model = model

    def select(self) -> SubagentProviderSelection:
        preset = (self._delegated_preset or "").strip()
        if preset and self._preset_snapshot_loader is not None:
            snapshot = self._preset_snapshot_loader(preset)
            return SubagentProviderSelection(
                provider=snapshot.provider,
                model=snapshot.model,
                provider_summary=f"preset:{preset}:{snapshot.model}",
                selection_mode="delegated_preset",
                provider_alias=preset,
            )
        return SubagentProviderSelection(
            provider=self._provider,
            model=self._model,
            provider_summary=f"inherit:{self._model}",
        )
