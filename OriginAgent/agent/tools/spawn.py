"""Spawn tool for creating background subagents."""

from contextvars import ContextVar
from typing import TYPE_CHECKING, Any

from OriginAgent.agent.subagent_policy import SubagentPolicy
from OriginAgent.agent.tools.base import Tool, tool_parameters
from OriginAgent.agent.tools.schema import StringSchema, tool_parameters_schema
from OriginAgent.security.capabilities import CapabilitySnapshot
from OriginAgent.security.policy import PolicyDeniedError

if TYPE_CHECKING:
    from OriginAgent.agent.subagent import SubagentManager


@tool_parameters(
    tool_parameters_schema(
        task=StringSchema("The task for the subagent to complete"),
        label=StringSchema("Optional short label for the task (for display)"),
        required=["task"],
    )
)
class SpawnTool(Tool):
    """Tool to spawn a subagent for background task execution."""

    def __init__(self, manager: "SubagentManager"):
        self._manager = manager
        self._origin_channel: ContextVar[str] = ContextVar("spawn_origin_channel", default="cli")
        self._origin_chat_id: ContextVar[str] = ContextVar("spawn_origin_chat_id", default="direct")
        self._session_key: ContextVar[str] = ContextVar("spawn_session_key", default="cli:direct")
        self._origin_message_id: ContextVar[str | None] = ContextVar(
            "spawn_origin_message_id",
            default=None,
        )
        self._capability_snapshot: ContextVar[CapabilitySnapshot | None] = ContextVar(
            "spawn_capability_snapshot",
            default=None,
        )
        self._parent_subagent_id: ContextVar[str | None] = ContextVar(
            "spawn_parent_subagent_id",
            default=None,
        )
        self._root_subagent_id: ContextVar[str | None] = ContextVar(
            "spawn_root_subagent_id",
            default=None,
        )
        self._subagent_depth: ContextVar[int] = ContextVar(
            "spawn_subagent_depth",
            default=0,
        )
        self._delegated_policy: ContextVar[SubagentPolicy | None] = ContextVar(
            "spawn_delegated_policy",
            default=None,
        )

    def set_context(self, channel: str, chat_id: str, effective_key: str | None = None) -> None:
        """Set the origin context for subagent announcements."""
        self._origin_channel.set(channel)
        self._origin_chat_id.set(chat_id)
        self._session_key.set(effective_key or f"{channel}:{chat_id}")

    def set_origin_message_id(self, message_id: str | None) -> None:
        """Set the source message id for downstream deduplication."""
        self._origin_message_id.set(message_id)

    def set_capability_snapshot(self, snapshot: CapabilitySnapshot | None) -> None:
        self._capability_snapshot.set(snapshot)

    def set_nested_policy(
        self,
        *,
        parent_subagent_id: str | None,
        root_subagent_id: str | None,
        subagent_depth: int,
        delegated_policy: SubagentPolicy | None,
    ) -> None:
        self._parent_subagent_id.set(parent_subagent_id)
        self._root_subagent_id.set(root_subagent_id)
        self._subagent_depth.set(subagent_depth)
        self._delegated_policy.set(delegated_policy)

    @property
    def name(self) -> str:
        return "spawn"

    @property
    def description(self) -> str:
        return (
            "Spawn a subagent to handle a task in the background. "
            "Use this for complex or time-consuming tasks that can run independently. "
            "The subagent will complete the task and report back when done. "
            "For deliverables or existing projects, inspect the workspace first "
            "and use a dedicated subdirectory when helpful."
        )

    async def execute(self, task: str, label: str | None = None, **kwargs: Any) -> str:
        """Spawn a subagent to execute the given task."""
        snapshot = self._capability_snapshot.get()
        if snapshot is None:
            raise PolicyDeniedError(
                "Spawning subagents requires an explicit capability snapshot",
                code="spawn_capability_snapshot_missing",
                boundary="spawn",
                policy_rule="capability_snapshot_required",
            )
        if not snapshot.can_spawn:
            raise PolicyDeniedError(
                "Spawning subagents is not allowed by the current capability snapshot",
                code="spawn_capability_denied",
                boundary="spawn",
                policy_rule="capability_spawn_denied",
            )
        running = self._manager.get_running_count()
        limit = self._manager.max_concurrent_subagents
        if running >= limit:
            return (
                f"Cannot spawn subagent: concurrency limit reached "
                f"({running}/{limit} running). Wait for a running subagent "
                f"to complete before spawning a new one."
            )
        delegated_policy = self._delegated_policy.get()
        child_policy = delegated_policy.child_policy() if delegated_policy is not None else None
        return await self._manager.spawn(
            task=task,
            label=label,
            origin_channel=self._origin_channel.get(),
            origin_chat_id=self._origin_chat_id.get(),
            session_key=self._session_key.get(),
            origin_message_id=self._origin_message_id.get(),
            parent_capability_snapshot=snapshot,
            grant_id=None,
            delegated_policy=child_policy,
            parent_subagent_id=self._parent_subagent_id.get(),
            root_subagent_id=self._root_subagent_id.get(),
            subagent_depth=self._subagent_depth.get() + 1,
        )
