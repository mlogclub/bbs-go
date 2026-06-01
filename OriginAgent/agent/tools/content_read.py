"""Platform-aware content reading tool."""

from __future__ import annotations

import json
from typing import Any

from loguru import logger

from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.tools.limits import ToolLimits
from OriginAgent.agent.tools.schema import IntegerSchema, StringSchema, tool_parameters_schema
from OriginAgent.integrations.content_read.reader import (
    CONTENT_READ_PROVIDERS,
)
from OriginAgent.security.policy import PolicyDeniedError


class ContentReadTool(Tool):
    """Read a URL through a platform-specific provider and return structured content."""

    name = "content_read"
    description = (
        "Deprecated compatibility wrapper for structured URL reads. "
        "Use web_fetch with mode=structured instead."
    )

    def __init__(
        self,
        *,
        config: Any | None = None,
        proxy: str | None = None,
        user_agent: str | None = None,
        limits: ToolLimits | None = None,
    ) -> None:
        self.config = config
        self.proxy = proxy
        self.user_agent = user_agent
        self._limits = limits or ToolLimits()
        self.max_chars = (
            getattr(config, "max_chars", None) or self._limits.web_fetch_max_chars
        )

    @property
    def parameters(self) -> dict[str, Any]:
        return tool_parameters_schema(
            url=StringSchema("URL to read"),
            provider={
                "type": "string",
                "description": "Provider to use; auto detects from URL",
                "enum": list(CONTENT_READ_PROVIDERS),
                "default": "auto",
            },
            max_chars=IntegerSchema(
                self.max_chars,
                description=f"Maximum content characters to return (default {self.max_chars:,})",
                minimum=100,
            ),
            required=["url"],
        )

    @property
    def read_only(self) -> bool:
        return True

    async def execute(
        self,
        url: str,
        provider: str = "auto",
        max_chars: int | None = None,
        **kwargs: Any,
    ) -> str:
        if "maxChars" in kwargs and max_chars is None:
            max_chars = kwargs.pop("maxChars")
        cleaned_url = url.strip(" \t\r\n`\"'")
        try:
            from OriginAgent.agent.tools.web import WebFetchTool
            from OriginAgent.config.schema import WebFetchConfig

            logger.warning(
                "content_read is deprecated; use web_fetch(mode='structured') instead"
            )
            web_fetch = WebFetchTool(
                config=WebFetchConfig(
                    use_jina_reader=bool(getattr(self.config, "use_jina_reader", True))
                ),
                proxy=self.proxy,
                user_agent=self.user_agent,
                max_chars=self.max_chars,
                limits=self._limits,
                content_read_config=self.config,
            )
            return await web_fetch.execute(
                cleaned_url,
                mode="structured",
                provider=provider,
                max_chars=max_chars or self.max_chars,
            )
        except PolicyDeniedError:
            raise
        except Exception as exc:
            return json.dumps(
                {"error": f"content_read failed: {type(exc).__name__}: {exc}", "url": cleaned_url},
                ensure_ascii=False,
            )
