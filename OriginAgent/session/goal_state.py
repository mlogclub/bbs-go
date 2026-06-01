"""Session metadata helpers for sustained goals."""

from __future__ import annotations

import json
from typing import Any, Mapping, MutableMapping

from OriginAgent.session.manager import SessionManager

GOAL_STATE_KEY = "goal_state"
_LEGACY_GOAL_STATE_SESSION_KEY = "thread_goal"
_MAX_OBJECTIVE_IN_RUNTIME = 4000
_MAX_OBJECTIVE_WS = 600


def _session_goal_raw(metadata: Mapping[str, Any] | None) -> Any:
    if not metadata:
        return None
    if GOAL_STATE_KEY in metadata:
        return metadata.get(GOAL_STATE_KEY)
    return metadata.get(_LEGACY_GOAL_STATE_SESSION_KEY)


def discard_legacy_goal_state_key(metadata: MutableMapping[str, Any]) -> None:
    """Remove the legacy metadata key after writing the canonical key."""
    metadata.pop(_LEGACY_GOAL_STATE_SESSION_KEY, None)


def goal_state_raw(metadata: Mapping[str, Any] | None) -> Any:
    """Return goal state from the canonical key, falling back to the legacy key."""
    return _session_goal_raw(metadata)


def parse_goal_state(blob: Any) -> dict[str, Any] | None:
    if blob is None:
        return None
    if isinstance(blob, dict):
        return blob
    if isinstance(blob, str):
        try:
            parsed = json.loads(blob)
        except json.JSONDecodeError:
            return None
        return parsed if isinstance(parsed, dict) else None
    return None


def sustained_goal_active(metadata: Mapping[str, Any] | None) -> bool:
    """Return true when session metadata carries an active sustained goal."""
    goal = parse_goal_state(goal_state_raw(metadata))
    return isinstance(goal, dict) and goal.get("status") == "active"


def goal_state_runtime_lines(metadata: Mapping[str, Any] | None) -> list[str]:
    """Lines appended inside Runtime Context when a goal is active."""
    if not metadata:
        return []
    goal = parse_goal_state(_session_goal_raw(metadata))
    if not isinstance(goal, dict) or goal.get("status") != "active":
        return []
    objective = str(goal.get("objective") or "").strip()
    if not objective:
        return ["Goal: active (no objective text stored)."]
    if len(objective) > _MAX_OBJECTIVE_IN_RUNTIME:
        objective = objective[:_MAX_OBJECTIVE_IN_RUNTIME].rstrip() + "\n... (truncated)"
    lines = ["Goal (active):", objective]
    summary = str(goal.get("ui_summary") or "").strip()
    if summary:
        lines.append(f"Summary: {summary}")
    return lines


def goal_state_ws_blob(metadata: Mapping[str, Any] | None) -> dict[str, Any]:
    """Return a JSON-safe WebSocket snapshot for one chat's goal state."""
    goal = parse_goal_state(_session_goal_raw(metadata)) if metadata else None
    if isinstance(goal, dict) and goal.get("status") == "active":
        objective = str(goal.get("objective") or "").strip()
        if len(objective) > _MAX_OBJECTIVE_WS:
            objective = objective[:_MAX_OBJECTIVE_WS].rstrip() + "..."
        summary = str(goal.get("ui_summary") or "").strip()[:120]
        blob: dict[str, Any] = {"active": True}
        if objective:
            blob["objective"] = objective
        if summary:
            blob["ui_summary"] = summary
        return blob
    return {"active": False}


def runner_wall_llm_timeout_s(
    sessions: SessionManager,
    session_key: str | None,
    *,
    metadata: Mapping[str, Any] | None = None,
) -> float | None:
    """Return an LLM wall timeout override for active goals."""
    meta: Mapping[str, Any] | None = metadata
    if meta is None and session_key:
        meta = sessions.get_or_create(session_key).metadata
    return 0.0 if sustained_goal_active(meta) else None
