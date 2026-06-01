"""Sustained goal tools for long-running OriginAgent tasks."""

from __future__ import annotations

from contextvars import ContextVar
from datetime import datetime
from typing import Any

from OriginAgent.agent.tools.base import Tool, tool_parameters
from OriginAgent.agent.tools.schema import StringSchema, tool_parameters_schema
from OriginAgent.bus.events import OutboundMessage
from OriginAgent.session.goal_state import (
    GOAL_STATE_KEY,
    discard_legacy_goal_state_key,
    goal_state_raw,
    goal_state_ws_blob,
    parse_goal_state,
)
from OriginAgent.session.manager import SessionManager


def _iso_now() -> str:
    return datetime.now().isoformat()


class _GoalToolsMixin:
    """Shared session lookup and WebSocket sync."""

    def __init__(self, sessions: SessionManager, bus: Any | None = None) -> None:
        self._sessions = sessions
        self._bus = bus
        self._channel: ContextVar[str] = ContextVar("goal_channel", default="")
        self._chat_id: ContextVar[str] = ContextVar("goal_chat_id", default="")
        self._session_key: ContextVar[str] = ContextVar("goal_session_key", default="")

    def set_context(
        self,
        channel: str,
        chat_id: str,
        *,
        session_key: str | None = None,
    ) -> None:
        self._channel.set(channel)
        self._chat_id.set(chat_id)
        self._session_key.set(session_key or f"{channel}:{chat_id}")

    def _session(self):
        key = self._session_key.get()
        if not key:
            return None
        return self._sessions.get_or_create(key)

    async def _publish_goal_state_ws(self, metadata: dict[str, Any]) -> None:
        if self._bus is None or self._channel.get() != "websocket":
            return
        chat_id = self._chat_id.get().strip()
        if not chat_id:
            return
        await self._bus.publish_outbound(
            OutboundMessage(
                channel="websocket",
                chat_id=chat_id,
                content="",
                metadata={
                    "_goal_state_sync": True,
                    "goal_state": goal_state_ws_blob(metadata),
                },
            )
        )


@tool_parameters(
    tool_parameters_schema(
        goal=StringSchema(
            "Sustained objective for this chat thread. Make it idempotent, "
            "self-contained, bounded, and explicit about done-ness.",
            max_length=12_000,
        ),
        ui_summary=StringSchema(
            "Optional one-line label for UI/session lists, at most 120 characters.",
            max_length=120,
            nullable=True,
        ),
        required=["goal"],
        additional_properties=False,
    )
)
class LongTaskTool(Tool, _GoalToolsMixin):
    """Mark the current thread as focused on a sustained objective."""

    def __init__(self, sessions: SessionManager, bus: Any | None = None) -> None:
        _GoalToolsMixin.__init__(self, sessions, bus)

    @property
    def name(self) -> str:
        return "long_task"

    @property
    def description(self) -> str:
        return (
            "Record a long-running goal for this chat. The active goal is mirrored "
            "into Runtime Context every turn until complete_goal is called."
        )

    async def execute(self, goal: str, ui_summary: str | None = None, **_: Any) -> str:
        sess = self._session()
        if sess is None:
            return "Error: long_task requires an active chat session."
        prior = parse_goal_state(goal_state_raw(sess.metadata))
        if isinstance(prior, dict) and prior.get("status") == "active":
            return (
                "Error: a sustained goal is already active. "
                "Use complete_goal when finished, or ask the user before replacing it."
            )

        summary = (ui_summary or "").strip()[:120]
        sess.metadata[GOAL_STATE_KEY] = {
            "status": "active",
            "objective": goal.strip(),
            "ui_summary": summary,
            "started_at": _iso_now(),
        }
        discard_legacy_goal_state_key(sess.metadata)
        self._sessions.save(sess)
        await self._publish_goal_state_ws(sess.metadata)
        suffix = f"\nSummary line: {summary}" if summary else ""
        return (
            "Goal recorded. Keep working toward the objective using ordinary tools. "
            "When fully done, call complete_goal with a short recap."
            f"{suffix}"
        )


@tool_parameters(
    tool_parameters_schema(
        recap=StringSchema(
            "Brief recap for the user. If the goal was cancelled or replaced, say so honestly.",
            max_length=8000,
            nullable=True,
        ),
        required=[],
        additional_properties=False,
    )
)
class CompleteGoalTool(Tool, _GoalToolsMixin):
    """Mark the active sustained goal as complete."""

    def __init__(self, sessions: SessionManager, bus: Any | None = None) -> None:
        _GoalToolsMixin.__init__(self, sessions, bus)

    @property
    def name(self) -> str:
        return "complete_goal"

    @property
    def description(self) -> str:
        return (
            "End bookkeeping for the active sustained goal after it is delivered, "
            "cancelled, redirected, or replaced."
        )

    async def execute(self, recap: str | None = None, **_: Any) -> str:
        sess = self._session()
        if sess is None:
            return "Error: complete_goal requires an active chat session."
        prior = parse_goal_state(goal_state_raw(sess.metadata))
        if not isinstance(prior, dict) or prior.get("status") != "active":
            return "No active goal to complete."

        ended = _iso_now()
        sess.metadata[GOAL_STATE_KEY] = {
            **prior,
            "status": "completed",
            "completed_at": ended,
            "recap": (recap or "").strip(),
        }
        discard_legacy_goal_state_key(sess.metadata)
        self._sessions.save(sess)
        await self._publish_goal_state_ws(sess.metadata)
        tail = (recap or "").strip()
        if tail:
            return f"Goal marked complete ({ended}). Recap:\n{tail}"
        return f"Goal marked complete ({ended})."
