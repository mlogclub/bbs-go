"""Feedback calibration for governed self-evolution signals."""

from __future__ import annotations

import json
import os
from contextlib import suppress
from dataclasses import asdict, dataclass
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

from filelock import FileLock

from OriginAgent.agent.evolution import (
    AUTO_EVOLUTION_ORIGIN,
    OpportunitySignalStore,
)
from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore, safe_append_outcome
from OriginAgent.utils.helpers import ensure_dir, truncate_text

FEEDBACK_STATE_RELATIVE = Path("memory") / "evolution_feedback_state.json"
NEGATIVE_EVENT_TYPES = {"review_rejected", "review_failed"}
POSITIVE_EVENT_TYPES = {"review_approved", "promoted"}


@dataclass(frozen=True)
class EvolutionFeedbackResult:
    """Summary for one calibration pass."""

    processed_events: int = 0
    feedback_applied: int = 0
    negative_feedback_applied: int = 0
    positive_feedback_applied: int = 0
    suppressed_signals: int = 0
    skipped_events: int = 0
    cooldown_skipped_events: int = 0
    last_event_id: str = ""
    last_calibrated_at: str | None = None

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class EvolutionFeedbackCalibrator:
    """Apply review, rollback, and promotion outcomes back onto signals."""

    def __init__(self, workspace: Path, config: Any | None = None) -> None:
        self.workspace = Path(workspace)
        self.config = config
        self.memory_dir = ensure_dir(self.workspace / "memory")
        self.state_path = self.workspace / FEEDBACK_STATE_RELATIVE
        self._lock = FileLock(str(self.memory_dir / ".evolution_feedback.lock"))
        self.signals = OpportunitySignalStore(self.workspace)
        self.outcomes = EvolutionOutcomeStore(self.workspace)

    def run(self) -> EvolutionFeedbackResult:
        """Apply unprocessed outcome feedback.

        This method is idempotent by event id. It is intentionally conservative:
        missing opportunity links are counted as skipped rather than inferred
        from weak artifact matches.
        """

        if not _config_bool(self.config, "feedback_calibration_enabled", True):
            return self._status_result()
        with self._lock:
            now_dt = datetime.now(timezone.utc)
            state = self._read_state_unlocked()
            processed_ids = {
                str(item)
                for item in state.get("processed_event_ids", [])
                if str(item)
            }
            cooldowns = _active_cooldowns(state.get("cooldowns"), now_dt)
            stats = {
                "processed_events": 0,
                "feedback_applied": 0,
                "negative_feedback_applied": 0,
                "positive_feedback_applied": 0,
                "suppressed_signals": 0,
                "skipped_events": 0,
                "cooldown_skipped_events": 0,
                "last_event_id": "",
            }
            for event in self.outcomes.read_all():
                event_id = str(event.get("event_id") or "")
                event_type = str(event.get("type") or "")
                if not event_id or event_id in processed_ids or event_type == "feedback_applied":
                    continue
                feedback = self._feedback_for_event(event)
                if feedback is None:
                    continue
                processed_ids.add(event_id)
                stats["processed_events"] += 1
                stats["last_event_id"] = event_id
                opportunity_id = feedback["opportunity_id"]
                if not opportunity_id:
                    stats["skipped_events"] += 1
                    continue
                cooldown = _cooldown_for(cooldowns, opportunity_id, now_dt)
                if feedback["polarity"] == "positive" and cooldown is not None:
                    stats["skipped_events"] += 1
                    stats["cooldown_skipped_events"] += 1
                    safe_append_outcome(
                        self.outcomes,
                        "feedback_applied",
                        opportunity_id=opportunity_id,
                        proposal_id=str(event.get("proposal_id") or ""),
                        artifact_type=str(event.get("artifact_type") or ""),
                        artifact_name=str(event.get("artifact_name") or ""),
                        artifact_path=str(event.get("artifact_path") or ""),
                        calibration_result={
                            "source_event_id": event_id,
                            "source_event_type": event_type,
                            "polarity": feedback["polarity"],
                            "skipped_by_cooldown": True,
                            "cooldown_until": str(cooldown.get("until") or ""),
                            "cooldown_reason": str(cooldown.get("reason") or ""),
                        },
                        metadata={
                            "reason": "Positive feedback skipped during active cooldown.",
                            "origin": AUTO_EVOLUTION_ORIGIN,
                        },
                    )
                    continue
                updated = self.signals.apply_feedback(
                    opportunity_id,
                    multiplier=feedback["multiplier"],
                    delta=feedback["delta"],
                    suppress=feedback["suppress"],
                    suppression_reason=feedback["reason"],
                    suppress_after_negative_count=_config_int(
                        self.config,
                        "feedback_suppress_after_negative_count",
                        3,
                    ),
                    risk_level=feedback["risk_level"],
                    verification_status=feedback["verification_status"],
                )
                if updated is None:
                    stats["skipped_events"] += 1
                    continue
                stats["feedback_applied"] += 1
                if feedback["polarity"] == "negative":
                    stats["negative_feedback_applied"] += 1
                    cooldowns[opportunity_id] = {
                        "until": (now_dt + timedelta(
                            days=max(1, _config_int(self.config, "feedback_cooldown_days", 14))
                        )).isoformat(),
                        "reason": feedback["reason"],
                        "source_event_id": event_id,
                        "source_event_type": event_type,
                    }
                elif feedback["polarity"] == "positive":
                    stats["positive_feedback_applied"] += 1
                if updated.status == "suppressed":
                    stats["suppressed_signals"] += 1
                safe_append_outcome(
                    self.outcomes,
                    "feedback_applied",
                    opportunity_id=opportunity_id,
                    proposal_id=str(event.get("proposal_id") or ""),
                    artifact_type=str(event.get("artifact_type") or ""),
                    artifact_name=str(event.get("artifact_name") or ""),
                    artifact_path=str(event.get("artifact_path") or ""),
                    feedback_score=updated.priority_score,
                    calibration_result={
                        "source_event_id": event_id,
                        "source_event_type": event_type,
                        "polarity": feedback["polarity"],
                        "multiplier": feedback["multiplier"],
                        "delta": feedback["delta"],
                        "suppress": feedback["suppress"],
                        "priority_score": updated.priority_score,
                        "negative_count": updated.feedback_negative_count,
                        "positive_count": updated.feedback_positive_count,
                    },
                    metadata={"reason": feedback["reason"], "origin": AUTO_EVOLUTION_ORIGIN},
                )
            now = now_dt.isoformat()
            trends = self._feedback_trends(now_dt)
            state.update({
                "processed_event_ids": sorted(processed_ids),
                "cooldowns": cooldowns,
                "feedback_trends": trends,
                "last_calibrated_at": now,
                "last_result": {
                    **stats,
                    "last_calibrated_at": now,
                },
            })
            self._write_state_unlocked(state)
            return EvolutionFeedbackResult(
                last_calibrated_at=now,
                **stats,
            )

    def status(self) -> dict[str, Any]:
        with self._lock:
            state = self._read_state_unlocked()
        last = state.get("last_result") if isinstance(state.get("last_result"), dict) else {}
        now_dt = datetime.now(timezone.utc)
        cooldowns = _active_cooldowns(state.get("cooldowns"), now_dt)
        trends = state.get("feedback_trends") if isinstance(state.get("feedback_trends"), dict) else {}
        trend_counts = _feedback_trend_counts(trends)
        feedback_events = [
            event for event in self.outcomes.read_all()
            if str(event.get("type") or "") == "feedback_applied"
        ]
        polarity_counts: dict[str, int] = {}
        for event in feedback_events:
            result = event.get("calibration_result")
            polarity = str(result.get("polarity") or "") if isinstance(result, dict) else ""
            if polarity:
                polarity_counts[polarity] = polarity_counts.get(polarity, 0) + 1
        return {
            "enabled": _config_bool(self.config, "feedback_calibration_enabled", True),
            "processed_event_count": len([
                item for item in state.get("processed_event_ids", [])
                if str(item)
            ]),
            "feedback_event_count": len(feedback_events),
            "feedback_polarity_counts": polarity_counts,
            "cooldown_count": len(cooldowns),
            "next_cooldown_expires_at": _next_cooldown_expiry(cooldowns),
            "feedback_trend_window_days": max(
                1,
                _config_int(self.config, "feedback_trend_window_days", 14),
            ),
            "feedback_trend_counts": trend_counts,
            "feedback_trends": trends,
            "last_calibrated_at": state.get("last_calibrated_at"),
            "last_result": last or None,
        }

    def _status_result(self) -> EvolutionFeedbackResult:
        status = self.status()
        last = status.get("last_result") if isinstance(status.get("last_result"), dict) else {}
        return EvolutionFeedbackResult(
            processed_events=_safe_int(last.get("processed_events"), 0),
            feedback_applied=_safe_int(last.get("feedback_applied"), 0),
            negative_feedback_applied=_safe_int(last.get("negative_feedback_applied"), 0),
            positive_feedback_applied=_safe_int(last.get("positive_feedback_applied"), 0),
            suppressed_signals=_safe_int(last.get("suppressed_signals"), 0),
            skipped_events=_safe_int(last.get("skipped_events"), 0),
            cooldown_skipped_events=_safe_int(last.get("cooldown_skipped_events"), 0),
            last_event_id=str(last.get("last_event_id") or ""),
            last_calibrated_at=(
                str(status.get("last_calibrated_at"))
                if status.get("last_calibrated_at") is not None
                else None
            ),
        )

    def _feedback_for_event(self, event: dict[str, Any]) -> dict[str, Any] | None:
        event_type = str(event.get("type") or "")
        rollback_status = str(event.get("rollback_status") or "")
        opportunity_id = self._opportunity_id_for_event(event)
        if event_type in NEGATIVE_EVENT_TYPES:
            return {
                "opportunity_id": opportunity_id,
                "polarity": "negative",
                "multiplier": _config_float(self.config, "feedback_reject_multiplier", 0.8),
                "delta": 0.0,
                "suppress": False,
                "risk_level": "",
                "verification_status": "feedback_negative",
                "reason": f"{event_type} lowered signal confidence.",
            }
        if event_type == "rolled_back" and rollback_status == "succeeded":
            return {
                "opportunity_id": opportunity_id,
                "polarity": "negative",
                "multiplier": _config_float(self.config, "feedback_rollback_multiplier", 0.6),
                "delta": _config_float(self.config, "feedback_rollback_delta", -0.1),
                "suppress": True,
                "risk_level": "high",
                "verification_status": "rolled_back",
                "reason": "Rollback suppresses the source opportunity until manually reconsidered.",
            }
        if event_type in POSITIVE_EVENT_TYPES and event_type != "promoted":
            metadata = event.get("metadata") if isinstance(event.get("metadata"), dict) else {}
            reason = str(metadata.get("reason") or "").strip()
            if reason == "auto_evolution verified low-risk workflow proposal":
                return None
            return {
                "opportunity_id": opportunity_id,
                "polarity": "positive",
                "multiplier": _config_float(self.config, "feedback_positive_multiplier", 1.02),
                "delta": 0.0,
                "suppress": False,
                "risk_level": "",
                "verification_status": "",
                "reason": f"{event_type} slightly reinforced signal confidence.",
            }
        return None

    def _opportunity_id_for_event(self, event: dict[str, Any]) -> str:
        direct = str(event.get("opportunity_id") or "")
        if direct:
            return direct
        metadata = event.get("metadata") if isinstance(event.get("metadata"), dict) else {}
        snapshot_id = str(metadata.get("snapshot_id") or "")
        if snapshot_id:
            with suppress(Exception):
                from OriginAgent.agent.evolution_snapshots import EvolutionSnapshotStore

                for snapshot in EvolutionSnapshotStore(self.workspace).list_snapshots():
                    if str(snapshot.get("snapshot_id") or "") == snapshot_id:
                        return str(snapshot.get("opportunity_id") or "")
        proposal_id = str(event.get("proposal_id") or metadata.get("proposal_id") or "")
        if proposal_id:
            with suppress(Exception):
                from OriginAgent.agent.background_review import ReviewProposalStore

                record = ReviewProposalStore(self.workspace).get(proposal_id)
                payload = record.get("payload") if isinstance(record, dict) else {}
                evolution = payload.get("evolution") if isinstance(payload, dict) else {}
                if isinstance(evolution, dict):
                    return str(evolution.get("opportunity_id") or "")
        return ""

    def _feedback_trends(self, now: datetime) -> dict[str, Any]:
        window_days = max(1, _config_int(self.config, "feedback_trend_window_days", 14))
        cutoff = now - timedelta(days=window_days)
        trends: dict[str, dict[str, Any]] = {}
        for event in self.outcomes.read_all():
            if str(event.get("type") or "") != "feedback_applied":
                continue
            event_time = _parse_datetime(str(event.get("timestamp") or ""))
            if event_time is None or event_time < cutoff:
                continue
            opportunity_id = str(event.get("opportunity_id") or "")
            if not opportunity_id:
                continue
            result = event.get("calibration_result")
            if not isinstance(result, dict):
                continue
            polarity = str(result.get("polarity") or "")
            if polarity not in {"negative", "positive"}:
                continue
            trend = trends.setdefault(
                opportunity_id,
                {
                    "window_days": window_days,
                    "positive": 0,
                    "negative": 0,
                    "skipped_positive": 0,
                    "skipped_negative": 0,
                    "net": 0,
                    "last_polarity": "",
                    "last_feedback_at": None,
                    "updated_at": now.isoformat(),
                },
            )
            skipped = bool(result.get("skipped_by_cooldown"))
            if skipped:
                key = "skipped_positive" if polarity == "positive" else "skipped_negative"
                trend[key] += 1
            else:
                trend[polarity] += 1
            if trend["last_feedback_at"] is None or event_time.isoformat() >= str(trend["last_feedback_at"]):
                trend["last_polarity"] = polarity
                trend["last_feedback_at"] = event_time.isoformat()
        for trend in trends.values():
            trend["net"] = int(trend["positive"]) - int(trend["negative"])
        return trends

    def _read_state_unlocked(self) -> dict[str, Any]:
        with suppress(FileNotFoundError):
            try:
                raw = json.loads(self.state_path.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError):
                raw = None
            if isinstance(raw, dict):
                return raw
        return {
            "processed_event_ids": [],
            "cooldowns": {},
            "feedback_trends": {},
            "last_calibrated_at": None,
            "last_result": None,
        }

    def _write_state_unlocked(self, state: dict[str, Any]) -> None:
        tmp_path = self.state_path.with_suffix(self.state_path.suffix + ".tmp")
        try:
            self.state_path.parent.mkdir(parents=True, exist_ok=True)
            with tmp_path.open("w", encoding="utf-8") as handle:
                handle.write(json.dumps(_redact_json(state), ensure_ascii=False, sort_keys=True, indent=2) + "\n")
                handle.flush()
                os.fsync(handle.fileno())
            os.replace(tmp_path, self.state_path)
        except BaseException:
            tmp_path.unlink(missing_ok=True)
            raise


def feedback_status(workspace: Path, config: Any | None = None) -> dict[str, Any]:
    try:
        return EvolutionFeedbackCalibrator(workspace, config).status()
    except Exception:
        return {
            "enabled": _config_bool(config, "feedback_calibration_enabled", True),
            "processed_event_count": 0,
            "feedback_event_count": 0,
            "feedback_polarity_counts": {},
            "cooldown_count": 0,
            "next_cooldown_expires_at": None,
            "feedback_trend_window_days": max(
                1,
                _config_int(config, "feedback_trend_window_days", 14),
            ),
            "feedback_trend_counts": {},
            "feedback_trends": {},
            "last_calibrated_at": None,
            "last_result": None,
        }


def _redact_json(value: Any) -> Any:
    if isinstance(value, str):
        return truncate_text(value, 1000)
    if isinstance(value, list):
        return [_redact_json(item) for item in value]
    if isinstance(value, dict):
        return {str(key): _redact_json(item) for key, item in value.items()}
    return value


def _safe_int(value: Any, default: int) -> int:
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default


def _parse_datetime(value: str) -> datetime | None:
    if not value:
        return None
    with suppress(ValueError):
        parsed = datetime.fromisoformat(value)
        if parsed.tzinfo is None:
            return parsed.replace(tzinfo=timezone.utc)
        return parsed.astimezone(timezone.utc)
    return None


def _active_cooldowns(value: Any, now: datetime) -> dict[str, dict[str, Any]]:
    if not isinstance(value, dict):
        return {}
    active: dict[str, dict[str, Any]] = {}
    for opportunity_id, raw in value.items():
        if not str(opportunity_id):
            continue
        if not isinstance(raw, dict):
            continue
        until = _parse_datetime(str(raw.get("until") or ""))
        if until is None or until <= now:
            continue
        active[str(opportunity_id)] = {
            "until": until.isoformat(),
            "reason": truncate_text(str(raw.get("reason") or ""), 512),
            "source_event_id": str(raw.get("source_event_id") or ""),
            "source_event_type": str(raw.get("source_event_type") or ""),
        }
    return active


def _cooldown_for(
    cooldowns: dict[str, dict[str, Any]],
    opportunity_id: str,
    now: datetime,
) -> dict[str, Any] | None:
    cooldown = cooldowns.get(opportunity_id)
    if not cooldown:
        return None
    until = _parse_datetime(str(cooldown.get("until") or ""))
    if until is None or until <= now:
        return None
    return cooldown


def _next_cooldown_expiry(cooldowns: dict[str, dict[str, Any]]) -> str | None:
    expiries = [
        str(cooldown.get("until") or "")
        for cooldown in cooldowns.values()
        if str(cooldown.get("until") or "")
    ]
    return min(expiries) if expiries else None


def _feedback_trend_counts(trends: Any) -> dict[str, int]:
    if not isinstance(trends, dict):
        return {}
    if not trends:
        return {}
    counts = {
        "positive": 0,
        "negative": 0,
        "skipped_positive": 0,
        "skipped_negative": 0,
        "net": 0,
    }
    for trend in trends.values():
        if not isinstance(trend, dict):
            continue
        counts["positive"] += _safe_int(trend.get("positive"), 0)
        counts["negative"] += _safe_int(trend.get("negative"), 0)
        counts["skipped_positive"] += _safe_int(trend.get("skipped_positive"), 0)
        counts["skipped_negative"] += _safe_int(trend.get("skipped_negative"), 0)
    counts["net"] = counts["positive"] - counts["negative"]
    return {key: value for key, value in counts.items() if value != 0 or key == "net"}


def _config_bool(config: Any | None, attr: str, default: bool) -> bool:
    return bool(getattr(config, attr, default) if config is not None else default)


def _config_int(config: Any | None, attr: str, default: int) -> int:
    return _safe_int(getattr(config, attr, default) if config is not None else default, default)


def _config_float(config: Any | None, attr: str, default: float) -> float:
    try:
        return float(getattr(config, attr, default) if config is not None else default)
    except (TypeError, ValueError):
        return default
