"""Built-in slash command handlers."""

from __future__ import annotations

import asyncio
import os
import shlex
import sys
import time
from contextlib import suppress
from dataclasses import dataclass

from OriginAgent import __version__
from OriginAgent.agent.self_model import SelfModelRenderer, SelfModelService
from OriginAgent.bus.events import OutboundMessage
from OriginAgent.command.router import CommandContext, CommandRouter
from OriginAgent.utils.helpers import build_status_content
from OriginAgent.utils.restart import set_restart_notice_to_env


@dataclass(frozen=True)
class BuiltinCommandSpec:
    command: str
    title: str
    description: str
    icon: str
    arg_hint: str = ""

    def as_dict(self) -> dict[str, str]:
        return {
            "command": self.command,
            "title": self.title,
            "description": self.description,
            "icon": self.icon,
            "arg_hint": self.arg_hint,
        }


BUILTIN_COMMAND_SPECS: tuple[BuiltinCommandSpec, ...] = (
    BuiltinCommandSpec(
        "/new",
        "New chat",
        "Stop the current task and start a fresh conversation.",
        "square-pen",
    ),
    BuiltinCommandSpec(
        "/stop",
        "Stop current task",
        "Cancel the active agent turn for this chat.",
        "square",
    ),
    BuiltinCommandSpec(
        "/restart",
        "Restart OriginAgent",
        "Restart the bot process in place.",
        "rotate-cw",
    ),
    BuiltinCommandSpec(
        "/status",
        "Show status",
        "Display runtime, provider, and channel status.",
        "activity",
    ),
    BuiltinCommandSpec(
        "/self",
        "Show self model",
        "Display a read-only capability and limitation summary.",
        "brain",
    ),
    BuiltinCommandSpec(
        "/goal",
        "Start long-running goal",
        "Treat the request as a sustained goal until complete_goal is called.",
        "target",
        "<goal>",
    ),
    BuiltinCommandSpec(
        "/model",
        "Switch model",
        "Show or switch the active model preset.",
        "badge",
        "[preset]",
    ),
    BuiltinCommandSpec(
        "/pairing",
        "Manage pairing",
        "List, approve, deny, or revoke DM pairing requests.",
        "key-round",
        "[list|approve|deny|revoke]",
    ),
    BuiltinCommandSpec(
        "/mcp",
        "Show MCP servers",
        "List configured MCP servers and registered capabilities.",
        "server",
    ),
    BuiltinCommandSpec(
        "/skill",
        "Show skills",
        "List available agent skills and where they come from.",
        "graduation-cap",
    ),
    BuiltinCommandSpec(
        "/domain",
        "Show domain packs",
        "List installed domain packs and their availability.",
        "boxes",
    ),
    BuiltinCommandSpec(
        "/history",
        "Show conversation history",
        "Print the last N persisted conversation messages.",
        "history",
        "[n]",
    ),
    BuiltinCommandSpec(
        "/reviews",
        "Show learning proposals",
        "List and review background learning proposals.",
        "list-checks",
        "[n|show|apply|reject]",
    ),
    BuiltinCommandSpec(
        "/dream",
        "Run Dream",
        "Manually trigger memory consolidation.",
        "sparkles",
    ),
    BuiltinCommandSpec(
        "/dream-log",
        "Show Dream log",
        "Show what the last Dream consolidation changed.",
        "book-open",
    ),
    BuiltinCommandSpec(
        "/dream-restore",
        "Restore memory",
        "Revert memory to a previous Dream snapshot.",
        "undo-2",
    ),
    BuiltinCommandSpec(
        "/help",
        "Show help",
        "List available slash commands.",
        "circle-help",
    ),
)


def builtin_command_palette() -> list[dict[str, str]]:
    """Return structured command metadata for UI command palettes."""
    return [spec.as_dict() for spec in BUILTIN_COMMAND_SPECS]


async def cmd_stop(ctx: CommandContext) -> OutboundMessage:
    """Cancel all active tasks and subagents for the session."""
    loop = ctx.loop
    msg = ctx.msg
    total = await loop._cancel_active_tasks(ctx.key)
    content = f"Stopped {total} task(s)." if total else "No active task to stop."
    return OutboundMessage(
        channel=msg.channel, chat_id=msg.chat_id, content=content,
        metadata=dict(msg.metadata or {})
    )


async def cmd_restart(ctx: CommandContext) -> OutboundMessage:
    """Restart the process in-place via os.execv."""
    msg = ctx.msg
    set_restart_notice_to_env(
        channel=msg.channel,
        chat_id=msg.chat_id,
        metadata=dict(msg.metadata or {}),
    )

    async def _do_restart():
        await asyncio.sleep(1)
        os.execv(sys.executable, [sys.executable, "-m", "OriginAgent"] + sys.argv[1:])

    asyncio.create_task(_do_restart())
    return OutboundMessage(
        channel=msg.channel, chat_id=msg.chat_id, content="Restarting...",
        metadata=dict(msg.metadata or {})
    )


async def cmd_status(ctx: CommandContext) -> OutboundMessage:
    """Build an outbound status message for a session."""
    loop = ctx.loop
    session = ctx.session or loop.sessions.get_or_create(ctx.key)
    ctx_est = 0
    with suppress(Exception):
        ctx_est, _ = loop.consolidator.estimate_session_prompt_tokens(session)
    if ctx_est <= 0:
        ctx_est = loop._last_usage.get("prompt_tokens", 0)

    # Fetch web search provider usage (best-effort, never blocks the response)
    search_usage_text: str | None = None
    # Never let usage fetch break /status
    with suppress(Exception):
        from OriginAgent.utils.searchusage import fetch_search_usage
        web_cfg = getattr(loop, "web_config", None)
        search_cfg = getattr(web_cfg, "search", None) if web_cfg else None
        if search_cfg is not None:
            provider = getattr(search_cfg, "provider", "duckduckgo")
            api_key = getattr(search_cfg, "api_key", "") or None
            usage = await fetch_search_usage(provider=provider, api_key=api_key)
            search_usage_text = usage.format()
    active_tasks = loop._active_tasks.get(ctx.key, [])
    task_count = sum(1 for t in active_tasks if not t.done())
    with suppress(Exception):
        task_count += loop.subagents.get_running_count_by_session(ctx.key)
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=build_status_content(
            version=__version__, model=loop.model,
            start_time=loop._start_time, last_usage=loop._last_usage,
            context_window_tokens=loop.context_window_tokens,
            session_msg_count=len(session.get_history(max_messages=0)),
            context_tokens_estimate=ctx_est,
            search_usage_text=search_usage_text,
            active_task_count=task_count,
            max_completion_tokens=getattr(
                getattr(loop.provider, "generation", None), "max_tokens", 8192
            ),
        ),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


async def cmd_self(ctx: CommandContext) -> OutboundMessage:
    """Render the current read-only self-model summary."""
    loop = ctx.loop
    service = SelfModelService(
        loop.workspace,
        registry=loop.tools,
        sessions=loop.sessions,
        pending_queues=loop._pending_queues,
        cron_service=loop.cron_service,
        audit_mode=loop._tool_audit_config.mode,
        runtime_profile=getattr(loop, "_runtime_profile", "default"),
        domain_pack_manager=loop.domain_packs,
        background_review_service=loop.background_review,
        curator_service=loop.curator,
        skills_loader=loop.context.skills,
        memory_store=loop.context.memory,
    )
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=SelfModelRenderer().render(service.build()),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


async def cmd_new(ctx: CommandContext) -> OutboundMessage:
    """Stop active task and start a fresh session."""
    loop = ctx.loop
    await loop._cancel_active_tasks(ctx.key)
    session = ctx.session or loop.sessions.get_or_create(ctx.key)
    snapshot = session.messages[session.last_consolidated:]
    session.clear()
    loop.sessions.save(session)
    loop.sessions.invalidate(session.key)
    if snapshot:
        loop._schedule_background(loop.consolidator.archive(snapshot))
    return OutboundMessage(
        channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
        content="New session started.",
        metadata=dict(ctx.msg.metadata or {})
    )


_GOAL_PROMPT_TEMPLATE = """The user invoked `/goal` to start a sustained objective.

Inspect or clarify if needed, then call `long_task` with the refined objective and optional short `ui_summary`. Work proceeds as normal assistant turns using ordinary tools. When the objective is fully done and verified, call `complete_goal` with a brief recap. If the user later cancels or changes direction, still call `complete_goal` with an honest recap before starting a replacement goal. Do not use `long_task` / `complete_goal` for trivial one-shot answers.

Goal:
{goal}
"""


async def cmd_goal(ctx: CommandContext) -> OutboundMessage | None:
    """Rewrite /goal into a normal agent turn that nudges long_task use."""
    goal = ctx.args.strip()
    if not goal:
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content="Usage: /goal <long-running task description>",
            metadata=dict(ctx.msg.metadata or {}),
        )

    active_tasks = ctx.loop._active_tasks.get(ctx.key, [])
    if any(not task.done() for task in active_tasks):
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=(
                "A task is already running in this chat. "
                "Use `/stop` first, then send `/goal <long-running task description>` again."
            ),
            metadata=dict(ctx.msg.metadata or {}),
        )

    ctx.msg.metadata.update(
        {
            "original_command": "/goal",
            "original_content": ctx.raw,
            "goal_started_at": time.time(),
        }
    )
    ctx.msg.content = _GOAL_PROMPT_TEMPLATE.format(goal=goal)
    return None


def _model_preset_names(loop) -> list[str]:
    names = set(getattr(loop, "model_presets", {}) or {})
    if names:
        names.add("default")
    return sorted(names)


def _format_model_status(loop) -> str:
    names = _model_preset_names(loop)
    active = getattr(loop, "model_preset", None) or "default"
    if not names:
        return f"Current model: `{loop.model}`.\n\nNo model presets are configured."
    lines = [f"Current preset: `{active}`", f"Current model: `{loop.model}`", "", "Available presets:"]
    for name in names:
        marker = "*" if name == active else "-"
        lines.append(f"{marker} `{name}`")
    return "\n".join(lines)


async def cmd_model(ctx: CommandContext) -> OutboundMessage:
    """Show or switch runtime model preset."""
    name = ctx.args.strip()
    if not name:
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=_format_model_status(ctx.loop),
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )
    try:
        ctx.loop.set_model_preset(name)
    except Exception as exc:
        names = ", ".join(f"`{n}`" for n in _model_preset_names(ctx.loop)) or "(none)"
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=f"Could not switch model preset: {exc}\n\nAvailable presets: {names}",
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=f"Switched model preset to `{ctx.loop.model_preset}`.\nCurrent model: `{ctx.loop.model}`",
        metadata={
            **dict(ctx.msg.metadata or {}),
            "render_as": "text",
            "model_preset": ctx.loop.model_preset,
            "model": ctx.loop.model,
        },
    )


def _pairing_config(ctx: CommandContext):
    config = getattr(ctx.loop, "pairing_config", None)
    if config is not None:
        return config
    with suppress(Exception):
        from OriginAgent.config.loader import load_config

        return load_config().security.pairing
    from OriginAgent.config.schema import PairingConfig

    return PairingConfig()


def _pairing_command_allowed(ctx: CommandContext, subcommand: str) -> bool:
    config = _pairing_config(ctx)
    if str(ctx.msg.channel) in set(config.approval_channels):
        return True
    if config.allow_self_approve:
        return True
    return False


async def cmd_pairing(ctx: CommandContext) -> OutboundMessage:
    """List, approve, deny or revoke pairing requests when pairing is enabled."""
    from OriginAgent.pairing import PAIRING_COMMAND_META_KEY, handle_pairing_command

    config = _pairing_config(ctx)
    meta = {**dict(ctx.msg.metadata or {}), PAIRING_COMMAND_META_KEY: True, "render_as": "text"}
    if not config.enabled:
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content="Pairing is disabled. Enable `security.pairing.enabled` to use `/pairing`.",
            metadata=meta,
        )
    subcommand = (ctx.args.strip().split() or ["list"])[0]
    if not _pairing_command_allowed(ctx, subcommand):
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content="Pairing approval commands are restricted to trusted approval channels.",
            metadata=meta,
        )
    reply = handle_pairing_command(ctx.msg.channel, ctx.args)
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=reply,
        metadata=meta,
    )


def _format_mcp_capabilities(items: list[dict], *, limit: int = 10) -> str:
    registered = [item for item in items if item.get("status") == "registered"]
    names = [str(item.get("wrapped_name") or item.get("name") or "").strip() for item in registered]
    names = [name for name in names if name]
    if not names:
        return "none"
    shown = names[:limit]
    suffix = f", +{len(names) - limit} more" if len(names) > limit else ""
    return ", ".join(f"`{name}`" for name in shown) + suffix


def _format_mcp_status(loop) -> str:
    configured = getattr(loop, "_mcp_servers", {}) or {}
    snapshot = getattr(loop, "_mcp_snapshot", {}) or {}
    connected = getattr(loop, "_mcp_connected", False)
    connecting = getattr(loop, "_mcp_connecting", False)

    if not configured:
        return (
            "## MCP Servers\n\n"
            "No MCP servers are configured.\n\n"
            "Add entries under `tools.mcp_servers` in the OriginAgent config, then restart the gateway."
        )

    lines = [
        "## MCP Servers",
        "",
        f"- Configured: {len(configured)}",
        f"- Connected: {len(getattr(loop, '_mcp_stacks', {}) or {})}",
        f"- State: {'connecting' if connecting else 'connected' if connected else 'not connected yet'}",
        "",
    ]

    for name in sorted(configured):
        cfg = configured[name]
        snap = snapshot.get(name, {})
        status = snap.get("status") or ("connected" if name in getattr(loop, "_mcp_stacks", {}) else "pending")
        transport = snap.get("transport") or getattr(cfg, "type", "") or ("stdio" if getattr(cfg, "command", "") else "streamableHttp" if getattr(cfg, "url", "") else "unknown")
        tools = snap.get("tools") or []
        resources = snap.get("resources") or []
        prompts = snap.get("prompts") or []
        registered_count = snap.get("registered_count")
        if registered_count is None:
            registered_count = sum(
                1
                for item in [*tools, *resources, *prompts]
                if isinstance(item, dict) and item.get("status") == "registered"
            )

        lines.extend(
            [
                f"### `{name}`",
                f"- Status: {status}",
                f"- Transport: {transport}",
                f"- Registered capabilities: {registered_count}",
                f"- Tools ({sum(1 for item in tools if item.get('status') == 'registered')}/{len(tools)}): {_format_mcp_capabilities(tools)}",
                f"- Resources ({sum(1 for item in resources if item.get('status') == 'registered')}/{len(resources)}): {_format_mcp_capabilities(resources)}",
                f"- Prompts ({sum(1 for item in prompts if item.get('status') == 'registered')}/{len(prompts)}): {_format_mcp_capabilities(prompts)}",
            ]
        )
        error = str(snap.get("error") or "").strip()
        if error:
            lines.append(f"- Error: {error}")
        lines.append("")

    return "\n".join(lines).rstrip()


async def cmd_mcp(ctx: CommandContext) -> OutboundMessage:
    """List configured MCP servers and registered capabilities."""
    loop = ctx.loop
    with suppress(Exception):
        await loop._connect_mcp()
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=_format_mcp_status(loop),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


def _format_skill_status(loop) -> str:
    loader = loop.context.skills
    all_skills = sorted(
        loader.list_skill_records(filter_unavailable=False),
        key=lambda item: (item.get("source", ""), item.get("name", "")),
    )
    if not all_skills:
        return "## Skills\n\nNo skills are available."

    lines = ["## Skills", ""]
    available_count = 0
    workspace_count = 0
    builtin_count = 0
    rows: list[str] = []

    for entry in all_skills:
        name = entry["name"]
        source = entry.get("source", "unknown")
        if source == "workspace":
            workspace_count += 1
        elif source == "builtin":
            builtin_count += 1
        meta = loader._get_skill_meta(name)
        available = loader._check_requirements(meta)
        if available:
            available_count += 1
        desc = loader._get_skill_description(name)
        missing = loader._get_missing_requirements(meta) if not available else ""
        suffix = f" unavailable: {missing}" if missing else " unavailable" if not available else "available"
        lifecycle = entry.get("lifecycle_status") or "unknown"
        verification = entry.get("verification_status") or "unknown"
        always = "always" if entry.get("effective_always") else "manual"
        rows.append(
            f"- `{name}` [{source}] — {desc} "
            f"({suffix}; lifecycle={lifecycle}; verification={verification}; {always})"
        )

    always = loader.get_always_skills()
    lines.extend(
        [
            f"- Total: {len(all_skills)}",
            f"- Available: {available_count}",
            f"- Workspace: {workspace_count}",
            f"- Built-in: {builtin_count}",
        ]
    )
    if always:
        lines.append(f"- Always loaded: {', '.join(f'`{name}`' for name in always)}")
    lines.extend(["", *rows])
    return "\n".join(lines)


def _format_skill_detail(record: dict | None) -> str:
    if record is None:
        return "Skill was not found."
    lines = [
        "## Skill",
        "",
        f"- Name: `{record.get('name') or ''}`",
        f"- Source: {record.get('source') or 'unknown'}",
        f"- Lifecycle: {record.get('lifecycle_status') or 'unknown'}",
        f"- Verification: {record.get('verification_status') or 'unknown'}",
        f"- Always: {'yes' if record.get('effective_always') else 'no'}",
        f"- Path: `{record.get('path') or ''}`",
    ]
    proposal_id = str(record.get("review_proposal_id") or "").strip()
    if proposal_id:
        lines.append(f"- Review proposal: `{proposal_id}`")
    reviewed_at = str(record.get("reviewed_at") or "").strip()
    if reviewed_at:
        lines.append(f"- Reviewed: {reviewed_at}")
    desc = str(record.get("description") or "").strip()
    if desc:
        lines.extend(["", "### Description", "", desc])
    preview = str(record.get("body_preview") or "").strip()
    if preview:
        lines.extend(["", "### Preview", "", preview])
    disabled = str(record.get("disabled_reason") or "").strip()
    if disabled:
        lines.extend(["", f"Action note: {disabled}"])
    return "\n".join(lines)


def _format_skill_lifecycle_result(result: object) -> str:
    if hasattr(result, "to_json"):
        data = result.to_json()
    elif isinstance(result, dict):
        data = result
    else:
        data = {"ok": False, "message": str(result)}
    lines = [
        f"Skill `{data.get('skill_name') or ''}`: {data.get('message') or data.get('status')}",
        f"- Status: {data.get('status') or 'unknown'}",
    ]
    skill = data.get("skill")
    if isinstance(skill, dict):
        lines.append(f"- Verification: {skill.get('verification_status') or 'unknown'}")
        lines.append(f"- Always: {'yes' if skill.get('effective_always') else 'no'}")
        path = skill.get("path")
        if path:
            lines.append(f"- Path: `{path}`")
    error = data.get("error")
    if error:
        lines.append(f"- Error: {error}")
    return "\n".join(lines)


def _skill_read_only_result(name: str, action: str, record: dict | None):
    from OriginAgent.agent.skill_lifecycle import SkillLifecycleResult

    return SkillLifecycleResult(
        skill_name=name,
        status=str((record or {}).get("lifecycle_status") or "missing"),
        action=action,
        ok=False,
        message=str((record or {}).get("disabled_reason") or "Only workspace skills can be changed in P9."),
        skill=record,
        error="read_only" if record else "not_found",
    )


async def cmd_skill(ctx: CommandContext) -> OutboundMessage:
    """List and govern skills."""
    raw_args = ctx.args.strip()
    parts = raw_args.split(maxsplit=2)
    loader = ctx.loop.context.skills
    if parts and parts[0] in {"show", "verify", "activate", "deprecate", "reject", "always"}:
        action = parts[0]
        name = parts[1] if len(parts) >= 2 else ""
        reason = parts[2] if len(parts) >= 3 else ""
        if not name:
            content = (
                "Usage: /skill show <skill_name>, /skill verify <skill_name> [reason], "
                "/skill activate <skill_name> [reason], /skill deprecate <skill_name> [reason], "
                "/skill reject <skill_name> [reason], or /skill always <skill_name> on|off [reason]"
            )
        elif action == "show":
            content = _format_skill_detail(loader.get_skill_record(name))
        else:
            record = loader.get_skill_record(name)
            if record is None or record.get("source") != "workspace":
                content = _format_skill_lifecycle_result(_skill_read_only_result(name, action, record))
            elif action == "always":
                always_parts = reason.split(maxsplit=1)
                setting = always_parts[0].lower() if always_parts else ""
                note = always_parts[1] if len(always_parts) > 1 else ""
                if setting not in {"on", "off"}:
                    content = "Usage: /skill always <skill_name> on|off [reason]"
                else:
                    content = _format_skill_lifecycle_result(
                        loader.lifecycle.transition(
                            name,
                            action="always",
                            enabled=setting == "on",
                            reason=note,
                        )
                    )
            else:
                content = _format_skill_lifecycle_result(
                    loader.lifecycle.transition(name, action=action, reason=reason)
                )
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=content,
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )

    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=_format_skill_status(ctx.loop),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


def _format_domain_status(loop) -> str:
    records = _domain_records(loop, limit=200)
    if records is None:
        return "## Domain Packs\n\nNo domain pack manager is configured."
    if not records:
        return "## Domain Packs\n\nNo domain packs are installed."

    lines = ["## Domain Packs", ""]
    available_count = sum(1 for pack in records if pack["status"] == "available")
    active_count = sum(1 for pack in records if pack["active"])
    invalid_count = sum(1 for pack in records if pack["status"] == "invalid")
    lines.extend(
        [
            f"- Total: {len(records)}",
            f"- Available: {available_count}",
            f"- Active: {active_count}",
            f"- Invalid: {invalid_count}",
            "",
        ]
    )
    for pack in records:
        active = "active" if pack["active"] else "inactive"
        reason = f"; reason: {pack['unavailable_reason']}" if pack.get("unavailable_reason") else ""
        enabled = "enabled" if pack.get("enabled") else "disabled"
        verification = str(pack.get("verification_status") or "unknown")
        override = "; overrides builtin" if pack.get("overrides_builtin") else ""
        lines.append(
            f"- `{pack['id']}` [{pack['source']}] — {pack['name']} "
            f"(status: {pack['status']}, {active}, {enabled}, verification={verification}{override}{reason})"
        )
        if pack.get("skills"):
            available_skills = [skill for skill in pack["skills"] if skill.get("status") == "available"]
            skipped_skills = [skill for skill in pack["skills"] if skill.get("status") == "skipped"]
            available_names = [
                skill.get("virtual_id") or skill.get("id")
                for skill in available_skills
            ]
            available_suffix = (
                " — " + ", ".join(f"`{name}`" for name in available_names[:5])
                if available_names
                else ""
            )
            if len(available_names) > 5:
                available_suffix += f", +{len(available_names) - 5} more"
            lines.append(
                f"  Skills: declared {len(pack['skills'])}, available {len(available_skills)}, "
                f"skipped {len(skipped_skills)}{available_suffix}"
            )
            for skill in skipped_skills[:3]:
                lines.append(f"  - skipped skill `{skill.get('id')}`: {skill.get('unavailable_reason')}")
        if pack.get("workflows"):
            available_workflows = [workflow for workflow in pack["workflows"] if workflow.get("status") == "available"]
            skipped_workflows = [workflow for workflow in pack["workflows"] if workflow.get("status") == "skipped"]
            lines.append(
                f"  Workflows: declared {len(pack['workflows'])}, available {len(available_workflows)}, skipped {len(skipped_workflows)}"
            )
        if pack.get("tools"):
            manifest_skipped = [tool for tool in pack["tools"] if tool.get("status") == "skipped"]
            runtime_records = pack.get("runtime_tools", [])
            registered = [record for record in runtime_records if record.get("status") == "registered"]
            runtime_skipped = [record for record in runtime_records if record.get("status") == "skipped"]
            skipped_count = len(manifest_skipped) + len(runtime_skipped)
            lines.append(
                f"  Tools: declared {len(pack['tools'])}, registered {len(registered)}, "
                f"skipped {skipped_count}"
            )
            for tool in manifest_skipped[:3]:
                lines.append(f"  - skipped tool `{tool.get('id')}`: {tool.get('unavailable_reason')}")
            for record in runtime_skipped[:3]:
                lines.append(f"  - skipped tool `{record.get('tool_id')}`: {record.get('reason')}")
    return "\n".join(lines)


def _domain_governance(loop):
    from OriginAgent.agent.domain_pack_governance import DomainPackGovernanceService

    context = getattr(loop, "context", None)
    manager = getattr(loop, "domain_packs", None) or getattr(context, "domain_packs", None)
    workspace = getattr(manager, "workspace", None) or getattr(loop, "workspace", None)
    if workspace is None:
        return None
    return DomainPackGovernanceService(workspace, domain_pack_manager=manager)


def _domain_records(loop, *, limit: int = 200) -> list[dict] | None:
    service = _domain_governance(loop)
    if service is None:
        return None
    records = service.list_records(limit=limit)
    manager = getattr(service, "_domain_pack_manager", None)
    runtime_by_pack: dict[str, list[dict[str, str]]] = {}
    if manager is not None and hasattr(manager, "domain_tool_runtime_records"):
        for record in manager.domain_tool_runtime_records():
            runtime_by_pack.setdefault(record.pack_id, []).append(
                {
                    "tool_id": record.tool_id,
                    "status": record.status,
                    "reason": record.reason,
                }
            )
    for record in records:
        record["runtime_tools"] = runtime_by_pack.get(str(record.get("id") or ""), [])
    return records


def _format_domain_detail(record: dict | None) -> str:
    if record is None:
        return "Domain pack was not found."
    lines = [
        "## Domain Pack",
        "",
        f"- ID: `{record.get('id') or ''}`",
        f"- Name: {record.get('name') or ''}",
        f"- Version: {record.get('version') or ''}",
        f"- Source: {record.get('source') or ''}",
        f"- Status: {record.get('status') or ''}",
        f"- Enabled: {'yes' if record.get('enabled') else 'no'}",
        f"- Active: {'yes' if record.get('active') else 'no'}",
        f"- Active in config: {'yes' if record.get('active_requested') else 'no'}",
        f"- Verification: {record.get('verification_status') or 'unknown'}",
        f"- Overrides builtin: {'yes' if record.get('overrides_builtin') else 'no'}",
        f"- Path: `{record.get('path') or ''}`",
    ]
    description = str(record.get("description") or "").strip()
    if description:
        lines.extend(["", "### Description", "", description])
    lines.extend(["", "### Validation", "", str(record.get("validation_summary") or "")])
    if record.get("skills"):
        skill_names = ", ".join(f"`{item.get('id')}`" for item in record["skills"])
        lines.extend(["", f"Skills: {skill_names}"])
    if record.get("workflows"):
        workflow_names = ", ".join(f"`{item.get('id')}`" for item in record["workflows"])
        lines.extend(["", f"Workflows: {workflow_names}"])
    if record.get("dependencies", {}).get("packs"):
        dependency_names = ", ".join(f"`{item}`" for item in record["dependencies"]["packs"])
        lines.extend(["", f"Dependencies: {dependency_names}"])
    last_eval = record.get("last_eval_result")
    if isinstance(last_eval, dict):
        lines.extend([
            "",
            "### Last Eval",
            "",
            f"- Status: {last_eval.get('status') or 'unknown'}",
            f"- Checks: {len(last_eval.get('checks') or [])}",
            f"- Warnings: {len(last_eval.get('warnings') or [])}",
            f"- Errors: {len(last_eval.get('errors') or [])}",
        ])
    return "\n".join(lines)


def _format_domain_result(result: object) -> str:
    if hasattr(result, "to_json"):
        data = result.to_json()
    elif isinstance(result, dict):
        data = result
    else:
        data = {"ok": False, "message": str(result)}
    lines = [
        f"Domain pack `{data.get('pack_id') or ''}`: {data.get('message') or data.get('status')}",
        f"- Status: {data.get('status') or 'unknown'}",
        f"- Action: {data.get('action') or 'unknown'}",
    ]
    artifact = data.get("artifact")
    if isinstance(artifact, dict):
        if artifact.get("skill_name"):
            lines.append(f"- Skill: `{artifact.get('skill_name')}`")
        if artifact.get("workflow_name"):
            lines.append(f"- Workflow: `{artifact.get('workflow_name')}`")
        if artifact.get("path"):
            lines.append(f"- Path: `{artifact.get('path')}`")
    eval_result = data.get("eval_result")
    if isinstance(eval_result, dict):
        lines.append(f"- Eval status: {eval_result.get('status') or 'unknown'}")
        lines.append(f"- Eval checks: {len(eval_result.get('checks') or [])}")
        if eval_result.get("errors"):
            lines.extend(f"  - {item}" for item in eval_result["errors"][:5])
    error = data.get("error")
    if error:
        lines.append(f"- Error: {error}")
    return "\n".join(lines)


def _domain_command_args(raw_args: str) -> list[str]:
    if not raw_args.strip():
        return []
    try:
        return shlex.split(raw_args, posix=False)
    except ValueError:
        return raw_args.split()


async def cmd_domain(ctx: CommandContext) -> OutboundMessage:
    """List and manage installed domain packs."""
    raw_args = ctx.args.strip()
    args = _domain_command_args(raw_args)
    service = _domain_governance(ctx.loop)
    if service is None:
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content="Domain pack governance is not available.",
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )

    if args and args[0] in {
        "show",
        "install",
        "upgrade",
        "enable",
        "disable",
        "activate",
        "deactivate",
        "uninstall",
        "eval",
    }:
        action = args[0]
        content = ""
        if action == "show":
            pack_id = args[1] if len(args) >= 2 else ""
            record = service.get_record(pack_id) if pack_id else None
            content = _format_domain_detail(record) if record else "Domain pack was not found."
        elif action == "install":
            if len(args) < 2:
                content = "Usage: /domains install <source_path> [reason]"
            else:
                source_path = args[1]
                reason = args[2] if len(args) >= 3 else ""
                content = _format_domain_result(service.install(source_path, reason=reason))
        elif action == "upgrade":
            if len(args) < 3:
                content = "Usage: /domains upgrade <pack_id> <source_path> [reason]"
            else:
                pack_id = args[1]
                source_path = args[2]
                reason = args[3] if len(args) >= 4 else ""
                content = _format_domain_result(service.upgrade(pack_id, source_path, reason=reason))
        elif action == "enable":
            if len(args) < 2:
                content = "Usage: /domains enable <pack_id> [reason]"
            else:
                pack_id = args[1]
                reason = args[2] if len(args) >= 3 else ""
                content = _format_domain_result(service.set_enabled(pack_id, enabled=True, reason=reason))
        elif action == "disable":
            if len(args) < 2:
                content = "Usage: /domains disable <pack_id> [reason]"
            else:
                pack_id = args[1]
                reason = args[2] if len(args) >= 3 else ""
                content = _format_domain_result(service.set_enabled(pack_id, enabled=False, reason=reason))
        elif action == "activate":
            if len(args) < 2:
                content = "Usage: /domains activate <pack_id> [reason]"
            else:
                pack_id = args[1]
                reason = args[2] if len(args) >= 3 else ""
                content = _format_domain_result(service.set_active(pack_id, active=True, reason=reason))
        elif action == "deactivate":
            if len(args) < 2:
                content = "Usage: /domains deactivate <pack_id> [reason]"
            else:
                pack_id = args[1]
                reason = args[2] if len(args) >= 3 else ""
                content = _format_domain_result(service.set_active(pack_id, active=False, reason=reason))
        elif action == "uninstall":
            if len(args) < 2:
                content = "Usage: /domains uninstall <pack_id> [reason]"
            else:
                pack_id = args[1]
                reason = args[2] if len(args) >= 3 else ""
                content = _format_domain_result(service.uninstall(pack_id, reason=reason))
        elif action == "eval":
            if len(args) < 2:
                content = "Usage: /domains eval <pack_id>"
            else:
                content = _format_domain_result(service.eval_pack(args[1]))
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=content,
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )

    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=_format_domain_status(ctx.loop),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


async def cmd_dream(ctx: CommandContext) -> OutboundMessage:
    """Manually trigger a Dream consolidation run."""
    import time

    loop = ctx.loop
    msg = ctx.msg

    async def _run_dream():
        t0 = time.monotonic()
        try:
            did_work = await loop.dream.run()
            elapsed = time.monotonic() - t0
            if did_work:
                content = f"Dream completed in {elapsed:.1f}s."
            else:
                content = "Dream: nothing to process."
        except Exception as e:
            elapsed = time.monotonic() - t0
            content = f"Dream failed after {elapsed:.1f}s: {e}"
        await loop.bus.publish_outbound(OutboundMessage(
            channel=msg.channel, chat_id=msg.chat_id, content=content,
        ))

    asyncio.create_task(_run_dream())
    return OutboundMessage(
        channel=msg.channel, chat_id=msg.chat_id, content="Dreaming...",
    )


def _extract_changed_files(diff: str) -> list[str]:
    """Extract changed file paths from a unified diff."""
    files: list[str] = []
    seen: set[str] = set()
    for line in diff.splitlines():
        if not line.startswith("diff --git "):
            continue
        parts = line.split()
        if len(parts) < 4:
            continue
        path = parts[3]
        if path.startswith("b/"):
            path = path[2:]
        if path in seen:
            continue
        seen.add(path)
        files.append(path)
    return files


def _format_changed_files(diff: str) -> str:
    files = _extract_changed_files(diff)
    if not files:
        return "No tracked memory files changed."
    return ", ".join(f"`{path}`" for path in files)


def _format_dream_log_content(commit, diff: str, *, requested_sha: str | None = None) -> str:
    files_line = _format_changed_files(diff)
    lines = [
        "## Dream Update",
        "",
        "Here is the selected Dream memory change." if requested_sha else "Here is the latest Dream memory change.",
        "",
        f"- Commit: `{commit.sha}`",
        f"- Time: {commit.timestamp}",
        f"- Changed files: {files_line}",
    ]
    if diff:
        lines.extend([
            "",
            f"Use `/dream-restore {commit.sha}` to undo this change.",
            "",
            "```diff",
            diff.rstrip(),
            "```",
        ])
    else:
        lines.extend([
            "",
            "Dream recorded this version, but there is no file diff to display.",
        ])
    return "\n".join(lines)


def _format_dream_restore_list(commits: list) -> str:
    lines = [
        "## Dream Restore",
        "",
        "Choose a Dream memory version to restore. Latest first:",
        "",
    ]
    for c in commits:
        lines.append(f"- `{c.sha}` {c.timestamp} - {c.message.splitlines()[0]}")
    lines.extend([
        "",
        "Preview a version with `/dream-log <sha>` before restoring it.",
        "Restore a version with `/dream-restore <sha>`.",
    ])
    return "\n".join(lines)


async def cmd_dream_log(ctx: CommandContext) -> OutboundMessage:
    """Show what the last Dream changed.

    Default: diff of the latest commit (HEAD~1 vs HEAD).
    With /dream-log <sha>: diff of that specific commit.
    """
    store = ctx.loop.consolidator.store
    git = store.git

    if not git.is_initialized():
        if store.get_last_dream_cursor() == 0:
            msg = "Dream has not run yet. Run `/dream`, or wait for the next scheduled Dream cycle."
        else:
            msg = "Dream history is not available because memory versioning is not initialized."
        return OutboundMessage(
            channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
            content=msg, metadata={"render_as": "text"},
        )

    args = ctx.args.strip()

    if args:
        # Show diff of a specific commit
        sha = args.split()[0]
        result = git.show_commit_diff(sha)
        if not result:
            content = (
                f"Couldn't find Dream change `{sha}`.\n\n"
                "Use `/dream-restore` to list recent versions, "
                "or `/dream-log` to inspect the latest one."
            )
        else:
            commit, diff = result
            content = _format_dream_log_content(commit, diff, requested_sha=sha)
    else:
        # Default: show the latest commit's diff
        commits = git.log(max_entries=1)
        result = git.show_commit_diff(commits[0].sha) if commits else None
        if result:
            commit, diff = result
            content = _format_dream_log_content(commit, diff)
        else:
            content = "Dream memory has no saved versions yet."

    return OutboundMessage(
        channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
        content=content, metadata={"render_as": "text"},
    )


async def cmd_dream_restore(ctx: CommandContext) -> OutboundMessage:
    """Restore memory files from a previous dream commit.

    Usage:
        /dream-restore          — list recent commits
        /dream-restore <sha>    — revert a specific commit
    """
    store = ctx.loop.consolidator.store
    git = store.git
    if not git.is_initialized():
        return OutboundMessage(
            channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
            content="Dream history is not available because memory versioning is not initialized.",
        )

    args = ctx.args.strip()
    if not args:
        # Show recent commits for the user to pick
        commits = git.log(max_entries=10)
        if not commits:
            content = "Dream memory has no saved versions to restore yet."
        else:
            content = _format_dream_restore_list(commits)
    else:
        sha = args.split()[0]
        result = git.show_commit_diff(sha)
        changed_files = _format_changed_files(result[1]) if result else "the tracked memory files"
        new_sha = git.revert(sha)
        if new_sha:
            content = (
                f"Restored Dream memory to the state before `{sha}`.\n\n"
                f"- New safety commit: `{new_sha}`\n"
                f"- Restored files: {changed_files}\n\n"
                f"Use `/dream-log {new_sha}` to inspect the restore diff."
            )
        else:
            content = (
                f"Couldn't restore Dream change `{sha}`.\n\n"
                "It may not exist, or it may be the first saved version with no earlier state to restore."
            )
    return OutboundMessage(
        channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
        content=content, metadata={"render_as": "text"},
    )


_HISTORY_DEFAULT_COUNT = 10
_HISTORY_MAX_COUNT = 50
_HISTORY_MAX_CONTENT_CHARS = 200


def _format_history_message(msg: dict) -> str | None:
    """Format a single history message for display. Returns None to skip."""
    role = msg.get("role")
    if role not in ("user", "assistant"):
        return None
    content = msg.get("content") or ""
    if isinstance(content, list):
        parts = [b.get("text", "") for b in content if isinstance(b, dict) and b.get("type") == "text"]
        content = " ".join(parts)
    content = str(content).strip()
    if not content:
        return None
    if len(content) > _HISTORY_MAX_CONTENT_CHARS:
        content = content[:_HISTORY_MAX_CONTENT_CHARS] + "…"
    label = "👤 You" if role == "user" else "🤖 Bot"
    return f"{label}: {content}"


async def cmd_history(ctx: CommandContext) -> OutboundMessage:
    """Show the last N messages of the current session (default 10, max 50).

    Usage: /history [count]
    """
    count = _HISTORY_DEFAULT_COUNT
    if ctx.args.strip():
        try:
            count = max(1, min(int(ctx.args.strip()), _HISTORY_MAX_COUNT))
        except ValueError:
            return OutboundMessage(
                channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
                content="Usage: /history [count] — e.g. /history 5 (default: 10, max: 50)",
                metadata=dict(ctx.msg.metadata or {}),
            )

    session = ctx.session or ctx.loop.sessions.get_or_create(ctx.key)
    history = session.get_history(max_messages=0)
    visible = [_format_history_message(m) for m in history]
    visible = [m for m in visible if m is not None]
    recent = visible[-count:]

    if not recent:
        return OutboundMessage(
            channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
            content="No conversation history yet.",
            metadata=dict(ctx.msg.metadata or {}),
        )

    header = f"Last {len(recent)} message(s):\n"
    return OutboundMessage(
        channel=ctx.msg.channel, chat_id=ctx.msg.chat_id,
        content=header + "\n".join(recent),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


_REVIEWS_DEFAULT_COUNT = 10
_REVIEWS_MAX_COUNT = 50


def _reviews_store(ctx: CommandContext):
    service = getattr(ctx.loop, "background_review", None)
    return getattr(service, "store", None), service


def _format_review_record(record: dict) -> str:
    proposal_id = str(record.get("id") or "unknown")
    origin = str(record.get("origin") or "background_review")
    proposal_type = str(record.get("proposal_type") or record.get("type") or "unknown")
    domain_id = str(record.get("domain_id") or "core")
    status = str(record.get("status") or "pending")
    title = str(record.get("title") or "(untitled)")
    content = str(record.get("content") or "").strip()
    if len(content) > 260:
        content = content[:260] + "..."
    created_at = str(record.get("created_at") or "")
    confidence = record.get("confidence")
    confidence_text = f", confidence={confidence}" if confidence is not None else ""
    lines = [
        f"- `{proposal_id}` [{status}] {proposal_type}/{domain_id} ({origin}): {title}",
        f"  created={created_at}{confidence_text}",
    ]
    subject = str(record.get("subject_label") or "").strip()
    suggested_action = str(record.get("suggested_action") or "").strip()
    if subject:
        lines.append(f"  subject={subject}")
    if suggested_action:
        lines.append(f"  suggested_action={suggested_action}")
    if content:
        lines.append(f"  {content}")
    return "\n".join(lines)


def _format_review_detail(record: dict) -> str:
    lines = [
        "## Background Review Proposal",
        "",
        f"- ID: `{record.get('id') or 'unknown'}`",
        f"- Status: {record.get('status') or 'pending'}",
        f"- Origin: {record.get('origin') or 'background_review'}",
        f"- Type: {record.get('proposal_type') or record.get('type') or 'unknown'}",
        f"- Domain: {record.get('domain_id') or 'core'}",
        f"- Created: {record.get('created_at') or ''}",
        f"- Session: {record.get('session_key') or ''}",
        "",
        f"### {record.get('title') or '(untitled)'}",
        "",
        str(record.get("content") or "").strip() or "(empty)",
    ]
    subject = str(record.get("subject_label") or "").strip()
    if subject:
        lines.insert(7, f"- Subject: {subject}")
    suggested_action = str(record.get("suggested_action") or "").strip()
    if suggested_action:
        lines.insert(8 if subject else 7, f"- Suggested Action: {suggested_action}")
    rationale = str(record.get("rationale") or "").strip()
    if rationale:
        lines.extend(["", "### Rationale", "", rationale])
    evidence = record.get("evidence")
    if isinstance(evidence, list) and evidence:
        lines.extend(["", "### Evidence", ""])
        lines.extend(f"- {item}" for item in evidence if str(item).strip())
    review_reason = str(record.get("review_reason") or "").strip()
    if review_reason:
        lines.extend(["", "### Review Reason", "", review_reason])
    fact_id = record.get("applied_fact_id")
    if isinstance(fact_id, str) and fact_id:
        lines.extend(["", f"Applied fact: `{fact_id}`"])
    skill_path = record.get("applied_skill_path")
    if isinstance(skill_path, str) and skill_path:
        lines.extend(["", f"Applied skill: `{skill_path}`"])
    workflow_path = record.get("applied_workflow_path")
    if isinstance(workflow_path, str) and workflow_path:
        lines.extend(["", f"Applied workflow: `{workflow_path}`"])
    unsupported = str(record.get("unsupported_reason") or "").strip()
    if unsupported and not record.get("can_apply"):
        lines.extend(["", f"Apply: {unsupported}"])
    return "\n".join(lines)


def _format_review_result(result: object) -> str:
    if hasattr(result, "to_json"):
        data = result.to_json()
    elif isinstance(result, dict):
        data = result
    else:
        data = {"ok": False, "message": str(result)}
    lines = [
        f"Review `{data.get('proposal_id') or ''}`: {data.get('message') or data.get('status')}",
        f"- Status: {data.get('status') or 'unknown'}",
    ]
    fact_id = data.get("fact_id")
    if fact_id:
        lines.append(f"- Fact: `{fact_id}`")
    artifact = data.get("artifact")
    if isinstance(artifact, dict):
        skill_name = artifact.get("skill_name")
        artifact_path = artifact.get("path")
        workflow_name = artifact.get("workflow_name")
        if skill_name:
            lines.append(f"- Skill: `{skill_name}`")
        if workflow_name:
            lines.append(f"- Workflow: `{workflow_name}`")
        if artifact_path:
            lines.append(f"- Path: `{artifact_path}`")
    error = data.get("error")
    if error:
        lines.append(f"- Error: {error}")
    return "\n".join(lines)


async def cmd_reviews(ctx: CommandContext) -> OutboundMessage:
    """Show and manage pending background review proposals.

    Usage:
      /reviews [count]
      /reviews show <proposal_id>
      /reviews apply|approve <proposal_id>
      /reviews reject|defer <proposal_id> [reason]
    """
    raw_args = ctx.args.strip()
    args = raw_args.split(maxsplit=2)
    store, service = _reviews_store(ctx)
    if store is None:
        content = "Background review proposal store is not available."
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=content,
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )

    if args and args[0] in {"show", "apply", "approve", "reject", "defer"}:
        action = args[0]
        proposal_id = args[1] if len(args) >= 2 else ""
        reason = args[2] if len(args) >= 3 else ""
        if not proposal_id:
            content = (
                "Usage: /reviews show <proposal_id>, "
                "/reviews apply <proposal_id>, "
                "/reviews reject <proposal_id> [reason], "
                "or /reviews defer <proposal_id> [reason]"
            )
        elif action == "show":
            record = store.get(proposal_id)
            content = _format_review_detail(record) if record else "Review proposal was not found."
        elif action in {"apply", "approve"}:
            content = _format_review_result(store.apply(proposal_id, reason=reason))
        elif action == "reject":
            content = _format_review_result(store.reject(proposal_id, reason=reason))
        else:
            content = _format_review_result(store.defer(proposal_id, reason=reason))
        return OutboundMessage(
            channel=ctx.msg.channel,
            chat_id=ctx.msg.chat_id,
            content=content,
            metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
        )

    count = _REVIEWS_DEFAULT_COUNT
    if raw_args:
        try:
            count = max(1, min(int(raw_args), _REVIEWS_MAX_COUNT))
        except ValueError:
            return OutboundMessage(
                channel=ctx.msg.channel,
                chat_id=ctx.msg.chat_id,
                content=(
                    "Usage: /reviews [count], /reviews show <proposal_id>, "
                    "/reviews apply <proposal_id>, /reviews reject <proposal_id> [reason]"
                ),
                metadata=dict(ctx.msg.metadata or {}),
            )

    records = store.recent(count)
    if not records:
        enabled = bool(getattr(service, "enabled", False))
        suffix = " It is currently disabled." if not enabled else ""
        content = "No background review proposals yet." + suffix
    else:
        stats = store.stats()
        lines = [
            "## Background Review Proposals",
            "",
            f"- Showing: {len(records)}",
            f"- Total stored: {stats.get('proposal_count', 0)}",
            f"- Pending: {stats.get('pending_count', 0)}",
            "",
        ]
        lines.extend(_format_review_record(record) for record in records)
        content = "\n".join(lines)
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=content,
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


async def cmd_help(ctx: CommandContext) -> OutboundMessage:
    """Return available slash commands."""
    return OutboundMessage(
        channel=ctx.msg.channel,
        chat_id=ctx.msg.chat_id,
        content=build_help_text(),
        metadata={**dict(ctx.msg.metadata or {}), "render_as": "text"},
    )


def build_help_text() -> str:
    """Build canonical help text shared across channels."""
    lines = ["OriginAgent commands:"]
    for spec in BUILTIN_COMMAND_SPECS:
        command = spec.command
        if spec.arg_hint:
            command = f"{command} {spec.arg_hint}"
        lines.append(f"{command} — {spec.description}")
    return "\n".join(lines)


def register_builtin_commands(router: CommandRouter) -> None:
    """Register the default set of slash commands."""
    router.priority("/stop", cmd_stop)
    router.priority("/restart", cmd_restart)
    router.priority("/status", cmd_status)
    router.exact("/new", cmd_new)
    router.exact("/status", cmd_status)
    router.exact("/self", cmd_self)
    router.exact("/goal", cmd_goal)
    router.prefix("/goal ", cmd_goal)
    router.exact("/model", cmd_model)
    router.prefix("/model ", cmd_model)
    router.exact("/pairing", cmd_pairing)
    router.prefix("/pairing ", cmd_pairing)
    router.exact("/mcp", cmd_mcp)
    router.exact("/skill", cmd_skill)
    router.exact("/skills", cmd_skill)
    router.prefix("/skill ", cmd_skill)
    router.prefix("/skills ", cmd_skill)
    router.exact("/domain", cmd_domain)
    router.exact("/domains", cmd_domain)
    router.prefix("/domain ", cmd_domain)
    router.prefix("/domains ", cmd_domain)
    router.exact("/history", cmd_history)
    router.prefix("/history ", cmd_history)
    router.exact("/reviews", cmd_reviews)
    router.prefix("/reviews ", cmd_reviews)
    router.exact("/dream", cmd_dream)
    router.exact("/dream-log", cmd_dream_log)
    router.prefix("/dream-log ", cmd_dream_log)
    router.exact("/dream-restore", cmd_dream_restore)
    router.prefix("/dream-restore ", cmd_dream_restore)
    router.exact("/help", cmd_help)
