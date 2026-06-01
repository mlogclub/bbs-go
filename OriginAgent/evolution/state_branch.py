"""Copy-on-write facts state branches for staged evolution modules."""

from __future__ import annotations

import hashlib
import json
import os
import uuid
from contextlib import suppress
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from filelock import FileLock

from OriginAgent.agent.facts import FactRecord, FactStore, render_memory_md
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger, canonical_dump

BRANCH_SCHEMA_VERSION = "originagent.evolution.state_branch.v1"


@dataclass(frozen=True)
class EvolutionMergeConflict:
    fact_id: str
    reason: str
    stable_hash: str
    base_hash: str

    def to_dict(self) -> dict[str, str]:
        return {
            "fact_id": self.fact_id,
            "reason": self.reason,
            "stable_hash": self.stable_hash,
            "base_hash": self.base_hash,
        }


@dataclass(frozen=True)
class EvolutionMergePreview:
    ok: bool
    branch_id: str
    added: tuple[str, ...] = ()
    modified: tuple[str, ...] = ()
    deprecated: tuple[str, ...] = ()
    conflicts: tuple[EvolutionMergeConflict, ...] = ()


@dataclass(frozen=True)
class EvolutionStateBranchResult:
    ok: bool
    status: str
    branch_id: str = ""
    artifact_digest: str = ""
    changed_fact_ids: tuple[str, ...] = ()
    conflicts: tuple[EvolutionMergeConflict, ...] = ()
    events: tuple[EvolutionEvent, ...] = ()
    error: str = ""


class EvolutionStateBranchStore:
    """Manage local CoW facts branches without executing module code."""

    def __init__(
        self,
        workspace: Path,
        ledger: EvolutionLedger | None = None,
        lock_path: Path | None = None,
        lock_factory: Callable[[], FileLock] | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self.memory_dir = self.workspace / "memory"
        self.branch_root = self.memory_dir / "evolution_branches"
        self.ledger = ledger or EvolutionLedger(self.workspace)
        self._lock_path = Path(lock_path) if lock_path is not None else self.memory_dir / ".lock"
        self._lock_factory = lock_factory
        self.fact_store = FactStore(
            self.workspace,
            lock_factory=self._locked,
        )

    def create_branch(self, artifact_digest: str, *, actor: str = "user") -> EvolutionStateBranchResult:
        with self._locked():
            verified = self._find_verified_event(artifact_digest)
            if verified is None:
                return EvolutionStateBranchResult(
                    ok=False,
                    status="failed",
                    artifact_digest=artifact_digest,
                    error="artifact_digest has no verified module event",
                )
            branch_id = f"branch_{uuid.uuid4().hex}"
            branch_dir = self._branch_dir(branch_id)
            stable_facts = self.fact_store.read_all_unlocked()
            metadata = {
                "schema_version": BRANCH_SCHEMA_VERSION,
                "branch_id": branch_id,
                "artifact_digest": artifact_digest,
                "created_at": datetime.now(timezone.utc).isoformat(),
                "status": "active",
                "base_event_hash": verified["terminal_event_hash"],
                "base_facts_hash": _facts_hash(stable_facts),
                "base_fact_hashes": _fact_hashes(stable_facts),
            }
            branch_dir.mkdir(parents=True, exist_ok=False)
            _write_json_atomic(branch_dir / "branch.json", metadata)
            _write_fact_records_atomic(branch_dir / "facts_overlay.jsonl", [])
            _write_tombstones_atomic(branch_dir / "facts_tombstones.jsonl", [])
            event = self.ledger.append(
                EvolutionEvent.new(
                    EventType.STATE_BRANCH_CREATED,
                    actor=actor,
                    artifact_digest=artifact_digest,
                    state_branch_id=branch_id,
                    result={
                        "status": "created",
                        "branch_id": branch_id,
                        "artifact_digest": artifact_digest,
                        "base_event_hash": metadata["base_event_hash"],
                        "base_facts_hash": metadata["base_facts_hash"],
                    },
                )
            )
            return EvolutionStateBranchResult(
                ok=True,
                status="created",
                branch_id=branch_id,
                artifact_digest=artifact_digest,
                events=(event,),
            )

    def read_facts(self, branch_id: str | None = None) -> list[FactRecord]:
        """Return FactRecord copies; mutating them never mutates stored facts."""

        with self._locked():
            stable = [_copy_fact(record) for record in self.fact_store.read_all_unlocked()]
            if branch_id is None:
                return stable
            metadata = self._load_branch(branch_id)
            if metadata.get("status") not in {"active", "merged", "discarded"}:
                raise ValueError("unsupported branch status")
            return self._compose_facts_unlocked(branch_id, stable)

    def upsert_fact(
        self,
        branch_id: str,
        content: str,
        *,
        category: str = "note",
        scope: str = "general",
        owner: str = "unknown",
        source_cursors: list[int] | None = None,
        source_excerpt: str = "",
        confidence: float = 1.0,
        expires_at: str | None = None,
        requires_confirmation: bool | None = None,
        status: str | None = None,
        supersedes_fact_id: str | None = None,
    ) -> EvolutionStateBranchResult:
        try:
            with self._locked():
                metadata = self._load_active_branch(branch_id)
                before = {record.fact_id: record.to_dict() for record in self._compose_facts_unlocked(branch_id)}
                records = self._compose_facts_unlocked(branch_id)
                self.fact_store.upsert_fact_in_records_unlocked(
                    records,
                    content,
                    category=category,
                    scope=scope,
                    owner=owner,
                    source_cursors=source_cursors,
                    source_excerpt=source_excerpt,
                    confidence=confidence,
                    expires_at=expires_at,
                    requires_confirmation=requires_confirmation,
                    status=status,
                    supersedes_fact_id=supersedes_fact_id,
                )
                changed = [
                    _copy_fact(record)
                    for record in records
                    if before.get(record.fact_id) != record.to_dict()
                ]
                overlay = {record.fact_id: record for record in self._read_overlay(branch_id)}
                for record in changed:
                    overlay[record.fact_id] = record
                self._write_overlay(branch_id, overlay.values())
                changed_ids = tuple(record.fact_id for record in changed)
                return EvolutionStateBranchResult(
                    ok=True,
                    status="updated",
                    branch_id=branch_id,
                    artifact_digest=str(metadata.get("artifact_digest") or ""),
                    changed_fact_ids=changed_ids,
                )
        except Exception as exc:
            return EvolutionStateBranchResult(ok=False, status="failed", branch_id=branch_id, error=str(exc))

    def deprecate_fact(
        self,
        branch_id: str,
        fact_id: str,
        *,
        reason: str | None = None,
        superseded_by: str | None = None,
    ) -> EvolutionStateBranchResult:
        try:
            with self._locked():
                metadata = self._load_active_branch(branch_id)
                overlay = {record.fact_id: record for record in self._read_overlay(branch_id)}
                if fact_id in overlay:
                    record = overlay[fact_id]
                    record.status = "deprecated"
                    record.updated_at = datetime.now().isoformat()
                    overlay[fact_id] = record
                    self._write_overlay(branch_id, overlay.values())
                else:
                    stable_ids = {record.fact_id for record in self.fact_store.read_all_unlocked()}
                    if fact_id not in stable_ids:
                        raise ValueError(f"unknown fact_id: {fact_id}")
                    tombstones = set(self._read_tombstones(branch_id))
                    tombstones.add(fact_id)
                    self._write_tombstones(branch_id, sorted(tombstones))
                return EvolutionStateBranchResult(
                    ok=True,
                    status="updated",
                    branch_id=branch_id,
                    artifact_digest=str(metadata.get("artifact_digest") or ""),
                    changed_fact_ids=(fact_id,),
                )
        except Exception as exc:
            return EvolutionStateBranchResult(ok=False, status="failed", branch_id=branch_id, error=str(exc))

    def preview_merge(self, branch_id: str) -> EvolutionMergePreview:
        try:
            with self._locked():
                metadata = self._load_active_branch(branch_id)
                overlay = self._read_overlay(branch_id)
                tombstones = tuple(self._read_tombstones(branch_id))
                stable = self.fact_store.read_all_unlocked()
                return self._preview_unlocked(branch_id, metadata, stable, overlay, tombstones)
        except Exception as exc:
            return EvolutionMergePreview(
                ok=False,
                branch_id=branch_id,
                conflicts=(
                    EvolutionMergeConflict(
                        fact_id="",
                        reason=str(exc),
                        stable_hash="",
                        base_hash="",
                    ),
                ),
            )

    def merge_branch(self, branch_id: str, *, actor: str = "user") -> EvolutionStateBranchResult:
        with self._locked():
            try:
                metadata = self._load_active_branch(branch_id)
                stable = self.fact_store.read_all_unlocked()
                overlay = self._read_overlay(branch_id)
                tombstones = tuple(self._read_tombstones(branch_id))
                preview = self._preview_unlocked(branch_id, metadata, stable, overlay, tombstones)
                artifact_digest = str(metadata.get("artifact_digest") or "")
                if not preview.ok:
                    event = self.ledger.append(
                        EvolutionEvent.new(
                            EventType.MODULE_FAILED,
                            actor=actor,
                            artifact_digest=artifact_digest,
                            state_branch_id=branch_id,
                            result={
                                "status": "merge_conflict",
                                "branch_id": branch_id,
                                "conflicts": [conflict.to_dict() for conflict in preview.conflicts],
                            },
                        )
                    )
                    return EvolutionStateBranchResult(
                        ok=False,
                        status="merge_conflict",
                        branch_id=branch_id,
                        artifact_digest=artifact_digest,
                        conflicts=preview.conflicts,
                        events=(event,),
                        error="merge conflict",
                    )
                merged = self._apply_merge(stable, overlay, tombstones, metadata)
                self.fact_store._write_records_unlocked(merged)
                _write_text_atomic(self.memory_dir / "MEMORY.md", render_memory_md(merged))
                metadata["status"] = "merged"
                metadata["merged_at"] = datetime.now(timezone.utc).isoformat()
                _write_json_atomic(self._branch_dir(branch_id) / "branch.json", metadata)
                changed_ids = tuple(sorted(set(preview.added + preview.modified + preview.deprecated)))
                event = self.ledger.append(
                    EvolutionEvent.new(
                        EventType.STATE_BRANCH_MERGED,
                        actor=actor,
                        artifact_digest=artifact_digest,
                        state_branch_id=branch_id,
                        result={
                            "status": "merged",
                            "branch_id": branch_id,
                            "changed_fact_ids": list(changed_ids),
                        },
                    )
                )
                return EvolutionStateBranchResult(
                    ok=True,
                    status="merged",
                    branch_id=branch_id,
                    artifact_digest=artifact_digest,
                    changed_fact_ids=changed_ids,
                    events=(event,),
                )
            except Exception as exc:
                return EvolutionStateBranchResult(ok=False, status="failed", branch_id=branch_id, error=str(exc))

    def discard_branch(self, branch_id: str, *, actor: str = "user") -> EvolutionStateBranchResult:
        with self._locked():
            try:
                metadata = self._load_active_branch(branch_id)
                metadata["status"] = "discarded"
                metadata["discarded_at"] = datetime.now(timezone.utc).isoformat()
                _write_json_atomic(self._branch_dir(branch_id) / "branch.json", metadata)
                artifact_digest = str(metadata.get("artifact_digest") or "")
                event = self.ledger.append(
                    EvolutionEvent.new(
                        EventType.STATE_BRANCH_DISCARDED,
                        actor=actor,
                        artifact_digest=artifact_digest,
                        state_branch_id=branch_id,
                        result={"status": "discarded", "branch_id": branch_id},
                    )
                )
                return EvolutionStateBranchResult(
                    ok=True,
                    status="discarded",
                    branch_id=branch_id,
                    artifact_digest=artifact_digest,
                    events=(event,),
                )
            except Exception as exc:
                return EvolutionStateBranchResult(ok=False, status="failed", branch_id=branch_id, error=str(exc))

    def _locked(self) -> FileLock:
        if self._lock_factory is not None:
            return self._lock_factory()
        self.memory_dir.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))

    def _branch_dir(self, branch_id: str) -> Path:
        return self.branch_root / branch_id

    def _load_branch(self, branch_id: str) -> dict[str, Any]:
        path = self._branch_dir(branch_id) / "branch.json"
        data = json.loads(path.read_text(encoding="utf-8"))
        if not isinstance(data, dict):
            raise ValueError("branch.json must be a mapping")
        if data.get("schema_version") != BRANCH_SCHEMA_VERSION:
            raise ValueError("unsupported branch schema_version")
        if data.get("branch_id") != branch_id:
            raise ValueError("branch_id mismatch")
        return data

    def _load_active_branch(self, branch_id: str) -> dict[str, Any]:
        metadata = self._load_branch(branch_id)
        if metadata.get("status") != "active":
            raise ValueError(f"branch is not active: {metadata.get('status')}")
        return metadata

    def _compose_facts_unlocked(
        self,
        branch_id: str,
        stable: list[FactRecord] | None = None,
    ) -> list[FactRecord]:
        stable_records = [_copy_fact(record) for record in (stable or self.fact_store.read_all_unlocked())]
        overlay = {record.fact_id: _copy_fact(record) for record in self._read_overlay(branch_id)}
        tombstones = set(self._read_tombstones(branch_id))
        composed: dict[str, FactRecord] = {
            record.fact_id: record
            for record in stable_records
            if record.fact_id not in tombstones
        }
        composed.update(overlay)
        return list(composed.values())

    def _read_overlay(self, branch_id: str) -> list[FactRecord]:
        return _read_fact_records(self._branch_dir(branch_id) / "facts_overlay.jsonl")

    def _write_overlay(self, branch_id: str, records: Any) -> None:
        _write_fact_records_atomic(self._branch_dir(branch_id) / "facts_overlay.jsonl", list(records))

    def _read_tombstones(self, branch_id: str) -> list[str]:
        path = self._branch_dir(branch_id) / "facts_tombstones.jsonl"
        tombstones: list[str] = []
        with suppress(FileNotFoundError):
            for line in path.read_text(encoding="utf-8").splitlines():
                if not line.strip():
                    continue
                data = json.loads(line)
                if not isinstance(data, dict) or not isinstance(data.get("fact_id"), str):
                    raise ValueError("invalid tombstone record")
                tombstones.append(data["fact_id"])
        return tombstones

    def _write_tombstones(self, branch_id: str, fact_ids: list[str]) -> None:
        _write_tombstones_atomic(self._branch_dir(branch_id) / "facts_tombstones.jsonl", fact_ids)

    def _preview_unlocked(
        self,
        branch_id: str,
        metadata: dict[str, Any],
        stable: list[FactRecord],
        overlay: list[FactRecord],
        tombstones: tuple[str, ...],
    ) -> EvolutionMergePreview:
        base_hashes = metadata.get("base_fact_hashes") or {}
        if not isinstance(base_hashes, dict):
            raise ValueError("base_fact_hashes must be a mapping")
        stable_by_id = {record.fact_id: record for record in stable}
        stable_hashes = _fact_hashes(stable)
        added: list[str] = []
        modified: list[str] = []
        deprecated: list[str] = []
        conflicts: list[EvolutionMergeConflict] = []

        for record in overlay:
            record_hash = _fact_hash(record)
            base_hash = str(base_hashes.get(record.fact_id) or "")
            stable_hash = str(stable_hashes.get(record.fact_id) or "")
            if not base_hash:
                if record.fact_id not in stable_by_id:
                    added.append(record.fact_id)
                elif stable_hash == record_hash:
                    continue
                else:
                    conflicts.append(
                        EvolutionMergeConflict(
                            fact_id=record.fact_id,
                            reason="stable_fact_id_collision",
                            stable_hash=stable_hash,
                            base_hash=base_hash,
                        )
                    )
                continue
            if stable_hash == base_hash:
                if record_hash != base_hash:
                    modified.append(record.fact_id)
                continue
            if stable_hash == record_hash:
                continue
            conflicts.append(
                EvolutionMergeConflict(
                    fact_id=record.fact_id,
                    reason="stable_fact_changed",
                    stable_hash=stable_hash,
                    base_hash=base_hash,
                )
            )

        for fact_id in tombstones:
            base_hash = str(base_hashes.get(fact_id) or "")
            stable = stable_by_id.get(fact_id)
            stable_hash = str(stable_hashes.get(fact_id) or "")
            if stable_hash == base_hash:
                deprecated.append(fact_id)
                continue
            if stable is not None and stable.status == "deprecated":
                continue
            conflicts.append(
                EvolutionMergeConflict(
                    fact_id=fact_id,
                    reason="stable_fact_changed",
                    stable_hash=stable_hash,
                    base_hash=base_hash,
                )
            )

        return EvolutionMergePreview(
            ok=not conflicts,
            branch_id=branch_id,
            added=tuple(sorted(set(added))),
            modified=tuple(sorted(set(modified))),
            deprecated=tuple(sorted(set(deprecated))),
            conflicts=tuple(conflicts),
        )

    def _apply_merge(
        self,
        stable: list[FactRecord],
        overlay: list[FactRecord],
        tombstones: tuple[str, ...],
        metadata: dict[str, Any],
    ) -> list[FactRecord]:
        base_hashes = metadata.get("base_fact_hashes") or {}
        merged = {record.fact_id: _copy_fact(record) for record in stable}
        for record in overlay:
            base_hash = str(base_hashes.get(record.fact_id) or "")
            current = merged.get(record.fact_id)
            if current is not None and not base_hash and _fact_hash(current) == _fact_hash(record):
                continue
            merged[record.fact_id] = _copy_fact(record)
        now = datetime.now().isoformat()
        for fact_id in tombstones:
            record = merged.get(fact_id)
            if record is not None:
                record.status = "deprecated"
                record.updated_at = now
        return list(merged.values())

    def _find_verified_event(self, artifact_digest: str) -> dict[str, str] | None:
        """Return the matching verified event by scanning the ledger in O(n)."""

        verification = self.ledger.verify_chain()
        if not verification.ok:
            return None
        if not self.ledger.event_path.exists():
            return None
        match: dict[str, str] | None = None
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
                        match = {
                            "event_hash": str(event.get("event_hash") or ""),
                            "terminal_event_hash": verification.terminal_event_hash or "",
                        }
        return match


def _copy_fact(record: FactRecord) -> FactRecord:
    return FactRecord.from_dict(record.to_dict())


def _fact_hash(record: FactRecord) -> str:
    return hashlib.sha256(canonical_dump(record.to_dict())).hexdigest()


def _fact_hashes(records: list[FactRecord]) -> dict[str, str]:
    return {record.fact_id: _fact_hash(record) for record in records}


def _facts_hash(records: list[FactRecord]) -> str:
    payload = {
        "facts": [
            record.to_dict()
            for record in sorted(records, key=lambda fact: fact.fact_id)
        ]
    }
    return hashlib.sha256(canonical_dump(payload)).hexdigest()


def _read_fact_records(path: Path) -> list[FactRecord]:
    records: list[FactRecord] = []
    with suppress(FileNotFoundError):
        for line in path.read_text(encoding="utf-8").splitlines():
            if not line.strip():
                continue
            data = json.loads(line)
            if not isinstance(data, dict):
                raise ValueError("fact record must be a mapping")
            records.append(FactRecord.from_dict(data))
    return records


def _write_fact_records_atomic(path: Path, records: list[FactRecord]) -> None:
    text = "".join(
        canonical_dump(record.to_dict()).decode("utf-8") + "\n"
        for record in records
    )
    _write_text_atomic(path, text)


def _write_tombstones_atomic(path: Path, fact_ids: list[str]) -> None:
    text = "".join(
        canonical_dump({"fact_id": fact_id}).decode("utf-8") + "\n"
        for fact_id in fact_ids
    )
    _write_text_atomic(path, text)


def _write_json_atomic(path: Path, data: dict[str, Any]) -> None:
    _write_text_atomic(path, canonical_dump(data).decode("utf-8") + "\n")


def _write_text_atomic(path: Path, text: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = path.with_name(f".{path.name}.tmp")
    try:
        with open(tmp_path, "w", encoding="utf-8") as handle:
            handle.write(text)
            handle.flush()
            os.fsync(handle.fileno())
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
