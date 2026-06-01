"""Promotion gate decisions for governed self-evolution proposals."""

from __future__ import annotations

from dataclasses import asdict, dataclass, field
from typing import Any


@dataclass(frozen=True)
class PromotionGateResult:
    """Normalized gate result attached to auto-evolution proposal payloads."""

    decision: str
    suggested_action: str
    risk_level: str
    reasons: list[str] = field(default_factory=list)
    static_gate_decision: str = ""
    sandbox_status: str = ""
    auto_verify_eligible: bool = False
    issue_counts: dict[str, int] = field(default_factory=dict)

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class PromotionGate:
    """Convert static checks, sandbox results, and policy thresholds into one decision."""

    def __init__(self, config: Any | None = None) -> None:
        self.config = config

    def evaluate(self, payload: dict[str, Any], *, proposal_type: str) -> PromotionGateResult:
        normalized_type = str(proposal_type or "").strip().lower()
        if normalized_type == "workflow":
            return self._evaluate_workflow(payload)
        if normalized_type == "skill":
            return self._evaluate_skill(payload)
        return PromotionGateResult(
            decision="manual_review",
            suggested_action="review_required",
            risk_level="medium",
            reasons=[f"Proposal type `{normalized_type or 'unknown'}` is not auto-promotable."],
        )

    def _evaluate_workflow(self, payload: dict[str, Any]) -> PromotionGateResult:
        static_gate = _mapping(payload.get("static_gate"))
        sandbox = _mapping(payload.get("sandbox"))
        evolution = _mapping(payload.get("evolution"))
        static_decision = str(static_gate.get("decision") or "unknown")
        sandbox_status = str(sandbox.get("status") or "unknown")
        risk_level = str(evolution.get("risk_level") or "low").strip().lower() or "low"
        issue_counts = _merged_issue_counts(static_gate, sandbox)
        reasons: list[str] = []
        decision = "pass"

        if static_decision == "reject":
            decision = "blocked"
            reasons.append("StaticGate rejected the workflow payload.")
        elif static_decision == "requires_manual_review":
            decision = "manual_review"
            reasons.append("StaticGate found issues that require review.")
        elif static_decision not in {"pass", "unknown"}:
            decision = "manual_review"
            reasons.append(f"StaticGate returned `{static_decision}`.")

        if sandbox_status == "failed":
            decision = "blocked"
            reasons.append("Sandbox evaluation failed the workflow payload.")
        elif sandbox_status == "blocked" and decision != "blocked":
            decision = "manual_review"
            reasons.append("Sandbox blocked one or more workflow steps.")
        elif sandbox_status not in {"passed", "skipped", "unknown"} and decision == "pass":
            decision = "manual_review"
            reasons.append(f"Sandbox returned `{sandbox_status}`.")

        if risk_level != "low" and decision == "pass":
            decision = "manual_review"
            reasons.append(f"Risk level `{risk_level}` requires human review.")

        auto_verify_eligible = self._workflow_auto_verify_eligible(evolution, decision=decision)
        suggested_action = "reject" if decision == "blocked" else "review_required"
        if auto_verify_eligible:
            suggested_action = "auto_apply"
            reasons.append("Workflow meets low-risk auto-verify thresholds.")
        elif decision == "pass":
            reasons.append("Workflow passed the gate but still requires review by current policy.")

        return PromotionGateResult(
            decision=decision,
            suggested_action=suggested_action,
            risk_level=risk_level,
            reasons=reasons,
            static_gate_decision=static_decision,
            sandbox_status=sandbox_status,
            auto_verify_eligible=auto_verify_eligible,
            issue_counts=issue_counts,
        )

    def _evaluate_skill(self, payload: dict[str, Any]) -> PromotionGateResult:
        static_gate = _mapping(payload.get("static_gate"))
        evolution = _mapping(payload.get("evolution"))
        static_decision = str(static_gate.get("decision") or "unknown")
        risk_level = str(evolution.get("risk_level") or "medium").strip().lower() or "medium"
        issue_counts = _merged_issue_counts(static_gate, {})
        reasons: list[str] = []
        decision = "pass"
        if static_decision == "reject":
            decision = "blocked"
            reasons.append("StaticGate rejected the skill draft.")
        elif static_decision == "requires_manual_review":
            decision = "manual_review"
            reasons.append("StaticGate found issues that require review.")
        elif static_decision not in {"pass", "unknown"}:
            decision = "manual_review"
            reasons.append(f"StaticGate returned `{static_decision}`.")
        if decision == "pass":
            reasons.append("Read-only skill draft passed the gate but must remain review-required.")
        return PromotionGateResult(
            decision=decision,
            suggested_action="reject" if decision == "blocked" else "review_required",
            risk_level=risk_level,
            reasons=reasons,
            static_gate_decision=static_decision,
            sandbox_status="not_applicable",
            auto_verify_eligible=False,
            issue_counts=issue_counts,
        )

    def _workflow_auto_verify_eligible(self, evolution: dict[str, Any], *, decision: str) -> bool:
        if decision != "pass":
            return False
        if not bool(getattr(self.config, "auto_verify_workflows", False)):
            return False
        try:
            priority = float(evolution.get("priority_score") or 0.0)
            seen_count = int(evolution.get("seen_count") or 0)
            threshold = float(getattr(self.config, "workflow_auto_verify_threshold", 0.9) or 0.9)
            min_seen = int(getattr(self.config, "workflow_auto_verify_min_seen_count", 5) or 5)
        except (TypeError, ValueError):
            return False
        return round(priority, 2) >= threshold and seen_count >= min_seen


def _mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def _merged_issue_counts(static_gate: dict[str, Any], sandbox: dict[str, Any]) -> dict[str, int]:
    counts: dict[str, int] = {}
    _merge_counts(counts, static_gate.get("issue_counts"))
    sandbox_issues = sandbox.get("issues")
    if isinstance(sandbox_issues, list):
        for issue in sandbox_issues:
            if not isinstance(issue, dict):
                continue
            severity = str(issue.get("severity") or "").strip()
            if severity:
                counts[severity] = counts.get(severity, 0) + 1
    return counts


def _merge_counts(target: dict[str, int], raw: Any) -> None:
    if not isinstance(raw, dict):
        return
    for key, value in raw.items():
        try:
            count = int(value)
        except (TypeError, ValueError):
            continue
        if count:
            normalized = str(key)
            target[normalized] = target.get(normalized, 0) + count
