"""Delegated runtime policy for subagent execution."""

from __future__ import annotations

from dataclasses import dataclass

from OriginAgent.security.capabilities import CapabilitySnapshot

_DEFAULT_ALLOWED_TOOL_NAMES = frozenset(
    {
        "read_file",
        "list_dir",
        "glob",
        "grep",
    }
)


@dataclass(frozen=True)
class SubagentPolicy:
    """Explicit delegated runtime policy for phase-1 subagents."""

    capability_snapshot: CapabilitySnapshot
    allowed_tool_names: frozenset[str]
    allow_web: bool = False
    allow_nested_spawn: bool = False
    max_subagent_depth: int = 1
    max_children_per_subagent: int = 0
    child_allowed_tool_names: frozenset[str] | None = None

    @classmethod
    def default_for_parent(
        cls,
        parent_snapshot: CapabilitySnapshot | None,
    ) -> "SubagentPolicy":
        parent = parent_snapshot or CapabilitySnapshot.system_default()
        snapshot = parent.derive_subagent()
        allowed = _DEFAULT_ALLOWED_TOOL_NAMES if snapshot.can_read_files else frozenset()
        return cls(
            capability_snapshot=snapshot,
            allowed_tool_names=allowed,
            allow_web=False,
            allow_nested_spawn=False,
            max_subagent_depth=1,
            max_children_per_subagent=0,
            child_allowed_tool_names=None,
        )

    def allows(self, tool_name: str) -> bool:
        return tool_name in self.allowed_tool_names

    def child_policy(self) -> "SubagentPolicy":
        """Return the downgraded policy for one nested delegated level."""

        child_snapshot = self.capability_snapshot.derive_subagent()
        child_tools = self.child_allowed_tool_names
        if child_tools is None:
            child_tools = _DEFAULT_ALLOWED_TOOL_NAMES if child_snapshot.can_read_files else frozenset()
        return SubagentPolicy(
            capability_snapshot=child_snapshot,
            allowed_tool_names=child_tools,
            allow_web=False,
            allow_nested_spawn=self.allow_nested_spawn and self.max_subagent_depth > 1,
            max_subagent_depth=max(0, self.max_subagent_depth - 1),
            max_children_per_subagent=self.max_children_per_subagent,
            child_allowed_tool_names=self.child_allowed_tool_names,
        )
