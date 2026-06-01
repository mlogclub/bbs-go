"""Local evolution module staging manager."""

from __future__ import annotations

import json
import shutil
import uuid
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable, Iterable, Literal, Mapping

from filelock import FileLock

from OriginAgent.evolution.activation import EvolutionActivationResult, EvolutionModuleActivator
from OriginAgent.evolution.capability_gate import EvolutionCapabilityGate, EvolutionCapabilityResult
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger, LedgerStatus, canonical_dump
from OriginAgent.evolution.package import (
    EvolutionPackage,
    copy_artifact,
    load_package,
)
from OriginAgent.evolution.recovery import EvolutionRecoveryManager, EvolutionRecoveryResult
from OriginAgent.evolution.state_branch import (
    EvolutionStateBranchResult,
    EvolutionStateBranchStore,
)
from OriginAgent.evolution.telemetry import (
    EvolutionProofBundleResult,
    EvolutionTelemetryRecorder,
    EvolutionTelemetryResult,
    EvolutionTokenBudgetResult,
)
from OriginAgent.evolution.verifier import EvolutionModuleVerifier, EvolutionVerificationReport
from OriginAgent.security.capabilities import CapabilitySnapshot

STAGING_SCHEMA_VERSION = "originagent.evolution.staging.v1"
StageStatus = Literal["staged", "already_staged", "failed", "dirty_rollback_blocked"]
VerificationStatus = Literal["verified", "failed"]


@dataclass(frozen=True)
class EvolutionStageResult:
    ok: bool
    status: StageStatus
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    artifact_digest: str = ""
    staging_path: str = ""
    events: tuple[EvolutionEvent, ...] = ()
    error: str = ""


@dataclass(frozen=True)
class EvolutionVerificationResult:
    ok: bool
    status: VerificationStatus
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    artifact_digest: str = ""
    staging_path: str = ""
    checks: tuple[dict[str, object], ...] = ()
    events: tuple[EvolutionEvent, ...] = ()
    error: str = ""


class EvolutionModuleManager:
    """Stage local evolution module packages without activating or executing them."""

    def __init__(
        self,
        workspace: Path,
        ledger: EvolutionLedger | None = None,
        lock_path: Path | None = None,
        config_loader: Callable[[], Any] | None = None,
        config_saver: Callable[[Any], None] | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self.ledger = ledger or EvolutionLedger(self.workspace)
        self._config_loader = config_loader
        self._config_saver = config_saver
        memory_dir = self.workspace / "memory"
        self.staging_root = memory_dir / "evolution_staging"
        self._lock_path = Path(lock_path) if lock_path is not None else memory_dir / ".evolution_manager.lock"
        with self._locked():
            self._cleanup_stale_tmp_dirs_unlocked()

    def stage(self, source_path: str | Path, *, actor: str = "user") -> EvolutionStageResult:
        """Validate and copy a local module package into the staging store."""

        source_name = _source_name(source_path)
        relative_staging_path = ""
        events: list[EvolutionEvent] = []
        package: EvolutionPackage | None = None
        try:
            with self._locked():
                package = load_package(source_path)
                source_name = package.source_name
                manifest = package.manifest
                target_dir = self.staging_root / package.artifact_digest
                relative_staging_path = target_dir.relative_to(self.workspace).as_posix()
                result_base = {
                    "source_name": source_name,
                    "staging_path": relative_staging_path,
                }
                common = {
                    "actor": actor,
                    "module_id": manifest.module_id,
                    "module_version": manifest.version,
                    "module_type": manifest.module_type,
                    "artifact_digest": package.artifact_digest,
                }
                dirty_activation = self._dirty_activation_for_module_id_unlocked(manifest.module_id)
                if dirty_activation is not None:
                    error = "module_id is blocked by dirty rollback"
                    failed = self.ledger.append(
                        EvolutionEvent.new(
                            EventType.MODULE_FAILED,
                            **common,
                            result={
                                **result_base,
                                "status": "dirty_rollback_blocked",
                                "error": error,
                                "dirty_artifact_digest": str(dirty_activation.get("artifact_digest") or ""),
                            },
                        )
                    )
                    events.append(failed)
                    return EvolutionStageResult(
                        ok=False,
                        status="dirty_rollback_blocked",
                        module_id=manifest.module_id,
                        module_type=manifest.module_type,
                        module_version=manifest.version,
                        artifact_digest=package.artifact_digest,
                        staging_path=relative_staging_path,
                        events=tuple(events),
                        error=error,
                    )
                proposed = self.ledger.append(
                    EvolutionEvent.new(
                        EventType.MODULE_PROPOSED,
                        **common,
                        result={**result_base, "status": "proposed"},
                    )
                )
                events.append(proposed)
                validated = self.ledger.append(
                    EvolutionEvent.new(
                        EventType.MODULE_MANIFEST_VALIDATED,
                        **common,
                        result={**result_base, "status": "validated"},
                    )
                )
                events.append(validated)

                existing = self._has_valid_staging_unlocked(package.artifact_digest)
                status: StageStatus = "already_staged" if existing else "staged"
                installed_event = EvolutionEvent.new(
                    EventType.MODULE_INSTALLED_STAGING,
                    **common,
                    result={**result_base, "status": status},
                )

                if existing:
                    installed = self.ledger.append(installed_event)
                    events.append(installed)
                    return _stage_result(
                        ok=True,
                        status="already_staged",
                        package=package,
                        staging_path=relative_staging_path,
                        events=events,
                    )

                if target_dir.exists():
                    _cleanup_path(target_dir)
                self._stage_package_unlocked(
                    package,
                    target_dir,
                    ledger_event_id=installed_event.event_id,
                )
                try:
                    installed = self.ledger.append(installed_event)
                except Exception:
                    _cleanup_path(target_dir)
                    raise
                events.append(installed)
                return _stage_result(
                    ok=True,
                    status="staged",
                    package=package,
                    staging_path=relative_staging_path,
                    events=events,
                )
        except Exception as exc:
            public_error = _safe_error(exc, source_path, self.workspace)
            failed = self.ledger.append(
                EvolutionEvent.new(
                    EventType.MODULE_FAILED,
                    actor=actor,
                    module_id=package.manifest.module_id if package is not None else "",
                    module_version=package.manifest.version if package is not None else "",
                    module_type=package.manifest.module_type if package is not None else "",
                    artifact_digest=package.artifact_digest if package is not None else "",
                    result={
                        "source_name": source_name,
                        "staging_path": relative_staging_path,
                        "status": "failed",
                        "error": public_error,
                    },
                )
            )
            events.append(failed)
            return EvolutionStageResult(
                ok=False,
                status="failed",
                module_id=package.manifest.module_id if package is not None else "",
                module_type=package.manifest.module_type if package is not None else "",
                module_version=package.manifest.version if package is not None else "",
                artifact_digest=package.artifact_digest if package is not None else "",
                events=tuple(events),
                error=public_error,
            )

    def verify(self, artifact_digest: str, *, actor: str = "user") -> EvolutionVerificationResult:
        """Run static verification against a staged module artifact."""

        with self._locked():
            report = EvolutionModuleVerifier(
                self.workspace,
                staging_root=self.staging_root,
            ).verify(artifact_digest)
            common = {
                "actor": actor,
                "module_id": report.module_id,
                "module_version": report.module_version,
                "module_type": report.module_type,
                "artifact_digest": artifact_digest,
            }
            events: list[EvolutionEvent] = []
            permission_checked = self.ledger.append(
                EvolutionEvent.new(
                    EventType.MODULE_PERMISSION_CHECKED,
                    **common,
                    result={
                        "status": "permission_check_completed",
                        "staging_path": report.staging_path,
                        "permissions_evaluated": report.permissions_evaluated or {},
                        "permissions_denied": list(report.permissions_denied),
                        "unknown_keys_rejected": list(report.unknown_keys_rejected),
                        "checks": list(report.checks),
                    },
                )
            )
            events.append(permission_checked)
            terminal_event_type = EventType.MODULE_VERIFIED if report.ok else EventType.MODULE_FAILED
            terminal_status = "verified" if report.ok else "failed"
            terminal = self.ledger.append(
                EvolutionEvent.new(
                    terminal_event_type,
                    **common,
                    result={
                        "status": terminal_status,
                        "staging_path": report.staging_path,
                        "checks": list(report.checks),
                        "error": report.error,
                    },
                )
            )
            events.append(terminal)
            return _verification_result(report, events)

    def create_state_branch(
        self,
        artifact_digest: str,
        *,
        actor: str = "user",
    ) -> EvolutionStateBranchResult:
        return EvolutionStateBranchStore(self.workspace, ledger=self.ledger).create_branch(
            artifact_digest,
            actor=actor,
        )

    def merge_state_branch(
        self,
        branch_id: str,
        *,
        actor: str = "user",
    ) -> EvolutionStateBranchResult:
        return EvolutionStateBranchStore(self.workspace, ledger=self.ledger).merge_branch(
            branch_id,
            actor=actor,
        )

    def discard_state_branch(
        self,
        branch_id: str,
        *,
        actor: str = "user",
    ) -> EvolutionStateBranchResult:
        return EvolutionStateBranchStore(self.workspace, ledger=self.ledger).discard_branch(
            branch_id,
            actor=actor,
        )

    def activate_module(
        self,
        artifact_digest: str,
        *,
        actor: str = "user",
    ) -> EvolutionActivationResult:
        return EvolutionModuleActivator(
            self.workspace,
            ledger=self.ledger,
            config_loader=self._config_loader,
            config_saver=self._config_saver,
        ).activate(artifact_digest, actor=actor)

    def rollback_module(
        self,
        artifact_digest: str,
        *,
        actor: str = "user",
    ) -> EvolutionActivationResult:
        return EvolutionModuleActivator(
            self.workspace,
            ledger=self.ledger,
            config_loader=self._config_loader,
            config_saver=self._config_saver,
        ).rollback(artifact_digest, actor=actor)

    def capability_snapshot(
        self,
        artifact_digest: str,
        *,
        base_snapshot: CapabilitySnapshot | None = None,
    ) -> EvolutionCapabilityResult:
        return EvolutionCapabilityGate(self.workspace, ledger=self.ledger).snapshot_for_artifact(
            artifact_digest,
            base_snapshot=base_snapshot,
        )

    def record_telemetry(
        self,
        artifact_digest: str,
        *,
        event_kind: str,
        status: str,
        actor: str = "user",
        **kwargs: Any,
    ) -> EvolutionTelemetryResult:
        return EvolutionTelemetryRecorder(self.workspace, ledger=self.ledger).record(
            artifact_digest,
            event_kind=event_kind,
            status=status,
            actor=actor,
            **kwargs,
        )

    def preflight_token_budget(
        self,
        artifact_digest: str,
        *,
        payload_texts: Iterable[str] = (),
        estimated_tokens: int | None = None,
        actor: str = "user",
    ) -> EvolutionTokenBudgetResult:
        return EvolutionTelemetryRecorder(self.workspace, ledger=self.ledger).preflight_token_budget(
            artifact_digest,
            payload_texts=payload_texts,
            estimated_tokens=estimated_tokens,
            actor=actor,
        )

    def record_postflight_usage(
        self,
        artifact_digest: str,
        *,
        budget_token: str,
        usage: Mapping[str, Any],
        actor: str = "user",
    ) -> EvolutionTokenBudgetResult:
        return EvolutionTelemetryRecorder(self.workspace, ledger=self.ledger).record_postflight_usage(
            artifact_digest,
            budget_token=budget_token,
            usage=usage,
            actor=actor,
        )

    def build_proof_bundle(
        self,
        artifact_digest: str,
        *,
        actor: str = "user",
    ) -> EvolutionProofBundleResult:
        return EvolutionTelemetryRecorder(self.workspace, ledger=self.ledger).build_proof_bundle(
            artifact_digest,
            actor=actor,
        )

    def record_teardown(
        self,
        artifact_digest: str,
        *,
        succeeded: bool,
        reason: str = "",
        residual_resources: Iterable[str] = (),
        actor: str = "user",
    ) -> EvolutionRecoveryResult:
        return EvolutionRecoveryManager(self.workspace, ledger=self.ledger).record_teardown(
            artifact_digest,
            succeeded=succeeded,
            reason=reason,
            residual_resources=residual_resources,
            actor=actor,
        )

    def force_clean_module(
        self,
        artifact_digest: str,
        *,
        reason: str,
        actor: str = "user",
    ) -> EvolutionRecoveryResult:
        return EvolutionRecoveryManager(self.workspace, ledger=self.ledger).force_clean(
            artifact_digest,
            reason=reason,
            actor=actor,
        )

    def runtime_status(self, verify_signatures: bool = False) -> LedgerStatus:
        return self.ledger.status(verify_signatures=verify_signatures)

    def _locked(self) -> FileLock:
        self.staging_root.parent.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))

    def _cleanup_stale_tmp_dirs_unlocked(self) -> None:
        if not self.staging_root.exists():
            return
        for child in self.staging_root.iterdir():
            if child.name.startswith(".tmp-"):
                _cleanup_path(child)

    def _has_valid_staging_unlocked(self, artifact_digest: str) -> bool:
        target_dir = self.staging_root / artifact_digest
        if not target_dir.exists():
            return False
        staging_json = target_dir / "staging.json"
        if not staging_json.exists():
            _cleanup_path(target_dir)
            return False
        try:
            data = json.loads(staging_json.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            _cleanup_path(target_dir)
            return False
        if data.get("artifact_digest") != artifact_digest:
            _cleanup_path(target_dir)
            return False
        if not (target_dir / "artifact").is_dir():
            _cleanup_path(target_dir)
            return False
        return True

    def _dirty_activation_for_module_id_unlocked(self, module_id: str) -> dict[str, Any] | None:
        activation_root = self.workspace / "memory" / "evolution_activations"
        if not activation_root.exists():
            return None
        for activation_json in activation_root.glob("*/activation.json"):
            try:
                data = json.loads(activation_json.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError):
                continue
            if (
                isinstance(data, dict)
                and data.get("module_id") == module_id
                and data.get("status") == "dirty_rollback"
            ):
                return data
        return None

    def _stage_package_unlocked(
        self,
        package: EvolutionPackage,
        target_dir: Path,
        *,
        ledger_event_id: str,
    ) -> None:
        self.staging_root.mkdir(parents=True, exist_ok=True)
        tmp_dir = self.staging_root / f".tmp-{uuid.uuid4().hex}"
        _cleanup_path(tmp_dir)
        try:
            copy_artifact(package.source_path, tmp_dir / "artifact")
            metadata = {
                "schema_version": STAGING_SCHEMA_VERSION,
                "module_id": package.manifest.module_id,
                "module_type": package.manifest.module_type,
                "version": package.manifest.version,
                "artifact_digest": package.artifact_digest,
                "staged_at": datetime.now(timezone.utc).isoformat(),
                "ledger_event_id": ledger_event_id,
            }
            (tmp_dir / "staging.json").write_bytes(canonical_dump(metadata) + b"\n")
            tmp_dir.rename(target_dir)
        except Exception:
            _cleanup_path(tmp_dir)
            raise


def _stage_result(
    *,
    ok: bool,
    status: StageStatus,
    package: EvolutionPackage,
    staging_path: str,
    events: list[EvolutionEvent],
) -> EvolutionStageResult:
    return EvolutionStageResult(
        ok=ok,
        status=status,
        module_id=package.manifest.module_id,
        module_type=package.manifest.module_type,
        module_version=package.manifest.version,
        artifact_digest=package.artifact_digest,
        staging_path=staging_path,
        events=tuple(events),
    )


def _verification_result(
    report: EvolutionVerificationReport,
    events: list[EvolutionEvent],
) -> EvolutionVerificationResult:
    return EvolutionVerificationResult(
        ok=report.ok,
        status=report.status,
        module_id=report.module_id,
        module_type=report.module_type,
        module_version=report.module_version,
        artifact_digest=report.artifact_digest,
        staging_path=report.staging_path,
        checks=report.checks,
        events=tuple(events),
        error=report.error,
    )


def _source_name(source_path: str | Path) -> str:
    return Path(source_path).expanduser().resolve(strict=False).name


def _safe_error(exc: Exception, source_path: str | Path, workspace: Path) -> str:
    message = str(exc)
    for sensitive in (
        str(Path(source_path).expanduser().resolve(strict=False)),
        str(Path(workspace).expanduser().resolve(strict=False)),
    ):
        if sensitive:
            message = message.replace(sensitive, "<path>")
    return message


def _cleanup_path(path: Path) -> None:
    if path.is_dir() and not path.is_symlink():
        shutil.rmtree(path, ignore_errors=True)
    else:
        try:
            path.unlink()
        except FileNotFoundError:
            pass
