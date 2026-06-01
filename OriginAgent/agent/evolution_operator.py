"""Operator-facing read models for governed evolution."""

from __future__ import annotations

from contextlib import suppress
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

from OriginAgent.agent.evolution import AUTO_EVOLUTION_ORIGIN, OpportunitySignal, OpportunitySignalStore
from OriginAgent.agent.evolution_feedback import feedback_status
from OriginAgent.agent.evolution_health import evolution_health_score
from OriginAgent.agent.evolution_health_history import EvolutionHealthHistoryStore
from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore, proposal_outcome_context, safe_append_outcome
from OriginAgent.agent.evolution_sandbox import SandboxEvaluator, sandbox_status_counts, trial_policy_status
from OriginAgent.agent.evolution_schema import validate_evolution_stores
from OriginAgent.agent.evolution_trial import TrialRunner
from OriginAgent.agent.evolution_trial_logs import EvolutionTrialLogStore
from OriginAgent.agent.evolution_dependencies import EvolutionDependencyStore
from OriginAgent.utils.helpers import truncate_text

_STEP_OUTPUT_PREVIEW_CHARS = 500
_MAX_RECOMMENDATIONS = 10
_REPORT_MAX_RECOMMENDATIONS = 6


@dataclass(frozen=True)
class RetryTrialResult:
    ok: bool
    proposal_id: str
    status: str
    message: str
    trial: dict[str, Any] | None = None
    proposal: dict[str, Any] | None = None
    error: str = ""

    def to_json(self) -> dict[str, Any]:
        return {
            "ok": self.ok,
            "proposal_id": self.proposal_id,
            "status": self.status,
            "message": self.message,
            "trial": self.trial,
            "proposal": self.proposal,
            "error": self.error,
        }


class EvolutionOperator:
    """Build operator views and guarded operations for evolution artifacts."""

    def __init__(self, workspace: Path, config: Any | None = None) -> None:
        self.workspace = Path(workspace)
        self.config = config

    def inspect_signal(self, opportunity_id: str) -> dict[str, Any]:
        signal_id = str(opportunity_id or "").strip()
        for signal in OpportunitySignalStore(self.workspace).read_all():
            if signal.opportunity_id == signal_id:
                return {
                    "found": True,
                    "signal": signal.to_record(),
                    "recommendations": _signal_recommendations(signal.to_record()),
                }
        return {
            "found": False,
            "error": "signal_not_found",
            "opportunity_id": signal_id,
        }

    def inspect_proposal(self, proposal_id: str) -> dict[str, Any]:
        from OriginAgent.agent.background_review import ReviewProposalStore

        proposal_key = str(proposal_id or "").strip()
        record = ReviewProposalStore(self.workspace).get(proposal_key)
        if record is None:
            return {
                "found": False,
                "error": "proposal_not_found",
                "proposal_id": proposal_key,
            }
        payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
        insights = build_operator_insights(
            payload,
            proposal_type=str(record.get("proposal_type") or ""),
            config=self.config,
        )
        return {
            "found": True,
            "proposal_id": proposal_key,
            "status": str(record.get("status") or "pending"),
            "proposal_type": str(record.get("proposal_type") or ""),
            "origin": str(record.get("origin") or ""),
            "title": str(record.get("title") or ""),
            "summary": _proposal_summary(record, payload, insights),
            "operator_insights": insights,
            "can_retry_trial": _can_retry_trial(record),
        }

    def explain_health(self) -> dict[str, Any]:
        outcome_stats = EvolutionOutcomeStore(self.workspace).stats()
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
        history = EvolutionHealthHistoryStore(self.workspace).summary()
        return {
            "health": health,
            "history": history,
            "recommendations": self.list_recommendations(
                health=health,
                health_history=history,
                outcome_stats=outcome_stats,
                dependency_stats=dependency_stats,
                feedback_stats=feedback_stats,
                sandbox_counts=sandbox_counts,
            ),
        }

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
        from OriginAgent.agent.background_review import ReviewProposalStore

        outcome_stats = outcome_stats or EvolutionOutcomeStore(self.workspace).stats()
        dependency_stats = dependency_stats or EvolutionDependencyStore(self.workspace).stats()
        feedback_stats = feedback_stats or feedback_status(self.workspace, self.config)
        sandbox_counts = sandbox_counts or sandbox_status_counts(self.workspace)
        if health is None:
            health = evolution_health_score(
                outcome_stats=outcome_stats,
                dependency_stats=dependency_stats,
                feedback_stats=feedback_stats,
                sandbox_counts=sandbox_counts,
                trial_status=trial_policy_status(self.config),
            )
        if health_history is None:
            health_history = EvolutionHealthHistoryStore(self.workspace).summary()

        items: list[dict[str, Any]] = []
        if str(health_history.get("trend") or "") == "degrading":
            items.append(_recommendation(
                code="health_trend_degrading",
                severity="warning",
                action="consider_conservative_or_dry_run",
                action_kind="explain_evolution_health",
                target_type="health",
                risk_level="medium",
                message="Evolution health is trending down; consider conservative mode or dry_run=true until proposals are reviewed.",
                preview="No state changes. Shows health score, trend, and reasons.",
            ))
        if int(dependency_stats.get("stale_reference_count") or 0) > 0:
            items.append(_recommendation(
                code="stale_dependencies",
                severity="warning",
                action="run_maintenance",
                action_kind="run_maintenance",
                target_type="maintenance",
                requires_manual_override=True,
                risk_level="low",
                message="Stale evolution dependencies exist; run maintenance before promoting or rolling back related artifacts.",
                preview="Would prune stale dependency references and enforce bounded evolution store retention.",
            ))
        blocked_or_failed = int(sandbox_counts.get("blocked") or 0) + int(sandbox_counts.get("failed") or 0)
        if blocked_or_failed:
            items.append(_recommendation(
                code="sandbox_attention_needed",
                severity="warning",
                action="inspect_or_retry_trial",
                action_kind="inspect_evolution_proposal",
                target_type="proposal",
                risk_level="medium",
                message=f"{blocked_or_failed} auto-evolution proposal(s) were blocked or failed by sandbox checks.",
                preview="No state changes. Inspect blocked proposals before retrying any trial.",
            ))
        rollback_succeeded = int(_mapping(outcome_stats.get("rollback_status_counts")).get("succeeded") or 0)
        if rollback_succeeded:
            items.append(_recommendation(
                code="rollback_feedback",
                severity="warning",
                action="inspect_suppression_candidates",
                action_kind="inspect_signal",
                target_type="signal",
                risk_level="medium",
                message="Recent successful rollbacks occurred; inspect related signals and consider suppression.",
                preview="No state changes. Inspect rollback-linked signals before deciding whether to suppress them.",
            ))
        trend_counts = _mapping(feedback_stats.get("feedback_trend_counts"))
        if int(trend_counts.get("negative") or 0) > 0:
            items.append(_recommendation(
                code="negative_feedback",
                severity="info",
                action="inspect_feedback_adjusted_signals",
                action_kind="inspect_signal",
                target_type="signal",
                risk_level="low",
                message="Recent negative feedback lowered one or more opportunity signals.",
                preview="No state changes. Review feedback-adjusted signals and their current priority.",
            ))

        store = ReviewProposalStore(self.workspace)
        for record in store.list_records(origin=AUTO_EVOLUTION_ORIGIN, status="pending", limit=20):
            payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
            insights = build_operator_insights(
                payload,
                proposal_type=str(record.get("proposal_type") or ""),
                config=self.config,
            )
            recommended_action = str(insights.get("recommended_action") or "")
            if recommended_action in {"review_required", "reject"}:
                items.append(_recommendation(
                    code="pending_evolution_proposal",
                    severity="info" if recommended_action == "review_required" else "warning",
                    action="inspect_proposal",
                    action_kind="inspect_evolution_proposal",
                    target_type="proposal",
                    target_id=str(record.get("id") or ""),
                    risk_level=str(_mapping(insights.get("risk_summary")).get("level") or "medium"),
                    message=f"Review auto-evolution proposal `{record.get('id')}`: {record.get('title')}",
                    preview="No state changes. Shows gate, trial, risk, and activation policy summary.",
                    proposal_id=str(record.get("id") or ""),
                    opportunity_id=str(_mapping(payload.get("evolution")).get("opportunity_id") or ""),
                ))
            if len(items) >= _MAX_RECOMMENDATIONS:
                break

        for signal in OpportunitySignalStore(self.workspace).read_all():
            for item in _signal_recommendations(signal.to_record()):
                items.append(item)
                if len(items) >= _MAX_RECOMMENDATIONS:
                    break
            if len(items) >= _MAX_RECOMMENDATIONS:
                break
        return items[:_MAX_RECOMMENDATIONS]

    def preview_action(
        self,
        action_kind: str,
        *,
        target_id: str = "",
        reason: str = "",
        fixtures: dict[str, str] | None = None,
        force_cleanup: bool = False,
    ) -> dict[str, Any]:
        """Preview an operator action without mutating governed evolution state."""

        action = _normalize_action_kind(action_kind)
        target = str(target_id or "").strip()
        if action == "suppress_signal":
            return self._preview_signal_update(
                target,
                next_status="suppressed",
                reason=reason or "Manual evolution control-plane suppression.",
            )
        if action == "resume_signal":
            return self._preview_signal_update(
                target,
                next_status="open",
                reason=reason or "Manual evolution control-plane resume.",
            )
        if action in {"run_maintenance", "force_cleanup"}:
            return self._preview_maintenance(force_cleanup=force_cleanup or action == "force_cleanup")
        if action == "run_feedback_calibration":
            return self._preview_feedback_calibration()
        if action == "retry_trial":
            return self._preview_retry_trial(target, fixtures=fixtures or {})
        if action == "inspect_evolution_proposal":
            return {
                "ok": True,
                "action_kind": action,
                "target_type": "proposal",
                "target_id": target,
                "requires_manual_override": False,
                "will_write": False,
                "preview": self.inspect_proposal(target),
            }
        if action == "inspect_signal":
            return {
                "ok": True,
                "action_kind": action,
                "target_type": "signal",
                "target_id": target,
                "requires_manual_override": False,
                "will_write": False,
                "preview": self.inspect_signal(target),
            }
        return {
            "ok": False,
            "action_kind": action,
            "target_id": target,
            "requires_manual_override": False,
            "will_write": False,
            "error": "unsupported_action",
            "message": f"Unsupported evolution operator action `{action_kind}`.",
        }

    def generate_report(self, *, period_days: int = 7) -> str:
        """Return a read-only Markdown evolution operator report."""

        days = max(1, min(int(period_days or 7), 90))
        now = datetime.now(timezone.utc)
        cutoff = now - timedelta(days=days)
        outcomes = EvolutionOutcomeStore(self.workspace)
        recent_events = [
            event for event in outcomes.read_all()
            if (_parse_datetime(str(event.get("timestamp") or "")) or now) >= cutoff
        ]
        outcome_type_counts = _counts(event.get("type") for event in recent_events)
        sandbox_counts = _counts(event.get("sandbox_status") for event in recent_events)
        review_counts = _counts(event.get("review_status") for event in recent_events)
        rollback_counts = _counts(event.get("rollback_status") for event in recent_events)
        signals = OpportunitySignalStore(self.workspace).read_all()
        open_signals = len([signal for signal in signals if signal.status == "open"])
        converted_signals = len([signal for signal in signals if signal.status == "converted"])
        suppressed_signals = len([signal for signal in signals if signal.status == "suppressed"])
        health = self.explain_health()
        health_summary = _mapping(health.get("health"))
        history = _mapping(health.get("history"))
        recommendations = self.list_recommendations()[:_REPORT_MAX_RECOMMENDATIONS]
        schema_validation = validate_evolution_stores(self.workspace)
        lines = [
            f"# Evolution Operator Report ({days}d)",
            "",
            "## Summary",
            f"- Health: {health_summary.get('score', 'unknown')} ({health_summary.get('level', 'unknown')})",
            f"- Health trend: {history.get('trend', 'unknown')} (delta {history.get('score_delta', 0)})",
            f"- Signals: {open_signals} open, {converted_signals} converted, {suppressed_signals} suppressed",
            f"- Recent outcome events: {len(recent_events)}",
            f"- Schema validation: {'ok' if schema_validation.get('ok') else 'attention needed'}",
            "",
            "## Recent Outcomes",
            _format_counts(outcome_type_counts),
            "",
            "## Sandbox And Review",
            f"- Sandbox: {_inline_counts(sandbox_counts)}",
            f"- Review: {_inline_counts(review_counts)}",
            f"- Rollback: {_inline_counts(rollback_counts)}",
            f"- Schema issues: {_inline_counts(_mapping(schema_validation.get('issue_counts')))}",
            "",
            "## Health Reasons",
        ]
        reasons = health_summary.get("reasons") if isinstance(health_summary.get("reasons"), list) else []
        lines.extend([f"- {reason}" for reason in reasons] or ["- No health reasons reported."])
        lines.extend(["", "## Recommendations"])
        if recommendations:
            for item in recommendations:
                lines.append(
                    f"- [{item.get('severity')}] {item.get('code')}: {item.get('message')} "
                    f"(action: {item.get('action_kind')})"
                )
        else:
            lines.append("- No operator recommendations.")
        lines.extend([
            "",
            "## Safety Boundaries",
            "- No recommendation automatically activates artifacts.",
            "- Preview actions do not write proposal, outcome, trial log, or signal stores.",
            "- Write actions still require `learning.evolution.allow_manual_override=true`.",
        ])
        return "\n".join(lines).strip() + "\n"

    def _preview_signal_update(self, opportunity_id: str, *, next_status: str, reason: str) -> dict[str, Any]:
        signal = _find_signal(OpportunitySignalStore(self.workspace).read_all(), opportunity_id)
        if signal is None:
            return {
                "ok": False,
                "action_kind": "suppress_signal" if next_status == "suppressed" else "resume_signal",
                "target_type": "signal",
                "target_id": opportunity_id,
                "requires_manual_override": True,
                "will_write": False,
                "error": "signal_not_found",
                "message": "Opportunity signal was not found.",
            }
        if next_status == "suppressed" and signal.status == "converted":
            return {
                "ok": False,
                "action_kind": "suppress_signal",
                "target_type": "signal",
                "target_id": opportunity_id,
                "requires_manual_override": True,
                "will_write": False,
                "error": "signal_converted",
                "message": "Converted opportunity signals cannot be manually suppressed.",
            }
        if next_status == "open" and signal.status != "suppressed":
            return {
                "ok": False,
                "action_kind": "resume_signal",
                "target_type": "signal",
                "target_id": opportunity_id,
                "requires_manual_override": True,
                "will_write": False,
                "error": "signal_not_suppressed",
                "message": "Only suppressed opportunity signals can be resumed.",
            }
        return {
            "ok": True,
            "action_kind": "suppress_signal" if next_status == "suppressed" else "resume_signal",
            "target_type": "signal",
            "target_id": signal.opportunity_id,
            "requires_manual_override": True,
            "will_write": False,
            "preview": {
                "current_status": signal.status,
                "next_status": next_status,
                "current_verification_status": signal.verification_status,
                "next_verification_status": "manual_suppressed" if next_status == "suppressed" else "manual_resumed",
                "reason": truncate_text(reason, 512),
                "priority_score": signal.priority_score,
                "risk_level": signal.risk_level,
            },
        }

    def _preview_maintenance(self, *, force_cleanup: bool) -> dict[str, Any]:
        outcome_preview = _retention_preview(
            EvolutionOutcomeStore(self.workspace).read_all(),
            retention_days=_config_int(self.config, "outcome_retention_days", 90),
            permanent_types={"promoted", "rolled_back", "review_rejected", "review_failed"},
        )
        dependency_stats = EvolutionDependencyStore(self.workspace).stats()
        trial_config = getattr(self.config, "trial", None)
        trial_log_preview = _retention_preview(
            EvolutionTrialLogStore(self.workspace).read_all(),
            retention_days=_config_int(trial_config, "trial_log_retention_days", 30),
            max_records=_config_int(trial_config, "max_retained_trial_logs", 10),
        )
        health_history_preview = _retention_preview(
            EvolutionHealthHistoryStore(self.workspace).read_all(),
            retention_days=_config_int(self.config, "health_history_retention_days", 90),
            max_records=_config_int(self.config, "max_health_history_snapshots", 100),
        )
        schema_validation = validate_evolution_stores(self.workspace)
        return {
            "ok": True,
            "action_kind": "force_cleanup" if force_cleanup else "run_maintenance",
            "target_type": "maintenance",
            "target_id": "",
            "requires_manual_override": True,
            "will_write": False,
            "preview": {
                "outcome_retention": outcome_preview,
                "dependency_cleanup": {
                    "enabled": force_cleanup or bool(getattr(self.config, "dependency_stale_cleanup_enabled", True)),
                    "stale_reference_count": dependency_stats.get("stale_reference_count", 0),
                    "tracked_artifacts": dependency_stats.get("tracked_artifacts", 0),
                },
                "trial_log_retention": trial_log_preview,
                "health_history_retention": health_history_preview,
                "schema_validation": {
                    "ok": schema_validation.get("ok"),
                    "record_counts": schema_validation.get("record_counts", {}),
                    "issue_counts": schema_validation.get("issue_counts", {}),
                },
                "would_append_health_snapshot": True,
            },
        }

    def _preview_feedback_calibration(self) -> dict[str, Any]:
        status = feedback_status(self.workspace, self.config)
        outcomes = EvolutionOutcomeStore(self.workspace).read_all()
        candidates = [
            event for event in outcomes
            if str(event.get("type") or "") in {"review_rejected", "review_failed", "review_approved", "rolled_back"}
        ]
        return {
            "ok": True,
            "action_kind": "run_feedback_calibration",
            "target_type": "feedback",
            "target_id": "",
            "requires_manual_override": True,
            "will_write": False,
            "preview": {
                "candidate_feedback_events": len(candidates),
                "processed_event_count": status.get("processed_event_count", 0),
                "cooldown_count": status.get("cooldown_count", 0),
                "feedback_trend_counts": status.get("feedback_trend_counts", {}),
            },
        }

    def _preview_retry_trial(self, proposal_id: str, *, fixtures: dict[str, str]) -> dict[str, Any]:
        from OriginAgent.agent.background_review import ReviewProposalStore

        record = ReviewProposalStore(self.workspace).get(proposal_id)
        if record is None:
            return {
                "ok": False,
                "action_kind": "retry_trial",
                "target_type": "proposal",
                "target_id": proposal_id,
                "requires_manual_override": True,
                "will_write": False,
                "error": "proposal_not_found",
                "message": "Review proposal was not found.",
            }
        if not _can_retry_trial(record):
            return {
                "ok": False,
                "action_kind": "retry_trial",
                "target_type": "proposal",
                "target_id": proposal_id,
                "requires_manual_override": True,
                "will_write": False,
                "error": "unsupported_proposal",
                "message": "Only pending auto-evolution workflow proposals can retry trial.",
            }
        payload = dict(record.get("payload") if isinstance(record.get("payload"), dict) else {})
        payload["review_proposal_id"] = proposal_id
        gate = SandboxEvaluator(self.workspace, self.config).evaluate_trial_workflow_payload(payload)
        return {
            "ok": True,
            "action_kind": "retry_trial",
            "target_type": "proposal",
            "target_id": proposal_id,
            "requires_manual_override": True,
            "will_write": False,
            "preview": {
                "trial_gate_status": str(gate.get("status") or ""),
                "steps_checked": _mapping(gate.get("replay_summary")).get("steps_checked", 0),
                "blocked_steps": _mapping(gate.get("replay_summary")).get("blocked_steps", 0),
                "fixture_count": len(fixtures),
                "would_update_proposal_payload": True,
                "would_write_trial_log": True,
                "would_write_outcome": True,
                "gate": {
                    "status": str(gate.get("status") or ""),
                    "issues": gate.get("issues") if isinstance(gate.get("issues"), list) else [],
                    "policy": gate.get("policy") if isinstance(gate.get("policy"), dict) else {},
                },
            },
        }

    def retry_trial(
        self,
        proposal_id: str,
        *,
        fixtures: dict[str, str] | None = None,
        actor: str = "manual_override",
    ) -> RetryTrialResult:
        from OriginAgent.agent.background_review import ReviewProposalStore

        proposal_key = str(proposal_id or "").strip()
        store = ReviewProposalStore(self.workspace)
        record = store.get(proposal_key)
        if record is None:
            return RetryTrialResult(
                ok=False,
                proposal_id=proposal_key,
                status="missing",
                message="Review proposal was not found.",
                error="proposal_not_found",
            )
        if not _can_retry_trial(record):
            return RetryTrialResult(
                ok=False,
                proposal_id=proposal_key,
                status="unsupported",
                message="Only pending auto-evolution workflow proposals can retry trial.",
                proposal=record,
                error="unsupported_proposal",
            )
        payload = dict(record.get("payload") if isinstance(record.get("payload"), dict) else {})
        payload["review_proposal_id"] = proposal_key
        trial = TrialRunner(self.workspace, self.config).run_workflow_payload(payload, fixtures=fixtures or {})
        compact_trial = compact_trial_result(trial)
        payload["trial"] = compact_trial
        payload["operator_insights"] = build_operator_insights(
            payload,
            proposal_type=str(record.get("proposal_type") or ""),
            config=self.config,
        )
        updated = store.update_payload(proposal_key, payload)
        if updated is None:
            context = proposal_outcome_context(record)
            context["sandbox_status"] = str(compact_trial.get("status") or "")
            safe_append_outcome(
                EvolutionOutcomeStore(self.workspace),
                "trial_retry_failed",
                **context,
                metadata={
                    "actor": actor,
                    "trial_id": str(compact_trial.get("log_id") or ""),
                    "error": "payload_update_failed",
                },
            )
            return RetryTrialResult(
                ok=False,
                proposal_id=proposal_key,
                status="failed",
                message="Trial retry ran but proposal payload could not be updated.",
                trial=compact_trial,
                proposal=record,
                error="payload_update_failed",
            )
        context = proposal_outcome_context(updated or record)
        context["sandbox_status"] = str(compact_trial.get("status") or "")
        safe_append_outcome(
            EvolutionOutcomeStore(self.workspace),
            "trial_retried",
            **context,
            metadata={
                "actor": actor,
                "trial_id": str(compact_trial.get("log_id") or ""),
                "summary": compact_trial.get("summary") if isinstance(compact_trial.get("summary"), dict) else {},
            },
        )
        return RetryTrialResult(
            ok=True,
            proposal_id=proposal_key,
            status=str(compact_trial.get("status") or "unknown"),
            message="Trial retry completed and proposal payload was updated.",
            trial=compact_trial,
            proposal=updated,
        )


def build_operator_insights(
    payload: dict[str, Any],
    *,
    proposal_type: str,
    config: Any | None = None,
) -> dict[str, Any]:
    """Summarize proposal evidence for human operators."""

    static_gate = _mapping(payload.get("static_gate"))
    sandbox = _mapping(payload.get("sandbox"))
    trial = _mapping(payload.get("trial"))
    gate = _mapping(payload.get("promotion_gate"))
    evolution = _mapping(payload.get("evolution"))
    proposal_kind = str(proposal_type or payload.get("subject_type") or "").strip().lower()

    risk_level = str(gate.get("risk_level") or evolution.get("risk_level") or "medium")
    recommended_action = str(gate.get("suggested_action") or "review_required")
    return {
        "trial_summary": _trial_summary(trial or sandbox),
        "risk_summary": {
            "level": risk_level,
            "static_gate_decision": str(static_gate.get("decision") or ""),
            "promotion_gate_decision": str(gate.get("decision") or ""),
            "sandbox_status": str((trial or sandbox).get("status") or gate.get("sandbox_status") or ""),
            "issue_counts": dict(gate.get("issue_counts") or static_gate.get("issue_counts") or {}),
            "issues": _issue_messages(static_gate, sandbox, trial),
        },
        "health_impact": _estimate_health_impact(evolution, gate, sandbox, trial),
        "recommended_action": recommended_action,
        "why_not_auto_active": _why_not_auto_active(
            proposal_kind=proposal_kind,
            gate=gate,
            config=config,
        ),
    }


def compact_trial_result(result: dict[str, Any]) -> dict[str, Any]:
    step_results = result.get("step_results") if isinstance(result.get("step_results"), list) else []
    compact_steps: list[dict[str, Any]] = []
    for step in step_results[:20]:
        if not isinstance(step, dict):
            continue
        output = str(step.get("output") or "")
        compact_steps.append({
            "index": step.get("index"),
            "title": truncate_text(str(step.get("title") or ""), 120),
            "tool": str(step.get("tool") or ""),
            "status": str(step.get("status") or ""),
            "executed": bool(step.get("executed")),
            "read_only": bool(step.get("read_only", True)),
            "isolated_workspace": bool(step.get("isolated_workspace", True)),
            "output": truncate_text(output, _STEP_OUTPUT_PREVIEW_CHARS),
            "output_truncated": len(output) > _STEP_OUTPUT_PREVIEW_CHARS,
            "issues": step.get("issues") if isinstance(step.get("issues"), list) else [],
        })
    return {
        "status": str(result.get("status") or "unknown"),
        "mode": str(result.get("mode") or "trial"),
        "read_only": bool(result.get("read_only", True)),
        "isolated_workspace": bool(result.get("isolated_workspace", True)),
        "gate_status": str(result.get("gate_status") or ""),
        "log_id": str(result.get("log_id") or ""),
        "summary": result.get("summary") if isinstance(result.get("summary"), dict) else {},
        "policy": result.get("policy") if isinstance(result.get("policy"), dict) else {},
        "step_results": compact_steps,
    }


def _trial_summary(source: dict[str, Any]) -> dict[str, Any]:
    step_results = source.get("step_results") if isinstance(source.get("step_results"), list) else []
    failed_steps = [
        {
            "index": step.get("index"),
            "title": step.get("title"),
            "tool": step.get("tool"),
            "status": step.get("status"),
        }
        for step in step_results
        if isinstance(step, dict) and str(step.get("status") or "") not in {"", "passed", "skipped"}
    ][:5]
    return {
        "status": str(source.get("status") or "not_run"),
        "mode": str(source.get("mode") or "sandbox"),
        "read_only": bool(source.get("read_only", True)),
        "isolated_workspace": bool(source.get("isolated_workspace", True)),
        "replay_summary": source.get("replay_summary") if isinstance(source.get("replay_summary"), dict) else {},
        "summary": source.get("summary") if isinstance(source.get("summary"), dict) else {},
        "failed_steps": failed_steps,
    }


def _issue_messages(*sources: dict[str, Any]) -> list[str]:
    messages: list[str] = []
    for source in sources:
        issues = source.get("issues") if isinstance(source.get("issues"), list) else []
        for issue in issues:
            if not isinstance(issue, dict):
                continue
            severity = str(issue.get("severity") or "")
            if severity not in {"pending", "reject", "warning"}:
                continue
            message = truncate_text(str(issue.get("message") or issue.get("code") or ""), 300)
            if message:
                messages.append(message)
    return messages[:8]


def _estimate_health_impact(
    evolution: dict[str, Any],
    gate: dict[str, Any],
    sandbox: dict[str, Any],
    trial: dict[str, Any],
) -> dict[str, Any]:
    status = str((trial or sandbox).get("status") or "")
    decision = str(gate.get("decision") or "")
    priority = _safe_float(evolution.get("priority_score"), 0.0)
    if decision == "blocked" or status in {"blocked", "failed"}:
        direction = "negative"
        estimate = -5
    elif decision == "pass" and status in {"passed", "skipped", ""}:
        direction = "neutral_positive"
        estimate = 1 if priority >= 0.8 else 0
    else:
        direction = "neutral"
        estimate = 0
    return {
        "estimated_score_delta": estimate,
        "direction": direction,
        "basis": {
            "priority_score": priority,
            "gate_decision": decision,
            "trial_or_sandbox_status": status,
        },
    }


def _why_not_auto_active(*, proposal_kind: str, gate: dict[str, Any], config: Any | None) -> list[str]:
    reasons = [
        "OriginAgent governance does not auto-activate evolution artifacts.",
        "`verified` is a validation state, not an activation state.",
    ]
    if proposal_kind == "skill":
        reasons.append("Evolution-generated skills must remain `always: false` and require review before activation.")
    elif proposal_kind == "workflow":
        reasons.append("Workflow proposals may be verified, but active use still requires an explicit review/apply path.")
    if bool(gate.get("auto_verify_eligible")):
        reasons.append("This proposal may be auto-verified when policy allows, but still will not become active automatically.")
    elif not bool(getattr(config, "auto_verify_workflows", False) if config is not None else False):
        reasons.append("Auto-verification is disabled by current policy.")
    return reasons


def _signal_recommendations(signal: dict[str, Any]) -> list[dict[str, Any]]:
    recommendations: list[dict[str, Any]] = []
    opportunity_id = str(signal.get("opportunity_id") or "")
    negative_count = int(signal.get("feedback_negative_count") or 0)
    if negative_count >= 2 and str(signal.get("status") or "") != "suppressed":
        recommendations.append(_recommendation(
            code="signal_negative_feedback",
            severity="warning",
            action="consider_suppress_signal",
            action_kind="suppress_signal",
            target_type="signal",
            target_id=opportunity_id,
            requires_manual_override=True,
            risk_level=str(signal.get("risk_level") or "medium"),
            message="Opportunity signal has repeated negative feedback; consider suppressing it.",
            preview="Would mark the signal suppressed and record a manual suppression outcome.",
            opportunity_id=opportunity_id,
        ))
    if str(signal.get("verification_status") or "") == "rolled_back":
        recommendations.append(_recommendation(
            code="signal_rolled_back",
            severity="warning",
            action="keep_suppressed_or_rework",
            action_kind="inspect_signal",
            target_type="signal",
            target_id=opportunity_id,
            risk_level="high",
            message="Opportunity came from a rolled-back artifact; keep it suppressed or rework before retrying.",
            preview="No state changes. Inspect rollback feedback before reworking this signal.",
            opportunity_id=opportunity_id,
        ))
    return recommendations


def _can_retry_trial(record: dict[str, Any]) -> bool:
    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    evolution = _mapping(payload.get("evolution"))
    return (
        str(record.get("status") or "pending") == "pending"
        and str(record.get("proposal_type") or "") == "workflow"
        and (
            str(record.get("origin") or "") == AUTO_EVOLUTION_ORIGIN
            or str(evolution.get("origin") or "") == AUTO_EVOLUTION_ORIGIN
        )
    )


def _recommendation(
    *,
    code: str,
    severity: str,
    action: str,
    message: str,
    action_kind: str = "",
    target_type: str = "",
    target_id: str = "",
    requires_manual_override: bool = False,
    risk_level: str = "low",
    preview: str = "",
    proposal_id: str = "",
    opportunity_id: str = "",
) -> dict[str, Any]:
    normalized_action = _normalize_action_kind(action_kind or action)
    effective_target = target_id or proposal_id or opportunity_id
    if not target_type:
        if proposal_id:
            target_type = "proposal"
        elif opportunity_id:
            target_type = "signal"
        else:
            target_type = "system"
    result = {
        "code": code,
        "severity": severity,
        "action": action,
        "action_kind": normalized_action,
        "target_type": target_type,
        "target_id": effective_target,
        "requires_manual_override": bool(
            requires_manual_override
            or normalized_action in {
                "suppress_signal",
                "resume_signal",
                "run_maintenance",
                "force_cleanup",
                "run_feedback_calibration",
                "retry_trial",
            }
        ),
        "risk_level": risk_level or "low",
        "preview": preview,
        "suggested_my_action": _suggested_my_action(normalized_action, effective_target),
        "message": message,
    }
    if proposal_id:
        result["proposal_id"] = proposal_id
    if opportunity_id:
        result["opportunity_id"] = opportunity_id
    return result


def _proposal_summary(
    record: dict[str, Any],
    payload: dict[str, Any],
    insights: dict[str, Any],
) -> dict[str, Any]:
    risk = _mapping(insights.get("risk_summary"))
    trial = _mapping(insights.get("trial_summary"))
    evolution = _mapping(payload.get("evolution"))
    gate = _mapping(payload.get("promotion_gate"))
    return {
        "proposal_id": str(record.get("id") or ""),
        "status": str(record.get("status") or "pending"),
        "proposal_type": str(record.get("proposal_type") or ""),
        "origin": str(record.get("origin") or ""),
        "title": str(record.get("title") or ""),
        "subject_id": str(payload.get("subject_id") or payload.get("workflow_name") or payload.get("skill_name") or ""),
        "subject_path": str(payload.get("subject_path") or ""),
        "opportunity_id": str(evolution.get("opportunity_id") or ""),
        "priority_score": evolution.get("priority_score"),
        "risk_level": str(risk.get("level") or "medium"),
        "recommended_action": str(insights.get("recommended_action") or "review_required"),
        "promotion_gate_decision": str(risk.get("promotion_gate_decision") or gate.get("decision") or ""),
        "static_gate_decision": str(risk.get("static_gate_decision") or ""),
        "sandbox_status": str(risk.get("sandbox_status") or ""),
        "trial_status": str(trial.get("status") or "not_run"),
        "failed_steps": trial.get("failed_steps") if isinstance(trial.get("failed_steps"), list) else [],
        "issue_counts": risk.get("issue_counts") if isinstance(risk.get("issue_counts"), dict) else {},
        "issues": risk.get("issues") if isinstance(risk.get("issues"), list) else [],
        "why_not_auto_active": (
            insights.get("why_not_auto_active")
            if isinstance(insights.get("why_not_auto_active"), list)
            else []
        ),
    }


def _normalize_action_kind(value: str) -> str:
    action = str(value or "").strip()
    aliases = {
        "inspect_proposal": "inspect_evolution_proposal",
        "inspect_or_retry_trial": "inspect_evolution_proposal",
        "inspect_feedback_adjusted_signals": "inspect_signal",
        "inspect_suppression_candidates": "inspect_signal",
        "consider_suppress_signal": "suppress_signal",
        "consider_conservative_or_dry_run": "explain_evolution_health",
    }
    return aliases.get(action, action)


def _suggested_my_action(action_kind: str, target_id: str) -> str:
    if not action_kind:
        return ""
    if action_kind in {"inspect_signal", "suppress_signal", "resume_signal"} and target_id:
        return f"my action={action_kind} key={target_id}"
    if action_kind in {"inspect_evolution_proposal", "retry_trial"} and target_id:
        return f"my action={action_kind} key={target_id}"
    if action_kind in {
        "explain_evolution_health",
        "list_evolution_recommendations",
        "run_maintenance",
        "force_cleanup",
        "run_feedback_calibration",
        "generate_evolution_report",
    }:
        return f"my action={action_kind}"
    return f"my action={action_kind}"


def _find_signal(signals: list[OpportunitySignal], opportunity_id: str) -> OpportunitySignal | None:
    target = str(opportunity_id or "").strip()
    if not target:
        return None
    for signal in signals:
        if signal.opportunity_id == target:
            return signal
    return None


def _retention_preview(
    records: list[dict[str, Any]],
    *,
    retention_days: int,
    max_records: int | None = None,
    permanent_types: set[str] | None = None,
) -> dict[str, Any]:
    now = datetime.now(timezone.utc)
    cutoff = now - timedelta(days=max(1, int(retention_days or 1)))
    permanent = permanent_types or set()
    recent: list[dict[str, Any]] = []
    removed_by_age = 0
    for record in records:
        timestamp = _parse_datetime(str(record.get("timestamp") or "")) or now
        if timestamp >= cutoff or str(record.get("type") or "") in permanent:
            recent.append(record)
        else:
            removed_by_age += 1
    recent.sort(key=lambda item: str(item.get("timestamp") or ""))
    removed_by_count = 0
    if max_records is not None:
        max_count = max(0, int(max_records or 0))
        if len(recent) > max_count:
            removed_by_count = len(recent) - max_count
            recent = recent[-max_count:] if max_count else []
    return {
        "retention_days": max(1, int(retention_days or 1)),
        "max_records": max_records,
        "current_count": len(records),
        "would_retain_count": len(recent),
        "would_remove_count": removed_by_age + removed_by_count,
        "would_remove_by_age": removed_by_age,
        "would_remove_by_count": removed_by_count,
        "cutoff": cutoff.isoformat(),
    }


def _parse_datetime(value: str) -> datetime | None:
    if not value:
        return None
    with suppress(ValueError):
        parsed = datetime.fromisoformat(value)
        if parsed.tzinfo is None:
            return parsed.replace(tzinfo=timezone.utc)
        return parsed.astimezone(timezone.utc)
    return None


def _counts(values: Any) -> dict[str, int]:
    counts: dict[str, int] = {}
    for value in values:
        key = str(value or "").strip()
        if not key:
            continue
        counts[key] = counts.get(key, 0) + 1
    return dict(sorted(counts.items()))


def _format_counts(counts: dict[str, int]) -> str:
    if not counts:
        return "- none"
    return "\n".join(f"- {key}: {value}" for key, value in counts.items())


def _inline_counts(counts: dict[str, int]) -> str:
    if not counts:
        return "none"
    return ", ".join(f"{key}={value}" for key, value in counts.items())


def _config_int(config: Any | None, attr: str, default: int) -> int:
    if isinstance(config, dict):
        value = config.get(attr, default)
    else:
        value = getattr(config, attr, default) if config is not None else default
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default


def _mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def _safe_float(value: Any, default: float) -> float:
    if isinstance(value, bool):
        return default
    try:
        return float(value)
    except (TypeError, ValueError):
        return default
