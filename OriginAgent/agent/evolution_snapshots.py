"""Governed snapshots and rollback for auto-evolution artifacts."""

from __future__ import annotations

import hashlib
import json
import uuid
from contextlib import suppress
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from filelock import FileLock

from OriginAgent.agent.evolution import AUTO_EVOLUTION_ORIGIN
from OriginAgent.agent.evolution_dependencies import EvolutionDependencyStore
from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore, safe_append_outcome
from OriginAgent.agent.skill_artifacts import validate_skill_artifact_dir
from OriginAgent.agent.workflow_artifacts import validate_workflow_artifact_dir
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger
from OriginAgent.utils.helpers import ensure_dir, truncate_text

SNAPSHOT_SCHEMA_VERSION = "originagent.evolution.snapshot.v1"
SNAPSHOT_ROOT_RELATIVE = Path("memory") / "evolution_snapshots"


@dataclass(frozen=True)
class EvolutionSnapshotRecord:
    """Metadata for one immutable artifact snapshot."""

    snapshot_id: str
    schema_version: str
    created_at: str
    artifact_type: str
    artifact_name: str
    artifact_path: str
    snapshot_path: str
    content_hash: str
    proposal_id: str = ""
    opportunity_id: str = ""
    source_event_id: str = ""
    created_by: str = AUTO_EVOLUTION_ORIGIN
    reason: str = ""

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class EvolutionRollbackResult:
    """Outcome for a governed rollback attempt."""

    ok: bool
    status: str
    message: str
    artifact_type: str = ""
    artifact_name: str = ""
    artifact_path: str = ""
    snapshot_id: str = ""
    old_version: str = ""
    new_version: str = ""
    event: dict[str, Any] | None = None
    error: str = ""
    dependency_blockers: list[dict[str, Any]] = field(default_factory=list)

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class EvolutionSnapshotStore:
    """Store immutable copies of governed workflow and skill artifacts."""

    def __init__(self, workspace: Path) -> None:
        self.workspace = Path(workspace)
        self.root = ensure_dir(self.workspace / SNAPSHOT_ROOT_RELATIVE)
        self._lock = FileLock(str(self.root / ".evolution_snapshots.lock"))

    def create_snapshot(
        self,
        artifact: dict[str, Any],
        *,
        proposal_id: str = "",
        opportunity_id: str = "",
        source_event_id: str = "",
        reason: str = "",
        created_by: str = AUTO_EVOLUTION_ORIGIN,
    ) -> dict[str, Any] | None:
        artifact_type = _artifact_type(artifact)
        artifact_name = _artifact_name(artifact)
        artifact_path = _artifact_path(artifact)
        if artifact_type not in {"workflow", "skill"} or not artifact_name or not artifact_path:
            return None
        source_file = self.workspace / artifact_path
        try:
            content = source_file.read_text(encoding="utf-8")
        except OSError:
            return None
        content_hash = _content_hash(content)
        snapshot_dir = self.root / artifact_type / artifact_name / content_hash
        artifact_filename = _artifact_filename(artifact_type)
        record = EvolutionSnapshotRecord(
            snapshot_id=f"evo_snapshot_{uuid.uuid4().hex}",
            schema_version=SNAPSHOT_SCHEMA_VERSION,
            created_at=datetime.now(timezone.utc).isoformat(),
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            artifact_path=artifact_path.replace("\\", "/"),
            snapshot_path=_relative_to_workspace(snapshot_dir / artifact_filename, self.workspace),
            content_hash=content_hash,
            proposal_id=_clean(proposal_id, 256),
            opportunity_id=_clean(opportunity_id, 256),
            source_event_id=_clean(source_event_id, 256),
            created_by=_clean(created_by, 128) or AUTO_EVOLUTION_ORIGIN,
            reason=_clean(reason, 512),
        )
        with self._lock:
            snapshot_dir.mkdir(parents=True, exist_ok=True)
            snapshot_file = snapshot_dir / artifact_filename
            metadata_file = snapshot_dir / "version_metadata.json"
            if not snapshot_file.exists():
                snapshot_file.write_text(content, encoding="utf-8")
            if not metadata_file.exists():
                metadata_file.write_text(
                    json.dumps(record.to_json(), ensure_ascii=False, sort_keys=True, indent=2) + "\n",
                    encoding="utf-8",
                )
            else:
                try:
                    existing = json.loads(metadata_file.read_text(encoding="utf-8"))
                except (OSError, json.JSONDecodeError):
                    existing = None
                if isinstance(existing, dict):
                    return existing
        return record.to_json()

    def latest_snapshot(
        self,
        *,
        artifact_type: str,
        artifact_name: str,
        before_hash: str | None = None,
    ) -> dict[str, Any] | None:
        records = [
            record for record in self.list_snapshots(artifact_type=artifact_type, artifact_name=artifact_name)
            if not before_hash or record.get("content_hash") != before_hash
        ]
        if not records:
            return None
        records.sort(key=lambda record: str(record.get("created_at") or ""), reverse=True)
        return records[0]

    def list_snapshots(
        self,
        *,
        artifact_type: str | None = None,
        artifact_name: str | None = None,
    ) -> list[dict[str, Any]]:
        base = self.root
        if artifact_type:
            base = base / artifact_type
            if artifact_name:
                base = base / artifact_name
        records: list[dict[str, Any]] = []
        with suppress(FileNotFoundError):
            for metadata_file in base.rglob("version_metadata.json"):
                try:
                    raw = json.loads(metadata_file.read_text(encoding="utf-8"))
                except (OSError, json.JSONDecodeError):
                    continue
                if isinstance(raw, dict):
                    records.append(raw)
        return records

    def stats(self) -> dict[str, Any]:
        records = self.list_snapshots()
        type_counts: dict[str, int] = {}
        last_created_at = None
        for record in records:
            artifact_type = str(record.get("artifact_type") or "")
            if artifact_type:
                type_counts[artifact_type] = type_counts.get(artifact_type, 0) + 1
            created_at = record.get("created_at")
            if isinstance(created_at, str) and (last_created_at is None or created_at > last_created_at):
                last_created_at = created_at
        return {
            "snapshot_count": len(records),
            "snapshot_type_counts": type_counts,
            "last_snapshot_at": last_created_at,
        }


class EvolutionRollbackService:
    """Restore workspace workflow/skill artifacts from governed snapshots."""

    def __init__(self, workspace: Path) -> None:
        self.workspace = Path(workspace)
        self.snapshots = EvolutionSnapshotStore(self.workspace)
        self.outcomes = EvolutionOutcomeStore(self.workspace)

    def rollback(
        self,
        *,
        artifact_type: str,
        artifact_name: str,
        snapshot_id: str | None = None,
        reason: str = "",
        actor: str = "user",
        force: bool = False,
    ) -> EvolutionRollbackResult:
        artifact_type = str(artifact_type or "").strip().lower()
        artifact_name = str(artifact_name or "").strip()
        if artifact_type not in {"workflow", "skill"} or not artifact_name:
            return EvolutionRollbackResult(
                ok=False,
                status="invalid",
                message="artifact_type and artifact_name are required.",
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                error="invalid_target",
            )
        target_file = _target_file(self.workspace, artifact_type, artifact_name)
        if not target_file.is_file():
            return EvolutionRollbackResult(
                ok=False,
                status="missing",
                message="Target artifact does not exist.",
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                error="target_missing",
            )
        current_content = target_file.read_text(encoding="utf-8")
        current_hash = _content_hash(current_content)
        snapshot = self._select_snapshot(
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            snapshot_id=snapshot_id,
            before_hash=current_hash,
        )
        if snapshot is None:
            return EvolutionRollbackResult(
                ok=False,
                status="missing_snapshot",
                message="No rollback snapshot is available.",
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                old_version=current_hash,
                error="snapshot_missing",
            )
        snapshot_file = self.workspace / str(snapshot.get("snapshot_path") or "")
        try:
            snapshot_content = snapshot_file.read_text(encoding="utf-8")
        except OSError as exc:
            return EvolutionRollbackResult(
                ok=False,
                status="failed",
                message="Could not read rollback snapshot.",
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                snapshot_id=str(snapshot.get("snapshot_id") or ""),
                old_version=current_hash,
                error=str(exc),
            )
        new_hash = _content_hash(snapshot_content)
        if new_hash == current_hash:
            return EvolutionRollbackResult(
                ok=True,
                status="already_current",
                message="Target artifact already matches the requested snapshot.",
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                snapshot_id=str(snapshot.get("snapshot_id") or ""),
                old_version=current_hash,
                new_version=new_hash,
            )

        dependency_store = EvolutionDependencyStore(self.workspace)
        blockers = dependency_store.rollback_blockers(
            artifact_type=artifact_type,
            artifact_name=artifact_name,
        )
        if blockers and not force:
            message = "Rollback blocked because other governed artifacts depend on this artifact."
            safe_append_outcome(
                self.outcomes,
                "rolled_back",
                opportunity_id=str(snapshot.get("opportunity_id") or ""),
                proposal_id=str(snapshot.get("proposal_id") or ""),
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                old_version=current_hash,
                new_version=new_hash,
                rollback_status="blocked",
                metadata={"snapshot_id": snapshot.get("snapshot_id"), "reason": reason, "blockers": blockers},
            )
            return EvolutionRollbackResult(
                ok=False,
                status="blocked_by_dependencies",
                message=message,
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                snapshot_id=str(snapshot.get("snapshot_id") or ""),
                old_version=current_hash,
                new_version=new_hash,
                error="dependency_blocked",
                dependency_blockers=blockers,
            )

        self.snapshots.create_snapshot(
            _artifact_for_target(artifact_type, artifact_name, target_file, self.workspace),
            proposal_id=str(snapshot.get("proposal_id") or ""),
            opportunity_id=str(snapshot.get("opportunity_id") or ""),
            reason="pre_rollback",
            created_by=actor,
        )
        started_event = self._append_ledger_event(
            EventType.MODULE_ROLLBACK_STARTED,
            actor=actor,
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            source_event_id=str(snapshot.get("snapshot_id") or ""),
            result={"snapshot_id": snapshot.get("snapshot_id"), "old_version": current_hash, "new_version": new_hash},
        )
        try:
            target_file.write_text(snapshot_content, encoding="utf-8")
            self._validate_restored(artifact_type, artifact_name)
        except Exception as exc:
            target_file.write_text(current_content, encoding="utf-8")
            failed_event = self._append_ledger_event(
                EventType.MODULE_ROLLBACK_FAILED,
                actor=actor,
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                source_event_id=started_event.get("event_id") if isinstance(started_event, dict) else "",
                result={"snapshot_id": snapshot.get("snapshot_id"), "error": str(exc)},
            )
            safe_append_outcome(
                self.outcomes,
                "rolled_back",
                opportunity_id=str(snapshot.get("opportunity_id") or ""),
                proposal_id=str(snapshot.get("proposal_id") or ""),
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                old_version=current_hash,
                new_version=new_hash,
                rollback_status="failed",
                metadata={"snapshot_id": snapshot.get("snapshot_id"), "event": failed_event, "reason": reason},
            )
            return EvolutionRollbackResult(
                ok=False,
                status="failed",
                message="Rollback failed and the previous artifact was restored.",
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
                snapshot_id=str(snapshot.get("snapshot_id") or ""),
                old_version=current_hash,
                new_version=new_hash,
                event=failed_event,
                error=str(exc),
            )

        succeeded_event = self._append_ledger_event(
            EventType.MODULE_ROLLBACK_SUCCEEDED,
            actor=actor,
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            source_event_id=started_event.get("event_id") if isinstance(started_event, dict) else "",
            result={"snapshot_id": snapshot.get("snapshot_id"), "old_version": current_hash, "new_version": new_hash},
        )
        safe_append_outcome(
            self.outcomes,
            "rolled_back",
            opportunity_id=str(snapshot.get("opportunity_id") or ""),
            proposal_id=str(snapshot.get("proposal_id") or ""),
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            artifact_path=_relative_to_workspace(target_file, self.workspace),
            old_version=current_hash,
            new_version=new_hash,
            rollback_status="succeeded",
            metadata={"snapshot_id": snapshot.get("snapshot_id"), "event": succeeded_event, "reason": reason},
        )
        try:
            dependency_store.update_artifact(
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                artifact_path=_relative_to_workspace(target_file, self.workspace),
            )
        except Exception:
            pass
        return EvolutionRollbackResult(
            ok=True,
            status="rolled_back",
            message="Artifact restored from snapshot.",
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            artifact_path=_relative_to_workspace(target_file, self.workspace),
            snapshot_id=str(snapshot.get("snapshot_id") or ""),
            old_version=current_hash,
            new_version=new_hash,
            event=succeeded_event,
        )

    def _select_snapshot(
        self,
        *,
        artifact_type: str,
        artifact_name: str,
        snapshot_id: str | None,
        before_hash: str,
    ) -> dict[str, Any] | None:
        records = self.snapshots.list_snapshots(artifact_type=artifact_type, artifact_name=artifact_name)
        if snapshot_id:
            for record in records:
                if record.get("snapshot_id") == snapshot_id or record.get("content_hash") == snapshot_id:
                    return record
            return None
        return self.snapshots.latest_snapshot(
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            before_hash=before_hash,
        )

    def _validate_restored(self, artifact_type: str, artifact_name: str) -> None:
        if artifact_type == "workflow":
            valid, message = validate_workflow_artifact_dir(
                self.workspace / "workflows" / artifact_name,
                workspace=self.workspace,
                expected_name=artifact_name,
            )
        else:
            valid, message = validate_skill_artifact_dir(
                self.workspace / "skills" / artifact_name,
                workspace=self.workspace,
                expected_name=artifact_name,
            )
        if not valid:
            raise ValueError(message)

    def _append_ledger_event(
        self,
        event_type: EventType,
        *,
        actor: str,
        artifact_type: str,
        artifact_name: str,
        source_event_id: str,
        result: dict[str, Any],
    ) -> dict[str, Any]:
        try:
            event = EvolutionLedger(self.workspace).append(EvolutionEvent.new(
                event_type,
                actor=actor,
                module_id=artifact_name,
                module_type=artifact_type,
                source_event_stream="evolution",
                source_event_id=source_event_id,
                result=result,
            ))
            return event.to_dict()
        except Exception as exc:
            return {"error": str(exc)}


def snapshot_artifact_if_governed(
    workspace: Path,
    record: dict[str, Any],
    event: dict[str, Any],
) -> dict[str, Any] | None:
    """Create a snapshot for auto-evolution workflow/skill artifacts."""

    if not _is_governed_record(record):
        return None
    artifact = event.get("artifact") if isinstance(event.get("artifact"), dict) else {}
    if _artifact_type(artifact) not in {"workflow", "skill"}:
        return None
    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
    snapshot = EvolutionSnapshotStore(workspace).create_snapshot(
        artifact,
        proposal_id=str(record.get("id") or event.get("proposal_id") or ""),
        opportunity_id=str(evolution.get("opportunity_id") or ""),
        source_event_id=str(event.get("event_id") or ""),
        reason="promotion_snapshot",
        created_by=str(evolution.get("origin") or AUTO_EVOLUTION_ORIGIN),
    )
    if snapshot is not None:
        safe_append_outcome(
            EvolutionOutcomeStore(workspace),
            "snapshot_created",
            opportunity_id=str(evolution.get("opportunity_id") or ""),
            proposal_id=str(record.get("id") or event.get("proposal_id") or ""),
            artifact_type=str(snapshot.get("artifact_type") or ""),
            artifact_name=str(snapshot.get("artifact_name") or ""),
            artifact_path=str(snapshot.get("artifact_path") or ""),
            new_version=str(snapshot.get("content_hash") or ""),
            metadata={"snapshot_id": snapshot.get("snapshot_id"), "snapshot_path": snapshot.get("snapshot_path")},
        )
    return snapshot


def _is_governed_record(record: dict[str, Any]) -> bool:
    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
    return (
        str(record.get("origin") or "").strip().lower() == AUTO_EVOLUTION_ORIGIN
        or str(evolution.get("origin") or "").strip().lower() == AUTO_EVOLUTION_ORIGIN
    )


def _artifact_type(artifact: dict[str, Any]) -> str:
    raw = str(artifact.get("artifact_type") or "").strip().lower()
    if raw:
        return raw
    if artifact.get("skill_name"):
        return "skill"
    if artifact.get("workflow_name"):
        return "workflow"
    return ""


def _artifact_name(artifact: dict[str, Any]) -> str:
    artifact_type = _artifact_type(artifact)
    if artifact_type == "workflow":
        return _clean(artifact.get("workflow_name"), 128)
    if artifact_type == "skill":
        return _clean(artifact.get("skill_name"), 128)
    return ""


def _artifact_path(artifact: dict[str, Any]) -> str:
    return str(artifact.get("path") or "").replace("\\", "/")


def _artifact_filename(artifact_type: str) -> str:
    return "workflow.yaml" if artifact_type == "workflow" else "SKILL.md"


def _target_file(workspace: Path, artifact_type: str, artifact_name: str) -> Path:
    if artifact_type == "workflow":
        return workspace / "workflows" / artifact_name / "workflow.yaml"
    return workspace / "skills" / artifact_name / "SKILL.md"


def _artifact_for_target(artifact_type: str, artifact_name: str, target_file: Path, workspace: Path) -> dict[str, Any]:
    relative = _relative_to_workspace(target_file, workspace)
    if artifact_type == "workflow":
        return {
            "artifact_type": "workflow",
            "workflow_name": artifact_name,
            "path": relative,
        }
    return {
        "artifact_type": "skill",
        "skill_name": artifact_name,
        "path": relative,
    }


def _content_hash(content: str) -> str:
    return hashlib.sha256(content.encode("utf-8")).hexdigest()


def _relative_to_workspace(path: Path, workspace: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return path.as_posix()


def _clean(value: Any, max_chars: int) -> str:
    return truncate_text(str(value or "").strip(), max_chars)
