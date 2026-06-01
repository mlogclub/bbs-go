"""Dirty rollback recovery bookkeeping for evolution activations."""

from __future__ import annotations

import json
import os
import uuid
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Iterable

from filelock import FileLock

from OriginAgent.evolution.activation import ACTIVATION_SCHEMA_VERSION
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger, canonical_dump
from OriginAgent.evolution.telemetry import sanitize_telemetry_text
from OriginAgent.utils.helpers import truncate_text


@dataclass(frozen=True)
class EvolutionRecoveryResult:
    ok: bool
    status: str
    artifact_digest: str = ""
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    events: tuple[EvolutionEvent, ...] = ()
    error: str = ""


class EvolutionRecoveryManager:
    """Record cleanup outcomes without executing third-party teardown code."""

    def __init__(
        self,
        workspace: Path,
        ledger: EvolutionLedger | None = None,
        lock_path: Path | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self.memory_dir = self.workspace / "memory"
        self.activation_root = self.memory_dir / "evolution_activations"
        self.ledger = ledger or EvolutionLedger(self.workspace)
        self._lock_path = Path(lock_path) if lock_path is not None else self.memory_dir / ".evolution_activation.lock"

    def record_teardown(
        self,
        artifact_digest: str,
        *,
        succeeded: bool,
        reason: str = "",
        residual_resources: Iterable[str] = (),
        actor: str = "user",
    ) -> EvolutionRecoveryResult:
        with self._locked():
            metadata = self._read_metadata(artifact_digest)
            if metadata is None:
                return EvolutionRecoveryResult(
                    ok=False,
                    status="not_active",
                    artifact_digest=artifact_digest,
                    error="activation metadata is missing",
                )
            if metadata.get("status") != "dirty_rollback":
                return _result_from_metadata(
                    metadata,
                    ok=False,
                    status="not_dirty_rollback",
                    error=f"activation is not dirty_rollback: {metadata.get('status')}",
                )
            safe_reason = _safe(reason, self.workspace)
            safe_residuals = [_safe(item, self.workspace) for item in residual_resources]
            started = self.ledger.append(
                EvolutionEvent.new(
                    EventType.TEARDOWN_STARTED,
                    **_event_common(actor, metadata),
                    result={
                        "status": "teardown_started",
                        "reason": safe_reason,
                        "residual_resources": safe_residuals,
                    },
                )
            )
            now = datetime.now(timezone.utc).isoformat()
            attempts = list(metadata.get("teardown_attempts") or [])
            attempts.append(
                {
                    "attempted_at": now,
                    "succeeded": bool(succeeded),
                    "reason": safe_reason,
                    "residual_resources": safe_residuals,
                }
            )
            metadata["teardown_attempts"] = attempts
            events = [started]
            if succeeded:
                metadata["status"] = "rolled_back"
                metadata["rolled_back_at"] = metadata.get("rolled_back_at") or now
                metadata["residual_resources"] = []
                event_type = EventType.TEARDOWN_SUCCEEDED
                status = "teardown_succeeded"
                ok = True
                error = ""
            else:
                metadata["status"] = "dirty_rollback"
                metadata["dirty_reason"] = safe_reason
                metadata["residual_resources"] = safe_residuals
                event_type = EventType.TEARDOWN_FAILED
                status = "teardown_failed"
                ok = False
                error = safe_reason or "teardown failed"
            terminal = self.ledger.append(
                EvolutionEvent.new(
                    event_type,
                    **_event_common(actor, metadata),
                    result={
                        "status": status,
                        "reason": safe_reason,
                        "residual_resources": safe_residuals,
                    },
                )
            )
            events.append(terminal)
            _write_json_atomic(self._metadata_path(artifact_digest), metadata)
            return _result_from_metadata(metadata, ok=ok, status=status, events=tuple(events), error=error)

    def force_clean(
        self,
        artifact_digest: str,
        *,
        reason: str,
        actor: str = "user",
    ) -> EvolutionRecoveryResult:
        with self._locked():
            metadata = self._read_metadata(artifact_digest)
            if metadata is None:
                return EvolutionRecoveryResult(
                    ok=False,
                    status="not_active",
                    artifact_digest=artifact_digest,
                    error="activation metadata is missing",
                )
            if metadata.get("status") != "dirty_rollback":
                return _result_from_metadata(
                    metadata,
                    ok=False,
                    status="not_dirty_rollback",
                    error=f"activation is not dirty_rollback: {metadata.get('status')}",
                )
            safe_reason = _safe(reason, self.workspace)
            requested = self.ledger.append(
                EvolutionEvent.new(
                    EventType.MODULE_FORCE_CLEAN_REQUESTED,
                    **_event_common(actor, metadata),
                    result={"status": "force_clean_requested", "reason": safe_reason},
                )
            )
            now = datetime.now(timezone.utc).isoformat()
            metadata["status"] = "rolled_back"
            metadata["rolled_back_at"] = metadata.get("rolled_back_at") or now
            metadata["force_cleaned_at"] = now
            metadata["dirty_reason"] = safe_reason
            metadata["residual_resources"] = []
            _write_json_atomic(self._metadata_path(artifact_digest), metadata)
            succeeded = self.ledger.append(
                EvolutionEvent.new(
                    EventType.MODULE_FORCE_CLEAN_SUCCEEDED,
                    **_event_common(actor, metadata),
                    result={"status": "force_clean_succeeded", "reason": safe_reason},
                )
            )
            return _result_from_metadata(
                metadata,
                ok=True,
                status="force_clean_succeeded",
                events=(requested, succeeded),
            )

    def _read_metadata(self, artifact_digest: str) -> dict[str, Any] | None:
        path = self._metadata_path(artifact_digest)
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
        return data

    def _metadata_path(self, artifact_digest: str) -> Path:
        return self.activation_root / artifact_digest / "activation.json"

    def _locked(self) -> FileLock:
        self.memory_dir.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))


def _event_common(actor: str, metadata: dict[str, Any]) -> dict[str, str]:
    return {
        "actor": actor,
        "module_id": str(metadata.get("module_id") or ""),
        "module_version": str(metadata.get("version") or ""),
        "module_type": str(metadata.get("module_type") or ""),
        "artifact_digest": str(metadata.get("artifact_digest") or ""),
    }


def _result_from_metadata(
    metadata: dict[str, Any],
    *,
    ok: bool,
    status: str,
    events: tuple[EvolutionEvent, ...] = (),
    error: str = "",
) -> EvolutionRecoveryResult:
    return EvolutionRecoveryResult(
        ok=ok,
        status=status,
        artifact_digest=str(metadata.get("artifact_digest") or ""),
        module_id=str(metadata.get("module_id") or ""),
        module_type=str(metadata.get("module_type") or ""),
        module_version=str(metadata.get("version") or ""),
        events=events,
        error=error,
    )


def _safe(value: Any, workspace: Path) -> str:
    return truncate_text(sanitize_telemetry_text(str(value or ""), workspace), 300)[:300]


def _write_json_atomic(path: Path, data: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = path.with_name(f".{path.name}.tmp-{uuid.uuid4().hex}")
    try:
        with tmp_path.open("wb") as handle:
            handle.write(canonical_dump(data))
            handle.write(b"\n")
        os.replace(tmp_path, path)
    finally:
        if tmp_path.exists():
            tmp_path.unlink()
