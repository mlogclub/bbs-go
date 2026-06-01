"""Normalized, sanitized presence evidence ingestion."""

from __future__ import annotations

import uuid
from dataclasses import dataclass, field, replace
from datetime import datetime, timedelta, timezone
from typing import Any

from .presence import (
    TRUSTED_PERSON_SOURCES,
    UNKNOWN_OCCUPANCY_MIN_CONFIDENCE,
    VALID_ROLES,
    VALID_SOURCES,
    OccupancyState,
    PersonPresence,
    PresenceStore,
)

VALID_STATUS_HINTS = {"home", "away", "unknown", "activity", "none"}
NON_PERSON_BINDING_SOURCES = {"motion", "door_sensor", "door_lock", "system", "schedule"}
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


@dataclass
class PresenceSignal:
    signal_id: str
    source: str
    observed_at: str
    status_hint: str
    confidence: float | None = None
    person_id: str | None = None
    role: str = "unknown"
    zone: str | None = None
    expires_at: str | None = None
    metadata: dict[str, str] = field(default_factory=dict)

    def __post_init__(self) -> None:
        self.signal_id = _required_str(self.signal_id, "signal_id")
        self.source = _normalize_source(self.source)
        self.observed_at = _required_str(self.observed_at, "observed_at")
        self.status_hint = _normalize_status_hint(self.status_hint)
        self.role = _normalize_role(self.role)
        if self.status_hint == "none" and self.source != "manual":
            raise ValueError('status_hint="none" is only valid for manual source')
        if self.person_id is not None:
            self.person_id = self.person_id.strip() or None
        if self.zone is not None:
            self.zone = self.zone.strip() or None
        if self.confidence is not None:
            self.confidence = _clamp_confidence(self.confidence)
        self.metadata = sanitize_presence_metadata(self.metadata)


class PresenceSignalIngestor:
    def __init__(self, presence_store: PresenceStore):
        self.presence_store = presence_store

    def ingest(self, signal: PresenceSignal) -> PersonPresence | OccupancyState | None:
        signal.metadata = sanitize_presence_metadata(signal.metadata)
        signal = sanitize_signal(signal)
        if _is_expired(signal.expires_at):
            return None

        expires_at = signal.expires_at or default_expires_at(signal.source, signal.observed_at)
        confidence = default_confidence(signal.source, signal.status_hint, signal.confidence)

        if signal.status_hint == "none":
            return None

        if signal.source in NON_PERSON_BINDING_SOURCES and signal.person_id:
            raise ValueError(f"{signal.source!r} signal cannot bind person presence")

        if signal.source == "schedule":
            return None

        if signal.person_id:
            if signal.source not in TRUSTED_PERSON_SOURCES:
                raise ValueError(f"{signal.source!r} signal cannot bind person presence")
            if signal.status_hint == "activity":
                raise ValueError("activity signals cannot bind person presence")
            return self.presence_store.upsert_presence(
                signal.person_id,
                role=signal.role,
                status=signal.status_hint,
                source=signal.source,
                confidence=confidence,
                expires_at=expires_at,
                now=_parse_datetime(signal.observed_at),
            )

        if signal.status_hint in {"activity", "unknown", "home"}:
            if confidence < UNKNOWN_OCCUPANCY_MIN_CONFIDENCE:
                return None
            return self.presence_store.mark_unknown_occupancy(
                source=signal.source,
                confidence=confidence,
                expires_at=expires_at,
                now=_parse_datetime(signal.observed_at),
            )

        raise ValueError("away signals require a trusted person_id")


def new_presence_signal(
    *,
    source: str,
    status_hint: str,
    person_id: str | None = None,
    role: str = "unknown",
    zone: str | None = None,
    observed_at: str | None = None,
    confidence: float | None = None,
    expires_at: str | None = None,
    metadata: dict[str, Any] | None = None,
) -> PresenceSignal:
    observed = observed_at or _format_datetime(datetime.now(timezone.utc))
    return PresenceSignal(
        signal_id=f"presence_signal_{uuid.uuid4().hex[:12]}",
        source=source,
        observed_at=observed,
        status_hint=status_hint,
        confidence=default_confidence(source, status_hint, confidence),
        person_id=person_id,
        role=role,
        zone=zone,
        expires_at=expires_at or default_expires_at(source, observed),
        metadata=sanitize_presence_metadata(metadata or {}),
    )


def sanitize_signal(signal: PresenceSignal) -> PresenceSignal:
    return replace(signal, metadata=sanitize_presence_metadata(signal.metadata))


def sanitize_presence_metadata(metadata: dict[str, Any]) -> dict[str, str]:
    sanitized: dict[str, str] = {}
    for key, value in metadata.items():
        if not isinstance(key, str):
            continue
        normalized_key = key.strip()
        if not normalized_key:
            continue
        if normalized_key.casefold() in FORBIDDEN_METADATA_KEYS:
            continue
        if value is None:
            continue
        sanitized[normalized_key] = str(value)
    return sanitized


sanitize_metadata = sanitize_presence_metadata


def default_confidence(source: str, status_hint: str, confidence: float | None) -> float:
    source = _normalize_source(source)
    status_hint = _normalize_status_hint(status_hint)
    if confidence is None:
        confidence = _source_status_default_confidence(source, status_hint)
    confidence = _clamp_confidence(confidence)
    if source == "schedule":
        return min(confidence, 0.30)
    return confidence


def default_expires_at(source: str, observed_at: str) -> str:
    observed = _parse_datetime(observed_at)
    ttl = _source_default_ttl(_normalize_source(source))
    return _format_datetime(observed + ttl)


def _source_status_default_confidence(source: str, status_hint: str) -> float:
    if source == "manual":
        return 0.95
    if source == "voice":
        return 0.85
    if source == "phone_geofence":
        return 0.85 if status_hint == "away" else 0.75
    if source == "wifi_presence":
        return 0.80 if status_hint == "away" else 0.70
    if source in {"door_sensor", "door_lock"}:
        return 0.60
    if source == "motion":
        return 0.65
    if source == "schedule":
        return 0.30
    return 0.50


def _source_default_ttl(source: str) -> timedelta:
    if source == "manual":
        return timedelta(hours=12)
    if source == "voice":
        return timedelta(hours=4)
    if source == "phone_geofence":
        return timedelta(minutes=30)
    if source == "wifi_presence":
        return timedelta(minutes=20)
    if source in {"door_sensor", "door_lock"}:
        return timedelta(minutes=10)
    if source == "motion":
        return timedelta(minutes=5)
    if source == "schedule":
        return timedelta(hours=1)
    return timedelta(minutes=5)


def _normalize_source(source: str | None) -> str:
    normalized = (source or "").strip().lower()
    if normalized not in VALID_SOURCES:
        raise ValueError(f"invalid presence signal source: {source!r}")
    return normalized


def _normalize_status_hint(status_hint: str | None) -> str:
    normalized = (status_hint or "").strip().lower()
    if normalized not in VALID_STATUS_HINTS:
        raise ValueError(f"invalid presence status_hint: {status_hint!r}")
    return normalized


def _normalize_role(role: str | None) -> str:
    normalized = (role or "unknown").strip().lower()
    if normalized not in VALID_ROLES:
        raise ValueError(f"invalid presence role: {role!r}")
    return normalized


def _required_str(value: str | None, field_name: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{field_name} must be a non-empty string")
    return value.strip()


def _clamp_confidence(value: Any) -> float:
    try:
        confidence = float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("confidence must be numeric") from exc
    return max(0.0, min(1.0, confidence))


def _is_expired(expires_at: str | None) -> bool:
    if not expires_at:
        return False
    expires = _parse_datetime(expires_at)
    now = datetime.now(timezone.utc) if expires.tzinfo is not None else datetime.now()
    return expires <= now


def _parse_datetime(value: str) -> datetime:
    try:
        parsed = datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError as exc:
        raise ValueError(f"invalid datetime: {value!r}") from exc
    if parsed.tzinfo is None:
        return parsed
    return parsed.astimezone(timezone.utc)


def _format_datetime(value: datetime) -> str:
    if value.tzinfo is not None:
        value = value.astimezone(timezone.utc)
    return value.isoformat()
