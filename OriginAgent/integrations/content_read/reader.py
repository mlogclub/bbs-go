"""Provider dispatcher for the content_read tool."""

from __future__ import annotations

from dataclasses import dataclass, field
from urllib.parse import urlparse

from OriginAgent.integrations.content_read.providers.generic import fetch_generic
from OriginAgent.integrations.content_read.providers.github import fetch_github, is_github_url
from OriginAgent.integrations.content_read.providers.hackernews import (
    fetch_hackernews,
    is_hackernews_url,
)
from OriginAgent.integrations.content_read.providers.rss import fetch_rss, is_rss_like_url
from OriginAgent.integrations.content_read.types import ContentReadResult

CONTENT_READ_PROVIDERS = ("auto", "generic", "rss", "github", "hackernews")


class ContentReadError(RuntimeError):
    """Raised for expected content_read failures."""


@dataclass
class ContentReader:
    """Routes a URL to a low-side-effect content provider."""

    proxy: str | None = None
    user_agent: str | None = None
    enabled_providers: set[str] = field(
        default_factory=lambda: {"generic", "rss", "github", "hackernews"}
    )
    use_jina_reader: bool = True
    rss_entry_limit: int = 10
    hackernews_comment_limit: int = 20

    def detect_provider(self, url: str) -> str:
        parsed = urlparse(url)
        if parsed.scheme not in {"http", "https"} or not parsed.netloc:
            raise ContentReadError("content_read requires an http(s) URL with a hostname")
        if is_github_url(url):
            return "github"
        if is_hackernews_url(url):
            return "hackernews"
        if is_rss_like_url(url):
            return "rss"
        return "generic"

    async def read(self, url: str, provider: str = "auto") -> ContentReadResult:
        normalized = provider.strip().lower() or "auto"
        if normalized not in CONTENT_READ_PROVIDERS:
            raise ContentReadError(
                f"Unsupported content_read provider '{provider}'. "
                f"Supported providers: {', '.join(CONTENT_READ_PROVIDERS)}"
            )
        resolved = self.detect_provider(url) if normalized == "auto" else normalized
        if resolved not in self.enabled_providers:
            raise ContentReadError(f"content_read provider '{resolved}' is disabled")
        if resolved == "github":
            return await fetch_github(url, proxy=self.proxy, user_agent=self.user_agent)
        if resolved == "hackernews":
            return await fetch_hackernews(
                url,
                proxy=self.proxy,
                user_agent=self.user_agent,
                max_comments=self.hackernews_comment_limit,
            )
        if resolved == "rss":
            return await fetch_rss(
                url,
                proxy=self.proxy,
                user_agent=self.user_agent,
                limit=self.rss_entry_limit,
            )
        return await fetch_generic(
            url,
            proxy=self.proxy,
            user_agent=self.user_agent,
            use_jina_reader=self.use_jina_reader,
        )
