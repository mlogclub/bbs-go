"""Runtime capability snapshots for active evolution modules."""

from __future__ import annotations

import json
from dataclasses import dataclass
from pathlib import Path
from typing import Any

from OriginAgent.evolution.activation import ACTIVATION_SCHEMA_VERSION
from OriginAgent.evolution.events import EventType
from OriginAgent.evolution.ledger import EvolutionLedger
from OriginAgent.evolution.package import read_package_manifest
from OriginAgent.evolution.verifier import EvolutionModuleVerifier
from OriginAgent.security.capabilities import CapabilitySnapshot, intersect_capability_snapshots


@dataclass(frozen=True)
class EvolutionCapabilityResult:
    ok: bool
    status: str
    artifact_digest: str = ""
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    snapshot: CapabilitySnapshot | None = None
    error: str = ""


class EvolutionCapabilityGate:
    """Build least-privilege runtime capabilities for active evolution modules."""

    def __init__(self, workspace: Path, ledger: EvolutionLedger | None = None) -> None:
        self.workspace = Path(workspace)
        self.memory_dir = self.workspace / "memory"
        self.staging_root = self.memory_dir / "evolution_staging"
        self.activation_root = self.memory_dir / "evolution_activations"
        self.ledger = ledger or EvolutionLedger(self.workspace)

    def snapshot_for_artifact(
        self,
        artifact_digest: str,
        *,
        base_snapshot: CapabilitySnapshot | None = None,
    ) -> EvolutionCapabilityResult:
        activation = self._read_active_activation(artifact_digest)
        if activation is None:
            return EvolutionCapabilityResult(
                ok=False,
                status="not_active",
                artifact_digest=artifact_digest,
                error="artifact is not active",
            )
        if not self._has_verified_event(artifact_digest):
            return _result_from_activation(
                activation,
                ok=False,
                status="unverified",
                error="artifact_digest has no verified module event",
            )

        report = EvolutionModuleVerifier(self.workspace, staging_root=self.staging_root).verify(artifact_digest)
        if not report.ok:
            return _result_from_activation(
                activation,
                ok=False,
                status="verification_failed",
                error=report.error,
            )

        artifact_dir = self.staging_root / artifact_digest / "artifact"
        try:
            manifest = read_package_manifest(artifact_dir)
        except Exception as exc:
            return _result_from_activation(
                activation,
                ok=False,
                status="manifest_failed",
                error=_safe_error(exc, self.workspace),
            )

        module_snapshot = _snapshot_from_permissions(manifest.permissions)
        snapshot = (
            intersect_capability_snapshots(
                base_snapshot,
                module_snapshot,
                source=base_snapshot.source,
                trigger=base_snapshot.trigger,
            )
            if base_snapshot is not None
            else module_snapshot
        )
        return EvolutionCapabilityResult(
            ok=True,
            status="ready",
            artifact_digest=artifact_digest,
            module_id=manifest.module_id,
            module_type=manifest.module_type,
            module_version=manifest.version,
            snapshot=snapshot,
        )

    def snapshot_for_domain_pack(
        self,
        pack_id: str,
        *,
        base_snapshot: CapabilitySnapshot | None = None,
    ) -> CapabilitySnapshot | None:
        artifact_digest = self._active_domain_pack_index().get(pack_id)
        if not artifact_digest:
            return None
        result = self.snapshot_for_artifact(artifact_digest, base_snapshot=base_snapshot)
        return result.snapshot if result.ok else None

    def _active_domain_pack_index(self) -> dict[str, str]:
        rows: list[tuple[str, str, str]] = []
        if not self.activation_root.is_dir():
            return {}
        for digest_dir in self.activation_root.iterdir():
            activation_json = digest_dir / "activation.json"
            if not activation_json.exists():
                continue
            try:
                data = json.loads(activation_json.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError):
                continue
            if (
                isinstance(data, dict)
                and data.get("schema_version") == ACTIVATION_SCHEMA_VERSION
                and data.get("status") == "active"
                and data.get("module_type") == "domain_pack"
                and isinstance(data.get("module_id"), str)
            ):
                rows.append((data["module_id"], str(data.get("activated_at") or ""), digest_dir.name))
        index: dict[str, str] = {}
        for module_id, _activated_at, artifact_digest in sorted(rows):
            index[module_id] = artifact_digest
        return index

    def _read_active_activation(self, artifact_digest: str) -> dict[str, Any] | None:
        path = self.activation_root / artifact_digest / "activation.json"
        if not path.exists():
            return None
        try:
            data = json.loads(path.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            return None
        if not isinstance(data, dict):
            return None
        if data.get("schema_version") != ACTIVATION_SCHEMA_VERSION:
            return None
        if data.get("artifact_digest") != artifact_digest:
            return None
        if data.get("status") != "active":
            return None
        return data

    def _has_verified_event(self, artifact_digest: str) -> bool:
        verification = self.ledger.verify_chain()
        if not verification.ok or not self.ledger.event_path.exists():
            return False
        with self.ledger._locked():
            with self.ledger.event_path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    if not line.strip():
                        continue
                    event = json.loads(line)
                    if (
                        event.get("event_type") == EventType.MODULE_VERIFIED.value
                        and event.get("artifact_digest") == artifact_digest
                    ):
                        return True
        return False


def _snapshot_from_permissions(permissions: dict[str, Any]) -> CapabilitySnapshot:
    return CapabilitySnapshot(
        version=1,
        source="system",
        trigger="system",
        can_exec=bool(permissions.get("exec", False)),
        can_read_files=bool(permissions.get("read_files", False)),
        can_write_files=bool(permissions.get("write_files", False)),
        can_send_cross_target=bool(permissions.get("send_cross_target", False)),
        can_create_cron=bool(permissions.get("create_cron", False)),
        can_spawn=bool(permissions.get("spawn", False)),
        allowed_device_domains=tuple(str(item) for item in permissions.get("device_domains") or ()),
        allowed_mcp_scopes=tuple(str(item) for item in permissions.get("mcp_scopes") or ()),
    )


def _result_from_activation(
    activation: dict[str, Any],
    *,
    ok: bool,
    status: str,
    error: str = "",
) -> EvolutionCapabilityResult:
    return EvolutionCapabilityResult(
        ok=ok,
        status=status,
        artifact_digest=str(activation.get("artifact_digest") or ""),
        module_id=str(activation.get("module_id") or ""),
        module_type=str(activation.get("module_type") or ""),
        module_version=str(activation.get("version") or ""),
        error=error,
    )


def _safe_error(exc: Exception, workspace: Path) -> str:
    message = str(exc)
    sensitive = str(Path(workspace).expanduser().resolve(strict=False))
    return message.replace(sensitive, "<path>") if sensitive else message
