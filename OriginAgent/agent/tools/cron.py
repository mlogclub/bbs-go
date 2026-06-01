"""Cron tool for scheduling reminders and tasks."""

from contextvars import ContextVar
from datetime import datetime
from typing import Any

from OriginAgent.agent.tools.base import Tool, tool_parameters
from OriginAgent.agent.tools.schema import (
    BooleanSchema,
    IntegerSchema,
    StringSchema,
    tool_parameters_schema,
)
from OriginAgent.cron.service import CronService
from OriginAgent.cron.types import CronJob, CronJobState, CronSchedule
from OriginAgent.security.capabilities import CapabilitySnapshot
from OriginAgent.security.policy import PolicyDeniedError

_CRON_PARAMETERS = tool_parameters_schema(
    action=StringSchema("Action to perform", enum=["add", "list", "remove"]),
    name=StringSchema(
        "Optional short human-readable label for the job "
        "(e.g., 'weather-monitor', 'daily-standup'). Defaults to first 30 chars of message."
    ),
    message=StringSchema(
        "REQUIRED when action='add'. Instruction for the agent to execute when the job triggers "
        "(e.g., 'Send a reminder to WeChat: xxx' or 'Check system status and report'). "
        "Not used for action='list' or action='remove'."
    ),
    every_seconds=IntegerSchema(0, description="Interval in seconds (for recurring tasks)"),
    cron_expr=StringSchema("Cron expression like '0 9 * * *' (for scheduled tasks)"),
    tz=StringSchema(
        "Optional IANA timezone for cron expressions (e.g. 'America/Vancouver'). "
        "When omitted with cron_expr, the tool's default timezone applies."
    ),
    at=StringSchema(
        "ISO datetime for one-time execution (e.g. '2026-02-12T10:30:00'). "
        "Naive values use the tool's default timezone."
    ),
    deliver=BooleanSchema(
        description="Whether to deliver the execution result to the user channel (default true)",
        default=True,
    ),
    job_id=StringSchema("REQUIRED when action='remove'. Job ID to remove (obtain via action='list')."),
    required=["action"],
    description=(
        "Action-specific parameters: add requires a non-empty message plus one schedule "
        "(every_seconds, cron_expr, or at); remove requires job_id; list only needs action. "
        "Per-action requirements are enforced at runtime (see field descriptions) so the "
        "top-level schema stays compatible with providers (e.g. OpenAI Codex/Responses) that "
        "reject oneOf/anyOf/allOf/enum/not at the root of function parameters."
    ),
    additional_properties=False,
)


@tool_parameters(_CRON_PARAMETERS)
class CronTool(Tool):
    """Tool to schedule reminders and recurring tasks."""

    def __init__(self, cron_service: CronService, default_timezone: str = "UTC"):
        self._cron = cron_service
        self._default_timezone = default_timezone
        self._channel: ContextVar[str] = ContextVar("cron_channel", default="")
        self._chat_id: ContextVar[str] = ContextVar("cron_chat_id", default="")
        self._metadata: ContextVar[dict] = ContextVar("cron_metadata", default={})
        self._session_key: ContextVar[str] = ContextVar("cron_session_key", default="")
        self._in_cron_context: ContextVar[bool] = ContextVar("cron_in_context", default=False)
        self._capability_snapshot: ContextVar[CapabilitySnapshot | None] = ContextVar(
            "cron_capability_snapshot",
            default=None,
        )

    def set_context(
        self, channel: str, chat_id: str,
        metadata: dict | None = None, session_key: str | None = None,
    ) -> None:
        """Set the current session context for delivery."""
        self._channel.set(channel)
        self._chat_id.set(chat_id)
        self._metadata.set(metadata or {})
        self._session_key.set(session_key or f"{channel}:{chat_id}")

    def set_cron_context(self, active: bool):
        """Mark whether the tool is executing inside a cron job callback."""
        return self._in_cron_context.set(active)

    def reset_cron_context(self, token) -> None:
        """Restore previous cron context."""
        self._in_cron_context.reset(token)

    def set_capability_snapshot(self, snapshot: CapabilitySnapshot | None) -> None:
        self._capability_snapshot.set(snapshot)

    @staticmethod
    def _validate_timezone(tz: str) -> str | None:
        from zoneinfo import ZoneInfo

        try:
            ZoneInfo(tz)
        except (KeyError, Exception):
            return f"Error: unknown timezone '{tz}'"
        return None

    def _display_timezone(self, schedule: CronSchedule) -> str:
        """Pick the most human-meaningful timezone for display."""
        return schedule.tz or self._default_timezone

    @staticmethod
    def _format_timestamp(ms: int, tz_name: str) -> str:
        from zoneinfo import ZoneInfo

        dt = datetime.fromtimestamp(ms / 1000, tz=ZoneInfo(tz_name))
        return f"{dt.isoformat()} ({tz_name})"

    @property
    def name(self) -> str:
        return "cron"

    @property
    def description(self) -> str:
        return (
            "Schedule reminders and recurring tasks. Actions: add, list, remove. "
            f"If tz is omitted, cron expressions and naive ISO times default to {self._default_timezone}."
        )

    def validate_params(self, params: dict[str, Any]) -> list[str]:
        errors = super().validate_params(params)
        action = params.get("action")
        if action == "add":
            errors.extend(self._validate_add_params(params))
        if action == "remove" and not str(params.get("job_id") or "").strip():
            errors.append("job_id is required when action='remove'")
        return errors

    def _validate_add_params(self, params: dict[str, Any]) -> list[str]:
        errors: list[str] = []
        if not str(params.get("message") or "").strip():
            errors.append("message is required when action='add'")
        schedules = [
            name for name in ("every_seconds", "cron_expr", "at")
            if params.get(name) not in (None, "")
        ]
        if len(schedules) != 1:
            errors.append("exactly one of every_seconds, cron_expr, or at is required when action='add'")
        if params.get("every_seconds") not in (None, ""):
            every_seconds = params.get("every_seconds")
            if not isinstance(every_seconds, int) or isinstance(every_seconds, bool):
                errors.append("every_seconds must be an integer")
            elif every_seconds < 60:
                errors.append("every_seconds must be >= 60")
        if params.get("tz") and not params.get("cron_expr"):
            errors.append("tz can only be used with cron_expr")
        if params.get("cron_expr"):
            effective_tz = str(params.get("tz") or self._default_timezone)
            if err := self._validate_timezone(effective_tz):
                errors.append(err.removeprefix("Error: "))
            else:
                try:
                    from zoneinfo import ZoneInfo

                    from croniter import croniter

                    croniter(str(params["cron_expr"]), datetime.now(ZoneInfo(effective_tz)))
                except Exception as e:
                    errors.append(f"invalid cron_expr: {e}")
        if params.get("at"):
            try:
                from zoneinfo import ZoneInfo

                default_tz = ZoneInfo(self._default_timezone)
                dt = datetime.fromisoformat(str(params["at"]))
                if dt.tzinfo is None:
                    dt = dt.replace(tzinfo=default_tz)
                else:
                    dt = dt.astimezone(default_tz)
                now = datetime.now(default_tz)
                if dt <= now:
                    errors.append("at must be in the future")
            except ValueError:
                errors.append(
                    f"invalid ISO datetime format '{params['at']}'. Expected format: YYYY-MM-DDTHH:MM:SS"
                )
            except Exception as e:
                errors.append(f"invalid at datetime: {e}")
        return errors

    async def execute(
        self,
        action: str,
        name: str | None = None,
        message: str = "",
        every_seconds: int | None = None,
        cron_expr: str | None = None,
        tz: str | None = None,
        at: str | None = None,
        job_id: str | None = None,
        deliver: bool = True,
        **kwargs: Any,
    ) -> str:
        params = dict(kwargs)
        params["action"] = action
        if name is not None:
            params["name"] = name
        if message or action == "add":
            params["message"] = message
        if every_seconds is not None:
            params["every_seconds"] = every_seconds
        if cron_expr is not None:
            params["cron_expr"] = cron_expr
        if tz is not None:
            params["tz"] = tz
        if at is not None:
            params["at"] = at
        if job_id is not None:
            params["job_id"] = job_id
        if deliver is not True:
            params["deliver"] = deliver

        errors = self.validate_params(params)
        if errors:
            return "Error: Invalid parameters for cron: " + "; ".join(errors)

        if action == "add":
            if self._in_cron_context.get():
                return "Error: cannot schedule new jobs from within a cron job execution"
            return self._add_job(name, message, every_seconds, cron_expr, tz, at, deliver)
        elif action == "list":
            return self._list_jobs()
        elif action == "remove":
            return self._remove_job(job_id)
        return f"Unknown action: {action}"

    def _add_job(
        self,
        name: str | None,
        message: str,
        every_seconds: int | None,
        cron_expr: str | None,
        tz: str | None,
        at: str | None,
        deliver: bool = True,
    ) -> str:
        if not message:
            return (
                "Error: cron action='add' requires a non-empty 'message' parameter "
                "describing what to do when the job triggers "
                "(e.g. the reminder text)."
            )
        channel = self._channel.get()
        chat_id = self._chat_id.get()
        if not channel or not chat_id:
            return "Error: no session context (channel/chat_id)"
        snapshot = self._capability_snapshot.get()
        if snapshot is None:
            raise PolicyDeniedError(
                "Creating cron jobs requires an explicit capability snapshot",
                code="cron_capability_snapshot_missing",
                boundary="cron",
                policy_rule="capability_snapshot_required",
            )
        if not snapshot.can_create_cron:
            raise PolicyDeniedError(
                "Creating cron jobs is not allowed by the current capability snapshot",
                code="cron_capability_denied",
                boundary="cron",
                policy_rule="capability_cron_denied",
            )
        if self._message_requires_high_capability(message):
            raise PolicyDeniedError(
                "This cron job requests high-capability tools and needs an explicit grant",
                code="cron_high_capability_denied",
                boundary="cron",
                policy_rule="cron_high_capability_requires_grant",
            )
        if tz and not cron_expr:
            return "Error: tz can only be used with cron_expr"
        if tz:
            if err := self._validate_timezone(tz):
                return err

        # Build schedule
        delete_after = False
        if every_seconds:
            schedule = CronSchedule(kind="every", every_ms=every_seconds * 1000)
        elif cron_expr:
            effective_tz = tz or self._default_timezone
            if err := self._validate_timezone(effective_tz):
                return err
            schedule = CronSchedule(kind="cron", expr=cron_expr, tz=effective_tz)
        elif at:
            from zoneinfo import ZoneInfo

            try:
                dt = datetime.fromisoformat(at)
            except ValueError:
                return f"Error: invalid ISO datetime format '{at}'. Expected format: YYYY-MM-DDTHH:MM:SS"
            if dt.tzinfo is None:
                if err := self._validate_timezone(self._default_timezone):
                    return err
                dt = dt.replace(tzinfo=ZoneInfo(self._default_timezone))
            at_ms = int(dt.timestamp() * 1000)
            schedule = CronSchedule(kind="at", at_ms=at_ms)
            delete_after = True
        else:
            return "Error: either every_seconds, cron_expr, or at is required"

        job = self._cron.add_job(
            name=name or message[:30],
            schedule=schedule,
            message=message,
            deliver=deliver,
            channel=channel,
            to=chat_id,
            delete_after_run=delete_after,
            channel_meta=self._metadata.get(),
            session_key=self._session_key.get() or None,
            capability_snapshot=CapabilitySnapshot.scheduled_default().to_dict(),
        )
        return f"Created job '{job.name}' (id: {job.id})"

    @staticmethod
    def _message_requires_high_capability(message: str) -> bool:
        lowered = message.casefold()
        risky = (
            "exec",
            "read_file",
            "write_file",
            "edit_file",
            "notebook_edit",
            "message",
            "originagent_device_",
            "spawn",
        )
        return any(token in lowered for token in risky)

    def _format_timing(self, schedule: CronSchedule) -> str:
        """Format schedule as a human-readable timing string."""
        if schedule.kind == "cron":
            tz = f" ({schedule.tz})" if schedule.tz else ""
            return f"cron: {schedule.expr}{tz}"
        if schedule.kind == "every" and schedule.every_ms:
            ms = schedule.every_ms
            if ms % 3_600_000 == 0:
                return f"every {ms // 3_600_000}h"
            if ms % 60_000 == 0:
                return f"every {ms // 60_000}m"
            if ms % 1000 == 0:
                return f"every {ms // 1000}s"
            return f"every {ms}ms"
        if schedule.kind == "at" and schedule.at_ms:
            return f"at {self._format_timestamp(schedule.at_ms, self._display_timezone(schedule))}"
        return schedule.kind

    def _format_state(self, state: CronJobState, schedule: CronSchedule) -> list[str]:
        """Format job run state as display lines."""
        lines: list[str] = []
        display_tz = self._display_timezone(schedule)
        if state.last_run_at_ms:
            info = (
                f"  Last run: {self._format_timestamp(state.last_run_at_ms, display_tz)}"
                f" — {state.last_status or 'unknown'}"
            )
            if state.last_error:
                info += f" ({state.last_error})"
            lines.append(info)
        if state.next_run_at_ms:
            lines.append(f"  Next run: {self._format_timestamp(state.next_run_at_ms, display_tz)}")
        return lines

    @staticmethod
    def _system_job_purpose(job: CronJob) -> str:
        if job.name == "dream":
            return "Dream memory consolidation for long-term memory."
        return "System-managed internal job."

    def _list_jobs(self) -> str:
        jobs = self._cron.list_jobs()
        if not jobs:
            return "No scheduled jobs."
        lines = []
        for j in jobs:
            timing = self._format_timing(j.schedule)
            parts = [f"- {j.name} (id: {j.id}, {timing})"]
            if j.payload.kind == "system_event":
                parts.append(f"  Purpose: {self._system_job_purpose(j)}")
                parts.append("  Protected: visible for inspection, but cannot be removed.")
            parts.extend(self._format_state(j.state, j.schedule))
            lines.append("\n".join(parts))
        return "Scheduled jobs:\n" + "\n".join(lines)

    def _remove_job(self, job_id: str | None) -> str:
        if not job_id:
            return "Error: job_id is required for remove"
        result = self._cron.remove_job(job_id)
        if result == "removed":
            return f"Removed job {job_id}"
        if result == "protected":
            job = self._cron.get_job(job_id)
            if job and job.name == "dream":
                return (
                    "Cannot remove job `dream`.\n"
                    "This is a system-managed Dream memory consolidation job for long-term memory.\n"
                    "It remains visible so you can inspect it, but it cannot be removed."
                )
            return (
                f"Cannot remove job `{job_id}`.\n"
                "This is a protected system-managed cron job."
            )
        return f"Job {job_id} not found"
