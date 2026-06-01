"""Health score for governed self-evolution observability."""

from __future__ import annotations

from typing import Any


def evolution_health_score(
    *,
    outcome_stats: dict[str, Any],
    dependency_stats: dict[str, Any],
    feedback_stats: dict[str, Any],
    sandbox_counts: dict[str, int],
    trial_status: dict[str, Any],
) -> dict[str, Any]:
    """Return a compact 0-100 health summary from existing evolution telemetry."""

    score = 100
    reasons: list[str] = []

    sandbox_total = sum(max(0, _safe_int(value, 0)) for value in sandbox_counts.values())
    blocked_or_failed = _safe_int(sandbox_counts.get("blocked"), 0) + _safe_int(sandbox_counts.get("failed"), 0)
    if sandbox_total:
        success_rate = max(0.0, min(1.0, _safe_int(sandbox_counts.get("passed"), 0) / sandbox_total))
        if success_rate >= 0.8:
            reasons.append(f"+ sandbox pass rate {round(success_rate * 100)}%")
        else:
            penalty = min(20, round((1.0 - success_rate) * 20))
            score -= penalty
            reasons.append(f"- sandbox pass rate {round(success_rate * 100)}%")
    if blocked_or_failed:
        penalty = min(15, blocked_or_failed * 5)
        score -= penalty
        reasons.append(f"- sandbox blocked or failed proposals: {blocked_or_failed}")

    rollback_counts = outcome_stats.get("rollback_status_counts")
    rollback_succeeded = _safe_int(_mapping(rollback_counts).get("succeeded"), 0)
    rollback_blocked = _safe_int(_mapping(rollback_counts).get("blocked"), 0)
    if rollback_succeeded:
        score -= min(15, rollback_succeeded * 7)
        reasons.append(f"- successful rollbacks: {rollback_succeeded}")
    else:
        reasons.append("+ no successful rollbacks")
    if rollback_blocked:
        score -= min(10, rollback_blocked * 4)
        reasons.append(f"- rollback dependency blocks: {rollback_blocked}")

    dependency_conflicts = _safe_int(dependency_stats.get("rollback_blocked_artifacts"), 0)
    stale_refs = _safe_int(dependency_stats.get("stale_reference_count"), 0)
    if dependency_conflicts:
        score -= min(15, dependency_conflicts * 5)
        reasons.append(f"- dependency conflicts: {dependency_conflicts}")
    else:
        reasons.append("+ no dependency conflicts")
    if stale_refs:
        score -= min(10, stale_refs * 2)
        reasons.append(f"- stale dependencies: {stale_refs}")

    trend_counts = _mapping(feedback_stats.get("feedback_trend_counts"))
    negative = _safe_int(trend_counts.get("negative"), 0)
    positive = _safe_int(trend_counts.get("positive"), 0)
    skipped_positive = _safe_int(trend_counts.get("skipped_positive"), 0)
    net = positive - negative
    if negative:
        score -= min(20, negative * 5)
        reasons.append(f"- recent negative feedback: {negative}")
    if net > 0:
        score += min(5, net * 2)
        reasons.append(f"+ positive feedback trend: +{net}")
    elif net < 0:
        score -= min(10, abs(net) * 3)
        reasons.append(f"- feedback trend: {net}")
    if skipped_positive:
        reasons.append(f"- positive feedback in cooldown: {skipped_positive}")

    cooldown_count = _safe_int(feedback_stats.get("cooldown_count"), 0)
    if cooldown_count:
        score -= min(10, cooldown_count * 3)
        reasons.append(f"- signals in cooldown: {cooldown_count}")

    if not bool(trial_status.get("enabled", True)):
        score -= 10
        reasons.append("- trial mode disabled")
    if not bool(trial_status.get("isolated_workspace", True)):
        score -= 25
        reasons.append("- trial workspace is not isolated")
    if not bool(trial_status.get("read_only_tools_only", True)):
        score -= 25
        reasons.append("- trial allows non-read-only tools")
    if (
        bool(trial_status.get("enabled", True))
        and bool(trial_status.get("isolated_workspace", True))
        and bool(trial_status.get("read_only_tools_only", True))
    ):
        reasons.append("+ trial isolation enforced")

    score = max(0, min(100, score))
    return {
        "score": score,
        "level": _level(score),
        "reasons": reasons[:8],
    }


def default_evolution_health() -> dict[str, Any]:
    return {
        "score": 100,
        "level": "healthy",
        "reasons": ["+ no successful rollbacks", "+ no dependency conflicts", "+ trial isolation enforced"],
    }


def _level(score: int) -> str:
    if score >= 80:
        return "healthy"
    if score >= 60:
        return "degraded"
    return "unhealthy"


def _mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def _safe_int(value: Any, default: int) -> int:
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default
