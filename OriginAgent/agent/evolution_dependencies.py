"""Static dependency graph for governed evolution artifacts."""

from __future__ import annotations

import hashlib
import json
import os
import re
from contextlib import suppress
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

import yaml
from filelock import FileLock

from OriginAgent.agent.metadata import read_originagent_metadata
from OriginAgent.utils.helpers import ensure_dir, truncate_text

DEPENDENCY_SCHEMA_VERSION = "originagent.evolution.dependencies.v1"
DEPENDENCY_STORE_RELATIVE = Path("memory") / "evolution_dependencies.jsonl"

_ARTIFACT_TYPES = {"workflow", "skill"}
_NAME_RE = re.compile(r"^[a-z0-9][a-z0-9-]{0,63}$")
_EXPLICIT_REF_RE = re.compile(
    r"(?i)\b(?P<type>workflow|skill)\s*:\s*(?P<name>[a-z0-9][a-z0-9-]{0,63})\b"
)
_QUOTED_REF_RE = re.compile(
    r"(?i)\b(?P<type>workflow|skill)\s+[`\"'](?P<name>[a-z0-9][a-z0-9-]{0,63})[`\"']"
)
_PATH_REF_RE = re.compile(
    r"(?i)\b(?P<root>workflows|skills)/(?P<name>[a-z0-9][a-z0-9-]{0,63})/"
    r"(?P<file>workflow\.yaml|SKILL\.md)\b"
)
_EVIDENCE_CHARS = 160


@dataclass(frozen=True)
class DependencyReference:
    """One artifact reference discovered in static text."""

    type: str
    name: str
    source: str = "static"
    evidence: str = ""

    def to_json(self) -> dict[str, str]:
        return asdict(self)


@dataclass(frozen=True)
class DependencyRecord:
    """Current dependency state for one workflow or skill artifact."""

    schema_version: str
    artifact_type: str
    artifact_name: str
    artifact_path: str
    depends_on: list[dict[str, str]]
    referenced_by: list[dict[str, str]]
    updated_at: str
    version: str

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class EvolutionDependencyStore:
    """Maintain a compact JSONL dependency graph for evolution artifacts."""

    def __init__(self, workspace: Path, *, path: Path | None = None) -> None:
        self.workspace = Path(workspace)
        memory_dir = ensure_dir(self.workspace / "memory")
        self.path = path or (self.workspace / DEPENDENCY_STORE_RELATIVE)
        self._lock = FileLock(str(memory_dir / ".evolution_dependencies.lock"))

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
        return self._with_reverse_links(records)

    def update_artifact_from_event(self, artifact: dict[str, Any]) -> dict[str, Any] | None:
        artifact_type = _artifact_type(artifact)
        artifact_name = _artifact_name(artifact)
        artifact_path = _artifact_path(artifact)
        if artifact_type not in _ARTIFACT_TYPES or not artifact_name:
            return None
        return self.update_artifact(
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            artifact_path=artifact_path,
        )

    def update_artifact(
        self,
        *,
        artifact_type: str,
        artifact_name: str,
        artifact_path: str = "",
    ) -> dict[str, Any] | None:
        artifact_type = _clean_type(artifact_type)
        artifact_name = _clean_name(artifact_name)
        if artifact_type not in _ARTIFACT_TYPES or not artifact_name:
            return None
        path = self.workspace / (artifact_path or _default_artifact_path(artifact_type, artifact_name))
        try:
            content = path.read_text(encoding="utf-8")
        except OSError:
            return None
        relative_path = _relative_to_workspace(path, self.workspace)
        depends_on = [
            ref.to_json()
            for ref in extract_static_dependencies(
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                content=content,
            )
        ]
        record = DependencyRecord(
            schema_version=DEPENDENCY_SCHEMA_VERSION,
            artifact_type=artifact_type,
            artifact_name=artifact_name,
            artifact_path=relative_path,
            depends_on=depends_on,
            referenced_by=[],
            updated_at=datetime.now(timezone.utc).isoformat(),
            version=_content_hash(content),
        ).to_json()
        with self._lock:
            records = [
                existing
                for existing in self.read_all()
                if not _same_artifact(existing, artifact_type, artifact_name)
            ]
            records.append(record)
            records = self._with_reverse_links(records)
            self._write_all_unlocked(records)
            for existing in records:
                if _same_artifact(existing, artifact_type, artifact_name):
                    return existing
        return record

    def remove_artifact(self, *, artifact_type: str, artifact_name: str) -> None:
        artifact_type = _clean_type(artifact_type)
        artifact_name = _clean_name(artifact_name)
        if artifact_type not in _ARTIFACT_TYPES or not artifact_name:
            return
        with self._lock:
            records = [
                record
                for record in self.read_all()
                if not _same_artifact(record, artifact_type, artifact_name)
            ]
            self._write_all_unlocked(self._with_reverse_links(records))

    def referenced_by(self, *, artifact_type: str, artifact_name: str) -> list[dict[str, str]]:
        artifact_type = _clean_type(artifact_type)
        artifact_name = _clean_name(artifact_name)
        refs: dict[tuple[str, str], dict[str, str]] = {}
        for record in self.read_all():
            source_type = _clean_type(record.get("artifact_type"))
            source_name = _clean_name(record.get("artifact_name"))
            if source_type not in _ARTIFACT_TYPES or not source_name:
                continue
            for dep in _references(record.get("depends_on")):
                if dep.type == artifact_type and dep.name == artifact_name:
                    refs[(source_type, source_name)] = {
                        "type": source_type,
                        "name": source_name,
                        "source": dep.source,
                        "evidence": dep.evidence,
                    }
        return sorted(refs.values(), key=lambda item: (item["type"], item["name"]))

    def rollback_blockers(self, *, artifact_type: str, artifact_name: str) -> list[dict[str, str]]:
        return self.referenced_by(artifact_type=artifact_type, artifact_name=artifact_name)

    def prune_stale_references(self) -> dict[str, Any]:
        """Drop dependency records and edges that point at missing or retired artifacts."""

        records = self.read_all()
        retained: list[dict[str, Any]] = []
        removed_records = 0
        removed_edges = 0
        for record in records:
            artifact_type = _clean_type(record.get("artifact_type"))
            artifact_name = _clean_name(record.get("artifact_name"))
            artifact_path = _normalize_path(record.get("artifact_path"))
            if (
                artifact_type not in _ARTIFACT_TYPES
                or not artifact_name
                or _artifact_is_stale(self.workspace, artifact_type, artifact_name, artifact_path)
            ):
                removed_records += 1
                continue
            next_record = dict(record)
            refs = []
            for dep in _references(record.get("depends_on")):
                if _artifact_is_stale(self.workspace, dep.type, dep.name, ""):
                    removed_edges += 1
                    continue
                refs.append(dep.to_json())
            next_record["depends_on"] = refs
            retained.append(next_record)
        retained = self._with_reverse_links(retained)
        if removed_records or removed_edges:
            with self._lock:
                self._write_all_unlocked(retained)
        return {
            "removed_artifacts": removed_records,
            "removed_edges": removed_edges,
            "remaining_artifacts": len(retained),
            "stale_reference_count": self.stats()["stale_reference_count"],
        }

    def stats(self) -> dict[str, Any]:
        records = self.read_all()
        dependency_edges = sum(len(_references(record.get("depends_on"))) for record in records)
        blocked = sum(1 for record in records if _references(record.get("referenced_by")))
        stale = 0
        for record in records:
            if not (self.workspace / str(record.get("artifact_path") or "")).is_file():
                stale += 1
            for dep in _references(record.get("depends_on")):
                if not (self.workspace / _default_artifact_path(dep.type, dep.name)).is_file():
                    stale += 1
        return {
            "tracked_artifacts": len(records),
            "dependency_edges": dependency_edges,
            "rollback_blocked_artifacts": blocked,
            "stale_reference_count": stale,
        }

    def _write_all_unlocked(self, records: list[dict[str, Any]]) -> None:
        tmp_path = self.path.with_suffix(self.path.suffix + ".tmp")
        self.path.parent.mkdir(parents=True, exist_ok=True)
        try:
            with tmp_path.open("w", encoding="utf-8") as handle:
                for record in sorted(records, key=lambda item: (
                    str(item.get("artifact_type") or ""),
                    str(item.get("artifact_name") or ""),
                )):
                    handle.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")
                handle.flush()
                os.fsync(handle.fileno())
            os.replace(tmp_path, self.path)
        except BaseException:
            tmp_path.unlink(missing_ok=True)
            raise

    @staticmethod
    def _with_reverse_links(records: list[dict[str, Any]]) -> list[dict[str, Any]]:
        cleaned: list[dict[str, Any]] = []
        for record in records:
            if not isinstance(record, dict):
                continue
            artifact_type = _clean_type(record.get("artifact_type"))
            artifact_name = _clean_name(record.get("artifact_name"))
            if artifact_type not in _ARTIFACT_TYPES or not artifact_name:
                continue
            next_record = dict(record)
            next_record["schema_version"] = str(next_record.get("schema_version") or DEPENDENCY_SCHEMA_VERSION)
            next_record["artifact_type"] = artifact_type
            next_record["artifact_name"] = artifact_name
            next_record["artifact_path"] = _normalize_path(next_record.get("artifact_path"))
            next_record["depends_on"] = [ref.to_json() for ref in _references(next_record.get("depends_on"))]
            next_record["referenced_by"] = []
            cleaned.append(next_record)

        reverse: dict[tuple[str, str], dict[tuple[str, str], dict[str, str]]] = {}
        for record in cleaned:
            source_type = str(record.get("artifact_type") or "")
            source_name = str(record.get("artifact_name") or "")
            for dep in _references(record.get("depends_on")):
                reverse.setdefault((dep.type, dep.name), {})[(source_type, source_name)] = {
                    "type": source_type,
                    "name": source_name,
                    "source": dep.source,
                    "evidence": dep.evidence,
                }

        for record in cleaned:
            key = (str(record.get("artifact_type") or ""), str(record.get("artifact_name") or ""))
            record["referenced_by"] = sorted(
                reverse.get(key, {}).values(),
                key=lambda item: (item["type"], item["name"]),
            )
        return cleaned


def extract_static_dependencies(
    *,
    artifact_type: str,
    artifact_name: str,
    content: str,
) -> list[DependencyReference]:
    """Extract explicit workflow/skill references from artifact text."""

    texts = _dependency_texts(artifact_type, content)
    refs: dict[tuple[str, str], DependencyReference] = {}
    for text in texts:
        for pattern in (_EXPLICIT_REF_RE, _QUOTED_REF_RE):
            for match in pattern.finditer(text):
                ref = _reference_from_match(match, text)
                if ref is not None:
                    refs.setdefault((ref.type, ref.name), ref)
        for match in _PATH_REF_RE.finditer(text):
            root = match.group("root").casefold()
            ref_type = "workflow" if root == "workflows" else "skill"
            ref = _make_reference(ref_type, match.group("name"), match, text)
            if ref is not None:
                refs.setdefault((ref.type, ref.name), ref)

    own_type = _clean_type(artifact_type)
    own_name = _clean_name(artifact_name)
    return [
        ref
        for key, ref in sorted(refs.items())
        if key != (own_type, own_name)
    ]


def _dependency_texts(artifact_type: str, content: str) -> list[str]:
    if _clean_type(artifact_type) != "workflow":
        return [content]
    try:
        data = yaml.safe_load(content)
    except yaml.YAMLError:
        return [content]
    if not isinstance(data, dict):
        return [content]
    texts: list[str] = []
    for key in ("description", "body"):
        value = data.get(key)
        if isinstance(value, str):
            texts.append(value)
    steps = data.get("steps")
    if isinstance(steps, list):
        for step in steps:
            if not isinstance(step, dict):
                continue
            for key in ("title", "instruction", "risk"):
                value = step.get(key)
                if isinstance(value, str):
                    texts.append(value)
    return texts or [content]


def _reference_from_match(match: re.Match[str], text: str) -> DependencyReference | None:
    return _make_reference(match.group("type"), match.group("name"), match, text)


def _make_reference(
    artifact_type: Any,
    artifact_name: Any,
    match: re.Match[str],
    text: str,
) -> DependencyReference | None:
    ref_type = _clean_type(artifact_type)
    ref_name = _clean_name(artifact_name)
    if ref_type not in _ARTIFACT_TYPES or not ref_name:
        return None
    return DependencyReference(
        type=ref_type,
        name=ref_name,
        source="static",
        evidence=_evidence(text, match.start(), match.end()),
    )


def _references(value: Any) -> list[DependencyReference]:
    if not isinstance(value, list):
        return []
    refs: list[DependencyReference] = []
    seen: set[tuple[str, str]] = set()
    for item in value:
        if not isinstance(item, dict):
            continue
        ref_type = _clean_type(item.get("type"))
        ref_name = _clean_name(item.get("name"))
        if ref_type not in _ARTIFACT_TYPES or not ref_name or (ref_type, ref_name) in seen:
            continue
        seen.add((ref_type, ref_name))
        refs.append(DependencyReference(
            type=ref_type,
            name=ref_name,
            source=truncate_text(str(item.get("source") or "static").strip(), 64),
            evidence=truncate_text(str(item.get("evidence") or "").strip(), _EVIDENCE_CHARS),
        ))
    return refs


def _artifact_type(artifact: dict[str, Any]) -> str:
    raw = _clean_type(artifact.get("artifact_type"))
    if raw:
        return raw
    if artifact.get("workflow_name"):
        return "workflow"
    if artifact.get("skill_name"):
        return "skill"
    return ""


def _artifact_name(artifact: dict[str, Any]) -> str:
    artifact_type = _artifact_type(artifact)
    if artifact_type == "workflow":
        return _clean_name(artifact.get("workflow_name"))
    if artifact_type == "skill":
        return _clean_name(artifact.get("skill_name"))
    return ""


def _artifact_path(artifact: dict[str, Any]) -> str:
    return _normalize_path(artifact.get("path"))


def _same_artifact(record: dict[str, Any], artifact_type: str, artifact_name: str) -> bool:
    return (
        _clean_type(record.get("artifact_type")) == artifact_type
        and _clean_name(record.get("artifact_name")) == artifact_name
    )


def _default_artifact_path(artifact_type: str, artifact_name: str) -> str:
    if artifact_type == "workflow":
        return f"workflows/{artifact_name}/workflow.yaml"
    return f"skills/{artifact_name}/SKILL.md"


def _artifact_is_stale(workspace: Path, artifact_type: str, artifact_name: str, artifact_path: str) -> bool:
    path = workspace / (artifact_path or _default_artifact_path(artifact_type, artifact_name))
    if not path.is_file():
        return True
    try:
        content = path.read_text(encoding="utf-8")
    except OSError:
        return True
    if artifact_type == "skill":
        return _skill_lifecycle_status(content) in {"deprecated", "rejected"}
    if artifact_type == "workflow":
        return _workflow_proposal_status(content) in {"archived", "deprecated", "rejected"}
    return False


def _skill_lifecycle_status(content: str) -> str:
    if not content.startswith("---"):
        return ""
    parts = content.split("---", 2)
    if len(parts) < 3:
        return ""
    try:
        frontmatter = yaml.safe_load(parts[1])
    except yaml.YAMLError:
        return ""
    if not isinstance(frontmatter, dict):
        return ""
    metadata = frontmatter.get("metadata")
    originagent = read_originagent_metadata(metadata)
    return str(originagent.get("lifecycle_status") or "").strip().casefold()


def _workflow_proposal_status(content: str) -> str:
    try:
        data = yaml.safe_load(content)
    except yaml.YAMLError:
        return ""
    if not isinstance(data, dict):
        return ""
    metadata = data.get("metadata")
    originagent = read_originagent_metadata(metadata)
    return str(originagent.get("proposal_status") or "").strip().casefold()


def _evidence(text: str, start: int, end: int) -> str:
    half = _EVIDENCE_CHARS // 2
    left = max(0, start - half)
    right = min(len(text), end + half)
    return truncate_text(" ".join(text[left:right].split()), _EVIDENCE_CHARS)


def _relative_to_workspace(path: Path, workspace: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return path.as_posix()


def _content_hash(content: str) -> str:
    return hashlib.sha256(content.encode("utf-8")).hexdigest()


def _clean_type(value: Any) -> str:
    return str(value or "").strip().casefold()


def _clean_name(value: Any) -> str:
    candidate = str(value or "").strip().casefold()
    return candidate if _NAME_RE.fullmatch(candidate) else ""


def _normalize_path(value: Any) -> str:
    return str(value or "").strip().replace("\\", "/")
