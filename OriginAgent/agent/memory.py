"""Memory system: pure file I/O store, lightweight Consolidator, and Dream processor."""

from __future__ import annotations

import asyncio
import hashlib
import json
import os
import re
import shutil
import weakref
from contextlib import suppress
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
from typing import TYPE_CHECKING, Any, Callable, Iterator

import tiktoken
from filelock import FileLock
from loguru import logger

from OriginAgent.agent.facts import (
    DreamFactApplyResult,
    DreamFactProposalBatch,
    FactRecord,
    FactStore,
    MAX_DEPRECATIONS_PER_BATCH,
    ValidationIssue,
    canonical_key_for_fact,
    domain_id_for_fact,
    parse_fact_proposal_response,
    render_memory_md as render_facts_memory_md,
    validate_deprecation_proposal,
    validate_fact_proposal,
)
from OriginAgent.agent.auxiliary_llm import call_llm
from OriginAgent.agent.evolution import (
    OpportunitySignalStore,
    detect_skill_opportunity_candidates,
    detect_workflow_opportunity_candidates,
)
from OriginAgent.agent.runner import AgentRunner, AgentRunSpec
from OriginAgent.agent.tools.registry import ToolRegistry
from OriginAgent.session.manager import Session
from OriginAgent.utils.gitstore import GitStore
from OriginAgent.utils.helpers import (
    ensure_dir,
    estimate_message_tokens,
    estimate_prompt_tokens_chain,
    find_legal_message_start,
    strip_think,
    truncate_text,
)
from OriginAgent.utils.prompt_templates import render_template

if TYPE_CHECKING:
    from OriginAgent.agent.auxiliary_llm import AuxiliaryLLMRouter
    from OriginAgent.providers.base import LLMProvider
    from OriginAgent.session.manager import SessionManager


_PRIVATE_KEY_RE = re.compile(
    r"-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----.*?-----END [A-Z0-9 ]*PRIVATE KEY-----",
    re.DOTALL,
)
_BEARER_TOKEN_RE = re.compile(r"(?i)\bBearer\s+[A-Za-z0-9._~+/=-]{8,}")
_OPENAI_KEY_RE = re.compile(r"\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b")
_GITHUB_TOKEN_RE = re.compile(r"\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{20,}\b")
_SECRET_ASSIGNMENT_RE = re.compile(
    r"(?i)\b(api[_-]?key|token|secret|password)\b(\s*[:=]\s*)([\"']?)"
    r"[^\"'\s,;]{8,}([\"']?)"
)
_EMAIL_RE = re.compile(r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b")
_CHINA_ID_RE = re.compile(
    r"(?<!\d)\d{6}(?:18|19|20)\d{2}(?:0[1-9]|1[0-2])"
    r"(?:0[1-9]|[12]\d|3[01])\d{3}[\dXx](?!\d)"
)
_LONG_NUMBER_RE = re.compile(r"(?<!\d)\d{16,}(?!\d)")


def redact_memory_text(text: str) -> str:
    """Redact high-risk secrets and PII before memory history is persisted."""
    if not text:
        return text
    text = _PRIVATE_KEY_RE.sub("[REDACTED_PRIVATE_KEY]", text)
    text = _BEARER_TOKEN_RE.sub("[REDACTED_BEARER_TOKEN]", text)
    text = _OPENAI_KEY_RE.sub("[REDACTED_SECRET]", text)
    text = _GITHUB_TOKEN_RE.sub("[REDACTED_SECRET]", text)
    text = _SECRET_ASSIGNMENT_RE.sub(
        lambda match: (
            f"{match.group(1)}{match.group(2)}"
            f"{match.group(3)}[REDACTED_SECRET]{match.group(4)}"
        ),
        text,
    )
    text = _EMAIL_RE.sub("[REDACTED_EMAIL]", text)
    text = _CHINA_ID_RE.sub("[REDACTED_ID]", text)
    text = _LONG_NUMBER_RE.sub("[REDACTED_LONG_NUMBER]", text)
    return text


# ---------------------------------------------------------------------------
# MemoryStore — pure file I/O layer
# ---------------------------------------------------------------------------

class MemoryStore:
    """Pure file I/O for memory files: MEMORY.md, history.jsonl, SOUL.md, USER.md."""

    _DEFAULT_MAX_HISTORY = 1000
    _LEGACY_ENTRY_START_RE = re.compile(r"^\[(\d{4}-\d{2}-\d{2}[^\]]*)\]\s*")
    _LEGACY_TIMESTAMP_RE = re.compile(r"^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2})\]\s*")
    _LEGACY_RAW_MESSAGE_RE = re.compile(
        r"^\[\d{4}-\d{2}-\d{2}[^\]]*\]\s+[A-Z][A-Z0-9_]*(?:\s+\[tools:\s*[^\]]+\])?:"
    )

    def __init__(self, workspace: Path, max_history_entries: int = _DEFAULT_MAX_HISTORY):
        self.workspace = workspace
        self.max_history_entries = max_history_entries
        self.memory_dir = ensure_dir(workspace / "memory")
        self._lock_file = self.memory_dir / ".lock"
        self.memory_file = self.memory_dir / "MEMORY.md"
        self.facts_file = self.memory_dir / "facts.jsonl"
        self.history_file = self.memory_dir / "history.jsonl"
        self.legacy_history_file = self.memory_dir / "HISTORY.md"
        self.soul_file = workspace / "SOUL.md"
        self.user_file = workspace / "USER.md"
        self._cursor_file = self.memory_dir / ".cursor"
        self._dream_cursor_file = self.memory_dir / ".dream_cursor"
        self._corruption_logged = False  # rate-limit non-int cursor warning
        self._oversize_logged = False  # rate-limit oversized-entry warning
        self.fact_store = FactStore(
            workspace=self.workspace,
            facts_file=self.facts_file,
            lock_factory=self._locked,
            redactor=redact_memory_text,
        )
        self._git = GitStore(workspace, tracked_files=[
            "SOUL.md", "USER.md", "memory/MEMORY.md",
            "memory/facts.jsonl", "memory/.dream_cursor",
        ])
        self._maybe_migrate_legacy_history()

    @property
    def git(self) -> GitStore:
        return self._git

    # -- generic helpers -----------------------------------------------------

    def _locked(self) -> FileLock:
        """Return the cross-process lock for memory state changes.

        This lock covers memory/history files under memory/ as well as
        workspace-level memory files such as SOUL.md and USER.md.
        """
        return FileLock(str(self._lock_file))

    @staticmethod
    def read_file(path: Path) -> str:
        try:
            return path.read_text(encoding="utf-8")
        except FileNotFoundError:
            return ""

    @staticmethod
    def _fsync_parent(path: Path) -> None:
        with suppress(PermissionError, OSError):
            fd = os.open(str(path.parent), os.O_RDONLY)
            try:
                os.fsync(fd)
            finally:
                os.close(fd)

    def _write_text_atomic(self, path: Path, text: str) -> None:
        """Atomically write text and fsync both the file and parent directory."""
        path.parent.mkdir(parents=True, exist_ok=True)
        tmp_path = path.with_name(f".{path.name}.tmp")
        try:
            with open(tmp_path, "w", encoding="utf-8") as f:
                f.write(text)
                f.flush()
                os.fsync(f.fileno())
            os.replace(tmp_path, path)
            self._fsync_parent(path)
        except BaseException:
            tmp_path.unlink(missing_ok=True)
            raise

    def _maybe_migrate_legacy_history(self) -> None:
        """One-time upgrade from legacy HISTORY.md to history.jsonl.

        The migration is best-effort and prioritizes preserving as much content
        as possible over perfect parsing.
        """
        if not self.legacy_history_file.exists():
            return
        with self._locked():
            if self.history_file.exists() and self.history_file.stat().st_size > 0:
                return

            try:
                legacy_text = self.legacy_history_file.read_text(
                    encoding="utf-8",
                    errors="replace",
                )
            except OSError:
                logger.exception("Failed to read legacy HISTORY.md for migration")
                return

            entries = self._parse_legacy_history(legacy_text)
            try:
                if entries:
                    self._write_entries(entries)
                    last_cursor = entries[-1]["cursor"]
                    self._write_text_atomic(self._cursor_file, str(last_cursor))
                    # Default to "already processed" so upgrades do not replay the
                    # user's entire historical archive into Dream on first start.
                    self._write_text_atomic(self._dream_cursor_file, str(last_cursor))

                backup_path = self._next_legacy_backup_path()
                self.legacy_history_file.replace(backup_path)
                self._fsync_parent(self.legacy_history_file)
                logger.info(
                    "Migrated legacy HISTORY.md to history.jsonl ({} entries)",
                    len(entries),
                )
            except Exception:
                logger.exception("Failed to migrate legacy HISTORY.md")

    def _parse_legacy_history(self, text: str) -> list[dict[str, Any]]:
        normalized = text.replace("\r\n", "\n").replace("\r", "\n").strip()
        if not normalized:
            return []

        fallback_timestamp = self._legacy_fallback_timestamp()
        entries: list[dict[str, Any]] = []
        chunks = self._split_legacy_history_chunks(normalized)

        for cursor, chunk in enumerate(chunks, start=1):
            timestamp = fallback_timestamp
            content = chunk
            match = self._LEGACY_TIMESTAMP_RE.match(chunk)
            if match:
                timestamp = match.group(1)
                remainder = chunk[match.end():].lstrip()
                if remainder:
                    content = remainder

            entries.append({
                "cursor": cursor,
                "timestamp": timestamp,
                "content": content,
            })
        return entries

    def _split_legacy_history_chunks(self, text: str) -> list[str]:
        lines = text.split("\n")
        chunks: list[str] = []
        current: list[str] = []
        saw_blank_separator = False

        for line in lines:
            if saw_blank_separator and line.strip() and current:
                chunks.append("\n".join(current).strip())
                current = [line]
                saw_blank_separator = False
                continue
            if self._should_start_new_legacy_chunk(line, current):
                chunks.append("\n".join(current).strip())
                current = [line]
                saw_blank_separator = False
                continue
            current.append(line)
            saw_blank_separator = not line.strip()

        if current:
            chunks.append("\n".join(current).strip())
        return [chunk for chunk in chunks if chunk]

    def _should_start_new_legacy_chunk(self, line: str, current: list[str]) -> bool:
        if not current:
            return False
        if not self._LEGACY_ENTRY_START_RE.match(line):
            return False
        if self._is_raw_legacy_chunk(current) and self._LEGACY_RAW_MESSAGE_RE.match(line):
            return False
        return True

    def _is_raw_legacy_chunk(self, lines: list[str]) -> bool:
        first_nonempty = next((line for line in lines if line.strip()), "")
        match = self._LEGACY_TIMESTAMP_RE.match(first_nonempty)
        if not match:
            return False
        return first_nonempty[match.end():].lstrip().startswith("[RAW]")

    def _legacy_fallback_timestamp(self) -> str:
        try:
            return datetime.fromtimestamp(
                self.legacy_history_file.stat().st_mtime,
            ).strftime("%Y-%m-%d %H:%M")
        except OSError:
            return datetime.now().strftime("%Y-%m-%d %H:%M")

    def _next_legacy_backup_path(self) -> Path:
        candidate = self.memory_dir / "HISTORY.md.bak"
        suffix = 2
        while candidate.exists():
            candidate = self.memory_dir / f"HISTORY.md.bak.{suffix}"
            suffix += 1
        return candidate

    # -- MEMORY.md (long-term facts) -----------------------------------------

    def read_memory(self) -> str:
        return self.read_file(self.memory_file)

    def write_memory(self, content: str) -> None:
        with self._locked():
            self._write_text_atomic(self.memory_file, content)

    def rebuild_memory_from_facts(self) -> None:
        with self._locked():
            content = self.fact_store.render_memory_md_unlocked()
            self._write_text_atomic(self.memory_file, content)

    def upsert_fact_and_rebuild_memory(
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
            fact = self.fact_store.upsert_fact_unlocked(
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
            memory_md = self.fact_store.render_memory_md_unlocked()
            self._write_text_atomic(self.memory_file, memory_md)
            return fact

    def update_fact_confidence_and_rebuild_memory(
        self,
        fact_id: str,
        confidence: float,
    ) -> FactRecord | None:
        with self._locked():
            records = self.fact_store.read_all_unlocked()
            fact = self.fact_store.update_confidence_in_records_unlocked(
                records,
                fact_id,
                confidence,
            )
            if fact is None:
                return None
            self.fact_store._write_records_unlocked(records)
            self._write_text_atomic(self.memory_file, render_facts_memory_md(records))
            return fact

    def scale_fact_confidence_for_canonical_key_and_rebuild_memory(
        self,
        canonical_key: str,
        multiplier: float,
        *,
        status: str = "active",
    ) -> FactRecord | None:
        canonical_key = str(canonical_key or "").strip()
        if not canonical_key:
            return None
        multiplier = max(0.0, min(1.0, float(multiplier)))
        with self._locked():
            records = self.fact_store.read_all_unlocked()
            target = next(
                (
                    record
                    for record in records
                    if record.canonical_key == canonical_key
                    and record.status == status
                ),
                None,
            )
            if target is None:
                return None
            fact = self.fact_store.update_confidence_in_records_unlocked(
                records,
                target.fact_id,
                target.confidence * multiplier,
            )
            if fact is None:
                return None
            self.fact_store._write_records_unlocked(records)
            self._write_text_atomic(self.memory_file, render_facts_memory_md(records))
            return fact

    def decay_fact_confidence_and_rebuild_memory(self) -> int:
        with self._locked():
            records = self.fact_store.read_all_unlocked()
            changed = self.fact_store.decay_confidence_in_records_unlocked(records)
            if not changed:
                return 0
            self.fact_store._write_records_unlocked(records)
            self._write_text_atomic(self.memory_file, render_facts_memory_md(records))
            return changed

    def apply_fact_proposals_and_rebuild_memory(
        self,
        batch: DreamFactProposalBatch,
        *,
        history_entries: list[dict[str, Any]],
    ) -> DreamFactApplyResult:
        cursors = [
            entry.get("cursor")
            for entry in history_entries
            if isinstance(entry.get("cursor"), int)
            and not isinstance(entry.get("cursor"), bool)
        ]
        if not cursors:
            result = DreamFactApplyResult(parse_rejected=list(batch.parse_rejected))
            issue = ValidationIssue(
                code="empty_history_batch",
                severity="reject",
                message="No valid history cursors were available for validation.",
            )
            for proposal in [*batch.facts_to_upsert, *batch.facts_to_deprecate]:
                result.rejected.append((proposal, [issue]))
            return result

        batch_cursor_min = min(cursors)
        batch_cursor_max = max(cursors)
        with self._locked():
            records = self.fact_store.read_all_unlocked()
            result = DreamFactApplyResult(parse_rejected=list(batch.parse_rejected))
            active_budget_remaining = self.fact_store.config.auto_active_budget
            deprecation_count = 0

            calibrated_upserts = [
                (
                    self.fact_store.calibrate_confidence(
                        "fact",
                        domain_id_for_fact(proposal.scope, proposal.content),
                        proposal.confidence,
                    ),
                    proposal,
                )
                for proposal in batch.facts_to_upsert
            ]
            calibrated_upserts.sort(key=lambda item: item[0], reverse=True)

            for calibrated_confidence, proposal in calibrated_upserts:
                proposal.confidence = calibrated_confidence
                validation = validate_fact_proposal(
                    proposal,
                    existing_facts=records,
                    history_entries=history_entries,
                    batch_cursor_min=batch_cursor_min,
                    batch_cursor_max=batch_cursor_max,
                )
                issues = list(validation.issues)
                decision = validation.decision
                existing_active_fact = None
                if decision == "active":
                    canonical_key = canonical_key_for_fact(
                        proposal.content,
                        proposal.owner,
                        proposal.category,
                        proposal.scope,
                    )
                    existing_active_fact = next(
                        (
                            record
                            for record in records
                            if record.canonical_key == canonical_key
                            and record.status == "active"
                        ),
                        None,
                    )
                if (
                    decision == "active"
                    and existing_active_fact is None
                    and validation.confidence
                    < self.fact_store.config.auto_active_confidence_threshold
                ):
                    issues.append(ValidationIssue(
                        code="active_confidence_below_threshold",
                        severity="pending",
                        message="Fact confidence is below the automatic activation threshold.",
                    ))
                    decision = "pending_confirmation"
                elif (
                    decision == "active"
                    and existing_active_fact is None
                    and active_budget_remaining < validation.confidence
                ):
                    issues.append(ValidationIssue(
                        code="active_confidence_budget_exceeded",
                        severity="pending",
                        message="Automatic active fact confidence budget exceeded for this Dream batch.",
                    ))
                    decision = "pending_confirmation"
                if decision == "reject":
                    result.rejected.append((proposal, issues))
                    continue

                status = "active" if decision == "active" else "pending_confirmation"
                requires_confirmation = decision != "active"
                fact = self.fact_store.upsert_fact_in_records_unlocked(
                    records,
                    proposal.content,
                    category=proposal.category,
                    scope=proposal.scope,
                    owner=proposal.owner,
                    source_cursors=proposal.source_cursors,
                    source_excerpt=proposal.source_excerpt,
                    confidence=validation.confidence,
                    expires_at=proposal.expires_at,
                    requires_confirmation=requires_confirmation,
                    status=status,
                    supersedes_fact_id=proposal.supersedes_fact_id,
                )
                if decision == "active":
                    if existing_active_fact is None:
                        active_budget_remaining -= validation.confidence
                    result.accepted.append(fact)
                else:
                    result.pending.append(fact)

            for proposal in batch.facts_to_deprecate:
                issues: list[ValidationIssue] = []
                if deprecation_count >= MAX_DEPRECATIONS_PER_BATCH:
                    issues.append(ValidationIssue(
                        code="deprecation_limit_exceeded",
                        severity="reject",
                        message="Automatic deprecation limit exceeded for this Dream batch.",
                    ))
                    result.rejected.append((proposal, issues))
                    continue
                validation = validate_deprecation_proposal(
                    proposal,
                    existing_facts=records,
                    history_entries=history_entries,
                    batch_cursor_min=batch_cursor_min,
                    batch_cursor_max=batch_cursor_max,
                )
                if validation.decision == "reject":
                    result.rejected.append((proposal, list(validation.issues)))
                    continue
                changed = self.fact_store._deprecate_fact_in_records(
                    records,
                    proposal.fact_id,
                    updated_at=datetime.now().isoformat(),
                )
                if changed:
                    deprecation_count += 1
                    result.deprecated.append(proposal.fact_id)
                else:
                    result.rejected.append((
                        proposal,
                        [ValidationIssue(
                            code="unknown_fact_id",
                            severity="reject",
                            message="Deprecation fact_id does not exist.",
                        )],
                    ))

            self.fact_store._write_records_unlocked(records)
            self._write_text_atomic(self.memory_file, render_facts_memory_md(records))
            return result

    def seed_legacy_memory_as_note(self) -> FactRecord | None:
        with self._locked():
            if self.fact_store.read_all_unlocked():
                return None
            legacy_memory = self.read_file(self.memory_file)
            if not legacy_memory.strip():
                return None
            fact = self.fact_store.upsert_fact_unlocked(
                legacy_memory,
                category="note",
                scope="legacy.memory",
                owner="system",
                source_cursors=[],
                source_excerpt="legacy MEMORY.md import",
                confidence=0.5,
                status="active",
            )
            memory_md = self.fact_store.render_memory_md_unlocked()
            self._write_text_atomic(self.memory_file, memory_md)
            return fact

    # -- SOUL.md -------------------------------------------------------------

    def read_soul(self) -> str:
        return self.read_file(self.soul_file)

    def write_soul(self, content: str) -> None:
        with self._locked():
            self._write_text_atomic(self.soul_file, content)

    # -- USER.md -------------------------------------------------------------

    def read_user(self) -> str:
        return self.read_file(self.user_file)

    def write_user(self, content: str) -> None:
        with self._locked():
            self._write_text_atomic(self.user_file, content)

    # -- context injection (used by context.py) ------------------------------

    def get_memory_context(self) -> str:
        long_term = self.read_memory()
        return f"## Long-term Memory\n{long_term}" if long_term else ""

    # -- history.jsonl — append-only, JSONL format ---------------------------

    def append_history(self, entry: str, *, max_chars: int | None = None) -> int:
        """Append *entry* to history.jsonl and return its auto-incrementing cursor.

        Entries are passed through `strip_think` to drop template-level leaks
        (e.g. unclosed `<think` prefixes, `<channel|>` markers) before being
        persisted. If the cleaned content is empty but the raw entry wasn't,
        the record is persisted with an empty string rather than falling back
        to the raw leak — otherwise `strip_think`'s guarantees would be
        undone by history replay / consolidation downstream.

        A defensive cap (*max_chars*, default ``_HISTORY_ENTRY_HARD_CAP``) is
        applied as a final safety net: individual callers should cap their own
        content more tightly; this default only exists to catch unintentional
        large writes (e.g. an LLM echoing its input back as a "summary").
        """
        limit = max_chars if max_chars is not None else _HISTORY_ENTRY_HARD_CAP
        ts = datetime.now().strftime("%Y-%m-%d %H:%M")
        raw = entry.rstrip()
        if len(raw) > limit:
            if not self._oversize_logged:
                self._oversize_logged = True
                logger.warning(
                    "history entry exceeds {} chars ({}); truncating. "
                    "Usually means a caller forgot its own cap; "
                    "further occurrences suppressed.",
                    limit, len(raw),
                )
        content = strip_think(raw)
        content = redact_memory_text(content)
        if len(content) > limit:
            content = truncate_text(content, limit)
        if raw and not content:
            stripped_to_empty = True
        else:
            stripped_to_empty = False

        with self._locked():
            cursor = self._next_cursor_unlocked()
            if stripped_to_empty:
                logger.debug(
                    "history entry {} stripped to empty (likely template leak); "
                    "persisting empty content to avoid re-polluting context",
                    cursor,
                )
            record = {"cursor": cursor, "timestamp": ts, "content": content}
            self.history_file.parent.mkdir(parents=True, exist_ok=True)
            with open(self.history_file, "a", encoding="utf-8") as f:
                f.write(json.dumps(record, ensure_ascii=False) + "\n")
                f.flush()
                os.fsync(f.fileno())
            self._write_text_atomic(self._cursor_file, str(cursor))
        return cursor

    @staticmethod
    def _valid_cursor(value: Any) -> int | None:
        """Int cursors only — reject bool (``isinstance(True, int)`` is True)."""
        if isinstance(value, bool) or not isinstance(value, int):
            return None
        return value

    def _iter_valid_entries(self) -> Iterator[tuple[dict[str, Any], int]]:
        """Yield ``(entry, cursor)`` for entries with int cursors; warn once on corruption."""
        poisoned: Any = None
        for entry in self._read_entries():
            if not isinstance(entry, dict):
                poisoned = entry
                continue
            raw = entry.get("cursor")
            if raw is None:
                continue
            cursor = self._valid_cursor(raw)
            if cursor is None:
                poisoned = raw
                continue
            yield entry, cursor
        if poisoned is not None and not self._corruption_logged:
            self._corruption_logged = True
            logger.warning(
                "history.jsonl contains a non-int cursor ({!r}); dropping it. "
                "Usually caused by an external writer; further occurrences suppressed.",
                poisoned,
            )

    def _max_history_cursor_unlocked(self) -> int:
        """Return the max valid cursor currently present in history.jsonl."""
        return max((c for _, c in self._iter_valid_entries()), default=0)

    def _next_cursor_unlocked(self) -> int:
        """Read cursor state and return the next value. Caller must hold lock."""
        cursor_file_value = 0
        if self._cursor_file.exists():
            with suppress(ValueError, OSError):
                cursor_file_value = int(self._cursor_file.read_text(encoding="utf-8").strip())
        return max(cursor_file_value, self._max_history_cursor_unlocked()) + 1

    def _next_cursor(self) -> int:
        """Read the current cursor counter and return the next value."""
        with self._locked():
            return self._next_cursor_unlocked()

    def read_unprocessed_history(self, since_cursor: int) -> list[dict[str, Any]]:
        """Return history entries with a valid cursor > *since_cursor*."""
        return [e for e, c in self._iter_valid_entries() if c > since_cursor]

    def compact_history(self, processed_through: int) -> None:
        """Drop oldest processed entries while preserving unsafe/unknown lines."""
        with self._locked():
            self._compact_history_unlocked(processed_through)

    def _compact_history_unlocked(self, processed_through: int) -> None:
        """Compact raw history lines without dropping invalid or unprocessed data."""
        if self.max_history_entries <= 0:
            return
        raw_lines = self._read_history_lines_unlocked()
        if len(raw_lines) <= self.max_history_entries:
            return

        delete_budget = len(raw_lines) - self.max_history_entries
        kept: list[str] = []

        for raw_line in raw_lines:
            cursor: int | None = None
            try:
                parsed = json.loads(raw_line.strip())
            except json.JSONDecodeError:
                kept.append(raw_line)
                continue
            if isinstance(parsed, dict):
                cursor = self._valid_cursor(parsed.get("cursor"))
            if cursor is None or cursor > processed_through:
                kept.append(raw_line)
                continue
            if delete_budget > 0:
                delete_budget -= 1
                continue
            kept.append(raw_line)

        self._write_history_lines_unlocked(kept)

    # -- JSONL helpers -------------------------------------------------------

    def _read_history_lines_unlocked(self) -> list[str]:
        with suppress(FileNotFoundError):
            with open(self.history_file, "r", encoding="utf-8") as f:
                return f.readlines()
        return []

    def _read_entries(self) -> list[Any]:
        """Read all entries from history.jsonl."""
        entries: list[Any] = []
        with suppress(FileNotFoundError):
            with open(self.history_file, "r", encoding="utf-8") as f:
                for line in f:
                    line = line.strip()
                    if line:
                        try:
                            entries.append(json.loads(line))
                        except json.JSONDecodeError:
                            continue

        return entries

    def _read_last_entry(self) -> dict[str, Any] | None:
        """Read the last entry from the JSONL file efficiently."""
        try:
            with open(self.history_file, "rb") as f:
                f.seek(0, 2)
                size = f.tell()
                if size == 0:
                    return None
                read_size = min(size, 4096)
                f.seek(size - read_size)
                data = f.read().decode("utf-8")
                lines = [line for line in data.split("\n") if line.strip()]
                if not lines:
                    return None
                return json.loads(lines[-1])
        except (FileNotFoundError, json.JSONDecodeError, UnicodeDecodeError):
            return None

    def _write_entries(self, entries: list[dict[str, Any]]) -> None:
        """Overwrite history.jsonl with the given entries (atomic write)."""
        self._write_history_lines_unlocked([
            json.dumps(entry, ensure_ascii=False) + "\n"
            for entry in entries
        ])

    def _write_history_lines_unlocked(self, lines: list[str]) -> None:
        """Overwrite history.jsonl with raw lines. Caller should hold lock when needed."""
        tmp_path = self.history_file.with_suffix(self.history_file.suffix + ".tmp")
        try:
            with open(tmp_path, "w", encoding="utf-8") as f:
                for line in lines:
                    f.write(line)
                f.flush()
                os.fsync(f.fileno())
            os.replace(tmp_path, self.history_file)
            self._fsync_parent(self.history_file)
        except BaseException:
            tmp_path.unlink(missing_ok=True)
            raise

    # -- dream cursor --------------------------------------------------------

    def get_last_dream_cursor(self) -> int:
        if self._dream_cursor_file.exists():
            with suppress(ValueError, OSError):
                return int(self._dream_cursor_file.read_text(encoding="utf-8").strip())
        return 0

    def set_last_dream_cursor(self, cursor: int) -> None:
        with self._locked():
            self._write_text_atomic(self._dream_cursor_file, str(cursor))

    def mark_dream_processed(self, cursor: int) -> None:
        """Mark Dream history as processed, then compact only processed entries."""
        with self._locked():
            self._write_text_atomic(self._dream_cursor_file, str(cursor))
            try:
                self._compact_history_unlocked(cursor)
            except Exception:
                logger.exception(
                    "Dream cursor advanced to {}, but history compact failed",
                    cursor,
                )

    # -- message formatting utility ------------------------------------------

    @staticmethod
    def _format_messages(messages: list[dict]) -> str:
        lines = []
        for message in messages:
            if not message.get("content"):
                continue
            tools = f" [tools: {', '.join(message['tools_used'])}]" if message.get("tools_used") else ""
            lines.append(
                f"[{message.get('timestamp', '?')[:16]}] {message['role'].upper()}{tools}: {message['content']}"
            )
        return "\n".join(lines)

    def raw_archive(self, messages: list[dict], *, max_chars: int | None = None) -> None:
        """Fallback: dump raw messages to history.jsonl without LLM summarization."""
        limit = max_chars if max_chars is not None else _RAW_ARCHIVE_MAX_CHARS
        formatted = self._format_messages(messages)
        formatted = redact_memory_text(formatted)
        formatted = truncate_text(formatted, limit)
        self.append_history(
            f"[RAW] {len(messages)} messages\n"
            f"{formatted}"
        )
        logger.warning(
            "Memory consolidation degraded: raw-archived {} messages", len(messages)
        )



# ---------------------------------------------------------------------------
# Consolidator — lightweight token-budget triggered consolidation
# ---------------------------------------------------------------------------


# Individual history.jsonl writers cap their own payloads tightly; the
# _HISTORY_ENTRY_HARD_CAP at append_history() is a belt-and-suspenders default
# that catches any new caller that forgot to set its own cap.
_RAW_ARCHIVE_MAX_CHARS = 16_000       # fallback dump (LLM failed)
_ARCHIVE_SUMMARY_MAX_CHARS = 8_000    # LLM-produced consolidation summary
_HISTORY_ENTRY_HARD_CAP = 64_000      # emergency cap in append_history
_RECENT_SUMMARIES_KEY = "_recent_summaries"
_LEGACY_LAST_SUMMARY_KEY = "_last_summary"
_MAX_RECENT_SUMMARIES = 5


@dataclass(frozen=True)
class ArchiveResult:
    summary: str
    history_cursor: int


def _valid_recent_summary_entries(value: Any) -> list[dict[str, Any]]:
    if not isinstance(value, list):
        return []
    entries: list[dict[str, Any]] = []
    for item in value:
        if not isinstance(item, dict):
            continue
        text = item.get("text")
        if not isinstance(text, str) or not text:
            continue
        entries.append(item)
    return entries


def record_recent_summary(
    session: Session,
    result: ArchiveResult | None,
    *,
    last_active: datetime | None = None,
    created_at: datetime | None = None,
) -> bool:
    if (
        not isinstance(result, ArchiveResult)
        or not result.summary
        or result.summary == "(nothing)"
    ):
        return False
    created = created_at or datetime.now()
    active = last_active or session.updated_at
    entries = _valid_recent_summary_entries(session.metadata.get(_RECENT_SUMMARIES_KEY))
    entries.append({
        "text": result.summary,
        "history_cursor": result.history_cursor,
        "created_at": created.isoformat(),
        "last_active": active.isoformat(),
    })
    session.metadata[_RECENT_SUMMARIES_KEY] = entries[-_MAX_RECENT_SUMMARIES:]
    session.metadata.pop(_LEGACY_LAST_SUMMARY_KEY, None)
    return True


def session_summary_text(session: Session, *, max_chars: int = 6000) -> str | None:
    entries = _valid_recent_summary_entries(session.metadata.get(_RECENT_SUMMARIES_KEY))
    if entries:
        sections = ["## Recent Session Summaries"]
        for entry in entries[-_MAX_RECENT_SUMMARIES:]:
            cursor = entry.get("history_cursor")
            if isinstance(cursor, bool) or not isinstance(cursor, int) or cursor <= 0:
                label = "[history cursor unknown]"
            else:
                label = f"[history cursor {cursor}]"
            sections.append(f"{label}\n{entry['text']}")
        return truncate_text("\n\n".join(sections), max_chars)

    legacy = session.metadata.get(_LEGACY_LAST_SUMMARY_KEY)
    if isinstance(legacy, dict):
        text = legacy.get("text")
        last_active = legacy.get("last_active")
        if isinstance(text, str) and text:
            if isinstance(last_active, str) and last_active:
                return truncate_text(
                    f"Previous conversation summary (last active {last_active}):\n{text}",
                    max_chars,
                )
            return truncate_text(text, max_chars)
    if isinstance(legacy, str) and legacy:
        return truncate_text(legacy, max_chars)
    return None


class Consolidator:
    """Lightweight consolidation: summarizes evicted messages into history.jsonl."""

    _MAX_CONSOLIDATION_ROUNDS = 5

    _SAFETY_BUFFER = 1024  # extra headroom for tokenizer estimation drift

    def __init__(
        self,
        store: MemoryStore,
        provider: LLMProvider,
        model: str,
        sessions: SessionManager,
        context_window_tokens: int,
        build_messages: Callable[..., list[dict[str, Any]]],
        get_tool_definitions: Callable[[], list[dict[str, Any]]],
        max_completion_tokens: int = 4096,
        consolidation_ratio: float = 0.5,
        auxiliary_router: AuxiliaryLLMRouter | None = None,
    ):
        self.store = store
        self.provider = provider
        self.model = model
        self.auxiliary_router = auxiliary_router
        self.sessions = sessions
        self.context_window_tokens = context_window_tokens
        self.max_completion_tokens = max_completion_tokens
        self.consolidation_ratio = consolidation_ratio
        self._build_messages = build_messages
        self._get_tool_definitions = get_tool_definitions
        self._locks: weakref.WeakValueDictionary[str, asyncio.Lock] = (
            weakref.WeakValueDictionary()
        )

    def set_provider(
        self,
        provider: LLMProvider,
        model: str,
        context_window_tokens: int,
    ) -> None:
        self.provider = provider
        self.model = model
        self.context_window_tokens = context_window_tokens
        self.max_completion_tokens = provider.generation.max_tokens
        if self.auxiliary_router is not None:
            self.auxiliary_router.set_primary(provider, model)

    def get_lock(self, session_key: str) -> asyncio.Lock:
        """Return the shared consolidation lock for one session."""
        return self._locks.setdefault(session_key, asyncio.Lock())

    def pick_consolidation_boundary(
        self,
        session: Session,
        tokens_to_remove: int,
    ) -> tuple[int, int] | None:
        """Pick a user-turn boundary that removes enough old prompt tokens."""
        start = session.last_consolidated
        if start >= len(session.messages) or tokens_to_remove <= 0:
            return None

        removed_tokens = 0
        last_boundary: tuple[int, int] | None = None
        for idx in range(start, len(session.messages)):
            message = session.messages[idx]
            if idx > start and message.get("role") == "user":
                last_boundary = (idx, removed_tokens)
                if removed_tokens >= tokens_to_remove:
                    return last_boundary
            removed_tokens += estimate_message_tokens(message)

        return last_boundary

    @staticmethod
    def _full_unconsolidated_history(
        session: Session,
        *,
        include_timestamps: bool = False,
    ) -> list[dict[str, Any]]:
        """Return the whole unconsolidated tail for consolidation decisions."""
        unconsolidated_count = len(session.messages) - session.last_consolidated
        if unconsolidated_count <= 0:
            return []
        return session.get_history(
            max_messages=unconsolidated_count,
            include_timestamps=include_timestamps,
        )

    @staticmethod
    def _replay_overflow_boundary(
        session: Session,
        replay_max_messages: int | None,
    ) -> int | None:
        if not replay_max_messages or replay_max_messages <= 0:
            return None
        tail = list(enumerate(session.messages[session.last_consolidated:], session.last_consolidated))
        if len(tail) <= replay_max_messages:
            return None

        sliced = tail[-replay_max_messages:]
        for i, (_idx, message) in enumerate(sliced):
            if message.get("role") == "user":
                start = i
                if i > 0 and sliced[i - 1][1].get("_channel_delivery"):
                    start = i - 1
                sliced = sliced[start:]
                break

        legal_start = find_legal_message_start([message for _idx, message in sliced])
        if legal_start:
            sliced = sliced[legal_start:]
        if not sliced:
            return len(session.messages)

        first_visible_idx = sliced[0][0]
        if first_visible_idx <= session.last_consolidated:
            return None
        return first_visible_idx

    async def _consolidate_replay_overflow(
        self,
        session: Session,
        replay_max_messages: int | None,
    ) -> ArchiveResult | None:
        """Archive messages that would be hidden by the replay message window."""
        end_idx = self._replay_overflow_boundary(session, replay_max_messages)
        if end_idx is None:
            return None
        chunk = session.messages[session.last_consolidated:end_idx]
        if not chunk:
            return None
        logger.info(
            "Replay-window consolidation for {}: chunk={} msgs, replay_max={}",
            session.key,
            len(chunk),
            replay_max_messages,
        )
        result = await self.archive(chunk)
        session.last_consolidated = end_idx
        self.sessions.save(session)
        return result

    def _persist_recent_summary(self, session: Session, result: ArchiveResult | None) -> None:
        if record_recent_summary(session, result):
            self.sessions.save(session)

    def estimate_session_prompt_tokens(
        self,
        session: Session,
    ) -> tuple[int, str]:
        """Estimate prompt size from the full unconsolidated session tail."""
        history = self._full_unconsolidated_history(session, include_timestamps=True)
        channel, chat_id = (session.key.split(":", 1) if ":" in session.key else (None, None))
        # Include archived summary in estimation so the budget accounts for it.
        summary = session_summary_text(session)
        probe_messages = self._build_messages(
            history=history,
            current_message="[token-probe]",
            channel=channel,
            chat_id=chat_id,
            sender_id=None,
            session_summary=summary,
        )
        return estimate_prompt_tokens_chain(
            self.provider,
            self.model,
            probe_messages,
            self._get_tool_definitions(),
        )

    @property
    def _input_token_budget(self) -> int:
        """Available input token budget for consolidation LLM."""
        return self.context_window_tokens - self.max_completion_tokens - self._SAFETY_BUFFER

    def _truncate_to_token_budget(self, text: str) -> str:
        """Truncate text so it fits within the consolidation LLM's token budget."""
        budget = self._input_token_budget
        if budget <= 0:
            return truncate_text(text, _RAW_ARCHIVE_MAX_CHARS)
        try:
            enc = tiktoken.get_encoding("cl100k_base")
            tokens = enc.encode(text)
            if len(tokens) <= budget:
                return text
            return enc.decode(tokens[:budget]) + "\n... (truncated)"
        except Exception:
            return truncate_text(text, budget * 4)

    async def archive(self, messages: list[dict]) -> ArchiveResult | None:
        """Summarize messages via LLM and append to history.jsonl.

        Returns the summary text and history cursor on success, None if nothing to archive.
        """
        if not messages:
            return None
        try:
            formatted = MemoryStore._format_messages(messages)
            formatted = self._truncate_to_token_budget(formatted)
            response = await call_llm(
                task="consolidation",
                router=self.auxiliary_router,
                provider=self.provider,
                model=self.model,
                messages=[
                    {
                        "role": "system",
                        "content": render_template(
                            "agent/consolidator_archive.md",
                            strip=True,
                        ),
                    },
                    {"role": "user", "content": formatted},
                ],
                tools=None,
                tool_choice=None,
            )
            if response.finish_reason == "error":
                raise RuntimeError(f"LLM returned error: {response.content}")
            summary = redact_memory_text(response.content or "[no summary]")
            cursor = self.store.append_history(summary, max_chars=_ARCHIVE_SUMMARY_MAX_CHARS)
            return ArchiveResult(summary=summary, history_cursor=cursor)
        except Exception:
            logger.warning("Consolidation LLM call failed, raw-dumping to history")
            self.store.raw_archive(messages)
            return None

    async def maybe_consolidate_by_tokens(
        self,
        session: Session,
        *,
        replay_max_messages: int | None = None,
    ) -> None:
        """Loop: archive old messages until prompt fits within safe budget.

        The budget reserves space for completion tokens and a safety buffer
        so the LLM request never exceeds the context window.
        """
        if not session.messages or self.context_window_tokens <= 0:
            return

        lock = self.get_lock(session.key)
        async with lock:
            budget = self._input_token_budget
            target = int(budget * self.consolidation_ratio)
            last_result = await self._consolidate_replay_overflow(
                session,
                replay_max_messages,
            )
            try:
                estimated, source = self.estimate_session_prompt_tokens(
                    session,
                )
            except Exception:
                logger.exception("Token estimation failed for {}", session.key)
                estimated, source = 0, "error"
            if estimated <= 0:
                self._persist_recent_summary(session, last_result)
                return
            if estimated < budget:
                unconsolidated_count = len(session.messages) - session.last_consolidated
                logger.debug(
                    "Token consolidation idle {}: {}/{} via {}, msgs={}",
                    session.key,
                    estimated,
                    self.context_window_tokens,
                    source,
                    unconsolidated_count,
                )
                self._persist_recent_summary(session, last_result)
                return

            for round_num in range(self._MAX_CONSOLIDATION_ROUNDS):
                if estimated <= target:
                    break

                boundary = self.pick_consolidation_boundary(session, max(1, estimated - target))
                if boundary is None:
                    logger.debug(
                        "Token consolidation: no safe boundary for {} (round {})",
                        session.key,
                        round_num,
                    )
                    break

                end_idx = boundary[0]

                chunk = session.messages[session.last_consolidated:end_idx]
                if not chunk:
                    break

                logger.info(
                    "Token consolidation round {} for {}: {}/{} via {}, chunk={} msgs",
                    round_num,
                    session.key,
                    estimated,
                    self.context_window_tokens,
                    source,
                    len(chunk),
                )
                result = await self.archive(chunk)
                # Advance the cursor either way: on success the chunk was
                # summarized; on failure archive() already raw-archived it as
                # a breadcrumb. Re-archiving the same chunk on the next call
                # would just emit duplicate [RAW] entries.
                if result:
                    last_result = result
                session.last_consolidated = end_idx
                self.sessions.save(session)
                if not result:
                    # LLM is degraded — stop hammering it this call;
                    # the next invocation can retry a fresh chunk.
                    break

                try:
                    estimated, source = self.estimate_session_prompt_tokens(
                        session,
                    )
                except Exception:
                    logger.exception("Token estimation failed for {}", session.key)
                    estimated, source = 0, "error"
                if estimated <= 0:
                    break

            # Persist the latest summary to session metadata so it can be injected
            # into the runtime context on the next prepare_session() call, aligning
            # the summary injection strategy with AutoCompact._archive().
            self._persist_recent_summary(session, last_result)


# ---------------------------------------------------------------------------
# Dream — heavyweight cron-scheduled memory consolidation
# ---------------------------------------------------------------------------


# Single source of truth for the staleness threshold used in _annotate_with_ages
# *and* in the Phase 1 prompt template (passed as `stale_threshold_days`).
# Keep code and prompt aligned — if you bump this, the LLM's instruction string
# updates automatically.
_STALE_THRESHOLD_DAYS = 14


class MemoryWorkspaceSnapshot:
    """File-level snapshot for Dream's writable memory surface."""

    _TRACKED_FILES = (
        Path("SOUL.md"),
        Path("USER.md"),
        Path("memory/MEMORY.md"),
        Path("memory/facts.jsonl"),
    )
    _TRACKED_DIRS = (Path("skills"),)

    def __init__(self, workspace: Path):
        self.workspace = workspace
        self.files, self.dirs = self._capture_state()

    @staticmethod
    def _sha256(data: bytes) -> str:
        return hashlib.sha256(data).hexdigest()

    def _capture_state(self) -> tuple[dict[Path, bytes], set[Path]]:
        files: dict[Path, bytes] = {}
        dirs: set[Path] = set()

        for rel in self._TRACKED_FILES:
            path = self.workspace / rel
            if path.is_file():
                files[rel] = path.read_bytes()

        for root_rel in self._TRACKED_DIRS:
            root = self.workspace / root_rel
            if not root.exists():
                continue
            if root.is_dir():
                dirs.add(root_rel)
                for child in root.rglob("*"):
                    rel = child.relative_to(self.workspace)
                    if child.is_dir():
                        dirs.add(rel)
                    elif child.is_file():
                        files[rel] = child.read_bytes()

        return files, dirs

    def _file_hashes(self, files: dict[Path, bytes]) -> dict[Path, str]:
        return {rel: self._sha256(data) for rel, data in files.items()}

    def _remove_path(self, path: Path) -> None:
        if path.is_dir():
            shutil.rmtree(path)
        else:
            path.unlink(missing_ok=True)

    def restore(self) -> bool:
        """Restore the snapshot and verify that the tracked tree matches exactly."""
        try:
            # Remove files and directories created under tracked directories.
            for root_rel in self._TRACKED_DIRS:
                root = self.workspace / root_rel
                if not root.exists():
                    continue
                for child in sorted(root.rglob("*"), key=lambda p: len(p.parts), reverse=True):
                    rel = child.relative_to(self.workspace)
                    if child.is_file() and rel not in self.files:
                        child.unlink(missing_ok=True)
                    elif child.is_dir() and rel not in self.dirs:
                        with suppress(OSError):
                            child.rmdir()
                if root_rel not in self.dirs and root.exists():
                    with suppress(OSError):
                        root.rmdir()

            # Remove tracked top-level files that did not exist in the snapshot.
            for rel in self._TRACKED_FILES:
                if rel not in self.files:
                    path = self.workspace / rel
                    if path.exists():
                        self._remove_path(path)

            # Restore directories before files so empty pre-existing skill dirs survive.
            for rel in sorted(self.dirs, key=lambda p: len(p.parts)):
                (self.workspace / rel).mkdir(parents=True, exist_ok=True)

            for rel, data in self.files.items():
                path = self.workspace / rel
                if path.exists() and path.is_dir():
                    shutil.rmtree(path)
                path.parent.mkdir(parents=True, exist_ok=True)
                path.write_bytes(data)

            current_files, current_dirs = self._capture_state()
            if current_dirs != self.dirs or self._file_hashes(current_files) != self._file_hashes(self.files):
                logger.error("Dream snapshot restore verification failed")
                return False
            return True
        except Exception:
            logger.exception("Dream snapshot restore failed")
            return False


class Dream:
    """Two-phase memory processor for structured long-term memory.

    Phase 1 proposes structured facts that deterministic validators route into
    facts.jsonl and the derived MEMORY.md. Phase 2 may maintain SOUL.md,
    USER.md, or skills, but generated memory files are guarded against edits.
    """

    # Caps on prompt-bound inputs so Dream's LLM calls never exceed the model's
    # context window just because a file (or a legacy large history entry) grew
    # unexpectedly. Each file still appears in full via read_file when the agent
    # needs it in Phase 2 — these caps only bound the Phase 1/2 prompt preview.
    _MEMORY_FILE_MAX_CHARS = 32_000
    _SOUL_FILE_MAX_CHARS = 16_000
    _USER_FILE_MAX_CHARS = 16_000
    _HISTORY_ENTRY_PREVIEW_MAX_CHARS = 4_000

    def __init__(
        self,
        store: MemoryStore,
        provider: LLMProvider,
        model: str,
        max_batch_size: int = 20,
        max_iterations: int = 10,
        max_tool_result_chars: int = 16_000,
        annotate_line_ages: bool = True,
        auxiliary_router: AuxiliaryLLMRouter | None = None,
        evolution_config: Any | None = None,
    ):
        self.store = store
        self.provider = provider
        self.model = model
        self.auxiliary_router = auxiliary_router
        self.max_batch_size = max_batch_size
        self.max_iterations = max_iterations
        self.max_tool_result_chars = max_tool_result_chars
        # Kill switch for the git-blame-based per-line age annotation in Phase 1.
        # Default True keeps the #3212 behavior; set False to feed MEMORY.md raw
        # (e.g. if a specific LLM reacts poorly to the `← Nd` suffix).
        self.annotate_line_ages = annotate_line_ages
        self.evolution_config = evolution_config
        self.opportunity_signals = OpportunitySignalStore(store.workspace)
        runner_provider = (
            auxiliary_router.task_provider("dream_phase2")
            if auxiliary_router is not None
            else provider
        )
        self._runner = AgentRunner(runner_provider)
        self._tools = self._build_tools()

    def set_provider(self, provider: LLMProvider, model: str) -> None:
        self.provider = provider
        self.model = model
        if self.auxiliary_router is not None:
            self.auxiliary_router.set_primary(provider, model)
            self._runner.provider = self.auxiliary_router.task_provider("dream_phase2")
        else:
            self._runner.provider = provider

    # -- tool registry -------------------------------------------------------

    def _build_tools(self) -> ToolRegistry:
        """Build a minimal tool registry for the Dream agent."""
        from OriginAgent.agent.skills import BUILTIN_SKILLS_DIR
        from OriginAgent.agent.tools.file_state import FileStates
        from OriginAgent.agent.tools.filesystem import EditFileTool, ReadFileTool, WriteFileTool

        tools = ToolRegistry()
        workspace = self.store.workspace
        # Allow reading builtin skills for reference during skill creation
        extra_read = [BUILTIN_SKILLS_DIR] if BUILTIN_SKILLS_DIR.exists() else None
        # Dream gets its own FileStates so its caches stay isolated from the
        # main loop's sessions (issue #3571).
        file_states = FileStates()
        tools.register(ReadFileTool(
            workspace=workspace,
            allowed_dir=workspace,
            extra_allowed_dirs=extra_read,
            file_states=file_states,
        ))
        tools.register(EditFileTool(workspace=workspace, allowed_dir=workspace, file_states=file_states))
        # write_file resolves relative paths from workspace root, but can only
        # write under skills/ so the prompt can safely use skills/<name>/SKILL.md.
        skills_dir = workspace / "skills"
        skills_dir.mkdir(parents=True, exist_ok=True)
        tools.register(WriteFileTool(workspace=workspace, allowed_dir=skills_dir, file_states=file_states))
        return tools

    # -- skill listing --------------------------------------------------------

    def _list_existing_skills(self) -> list[str]:
        """List existing skills as 'name — description' for dedup context."""
        import re as _re

        from OriginAgent.agent.skills import BUILTIN_SKILLS_DIR

        desc_re = _re.compile(r"^description:\s*(.+)$", _re.MULTILINE | _re.IGNORECASE)
        entries: dict[str, str] = {}
        for base in (self.store.workspace / "skills", BUILTIN_SKILLS_DIR):
            if not base.exists():
                continue
            for d in base.iterdir():
                if not d.is_dir():
                    continue
                skill_md = d / "SKILL.md"
                if not skill_md.exists():
                    continue
                # Prefer workspace skills over builtin (same name)
                if d.name in entries and base == BUILTIN_SKILLS_DIR:
                    continue
                content = skill_md.read_text(encoding="utf-8")[:500]
                m = desc_re.search(content)
                desc = m.group(1).strip() if m else "(no description)"
                entries[d.name] = desc
        return [f"{name} — {desc}" for name, desc in sorted(entries.items())]

    # -- main entry ----------------------------------------------------------

    def _annotate_with_ages(self, content: str) -> str:
        """Append per-line age suffixes to MEMORY.md content.

        Each non-blank line whose age exceeds ``_STALE_THRESHOLD_DAYS`` gets a
        suffix like ``← 30d`` indicating days since last modification.
        Returns the original content unchanged if git is unavailable,
        annotate fails, or the line count doesn't match the age count
        (which can happen with an uncommitted working-tree edit — better to
        skip annotation than to tag the wrong line).
        SOUL.md and USER.md are never annotated.
        """
        file_path = "memory/MEMORY.md"
        try:
            ages = self.store.git.line_ages(file_path)
        except Exception:
            logger.debug("line_ages failed for {}", file_path)
            return content
        if not ages:
            return content

        had_trailing = content.endswith("\n")
        lines = content.splitlines()
        # If HEAD-blob line count disagrees with the working-tree content we
        # received, ages would be assigned to the wrong lines — skip entirely
        # and feed the LLM un-annotated content rather than misleading data.
        if len(lines) != len(ages):
            logger.debug(
                "line_ages length mismatch for {} (lines={}, ages={}); skipping annotation",
                file_path, len(lines), len(ages),
            )
            return content

        annotated: list[str] = []
        for line, age in zip(lines, ages):
            if not line.strip():
                annotated.append(line)
                continue
            if age.age_days > _STALE_THRESHOLD_DAYS:
                annotated.append(f"{line}  \u2190 {age.age_days}d")
            else:
                annotated.append(line)
        result = "\n".join(annotated)
        if had_trailing:
            result += "\n"
        return result

    def _format_current_facts(self) -> str:
        facts = self.store.fact_store.read_all()
        payload = [
            {
                "fact_id": fact.fact_id,
                "status": fact.status,
                "category": fact.category,
                "scope": fact.scope,
                "content": fact.content,
            }
            for fact in facts
        ]
        if not payload:
            return "[]"
        return truncate_text(
            json.dumps(payload, ensure_ascii=False, indent=2),
            self._MEMORY_FILE_MAX_CHARS,
        )

    def _memory_fact_hashes(self) -> tuple[str | None, str | None]:
        return (
            self._file_hash(self.store.memory_file),
            self._file_hash(self.store.facts_file),
        )

    @staticmethod
    def _file_hash(path: Path) -> str | None:
        try:
            return hashlib.sha256(path.read_bytes()).hexdigest()
        except FileNotFoundError:
            return None

    @staticmethod
    def _format_apply_result(result: DreamFactApplyResult) -> str:
        return (
            "active={active} pending={pending} rejected={rejected} "
            "parse_rejected={parse_rejected} deprecated={deprecated}"
        ).format(
            active=len(result.accepted),
            pending=len(result.pending),
            rejected=len(result.rejected),
            parse_rejected=len(result.parse_rejected),
            deprecated=len(result.deprecated),
        )

    def _record_opportunity_signals(self, batch: list[dict[str, Any]]) -> int:
        retention_days = int(getattr(self.evolution_config, "signal_retention_days", 30) or 30)
        candidates = detect_workflow_opportunity_candidates(
            batch,
            min_evidence_sources=1,
        )
        candidates.extend(detect_skill_opportunity_candidates(
            batch,
            min_evidence_sources=1,
        ))
        if not candidates:
            return 0
        updated = self.opportunity_signals.upsert_candidates(
            candidates,
            retention_days=retention_days,
        )
        return len(updated)

    async def run(self) -> bool:
        """Process unprocessed history entries. Returns True if work was done."""
        from OriginAgent.agent.skills import BUILTIN_SKILLS_DIR

        last_cursor = self.store.get_last_dream_cursor()
        entries = self.store.read_unprocessed_history(since_cursor=last_cursor)
        if not entries:
            return False

        batch = entries[: self.max_batch_size]
        logger.info(
            "Dream: processing {} entries (cursor {}→{}), batch={}",
            len(entries), last_cursor, batch[-1]["cursor"], len(batch),
        )
        snapshot = MemoryWorkspaceSnapshot(self.store.workspace)

        # Build history text for LLM — cap each entry so a legacy oversized
        # record (e.g. pre-#3412 raw_archive dump) can't blow up the prompt.
        history_text = "\n".join(
            f"[cursor {e['cursor']}] [{e['timestamp']}] "
            f"{truncate_text(e['content'], self._HISTORY_ENTRY_PREVIEW_MAX_CHARS)}"
            for e in batch
        )
        try:
            signal_count = self._record_opportunity_signals(batch)
            if signal_count:
                logger.info("Dream recorded {} evolution opportunity signal(s)", signal_count)
        except Exception:
            logger.exception("Dream opportunity signal collection failed")

        # Current file contents + per-line age annotations (MEMORY.md only).
        # Each file is capped in the *prompt preview* only; Phase 2 still sees
        # the full file via the read_file tool.
        current_date = datetime.now().strftime("%Y-%m-%d")
        raw_memory = self.store.read_memory() or "(empty)"
        annotated_memory = (
            self._annotate_with_ages(raw_memory)
            if self.annotate_line_ages
            else raw_memory
        )
        current_memory = truncate_text(annotated_memory, self._MEMORY_FILE_MAX_CHARS)
        current_soul = truncate_text(
            self.store.read_soul() or "(empty)", self._SOUL_FILE_MAX_CHARS,
        )
        current_user = truncate_text(
            self.store.read_user() or "(empty)", self._USER_FILE_MAX_CHARS,
        )

        file_context = (
            f"## Current Date\n{current_date}\n\n"
            f"## Current MEMORY.md ({len(current_memory)} chars)\n{current_memory}\n\n"
            f"## Current SOUL.md ({len(current_soul)} chars)\n{current_soul}\n\n"
            f"## Current USER.md ({len(current_user)} chars)\n{current_user}"
        )
        facts_context = self._format_current_facts()

        # Phase 1: propose structured facts.
        phase1_prompt = (
            f"## Conversation History\n{history_text}\n\n"
            f"## Current Facts\n{facts_context}\n\n"
            f"{file_context}"
        )

        try:
            phase1_response = await call_llm(
                task="dream_phase1",
                router=self.auxiliary_router,
                provider=self.provider,
                model=self.model,
                messages=[
                    {
                        "role": "system",
                        "content": render_template(
                            "agent/dream_phase1.md",
                            strip=True,
                            stale_threshold_days=_STALE_THRESHOLD_DAYS,
                        ),
                    },
                    {"role": "user", "content": phase1_prompt},
                ],
                tools=None,
                tool_choice=None,
            )
            proposal_json = phase1_response.content or ""
            logger.debug(
                "Dream Phase 1 fact proposal JSON ({} chars): {}",
                len(proposal_json),
                proposal_json[:500],
            )
        except Exception:
            logger.exception("Dream Phase 1 failed")
            return False

        try:
            proposal_batch = parse_fact_proposal_response(proposal_json)
        except Exception:
            logger.exception("Dream Phase 1 returned invalid fact proposal JSON")
            if not snapshot.restore():
                logger.error("Dream parse failure: snapshot restore failed")
            return False

        decayed_fact_count = self.store.decay_fact_confidence_and_rebuild_memory()
        if decayed_fact_count:
            logger.info("Dream decayed confidence for {} active fact(s)", decayed_fact_count)

        try:
            apply_result = self.store.apply_fact_proposals_and_rebuild_memory(
                proposal_batch,
                history_entries=batch,
            )
        except Exception:
            logger.exception("Dream fact proposal apply failed")
            if not snapshot.restore():
                logger.error("Dream fact apply failure: snapshot restore failed")
            return False
        logger.info(
            "Dream fact proposals: active={} pending={} rejected={} deprecated={}",
            len(apply_result.accepted),
            len(apply_result.pending),
            len(apply_result.rejected) + len(apply_result.parse_rejected),
            len(apply_result.deprecated),
        )
        post_apply_hashes = self._memory_fact_hashes()

        # Phase 2: Delegate to AgentRunner for non-MEMORY maintenance only.
        existing_skills = self._list_existing_skills()
        skills_section = ""
        if existing_skills:
            skills_section = (
                "\n\n## Existing Skills\n"
                + "\n".join(f"- {s}" for s in existing_skills)
            )
        phase2_prompt = (
            f"## Fact Proposal Apply Result\n"
            f"{self._format_apply_result(apply_result)}\n\n"
            f"## Fact Proposal JSON\n{proposal_json}\n\n"
            f"{file_context}{skills_section}"
        )

        tools = self._tools
        skill_creator_path = BUILTIN_SKILLS_DIR / "skill-creator" / "SKILL.md"
        messages: list[dict[str, Any]] = [
            {
                "role": "system",
                "content": render_template(
                    "agent/dream_phase2.md",
                    strip=True,
                    skill_creator_path=str(skill_creator_path),
                ),
            },
            {"role": "user", "content": phase2_prompt},
        ]

        try:
            result = await self._runner.run(AgentRunSpec(
                initial_messages=messages,
                tools=tools,
                model=self.model,
                max_iterations=self.max_iterations,
                max_tool_result_chars=self.max_tool_result_chars,
                fail_on_tool_error=False,
            ))
            logger.debug(
                "Dream Phase 2 complete: stop_reason={}, tool_events={}",
                result.stop_reason, len(result.tool_events),
            )
            for ev in (result.tool_events or []):
                logger.info("Dream tool_event: name={}, status={}, detail={}", ev.get("name"), ev.get("status"), ev.get("detail", "")[:200])
        except Exception:
            logger.exception("Dream Phase 2 failed")
            result = None

        # Build changelog from tool events
        changelog: list[str] = []
        if result and result.tool_events:
            for event in result.tool_events:
                if event["status"] == "ok":
                    changelog.append(f"{event['name']}: {event['detail']}")
        fact_changes = (
            len(apply_result.accepted)
            + len(apply_result.pending)
            + len(apply_result.deprecated)
            + decayed_fact_count
        )
        if fact_changes:
            fact_summary = self._format_apply_result(apply_result)
            if decayed_fact_count:
                fact_summary = f"{fact_summary} decayed={decayed_fact_count}"
            changelog.insert(0, f"facts: {fact_summary}")

        # Only advance cursor on successful completion to prevent silent loss
        if result and result.stop_reason == "completed":
            if self._memory_fact_hashes() != post_apply_hashes:
                if not snapshot.restore():
                    logger.error(
                        "Dream Phase 2 modified generated memory state and "
                        "snapshot restore failed; cursor NOT advanced",
                    )
                    return False
                logger.warning(
                    "Dream Phase 2 modified memory/MEMORY.md or "
                    "memory/facts.jsonl; restored snapshot and cursor NOT advanced",
                )
                return False
            new_cursor = batch[-1]["cursor"]
            self.store.mark_dream_processed(new_cursor)
            logger.info(
                "Dream done: {} change(s), cursor advanced to {}",
                len(changelog), new_cursor,
            )
        else:
            reason = result.stop_reason if result else "exception"
            if not snapshot.restore():
                logger.error(
                    "Dream incomplete ({}): snapshot restore failed; "
                    "cursor NOT advanced",
                    reason,
                )
                return False
            logger.warning(
                "Dream incomplete ({}): cursor NOT advanced, will retry next cron cycle",
                reason,
            )
            return False

        # Git auto-commit (only when there are actual changes)
        if changelog and self.store.git.is_initialized():
            ts = batch[-1]["timestamp"]
            summary = f"dream: {ts}, {len(changelog)} change(s)"
            commit_msg = f"{summary}\n\n{proposal_json.strip()}"
            sha = self.store.git.auto_commit(commit_msg)
            if sha:
                logger.info("Dream commit: {}", sha)

        return True
