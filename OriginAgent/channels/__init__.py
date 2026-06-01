"""Chat channels module with plugin architecture."""

from OriginAgent.channels.base import BaseChannel
from OriginAgent.channels.manager import ChannelManager

__all__ = ["BaseChannel", "ChannelManager"]
