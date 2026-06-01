"""Base channel interface for chat platforms."""

from __future__ import annotations

from abc import ABC, abstractmethod
from pathlib import Path
from typing import Any

from loguru import logger

from OriginAgent.bus.events import InboundMessage, OutboundMessage
from OriginAgent.bus.queue import MessageBus
from OriginAgent.config.schema import PairingConfig
from OriginAgent.pairing import PAIRING_CODE_META_KEY, format_pairing_reply, generate_code, is_approved


class BaseChannel(ABC):
    """
    Abstract base class for chat channel implementations.

    Each channel (Telegram, Discord, etc.) should implement this interface
    to integrate with the OriginAgent message bus.
    """

    name: str = "base"
    display_name: str = "Base"
    transcription_provider: str = "groq"
    transcription_api_key: str = ""
    transcription_api_base: str = ""
    transcription_language: str | None = None
    send_progress: bool = True
    send_tool_hints: bool = False
    show_reasoning: bool = True

    def __init__(self, config: Any, bus: MessageBus):
        """
        Initialize the channel.

        Args:
            config: Channel-specific configuration.
            bus: The message bus for communication.
        """
        self.config = config
        self.logger = logger.bind(channel=self.name)
        self.bus = bus
        self._running = False
        self.pairing_config = PairingConfig()

    async def transcribe_audio(self, file_path: str | Path) -> str:
        """Transcribe an audio file via Whisper (OpenAI or Groq). Returns empty string on failure."""
        if not self.transcription_api_key:
            return ""
        try:
            if self.transcription_provider == "openai":
                from OriginAgent.providers.transcription import OpenAITranscriptionProvider
                provider = OpenAITranscriptionProvider(
                    api_key=self.transcription_api_key,
                    api_base=self.transcription_api_base or None,
                    language=self.transcription_language or None,
                )
            else:
                from OriginAgent.providers.transcription import GroqTranscriptionProvider
                provider = GroqTranscriptionProvider(
                    api_key=self.transcription_api_key,
                    api_base=self.transcription_api_base or None,
                    language=self.transcription_language or None,
                )
            return await provider.transcribe(file_path)
        except Exception:
            self.logger.exception("Audio transcription failed")
            return ""

    async def login(self, force: bool = False) -> bool:
        """
        Perform channel-specific interactive login (e.g. QR code scan).

        Args:
            force: If True, ignore existing credentials and force re-authentication.

        Returns True if already authenticated or login succeeds.
        Override in subclasses that support interactive login.
        """
        return True

    @abstractmethod
    async def start(self) -> None:
        """
        Start the channel and begin listening for messages.

        This should be a long-running async task that:
        1. Connects to the chat platform
        2. Listens for incoming messages
        3. Forwards messages to the bus via _handle_message()
        """
        pass

    @abstractmethod
    async def stop(self) -> None:
        """Stop the channel and clean up resources."""
        pass

    @abstractmethod
    async def send(self, msg: OutboundMessage) -> None:
        """
        Send a message through this channel.

        Args:
            msg: The message to send.

        Implementations should raise on delivery failure so the channel manager
        can apply any retry policy in one place.
        """
        pass

    async def send_delta(self, chat_id: str, delta: str, metadata: dict[str, Any] | None = None) -> None:
        """Deliver a streaming text chunk.

        Override in subclasses to enable streaming. Implementations should
        raise on delivery failure so the channel manager can retry.

        Streaming contract: ``_stream_delta`` is a chunk, ``_stream_end`` ends
        the current segment, and stateful implementations must key buffers by
        ``_stream_id`` rather than only by ``chat_id``.
        """
        pass

    async def send_reasoning_delta(
        self, chat_id: str, delta: str, metadata: dict[str, Any] | None = None
    ) -> None:
        """Stream a chunk of model reasoning/thinking content.

        Default is no-op. Channels with a native low-emphasis primitive
        override this to render reasoning as a subordinate trace.
        """
        return

    async def send_reasoning_end(
        self, chat_id: str, metadata: dict[str, Any] | None = None
    ) -> None:
        """Mark the end of a reasoning stream segment."""
        return

    async def send_reasoning(self, msg: OutboundMessage) -> None:
        """Deliver a complete reasoning block via the delta/end primitives."""
        if not msg.content:
            return
        meta = dict(msg.metadata or {})
        meta.setdefault("_reasoning_delta", True)
        await self.send_reasoning_delta(msg.chat_id, msg.content, meta)
        end_meta = dict(meta)
        end_meta.pop("_reasoning_delta", None)
        end_meta["_reasoning_end"] = True
        await self.send_reasoning_end(msg.chat_id, end_meta)

    @property
    def supports_streaming(self) -> bool:
        """True when config enables streaming AND this subclass implements send_delta."""
        cfg = self.config
        streaming = cfg.get("streaming", False) if isinstance(cfg, dict) else getattr(cfg, "streaming", False)
        return bool(streaming) and type(self).send_delta is not BaseChannel.send_delta

    def is_allowed(self, sender_id: str) -> bool:
        """Check sender permission: star > allowlist > opt-in pairing store > deny."""
        if isinstance(self.config, dict):
            if "allow_from" in self.config:
                allow_list = self.config.get("allow_from")
            else:
                allow_list = self.config.get("allowFrom", [])
        else:
            allow_list = getattr(self.config, "allow_from", [])
        if "*" in allow_list:
            return True
        if str(sender_id) in allow_list:
            return True
        if getattr(self.pairing_config, "enabled", False) and is_approved(self.name, str(sender_id)):
            return True
        if not allow_list:
            self.logger.warning("allow_from is empty — all access denied")
        return False

    async def _handle_message(
        self,
        sender_id: str,
        chat_id: str,
        content: str,
        media: list[str] | None = None,
        metadata: dict[str, Any] | None = None,
        session_key: str | None = None,
        is_dm: bool = False,
    ) -> None:
        """
        Handle an incoming message from the chat platform.

        This method checks permissions and forwards to the bus.

        Args:
            sender_id: The sender's identifier.
            chat_id: The chat/channel identifier.
            content: Message text content.
            media: Optional list of media URLs.
            metadata: Optional channel-specific metadata.
            session_key: Optional session key override (e.g. thread-scoped sessions).
        """
        if not self.is_allowed(sender_id):
            if is_dm and getattr(self.pairing_config, "enabled", False):
                code = generate_code(
                    self.name,
                    str(sender_id),
                    ttl=getattr(self.pairing_config, "ttl_seconds", 600),
                )
                await self.send(
                    OutboundMessage(
                        channel=self.name,
                        chat_id=str(chat_id),
                        content=format_pairing_reply(code),
                        metadata={PAIRING_CODE_META_KEY: code},
                    )
                )
                self.logger.info("Sent pairing code <redacted> to sender {} in chat {}", sender_id, chat_id)
                return
            self.logger.warning(
                "Access denied for sender {}. "
                "Add them to allowFrom list in config to grant access.",
                sender_id,
            )
            return

        meta = metadata or {}
        if self.supports_streaming:
            meta = {**meta, "_wants_stream": True}

        msg = InboundMessage(
            channel=self.name,
            sender_id=str(sender_id),
            chat_id=str(chat_id),
            content=content,
            media=media or [],
            metadata=meta,
            session_key_override=session_key,
        )

        await self.bus.publish_inbound(msg)

    @classmethod
    def default_config(cls) -> dict[str, Any]:
        """Return default config for onboard. Override in plugins to auto-populate config.json."""
        return {"enabled": False}

    @property
    def is_running(self) -> bool:
        """Check if the channel is running."""
        return self._running
