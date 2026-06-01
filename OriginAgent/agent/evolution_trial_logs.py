"""Bounded trial execution logs for governed evolution."""

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

TRIAL_LOG_SCHEMA_VERSION = "originagent.evolution.trial_log.v1"
TRIAL_LOG_STORE_RELATIVE = Path("memory") / "evolution_trial_logs.jsonl"
_SUMMARY_CHARS = 300


class EvolutionTrialLogStore:
    """Append compact trial logs while bounding detailed output retention."""

    def __init__(self, workspace: Path, *, path: Path | None = None) -> None:
        self.workspace = Path(workspace)
        memory_dir = ensure_dir(self.workspace / "memory")
        self.path = path or (self.workspace / TRIAL_LOG_STORE_RELATIVE)
        self._lock = FileLock(str(memory_dir / ".evolution_trial_logs.lock"))

    def append_log(
        self,
        *,
        opportunity_id: str = "",
        proposal_id: str = "",
        artifact_type: str = "",
        artifact_name: str = "",
        artifact_path: str = "",
        status: str = "unknown",
        step_logs: list[dict[str, Any]] | None = None,
        summary: dict[str, Any] | None = None,
        metadata: dict[str, Any] | None = None,
        max_step_output_chars: int = 2000,
        timestamp: datetime | None = None,
    ) -> dict[str, Any]:
        record = {
            "schema_version": TRIAL_LOG_SCHEMA_VERSION,
            "trial_id": f"trial_{uuid.uuid4().hex}",
            "timestamp": _timestamp(timestamp),
            "opportunity_id": _clean(opportunity_id, 256),
            "proposal_id": _clean(proposal_id, 256),
            "artifact_type": _clean(artifact_type, 64),
            "artifact_name": _clean(artifact_name, 256),
            "artifact_path": _clean(_normalize_path(artifact_path), 512),
            "status": _clean(status, 128) or "unknown",
            "summary": _redact_json(summary or {}),
            "step_logs": [
                _bounded_step_log(item, max_step_output_chars=max_step_output_chars)
                for item in (step_logs or [])
                if isinstance(item, dict)
            ],
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
        max_records: int = 10,
        retention_days: int = 30,
        now: datetime | None = None,
    ) -> dict[str, Any]:
        """Keep only recent trial logs and cap detailed log count."""

        now_dt = _timestamp_dt(now)
        cutoff = now_dt - timedelta(days=max(1, int(retention_days or 30)))
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
            "retention_days": max(1, int(retention_days or 30)),
            "max_retained_trial_logs": max_count,
            "retained_count": len(retained),
            "removed_count": removed_count,
            "cutoff": cutoff.isoformat(),
        }

    def stats(self) -> dict[str, Any]:
        records = self.read_all()
        status_counts: dict[str, int] = {}
        last_trial_at = None
        truncated_steps = 0
        for record in records:
            status = str(record.get("status") or "").strip()
            if status:
                status_counts[status] = status_counts.get(status, 0) + 1
            timestamp = record.get("timestamp")
            if isinstance(timestamp, str) and (last_trial_at is None or timestamp > last_trial_at):
                last_trial_at = timestamp
            steps = record.get("step_logs") if isinstance(record.get("step_logs"), list) else []
            truncated_steps += len([
                step for step in steps
                if isinstance(step, dict) and bool(step.get("output_truncated"))
            ])
        return {
            "trial_log_count": len(records),
            "trial_log_status_counts": status_counts,
            "last_trial_at": last_trial_at,
            "truncated_step_output_count": truncated_steps,
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


def trial_log_policy_status(config: Any | None = None) -> dict[str, Any]:
    trial_config = getattr(config, "trial", None)
    return {
        "max_step_output_chars": _config_int(trial_config, "max_step_output_chars", 2000),
        "max_retained_trial_logs": _config_int(trial_config, "max_retained_trial_logs", 10),
        "trial_log_retention_days": _config_int(trial_config, "trial_log_retention_days", 30),
    }


def _bounded_step_log(step: dict[str, Any], *, max_step_output_chars: int) -> dict[str, Any]:
    raw_output = _redact_text(str(step.get("output") or ""))
    max_chars = max(0, int(max_step_output_chars or 0))
    output = raw_output[:max_chars] if max_chars else ""
    truncated = bool(max_chars and len(raw_output) > max_chars)
    return {
        "index": _config_int(step, "index", 0),
        "title": _clean(step.get("title"), 160),
        "tool": _clean(step.get("tool"), 128),
        "status": _clean(step.get("status"), 128) or "unknown",
        "output": output,
        "output_truncated": truncated,
        "output_chars": len(raw_output),
        "output_summary": truncate_text(raw_output, _SUMMARY_CHARS) if raw_output else "",
        "issues": _redact_json(step.get("issues") if isinstance(step.get("issues"), list) else []),
        "metadata": _redact_json(step.get("metadata") if isinstance(step.get("metadata"), dict) else {}),
    }


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


def _normalize_path(path: str) -> str:
    return str(path or "").replace("\\", "/")


def _redact_json(value: Any) -> Any:
    if isinstance(value, str):
        return truncate_text(_redact_text(value), 1000)
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
