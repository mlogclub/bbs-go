"""Maintenance helpers for governed self-evolution stores."""

from __future__ import annotations

from pathlib import Path
from typing import Any

from loguru import logger

from OriginAgent.agent.evolution_dependencies import EvolutionDependencyStore
from OriginAgent.agent.evolution_config_overlay import apply_config_overlay, self_tune_evolution_config
from OriginAgent.agent.evolution_feedback import feedback_status
from OriginAgent.agent.evolution_health import evolution_health_score
from OriginAgent.agent.evolution_health_history import EvolutionHealthHistoryStore
from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore
from OriginAgent.agent.evolution_sandbox import sandbox_status_counts, trial_policy_status
from OriginAgent.agent.evolution_schema import validate_evolution_stores
from OriginAgent.agent.evolution_trial_logs import EvolutionTrialLogStore


def run_evolution_maintenance(
    workspace: Path,
    config: Any | None = None,
    *,
    force_cleanup: bool = False,
) -> dict[str, Any]:
    """Run bounded cleanup for append-only evolution stores."""

    workspace = Path(workspace)
    config = apply_config_overlay(workspace, config)
    outcome_retention_days = max(
        1,
        int(getattr(config, "outcome_retention_days", 90) if config is not None else 90),
    )
    outcome_archive_enabled = bool(
        getattr(config, "outcome_archive_enabled", True) if config is not None else True
    )
    dependency_cleanup_enabled = bool(
        getattr(config, "dependency_stale_cleanup_enabled", True) if config is not None else True
    )
    health_history_retention_days = max(
        1,
        int(getattr(config, "health_history_retention_days", 90) if config is not None else 90),
    )
    max_health_history_snapshots = max(
        0,
        int(getattr(config, "max_health_history_snapshots", 100) if config is not None else 100),
    )
    trial_config = getattr(config, "trial", None) if config is not None else None
    trial_log_retention_days = max(
        1,
        int(getattr(trial_config, "trial_log_retention_days", 30) if trial_config is not None else 30),
    )
    max_retained_trial_logs = max(
        0,
        int(getattr(trial_config, "max_retained_trial_logs", 10) if trial_config is not None else 10),
    )
    outcomes = EvolutionOutcomeStore(workspace)
    dependencies = EvolutionDependencyStore(workspace)
    trial_logs = EvolutionTrialLogStore(workspace)
    health_history = EvolutionHealthHistoryStore(workspace)
    maintenance: dict[str, Any] = {
        "outcome_retention_days": outcome_retention_days,
        "outcome_archive_enabled": outcome_archive_enabled,
        "dependency_stale_cleanup_enabled": dependency_cleanup_enabled,
        "health_history_retention_days": health_history_retention_days,
        "max_health_history_snapshots": max_health_history_snapshots,
        "trial_log_retention_days": trial_log_retention_days,
        "max_retained_trial_logs": max_retained_trial_logs,
    }
    try:
        maintenance["outcome_retention"] = outcomes.enforce_retention(
            retention_days=outcome_retention_days,
            archive=outcome_archive_enabled,
        )
    except Exception:
        logger.exception("Evolution outcome retention failed")
    if force_cleanup or dependency_cleanup_enabled:
        try:
            maintenance["dependency_cleanup"] = dependencies.prune_stale_references()
        except Exception:
            logger.exception("Evolution dependency cleanup failed")
    try:
        maintenance["trial_log_retention"] = trial_logs.enforce_retention(
            max_records=max_retained_trial_logs,
            retention_days=trial_log_retention_days,
        )
    except Exception:
        logger.exception("Evolution trial log retention failed")
    try:
        outcome_stats = outcomes.stats()
        dependency_stats = dependencies.stats()
        feedback_stats = feedback_status(workspace, config)
        sandbox_counts = sandbox_status_counts(workspace)
        health = evolution_health_score(
            outcome_stats=outcome_stats,
            dependency_stats=dependency_stats,
            feedback_stats=feedback_stats,
            sandbox_counts=sandbox_counts,
            trial_status=trial_policy_status(config),
        )
        snapshot = health_history.append_snapshot(
            health,
            metadata={"source": "evolution_maintenance"},
        )
        maintenance["health_history_snapshot"] = {
            "snapshot_id": snapshot.get("snapshot_id"),
            "score": snapshot.get("score"),
            "level": snapshot.get("level"),
            "timestamp": snapshot.get("timestamp"),
        }
        maintenance["health_history_retention"] = health_history.enforce_retention(
            max_records=max_health_history_snapshots,
            retention_days=health_history_retention_days,
        )
        maintenance["config_self_tuning"] = self_tune_evolution_config(
            workspace,
            config,
            health=health,
            outcome_stats=outcome_stats,
            sandbox_counts=sandbox_counts,
            feedback_stats=feedback_stats,
        )
    except Exception:
        logger.exception("Evolution health history maintenance failed")
    try:
        validation = validate_evolution_stores(workspace)
        maintenance["schema_validation"] = {
            "ok": validation.get("ok"),
            "record_counts": validation.get("record_counts", {}),
            "issue_counts": validation.get("issue_counts", {}),
        }
    except Exception:
        logger.exception("Evolution schema validation failed")
    return maintenance
