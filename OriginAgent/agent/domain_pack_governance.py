"""Local governance for workspace and builtin domain packs."""

from __future__ import annotations

import importlib.util
import json
import shutil
import sys
import uuid
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

import yaml
from filelock import FileLock

from OriginAgent.agent.domain_packs import (
    DomainEvalDeclaration,
    DomainFileDeclaration,
    DomainPack,
    DomainPackManager,
    DomainPackRuntimeConfig,
    DomainPackValidator,
    DomainSkillDeclaration,
    DomainToolDeclaration,
    DomainWorkflowDeclaration,
)
from OriginAgent.agent.metadata import read_originagent_metadata, set_originagent_metadata
from OriginAgent.agent.skill_lifecycle import _metadata_originagent, _read_skill_markdown, _set_metadata_originagent, _write_skill_markdown
from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.workflow_artifacts import validate_workflow_artifact_dir
from OriginAgent.utils.helpers import truncate_text

DOMAIN_PACK_EVENTS_RELATIVE = Path("memory") / "domain_pack_events.jsonl"
_REASON_MAX_CHARS = 1000
_SUPPORTED_EVAL_KINDS = frozenset({"manifest", "skill", "tool", "workflow"})
_MANAGED_ARTIFACT_TYPES = frozenset({"skill", "workflow"})


@dataclass(frozen=True)
class DomainPackEvent:
    event_id: str
    pack_id: str
    action: str
    created_at: str
    reason: str = ""
    actor: str = "user"
    source: str = ""
    previous: dict[str, Any] | None = None
    next: dict[str, Any] | None = None
    result: dict[str, Any] | None = None
    overrides_builtin: bool = False
    validation_summary: str = ""
    eval_summary: dict[str, Any] | None = None
    artifact_paths: list[str] | None = None
    review_proposal_id: str = ""

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class DomainPackGovernanceResult:
    pack_id: str
    status: str
    action: str
    ok: bool
    message: str
    pack: dict[str, Any] | None = None
    event: dict[str, Any] | None = None
    eval_result: dict[str, Any] | None = None
    artifact: dict[str, Any] | None = None
    error: str = ""

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class DomainPackGovernanceService:
    """Govern local domain packs and append audit events."""

    def __init__(
        self,
        workspace: Path,
        *,
        domain_pack_manager: DomainPackManager | None = None,
        config_loader: Callable[[], Any] | None = None,
        config_saver: Callable[[Any], None] | None = None,
        event_path: Path | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self._domain_pack_manager = domain_pack_manager
        self._config_loader = config_loader
        self._config_saver = config_saver
        self.event_path = event_path or (self.workspace / DOMAIN_PACK_EVENTS_RELATIVE)
        self._lock_path = self.event_path.parent / ".domain_pack_governance.lock"

    def _locked(self) -> FileLock:
        self.event_path.parent.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))

    def _load_config(self) -> Any:
        if self._config_loader is not None:
            return self._config_loader()
        from OriginAgent.config.loader import load_config

        return load_config()

    def _save_config(self, config: Any) -> None:
        if self._config_saver is not None:
            self._config_saver(config)
            return
        from OriginAgent.config.loader import save_config

        save_config(config)

    def _manager(self, config: Any | None = None) -> DomainPackManager:
        if self._domain_pack_manager is not None:
            if config is not None:
                self._domain_pack_manager.config = DomainPackRuntimeConfig.from_config(
                    config.agents.defaults.domain_packs
                )
            self._domain_pack_manager.refresh()
            return self._domain_pack_manager
        if config is None:
            config = self._load_config()
        return DomainPackManager(
            self.workspace,
            config=config.agents.defaults.domain_packs,
        )

    def _read_manager(self) -> DomainPackManager:
        if self._domain_pack_manager is not None and self._config_loader is None:
            return self._domain_pack_manager
        config = self._load_config()
        return self._manager(config)

    def _refresh_manager(self, config: Any) -> DomainPackManager:
        return self._manager(config)

    def _read_events_unlocked(self) -> list[dict[str, Any]]:
        rows: list[dict[str, Any]] = []
        try:
            with self.event_path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        raw = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(raw, dict):
                        rows.append(raw)
        except FileNotFoundError:
            return []
        except OSError:
            return []
        return rows

    def _append_event_unlocked(
        self,
        *,
        pack_id: str,
        action: str,
        reason: str,
        actor: str,
        source: str,
        previous: dict[str, Any] | None,
        next_state: dict[str, Any] | None,
        result: dict[str, Any] | None,
        overrides_builtin: bool,
        validation_summary: str,
        eval_summary: dict[str, Any] | None,
        artifact_paths: list[str] | None,
        review_proposal_id: str,
    ) -> dict[str, Any]:
        event = DomainPackEvent(
            event_id=f"domain_pack_event_{uuid.uuid4().hex}",
            pack_id=pack_id,
            action=action,
            created_at=datetime.now(timezone.utc).isoformat(),
            reason=_clean_text(reason, _REASON_MAX_CHARS),
            actor=_clean_text(actor, 128) or "user",
            source=source,
            previous=previous,
            next=next_state,
            result=result,
            overrides_builtin=overrides_builtin,
            validation_summary=_clean_text(validation_summary, 1000),
            eval_summary=eval_summary,
            artifact_paths=artifact_paths or [],
            review_proposal_id=_clean_text(review_proposal_id, 256),
        ).to_json()
        self.event_path.parent.mkdir(parents=True, exist_ok=True)
        with self.event_path.open("a", encoding="utf-8") as handle:
            handle.write(json.dumps(event, ensure_ascii=False) + "\n")
        return event

    def list_records(
        self,
        *,
        source: str | None = None,
        status: str | None = None,
        limit: int = 50,
    ) -> list[dict[str, Any]]:
        limit = max(1, min(int(limit or 50), 200))
        manager = self._read_manager()
        packs = manager.list_packs()
        with self._locked():
            latest = self._latest_events_by_pack_unlocked()
            evals = self._latest_eval_events_by_pack_unlocked()
            records = [self._decorate_pack(pack, latest.get(pack.id), evals.get(pack.id)) for pack in packs]
        if source:
            source = source.strip().lower()
            records = [record for record in records if str(record.get("source") or "") == source]
        if status:
            status = status.strip().lower()
            records = [record for record in records if str(record.get("status") or "") == status]
        return records[:limit]

    def get_record(self, pack_id: str) -> dict[str, Any] | None:
        manager = self._read_manager()
        pack = manager.get_pack(pack_id)
        if pack is None:
            return None
        with self._locked():
            latest = self._latest_events_by_pack_unlocked()
            evals = self._latest_eval_events_by_pack_unlocked()
            return self._decorate_pack(pack, latest.get(pack.id), evals.get(pack.id))

    def stats(self) -> dict[str, Any]:
        records = self.list_records(limit=500)
        status_counts: dict[str, int] = {}
        eval_status_counts: dict[str, int] = {}
        active_count = 0
        override_count = 0
        builtin_count = 0
        workspace_count = 0
        last_event_at = None
        with self._locked():
            for event in self._read_events_unlocked():
                created_at = str(event.get("created_at") or "")
                if created_at and (last_event_at is None or created_at > last_event_at):
                    last_event_at = created_at
        for record in records:
            state = str(record.get("status") or "unknown")
            status_counts[state] = status_counts.get(state, 0) + 1
            eval_status = str((record.get("last_eval_result") or {}).get("status") or "")
            if eval_status:
                eval_status_counts[eval_status] = eval_status_counts.get(eval_status, 0) + 1
            if record.get("active"):
                active_count += 1
            if record.get("overrides_builtin"):
                override_count += 1
            if record.get("source") == "builtin":
                builtin_count += 1
            if record.get("source") == "workspace":
                workspace_count += 1
        return {
            "workspace_domain_pack_count": workspace_count,
            "builtin_domain_pack_count": builtin_count,
            "domain_pack_status_counts": status_counts,
            "active_domain_pack_count": active_count,
            "domain_pack_override_count": override_count,
            "domain_pack_eval_status_counts": eval_status_counts,
            "last_domain_pack_event_at": last_event_at,
        }

    def install(
        self,
        source_path: str,
        *,
        reason: str = "",
        actor: str = "user",
    ) -> DomainPackGovernanceResult:
        with self._locked():
            config = self._load_config()
            resolved_source = _resolve_local_dir(source_path)
            if resolved_source is None:
                return DomainPackGovernanceResult(
                    pack_id="",
                    status="failed",
                    action="install",
                    ok=False,
                    message="Source path is invalid.",
                    error="invalid_source",
                )
            validated = DomainPackValidator(
                runtime_config=DomainPackRuntimeConfig(),
                strict_declarations=True,
            ).validate_pack(
                resolved_source,
                source="workspace",
            )
            if validated.status == "invalid":
                return DomainPackGovernanceResult(
                    pack_id=validated.id,
                    status="failed",
                    action="install",
                    ok=False,
                    message="Domain pack source failed validation.",
                    error=validated.validation_summary or validated.unavailable_reason or "invalid_domain_pack",
                )
            target_dir = self.workspace / "domain_packs" / validated.id
            if target_dir.exists():
                return DomainPackGovernanceResult(
                    pack_id=validated.id,
                    status="failed",
                    action="install",
                    ok=False,
                    message="Workspace domain pack already exists.",
                    error="pack_exists",
                )

            previous = None
            staging_dir = target_dir.parent / f".{validated.id}.install.{uuid.uuid4().hex}"
            try:
                self._copy_pack_to_staging(resolved_source, staging_dir)
                installed = DomainPackValidator(
                    runtime_config=DomainPackRuntimeConfig(),
                    strict_declarations=True,
                ).validate_pack(
                    staging_dir,
                    source="workspace",
                )
                if installed.status == "invalid":
                    raise ValueError(installed.validation_summary or installed.unavailable_reason)
                target_dir.parent.mkdir(parents=True, exist_ok=True)
                staging_dir.rename(target_dir)
                defaults = config.agents.defaults.domain_packs
                defaults.disabled = [item for item in defaults.disabled if item != installed.id]
                defaults.active = [item for item in defaults.active if item != installed.id]
                self._save_config(config)
                manager = self._refresh_manager(config)
                latest = self._latest_events_by_pack_unlocked()
                evals = self._latest_eval_events_by_pack_unlocked()
                record = self._decorate_pack(manager.get_pack(installed.id), latest.get(installed.id), evals.get(installed.id))
                event = self._append_event_unlocked(
                    pack_id=installed.id,
                    action="install",
                    reason=reason,
                    actor=actor,
                    source=str(resolved_source),
                    previous=previous,
                    next_state=_record_snapshot(record),
                    result={"status": "installed"},
                    overrides_builtin=bool(record and record.get("overrides_builtin")),
                    validation_summary=str(record.get("validation_summary") if record else installed.validation_summary),
                    eval_summary=None,
                    artifact_paths=[f"domain_packs/{installed.id}"],
                    review_proposal_id="",
                )
                return DomainPackGovernanceResult(
                    pack_id=installed.id,
                    status="installed",
                    action="install",
                    ok=True,
                    message="Domain pack installed.",
                    pack=record,
                    event=event,
                )
            except Exception as exc:
                _cleanup_dir(staging_dir)
                return DomainPackGovernanceResult(
                    pack_id=validated.id,
                    status="failed",
                    action="install",
                    ok=False,
                    message="Failed to install domain pack.",
                    error=str(exc),
                )

    def upgrade(
        self,
        pack_id: str,
        source_path: str,
        *,
        reason: str = "",
        actor: str = "user",
    ) -> DomainPackGovernanceResult:
        with self._locked():
            config = self._load_config()
            manager = self._refresh_manager(config)
            current = manager.get_pack(pack_id)
            if current is None:
                return _pack_not_found_result(pack_id, "upgrade")
            if current.source != "workspace":
                return _read_only_result(pack_id, "upgrade", current)

            resolved_source = _resolve_local_dir(source_path)
            if resolved_source is None:
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="failed",
                    action="upgrade",
                    ok=False,
                    message="Source path is invalid.",
                    error="invalid_source",
                )
            validated = DomainPackValidator(
                runtime_config=DomainPackRuntimeConfig(),
                strict_declarations=True,
            ).validate_pack(
                resolved_source,
                source="workspace",
            )
            if validated.id != pack_id:
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="failed",
                    action="upgrade",
                    ok=False,
                    message="Source pack id does not match target pack id.",
                    error="pack_id_mismatch",
                )
            if validated.status == "invalid":
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="failed",
                    action="upgrade",
                    ok=False,
                    message="Domain pack source failed validation.",
                    error=validated.validation_summary or validated.unavailable_reason or "invalid_domain_pack",
                )

            previous = self._decorate_pack(current, None, None)
            target_dir = current.path
            staging_dir = target_dir.parent / f".{pack_id}.upgrade.{uuid.uuid4().hex}"
            backup_dir = target_dir.parent / f".{pack_id}.backup.{uuid.uuid4().hex}"
            try:
                self._copy_pack_to_staging(resolved_source, staging_dir)
                upgraded = DomainPackValidator(
                    runtime_config=DomainPackRuntimeConfig(),
                    strict_declarations=True,
                ).validate_pack(
                    staging_dir,
                    source="workspace",
                )
                if upgraded.status == "invalid":
                    raise ValueError(upgraded.validation_summary or upgraded.unavailable_reason)
                target_dir.rename(backup_dir)
                try:
                    staging_dir.rename(target_dir)
                except Exception:
                    if not target_dir.exists() and backup_dir.exists():
                        backup_dir.rename(target_dir)
                    raise
                _cleanup_dir(backup_dir)
                manager = self._refresh_manager(config)
                latest = self._latest_events_by_pack_unlocked()
                evals = self._latest_eval_events_by_pack_unlocked()
                record = self._decorate_pack(manager.get_pack(pack_id), latest.get(pack_id), evals.get(pack_id))
                event = self._append_event_unlocked(
                    pack_id=pack_id,
                    action="upgrade",
                    reason=reason,
                    actor=actor,
                    source=str(resolved_source),
                    previous=_record_snapshot(previous),
                    next_state=_record_snapshot(record),
                    result={"status": "upgraded"},
                    overrides_builtin=bool(record and record.get("overrides_builtin")),
                    validation_summary=str(record.get("validation_summary") if record else upgraded.validation_summary),
                    eval_summary=None,
                    artifact_paths=[f"domain_packs/{pack_id}"],
                    review_proposal_id="",
                )
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="upgraded",
                    action="upgrade",
                    ok=True,
                    message="Domain pack upgraded.",
                    pack=record,
                    event=event,
                )
            except Exception as exc:
                _cleanup_dir(staging_dir)
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="failed",
                    action="upgrade",
                    ok=False,
                    message="Failed to upgrade domain pack.",
                    error=str(exc),
                )

    def uninstall(
        self,
        pack_id: str,
        *,
        reason: str = "",
        actor: str = "user",
    ) -> DomainPackGovernanceResult:
        with self._locked():
            config = self._load_config()
            manager = self._refresh_manager(config)
            current = manager.get_pack(pack_id)
            if current is None:
                return _pack_not_found_result(pack_id, "uninstall")
            if current.source != "workspace":
                return _read_only_result(pack_id, "uninstall", current)
            previous = self._decorate_pack(current, None, None)
            trash_dir = current.path.parent / f".{pack_id}.delete.{uuid.uuid4().hex}"
            try:
                current.path.rename(trash_dir)
                defaults = config.agents.defaults.domain_packs
                defaults.active = [item for item in defaults.active if item != pack_id]
                defaults.disabled = [item for item in defaults.disabled if item != pack_id]
                self._save_config(config)
                _cleanup_dir(trash_dir)
                event = self._append_event_unlocked(
                    pack_id=pack_id,
                    action="uninstall",
                    reason=reason,
                    actor=actor,
                    source=str(current.path),
                    previous=_record_snapshot(previous),
                    next_state=None,
                    result={"status": "uninstalled"},
                    overrides_builtin=bool(previous and previous.get("overrides_builtin")),
                    validation_summary=str(previous.get("validation_summary") if previous else ""),
                    eval_summary=None,
                    artifact_paths=[f"domain_packs/{pack_id}"],
                    review_proposal_id="",
                )
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="uninstalled",
                    action="uninstall",
                    ok=True,
                    message="Domain pack uninstalled.",
                    pack=None,
                    event=event,
                )
            except Exception as exc:
                if trash_dir.exists() and not current.path.exists():
                    trash_dir.rename(current.path)
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status="failed",
                    action="uninstall",
                    ok=False,
                    message="Failed to uninstall domain pack.",
                    error=str(exc),
                )

    def set_enabled(
        self,
        pack_id: str,
        *,
        enabled: bool,
        reason: str = "",
        actor: str = "user",
    ) -> DomainPackGovernanceResult:
        return self._update_config_state(
            pack_id,
            action="enable" if enabled else "disable",
            reason=reason,
            actor=actor,
            mutate=lambda defaults: _set_membership(defaults.disabled, pack_id, present=not enabled),
            idempotent=lambda record: bool(record.get("enabled")) is enabled,
            ok_message="Domain pack enabled." if enabled else "Domain pack disabled.",
            idempotent_message="Domain pack is already enabled." if enabled else "Domain pack is already disabled.",
        )

    def set_active(
        self,
        pack_id: str,
        *,
        active: bool,
        reason: str = "",
        actor: str = "user",
    ) -> DomainPackGovernanceResult:
        return self._update_config_state(
            pack_id,
            action="activate" if active else "deactivate",
            reason=reason,
            actor=actor,
            mutate=lambda defaults: _set_membership(defaults.active, pack_id, present=active),
            idempotent=lambda record: bool(record.get("active_requested")) is active,
            ok_message="Domain pack activated in config." if active else "Domain pack deactivated in config.",
            idempotent_message="Domain pack is already active in config." if active else "Domain pack is already inactive in config.",
        )

    def eval_pack(self, pack_id: str, *, actor: str = "user") -> DomainPackGovernanceResult:
        with self._locked():
            config = self._load_config()
            manager = self._refresh_manager(config)
            pack = manager.get_pack(pack_id)
            if pack is None:
                return _pack_not_found_result(pack_id, "eval")
            latest = self._latest_events_by_pack_unlocked()
            evals = self._latest_eval_events_by_pack_unlocked()
            record = self._decorate_pack(pack, latest.get(pack_id), evals.get(pack_id))
            eval_result = self._run_eval(pack)
            event = self._append_event_unlocked(
                pack_id=pack_id,
                action="eval",
                reason="",
                actor=actor,
                source=str(pack.path),
                previous=_record_snapshot(record),
                next_state=_record_snapshot(record),
                result={"status": eval_result.get("status")},
                overrides_builtin=bool(record.get("overrides_builtin")),
                validation_summary=str(record.get("validation_summary") or ""),
                eval_summary=eval_result,
                artifact_paths=[f"domain_packs/{pack_id}"],
                review_proposal_id="",
            )
            return DomainPackGovernanceResult(
                pack_id=pack_id,
                status=str(eval_result.get("status") or "unknown"),
                action="eval",
                ok=bool(eval_result.get("ok")),
                message="Domain pack evaluation completed.",
                pack=self._decorate_pack(pack, event, event),
                event=event,
                eval_result=eval_result,
            )

    def move_to_domain_capability(self, proposal: dict[str, Any]) -> tuple[bool, str]:
        domain_id = str(proposal.get("domain_id") or "").strip()
        payload = proposal.get("payload") if isinstance(proposal.get("payload"), dict) else {}
        subject_type = str(payload.get("subject_type") or "").strip()
        if subject_type not in _MANAGED_ARTIFACT_TYPES:
            return False, "Only workspace skill and workflow artifacts can move to a domain pack."
        if not domain_id or domain_id == "core":
            return False, "move_to_domain requires a non-core target domain pack."
        try:
            manager = DomainPackManager(self.workspace)
            pack = manager.get_pack(domain_id)
        except Exception:
            pack = None
        if pack is None:
            return False, "Target domain pack was not found."
        if pack.source != "workspace":
            return False, "Target builtin domain pack is read-only in P11."
        if pack.status == "invalid":
            return False, "Target workspace domain pack is invalid."
        return True, ""

    def move_artifact_to_domain(
        self,
        proposal: dict[str, Any],
        *,
        reason: str = "",
        actor: str = "curator",
        review_proposal_id: str = "",
    ) -> DomainPackGovernanceResult:
        with self._locked():
            allowed, unsupported_reason = self.move_to_domain_capability(proposal)
            if not allowed:
                return DomainPackGovernanceResult(
                    pack_id=str(proposal.get("domain_id") or ""),
                    status="pending",
                    action="move_to_domain",
                    ok=False,
                    message=unsupported_reason,
                    error="unsupported",
                )

            domain_id = str(proposal.get("domain_id") or "").strip()
            payload = proposal.get("payload") if isinstance(proposal.get("payload"), dict) else {}
            subject_type = str(payload.get("subject_type") or "").strip()
            subject_id = str(
                payload.get("subject_id")
                or payload.get("skill_name")
                or payload.get("workflow_name")
                or ""
            ).strip()
            manager = DomainPackManager(self.workspace)
            target_pack = manager.get_pack(domain_id)
            if target_pack is None:
                return _pack_not_found_result(domain_id, "move_to_domain")

            if subject_type == "skill":
                return self._move_skill_to_domain_unlocked(
                    target_pack,
                    skill_name=subject_id,
                    reason=reason,
                    actor=actor,
                    review_proposal_id=review_proposal_id,
                )
            if subject_type == "workflow":
                return self._move_workflow_to_domain_unlocked(
                    target_pack,
                    workflow_name=subject_id,
                    reason=reason,
                    actor=actor,
                    review_proposal_id=review_proposal_id,
                )
            return DomainPackGovernanceResult(
                pack_id=domain_id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Unsupported move_to_domain subject.",
                error="unsupported_subject",
            )

    def _move_skill_to_domain_unlocked(
        self,
        target_pack: DomainPack,
        *,
        skill_name: str,
        reason: str,
        actor: str,
        review_proposal_id: str,
    ) -> DomainPackGovernanceResult:
        source_file = self.workspace / "skills" / skill_name / "SKILL.md"
        if not source_file.is_file():
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Workspace skill was not found.",
                error="missing_workspace_skill",
            )
        target_dir = target_pack.path / "skills" / skill_name
        target_file = target_dir / "SKILL.md"
        if target_file.exists():
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Target domain pack already contains that skill.",
                error="target_skill_exists",
            )
        frontmatter, body = _read_skill_markdown(source_file)
        metadata = _metadata_originagent(frontmatter)
        frontmatter["always"] = False
        metadata.update(
            {
                "always": False,
                "migrated_from_workspace": True,
                "original_path": f"skills/{skill_name}/SKILL.md",
                "migrated_by": "curator",
                "migration_review_proposal_id": review_proposal_id,
                "managed_by_domain_pack": True,
            }
        )
        frontmatter["metadata"] = _set_metadata_originagent(frontmatter.get("metadata"), metadata)
        manifest_backup = _read_yaml_file(target_pack.path / "domain_pack.yaml")
        try:
            target_dir.mkdir(parents=True, exist_ok=False)
            _write_skill_markdown(target_file, frontmatter, body)
            updated_manifest = _append_manifest_id(manifest_backup, "skills", skill_name)
            _write_yaml_file(target_pack.path / "domain_pack.yaml", updated_manifest)
            validation = DomainPackValidator(
                runtime_config=DomainPackRuntimeConfig(),
                strict_declarations=True,
            ).validate_pack(
                target_pack.path,
                source="workspace",
            )
            if validation.status == "invalid":
                raise ValueError(validation.validation_summary or validation.unavailable_reason)
            source_dir = source_file.parent
            source_file.unlink()
            source_dir.rmdir()
            manager = DomainPackManager(self.workspace)
            latest = self._latest_events_by_pack_unlocked()
            evals = self._latest_eval_events_by_pack_unlocked()
            record = self._decorate_pack(manager.get_pack(target_pack.id), latest.get(target_pack.id), evals.get(target_pack.id))
            artifact = {
                "artifact_type": "skill",
                "skill_name": skill_name,
                "path": f"domain_packs/{target_pack.id}/skills/{skill_name}/SKILL.md",
                "validation": "Moved into workspace domain pack.",
            }
            event = self._append_event_unlocked(
                pack_id=target_pack.id,
                action="move_to_domain",
                reason=reason,
                actor=actor,
                source=str(source_file),
                previous={"artifact_type": "skill", "path": f"skills/{skill_name}/SKILL.md"},
                next_state={"artifact_type": "skill", "path": artifact["path"]},
                result={"status": "moved"},
                overrides_builtin=bool(record and record.get("overrides_builtin")),
                validation_summary=str(record.get("validation_summary") if record else validation.validation_summary),
                eval_summary=None,
                artifact_paths=[artifact["path"]],
                review_proposal_id=review_proposal_id,
            )
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="moved",
                action="move_to_domain",
                ok=True,
                message="Workspace skill moved into domain pack.",
                pack=record,
                event=event,
                artifact=artifact,
            )
        except Exception as exc:
            _cleanup_dir(target_dir)
            _write_yaml_file(target_pack.path / "domain_pack.yaml", manifest_backup)
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Failed to move workspace skill into domain pack.",
                error=str(exc),
            )

    def _move_workflow_to_domain_unlocked(
        self,
        target_pack: DomainPack,
        *,
        workflow_name: str,
        reason: str,
        actor: str,
        review_proposal_id: str,
    ) -> DomainPackGovernanceResult:
        source_file = self.workspace / "workflows" / workflow_name / "workflow.yaml"
        if not source_file.is_file():
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Workspace workflow was not found.",
                error="missing_workspace_workflow",
            )
        target_dir = target_pack.path / "workflows" / workflow_name
        target_file = target_dir / "workflow.yaml"
        if target_file.exists():
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Target domain pack already contains that workflow.",
                error="target_workflow_exists",
            )
        try:
            data = yaml.safe_load(source_file.read_text(encoding="utf-8"))
        except (OSError, yaml.YAMLError) as exc:
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Workspace workflow could not be read.",
                error=str(exc),
            )
        if not isinstance(data, dict):
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Workspace workflow is invalid.",
                error="invalid_workspace_workflow",
            )
        metadata = data.get("metadata")
        if not isinstance(metadata, dict):
            metadata = {}
            data["metadata"] = metadata
        originagent = read_originagent_metadata(metadata)
        originagent.update(
            {
                "migrated_from_workspace": True,
                "original_path": f"workflows/{workflow_name}/workflow.yaml",
                "migrated_by": "curator",
                "migration_review_proposal_id": review_proposal_id,
                "managed_by_domain_pack": True,
            }
        )
        data["metadata"] = set_originagent_metadata(metadata, originagent)
        manifest_backup = _read_yaml_file(target_pack.path / "domain_pack.yaml")
        try:
            target_dir.mkdir(parents=True, exist_ok=False)
            target_file.write_text(yaml.safe_dump(data, allow_unicode=True, sort_keys=False), encoding="utf-8")
            updated_manifest = _append_manifest_id(manifest_backup, "workflows", workflow_name)
            _write_yaml_file(target_pack.path / "domain_pack.yaml", updated_manifest)
            validation = DomainPackValidator(
                runtime_config=DomainPackRuntimeConfig(),
                strict_declarations=True,
            ).validate_pack(
                target_pack.path,
                source="workspace",
            )
            if validation.status == "invalid":
                raise ValueError(validation.validation_summary or validation.unavailable_reason)
            source_dir = source_file.parent
            source_file.unlink()
            source_dir.rmdir()
            manager = DomainPackManager(self.workspace)
            latest = self._latest_events_by_pack_unlocked()
            evals = self._latest_eval_events_by_pack_unlocked()
            record = self._decorate_pack(manager.get_pack(target_pack.id), latest.get(target_pack.id), evals.get(target_pack.id))
            artifact = {
                "artifact_type": "workflow",
                "workflow_name": workflow_name,
                "path": f"domain_packs/{target_pack.id}/workflows/{workflow_name}/workflow.yaml",
                "validation": "Moved into workspace domain pack.",
            }
            event = self._append_event_unlocked(
                pack_id=target_pack.id,
                action="move_to_domain",
                reason=reason,
                actor=actor,
                source=str(source_file),
                previous={"artifact_type": "workflow", "path": f"workflows/{workflow_name}/workflow.yaml"},
                next_state={"artifact_type": "workflow", "path": artifact["path"]},
                result={"status": "moved"},
                overrides_builtin=bool(record and record.get("overrides_builtin")),
                validation_summary=str(record.get("validation_summary") if record else validation.validation_summary),
                eval_summary=None,
                artifact_paths=[artifact["path"]],
                review_proposal_id=review_proposal_id,
            )
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="moved",
                action="move_to_domain",
                ok=True,
                message="Workspace workflow moved into domain pack.",
                pack=record,
                event=event,
                artifact=artifact,
            )
        except Exception as exc:
            _cleanup_dir(target_dir)
            _write_yaml_file(target_pack.path / "domain_pack.yaml", manifest_backup)
            return DomainPackGovernanceResult(
                pack_id=target_pack.id,
                status="failed",
                action="move_to_domain",
                ok=False,
                message="Failed to move workspace workflow into domain pack.",
                error=str(exc),
            )

    def _latest_events_by_pack_unlocked(self) -> dict[str, dict[str, Any]]:
        latest: dict[str, dict[str, Any]] = {}
        for event in self._read_events_unlocked():
            pack_id = str(event.get("pack_id") or "")
            if pack_id:
                latest[pack_id] = event
        return latest

    def _latest_eval_events_by_pack_unlocked(self) -> dict[str, dict[str, Any]]:
        latest: dict[str, dict[str, Any]] = {}
        for event in self._read_events_unlocked():
            if str(event.get("action") or "") != "eval":
                continue
            pack_id = str(event.get("pack_id") or "")
            if pack_id:
                latest[pack_id] = event
        return latest

    def _decorate_pack(
        self,
        pack: DomainPack | None,
        latest_event: dict[str, Any] | None,
        latest_eval_event: dict[str, Any] | None,
    ) -> dict[str, Any] | None:
        if pack is None:
            return None
        return {
            "id": pack.id,
            "name": pack.name,
            "version": pack.version,
            "path": str(pack.path),
            "source": pack.source,
            "status": pack.status,
            "enabled": pack.enabled,
            "active": pack.active,
            "active_requested": pack.active_requested,
            "verification_status": pack.verification_status,
            "overrides_builtin": pack.overrides_builtin,
            "description": pack.description,
            "validation_summary": pack.validation_summary or pack.unavailable_reason or "Domain pack is valid.",
            "unavailable_reason": pack.unavailable_reason,
            "capabilities": list(pack.capabilities),
            "triggers": list(pack.triggers),
            "dependencies": {"packs": list(pack.dependencies.packs)},
            "source_info": asdict(pack.source_info),
            "skills": [_skill_decl_json(item) for item in pack.skills],
            "workflows": [_workflow_decl_json(item) for item in pack.workflows],
            "policies": [_file_decl_json(item) for item in pack.policies],
            "schemas": [_file_decl_json(item) for item in pack.schemas],
            "tools": [_tool_decl_json(item) for item in pack.tools],
            "evals": [_eval_decl_json(item) for item in pack.evals],
            "can_install": False,
            "can_upgrade": pack.source == "workspace",
            "can_uninstall": pack.source == "workspace",
            "can_enable": not pack.enabled,
            "can_disable": pack.enabled,
            "can_activate": not pack.active_requested,
            "can_deactivate": pack.active_requested,
            "can_eval": True,
            "disabled_reason": "" if pack.source == "workspace" else "Builtin domain packs are read-only for install, upgrade, uninstall, and move.",
            "last_event": latest_event,
            "last_eval_result": (
                latest_eval_event.get("eval_summary")
                if isinstance(latest_eval_event, dict)
                else None
            ),
        }

    def _update_config_state(
        self,
        pack_id: str,
        *,
        action: str,
        reason: str,
        actor: str,
        mutate: Callable[[Any], list[str]],
        idempotent: Callable[[dict[str, Any]], bool],
        ok_message: str,
        idempotent_message: str,
    ) -> DomainPackGovernanceResult:
        with self._locked():
            config = self._load_config()
            manager = self._refresh_manager(config)
            current = manager.get_pack(pack_id)
            if current is None:
                return _pack_not_found_result(pack_id, action)
            previous = self._decorate_pack(current, None, None)
            if previous is not None and idempotent(previous):
                return DomainPackGovernanceResult(
                    pack_id=pack_id,
                    status=str(previous.get("status") or "unknown"),
                    action=action,
                    ok=True,
                    message=idempotent_message,
                    pack=previous,
                )
            defaults = config.agents.defaults.domain_packs
            updated = mutate(defaults)
            if action == "disable":
                defaults.disabled = updated
            elif action == "enable":
                defaults.disabled = updated
            elif action == "activate":
                defaults.active = updated
            elif action == "deactivate":
                defaults.active = updated
            self._save_config(config)
            manager = self._refresh_manager(config)
            record = self._decorate_pack(manager.get_pack(pack_id), None, None)
            event = self._append_event_unlocked(
                pack_id=pack_id,
                action=action,
                reason=reason,
                actor=actor,
                source=str(current.path),
                previous=_record_snapshot(previous),
                next_state=_record_snapshot(record),
                result={"status": action},
                overrides_builtin=bool(record and record.get("overrides_builtin")),
                validation_summary=str(record.get("validation_summary") if record else current.validation_summary),
                eval_summary=None,
                artifact_paths=[f"domain_packs/{pack_id}"],
                review_proposal_id="",
            )
            return DomainPackGovernanceResult(
                pack_id=pack_id,
                status=action,
                action=action,
                ok=True,
                message=ok_message,
                pack=record,
                event=event,
            )

    def _run_eval(self, pack: DomainPack) -> dict[str, Any]:
        checks: list[dict[str, Any]] = []
        warnings: list[str] = []
        errors: list[str] = []
        validator = DomainPackValidator(
            runtime_config=DomainPackRuntimeConfig(),
            strict_declarations=True,
        )
        validated = validator.validate_pack(pack.path, source=pack.source)
        if validated.status == "invalid":
            errors.append(validated.validation_summary or validated.unavailable_reason or "Invalid domain pack.")
        elif validated.status == "unavailable" and validated.unavailable_reason:
            warnings.append(validated.unavailable_reason)
        checks.append(
            {
                "kind": "manifest",
                "ok": validated.status != "invalid",
                "message": validated.validation_summary or validated.unavailable_reason or "Domain pack is valid.",
            }
        )
        declarations = list(pack.evals) if pack.evals else [
            DomainEvalDeclaration(id="manifest", kind="manifest"),
            DomainEvalDeclaration(id="skills", kind="skill"),
            DomainEvalDeclaration(id="tools", kind="tool"),
            DomainEvalDeclaration(id="workflows", kind="workflow"),
        ]
        for declaration in declarations:
            if declaration.kind == "manifest":
                continue
            if declaration.kind not in _SUPPORTED_EVAL_KINDS:
                errors.append(f"Unsupported eval kind `{declaration.kind}`.")
                continue
            if declaration.kind == "skill":
                items = _filter_named_declarations(pack.skills, declaration.target)
                for item in items:
                    ok = item.status == "available" and item.path is not None and item.path.exists()
                    message = "Skill artifact is present." if ok else item.unavailable_reason or "Missing skill artifact."
                    checks.append({"kind": "skill", "target": item.id, "ok": ok, "message": message})
                    if not ok:
                        errors.append(f"skill `{item.id}`: {message}")
            elif declaration.kind == "workflow":
                items = _filter_named_declarations(pack.workflows, declaration.target)
                for item in items:
                    ok = False
                    message = item.unavailable_reason or "Missing workflow artifact."
                    if item.status == "available" and item.path is not None:
                        ok, message = validate_workflow_artifact_dir(
                            item.path.parent,
                            workspace=pack.path,
                            expected_name=item.id,
                        )
                    checks.append({"kind": "workflow", "target": item.id, "ok": ok, "message": message})
                    if not ok:
                        errors.append(f"workflow `{item.id}`: {message}")
            elif declaration.kind == "tool":
                items = _filter_named_declarations(pack.tools, declaration.target)
                for item in items:
                    ok, message = _tool_eval_result(pack, item)
                    checks.append({"kind": "tool", "target": item.id, "ok": ok, "message": message})
                    if not ok:
                        errors.append(f"tool `{item.id}`: {message}")
        status = "error" if errors else "warning" if warnings else "ok"
        return {
            "ok": not errors,
            "status": status,
            "checks": checks,
            "warnings": warnings,
            "errors": errors,
            "pack": pack.id,
        }

    @staticmethod
    def _copy_pack_to_staging(source_dir: Path, staging_dir: Path) -> None:
        _cleanup_dir(staging_dir)
        shutil.copytree(source_dir, staging_dir)


def summarize_domain_pack_governance(workspace: Path, manager: DomainPackManager | None = None) -> dict[str, Any]:
    service = DomainPackGovernanceService(workspace, domain_pack_manager=manager)
    return service.stats()


def _resolve_local_dir(source_path: str) -> Path | None:
    raw = str(source_path or "").strip()
    if not raw:
        return None
    try:
        path = Path(raw).expanduser().resolve()
    except OSError:
        return None
    if not path.exists() or not path.is_dir():
        return None
    if not (path / "domain_pack.yaml").is_file():
        return None
    return path


def _cleanup_dir(path: Path) -> None:
    try:
        if path.exists():
            shutil.rmtree(path)
    except OSError:
        pass


def _read_only_result(pack_id: str, action: str, pack: DomainPack) -> DomainPackGovernanceResult:
    return DomainPackGovernanceResult(
        pack_id=pack_id,
        status=pack.status,
        action=action,
        ok=False,
        message="Builtin domain pack is read-only for this action.",
        error="read_only",
    )


def _pack_not_found_result(pack_id: str, action: str) -> DomainPackGovernanceResult:
    return DomainPackGovernanceResult(
        pack_id=pack_id,
        status="missing",
        action=action,
        ok=False,
        message="Domain pack was not found.",
        error="not_found",
    )


def _set_membership(items: list[str], value: str, *, present: bool) -> list[str]:
    cleaned = [item for item in items if item != value]
    if present:
        cleaned.append(value)
    return cleaned


def _read_yaml_file(path: Path) -> dict[str, Any]:
    raw = yaml.safe_load(path.read_text(encoding="utf-8"))
    return raw if isinstance(raw, dict) else {}


def _write_yaml_file(path: Path, data: dict[str, Any]) -> None:
    path.write_text(yaml.safe_dump(data, allow_unicode=True, sort_keys=False), encoding="utf-8")


def _append_manifest_id(manifest: dict[str, Any], key: str, item_id: str) -> dict[str, Any]:
    updated = dict(manifest)
    raw = updated.get(key)
    if raw is None:
        updated[key] = [item_id]
        return updated
    if not isinstance(raw, list):
        raise ValueError(f"domain_pack.yaml field `{key}` must be a list")
    existing_ids: list[str] = []
    for item in raw:
        if isinstance(item, str):
            existing_ids.append(item.strip())
        elif isinstance(item, dict):
            existing_ids.append(str(item.get("id") or item.get("name") or "").strip())
    if item_id in existing_ids:
        return updated
    raw.append(item_id)
    updated[key] = raw
    return updated


def _record_snapshot(record: dict[str, Any] | None) -> dict[str, Any] | None:
    if record is None:
        return None
    return {
        "id": record.get("id"),
        "status": record.get("status"),
        "enabled": record.get("enabled"),
        "active": record.get("active"),
        "active_requested": record.get("active_requested"),
        "verification_status": record.get("verification_status"),
        "overrides_builtin": record.get("overrides_builtin"),
    }


def _clean_text(value: Any, max_chars: int) -> str:
    return truncate_text(str(value or "").strip(), max_chars)


def _skill_decl_json(item: DomainSkillDeclaration) -> dict[str, Any]:
    return {
        "id": item.id,
        "virtual_id": item.virtual_id,
        "path": str(item.path) if item.path is not None else "",
        "status": item.status,
        "unavailable_reason": item.unavailable_reason,
    }


def _workflow_decl_json(item: DomainWorkflowDeclaration) -> dict[str, Any]:
    return {
        "id": item.id,
        "path": str(item.path) if item.path is not None else "",
        "status": item.status,
        "unavailable_reason": item.unavailable_reason,
    }


def _file_decl_json(item: DomainFileDeclaration) -> dict[str, Any]:
    return {
        "id": item.id,
        "path": str(item.path) if item.path is not None else "",
        "status": item.status,
        "unavailable_reason": item.unavailable_reason,
    }


def _tool_decl_json(item: DomainToolDeclaration) -> dict[str, Any]:
    return {
        "id": item.id,
        "module": item.module,
        "class_name": item.class_name,
        "permissions": list(item.permissions),
        "audit": item.audit,
        "module_path": str(item.module_path) if item.module_path is not None else "",
        "status": item.status,
        "unavailable_reason": item.unavailable_reason,
    }


def _eval_decl_json(item: DomainEvalDeclaration) -> dict[str, Any]:
    return {
        "id": item.id,
        "kind": item.kind,
        "target": item.target,
        "status": item.status,
        "unavailable_reason": item.unavailable_reason,
    }


def _filter_named_declarations(items: tuple[Any, ...], target: str) -> list[Any]:
    if not target:
        return list(items)
    return [item for item in items if str(getattr(item, "id", "")) == target]


def _tool_eval_result(pack: DomainPack, declaration: DomainToolDeclaration) -> tuple[bool, str]:
    if declaration.status != "available" or declaration.module_path is None:
        return False, declaration.unavailable_reason or "Tool declaration is invalid."
    module_name = f"_originagent_domain_pack_eval_{pack.id}_{declaration.id}_{uuid.uuid4().hex}"
    try:
        spec = importlib.util.spec_from_file_location(module_name, declaration.module_path.resolve())
        if spec is None or spec.loader is None:
            return False, "could not create module spec"
        module = importlib.util.module_from_spec(spec)
        sys.modules[module_name] = module
        spec.loader.exec_module(module)
        attr = getattr(module, declaration.class_name, None)
        if not isinstance(attr, type):
            return False, f"class {declaration.class_name} not found"
        if not issubclass(attr, Tool) or attr is Tool:
            return False, f"class {declaration.class_name} is not a Tool"
        return True, "Tool module imported successfully."
    except Exception as exc:
        return False, f"import failed: {exc}"
    finally:
        sys.modules.pop(module_name, None)
