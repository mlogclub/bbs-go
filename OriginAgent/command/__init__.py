"""Slash command routing and built-in handlers."""

from OriginAgent.command.builtin import register_builtin_commands
from OriginAgent.command.router import CommandContext, CommandRouter

__all__ = ["CommandContext", "CommandRouter", "register_builtin_commands"]
