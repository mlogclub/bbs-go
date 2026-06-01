"""Message tool for sending messages to users."""

import stat
from contextvars import ContextVar
from pathlib import Path
from typing import Any, Awaitable, Callable

from OriginAgent.agent.tools.base import Tool, tool_parameters
from OriginAgent.agent.tools.schema import ArraySchema, StringSchema, tool_parameters_schema
from OriginAgent.bus.events import OutboundMessage
from OriginAgent.config.paths import get_media_dir, get_workspace_path
from OriginAgent.security.policy import PolicyDeniedError
from OriginAgent.utils.helpers import detect_image_mime


def _is_under(path: Path, root: Path) -> bool:
    try:
        path.relative_to(root)
        return True
    except ValueError:
        return False


_ALLOWED_ATTACHMENT_MIME_TYPES = {"application/pdf"}
_BLOCKED_ATTACHMENT_SUFFIXES = {
    ".json",
    ".jsonl",
    ".env",
    ".pem",
    ".key",
    ".crt",
    ".sqlite",
    ".db",
}
_BLOCKED_ATTACHMENT_NAMES = {
    "config.json",
    "history.jsonl",
    ".dream_cursor",
    ".cursor",
}


def _detect_allowed_attachment_type(path: Path) -> str | None:
    try:
        with path.open("rb") as handle:
            header = handle.read(512)
    except OSError:
        return None
    image_mime = detect_image_mime(header)
    if image_mime and image_mime.startswith("image/"):
        return image_mime
    if header.startswith(b"%PDF-"):
        return "application/pdf"
    return None


@tool_parameters(
    tool_parameters_schema(
        content=StringSchema(
            "Message content for proactive or cross-channel delivery. "
            "Do not use this for a normal reply in the current chat."
        ),
        channel=StringSchema(
            "Optional target channel for cross-channel/proactive delivery. "
            "Do not set this to the current runtime channel for a normal reply."
        ),
        chat_id=StringSchema(
            "Optional target chat/user ID for cross-channel/proactive delivery. "
            "Do not set this to the current runtime chat for a normal reply."
        ),
        media=ArraySchema(
            StringSchema(""),
            description=(
                "Optional list of existing file paths to attach for proactive or cross-channel delivery. "
                "Do not use this to resend generate_image outputs in the current chat."
            ),
            max_items=8,
        ),
        buttons=ArraySchema(
            ArraySchema(StringSchema("Button label", max_length=64), max_items=5),
            description="Optional: inline keyboard buttons as list of rows, each row is list of button labels.",
            max_items=5,
        ),
        required=["content"],
        additional_properties=False,
    )
)
class MessageTool(Tool):
    """Tool to send messages to users on chat channels."""

    def __init__(
        self,
        send_callback: Callable[[OutboundMessage], Awaitable[None]] | None = None,
        default_channel: str = "",
        default_chat_id: str = "",
        default_message_id: str | None = None,
        workspace: str | Path | None = None,
        attachment_roots: list[Path] | None = None,
        max_media_count: int = 8,
        max_media_bytes: int = 50 * 1024 * 1024,
    ):
        self._send_callback = send_callback
        workspace_path = Path(workspace).expanduser() if workspace is not None else get_workspace_path()
        self._workspace = workspace_path.resolve()
        roots = [self._workspace, get_media_dir().resolve(), *(attachment_roots or [])]
        resolved_roots: list[Path] = []
        for root in roots:
            resolved = Path(root).expanduser().resolve()
            if resolved not in resolved_roots:
                resolved_roots.append(resolved)
        self._attachment_roots = tuple(resolved_roots)
        self._max_media_count = max_media_count
        self._max_media_bytes = max_media_bytes
        self._default_channel: ContextVar[str] = ContextVar("message_default_channel", default=default_channel)
        self._default_chat_id: ContextVar[str] = ContextVar("message_default_chat_id", default=default_chat_id)
        self._default_message_id: ContextVar[str | None] = ContextVar(
            "message_default_message_id",
            default=default_message_id,
        )
        self._default_metadata: ContextVar[dict[str, Any]] = ContextVar(
            "message_default_metadata",
            default={},
        )
        self._sent_in_turn_var: ContextVar[bool] = ContextVar("message_sent_in_turn", default=False)
        self._turn_delivered_media_paths_var: ContextVar[tuple[str, ...]] = ContextVar(
            "message_turn_delivered_media_paths",
            default=(),
        )
        self._record_channel_delivery_var: ContextVar[bool] = ContextVar(
            "message_record_channel_delivery",
            default=False,
        )
        self._allow_cross_target_var: ContextVar[bool] = ContextVar(
            "message_allow_cross_target",
            default=False,
        )

    def set_context(
        self,
        channel: str,
        chat_id: str,
        message_id: str | None = None,
        metadata: dict[str, Any] | None = None,
    ) -> None:
        """Set the current message context."""
        self._default_channel.set(channel)
        self._default_chat_id.set(chat_id)
        self._default_message_id.set(message_id)
        self._default_metadata.set(metadata or {})

    def set_send_callback(self, callback: Callable[[OutboundMessage], Awaitable[None]]) -> None:
        """Set the callback for sending messages."""
        self._send_callback = callback

    def start_turn(self) -> None:
        """Reset per-turn send tracking."""
        self._sent_in_turn = False
        self._turn_delivered_media_paths_var.set(())

    def turn_delivered_media_paths(self) -> list[str]:
        """Return local media paths delivered by this tool during the current turn."""
        return list(self._turn_delivered_media_paths_var.get())

    def set_record_channel_delivery(self, active: bool):
        """Mark tool-sent messages as proactive channel deliveries."""
        return self._record_channel_delivery_var.set(active)

    def reset_record_channel_delivery(self, token) -> None:
        """Restore previous proactive delivery recording state."""
        self._record_channel_delivery_var.reset(token)

    def set_cross_target_grant(self, active: bool):
        return self._allow_cross_target_var.set(active)

    def reset_cross_target_grant(self, token) -> None:
        self._allow_cross_target_var.reset(token)

    @property
    def _sent_in_turn(self) -> bool:
        return self._sent_in_turn_var.get()

    @_sent_in_turn.setter
    def _sent_in_turn(self, value: bool) -> None:
        self._sent_in_turn_var.set(value)

    @property
    def name(self) -> str:
        return "message"

    @property
    def description(self) -> str:
        return (
            "Proactively send a message to a user/channel, optionally with file attachments. "
            "Use this for reminders, cross-channel delivery, or explicit proactive sends. "
            "Do not use this for the normal reply in the current chat: answer naturally instead. "
            "If channel/chat_id would target the current runtime conversation, do not call this tool "
            "unless the user explicitly asked you to proactively send an existing file attachment. "
            "When generate_image creates images in the current chat, the final assistant reply "
            "automatically attaches them; do not call message just to announce or resend them. "
            "For proactive attachment delivery, use the 'media' parameter with file paths. "
            "Do NOT use read_file to send files — that only reads content for your own analysis."
        )

    def _validate_buttons(self, buttons: list[list[str]] | None) -> str | None:
        if buttons is None:
            return None
        if not isinstance(buttons, list) or any(
            not isinstance(row, list) or any(not isinstance(label, str) for label in row)
            for row in buttons
        ):
            return "Error: buttons must be a list of list of strings"
        if len(buttons) > 5:
            return "Error: buttons may contain at most 5 rows"
        for row in buttons:
            if len(row) > 5:
                return "Error: each button row may contain at most 5 labels"
            for label in row:
                if len(label) > 64:
                    return "Error: button labels may be at most 64 characters"
        return None

    def _resolve_media_attachment(self, value: str) -> str:
        if not isinstance(value, str) or not value.strip():
            raise ValueError("media entries must be non-empty file paths")
        raw = value.strip()
        if raw.lower().startswith(("http://", "https://")):
            raise ValueError("HTTP/HTTPS media attachments are not allowed")

        candidate = Path(raw).expanduser()
        if not candidate.is_absolute():
            candidate = self._workspace / candidate

        try:
            initial_stat = candidate.lstat()
        except OSError as exc:
            raise ValueError(f"media attachment not found: {value}") from exc
        if stat.S_ISLNK(initial_stat.st_mode):
            raise ValueError("media attachments may not be symlinks")

        try:
            resolved = candidate.resolve(strict=True)
        except OSError as exc:
            raise ValueError(f"media attachment could not be resolved: {value}") from exc

        if not any(_is_under(resolved, root) for root in self._attachment_roots):
            raise ValueError("media attachment is outside allowed attachment roots")

        try:
            resolved_stat = resolved.stat()
        except OSError as exc:
            raise ValueError(f"media attachment could not be inspected: {value}") from exc
        if not stat.S_ISREG(resolved_stat.st_mode):
            raise ValueError("media attachment must be a regular file")
        if resolved_stat.st_size > self._max_media_bytes:
            raise ValueError(
                f"media attachment exceeds maximum size ({self._max_media_bytes} bytes)"
            )
        if resolved.name.casefold() in _BLOCKED_ATTACHMENT_NAMES:
            raise ValueError("attachment type is not allowed by message media policy")
        if resolved.suffix.casefold() in _BLOCKED_ATTACHMENT_SUFFIXES:
            raise ValueError("attachment type is not allowed by message media policy")
        mime = _detect_allowed_attachment_type(resolved)
        if not (
            (mime and mime.startswith("image/"))
            or mime in _ALLOWED_ATTACHMENT_MIME_TYPES
        ):
            raise ValueError("attachment type is not allowed by message media policy")
        return str(resolved)

    async def execute(
        self,
        content: str,
        channel: str | None = None,
        chat_id: str | None = None,
        message_id: str | None = None,
        media: list[str] | None = None,
        buttons: list[list[str]] | None = None,
        **kwargs: Any
    ) -> str:
        from OriginAgent.utils.helpers import strip_think
        content = strip_think(content)

        button_error = self._validate_buttons(buttons)
        if button_error:
            return button_error
        default_channel = self._default_channel.get()
        default_chat_id = self._default_chat_id.get()
        channel = channel or default_channel
        chat_id = chat_id or default_chat_id
        # Only inherit default message_id when targeting the same channel+chat.
        # Cross-chat sends must not carry the original message_id, because
        # some channels (e.g. Feishu) use it to determine the target
        # conversation via their Reply API, which would route the message
        # to the wrong chat entirely.
        has_runtime_target = bool(default_channel and default_chat_id)
        same_target = channel == default_channel and chat_id == default_chat_id
        if same_target:
            message_id = message_id or self._default_message_id.get()
        else:
            message_id = None
            if has_runtime_target and not self._allow_cross_target_var.get():
                raise PolicyDeniedError(
                    "Cross-target message sends require a runtime grant",
                    code="cross_target_denied",
                    boundary="message",
                    policy_rule="message_cross_target_grant_required",
                )

        if not channel or not chat_id:
            return "Error: No target channel/chat specified"

        if not self._send_callback:
            return "Error: Message sending not configured"

        if media:
            if not isinstance(media, list):
                return "Error: media must be a list of file paths"
            if len(media) > self._max_media_count:
                return f"Error: media may contain at most {self._max_media_count} attachments"
            try:
                media = [self._resolve_media_attachment(p) for p in media]
            except ValueError as e:
                return f"Error: {e}"

        metadata = dict(self._default_metadata.get()) if same_target else {}
        if message_id:
            metadata["message_id"] = message_id
        if self._record_channel_delivery_var.get() or media:
            metadata["_record_channel_delivery"] = True

        msg = OutboundMessage(
            channel=channel,
            chat_id=chat_id,
            content=content,
            media=media or [],
            buttons=buttons or [],
            metadata=metadata,
        )

        try:
            await self._send_callback(msg)
            if channel == default_channel and chat_id == default_chat_id:
                self._sent_in_turn = True
            if media:
                existing = self._turn_delivered_media_paths_var.get()
                self._turn_delivered_media_paths_var.set(tuple(dict.fromkeys([*existing, *media])))
            media_info = f" with {len(media)} attachments" if media else ""
            button_info = f" with {sum(len(row) for row in buttons)} button(s)" if buttons else ""
            return f"Message sent to {channel}:{chat_id}{media_info}{button_info}"
        except Exception as e:
            return f"Error sending message: {str(e)}"
