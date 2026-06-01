"""Agent loop: the core processing engine."""

from __future__ import annotations

import asyncio
import dataclasses
import os
import time
from contextlib import AsyncExitStack, nullcontext, suppress
from dataclasses import dataclass, field
from enum import Enum, auto
from pathlib import Path
from typing import TYPE_CHECKING, Any, Awaitable, Callable

from loguru import logger

from OriginAgent.agent import model_presets as preset_helpers
from OriginAgent.agent.agent_runtime_context import (
    build_bus_progress_callback,
    build_retry_wait_callback,
    runtime_chat_id,
    set_tool_context as set_tools_runtime_context,
    snapshot_for_trigger,
)
from OriginAgent.agent.agent_tool_setup import (
    build_domain_tool_context_extras,
    build_tool_context,
    register_default_tools,
    register_domain_tools,
    register_plugin_tools,
    should_register_exec,
)
from OriginAgent.agent.active_intents import ActiveIntentConfig, ActiveIntentService
from OriginAgent.agent.reminders import ReminderStore
from OriginAgent.agent.agent_turn_persist import TurnPersistManager
from OriginAgent.agent.autocompact import AutoCompact
from OriginAgent.agent.auxiliary_llm import AuxiliaryLLMRouter
from OriginAgent.agent.background_review import BackgroundReviewService
from OriginAgent.agent.context import ContextBuilder
from OriginAgent.agent.curator import CuratorService
from OriginAgent.agent.domain_packs import DomainPackManager
from OriginAgent.agent.hook import AgentHook, CompositeHook
from OriginAgent.agent.identity import ActorResolver, RuntimeContext
from OriginAgent.agent.introspection.service import RuntimeIntrospectionService
from OriginAgent.agent.memory import Consolidator, Dream
from OriginAgent.agent.progress_hook import AgentProgressHook
from OriginAgent.agent.runner import _MAX_INJECTIONS_PER_TURN, AgentRunner, AgentRunSpec
from OriginAgent.agent.subagent import SubagentManager
from OriginAgent.agent.tools.ask import (
    ask_user_options_from_messages,
    ask_user_outbound,
    ask_user_tool_result_messages,
    pending_ask_user_id,
)
from OriginAgent.agent.tools.audit import JsonlToolAuditSink, ToolAuditConfig
from OriginAgent.agent.confirmation import PendingConfirmationStore
from OriginAgent.agent.tools.file_state import FileStateStore, bind_file_states, reset_file_states
from OriginAgent.agent.tools.message import MessageTool
from OriginAgent.agent.tools.registry import ToolRegistry
from OriginAgent.agent.tools.self import MyTool
from OriginAgent.bus.events import InboundMessage, OutboundMessage
from OriginAgent.bus.queue import MessageBus
from OriginAgent.command import CommandContext, CommandRouter, register_builtin_commands
from OriginAgent.config.schema import AgentDefaults
from OriginAgent.providers.base import LLMProvider
from OriginAgent.providers.factory import ProviderSnapshot
from OriginAgent.security.capabilities import CapabilitySnapshot
from OriginAgent.security.grants import CapabilityGrantStore
from OriginAgent.session.cold_archive import SessionColdArchiveStore
from OriginAgent.session.goal_state import goal_state_ws_blob, runner_wall_llm_timeout_s
from OriginAgent.session.manager import Session, SessionManager
from OriginAgent.session.search_index import SessionSearchIndexService
from OriginAgent.utils.artifacts import generated_image_paths_from_messages
from OriginAgent.utils.document import extract_documents
from OriginAgent.utils.image_generation_intent import image_generation_prompt
from OriginAgent.utils.runtime import EMPTY_FINAL_RESPONSE_MESSAGE
from OriginAgent.utils.session_attachments import merge_turn_media_into_last_assistant
from OriginAgent.utils.webui_titles import mark_webui_session, maybe_generate_webui_title_after_turn
from OriginAgent.utils.webui_transcript import append_transcript_object, delete_webui_transcript
from OriginAgent.utils.webui_turn_helpers import publish_turn_run_status, websocket_turn_latency_ms

if TYPE_CHECKING:
    from OriginAgent.config.schema import (
        AuxiliaryConfig,
        ChannelsConfig,
        Config,
        DomainPacksConfig,
        ExecToolConfig,
        BackgroundReviewConfig,
        CuratorConfig,
        EvolutionConfig,
        ModelPresetConfig,
        ProviderConfig,
        ToolsConfig,
        WebToolsConfig,
    )
    from OriginAgent.cron.service import CronService


UNIFIED_SESSION_KEY = "unified:default"
_SENSITIVE_TOOL_LOG_PREFIXES = (
    "originagent_device_",
)
_SENSITIVE_TOOL_LOG_NAMES = {
    "exec",
    "message",
    "web_fetch",
}


class TurnState(Enum):
    RESTORE = auto()
    COMPACT = auto()
    COMMAND = auto()
    BUILD = auto()
    RUN = auto()
    SAVE = auto()
    RESPOND = auto()
    DONE = auto()


def _is_sensitive_tool_log(name: str) -> bool:
    return name in _SENSITIVE_TOOL_LOG_NAMES or any(
        name.startswith(prefix) for prefix in _SENSITIVE_TOOL_LOG_PREFIXES
    )


def _should_register_exec(config: Any) -> bool:
    return should_register_exec(config)


@dataclass
class StateTraceEntry:
    state: TurnState
    started_at: float
    duration_ms: float
    event: str
    error: str | None = None


@dataclass
class TurnContext:
    msg: InboundMessage
    session_key: str
    state: TurnState
    turn_id: str
    session: Session | None = None

    history: list[dict[str, Any]] = field(default_factory=list)
    initial_messages: list[dict[str, Any]] = field(default_factory=list)

    final_content: str | None = None
    tools_used: list[str] = field(default_factory=list)
    all_messages: list[dict[str, Any]] = field(default_factory=list)
    stop_reason: str = ""
    had_injections: bool = False

    user_persisted_early: bool = False
    save_skip: int = 0

    outbound: OutboundMessage | None = None
    generated_media: list[str] = field(default_factory=list)

    on_progress: Callable[..., Awaitable[None]] | None = None
    on_stream: Callable[[str], Awaitable[None]] | None = None
    on_stream_end: Callable[..., Awaitable[None]] | None = None
    on_retry_wait: Callable[[str], Awaitable[None]] | None = None

    pending_queue: asyncio.Queue | None = None
    pending_summary: str | None = None
    runtime_context: RuntimeContext | None = None
    capability_snapshot: CapabilitySnapshot | None = None

    trace: list[StateTraceEntry] = field(default_factory=list)


class AgentLoop:
    """
    The agent loop is the core processing engine.

    It:
    1. Receives messages from the bus
    2. Builds context with history, memory, skills
    3. Calls the LLM
    4. Executes tool calls
    5. Sends responses back
    """

    _RUNTIME_CHECKPOINT_KEY = "runtime_checkpoint"
    _PENDING_USER_TURN_KEY = "pending_user_turn"

    # Event-driven state transition table.
    # Handlers return an event string; the driver looks up the next state here.
    _TRANSITIONS: dict[tuple[TurnState, str], TurnState] = {
        (TurnState.RESTORE, "ok"): TurnState.COMPACT,
        (TurnState.COMPACT, "ok"): TurnState.COMMAND,
        (TurnState.COMMAND, "dispatch"): TurnState.BUILD,
        (TurnState.COMMAND, "shortcut"): TurnState.DONE,
        (TurnState.BUILD, "ok"): TurnState.RUN,
        (TurnState.RUN, "ok"): TurnState.SAVE,
        (TurnState.SAVE, "ok"): TurnState.RESPOND,
        (TurnState.RESPOND, "ok"): TurnState.DONE,
    }

    def __init__(
        self,
        bus: MessageBus,
        provider: LLMProvider,
        workspace: Path,
        model: str | None = None,
        max_iterations: int | None = None,
        context_window_tokens: int | None = None,
        context_block_limit: int | None = None,
        max_tool_result_chars: int | None = None,
        provider_retry_mode: str = "standard",
        tool_hint_max_length: int | None = None,
        web_config: WebToolsConfig | None = None,
        exec_config: ExecToolConfig | None = None,
        cron_service: CronService | None = None,
        restrict_to_workspace: bool = False,
        session_manager: SessionManager | None = None,
        mcp_servers: dict | None = None,
        channels_config: ChannelsConfig | None = None,
        timezone: str | None = None,
        runtime_profile: str = "default",
        session_ttl_minutes: int = 0,
        consolidation_ratio: float = 0.5,
        max_messages: int = 120,
        hooks: list[AgentHook] | None = None,
        unified_session: bool = False,
        disabled_skills: list[str] | None = None,
        tools_config: ToolsConfig | None = None,
        image_generation_provider_config: ProviderConfig | None = None,
        image_generation_provider_configs: dict[str, ProviderConfig] | None = None,
        provider_snapshot_loader: Callable[..., ProviderSnapshot] | None = None,
        provider_signature: tuple[object, ...] | None = None,
        model_presets: dict[str, ModelPresetConfig] | None = None,
        model_preset: str | None = None,
        preset_snapshot_loader: preset_helpers.PresetSnapshotLoader | None = None,
        runtime_model_publisher: Callable[[str, str | None], None] | None = None,
        device_action_executor: Any | None = None,
        device_tools_real_mode: bool = False,
        device_registry: Any | None = None,
        domain_runtime_overrides: dict[str, Any] | None = None,
        actor_resolver: ActorResolver | None = None,
        tool_audit_config: ToolAuditConfig | None = None,
        pairing_config: Any | None = None,
        auxiliary_config: "AuxiliaryConfig | None" = None,
        auxiliary_source_config: "Config | None" = None,
        auxiliary_provider_factory: Callable[["ModelPresetConfig"], LLMProvider] | None = None,
        primary_provider_name: str | None = None,
        domain_packs_config: "DomainPacksConfig | None" = None,
        domain_pack_manager: DomainPackManager | None = None,
        learning_config: "BackgroundReviewConfig | None" = None,
        learning_config_loader: Callable[[], "BackgroundReviewConfig"] | None = None,
        curator_config: "CuratorConfig | None" = None,
        curator_config_loader: Callable[[], "CuratorConfig"] | None = None,
        evolution_config: "EvolutionConfig | None" = None,
        evolution_config_loader: Callable[[], "EvolutionConfig"] | None = None,
        cold_archive_enabled: bool = True,
        tool_concurrency_limit: int | None = None,
        allow_agent_initiated_messages: bool | None = None,
        active_intent_interval_seconds: int | None = None,
        active_intent_session_cooldown_seconds: int | None = None,
        active_intent_intent_cooldown_seconds: int | None = None,
        active_intent_max_messages_per_session_per_pass: int | None = None,
    ):
        from OriginAgent.config.schema import ExecToolConfig, ToolsConfig, WebToolsConfig

        _tc = tools_config or ToolsConfig()
        defaults = AgentDefaults()
        self.bus = bus
        self.channels_config = channels_config
        self.provider = provider
        self._provider_snapshot_loader = provider_snapshot_loader
        self._preset_snapshot_loader = preset_snapshot_loader
        self._runtime_model_publisher = runtime_model_publisher
        self._provider_signature = provider_signature
        self._default_selection_signature = preset_helpers.default_selection_signature(
            provider_signature
        )
        self.workspace = workspace
        self.model = model or provider.get_default_model()
        self.auxiliary_router = AuxiliaryLLMRouter(
            primary_provider=provider,
            primary_model=self.model,
            auxiliary_config=auxiliary_config or defaults.auxiliary,
            config=auxiliary_source_config,
            provider_factory=auxiliary_provider_factory,
            primary_provider_name=primary_provider_name,
        )
        self.model_presets: dict[str, ModelPresetConfig] = model_presets or {}
        self.model_preset = model_preset or ("default" if self.model_presets else None)
        self.max_iterations = (
            max_iterations if max_iterations is not None else defaults.max_tool_iterations
        )
        self.context_window_tokens = (
            context_window_tokens
            if context_window_tokens is not None
            else defaults.context_window_tokens
        )
        self.context_block_limit = context_block_limit
        self.max_tool_result_chars = (
            max_tool_result_chars
            if max_tool_result_chars is not None
            else defaults.max_tool_result_chars
        )
        self.provider_retry_mode = provider_retry_mode
        self.tool_hint_max_length = (
            tool_hint_max_length if tool_hint_max_length is not None
            else defaults.tool_hint_max_length
        )
        self.web_config = web_config or WebToolsConfig()
        self.exec_config = exec_config or ExecToolConfig()
        self.tools_config = _tc
        self.evolution_config = evolution_config or defaults.learning.evolution
        self.session_search_index = SessionSearchIndexService(
            workspace,
            backend=_tc.session_search.backend,
            semantic_enabled=bool(_tc.session_search.semantic_enabled and _tc.session_search.enabled),
            rebuild_on_start=_tc.session_search.rebuild_on_start,
        )
        self.pairing_config = pairing_config
        self._image_generation_provider_configs = dict(image_generation_provider_configs or {})
        if (
            image_generation_provider_config is not None
            and "openrouter" not in self._image_generation_provider_configs
        ):
            self._image_generation_provider_configs["openrouter"] = image_generation_provider_config
        self.cron_service = cron_service
        self.restrict_to_workspace = restrict_to_workspace
        self._runtime_profile = runtime_profile
        self._start_time = time.time()
        self._last_usage: dict[str, int] = {}
        self._extra_hooks: list[AgentHook] = hooks or []

        self.domain_packs = domain_pack_manager or DomainPackManager(
            workspace,
            config=domain_packs_config or defaults.domain_packs,
        )
        self.background_review = BackgroundReviewService(
            workspace=workspace,
            provider=provider,
            model=self.model,
            router=self.auxiliary_router,
            config=learning_config or defaults.learning.background_review,
            config_loader=learning_config_loader,
            domain_pack_manager=self.domain_packs,
        )
        self.curator = CuratorService(
            workspace=workspace,
            config=curator_config or defaults.learning.curator,
            config_loader=curator_config_loader,
            evolution_config=self.evolution_config,
            evolution_config_loader=evolution_config_loader,
            domain_pack_manager=self.domain_packs,
        )
        self.sessions = session_manager or SessionManager(workspace)
        self.session_cold_archive = (
            SessionColdArchiveStore(workspace) if cold_archive_enabled else None
        )
        self._persist = TurnPersistManager(self.max_tool_result_chars, self.sessions)
        self._tool_audit_config = ToolAuditConfig.from_config(tool_audit_config or _tc.audit)
        self.tools = ToolRegistry(
            audit_sink=JsonlToolAuditSink(workspace),
            audit_config=self._tool_audit_config,
        )
        self._domain_runtime_overrides = dict(domain_runtime_overrides or {})
        if device_action_executor is not None:
            self._domain_runtime_overrides.setdefault("device_action_executor", device_action_executor)
        if device_registry is not None:
            self._domain_runtime_overrides.setdefault("device_registry", device_registry)
        self.actor_resolver = actor_resolver or ActorResolver()
        # One file-read/write tracker per logical session. The tool registry is
        # shared by this loop, so tools resolve the active state via contextvars.
        self._file_state_store = FileStateStore()
        self.runner = AgentRunner(provider)
        self.subagents = SubagentManager(
            provider=provider,
            workspace=workspace,
            bus=bus,
            model=self.model,
            web_config=self.web_config,
            content_read_config=_tc.content_read,
            max_tool_result_chars=self.max_tool_result_chars,
            exec_config=self.exec_config,
            restrict_to_workspace=restrict_to_workspace,
            disabled_skills=disabled_skills,
            max_iterations=self.max_iterations,
            grant_store=CapabilityGrantStore(workspace),
            preset_snapshot_loader=self._preset_snapshot_loader,
        )
        self._unified_session = unified_session
        self._max_messages = max_messages if max_messages > 0 else 120
        self._running = False
        self._mcp_servers = mcp_servers or {}
        self._mcp_stacks: dict[str, AsyncExitStack] = {}
        self._mcp_snapshot: dict[str, Any] = {}
        self._mcp_state = "disconnected"
        self._mcp_connected = False
        self._mcp_connecting = False
        self._mcp_lifecycle_lock = asyncio.Lock()
        self._mcp_ready: asyncio.Future[bool] | None = None
        self._mcp_shutdown_event: asyncio.Event | None = None
        self._mcp_runtime_task: asyncio.Task[None] | None = None
        self._active_intent_task: asyncio.Task[None] | None = None
        self._mcp_startup_error: BaseException | None = None
        self._active_tasks: dict[str, list[asyncio.Task]] = {}  # session_key -> tasks
        self._background_tasks: set[asyncio.Task] = set()
        self._session_locks: dict[str, asyncio.Lock] = {}
        # Per-session pending queues for mid-turn message injection.
        # When a session has an active task, new messages for that session
        # are routed here instead of creating a new task.
        self._pending_queues: dict[str, asyncio.Queue] = {}
        self.context = ContextBuilder(
            workspace,
            timezone=timezone,
            disabled_skills=disabled_skills,
            domain_pack_manager=self.domain_packs,
            audit_mode=self._tool_audit_config.mode,
            runtime_profile=self._runtime_profile,
            registry=self.tools,
            sessions=self.sessions,
            pending_queues=self._pending_queues,
            cron_service=self.cron_service,
            background_review_service=self.background_review,
            curator_service=self.curator,
        )
        self._confirmation_store = PendingConfirmationStore(workspace)
        self._reminder_store = ReminderStore(workspace)
        self._active_intent_config = ActiveIntentConfig(
            enabled=(
                defaults.allow_agent_initiated_messages
                if allow_agent_initiated_messages is None
                else allow_agent_initiated_messages
            ),
            interval_seconds=(
                defaults.active_intent_interval_seconds
                if active_intent_interval_seconds is None
                else active_intent_interval_seconds
            ),
            session_cooldown_seconds=(
                defaults.active_intent_session_cooldown_seconds
                if active_intent_session_cooldown_seconds is None
                else active_intent_session_cooldown_seconds
            ),
            intent_cooldown_seconds=(
                defaults.active_intent_intent_cooldown_seconds
                if active_intent_intent_cooldown_seconds is None
                else active_intent_intent_cooldown_seconds
            ),
            max_messages_per_session_per_pass=(
                defaults.active_intent_max_messages_per_session_per_pass
                if active_intent_max_messages_per_session_per_pass is None
                else active_intent_max_messages_per_session_per_pass
            ),
        )
        self.active_intents = ActiveIntentService(
            workspace=workspace,
            bus=bus,
            sessions=self.sessions,
            confirmation_store=self._confirmation_store,
            fact_store=self.context.memory.fact_store,
            reminder_store=self._reminder_store,
            config=self._active_intent_config,
        )
        self.introspection = RuntimeIntrospectionService(
            loop=self,
            workspace=workspace,
            registry=self.tools,
            sessions=self.sessions,
            pending_queues=self._pending_queues,
            cron_service=self.cron_service,
            confirmation_store=self._confirmation_store,
            reminder_store=self._reminder_store,
            audit_mode=self._tool_audit_config.mode,
            runtime_profile=self._runtime_profile,
            domain_pack_manager=self.domain_packs,
            background_review_service=self.background_review,
            curator_service=self.curator,
            session_search_index_service=self.session_search_index,
            evolution_config=self.evolution_config,
        )
        # ORIGINAGENT_MAX_CONCURRENT_REQUESTS: <=0 means unlimited; default 3.
        _max = int(os.environ.get("ORIGINAGENT_MAX_CONCURRENT_REQUESTS", "3"))
        self._concurrency_gate: asyncio.Semaphore | None = (
            asyncio.Semaphore(_max) if _max > 0 else None
        )
        self._tool_concurrency_limit = (
            tool_concurrency_limit
            if tool_concurrency_limit is not None
            else max(1, int(os.environ.get("ORIGINAGENT_MAX_CONCURRENT_TOOLS", "4")))
        )
        self.consolidator = Consolidator(
            store=self.context.memory,
            provider=provider,
            model=self.model,
            auxiliary_router=self.auxiliary_router,
            sessions=self.sessions,
            context_window_tokens=self.context_window_tokens,
            build_messages=self.context.build_messages,
            get_tool_definitions=self.tools.get_definitions,
            max_completion_tokens=provider.generation.max_tokens,
            consolidation_ratio=consolidation_ratio,
        )
        self.auto_compact = AutoCompact(
            sessions=self.sessions,
            consolidator=self.consolidator,
            session_ttl_minutes=session_ttl_minutes,
            cold_archive=self.session_cold_archive,
        )
        self.dream = Dream(
            store=self.context.memory,
            provider=provider,
            model=self.model,
            auxiliary_router=self.auxiliary_router,
            evolution_config=self.evolution_config,
        )
        self._register_default_tools()
        if _tc.my.enable:
            self.tools.register(
                MyTool(
                    loop=self,
                    modify_allowed=_tc.my.allow_set,
                    introspection_service=self.introspection,
                )
            )
        self._runtime_vars: dict[str, Any] = {}
        self._capability_snapshot: CapabilitySnapshot | None = None
        self._current_iteration: int = 0
        self.commands = CommandRouter()
        register_builtin_commands(self.commands)

    @classmethod
    def from_config(
        cls,
        config: Any,
        bus: MessageBus | None = None,
        **extra: Any,
    ) -> AgentLoop:
        """Create an AgentLoop from config with the common parameter set.

        Extra keyword arguments are forwarded to ``AgentLoop.__init__``,
        allowing callers to override or extend the standard config-derived
        parameters (e.g. ``cron_service``, ``session_manager``).
        """
        from OriginAgent.config.profiles import apply_runtime_profile
        from OriginAgent.providers.factory import make_provider

        config = apply_runtime_profile(config)
        if bus is None:
            bus = MessageBus()
        defaults = config.agents.defaults
        resolved = config.resolve_preset()
        provider = extra.pop("provider", None) or make_provider(config, resolved)
        model = extra.pop("model", None) or resolved.model
        primary_provider_name = config.get_provider_name(model) or resolved.provider
        context_window_tokens = (
            extra.pop("context_window_tokens", None)
            or resolved.context_window_tokens
            or defaults.context_window_tokens
        )
        model_presets = preset_helpers.configured_model_presets(config)
        provider_snapshot_loader = extra.get("provider_snapshot_loader")
        preset_snapshot_loader = extra.pop(
            "preset_snapshot_loader",
            preset_helpers.make_preset_snapshot_loader(config, provider_snapshot_loader),
        )
        domain_pack_manager = extra.get("domain_pack_manager") or DomainPackManager(
            config.workspace_path,
            config=defaults.domain_packs,
        )
        extra["domain_pack_manager"] = domain_pack_manager
        domain_runtime_overrides = dict(extra.pop("domain_runtime_overrides", {}) or {})
        if "device_action_executor" in extra:
            domain_runtime_overrides.setdefault(
                "device_action_executor",
                extra.pop("device_action_executor"),
            )
        if "device_registry" in extra:
            domain_runtime_overrides.setdefault("device_registry", extra.pop("device_registry"))

        def _background_review_config_loader():
            from OriginAgent.config.loader import load_config

            return load_config().agents.defaults.learning.background_review

        def _curator_config_loader():
            from OriginAgent.config.loader import load_config

            return load_config().agents.defaults.learning.curator

        def _evolution_config_loader():
            from OriginAgent.config.loader import load_config

            return load_config().agents.defaults.learning.evolution

        return cls(
            bus=bus,
            provider=provider,
            workspace=config.workspace_path,
            model=model,
            max_iterations=defaults.max_tool_iterations,
            context_window_tokens=context_window_tokens,
            context_block_limit=defaults.context_block_limit,
            max_tool_result_chars=defaults.max_tool_result_chars,
            provider_retry_mode=defaults.provider_retry_mode,
            tool_hint_max_length=defaults.tool_hint_max_length,
            web_config=config.tools.web,
            exec_config=config.tools.exec,
            restrict_to_workspace=config.tools.restrict_to_workspace,
            mcp_servers=config.tools.mcp_servers,
            channels_config=config.channels,
            timezone=defaults.timezone,
            runtime_profile=config.runtime.profile,
            unified_session=defaults.unified_session,
            disabled_skills=defaults.disabled_skills,
            session_ttl_minutes=defaults.session_ttl_minutes,
            cold_archive_enabled=defaults.cold_archive_enabled,
            consolidation_ratio=defaults.consolidation_ratio,
            max_messages=defaults.max_messages,
            tools_config=config.tools,
            model_presets=model_presets,
            model_preset=defaults.model_preset or "default",
            preset_snapshot_loader=preset_snapshot_loader,
            domain_runtime_overrides=domain_runtime_overrides,
            tool_audit_config=config.tools.audit,
            pairing_config=config.security.pairing,
            auxiliary_config=defaults.auxiliary,
            auxiliary_source_config=config,
            primary_provider_name=primary_provider_name,
            domain_packs_config=defaults.domain_packs,
            learning_config=defaults.learning.background_review,
            learning_config_loader=_background_review_config_loader,
            curator_config=defaults.learning.curator,
            curator_config_loader=_curator_config_loader,
            evolution_config=defaults.learning.evolution,
            evolution_config_loader=_evolution_config_loader,
            **extra,
        )

    def _sync_subagent_runtime_limits(self) -> None:
        """Keep subagent runtime limits aligned with mutable loop settings."""
        self.subagents.max_iterations = self.max_iterations

    def _archive_session_file_cap(
        self,
        messages: list[dict[str, Any]],
        *,
        session_key: str,
        reason: str,
    ) -> None:
        if self.session_cold_archive is not None:
            self.session_cold_archive.archive(session_key, messages, reason=reason)
        self.context.memory.raw_archive(messages)

    def _apply_provider_snapshot(self, snapshot: ProviderSnapshot) -> None:
        """Swap model/provider for future turns without disturbing an active one."""
        provider = snapshot.provider
        model = snapshot.model
        context_window_tokens = snapshot.context_window_tokens
        if self.provider is provider and self.model == model:
            return
        old_model = self.model
        self.provider = provider
        self.model = model
        self.context_window_tokens = context_window_tokens
        self.runner.provider = provider
        self.subagents.set_provider(provider, model)
        self.auxiliary_router.set_primary(provider, model)
        self.background_review.set_provider(provider, model)
        self.consolidator.set_provider(provider, model, context_window_tokens)
        self.dream.set_provider(provider, model)
        self._provider_signature = snapshot.signature
        self._default_selection_signature = preset_helpers.default_selection_signature(
            snapshot.signature
        )
        logger.info("Runtime model switched for next turn: {} -> {}", old_model, model)
        if self._runtime_model_publisher:
            self._runtime_model_publisher(model, self.model_preset)

    def _refresh_provider_snapshot(self) -> None:
        if self.model_preset and self.model_preset != "default":
            if self._preset_snapshot_loader is None:
                return
            try:
                snapshot = self._preset_snapshot_loader(self.model_preset)
            except Exception:
                logger.exception("Failed to refresh model preset config")
                return
            if snapshot.signature == self._provider_signature:
                return
            self._apply_provider_snapshot(snapshot)
            return

        if self._provider_snapshot_loader is None:
            return
        try:
            snapshot = self._provider_snapshot_loader()
        except Exception:
            logger.exception("Failed to refresh provider config")
            return
        if snapshot.signature == self._provider_signature:
            return
        self.model_preset = "default"
        self._apply_provider_snapshot(snapshot)

    def set_model_preset(self, name: str) -> None:
        """Switch the active runtime model preset for subsequent turns."""
        preset_name = preset_helpers.normalize_preset_name(name, self.model_presets)
        snapshot = preset_helpers.build_runtime_preset_snapshot(
            name=preset_name,
            presets=self.model_presets,
            provider=self.provider,
            loader=self._preset_snapshot_loader,
        )
        self.model_preset = preset_name
        self._apply_provider_snapshot(snapshot)

    def _register_default_tools(self) -> None:
        """Register the default set of tools."""
        register_default_tools(
            self.tools,
            workspace=self.workspace,
            bus=self.bus,
            config=self.tools_config,
            web_config=self.web_config,
            exec_config=self.exec_config,
            restrict_to_workspace=self.restrict_to_workspace,
            sessions=self.sessions,
            pending_queues=self._pending_queues,
            cron_service=self.cron_service,
            audit_config=self._tool_audit_config,
            domain_pack_manager=self.domain_packs,
            background_review_service=self.background_review,
            curator_service=self.curator,
            session_search_index_service=self.session_search_index,
            subagent_manager=self.subagents,
            file_state_store=self._file_state_store,
            provider_snapshot_loader=self._provider_snapshot_loader,
            image_generation_provider_configs=self._image_generation_provider_configs,
            timezone=self.context.timezone or "UTC",
            runtime_profile=self._runtime_profile,
            introspection_service=self.introspection,
            confirmation_store=self._confirmation_store,
            domain_runtime_overrides=self._domain_runtime_overrides,
            evolution_config=self.evolution_config,
        )

    def _build_tool_context(self):
        return build_tool_context(
            config=self.tools_config,
            workspace=self.workspace,
            bus=self.bus,
            subagent_manager=self.subagents,
            cron_service=self.cron_service,
            sessions=self.sessions,
            file_state_store=self._file_state_store,
            provider_snapshot_loader=self._provider_snapshot_loader,
            image_generation_provider_configs=self._image_generation_provider_configs,
            timezone=self.context.timezone or "UTC",
            audit_config=self._tool_audit_config,
            context_extras=build_domain_tool_context_extras(
                domain_pack_manager=self.domain_packs,
                workspace=self.workspace,
                config=self.tools_config,
                overrides=self._domain_runtime_overrides,
            ),
        )

    def _register_domain_tools(self) -> None:
        """Load tools declared by active domain packs without replacing core tools."""
        register_domain_tools(
            self.tools,
            domain_pack_manager=self.domain_packs,
            context=self._build_tool_context(),
        )

    def _register_plugin_tools(self) -> None:
        """Load external OriginAgent tool plugins without replacing core tools."""
        register_plugin_tools(self.tools, context=self._build_tool_context())

    async def _connect_mcp(self) -> None:
        """Connect to configured MCP servers (one-time, lazy)."""
        if not self._mcp_servers:
            return
        while True:
            ready: asyncio.Future[bool] | None = None
            runtime_task: asyncio.Task[None] | None = None
            async with self._mcp_lifecycle_lock:
                if self._mcp_state == "connected":
                    return
                if self._mcp_state == "connecting":
                    ready = self._mcp_ready
                elif self._mcp_state == "closing":
                    runtime_task = self._mcp_runtime_task
                else:
                    ready = asyncio.get_running_loop().create_future()
                    self._mcp_state = "connecting"
                    self._mcp_connected = False
                    self._mcp_connecting = True
                    self._mcp_startup_error = None
                    self._mcp_ready = ready
                    self._mcp_shutdown_event = asyncio.Event()
                    self._mcp_runtime_task = asyncio.create_task(
                        self._run_mcp_runtime(ready, self._mcp_shutdown_event),
                        name="originagent-mcp-runtime",
                    )
            if runtime_task is not None:
                with suppress(Exception):
                    await asyncio.shield(runtime_task)
                continue
            if ready is not None:
                with suppress(Exception):
                    await asyncio.shield(ready)
                return
            return

    async def _run_mcp_runtime(
        self,
        ready: asyncio.Future[bool],
        shutdown_event: asyncio.Event,
    ) -> None:
        """Own the MCP connection lifecycle inside a single task."""
        from OriginAgent.agent.tools.mcp import connect_mcp_servers

        stacks: dict[str, AsyncExitStack] = {}
        clear_snapshot_on_exit = False
        try:
            stacks = await connect_mcp_servers(
                self._mcp_servers,
                self.tools,
                snapshot_out=self._mcp_snapshot,
            )
            if not stacks:
                logger.warning("No MCP servers connected successfully (will retry next message)")
                async with self._mcp_lifecycle_lock:
                    self._mcp_stacks = {}
                    self._mcp_connected = False
                    self._mcp_connecting = False
                    self._mcp_state = "disconnected"
                    if not ready.done():
                        ready.set_result(False)
                return

            async with self._mcp_lifecycle_lock:
                self._mcp_stacks = stacks
                self._mcp_connected = True
                self._mcp_connecting = False
                self._mcp_state = "connected"
                self._mcp_startup_error = None
                if not ready.done():
                    ready.set_result(True)

            await shutdown_event.wait()
            clear_snapshot_on_exit = True
        except asyncio.CancelledError:
            clear_snapshot_on_exit = True
            logger.warning("MCP runtime cancelled (will retry next message)")
            async with self._mcp_lifecycle_lock:
                self._mcp_stacks.clear()
                self._mcp_snapshot.clear()
                self._mcp_connected = False
                self._mcp_connecting = False
                self._mcp_state = "disconnected"
                if not ready.done():
                    ready.set_result(False)
            raise
        except BaseException as e:
            clear_snapshot_on_exit = True
            logger.warning("Failed to connect MCP servers (will retry next message): {}", e)
            async with self._mcp_lifecycle_lock:
                self._mcp_stacks.clear()
                self._mcp_snapshot.clear()
                self._mcp_connected = False
                self._mcp_connecting = False
                self._mcp_state = "disconnected"
                self._mcp_startup_error = e
                if not ready.done():
                    ready.set_result(False)
            return
        finally:
            for name, stack in stacks.items():
                try:
                    await stack.aclose()
                except (RuntimeError, BaseExceptionGroup):
                    logger.debug("MCP server '{}' cleanup error (can be ignored)", name)
            async with self._mcp_lifecycle_lock:
                if self._mcp_runtime_task is asyncio.current_task():
                    self._mcp_stacks.clear()
                    if clear_snapshot_on_exit:
                        self._mcp_snapshot.clear()
                    self._mcp_connected = False
                    self._mcp_connecting = False
                    self._mcp_state = "disconnected"
                    self._mcp_runtime_task = None
                    self._mcp_ready = None
                    self._mcp_shutdown_event = None

    def _set_tool_context(
        self, channel: str, chat_id: str,
        message_id: str | None = None, metadata: dict | None = None,
        session_key: str | None = None,
        actor_id: str | None = None,
        trigger: str | None = None,
        capability_snapshot: CapabilitySnapshot | None = None,
        runtime_context: RuntimeContext | None = None,
    ) -> None:
        """Update context for all tools that need routing info."""
        snapshot = capability_snapshot or self._capability_snapshot
        set_tools_runtime_context(
            self.tools,
            channel=channel,
            chat_id=chat_id,
            message_id=message_id,
            metadata=metadata,
            session_key=session_key,
            actor_id=actor_id,
            trigger=trigger,
            capability_snapshot=snapshot,
            runtime_context=runtime_context,
            unified_session=self._unified_session,
            unified_session_key=UNIFIED_SESSION_KEY,
        )

    @staticmethod
    def _strip_think(text: str | None) -> str | None:
        """Remove <think>…</think> blocks that some models embed in content."""
        if not text:
            return None
        from OriginAgent.utils.helpers import strip_think

        return strip_think(text) or None

    @staticmethod
    def _runtime_chat_id(msg: InboundMessage) -> str:
        """Return the chat id shown in runtime metadata for the model."""
        return runtime_chat_id(msg)

    @staticmethod
    def _snapshot_for_trigger(trigger: str | None) -> CapabilitySnapshot:
        return snapshot_for_trigger(trigger)

    def _resolve_runtime_context(
        self,
        msg: InboundMessage,
        *,
        channel: str | None = None,
        chat_id: str | None = None,
        session_key: str | None = None,
    ) -> RuntimeContext:
        return self.actor_resolver.resolve_runtime_context(
            channel=msg.channel,
            chat_id=msg.chat_id,
            sender_id=msg.sender_id,
            metadata=msg.metadata or {},
            session_key=session_key,
            routing_channel=channel,
            routing_chat_id=chat_id,
        )

    def _tool_hint(self, tool_calls: list) -> str:
        """Format tool calls as concise hints with smart abbreviation."""
        from OriginAgent.utils.tool_hints import format_tool_hints

        return format_tool_hints(tool_calls, max_length=self.tool_hint_max_length)

    async def _build_bus_progress_callback(
        self, msg: InboundMessage
    ) -> Callable[..., Awaitable[None]]:
        """Build a progress callback that publishes to the message bus."""
        return await build_bus_progress_callback(self.bus, msg)

    async def _build_retry_wait_callback(
        self, msg: InboundMessage
    ) -> Callable[[str], Awaitable[None]]:
        """Build a retry-wait callback that publishes to the message bus."""
        return await build_retry_wait_callback(self.bus, msg)

    def _turn_persist_manager(self) -> TurnPersistManager:
        manager = getattr(self, "_persist", None)
        if isinstance(manager, TurnPersistManager):
            return manager
        max_chars = getattr(self, "max_tool_result_chars", AgentDefaults().max_tool_result_chars)
        manager = TurnPersistManager(max_chars, getattr(self, "sessions", None))
        try:
            self._persist = manager
        except Exception:
            pass
        return manager

    def _persist_user_message_early(
        self,
        msg: InboundMessage,
        session: Session,
        pending_ask_id: str | None,
        **kwargs: Any,
    ) -> bool:
        """Persist the triggering user message before the turn starts.

        Returns True if the message was persisted.
        """
        return AgentLoop._turn_persist_manager(self).persist_user_message_early(
            msg,
            session,
            pending_ask_id,
            **kwargs,
        )

    def _build_initial_messages(
        self,
        msg: InboundMessage,
        session: Session,
        history: list[dict[str, Any]],
        pending_ask_id: str | None,
        pending_summary: str | None,
    ) -> list[dict[str, Any]]:
        """Build the initial message list for the LLM turn."""
        if pending_ask_id:
            system_prompt = self.context.build_system_prompt(
                channel=msg.channel,
                session_summary=pending_summary,
            )
            messages = ask_user_tool_result_messages(
                system_prompt,
                history,
                pending_ask_id,
                image_generation_prompt(msg.content, msg.metadata),
            )
            messages.append({
                "role": "user",
                "content": [
                    self.context.build_runtime_context_block(
                        msg.channel,
                        self._runtime_chat_id(msg),
                        self.context.timezone,
                        sender_id=msg.sender_id,
                        session_metadata=session.metadata,
                    ),
                    *self.context.build_reference_context_blocks(
                        session_summary=pending_summary,
                    ),
                ],
            })
            return messages
        return self.context.build_messages(
            history=history,
            current_message=image_generation_prompt(msg.content, msg.metadata),
            media=msg.media if msg.media else None,
            channel=msg.channel,
            chat_id=self._runtime_chat_id(msg),
            sender_id=msg.sender_id,
            session_summary=pending_summary,
            session_metadata=session.metadata,
        )

    def _is_webui_message(self, msg: InboundMessage) -> bool:
        return msg.channel == "websocket" and msg.metadata.get("webui") is True

    def _append_webui_command_transcript(
        self,
        msg: InboundMessage,
        content: str,
    ) -> None:
        """Persist a command response to the WebUI transcript."""
        if not self._is_webui_message(msg):
            return
        try:
            append_transcript_object(
                f"websocket:{msg.chat_id}",
                {"event": "message", "chat_id": msg.chat_id, "text": content},
            )
        except (TypeError, ValueError, OSError) as e:
            logger.warning("webui command transcript append failed: {}", e)

    def _persist_shortcut_command_turn(
        self,
        msg: InboundMessage,
        session_key: str,
        result: OutboundMessage,
    ) -> None:
        """Persist slash-command turns that bypass the normal RUN/SAVE states."""
        raw = msg.content.strip()
        if raw.lower() == "/new":
            if self._is_webui_message(msg):
                delete_webui_transcript(session_key)
            return
        session = self.sessions.get_or_create(session_key)
        mark_webui_session(session, msg.metadata)
        self._persist_user_message_early(
            msg,
            session,
            pending_ask_id=None,
            _command=True,
        )
        if result.content.strip():
            session.add_message("assistant", result.content, _command=True)
        self._clear_pending_user_turn(session)
        self.sessions.save(session)
        self._append_webui_command_transcript(msg, result.content)

    async def _dispatch_command_inline(
        self,
        msg: InboundMessage,
        key: str,
        raw: str,
        dispatch_fn: Callable[[CommandContext], Awaitable[OutboundMessage | None]],
    ) -> None:
        """Dispatch a command directly from the run() loop and publish the result."""
        session = self.sessions.get_or_create(key)
        ctx = CommandContext(msg=msg, session=session, key=key, raw=raw, loop=self)
        result = await dispatch_fn(ctx)
        if result:
            self._persist_shortcut_command_turn(msg, key, result)
            if self._is_webui_message(msg):
                result.metadata["_webui_transcript_recorded"] = True
            await self.bus.publish_outbound(result)
            if msg.channel == "websocket":
                await self.bus.publish_outbound(
                    OutboundMessage(
                        channel=msg.channel,
                        chat_id=msg.chat_id,
                        content="",
                        metadata={
                            **dict(msg.metadata or {}),
                            "_turn_end": True,
                            "goal_state": goal_state_ws_blob(
                                self.sessions.get_or_create(key).metadata
                            ),
                        },
                    )
                )
        else:
            logger.warning("Command '{}' matched but dispatch returned None", raw)

    async def _cancel_active_tasks(self, key: str) -> int:
        """Cancel and await all active tasks and subagents for *key*.

        Returns the total number of cancelled tasks + subagents.
        """
        tasks = self._active_tasks.pop(key, [])
        cancelled = sum(1 for t in tasks if not t.done() and t.cancel())
        for t in tasks:
            with suppress(asyncio.CancelledError, Exception):
                await t
        sub_cancelled = await self.subagents.cancel_by_session(key)
        return cancelled + sub_cancelled

    def _effective_session_key(self, msg: InboundMessage) -> str:
        """Return the session key used for task routing and mid-turn injections."""
        if self._unified_session and not msg.session_key_override:
            return UNIFIED_SESSION_KEY
        return msg.session_key

    def _replay_token_budget(self) -> int:
        """Derive a token budget for session history replay from the context window."""
        if self.context_window_tokens <= 0:
            return 0
        max_output = getattr(getattr(self.provider, "generation", None), "max_tokens", 4096)
        try:
            reserved_output = int(max_output)
        except (TypeError, ValueError):
            reserved_output = 4096
        budget = self.context_window_tokens - max(1, reserved_output) - 1024
        return budget if budget > 0 else max(128, self.context_window_tokens // 2)

    async def _run_agent_loop(
        self,
        initial_messages: list[dict],
        on_progress: Callable[..., Awaitable[None]] | None = None,
        on_stream: Callable[[str], Awaitable[None]] | None = None,
        on_stream_end: Callable[..., Awaitable[None]] | None = None,
        on_retry_wait: Callable[[str], Awaitable[None]] | None = None,
        *,
        session: Session | None = None,
        channel: str = "cli",
        chat_id: str = "direct",
        message_id: str | None = None,
        metadata: dict[str, Any] | None = None,
        session_key: str | None = None,
        pending_queue: asyncio.Queue | None = None,
        actor_id: str | None = None,
        trigger: str | None = None,
        capability_snapshot: CapabilitySnapshot | None = None,
    ) -> tuple[str | None, list[str], list[dict], str, bool]:
        """Run the agent iteration loop.

        *on_stream*: called with each content delta during streaming.
        *on_stream_end(resuming)*: called when a streaming session finishes.
        ``resuming=True`` means tool calls follow (spinner should restart);
        ``resuming=False`` means this is the final response.

        Returns (final_content, tools_used, messages, stop_reason, had_injections).
        """
        self._sync_subagent_runtime_limits()
        self._capability_snapshot = capability_snapshot or self._snapshot_for_trigger(trigger)
        if hasattr(self.tools, "set_capability_snapshot"):
            self.tools.set_capability_snapshot(self._capability_snapshot)

        loop_hook = AgentProgressHook(
            on_progress=on_progress,
            on_stream=on_stream,
            on_stream_end=on_stream_end,
            channel=channel,
            chat_id=chat_id,
            message_id=message_id,
            metadata=metadata,
            session_key=session_key,
            tool_hint_max_length=self.tool_hint_max_length,
            set_tool_context=self._set_tool_context,
            on_iteration=lambda iteration: setattr(self, "_current_iteration", iteration),
            actor_id=actor_id,
            trigger=trigger,
            capability_snapshot=self._capability_snapshot,
            sensitive_tool_log_names=_SENSITIVE_TOOL_LOG_NAMES,
            sensitive_tool_log_prefixes=_SENSITIVE_TOOL_LOG_PREFIXES,
        )
        hook: AgentHook = (
            CompositeHook([loop_hook] + self._extra_hooks) if self._extra_hooks else loop_hook
        )

        async def _checkpoint(payload: dict[str, Any]) -> None:
            if session is None:
                return
            self._set_runtime_checkpoint(session, payload)

        async def _drain_pending(*, limit: int = _MAX_INJECTIONS_PER_TURN) -> list[dict[str, Any]]:
            """Drain follow-up messages from the pending queue.

            When no messages are immediately available but sub-agents
            spawned in this dispatch are still running, blocks until at
            least one result arrives (or timeout).  This keeps the runner
            loop alive so subsequent sub-agent completions are consumed
            in-order rather than dispatched separately.
            """
            if pending_queue is None:
                return []

            def _to_user_message(pending_msg: InboundMessage) -> dict[str, Any]:
                content = pending_msg.content
                media = pending_msg.media if pending_msg.media else None
                if media:
                    content, media = extract_documents(content, media)
                    media = media or None
                runtime_block = self.context.build_runtime_context_block(
                    pending_msg.channel,
                    self._runtime_chat_id(pending_msg),
                    self.context.timezone,
                )
                if (
                    pending_msg.sender_id == "subagent"
                    or pending_msg.metadata.get("injected_event") == "subagent_result"
                ):
                    merged: list[dict[str, Any]] = [
                        runtime_block,
                        self.context.build_internal_event_block("subagent_result", content),
                    ]
                elif pending_msg.metadata.get("injected_event") == "active_intent":
                    merged = [
                        runtime_block,
                        self.context.build_internal_event_block("active_intent", content),
                    ]
                else:
                    merged = [
                        runtime_block,
                        *self.context._build_user_content(content, media),
                    ]
                return {"role": "user", "content": merged}

            items: list[dict[str, Any]] = []
            while len(items) < limit:
                try:
                    items.append(_to_user_message(pending_queue.get_nowait()))
                except asyncio.QueueEmpty:
                    break

            # Block if nothing drained but sub-agents spawned in this dispatch
            # are still running.  Keeps the runner loop alive so subsequent
            # completions are injected in-order rather than dispatched separately.
            if (not items
                    and session is not None
                    and self.subagents.get_running_count_by_session(session.key) > 0):
                try:
                    msg = await asyncio.wait_for(pending_queue.get(), timeout=300)
                except asyncio.TimeoutError:
                    logger.warning(
                        "Timeout waiting for sub-agent completion in session {}",
                        session.key,
                    )
                    return items
                items.append(_to_user_message(msg))
                while len(items) < limit:
                    try:
                        items.append(_to_user_message(pending_queue.get_nowait()))
                    except asyncio.QueueEmpty:
                        break

            return items

        active_session_key = session.key if session else session_key
        file_state_token = bind_file_states(self._file_state_store.for_session(active_session_key))
        try:
            result = await self.runner.run(AgentRunSpec(
                initial_messages=initial_messages,
                tools=self.tools,
                model=self.model,
                max_iterations=self.max_iterations,
                max_tool_result_chars=self.max_tool_result_chars,
                hook=hook,
                error_message="Sorry, I encountered an error calling the AI model.",
                concurrent_tools=True,
                tool_concurrency_limit=self._tool_concurrency_limit,
                workspace=self.workspace,
                session_key=session.key if session else None,
                context_window_tokens=self.context_window_tokens,
                context_block_limit=self.context_block_limit,
                provider_retry_mode=self.provider_retry_mode,
                progress_callback=on_progress,
                stream_progress_deltas=on_stream is not None,
                retry_wait_callback=on_retry_wait,
                checkpoint_callback=_checkpoint,
                injection_callback=_drain_pending,
                llm_timeout_s=runner_wall_llm_timeout_s(
                    self.sessions,
                    session_key,
                    metadata=session.metadata if session is not None else None,
                ),
            ))
        finally:
            reset_file_states(file_state_token)
        self._last_usage = result.usage
        if result.stop_reason == "max_iterations":
            logger.warning("Max iterations ({}) reached", self.max_iterations)
            # Push final content through stream so streaming channels (e.g. Feishu)
            # update the card instead of leaving it empty.
            if on_stream and on_stream_end:
                await on_stream(result.final_content or "")
                await on_stream_end(resuming=False)
        elif result.stop_reason == "error":
            logger.error("LLM returned error: {}", (result.final_content or "")[:200])
        return result.final_content, result.tools_used, result.messages, result.stop_reason, result.had_injections

    async def run(self) -> None:
        """Run the agent loop, dispatching messages as tasks to stay responsive to /stop."""
        self._running = True
        await self._connect_mcp()
        self._schedule_session_search_refresh(force=self.session_search_index.rebuild_on_start)
        self._start_active_intent_loop()
        logger.info("Agent loop started")

        while self._running:
            try:
                msg = await asyncio.wait_for(self.bus.consume_inbound(), timeout=1.0)
            except asyncio.TimeoutError:
                self.auto_compact.check_expired(
                    self._schedule_background,
                    active_session_keys=self._pending_queues.keys(),
                )
                continue
            except asyncio.CancelledError:
                # Preserve real task cancellation so shutdown can complete cleanly.
                # Only ignore non-task CancelledError signals that may leak from integrations.
                if not self._running or asyncio.current_task().cancelling():
                    raise
                continue
            except Exception as e:
                logger.warning("Error consuming inbound message: {}, continuing...", e)
                continue

            raw = msg.content.strip()
            effective_key = self._effective_session_key(msg)
            if self.commands.is_priority(raw):
                await self._dispatch_command_inline(
                    msg, effective_key, raw,
                    self.commands.dispatch_priority,
                )
                continue
            # If this session already has an active pending queue (i.e. a task
            # is processing this session), route the message there for mid-turn
            # injection instead of creating a competing task.
            if effective_key in self._pending_queues:
                # Non-priority commands must not be queued for injection;
                # dispatch them directly (same pattern as priority commands).
                if self.commands.is_dispatchable_command(raw):
                    await self._dispatch_command_inline(
                        msg, effective_key, raw,
                        self.commands.dispatch,
                    )
                    continue
                pending_msg = msg
                if effective_key != msg.session_key:
                    pending_msg = dataclasses.replace(
                        msg,
                        session_key_override=effective_key,
                    )
                try:
                    self._pending_queues[effective_key].put_nowait(pending_msg)
                except asyncio.QueueFull:
                    logger.warning(
                        "Pending queue full for session {}, falling back to queued task",
                        effective_key,
                    )
                else:
                    logger.info(
                        "Routed follow-up message to pending queue for session {}",
                        effective_key,
                    )
                    continue
            # Compute the effective session key before dispatching
            # This ensures /stop command can find tasks correctly when unified session is enabled
            task = asyncio.create_task(self._dispatch(msg))
            self._active_tasks.setdefault(effective_key, []).append(task)
            task.add_done_callback(
                lambda t, k=effective_key: self._active_tasks.get(k, [])
                and self._active_tasks[k].remove(t)
                if t in self._active_tasks.get(k, [])
                else None
            )

    async def _dispatch(self, msg: InboundMessage) -> None:
        """Process a message: per-session serial, cross-session concurrent."""
        session_key = self._effective_session_key(msg)
        if session_key != msg.session_key:
            msg = dataclasses.replace(msg, session_key_override=session_key)
        lock = self._session_locks.setdefault(session_key, asyncio.Lock())
        gate = self._concurrency_gate or nullcontext()

        # Register a pending queue so follow-up messages for this session are
        # routed here (mid-turn injection) instead of spawning a new task.
        pending = asyncio.Queue(maxsize=20)
        self._pending_queues[session_key] = pending

        try:
            async with lock, gate:
                try:
                    on_stream = on_stream_end = None
                    if msg.metadata.get("_wants_stream"):
                        # Split one answer into distinct stream segments.
                        stream_base_id = f"{msg.session_key}:{time.time_ns()}"
                        stream_segment = 0

                        def _current_stream_id() -> str:
                            return f"{stream_base_id}:{stream_segment}"

                        async def on_stream(delta: str) -> None:
                            meta = dict(msg.metadata or {})
                            meta["_stream_delta"] = True
                            meta["_stream_id"] = _current_stream_id()
                            await self.bus.publish_outbound(OutboundMessage(
                                channel=msg.channel, chat_id=msg.chat_id,
                                content=delta,
                                metadata=meta,
                            ))

                        async def on_stream_end(*, resuming: bool = False) -> None:
                            nonlocal stream_segment
                            meta = dict(msg.metadata or {})
                            meta["_stream_end"] = True
                            meta["_resuming"] = resuming
                            meta["_stream_id"] = _current_stream_id()
                            await self.bus.publish_outbound(OutboundMessage(
                                channel=msg.channel, chat_id=msg.chat_id,
                                content="",
                                metadata=meta,
                            ))
                            stream_segment += 1

                    response = await self._process_message(
                        msg, on_stream=on_stream, on_stream_end=on_stream_end,
                        pending_queue=pending,
                    )
                    if response is not None:
                        await self.bus.publish_outbound(response)
                    elif msg.channel == "cli":
                        await self.bus.publish_outbound(OutboundMessage(
                            channel=msg.channel, chat_id=msg.chat_id,
                            content="", metadata=msg.metadata or {},
                        ))
                    if msg.channel == "websocket":
                        # Signal that the turn is fully complete (all tools executed,
                        # final text streamed).  This lets WS clients know when to
                        # definitively stop the loading indicator.
                        await self.bus.publish_outbound(OutboundMessage(
                            channel=msg.channel, chat_id=msg.chat_id,
                            content="",
                            metadata={
                                **msg.metadata,
                                "_turn_end": True,
                                "latency_ms": msg.metadata.get("webui_turn_latency_ms"),
                                "goal_state": goal_state_ws_blob(
                                    self.sessions.get_or_create(session_key).metadata
                                ),
                            },
                        ))
                        if msg.metadata.get("webui") is True:
                            async def _generate_title_and_notify() -> None:
                                generated = await maybe_generate_webui_title_after_turn(
                                    channel=msg.channel,
                                    metadata=msg.metadata,
                                    sessions=self.sessions,
                                    session_key=session_key,
                                    provider=self.provider,
                                    model=self.model,
                                )
                                if generated:
                                    await self.bus.publish_outbound(OutboundMessage(
                                        channel=msg.channel,
                                        chat_id=msg.chat_id,
                                        content="",
                                        metadata={**msg.metadata, "_session_updated": True},
                                    ))

                            self._schedule_background(_generate_title_and_notify())
                except asyncio.CancelledError:
                    logger.info("Task cancelled for session {}", session_key)
                    # Preserve partial context from the interrupted turn so
                    # the user does not lose tool results and assistant
                    # messages accumulated before /stop.  The checkpoint was
                    # already persisted to session metadata by
                    # _emit_checkpoint during tool execution; materializing
                    # it into session history now makes it visible in the
                    # next conversation turn.
                    try:
                        key = self._effective_session_key(msg)
                        session = self.sessions.get_or_create(key)
                        if self._restore_runtime_checkpoint(session):
                            self._clear_pending_user_turn(session)
                            self.sessions.save(session)
                            logger.info(
                                "Restored partial context for cancelled session {}",
                                key,
                            )
                    except Exception:
                        logger.debug(
                            "Could not restore checkpoint for cancelled session {}",
                            session_key,
                            exc_info=True,
                        )
                    raise
                except Exception:
                    logger.exception("Error processing message for session {}", session_key)
                    await self.bus.publish_outbound(OutboundMessage(
                        channel=msg.channel, chat_id=msg.chat_id,
                        content="Sorry, I encountered an error.",
                    ))
        finally:
            # Drain any messages still in the pending queue and re-publish
            # them to the bus so they are processed as fresh inbound messages
            # rather than silently lost.
            queue = self._pending_queues.pop(session_key, None)
            if queue is not None:
                leftover = 0
                while True:
                    try:
                        item = queue.get_nowait()
                    except asyncio.QueueEmpty:
                        break
                    await self.bus.publish_inbound(item)
                    leftover += 1
                if leftover:
                    logger.info(
                        "Re-published {} leftover message(s) to bus for session {}",
                        leftover, session_key,
                    )

    async def close_mcp(self) -> None:
        """Drain pending background archives, then close MCP connections."""
        if self._active_intent_task is not None:
            self._active_intent_task.cancel()
            with suppress(asyncio.CancelledError):
                await asyncio.shield(self._active_intent_task)
            self._active_intent_task = None
        if self._background_tasks:
            await asyncio.gather(*self._background_tasks, return_exceptions=True)
            self._background_tasks.clear()
        runtime_task: asyncio.Task[None] | None = None
        async with self._mcp_lifecycle_lock:
            runtime_task = self._mcp_runtime_task
            shutdown_event = self._mcp_shutdown_event
            if runtime_task is None:
                self._mcp_connected = False
                self._mcp_connecting = False
                self._mcp_state = "disconnected"
                self._mcp_stacks.clear()
                self._mcp_snapshot.clear()
                return
            self._mcp_connected = False
            self._mcp_connecting = False
            self._mcp_state = "closing"
            if shutdown_event is not None:
                shutdown_event.set()
        with suppress(Exception):
            await asyncio.shield(runtime_task)

    def _schedule_background(self, coro) -> None:
        """Schedule a coroutine as a tracked background task (drained on shutdown)."""
        task = asyncio.create_task(coro)
        self._background_tasks.add(task)
        task.add_done_callback(self._background_tasks.discard)

    def _start_active_intent_loop(self) -> None:
        if not self._active_intent_config.enabled or self._active_intent_task is not None:
            return
        self._active_intent_task = asyncio.create_task(self._active_intent_loop())

    async def _active_intent_loop(self) -> None:
        interval = max(1, int(self._active_intent_config.interval_seconds))
        try:
            while self._running:
                await asyncio.sleep(interval)
                for session_key in self.active_intents.session_keys():
                    active_tasks = self._active_tasks.get(session_key, [])
                    active_count = sum(1 for task in active_tasks if not task.done())
                    running_subagents = self.subagents.get_running_count_by_session(session_key)
                    await self.active_intents.process_session(
                        session_key,
                        active_task_count=active_count,
                        running_subagents=running_subagents,
                    )
        except asyncio.CancelledError:
            raise

    def stop(self) -> None:
        """Stop the agent loop."""
        self._running = False
        logger.info("Agent loop stopping")

    async def _process_system_message(
        self,
        msg: InboundMessage,
        session_key: str | None = None,
        on_progress: Callable[..., Awaitable[None]] | None = None,
        on_stream: Callable[[str], Awaitable[None]] | None = None,
        on_stream_end: Callable[..., Awaitable[None]] | None = None,
        pending_queue: asyncio.Queue | None = None,
        capability_snapshot: CapabilitySnapshot | None = None,
    ) -> OutboundMessage | None:
        """Process a system inbound message (e.g. subagent announce)."""
        channel, chat_id = (
            msg.chat_id.split(":", 1) if ":" in msg.chat_id else ("cli", msg.chat_id)
        )
        logger.info("Processing system message from {}", msg.sender_id)
        key = msg.session_key_override or f"{channel}:{chat_id}"
        session = self.sessions.get_or_create(key)
        if self._restore_runtime_checkpoint(session):
            self.sessions.save(session)
        if self._restore_pending_user_turn(session):
            self.sessions.save(session)

        session, pending = self.auto_compact.prepare_session(session, key)
        if pending:
            logger.info("Memory compact triggered for session {}", key)

        await self.consolidator.maybe_consolidate_by_tokens(
            session,
            replay_max_messages=self._max_messages,
        )
        event_kind = str(msg.metadata.get("injected_event") or "").strip()
        is_subagent = msg.sender_id == "subagent" or event_kind == "subagent_result"
        is_active_intent = event_kind == "active_intent"
        persisted_subagent = False
        if is_subagent and self._persist_subagent_followup(session, msg):
            persisted_subagent = True
            logger.debug("Subagent result persisted for session {}", key)
            self.sessions.save(session)
        runtime_context = self._resolve_runtime_context(
            msg,
            channel=channel,
            chat_id=chat_id,
            session_key=key,
        )
        snapshot = capability_snapshot or self._snapshot_for_trigger(runtime_context.trigger)
        self._set_tool_context(
            channel, chat_id, msg.metadata.get("message_id"),
            msg.metadata,
            session_key=key,
            capability_snapshot=snapshot,
            runtime_context=runtime_context,
        )
        _hist_kwargs: dict[str, Any] = {
            "max_messages": self._max_messages,
            "max_tokens": self._replay_token_budget(),
            "include_timestamps": True,
        }
        history = session.get_history(**_hist_kwargs)
        history_for_model = list(history)
        if is_subagent and persisted_subagent:
            for index in range(len(history_for_model) - 1, -1, -1):
                candidate = history_for_model[index]
                if (
                    candidate.get("role") == "assistant"
                    and candidate.get("content") == msg.content
                ):
                    history_for_model.pop(index)
                    break

        messages = self.context.build_messages(
            history=history_for_model,
            current_message=None if (is_subagent or is_active_intent) else msg.content,
            channel=channel,
            chat_id=chat_id,
            current_role="user",
            sender_id=msg.sender_id,
            session_summary=pending,
            session_metadata=session.metadata,
            internal_event=(
                ("subagent_result", msg.content)
                if is_subagent
                else ("active_intent", msg.content) if is_active_intent else None
            ),
        )
        final_content, _, all_msgs, stop_reason, _ = await self._run_agent_loop(
            messages, session=session, channel=channel, chat_id=chat_id,
            message_id=msg.metadata.get("message_id"),
            metadata=msg.metadata,
            session_key=key,
            pending_queue=pending_queue,
            actor_id=runtime_context.actor_id,
            trigger=runtime_context.trigger,
            capability_snapshot=snapshot,
        )
        save_skip = 1 + len(history_for_model) + (1 if (is_subagent or is_active_intent) else 0)
        self._save_turn(session, all_msgs, save_skip)
        session.enforce_file_cap(on_archive=self._archive_session_file_cap)
        self._clear_runtime_checkpoint(session)
        self.sessions.save(session)
        self._schedule_background(
            self.consolidator.maybe_consolidate_by_tokens(
                session,
                replay_max_messages=self._max_messages,
            )
        )
        options = ask_user_options_from_messages(all_msgs) if stop_reason == "ask_user" else []
        content, buttons = ask_user_outbound(
            final_content or "Background task completed.",
            options,
            channel,
        )
        outbound_metadata: dict[str, Any] = {}
        if channel == "slack" and key.startswith("slack:") and key.count(":") >= 2:
            outbound_metadata["slack"] = {"thread_ts": key.split(":", 2)[2]}
        if origin_message_id := msg.metadata.get("origin_message_id"):
            outbound_metadata["origin_message_id"] = origin_message_id
        if channel == "websocket":
            outbound_metadata["goal_state"] = goal_state_ws_blob(session.metadata)
        return OutboundMessage(
            channel=channel,
            chat_id=chat_id,
            content=content,
            buttons=buttons,
            metadata=outbound_metadata,
        )

    async def _process_message(
        self,
        msg: InboundMessage,
        session_key: str | None = None,
        on_progress: Callable[..., Awaitable[None]] | None = None,
        on_stream: Callable[[str], Awaitable[None]] | None = None,
        on_stream_end: Callable[..., Awaitable[None]] | None = None,
        pending_queue: asyncio.Queue | None = None,
        capability_snapshot: CapabilitySnapshot | None = None,
    ) -> OutboundMessage | None:
        """Process a single inbound message and return the response."""
        self._refresh_provider_snapshot()

        if msg.channel == "system":
            return await self._process_system_message(
                msg,
                session_key=session_key,
                on_progress=on_progress,
                on_stream=on_stream,
                on_stream_end=on_stream_end,
                pending_queue=pending_queue,
                capability_snapshot=capability_snapshot,
            )

        key = session_key or msg.session_key
        ctx = TurnContext(
            msg=msg,
            session=None,
            session_key=key,
            state=TurnState.RESTORE,
            turn_id=f"{key}:{time.time_ns()}",
            on_progress=on_progress,
            on_stream=on_stream,
            on_stream_end=on_stream_end,
            pending_queue=pending_queue,
            capability_snapshot=capability_snapshot,
        )

        while ctx.state is not TurnState.DONE:
            handler_name = f"_state_{ctx.state.name.lower()}"
            handler = getattr(self, handler_name, None)
            if handler is None:
                raise RuntimeError(f"Missing state handler for {ctx.state}")

            t0 = time.perf_counter()
            try:
                event = await handler(ctx)
            except Exception:
                duration = (time.perf_counter() - t0) * 1000
                ctx.trace.append(
                    StateTraceEntry(
                        state=ctx.state,
                        started_at=t0,
                        duration_ms=duration,
                        event="",
                        error="exception",
                    )
                )
                raise

            duration = (time.perf_counter() - t0) * 1000
            ctx.trace.append(
                StateTraceEntry(
                    state=ctx.state,
                    started_at=t0,
                    duration_ms=duration,
                    event=event,
                )
            )
            logger.debug(
                "[turn {}] State {} took {:.1f}ms -> event {}",
                ctx.turn_id,
                ctx.state.name,
                duration,
                event,
            )

            next_state = self._TRANSITIONS.get((ctx.state, event))
            if next_state is None:
                raise RuntimeError(
                    f"[turn {ctx.turn_id}] No transition from {ctx.state} "
                    f"on event {event!r}"
                )
            ctx.state = next_state

        logger.debug(
            "[turn {}] Turn completed after {} states",
            ctx.turn_id,
            len(ctx.trace),
        )
        return ctx.outbound

    def _assemble_outbound(
        self,
        msg: InboundMessage,
        final_content: str,
        all_msgs: list[dict[str, Any]],
        stop_reason: str,
        had_injections: bool,
        generated_media: list[str],
        on_stream: Callable[[str], Awaitable[None]] | None,
    ) -> OutboundMessage | None:
        """Assemble the final outbound message from turn results."""
        # MessageTool suppression
        if (mt := self.tools.get("message")) and isinstance(mt, MessageTool) and mt._sent_in_turn:
            if not had_injections or stop_reason == "empty_final_response":
                return None

        preview = final_content[:120] + "..." if len(final_content) > 120 else final_content
        logger.info("Response to {}:{}: {}", msg.channel, msg.sender_id, preview)

        meta = dict(msg.metadata or {})
        content, buttons = ask_user_outbound(
            final_content,
            ask_user_options_from_messages(all_msgs) if stop_reason == "ask_user" else [],
            msg.channel,
        )
        if on_stream is not None and stop_reason not in {"ask_user", "error", "tool_error"}:
            meta["_streamed"] = True
        if msg.channel == "websocket":
            meta["goal_state"] = goal_state_ws_blob(
                self.sessions.get_or_create(self._effective_session_key(msg)).metadata
            )

        return OutboundMessage(
            channel=msg.channel,
            chat_id=msg.chat_id,
            content=content,
            media=generated_media,
            metadata=meta,
            buttons=buttons,
        )

    async def _state_restore(self, ctx: TurnContext) -> TurnState:
        """Restore checkpoint / pending user turn; extract documents."""
        msg = ctx.msg

        if msg.media:
            new_content, image_only = extract_documents(msg.content, msg.media)
            ctx.msg = dataclasses.replace(msg, content=new_content, media=image_only)
            msg = ctx.msg

        preview = msg.content[:80] + "..." if len(msg.content) > 80 else msg.content
        logger.info("Processing message from {}:{}: {}", msg.channel, msg.sender_id, preview)

        # Session is already fetched by the caller (_process_message) but
        # ensure it exists in case this handler is invoked independently.
        if ctx.session is None:
            ctx.session = self.sessions.get_or_create(ctx.session_key)
        mark_webui_session(ctx.session, msg.metadata)

        if self._restore_runtime_checkpoint(ctx.session):
            self.sessions.save(ctx.session)
        if self._restore_pending_user_turn(ctx.session):
            self.sessions.save(ctx.session)

        return "ok"

    async def _state_compact(self, ctx: TurnContext) -> str:
        ctx.session, pending = self.auto_compact.prepare_session(ctx.session, ctx.session_key)
        ctx.pending_summary = pending
        return "ok"

    async def _state_command(self, ctx: TurnContext) -> str:
        raw = ctx.msg.content.strip()
        cmd_ctx = CommandContext(
            msg=ctx.msg, session=ctx.session, key=ctx.session_key, raw=raw, loop=self
        )
        result = await self.commands.dispatch(cmd_ctx)
        if result is not None:
            ctx.outbound = result
            # Shortcut commands skip BUILD/RUN/SAVE, so persist both sides of
            # the turn here.  Otherwise the live WebUI may briefly receive the
            # command response, then lose it when history hydration follows
            # turn_end/session updates.  Keep these rows out of future LLM
            # context with the _command marker.
            self._persist_shortcut_command_turn(ctx.msg, ctx.session_key, result)
            if self._is_webui_message(ctx.msg):
                result.metadata["_webui_transcript_recorded"] = True
            return "shortcut"
        return "dispatch"

    async def _state_build(self, ctx: TurnContext) -> str:
        await self.consolidator.maybe_consolidate_by_tokens(
            ctx.session,
            replay_max_messages=self._max_messages,
        )
        runtime_context = self._resolve_runtime_context(ctx.msg, session_key=ctx.session_key)
        ctx.runtime_context = runtime_context
        snapshot = ctx.capability_snapshot or self._snapshot_for_trigger(runtime_context.trigger)
        self._set_tool_context(
            ctx.msg.channel,
            ctx.msg.chat_id,
            ctx.msg.metadata.get("message_id"),
            ctx.msg.metadata,
            session_key=ctx.session_key,
            capability_snapshot=snapshot,
            runtime_context=runtime_context,
        )
        if message_tool := self.tools.get("message"):
            if isinstance(message_tool, MessageTool):
                message_tool.start_turn()

        _hist_kwargs: dict[str, Any] = {
            "max_messages": self._max_messages,
            "max_tokens": self._replay_token_budget(),
            "include_timestamps": True,
        }
        ctx.history = ctx.session.get_history(**_hist_kwargs)

        pending_ask_id = pending_ask_user_id(ctx.history)
        ctx.initial_messages = self._build_initial_messages(
            ctx.msg, ctx.session, ctx.history, pending_ask_id, ctx.pending_summary
        )
        ctx.user_persisted_early = self._persist_user_message_early(
            ctx.msg, ctx.session, pending_ask_id
        )
        if ctx.user_persisted_early:
            self._schedule_session_search_refresh(sources=["sessions"])

        if ctx.on_progress is None:
            ctx.on_progress = await self._build_bus_progress_callback(ctx.msg)
        if ctx.on_retry_wait is None:
            ctx.on_retry_wait = await self._build_retry_wait_callback(ctx.msg)

        return "ok"

    async def _state_run(self, ctx: TurnContext) -> str:
        runtime_context = ctx.runtime_context or self._resolve_runtime_context(
            ctx.msg,
            session_key=ctx.session_key,
        )
        ctx.runtime_context = runtime_context
        snapshot = ctx.capability_snapshot or self._snapshot_for_trigger(runtime_context.trigger)
        await publish_turn_run_status(self.bus, ctx.msg, "running")
        try:
            result = await self._run_agent_loop(
                ctx.initial_messages,
                on_progress=ctx.on_progress,
                on_stream=ctx.on_stream,
                on_stream_end=ctx.on_stream_end,
                on_retry_wait=ctx.on_retry_wait,
                session=ctx.session,
                channel=ctx.msg.channel,
                chat_id=ctx.msg.chat_id,
                message_id=ctx.msg.metadata.get("message_id"),
                metadata=ctx.msg.metadata,
                session_key=ctx.session_key,
                pending_queue=ctx.pending_queue,
                actor_id=runtime_context.actor_id,
                trigger=runtime_context.trigger,
                capability_snapshot=snapshot,
            )
        finally:
            if ctx.msg.channel == "websocket":
                latency = websocket_turn_latency_ms(str(ctx.msg.chat_id))
                if latency is not None:
                    ctx.msg.metadata["webui_turn_latency_ms"] = latency
            await publish_turn_run_status(self.bus, ctx.msg, "idle")
        final_content, tools_used, all_msgs, stop_reason, had_injections = result
        ctx.final_content = final_content
        ctx.tools_used = tools_used
        ctx.all_messages = all_msgs
        ctx.stop_reason = stop_reason
        ctx.had_injections = had_injections
        return "ok"

    async def _state_save(self, ctx: TurnContext) -> str:
        if ctx.final_content is None or not ctx.final_content.strip():
            ctx.final_content = EMPTY_FINAL_RESPONSE_MESSAGE

        ctx.save_skip = 1 + len(ctx.history) + (1 if ctx.user_persisted_early else 0)
        skip_msgs = ctx.all_messages[ctx.save_skip:]
        ctx.generated_media = generated_image_paths_from_messages(skip_msgs)
        message_tool = self.tools.get("message")
        extra_media = (
            message_tool.turn_delivered_media_paths()
            if hasattr(message_tool, "turn_delivered_media_paths")
            else []
        )
        merge_turn_media_into_last_assistant(ctx.all_messages, ctx.generated_media, extra_media)

        self._save_turn(ctx.session, ctx.all_messages, ctx.save_skip)
        ctx.session.enforce_file_cap(on_archive=self._archive_session_file_cap)
        self._clear_pending_user_turn(ctx.session)
        self._clear_runtime_checkpoint(ctx.session)
        self.sessions.save(ctx.session)
        self._schedule_session_search_refresh(sources=["sessions", "history"])
        self._schedule_background(
            self.consolidator.maybe_consolidate_by_tokens(
                ctx.session,
                replay_max_messages=self._max_messages,
            )
        )
        self._schedule_background_review(ctx)
        self._schedule_curator_review(ctx)
        return "ok"

    def _schedule_session_search_refresh(
        self,
        *,
        sources: list[str] | None = None,
        force: bool = False,
    ) -> None:
        if not getattr(self.session_search_index, "enabled", False):
            return
        self._schedule_background(
            self._refresh_session_search_index(sources=sources, force=force)
        )

    async def _refresh_session_search_index(
        self,
        *,
        sources: list[str] | None = None,
        force: bool = False,
    ) -> None:
        await asyncio.to_thread(
            self.session_search_index.refresh_incremental,
            sources=sources,
            force=force,
        )

    def _schedule_background_review(self, ctx: TurnContext) -> None:
        """Schedule a controlled learning review for successful foreground turns."""
        if ctx.session is None:
            return
        self.background_review.refresh_config()
        if not self.background_review.enabled:
            return
        if ctx.stop_reason in {"ask_user", "error", "tool_error"}:
            return
        if ctx.msg.channel == "system" or ctx.msg.sender_id == "subagent":
            return
        if not (ctx.final_content or "").strip():
            return

        max_recent = int(
            getattr(self.background_review.config, "max_recent_messages", 12) or 12
        )
        messages = [
            dict(message)
            for message in ctx.session.messages
            if not message.get("_command")
        ][-max_recent:]
        self._schedule_background(
            self.background_review.review_turn(
                session_key=ctx.session_key,
                turn_id=ctx.turn_id,
                channel=ctx.msg.channel,
                chat_id=ctx.msg.chat_id,
                message_id=ctx.msg.metadata.get("message_id"),
                messages=messages,
            )
        )

    def _schedule_curator_review(self, ctx: TurnContext) -> None:
        """Schedule deterministic curator review after successful foreground turns."""
        if ctx.session is None:
            return
        self.curator.refresh_config()
        if not self.curator.enabled:
            return
        if ctx.stop_reason in {"ask_user", "error", "tool_error"}:
            return
        if ctx.msg.channel == "system" or ctx.msg.sender_id == "subagent":
            return
        if not (ctx.final_content or "").strip():
            return
        self._schedule_background(
            self.curator.review_workspace(
                session_key=ctx.session_key,
                turn_id=ctx.turn_id,
            )
        )

    async def _state_respond(self, ctx: TurnContext) -> str:
        ctx.outbound = self._assemble_outbound(
            ctx.msg,
            ctx.final_content,
            ctx.all_messages,
            ctx.stop_reason,
            ctx.had_injections,
            ctx.generated_media,
            ctx.on_stream,
        )
        return "ok"

    def _sanitize_persisted_blocks(
        self,
        content: list[dict[str, Any]],
        *,
        should_truncate_text: bool = False,
        drop_runtime: bool = False,
    ) -> list[dict[str, Any]]:
        """Strip volatile multimodal payloads before writing session history."""
        return AgentLoop._turn_persist_manager(self).sanitize_persisted_blocks(
            content,
            should_truncate_text=should_truncate_text,
            drop_runtime=drop_runtime,
        )

    def _save_turn(self, session: Session, messages: list[dict], skip: int) -> None:
        """Save new-turn messages into session, truncating large tool results."""
        AgentLoop._turn_persist_manager(self).save_turn(session, messages, skip)

    def _persist_subagent_followup(self, session: Session, msg: InboundMessage) -> bool:
        """Persist subagent follow-ups before prompt assembly so history stays durable.

        Returns True if a new entry was appended; False if the follow-up was
        deduped (same ``subagent_task_id`` already in session) or carries no
        content worth persisting.
        """
        return AgentLoop._turn_persist_manager(self).persist_subagent_followup(session, msg)

    def _set_runtime_checkpoint(self, session: Session, payload: dict[str, Any]) -> None:
        """Persist the latest in-flight turn state into session metadata."""
        AgentLoop._turn_persist_manager(self).set_checkpoint(session, payload)

    def _mark_pending_user_turn(self, session: Session) -> None:
        AgentLoop._turn_persist_manager(self).mark_pending_user_turn(session)

    def _clear_pending_user_turn(self, session: Session) -> None:
        AgentLoop._turn_persist_manager(self).clear_pending_user_turn(session)

    def _clear_runtime_checkpoint(self, session: Session) -> None:
        AgentLoop._turn_persist_manager(self).clear_checkpoint(session)

    @staticmethod
    def _checkpoint_message_key(message: dict[str, Any]) -> tuple[Any, ...]:
        return TurnPersistManager.checkpoint_message_key(message)

    def _restore_runtime_checkpoint(self, session: Session) -> bool:
        """Materialize an unfinished turn into session history before a new request."""
        return AgentLoop._turn_persist_manager(self).restore_checkpoint(session)

    def _restore_pending_user_turn(self, session: Session) -> bool:
        """Close a turn that only persisted the user message before crashing."""
        return AgentLoop._turn_persist_manager(self).restore_pending_user_turn(session)

    async def process_direct(
        self,
        content: str,
        session_key: str = "cli:direct",
        channel: str = "cli",
        chat_id: str = "direct",
        media: list[str] | None = None,
        on_progress: Callable[..., Awaitable[None]] | None = None,
        on_stream: Callable[[str], Awaitable[None]] | None = None,
        on_stream_end: Callable[..., Awaitable[None]] | None = None,
    ) -> OutboundMessage | None:
        """Process a message directly and return the outbound payload."""
        await self._connect_mcp()
        msg = InboundMessage(
            channel=channel, sender_id="user", chat_id=chat_id,
            content=content, media=media or [],
        )
        return await self._process_message(
            msg,
            session_key=session_key,
            on_progress=on_progress,
            on_stream=on_stream,
            on_stream_end=on_stream_end,
        )
