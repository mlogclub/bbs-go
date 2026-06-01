"""Bounded health history for governed evolution monitoring."""

from __future__ import annotations

import json
import os
import uuid
from contextlib import suppress
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

from filelock import FileLock

from OriginAgent.utils.helpers import ensure_dir, truncate_text

HEALTH_HISTORY_SCHEMA_VERSION = "originagent.evolution.health_history.v1"
HEALTH_HISTORY_STORE_RELATIVE = Path("memory") / "evolution_health_history.jsonl"


class EvolutionHealthHistoryStore:
    """Append compact evolution health snapshots with bounded retention."""

    def __init__(self, workspace: Path, *, path: Path | None = None) -> None:
        self.workspace = Path(workspace)
        memory_dir = ensure_dir(self.workspace / "memory")
        self.path = path or (self.workspace / HEALTH_HISTORY_STORE_RELATIVE)
        self._lock = FileLock(str(memory_dir / ".evolution_health_history.lock"))

    def append_snapshot(
        self,
        health: dict[str, Any],
        *,
        metadata: dict[str, Any] | None = None,
        timestamp: datetime | None = None,
    ) -> dict[str, Any]:
        record = {
            "schema_version": HEALTH_HISTORY_SCHEMA_VERSION,
            "snapshot_id": f"evo_health_{uuid.uuid4().hex}",
            "timestamp": _timestamp(timestamp),
            "score": _score(health.get("score")),
            "level": _clean(health.get("level"), 64) or "unknown",
            "reasons": [
                _clean(reason, 300)
                for reason in (health.get("reasons") if isinstance(health.get("reasons"), list) else [])
            ][:8],
            "metadata": _redact_json(metadata or {}),
        }
        with self._lock:
            self.path.parent.mkdir(parents=True, exist_ok=True)
            with self.path.open("a", encoding="utf-8") as handle:
                handle.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")
        return record

    def read_all(self) -> list[dict[str, Any]]:
        records: list[dict[str, Any]] = []
        with suppress(FileNotFoundError):
            with self.path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        raw = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(raw, dict):
                        records.append(raw)
        return records

    def enforce_retention(
        self,
        *,
        max_records: int = 100,
        retention_days: int = 90,
        now: datetime | None = None,
    ) -> dict[str, Any]:
        now_dt = _timestamp_dt(now)
        cutoff = now_dt - timedelta(days=max(1, int(retention_days or 90)))
        max_count = max(0, int(max_records or 0))
        records = self.read_all()
        recent = [
            record for record in records
            if (_parse_datetime(str(record.get("timestamp") or "")) or now_dt) >= cutoff
        ]
        recent.sort(key=lambda record: str(record.get("timestamp") or ""))
        retained = recent[-max_count:] if max_count else []
        removed_count = len(records) - len(retained)
        if removed_count:
            self.compact_copy(retained)
        return {
            "retention_days": max(1, int(retention_days or 90)),
            "max_health_history_snapshots": max_count,
            "retained_count": len(retained),
            "removed_count": removed_count,
            "cutoff": cutoff.isoformat(),
        }

    def summary(self) -> dict[str, Any]:
        records = sorted(self.read_all(), key=lambda record: str(record.get("timestamp") or ""))
        if not records:
            return {
                "snapshot_count": 0,
                "latest_score": None,
                "latest_level": None,
                "previous_score": None,
                "score_delta": 0,
                "trend": "unknown",
                "last_snapshot_at": None,
            }
        latest = records[-1]
        previous = records[-2] if len(records) > 1 else None
        latest_score = _score(latest.get("score"))
        previous_score = _score(previous.get("score")) if previous is not None else None
        delta = latest_score - previous_score if previous_score is not None else 0
        return {
            "snapshot_count": len(records),
            "latest_score": latest_score,
            "latest_level": _clean(latest.get("level"), 64) or "unknown",
            "previous_score": previous_score,
            "score_delta": delta,
            "trend": _trend(delta, has_previous=previous_score is not None),
            "last_snapshot_at": latest.get("timestamp") if isinstance(latest.get("timestamp"), str) else None,
        }

    def compact_copy(self, records: list[dict[str, Any]]) -> None:
        tmp_path = self.path.with_suffix(self.path.suffix + ".tmp")
        with self._lock:
            self.path.parent.mkdir(parents=True, exist_ok=True)
            try:
                with tmp_path.open("w", encoding="utf-8") as handle:
                    for record in records:
                        handle.write(json.dumps(_redact_json(record), ensure_ascii=False, sort_keys=True) + "\n")
                    handle.flush()
                    os.fsync(handle.fileno())
                os.replace(tmp_path, self.path)
            except BaseException:
                tmp_path.unlink(missing_ok=True)
                raise


def health_history_policy_status(config: Any | None = None) -> dict[str, Any]:
    return {
        "health_history_retention_days": _config_int(config, "health_history_retention_days", 90),
        "max_health_history_snapshots": _config_int(config, "max_health_history_snapshots", 100),
    }


def _trend(delta: int, *, has_previous: bool) -> str:
    if not has_previous:
        return "unknown"
    if delta >= 3:
        return "improving"
    if delta <= -3:
        return "degrading"
    return "stable"


def _timestamp(value: datetime | None) -> str:
    return _timestamp_dt(value).isoformat()


def _timestamp_dt(value: datetime | None) -> datetime:
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
        if parsed.tzinfo is None:
            return parsed.replace(tzinfo=timezone.utc)
        return parsed.astimezone(timezone.utc)
    return None


def _score(value: Any) -> int:
    if isinstance(value, bool):
        return 0
    try:
        return max(0, min(100, int(value)))
    except (TypeError, ValueError):
        return 0


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


def _clean(value: Any, max_chars: int) -> str:
    if value is None:
        return ""
    return truncate_text(_redact_text(str(value).strip()), max_chars)


def _redact_json(value: Any) -> Any:
    if isinstance(value, str):
        return _clean(value, 1000)
    if isinstance(value, list):
        return [_redact_json(item) for item in value]
    if isinstance(value, dict):
        return {str(key): _redact_json(item) for key, item in value.items()}
    return value


def _redact_text(text: str) -> str:
    try:
        from OriginAgent.agent.memory import redact_memory_text

        return redact_memory_text(text)
    except Exception:
        return text
