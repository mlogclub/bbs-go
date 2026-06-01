"""Tool wrapper for searching persisted OriginAgent conversation history."""

from __future__ import annotations

from pathlib import Path
from typing import Any

from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.tools.schema import ArraySchema, IntegerSchema, StringSchema, tool_parameters_schema
from OriginAgent.session.search import SessionSearchService


class SessionSearchTool(Tool):
    """Read-only structured search over session/history JSONL sources."""

    def __init__(self, workspace: Path, *, config: Any | None = None, index_service: Any | None = None):
        self._workspace = Path(workspace)
        self._config = config
        refresh_ms = getattr(config, "max_tool_refresh_ms", 500)
        if refresh_ms is None:
            refresh_ms = 500
        self._service = SessionSearchService(
            self._workspace,
            index_service=index_service,
            index_backend=str(getattr(config, "backend", "auto")),
            semantic_enabled=bool(getattr(config, "semantic_enabled", True)),
            max_tool_refresh_ms=int(refresh_ms),
        )

    @property
    def name(self) -> str:
        return "session_search"

    @property
    def description(self) -> str:
        return (
            "Search previous OriginAgent conversations, memory/history archives, cold archives, and WebUI "
            "transcripts. Defaults to case-insensitive literal text matching; pass "
            "mode='hybrid' or mode='semantic' for indexed multilingual recall when "
            "enabled. Use this for "
            "questions about prior discussions, previous plans, or what was said before; "
            "use grep for arbitrary project files. Prefer supplying since/until to reduce "
            "history scanning."
        )

    @property
    def read_only(self) -> bool:
        return True

    @property
    def parameters(self) -> dict[str, Any]:
        return tool_parameters_schema(
            required=["query"],
            additional_properties=False,
            query=StringSchema(
                "Literal text to search for in historical conversation records.",
                min_length=1,
                max_length=500,
            ),
            roles=ArraySchema(
                StringSchema(
                    "Role filter.",
                    enum=["user", "assistant", "tool", "system", "archive"],
                ),
                description="Optional roles to include.",
                max_items=8,
            ),
            sources=ArraySchema(
                StringSchema(
                    "History source.",
                    enum=["sessions", "history", "webui", "facts", "cold"],
                ),
                description="Optional sources to search. Defaults to sessions/history/webui.",
                max_items=5,
            ),
            mode=StringSchema(
                "Search mode. literal preserves exact legacy behavior; hybrid combines literal and indexed multilingual recall; semantic uses indexed multilingual recall.",
                enum=["literal", "hybrid", "semantic"],
            ),
            session_key=StringSchema(
                "Optional exact session key such as 'websocket:chat1'.",
                max_length=200,
            ),
            channel=StringSchema(
                "Optional channel filter, used with chat_id when known.",
                max_length=80,
            ),
            chat_id=StringSchema(
                "Optional chat id filter, used with channel when known.",
                max_length=160,
            ),
            since=StringSchema(
                "Optional inclusive start time, YYYY-MM-DD or ISO datetime.",
                max_length=80,
            ),
            until=StringSchema(
                "Optional inclusive end time, YYYY-MM-DD or ISO datetime.",
                max_length=80,
            ),
            limit=IntegerSchema(
                description="Maximum results to return. Defaults to 10; values above 50 are clamped.",
                minimum=1,
            ),
        )

    async def execute(
        self,
        *,
        query: str,
        roles: list[str] | None = None,
        sources: list[str] | None = None,
        session_key: str | None = None,
        channel: str | None = None,
        chat_id: str | None = None,
        since: str | None = None,
        until: str | None = None,
        limit: int | None = None,
        mode: str = "literal",
    ) -> dict[str, Any]:
        return self._service.search(
            query=query,
            roles=roles,
            sources=sources,
            session_key=session_key,
            channel=channel,
            chat_id=chat_id,
            since=since,
            until=until,
            limit=limit,
            mode=mode,
        )
