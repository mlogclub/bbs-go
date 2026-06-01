"""Capability grants for delegated runtime tasks."""

from __future__ import annotations

import hashlib
import json
import os
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Literal

from OriginAgent.cron.types import CronPayload
from OriginAgent.security.capabilities import CapabilitySnapshot, CapabilitySource, CapabilityTrigger
from OriginAgent.security.policy import PolicyDeniedError

GrantSource = Literal["admin_config", "user_confirmation", "test"]
_GRANT_ERROR_MESSAGE = "Capability grant is missing, expired, or revoked."


@dataclass(frozen=True)
class CapabilityGrant:
    grant_id: str
    created_by: str
    created_at: str
    expires_at: str | None = None
    revoked_at: str | None = None
    source: GrantSource = "test"
    can_exec: bool = False
    can_read_files: bool = False
    can_write_files: bool = False
    can_send_cross_target: bool = False
    can_create_cron: bool = False
    can_spawn: bool = False
    allowed_device_domains: tuple[str, ...] = ()
    allowed_mcp_scopes: tuple[str, ...] = ()

    def is_expired(self, now: datetime | None = None) -> bool:
        if not self.expires_at:
            return False
        expires = _parse_datetime(self.expires_at)
        if expires is None:
            return True
        now = _normalize_datetime(now or datetime.now(timezone.utc))
        return expires <= now

    def is_revoked(self) -> bool:
        return bool(self.revoked_at)

    def is_active(self, now: datetime | None = None) -> bool:
        return not self.is_revoked() and not self.is_expired(now)

    def to_snapshot(self, *, trigger: CapabilityTrigger = "scheduled") -> CapabilitySnapshot:
        # C5a grants are cron-only; subagent grants need a separate C5b conversion path.
        return self._to_snapshot_for(source="cron", trigger=trigger)

    def to_subagent_snapshot(self) -> CapabilitySnapshot:
        return self._to_snapshot_for(source="subagent", trigger="subagent")

    def _to_snapshot_for(
        self,
        *,
        source: CapabilitySource,
        trigger: CapabilityTrigger,
    ) -> CapabilitySnapshot:
        return CapabilitySnapshot(
            version=1,
            source=source,
            trigger=trigger,
            can_exec=self.can_exec,
            can_read_files=self.can_read_files,
            can_write_files=self.can_write_files,
            can_send_cross_target=self.can_send_cross_target,
            can_create_cron=self.can_create_cron,
            can_spawn=self.can_spawn,
            allowed_device_domains=tuple(self.allowed_device_domains),
            allowed_mcp_scopes=tuple(self.allowed_mcp_scopes),
        )

    def summary(self) -> dict[str, Any]:
        flags = [
            self.can_exec,
            self.can_read_files,
            self.can_write_files,
            self.can_send_cross_target,
            self.can_create_cron,
            self.can_spawn,
        ]
        return {
            "grant_ref": _grant_ref(self.grant_id),
            "source": self.source,
            "active": self.is_active(),
            "expired": self.is_expired(),
            "revoked": self.is_revoked(),
            "expires_at": self.expires_at,
            "enabled_flags_count": sum(1 for flag in flags if flag),
            "allowed_device_domains": list(self.allowed_device_domains),
            "allowed_mcp_scopes": list(self.allowed_mcp_scopes),
        }

    def to_dict(self) -> dict[str, Any]:
        data = asdict(self)
        data["allowed_device_domains"] = list(self.allowed_device_domains)
        data["allowed_mcp_scopes"] = list(self.allowed_mcp_scopes)
        return data

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> "CapabilityGrant":
        return cls(
            grant_id=str(data.get("grant_id") or data.get("grantId") or ""),
            created_by=str(data.get("created_by") or data.get("createdBy") or ""),
            created_at=str(data.get("created_at") or data.get("createdAt") or ""),
            expires_at=data.get("expires_at") or data.get("expiresAt"),
            revoked_at=data.get("revoked_at") or data.get("revokedAt"),
            source=data.get("source") or "test",
            can_exec=bool(data.get("can_exec") or data.get("canExec")),
            can_read_files=bool(data.get("can_read_files") or data.get("canReadFiles")),
            can_write_files=bool(data.get("can_write_files") or data.get("canWriteFiles")),
            can_send_cross_target=bool(
                data.get("can_send_cross_target") or data.get("canSendCrossTarget")
            ),
            can_create_cron=bool(data.get("can_create_cron") or data.get("canCreateCron")),
            can_spawn=bool(data.get("can_spawn") or data.get("canSpawn")),
            allowed_device_domains=tuple(
                data.get("allowed_device_domains")
                or data.get("allowedDeviceDomains")
                or ()
            ),
            allowed_mcp_scopes=tuple(
                data.get("allowed_mcp_scopes")
                or data.get("allowedMcpScopes")
                or ()
            ),
        )


class CapabilityGrantStore:
    def __init__(self, workspace: Path):
        self.workspace = Path(workspace)
        self.path = self.workspace / "memory" / "security" / "capability_grants.json"

    def put(self, grant: CapabilityGrant) -> None:
        grants = {item.grant_id: item for item in self.list_all()}
        grants[grant.grant_id] = grant
        self._save(list(grants.values()))

    def get(self, grant_id: str) -> CapabilityGrant | None:
        for grant in self.list_all():
            if grant.grant_id == grant_id:
                return grant
        return None

    def revoke(self, grant_id: str, *, revoked_at: datetime | None = None) -> bool:
        grants = self.list_all()
        found = False
        revoked = _format_datetime(revoked_at or datetime.now(timezone.utc))
        updated: list[CapabilityGrant] = []
        for grant in grants:
            if grant.grant_id == grant_id:
                found = True
                updated.append(
                    CapabilityGrant(
                        grant_id=grant.grant_id,
                        created_by=grant.created_by,
                        created_at=grant.created_at,
                        expires_at=grant.expires_at,
                        revoked_at=revoked,
                        source=grant.source,
                        can_exec=grant.can_exec,
                        can_read_files=grant.can_read_files,
                        can_write_files=grant.can_write_files,
                        can_send_cross_target=grant.can_send_cross_target,
                        can_create_cron=grant.can_create_cron,
                        can_spawn=grant.can_spawn,
                        allowed_device_domains=grant.allowed_device_domains,
                        allowed_mcp_scopes=grant.allowed_mcp_scopes,
                    )
                )
            else:
                updated.append(grant)
        if found:
            self._save(updated)
        return found

    def list_all(self) -> list[CapabilityGrant]:
        try:
            data = json.loads(self.path.read_text(encoding="utf-8"))
        except FileNotFoundError:
            return []
        except (json.JSONDecodeError, OSError, TypeError):
            return []
        grants = data.get("grants", []) if isinstance(data, dict) else []
        result: list[CapabilityGrant] = []
        for item in grants:
            if not isinstance(item, dict):
                continue
            grant = CapabilityGrant.from_dict(item)
            if grant.grant_id:
                result.append(grant)
        return result

    def list_active(self) -> list[CapabilityGrant]:
        return [grant for grant in self.list_all() if grant.is_active()]

    def _save(self, grants: list[CapabilityGrant]) -> None:
        payload = {
            "version": 1,
            "grants": [grant.to_dict() for grant in grants],
        }
        self.path.parent.mkdir(parents=True, exist_ok=True)
        tmp = self.path.with_suffix(self.path.suffix + ".tmp")
        try:
            with open(tmp, "w", encoding="utf-8") as handle:
                json.dump(payload, handle, indent=2, ensure_ascii=False)
                handle.write("\n")
                handle.flush()
                os.fsync(handle.fileno())
            os.replace(tmp, self.path)
        except BaseException:
            tmp.unlink(missing_ok=True)
            raise


def snapshot_for_cron_payload(
    payload: CronPayload,
    grant_store: CapabilityGrantStore,
) -> CapabilitySnapshot:
    if not payload.grant_id:
        return CapabilitySnapshot.scheduled_default()
    grant = grant_store.get(payload.grant_id)
    if grant is None:
        _raise_grant_denied("capability_grant_missing")
    if grant.is_revoked():
        _raise_grant_denied("capability_grant_revoked")
    if grant.is_expired():
        _raise_grant_denied("capability_grant_expired")
    return grant.to_snapshot(trigger="scheduled")


def _raise_grant_denied(policy_rule: str) -> None:
    raise PolicyDeniedError(
        _GRANT_ERROR_MESSAGE,
        code=policy_rule,
        boundary="cron",
        policy_rule=policy_rule,
    )


def _grant_ref(grant_id: str) -> str:
    return hashlib.sha256(grant_id.encode("utf-8")).hexdigest()[:12]


def _parse_datetime(value: str) -> datetime | None:
    try:
        parsed = datetime.fromisoformat(value.replace("Z", "+00:00"))
    except (TypeError, ValueError):
        return None
    return _normalize_datetime(parsed)


def _normalize_datetime(value: datetime) -> datetime:
    if value.tzinfo is None:
        return value.replace(tzinfo=timezone.utc)
    return value.astimezone(timezone.utc)


def _format_datetime(value: datetime) -> str:
    return _normalize_datetime(value).isoformat()
