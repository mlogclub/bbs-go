"""Shared artifact metadata namespace helpers."""

from __future__ import annotations

import json
from typing import Any

ORIGINAGENT_METADATA_KEY = "OriginAgent"
# Legacy compatibility for skills adapted from OpenClaw. New artifacts should
# write only ORIGINAGENT_METADATA_KEY, but readers keep this fallback so older
# workspace and registry skills remain usable.
LEGACY_METADATA_KEYS = ("openclaw",)


def read_originagent_metadata(raw: Any) -> dict[str, Any]:
    """Read current OriginAgent metadata, falling back to legacy namespaces."""

    data = _metadata_mapping(raw)
    for key in (ORIGINAGENT_METADATA_KEY, *LEGACY_METADATA_KEYS):
        value = data.get(key)
        if isinstance(value, dict):
            return dict(value)
    return {}


def set_originagent_metadata(raw: Any, originagent: dict[str, Any]) -> dict[str, Any]:
    """Return a metadata mapping with the current OriginAgent namespace set."""

    metadata = dict(raw) if isinstance(raw, dict) else {}
    metadata[ORIGINAGENT_METADATA_KEY] = dict(originagent)
    return metadata


def originagent_metadata(originagent: dict[str, Any]) -> dict[str, Any]:
    """Build a new metadata mapping under the current OriginAgent namespace."""

    return {ORIGINAGENT_METADATA_KEY: dict(originagent)}


def _metadata_mapping(raw: Any) -> dict[str, Any]:
    if isinstance(raw, dict):
        return raw
    if isinstance(raw, str):
        try:
            parsed = json.loads(raw)
        except (json.JSONDecodeError, TypeError):
            return {}
        return parsed if isinstance(parsed, dict) else {}
    return {}
