"""Lifecycle governance for workspace skills."""

from __future__ import annotations

import json
import re
import uuid
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

import yaml
from filelock import FileLock
from loguru import logger

from OriginAgent.agent.metadata import read_originagent_metadata, set_originagent_metadata
from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.utils.helpers import truncate_text

SKILL_LIFECYCLE_EVENTS_RELATIVE = Path("memory") / "skill_lifecycle_events.jsonl"

_FRONTMATTER_RE = re.compile(r"^---\s*\r?\n(.*?)\r?\n---\s*\r?\n?", re.DOTALL)
_SAFE_WORKSPACE_SKILL_RE = re.compile(r"^[a-z0-9][a-z0-9-]{0,63}$")
_TERMINAL_LIFECYCLE_STATUSES = {"deprecated", "rejected"}
_BODY_PREVIEW_MAX_CHARS = 900
_REASON_MAX_CHARS = 1000


@dataclass(frozen=True)
class SkillLifecycleEvent:
    """One append-only skill lifecycle decision event."""

    event_id: str
    skill_name: str
    action: str
    created_at: str
    reason: str = ""
    actor: str = "user"
    previous: dict[str, Any] | None = None
    next: dict[str, Any] | None = None
    artifact_path: str = ""
    review_proposal_id: str = ""

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class SkillLifecycleResult:
    """Outcome for a lifecycle operation."""

    skill_name: str
    status: str
    action: str
    ok: bool
    message: str
    skill: dict[str, Any] | None = None
    event: dict[str, Any] | None = None
    error: str = ""

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class SkillLifecycleStore:
    """Append-only lifecycle store plus deterministic frontmatter updates."""

    def __init__(self, workspace: Path, *, event_path: Path | None = None) -> None:
        self.workspace = Path(workspace)
        self.event_path = event_path or (self.workspace / SKILL_LIFECYCLE_EVENTS_RELATIVE)
        self._lock_path = self.event_path.parent / ".skill_lifecycle.lock"

    def _locked(self) -> FileLock:
        self.event_path.parent.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))

    def list_records(self, entries: list[dict[str, str]]) -> list[dict[str, Any]]:
        with self._locked():
            latest = self._latest_events_unlocked()
            return [self.decorate_entry(entry, latest_event=latest.get(entry.get("name", ""))) for entry in entries]

    def get_record(self, entry: dict[str, str] | None) -> dict[str, Any] | None:
        if entry is None:
            return None
        with self._locked():
            latest = self._latest_events_unlocked()
            return self.decorate_entry(entry, latest_event=latest.get(entry.get("name", "")))

    def stats(self, entries: list[dict[str, str]]) -> dict[str, Any]:
        records = self.list_records(entries)
        lifecycle_counts: dict[str, int] = {}
        verification_counts: dict[str, int] = {}
        always_workspace_count = 0
        workspace_count = 0
        for record in records:
            lifecycle = str(record.get("lifecycle_status") or "unknown")
            verification = str(record.get("verification_status") or "unknown")
            lifecycle_counts[lifecycle] = lifecycle_counts.get(lifecycle, 0) + 1
            verification_counts[verification] = verification_counts.get(verification, 0) + 1
            if record.get("source") == "workspace":
                workspace_count += 1
                if record.get("effective_always"):
                    always_workspace_count += 1
        return {
            "skills_count": len(records),
            "workspace_skills_count": workspace_count,
            "skill_lifecycle_status_counts": lifecycle_counts,
            "skill_verification_status_counts": verification_counts,
            "unverified_skill_count": verification_counts.get("unverified", 0),
            "deprecated_skill_count": lifecycle_counts.get("deprecated", 0),
            "rejected_skill_count": lifecycle_counts.get("rejected", 0),
            "always_workspace_skill_count": always_workspace_count,
        }

    def transition(
        self,
        skill_name: str,
        *,
        action: str,
        reason: str = "",
        actor: str = "user",
        enabled: bool | None = None,
    ) -> SkillLifecycleResult:
        skill_name = skill_name.strip()
        action = action.strip().lower()
        if not _is_safe_workspace_skill_name(skill_name):
            return SkillLifecycleResult(
                skill_name=skill_name,
                status="missing",
                action=action,
                ok=False,
                message="Skill name is invalid.",
                error="invalid_skill_name",
            )

        with self._locked():
            skill_file = self.workspace / "skills" / skill_name / "SKILL.md"
            if not skill_file.is_file():
                return SkillLifecycleResult(
                    skill_name=skill_name,
                    status="missing",
                    action=action,
                    ok=False,
                    message="Workspace skill was not found.",
                    error="not_found",
                )
            entry = {"name": skill_name, "path": str(skill_file), "source": "workspace"}
            latest = self._latest_events_unlocked().get(skill_name)
            current = self.decorate_entry(entry, latest_event=latest)
            if not current.get("can_manage_lifecycle"):
                return SkillLifecycleResult(
                    skill_name=skill_name,
                    status=str(current.get("lifecycle_status") or "unknown"),
                    action=action,
                    ok=False,
                    message=str(current.get("disabled_reason") or "Skill is not lifecycle-manageable."),
                    skill=current,
                    error="read_only",
                )

            target = self._target_state(current, action=action, enabled=enabled)
            if target.get("error"):
                return SkillLifecycleResult(
                    skill_name=skill_name,
                    status=str(current.get("lifecycle_status") or "unknown"),
                    action=action,
                    ok=False,
                    message=str(target["message"]),
                    skill=current,
                    error=str(target["error"]),
                )
            if target.get("idempotent"):
                return SkillLifecycleResult(
                    skill_name=skill_name,
                    status=str(current.get("lifecycle_status") or "unknown"),
                    action=action,
                    ok=True,
                    message=str(target["message"]),
                    skill=current,
                )

            frontmatter, body = _read_skill_markdown(skill_file)
            metadata = _metadata_originagent(frontmatter)
            next_state = dict(target["next"])
            now = datetime.now(timezone.utc).isoformat()
            event_id = f"skill_lifecycle_{uuid.uuid4().hex}"
            cleaned_actor = _clean(actor, 128)
            metadata.update({
                "proposal_status": _clean(metadata.get("proposal_status") or "unknown", 128),
                "verification_status": next_state["verification_status"],
                "lifecycle_status": next_state["lifecycle_status"],
                "version": _clean(metadata.get("version") or "1", 128),
                "supersedes_skill": _clean(metadata.get("supersedes_skill") or "", 128),
                "reviewed_by": cleaned_actor,
                "reviewed_at": now,
                "last_lifecycle_event_id": event_id,
            })
            if "always" in next_state:
                metadata["always"] = bool(next_state["always"])
                frontmatter["always"] = bool(next_state["always"])
            frontmatter.setdefault("name", skill_name)
            frontmatter.setdefault("description", skill_name)
            frontmatter["metadata"] = _set_metadata_originagent(frontmatter.get("metadata"), metadata)
            _write_skill_markdown(skill_file, frontmatter, body)

            updated = self.decorate_entry(entry)
            event = SkillLifecycleEvent(
                event_id=event_id,
                skill_name=skill_name,
                action=action if action != "always" else ("always_on" if next_state.get("always") else "always_off"),
                created_at=now,
                reason=_clean(reason, _REASON_MAX_CHARS),
                actor=cleaned_actor,
                previous=_state_snapshot(current),
                next=_state_snapshot(updated),
                artifact_path=f"skills/{skill_name}/SKILL.md",
                review_proposal_id=_clean(updated.get("review_proposal_id") or "", 256),
            ).to_json()
            self._append_event_unlocked(event)
            updated = self.decorate_entry(entry, latest_event=event)
            return SkillLifecycleResult(
                skill_name=skill_name,
                status=str(updated.get("lifecycle_status") or "unknown"),
                action=action,
                ok=True,
                message=str(target["message"]),
                skill=updated,
                event=event,
            )

    def _target_state(
        self,
        record: dict[str, Any],
        *,
        action: str,
        enabled: bool | None,
    ) -> dict[str, Any]:
        lifecycle = str(record.get("lifecycle_status") or "unknown")
        verification = str(record.get("verification_status") or "unknown")
        always = bool(record.get("always"))

        if action == "verify":
            if lifecycle in _TERMINAL_LIFECYCLE_STATUSES:
                return _transition_error("terminal_status", f"Skill is already {lifecycle}.")
            if lifecycle != "proposed":
                return _transition_error("invalid_transition", "Only proposed skills can be verified.")
            if verification == "verified":
                return _transition_idempotent("Skill is already verified.")
            return {
                "message": "Skill verified.",
                "next": {"lifecycle_status": "proposed", "verification_status": "verified"},
            }

        if action == "activate":
            if lifecycle == "active" and verification == "verified":
                return _transition_idempotent("Skill is already active.")
            if lifecycle in _TERMINAL_LIFECYCLE_STATUSES:
                return _transition_error("terminal_status", f"Skill is already {lifecycle}.")
            if lifecycle != "proposed" or verification != "verified":
                return _transition_error("invalid_transition", "Only verified proposed skills can be activated.")
            return {
                "message": "Skill activated.",
                "next": {"lifecycle_status": "active", "verification_status": "verified"},
            }

        if action == "reject":
            if lifecycle == "rejected":
                return _transition_idempotent("Skill is already rejected.")
            if lifecycle != "proposed":
                return _transition_error("invalid_transition", "Only proposed skills can be rejected.")
            return {
                "message": "Skill rejected.",
                "next": {"lifecycle_status": "rejected", "verification_status": "rejected"},
            }

        if action == "deprecate":
            if lifecycle == "deprecated":
                return _transition_idempotent("Skill is already deprecated.")
            if lifecycle == "rejected":
                return _transition_error("terminal_status", "Rejected skills cannot be deprecated.")
            return {
                "message": "Skill deprecated.",
                "next": {"lifecycle_status": "deprecated", "verification_status": verification},
            }

        if action == "always":
            next_enabled = bool(enabled)
            if next_enabled:
                if lifecycle != "active" or verification != "verified":
                    return _transition_error(
                        "invalid_transition",
                        "Only active verified workspace skills can be always loaded.",
                    )
                if always:
                    return _transition_idempotent("Skill is already always loaded.")
                return {
                    "message": "Skill always loading enabled.",
                    "next": {
                        "lifecycle_status": lifecycle,
                        "verification_status": verification,
                        "always": True,
                    },
                }
            if not always:
                return _transition_idempotent("Skill is already not always loaded.")
            return {
                "message": "Skill always loading disabled.",
                "next": {
                    "lifecycle_status": lifecycle,
                    "verification_status": verification,
                    "always": False,
                },
            }

        return _transition_error("unknown_action", "Unknown skill lifecycle action.")

    def decorate_entry(
        self,
        entry: dict[str, str],
        *,
        latest_event: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        name = str(entry.get("name") or "")
        source = str(entry.get("source") or "unknown")
        path = Path(str(entry.get("path") or ""))
        frontmatter, body = _read_skill_markdown(path)
        originagent = _metadata_originagent(frontmatter)
        lifecycle, verification = _derive_status(source, originagent)
        event_next = latest_event.get("next") if isinstance(latest_event, dict) else None
        if isinstance(event_next, dict):
            lifecycle = str(event_next.get("lifecycle_status") or lifecycle)
            verification = str(event_next.get("verification_status") or verification)
        always = bool(frontmatter.get("always") or originagent.get("always"))
        if isinstance(event_next, dict) and "always" in event_next:
            always = bool(event_next.get("always"))
        effective_always = _effective_always(source, lifecycle, verification, always)
        description = str(frontmatter.get("description") or name)
        record = {
            **entry,
            "description": _clean(description, 512),
            "body_preview": truncate_text(redact_memory_text(body.strip()), _BODY_PREVIEW_MAX_CHARS),
            "proposal_status": _clean(originagent.get("proposal_status") or "unknown", 128),
            "verification_status": verification,
            "lifecycle_status": lifecycle,
            "version": _clean(originagent.get("version") or "1", 128),
            "supersedes_skill": _clean(originagent.get("supersedes_skill") or "", 128),
            "reviewed_by": _clean(originagent.get("reviewed_by") or "", 128),
            "reviewed_at": _clean(originagent.get("reviewed_at") or "", 128),
            "review_proposal_id": _clean(originagent.get("review_proposal_id") or "", 256),
            "domain_id": _clean(originagent.get("domain_id") or "core", 128),
            "created_by": _clean(originagent.get("created_by") or "", 128),
            "last_lifecycle_event_id": _clean(originagent.get("last_lifecycle_event_id") or "", 128),
            "always": always,
            "effective_always": effective_always,
            "can_manage_lifecycle": source == "workspace",
            "disabled_reason": "" if source == "workspace" else "Only workspace skills can be changed in P9.",
            "last_event": _redact_json_value(latest_event) if isinstance(latest_event, dict) else None,
        }
        record.update(_action_capabilities(record))
        return record

    def _latest_events_unlocked(self) -> dict[str, dict[str, Any]]:
        latest: dict[str, dict[str, Any]] = {}
        for event in self._iter_events_unlocked():
            name = str(event.get("skill_name") or "")
            if name:
                latest[name] = event
        return latest

    def _iter_events_unlocked(self) -> list[dict[str, Any]]:
        rows: list[dict[str, Any]] = []
        try:
            with self.event_path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(data, dict):
                        rows.append(data)
        except FileNotFoundError:
            return []
        except OSError:
            logger.exception("Failed to read skill lifecycle events")
            return []
        return rows

    def _append_event_unlocked(self, event: dict[str, Any]) -> None:
        self.event_path.parent.mkdir(parents=True, exist_ok=True)
        with self.event_path.open("a", encoding="utf-8") as handle:
            handle.write(json.dumps(event, ensure_ascii=False) + "\n")


def summarize_skill_lifecycle(workspace: Path, entries: list[dict[str, str]]) -> dict[str, Any]:
    return SkillLifecycleStore(workspace).stats(entries)


def _derive_status(source: str, metadata: dict[str, Any]) -> tuple[str, str]:
    verification = str(metadata.get("verification_status") or "unknown")
    lifecycle = str(metadata.get("lifecycle_status") or "")
    if lifecycle:
        return lifecycle, verification
    if source != "workspace":
        return "active", verification
    if metadata.get("created_by") == "background_review" and verification == "unverified":
        return "proposed", verification
    return "active", verification


def _effective_always(source: str, lifecycle: str, verification: str, always: bool) -> bool:
    if not always:
        return False
    if source != "workspace":
        return True
    if lifecycle in {"proposed", "deprecated", "rejected"}:
        return False
    if verification == "unverified":
        return False
    return True


def _action_capabilities(record: dict[str, Any]) -> dict[str, Any]:
    if record.get("source") != "workspace":
        return {
            "can_verify": False,
            "can_activate": False,
            "can_deprecate": False,
            "can_reject": False,
            "can_toggle_always": False,
        }
    lifecycle = str(record.get("lifecycle_status") or "unknown")
    verification = str(record.get("verification_status") or "unknown")
    raw_always = bool(record.get("always"))
    return {
        "can_verify": lifecycle == "proposed" and verification != "verified",
        "can_activate": lifecycle == "proposed" and verification == "verified",
        "can_deprecate": lifecycle not in _TERMINAL_LIFECYCLE_STATUSES,
        "can_reject": lifecycle == "proposed",
        "can_toggle_always": raw_always or (lifecycle == "active" and verification == "verified"),
    }


def _state_snapshot(record: dict[str, Any]) -> dict[str, Any]:
    return {
        "lifecycle_status": record.get("lifecycle_status"),
        "verification_status": record.get("verification_status"),
        "always": bool(record.get("always")),
        "effective_always": bool(record.get("effective_always")),
    }


def _transition_error(error: str, message: str) -> dict[str, Any]:
    return {"error": error, "message": message}


def _transition_idempotent(message: str) -> dict[str, Any]:
    return {"idempotent": True, "message": message}


def _read_skill_markdown(path: Path) -> tuple[dict[str, Any], str]:
    try:
        content = path.read_text(encoding="utf-8")
    except OSError:
        return {}, ""
    match = _FRONTMATTER_RE.match(content)
    if not match:
        return {}, content
    try:
        parsed = yaml.safe_load(match.group(1))
    except yaml.YAMLError:
        return {}, content[match.end():].strip()
    frontmatter = parsed if isinstance(parsed, dict) else {}
    return {str(key): value for key, value in frontmatter.items()}, content[match.end():].strip()


def _write_skill_markdown(path: Path, frontmatter: dict[str, Any], body: str) -> None:
    text = "---\n" + yaml.safe_dump(frontmatter, allow_unicode=True, sort_keys=False) + "---\n\n" + body.rstrip() + "\n"
    path.write_text(text, encoding="utf-8")


def _metadata_originagent(frontmatter: dict[str, Any]) -> dict[str, Any]:
    return read_originagent_metadata(frontmatter.get("metadata"))


def _set_metadata_originagent(raw: Any, originagent: dict[str, Any]) -> dict[str, Any]:
    return set_originagent_metadata(raw, originagent)


def _is_safe_workspace_skill_name(name: str) -> bool:
    if not name or name in {".", ".."}:
        return False
    if "/" in name or "\\" in name or ".." in name or ":" in name:
        return False
    return bool(_SAFE_WORKSPACE_SKILL_RE.fullmatch(name))


def _clean(value: Any, max_chars: int) -> str:
    return truncate_text(redact_memory_text(str(value or "").strip()), max_chars)


def _redact_json_value(value: Any) -> Any:
    if isinstance(value, str):
        return _clean(value, 1000)
    if isinstance(value, list):
        return [_redact_json_value(item) for item in value]
    if isinstance(value, dict):
        return {str(key): _redact_json_value(item) for key, item in value.items()}
    return value
