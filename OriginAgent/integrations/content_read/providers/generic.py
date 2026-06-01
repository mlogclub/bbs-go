"""Generic web content provider inspired by feedgrab's Jina fallback."""

from __future__ import annotations

import html
import re
from urllib.parse import quote

import httpx

from OriginAgent.integrations.content_read.types import ContentReadResult
from OriginAgent.security.network import validate_resolved_url


async def fetch_generic(
    url: str,
    *,
    proxy: str | None = None,
    user_agent: str | None = None,
    use_jina_reader: bool = True,
) -> ContentReadResult:
    if use_jina_reader:
        jina = await _fetch_jina(url, proxy=proxy, user_agent=user_agent)
        if jina is not None:
            return jina
    return await _fetch_html(url, proxy=proxy, user_agent=user_agent)


async def _fetch_jina(
    url: str,
    *,
    proxy: str | None,
    user_agent: str | None,
) -> ContentReadResult | None:
    headers = {"Accept": "text/plain", "User-Agent": user_agent or "OriginAgent"}
    endpoint = f"https://r.jina.ai/{quote(url, safe=':/?=&%')}"
    try:
        async with httpx.AsyncClient(proxy=proxy, follow_redirects=True, timeout=20.0) as client:
            response = await client.get(endpoint, headers=headers)
            ok, err = validate_resolved_url(str(response.url))
            if not ok:
                raise ValueError(f"Redirect blocked: {err}")
            response.raise_for_status()
        text = response.text.strip()
        if not text:
            return None
        return ContentReadResult(
            source_type="web",
            title=_title_from_markdown(text) or url,
            url=url,
            content=text,
            metadata={"extractor": "jina"},
        )
    except Exception:
        return None


async def _fetch_html(
    url: str,
    *,
    proxy: str | None,
    user_agent: str | None,
) -> ContentReadResult:
    headers = {"User-Agent": user_agent or "OriginAgent"}
    async with httpx.AsyncClient(proxy=proxy, follow_redirects=True, timeout=15.0) as client:
        response = await client.get(url, headers=headers)
        ok, err = validate_resolved_url(str(response.url))
        if not ok:
            raise ValueError(f"Redirect blocked: {err}")
        response.raise_for_status()
    raw = response.text
    title = _title_from_html(raw) or str(response.url)
    text = _html_to_text(raw)
    return ContentReadResult(
        source_type="web",
        title=title,
        url=str(response.url),
        content=text,
        metadata={"extractor": "html"},
    )


def _title_from_markdown(text: str) -> str:
    for line in text.splitlines():
        stripped = line.strip()
        if stripped.startswith("Title:"):
            return stripped.removeprefix("Title:").strip()
        if stripped.startswith("#"):
            return stripped.lstrip("#").strip()
    return ""


def _title_from_html(raw: str) -> str:
    match = re.search(r"<title[^>]*>(.*?)</title>", raw, flags=re.I | re.S)
    return html.unescape(match.group(1)).strip() if match else ""


def _html_to_text(raw: str) -> str:
    text = re.sub(r"<(script|style|nav|footer|header)[\s\S]*?</\1>", "", raw, flags=re.I)
    text = re.sub(r"<br\s*/?>", "\n", text, flags=re.I)
    text = re.sub(r"</p>", "\n\n", text, flags=re.I)
    text = re.sub(r"<[^>]+>", "", text)
    text = html.unescape(text)
    text = re.sub(r"[ \t]+", " ", text)
    return re.sub(r"\n{3,}", "\n\n", text).strip()
