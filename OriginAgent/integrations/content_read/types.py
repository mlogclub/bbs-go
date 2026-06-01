"""Shared types for content_read providers."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


@dataclass
class ContentReadResult:
    """Structured content returned by a platform provider."""

    source_type: str
    title: str
    url: str
    content: str
    metadata: dict[str, Any] = field(default_factory=dict)

    def to_payload(self, max_chars: int) -> dict[str, Any]:
        text = self.content or ""
        truncated = len(text) > max_chars
        if truncated:
            text = text[:max_chars]
        return {
            "source_type": self.source_type,
            "title": self.title,
            "url": self.url,
            "content": text,
            "metadata": self.metadata,
            "truncated": truncated,
        }
