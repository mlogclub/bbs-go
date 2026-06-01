"""Runtime context helpers for AgentLoop."""

from __future__ import annotations

from typing import Any, Awaitable, Callable

from OriginAgent.agent.identity import RuntimeContext
from OriginAgent.agent.tools.context import RequestContext
from OriginAgent.bus.events import InboundMessage, OutboundMessage
from OriginAgent.bus.queue import MessageBus
from OriginAgent.security.capabilities import CapabilitySnapshot


def runtime_chat_id(msg: InboundMessage) -> str:
    """Return the chat id shown in runtime metadata for the model."""
    return str(msg.metadata.get("context_chat_id") or msg.chat_id)


def snapshot_for_trigger(trigger: str | None) -> CapabilitySnapshot:
    if trigger == "scheduled":
        return CapabilitySnapshot.scheduled_default()
    if trigger == "subagent":
        return CapabilitySnapshot.system_default().derive_subagent()
    if trigger == "system":
        return CapabilitySnapshot.system_default()
    return CapabilitySnapshot.user_turn()


def set_tool_context(
    tools: Any,
    *,
    channel: str,
    chat_id: str,
    message_id: str | None = None,
    metadata: dict | None = None,
    session_key: str | None = None,
    actor_id: str | None = None,
    trigger: str | None = None,
    capability_snapshot: CapabilitySnapshot | None = None,
    runtime_context: RuntimeContext | None = None,
    unified_session: bool = False,
    unified_session_key: str = "unified:default",
) -> None:
    """Update context for all tools that need routing info."""
    if runtime_context is not None:
        channel = runtime_context.channel
        chat_id = runtime_context.chat_id
        session_key = runtime_context.session_key
        actor_id = runtime_context.actor_id
        trigger = runtime_context.trigger

    if session_key is not None:
        effective_key = session_key
    elif unified_session:
        effective_key = unified_session_key
    else:
        effective_key = f"{channel}:{chat_id}"

    raw_tools = getattr(tools, "_tools", None)
    if isinstance(raw_tools, dict):
        context_tool_names = [
            name
            for name, tool in raw_tools.items()
            if hasattr(tool, "set_context") or hasattr(tool, "set_capability_snapshot")
        ]
    else:
        candidates = list(getattr(tools, "tool_names", ()) or ())
        if not candidates:
            candidates = ["spawn", "cron", "long_task", "complete_goal", "message", "my"]
        context_tool_names = []
        for name in dict.fromkeys(candidates):
            tool = tools.get(name)
            if tool is not None and (
                hasattr(tool, "set_context") or hasattr(tool, "set_capability_snapshot")
            ):
                context_tool_names.append(name)

    request_ctx = RequestContext(
        channel=channel,
        chat_id=chat_id,
        message_id=message_id,
        session_key=effective_key,
        metadata=metadata or {},
        actor_id=actor_id,
        trigger=trigger,
        capability_snapshot=capability_snapshot,
    )
    if hasattr(tools, "set_capability_snapshot"):
        tools.set_capability_snapshot(capability_snapshot)
    if hasattr(tools, "set_audit_context"):
        tools.set_audit_context(actor_id=actor_id, session_key=effective_key)
    for name in context_tool_names:
        if tool := tools.get(name):
            permissions = tuple(getattr(tool, "_domain_tool_permissions", ()) or ())
            if hasattr(tool, "set_capability_snapshot"):
                tool.set_capability_snapshot(capability_snapshot)
            if hasattr(tool, "set_context"):
                if any(permission.startswith("device:") for permission in permissions):
                    if actor_id is not None and trigger is not None:
                        tool.set_context(actor_id, trigger)
                elif name == "spawn":
                    tool.set_context(channel, chat_id, effective_key=effective_key)
                    if hasattr(tool, "set_origin_message_id"):
                        tool.set_origin_message_id(message_id)
                elif name == "cron":
                    tool.set_context(channel, chat_id, metadata=metadata, session_key=session_key)
                elif name in {"long_task", "complete_goal"}:
                    tool.set_context(channel, chat_id, session_key=effective_key)
                elif name == "message":
                    tool.set_context(channel, chat_id, message_id, metadata=metadata)
                elif name == "my":
                    tool.set_context(channel, chat_id)
                else:
                    try:
                        tool.set_context(request_ctx)
                    except TypeError:
                        tool.set_context(channel, chat_id)


async def build_bus_progress_callback(
    bus: MessageBus,
    msg: InboundMessage,
) -> Callable[..., Awaitable[None]]:
    """Build a progress callback that publishes to the message bus."""

    async def _bus_progress(
        content: str,
        *,
        tool_hint: bool = False,
        tool_events: list[dict[str, Any]] | None = None,
        reasoning: bool = False,
        reasoning_end: bool = False,
    ) -> None:
        meta = dict(msg.metadata or {})
        meta["_progress"] = True
        meta["_tool_hint"] = tool_hint
        if reasoning:
            meta["_reasoning_delta"] = True
        if reasoning_end:
            meta["_reasoning_end"] = True
        if tool_events:
            meta["_tool_events"] = tool_events
        await bus.publish_outbound(
            OutboundMessage(
                channel=msg.channel,
                chat_id=msg.chat_id,
                content=content,
                metadata=meta,
            )
        )

    return _bus_progress


async def build_retry_wait_callback(
    bus: MessageBus,
    msg: InboundMessage,
) -> Callable[[str], Awaitable[None]]:
    """Build a retry-wait callback that publishes to the message bus."""

    async def _on_retry_wait(content: str) -> None:
        meta = dict(msg.metadata or {})
        meta["_retry_wait"] = True
        await bus.publish_outbound(
            OutboundMessage(
                channel=msg.channel,
                chat_id=msg.chat_id,
                content=content,
                metadata=meta,
            )
        )

    return _on_retry_wait
