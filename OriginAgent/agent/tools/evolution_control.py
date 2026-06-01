"""Dedicated tool entry point for the governed evolution control plane."""

from __future__ import annotations

from pathlib import Path
from typing import Any

from OriginAgent.agent.evolution_control_plane import EvolutionControlPlane
from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.tools.context import RequestContext


class EvolutionControlTool(Tool):
    """Expose governed evolution control-plane actions through a stable tool."""

    name = "originagent_evolution_control"

    def __init__(self, *, workspace: Path, evolution_config: Any | None = None) -> None:
        self._workspace = Path(workspace)
        self._evolution_config = evolution_config
        self._actor = "tool"
        self._source = self.name

    @property
    def description(self) -> str:
        return (
            "Operate the governed evolution control plane. Supports read-only status, "
            "signal/proposal inspection, structured action discovery, non-mutating "
            "previews, and guarded executions. Write executions still require "
            "learning.evolution.allow_manual_override=true."
        )

    @property
    def parameters(self) -> dict[str, Any]:
        return {
            "type": "object",
            "properties": {
                "operation": {
                    "type": "string",
                    "enum": [
                        "status",
                        "list_actions",
                        "list_recommendations",
                        "list_signals",
                        "list_proposals",
                        "inspect_signal",
                        "inspect_proposal",
                        "list_config_overlay",
                        "preview_action",
                        "execute_action",
                        "generate_report",
                        "explain_health",
                        "validate_schema",
                    ],
                },
                "action_kind": {"type": "string"},
                "target_id": {"type": "string"},
                "reason": {"type": "string"},
                "fixtures": {
                    "type": "object",
                    "additionalProperties": {"type": "string"},
                },
                "force_cleanup": {"type": "boolean"},
                "artifact_type": {"type": "string"},
                "artifact_name": {"type": "string"},
                "snapshot_id": {"type": "string"},
                "period_days": {"type": "integer", "minimum": 1, "maximum": 90},
                "status": {"type": "string"},
                "kind": {"type": "string"},
                "proposal_type": {"type": "string"},
                "limit": {"type": "integer", "minimum": 1, "maximum": 200},
            },
            "required": ["operation"],
            "additionalProperties": False,
        }

    def set_context(self, ctx: RequestContext) -> None:
        self._actor = str(ctx.actor_id or ctx.session_key or ctx.chat_id or "tool")
        self._source = f"{self.name}:{ctx.trigger or ctx.channel or 'unknown'}"

    async def execute(
        self,
        operation: str,
        action_kind: str = "",
        target_id: str = "",
        reason: str = "",
        fixtures: dict[str, str] | None = None,
        force_cleanup: bool = False,
        artifact_type: str = "",
        artifact_name: str = "",
        snapshot_id: str = "",
        period_days: int = 7,
        status: str = "",
        kind: str = "",
        proposal_type: str = "",
        limit: int = 50,
    ) -> dict[str, Any] | str:
        plane = EvolutionControlPlane(self._workspace, self._evolution_config)
        op = str(operation or "").strip()
        if op == "status":
            return plane.status()
        if op == "list_actions":
            return plane.list_actions()
        if op == "list_recommendations":
            return {"ok": True, "recommendations": plane.list_recommendations()}
        if op == "list_signals":
            return plane.list_signals(status=status or None, kind=kind or None, limit=limit)
        if op == "list_proposals":
            return plane.list_proposals(status=status or None, proposal_type=proposal_type or None, limit=limit)
        if op == "inspect_signal":
            return plane.inspect_signal(target_id)
        if op == "inspect_proposal":
            return plane.inspect_proposal(target_id)
        if op == "list_config_overlay":
            return plane.execute_action("list_config_overlay", actor=self._actor, source=self._source)
        if op == "explain_health":
            return plane.explain_health()
        if op == "generate_report":
            return plane.generate_report(period_days=period_days)
        if op == "validate_schema":
            return plane.execute_action("validate_schema", actor=self._actor, source=self._source)
        if op == "preview_action":
            if not action_kind:
                return {"ok": False, "error": "missing_action_kind", "message": "preview_action requires action_kind."}
            return plane.preview_action(
                action_kind,
                target_id=target_id,
                reason=reason,
                fixtures=fixtures or {},
                force_cleanup=force_cleanup,
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                snapshot_id=snapshot_id,
                period_days=period_days,
            )
        if op == "execute_action":
            if not action_kind:
                return {"ok": False, "error": "missing_action_kind", "message": "execute_action requires action_kind."}
            return plane.execute_action(
                action_kind,
                target_id=target_id,
                reason=reason,
                fixtures=fixtures or {},
                force_cleanup=force_cleanup,
                artifact_type=artifact_type,
                artifact_name=artifact_name,
                snapshot_id=snapshot_id,
                period_days=period_days,
                actor=self._actor,
                source=self._source,
            )
        return {"ok": False, "error": "unsupported_operation", "message": f"Unsupported operation `{operation}`."}
