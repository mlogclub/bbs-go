"""Message bus module for decoupled channel-agent communication."""

from OriginAgent.bus.events import InboundMessage, OutboundMessage
from OriginAgent.bus.queue import MessageBus

__all__ = ["MessageBus", "InboundMessage", "OutboundMessage"]
