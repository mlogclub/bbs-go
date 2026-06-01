"""Device registry and audit-safe device identity helpers."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from typing import Iterable

from OriginAgent.security.policy import PolicyDeniedError


@dataclass(frozen=True)
class DeviceRecord:
    device_id: str
    domain: str
    room: str | None
    display_name: str | None = None
    audit_label: str | None = None
    device_ref: str | None = None


class DeviceRegistry:
    def __init__(self, records: Iterable[DeviceRecord] = ()):
        self._records = tuple(records)

    def resolve(
        self,
        *,
        actor_id: str,
        domain: str,
        room: str | None,
        device_ref: str,
    ) -> DeviceRecord:
        normalized_ref = _norm(device_ref)
        normalized_domain = _norm(domain)
        normalized_room = _norm_optional(room)
        matches = [
            record for record in self._records
            if _norm(record.domain) == normalized_domain
            and (normalized_room is None or _norm_optional(record.room) == normalized_room)
            and normalized_ref in {
                _norm(record.device_ref or ""),
                _norm(record.display_name or ""),
                _norm(record.audit_label or ""),
            }
        ]
        if len(matches) == 1:
            return matches[0]
        if not matches:
            raise PolicyDeniedError(
                "Unknown device reference",
                code="unknown_device_ref",
                boundary="device",
                policy_rule="device_registry_required",
            )
        raise PolicyDeniedError(
            "Ambiguous device reference",
            code="ambiguous_device_ref",
            boundary="device",
            policy_rule="device_registry_disambiguation_required",
        )


def sanitize_device_scope(scope: str | None) -> str | None:
    if not scope:
        return None
    parts = [part for part in str(scope).split(".") if part]
    if len(parts) >= 4 and parts[0] == "home":
        return ".".join([*parts[:-1], "<device>"])
    if len(parts) >= 3 and parts[0] == "home":
        return ".".join([*parts[:-1], "<device>"])
    return scope


class SmartHomeScopeRedactor:
    def __call__(self, scope: str | None) -> str | None:
        return sanitize_device_scope(scope)


DEVICE_SCOPE_REDACTOR = SmartHomeScopeRedactor()


def _norm(value: str | None) -> str:
    return str(value or "").strip().casefold()


def _norm_optional(value: str | None) -> str | None:
    normalized = _norm(value)
    return normalized or None
