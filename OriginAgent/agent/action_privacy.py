"""Generic action metadata redaction helpers."""

from __future__ import annotations

FORBIDDEN_METADATA_KEYS = {
    "mac",
    "ip",
    "ssid",
    "bssid",
    "image",
    "photo",
    "face_embedding",
    "voice_embedding",
    "raw_audio",
    "raw_video",
    "access_log",
    "door_log",
    "device_fingerprint",
}


def sanitize_metadata(metadata: dict | None) -> dict[str, str]:
    sanitized: dict[str, str] = {}
    if not isinstance(metadata, dict):
        return sanitized
    forbidden = {key.casefold() for key in FORBIDDEN_METADATA_KEYS}
    for key, value in metadata.items():
        normalized_key = str(key or "").strip()
        if not normalized_key or normalized_key.casefold() in forbidden:
            continue
        if value is None:
            continue
        text = str(value).strip()
        if text:
            sanitized[normalized_key] = text
    return sanitized
