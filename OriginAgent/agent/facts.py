"""Structured long-term memory facts.

`facts.jsonl` is a current-state store, not an append-only event log. Bad JSON
lines are ignored on reads and dropped the next time the store rewrites the
file. Fact `content` is the remembered human-readable fact and is not redacted;
`source_excerpt` is supporting evidence and is redacted before persistence.
"""

from __future__ import annotations

from collections import Counter
import hashlib
import json
import os
import re
import uuid
from contextlib import suppress
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from filelock import FileLock
from loguru import logger

from OriginAgent.utils.helpers import ensure_dir


@dataclass(frozen=True)
class FactStoreConfig:
    category_order: tuple[str, ...]
    legacy_categories: tuple[str, ...] = ()
    valid_owners: tuple[str, ...] = ("user", "assistant", "system", "unknown")
    high_risk_categories: tuple[str, ...] = ("policy", "safety")
    conflict_categories: tuple[str, ...] = ("preference", "routine", "policy", "safety", "temporary")
    high_risk_keywords: tuple[str, ...] = ()
    temporary_language: tuple[str, ...] = ()
    uncertain_language: tuple[str, ...] = ()
    confidence_decay_factor: float = 0.98
    min_confidence: float = 0.3
    decay_start_days: int = 30
    auto_active_confidence_threshold: float = 0.8
    auto_active_budget: float = 5.0
    reject_confidence_multiplier: float = 0.7
    calibration_min_count: int = 10


DEFAULT_FACT_STORE_CONFIG = FactStoreConfig(
    category_order=(
        "preference",
        "routine",
        "policy",
        "safety",
        "temporary",
        "note",
    ),
    legacy_categories=("household", "device"),
    valid_owners=("user", "assistant", "system", "unknown", "household"),
    high_risk_keywords=(
        "security",
        "camera",
        "gas",
        "medication",
        "payment",
        "password",
        "key",
        "permission",
        "token",
    ),
    temporary_language=(
        "today",
        "tomorrow",
        "this week",
        "temporary",
        "for now",
        "just this time",
        "今天",
        "明天",
        "这周",
        "本周",
        "临时",
        "暂时",
        "先",
        "这次",
    ),
    uncertain_language=(
        "maybe",
        "usually",
        "sometimes",
        "probably",
        "around",
        "roughly",
        "可能",
        "一般",
        "有时",
        "大概",
        "差不多",
        "左右",
        "偶尔",
    ),
)

CATEGORY_ORDER = DEFAULT_FACT_STORE_CONFIG.category_order + DEFAULT_FACT_STORE_CONFIG.legacy_categories
VALID_CATEGORIES = set(CATEGORY_ORDER)
VALID_OWNERS = set(DEFAULT_FACT_STORE_CONFIG.valid_owners)
VALID_STATUSES = {
    "active",
    "deprecated",
    "contradicted",
    "pending_confirmation",
}
HIGH_RISK_CATEGORIES = set(DEFAULT_FACT_STORE_CONFIG.high_risk_categories)
CONFLICT_CATEGORIES = set(DEFAULT_FACT_STORE_CONFIG.conflict_categories)
MAX_DEPRECATIONS_PER_BATCH = 3
HIGH_RISK_KEYWORDS = DEFAULT_FACT_STORE_CONFIG.high_risk_keywords
TEMPORARY_LANGUAGE = DEFAULT_FACT_STORE_CONFIG.temporary_language
UNCERTAIN_LANGUAGE = DEFAULT_FACT_STORE_CONFIG.uncertain_language
HIGH_RISK_DEVICE_DOMAINS = {"lock", "security", "camera", "gas", "presence"}


@dataclass
class ParseRejectedProposal:
    section: str
    index: int | None
    raw: Any
    error: str


@dataclass
class FactProposal:
    content: str
    category: str
    scope: str
    owner: str
    source_cursors: list[int]
    source_excerpt: str
    confidence: float = 0.7
    expires_at: str | None = None
    supersedes_fact_id: str | None = None
    requires_confirmation: bool | None = None
    status: str | None = None
    reason: str = ""


@dataclass
class FactDeprecationProposal:
    fact_id: str
    reason: str
    source_cursors: list[int] = field(default_factory=list)


@dataclass
class DreamFactProposalBatch:
    facts_to_upsert: list[FactProposal] = field(default_factory=list)
    facts_to_deprecate: list[FactDeprecationProposal] = field(default_factory=list)
    memory_render_hints: list[Any] = field(default_factory=list)
    parse_rejected: list[ParseRejectedProposal] = field(default_factory=list)


@dataclass
class ValidationIssue:
    code: str
    severity: str
    message: str


@dataclass
class ValidationResult:
    decision: str
    confidence: float
    issues: list[ValidationIssue] = field(default_factory=list)


@dataclass
class DreamFactApplyResult:
    accepted: list[FactRecord] = field(default_factory=list)
    pending: list[FactRecord] = field(default_factory=list)
    rejected: list[tuple[FactProposal | FactDeprecationProposal, list[ValidationIssue]]] = (
        field(default_factory=list)
    )
    deprecated: list[str] = field(default_factory=list)
    parse_rejected: list[ParseRejectedProposal] = field(default_factory=list)


@dataclass
class FactRecord:
    fact_id: str
    content: str
    canonical_key: str
    category: str
    scope: str
    owner: str
    source_cursors: list[int]
    source_excerpt: str
    confidence: float
    status: str
    created_at: str
    updated_at: str
    last_seen_at: str
    expires_at: str | None
    supersedes_fact_id: str | None
    requires_confirmation: bool

    @classmethod
    def from_dict(cls, raw: dict[str, Any]) -> "FactRecord":
        required_strings = (
            "fact_id",
            "content",
            "canonical_key",
            "category",
            "scope",
            "owner",
            "source_excerpt",
            "status",
            "created_at",
            "updated_at",
            "last_seen_at",
        )
        if any(not isinstance(raw.get(key), str) for key in required_strings):
            raise ValueError("fact record missing required string fields")
        source_cursors = _normalize_source_cursors(raw.get("source_cursors"))
        confidence = _normalize_confidence(raw.get("confidence", 1.0))
        expires_at = raw.get("expires_at")
        supersedes_fact_id = raw.get("supersedes_fact_id")
        requires_confirmation = raw.get("requires_confirmation", False)
        if expires_at is not None and not isinstance(expires_at, str):
            raise ValueError("expires_at must be a string or null")
        if supersedes_fact_id is not None and not isinstance(supersedes_fact_id, str):
            raise ValueError("supersedes_fact_id must be a string or null")
        if not isinstance(requires_confirmation, bool):
            raise ValueError("requires_confirmation must be bool")
        return cls(
            fact_id=raw["fact_id"],
            content=raw["content"],
            canonical_key=raw["canonical_key"],
            category=_normalize_category(raw["category"]),
            scope=_normalize_scope(raw["scope"]),
            owner=_normalize_owner(raw["owner"]),
            source_cursors=source_cursors,
            source_excerpt=raw["source_excerpt"],
            confidence=confidence,
            status=_normalize_status(raw["status"]),
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            last_seen_at=raw["last_seen_at"],
            expires_at=expires_at,
            supersedes_fact_id=supersedes_fact_id,
            requires_confirmation=requires_confirmation,
        )

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


def normalize_fact_content(content: str) -> str:
    text = content.strip().lower()
    return re.sub(r"\s+", " ", text)


def canonical_key_for_fact(
    content: str,
    owner: str,
    category: str,
    scope: str,
) -> str:
    payload = "\0".join([
        normalize_fact_content(content),
        owner.strip().lower(),
        category.strip().lower(),
        scope.strip().lower(),
    ])
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def domain_id_for_fact(scope: str, content: str) -> str:
    return _infer_domain_key(scope, content)


def parse_fact_proposal_response(text: str) -> DreamFactProposalBatch:
    raw_text = _extract_json_text(text)
    try:
        parsed = json.loads(raw_text)
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid fact proposal JSON: {exc}") from exc
    if not isinstance(parsed, dict):
        raise ValueError("fact proposal response must be a JSON object")

    batch = DreamFactProposalBatch(
        memory_render_hints=(
            parsed.get("memory_render_hints")
            if isinstance(parsed.get("memory_render_hints"), list)
            else []
        ),
    )
    upserts = parsed.get("facts_to_upsert", [])
    deprecations = parsed.get("facts_to_deprecate", [])
    if not isinstance(upserts, list):
        batch.parse_rejected.append(ParseRejectedProposal(
            section="facts_to_upsert",
            index=None,
            raw=upserts,
            error="facts_to_upsert must be a list",
        ))
        upserts = []
    if not isinstance(deprecations, list):
        batch.parse_rejected.append(ParseRejectedProposal(
            section="facts_to_deprecate",
            index=None,
            raw=deprecations,
            error="facts_to_deprecate must be a list",
        ))
        deprecations = []

    for index, raw in enumerate(upserts):
        try:
            batch.facts_to_upsert.append(_parse_fact_proposal(raw))
        except ValueError as exc:
            batch.parse_rejected.append(ParseRejectedProposal(
                section="facts_to_upsert",
                index=index,
                raw=raw,
                error=str(exc),
            ))

    for index, raw in enumerate(deprecations):
        try:
            batch.facts_to_deprecate.append(_parse_deprecation_proposal(raw))
        except ValueError as exc:
            batch.parse_rejected.append(ParseRejectedProposal(
                section="facts_to_deprecate",
                index=index,
                raw=raw,
                error=str(exc),
            ))

    return batch


def validate_fact_proposal(
    proposal: FactProposal,
    *,
    existing_facts: list[FactRecord],
    history_entries: list[dict[str, Any]],
    batch_cursor_min: int,
    batch_cursor_max: int,
) -> ValidationResult:
    issues: list[ValidationIssue] = []
    try:
        category = _normalize_category(proposal.category)
        scope = _normalize_scope(proposal.scope)
        _normalize_owner(proposal.owner)
        if proposal.status is not None:
            _normalize_status(proposal.status)
        if (
            proposal.requires_confirmation is not None
            and not isinstance(proposal.requires_confirmation, bool)
        ):
            raise ValueError("requires_confirmation must be bool")
    except ValueError as exc:
        issues.append(_issue("invalid_field", "reject", str(exc)))
        return ValidationResult(
            decision="reject",
            confidence=_normalize_confidence(proposal.confidence),
            issues=issues,
        )

    confidence = _normalize_confidence(proposal.confidence)
    history_cursors = _history_cursor_set(history_entries)
    if not proposal.source_cursors:
        issues.append(_issue(
            "missing_source_cursors",
            "reject",
            "Fact proposals must cite at least one source cursor.",
        ))
    elif any(
        cursor < batch_cursor_min
        or cursor > batch_cursor_max
        or cursor not in history_cursors
        for cursor in proposal.source_cursors
    ):
        issues.append(_issue(
            "source_cursor_out_of_batch",
            "reject",
            "All source cursors must be in the current Dream batch.",
        ))

    if not proposal.source_excerpt.strip():
        issues.append(_issue(
            "missing_source_excerpt",
            "reject",
            "Fact proposals must include a supporting source excerpt.",
        ))
    elif not _source_excerpt_matches_history(proposal, history_entries):
        issues.append(_issue(
            "source_excerpt_not_found",
            "reject",
            "source_excerpt must be present in the cited history cursor content.",
        ))

    combined_text = " ".join([
        proposal.content,
        proposal.source_excerpt,
        proposal.scope,
        proposal.reason,
    ])
    if (
        category in HIGH_RISK_CATEGORIES
        or _contains_any(combined_text, HIGH_RISK_KEYWORDS)
        or _requires_high_risk_confirmation(scope=scope, content=proposal.content)
    ):
        issues.append(_issue(
            "high_risk_memory",
            "pending",
            "High-risk security facts require confirmation.",
        ))

    if category == "temporary" and not proposal.expires_at:
        issues.append(_issue(
            "temporary_missing_expires_at",
            "pending",
            "Temporary facts without an expiration require confirmation.",
        ))

    if category != "temporary" and _contains_any(combined_text, TEMPORARY_LANGUAGE):
        issues.append(_issue(
            "temporary_language_non_temporary",
            "pending",
            "Temporary language in a non-temporary fact requires confirmation.",
        ))

    if _contains_any(combined_text, UNCERTAIN_LANGUAGE):
        confidence = min(confidence, 0.7)
        issues.append(_issue(
            "uncertain_language",
            "warn",
            "Uncertain language capped confidence at 0.7.",
        ))

    if category in CONFLICT_CATEGORIES:
        for record in existing_facts:
            if (
                record.status == "active"
                and record.category == category
                and record.scope == scope
                and normalize_fact_content(record.content)
                != normalize_fact_content(proposal.content)
            ):
                issues.append(_issue(
                    "possible_conflict",
                    "pending",
                    "Existing active fact in the same category/scope differs.",
                ))
                break

    if any(issue.severity == "reject" for issue in issues):
        decision = "reject"
    elif any(issue.severity == "pending" for issue in issues):
        decision = "pending_confirmation"
    else:
        decision = "active"
    return ValidationResult(decision=decision, confidence=confidence, issues=issues)


def validate_deprecation_proposal(
    proposal: FactDeprecationProposal,
    *,
    existing_facts: list[FactRecord],
    history_entries: list[dict[str, Any]],
    batch_cursor_min: int,
    batch_cursor_max: int,
) -> ValidationResult:
    issues: list[ValidationIssue] = []
    target = next((fact for fact in existing_facts if fact.fact_id == proposal.fact_id), None)
    if target is None:
        issues.append(_issue(
            "unknown_fact_id",
            "reject",
            "Deprecation fact_id does not exist.",
        ))
    if not proposal.reason.strip():
        issues.append(_issue(
            "missing_deprecation_reason",
            "reject",
            "Deprecation proposals must include a non-empty reason.",
        ))

    history_cursors = _history_cursor_set(history_entries)
    if not proposal.source_cursors or any(
        cursor < batch_cursor_min
        or cursor > batch_cursor_max
        or cursor not in history_cursors
        for cursor in proposal.source_cursors
    ):
        issues.append(_issue(
            "missing_current_source_for_deprecation",
            "reject",
            "Deprecations must include source_cursors in the current Dream batch.",
        ))

    if (
        target is not None
        and target.category in HIGH_RISK_CATEGORIES
        and target.status == "active"
    ):
        issues.append(_issue(
            "active_high_risk_deprecation",
            "reject",
            "Active policy/safety facts cannot be deprecated automatically.",
        ))

    decision = "reject" if any(issue.severity == "reject" for issue in issues) else "active"
    return ValidationResult(decision=decision, confidence=1.0, issues=issues)


def render_memory_md(records: list[FactRecord]) -> str:
    active_by_category: dict[str, list[FactRecord]] = {
        category: [] for category in CATEGORY_ORDER
    }
    pending: list[FactRecord] = []
    for record in records:
        if record.status == "pending_confirmation":
            pending.append(record)
            continue
        if record.status != "active":
            continue
        active_by_category.setdefault(record.category, []).append(record)

    lines = [
        "# Long-term Memory",
        "",
        "> Generated from memory/facts.jsonl. Do not edit directly.",
    ]

    for category in CATEGORY_ORDER:
        facts = sorted(
            active_by_category.get(category, []),
            key=lambda fact: (
                fact.scope.casefold(),
                fact.content.casefold(),
                fact.fact_id,
            ),
        )
        if not facts:
            continue
        lines.extend(["", f"## {_section_title(category)}"])
        for fact in facts:
            lines.extend(_render_fact_lines(fact))

    pending = sorted(
        pending,
        key=lambda fact: (
            CATEGORY_ORDER.index(fact.category)
            if fact.category in VALID_CATEGORIES
            else len(CATEGORY_ORDER),
            fact.scope.casefold(),
            fact.content.casefold(),
            fact.fact_id,
        ),
    )
    if pending:
        lines.extend(["", "## Pending Confirmation"])
        for fact in pending:
            lines.extend(_render_fact_lines(fact, include_category=True))

    return "\n".join(lines).rstrip() + "\n"


def summarize_facts(
    workspace: Path,
    *,
    fact_store: FactStore | None = None,
) -> dict[str, Any]:
    """Return redacted fact counts without exposing raw content."""

    store = fact_store or FactStore(Path(workspace))
    with store._locked():
        records = store.read_all_unlocked()

    category_counts: Counter[str] = Counter()
    domain_counts: Counter[str] = Counter()
    active_count = 0
    pending_confirmation_count = 0
    for record in records:
        if record.status == "active":
            active_count += 1
            category_counts[record.category] += 1
            domain_counts[_fact_domain_key(record)] += 1
        elif record.status == "pending_confirmation":
            pending_confirmation_count += 1
    return {
        "active_count": active_count,
        "pending_confirmation_count": pending_confirmation_count,
        "category_counts": dict(category_counts),
        "domain_counts": dict(domain_counts),
    }


class FactStore:
    def __init__(
        self,
        workspace: Path,
        *,
        facts_file: Path | None = None,
        lock_factory: Callable[[], FileLock] | None = None,
        redactor: Callable[[str], str] | None = None,
        config: FactStoreConfig | None = None,
    ):
        self.workspace = workspace
        self.memory_dir = ensure_dir(workspace / "memory")
        self.facts_file = facts_file or self.memory_dir / "facts.jsonl"
        self.calibration_file = self.memory_dir / "confidence_calibration.json"
        self._lock_file = self.memory_dir / ".lock"
        self._lock_factory = lock_factory
        self._redactor = redactor or _default_redactor
        self.config = config or DEFAULT_FACT_STORE_CONFIG

    def _locked(self) -> FileLock:
        if self._lock_factory is not None:
            return self._lock_factory()
        return FileLock(str(self._lock_file))

    def read_all(self) -> list[FactRecord]:
        with self._locked():
            return self.read_all_unlocked()

    def read_all_unlocked(self) -> list[FactRecord]:
        records: list[FactRecord] = []
        with suppress(FileNotFoundError):
            with open(self.facts_file, "r", encoding="utf-8") as f:
                for line_no, line in enumerate(f, start=1):
                    raw_line = line.strip()
                    if not raw_line:
                        continue
                    try:
                        parsed = json.loads(raw_line)
                        if not isinstance(parsed, dict):
                            raise ValueError("fact line is not an object")
                        records.append(FactRecord.from_dict(parsed))
                    except (json.JSONDecodeError, ValueError, TypeError):
                        logger.warning(
                            "Skipping invalid facts.jsonl line {} in {}",
                            line_no,
                            self.facts_file,
                        )
                        continue
        return records

    def list_active(
        self,
        *,
        category: str | None = None,
        scope_prefix: str | None = None,
        include_pending: bool = False,
    ) -> list[FactRecord]:
        with self._locked():
            return self.list_active_unlocked(
                category=category,
                scope_prefix=scope_prefix,
                include_pending=include_pending,
            )

    def list_active_unlocked(
        self,
        *,
        category: str | None = None,
        scope_prefix: str | None = None,
        include_pending: bool = False,
    ) -> list[FactRecord]:
        target_category = _normalize_category(category) if category else None
        target_scope = _normalize_scope(scope_prefix) if scope_prefix else None
        statuses = {"active", "pending_confirmation"} if include_pending else {"active"}
        facts = [
            record
            for record in self.read_all_unlocked()
            if record.status in statuses
        ]
        if target_category:
            facts = [record for record in facts if record.category == target_category]
        if target_scope:
            facts = [
                record
                for record in facts
                if record.scope == target_scope or record.scope.startswith(f"{target_scope}.")
            ]
        return sorted(
            facts,
            key=lambda fact: (
                CATEGORY_ORDER.index(fact.category)
                if fact.category in VALID_CATEGORIES
                else len(CATEGORY_ORDER),
                fact.scope.casefold(),
                fact.content.casefold(),
                fact.fact_id,
            ),
        )

    def find_by_canonical_key(
        self,
        canonical_key: str,
        *,
        status: str | None = None,
    ) -> FactRecord | None:
        with self._locked():
            return self.find_by_canonical_key_unlocked(canonical_key, status=status)

    def find_by_canonical_key_unlocked(
        self,
        canonical_key: str,
        *,
        status: str | None = None,
    ) -> FactRecord | None:
        canonical_key = str(canonical_key or "").strip()
        if not canonical_key:
            return None
        target_status = _normalize_status(status) if status else None
        for record in self.read_all_unlocked():
            if record.canonical_key != canonical_key:
                continue
            if target_status and record.status != target_status:
                continue
            return record
        return None

    def update_confidence(self, fact_id: str, confidence: float) -> FactRecord | None:
        with self._locked():
            records = self.read_all_unlocked()
            updated = self.update_confidence_in_records_unlocked(
                records,
                fact_id,
                confidence,
            )
            if updated is not None:
                self._write_records_unlocked(records)
            return updated

    def update_confidence_in_records_unlocked(
        self,
        records: list[FactRecord],
        fact_id: str,
        confidence: float,
    ) -> FactRecord | None:
        fact_id = str(fact_id or "").strip()
        if not fact_id:
            return None
        new_confidence = _normalize_confidence(confidence)
        for record in records:
            if record.fact_id != fact_id:
                continue
            if record.confidence == new_confidence:
                return record
            record.confidence = new_confidence
            record.updated_at = datetime.now().isoformat()
            return record
        return None

    def decay_confidence(
        self,
        *,
        factor: float | None = None,
        min_confidence: float | None = None,
        decay_start_days: int | None = None,
        now: datetime | None = None,
    ) -> int:
        with self._locked():
            records = self.read_all_unlocked()
            changed = self.decay_confidence_in_records_unlocked(
                records,
                factor=factor,
                min_confidence=min_confidence,
                decay_start_days=decay_start_days,
                now=now,
            )
            if changed:
                self._write_records_unlocked(records)
            return changed

    def decay_confidence_in_records_unlocked(
        self,
        records: list[FactRecord],
        *,
        factor: float | None = None,
        min_confidence: float | None = None,
        decay_start_days: int | None = None,
        now: datetime | None = None,
    ) -> int:
        factor = _normalize_decay_factor(
            self.config.confidence_decay_factor if factor is None else factor
        )
        min_confidence = _normalize_confidence(
            self.config.min_confidence if min_confidence is None else min_confidence
        )
        start_days = (
            self.config.decay_start_days
            if decay_start_days is None
            else int(decay_start_days)
        )
        start_days = max(0, start_days)
        now_dt = now or datetime.now()
        changed = 0
        for record in records:
            if record.status != "active":
                continue
            last_seen = _parse_datetime(record.last_seen_at)
            if last_seen is None:
                continue
            days_since_seen = _days_between(last_seen, now_dt)
            if days_since_seen <= start_days:
                continue
            new_confidence = max(
                record.confidence * (factor ** days_since_seen),
                min_confidence,
            )
            new_confidence = _normalize_confidence(new_confidence)
            if new_confidence == record.confidence:
                continue
            record.confidence = new_confidence
            record.updated_at = now_dt.isoformat()
            changed += 1
        return changed

    def calibrate_confidence(
        self,
        proposal_type: str,
        domain_id: str,
        raw_confidence: float,
    ) -> float:
        confidence = _normalize_confidence(raw_confidence)
        key = _calibration_key(proposal_type, domain_id)
        if not key:
            return confidence
        calibration = self._read_confidence_calibration()
        raw_entry = calibration.get(key)
        if not isinstance(raw_entry, dict):
            return confidence
        count = _int_value(raw_entry.get("count"), default=0)
        if count <= self.config.calibration_min_count:
            return confidence
        bias = _float_value(raw_entry.get("bias"), default=0.0)
        return max(0.1, min(0.99, confidence + bias))

    def _read_confidence_calibration(self) -> dict[str, Any]:
        try:
            parsed = json.loads(self.calibration_file.read_text(encoding="utf-8"))
        except (FileNotFoundError, json.JSONDecodeError, OSError):
            return {}
        return parsed if isinstance(parsed, dict) else {}

    def upsert_fact(
        self,
        content: str,
        *,
        category: str = "note",
        scope: str = "general",
        owner: str = "unknown",
        source_cursors: list[int] | None = None,
        source_excerpt: str = "",
        confidence: float = 1.0,
        expires_at: str | None = None,
        requires_confirmation: bool | None = None,
        status: str | None = None,
        supersedes_fact_id: str | None = None,
    ) -> FactRecord:
        with self._locked():
            return self.upsert_fact_unlocked(
                content,
                category=category,
                scope=scope,
                owner=owner,
                source_cursors=source_cursors,
                source_excerpt=source_excerpt,
                confidence=confidence,
                expires_at=expires_at,
                requires_confirmation=requires_confirmation,
                status=status,
                supersedes_fact_id=supersedes_fact_id,
            )

    def upsert_fact_unlocked(
        self,
        content: str,
        *,
        category: str = "note",
        scope: str = "general",
        owner: str = "unknown",
        source_cursors: list[int] | None = None,
        source_excerpt: str = "",
        confidence: float = 1.0,
        expires_at: str | None = None,
        requires_confirmation: bool | None = None,
        status: str | None = None,
        supersedes_fact_id: str | None = None,
    ) -> FactRecord:
        records = self.read_all_unlocked()
        fact = self.upsert_fact_in_records_unlocked(
            records,
            content,
            category=category,
            scope=scope,
            owner=owner,
            source_cursors=source_cursors,
            source_excerpt=source_excerpt,
            confidence=confidence,
            expires_at=expires_at,
            requires_confirmation=requires_confirmation,
            status=status,
            supersedes_fact_id=supersedes_fact_id,
        )
        self._write_records_unlocked(records)
        return fact

    def upsert_fact_in_records_unlocked(
        self,
        records: list[FactRecord],
        content: str,
        *,
        category: str = "note",
        scope: str = "general",
        owner: str = "unknown",
        source_cursors: list[int] | None = None,
        source_excerpt: str = "",
        confidence: float = 1.0,
        expires_at: str | None = None,
        requires_confirmation: bool | None = None,
        status: str | None = None,
        supersedes_fact_id: str | None = None,
    ) -> FactRecord:
        content = content.strip()
        if not content:
            raise ValueError("fact content cannot be empty")
        explicit_status = status
        explicit_requires_confirmation = requires_confirmation
        category = _normalize_category(category)
        scope = _normalize_scope(scope)
        owner = _normalize_owner(owner)
        source_cursors = _normalize_source_cursors(source_cursors)
        confidence = _normalize_confidence(confidence)
        status, requires_confirmation = _resolve_status_and_confirmation(
            category=category,
            status=status,
            requires_confirmation=requires_confirmation,
            config=self.config,
        )
        redacted_excerpt = self._redactor(source_excerpt.strip()) if source_excerpt else ""
        now = datetime.now().isoformat()
        canonical_key = canonical_key_for_fact(content, owner, category, scope)

        existing = next(
            (
                record
                for record in records
                if record.canonical_key == canonical_key
                and record.status in {"active", "pending_confirmation"}
            ),
            None,
        )
        if existing is not None:
            existing.updated_at = now
            existing.last_seen_at = now
            existing.source_cursors = _merge_source_cursors(
                existing.source_cursors,
                source_cursors,
            )
            if redacted_excerpt:
                existing.source_excerpt = redacted_excerpt
            existing.confidence = max(existing.confidence, confidence)
            if expires_at is not None:
                existing.expires_at = expires_at
            if supersedes_fact_id is not None:
                existing.supersedes_fact_id = supersedes_fact_id
            if (
                explicit_status is not None
                or explicit_requires_confirmation is not None
            ):
                existing.status = status
                existing.requires_confirmation = requires_confirmation
            return existing

        if supersedes_fact_id:
            self._deprecate_fact_in_records(
                records,
                supersedes_fact_id,
                updated_at=now,
            )
        record = FactRecord(
            fact_id=self._new_fact_id(records),
            content=content,
            canonical_key=canonical_key,
            category=category,
            scope=scope,
            owner=owner,
            source_cursors=source_cursors,
            source_excerpt=redacted_excerpt,
            confidence=confidence,
            status=status,
            created_at=now,
            updated_at=now,
            last_seen_at=now,
            expires_at=expires_at,
            supersedes_fact_id=supersedes_fact_id,
            requires_confirmation=requires_confirmation,
        )
        records.append(record)
        return record

    def deprecate_fact(
        self,
        fact_id: str,
        *,
        reason: str | None = None,
        superseded_by: str | None = None,
    ) -> bool:
        with self._locked():
            return self.deprecate_fact_unlocked(
                fact_id,
                reason=reason,
                superseded_by=superseded_by,
            )

    def deprecate_fact_unlocked(
        self,
        fact_id: str,
        *,
        reason: str | None = None,
        superseded_by: str | None = None,
    ) -> bool:
        records = self.read_all_unlocked()
        changed = self._deprecate_fact_in_records(
            records,
            fact_id,
            updated_at=datetime.now().isoformat(),
        )
        if changed:
            self._write_records_unlocked(records)
        return changed

    def render_memory_md(self) -> str:
        with self._locked():
            return self.render_memory_md_unlocked()

    def render_memory_md_unlocked(self) -> str:
        return render_memory_md(self.read_all_unlocked())

    def _write_records_unlocked(self, records: list[FactRecord]) -> None:
        text = "".join(
            json.dumps(record.to_dict(), ensure_ascii=False, sort_keys=True) + "\n"
            for record in records
        )
        _write_text_atomic(self.facts_file, text)

    def _new_fact_id(self, records: list[FactRecord]) -> str:
        existing = {record.fact_id for record in records}
        while True:
            fact_id = f"fact_{uuid.uuid4().hex[:12]}"
            if fact_id not in existing:
                return fact_id

    @staticmethod
    def _deprecate_fact_in_records(
        records: list[FactRecord],
        fact_id: str,
        *,
        updated_at: str,
    ) -> bool:
        for record in records:
            if record.fact_id != fact_id:
                continue
            record.status = "deprecated"
            record.updated_at = updated_at
            return True
        return False


def _default_redactor(text: str) -> str:
    from OriginAgent.agent.memory import redact_memory_text

    return redact_memory_text(text)


def _extract_json_text(text: str) -> str:
    stripped = text.strip()
    fenced = re.fullmatch(r"```(?:json)?\s*(.*?)\s*```", stripped, re.DOTALL | re.IGNORECASE)
    if fenced:
        return fenced.group(1).strip()
    return stripped


def _parse_fact_proposal(raw: Any) -> FactProposal:
    if not isinstance(raw, dict):
        raise ValueError("fact proposal must be an object")
    content = _required_str(raw, "content")
    category = _required_str(raw, "category")
    scope = _required_str(raw, "scope")
    owner = _required_str(raw, "owner")
    source_cursors = _parse_cursor_list(raw.get("source_cursors", []), "source_cursors")
    source_excerpt = _optional_str(raw, "source_excerpt", "")
    confidence = _optional_float(raw, "confidence", 0.7)
    expires_at = _optional_nullable_str(raw, "expires_at")
    supersedes_fact_id = _optional_nullable_str(raw, "supersedes_fact_id")
    status = _optional_nullable_str(raw, "status")
    requires_confirmation = raw.get("requires_confirmation")
    if requires_confirmation is not None and not isinstance(requires_confirmation, bool):
        raise ValueError("requires_confirmation must be bool or null")
    reason = _optional_str(raw, "reason", "")
    return FactProposal(
        content=content,
        category=category,
        scope=scope,
        owner=owner,
        source_cursors=source_cursors,
        source_excerpt=source_excerpt,
        confidence=confidence,
        expires_at=expires_at,
        supersedes_fact_id=supersedes_fact_id,
        requires_confirmation=requires_confirmation,
        status=status,
        reason=reason,
    )


def _parse_deprecation_proposal(raw: Any) -> FactDeprecationProposal:
    if not isinstance(raw, dict):
        raise ValueError("deprecation proposal must be an object")
    return FactDeprecationProposal(
        fact_id=_required_str(raw, "fact_id"),
        reason=_optional_str(raw, "reason", ""),
        source_cursors=_parse_cursor_list(raw.get("source_cursors", []), "source_cursors"),
    )


def _required_str(raw: dict[str, Any], key: str) -> str:
    value = raw.get(key)
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{key} must be a non-empty string")
    return value.strip()


def _optional_str(raw: dict[str, Any], key: str, default: str) -> str:
    value = raw.get(key, default)
    if value is None:
        return default
    if not isinstance(value, str):
        raise ValueError(f"{key} must be a string")
    return value.strip()


def _optional_nullable_str(raw: dict[str, Any], key: str) -> str | None:
    value = raw.get(key)
    if value is None:
        return None
    if not isinstance(value, str):
        raise ValueError(f"{key} must be a string or null")
    return value.strip() or None


def _optional_float(raw: dict[str, Any], key: str, default: float) -> float:
    value = raw.get(key, default)
    if isinstance(value, bool):
        raise ValueError(f"{key} must be numeric")
    try:
        return float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError(f"{key} must be numeric") from exc


def _parse_cursor_list(value: Any, key: str) -> list[int]:
    if not isinstance(value, list):
        raise ValueError(f"{key} must be a list")
    cursors: list[int] = []
    for item in value:
        if isinstance(item, bool) or not isinstance(item, int):
            raise ValueError(f"{key} must contain integer cursors")
        if item > 0:
            cursors.append(item)
    return sorted(set(cursors))


def _history_cursor_set(history_entries: list[dict[str, Any]]) -> set[int]:
    cursors: set[int] = set()
    for entry in history_entries:
        raw = entry.get("cursor")
        if isinstance(raw, bool) or not isinstance(raw, int):
            continue
        cursors.add(raw)
    return cursors


def _source_excerpt_matches_history(
    proposal: FactProposal,
    history_entries: list[dict[str, Any]],
) -> bool:
    by_cursor = {
        entry.get("cursor"): str(entry.get("content", ""))
        for entry in history_entries
    }
    excerpt = _normalize_excerpt_for_match(proposal.source_excerpt)
    if not excerpt:
        return False
    for cursor in proposal.source_cursors:
        content = _normalize_excerpt_for_match(by_cursor.get(cursor, ""))
        if excerpt in content:
            return True
    return False


def _normalize_excerpt_for_match(text: str) -> str:
    return re.sub(r"\s+", " ", text.strip().casefold())


def _contains_any(text: str, needles: tuple[str, ...]) -> bool:
    lower = text.casefold()
    return any(needle.casefold() in lower for needle in needles)


def _requires_high_risk_confirmation(*, scope: str, content: str) -> bool:
    return _infer_domain_key(scope, content) in HIGH_RISK_DEVICE_DOMAINS


def _fact_domain_key(record: FactRecord) -> str:
    domain = _infer_domain_key(record.scope, record.content)
    if domain and domain != "general":
        return domain
    scope_head = record.scope.split(".", 1)[0].strip().lower()
    return scope_head or "general"


def _infer_domain_key(scope: str, content: str) -> str:
    text = f"{scope or ''} {content or ''}".casefold()
    for domain in HIGH_RISK_DEVICE_DOMAINS:
        if domain in text:
            return domain
    parts = [part for part in str(scope or "").strip().lower().split(".") if part]
    return parts[0] if parts else "general"


def _issue(code: str, severity: str, message: str) -> ValidationIssue:
    return ValidationIssue(code=code, severity=severity, message=message)


def _normalize_category(category: str | None) -> str:
    normalized = (category or "note").strip().lower()
    if normalized not in VALID_CATEGORIES:
        raise ValueError(f"invalid fact category: {category!r}")
    return normalized


def _normalize_owner(owner: str | None) -> str:
    normalized = (owner or "unknown").strip().lower()
    if normalized not in VALID_OWNERS:
        raise ValueError(f"invalid fact owner: {owner!r}")
    return normalized


def _normalize_scope(scope: str | None) -> str:
    normalized = (scope or "general").strip().lower()
    if not normalized:
        normalized = "general"
    return re.sub(r"\s+", ".", normalized)


def _normalize_status(status: str | None) -> str:
    normalized = (status or "active").strip().lower()
    if normalized not in VALID_STATUSES:
        raise ValueError(f"invalid fact status: {status!r}")
    return normalized


def _normalize_confidence(value: Any) -> float:
    try:
        confidence = float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("confidence must be numeric") from exc
    return max(0.0, min(1.0, confidence))


def _normalize_decay_factor(value: Any) -> float:
    try:
        factor = float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("confidence decay factor must be numeric") from exc
    return max(0.0, min(1.0, factor))


def _parse_datetime(value: str) -> datetime | None:
    text = str(value or "").strip()
    if not text:
        return None
    if text.endswith("Z"):
        text = f"{text[:-1]}+00:00"
    try:
        return datetime.fromisoformat(text)
    except ValueError:
        return None


def _days_between(start: datetime, end: datetime) -> int:
    if start.tzinfo is not None and end.tzinfo is None:
        end = end.replace(tzinfo=timezone.utc)
    elif start.tzinfo is None and end.tzinfo is not None:
        start = start.replace(tzinfo=end.tzinfo)
    delta = end - start
    return max(0, delta.days)


def _calibration_key(proposal_type: str, domain_id: str) -> str:
    proposal = str(proposal_type or "").strip().lower()
    domain = str(domain_id or "").strip().lower()
    if not proposal or not domain:
        return ""
    return f"{proposal}:{domain}"


def _int_value(value: Any, *, default: int) -> int:
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default


def _float_value(value: Any, *, default: float) -> float:
    if isinstance(value, bool):
        return default
    try:
        return float(value)
    except (TypeError, ValueError):
        return default


def _normalize_source_cursors(value: list[int] | Any | None) -> list[int]:
    if value is None:
        return []
    if not isinstance(value, list):
        raise ValueError("source_cursors must be a list")
    cursors: set[int] = set()
    for item in value:
        if isinstance(item, bool) or not isinstance(item, int):
            continue
        if item > 0:
            cursors.add(item)
    return sorted(cursors)


def _merge_source_cursors(existing: list[int], incoming: list[int]) -> list[int]:
    return sorted(set(existing) | set(incoming))


def _resolve_status_and_confirmation(
    *,
    category: str,
    status: str | None,
    requires_confirmation: bool | None,
    config: FactStoreConfig = DEFAULT_FACT_STORE_CONFIG,
) -> tuple[str, bool]:
    if requires_confirmation is not None and not isinstance(requires_confirmation, bool):
        raise ValueError("requires_confirmation must be bool")
    requested_status = _normalize_status(status) if status is not None else None
    if category in set(config.high_risk_categories):
        if requires_confirmation is False and requested_status == "active":
            return "active", False
        if requested_status in {"deprecated", "contradicted"}:
            return requested_status, True
        return "pending_confirmation", True
    return requested_status or "active", bool(requires_confirmation)


def _section_title(category: str) -> str:
    return {
        "preference": "Preferences",
        "routine": "Routines",
        "household": "Household",
        "device": "Devices",
        "policy": "Policies",
        "safety": "Safety",
        "temporary": "Temporary",
        "note": "Notes",
    }[category]


def _render_fact_lines(
    fact: FactRecord,
    *,
    include_category: bool = False,
) -> list[str]:
    content_lines = fact.content.splitlines() or [""]
    lines = [f"- {content_lines[0]}"]
    for line in content_lines[1:]:
        lines.append(f"  {line}")
    if include_category:
        lines.append(f"  - category: {fact.category}")
    lines.append(f"  - scope: {fact.scope}")
    lines.append(f"  - confidence: {_format_confidence(fact.confidence)}")
    if fact.expires_at:
        lines.append(f"  - expires_at: {fact.expires_at}")
    source = _format_source(fact.source_cursors)
    if source:
        lines.append(f"  - source: {source}")
    return lines


def _format_confidence(value: float) -> str:
    if value == int(value):
        return f"{value:.1f}"
    return f"{value:.2f}".rstrip("0").rstrip(".")


def _format_source(source_cursors: list[int]) -> str:
    recent = source_cursors[-3:]
    if not recent:
        return ""
    if len(recent) == 1:
        return f"cursor {recent[0]}"
    return "cursors " + ", ".join(str(cursor) for cursor in recent)


def _write_text_atomic(path: Path, text: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = path.with_name(f".{path.name}.tmp")
    try:
        with open(tmp_path, "w", encoding="utf-8") as f:
            f.write(text)
            f.flush()
            os.fsync(f.fileno())
        os.replace(tmp_path, path)
        _fsync_parent(path)
    except BaseException:
        tmp_path.unlink(missing_ok=True)
        raise


def _fsync_parent(path: Path) -> None:
    with suppress(PermissionError, OSError):
        fd = os.open(str(path.parent), os.O_RDONLY)
        try:
            os.fsync(fd)
        finally:
            os.close(fd)
