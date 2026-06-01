"""Current-state presence storage and conservative occupancy resolution."""

from __future__ import annotations

import json
import os
from contextlib import suppress
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from filelock import FileLock
from loguru import logger

from OriginAgent.utils.helpers import ensure_dir

VALID_ROLES = {"admin", "resident", "guest", "child", "elder", "unknown"}
VALID_PERSON_STATUSES = {"home", "away", "unknown"}
VALID_OCCUPANCY_STATUSES = {"occupied", "empty", "unknown"}
TRUSTED_PERSON_SOURCES = {
    "manual",
    # Conversation/current-session identity self-report. This is not voiceprint,
    # speaker recognition, or any biometric identity binding.
    "voice",
    "phone_geofence",
    "wifi_presence",
    "trusted_identity",
}
VALID_SOURCES = TRUSTED_PERSON_SOURCES | {
    "motion",
    "door_lock",
    "door_sensor",
    "schedule",
    "system",
}
HOME_CONFIDENCE_THRESHOLD = 0.7
AWAY_CONFIDENCE_THRESHOLD = 0.8
UNKNOWN_OCCUPANCY_MIN_CONFIDENCE = 0.5


@dataclass
class PersonPresence:
    person_id: str
    role: str
    status: str
    source: str
    confidence: float
    updated_at: str
    expires_at: str | None = None

    @classmethod
    def from_dict(cls, raw: dict[str, Any]) -> "PersonPresence":
        person_id = _required_str(raw, "person_id")
        role = _normalize_role(_required_str(raw, "role"))
        status = _normalize_person_status(_required_str(raw, "status"))
        source = _normalize_source(_required_str(raw, "source"))
        updated_at = _required_str(raw, "updated_at")
        expires_at = raw.get("expires_at")
        if expires_at is not None and not isinstance(expires_at, str):
            raise ValueError("expires_at must be a string or null")
        return cls(
            person_id=person_id,
            role=role,
            status=status,
            source=source,
            confidence=_normalize_confidence(raw.get("confidence", 1.0)),
            updated_at=updated_at,
            expires_at=expires_at,
        )

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass
class OccupancyState:
    status: str
    source: str = "resolver"
    confidence: float = 0.0
    updated_at: str = ""
    expires_at: str | None = None

    @classmethod
    def from_dict(cls, raw: dict[str, Any]) -> "OccupancyState":
        status = _normalize_occupancy_status(_required_str(raw, "status"))
        source = _normalize_source(_required_str(raw, "source"))
        updated_at = _required_str(raw, "updated_at")
        expires_at = raw.get("expires_at")
        if expires_at is not None and not isinstance(expires_at, str):
            raise ValueError("expires_at must be a string or null")
        return cls(
            status=status,
            source=source,
            confidence=_normalize_confidence(raw.get("confidence", 1.0)),
            updated_at=updated_at,
            expires_at=expires_at,
        )

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


class PresenceStore:
    """Current-state store for normalized presence only.

    The file intentionally contains no raw device identifiers, MAC/IP values,
    screenshots, embeddings, biometric material, or raw access-control logs.
    """

    def __init__(
        self,
        workspace: Path,
        *,
        presence_file: Path | None = None,
        lock_factory: Callable[[], FileLock] | None = None,
    ):
        self.workspace = workspace
        self.memory_dir = ensure_dir(workspace / "memory")
        self.presence_file = presence_file or self.memory_dir / "presence.json"
        self._lock_file = self.memory_dir / ".lock"
        self._lock_factory = lock_factory

    def _locked(self) -> FileLock:
        if self._lock_factory is not None:
            return self._lock_factory()
        return FileLock(str(self._lock_file))

    def read_state(self) -> dict[str, Any]:
        with self._locked():
            return self.read_state_unlocked()

    def read_state_unlocked(self) -> dict[str, Any]:
        try:
            raw = json.loads(self.presence_file.read_text(encoding="utf-8"))
            if not isinstance(raw, dict):
                raise ValueError("presence state must be an object")
            people = raw.get("people", {})
            if not isinstance(people, dict):
                raise ValueError("people must be an object")
            normalized_people: dict[str, dict[str, Any]] = {}
            for person_id, record in people.items():
                if not isinstance(person_id, str) or not isinstance(record, dict):
                    raise ValueError("invalid people record")
                presence = PersonPresence.from_dict(record)
                if presence.person_id != person_id:
                    raise ValueError("person_id key mismatch")
                normalized_people[person_id] = presence.to_dict()
            state: dict[str, Any] = {"people": normalized_people}
            unknown = raw.get("unknown_occupancy")
            if unknown is not None:
                if not isinstance(unknown, dict):
                    raise ValueError("unknown_occupancy must be an object or null")
                state["unknown_occupancy"] = OccupancyState.from_dict(unknown).to_dict()
            return state
        except FileNotFoundError:
            return self._empty_state()
        except (OSError, json.JSONDecodeError, ValueError, TypeError):
            logger.warning("Failed to read presence state from {}; treating as unknown", self.presence_file)
            return self._empty_state()

    def upsert_presence(
        self,
        person_id: str,
        *,
        role: str,
        status: str,
        source: str,
        confidence: float = 1.0,
        expires_at: str | None = None,
        now: datetime | None = None,
    ) -> PersonPresence:
        """Upsert a concrete person's presence from a trusted identity source."""

        person_id = person_id.strip()
        if not person_id:
            raise ValueError("person_id must be non-empty")
        source = _normalize_source(source)
        if source not in TRUSTED_PERSON_SOURCES:
            raise ValueError(f"source cannot bind person presence: {source!r}")
        presence = PersonPresence(
            person_id=person_id,
            role=_normalize_role(role),
            status=_normalize_person_status(status),
            source=source,
            confidence=_normalize_confidence(confidence),
            updated_at=_format_timestamp(now),
            expires_at=expires_at,
        )
        with self._locked():
            state = self.read_state_unlocked()
            people = state.setdefault("people", {})
            people[person_id] = presence.to_dict()
            self._write_state_unlocked(state)
        return presence

    def mark_unknown_occupancy(
        self,
        *,
        source: str,
        confidence: float = 1.0,
        expires_at: str | None = None,
        now: datetime | None = None,
    ) -> OccupancyState:
        """Record an active unknown occupancy signal without binding identity."""

        source = _normalize_source(source)
        signal = OccupancyState(
            status="unknown",
            source=source,
            confidence=_normalize_confidence(confidence),
            updated_at=_format_timestamp(now),
            expires_at=expires_at,
        )
        with self._locked():
            state = self.read_state_unlocked()
            state["unknown_occupancy"] = signal.to_dict()
            self._write_state_unlocked(state)
        return signal

    def resolve_person(
        self,
        person_id: str,
        *,
        now: datetime | None = None,
    ) -> PersonPresence | None:
        state = self.read_state()
        raw = state.get("people", {}).get(person_id)
        if not isinstance(raw, dict):
            return None
        try:
            presence = PersonPresence.from_dict(raw)
        except ValueError:
            return None
        if _is_expired(presence.expires_at, now) or (
            presence.confidence < _person_status_threshold(presence.status)
        ):
            return PersonPresence(
                person_id=presence.person_id,
                role=presence.role,
                status="unknown",
                source=presence.source,
                confidence=0.0,
                updated_at=presence.updated_at,
                expires_at=presence.expires_at,
            )
        return presence

    def resolve_occupancy(self, *, now: datetime | None = None) -> OccupancyState:
        state = self.read_state()
        people = []
        for raw in state.get("people", {}).values():
            if not isinstance(raw, dict):
                continue
            with suppress(ValueError):
                people.append(PersonPresence.from_dict(raw))

        for person in people:
            if (
                person.status == "home"
                and person.confidence >= HOME_CONFIDENCE_THRESHOLD
                and not _is_expired(person.expires_at, now)
            ):
                return OccupancyState(
                    status="occupied",
                    source="person_presence",
                    confidence=person.confidence,
                    updated_at=person.updated_at,
                )

        unknown_signal = self._active_unknown_signal(state, now=now)
        if unknown_signal is not None:
            return OccupancyState(
                status="unknown",
                source=unknown_signal.source,
                confidence=unknown_signal.confidence,
                updated_at=unknown_signal.updated_at,
                expires_at=unknown_signal.expires_at,
            )

        known_residents = [
            person
            for person in people
            if person.role in {"admin", "resident"}
            and not _is_expired(person.expires_at, now)
        ]
        if not known_residents:
            return self._unknown()
        if all(
            person.status == "away" and person.confidence >= AWAY_CONFIDENCE_THRESHOLD
            for person in known_residents
        ):
            return OccupancyState(
                status="empty",
                source="person_presence",
                confidence=min(person.confidence for person in known_residents),
                updated_at=max(person.updated_at for person in known_residents),
            )
        return self._unknown()

    def _active_unknown_signal(
        self,
        state: dict[str, Any],
        *,
        now: datetime | None,
    ) -> OccupancyState | None:
        raw = state.get("unknown_occupancy")
        if not isinstance(raw, dict):
            return None
        try:
            signal = OccupancyState.from_dict(raw)
        except ValueError:
            return None
        if _is_expired(signal.expires_at, now):
            return None
        if signal.confidence < UNKNOWN_OCCUPANCY_MIN_CONFIDENCE:
            return None
        return signal

    def _write_state_unlocked(self, state: dict[str, Any]) -> None:
        text = json.dumps(state, ensure_ascii=False, indent=2, sort_keys=True) + "\n"
        _write_text_atomic(self.presence_file, text)

    @staticmethod
    def _empty_state() -> dict[str, Any]:
        return {"people": {}}

    @staticmethod
    def _unknown() -> OccupancyState:
        return OccupancyState(status="unknown", updated_at="")


def _required_str(raw: dict[str, Any], key: str) -> str:
    value = raw.get(key)
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{key} must be a non-empty string")
    return value.strip()


def _normalize_role(role: str | None) -> str:
    normalized = (role or "").strip().lower()
    if normalized not in VALID_ROLES:
        raise ValueError(f"invalid presence role: {role!r}")
    return normalized


def _normalize_person_status(status: str | None) -> str:
    normalized = (status or "").strip().lower()
    if normalized not in VALID_PERSON_STATUSES:
        raise ValueError(f"invalid presence status: {status!r}")
    return normalized


def _normalize_occupancy_status(status: str | None) -> str:
    normalized = (status or "").strip().lower()
    if normalized not in VALID_OCCUPANCY_STATUSES:
        raise ValueError(f"invalid occupancy status: {status!r}")
    return normalized


def _normalize_source(source: str | None) -> str:
    normalized = (source or "").strip().lower()
    if normalized not in VALID_SOURCES:
        raise ValueError(f"invalid presence source: {source!r}")
    return normalized


def _person_status_threshold(status: str) -> float:
    if status == "home":
        return HOME_CONFIDENCE_THRESHOLD
    if status == "away":
        return AWAY_CONFIDENCE_THRESHOLD
    return 1.0


def _normalize_confidence(value: Any) -> float:
    try:
        confidence = float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("confidence must be numeric") from exc
    return max(0.0, min(1.0, confidence))


def _format_timestamp(now: datetime | None = None) -> str:
    ts = now or datetime.now(timezone.utc)
    if ts.tzinfo is None:
        ts = ts.replace(tzinfo=timezone.utc)
    return ts.astimezone(timezone.utc).isoformat()


def _is_expired(expires_at: str | None, now: datetime | None = None) -> bool:
    if not expires_at:
        return False
    expires_ts = _timestamp(expires_at)
    if expires_ts is None:
        return True
    if now is None:
        current = datetime.now(timezone.utc).timestamp()
    elif now.tzinfo is None:
        current = now.replace(tzinfo=timezone.utc).timestamp()
    else:
        current = now.astimezone(timezone.utc).timestamp()
    return expires_ts <= current


def _timestamp(value: str) -> float | None:
    try:
        normalized = value.replace("Z", "+00:00")
        parsed = datetime.fromisoformat(normalized)
        if parsed.tzinfo is None:
            return parsed.replace(tzinfo=timezone.utc).timestamp()
        return parsed.astimezone(timezone.utc).timestamp()
    except ValueError:
        return None


def _write_text_atomic(path: Path, text: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = path.with_name(f".{path.name}.tmp")
    try:
        with open(tmp_path, "w", encoding="utf-8") as f:
            f.write(text)
            f.flush()
            os.fsync(f.fileno())
        os.replace(tmp_path, path)
        _fsync_parent(path)
    except BaseException:
        tmp_path.unlink(missing_ok=True)
        raise


def _fsync_parent(path: Path) -> None:
    with suppress(PermissionError, OSError):
        fd = os.open(str(path.parent), os.O_RDONLY)
        try:
            os.fsync(fd)
        finally:
            os.close(fd)
