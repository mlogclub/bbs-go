"""Runtime actor identity resolution."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Literal


RuntimeTrigger = Literal["user_initiated", "scheduled", "system", "subagent"]
RuntimeSource = Literal["user_turn", "cron", "system", "subagent"]


@dataclass(frozen=True)
class RuntimeContext:
    actor_id: str
    trigger: RuntimeTrigger
    channel: str
    chat_id: str
    session_key: str | None = None
    source: RuntimeSource = "user_turn"


class ActorResolver:
    def resolve(
        self,
        *,
        channel: str,
        chat_id: str,
        sender_id: str,
        metadata: dict,
    ) -> str:
        return self.resolve_runtime_context(
            channel=channel,
            chat_id=chat_id,
            sender_id=sender_id,
            metadata=metadata,
        ).actor_id

    def resolve_runtime_context(
        self,
        *,
        channel: str,
        chat_id: str,
        sender_id: str,
        metadata: dict,
        session_key: str | None = None,
        routing_channel: str | None = None,
        routing_chat_id: str | None = None,
    ) -> RuntimeContext:
        sender = str(sender_id or "").strip()
        actor_id = sender or "unknown"
        normalized_channel = str(channel or "").strip()
        context_channel = str(routing_channel or normalized_channel or "").strip()
        context_chat_id = str(routing_chat_id if routing_chat_id is not None else chat_id or "")
        metadata = metadata if isinstance(metadata, dict) else {}
        if sender == "subagent" or metadata.get("injected_event") == "subagent_result":
            trigger: RuntimeTrigger = "subagent"
            source: RuntimeSource = "subagent"
        elif normalized_channel == "cron":
            trigger = "scheduled"
            source = "cron"
        elif normalized_channel == "system":
            trigger = "system"
            source = "system"
        else:
            trigger = "user_initiated"
            source = "user_turn"
        return RuntimeContext(
            actor_id=actor_id,
            trigger=trigger,
            channel=context_channel,
            chat_id=context_chat_id,
            session_key=session_key,
            source=source,
        )
