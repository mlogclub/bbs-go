"""Configuration schema using Pydantic."""

from pathlib import Path
from typing import Any, Literal

from pydantic import AliasChoices, BaseModel, ConfigDict, Field, model_validator
from pydantic.alias_generators import to_camel
from pydantic_settings import BaseSettings

from OriginAgent.cron.types import CronSchedule


class Base(BaseModel):
    """Base model that accepts both camelCase and snake_case keys."""

    model_config = ConfigDict(alias_generator=to_camel, populate_by_name=True)

class ChannelsConfig(Base):
    """Configuration for chat channels.

    Built-in and plugin channel configs are stored as extra fields (dicts).
    Each channel parses its own config in __init__.
    Per-channel "streaming": true enables streaming output (requires send_delta impl).
    """

    model_config = ConfigDict(extra="allow")

    send_progress: bool = True  # stream agent's text progress to the channel
    send_tool_hints: bool = False  # stream tool-call hints (e.g. read_file("…"))
    show_reasoning: bool = True  # surface model reasoning when channel implements it
    send_max_retries: int = Field(default=3, ge=0, le=10)  # Max delivery attempts (initial send included)
    transcription_provider: str = "groq"  # Voice transcription backend: "groq" or "openai"
    transcription_language: str | None = Field(default=None, pattern=r"^[a-z]{2,3}$")  # Optional ISO-639-1 hint for audio transcription


class DreamConfig(Base):
    """Dream memory consolidation configuration."""

    _HOUR_MS = 3_600_000

    interval_h: int = Field(default=2, ge=1)  # Every 2 hours by default
    cron: str | None = Field(default=None, exclude=True)  # Legacy compatibility override
    model_override: str | None = Field(
        default=None,
        validation_alias=AliasChoices("modelOverride", "model", "model_override"),
    )  # Optional Dream-specific model override
    max_batch_size: int = Field(default=20, ge=1)  # Max history entries per run
    # Bumped from 10 to 15 in #3212 (exp002: +30% dedup, no accuracy loss; >15 plateaus).
    max_iterations: int = Field(default=15, ge=1)  # Max tool calls per Phase 2
    # Per-line git-blame age annotation in Phase 1 prompt (see #3212). Default
    # on — set to False to feed MEMORY.md raw if a specific LLM reacts poorly
    # to the `← Nd` suffix or you want deterministic, git-independent prompts.
    annotate_line_ages: bool = True

    def build_schedule(self, timezone: str) -> CronSchedule:
        """Build the runtime schedule, preferring the legacy cron override if present."""
        if self.cron:
            return CronSchedule(kind="cron", expr=self.cron, tz=timezone)
        return CronSchedule(kind="every", every_ms=self.interval_h * self._HOUR_MS)

    def describe_schedule(self) -> str:
        """Return a human-readable summary for logs and startup output."""
        if self.cron:
            return f"cron {self.cron} (legacy)"
        hours = self.interval_h
        return f"every {hours}h"


class InlineFallbackConfig(Base):
    """Inline fallback model candidate."""

    model: str
    provider: str = "auto"
    max_tokens: int | None = None
    context_window_tokens: int | None = None
    temperature: float | None = None
    reasoning_effort: str | None = None


class ModelPresetConfig(Base):
    """Named runtime model + provider configuration."""

    model: str
    provider: str = "auto"
    max_tokens: int | None = None
    context_window_tokens: int | None = None
    temperature: float | None = None
    reasoning_effort: str | None = None
    fallback_models: list[str | InlineFallbackConfig] = Field(default_factory=list)

    def to_generation_settings(self):
        from OriginAgent.providers.base import GenerationSettings

        return GenerationSettings(
            temperature=self.temperature,
            max_tokens=self.max_tokens,
            reasoning_effort=self.reasoning_effort,
        )


FallbackCandidate = str | InlineFallbackConfig


class AuxiliaryTaskConfig(Base):
    """Per-background-task auxiliary LLM routing settings."""

    model_override: str | None = Field(
        default=None,
        validation_alias=AliasChoices("modelOverride", "model", "model_override"),
    )
    timeout_s: float | None = Field(
        default=None,
        ge=0,
        validation_alias=AliasChoices("timeoutS", "timeout", "timeout_s"),
    )
    fallback_models: list[FallbackCandidate] = Field(default_factory=list)


class AuxiliaryConfig(Base):
    """Background LLM routing and fallback configuration."""

    enabled: bool = True
    payment_cooldown_s: int = Field(
        default=1800,
        ge=0,
        validation_alias=AliasChoices("paymentCooldownS", "payment_cooldown_s"),
    )
    transient_cooldown_s: int = Field(
        default=60,
        ge=0,
        validation_alias=AliasChoices("transientCooldownS", "transient_cooldown_s"),
    )
    tasks: dict[str, AuxiliaryTaskConfig] = Field(default_factory=dict)


class DomainPacksConfig(Base):
    """Domain pack discovery and prompt injection configuration."""

    enabled: bool = True
    disabled: list[str] = Field(default_factory=list)
    active: list[str] = Field(default_factory=list)
    max_capability_chars: int = Field(
        default=4000,
        ge=0,
        validation_alias=AliasChoices("maxCapabilityChars", "max_capability_chars"),
        serialization_alias="maxCapabilityChars",
    )


class BackgroundReviewConfig(Base):
    """Controlled background learning proposal generation."""

    enabled: bool = False
    max_recent_messages: int = Field(
        default=12,
        ge=1,
        le=100,
        validation_alias=AliasChoices("maxRecentMessages", "max_recent_messages"),
        serialization_alias="maxRecentMessages",
    )
    max_prompt_chars: int = Field(
        default=16000,
        ge=1000,
        validation_alias=AliasChoices("maxPromptChars", "max_prompt_chars"),
        serialization_alias="maxPromptChars",
    )
    max_proposals_per_turn: int = Field(
        default=8,
        ge=1,
        le=50,
        validation_alias=AliasChoices("maxProposalsPerTurn", "max_proposals_per_turn"),
        serialization_alias="maxProposalsPerTurn",
    )
    max_concurrent_reviews: int = Field(
        default=1,
        ge=1,
        le=8,
        validation_alias=AliasChoices("maxConcurrentReviews", "max_concurrent_reviews"),
        serialization_alias="maxConcurrentReviews",
    )
    allowed_proposal_types: list[str] = Field(
        default_factory=lambda: ["memory", "fact", "skill", "workflow"],
        validation_alias=AliasChoices("allowedProposalTypes", "allowed_proposal_types"),
        serialization_alias="allowedProposalTypes",
    )


class CuratorConfig(Base):
    """Deterministic curator proposal generation."""

    enabled: bool = False
    max_proposals_per_run: int = Field(
        default=12,
        ge=1,
        le=50,
        validation_alias=AliasChoices("maxProposalsPerRun", "max_proposals_per_run"),
        serialization_alias="maxProposalsPerRun",
    )


class EvolutionSandboxConfig(Base):
    """Read-only sandbox evaluation settings for governed evolution."""

    enabled: bool = True
    read_only_tools: list[str] = Field(
        default_factory=lambda: ["read_file", "glob", "grep"],
        validation_alias=AliasChoices("readOnlyTools", "read_only_tools"),
        serialization_alias="readOnlyTools",
    )
    timeout_seconds: int = Field(
        default=10,
        ge=1,
        le=120,
        validation_alias=AliasChoices("timeoutSeconds", "timeout_seconds"),
        serialization_alias="timeoutSeconds",
    )
    max_output_chars: int = Field(
        default=8000,
        ge=100,
        le=100_000,
        validation_alias=AliasChoices("maxOutputChars", "max_output_chars"),
        serialization_alias="maxOutputChars",
    )
    max_replay_samples: int = Field(
        default=3,
        ge=0,
        le=20,
        validation_alias=AliasChoices("maxReplaySamples", "max_replay_samples"),
        serialization_alias="maxReplaySamples",
    )
    cache_ttl_hours: int = Field(
        default=24,
        ge=0,
        le=168,
        validation_alias=AliasChoices("cacheTtlHours", "cache_ttl_hours"),
        serialization_alias="cacheTtlHours",
    )


class EvolutionTrialConfig(Base):
    """Isolated trial-mode settings for verified evolution artifacts."""

    enabled: bool = True
    isolated_workspace: bool = Field(
        default=True,
        validation_alias=AliasChoices("isolatedWorkspace", "isolated_workspace"),
        serialization_alias="isolatedWorkspace",
    )
    read_only_tools_only: bool = Field(
        default=True,
        validation_alias=AliasChoices("readOnlyToolsOnly", "read_only_tools_only"),
        serialization_alias="readOnlyToolsOnly",
    )
    blocked_tools: list[str] = Field(
        default_factory=lambda: ["write_file", "edit_file", "exec", "message", "cron", "spawn"],
        validation_alias=AliasChoices("blockedTools", "blocked_tools"),
        serialization_alias="blockedTools",
    )
    temp_dir: str = Field(
        default="",
        validation_alias=AliasChoices("tempDir", "temp_dir"),
        serialization_alias="tempDir",
    )
    max_step_output_chars: int = Field(
        default=2000,
        ge=100,
        le=100_000,
        validation_alias=AliasChoices("maxStepOutputChars", "max_step_output_chars"),
        serialization_alias="maxStepOutputChars",
    )
    max_retained_trial_logs: int = Field(
        default=10,
        ge=0,
        le=1000,
        validation_alias=AliasChoices("maxRetainedTrialLogs", "max_retained_trial_logs"),
        serialization_alias="maxRetainedTrialLogs",
    )
    trial_log_retention_days: int = Field(
        default=30,
        ge=1,
        le=3650,
        validation_alias=AliasChoices("trialLogRetentionDays", "trial_log_retention_days"),
        serialization_alias="trialLogRetentionDays",
    )


class EvolutionConfig(Base):
    """Governed self-evolution observability settings."""

    mode: Literal["conservative", "curated", "exploratory", "aggressive"] = "conservative"
    allow_manual_override: bool = Field(
        default=False,
        validation_alias=AliasChoices("allowManualOverride", "allow_manual_override"),
        serialization_alias="allowManualOverride",
    )
    dry_run: bool = Field(
        default=True,
        validation_alias=AliasChoices("dryRun", "dry_run"),
        serialization_alias="dryRun",
    )
    signal_retention_days: int = Field(
        default=30,
        ge=1,
        le=365,
        validation_alias=AliasChoices("signalRetentionDays", "signal_retention_days"),
        serialization_alias="signalRetentionDays",
    )
    outcome_retention_days: int = Field(
        default=90,
        ge=1,
        le=3650,
        validation_alias=AliasChoices("outcomeRetentionDays", "outcome_retention_days"),
        serialization_alias="outcomeRetentionDays",
    )
    outcome_archive_enabled: bool = Field(
        default=True,
        validation_alias=AliasChoices("outcomeArchiveEnabled", "outcome_archive_enabled"),
        serialization_alias="outcomeArchiveEnabled",
    )
    dependency_stale_cleanup_enabled: bool = Field(
        default=True,
        validation_alias=AliasChoices("dependencyStaleCleanupEnabled", "dependency_stale_cleanup_enabled"),
        serialization_alias="dependencyStaleCleanupEnabled",
    )
    health_history_retention_days: int = Field(
        default=90,
        ge=1,
        le=3650,
        validation_alias=AliasChoices("healthHistoryRetentionDays", "health_history_retention_days"),
        serialization_alias="healthHistoryRetentionDays",
    )
    max_health_history_snapshots: int = Field(
        default=100,
        ge=0,
        le=10_000,
        validation_alias=AliasChoices("maxHealthHistorySnapshots", "max_health_history_snapshots"),
        serialization_alias="maxHealthHistorySnapshots",
    )
    workflow_min_seen_count: int = Field(
        default=3,
        ge=1,
        le=100,
        validation_alias=AliasChoices("workflowMinSeenCount", "workflow_min_seen_count"),
        serialization_alias="workflowMinSeenCount",
    )
    workflow_min_evidence_sources: int = Field(
        default=2,
        ge=1,
        le=50,
        validation_alias=AliasChoices("workflowMinEvidenceSources", "workflow_min_evidence_sources"),
        serialization_alias="workflowMinEvidenceSources",
    )
    workflow_priority_threshold: float = Field(
        default=0.7,
        ge=0.0,
        le=1.0,
        validation_alias=AliasChoices("workflowPriorityThreshold", "workflow_priority_threshold"),
        serialization_alias="workflowPriorityThreshold",
    )
    max_high_score_signals: int = Field(
        default=5,
        ge=0,
        le=20,
        validation_alias=AliasChoices("maxHighScoreSignals", "max_high_score_signals"),
        serialization_alias="maxHighScoreSignals",
    )
    max_proposals_per_cycle: int = Field(
        default=3,
        ge=0,
        le=20,
        validation_alias=AliasChoices("maxProposalsPerCycle", "max_proposals_per_cycle"),
        serialization_alias="maxProposalsPerCycle",
    )
    static_gate_allowed_workflow_tools: list[str] = Field(
        default_factory=lambda: ["read_file", "glob", "grep", "web_fetch"],
        validation_alias=AliasChoices(
            "staticGateAllowedWorkflowTools",
            "static_gate_allowed_workflow_tools",
        ),
        serialization_alias="staticGateAllowedWorkflowTools",
    )
    static_gate_max_workflow_steps: int = Field(
        default=10,
        ge=1,
        le=25,
        validation_alias=AliasChoices("staticGateMaxWorkflowSteps", "static_gate_max_workflow_steps"),
        serialization_alias="staticGateMaxWorkflowSteps",
    )
    auto_verify_workflows: bool = Field(
        default=False,
        validation_alias=AliasChoices("autoVerifyWorkflows", "auto_verify_workflows"),
        serialization_alias="autoVerifyWorkflows",
    )
    workflow_auto_verify_threshold: float = Field(
        default=0.9,
        ge=0.0,
        le=1.0,
        validation_alias=AliasChoices("workflowAutoVerifyThreshold", "workflow_auto_verify_threshold"),
        serialization_alias="workflowAutoVerifyThreshold",
    )
    workflow_auto_verify_min_seen_count: int = Field(
        default=5,
        ge=1,
        le=100,
        validation_alias=AliasChoices("workflowAutoVerifyMinSeenCount", "workflow_auto_verify_min_seen_count"),
        serialization_alias="workflowAutoVerifyMinSeenCount",
    )
    skill_candidates_enabled: bool = Field(
        default=False,
        validation_alias=AliasChoices("skillCandidatesEnabled", "skill_candidates_enabled"),
        serialization_alias="skillCandidatesEnabled",
    )
    skill_min_seen_count: int = Field(
        default=5,
        ge=1,
        le=100,
        validation_alias=AliasChoices("skillMinSeenCount", "skill_min_seen_count"),
        serialization_alias="skillMinSeenCount",
    )
    skill_priority_threshold: float = Field(
        default=0.85,
        ge=0.0,
        le=1.0,
        validation_alias=AliasChoices("skillPriorityThreshold", "skill_priority_threshold"),
        serialization_alias="skillPriorityThreshold",
    )
    static_gate_allowed_skill_tools: list[str] = Field(
        default_factory=lambda: ["read_file", "glob", "grep"],
        validation_alias=AliasChoices(
            "staticGateAllowedSkillTools",
            "static_gate_allowed_skill_tools",
        ),
        serialization_alias="staticGateAllowedSkillTools",
    )
    max_skill_proposals_per_cycle: int = Field(
        default=1,
        ge=0,
        le=10,
        validation_alias=AliasChoices("maxSkillProposalsPerCycle", "max_skill_proposals_per_cycle"),
        serialization_alias="maxSkillProposalsPerCycle",
    )
    sandbox: EvolutionSandboxConfig = Field(
        default_factory=EvolutionSandboxConfig,
        validation_alias=AliasChoices("sandbox"),
        serialization_alias="sandbox",
    )
    trial: EvolutionTrialConfig = Field(
        default_factory=EvolutionTrialConfig,
        validation_alias=AliasChoices("trial"),
        serialization_alias="trial",
    )
    feedback_calibration_enabled: bool = Field(
        default=True,
        validation_alias=AliasChoices("feedbackCalibrationEnabled", "feedback_calibration_enabled"),
        serialization_alias="feedbackCalibrationEnabled",
    )
    feedback_reject_multiplier: float = Field(
        default=0.8,
        ge=0.05,
        le=1.0,
        validation_alias=AliasChoices("feedbackRejectMultiplier", "feedback_reject_multiplier"),
        serialization_alias="feedbackRejectMultiplier",
    )
    feedback_rollback_multiplier: float = Field(
        default=0.6,
        ge=0.05,
        le=1.0,
        validation_alias=AliasChoices("feedbackRollbackMultiplier", "feedback_rollback_multiplier"),
        serialization_alias="feedbackRollbackMultiplier",
    )
    feedback_rollback_delta: float = Field(
        default=-0.1,
        ge=-1.0,
        le=0.0,
        validation_alias=AliasChoices("feedbackRollbackDelta", "feedback_rollback_delta"),
        serialization_alias="feedbackRollbackDelta",
    )
    feedback_positive_multiplier: float = Field(
        default=1.02,
        ge=1.0,
        le=2.0,
        validation_alias=AliasChoices("feedbackPositiveMultiplier", "feedback_positive_multiplier"),
        serialization_alias="feedbackPositiveMultiplier",
    )
    feedback_suppress_after_negative_count: int = Field(
        default=3,
        ge=1,
        le=20,
        validation_alias=AliasChoices(
            "feedbackSuppressAfterNegativeCount",
            "feedback_suppress_after_negative_count",
        ),
        serialization_alias="feedbackSuppressAfterNegativeCount",
    )
    feedback_cooldown_days: int = Field(
        default=14,
        ge=1,
        le=365,
        validation_alias=AliasChoices("feedbackCooldownDays", "feedback_cooldown_days"),
        serialization_alias="feedbackCooldownDays",
    )
    feedback_trend_window_days: int = Field(
        default=14,
        ge=1,
        le=365,
        validation_alias=AliasChoices("feedbackTrendWindowDays", "feedback_trend_window_days"),
        serialization_alias="feedbackTrendWindowDays",
    )


class LearningConfig(Base):
    """Agent self-improvement and review configuration."""

    background_review: BackgroundReviewConfig = Field(
        default_factory=BackgroundReviewConfig,
        validation_alias=AliasChoices("backgroundReview", "background_review"),
        serialization_alias="backgroundReview",
    )
    curator: CuratorConfig = Field(
        default_factory=CuratorConfig,
        validation_alias=AliasChoices("curator"),
        serialization_alias="curator",
    )
    evolution: EvolutionConfig = Field(
        default_factory=EvolutionConfig,
        validation_alias=AliasChoices("evolution"),
        serialization_alias="evolution",
    )


class AgentDefaults(Base):
    """Default agent configuration."""

    workspace: str = "~/.originagent/workspace"
    model: str = "deepseek-chat"
    model_preset: str | None = Field(
        default=None,
        validation_alias=AliasChoices("modelPreset", "model_preset"),
        serialization_alias="modelPreset",
    )
    provider: str = (
        "deepseek"  # Provider name (e.g. "deepseek", "openrouter") or "auto" for auto-detection
    )
    max_tokens: int = 8192
    context_window_tokens: int = 65_536
    context_block_limit: int | None = None
    temperature: float = 0.1
    max_tool_iterations: int = 200
    max_concurrent_subagents: int = Field(default=1, ge=1)
    max_tool_result_chars: int = 16_000
    provider_retry_mode: Literal["standard", "persistent"] = "standard"
    tool_hint_max_length: int = Field(
        default=40,
        ge=20,
        le=500,
        validation_alias=AliasChoices("toolHintMaxLength"),
        serialization_alias="toolHintMaxLength",
    )  # Max characters for tool hint display (e.g. "$ cd …/project && npm test")
    reasoning_effort: str | None = None  # low / medium / high / adaptive / none — LLM thinking effort; None preserves the provider default
    fallback_models: list[FallbackCandidate] = Field(default_factory=list)
    auxiliary: AuxiliaryConfig = Field(default_factory=AuxiliaryConfig)
    domain_packs: DomainPacksConfig = Field(
        default_factory=DomainPacksConfig,
        validation_alias=AliasChoices("domainPacks", "domain_packs"),
        serialization_alias="domainPacks",
    )
    learning: LearningConfig = Field(default_factory=LearningConfig)
    timezone: str = "Asia/Shanghai"  # IANA timezone, e.g. "Asia/Shanghai", "America/New_York"
    bot_name: str = "OriginAgent"  # Display name shown in CLI prompts (e.g. "{name} is thinking...")
    bot_icon: str = "OA"  # Short icon (emoji or text) shown next to the bot name in CLI; "" to omit
    unified_session: bool = False  # Share one session across all channels (single-user multi-device)
    disabled_skills: list[str] = Field(default_factory=list)  # Skill names to exclude from loading (e.g. ["summarize", "skill-creator"])
    session_ttl_minutes: int = Field(
        default=0,
        ge=0,
        validation_alias=AliasChoices("idleCompactAfterMinutes", "sessionTtlMinutes"),
        serialization_alias="idleCompactAfterMinutes",
    )  # Auto-compact idle threshold in minutes (0 = disabled)
    cold_archive_enabled: bool = Field(
        default=True,
        validation_alias=AliasChoices("coldArchiveEnabled", "cold_archive_enabled"),
        serialization_alias="coldArchiveEnabled",
    )  # Preserve trimmed persisted session messages in a local cold archive.
    max_messages: int = Field(
        default=120,
        ge=0,
    )  # Max messages to replay from session history (0 = use default 120, respects token budget)
    allow_agent_initiated_messages: bool = Field(
        default=False,
        validation_alias=AliasChoices(
            "allowAgentInitiatedMessages",
            "allow_agent_initiated_messages",
        ),
        serialization_alias="allowAgentInitiatedMessages",
    )
    active_intent_interval_seconds: int = Field(
        default=30,
        ge=5,
        le=3600,
        validation_alias=AliasChoices(
            "activeIntentIntervalSeconds",
            "active_intent_interval_seconds",
        ),
        serialization_alias="activeIntentIntervalSeconds",
    )
    active_intent_session_cooldown_seconds: int = Field(
        default=300,
        ge=0,
        le=86400,
        validation_alias=AliasChoices(
            "activeIntentSessionCooldownSeconds",
            "active_intent_session_cooldown_seconds",
        ),
        serialization_alias="activeIntentSessionCooldownSeconds",
    )
    active_intent_intent_cooldown_seconds: int = Field(
        default=300,
        ge=0,
        le=86400,
        validation_alias=AliasChoices(
            "activeIntentIntentCooldownSeconds",
            "active_intent_intent_cooldown_seconds",
        ),
        serialization_alias="activeIntentIntentCooldownSeconds",
    )
    active_intent_max_messages_per_session_per_pass: int = Field(
        default=1,
        ge=1,
        le=10,
        validation_alias=AliasChoices(
            "activeIntentMaxMessagesPerSessionPerPass",
            "active_intent_max_messages_per_session_per_pass",
        ),
        serialization_alias="activeIntentMaxMessagesPerSessionPerPass",
    )
    consolidation_ratio: float = Field(
        default=0.5,
        ge=0.1,
        le=0.95,
        validation_alias=AliasChoices("consolidationRatio"),
        serialization_alias="consolidationRatio",
    )  # Consolidation target ratio (0.5 = 50% of budget retained after compression)
    dream: DreamConfig = Field(default_factory=DreamConfig)


class AgentsConfig(Base):
    """Agent configuration."""

    defaults: AgentDefaults = Field(default_factory=AgentDefaults)


class ProviderConfig(Base):
    """LLM provider configuration."""

    api_key: str | None = None
    api_base: str | None = None
    extra_headers: dict[str, str] | None = None  # Custom headers (e.g. APP-Code for AiHubMix)
    extra_body: dict[str, Any] | None = None  # Extra fields merged into every request body


class BedrockProviderConfig(ProviderConfig):
    """AWS Bedrock Runtime provider configuration."""

    region: str | None = None  # AWS region, falls back to AWS_REGION/AWS_DEFAULT_REGION/profile
    profile: str | None = None  # Optional AWS shared config profile


class ProvidersConfig(Base):
    """Configuration for LLM providers."""

    custom: ProviderConfig = Field(default_factory=ProviderConfig)  # Any OpenAI-compatible endpoint
    azure_openai: ProviderConfig = Field(default_factory=ProviderConfig)  # Azure OpenAI (model = deployment name)
    bedrock: BedrockProviderConfig = Field(default_factory=BedrockProviderConfig)  # AWS Bedrock Converse
    anthropic: ProviderConfig = Field(default_factory=ProviderConfig)
    openai: ProviderConfig = Field(default_factory=ProviderConfig)
    openrouter: ProviderConfig = Field(default_factory=ProviderConfig)
    huggingface: ProviderConfig = Field(default_factory=ProviderConfig)
    deepseek: ProviderConfig = Field(default_factory=ProviderConfig)
    groq: ProviderConfig = Field(default_factory=ProviderConfig)
    zhipu: ProviderConfig = Field(default_factory=ProviderConfig)
    dashscope: ProviderConfig = Field(default_factory=ProviderConfig)
    vllm: ProviderConfig = Field(default_factory=ProviderConfig)
    ollama: ProviderConfig = Field(default_factory=ProviderConfig)  # Ollama local models
    lm_studio: ProviderConfig = Field(default_factory=ProviderConfig)  # LM Studio local models
    ovms: ProviderConfig = Field(default_factory=ProviderConfig)  # OpenVINO Model Server (OVMS)
    gemini: ProviderConfig = Field(default_factory=ProviderConfig)
    moonshot: ProviderConfig = Field(default_factory=ProviderConfig)
    minimax: ProviderConfig = Field(default_factory=ProviderConfig)
    minimax_anthropic: ProviderConfig = Field(default_factory=ProviderConfig)  # MiniMax Anthropic endpoint (thinking)
    mistral: ProviderConfig = Field(default_factory=ProviderConfig)
    stepfun: ProviderConfig = Field(default_factory=ProviderConfig)  # Step Fun (阶跃星辰)
    xiaomi_mimo: ProviderConfig = Field(default_factory=ProviderConfig)  # Xiaomi MIMO (小米)
    longcat: ProviderConfig = Field(default_factory=ProviderConfig)  # LongCat
    aihubmix: ProviderConfig = Field(default_factory=ProviderConfig)  # AiHubMix API gateway
    siliconflow: ProviderConfig = Field(default_factory=ProviderConfig)  # SiliconFlow (硅基流动)
    volcengine: ProviderConfig = Field(default_factory=ProviderConfig)  # VolcEngine (火山引擎)
    volcengine_coding_plan: ProviderConfig = Field(default_factory=ProviderConfig)  # VolcEngine Coding Plan
    byteplus: ProviderConfig = Field(default_factory=ProviderConfig)  # BytePlus (VolcEngine international)
    byteplus_coding_plan: ProviderConfig = Field(default_factory=ProviderConfig)  # BytePlus Coding Plan
    openai_codex: ProviderConfig = Field(default_factory=ProviderConfig, exclude=True)  # OpenAI Codex (OAuth)
    github_copilot: ProviderConfig = Field(default_factory=ProviderConfig, exclude=True)  # Github Copilot (OAuth)
    qianfan: ProviderConfig = Field(default_factory=ProviderConfig)  # Qianfan (百度千帆)
    nvidia: ProviderConfig = Field(default_factory=ProviderConfig)  # NVIDIA NIM (nvapi- keys)
    atomic_chat: ProviderConfig = Field(default_factory=ProviderConfig)  # Atomic Chat local models


class HeartbeatConfig(Base):
    """Heartbeat service configuration."""

    enabled: bool = True
    interval_s: int = 30 * 60  # 30 minutes
    keep_recent_messages: int = 8


class ApiConfig(Base):
    """OpenAI-compatible API server configuration."""

    host: str = "127.0.0.1"  # Safer default: local-only bind.
    port: int = 8900
    timeout: float = 120.0  # Per-request timeout in seconds.


class GatewayConfig(Base):
    """Gateway/server configuration."""

    host: str = "127.0.0.1"  # Safer default: local-only bind.
    port: int = 18790
    heartbeat: HeartbeatConfig = Field(default_factory=HeartbeatConfig)


RuntimeProfile = Literal["default", "safe", "household_safe", "local_dev", "automation"]


class RuntimeConfig(Base):
    """Runtime profile configuration."""

    profile: RuntimeProfile = "default"


class PairingConfig(Base):
    """Opt-in DM pairing for approving channel senders."""

    enabled: bool = False
    ttl_seconds: int = Field(
        default=600,
        ge=60,
        le=86_400,
        validation_alias=AliasChoices("ttlSeconds", "ttl_seconds"),
        serialization_alias="ttlSeconds",
    )
    allow_self_approve: bool = Field(
        default=False,
        validation_alias=AliasChoices("allowSelfApprove", "allow_self_approve"),
        serialization_alias="allowSelfApprove",
    )
    approval_channels: list[str] = Field(
        default_factory=lambda: ["cli", "websocket"],
        validation_alias=AliasChoices("approvalChannels", "approval_channels"),
        serialization_alias="approvalChannels",
    )


class SecurityConfig(Base):
    """Security-related runtime controls."""

    pairing: PairingConfig = Field(default_factory=PairingConfig)


class WebSearchConfig(Base):
    """Web search tool configuration."""

    provider: str = "duckduckgo"  # brave, tavily, duckduckgo, searxng, jina, kagi, olostep
    api_key: str = ""
    base_url: str = ""  # SearXNG base URL
    max_results: int = 5
    timeout: int = 30  # Wall-clock timeout (seconds) for search operations


class WebFetchConfig(Base):
    """Web fetch tool configuration."""

    use_jina_reader: bool = True


class SessionSearchConfig(Base):
    """Local session_search retrieval configuration."""

    enabled: bool = True
    backend: Literal["auto", "literal", "sqlite_fts"] = "auto"
    semantic_enabled: bool = True
    max_tool_refresh_ms: int = Field(
        default=500,
        ge=0,
        le=10_000,
        validation_alias=AliasChoices("maxToolRefreshMs", "max_tool_refresh_ms"),
        serialization_alias="maxToolRefreshMs",
    )
    rebuild_on_start: bool = Field(
        default=False,
        validation_alias=AliasChoices("rebuildOnStart", "rebuild_on_start"),
        serialization_alias="rebuildOnStart",
    )


class WebToolsConfig(Base):
    """Web tools configuration."""

    enable: bool = True
    proxy: str | None = (
        None  # HTTP/SOCKS5 proxy URL, e.g. "http://127.0.0.1:7890" or "socks5://127.0.0.1:1080"
    )
    user_agent: str | None = None
    search: WebSearchConfig = Field(default_factory=WebSearchConfig)
    fetch: WebFetchConfig = Field(default_factory=WebFetchConfig)


class ContentReadToolConfig(Base):
    """Platform-aware content_read tool configuration."""

    enabled: bool = False
    providers: list[str] = Field(
        default_factory=lambda: ["generic", "rss", "github", "hackernews"]
    )
    max_chars: int = Field(default=50_000, ge=100)
    use_jina_reader: bool = True
    rss_entry_limit: int = Field(default=10, ge=1, le=50)
    hackernews_comment_limit: int = Field(default=20, ge=0, le=100)


class ExecToolConfig(Base):
    """Shell exec tool configuration."""

    enable: bool = True
    profile: Literal["secure", "local_dev", "disabled"] = "secure"
    allow_unsafe_exec: bool = False
    shell_syntax_policy: Literal["restricted", "shell"] = "restricted"
    timeout: int = 60
    path_append: str = ""
    sandbox: str = ""  # sandbox backend: "" (none) or "bwrap"
    allowed_env_keys: list[str] = Field(default_factory=list)  # Env var names to pass through to subprocess (e.g. ["GOPATH", "JAVA_HOME"])
    allow_patterns: list[str] = Field(default_factory=list)  # Regex patterns that bypass deny_patterns (e.g. [r"rm\s+-rf\s+/tmp/"])
    deny_patterns: list[str] = Field(default_factory=list)  # Extra regex patterns to block (appended to built-in list)

class MCPServerConfig(Base):
    """MCP server connection configuration (stdio or HTTP)."""

    type: Literal["stdio", "sse", "streamableHttp"] | None = None  # auto-detected if omitted
    command: str = ""  # Stdio: command to run (e.g. "npx")
    args: list[str] = Field(default_factory=list)  # Stdio: command arguments
    env: dict[str, str] = Field(default_factory=dict)  # Stdio: extra env vars
    url: str = ""  # HTTP/SSE: endpoint URL
    headers: dict[str, str] = Field(default_factory=dict)  # HTTP/SSE: custom headers
    tool_timeout: int = 30  # seconds before a tool call is cancelled
    enabled_tools: list[str] = Field(default_factory=lambda: ["*"])  # Only register these tools; accepts raw MCP names or wrapped mcp_<server>_<tool> names; ["*"] = all tools; [] = no tools

class MyToolConfig(Base):
    """Self-inspection tool configuration."""

    enable: bool = True  # register the `my` tool (agent runtime state inspection)
    allow_set: bool = False  # let `my` modify loop state (read-only if False)


class ImageGenerationToolConfig(Base):
    """Image generation tool configuration."""

    enabled: bool = False
    provider: str = "openrouter"
    model: str = "openai/gpt-5.4-image-2"
    default_aspect_ratio: str = "1:1"
    default_image_size: str = "1K"
    max_images_per_turn: int = Field(default=4, ge=1, le=8)
    save_dir: str = "generated"


class DeviceToolsConfig(Base):
    """Real-world device tool gateway configuration."""

    enabled: bool = False
    lighting_enabled: bool = False
    mode: Literal["dry_run", "real"] = "dry_run"
    backend: Literal["none", "fake", "lighting_client"] = "none"


class ToolAuditConfig(Base):
    """Privacy-preserving tool call audit configuration."""

    mode: Literal["off", "minimal", "security"] = "minimal"
    security_tools: tuple[str, ...] = (
        "exec",
        "message",
        "web_fetch",
        "content_read",
        "cron",
        "spawn",
        "originagent_device_*",
        "mcp_*",
    )
    security_on_policy_denial: bool = True


class ToolsConfig(Base):
    """Tools configuration."""

    session_search: SessionSearchConfig = Field(
        default_factory=SessionSearchConfig,
        validation_alias=AliasChoices("sessionSearch", "session_search"),
        serialization_alias="sessionSearch",
    )
    web: WebToolsConfig = Field(default_factory=WebToolsConfig)
    content_read: ContentReadToolConfig = Field(default_factory=ContentReadToolConfig)
    exec: ExecToolConfig = Field(default_factory=ExecToolConfig)
    my: MyToolConfig = Field(default_factory=MyToolConfig)
    image_generation: ImageGenerationToolConfig = Field(default_factory=ImageGenerationToolConfig)
    device: DeviceToolsConfig = Field(default_factory=DeviceToolsConfig)
    audit: ToolAuditConfig = Field(default_factory=ToolAuditConfig)
    restrict_to_workspace: bool = False  # restrict all tool access to workspace directory
    mcp_servers: dict[str, MCPServerConfig] = Field(default_factory=dict)
    ssrf_whitelist: list[str] = Field(default_factory=list)  # CIDR ranges to exempt from SSRF blocking (e.g. ["100.64.0.0/10"] for Tailscale)


class Config(BaseSettings):
    """Root configuration for OriginAgent."""

    agents: AgentsConfig = Field(default_factory=AgentsConfig)
    channels: ChannelsConfig = Field(default_factory=ChannelsConfig)
    providers: ProvidersConfig = Field(default_factory=ProvidersConfig)
    api: ApiConfig = Field(default_factory=ApiConfig)
    gateway: GatewayConfig = Field(default_factory=GatewayConfig)
    runtime: RuntimeConfig = Field(default_factory=RuntimeConfig)
    security: SecurityConfig = Field(default_factory=SecurityConfig)
    tools: ToolsConfig = Field(default_factory=ToolsConfig)
    model_presets: dict[str, ModelPresetConfig] = Field(
        default_factory=dict,
        validation_alias=AliasChoices("modelPresets", "model_presets"),
        serialization_alias="modelPresets",
    )

    @model_validator(mode="after")
    def _validate_model_presets(self) -> "Config":
        if "default" in self.model_presets:
            raise ValueError("model_presets must not define reserved preset 'default'")
        name = self.agents.defaults.model_preset
        if name and name != "default" and name not in self.model_presets:
            raise ValueError(f"model_preset {name!r} not found in model_presets")
        for fallback in self.agents.defaults.fallback_models:
            if isinstance(fallback, str) and fallback != "default" and fallback not in self.model_presets:
                raise ValueError(f"fallback_models entry {fallback!r} not found in model_presets")
        for preset_name, preset in self.model_presets.items():
            for fallback in preset.fallback_models:
                if isinstance(fallback, str) and fallback != "default" and fallback not in self.model_presets:
                    raise ValueError(
                        f"model_presets.{preset_name}.fallback_models entry {fallback!r} not found in model_presets"
                    )
        for task_name, task in self.agents.defaults.auxiliary.tasks.items():
            for fallback in task.fallback_models:
                if isinstance(fallback, str) and fallback != "default" and fallback not in self.model_presets:
                    raise ValueError(
                        f"auxiliary.tasks.{task_name}.fallback_models entry {fallback!r} not found in model_presets"
                    )
        return self

    def resolve_default_preset(self) -> ModelPresetConfig:
        defaults = self.agents.defaults
        return ModelPresetConfig(
            model=defaults.model,
            provider=defaults.provider,
            max_tokens=defaults.max_tokens,
            context_window_tokens=defaults.context_window_tokens,
            temperature=defaults.temperature,
            reasoning_effort=defaults.reasoning_effort,
            fallback_models=defaults.fallback_models,
        )

    def resolve_preset(self, name: str | None = None) -> ModelPresetConfig:
        name = name or self.agents.defaults.model_preset or "default"
        if name == "default":
            return self.resolve_default_preset()
        if name not in self.model_presets:
            raise KeyError(f"model_preset {name!r} not found in model_presets")
        return self.model_presets[name]

    @property
    def workspace_path(self) -> Path:
        """Get expanded workspace path."""
        return Path(self.agents.defaults.workspace).expanduser()

    def _match_provider(
        self, model: str | None = None
    ) -> tuple["ProviderConfig | None", str | None]:
        """Match provider config and its registry name. Returns (config, spec_name)."""
        from OriginAgent.providers.registry import PROVIDERS, find_by_name

        forced = self.agents.defaults.provider
        if forced != "auto":
            spec = find_by_name(forced)
            if spec:
                p = getattr(self.providers, spec.name, None)
                if p and (spec.is_oauth or spec.is_local or spec.is_direct or p.api_key):
                    return p, spec.name
            else:
                return None, None

        model_lower = (model or self.agents.defaults.model).lower()
        model_normalized = model_lower.replace("-", "_")
        model_prefix = model_lower.split("/", 1)[0] if "/" in model_lower else ""
        normalized_prefix = model_prefix.replace("-", "_")

        def _kw_matches(kw: str) -> bool:
            kw = kw.lower()
            return kw in model_lower or kw.replace("-", "_") in model_normalized

        # Explicit provider prefix wins — prevents `github-copilot/...codex` matching openai_codex.
        for spec in PROVIDERS:
            p = getattr(self.providers, spec.name, None)
            if p and model_prefix and normalized_prefix == spec.name:
                if spec.is_oauth or spec.is_local or spec.is_direct or p.api_key:
                    return p, spec.name

        # Match by keyword (order follows PROVIDERS registry)
        for spec in PROVIDERS:
            p = getattr(self.providers, spec.name, None)
            if p and any(_kw_matches(kw) for kw in spec.keywords):
                if spec.is_oauth or spec.is_local or spec.is_direct or p.api_key:
                    return p, spec.name

        # Fallback: configured local providers can route models without
        # provider-specific keywords (for example plain "llama3.2" on Ollama).
        # Prefer providers whose detect_by_base_keyword matches the configured api_base
        # (e.g. Ollama's "11434" in "http://localhost:11434") over plain registry order.
        local_fallback: tuple[ProviderConfig, str] | None = None
        for spec in PROVIDERS:
            if not spec.is_local:
                continue
            p = getattr(self.providers, spec.name, None)
            if not (p and p.api_base):
                continue
            if spec.detect_by_base_keyword and spec.detect_by_base_keyword in p.api_base:
                return p, spec.name
            if local_fallback is None:
                local_fallback = (p, spec.name)
        if local_fallback:
            return local_fallback

        # Fallback: gateways first, then others (follows registry order)
        # OAuth providers are NOT valid fallbacks — they require explicit model selection
        for spec in PROVIDERS:
            if spec.is_oauth:
                continue
            p = getattr(self.providers, spec.name, None)
            if p and p.api_key:
                return p, spec.name
        return None, None

    def get_provider(self, model: str | None = None) -> ProviderConfig | None:
        """Get matched provider config (api_key, api_base, extra_headers). Falls back to first available."""
        p, _ = self._match_provider(model)
        return p

    def get_provider_name(self, model: str | None = None) -> str | None:
        """Get the registry name of the matched provider (e.g. "deepseek", "openrouter")."""
        _, name = self._match_provider(model)
        return name

    def get_api_key(self, model: str | None = None) -> str | None:
        """Get API key for the given model. Falls back to first available key."""
        p = self.get_provider(model)
        return p.api_key if p else None

    def get_api_base(self, model: str | None = None) -> str | None:
        """Get API base URL for the given model, falling back to the provider default when present."""
        from OriginAgent.providers.registry import find_by_name

        p, name = self._match_provider(model)
        if p and p.api_base:
            return p.api_base
        if name:
            spec = find_by_name(name)
            if spec and spec.default_api_base:
                return spec.default_api_base
        return None

    model_config = ConfigDict(env_prefix="ORIGINAGENT_", env_nested_delimiter="__")
