"""Hacker News provider adapted from feedgrab's Firebase API fetcher."""

from __future__ import annotations

import html
import re
from datetime import datetime, timezone
from urllib.parse import parse_qs, urlparse

import httpx

from OriginAgent.integrations.content_read.types import ContentReadResult
from OriginAgent.security.network import validate_resolved_url

API_BASE = "https://hacker-news.firebaseio.com/v0"


def is_hackernews_url(url: str) -> bool:
    hostname = (urlparse(url).hostname or "").lower()
    return hostname == "news.ycombinator.com"


async def fetch_hackernews(
    url: str,
    *,
    proxy: str | None = None,
    user_agent: str | None = None,
    max_comments: int = 20,
) -> ContentReadResult:
    item_id = _parse_item_id(url)
    headers = {"User-Agent": user_agent or "OriginAgent"}
    async with httpx.AsyncClient(proxy=proxy, follow_redirects=True, timeout=15.0) as client:
        item = await _api_get(client, f"item/{item_id}", headers)
        if not item or item.get("dead") or item.get("deleted"):
            raise ValueError(f"Hacker News item {item_id} is missing or deleted")
        kids = [str(k) for k in item.get("kids", [])[: max(0, max_comments)]]
        comments = [
            comment
            for comment in await _api_get_many(client, [f"item/{kid}" for kid in kids], headers)
            if comment and not comment.get("dead") and not comment.get("deleted")
        ]

    title = item.get("title") or f"Hacker News item {item_id}"
    target_url = item.get("url") or f"https://news.ycombinator.com/item?id={item_id}"
    lines = [
        f"# {title}",
        "",
        f"HN: https://news.ycombinator.com/item?id={item_id}",
        f"URL: {target_url}",
        f"By: {item.get('by', '')}",
        f"Score: {item.get('score', 0)}",
        f"Time: {_format_time(item.get('time'))}",
    ]
    text = _html_to_text(item.get("text") or "")
    if text:
        lines.extend(["", text])
    if comments:
        lines.extend(["", "## Top comments"])
        for comment in comments:
            rendered = _render_comment(comment)
            if rendered:
                lines.extend(["", rendered])

    return ContentReadResult(
        source_type="hackernews",
        title=title,
        url=f"https://news.ycombinator.com/item?id={item_id}",
        content="\n".join(lines).strip(),
        metadata={
            "id": item_id,
            "by": item.get("by") or "",
            "score": item.get("score", 0),
            "type": item.get("type") or "",
            "comment_count": len(comments),
            "target_url": target_url,
        },
    )


def _parse_item_id(url: str) -> str:
    parsed = urlparse(url)
    qs = parse_qs(parsed.query)
    if parsed.path.rstrip("/") == "/item" and qs.get("id") and qs["id"][0].isdigit():
        return qs["id"][0]
    raise ValueError("Only Hacker News item URLs are supported in content_read MVP")


async def _api_get(
    client: httpx.AsyncClient,
    endpoint: str,
    headers: dict[str, str],
) -> dict | None:
    response = await client.get(f"{API_BASE}/{endpoint}.json", headers=headers)
    ok, err = validate_resolved_url(str(response.url))
    if not ok:
        raise ValueError(f"Redirect blocked: {err}")
    if response.status_code == 404:
        return None
    response.raise_for_status()
    return response.json()


async def _api_get_many(
    client: httpx.AsyncClient,
    endpoints: list[str],
    headers: dict[str, str],
) -> list[dict | None]:
    results: list[dict | None] = []
    for endpoint in endpoints:
        results.append(await _api_get(client, endpoint, headers))
    return results


def _render_comment(comment: dict) -> str:
    by = comment.get("by") or "[deleted]"
    body = _html_to_text(comment.get("text") or "")
    if not body:
        return ""
    return f"**@{by} · {_format_time(comment.get('time'))}**\n\n{body}".strip()


def _format_time(value: int | None) -> str:
    if not value:
        return ""
    return datetime.fromtimestamp(int(value), tz=timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _html_to_text(value: str) -> str:
    text = re.sub(r"<p>", "\n\n", value, flags=re.I)
    text = re.sub(r'<a\s+href="([^"]+)"[^>]*>(.*?)</a>', r"\2 (\1)", text, flags=re.I | re.S)
    text = re.sub(r"<[^>]+>", "", text)
    text = html.unescape(text)
    return re.sub(r"\n{3,}", "\n\n", text).strip()
