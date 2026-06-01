"""Append-only outcome trace for governed self-evolution."""

from __future__ import annotations

import json
import os
import uuid
from contextlib import suppress
from dataclasses import asdict, dataclass, field
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

from filelock import FileLock

from OriginAgent.utils.helpers import ensure_dir, truncate_text

OUTCOME_SCHEMA_VERSION = "originagent.evolution.outcome.v1"
OUTCOME_STORE_RELATIVE = Path("memory") / "evolution_outcomes.jsonl"
OUTCOME_ARCHIVE_RELATIVE = Path("memory") / "archive" / "evolution_outcomes_archive.jsonl"

_MAX_TEXT_CHARS = 1000
_PERMANENT_EVENT_TYPES = {"promoted", "rolled_back", "review_rejected", "review_failed"}


@dataclass(frozen=True)
class EvolutionOutcomeEvent:
    """One immutable event in the signal-to-promotion evidence chain."""

    event_id: str
    timestamp: str
    type: str
    schema_version: str = OUTCOME_SCHEMA_VERSION
    opportunity_id: str = ""
    proposal_id: str = ""
    artifact_type: str = ""
    artifact_name: str = ""
    artifact_path: str = ""
    old_version: str = ""
    new_version: str = ""
    gate_decision: str = ""
    sandbox_status: str = ""
    review_status: str = ""
    promotion_status: str = ""
    rollback_status: str = ""
    feedback_score: float | None = None
    calibration_result: dict[str, Any] = field(default_factory=dict)
    metadata: dict[str, Any] = field(default_factory=dict)

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class EvolutionOutcomeStore:
    """Append-only JSONL trace for governed evolution outcomes."""

    def __init__(self, workspace: Path, *, path: Path | None = None) -> None:
        self.workspace = Path(workspace)
        memory_dir = ensure_dir(self.workspace / "memory")
        self.path = path or (self.workspace / OUTCOME_STORE_RELATIVE)
        self._lock = FileLock(str(memory_dir / ".evolution_outcomes.lock"))

    def append_event(
        self,
        event_type: str,
        *,
        opportunity_id: str = "",
        proposal_id: str = "",
        artifact_type: str = "",
        artifact_name: str = "",
        artifact_path: str = "",
        old_version: str = "",
        new_version: str = "",
        gate_decision: str = "",
        sandbox_status: str = "",
        review_status: str = "",
        promotion_status: str = "",
        rollback_status: str = "",
        feedback_score: float | None = None,
        calibration_result: dict[str, Any] | None = None,
        metadata: dict[str, Any] | None = None,
        timestamp: datetime | None = None,
    ) -> dict[str, Any]:
        event = EvolutionOutcomeEvent(
            event_id=f"evo_outcome_{uuid.uuid4().hex}",
            timestamp=_timestamp(timestamp),
            type=_clean(event_type, 128),
            opportunity_id=_clean(opportunity_id, 256),
            proposal_id=_clean(proposal_id, 256),
            artifact_type=_clean(artifact_type, 64),
            artifact_name=_clean(artifact_name, 256),
            artifact_path=_clean(_normalize_path(artifact_path), 512),
            old_version=_clean(old_version, 256),
            new_version=_clean(new_version, 256),
            gate_decision=_clean(gate_decision, 128),
            sandbox_status=_clean(sandbox_status, 128),
            review_status=_clean(review_status, 128),
            promotion_status=_clean(promotion_status, 128),
            rollback_status=_clean(rollback_status, 128),
            feedback_score=_safe_score(feedback_score),
            calibration_result=_redact_json(calibration_result or {}),
            metadata=_redact_json(metadata or {}),
        ).to_json()
        self.append(event)
        return event

    def append(self, event: dict[str, Any]) -> None:
        record = _redact_json(event)
        with self._lock:
            self.path.parent.mkdir(parents=True, exist_ok=True)
            with self.path.open("a", encoding="utf-8") as handle:
                handle.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")

    def append_many(self, events: list[dict[str, Any]]) -> int:
        if not events:
            return 0
        records = [_redact_json(event) for event in events if isinstance(event, dict)]
        if not records:
            return 0
        with self._lock:
            self.path.parent.mkdir(parents=True, exist_ok=True)
            with self.path.open("a", encoding="utf-8") as handle:
                for record in records:
                    handle.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")
        return len(records)

    def read_all(self) -> list[dict[str, Any]]:
        records: list[dict[str, Any]] = []
        with suppress(FileNotFoundError):
            with self.path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        raw = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(raw, dict):
                        records.append(raw)
        return records

    def compact_copy(self, records: list[dict[str, Any]]) -> None:
        """Rewrite the store from known-good records.

        This is intentionally narrow and is used by maintenance tasks that need
        deterministic test fixtures. Normal runtime writes are append-only.
        """

        tmp_path = self.path.with_suffix(self.path.suffix + ".tmp")
        with self._lock:
            self.path.parent.mkdir(parents=True, exist_ok=True)
            try:
                with tmp_path.open("w", encoding="utf-8") as handle:
                    for record in records:
                        handle.write(json.dumps(_redact_json(record), ensure_ascii=False, sort_keys=True) + "\n")
                    handle.flush()
                    os.fsync(handle.fileno())
                os.replace(tmp_path, self.path)
            except BaseException:
                tmp_path.unlink(missing_ok=True)
                raise

    def enforce_retention(
        self,
        *,
        retention_days: int = 90,
        archive: bool = True,
        now: datetime | None = None,
    ) -> dict[str, Any]:
        """Move old non-critical outcome records into an archive copy."""

        now_dt = _timestamp_dt(now)
        cutoff = now_dt - timedelta(days=max(1, int(retention_days or 90)))
        records = self.read_all()
        retained: list[dict[str, Any]] = []
        archived: list[dict[str, Any]] = []
        for record in records:
            event_time = _parse_datetime(str(record.get("timestamp") or ""))
            if (
                event_time is None
                or event_time >= cutoff
                or str(record.get("type") or "") in _PERMANENT_EVENT_TYPES
            ):
                retained.append(record)
            else:
                archived.append(record)
        if archived:
            if archive:
                self._append_archive(archived)
            self.compact_copy(retained)
        return {
            "retention_days": max(1, int(retention_days or 90)),
            "archive_enabled": bool(archive),
            "retained_count": len(retained),
            "archived_count": len(archived),
            "archive_path": str(OUTCOME_ARCHIVE_RELATIVE).replace("\\", "/") if archived and archive else "",
            "cutoff": cutoff.isoformat(),
        }

    def stats(self) -> dict[str, Any]:
        records = self.read_all()
        type_counts: dict[str, int] = {}
        gate_counts: dict[str, int] = {}
        sandbox_counts: dict[str, int] = {}
        review_counts: dict[str, int] = {}
        promotion_counts: dict[str, int] = {}
        rollback_counts: dict[str, int] = {}
        last_event_at = None
        for record in records:
            _bump(type_counts, record.get("type"))
            _bump(gate_counts, record.get("gate_decision"))
            _bump(sandbox_counts, record.get("sandbox_status"))
            _bump(review_counts, record.get("review_status"))
            _bump(promotion_counts, record.get("promotion_status"))
            _bump(rollback_counts, record.get("rollback_status"))
            timestamp = record.get("timestamp")
            if isinstance(timestamp, str) and (last_event_at is None or timestamp > last_event_at):
                last_event_at = timestamp
        return {
            "outcome_event_count": len(records),
            "outcome_type_counts": type_counts,
            "gate_decision_counts": gate_counts,
            "sandbox_status_counts": sandbox_counts,
            "review_status_counts": review_counts,
            "promotion_status_counts": promotion_counts,
            "rollback_status_counts": rollback_counts,
            "last_outcome_at": last_event_at,
            "archive": self.archive_stats(),
        }

    def archive_stats(self) -> dict[str, Any]:
        archive_path = self.workspace / OUTCOME_ARCHIVE_RELATIVE
        records = 0
        last_archived_at = None
        with suppress(FileNotFoundError):
            with archive_path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        raw = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if not isinstance(raw, dict):
                        continue
                    records += 1
                    archived_at = raw.get("archived_at")
                    if isinstance(archived_at, str) and (
                        last_archived_at is None or archived_at > last_archived_at
                    ):
                        last_archived_at = archived_at
        return {
            "archived_outcome_count": records,
            "last_archived_at": last_archived_at,
        }

    def _append_archive(self, records: list[dict[str, Any]]) -> None:
        if not records:
            return
        archive_path = self.workspace / OUTCOME_ARCHIVE_RELATIVE
        now = datetime.now(timezone.utc).isoformat()
        with self._lock:
            archive_path.parent.mkdir(parents=True, exist_ok=True)
            with archive_path.open("a", encoding="utf-8") as handle:
                for record in records:
                    archived = {
                        "archived_at": now,
                        "source_path": str(OUTCOME_STORE_RELATIVE).replace("\\", "/"),
                        "record": _redact_json(record),
                    }
                    handle.write(json.dumps(archived, ensure_ascii=False, sort_keys=True) + "\n")


def proposal_outcome_context(record: dict[str, Any], event: dict[str, Any] | None = None) -> dict[str, str]:
    """Extract trace identifiers from a review proposal and optional review event."""

    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
    artifact = None
    if event is not None and isinstance(event.get("artifact"), dict):
        artifact = event["artifact"]
    elif isinstance(record.get("apply_artifact"), dict):
        artifact = record["apply_artifact"]
    if artifact is None:
        artifact = {}
    artifact_type = str(artifact.get("artifact_type") or payload.get("subject_type") or record.get("proposal_type") or "")
    artifact_name = str(
        artifact.get("workflow_name")
        or artifact.get("skill_name")
        or payload.get("workflow_name")
        or payload.get("skill_name")
        or payload.get("subject_id")
        or ""
    )
    artifact_path = str(artifact.get("path") or payload.get("subject_path") or "")
    promotion_gate = payload.get("promotion_gate") if isinstance(payload.get("promotion_gate"), dict) else {}
    static_gate = payload.get("static_gate") if isinstance(payload.get("static_gate"), dict) else {}
    sandbox = payload.get("sandbox") if isinstance(payload.get("sandbox"), dict) else {}
    return {
        "opportunity_id": str(evolution.get("opportunity_id") or ""),
        "proposal_id": str(record.get("id") or (event or {}).get("proposal_id") or ""),
        "artifact_type": artifact_type,
        "artifact_name": artifact_name,
        "artifact_path": artifact_path,
        "gate_decision": str(promotion_gate.get("decision") or static_gate.get("decision") or ""),
        "sandbox_status": str(sandbox.get("status") or ""),
    }


def safe_append_outcome(store: EvolutionOutcomeStore, event_type: str, **kwargs: Any) -> dict[str, Any] | None:
    """Best-effort outcome append that must never affect governed actions."""

    try:
        return store.append_event(event_type, **kwargs)
    except Exception:
        return None


def _timestamp(value: datetime | None) -> str:
    return _timestamp_dt(value).isoformat()


def _timestamp_dt(value: datetime | None) -> datetime:
    if value is None:
        return datetime.now(timezone.utc)
    if value.tzinfo is None:
        return value.replace(tzinfo=timezone.utc)
    return value.astimezone(timezone.utc)


def _parse_datetime(value: str) -> datetime | None:
    if not value:
        return None
    with suppress(ValueError):
        parsed = datetime.fromisoformat(value)
        if parsed.tzinfo is None:
            return parsed.replace(tzinfo=timezone.utc)
        return parsed.astimezone(timezone.utc)
    return None


def _safe_score(value: float | None) -> float | None:
    if value is None:
        return None
    try:
        return max(0.0, min(1.0, float(value)))
    except (TypeError, ValueError):
        return None


def _clean(value: Any, max_chars: int) -> str:
    if value is None:
        return ""
    return truncate_text(_redact_text(str(value).strip()), max_chars)


def _normalize_path(path: str) -> str:
    return str(path or "").replace("\\", "/")


def _redact_json(value: Any) -> Any:
    if isinstance(value, str):
        return _clean(value, _MAX_TEXT_CHARS)
    if isinstance(value, list):
        return [_redact_json(item) for item in value]
    if isinstance(value, dict):
        return {str(key): _redact_json(item) for key, item in value.items()}
    return value


def _bump(counts: dict[str, int], value: Any) -> None:
    key = str(value or "").strip()
    if key:
        counts[key] = counts.get(key, 0) + 1


def _redact_text(text: str) -> str:
    try:
        from OriginAgent.agent.memory import redact_memory_text

        return redact_memory_text(text)
    except Exception:
        return text
