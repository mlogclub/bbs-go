"""Function-style presence evidence adapters.

These adapters do not talk to real devices or perform identification. They
only normalize already-trusted low-privacy inputs into PresenceSignal objects.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any

from .presence_signals import PresenceSignal, new_presence_signal


@dataclass
class WifiDeviceEvent:
    registered_device_id: str
    person_id: str
    online: bool
    role: str = "resident"
    metadata: dict[str, Any] = field(default_factory=dict)


@dataclass
class PhoneGeofenceEvent:
    person_id: str
    state: str
    role: str = "resident"
    metadata: dict[str, Any] = field(default_factory=dict)


@dataclass
class DoorEvent:
    scope: str
    event: str
    metadata: dict[str, Any] = field(default_factory=dict)


@dataclass
class MotionEvent:
    zone: str
    active: bool
    metadata: dict[str, Any] = field(default_factory=dict)


@dataclass
class ScheduleHint:
    expected: str
    person_id: str | None = None
    zone: str | None = None
    metadata: dict[str, Any] = field(default_factory=dict)


class ManualPresenceAdapter:
    @staticmethod
    def home(person_id: str, role: str = "resident") -> PresenceSignal:
        return new_presence_signal(
            source="manual",
            status_hint="home",
            person_id=person_id,
            role=role,
        )

    @staticmethod
    def away(person_id: str, role: str = "resident") -> PresenceSignal:
        return new_presence_signal(
            source="manual",
            status_hint="away",
            person_id=person_id,
            role=role,
        )

    @staticmethod
    def nobody_home() -> PresenceSignal:
        return new_presence_signal(source="manual", status_hint="none")

    @staticmethod
    def someone_home(source: str = "manual") -> PresenceSignal:
        return new_presence_signal(source=source, status_hint="activity")


class VoiceSelfReportAdapter:
    """Self-report from the current authenticated session, not voiceprint."""

    @staticmethod
    def home(person_id: str, role: str = "resident") -> PresenceSignal:
        return new_presence_signal(
            source="voice",
            status_hint="home",
            person_id=person_id,
            role=role,
        )

    @staticmethod
    def away(person_id: str, role: str = "resident") -> PresenceSignal:
        return new_presence_signal(
            source="voice",
            status_hint="away",
            person_id=person_id,
            role=role,
        )


class WifiPresenceAdapter:
    @staticmethod
    def from_event(event: WifiDeviceEvent) -> PresenceSignal:
        return new_presence_signal(
            source="wifi_presence",
            status_hint="home" if event.online else "away",
            person_id=event.person_id,
            role=event.role,
            metadata={
                **event.metadata,
                "registered_device_id": event.registered_device_id,
            },
        )


class PhoneGeofenceAdapter:
    @staticmethod
    def from_event(event: PhoneGeofenceEvent) -> PresenceSignal:
        state = event.state.strip().lower()
        if state == "inside_home":
            status_hint = "home"
        elif state == "outside_home":
            status_hint = "away"
        else:
            raise ValueError(f"unsupported geofence state: {event.state!r}")
        return new_presence_signal(
            source="phone_geofence",
            status_hint=status_hint,
            person_id=event.person_id,
            role=event.role,
            metadata=event.metadata,
        )


class DoorSensorAdapter:
    @staticmethod
    def from_event(event: DoorEvent) -> PresenceSignal:
        return new_presence_signal(
            source="door_lock" if event.event.strip().lower() in {"locked", "unlocked"} else "door_sensor",
            status_hint="activity",
            zone=event.scope,
            metadata=event.metadata,
        )


class MotionAdapter:
    @staticmethod
    def from_event(event: MotionEvent) -> PresenceSignal:
        return new_presence_signal(
            source="motion",
            status_hint="activity" if event.active else "unknown",
            zone=event.zone,
            confidence=None if event.active else 0.0,
            metadata=event.metadata,
        )


class ScheduleHintAdapter:
    @staticmethod
    def from_hint(hint: ScheduleHint) -> PresenceSignal:
        expected = hint.expected.strip().lower()
        if expected not in {"home", "away", "unknown", "activity"}:
            raise ValueError(f"unsupported schedule expectation: {hint.expected!r}")
        return new_presence_signal(
            source="schedule",
            status_hint=expected,
            person_id=None,
            role="unknown",
            zone=hint.zone,
            metadata=hint.metadata,
        )
