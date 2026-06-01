"""Local activation and rollback for verified evolution modules.

Domain pack activation calls ``DomainPackGovernanceService`` while holding the
activation lock; Phase 1 does not allow governance code to call back into this
activator because the reverse lock order would risk deadlock.

Domain pack activation also writes ``domain_pack_events.jsonl`` through the
governance service in addition to ``evolution_events.jsonl``. That dual stream
is intentional bridge behavior.
"""

from __future__ import annotations

import json
import os
import re
import shutil
import uuid
from contextlib import suppress
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from filelock import FileLock

from OriginAgent.agent.domain_pack_governance import DomainPackGovernanceService
from OriginAgent.agent.skill_lifecycle import SkillLifecycleStore
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger, canonical_dump
from OriginAgent.evolution.package import EVOLUTION_MANIFEST_FILENAME, read_package_manifest
from OriginAgent.evolution.verifier import EvolutionModuleVerifier

ACTIVATION_SCHEMA_VERSION = "originagent.evolution.activation.v1"
_SAFE_SKILL_RE = re.compile(r"^[a-z0-9][a-z0-9-]{0,63}$")


@dataclass(frozen=True)
class EvolutionActivationResult:
    ok: bool
    status: str
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    artifact_digest: str = ""
    activation_id: str = ""
    target_paths: tuple[str, ...] = ()
    events: tuple[EvolutionEvent, ...] = ()
    error: str = ""


class EvolutionModuleActivator:
    """Activate or roll back verified staged modules without executing module code."""

    def __init__(
        self,
        workspace: Path,
        ledger: EvolutionLedger | None = None,
        lock_path: Path | None = None,
        config_loader: Callable[[], Any] | None = None,
        config_saver: Callable[[Any], None] | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self.memory_dir = self.workspace / "memory"
        self.staging_root = self.memory_dir / "evolution_staging"
        self.activation_root = self.memory_dir / "evolution_activations"
        self.branch_root = self.memory_dir / "evolution_branches"
        self.ledger = ledger or EvolutionLedger(self.workspace)
        self._lock_path = Path(lock_path) if lock_path is not None else self.memory_dir / ".evolution_activation.lock"
        self._config_loader = config_loader
        self._config_saver = config_saver

    def activate(self, artifact_digest: str, *, actor: str = "user") -> EvolutionActivationResult:
        with self._locked():
            metadata = self._read_activation_metadata(artifact_digest)
            if metadata and metadata.get("status") == "active":
                event = self._append_activation_event(
                    actor=actor,
                    metadata=metadata,
                    status="already_active",
                )
                return _result_from_metadata(
                    metadata,
                    ok=True,
                    status="already_active",
                    events=(event,),
                )

            context = self._load_activation_context(artifact_digest)
            if context["error"]:
                return self._failed_result(
                    artifact_digest=artifact_digest,
                    actor=actor,
                    status=str(context["status"]),
                    error=str(context["error"]),
                    context=context,
                )

            module_type = str(context["module_type"])
            if module_type == "workflow" or module_type == "tool":
                return self._failed_result(
                    artifact_digest=artifact_digest,
                    actor=actor,
                    status="unsupported",
                    error=f"activation is unsupported for module_type: {module_type}",
                    context=context,
                )
            if self._has_active_state_branch(artifact_digest):
                return self._failed_result(
                    artifact_digest=artifact_digest,
                    actor=actor,
                    status="active_state_branch",
                    error="artifact has an active state branch",
                    context=context,
                )

            if module_type == "skill":
                return self._activate_skill(context, actor=actor)
            if module_type == "domain_pack":
                return self._activate_domain_pack(context, actor=actor)
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="unsupported",
                error=f"activation is unsupported for module_type: {module_type}",
                context=context,
            )

    def rollback(self, artifact_digest: str, *, actor: str = "user") -> EvolutionActivationResult:
        with self._locked():
            metadata = self._read_activation_metadata(artifact_digest)
            if metadata is None:
                event = self._append_failed_event(
                    actor=actor,
                    artifact_digest=artifact_digest,
                    status="not_active",
                    error="activation metadata is missing",
                )
                return EvolutionActivationResult(
                    ok=False,
                    status="not_active",
                    artifact_digest=artifact_digest,
                    events=(event,),
                    error="activation metadata is missing",
                )
            if metadata.get("status") == "rolled_back":
                event = self._append_rollback_succeeded_event(
                    actor=actor,
                    metadata=metadata,
                    status="already_rolled_back",
                )
                return _result_from_metadata(
                    metadata,
                    ok=True,
                    status="already_rolled_back",
                    events=(event,),
                )
            if metadata.get("status") != "active":
                event = self._append_rollback_failed_event(
                    actor=actor,
                    metadata=metadata,
                    status="not_active",
                    error=f"activation is not active: {metadata.get('status')}",
                )
                return _result_from_metadata(
                    metadata,
                    ok=False,
                    status="not_active",
                    events=(event,),
                    error=f"activation is not active: {metadata.get('status')}",
                )

            started = self.ledger.append(
                EvolutionEvent.new(
                    EventType.MODULE_ROLLBACK_STARTED,
                    actor=actor,
                    module_id=str(metadata.get("module_id") or ""),
                    module_version=str(metadata.get("version") or ""),
                    module_type=str(metadata.get("module_type") or ""),
                    artifact_digest=artifact_digest,
                    result={
                        "status": "rollback_started",
                        "target_paths": list(_target_paths(metadata)),
                    },
                )
            )
            events = [started]

            module_type = str(metadata.get("module_type") or "")
            if module_type == "skill":
                failed = self._rollback_skill(metadata, actor=actor)
            elif module_type == "domain_pack":
                failed = self._rollback_domain_pack(metadata, actor=actor)
            else:
                failed = f"rollback is unsupported for module_type: {module_type}"

            if failed:
                event = self._append_rollback_failed_event(
                    actor=actor,
                    metadata=metadata,
                    status="rollback_failed",
                    error=failed,
                )
                events.append(event)
                dirty = self._record_dirty_rollback(
                    actor=actor,
                    metadata=metadata,
                    reason=failed,
                )
                events.append(dirty)
                return _result_from_metadata(
                    metadata,
                    ok=False,
                    status="dirty_rollback",
                    events=tuple(events),
                    error=failed,
                )

            succeeded_event = EvolutionEvent.new(
                EventType.MODULE_ROLLBACK_SUCCEEDED,
                actor=actor,
                module_id=str(metadata.get("module_id") or ""),
                module_version=str(metadata.get("version") or ""),
                module_type=module_type,
                artifact_digest=artifact_digest,
                result={
                    "status": "rolled_back",
                    "target_paths": list(_target_paths(metadata)),
                },
            )
            metadata["status"] = "rolled_back"
            metadata["rolled_back_at"] = datetime.now(timezone.utc).isoformat()
            metadata["rollback_event_id"] = succeeded_event.event_id
            _write_json_atomic(self._activation_metadata_path(artifact_digest), metadata)
            succeeded = self.ledger.append(succeeded_event)
            events.append(succeeded)
            return _result_from_metadata(
                metadata,
                ok=True,
                status="rolled_back",
                events=tuple(events),
            )

    def _activate_skill(self, context: dict[str, Any], *, actor: str) -> EvolutionActivationResult:
        module_id = str(context["module_id"])
        artifact_digest = str(context["artifact_digest"])
        if not _SAFE_SKILL_RE.fullmatch(module_id):
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="invalid_skill_name",
                error="skill module_id is not a safe workspace skill name",
                context=context,
            )
        artifact_dir = Path(context["artifact_dir"])
        if not (artifact_dir / "SKILL.md").is_file():
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="missing_skill_manifest",
                error="skill artifact must include SKILL.md",
                context=context,
            )
        target_rel = f"skills/{module_id}"
        target_dir = self.workspace / target_rel
        if target_dir.exists():
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="skill_exists",
                error="workspace skill already exists",
                context=context,
                target_paths=(target_rel,),
            )

        tmp_dir = target_dir.parent / f".{module_id}.evolution.{uuid.uuid4().hex}"
        try:
            _copy_skill_artifact(artifact_dir, tmp_dir)
            target_dir.parent.mkdir(parents=True, exist_ok=True)
            tmp_dir.rename(target_dir)
        except Exception as exc:
            _cleanup_path(tmp_dir)
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="copy_failed",
                error=str(exc),
                context=context,
                target_paths=(target_rel,),
            )

        lifecycle = SkillLifecycleStore(self.workspace)
        verified = lifecycle.transition(
            module_id,
            action="verify",
            actor=actor,
            reason="evolution activation",
        )
        if not verified.ok:
            _cleanup_path(target_dir)
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="skill_verify_failed",
                error=verified.error or verified.message,
                context=context,
                target_paths=(target_rel,),
            )

        activated = lifecycle.transition(
            module_id,
            action="activate",
            actor=actor,
            reason="evolution activation",
        )
        if not activated.ok:
            deprecated = lifecycle.transition(
                module_id,
                action="deprecate",
                actor=actor,
                reason="evolution activation failed",
            )
            suffix = "; deprecate attempted" if deprecated.ok else f"; deprecate failed: {deprecated.error}"
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="skill_activate_failed",
                error=(activated.error or activated.message) + suffix,
                context=context,
                target_paths=(target_rel,),
            )

        return self._write_activation_success(context, actor=actor, target_paths=(target_rel,))

    def _activate_domain_pack(self, context: dict[str, Any], *, actor: str) -> EvolutionActivationResult:
        module_id = str(context["module_id"])
        artifact_digest = str(context["artifact_digest"])
        artifact_dir = Path(context["artifact_dir"])
        if not (artifact_dir / "domain_pack.yaml").is_file():
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="missing_domain_pack_manifest",
                error="domain_pack artifact must include domain_pack.yaml",
                context=context,
            )
        target_rel = f"domain_packs/{module_id}"
        if (self.workspace / target_rel).exists():
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="domain_pack_exists",
                error="workspace domain pack already exists",
                context=context,
                target_paths=(target_rel,),
            )

        service = DomainPackGovernanceService(
            self.workspace,
            config_loader=self._config_loader,
            config_saver=self._config_saver,
        )
        installed = service.install(
            str(artifact_dir),
            actor=actor,
            reason="evolution activation",
        )
        if not installed.ok:
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="domain_pack_install_failed",
                error=installed.error or installed.message,
                context=context,
                target_paths=(target_rel,),
            )
        activated = service.set_active(
            module_id,
            active=True,
            actor=actor,
            reason="evolution activation",
        )
        if not activated.ok:
            with suppress(Exception):
                service.uninstall(module_id, actor=actor, reason="evolution activation rollback")
            return self._failed_result(
                artifact_digest=artifact_digest,
                actor=actor,
                status="domain_pack_activate_failed",
                error=activated.error or activated.message,
                context=context,
                target_paths=(target_rel,),
            )
        return self._write_activation_success(context, actor=actor, target_paths=(target_rel,))

    def _rollback_skill(self, metadata: dict[str, Any], *, actor: str) -> str:
        result = SkillLifecycleStore(self.workspace).transition(
            str(metadata.get("module_id") or ""),
            action="deprecate",
            actor=actor,
            reason="evolution rollback",
        )
        if result.ok:
            return ""
        return result.error or result.message

    def _rollback_domain_pack(self, metadata: dict[str, Any], *, actor: str) -> str:
        module_id = str(metadata.get("module_id") or "")
        service = DomainPackGovernanceService(
            self.workspace,
            config_loader=self._config_loader,
            config_saver=self._config_saver,
        )
        deactivated = service.set_active(
            module_id,
            active=False,
            actor=actor,
            reason="evolution rollback",
        )
        if not deactivated.ok:
            return deactivated.error or deactivated.message
        uninstalled = service.uninstall(module_id, actor=actor, reason="evolution rollback")
        if not uninstalled.ok:
            return uninstalled.error or uninstalled.message
        return ""

    def _record_dirty_rollback(
        self,
        *,
        actor: str,
        metadata: dict[str, Any],
        reason: str,
    ) -> EvolutionEvent:
        now = datetime.now(timezone.utc).isoformat()
        target_paths = list(_target_paths(metadata))
        metadata["status"] = "dirty_rollback"
        metadata["dirty_rollback_at"] = now
        metadata["dirty_reason"] = reason
        metadata["residual_resources"] = target_paths
        metadata.setdefault("teardown_attempts", [])
        metadata.setdefault("force_cleaned_at", "")
        _write_json_atomic(self._activation_metadata_path(str(metadata.get("artifact_digest") or "")), metadata)
        return self.ledger.append(
            EvolutionEvent.new(
                EventType.DIRTY_ROLLBACK,
                actor=actor,
                module_id=str(metadata.get("module_id") or ""),
                module_version=str(metadata.get("version") or ""),
                module_type=str(metadata.get("module_type") or ""),
                artifact_digest=str(metadata.get("artifact_digest") or ""),
                result={
                    "status": "dirty_rollback",
                    "reason": reason,
                    "residual_resources": target_paths,
                },
            )
        )

    def _write_activation_success(
        self,
        context: dict[str, Any],
        *,
        actor: str,
        target_paths: tuple[str, ...],
    ) -> EvolutionActivationResult:
        activation_id = f"activation_{uuid.uuid4().hex}"
        event = EvolutionEvent.new(
            EventType.MODULE_ACTIVATED,
            actor=actor,
            module_id=str(context["module_id"]),
            module_version=str(context["module_version"]),
            module_type=str(context["module_type"]),
            artifact_digest=str(context["artifact_digest"]),
            result={
                "status": "active",
                "staging_path": str(context["staging_path"]),
                "target_paths": list(target_paths),
            },
        )
        metadata = {
            "schema_version": ACTIVATION_SCHEMA_VERSION,
            "activation_id": activation_id,
            "artifact_digest": str(context["artifact_digest"]),
            "module_id": str(context["module_id"]),
            "module_type": str(context["module_type"]),
            "version": str(context["module_version"]),
            "status": "active",
            "activated_at": datetime.now(timezone.utc).isoformat(),
            "rolled_back_at": "",
            "target_paths": list(target_paths),
            "activation_event_id": event.event_id,
            "rollback_event_id": "",
        }
        _write_json_atomic(self._activation_metadata_path(str(context["artifact_digest"])), metadata)
        appended = self.ledger.append(event)
        return _result_from_metadata(metadata, ok=True, status="active", events=(appended,))

    def _load_activation_context(self, artifact_digest: str) -> dict[str, Any]:
        context: dict[str, Any] = {
            "artifact_digest": artifact_digest,
            "module_id": "",
            "module_type": "",
            "module_version": "",
            "staging_path": "",
            "artifact_dir": "",
            "status": "",
            "error": "",
        }
        staging_dir = self.staging_root / artifact_digest
        artifact_dir = staging_dir / "artifact"
        context["staging_path"] = _relative_to_workspace(staging_dir, self.workspace)
        context["artifact_dir"] = artifact_dir
        if not artifact_dir.is_dir():
            context["status"] = "missing_staging"
            context["error"] = "staging artifact is missing"
            return context
        if not self._has_verified_event(artifact_digest):
            context["status"] = "unverified"
            context["error"] = "artifact_digest has no verified module event"
            return context
        report = EvolutionModuleVerifier(
            self.workspace,
            staging_root=self.staging_root,
        ).verify(artifact_digest)
        context.update(
            {
                "module_id": report.module_id,
                "module_type": report.module_type,
                "module_version": report.module_version,
            }
        )
        if not report.ok:
            context["status"] = "verification_failed"
            context["error"] = report.error
            return context
        try:
            manifest = read_package_manifest(artifact_dir)
        except Exception as exc:
            context["status"] = "manifest_failed"
            context["error"] = str(exc)
            return context
        context.update(
            {
                "module_id": manifest.module_id,
                "module_type": manifest.module_type,
                "module_version": manifest.version,
            }
        )
        return context

    def _failed_result(
        self,
        *,
        artifact_digest: str,
        actor: str,
        status: str,
        error: str,
        context: dict[str, Any] | None = None,
        target_paths: tuple[str, ...] = (),
    ) -> EvolutionActivationResult:
        event = self._append_failed_event(
            actor=actor,
            artifact_digest=artifact_digest,
            status=status,
            error=error,
            context=context,
            target_paths=target_paths,
        )
        return EvolutionActivationResult(
            ok=False,
            status=status,
            module_id=str((context or {}).get("module_id") or ""),
            module_type=str((context or {}).get("module_type") or ""),
            module_version=str((context or {}).get("module_version") or ""),
            artifact_digest=artifact_digest,
            target_paths=target_paths,
            events=(event,),
            error=error,
        )

    def _append_activation_event(
        self,
        *,
        actor: str,
        metadata: dict[str, Any],
        status: str,
    ) -> EvolutionEvent:
        return self.ledger.append(
            EvolutionEvent.new(
                EventType.MODULE_ACTIVATED,
                actor=actor,
                module_id=str(metadata.get("module_id") or ""),
                module_version=str(metadata.get("version") or ""),
                module_type=str(metadata.get("module_type") or ""),
                artifact_digest=str(metadata.get("artifact_digest") or ""),
                result={
                    "status": status,
                    "target_paths": list(_target_paths(metadata)),
                },
            )
        )

    def _append_failed_event(
        self,
        *,
        actor: str,
        artifact_digest: str,
        status: str,
        error: str,
        context: dict[str, Any] | None = None,
        target_paths: tuple[str, ...] = (),
    ) -> EvolutionEvent:
        context = context or {}
        result: dict[str, Any] = {
            "status": status,
            "error": error,
        }
        staging_path = str(context.get("staging_path") or "")
        if staging_path:
            result["staging_path"] = staging_path
        if target_paths:
            result["target_paths"] = list(target_paths)
        return self.ledger.append(
            EvolutionEvent.new(
                EventType.MODULE_FAILED,
                actor=actor,
                module_id=str(context.get("module_id") or ""),
                module_version=str(context.get("module_version") or ""),
                module_type=str(context.get("module_type") or ""),
                artifact_digest=artifact_digest,
                result=result,
            )
        )

    def _append_rollback_succeeded_event(
        self,
        *,
        actor: str,
        metadata: dict[str, Any],
        status: str,
    ) -> EvolutionEvent:
        return self.ledger.append(
            EvolutionEvent.new(
                EventType.MODULE_ROLLBACK_SUCCEEDED,
                actor=actor,
                module_id=str(metadata.get("module_id") or ""),
                module_version=str(metadata.get("version") or ""),
                module_type=str(metadata.get("module_type") or ""),
                artifact_digest=str(metadata.get("artifact_digest") or ""),
                result={
                    "status": status,
                    "target_paths": list(_target_paths(metadata)),
                },
            )
        )

    def _append_rollback_failed_event(
        self,
        *,
        actor: str,
        metadata: dict[str, Any],
        status: str,
        error: str,
    ) -> EvolutionEvent:
        return self.ledger.append(
            EvolutionEvent.new(
                EventType.MODULE_ROLLBACK_FAILED,
                actor=actor,
                module_id=str(metadata.get("module_id") or ""),
                module_version=str(metadata.get("version") or ""),
                module_type=str(metadata.get("module_type") or ""),
                artifact_digest=str(metadata.get("artifact_digest") or ""),
                result={
                    "status": status,
                    "error": error,
                    "target_paths": list(_target_paths(metadata)),
                },
            )
        )

    def _has_verified_event(self, artifact_digest: str) -> bool:
        verification = self.ledger.verify_chain()
        if not verification.ok or not self.ledger.event_path.exists():
            return False
        with self.ledger._locked():
            with self.ledger.event_path.open("r", encoding="utf-8") as handle:
                return any(
                    bool(line.strip())
                    and (event := json.loads(line)).get("event_type") == EventType.MODULE_VERIFIED.value
                    and event.get("artifact_digest") == artifact_digest
                    for line in handle
                )

    def _has_active_state_branch(self, artifact_digest: str) -> bool:
        if not self.branch_root.exists():
            return False
        for branch_json in self.branch_root.glob("*/branch.json"):
            try:
                data = json.loads(branch_json.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError):
                continue
            if data.get("artifact_digest") == artifact_digest and data.get("status") == "active":
                return True
        return False

    def _read_activation_metadata(self, artifact_digest: str) -> dict[str, Any] | None:
        path = self._activation_metadata_path(artifact_digest)
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

    def _activation_metadata_path(self, artifact_digest: str) -> Path:
        return self.activation_root / artifact_digest / "activation.json"

    def _locked(self) -> FileLock:
        self.memory_dir.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))


def _copy_skill_artifact(source: Path, target: Path) -> None:
    _cleanup_path(target)
    target.mkdir(parents=True, exist_ok=True)
    for root, dirnames, filenames in os.walk(source):
        current = Path(root)
        dirnames[:] = sorted(dirnames)
        for filename in sorted(filenames):
            if filename == EVOLUTION_MANIFEST_FILENAME:
                continue
            file_path = current / filename
            if not file_path.is_file():
                continue
            rel_path = file_path.relative_to(source)
            destination = target / rel_path
            destination.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(file_path, destination)


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
    except BaseException:
        tmp_path.unlink(missing_ok=True)
        raise


def _cleanup_path(path: Path) -> None:
    if path.is_dir() and not path.is_symlink():
        shutil.rmtree(path, ignore_errors=True)
    else:
        with suppress(FileNotFoundError):
            path.unlink()


def _target_paths(metadata: dict[str, Any]) -> tuple[str, ...]:
    raw = metadata.get("target_paths") or ()
    if isinstance(raw, list):
        return tuple(str(item) for item in raw)
    return ()


def _result_from_metadata(
    metadata: dict[str, Any],
    *,
    ok: bool,
    status: str,
    events: tuple[EvolutionEvent, ...],
    error: str = "",
) -> EvolutionActivationResult:
    return EvolutionActivationResult(
        ok=ok,
        status=status,
        module_id=str(metadata.get("module_id") or ""),
        module_type=str(metadata.get("module_type") or ""),
        module_version=str(metadata.get("version") or ""),
        artifact_digest=str(metadata.get("artifact_digest") or ""),
        activation_id=str(metadata.get("activation_id") or ""),
        target_paths=_target_paths(metadata),
        events=events,
        error=error,
    )


def _relative_to_workspace(path: Path, workspace: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return path.name
