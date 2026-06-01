"""Immutable capability snapshots for delegated and scheduled work."""

from __future__ import annotations

from dataclasses import asdict, dataclass
from typing import Literal

CapabilitySource = Literal["user_turn", "cron", "subagent", "system"]
CapabilityTrigger = Literal["user_initiated", "scheduled", "system", "subagent"]


@dataclass(frozen=True)
class CapabilitySnapshot:
    version: int
    source: CapabilitySource
    trigger: CapabilityTrigger
    can_exec: bool
    can_read_files: bool
    can_write_files: bool
    can_send_cross_target: bool
    can_create_cron: bool
    can_spawn: bool
    allowed_device_domains: tuple[str, ...]
    allowed_mcp_scopes: tuple[str, ...]

    @classmethod
    def user_turn(cls) -> "CapabilitySnapshot":
        return cls(
            version=1,
            source="user_turn",
            trigger="user_initiated",
            can_exec=True,
            can_read_files=True,
            can_write_files=True,
            can_send_cross_target=False,
            can_create_cron=True,
            can_spawn=True,
            allowed_device_domains=("lighting",),
            allowed_mcp_scopes=("read",),
        )

    @classmethod
    def scheduled_default(cls) -> "CapabilitySnapshot":
        return cls(
            version=1,
            source="cron",
            trigger="scheduled",
            can_exec=False,
            can_read_files=False,
            can_write_files=False,
            can_send_cross_target=False,
            can_create_cron=False,
            can_spawn=False,
            allowed_device_domains=(),
            allowed_mcp_scopes=("read",),
        )

    @classmethod
    def system_default(cls) -> "CapabilitySnapshot":
        return cls(
            version=1,
            source="system",
            trigger="system",
            can_exec=False,
            can_read_files=False,
            can_write_files=False,
            can_send_cross_target=False,
            can_create_cron=False,
            can_spawn=False,
            allowed_device_domains=(),
            allowed_mcp_scopes=("read",),
        )

    def derive_subagent(self) -> "CapabilitySnapshot":
        """Return a least-privilege subagent snapshot derived from this one."""

        return CapabilitySnapshot(
            version=self.version,
            source="subagent",
            trigger="subagent" if self.trigger != "scheduled" else "scheduled",
            can_exec=False,
            can_read_files=self.can_read_files,
            can_write_files=False,
            can_send_cross_target=False,
            can_create_cron=False,
            can_spawn=False,
            allowed_device_domains=(),
            allowed_mcp_scopes=tuple(scope for scope in self.allowed_mcp_scopes if scope == "read"),
        )

    def to_dict(self) -> dict[str, object]:
        data = asdict(self)
        data["allowed_device_domains"] = list(self.allowed_device_domains)
        data["allowed_mcp_scopes"] = list(self.allowed_mcp_scopes)
        return data

    @classmethod
    def from_dict(cls, data: dict | None) -> "CapabilitySnapshot":
        if not data:
            return cls.scheduled_default()
        version = int(data.get("version") or 1)
        if version != 1:
            return cls.scheduled_default()
        return cls(
            version=version,
            source=data.get("source") or "cron",
            trigger=data.get("trigger") or "scheduled",
            can_exec=bool(data.get("can_exec", False)),
            can_read_files=bool(data.get("can_read_files", False)),
            can_write_files=bool(data.get("can_write_files", False)),
            can_send_cross_target=bool(data.get("can_send_cross_target", False)),
            can_create_cron=bool(data.get("can_create_cron", False)),
            can_spawn=bool(data.get("can_spawn", False)),
            allowed_device_domains=tuple(data.get("allowed_device_domains") or ()),
            allowed_mcp_scopes=tuple(data.get("allowed_mcp_scopes") or ()),
        )


def intersect_capability_snapshots(
    left: CapabilitySnapshot,
    right: CapabilitySnapshot,
    *,
    source: CapabilitySource | None = None,
    trigger: CapabilityTrigger | None = None,
) -> CapabilitySnapshot:
    """Return the explicit intersection of two capability snapshots."""

    return CapabilitySnapshot(
        version=left.version,
        source=source or left.source,
        trigger=trigger or left.trigger,
        can_exec=left.can_exec and right.can_exec,
        can_read_files=left.can_read_files and right.can_read_files,
        can_write_files=left.can_write_files and right.can_write_files,
        can_send_cross_target=left.can_send_cross_target and right.can_send_cross_target,
        can_create_cron=left.can_create_cron and right.can_create_cron,
        can_spawn=left.can_spawn and right.can_spawn,
        allowed_device_domains=tuple(
            sorted(set(left.allowed_device_domains) & set(right.allowed_device_domains))
        ),
        allowed_mcp_scopes=tuple(
            sorted(set(left.allowed_mcp_scopes) & set(right.allowed_mcp_scopes))
        ),
    )
