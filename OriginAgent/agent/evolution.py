"""Lightweight governed self-evolution opportunity signals."""

from __future__ import annotations

import hashlib
import json
import os
import re
from contextlib import suppress
from dataclasses import asdict, dataclass, field
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

from filelock import FileLock

from OriginAgent.agent.facts import ValidationIssue
from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore, safe_append_outcome
from OriginAgent.utils.helpers import ensure_dir, truncate_text


SIGNAL_KIND_WORKFLOW = "workflow_candidate"
SIGNAL_KIND_SKILL = "skill_candidate"
DEFAULT_EVOLUTION_MODE = "conservative"
DEFAULT_EVOLUTION_DRY_RUN = True
AUTO_EVOLUTION_ORIGIN = "auto_evolution"
_EVOLUTION_ALLOWED_PROPOSAL_MODES = {"curated", "exploratory", "aggressive"}
_DANGEROUS_WORKFLOW_TERMS = (
    "exec",
    "shell",
    "command",
    "write_file",
    "edit_file",
    "message",
    "cron",
    "spawn",
    "delete",
    "remove",
    "rm ",
    "powershell",
    "cmd.exe",
    "发送消息",
    "定时任务",
    "删除",
    "写入文件",
)
_DANGEROUS_WORKFLOW_WORD_RE = re.compile(
    r"(?i)(?<![a-z0-9_])"
    r"(?:exec|shell|command|write_file|edit_file|message|cron|spawn|delete|remove|rm|powershell|cmd\.exe)"
    r"(?![a-z0-9_])"
)
_DANGEROUS_WORKFLOW_CJK_TERMS = tuple(
    term for term in _DANGEROUS_WORKFLOW_TERMS if re.search(r"[\u4e00-\u9fff]", term)
)

# Governed evolution stays proposal-only here; activation remains in review/apply gates.
_SIGNAL_PREVIEW_MAX_CHARS = 500
_TARGET_MAX_CHARS = 80
_WORKFLOW_CONFIDENCE = 0.65
_MAX_EXPECTED_SEEN_COUNT = 5
_MAX_EVIDENCE_PER_SIGNAL = 12
_WORKFLOW_CUE_RE = re.compile(
    r"(?i)("
    r"workflow|repeat|repeated|recurring|every time|each time|automate|automation|"
    r"deploy|build|test|check|review|release|backup|sync|"
    r"流程|步骤|重复|反复|每次|自动|自动化|部署|构建|测试|检查|审核|发布|同步"
    r")"
)
_SKILL_CUE_RE = re.compile(
    r"(?i)("
    r"skill|reusable skill|troubleshoot|troubleshooting|diagnose|diagnostic|playbook|"
    r"analyze|analysis helper|codify|teach the agent|"
    r"技能|可复用技能|排查|诊断|分析助手|沉淀成技能|教会 agent|教会智能体"
    r")"
)
_SKILL_BODY_MAX_CHARS = 5000
_SKILL_CONFIDENCE = 0.7
_DANGEROUS_SKILL_TOOL_RE = re.compile(
    r"(?i)(?<![a-z0-9_])"
    r"(?:exec|shell|command|write_file|edit_file|message|cron|spawn|curl|wget|nc|telnet)"
    r"(?![a-z0-9_])"
)
_DANGEROUS_SKILL_INSTALL_RE = re.compile(
    r"(?i)(?:pip\s+install|npm\s+install|apt-get|apt\s+install|pnpm\s+install|yarn\s+add)"
)
_DANGEROUS_SKILL_COMMAND_HEADING_RE = re.compile(r"(?im)^\s*#\s*command\s*:")
_NOISE_RE = re.compile(r"(?i)\b(event|some event|remember this|user prefers)\b")
_SPACE_RE = re.compile(r"\s+")


@dataclass
class OpportunitySignal:
    """A non-actionable observation that may become a future reviewed proposal."""

    opportunity_id: str
    kind: str
    target_key: str
    title: str
    summary: str
    evidence_sources: list[dict[str, Any]] = field(default_factory=list)
    first_seen_at: str = ""
    last_seen_at: str = ""
    seen_count: int = 0
    priority_score: float = 0.0
    risk_level: str = "low"
    status: str = "open"
    converted_proposal_id: str | None = None
    verification_status: str = ""
    feedback_multiplier: float = 1.0
    feedback_score_offset: float = 0.0
    feedback_negative_count: int = 0
    feedback_positive_count: int = 0
    last_feedback_at: str = ""
    suppression_reason: str = ""

    @classmethod
    def from_record(cls, record: dict[str, Any]) -> "OpportunitySignal":
        return cls(
            opportunity_id=str(record.get("opportunity_id") or ""),
            kind=str(record.get("kind") or SIGNAL_KIND_WORKFLOW),
            target_key=str(record.get("target_key") or ""),
            title=str(record.get("title") or ""),
            summary=str(record.get("summary") or ""),
            evidence_sources=[
                item for item in record.get("evidence_sources") or [] if isinstance(item, dict)
            ],
            first_seen_at=str(record.get("first_seen_at") or ""),
            last_seen_at=str(record.get("last_seen_at") or ""),
            seen_count=_safe_int(record.get("seen_count"), 0),
            priority_score=_safe_float(record.get("priority_score"), 0.0),
            risk_level=str(record.get("risk_level") or "low"),
            status=str(record.get("status") or "open"),
            converted_proposal_id=(
                str(record.get("converted_proposal_id"))
                if record.get("converted_proposal_id") is not None
                else None
            ),
            verification_status=str(record.get("verification_status") or ""),
            feedback_multiplier=_safe_float(record.get("feedback_multiplier"), 1.0),
            feedback_score_offset=_safe_float(record.get("feedback_score_offset"), 0.0),
            feedback_negative_count=_safe_int(record.get("feedback_negative_count"), 0),
            feedback_positive_count=_safe_int(record.get("feedback_positive_count"), 0),
            last_feedback_at=str(record.get("last_feedback_at") or ""),
            suppression_reason=str(record.get("suppression_reason") or ""),
        )

    def to_record(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class OpportunitySignalCandidate:
    """A freshly detected observation before it is merged into the store."""

    kind: str
    target_key: str
    title: str
    summary: str
    evidence_sources: list[dict[str, Any]]
    risk_level: str = "low"

    @property
    def opportunity_id(self) -> str:
        target = _redact_text(self.target_key)
        digest = hashlib.sha256(f"{self.kind}:{target}".encode("utf-8")).hexdigest()
        return f"{self.kind}:{digest[:16]}"


class OpportunitySignalStore:
    """JSONL store for non-actionable self-evolution signals."""

    def __init__(self, workspace: Path):
        self.workspace = Path(workspace)
        self.memory_dir = ensure_dir(self.workspace / "memory")
        self.path = self.memory_dir / "opportunity_signals.jsonl"
        self._lock = FileLock(str(self.memory_dir / ".opportunity_signals.lock"))

    def read_all(self) -> list[OpportunitySignal]:
        records: list[OpportunitySignal] = []
        with suppress(FileNotFoundError):
            with open(self.path, "r", encoding="utf-8") as f:
                for line in f:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        raw = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(raw, dict):
                        signal = OpportunitySignal.from_record(raw)
                        if signal.opportunity_id:
                            records.append(signal)
        return records

    def upsert_candidates(
        self,
        candidates: list[OpportunitySignalCandidate],
        *,
        retention_days: int = 30,
        now: datetime | None = None,
    ) -> list[OpportunitySignal]:
        if not candidates:
            return []
        now_dt = _normalize_datetime(now)
        now_iso = now_dt.isoformat()
        with self._lock:
            by_id = {signal.opportunity_id: signal for signal in self.read_all()}
            changed: list[OpportunitySignal] = []
            created_ids: set[str] = set()
            for candidate in candidates:
                signal = by_id.get(candidate.opportunity_id)
                if signal is None:
                    signal = OpportunitySignal(
                        opportunity_id=candidate.opportunity_id,
                        kind=candidate.kind,
                        target_key=_clean_signal_text(candidate.target_key, _TARGET_MAX_CHARS),
                        title=_clean_signal_text(candidate.title, 160),
                        summary=_clean_signal_text(candidate.summary, 512),
                        evidence_sources=[],
                        first_seen_at=now_iso,
                        last_seen_at=now_iso,
                        risk_level=candidate.risk_level,
                    )
                    by_id[signal.opportunity_id] = signal
                    created_ids.add(signal.opportunity_id)
                signal.target_key = _clean_signal_text(candidate.target_key, _TARGET_MAX_CHARS) or signal.target_key
                signal.title = _clean_signal_text(candidate.title, 160) or signal.title
                signal.summary = _clean_signal_text(candidate.summary, 512) or signal.summary
                signal.risk_level = candidate.risk_level or signal.risk_level
                signal.last_seen_at = now_iso
                signal.evidence_sources = _merge_evidence_sources(
                    signal.evidence_sources,
                    candidate.evidence_sources,
                )
                signal.seen_count = len(signal.evidence_sources)
                signal.priority_score = _priority_score(signal)
                changed.append(signal)

            records = self._retained_records(list(by_id.values()), retention_days, now_dt)
            self._write_all_unlocked(records)
        self._trace_signal_changes(changed, created_ids=created_ids)
        return changed

    def select_workflow_candidates(
        self,
        config: Any | None = None,
        *,
        limit: int | None = None,
    ) -> list[OpportunitySignal]:
        threshold = _config_float(config, "workflow_priority_threshold", 0.7)
        min_seen = _config_int(config, "workflow_min_seen_count", 3)
        min_evidence = _config_int(config, "workflow_min_evidence_sources", 2)
        effective_limit = (
            max(0, limit)
            if limit is not None
            else max(0, _config_int(config, "max_proposals_per_cycle", 3))
        )
        if effective_limit <= 0:
            return []
        candidates = [
            signal for signal in self.read_all()
            if (
                signal.status == "open"
                and signal.kind == SIGNAL_KIND_WORKFLOW
                and signal.converted_proposal_id is None
                and signal.seen_count >= min_seen
                and len(signal.evidence_sources) >= min_evidence
                and signal.priority_score >= threshold
            )
        ]
        candidates.sort(key=lambda signal: (signal.priority_score, signal.last_seen_at), reverse=True)
        return candidates[:effective_limit]

    def select_skill_candidates(
        self,
        config: Any | None = None,
        *,
        limit: int | None = None,
    ) -> list[OpportunitySignal]:
        threshold = _config_float(config, "skill_priority_threshold", 0.85)
        min_seen = _config_int(config, "skill_min_seen_count", 5)
        effective_limit = (
            max(0, limit)
            if limit is not None
            else max(0, _config_int(config, "max_skill_proposals_per_cycle", 1))
        )
        if effective_limit <= 0:
            return []
        candidates = [
            signal for signal in self.read_all()
            if (
                signal.status == "open"
                and signal.kind == SIGNAL_KIND_SKILL
                and signal.converted_proposal_id is None
                and signal.seen_count >= min_seen
                and signal.priority_score >= threshold
            )
        ]
        candidates.sort(key=lambda signal: (signal.priority_score, signal.last_seen_at), reverse=True)
        return candidates[:effective_limit]

    def mark_converted(
        self,
        opportunity_id: str,
        proposal_id: str,
        *,
        verification_status: str = "",
    ) -> bool:
        if not opportunity_id or not proposal_id:
            return False
        with self._lock:
            records = self.read_all()
            changed = False
            now_iso = datetime.now(timezone.utc).isoformat()
            for signal in records:
                if signal.opportunity_id != opportunity_id:
                    continue
                signal.status = "converted"
                signal.converted_proposal_id = proposal_id
                if verification_status:
                    signal.verification_status = verification_status
                signal.last_seen_at = now_iso
                changed = True
                break
            if changed:
                self._write_all_unlocked(records)
            return changed

    def suppress_signal(
        self,
        opportunity_id: str,
        *,
        reason: str = "",
        now: datetime | None = None,
    ) -> OpportunitySignal | None:
        """Manually suppress an open opportunity signal."""

        if not opportunity_id:
            return None
        now_dt = _normalize_datetime(now)
        with self._lock:
            records = self.read_all()
            updated: OpportunitySignal | None = None
            for signal in records:
                if signal.opportunity_id != opportunity_id or signal.status == "converted":
                    continue
                signal.status = "suppressed"
                signal.suppression_reason = _clean_signal_text(
                    reason or "Manual evolution control-plane suppression.",
                    512,
                )
                signal.verification_status = "manual_suppressed"
                signal.last_feedback_at = now_dt.isoformat()
                updated = signal
                break
            if updated is None:
                return None
            self._write_all_unlocked(records)
        safe_append_outcome(
            EvolutionOutcomeStore(self.workspace),
            "signal_suppressed",
            opportunity_id=updated.opportunity_id,
            artifact_type=updated.kind,
            feedback_score=updated.priority_score,
            metadata={
                "origin": AUTO_EVOLUTION_ORIGIN,
                "reason": updated.suppression_reason,
                "source": "manual_override",
            },
            timestamp=now_dt,
        )
        return updated

    def resume_signal(
        self,
        opportunity_id: str,
        *,
        reason: str = "",
        now: datetime | None = None,
    ) -> OpportunitySignal | None:
        """Reopen a manually or feedback-suppressed opportunity signal."""

        if not opportunity_id:
            return None
        now_dt = _normalize_datetime(now)
        with self._lock:
            records = self.read_all()
            updated: OpportunitySignal | None = None
            for signal in records:
                if signal.opportunity_id != opportunity_id or signal.status != "suppressed":
                    continue
                signal.status = "open"
                signal.suppression_reason = ""
                signal.verification_status = "manual_resumed"
                signal.last_feedback_at = now_dt.isoformat()
                updated = signal
                break
            if updated is None:
                return None
            self._write_all_unlocked(records)
        safe_append_outcome(
            EvolutionOutcomeStore(self.workspace),
            "signal_resumed",
            opportunity_id=updated.opportunity_id,
            artifact_type=updated.kind,
            feedback_score=updated.priority_score,
            metadata={
                "origin": AUTO_EVOLUTION_ORIGIN,
                "reason": _clean_signal_text(reason or "Manual evolution control-plane resume.", 512),
                "source": "manual_override",
            },
            timestamp=now_dt,
        )
        return updated

    def apply_feedback(
        self,
        opportunity_id: str,
        *,
        multiplier: float = 1.0,
        delta: float = 0.0,
        suppress: bool = False,
        suppression_reason: str = "",
        suppress_after_negative_count: int = 0,
        risk_level: str = "",
        verification_status: str = "",
        now: datetime | None = None,
    ) -> OpportunitySignal | None:
        """Apply durable feedback calibration to one opportunity signal.

        The multiplier and offset are persisted so future evidence upserts keep
        the calibration instead of recalculating the signal back to its raw
        heuristic score.
        """

        if not opportunity_id:
            return None
        now_iso = _normalize_datetime(now).isoformat()
        with self._lock:
            records = self.read_all()
            updated: OpportunitySignal | None = None
            for signal in records:
                if signal.opportunity_id != opportunity_id:
                    continue
                safe_multiplier = max(0.05, min(2.0, _safe_float(multiplier, 1.0)))
                safe_delta = max(-1.0, min(1.0, _safe_float(delta, 0.0)))
                signal.feedback_multiplier = max(
                    0.05,
                    min(2.0, signal.feedback_multiplier * safe_multiplier),
                )
                signal.feedback_score_offset = max(
                    -1.0,
                    min(1.0, signal.feedback_score_offset + safe_delta),
                )
                if safe_multiplier < 1.0 or safe_delta < 0.0 or suppress:
                    signal.feedback_negative_count += 1
                elif safe_multiplier > 1.0 or safe_delta > 0.0:
                    signal.feedback_positive_count += 1
                threshold = max(0, _safe_int(suppress_after_negative_count, 0))
                if threshold > 0 and signal.feedback_negative_count >= threshold:
                    suppress = True
                    if not suppression_reason:
                        suppression_reason = "Repeated negative feedback reached suppression threshold."
                if suppress:
                    signal.status = "suppressed"
                    signal.suppression_reason = _clean_signal_text(suppression_reason, 512)
                if risk_level:
                    signal.risk_level = _clean_signal_text(risk_level, 64)
                if verification_status:
                    signal.verification_status = _clean_signal_text(verification_status, 128)
                signal.priority_score = _priority_score(signal)
                signal.last_feedback_at = now_iso
                updated = signal
                break
            if updated is None:
                return None
            self._write_all_unlocked(records)
            return updated

    def runtime_status(self, config: Any | None = None) -> dict[str, Any]:
        signals = self.read_all()
        open_signals = [signal for signal in signals if signal.status == "open"]
        converted_signals = [signal for signal in signals if signal.status == "converted"]
        suppressed_signals = [signal for signal in signals if signal.status == "suppressed"]
        feedback_adjusted = [
            signal for signal in signals
            if (
                signal.feedback_negative_count > 0
                or signal.feedback_positive_count > 0
                or signal.feedback_multiplier != 1.0
                or signal.feedback_score_offset != 0.0
            )
        ]
        threshold = _config_float(config, "workflow_priority_threshold", 0.7)
        min_seen = _config_int(config, "workflow_min_seen_count", 3)
        min_evidence = _config_int(config, "workflow_min_evidence_sources", 2)
        max_high = _config_int(config, "max_high_score_signals", 5)
        high = [
            signal for signal in open_signals
            if (
                signal.kind == SIGNAL_KIND_WORKFLOW
                and signal.seen_count >= min_seen
                and len(signal.evidence_sources) >= min_evidence
                and signal.priority_score >= threshold
            )
        ]
        eligible_skill_signals = [
            signal for signal in open_signals
            if (
                signal.kind == SIGNAL_KIND_SKILL
                and signal.seen_count >= _config_int(config, "skill_min_seen_count", 5)
                and signal.priority_score >= _config_float(config, "skill_priority_threshold", 0.85)
            )
        ]
        high.sort(key=lambda signal: (signal.priority_score, signal.last_seen_at), reverse=True)
        return {
            "mode": _config_str(config, "mode", DEFAULT_EVOLUTION_MODE),
            "dry_run": _config_bool(config, "dry_run", DEFAULT_EVOLUTION_DRY_RUN),
            "opportunity_signals_count": len(open_signals),
            "eligible_workflow_signals": len(high),
            "eligible_skill_signals": len(eligible_skill_signals),
            "skill_candidates_enabled": _config_bool(config, "skill_candidates_enabled", False),
            "converted_signals_count": len(converted_signals),
            "suppressed_signals_count": len(suppressed_signals),
            "feedback_adjusted_signals_count": len(feedback_adjusted),
            "feedback_negative_signals_count": len([
                signal for signal in feedback_adjusted if signal.feedback_negative_count > 0
            ]),
            "feedback_positive_signals_count": len([
                signal for signal in feedback_adjusted if signal.feedback_positive_count > 0
            ]),
            "pending_proposals_from_evolution": 0,
            "high_score_signals": [
                {
                    "kind": signal.kind,
                    "target": signal.target_key,
                    "priority_score": round(signal.priority_score, 3),
                }
                for signal in high[:max_high]
            ],
            "skill_auto_evolution_note": (
                ""
                if _config_bool(config, "skill_candidates_enabled", False)
                else "skill auto-evolution disabled, set evolution.skill_candidates_enabled=true to enable"
            ),
        }

    def _retained_records(
        self,
        records: list[OpportunitySignal],
        retention_days: int,
        now: datetime,
    ) -> list[OpportunitySignal]:
        cutoff = now - timedelta(days=max(1, retention_days))
        retained: list[OpportunitySignal] = []
        for signal in records:
            if signal.status in {"converted", "archived"}:
                retained.append(signal)
                continue
            last_seen = _parse_datetime(signal.last_seen_at)
            if last_seen is None or last_seen >= cutoff or signal.priority_score >= 0.7:
                retained.append(signal)
        return sorted(retained, key=lambda signal: (signal.kind, signal.target_key))

    def _write_all_unlocked(self, records: list[OpportunitySignal]) -> None:
        tmp_path = self.path.with_suffix(self.path.suffix + ".tmp")
        try:
            with open(tmp_path, "w", encoding="utf-8") as f:
                for signal in records:
                    f.write(json.dumps(signal.to_record(), ensure_ascii=False, sort_keys=True) + "\n")
                f.flush()
                os.fsync(f.fileno())
            os.replace(tmp_path, self.path)
        except BaseException:
            tmp_path.unlink(missing_ok=True)
            raise

    def _trace_signal_changes(self, signals: list[OpportunitySignal], *, created_ids: set[str]) -> None:
        if not signals:
            return
        outcomes = EvolutionOutcomeStore(self.workspace)
        for signal in signals:
            event_type = "signal_created" if signal.opportunity_id in created_ids else "signal_updated"
            safe_append_outcome(
                outcomes,
                event_type,
                opportunity_id=signal.opportunity_id,
                artifact_type=signal.kind,
                feedback_score=signal.priority_score,
                metadata={
                    "kind": signal.kind,
                    "target_key": signal.target_key,
                    "title": signal.title,
                    "status": signal.status,
                    "seen_count": signal.seen_count,
                    "risk_level": signal.risk_level,
                    "evidence_sources_count": len(signal.evidence_sources),
                },
                timestamp=_parse_datetime(signal.last_seen_at),
            )


def detect_workflow_opportunity_candidates(
    history_entries: list[dict[str, Any]],
    *,
    min_evidence_sources: int = 2,
) -> list[OpportunitySignalCandidate]:
    """Detect repeated workflow-like patterns without generating any proposal."""

    grouped: dict[str, list[dict[str, Any]]] = {}
    titles: dict[str, str] = {}
    summaries: dict[str, str] = {}
    for entry in history_entries:
        content = str(entry.get("content") or "").strip()
        if not _looks_like_workflow_signal(content):
            continue
        target = _target_key(content)
        if not target:
            continue
        evidence = {
            "cursor": entry.get("cursor"),
            "session_key": entry.get("session_key"),
            "timestamp": entry.get("timestamp"),
            "preview": truncate_text(content, _SIGNAL_PREVIEW_MAX_CHARS),
        }
        grouped.setdefault(target, []).append(evidence)
        titles.setdefault(target, _title_for_target(target))
        summaries.setdefault(target, f"Repeated workflow-like request pattern: {target}")

    candidates: list[OpportunitySignalCandidate] = []
    for target, evidence in grouped.items():
        unique = _dedupe_evidence(evidence)
        if len(unique) < max(1, min_evidence_sources):
            continue
        candidates.append(OpportunitySignalCandidate(
            kind=SIGNAL_KIND_WORKFLOW,
            target_key=target,
            title=titles[target],
            summary=summaries[target],
            evidence_sources=unique,
            risk_level="low",
        ))
    return candidates


def detect_skill_opportunity_candidates(
    history_entries: list[dict[str, Any]],
    *,
    min_evidence_sources: int = 2,
) -> list[OpportunitySignalCandidate]:
    """Detect repeated read-only skill-like patterns without creating proposals."""

    grouped: dict[str, list[dict[str, Any]]] = {}
    titles: dict[str, str] = {}
    summaries: dict[str, str] = {}
    for entry in history_entries:
        content = str(entry.get("content") or "").strip()
        if not _looks_like_skill_signal(content):
            continue
        target = _target_key(content)
        if not target:
            continue
        evidence = {
            "cursor": entry.get("cursor"),
            "session_key": entry.get("session_key"),
            "timestamp": entry.get("timestamp"),
            "preview": truncate_text(content, _SIGNAL_PREVIEW_MAX_CHARS),
        }
        grouped.setdefault(target, []).append(evidence)
        titles.setdefault(target, _skill_title_for_target(target))
        summaries.setdefault(target, f"Repeated read-only skill-like request pattern: {target}")

    candidates: list[OpportunitySignalCandidate] = []
    for target, evidence in grouped.items():
        unique = _dedupe_evidence(evidence)
        if len(unique) < max(1, min_evidence_sources):
            continue
        candidates.append(OpportunitySignalCandidate(
            kind=SIGNAL_KIND_SKILL,
            target_key=target,
            title=titles[target],
            summary=summaries[target],
            evidence_sources=unique,
            risk_level="medium",
        ))
    return candidates


def _looks_like_workflow_signal(content: str) -> bool:
    if len(content.strip()) < 12:
        return False
    if _NOISE_RE.search(content):
        return False
    return bool(_WORKFLOW_CUE_RE.search(content))


def _looks_like_skill_signal(content: str) -> bool:
    if len(content.strip()) < 12:
        return False
    if _NOISE_RE.search(content):
        return False
    return bool(_SKILL_CUE_RE.search(content))


def _target_key(content: str) -> str:
    text = _SPACE_RE.sub(" ", content).strip().lower()
    text = re.sub(r"^\[raw\]\s+\d+\s+messages\s+", "", text)
    text = re.sub(r"\[[^\]]+\]\s+(user|assistant|system)(?:\s+\[tools:[^\]]+\])?:", " ", text)
    words = re.findall(r"[\w\u4e00-\u9fff-]+", text)
    if not words:
        return ""
    return truncate_text(" ".join(words[:16]), _TARGET_MAX_CHARS).replace("\n... (truncated)", "")


def _title_for_target(target: str) -> str:
    cleaned = target.strip()
    if not cleaned:
        return "Workflow candidate"
    return f"Workflow candidate: {truncate_text(cleaned, 48).replace(chr(10), ' ')}"


def _skill_title_for_target(target: str) -> str:
    cleaned = target.strip()
    if not cleaned:
        return "Skill candidate"
    return f"Skill candidate: {truncate_text(cleaned, 48).replace(chr(10), ' ')}"


def _merge_evidence_sources(
    existing: list[dict[str, Any]],
    incoming: list[dict[str, Any]],
) -> list[dict[str, Any]]:
    merged = _dedupe_evidence([*existing, *incoming])
    return merged[-_MAX_EVIDENCE_PER_SIGNAL:]


def _dedupe_evidence(items: list[dict[str, Any]]) -> list[dict[str, Any]]:
    by_key: dict[str, dict[str, Any]] = {}
    for item in items:
        cursor = item.get("cursor")
        if cursor is not None:
            key = f"cursor:{cursor}"
        else:
            payload = json.dumps(item, ensure_ascii=False, sort_keys=True)
            key = hashlib.sha256(payload.encode("utf-8")).hexdigest()
        by_key[key] = {
            "cursor": item.get("cursor"),
            "session_key": item.get("session_key"),
            "timestamp": item.get("timestamp"),
            "preview": truncate_text(
                _redact_text(str(item.get("preview") or "")),
                _SIGNAL_PREVIEW_MAX_CHARS,
            ),
        }
    return list(by_key.values())


def evolution_allows_workflow_proposals(config: Any | None = None) -> bool:
    mode = _config_str(config, "mode", DEFAULT_EVOLUTION_MODE)
    return mode in _EVOLUTION_ALLOWED_PROPOSAL_MODES and not _config_bool(
        config,
        "dry_run",
        DEFAULT_EVOLUTION_DRY_RUN,
    )


def evolution_allows_skill_proposals(config: Any | None = None) -> bool:
    return (
        evolution_allows_workflow_proposals(config)
        and _config_bool(config, "skill_candidates_enabled", False)
    )


def build_workflow_payload_from_signal(
    signal: OpportunitySignal,
    *,
    config: Any | None = None,
) -> dict[str, Any]:
    workflow_name = _workflow_name_from_target(signal.target_key)
    steps = _workflow_steps_from_signal(signal, config=config)
    payload = {
        "subject_type": "workflow",
        "subject_id": workflow_name,
        "subject_path": f"workflows/{workflow_name}/workflow.yaml",
        "curator_key": f"auto-evolution-workflow:{signal.opportunity_id}",
        "target_state_hash": _stable_hash([
            signal.opportunity_id,
            signal.target_key,
            signal.seen_count,
            round(signal.priority_score, 3),
        ]),
        "suggested_action": "workflow",
        "impact_summary": signal.summary,
        "workflow_name": workflow_name,
        "description": f"Manual workflow candidate discovered from repeated usage: {signal.target_key}",
        "body": _workflow_body_from_signal(signal),
        "steps": steps,
        "evolution": {
            "origin": AUTO_EVOLUTION_ORIGIN,
            "opportunity_id": signal.opportunity_id,
            "kind": signal.kind,
            "priority_score": round(signal.priority_score, 3),
            "seen_count": signal.seen_count,
            "risk_level": signal.risk_level,
            "evidence_sources": _clean_evidence_sources(signal.evidence_sources),
        },
    }
    gate = static_gate_workflow_payload(payload, config=config)
    payload["static_gate"] = {
        "decision": gate["decision"],
        "issues": gate["issues"],
        "issue_counts": gate["issue_counts"],
    }
    return payload


def build_skill_payload_from_signal(
    signal: OpportunitySignal,
    *,
    config: Any | None = None,
) -> dict[str, Any]:
    skill_name = _skill_name_from_target(signal.target_key)
    body, body_truncated = _skill_body_from_signal(signal)
    payload = {
        "subject_type": "skill",
        "subject_id": skill_name,
        "subject_path": f"skills/{skill_name}/SKILL.md",
        "curator_key": f"auto-evolution-skill:{signal.opportunity_id}",
        "target_state_hash": _stable_hash([
            signal.opportunity_id,
            signal.target_key,
            signal.seen_count,
            round(signal.priority_score, 3),
        ]),
        "suggested_action": "skill",
        "impact_summary": signal.summary,
        "skill_name": skill_name,
        "description": f"Read-only skill candidate discovered from repeated usage: {signal.target_key}",
        "body": body,
        "evolution": {
            "origin": AUTO_EVOLUTION_ORIGIN,
            "opportunity_id": signal.opportunity_id,
            "kind": signal.kind,
            "priority_score": round(signal.priority_score, 3),
            "seen_count": signal.seen_count,
            "risk_level": signal.risk_level,
            "body_truncated": body_truncated,
            "evidence_sources": _clean_evidence_sources(signal.evidence_sources),
        },
    }
    gate = static_gate_skill_payload(payload, config=config)
    payload["static_gate"] = {
        "decision": gate["decision"],
        "issues": gate["issues"],
        "issue_counts": gate["issue_counts"],
    }
    return payload


def static_gate_workflow_payload(payload: dict[str, Any], *, config: Any | None = None) -> dict[str, Any]:
    issues: list[ValidationIssue] = []
    allowed_tools = {
        str(item).strip()
        for item in getattr(config, "static_gate_allowed_workflow_tools", []) or []
        if str(item).strip()
    } or {"read_file", "glob", "grep", "web_fetch"}
    max_steps = max(1, _config_int(config, "static_gate_max_workflow_steps", 10))
    steps = payload.get("steps")
    if not isinstance(steps, list):
        issues.append(_issue("workflow_steps_not_list", "reject", "Workflow steps must be a list."))
        steps = []
    if len(steps) > max_steps:
        issues.append(_issue(
            "workflow_too_many_steps",
            "pending",
            f"Workflow has {len(steps)} steps; configured maximum is {max_steps}.",
        ))
    for index, step in enumerate(steps):
        if not isinstance(step, dict):
            issues.append(_issue(
                "workflow_step_not_mapping",
                "reject",
                f"Workflow step {index + 1} must be a mapping.",
            ))
            continue
        text = " ".join(
            str(step.get(key) or "")
            for key in ("title", "instruction", "risk")
        )
        tool = str(step.get("tool") or "").strip()
        if tool and tool not in allowed_tools:
            issues.append(_issue(
                "workflow_tool_not_allowed",
                "pending",
                f"Workflow step {index + 1} references non-whitelisted tool `{tool}`.",
            ))
        lowered = text.casefold()
        if _contains_dangerous_workflow_term(lowered):
            issues.append(_issue(
                "workflow_dangerous_instruction",
                "pending",
                f"Workflow step {index + 1} mentions a sensitive or side-effecting action.",
            ))
    body = str(payload.get("body") or "")
    if _contains_dangerous_workflow_term(body):
        issues.append(_issue(
            "workflow_body_sensitive_action",
            "pending",
            "Workflow body mentions a sensitive or side-effecting action.",
        ))
    if any(issue.severity == "reject" for issue in issues):
        decision = "reject"
    elif issues:
        decision = "requires_manual_review"
    else:
        decision = "pass"
    return {
        "decision": decision,
        "issues": [asdict(issue) for issue in issues],
        "issue_counts": _issue_counts(issues),
    }


def static_gate_skill_payload(payload: dict[str, Any], *, config: Any | None = None) -> dict[str, Any]:
    issues: list[ValidationIssue] = []
    allowed_tools = {
        str(item).strip()
        for item in getattr(config, "static_gate_allowed_skill_tools", []) or []
        if str(item).strip()
    } or {"read_file", "glob", "grep"}
    body = str(payload.get("body") or "")
    lower_body = body.casefold()
    for tool in sorted(allowed_tools):
        lower_body = lower_body.replace(tool.casefold(), "")
    if _DANGEROUS_SKILL_TOOL_RE.search(lower_body):
        issues.append(_issue(
            "skill_tool_not_allowed",
            "pending",
            "Skill draft references a non-read-only or sensitive tool.",
        ))
    if _DANGEROUS_SKILL_INSTALL_RE.search(body):
        issues.append(_issue(
            "skill_install_command",
            "pending",
            "Skill draft includes an installation command.",
        ))
    if _DANGEROUS_SKILL_COMMAND_HEADING_RE.search(body):
        issues.append(_issue(
            "skill_command_heading",
            "pending",
            "Skill draft includes a command-oriented heading.",
        ))
    if _contains_dangerous_workflow_term(body):
        issues.append(_issue(
            "skill_sensitive_action",
            "pending",
            "Skill draft mentions a sensitive or side-effecting action.",
        ))
    evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
    if bool(evolution.get("body_truncated")):
        issues.append(_issue(
            "skill_body_truncated",
            "warning",
            "Skill body was truncated to the configured artifact size limit.",
        ))
    if any(issue.severity == "reject" for issue in issues):
        decision = "reject"
    elif issues:
        decision = "requires_manual_review"
    else:
        decision = "pass"
    return {
        "decision": decision,
        "issues": [asdict(issue) for issue in issues],
        "issue_counts": _issue_counts(issues),
    }


def _workflow_name_from_target(target: str) -> str:
    slug = re.sub(r"[^a-z0-9]+", "-", _redact_text(target).casefold()).strip("-")
    slug = re.sub(r"-{2,}", "-", slug)
    if not slug:
        digest = hashlib.sha256(target.encode("utf-8")).hexdigest()[:8]
        slug = f"workflow-{digest}"
    if not re.match(r"^[a-z0-9]", slug):
        slug = f"workflow-{slug}"
    return slug[:64].strip("-") or "workflow-candidate"


def _skill_name_from_target(target: str) -> str:
    slug = re.sub(r"[^a-z0-9]+", "-", _redact_text(target).casefold()).strip("-")
    slug = re.sub(r"-{2,}", "-", slug)
    if not slug:
        digest = hashlib.sha256(target.encode("utf-8")).hexdigest()[:8]
        slug = f"skill-{digest}"
    if not re.match(r"^[a-z0-9]", slug):
        slug = f"skill-{slug}"
    return slug[:64].strip("-") or "skill-candidate"


def _workflow_body_from_signal(signal: OpportunitySignal) -> str:
    lines = [
        signal.summary,
        "",
        "## Evidence",
    ]
    for item in signal.evidence_sources[:5]:
        preview = _redact_text(str(item.get("preview") or "")).strip()
        cursor = item.get("cursor")
        timestamp = item.get("timestamp")
        label = f"cursor={cursor}" if cursor is not None else "cursor=unknown"
        if timestamp:
            label = f"{label}, timestamp={timestamp}"
        if preview:
            lines.append(f"- {label}: {preview}")
        else:
            lines.append(f"- {label}")
    return "\n".join(lines).strip()


def _skill_body_from_signal(signal: OpportunitySignal) -> tuple[str, bool]:
    title = signal.title.replace("Skill candidate:", "").strip() or signal.target_key
    lines = [
        f"# {_human_title(title)}",
        "",
        "Use this skill when a repeated, read-only analysis or troubleshooting pattern matches the evidence below.",
        "",
        "## Guardrails",
        "- Use only read-only tools unless a human explicitly approves a separate proposal.",
        "- Do not install packages, change files, contact people, or schedule background work.",
        "- Keep findings grounded in the current workspace and the cited evidence.",
        "",
        "## Procedure",
        f"1. Confirm the user is asking about this pattern: {_redact_text(signal.target_key)}.",
        "2. Gather read-only context with allowed tools such as `read_file`, `glob`, or `grep`.",
        "3. Summarize the result and call out uncertainty instead of making state-changing changes.",
        "",
        "## Evidence",
    ]
    for item in signal.evidence_sources[:5]:
        preview = _redact_text(str(item.get("preview") or "")).strip()
        cursor = item.get("cursor")
        timestamp = item.get("timestamp")
        label = f"cursor={cursor}" if cursor is not None else "cursor=unknown"
        if timestamp:
            label = f"{label}, timestamp={timestamp}"
        if preview:
            lines.append(f"- {label}: {preview}")
        else:
            lines.append(f"- {label}")
    body = "\n".join(lines).strip()
    if len(body) <= _SKILL_BODY_MAX_CHARS:
        return body, False
    return body[:_SKILL_BODY_MAX_CHARS].rstrip(), True


def _workflow_steps_from_signal(signal: OpportunitySignal, *, config: Any | None = None) -> list[dict[str, Any]]:
    max_steps = max(1, _config_int(config, "static_gate_max_workflow_steps", 10))
    steps = [
        {
            "title": "Confirm workflow intent",
            "instruction": (
                "Review the repeated request pattern and confirm the manual workflow still matches "
                "the user's current intent."
            ),
            "risk": "low",
            "confirmation_required": False,
        },
        {
            "title": "Gather read-only context",
            "instruction": (
                "Use only read-only context gathering, such as reading relevant files or searching "
                "workspace text, before following the workflow."
            ),
            "risk": "low",
            "confirmation_required": False,
        },
        {
            "title": "Run the documented manual checks",
            "instruction": (
                f"Follow the repeated pattern: {truncate_text(_redact_text(signal.target_key), 220)}."
            ),
            "risk": "low",
            "confirmation_required": True,
        },
    ]
    return steps[:max_steps]


def _clean_evidence_sources(items: list[dict[str, Any]]) -> list[dict[str, Any]]:
    cleaned: list[dict[str, Any]] = []
    for item in items:
        if not isinstance(item, dict):
            continue
        cleaned.append({
            "cursor": item.get("cursor"),
            "session_key": item.get("session_key"),
            "timestamp": item.get("timestamp"),
            "preview": truncate_text(
                _redact_text(str(item.get("preview") or "")),
                _SIGNAL_PREVIEW_MAX_CHARS,
            ),
        })
    return cleaned


def _clean_signal_text(text: Any, max_chars: int) -> str:
    return truncate_text(_redact_text(str(text or "").strip()), max_chars)


def _contains_dangerous_workflow_term(text: str) -> bool:
    if _DANGEROUS_WORKFLOW_WORD_RE.search(text):
        return True
    lowered = text.casefold()
    return any(term.casefold() in lowered for term in _DANGEROUS_WORKFLOW_CJK_TERMS)


def _contains_dangerous_skill_term(text: str) -> bool:
    return bool(
        _DANGEROUS_SKILL_TOOL_RE.search(text)
        or _DANGEROUS_SKILL_INSTALL_RE.search(text)
        or _DANGEROUS_SKILL_COMMAND_HEADING_RE.search(text)
    )


def _human_title(text: str) -> str:
    words = re.findall(r"[\w\u4e00-\u9fff-]+", text.strip())
    if not words:
        return "Read Only Skill Candidate"
    return " ".join(words[:8]).replace("-", " ").title()


def _redact_text(text: str) -> str:
    try:
        from OriginAgent.agent.memory import redact_memory_text

        return redact_memory_text(text)
    except Exception:
        return text


def _issue(code: str, severity: str, message: str) -> ValidationIssue:
    return ValidationIssue(code=code, severity=severity, message=message)


def _issue_counts(issues: list[ValidationIssue]) -> dict[str, int]:
    counts: dict[str, int] = {}
    for issue in issues:
        counts[issue.severity] = counts.get(issue.severity, 0) + 1
    return counts


def _stable_hash(values: list[Any]) -> str:
    payload = json.dumps(values, sort_keys=True, ensure_ascii=False, separators=(",", ":"))
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def _priority_score(signal: OpportunitySignal) -> float:
    seen_component = min(signal.seen_count / _MAX_EXPECTED_SEEN_COUNT, 1.0)
    evidence_diversity = min(len({
        str(item.get("cursor") or item.get("timestamp") or item.get("preview") or "")
        for item in signal.evidence_sources
    }) / max(signal.seen_count, 1), 1.0)
    score = (
        0.4 * seen_component
        + 0.3 * (_SKILL_CONFIDENCE if signal.kind == SIGNAL_KIND_SKILL else _WORKFLOW_CONFIDENCE)
        + 0.3 * evidence_diversity
    )
    score = (score * signal.feedback_multiplier) + signal.feedback_score_offset
    return max(0.0, min(1.0, score))


def _normalize_datetime(value: datetime | None) -> datetime:
    if value is None:
        return datetime.now(timezone.utc)
    if value.tzinfo is None:
        return value.replace(tzinfo=timezone.utc)
    return value.astimezone(timezone.utc)


def _parse_datetime(value: str) -> datetime | None:
    if not value:
        return None
    with suppress(ValueError):
        parsed = datetime.fromisoformat(value)
        return _normalize_datetime(parsed)
    return None


def _safe_int(value: Any, default: int) -> int:
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default


def _safe_float(value: Any, default: float) -> float:
    try:
        return float(value)
    except (TypeError, ValueError):
        return default


def _config_str(config: Any | None, attr: str, default: str) -> str:
    value = getattr(config, attr, default) if config is not None else default
    return str(value or default)


def _config_bool(config: Any | None, attr: str, default: bool) -> bool:
    value = getattr(config, attr, default) if config is not None else default
    return bool(value)


def _config_int(config: Any | None, attr: str, default: int) -> int:
    return _safe_int(getattr(config, attr, default) if config is not None else default, default)


def _config_float(config: Any | None, attr: str, default: float) -> float:
    return _safe_float(getattr(config, attr, default) if config is not None else default, default)
