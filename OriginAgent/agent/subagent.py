"""Subagent manager for background task execution."""

import asyncio
import json
import time
import uuid
from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any, Callable

from loguru import logger

from OriginAgent.agent.agent_tool_setup import register_default_tools
from OriginAgent.agent.hook import AgentHook, AgentHookContext
from OriginAgent.agent.runner import AgentRunner, AgentRunSpec
from OriginAgent.agent.subagent_policy import SubagentPolicy
from OriginAgent.agent.subagent_provider import SubagentProviderSelector
from OriginAgent.agent.subagent_records import (
    JsonlSubagentRecordStore,
    SubagentLifecycleRecord,
    SubagentTaskRecord,
    SubagentToolRecord,
    summarize_payload,
    summarize_text,
)
from OriginAgent.agent.tools.registry import ToolRegistry
from OriginAgent.agent.tools.audit import JsonlToolAuditSink, ToolAuditConfig
from OriginAgent.bus.events import InboundMessage
from OriginAgent.bus.queue import MessageBus
from OriginAgent.config.schema import AgentDefaults, ExecToolConfig, ToolsConfig, WebToolsConfig
from OriginAgent.providers.base import LLMProvider
from OriginAgent.security.capabilities import CapabilitySnapshot, intersect_capability_snapshots
from OriginAgent.security.grants import CapabilityGrantStore
from OriginAgent.security.policy import PolicyDeniedError
from OriginAgent.utils.prompt_templates import render_template

_GRANT_ERROR_MESSAGE = "Capability grant is missing, expired, or revoked."


@dataclass(slots=True)
class SubagentStatus:
    """Real-time status of a running subagent."""

    task_id: str
    label: str
    task_description: str
    started_at: float          # time.monotonic()
    phase: str = "initializing"  # initializing | awaiting_tools | tools_completed | final_response | done | error
    iteration: int = 0
    tool_events: list = field(default_factory=list)   # [{name, status, detail}, ...]
    usage: dict = field(default_factory=dict)          # token usage
    stop_reason: str | None = None
    error: str | None = None
    root_subagent_id: str | None = None
    parent_subagent_id: str | None = None
    subagent_depth: int = 1


class _SubagentHook(AgentHook):
    """Hook for subagent execution — logs tool calls and updates status."""

    def __init__(self, task_id: str, status: SubagentStatus | None = None) -> None:
        super().__init__()
        self._task_id = task_id
        self._status = status

    async def before_execute_tools(self, context: AgentHookContext) -> None:
        for tool_call in context.tool_calls:
            args_str = json.dumps(tool_call.arguments, ensure_ascii=False)
            logger.debug(
                "Subagent [{}] executing: {} with arguments: {}",
                self._task_id, tool_call.name, args_str,
            )

    async def after_iteration(self, context: AgentHookContext) -> None:
        if self._status is None:
            return
        self._status.iteration = context.iteration
        self._status.tool_events = list(context.tool_events)
        self._status.usage = dict(context.usage)
        if context.error:
            self._status.error = str(context.error)


class SubagentManager:
    """Manages background subagent execution."""

    def __init__(
        self,
        provider: LLMProvider,
        workspace: Path,
        bus: MessageBus,
        max_tool_result_chars: int,
        model: str | None = None,
        web_config: "WebToolsConfig | None" = None,
        content_read_config: Any | None = None,
        exec_config: "ExecToolConfig | None" = None,
        restrict_to_workspace: bool = False,
        disabled_skills: list[str] | None = None,
        max_iterations: int | None = None,
        grant_store: CapabilityGrantStore | None = None,
        preset_snapshot_loader: Callable[[str], Any] | None = None,
        delegated_model_preset: str | None = None,
    ):
        defaults = AgentDefaults()
        self.provider = provider
        self.workspace = workspace
        self.bus = bus
        self.model = model or provider.get_default_model()
        self.web_config = web_config or WebToolsConfig()
        self.content_read_config = content_read_config
        self.max_tool_result_chars = max_tool_result_chars
        self.exec_config = exec_config or ExecToolConfig()
        self.restrict_to_workspace = restrict_to_workspace
        self.disabled_skills = set(disabled_skills or [])
        self.max_iterations = (
            max_iterations
            if max_iterations is not None
            else defaults.max_tool_iterations
        )
        self.grant_store = grant_store
        self.max_concurrent_subagents = defaults.max_concurrent_subagents
        self.tools_config = ToolsConfig()
        self._tool_audit_config = ToolAuditConfig.from_config(self.tools_config.audit)
        self._provider_selector = SubagentProviderSelector(
            provider=provider,
            model=self.model,
            preset_snapshot_loader=preset_snapshot_loader,
            delegated_preset=delegated_model_preset,
        )
        self.runner = AgentRunner(provider)
        self.records = JsonlSubagentRecordStore(workspace)
        self._running_tasks: dict[str, asyncio.Task[None]] = {}
        self._task_statuses: dict[str, SubagentStatus] = {}
        self._session_tasks: dict[str, set[str]] = {}  # session_key -> {task_id, ...}
        self._child_counts: dict[str, int] = {}

    def set_provider(self, provider: LLMProvider, model: str) -> None:
        self.provider = provider
        self.model = model
        self.runner.provider = provider
        self._provider_selector.set_runtime(provider, model)

    async def spawn(
        self,
        task: str,
        label: str | None = None,
        origin_channel: str = "cli",
        origin_chat_id: str = "direct",
        session_key: str | None = None,
        origin_message_id: str | None = None,
        capability_snapshot: CapabilitySnapshot | None = None,
        parent_capability_snapshot: CapabilitySnapshot | None = None,
        grant_id: str | None = None,
        delegated_policy: SubagentPolicy | None = None,
        parent_subagent_id: str | None = None,
        root_subagent_id: str | None = None,
        subagent_depth: int = 1,
    ) -> str:
        """Spawn a subagent to execute a task in the background.

        ``parent_capability_snapshot`` is the parent turn snapshot. The manager
        derives the subagent snapshot internally so grants cannot skip the
        least-privilege subagent boundary. ``capability_snapshot`` is retained
        as a compatibility alias for existing internal callers and has the same
        parent-snapshot meaning.
        """
        if self.get_running_count() >= self.max_concurrent_subagents:
            return (
                "Cannot spawn subagent: concurrency limit reached "
                f"({self.max_concurrent_subagents}). Please wait for a running subagent to finish."
            )
        if delegated_policy is not None:
            if subagent_depth > delegated_policy.max_subagent_depth:
                return (
                    "Cannot spawn nested subagent: maximum delegated depth reached "
                    f"({delegated_policy.max_subagent_depth})."
                )
            if parent_subagent_id is not None:
                child_limit = delegated_policy.max_children_per_subagent
                current_children = self._child_counts.get(parent_subagent_id, 0)
                if child_limit >= 0 and current_children >= child_limit:
                    return (
                        "Cannot spawn nested subagent: child delegation limit reached "
                        f"({current_children}/{child_limit})."
                    )
        task_id = str(uuid.uuid4())[:8]
        display_label = label or task[:30] + ("..." if len(task) > 30 else "")
        resolved_root = root_subagent_id or task_id
        origin = {
            "channel": origin_channel,
            "chat_id": origin_chat_id,
            "session_key": session_key,
            "root_subagent_id": resolved_root,
            "parent_subagent_id": parent_subagent_id,
            "subagent_depth": str(subagent_depth),
        }
        policy = delegated_policy or SubagentPolicy.default_for_parent(
            parent_capability_snapshot or capability_snapshot
        )
        effective_snapshot = (
            policy.capability_snapshot
            if delegated_policy is not None
            else self._snapshot_for_spawn(
                parent_snapshot=parent_capability_snapshot or capability_snapshot,
                grant_id=grant_id,
            )
        )

        status = SubagentStatus(
            task_id=task_id,
            label=display_label,
            task_description=task,
            started_at=time.monotonic(),
            root_subagent_id=resolved_root,
            parent_subagent_id=parent_subagent_id,
            subagent_depth=subagent_depth,
        )
        self._task_statuses[task_id] = status
        self.records.append_task(SubagentTaskRecord(
            subagent_id=task_id,
            root_subagent_id=resolved_root,
            parent_subagent_id=parent_subagent_id,
            subagent_depth=subagent_depth,
            parent_session_key=session_key,
            origin_channel=origin_channel,
            origin_chat_id=origin_chat_id,
            origin_message_id=origin_message_id,
            task_label=display_label,
            task_summary=summarize_text(task),
            delegated_profile_summary="pending",
            allowed_tools_summary=[],
            provider_summary="pending",
            isolation_mode="shared_process",
            terminal_status="spawned",
            started_at=self._utcnow(),
        ))
        self.records.append_lifecycle(SubagentLifecycleRecord(
            subagent_id=task_id,
            root_subagent_id=resolved_root,
            parent_subagent_id=parent_subagent_id,
            subagent_depth=subagent_depth,
            parent_session_key=session_key,
            state="spawned",
            detail=summarize_text(display_label, max_chars=120),
        ))
        if parent_subagent_id is not None:
            self._child_counts[parent_subagent_id] = self._child_counts.get(parent_subagent_id, 0) + 1

        bg_task = asyncio.create_task(
            self._run_subagent(
                task_id,
                task,
                display_label,
                origin,
                status,
                origin_message_id,
                effective_snapshot,
                policy,
            )
        )
        self._running_tasks[task_id] = bg_task
        if session_key:
            self._session_tasks.setdefault(session_key, set()).add(task_id)

        def _cleanup(_: asyncio.Task) -> None:
            self._running_tasks.pop(task_id, None)
            self._task_statuses.pop(task_id, None)
            if session_key and (ids := self._session_tasks.get(session_key)):
                ids.discard(task_id)
                if not ids:
                    del self._session_tasks[session_key]

        bg_task.add_done_callback(_cleanup)

        logger.info("Spawned subagent [{}]: {}", task_id, display_label)
        return f"Subagent [{display_label}] started (id: {task_id}). I'll notify you when it completes."

    def _snapshot_for_spawn(
        self,
        *,
        parent_snapshot: CapabilitySnapshot | None,
        grant_id: str | None = None,
    ) -> CapabilitySnapshot:
        parent = parent_snapshot or CapabilitySnapshot.system_default()
        base = parent.derive_subagent()
        if not grant_id:
            return base
        if self.grant_store is None:
            _raise_grant_denied("capability_grant_missing")
        grant = self.grant_store.get(grant_id)
        if grant is None:
            _raise_grant_denied("capability_grant_missing")
        if grant.is_revoked():
            _raise_grant_denied("capability_grant_revoked")
        if grant.is_expired():
            _raise_grant_denied("capability_grant_expired")
        return intersect_capability_snapshots(
            base,
            grant.to_subagent_snapshot(),
            source="subagent",
            trigger="subagent",
        )

    async def _run_subagent(
        self,
        task_id: str,
        task: str,
        label: str,
        origin: dict[str, str],
        status: SubagentStatus,
        origin_message_id: str | None = None,
        capability_snapshot: CapabilitySnapshot | None = None,
        delegated_policy: SubagentPolicy | None = None,
    ) -> None:
        """Execute the subagent task and announce the result."""
        logger.info("Subagent [{}] starting task: {}", task_id, label)
        last_phase = "running"

        async def _on_checkpoint(payload: dict) -> None:
            nonlocal last_phase
            status.phase = payload.get("phase", status.phase)
            status.iteration = payload.get("iteration", status.iteration)
            phase = str(payload.get("phase") or "").strip()
            if phase and phase != last_phase and phase in {"awaiting_tools", "tools_completed", "final_response"}:
                self.records.append_lifecycle(SubagentLifecycleRecord(
                    subagent_id=task_id,
                    root_subagent_id=status.root_subagent_id,
                    parent_subagent_id=status.parent_subagent_id,
                    subagent_depth=status.subagent_depth,
                    parent_session_key=origin.get("session_key"),
                    state=phase,  # type: ignore[arg-type]
                    detail=f"iteration={status.iteration}",
                ))
                last_phase = phase

        try:
            policy = delegated_policy or SubagentPolicy.default_for_parent(capability_snapshot)
            snapshot = policy.capability_snapshot
            provider_selection = self._provider_selector.select()
            self.runner.provider = provider_selection.provider
            tools = ToolRegistry(
                audit_sink=JsonlToolAuditSink(self.workspace),
                audit_config=self._tool_audit_config,
                capability_snapshot=snapshot,
                execution_observer=_SubagentToolObserver(
                    records=self.records,
                    subagent_id=task_id,
                    root_subagent_id=status.root_subagent_id,
                    parent_subagent_id=status.parent_subagent_id,
                    subagent_depth=status.subagent_depth,
                    parent_session_key=origin.get("session_key"),
                ),
            )
            tools.set_audit_context(
                actor_id="subagent",
                session_key=origin.get("session_key"),
                subagent_task_id=task_id,
                parent_session_key=origin.get("session_key"),
                origin_channel=origin.get("channel"),
                origin_chat_id=origin.get("chat_id"),
            )
            self.records.append_lifecycle(SubagentLifecycleRecord(
                subagent_id=task_id,
                root_subagent_id=status.root_subagent_id,
                parent_subagent_id=status.parent_subagent_id,
                subagent_depth=status.subagent_depth,
                parent_session_key=origin.get("session_key"),
                state="running",
                detail=summarize_text(
                    ",".join(sorted(policy.allowed_tool_names)) or "delegated",
                    max_chars=120,
                ),
            ))
            self.records.append_task(SubagentTaskRecord(
                subagent_id=task_id,
                root_subagent_id=status.root_subagent_id,
                parent_subagent_id=status.parent_subagent_id,
                subagent_depth=status.subagent_depth,
                parent_session_key=origin.get("session_key"),
                origin_channel=origin.get("channel"),
                origin_chat_id=origin.get("chat_id"),
                origin_message_id=origin_message_id,
                task_label=label,
                task_summary=summarize_text(task),
                delegated_profile_summary=self.records.delegated_profile_summary(
                    capability_snapshot=snapshot,
                    allow_web=policy.allow_web,
                    allowed_tool_names=policy.allowed_tool_names,
                ),
                allowed_tools_summary=self.records.allowed_tools_summary(policy.allowed_tool_names),
                provider_summary=provider_selection.provider_summary,
                isolation_mode="shared_process",
                terminal_status="running",
                started_at=self._utcnow(),
            ))
            register_default_tools(
                tools,
                workspace=self.workspace,
                bus=self.bus,
                config=self.tools_config,
                web_config=(self.web_config if policy.allow_web else WebToolsConfig(enable=False)),
                exec_config=self.exec_config,
                restrict_to_workspace=self.restrict_to_workspace,
                sessions=None,
                pending_queues={},
                cron_service=None,
                audit_config=self._tool_audit_config,
                domain_pack_manager=None,
                background_review_service=None,
                curator_service=None,
                session_search_index_service=None,
                subagent_manager=self,
                file_state_store=None,
                provider_snapshot_loader=None,
                image_generation_provider_configs={},
                timezone="UTC",
                runtime_profile="automation",
                introspection_service=None,
                confirmation_store=None,
                domain_runtime_overrides=None,
                evolution_config=None,
                allowed_tool_names=(
                    frozenset(set(policy.allowed_tool_names) | {"spawn"})
                    if policy.allow_nested_spawn
                    else policy.allowed_tool_names
                ),
            )
            self._configure_spawn_tool(
                tools=tools,
                origin=origin,
                origin_message_id=origin_message_id,
                policy=policy,
                status=status,
            )
            system_prompt = self._build_subagent_prompt()
            messages: list[dict[str, Any]] = [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": task},
            ]

            result = await self.runner.run(AgentRunSpec(
                initial_messages=messages,
                tools=tools,
                model=provider_selection.model,
                max_iterations=self.max_iterations,
                max_tool_result_chars=self.max_tool_result_chars,
                hook=_SubagentHook(task_id, status),
                max_iterations_message="Task completed but no final response was generated.",
                error_message=None,
                fail_on_tool_error=True,
                checkpoint_callback=_on_checkpoint,
            ))
            status.phase = "done"
            status.stop_reason = result.stop_reason
            if result.tool_events:
                self.records.append_lifecycle(SubagentLifecycleRecord(
                    subagent_id=task_id,
                    root_subagent_id=status.root_subagent_id,
                    parent_subagent_id=status.parent_subagent_id,
                    subagent_depth=status.subagent_depth,
                    parent_session_key=origin.get("session_key"),
                    state="tools_completed",
                    detail=summarize_text(result.stop_reason or "tools completed", max_chars=120),
                ))

            if result.stop_reason == "tool_error":
                status.tool_events = list(result.tool_events)
                self._record_terminal_task(
                    task_id=task_id,
                    label=label,
                    task=task,
                    origin=origin,
                    origin_message_id=origin_message_id,
                    terminal_status="failed",
                    stop_reason=result.stop_reason,
                    failure_summary=self._format_partial_progress(result),
                    policy=policy,
                    provider_summary=provider_selection.provider_summary,
                )
                await self._announce_result(
                    task_id, label, task,
                    self._format_partial_progress(result),
                    origin, "error", origin_message_id,
                )
            elif result.stop_reason == "error":
                self._record_terminal_task(
                    task_id=task_id,
                    label=label,
                    task=task,
                    origin=origin,
                    origin_message_id=origin_message_id,
                    terminal_status="failed",
                    stop_reason=result.stop_reason,
                    failure_summary=result.error or "subagent execution failed",
                    policy=policy,
                    provider_summary=provider_selection.provider_summary,
                )
                await self._announce_result(
                    task_id, label, task,
                    result.error or "Error: subagent execution failed.",
                    origin, "error", origin_message_id,
                )
            else:
                final_result = result.final_content or "Task completed but no final response was generated."
                logger.info("Subagent [{}] completed successfully", task_id)
                self.records.append_lifecycle(SubagentLifecycleRecord(
                    subagent_id=task_id,
                    root_subagent_id=status.root_subagent_id,
                    parent_subagent_id=status.parent_subagent_id,
                    subagent_depth=status.subagent_depth,
                    parent_session_key=origin.get("session_key"),
                    state="final_response",
                    detail=summarize_text(final_result, max_chars=120),
                ))
                self._record_terminal_task(
                    task_id=task_id,
                    label=label,
                    task=task,
                    origin=origin,
                    origin_message_id=origin_message_id,
                    terminal_status="completed",
                    stop_reason=result.stop_reason,
                    failure_summary=None,
                    policy=policy,
                    provider_summary=provider_selection.provider_summary,
                )
                await self._announce_result(task_id, label, task, final_result, origin, "ok", origin_message_id)

        except asyncio.CancelledError:
            status.phase = "error"
            status.error = "cancelled"
            self.records.append_lifecycle(SubagentLifecycleRecord(
                subagent_id=task_id,
                root_subagent_id=status.root_subagent_id,
                parent_subagent_id=status.parent_subagent_id,
                subagent_depth=status.subagent_depth,
                parent_session_key=origin.get("session_key"),
                state="cancelled",
                detail="cancelled",
            ))
            self.records.append_task(SubagentTaskRecord(
                subagent_id=task_id,
                root_subagent_id=status.root_subagent_id,
                parent_subagent_id=status.parent_subagent_id,
                subagent_depth=status.subagent_depth,
                parent_session_key=origin.get("session_key"),
                origin_channel=origin.get("channel"),
                origin_chat_id=origin.get("chat_id"),
                origin_message_id=origin_message_id,
                task_label=label,
                task_summary=summarize_text(task),
                delegated_profile_summary="cancelled",
                allowed_tools_summary=[],
                provider_summary="inherit:cancelled",
                isolation_mode="shared_process",
                terminal_status="cancelled",
                stop_reason="cancelled",
                failure_summary="cancelled",
                started_at=self._utcnow(),
                ended_at=self._utcnow(),
            ))
            raise
        except Exception as e:
            status.phase = "error"
            status.error = str(e)
            self.records.append_lifecycle(SubagentLifecycleRecord(
                subagent_id=task_id,
                root_subagent_id=status.root_subagent_id,
                parent_subagent_id=status.parent_subagent_id,
                subagent_depth=status.subagent_depth,
                parent_session_key=origin.get("session_key"),
                state="failed",
                detail=summarize_text(str(e), max_chars=120),
            ))
            self.records.append_task(SubagentTaskRecord(
                subagent_id=task_id,
                root_subagent_id=status.root_subagent_id,
                parent_subagent_id=status.parent_subagent_id,
                subagent_depth=status.subagent_depth,
                parent_session_key=origin.get("session_key"),
                origin_channel=origin.get("channel"),
                origin_chat_id=origin.get("chat_id"),
                origin_message_id=origin_message_id,
                task_label=label,
                task_summary=summarize_text(task),
                delegated_profile_summary="failed",
                allowed_tools_summary=[],
                provider_summary=f"inherit:{self.model}",
                isolation_mode="shared_process",
                terminal_status="failed",
                stop_reason="error",
                failure_summary=summarize_text(str(e), max_chars=240),
                started_at=self._utcnow(),
                ended_at=self._utcnow(),
            ))
            logger.exception("Subagent [{}] failed", task_id)
            await self._announce_result(task_id, label, task, f"Error: {e}", origin, "error", origin_message_id)
        finally:
            self.runner.provider = self.provider

    async def _announce_result(
        self,
        task_id: str,
        label: str,
        task: str,
        result: str,
        origin: dict[str, str],
        status: str,
        origin_message_id: str | None = None,
    ) -> None:
        """Announce the subagent result to the main agent via the message bus."""
        status_text = "completed successfully" if status == "ok" else "failed"

        announce_content = render_template(
            "agent/subagent_announce.md",
            label=label,
            status_text=status_text,
            task=task,
            result=result,
        )

        # Route as channel="system" only for internal bus dispatch.
        # LLM-facing content is wrapped by AgentLoop as an internal_event block.
        # Use session_key_override to align with the main agent's effective
        # session key (which accounts for unified sessions) so the result is
        # routed to the correct pending queue (mid-turn injection) instead of
        # being dispatched as a competing independent task.
        override = origin.get("session_key") or f"{origin['channel']}:{origin['chat_id']}"
        metadata: dict[str, Any] = {
            "injected_event": "subagent_result",
            "subagent_task_id": task_id,
        }
        if origin_message_id:
            metadata["origin_message_id"] = origin_message_id
        msg = InboundMessage(
            channel="system",
            sender_id="subagent",
            chat_id=f"{origin['channel']}:{origin['chat_id']}",
            content=announce_content,
            session_key_override=override,
            metadata=metadata,
        )

        await self.bus.publish_inbound(msg)
        logger.debug("Subagent [{}] announced result to {}:{}", task_id, origin['channel'], origin['chat_id'])

    @staticmethod
    def _format_partial_progress(result) -> str:
        completed = [e for e in result.tool_events if e["status"] == "ok"]
        failure = next((e for e in reversed(result.tool_events) if e["status"] == "error"), None)
        lines: list[str] = []
        if completed:
            lines.append("Completed steps:")
            for event in completed[-3:]:
                lines.append(f"- {event['name']}: {event['detail']}")
        if failure:
            if lines:
                lines.append("")
            lines.append("Failure:")
            lines.append(f"- {failure['name']}: {failure['detail']}")
        if result.error and not failure:
            if lines:
                lines.append("")
            lines.append("Failure:")
            lines.append(f"- {result.error}")
        return "\n".join(lines) or (result.error or "Error: subagent execution failed.")

    def _build_subagent_prompt(self) -> str:
        """Build a focused system prompt for the subagent."""
        from OriginAgent.agent.context import ContextBuilder
        from OriginAgent.agent.skills import SkillsLoader

        time_ctx = ContextBuilder._build_runtime_context(None, None)
        skills_summary = SkillsLoader(
            self.workspace,
            disabled_skills=self.disabled_skills,
        ).build_skills_summary()
        return render_template(
            "agent/subagent_system.md",
            time_ctx=time_ctx,
            workspace=str(self.workspace),
            skills_summary=skills_summary or "",
        )

    async def cancel_by_session(self, session_key: str) -> int:
        """Cancel all subagents for the given session. Returns count cancelled."""
        tasks = [self._running_tasks[tid] for tid in self._session_tasks.get(session_key, [])
                 if tid in self._running_tasks and not self._running_tasks[tid].done()]
        for t in tasks:
            t.cancel()
        if tasks:
            await asyncio.gather(*tasks, return_exceptions=True)
        return len(tasks)

    def get_running_count(self) -> int:
        """Return the number of currently running subagents."""
        return len(self._running_tasks)

    def get_running_count_by_session(self, session_key: str) -> int:
        """Return the number of currently running subagents for a session."""
        tids = self._session_tasks.get(session_key, set())
        return sum(
            1 for tid in tids
            if tid in self._running_tasks and not self._running_tasks[tid].done()
        )

    def runtime_status(self) -> dict[str, Any]:
        summary = self.records.recent_task_summary()
        summary["subagent_running_count"] = self.get_running_count()
        return summary

    @staticmethod
    def _utcnow() -> str:
        return datetime.now(timezone.utc).isoformat()

    def _record_terminal_task(
        self,
        *,
        task_id: str,
        label: str,
        task: str,
        origin: dict[str, str],
        origin_message_id: str | None,
        terminal_status: str,
        stop_reason: str | None,
        failure_summary: str | None,
        policy: SubagentPolicy,
        provider_summary: str,
    ) -> None:
        terminal_state = "completed" if terminal_status == "completed" else "failed"
        self.records.append_lifecycle(SubagentLifecycleRecord(
            subagent_id=task_id,
            root_subagent_id=origin.get("root_subagent_id"),
            parent_subagent_id=origin.get("parent_subagent_id"),
            subagent_depth=int(origin.get("subagent_depth", "1")),
            parent_session_key=origin.get("session_key"),
            state=terminal_state,  # type: ignore[arg-type]
            detail=summarize_text(stop_reason or terminal_status, max_chars=120),
        ))
        self.records.append_task(SubagentTaskRecord(
            subagent_id=task_id,
            root_subagent_id=origin.get("root_subagent_id"),
            parent_subagent_id=origin.get("parent_subagent_id"),
            subagent_depth=int(origin.get("subagent_depth", "1")),
            parent_session_key=origin.get("session_key"),
            origin_channel=origin.get("channel"),
            origin_chat_id=origin.get("chat_id"),
            origin_message_id=origin_message_id,
            task_label=label,
            task_summary=summarize_text(task),
            delegated_profile_summary=self.records.delegated_profile_summary(
                capability_snapshot=policy.capability_snapshot,
                allow_web=policy.allow_web,
                allowed_tool_names=policy.allowed_tool_names,
            ),
            allowed_tools_summary=self.records.allowed_tools_summary(policy.allowed_tool_names),
            provider_summary=provider_summary,
            isolation_mode="shared_process",
            terminal_status=terminal_status,  # type: ignore[arg-type]
            stop_reason=stop_reason,
            failure_summary=summarize_text(failure_summary, max_chars=240) if failure_summary else None,
            started_at=None,
            ended_at=self._utcnow(),
        ))

    def _configure_spawn_tool(
        self,
        *,
        tools: ToolRegistry,
        origin: dict[str, str],
        origin_message_id: str | None,
        policy: SubagentPolicy,
        status: SubagentStatus,
    ) -> None:
        tool = tools.get("spawn")
        if tool is None:
            return
        if hasattr(tool, "set_context"):
            tool.set_context(
                origin.get("channel") or "cli",
                origin.get("chat_id") or "direct",
                effective_key=origin.get("session_key"),
            )
        if hasattr(tool, "set_origin_message_id"):
            tool.set_origin_message_id(origin_message_id)
        if hasattr(tool, "set_capability_snapshot"):
            tool.set_capability_snapshot(policy.capability_snapshot)
        if hasattr(tool, "set_nested_policy"):
            tool.set_nested_policy(
                parent_subagent_id=status.task_id,
                root_subagent_id=status.root_subagent_id or status.task_id,
                subagent_depth=status.subagent_depth,
                delegated_policy=policy,
            )


class _SubagentToolObserver:
    def __init__(
        self,
        *,
        records: JsonlSubagentRecordStore,
        subagent_id: str,
        root_subagent_id: str | None,
        parent_subagent_id: str | None,
        subagent_depth: int,
        parent_session_key: str | None,
    ) -> None:
        self._records = records
        self._subagent_id = subagent_id
        self._root_subagent_id = root_subagent_id
        self._parent_subagent_id = parent_subagent_id
        self._subagent_depth = subagent_depth
        self._parent_session_key = parent_session_key

    def on_tool_result(
        self,
        *,
        name: str,
        params: dict[str, Any],
        status: str,
        start: float,
        error_kind: str | None = None,
        policy_rule: str | None = None,
        result: Any = None,
    ) -> None:
        now = datetime.now(timezone.utc)
        duration_ms = max(0, int((time.monotonic() - start) * 1000))
        started_at = (now - timedelta(milliseconds=duration_ms)).isoformat()
        mapped_status = {
            "success": "success",
            "policy_denied": "denied",
            "interrupted": "interrupted",
        }.get(status, "failed")
        self._records.append_tool(SubagentToolRecord(
            subagent_id=self._subagent_id,
            root_subagent_id=self._root_subagent_id,
            parent_subagent_id=self._parent_subagent_id,
            subagent_depth=self._subagent_depth,
            parent_session_key=self._parent_session_key,
            tool_name=name,
            status=mapped_status,  # type: ignore[arg-type]
            started_at=started_at,
            ended_at=now.isoformat(),
            duration_ms=duration_ms,
            argument_summary=summarize_payload(params),
            result_summary=summarize_payload(result) if status == "success" else summarize_text(
                policy_rule or error_kind or status,
                max_chars=200,
            ),
            policy_rule=policy_rule,
            error_kind=error_kind,
        ))


def _raise_grant_denied(policy_rule: str) -> None:
    raise PolicyDeniedError(
        _GRANT_ERROR_MESSAGE,
        code=policy_rule,
        boundary="spawn",
        policy_rule=policy_rule,
    )
