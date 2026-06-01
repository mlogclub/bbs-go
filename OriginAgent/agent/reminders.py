"""Durable one-shot reminder records for proactive follow-up."""

from __future__ import annotations

import json
import os
import uuid
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Literal

from filelock import FileLock

from OriginAgent.utils.helpers import ensure_dir, truncate_text

ReminderStatus = Literal["pending", "due", "fired", "completed", "cancelled", "expired"]
_SUMMARY_MAX_CHARS = 240


def _utcnow() -> datetime:
    return datetime.now(timezone.utc)


def _utcnow_iso() -> str:
    return _utcnow().isoformat()


def _normalize_iso(value: str) -> str:
    parsed = datetime.fromisoformat(str(value).strip())
    if parsed.tzinfo is None:
        parsed = parsed.replace(tzinfo=timezone.utc)
    else:
        parsed = parsed.astimezone(timezone.utc)
    return parsed.isoformat()


def _summarize(value: Any, max_chars: int = _SUMMARY_MAX_CHARS) -> str:
    text = str(value or "").strip()
    if not text:
        return ""
    return truncate_text(text, max_chars).replace("\n... (truncated)", " ...")


@dataclass(frozen=True)
class ReminderRecord:
    reminder_id: str
    session_key: str
    channel: str
    chat_id: str
    content: str
    due_at: str
    status: ReminderStatus = "pending"
    created_at: str = field(default_factory=_utcnow_iso)
    updated_at: str = field(default_factory=_utcnow_iso)
    fired_at: str | None = None
    completed_at: str | None = None
    cancelled_at: str | None = None
    last_delivery_at: str | None = None
    delivery_count: int = 0
    source: str = "user_request"

    @classmethod
    def create(
        cls,
        *,
        session_key: str,
        channel: str,
        chat_id: str,
        content: str,
        due_at: str,
        source: str = "user_request",
        reminder_id: str | None = None,
    ) -> "ReminderRecord":
        now = _utcnow_iso()
        return cls(
            reminder_id=reminder_id or str(uuid.uuid4())[:8],
            session_key=str(session_key).strip(),
            channel=str(channel).strip(),
            chat_id=str(chat_id).strip(),
            content=_summarize(content),
            due_at=_normalize_iso(due_at),
            status="pending",
            created_at=now,
            updated_at=now,
            source=str(source or "user_request").strip() or "user_request",
        )

    @classmethod
    def from_dict(cls, raw: dict[str, Any]) -> "ReminderRecord":
        return cls(
            reminder_id=str(raw.get("reminder_id") or "").strip(),
            session_key=str(raw.get("session_key") or "").strip(),
            channel=str(raw.get("channel") or "").strip(),
            chat_id=str(raw.get("chat_id") or "").strip(),
            content=_summarize(raw.get("content")),
            due_at=_normalize_iso(str(raw.get("due_at") or "").strip()),
            status=str(raw.get("status") or "pending"),  # type: ignore[arg-type]
            created_at=str(raw.get("created_at") or _utcnow_iso()),
            updated_at=str(raw.get("updated_at") or _utcnow_iso()),
            fired_at=str(raw.get("fired_at")).strip() if raw.get("fired_at") else None,
            completed_at=str(raw.get("completed_at")).strip() if raw.get("completed_at") else None,
            cancelled_at=str(raw.get("cancelled_at")).strip() if raw.get("cancelled_at") else None,
            last_delivery_at=(
                str(raw.get("last_delivery_at")).strip() if raw.get("last_delivery_at") else None
            ),
            delivery_count=max(0, int(raw.get("delivery_count") or 0)),
            source=str(raw.get("source") or "user_request").strip() or "user_request",
        )

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


class ReminderStore:
    """Durable JSONL-backed store for one-shot reminders."""

    def __init__(self, workspace: Path):
        self.workspace = Path(workspace)
        self.root = ensure_dir(self.workspace / "memory" / "active_intents")
        self.path = self.root / "reminders.jsonl"
        self.lock = FileLock(str(self.root / ".lock"))

    def read_all(self) -> list[ReminderRecord]:
        with self.lock:
            return self.read_all_unlocked()

    def read_all_unlocked(self) -> list[ReminderRecord]:
        records: list[ReminderRecord] = []
        if not self.path.exists():
            return records
        with self.path.open("r", encoding="utf-8") as handle:
            for line in handle:
                line = line.strip()
                if not line:
                    continue
                try:
                    payload = json.loads(line)
                except json.JSONDecodeError:
                    continue
                if isinstance(payload, dict):
                    try:
                        records.append(ReminderRecord.from_dict(payload))
                    except Exception:
                        continue
        return records

    def upsert(self, record: ReminderRecord) -> ReminderRecord:
        with self.lock:
            records = self.read_all_unlocked()
            replaced = False
            for idx, existing in enumerate(records):
                if existing.reminder_id == record.reminder_id:
                    records[idx] = record
                    replaced = True
                    break
            if not replaced:
                records.append(record)
            self._write_all_unlocked(records)
        return record

    def get(self, reminder_id: str) -> ReminderRecord | None:
        with self.lock:
            for record in self.read_all_unlocked():
                if record.reminder_id == reminder_id:
                    return record
        return None

    def list_due(self, *, now: datetime | None = None) -> list[ReminderRecord]:
        current = now or _utcnow()
        due: list[ReminderRecord] = []
        for record in self.read_all():
            if record.status not in {"pending", "due"}:
                continue
            due_at = datetime.fromisoformat(record.due_at)
            if due_at <= current:
                due.append(record)
        return sorted(due, key=lambda item: item.due_at)

    def mark_fired(self, reminder_id: str, *, now: datetime | None = None) -> ReminderRecord | None:
        current = (now or _utcnow()).isoformat()
        with self.lock:
            records = self.read_all_unlocked()
            updated = None
            for idx, record in enumerate(records):
                if record.reminder_id != reminder_id:
                    continue
                updated = ReminderRecord(
                    **{
                        **record.to_dict(),
                        "status": "fired",
                        "fired_at": current,
                        "last_delivery_at": current,
                        "updated_at": current,
                        "delivery_count": record.delivery_count + 1,
                    }
                )
                records[idx] = updated
                break
            if updated is not None:
                self._write_all_unlocked(records)
            return updated

    def mark_completed(self, reminder_id: str, *, now: datetime | None = None) -> ReminderRecord | None:
        return self._transition(reminder_id, status="completed", field_name="completed_at", now=now)

    def mark_cancelled(self, reminder_id: str, *, now: datetime | None = None) -> ReminderRecord | None:
        return self._transition(reminder_id, status="cancelled", field_name="cancelled_at", now=now)

    def stats(self, *, now: datetime | None = None) -> dict[str, Any]:
        current = now or _utcnow()
        records = self.read_all()
        counts: dict[str, int] = {}
        due_count = 0
        last_fired_at: str | None = None
        for record in records:
            counts[record.status] = counts.get(record.status, 0) + 1
            if record.status in {"pending", "due"} and datetime.fromisoformat(record.due_at) <= current:
                due_count += 1
            if record.fired_at and (last_fired_at is None or record.fired_at > last_fired_at):
                last_fired_at = record.fired_at
        return {
            "reminder_total": len(records),
            "reminder_status_counts": counts,
            "reminder_due_count": due_count,
            "reminder_last_fired_at": last_fired_at,
        }

    def _transition(
        self,
        reminder_id: str,
        *,
        status: ReminderStatus,
        field_name: str,
        now: datetime | None = None,
    ) -> ReminderRecord | None:
        current = (now or _utcnow()).isoformat()
        with self.lock:
            records = self.read_all_unlocked()
            updated = None
            for idx, record in enumerate(records):
                if record.reminder_id != reminder_id:
                    continue
                payload = {**record.to_dict(), "status": status, "updated_at": current, field_name: current}
                updated = ReminderRecord(**payload)
                records[idx] = updated
                break
            if updated is not None:
                self._write_all_unlocked(records)
            return updated

    def _write_all_unlocked(self, records: list[ReminderRecord]) -> None:
        ensure_dir(self.path.parent)
        with self.path.open("w", encoding="utf-8") as handle:
            for record in records:
                handle.write(json.dumps(record.to_dict(), ensure_ascii=False, sort_keys=True))
                handle.write("\n")
            handle.flush()
            os.fsync(handle.fileno())
