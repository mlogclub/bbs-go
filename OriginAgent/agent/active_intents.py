"""Bounded proactive active-intent support for idle sessions."""

from __future__ import annotations

import json
import os
import threading
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Literal

from OriginAgent.agent.confirmation import PendingConfirmationStore
from OriginAgent.agent.facts import FactStore
from OriginAgent.agent.reminders import ReminderStore
from OriginAgent.bus.events import InboundMessage
from OriginAgent.bus.queue import MessageBus
from OriginAgent.session.goal_state import goal_state_raw, parse_goal_state
from OriginAgent.session.manager import Session, SessionManager
from OriginAgent.utils.helpers import ensure_dir, truncate_text

ActiveIntentOutcome = Literal["emitted", "suppressed", "skipped"]
ActiveIntentType = Literal["goal_nudge", "pending_confirmation_nudge", "scheduled_reminder"]

_SUMMARY_MAX_CHARS = 240
_RECENT_SCAN_LIMIT = 200


def _utcnow() -> datetime:
    return datetime.now(timezone.utc)


def _utcnow_iso() -> str:
    return _utcnow().isoformat()


def _parse_iso(value: str | None) -> datetime | None:
    if not value:
        return None
    try:
        parsed = datetime.fromisoformat(value)
    except ValueError:
        return None
    if parsed.tzinfo is None:
        return parsed.replace(tzinfo=timezone.utc)
    return parsed.astimezone(timezone.utc)


def _summarize_text(value: Any, max_chars: int = _SUMMARY_MAX_CHARS) -> str:
    text = str(value or "").strip()
    if not text:
        return ""
    return truncate_text(text, max_chars).replace("\n... (truncated)", " ...")


@dataclass(frozen=True)
class ActiveIntentConfig:
    enabled: bool = False
    interval_seconds: int = 30
    session_cooldown_seconds: int = 300
    intent_cooldown_seconds: int = 300
    max_messages_per_session_per_pass: int = 1


@dataclass(frozen=True)
class ActiveIntentCandidate:
    intent_type: ActiveIntentType
    intent_id: str
    content: str
    source_type: str
    source_reference: str
    summary: str = ""


@dataclass(frozen=True)
class ActiveIntentRecord:
    timestamp: str
    session_key: str
    intent_type: str
    intent_id: str
    source_type: str
    source_reference: str
    outcome: ActiveIntentOutcome
    summary: str = ""
    suppression_reason: str | None = None

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


class JsonlActiveIntentLedger:
    """Append-only ledger for proactive emission and suppression decisions."""

    def __init__(self, workspace: Path):
        root = Path(workspace) / "memory" / "active_intents"
        self._records_path = root / "records.jsonl"
        self._lock = threading.Lock()

    def append(self, record: ActiveIntentRecord) -> None:
        with self._lock:
            ensure_dir(self._records_path.parent)
            with self._records_path.open("a", encoding="utf-8") as handle:
                handle.write(json.dumps(record.to_dict(), ensure_ascii=False, sort_keys=True))
                handle.write("\n")
                handle.flush()
                os.fsync(handle.fileno())

    def recent(self, limit: int = _RECENT_SCAN_LIMIT) -> list[dict[str, Any]]:
        if not self._records_path.exists():
            return []
        items: list[dict[str, Any]] = []
        try:
            with self._records_path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        payload = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(payload, dict):
                        items.append(payload)
        except Exception:
            return []
        if limit <= 0:
            return items
        return items[-limit:]


class ActiveIntentService:
    """Produce bounded internal proactive messages for eligible idle sessions."""

    def __init__(
        self,
        *,
        workspace: Path,
        bus: MessageBus,
        sessions: SessionManager,
        confirmation_store: PendingConfirmationStore,
        fact_store: FactStore,
        reminder_store: ReminderStore | None = None,
        config: ActiveIntentConfig,
    ) -> None:
        self.workspace = Path(workspace)
        self.bus = bus
        self.sessions = sessions
        self.confirmation_store = confirmation_store
        self.fact_store = fact_store
        self.reminder_store = reminder_store or ReminderStore(workspace)
        self.config = config
        self.ledger = JsonlActiveIntentLedger(workspace)

    def session_keys(self) -> list[str]:
        seen: set[str] = set()
        keys: list[str] = []
        for item in self.sessions.list_sessions():
            key = str(item.get("key") or "").strip()
            if key and key not in seen:
                seen.add(key)
                keys.append(key)
        return keys

    def eligible_session(
        self,
        session_key: str,
        *,
        active_task_count: int,
        running_subagents: int,
    ) -> tuple[bool, str | None]:
        if not self.config.enabled:
            return False, "disabled"
        if active_task_count > 0:
            return False, "active_tasks"
        if running_subagents > 0:
            return False, "running_subagents"
        return True, None

    async def process_session(
        self,
        session_key: str,
        *,
        active_task_count: int,
        running_subagents: int,
    ) -> list[ActiveIntentCandidate]:
        eligible, reason = self.eligible_session(
            session_key,
            active_task_count=active_task_count,
            running_subagents=running_subagents,
        )
        if not eligible:
            self._record_skip(session_key, reason or "ineligible")
            return []

        session = self.sessions.get_or_create(session_key)
        candidates = self._build_candidates(session)
        emitted: list[ActiveIntentCandidate] = []
        for candidate in candidates:
            allowed, suppression_reason = self._passes_cooldown(session_key, candidate.intent_id)
            if not allowed:
                self.ledger.append(ActiveIntentRecord(
                    timestamp=_utcnow_iso(),
                    session_key=session_key,
                    intent_type=candidate.intent_type,
                    intent_id=candidate.intent_id,
                    source_type=candidate.source_type,
                    source_reference=candidate.source_reference,
                    outcome="suppressed",
                    summary=candidate.summary,
                    suppression_reason=suppression_reason,
                ))
                continue
            await self.bus.publish_inbound(self._build_message(session, candidate))
            if candidate.intent_type == "scheduled_reminder" and candidate.source_reference:
                self.reminder_store.mark_fired(candidate.source_reference)
            self.ledger.append(ActiveIntentRecord(
                timestamp=_utcnow_iso(),
                session_key=session_key,
                intent_type=candidate.intent_type,
                intent_id=candidate.intent_id,
                source_type=candidate.source_type,
                source_reference=candidate.source_reference,
                outcome="emitted",
                summary=candidate.summary,
            ))
            emitted.append(candidate)
            if len(emitted) >= self.config.max_messages_per_session_per_pass:
                break
        return emitted

    def _build_candidates(self, session: Session) -> list[ActiveIntentCandidate]:
        candidates: list[ActiveIntentCandidate] = []
        goal_candidate = self._goal_candidate(session)
        if goal_candidate is not None:
            candidates.append(goal_candidate)
        confirmation_candidate = self._pending_confirmation_candidate(session)
        if confirmation_candidate is not None:
            candidates.append(confirmation_candidate)
        reminder_candidate = self._scheduled_reminder_candidate(session)
        if reminder_candidate is not None:
            candidates.append(reminder_candidate)
        return candidates

    def _goal_candidate(self, session: Session) -> ActiveIntentCandidate | None:
        goal = parse_goal_state(goal_state_raw(session.metadata))
        if not isinstance(goal, dict) or goal.get("status") != "active":
            return None
        objective = str(goal.get("objective") or "").strip()
        if not objective:
            return None
        summary = str(goal.get("ui_summary") or "").strip() or _summarize_text(objective, max_chars=80)
        started_at = str(goal.get("started_at") or "").strip()
        content = (
            "Active goal follow-up: there is still an unfinished sustained goal in this chat.\n"
            f"Goal: {summary}\n"
            "If work should continue, resume it. If the goal is no longer relevant, clarify or close it."
        )
        return ActiveIntentCandidate(
            intent_type="goal_nudge",
            intent_id=f"goal_nudge:{session.key}:{started_at or summary}",
            content=content,
            source_type="goal_state",
            source_reference=started_at or summary,
            summary=summary,
        )

    def _pending_confirmation_candidate(self, session: Session) -> ActiveIntentCandidate | None:
        pending = [
            item for item in self.confirmation_store.read_all()
            if item.status in {"pending", "notified"} and item.metadata.get("session_key") == session.key
        ]
        if pending:
            confirmation = pending[0]
            prompt = _summarize_text(confirmation.prompt, max_chars=160)
            return ActiveIntentCandidate(
                intent_type="pending_confirmation_nudge",
                intent_id=f"pending_confirmation:{session.key}:{confirmation.confirmation_id}",
                content=(
                    "Pending confirmation follow-up: there is a safety or fact confirmation still waiting.\n"
                    f"Pending item: {prompt}\n"
                    "If the user is ready, ask for confirmation or help them resolve it."
                ),
                source_type="pending_confirmation",
                source_reference=confirmation.confirmation_id,
                summary=prompt,
            )

        pending_facts = [
            record for record in self.fact_store.list_active(include_pending=True)
            if record.status == "pending_confirmation" and record.scope == session.key
        ]
        if not pending_facts:
            return None
        fact = pending_facts[0]
        summary = _summarize_text(fact.content, max_chars=160)
        return ActiveIntentCandidate(
            intent_type="pending_confirmation_nudge",
            intent_id=f"pending_fact:{session.key}:{fact.fact_id}",
            content=(
                "Pending fact follow-up: there is an unconfirmed fact associated with this chat.\n"
                f"Fact: {summary}\n"
                "If useful, ask the user to confirm, correct, or dismiss it."
            ),
            source_type="pending_fact",
            source_reference=fact.fact_id,
            summary=summary,
        )

    def _scheduled_reminder_candidate(self, session: Session) -> ActiveIntentCandidate | None:
        due = [
            record for record in self.reminder_store.list_due()
            if record.session_key == session.key
        ]
        if not due:
            return None
        reminder = due[0]
        summary = _summarize_text(reminder.content, max_chars=160)
        return ActiveIntentCandidate(
            intent_type="scheduled_reminder",
            intent_id=f"reminder:{reminder.reminder_id}",
            content=(
                "Scheduled reminder: the due time for this reminder has arrived.\n"
                f"Reminder: {summary}\n"
                "Deliver this reminder clearly to the user and keep the follow-up bounded to the reminder itself."
            ),
            source_type="reminder",
            source_reference=reminder.reminder_id,
            summary=summary,
        )

    def _passes_cooldown(self, session_key: str, intent_id: str) -> tuple[bool, str | None]:
        now = _utcnow()
        session_cutoff = self.config.session_cooldown_seconds
        intent_cutoff = self.config.intent_cooldown_seconds
        for record in reversed(self.ledger.recent()):
            outcome = str(record.get("outcome") or "")
            if outcome != "emitted":
                continue
            ts = _parse_iso(str(record.get("timestamp") or ""))
            if ts is None:
                continue
            delta = (now - ts).total_seconds()
            if str(record.get("session_key") or "") == session_key and delta < session_cutoff:
                return False, "session_cooldown"
            if str(record.get("intent_id") or "") == intent_id and delta < intent_cutoff:
                return False, "intent_cooldown"
        return True, None

    def _record_skip(self, session_key: str, reason: str) -> None:
        self.ledger.append(ActiveIntentRecord(
            timestamp=_utcnow_iso(),
            session_key=session_key,
            intent_type="goal_nudge",
            intent_id=f"skip:{session_key}:{reason}",
            source_type="runtime",
            source_reference="eligibility",
            outcome="skipped",
            suppression_reason=reason,
        ))

    @staticmethod
    def _build_message(session: Session, candidate: ActiveIntentCandidate) -> InboundMessage:
        channel, chat_id = (
            session.key.split(":", 1)
            if ":" in session.key
            else ("cli", session.key)
        )
        return InboundMessage(
            channel="system",
            sender_id="agent_active",
            chat_id=f"{channel}:{chat_id}",
            content=candidate.content,
            session_key_override=session.key,
            metadata={
                "injected_event": "active_intent",
                "_from_active": True,
                "active_intent_type": candidate.intent_type,
                "active_intent_id": candidate.intent_id,
            },
        )
