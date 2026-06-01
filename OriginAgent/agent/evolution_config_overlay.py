"""Governed self-tuning overlay for evolution config."""

from __future__ import annotations

import json
from copy import deepcopy
from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from filelock import FileLock

from OriginAgent.agent.evolution_outcomes import EvolutionOutcomeStore, safe_append_outcome
from OriginAgent.utils.helpers import ensure_dir

OVERLAY_SCHEMA_VERSION = "originagent.evolution.config_overlay.v1"
PATCH_SCHEMA_VERSION = "originagent.evolution.config_patch.v1"
OVERLAY_RELATIVE = Path("memory") / "evolution_config_overrides.json"
PATCH_LOG_RELATIVE = Path("memory") / "evolution_config_patches.jsonl"
CONFIG_PATCH_APPLIED = "config_patch_applied"
CONFIG_PATCH_REJECTED = "config_patch_rejected"

_MODE_RANK = {"conservative": 0, "curated": 1, "exploratory": 2, "aggressive": 3}
_MAX_THRESHOLD_STEP = 0.08
_SAFE_TOOL_SET = {"read_file", "glob", "grep"}
_BLOCKED_TRIAL_TOOLS = {"write_file", "edit_file", "exec", "message", "cron", "spawn"}


@dataclass(frozen=True)
class ConfigPatch:
    """One proposed evolution config override patch."""

    path: str
    value: Any
    reason: str
    old_value: Any = None

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class ConfigMutationResult:
    """Result of applying governed self-tuning patches."""

    ok: bool
    applied: list[dict[str, Any]] = field(default_factory=list)
    rejected: list[dict[str, Any]] = field(default_factory=list)
    overlay: dict[str, Any] = field(default_factory=dict)
    message: str = ""

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


class EvolutionConfigOverlayStore:
    """Read and write the self-tuning overlay without touching main config."""

    def __init__(self, workspace: Path) -> None:
        self.workspace = Path(workspace)
        self.memory_dir = ensure_dir(self.workspace / "memory")
        self.path = self.workspace / OVERLAY_RELATIVE
        self.patch_log_path = self.workspace / PATCH_LOG_RELATIVE
        self._lock = FileLock(str(self.memory_dir / ".evolution_config_overrides.lock"))

    def read_overlay(self) -> dict[str, Any]:
        try:
            with self.path.open("r", encoding="utf-8") as handle:
                raw = json.load(handle)
        except (FileNotFoundError, json.JSONDecodeError, OSError):
            return _empty_overlay()
        if not isinstance(raw, dict):
            return _empty_overlay()
        raw.setdefault("schema_version", OVERLAY_SCHEMA_VERSION)
        raw.setdefault("overrides", {})
        raw.setdefault("patch_count", 0)
        raw.setdefault("active", bool(raw.get("overrides")))
        if not isinstance(raw.get("overrides"), dict):
            raw["overrides"] = {}
        return raw

    def status(self) -> dict[str, Any]:
        overlay = self.read_overlay()
        overrides = overlay.get("overrides") if isinstance(overlay.get("overrides"), dict) else {}
        patches = self.read_patch_log(limit=20)
        return {
            "schema_version": OVERLAY_SCHEMA_VERSION,
            "active": bool(overrides),
            "override_count": len(_flatten(overrides)),
            "overrides": overrides,
            "patch_count": int(overlay.get("patch_count") or 0),
            "last_patch_at": str(overlay.get("last_patch_at") or ""),
            "last_actor": str(overlay.get("last_actor") or ""),
            "recent_patches": patches[-5:],
            "store_path": str(OVERLAY_RELATIVE).replace("\\", "/"),
            "patch_log_path": str(PATCH_LOG_RELATIVE).replace("\\", "/"),
        }

    def effective_config(self, config: Any | None) -> Any | None:
        """Return a deep-copied config with overlay values applied."""

        if config is None:
            return None
        clone = deepcopy(config)
        overrides = self.read_overlay().get("overrides")
        if not isinstance(overrides, dict):
            return clone
        for path, value in _flatten(overrides).items():
            _set_path(clone, path, value)
        return clone

    def apply_patches(
        self,
        config: Any | None,
        patches: list[ConfigPatch],
        *,
        actor: str = "auto_evolution",
        source: str = "self_tuning",
        evidence: dict[str, Any] | None = None,
    ) -> ConfigMutationResult:
        overlay = self.read_overlay()
        overrides = overlay.get("overrides") if isinstance(overlay.get("overrides"), dict) else {}
        rejected: list[dict[str, Any]] = []
        applied: list[dict[str, Any]] = []
        now = datetime.now(timezone.utc).isoformat()
        next_overrides = deepcopy(overrides)
        effective_before = self.effective_config(config)

        for patch in patches:
            old_value = _get_path(effective_before, patch.path)
            normalized = ConfigPatch(
                path=patch.path,
                value=patch.value,
                reason=patch.reason,
                old_value=old_value,
            )
            issue = ConfigMutationGate.validate(effective_before, normalized)
            record = {
                "schema_version": PATCH_SCHEMA_VERSION,
                "timestamp": now,
                "actor": actor,
                "source": source,
                "patch": normalized.to_json(),
                "evidence": evidence or {},
            }
            if issue:
                record["status"] = "rejected"
                record["error"] = issue
                rejected.append(record)
                self._append_patch_log(record)
                safe_append_outcome(
                    EvolutionOutcomeStore(self.workspace),
                    CONFIG_PATCH_REJECTED,
                    metadata={
                        "actor": actor,
                        "source": source,
                        "path": patch.path,
                        "error": issue,
                    },
                )
                continue
            _set_path_in_mapping(next_overrides, patch.path, patch.value)
            _set_path(effective_before, patch.path, patch.value)
            record["status"] = "applied"
            applied.append(record)
            self._append_patch_log(record)

        if applied:
            next_overlay = {
                "schema_version": OVERLAY_SCHEMA_VERSION,
                "active": True,
                "updated_at": now,
                "last_patch_at": now,
                "last_actor": actor,
                "patch_count": int(overlay.get("patch_count") or 0) + len(applied),
                "overrides": next_overrides,
            }
            self._write_overlay(next_overlay)
            for record in applied:
                safe_append_outcome(
                    EvolutionOutcomeStore(self.workspace),
                    CONFIG_PATCH_APPLIED,
                    metadata={
                        "actor": actor,
                        "source": source,
                        "path": record["patch"]["path"],
                        "old_value": record["patch"].get("old_value"),
                        "new_value": record["patch"].get("value"),
                    },
                )
            return ConfigMutationResult(
                ok=True,
                applied=applied,
                rejected=rejected,
                overlay=next_overlay,
                message=f"Applied {len(applied)} governed config override patch(es).",
            )
        return ConfigMutationResult(
            ok=not rejected,
            applied=[],
            rejected=rejected,
            overlay=overlay,
            message="No governed config override patches were applied.",
        )

    def clear(self, *, actor: str = "auto_evolution", source: str = "self_tuning") -> dict[str, Any]:
        now = datetime.now(timezone.utc).isoformat()
        overlay = {
            "schema_version": OVERLAY_SCHEMA_VERSION,
            "active": False,
            "updated_at": now,
            "last_patch_at": str(self.read_overlay().get("last_patch_at") or ""),
            "last_actor": actor,
            "patch_count": int(self.read_overlay().get("patch_count") or 0),
            "overrides": {},
        }
        self._write_overlay(overlay)
        record = {
            "schema_version": PATCH_SCHEMA_VERSION,
            "timestamp": now,
            "actor": actor,
            "source": source,
            "status": "cleared",
            "patch": {"path": "*", "old_value": None, "value": None, "reason": "Clear governed config overlay."},
            "evidence": {},
        }
        self._append_patch_log(record)
        return overlay

    def read_patch_log(self, *, limit: int = 50) -> list[dict[str, Any]]:
        records: list[dict[str, Any]] = []
        try:
            with self.patch_log_path.open("r", encoding="utf-8") as handle:
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
        except FileNotFoundError:
            return []
        effective_limit = max(1, min(int(limit or 50), 200))
        return records[-effective_limit:]

    def _write_overlay(self, overlay: dict[str, Any]) -> None:
        with self._lock:
            self.path.parent.mkdir(parents=True, exist_ok=True)
            tmp_path = self.path.with_suffix(self.path.suffix + ".tmp")
            with tmp_path.open("w", encoding="utf-8") as handle:
                json.dump(overlay, handle, ensure_ascii=False, indent=2, sort_keys=True)
                handle.write("\n")
            tmp_path.replace(self.path)

    def _append_patch_log(self, record: dict[str, Any]) -> None:
        with self._lock:
            self.patch_log_path.parent.mkdir(parents=True, exist_ok=True)
            with self.patch_log_path.open("a", encoding="utf-8") as handle:
                handle.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")


class ConfigMutationGate:
    """Allow only safety-tightening evolution config changes."""

    @staticmethod
    def validate(config: Any | None, patch: ConfigPatch) -> str:
        path = str(patch.path or "").strip()
        value = patch.value
        old = patch.old_value if patch.old_value is not None else _get_path(config, path)
        if path == "allow_manual_override":
            return "allow_manual_override can never be changed by self-tuning."
        if path == "dry_run":
            return "" if value is True else "self-tuning may only set dry_run=true."
        if path == "mode":
            old_rank = _MODE_RANK.get(str(old or "conservative"), 0)
            new_rank = _MODE_RANK.get(str(value or ""), 99)
            if new_rank <= old_rank:
                return ""
            return "self-tuning may only move mode to an equal or more conservative value."
        if path == "workflow_priority_threshold":
            try:
                old_float = float(old if old is not None else 0.7)
                new_float = float(value)
            except (TypeError, ValueError):
                return "workflow_priority_threshold must be numeric."
            if not 0.0 <= new_float <= 1.0:
                return "workflow_priority_threshold must stay between 0 and 1."
            if new_float < old_float:
                return "self-tuning may only raise workflow_priority_threshold."
            if new_float - old_float > _MAX_THRESHOLD_STEP:
                return "workflow_priority_threshold change is too large for one patch."
            return ""
        if path == "max_proposals_per_cycle":
            try:
                old_int = int(old if old is not None else 3)
                new_int = int(value)
            except (TypeError, ValueError):
                return "max_proposals_per_cycle must be an integer."
            if not 0 <= new_int <= 20:
                return "max_proposals_per_cycle must stay between 0 and 20."
            if new_int <= old_int:
                return ""
            return "self-tuning may only reduce max_proposals_per_cycle."
        if path == "skill_candidates_enabled":
            return "" if value is False else "self-tuning may only disable skill candidates."
        if path == "auto_verify_workflows":
            return "" if value is False else "self-tuning may only disable auto verification."
        if path == "trial.enabled":
            return "" if value is True else "self-tuning may not disable trial evaluation."
        if path == "trial.isolated_workspace":
            return "" if value is True else "self-tuning may only require trial isolation."
        if path == "trial.read_only_tools_only":
            return "" if value is True else "self-tuning may only require read-only trial tools."
        if path == "trial.blocked_tools":
            new_tools = {str(item).strip() for item in value or [] if str(item).strip()}
            old_tools = {str(item).strip() for item in old or [] if str(item).strip()}
            if not _BLOCKED_TRIAL_TOOLS.issubset(new_tools):
                return "trial.blocked_tools must include all side-effecting tools."
            if old_tools and not old_tools.issubset(new_tools):
                return "self-tuning may only add blocked trial tools."
            return ""
        if path == "sandbox.read_only_tools":
            new_tools = {str(item).strip() for item in value or [] if str(item).strip()}
            old_tools = {str(item).strip() for item in old or [] if str(item).strip()}
            if not new_tools:
                return "sandbox.read_only_tools cannot be empty."
            if not new_tools.issubset(_SAFE_TOOL_SET):
                return "sandbox.read_only_tools may only contain built-in read-only tools."
            if old_tools and not new_tools.issubset(old_tools):
                return "self-tuning may only reduce sandbox read-only tools."
            return ""
        return f"self-tuning cannot modify `{path}`."


class EvolutionConfigTuner:
    """Derive safety-tightening config patches from evolution telemetry."""

    def __init__(self, workspace: Path, config: Any | None = None) -> None:
        self.workspace = Path(workspace)
        self.config = config

    def propose(
        self,
        *,
        health: dict[str, Any],
        outcome_stats: dict[str, Any],
        sandbox_counts: dict[str, int],
        feedback_stats: dict[str, Any],
    ) -> tuple[list[ConfigPatch], dict[str, Any]]:
        score = _safe_int(health.get("score"), 100)
        mode = str(getattr(self.config, "mode", "conservative") if self.config is not None else "conservative")
        dry_run = bool(getattr(self.config, "dry_run", True) if self.config is not None else True)
        threshold = _safe_float(getattr(self.config, "workflow_priority_threshold", 0.7), 0.7)
        max_proposals = _safe_int(getattr(self.config, "max_proposals_per_cycle", 3), 3)
        rollback_succeeded = _safe_int(_mapping(outcome_stats.get("rollback_status_counts")).get("succeeded"), 0)
        blocked_or_failed = _safe_int(sandbox_counts.get("blocked"), 0) + _safe_int(sandbox_counts.get("failed"), 0)
        negative_feedback = _safe_int(_mapping(feedback_stats.get("feedback_trend_counts")).get("negative"), 0)
        pressure = rollback_succeeded + blocked_or_failed + negative_feedback
        patches: list[ConfigPatch] = []
        evidence = {
            "health_score": score,
            "health_level": str(health.get("level") or ""),
            "rollback_succeeded": rollback_succeeded,
            "sandbox_blocked_or_failed": blocked_or_failed,
            "negative_feedback": negative_feedback,
        }
        if score < 80 or pressure > 0:
            if not dry_run:
                patches.append(ConfigPatch("dry_run", True, "Health or feedback pressure requires preview mode."))
            if mode != "conservative":
                patches.append(ConfigPatch("mode", "conservative", "Health or feedback pressure requires conservative mode."))
            if threshold < 0.95:
                patches.append(ConfigPatch(
                    "workflow_priority_threshold",
                    round(min(0.95, threshold + 0.05), 2),
                    "Raise workflow threshold after failed, blocked, rollback, or negative feedback signals.",
                ))
            if max_proposals > 1:
                patches.append(ConfigPatch(
                    "max_proposals_per_cycle",
                    max(1, max_proposals - 1),
                    "Reduce proposal rate while evolution health is under pressure.",
                ))
        if score < 60:
            patches.extend([
                ConfigPatch("skill_candidates_enabled", False, "Unhealthy evolution disables skill candidates."),
                ConfigPatch("auto_verify_workflows", False, "Unhealthy evolution disables automatic verification."),
                ConfigPatch("trial.isolated_workspace", True, "Unhealthy evolution reasserts trial isolation."),
                ConfigPatch("trial.read_only_tools_only", True, "Unhealthy evolution reasserts read-only trials."),
                ConfigPatch(
                    "trial.blocked_tools",
                    sorted(_BLOCKED_TRIAL_TOOLS),
                    "Unhealthy evolution reasserts side-effecting trial blocklist.",
                ),
            ])
        return patches, evidence


def apply_config_overlay(workspace: Path, config: Any | None) -> Any | None:
    return EvolutionConfigOverlayStore(workspace).effective_config(config)


def self_tune_evolution_config(
    workspace: Path,
    config: Any | None,
    *,
    health: dict[str, Any],
    outcome_stats: dict[str, Any],
    sandbox_counts: dict[str, int],
    feedback_stats: dict[str, Any],
    actor: str = "auto_evolution",
    source: str = "evolution_maintenance",
) -> dict[str, Any]:
    store = EvolutionConfigOverlayStore(workspace)
    effective = store.effective_config(config)
    patches, evidence = EvolutionConfigTuner(workspace, effective).propose(
        health=health,
        outcome_stats=outcome_stats,
        sandbox_counts=sandbox_counts,
        feedback_stats=feedback_stats,
    )
    if not patches:
        return {
            "enabled": True,
            "applied_count": 0,
            "rejected_count": 0,
            "overlay": store.status(),
            "evidence": evidence,
        }
    result = store.apply_patches(
        effective,
        patches,
        actor=actor,
        source=source,
        evidence=evidence,
    )
    return {
        "enabled": True,
        "applied_count": len(result.applied),
        "rejected_count": len(result.rejected),
        "result": result.to_json(),
        "overlay": store.status(),
        "evidence": evidence,
    }


def _empty_overlay() -> dict[str, Any]:
    return {
        "schema_version": OVERLAY_SCHEMA_VERSION,
        "active": False,
        "overrides": {},
        "patch_count": 0,
        "last_patch_at": "",
        "last_actor": "",
    }


def _flatten(value: dict[str, Any], prefix: str = "") -> dict[str, Any]:
    flattened: dict[str, Any] = {}
    for key, item in value.items():
        path = f"{prefix}.{key}" if prefix else str(key)
        if isinstance(item, dict):
            flattened.update(_flatten(item, path))
        else:
            flattened[path] = item
    return flattened


def _get_path(obj: Any, path: str) -> Any:
    current = obj
    for part in str(path or "").split("."):
        if not part:
            continue
        if isinstance(current, dict):
            current = current.get(part)
        else:
            current = getattr(current, part, None)
        if current is None:
            return None
    return current


def _set_path(obj: Any, path: str, value: Any) -> None:
    parts = [part for part in str(path or "").split(".") if part]
    current = obj
    for part in parts[:-1]:
        current = current.get(part) if isinstance(current, dict) else getattr(current, part, None)
        if current is None:
            return
    if not parts:
        return
    if isinstance(current, dict):
        current[parts[-1]] = value
    elif current is not None and hasattr(current, parts[-1]):
        setattr(current, parts[-1], value)


def _set_path_in_mapping(mapping: dict[str, Any], path: str, value: Any) -> None:
    parts = [part for part in str(path or "").split(".") if part]
    current = mapping
    for part in parts[:-1]:
        item = current.get(part)
        if not isinstance(item, dict):
            item = {}
            current[part] = item
        current = item
    if parts:
        current[parts[-1]] = value


def _mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def _safe_int(value: Any, default: int) -> int:
    if isinstance(value, bool):
        return default
    try:
        return int(value)
    except (TypeError, ValueError):
        return default


def _safe_float(value: Any, default: float) -> float:
    if isinstance(value, bool):
        return default
    try:
        return float(value)
    except (TypeError, ValueError):
        return default
