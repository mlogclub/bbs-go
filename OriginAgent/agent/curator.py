"""Deterministic curator proposals for reviewed workspace artifacts."""

from __future__ import annotations

import asyncio
import hashlib
import json
import re
import uuid
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

import yaml
from loguru import logger

from OriginAgent.agent.background_review import ReviewProposal, ReviewProposalStore
from OriginAgent.agent.domain_packs import DomainPackManager
from OriginAgent.agent.evolution import (
    AUTO_EVOLUTION_ORIGIN,
    OpportunitySignalStore,
    build_skill_payload_from_signal,
    build_workflow_payload_from_signal,
    evolution_allows_skill_proposals,
    evolution_allows_workflow_proposals,
)
from OriginAgent.agent.evolution_config_overlay import apply_config_overlay
from OriginAgent.agent.evolution_outcomes import (
    EvolutionOutcomeStore,
    proposal_outcome_context,
    safe_append_outcome,
)
from OriginAgent.agent.evolution_dependencies import EvolutionDependencyStore
from OriginAgent.agent.evolution_gate import PromotionGate
from OriginAgent.agent.evolution_feedback import EvolutionFeedbackCalibrator
from OriginAgent.agent.evolution_maintenance import run_evolution_maintenance
from OriginAgent.agent.evolution_operator import build_operator_insights
from OriginAgent.agent.evolution_sandbox import SandboxEvaluator
from OriginAgent.agent.facts import CONFLICT_CATEGORIES, FactStore, normalize_fact_content
from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.agent.skill_lifecycle import _read_skill_markdown
from OriginAgent.agent.skills import SkillsLoader
from OriginAgent.agent.workflow_artifacts import validate_workflow_artifact_dir, write_workflow_artifact
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger

CURATOR_ORIGIN = "curator"
_CURATOR_SESSION_KEY = "curator:system"
_BODY_MAX_CHARS = 2400
_RATIONALE_MAX_CHARS = 1200
_TITLE_MAX_CHARS = 160
_EVIDENCE_MAX_CHARS = 500
_MAX_EVIDENCE = 5
_NORMALIZE_SPACE_RE = re.compile(r"\s+")


@dataclass(frozen=True)
class CuratorResult:
    """Runtime outcome for curator proposal generation."""

    status: str
    proposals_written: int = 0
    reason: str = ""
    evolution_candidates: int = 0
    evolution_proposals_prepared: int = 0
    evolution_dry_run: bool = True
    evolution_mode: str = "conservative"


class CuratorService:
    """Generate deterministic curator proposals into the shared review store."""

    def __init__(
        self,
        *,
        workspace: Path,
        config: Any | None = None,
        config_loader: Any | None = None,
        evolution_config: Any | None = None,
        evolution_config_loader: Any | None = None,
        domain_pack_manager: DomainPackManager | None = None,
        store: ReviewProposalStore | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self._config = config
        self._config_loader = config_loader
        self._evolution_config = evolution_config
        self._evolution_config_loader = evolution_config_loader
        self.domain_pack_manager = domain_pack_manager
        self.store = store or ReviewProposalStore(self.workspace)
        self.opportunity_signals = OpportunitySignalStore(self.workspace)
        self.outcomes = EvolutionOutcomeStore(self.workspace)
        self.dependencies = EvolutionDependencyStore(self.workspace)
        self.sandbox = SandboxEvaluator(self.workspace, self.evolution_config)
        self.promotion_gate = PromotionGate(self.evolution_config)
        self.feedback_calibrator = EvolutionFeedbackCalibrator(self.workspace, self.evolution_config)
        self._running = 0
        self._last_result: CuratorResult | None = None
        self._last_evolution_scan: dict[str, Any] = {}

    def refresh_config(self) -> None:
        if self._config_loader is not None:
            try:
                self._config = self._config_loader()
            except Exception:
                logger.exception("Failed to refresh curator config")
        if self._evolution_config_loader is not None:
            try:
                self._evolution_config = self._evolution_config_loader()
                self.sandbox = SandboxEvaluator(self.workspace, self._evolution_config)
                self.promotion_gate = PromotionGate(self._evolution_config)
                self.feedback_calibrator = EvolutionFeedbackCalibrator(self.workspace, self._evolution_config)
            except Exception:
                logger.exception("Failed to refresh evolution config")

    @property
    def config(self) -> Any:
        if self._config is None:
            from OriginAgent.config.schema import CuratorConfig

            self._config = CuratorConfig()
        return self._config

    @property
    def enabled(self) -> bool:
        return bool(getattr(self.config, "enabled", False))

    @property
    def evolution_config(self) -> Any:
        if self._evolution_config is None:
            from OriginAgent.config.schema import EvolutionConfig

            self._evolution_config = EvolutionConfig()
        return apply_config_overlay(self.workspace, self._evolution_config)

    def runtime_status(self) -> dict[str, Any]:
        stats = self.store.stats(origin=CURATOR_ORIGIN)
        return {
            "curator_enabled": self.enabled,
            "curator_running_count": self._running,
            "curator_proposal_count": stats["proposal_count"],
            "curator_pending_count": stats["pending_count"],
            "curator_last_created_at": stats["last_created_at"],
            "curator_last_result": asdict(self._last_result) if self._last_result is not None else None,
            "curator_type_counts": self.store.type_counts(origin=CURATOR_ORIGIN),
        }

    async def review_workspace(
        self,
        *,
        session_key: str,
        turn_id: str,
    ) -> CuratorResult:
        self.refresh_config()
        if not self.enabled:
            return self._remember(CuratorResult(status="skipped", reason="disabled"))
        if self._running > 0:
            return self._remember(CuratorResult(status="skipped", reason="concurrency_limit"))

        self._running += 1
        try:
            self._last_evolution_scan = {}
            self._run_feedback_calibration()
            proposals = self._build_proposals(session_key=session_key, turn_id=turn_id)
            written = await asyncio.to_thread(self.store.append_many, proposals)
            if written:
                self._trace_written_evolution_proposals(proposals)
                self._mark_evolution_proposals_converted(proposals)
            self._run_evolution_maintenance()
            scan = dict(self._last_evolution_scan)
            return self._remember(CuratorResult(
                status="ok",
                proposals_written=written,
                evolution_candidates=int(scan.get("candidates", 0) or 0),
                evolution_proposals_prepared=int(scan.get("prepared", 0) or 0),
                evolution_dry_run=bool(scan.get("dry_run", True)),
                evolution_mode=str(scan.get("mode") or "conservative"),
            ))
        except Exception as exc:
            logger.exception("Curator review failed")
            return self._remember(CuratorResult(status="error", reason=str(exc)))
        finally:
            self._running = max(0, self._running - 1)

    def _remember(self, result: CuratorResult) -> CuratorResult:
        self._last_result = result
        return result

    def _run_feedback_calibration(self) -> None:
        try:
            result = self.feedback_calibrator.run()
            self._last_evolution_scan["feedback_calibration"] = result.to_json()
        except Exception:
            logger.exception("Evolution feedback calibration failed")

    def _run_evolution_maintenance(self) -> None:
        self._last_evolution_scan["maintenance"] = run_evolution_maintenance(
            self.workspace,
            self.evolution_config,
        )

    def _build_proposals(self, *, session_key: str, turn_id: str) -> list[ReviewProposal]:
        now = datetime.now(timezone.utc).isoformat()
        limit = max(1, int(getattr(self.config, "max_proposals_per_run", 12) or 12))
        existing = self.store.iter_all()
        review_lookup = {
            str(record.get("id") or ""): record
            for record in existing
            if str(record.get("id") or "")
        }
        deduper = _ProposalDeduper(existing)
        proposals: list[ReviewProposal] = []

        def add_all(items: list[ReviewProposal]) -> None:
            for item in items:
                if len(proposals) >= limit:
                    return
                if deduper.should_add(item):
                    proposals.append(item)
                    deduper.track(item)

        skill_records = self._workspace_skill_records(review_lookup)
        add_all(self._skill_proposals(skill_records, session_key=session_key, turn_id=turn_id, created_at=now))
        if len(proposals) >= limit:
            return proposals

        workflow_records = self._workflow_records(review_lookup)
        add_all(self._workflow_proposals(workflow_records, session_key=session_key, turn_id=turn_id, created_at=now))
        if len(proposals) >= limit:
            return proposals

        add_all(self._fact_conflict_proposals(session_key=session_key, turn_id=turn_id, created_at=now))
        if len(proposals) >= limit:
            return proposals

        add_all(self._evolution_workflow_proposals(
            session_key=session_key,
            turn_id=turn_id,
            created_at=now,
            limit=limit - len(proposals),
        ))
        if len(proposals) >= limit:
            return proposals

        add_all(self._evolution_skill_proposals(
            session_key=session_key,
            turn_id=turn_id,
            created_at=now,
            limit=limit - len(proposals),
        ))
        return proposals

    def _evolution_workflow_proposals(
        self,
        *,
        session_key: str,
        turn_id: str,
        created_at: str,
        limit: int,
    ) -> list[ReviewProposal]:
        config = self.evolution_config
        mode = str(getattr(config, "mode", "conservative") or "conservative")
        dry_run = bool(getattr(config, "dry_run", True))
        signals = self.opportunity_signals.select_workflow_candidates(config, limit=max(0, limit))
        feedback_calibration = self._last_evolution_scan.get("feedback_calibration")
        self._last_evolution_scan = {
            "mode": mode,
            "dry_run": dry_run,
            "workflow_candidates": len(signals),
            "skill_candidates": int(self._last_evolution_scan.get("skill_candidates", 0) or 0),
            "candidates": len(signals) + int(self._last_evolution_scan.get("skill_candidates", 0) or 0),
            "prepared": 0,
            "workflow_prepared": 0,
            "skill_prepared": 0,
        }
        if feedback_calibration is not None:
            self._last_evolution_scan["feedback_calibration"] = feedback_calibration
        if not signals or not evolution_allows_workflow_proposals(config):
            return []

        proposals: list[ReviewProposal] = []
        sandbox = SandboxEvaluator(self.workspace, config)
        promotion_gate = PromotionGate(config)
        for signal in signals:
            payload = build_workflow_payload_from_signal(signal, config=config)
            payload["sandbox"] = sandbox.evaluate_workflow_payload(payload)
            gate = promotion_gate.evaluate(payload, proposal_type="workflow")
            payload["promotion_gate"] = gate.to_json()
            payload["operator_insights"] = build_operator_insights(
                payload,
                proposal_type="workflow",
                config=config,
            )
            if gate.decision == "blocked":
                self._trace_gate_evaluated(payload, proposal_type="workflow")
                continue
            evidence = []
            for item in signal.evidence_sources[:_MAX_EVIDENCE]:
                cursor = item.get("cursor")
                timestamp = item.get("timestamp")
                preview = str(item.get("preview") or "").strip()
                evidence.append(
                    f"cursor={cursor} timestamp={timestamp}: {preview}"
                    if preview
                    else f"cursor={cursor} timestamp={timestamp}"
                )
            proposals.append(self._proposal(
                session_key=session_key,
                turn_id=turn_id,
                created_at=created_at,
                proposal_type="workflow",
                domain_id="core",
                title=f"Create workflow from repeated pattern `{signal.target_key}`",
                content=(
                    "A high-scoring opportunity signal suggests this repeated interaction "
                    "should become a reviewed manual workflow."
                ),
                rationale=(
                    "Curator converted an auto-evolution opportunity signal into a normal "
                    "workflow review proposal. StaticGate results are attached in payload.static_gate."
                ),
                evidence=evidence,
                payload=payload,
                confidence=max(0.1, min(0.99, signal.priority_score)),
                origin=AUTO_EVOLUTION_ORIGIN,
            ))
        self._last_evolution_scan["workflow_prepared"] = len(proposals)
        self._last_evolution_scan["prepared"] = int(self._last_evolution_scan.get("prepared", 0) or 0) + len(proposals)
        return proposals

    def _evolution_skill_proposals(
        self,
        *,
        session_key: str,
        turn_id: str,
        created_at: str,
        limit: int,
    ) -> list[ReviewProposal]:
        config = self.evolution_config
        mode = str(getattr(config, "mode", "conservative") or "conservative")
        dry_run = bool(getattr(config, "dry_run", True))
        signals = self.opportunity_signals.select_skill_candidates(config, limit=max(0, limit))
        self._last_evolution_scan.update({
            "mode": mode,
            "dry_run": dry_run,
            "skill_candidates": len(signals),
            "candidates": int(self._last_evolution_scan.get("workflow_candidates", 0) or 0) + len(signals),
        })
        if not signals or not evolution_allows_skill_proposals(config):
            return []

        proposals: list[ReviewProposal] = []
        promotion_gate = PromotionGate(config)
        for signal in signals:
            payload = build_skill_payload_from_signal(signal, config=config)
            gate = promotion_gate.evaluate(payload, proposal_type="skill")
            payload["promotion_gate"] = gate.to_json()
            payload["operator_insights"] = build_operator_insights(
                payload,
                proposal_type="skill",
                config=config,
            )
            if gate.decision == "blocked":
                self._trace_gate_evaluated(payload, proposal_type="skill")
                continue
            evidence = _evidence_lines(signal.evidence_sources)
            proposals.append(self._proposal(
                session_key=session_key,
                turn_id=turn_id,
                created_at=created_at,
                proposal_type="skill",
                domain_id="core",
                title=f"Create read-only skill from repeated pattern `{signal.target_key}`",
                content=(
                    "A high-scoring opportunity signal suggests this repeated read-only pattern "
                    "could become a reviewed skill draft."
                ),
                rationale=(
                    "Curator converted an auto-evolution opportunity signal into a normal skill "
                    "review proposal. StaticGate results are attached in payload.static_gate."
                ),
                evidence=evidence,
                payload=payload,
                confidence=max(0.1, min(0.99, signal.priority_score)),
                origin=AUTO_EVOLUTION_ORIGIN,
            ))
        self._last_evolution_scan["skill_prepared"] = len(proposals)
        self._last_evolution_scan["prepared"] = int(self._last_evolution_scan.get("prepared", 0) or 0) + len(proposals)
        return proposals

    def _mark_evolution_proposals_converted(self, proposals: list[ReviewProposal]) -> None:
        for proposal in proposals:
            if proposal.origin != AUTO_EVOLUTION_ORIGIN:
                continue
            payload = proposal.payload if isinstance(proposal.payload, dict) else {}
            evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
            opportunity_id = str(evolution.get("opportunity_id") or "")
            verified = proposal.proposal_type == "workflow" and self._maybe_auto_verify_workflow(proposal)
            if opportunity_id:
                self.opportunity_signals.mark_converted(
                    opportunity_id,
                    proposal.id,
                    verification_status="verified" if verified else "",
                )

    def _trace_written_evolution_proposals(self, proposals: list[ReviewProposal]) -> None:
        for proposal in proposals:
            if proposal.origin != AUTO_EVOLUTION_ORIGIN:
                continue
            record = proposal.to_json()
            context = proposal_outcome_context(record)
            gate = payload_gate(record)
            safe_append_outcome(
                self.outcomes,
                "gate_evaluated",
                **context,
                feedback_score=proposal.confidence,
                metadata={
                    "proposal_type": proposal.proposal_type,
                    "domain_id": proposal.domain_id,
                    "origin": proposal.origin,
                    "promotion_gate": gate,
                },
            )
            safe_append_outcome(
                self.outcomes,
                "proposal_generated",
                **context,
                review_status=proposal.status,
                feedback_score=proposal.confidence,
                metadata={
                    "proposal_type": proposal.proposal_type,
                    "domain_id": proposal.domain_id,
                    "origin": proposal.origin,
                    "title": proposal.title,
                    "created_at": proposal.created_at,
                },
            )

    def _trace_gate_evaluated(self, payload: dict[str, Any], *, proposal_type: str) -> None:
        evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
        gate = payload_gate({"payload": payload})
        sandbox = payload.get("sandbox") if isinstance(payload.get("sandbox"), dict) else {}
        safe_append_outcome(
            self.outcomes,
            "gate_evaluated",
            opportunity_id=str(evolution.get("opportunity_id") or ""),
            artifact_type=proposal_type,
            artifact_name=str(payload.get("workflow_name") or payload.get("skill_name") or ""),
            artifact_path=str(payload.get("subject_path") or ""),
            gate_decision=str(gate.get("decision") or ""),
            sandbox_status=str(sandbox.get("status") or ""),
            feedback_score=_safe_float(evolution.get("priority_score")),
            metadata={
                "proposal_type": proposal_type,
                "origin": AUTO_EVOLUTION_ORIGIN,
                "promotion_gate": gate,
                "blocked_before_proposal": True,
            },
        )

    def _maybe_auto_verify_workflow(self, proposal: ReviewProposal) -> bool:
        config = self.evolution_config
        if not bool(getattr(config, "auto_verify_workflows", False)):
            return False
        if not evolution_allows_workflow_proposals(config):
            return False
        payload = proposal.payload if isinstance(proposal.payload, dict) else {}
        evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
        if str(evolution.get("origin") or "") != AUTO_EVOLUTION_ORIGIN:
            return False
        gate = payload.get("promotion_gate") if isinstance(payload.get("promotion_gate"), dict) else {}
        if not bool(gate.get("auto_verify_eligible")):
            return False

        now = datetime.now(timezone.utc).isoformat()
        record = proposal.to_json()
        metadata = {
            "verification_status": "verified",
            "created_by": AUTO_EVOLUTION_ORIGIN,
            "previous_version": None,
            "verified_by": AUTO_EVOLUTION_ORIGIN,
            "verified_at": now,
            "opportunity_id": str(evolution.get("opportunity_id") or ""),
        }
        try:
            artifact = write_workflow_artifact(record, self.workspace, metadata_overrides=metadata).to_json()
            with self.store._locked():
                event = self.store._append_event_unlocked(
                    proposal.id,
                    status="applied",
                    reason="auto_evolution verified low-risk workflow proposal",
                    artifact=artifact,
                )
            self._append_auto_promotion_event(
                proposal=proposal,
                artifact=artifact,
                review_event_id=str(event.get("event_id") or ""),
            )
            context = proposal_outcome_context(record, event)
            safe_append_outcome(
                self.outcomes,
                "promoted",
                **context,
                review_status="applied",
                promotion_status="verified",
                metadata={
                    "review_event_id": str(event.get("event_id") or ""),
                    "auto_verified": True,
                    "reason": "auto_evolution verified low-risk workflow proposal",
                },
            )
            return True
        except Exception:
            logger.exception("Auto-evolution workflow verification failed for {}", proposal.id)
            return False

    def _append_auto_promotion_event(
        self,
        *,
        proposal: ReviewProposal,
        artifact: dict[str, Any],
        review_event_id: str,
    ) -> None:
        payload = proposal.payload if isinstance(proposal.payload, dict) else {}
        evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
        workflow_name = str(artifact.get("workflow_name") or payload.get("workflow_name") or "")
        artifact_path = str(artifact.get("path") or payload.get("subject_path") or "")
        try:
            EvolutionLedger(self.workspace).append(EvolutionEvent.new(
                EventType.UNMAPPED,
                actor=AUTO_EVOLUTION_ORIGIN,
                module_id=workflow_name,
                module_type="workflow",
                source_event_stream="evolution",
                source_event_id=review_event_id or proposal.id,
                result={
                    "event_name": "evolution_auto_promotion",
                    "opportunity_id": str(evolution.get("opportunity_id") or ""),
                    "proposal_id": proposal.id,
                    "workflow_name": workflow_name,
                    "artifact_path": artifact_path,
                    "activated_by": AUTO_EVOLUTION_ORIGIN,
                    "promoted_at": datetime.now(timezone.utc).isoformat(),
                },
            ))
        except Exception:
            logger.exception("Failed to append auto-evolution promotion event for {}", proposal.id)

    def _workspace_skill_records(self, review_lookup: dict[str, dict[str, Any]]) -> list[dict[str, Any]]:
        loader = SkillsLoader(self.workspace, domain_pack_manager=self.domain_pack_manager)
        records: list[dict[str, Any]] = []
        for record in loader.list_skill_records(filter_unavailable=False):
            if str(record.get("source") or "") != "workspace":
                continue
            path = Path(str(record.get("path") or ""))
            frontmatter, body = _read_skill_markdown(path)
            frontmatter_name = str(frontmatter.get("name") or record.get("name") or "").strip()
            proposal_id = str(record.get("review_proposal_id") or "").strip()
            records.append({
                **record,
                "frontmatter_name": frontmatter_name or str(record.get("name") or ""),
                "description": str(frontmatter.get("description") or record.get("description") or record.get("name") or "").strip(),
                "body": body.strip(),
                "content_hash": _skill_content_hash(
                    frontmatter_name or str(record.get("name") or ""),
                    str(frontmatter.get("description") or record.get("description") or ""),
                    body,
                ),
                "created_at": str(review_lookup.get(proposal_id, {}).get("created_at") or ""),
                "subject_path": f"skills/{record.get('name')}/SKILL.md",
            })
        return records

    def _skill_proposals(
        self,
        records: list[dict[str, Any]],
        *,
        session_key: str,
        turn_id: str,
        created_at: str,
    ) -> list[ReviewProposal]:
        proposals: list[ReviewProposal] = []
        by_group: dict[tuple[str, str], list[dict[str, Any]]] = {}
        for record in records:
            group_key = (str(record.get("domain_id") or "core"), str(record.get("content_hash") or ""))
            by_group.setdefault(group_key, []).append(record)

        duplicate_info: dict[str, dict[str, Any]] = {}
        for group in by_group.values():
            if len(group) <= 1:
                continue
            ordered = sorted(group, key=_skill_priority_key)
            ambiguous = _skill_group_ambiguous(ordered)
            canonical = ordered[0]
            duplicate_info[str(canonical.get("name") or "")] = {
                "canonical": True,
                "ambiguous": ambiguous,
            }
            for record in ordered[1:]:
                duplicate_info[str(record.get("name") or "")] = {
                    "canonical": False,
                    "ambiguous": ambiguous,
                    "canonical_record": canonical,
                }
            if ambiguous:
                names = [str(item.get("name") or "") for item in ordered]
                proposals.append(self._proposal(
                    session_key=session_key,
                    turn_id=turn_id,
                    created_at=created_at,
                    proposal_type="merge_skill",
                    domain_id=str(canonical.get("domain_id") or "core"),
                    title=f"Review duplicate workspace skills: {', '.join(names[:2])}",
                    content=(
                        f"Multiple workspace skills in domain `{canonical.get('domain_id') or 'core'}` "
                        "have the same normalized content and cannot be safely auto-deprecated."
                    ),
                    rationale=(
                        "The duplicate skill group has matching lifecycle priority and creation ordering, "
                        "so curator cannot safely choose a single survivor in P10."
                    ),
                    evidence=[
                        f"{item.get('name')}: {item.get('subject_path')}" for item in ordered[:_MAX_EVIDENCE]
                    ],
                    payload={
                        "subject_type": "skill_group",
                        "subject_id": ",".join(names),
                        "subject_path": ", ".join(str(item.get("subject_path") or "") for item in ordered[:2]),
                        "curator_key": f"merge-skill:{canonical.get('domain_id') or 'core'}:{canonical.get('content_hash') or ''}",
                        "target_state_hash": _stable_hash(names),
                        "suggested_action": "merge_skill",
                        "impact_summary": "Manual skill merge review required.",
                        "skill_names": names,
                    },
                    confidence=0.78,
                ))

        for record in records:
            name = str(record.get("name") or "")
            lifecycle = str(record.get("lifecycle_status") or "")
            verification = str(record.get("verification_status") or "")
            created_by = str(record.get("created_by") or "")
            domain_id = str(record.get("domain_id") or "core")
            duplicate = duplicate_info.get(name)

            if (
                created_by == "background_review"
                and lifecycle == "proposed"
                and verification == "verified"
                and not (duplicate and not duplicate.get("canonical"))
                and not (duplicate and duplicate.get("ambiguous"))
            ):
                proposals.append(self._proposal(
                    session_key=session_key,
                    turn_id=turn_id,
                    created_at=created_at,
                    proposal_type="promote_skill",
                    domain_id=domain_id,
                    title=f"Promote verified skill `{name}`",
                    content=(
                        f"Workspace skill `{name}` is verified but still proposed. "
                        "Curator recommends activating it for normal discovery."
                    ),
                    rationale=(
                        "The skill is already verified and is not blocked by a stronger duplicate candidate."
                    ),
                    evidence=[
                        f"skill={name}",
                        f"lifecycle={lifecycle}",
                        f"verification={verification}",
                        str(record.get("subject_path") or ""),
                    ],
                    payload={
                        "subject_type": "skill",
                        "subject_id": name,
                        "subject_path": str(record.get("subject_path") or ""),
                        "curator_key": f"promote-skill:{name}",
                        "target_state_hash": _stable_hash([
                            name,
                            lifecycle,
                            verification,
                            str(record.get("last_lifecycle_event_id") or ""),
                        ]),
                        "suggested_action": "promote_skill",
                        "impact_summary": "Skill will become active/verified for default discovery.",
                        "skill_name": name,
                    },
                    confidence=0.84,
                ))

            if duplicate and not duplicate.get("canonical") and not duplicate.get("ambiguous"):
                canonical = dict(duplicate.get("canonical_record") or {})
                proposals.append(self._proposal(
                    session_key=session_key,
                    turn_id=turn_id,
                    created_at=created_at,
                    proposal_type="deprecate_skill",
                    domain_id=domain_id,
                    title=f"Deprecate duplicate skill `{name}`",
                    content=(
                        f"Workspace skill `{name}` duplicates `{canonical.get('name')}` "
                        "and is the weaker copy in this normalized skill group."
                    ),
                    rationale=(
                        "Curator selected a stronger canonical skill using lifecycle status, "
                        "verification status, proposal creation ordering, and skill name."
                    ),
                    evidence=[
                        f"duplicate={name}",
                        f"canonical={canonical.get('name')}",
                        str(record.get("subject_path") or ""),
                        str(canonical.get("subject_path") or ""),
                    ],
                    payload={
                        "subject_type": "skill",
                        "subject_id": name,
                        "subject_path": str(record.get("subject_path") or ""),
                        "curator_key": f"deprecate-skill:{name}",
                        "target_state_hash": _stable_hash([
                            name,
                            canonical.get("name") or "",
                            lifecycle,
                            verification,
                            str(record.get("last_lifecycle_event_id") or ""),
                        ]),
                        "suggested_action": "deprecate_skill",
                        "impact_summary": f"Skill `{name}` will be marked deprecated.",
                        "skill_name": name,
                        "canonical_skill_name": str(canonical.get("name") or ""),
                    },
                    confidence=0.88,
                ))

            if (
                lifecycle == "active"
                and verification == "verified"
                and domain_id
                and domain_id != "core"
                and self._domain_pack_available(domain_id)
            ):
                proposals.append(self._proposal(
                    session_key=session_key,
                    turn_id=turn_id,
                    created_at=created_at,
                    proposal_type="move_to_domain",
                    domain_id=domain_id,
                    title=f"Consider moving skill `{name}` into domain pack `{domain_id}`",
                    content=(
                        f"Workspace skill `{name}` is active and verified, and its domain pack "
                        f"`{domain_id}` is available."
                    ),
                    rationale=(
                        "This skill appears mature enough for domain-level governance, "
                        "but P10 keeps the move as a reviewed proposal only."
                    ),
                    evidence=[
                        str(record.get("subject_path") or ""),
                        f"domain={domain_id}",
                        f"lifecycle={lifecycle}",
                        f"verification={verification}",
                    ],
                    payload={
                        "subject_type": "skill",
                        "subject_id": name,
                        "subject_path": str(record.get("subject_path") or ""),
                        "curator_key": f"move-to-domain-skill:{name}:{domain_id}",
                        "target_state_hash": _stable_hash([
                            name,
                            domain_id,
                            lifecycle,
                            verification,
                        ]),
                        "suggested_action": "move_to_domain",
                        "impact_summary": f"Future phase could migrate `{name}` into domain pack `{domain_id}`.",
                        "skill_name": name,
                    },
                    confidence=0.73,
                ))
        return proposals

    def _workflow_records(self, review_lookup: dict[str, dict[str, Any]]) -> list[dict[str, Any]]:
        root = self.workspace / "workflows"
        records: list[dict[str, Any]] = []
        try:
            children = sorted(root.iterdir(), key=lambda path: path.name)
        except FileNotFoundError:
            return []
        except OSError:
            logger.exception("Failed to read workflow artifacts for curator")
            return []

        for child in children:
            if not child.is_dir():
                continue
            relative_path = f"workflows/{child.name}/workflow.yaml"
            valid, message = validate_workflow_artifact_dir(child, workspace=self.workspace)
            record: dict[str, Any] = {
                "name": child.name,
                "subject_path": relative_path,
                "valid": valid,
                "validation": message,
                "domain_id": "core",
                "created_by": "",
                "review_proposal_id": "",
                "created_at": "",
                "body": "",
                "steps": [],
                "content_hash": "",
            }
            if not valid:
                records.append(record)
                continue
            try:
                raw = yaml.safe_load((child / "workflow.yaml").read_text(encoding="utf-8"))
            except (OSError, yaml.YAMLError):
                records.append(record)
                continue
            if not isinstance(raw, dict):
                records.append(record)
                continue
            metadata = raw.get("metadata")
            originagent = metadata.get("OriginAgent") if isinstance(metadata, dict) else {}
            if not isinstance(originagent, dict):
                originagent = {}
            proposal_id = str(originagent.get("review_proposal_id") or "")
            body = str(raw.get("body") or "")
            steps = raw.get("steps") if isinstance(raw.get("steps"), list) else []
            record.update({
                "name": str(raw.get("name") or child.name),
                "domain_id": str(originagent.get("domain_id") or "core"),
                "created_by": str(originagent.get("created_by") or ""),
                "review_proposal_id": proposal_id,
                "created_at": str(review_lookup.get(proposal_id, {}).get("created_at") or ""),
                "body": body,
                "steps": steps,
                "content_hash": _workflow_content_hash(body, steps),
            })
            records.append(record)
        return records

    def _workflow_proposals(
        self,
        records: list[dict[str, Any]],
        *,
        session_key: str,
        turn_id: str,
        created_at: str,
    ) -> list[ReviewProposal]:
        proposals: list[ReviewProposal] = []
        valid_records = [record for record in records if record.get("valid")]
        for record in records:
            if record.get("valid"):
                continue
            name = str(record.get("name") or "unknown-workflow")
            proposals.append(self._proposal(
                session_key=session_key,
                turn_id=turn_id,
                created_at=created_at,
                proposal_type="archive_workflow",
                domain_id=str(record.get("domain_id") or "core"),
                title=f"Archive invalid workflow `{name}`",
                content=f"Workflow `{name}` no longer passes artifact validation.",
                rationale=(
                    "Curator found an invalid workspace workflow artifact. "
                    "P10 keeps this as a manual archive review."
                ),
                evidence=[
                    str(record.get("subject_path") or ""),
                    str(record.get("validation") or ""),
                ],
                payload={
                    "subject_type": "workflow",
                    "subject_id": name,
                    "subject_path": str(record.get("subject_path") or ""),
                    "curator_key": f"archive-workflow-invalid:{name}",
                    "target_state_hash": _stable_hash([
                        name,
                        str(record.get("validation") or ""),
                    ]),
                    "suggested_action": "archive_workflow",
                    "impact_summary": "Manual archive review for an invalid workflow artifact.",
                    "workflow_name": name,
                },
                confidence=0.83,
            ))

        groups: dict[tuple[str, str], list[dict[str, Any]]] = {}
        for record in valid_records:
            groups.setdefault(
                (str(record.get("domain_id") or "core"), str(record.get("content_hash") or "")),
                [],
            ).append(record)
        for group in groups.values():
            if len(group) <= 1:
                continue
            ordered = sorted(group, key=_workflow_priority_key)
            canonical = ordered[0]
            for record in ordered[1:]:
                name = str(record.get("name") or "")
                proposals.append(self._proposal(
                    session_key=session_key,
                    turn_id=turn_id,
                    created_at=created_at,
                    proposal_type="archive_workflow",
                    domain_id=str(record.get("domain_id") or "core"),
                    title=f"Archive duplicate workflow `{name}`",
                    content=(
                        f"Workflow `{name}` duplicates `{canonical.get('name')}` "
                        "and is not the canonical artifact in this normalized group."
                    ),
                    rationale=(
                        "Curator found multiple workflow artifacts with the same normalized "
                        "manual guide content."
                    ),
                    evidence=[
                        str(record.get("subject_path") or ""),
                        str(canonical.get("subject_path") or ""),
                    ],
                    payload={
                        "subject_type": "workflow",
                        "subject_id": name,
                        "subject_path": str(record.get("subject_path") or ""),
                        "curator_key": f"archive-workflow-duplicate:{name}",
                        "target_state_hash": _stable_hash([
                            name,
                            str(canonical.get("name") or ""),
                            str(record.get("content_hash") or ""),
                        ]),
                        "suggested_action": "archive_workflow",
                        "impact_summary": f"Workflow `{name}` appears redundant next to `{canonical.get('name')}`.",
                        "workflow_name": name,
                    },
                    confidence=0.8,
                ))

        for record in valid_records:
            domain_id = str(record.get("domain_id") or "core")
            if domain_id and domain_id != "core" and self._domain_pack_available(domain_id):
                name = str(record.get("name") or "")
                proposals.append(self._proposal(
                    session_key=session_key,
                    turn_id=turn_id,
                    created_at=created_at,
                    proposal_type="move_to_domain",
                    domain_id=domain_id,
                    title=f"Consider moving workflow `{name}` into domain pack `{domain_id}`",
                    content=(
                        f"Workflow `{name}` belongs to domain `{domain_id}` and that domain pack is available."
                    ),
                    rationale=(
                        "This workflow could later be governed by the domain pack ecosystem, "
                        "but P10 only records the suggestion."
                    ),
                    evidence=[
                        str(record.get("subject_path") or ""),
                        f"domain={domain_id}",
                    ],
                    payload={
                        "subject_type": "workflow",
                        "subject_id": name,
                        "subject_path": str(record.get("subject_path") or ""),
                        "curator_key": f"move-to-domain-workflow:{name}:{domain_id}",
                        "target_state_hash": _stable_hash([
                            name,
                            domain_id,
                            str(record.get("content_hash") or ""),
                        ]),
                        "suggested_action": "move_to_domain",
                        "impact_summary": f"Future phase could move `{name}` into domain pack `{domain_id}`.",
                        "workflow_name": name,
                    },
                    confidence=0.7,
                ))
        return proposals

    def _fact_conflict_proposals(
        self,
        *,
        session_key: str,
        turn_id: str,
        created_at: str,
    ) -> list[ReviewProposal]:
        store = FactStore(self.workspace)
        groups: dict[tuple[str, str, str], list[Any]] = {}
        for record in store.read_all():
            if record.status in {"deprecated", "contradicted"}:
                continue
            if record.category not in CONFLICT_CATEGORIES:
                continue
            groups.setdefault((record.scope, record.owner, record.category), []).append(record)

        proposals: list[ReviewProposal] = []
        for (scope, owner, category), records in groups.items():
            if len(records) <= 1:
                continue
            normalized = {normalize_fact_content(record.content) for record in records}
            if len(normalized) <= 1:
                continue
            ids = {record.fact_id for record in records}
            if any(record.supersedes_fact_id in ids for record in records if record.supersedes_fact_id):
                continue
            summary = ", ".join(record.fact_id for record in records[:3])
            proposals.append(self._proposal(
                session_key=session_key,
                turn_id=turn_id,
                created_at=created_at,
                proposal_type="fact_conflict",
                domain_id="core",
                title=f"Review conflicting facts for {category}/{scope}",
                content=(
                    f"OriginAgent found multiple active facts for `{category}` in scope `{scope}` "
                    "with conflicting normalized content."
                ),
                rationale=(
                    "Curator detected a conflict-prone fact group without an explicit supersedes relationship."
                ),
                evidence=[
                    f"{record.fact_id}: {redact_memory_text(record.content)}"[:_EVIDENCE_MAX_CHARS]
                    for record in records[:_MAX_EVIDENCE]
                ],
                payload={
                    "subject_type": "fact_group",
                    "subject_id": f"{scope}|{owner}|{category}",
                    "subject_path": "memory/facts.jsonl",
                    "curator_key": f"fact-conflict:{scope}:{owner}:{category}",
                    "target_state_hash": _stable_hash([
                        scope,
                        owner,
                        category,
                        *sorted(ids),
                        *sorted(normalized),
                    ]),
                    "suggested_action": "fact_conflict",
                    "impact_summary": f"Manual fact conflict review for {summary}.",
                    "fact_ids": sorted(ids),
                },
                confidence=0.82,
            ))
        return proposals

    def _proposal(
        self,
        *,
        session_key: str,
        turn_id: str,
        created_at: str,
        proposal_type: str,
        domain_id: str,
        title: str,
        content: str,
        rationale: str,
        evidence: list[str],
        payload: dict[str, Any],
        confidence: float,
        origin: str = CURATOR_ORIGIN,
    ) -> ReviewProposal:
        return ReviewProposal(
            id=f"review_{uuid.uuid4().hex}",
            origin=origin,
            created_at=created_at,
            session_key=session_key or _CURATOR_SESSION_KEY,
            turn_id=turn_id,
            proposal_type=proposal_type,
            domain_id=domain_id or "core",
            title=_clean_text(title, _TITLE_MAX_CHARS),
            content=_clean_text(content, _BODY_MAX_CHARS),
            rationale=_clean_text(rationale, _RATIONALE_MAX_CHARS),
            confidence=confidence,
            evidence=_clean_evidence(evidence),
            payload=_redact_payload(payload),
        )

    def _domain_pack_available(self, domain_id: str) -> bool:
        if not domain_id or domain_id == "core" or self.domain_pack_manager is None:
            return False
        pack = self.domain_pack_manager.get_pack(domain_id)
        return bool(pack is not None and pack.available)


class _ProposalDeduper:
    _ORIGINS = {CURATOR_ORIGIN, AUTO_EVOLUTION_ORIGIN}

    def __init__(self, records: list[dict[str, Any]]) -> None:
        self._pending: set[tuple[str, str]] = set()
        self._terminal: set[tuple[str, str, str]] = set()
        for record in records:
            if str(record.get("origin") or "background_review") not in self._ORIGINS:
                continue
            payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
            curator_key = str(payload.get("curator_key") or "").strip()
            target_state_hash = str(payload.get("target_state_hash") or "").strip()
            proposal_type = str(record.get("proposal_type") or record.get("type") or "").strip().lower()
            if not proposal_type or not curator_key:
                continue
            if str(record.get("status") or "pending") == "pending":
                self._pending.add((proposal_type, curator_key))
            elif target_state_hash:
                self._terminal.add((proposal_type, curator_key, target_state_hash))

    def should_add(self, proposal: ReviewProposal) -> bool:
        payload = proposal.payload if isinstance(proposal.payload, dict) else {}
        curator_key = str(payload.get("curator_key") or "").strip()
        target_state_hash = str(payload.get("target_state_hash") or "").strip()
        proposal_type = str(proposal.proposal_type or "").strip().lower()
        if not proposal_type or not curator_key or not target_state_hash:
            return False
        if (proposal_type, curator_key) in self._pending:
            return False
        if (proposal_type, curator_key, target_state_hash) in self._terminal:
            return False
        return True

    def track(self, proposal: ReviewProposal) -> None:
        payload = proposal.payload if isinstance(proposal.payload, dict) else {}
        curator_key = str(payload.get("curator_key") or "").strip()
        proposal_type = str(proposal.proposal_type or "").strip().lower()
        if proposal_type and curator_key:
            self._pending.add((proposal_type, curator_key))


def _clean_text(text: str, max_chars: int) -> str:
    return redact_memory_text(str(text or "").strip())[:max_chars].strip()


def _clean_evidence(items: list[str]) -> list[str]:
    cleaned: list[str] = []
    for item in items:
        text = _clean_text(item, _EVIDENCE_MAX_CHARS)
        if text:
            cleaned.append(text)
        if len(cleaned) >= _MAX_EVIDENCE:
            break
    return cleaned


def _evidence_lines(items: list[dict[str, Any]]) -> list[str]:
    evidence: list[str] = []
    for item in items[:_MAX_EVIDENCE]:
        cursor = item.get("cursor")
        timestamp = item.get("timestamp")
        preview = str(item.get("preview") or "").strip()
        evidence.append(
            f"cursor={cursor} timestamp={timestamp}: {preview}"
            if preview
            else f"cursor={cursor} timestamp={timestamp}"
        )
    return evidence


def _redact_payload(value: Any) -> Any:
    if isinstance(value, str):
        return redact_memory_text(value)
    if isinstance(value, list):
        return [_redact_payload(item) for item in value]
    if isinstance(value, dict):
        return {str(key): _redact_payload(item) for key, item in value.items()}
    return value


def payload_gate(record: dict[str, Any]) -> dict[str, Any]:
    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    gate = payload.get("promotion_gate") if isinstance(payload.get("promotion_gate"), dict) else {}
    return dict(gate)


def _safe_float(value: Any) -> float | None:
    try:
        return float(value)
    except (TypeError, ValueError):
        return None


def _skill_content_hash(name: str, description: str, body: str) -> str:
    payload = "\0".join([
        _normalize_text(name),
        _normalize_text(description),
        _normalize_text(body),
    ])
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def _workflow_content_hash(body: str, steps: list[Any]) -> str:
    normalized_steps = []
    for step in steps:
        if not isinstance(step, dict):
            continue
        normalized_steps.append({
            "title": _normalize_text(str(step.get("title") or "")),
            "instruction": _normalize_text(str(step.get("instruction") or "")),
            "risk": _normalize_text(str(step.get("risk") or "")),
            "confirmation_required": bool(step.get("confirmation_required")),
        })
    payload = json.dumps(
        {
            "body": _normalize_text(body),
            "steps": normalized_steps,
        },
        sort_keys=True,
        ensure_ascii=False,
    )
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def _normalize_text(text: str) -> str:
    redacted = redact_memory_text(text or "")
    return _NORMALIZE_SPACE_RE.sub(" ", redacted.strip().casefold())


def _stable_hash(values: list[Any]) -> str:
    payload = json.dumps(values, sort_keys=True, ensure_ascii=False, separators=(",", ":"))
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def _skill_priority_key(record: dict[str, Any]) -> tuple[int, int, str, str]:
    return (
        _skill_rank(str(record.get("lifecycle_status") or ""), str(record.get("verification_status") or "")),
        1 if not str(record.get("created_at") or "") else 0,
        str(record.get("created_at") or ""),
        str(record.get("name") or ""),
    )


def _skill_group_ambiguous(ordered: list[dict[str, Any]]) -> bool:
    if len(ordered) < 2:
        return False
    first = ordered[0]
    second = ordered[1]
    return (
        _skill_rank(str(first.get("lifecycle_status") or ""), str(first.get("verification_status") or ""))
        == _skill_rank(str(second.get("lifecycle_status") or ""), str(second.get("verification_status") or ""))
        and str(first.get("created_at") or "") == str(second.get("created_at") or "")
    )


def _skill_rank(lifecycle: str, verification: str) -> int:
    if lifecycle == "active" and verification == "verified":
        return 0
    if lifecycle == "proposed" and verification == "verified":
        return 1
    if lifecycle == "proposed" and verification == "unverified":
        return 2
    if lifecycle == "deprecated":
        return 3
    if lifecycle == "rejected":
        return 4
    return 5


def _workflow_priority_key(record: dict[str, Any]) -> tuple[int, str, str]:
    return (
        1 if not str(record.get("created_at") or "") else 0,
        str(record.get("created_at") or ""),
        str(record.get("name") or ""),
    )
