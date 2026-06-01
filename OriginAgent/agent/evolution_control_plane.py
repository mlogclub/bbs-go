"""Unified control plane for governed evolution operations."""

from __future__ import annotations

from dataclasses import asdict, dataclass, field
from pathlib import Path
from typing import Any

from OriginAgent.agent.evolution import AUTO_EVOLUTION_ORIGIN, OpportunitySignalStore
from OriginAgent.agent.evolution_config_overlay import EvolutionConfigOverlayStore, apply_config_overlay
from OriginAgent.agent.evolution_dependencies import EvolutionDependencyStore
from OriginAgent.agent.evolution_feedback import EvolutionFeedbackCalibrator, feedback_status
from OriginAgent.agent.evolution_health import evolution_health_score
from OriginAgent.agent.evolution_health_history import (
    EvolutionHealthHistoryStore,
    health_history_policy_status,
)
from OriginAgent.agent.evolution_maintenance import run_evolution_maintenance
from OriginAgent.agent.evolution_operator import EvolutionOperator
from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore, safe_append_outcome
from OriginAgent.agent.evolution_sandbox import sandbox_status_counts, trial_policy_status
from OriginAgent.agent.evolution_schema import validate_evolution_stores
from OriginAgent.agent.evolution_snapshots import EvolutionRollbackService, EvolutionSnapshotStore
from OriginAgent.agent.evolution_trial_logs import EvolutionTrialLogStore, trial_log_policy_status

EVOLUTION_MANUAL_OVERRIDE_DISABLED = (
    "Evolution manual override is disabled. Set evolution.allow_manual_override=true in config to enable."
)

READ_ACTIONS = frozenset({
    "status",
    "list_actions",
    "list_signals",
    "list_proposals",
    "inspect_signal",
    "inspect_evolution_proposal",
    "inspect_proposal",
    "explain_evolution_health",
    "list_evolution_recommendations",
    "list_recommendations",
    "generate_evolution_report",
    "validate_schema",
    "list_config_overlay",
})
PREVIEW_ACTIONS = frozenset({
    "preview_evolution_action",
    "suppress_signal",
    "resume_signal",
    "run_maintenance",
    "force_cleanup",
    "run_feedback_calibration",
    "retry_trial",
    "rollback_artifact",
    "validate_schema",
    "generate_evolution_report",
    "clear_config_overlay",
})
WRITE_ACTIONS = frozenset({
    "suppress_signal",
    "resume_signal",
    "run_maintenance",
    "force_cleanup",
    "run_feedback_calibration",
    "retry_trial",
    "rollback_artifact",
    "clear_config_overlay",
})

ACTION_SCHEMA_VERSION = "originagent.evolution.action.v1"
ACTION_RESULT_SCHEMA_VERSION = "originagent.evolution.action_result.v1"
CONTROL_EVENT_EXECUTED = "control_action_executed"
CONTROL_EVENT_DENIED = "control_action_denied"
CONTROL_EVENT_FAILED = "control_action_failed"


@dataclass(frozen=True)
class EvolutionPolicyDecision:
    """One centralized allow/deny decision for an evolution control-plane action."""

    allowed: bool
    action_kind: str
    permission: str
    reason: str = ""
    requires_manual_override: bool = False
    dry_run: bool = True
    mode: str = "conservative"

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class EvolutionActionDescriptor:
    """Stable UI/API contract for one governed evolution operator action."""

    action_id: str
    action_kind: str
    target_type: str
    target_id: str = ""
    permission: str = "read"
    risk_level: str = "low"
    previewable: bool = True
    executable: bool = True
    requires_manual_override: bool = False
    summary: str = ""
    source: str = "control_plane"
    parameters_schema: dict[str, Any] = field(default_factory=dict)
    suggested_my_action: str = ""
    policy: dict[str, Any] = field(default_factory=dict)
    schema_version: str = ACTION_SCHEMA_VERSION

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class EvolutionPolicy:
    """Centralized policy decisions for the evolution control plane."""

    def __init__(self, config: Any | None = None) -> None:
        self.config = config

    def decide(self, action_kind: str, *, execution: bool = False) -> EvolutionPolicyDecision:
        action = normalize_action_kind(action_kind)
        permission = action_permission(action)
        dry_run = bool(getattr(self.config, "dry_run", True) if self.config is not None else True)
        mode = str(getattr(self.config, "mode", "conservative") if self.config is not None else "conservative")
        if permission == "unknown":
            return EvolutionPolicyDecision(
                allowed=False,
                action_kind=action,
                permission=permission,
                reason=f"Unsupported evolution control-plane action `{action_kind}`.",
                requires_manual_override=False,
                dry_run=dry_run,
                mode=mode,
            )
        requires_override = permission in {"maintenance", "override", "rollback"}
        if execution and requires_override and not bool(getattr(self.config, "allow_manual_override", False)):
            return EvolutionPolicyDecision(
                allowed=False,
                action_kind=action,
                permission=permission,
                reason=EVOLUTION_MANUAL_OVERRIDE_DISABLED,
                requires_manual_override=True,
                dry_run=dry_run,
                mode=mode,
            )
        return EvolutionPolicyDecision(
            allowed=True,
            action_kind=action,
            permission=permission,
            reason="allowed",
            requires_manual_override=requires_override,
            dry_run=dry_run,
            mode=mode,
        )


class EvolutionControlPlane:
    """Read, preview, and guarded write API for governed evolution."""

    def __init__(self, workspace: Path, config: Any | None = None) -> None:
        self.workspace = Path(workspace)
        self.raw_config = config
        self.config = apply_config_overlay(self.workspace, config)
        self.policy = EvolutionPolicy(config)
        self.operator = EvolutionOperator(self.workspace, self.config)

    def status(self) -> dict[str, Any]:
        """Return the unified dashboard-ready evolution read model."""

        try:
            return self._status_unchecked()
        except Exception:
            return self._fallback_status()

    def list_signals(
        self,
        *,
        status: str | None = None,
        kind: str | None = None,
        limit: int = 50,
    ) -> dict[str, Any]:
        records = [signal.to_record() for signal in OpportunitySignalStore(self.workspace).read_all()]
        if status:
            normalized = str(status).strip().lower()
            records = [record for record in records if str(record.get("status") or "").lower() == normalized]
        if kind:
            normalized_kind = str(kind).strip().lower()
            records = [record for record in records if str(record.get("kind") or "").lower() == normalized_kind]
        records.sort(key=lambda record: (float(record.get("priority_score") or 0.0), str(record.get("last_seen_at") or "")), reverse=True)
        effective_limit = max(1, min(int(limit or 50), 200))
        return {
            "ok": True,
            "count": len(records),
            "signals": records[:effective_limit],
        }

    def list_proposals(
        self,
        *,
        status: str | None = None,
        proposal_type: str | None = None,
        limit: int = 50,
    ) -> dict[str, Any]:
        from OriginAgent.agent.background_review import ReviewProposalStore

        records = ReviewProposalStore(self.workspace).list_records(
            origin=AUTO_EVOLUTION_ORIGIN,
            status=status,
            proposal_type=proposal_type,
            limit=limit,
        )
        summaries: list[dict[str, Any]] = []
        for record in records:
            inspected = self.operator.inspect_proposal(str(record.get("id") or ""))
            summary = inspected.get("summary") if isinstance(inspected.get("summary"), dict) else {}
            summaries.append(summary or {
                "proposal_id": str(record.get("id") or ""),
                "status": str(record.get("status") or "pending"),
                "proposal_type": str(record.get("proposal_type") or ""),
                "title": str(record.get("title") or ""),
            })
        return {
            "ok": True,
            "count": len(summaries),
            "proposals": summaries,
        }

    def inspect_signal(self, opportunity_id: str) -> dict[str, Any]:
        return self.operator.inspect_signal(opportunity_id)

    def inspect_proposal(self, proposal_id: str) -> dict[str, Any]:
        return self.operator.inspect_proposal(proposal_id)

    def list_actions(self) -> dict[str, Any]:
        actions = []
        for action in sorted({normalize_action_kind(item) for item in READ_ACTIONS | WRITE_ACTIONS}):
            if action == "preview_evolution_action":
                continue
            descriptor = self.action_descriptor(action)
            item = descriptor.to_json()
            item["can_preview"] = item["previewable"]
            item["can_execute"] = item["executable"]
            actions.append(item)
        return {
            "ok": True,
            "schema_version": ACTION_SCHEMA_VERSION,
            "actions": actions,
            "safety_boundaries": safety_boundaries(),
        }

    def explain_health(self) -> dict[str, Any]:
        return self.operator.explain_health()

    def list_recommendations(
        self,
        *,
        health: dict[str, Any] | None = None,
        health_history: dict[str, Any] | None = None,
        outcome_stats: dict[str, Any] | None = None,
        dependency_stats: dict[str, Any] | None = None,
        feedback_stats: dict[str, Any] | None = None,
        sandbox_counts: dict[str, int] | None = None,
    ) -> list[dict[str, Any]]:
        return [
            self._recommendation_with_action(item)
            for item in self.operator.list_recommendations(
                health=health,
                health_history=health_history,
                outcome_stats=outcome_stats,
                dependency_stats=dependency_stats,
                feedback_stats=feedback_stats,
                sandbox_counts=sandbox_counts,
            )
        ]

    def generate_report(self, *, period_days: int = 7) -> str:
        return self.operator.generate_report(period_days=period_days)

    def action_descriptor(
        self,
        action_kind: str,
        *,
        target_id: str = "",
        risk_level: str = "",
        source: str = "control_plane",
    ) -> EvolutionActionDescriptor:
        """Return the stable action contract consumed by tools, CLIs, and UIs."""

        action = normalize_action_kind(action_kind)
        decision = self.policy.decide(action, execution=False)
        return build_action_descriptor(
            action,
            target_id=target_id,
            decision=decision,
            risk_level=str(risk_level or risk_level_for_action(action)),
            source=source,
        )

    def preview_action(
        self,
        action_kind: str,
        *,
        target_id: str = "",
        reason: str = "",
        fixtures: dict[str, str] | None = None,
        force_cleanup: bool = False,
        artifact_type: str = "",
        artifact_name: str = "",
        snapshot_id: str = "",
        period_days: int = 7,
    ) -> dict[str, Any]:
        action = normalize_action_kind(action_kind)
        decision = self.policy.decide(action, execution=False)
        if action == "rollback_artifact":
            result = self._preview_rollback(
                artifact_type=artifact_type,
                artifact_name=artifact_name or target_id,
                snapshot_id=snapshot_id,
                reason=reason,
            )
        elif action == "validate_schema":
            result = self._action_result(
                action_kind=action,
                target_type="schema",
                target_id="",
                will_write=False,
                preview={"schema_validation": validate_evolution_stores(self.workspace)},
            )
        elif action == "generate_evolution_report":
            result = self._action_result(
                action_kind=action,
                target_type="report",
                target_id="",
                will_write=False,
                preview={"period_days": _coerce_period_days(period_days), "would_generate_markdown": True},
            )
        elif action == "clear_config_overlay":
            result = self._action_result(
                action_kind=action,
                target_type="config_overlay",
                target_id="evolution_config",
                will_write=False,
                preview={
                    "would_clear_overlay": True,
                    "current_overlay": EvolutionConfigOverlayStore(self.workspace).status(),
                },
                message="Would clear governed evolution config overlay.",
            )
        else:
            result = self.operator.preview_action(
                action,
                target_id=target_id,
                reason=reason,
                fixtures=fixtures or {},
                force_cleanup=force_cleanup,
            )
        return self._with_policy(result, decision, will_execute=False)

    def execute_action(
        self,
        action_kind: str,
        *,
        target_id: str = "",
        reason: str = "",
        fixtures: dict[str, str] | None = None,
        force_cleanup: bool = False,
        artifact_type: str = "",
        artifact_name: str = "",
        snapshot_id: str = "",
        period_days: int = 7,
        actor: str = "control_plane",
        source: str = "control_plane",
    ) -> dict[str, Any]:
        action = normalize_action_kind(action_kind)
        decision = self.policy.decide(action, execution=True)
        if not decision.allowed:
            error = "unsupported_action" if decision.permission == "unknown" else "manual_override_disabled"
            result = self._with_policy(
                self._action_result(
                    ok=False,
                    action_kind=action,
                    target_type=target_type_for_action(action),
                    target_id=target_id or artifact_name,
                    will_write=False,
                    error=error,
                    message=decision.reason,
                ),
                decision,
                will_execute=True,
            )
            self._append_control_outcome(
                CONTROL_EVENT_DENIED,
                result,
                decision,
                actor=actor,
                source=source,
            )
            return result
        if action in READ_ACTIONS and action not in WRITE_ACTIONS:
            return self._with_policy(
                self._execute_read_action(action, target_id=target_id, period_days=period_days),
                decision,
                will_execute=True,
            )
        if action == "suppress_signal":
            signal = OpportunitySignalStore(self.workspace).suppress_signal(target_id, reason=reason)
            result = self._action_result(
                ok=signal is not None,
                action_kind=action,
                target_type="signal",
                target_id=target_id,
                will_write=True,
                result=signal.to_record() if signal is not None else None,
                error="" if signal is not None else "signal_not_found",
                message="Signal suppressed." if signal is not None else "Opportunity signal was not found or cannot be suppressed.",
            )
        elif action == "resume_signal":
            signal = OpportunitySignalStore(self.workspace).resume_signal(target_id, reason=reason)
            result = self._action_result(
                ok=signal is not None,
                action_kind=action,
                target_type="signal",
                target_id=target_id,
                will_write=True,
                result=signal.to_record() if signal is not None else None,
                error="" if signal is not None else "signal_not_found",
                message="Signal resumed." if signal is not None else "Opportunity signal was not found or is not suppressed.",
            )
        elif action in {"run_maintenance", "force_cleanup"}:
            result = self._action_result(
                action_kind=action,
                target_type="maintenance",
                target_id="",
                will_write=True,
                result=run_evolution_maintenance(
                    self.workspace,
                    self.config,
                    force_cleanup=force_cleanup or action == "force_cleanup",
                ),
            )
        elif action == "run_feedback_calibration":
            result = self._action_result(
                action_kind=action,
                target_type="feedback",
                target_id="",
                will_write=True,
                result=EvolutionFeedbackCalibrator(self.workspace, self.config).run().to_json(),
            )
        elif action == "retry_trial":
            retry = self.operator.retry_trial(target_id, fixtures=fixtures or {}, actor=actor or "control_plane")
            result = self._action_result(
                ok=retry.ok,
                action_kind=action,
                target_type="proposal",
                target_id=target_id,
                will_write=True,
                result=retry.to_json(),
                error=retry.error,
                message=retry.message,
            )
        elif action == "rollback_artifact":
            rollback = EvolutionRollbackService(self.workspace).rollback(
                artifact_type=artifact_type,
                artifact_name=artifact_name or target_id,
                snapshot_id=snapshot_id or None,
                reason=reason,
                actor=actor or "control_plane",
                force=force_cleanup,
            )
            result = self._action_result(
                ok=rollback.ok,
                action_kind=action,
                target_type="artifact",
                target_id=artifact_name or target_id,
                will_write=True,
                result=rollback.to_json(),
                error=rollback.error,
                message=rollback.message,
            )
        elif action == "clear_config_overlay":
            overlay = EvolutionConfigOverlayStore(self.workspace).clear(actor=actor or "control_plane", source=source)
            result = self._action_result(
                action_kind=action,
                target_type="config_overlay",
                target_id="evolution_config",
                will_write=True,
                result=overlay,
                message="Governed evolution config overlay cleared.",
            )
        else:
            result = self._action_result(
                ok=False,
                action_kind=action,
                target_type=target_type_for_action(action),
                target_id=target_id,
                will_write=False,
                error="unsupported_action",
                message=f"Unsupported evolution control-plane action `{action_kind}`.",
            )
        result = self._with_policy(result, decision, will_execute=True)
        self._append_control_outcome(
            CONTROL_EVENT_EXECUTED if result.get("ok") else CONTROL_EVENT_FAILED,
            result,
            decision,
            actor=actor,
            source=source,
        )
        return result

    def _status_unchecked(self) -> dict[str, Any]:
        signal_status = OpportunitySignalStore(self.workspace).runtime_status(self.config)
        from OriginAgent.agent.background_review import ReviewProposalStore

        proposal_store = ReviewProposalStore(self.workspace)
        proposal_stats = proposal_store.stats(origin=AUTO_EVOLUTION_ORIGIN)
        outcome_stats = EvolutionOutcomeStore(self.workspace).stats()
        snapshot_stats = EvolutionSnapshotStore(self.workspace).stats()
        dependency_stats = EvolutionDependencyStore(self.workspace).stats()
        feedback_stats = feedback_status(self.workspace, self.config)
        sandbox_counts = sandbox_status_counts(self.workspace)
        trial_status = trial_policy_status(self.config)
        health = evolution_health_score(
            outcome_stats=outcome_stats,
            dependency_stats=dependency_stats,
            feedback_stats=feedback_stats,
            sandbox_counts=sandbox_counts,
            trial_status=trial_status,
        )
        health_history = EvolutionHealthHistoryStore(self.workspace).summary()
        pending_records = proposal_store.list_records(origin=AUTO_EVOLUTION_ORIGIN, status="pending", limit=50)
        recent_records = proposal_store.list_records(origin=AUTO_EVOLUTION_ORIGIN, limit=50)
        applied_records = proposal_store.list_records(origin=AUTO_EVOLUTION_ORIGIN, status="applied", limit=50)
        issue_counts = _static_gate_issue_counts(pending_records)
        promotion_gate_counts = _promotion_gate_counts(recent_records)
        auto_verified = _auto_verified_workflow_count(applied_records)
        trial_log_store = EvolutionTrialLogStore(self.workspace)
        schema_validation = validate_evolution_stores(self.workspace)
        config_overlay = EvolutionConfigOverlayStore(self.workspace).status()
        return {
            **signal_status,
            "control_plane": {
                "version": "originagent.evolution.control_plane.v2",
                "actions_count": len(READ_ACTIONS | WRITE_ACTIONS),
                "manual_override_enabled": bool(getattr(self.config, "allow_manual_override", False) if self.config is not None else False),
                "safety_boundaries": safety_boundaries(),
            },
            "policy": {
                "mode": str(getattr(self.config, "mode", "conservative") if self.config is not None else "conservative"),
                "dry_run": bool(getattr(self.config, "dry_run", True) if self.config is not None else True),
                "allow_manual_override": bool(getattr(self.config, "allow_manual_override", False) if self.config is not None else False),
                "overlay_active": bool(config_overlay.get("active")),
                "permissions": permission_summary(self.config),
            },
            "config_overlay": config_overlay,
            "pending_proposals_from_evolution": proposal_stats["pending_count"],
            "proposal_count_from_evolution": proposal_stats["proposal_count"],
            "auto_verified_workflows_count": auto_verified,
            "maintenance": evolution_maintenance_policy(self.config),
            "outcomes": outcome_stats,
            "snapshots": snapshot_stats,
            "dependencies": dependency_stats,
            "feedback_calibration": feedback_stats,
            "evolution_health": health,
            "evolution_health_history": {
                **health_history_policy_status(self.config),
                **health_history,
            },
            "operator_recommendations": self.list_recommendations(
                health=health,
                health_history=health_history,
                outcome_stats=outcome_stats,
                dependency_stats=dependency_stats,
                feedback_stats=feedback_stats,
                sandbox_counts=sandbox_counts,
            ),
            "promotion_gate_decision_counts": promotion_gate_counts,
            "static_gate_issue_counts": issue_counts,
            "sandbox": {
                "enabled": bool(getattr(getattr(self.config, "sandbox", None), "enabled", True)),
                "passed_workflow_proposals": sandbox_counts.get("passed", 0),
                "failed_workflow_proposals": sandbox_counts.get("failed", 0),
                "blocked_workflow_proposals": sandbox_counts.get("blocked", 0),
            },
            "trial": trial_status,
            "trial_logs": {
                **trial_log_policy_status(self.config),
                **trial_log_store.stats(),
            },
            "schema_validation": {
                "ok": schema_validation.get("ok"),
                "record_counts": schema_validation.get("record_counts", {}),
                "issue_counts": schema_validation.get("issue_counts", {}),
            },
            "read_model": {
                "signals": {
                    "open": signal_status.get("opportunity_signals_count", 0),
                    "converted": signal_status.get("converted_signals_count", 0),
                    "suppressed": signal_status.get("suppressed_signals_count", 0),
                    "eligible_workflow": signal_status.get("eligible_workflow_signals", 0),
                    "eligible_skill": signal_status.get("eligible_skill_signals", 0),
                },
                "proposals": {
                    "total": proposal_stats["proposal_count"],
                    "pending": proposal_stats["pending_count"],
                    "auto_verified_workflows": auto_verified,
                },
                "health": health,
                "dependencies": dependency_stats,
                "snapshots": snapshot_stats,
            },
        }

    def _execute_read_action(self, action: str, *, target_id: str, period_days: int) -> dict[str, Any]:
        if action == "status":
            result = self.status()
        elif action == "list_signals":
            result = self.list_signals()
        elif action == "list_proposals":
            result = self.list_proposals()
        elif action == "list_actions":
            result = self.list_actions()
        elif action == "inspect_signal":
            result = self.inspect_signal(target_id)
        elif action in {"inspect_evolution_proposal", "inspect_proposal"}:
            result = self.inspect_proposal(target_id)
        elif action == "explain_evolution_health":
            result = self.explain_health()
        elif action in {"list_evolution_recommendations", "list_recommendations"}:
            result = self.list_recommendations()
        elif action == "generate_evolution_report":
            result = self.generate_report(period_days=period_days)
        elif action == "validate_schema":
            result = validate_evolution_stores(self.workspace)
        elif action == "list_config_overlay":
            result = EvolutionConfigOverlayStore(self.workspace).status()
        else:
            return self._action_result(
                ok=False,
                action_kind=action,
                target_type=target_type_for_action(action),
                target_id=target_id,
                will_write=False,
                error="unsupported_action",
            )
        return self._action_result(
            action_kind=action,
            target_type=target_type_for_action(action),
            target_id=target_id,
            will_write=False,
            result=result,
        )

    def _preview_rollback(
        self,
        *,
        artifact_type: str,
        artifact_name: str,
        snapshot_id: str,
        reason: str,
    ) -> dict[str, Any]:
        artifact_type = str(artifact_type or "").strip().lower()
        artifact_name = str(artifact_name or "").strip()
        snapshots = EvolutionSnapshotStore(self.workspace)
        snapshot_records = snapshots.list_snapshots(
            artifact_type=artifact_type or None,
            artifact_name=artifact_name or None,
        )
        selected = None
        if snapshot_id:
            for record in snapshot_records:
                if record.get("snapshot_id") == snapshot_id or record.get("content_hash") == snapshot_id:
                    selected = record
                    break
        elif artifact_type and artifact_name:
            selected = snapshots.latest_snapshot(artifact_type=artifact_type, artifact_name=artifact_name)
        blockers = []
        if artifact_type and artifact_name:
            blockers = EvolutionDependencyStore(self.workspace).rollback_blockers(
                artifact_type=artifact_type,
                artifact_name=artifact_name,
            )
        ok = bool(artifact_type in {"workflow", "skill"} and artifact_name and selected is not None)
        return self._action_result(
            ok=ok,
            action_kind="rollback_artifact",
            target_type="artifact",
            target_id=artifact_name,
            will_write=False,
            error="" if ok else "rollback_preview_unavailable",
            message=(
                "Rollback preview ready."
                if ok
                else "Rollback preview requires artifact_type, artifact_name, and an available snapshot."
            ),
            preview={
                "artifact_type": artifact_type,
                "artifact_name": artifact_name,
                "snapshot_id": str((selected or {}).get("snapshot_id") or snapshot_id or ""),
                "snapshot_path": str((selected or {}).get("snapshot_path") or ""),
                "content_hash": str((selected or {}).get("content_hash") or ""),
                "dependency_blockers": blockers,
                "would_write_artifact": ok,
                "would_write_outcome": ok,
                "reason": reason,
            },
        )

    def _fallback_status(self) -> dict[str, Any]:
        mode = str(getattr(self.config, "mode", "conservative") if self.config is not None else "conservative")
        dry_run = bool(getattr(self.config, "dry_run", True) if self.config is not None else True)
        return {
            "mode": mode,
            "dry_run": dry_run,
            "control_plane": {
                "version": "originagent.evolution.control_plane.v2",
                "manual_override_enabled": bool(getattr(self.config, "allow_manual_override", False) if self.config is not None else False),
                "safety_boundaries": safety_boundaries(),
            },
            "policy": {
                "mode": mode,
                "dry_run": dry_run,
                "allow_manual_override": bool(getattr(self.config, "allow_manual_override", False) if self.config is not None else False),
                "overlay_active": False,
                "permissions": permission_summary(self.config),
            },
            "config_overlay": {
                "schema_version": "originagent.evolution.config_overlay.v1",
                "active": False,
                "override_count": 0,
                "overrides": {},
                "patch_count": 0,
                "last_patch_at": "",
                "last_actor": "",
                "recent_patches": [],
                "store_path": "memory/evolution_config_overrides.json",
                "patch_log_path": "memory/evolution_config_patches.jsonl",
            },
            "opportunity_signals_count": 0,
            "converted_signals_count": 0,
            "suppressed_signals_count": 0,
            "feedback_adjusted_signals_count": 0,
            "feedback_negative_signals_count": 0,
            "feedback_positive_signals_count": 0,
            "pending_proposals_from_evolution": 0,
            "proposal_count_from_evolution": 0,
            "auto_verified_workflows_count": 0,
            "maintenance": evolution_maintenance_policy(self.config),
            "outcomes": {
                "outcome_event_count": 0,
                "outcome_type_counts": {},
                "gate_decision_counts": {},
                "sandbox_status_counts": {},
                "review_status_counts": {},
                "promotion_status_counts": {},
                "rollback_status_counts": {},
                "last_outcome_at": None,
                "archive": {"archived_outcome_count": 0, "last_archived_at": None},
            },
            "snapshots": {"snapshot_count": 0, "snapshot_type_counts": {}, "last_snapshot_at": None},
            "dependencies": {
                "tracked_artifacts": 0,
                "dependency_edges": 0,
                "rollback_blocked_artifacts": 0,
                "stale_reference_count": 0,
            },
            "feedback_calibration": {},
            "evolution_health": {
                "score": 100,
                "level": "healthy",
                "reasons": ["+ no successful rollbacks", "+ no dependency conflicts", "+ trial isolation enforced"],
            },
            "evolution_health_history": {
                "health_history_retention_days": 90,
                "max_health_history_snapshots": 100,
                "snapshot_count": 0,
                "latest_score": None,
                "latest_level": None,
                "previous_score": None,
                "score_delta": 0,
                "trend": "unknown",
                "last_snapshot_at": None,
            },
            "operator_recommendations": [],
            "promotion_gate_decision_counts": {},
            "static_gate_issue_counts": {},
            "sandbox": {
                "enabled": True,
                "passed_workflow_proposals": 0,
                "failed_workflow_proposals": 0,
                "blocked_workflow_proposals": 0,
            },
            "trial": {
                "enabled": True,
                "isolated_workspace": True,
                "read_only_tools_only": True,
                "allowed_tools": ["glob", "grep", "read_file"],
                "blocked_tools": ["cron", "edit_file", "exec", "message", "spawn", "write_file"],
                "temp_dir_configured": False,
            },
            "trial_logs": {
                "max_step_output_chars": 2000,
                "max_retained_trial_logs": 10,
                "trial_log_retention_days": 30,
                "trial_log_count": 0,
                "trial_log_status_counts": {},
                "last_trial_at": None,
                "truncated_step_output_count": 0,
            },
            "skill_candidates_enabled": False,
            "eligible_workflow_signals": 0,
            "eligible_skill_signals": 0,
            "high_score_signals": [],
            "schema_validation": {"ok": True, "record_counts": {}, "issue_counts": {}},
            "read_model": {},
        }

    def _recommendation_with_action(self, item: dict[str, Any]) -> dict[str, Any]:
        recommendation = dict(item)
        action = normalize_action_kind(str(recommendation.get("action_kind") or recommendation.get("action") or ""))
        target_id = str(
            recommendation.get("target_id")
            or recommendation.get("proposal_id")
            or recommendation.get("opportunity_id")
            or ""
        )
        risk_level = str(recommendation.get("risk_level") or risk_level_for_action(action))
        descriptor = self.action_descriptor(
            action,
            target_id=target_id,
            risk_level=risk_level,
            source=f"recommendation:{recommendation.get('code') or 'operator'}",
        )
        recommendation["action_descriptor"] = descriptor.to_json()
        recommendation["previewable"] = descriptor.previewable
        recommendation["executable"] = descriptor.executable
        recommendation["requires_manual_override"] = bool(
            recommendation.get("requires_manual_override") or descriptor.requires_manual_override
        )
        recommendation["suggested_next_step"] = (
            "preview_action"
            if descriptor.previewable and descriptor.requires_manual_override
            else "execute_or_inspect"
            if descriptor.executable
            else "inspect"
        )
        return recommendation

    def _append_control_outcome(
        self,
        event_type: str,
        result: dict[str, Any],
        decision: EvolutionPolicyDecision,
        *,
        actor: str,
        source: str,
    ) -> None:
        action = normalize_action_kind(str(result.get("action_kind") or decision.action_kind))
        metadata = {
            "actor": str(actor or "control_plane"),
            "source": str(source or "control_plane"),
            "action_kind": action,
            "target_type": str(result.get("target_type") or target_type_for_action(action)),
            "target_id": str(result.get("target_id") or ""),
            "policy_decision": "allowed" if decision.allowed else "denied",
            "permission": decision.permission,
            "requires_manual_override": decision.requires_manual_override,
            "result_status": "succeeded" if result.get("ok") else "failed",
            "error": str(result.get("error") or ""),
        }
        safe_append_outcome(
            EvolutionOutcomeStore(self.workspace),
            event_type,
            opportunity_id=str(result.get("target_id") or "") if target_type_for_action(action) == "signal" else "",
            proposal_id=str(result.get("target_id") or "") if target_type_for_action(action) == "proposal" else "",
            artifact_type="",
            artifact_name=str(result.get("target_id") or "") if target_type_for_action(action) == "artifact" else "",
            metadata=metadata,
        )

    @staticmethod
    def _action_result(
        *,
        action_kind: str,
        target_type: str,
        target_id: str,
        will_write: bool,
        ok: bool = True,
        preview: Any = None,
        result: Any = None,
        error: str = "",
        message: str = "",
    ) -> dict[str, Any]:
        return {
            "schema_version": ACTION_RESULT_SCHEMA_VERSION,
            "ok": bool(ok),
            "action_kind": normalize_action_kind(action_kind),
            "target_type": target_type,
            "target_id": str(target_id or ""),
            "requires_manual_override": normalize_action_kind(action_kind) in WRITE_ACTIONS,
            "will_write": bool(will_write),
            "preview": preview,
            "result": result,
            "error": error,
            "message": message,
        }

    @staticmethod
    def _with_policy(result: dict[str, Any], decision: EvolutionPolicyDecision, *, will_execute: bool) -> dict[str, Any]:
        merged = dict(result)
        action = normalize_action_kind(str(merged.get("action_kind") or decision.action_kind))
        merged["policy"] = decision.to_json()
        merged["allowed"] = decision.allowed
        merged["schema_version"] = str(merged.get("schema_version") or ACTION_RESULT_SCHEMA_VERSION)
        merged["action"] = build_action_descriptor(
            action,
            target_id=str(merged.get("target_id") or ""),
            decision=decision,
            risk_level=risk_level_for_action(action),
            source="execute" if will_execute else "preview",
        ).to_json()
        merged["requires_manual_override"] = bool(
            merged.get("requires_manual_override") or decision.requires_manual_override
        )
        if will_execute and not decision.allowed:
            merged["ok"] = False
            merged["will_write"] = False
            merged["error"] = merged.get("error") or "policy_denied"
            merged["message"] = merged.get("message") or decision.reason
        return merged


def normalize_action_kind(value: str) -> str:
    action = str(value or "").strip()
    aliases = {
        "inspect_proposal": "inspect_evolution_proposal",
        "inspect_or_retry_trial": "inspect_evolution_proposal",
        "inspect_feedback_adjusted_signals": "inspect_signal",
        "inspect_suppression_candidates": "inspect_signal",
        "consider_suppress_signal": "suppress_signal",
        "consider_conservative_or_dry_run": "explain_evolution_health",
        "list_actions": "list_actions",
        "rollback": "rollback_artifact",
        "validate_evolution_schema": "validate_schema",
    }
    return aliases.get(action, action)


def action_permission(action_kind: str) -> str:
    action = normalize_action_kind(action_kind)
    if action in {"status", "list_actions", "list_signals", "list_proposals", "inspect_signal", "inspect_evolution_proposal", "inspect_proposal", "explain_evolution_health", "list_evolution_recommendations", "list_recommendations", "generate_evolution_report", "validate_schema", "list_config_overlay"}:
        return "read"
    if action in {"run_maintenance", "force_cleanup", "run_feedback_calibration", "clear_config_overlay"}:
        return "maintenance"
    if action in {"suppress_signal", "resume_signal", "retry_trial"}:
        return "override"
    if action == "rollback_artifact":
        return "rollback"
    return "unknown"


def target_type_for_action(action_kind: str) -> str:
    action = normalize_action_kind(action_kind)
    if action in {"inspect_signal", "suppress_signal", "resume_signal"}:
        return "signal"
    if action in {"inspect_evolution_proposal", "inspect_proposal", "retry_trial"}:
        return "proposal"
    if action in {"run_maintenance", "force_cleanup"}:
        return "maintenance"
    if action == "run_feedback_calibration":
        return "feedback"
    if action == "rollback_artifact":
        return "artifact"
    if action in {"list_config_overlay", "clear_config_overlay"}:
        return "config_overlay"
    if action == "validate_schema":
        return "schema"
    if action == "generate_evolution_report":
        return "report"
    return "system"


def build_action_descriptor(
    action_kind: str,
    *,
    target_id: str,
    decision: EvolutionPolicyDecision,
    risk_level: str,
    source: str,
) -> EvolutionActionDescriptor:
    action = normalize_action_kind(action_kind)
    target = str(target_id or "").strip()
    return EvolutionActionDescriptor(
        action_id=action_id_for(action, target),
        action_kind=action,
        target_type=target_type_for_action(action),
        target_id=target,
        permission=decision.permission,
        risk_level=str(risk_level or risk_level_for_action(action)),
        previewable=action in PREVIEW_ACTIONS or action in READ_ACTIONS,
        executable=action in READ_ACTIONS or action in WRITE_ACTIONS,
        requires_manual_override=decision.requires_manual_override,
        summary=action_summary(action),
        source=str(source or "control_plane"),
        parameters_schema=action_parameters_schema(action),
        suggested_my_action=suggested_my_action(action, target),
        policy=decision.to_json(),
    )


def action_id_for(action_kind: str, target_id: str = "") -> str:
    action = normalize_action_kind(action_kind)
    target = str(target_id or "").strip()
    return f"{action}:{target}" if target else action


def risk_level_for_action(action_kind: str) -> str:
    action = normalize_action_kind(action_kind)
    if action == "rollback_artifact":
        return "high"
    if action in {"retry_trial", "suppress_signal", "resume_signal", "run_feedback_calibration", "clear_config_overlay"}:
        return "medium"
    return "low"


def action_summary(action_kind: str) -> str:
    action = normalize_action_kind(action_kind)
    summaries = {
        "status": "Return the governed evolution dashboard read model.",
        "list_actions": "List stable control-plane action descriptors.",
        "list_signals": "List governed evolution opportunity signals.",
        "list_proposals": "List pending or recent auto-evolution proposals.",
        "inspect_signal": "Inspect one opportunity signal and related operator guidance.",
        "inspect_evolution_proposal": "Inspect one auto-evolution review proposal.",
        "explain_evolution_health": "Explain evolution health score, trend, and recommendations.",
        "list_evolution_recommendations": "List structured operator recommendations.",
        "list_recommendations": "List structured operator recommendations.",
        "generate_evolution_report": "Generate a read-only Markdown evolution report.",
        "validate_schema": "Validate governed evolution stores without mutation.",
        "list_config_overlay": "Inspect the governed self-tuning config overlay.",
        "suppress_signal": "Manually suppress an opportunity signal.",
        "resume_signal": "Resume a manually or feedback-suppressed opportunity signal.",
        "run_maintenance": "Run governed evolution retention and cleanup maintenance.",
        "force_cleanup": "Run maintenance with force cleanup enabled.",
        "run_feedback_calibration": "Apply feedback calibration to opportunity signals.",
        "retry_trial": "Retry read-only isolated trial for an auto-evolution workflow proposal.",
        "rollback_artifact": "Rollback a workflow or skill from a governed snapshot.",
        "clear_config_overlay": "Clear governed self-tuning config overrides.",
    }
    return summaries.get(action, f"Unsupported evolution action `{action}`.")


def action_parameters_schema(action_kind: str) -> dict[str, Any]:
    action = normalize_action_kind(action_kind)
    schema: dict[str, Any] = {
        "type": "object",
        "properties": {
            "target_id": {"type": "string"},
            "reason": {"type": "string"},
        },
        "additionalProperties": True,
    }
    if action in {"inspect_signal", "suppress_signal", "resume_signal"}:
        schema["required"] = ["target_id"]
    elif action in {"inspect_evolution_proposal", "retry_trial"}:
        schema["required"] = ["target_id"]
        schema["properties"]["fixtures"] = {
            "type": "object",
            "additionalProperties": {"type": "string"},
        }
    elif action == "rollback_artifact":
        schema["required"] = ["artifact_type", "artifact_name"]
        schema["properties"].update({
            "artifact_type": {"type": "string", "enum": ["workflow", "skill"]},
            "artifact_name": {"type": "string"},
            "snapshot_id": {"type": "string"},
            "force_cleanup": {"type": "boolean"},
        })
    elif action == "generate_evolution_report":
        schema["properties"]["period_days"] = {"type": "integer", "minimum": 1, "maximum": 90}
    elif action in {"run_maintenance", "force_cleanup"}:
        schema["properties"]["force_cleanup"] = {"type": "boolean"}
    return schema


def suggested_my_action(action_kind: str, target_id: str = "") -> str:
    action = normalize_action_kind(action_kind)
    target = str(target_id or "").strip()
    read_aliases = {
        "status": "evolution_status",
        "list_actions": "list_evolution_actions",
        "list_signals": "list_evolution_signals",
        "list_proposals": "list_evolution_proposals",
        "list_recommendations": "list_evolution_recommendations",
        "validate_schema": "validate_evolution_schema",
        "list_config_overlay": "list_config_overlay",
    }
    if action in read_aliases:
        return f"my action={read_aliases[action]}"
    if action in {"inspect_signal", "suppress_signal", "resume_signal"} and target:
        return f"my action={action} key={target}"
    if action in {"inspect_evolution_proposal", "retry_trial"} and target:
        return f"my action={action} key={target}"
    if action == "rollback_artifact":
        return "my action=rollback_artifact value={artifact_type, artifact_name, snapshot_id}"
    if action == "generate_evolution_report":
        return "my action=generate_evolution_report value={period_days: 7}"
    return f"my action={action}"


def permission_summary(config: Any | None) -> dict[str, Any]:
    manual_override = bool(getattr(config, "allow_manual_override", False) if config is not None else False)
    return {
        "read": True,
        "preview": True,
        "maintenance": manual_override,
        "override": manual_override,
        "rollback": manual_override,
        "apply": False,
    }


def safety_boundaries() -> list[str]:
    return [
        "No auto active.",
        "No auto apply proposal.",
        "No auto suppress signal.",
        "No auto main config mutation.",
        "Self-tuning may only write safety-tightening governed config overlays.",
        "Trial stays read-only and isolated.",
        "Evolution-generated skills must not be always-on automatically.",
    ]


def evolution_maintenance_policy(config: Any | None) -> dict[str, Any]:
    trial = getattr(config, "trial", None) if config is not None else None
    return {
        "outcome_retention_days": int(getattr(config, "outcome_retention_days", 90) if config is not None else 90),
        "outcome_archive_enabled": bool(getattr(config, "outcome_archive_enabled", True) if config is not None else True),
        "dependency_stale_cleanup_enabled": bool(getattr(config, "dependency_stale_cleanup_enabled", True) if config is not None else True),
        "health_history_retention_days": int(getattr(config, "health_history_retention_days", 90) if config is not None else 90),
        "max_health_history_snapshots": int(getattr(config, "max_health_history_snapshots", 100) if config is not None else 100),
        "trial_log_retention_days": int(getattr(trial, "trial_log_retention_days", 30) if trial is not None else 30),
        "max_retained_trial_logs": int(getattr(trial, "max_retained_trial_logs", 10) if trial is not None else 10),
    }


def _static_gate_issue_counts(records: list[dict[str, Any]]) -> dict[str, int]:
    issue_counts: dict[str, int] = {}
    for record in records:
        payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
        gate = payload.get("promotion_gate") if isinstance(payload.get("promotion_gate"), dict) else {}
        if not gate:
            gate = payload.get("static_gate") if isinstance(payload.get("static_gate"), dict) else {}
        counts = gate.get("issue_counts") if isinstance(gate.get("issue_counts"), dict) else {}
        for severity, count in counts.items():
            try:
                issue_counts[str(severity)] = issue_counts.get(str(severity), 0) + int(count)
            except (TypeError, ValueError):
                continue
    return issue_counts


def _promotion_gate_counts(records: list[dict[str, Any]]) -> dict[str, int]:
    counts: dict[str, int] = {}
    for record in records:
        payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
        gate = payload.get("promotion_gate") if isinstance(payload.get("promotion_gate"), dict) else {}
        decision = str(gate.get("decision") or "").strip()
        if decision:
            counts[decision] = counts.get(decision, 0) + 1
    return counts


def _auto_verified_workflow_count(records: list[dict[str, Any]]) -> int:
    total = 0
    for record in records:
        payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
        evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
        if (
            str(record.get("proposal_type") or "") == "workflow"
            and str(evolution.get("origin") or "") == AUTO_EVOLUTION_ORIGIN
            and str(record.get("review_reason") or "") == "auto_evolution verified low-risk workflow proposal"
        ):
            total += 1
    return total


def _coerce_period_days(value: Any) -> int:
    if isinstance(value, bool):
        return 7
    try:
        return max(1, min(int(value or 7), 90))
    except (TypeError, ValueError):
        return 7
