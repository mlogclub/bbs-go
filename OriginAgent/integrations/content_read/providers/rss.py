"""RSS/Atom provider adapted from feedgrab's low-risk RSS fetcher."""

from __future__ import annotations

import html
import re
import xml.etree.ElementTree as ET
from typing import Any

import httpx

from OriginAgent.integrations.content_read.types import ContentReadResult
from OriginAgent.security.network import validate_resolved_url


def is_rss_like_url(url: str) -> bool:
    lowered = url.lower()
    return lowered.endswith((".xml", ".rss", ".atom")) or any(
        marker in lowered for marker in ("/rss", "/feed", "/atom")
    )


async def fetch_rss(
    url: str,
    *,
    proxy: str | None = None,
    user_agent: str | None = None,
    limit: int = 10,
) -> ContentReadResult:
    headers = {"User-Agent": user_agent or "OriginAgent"}
    async with httpx.AsyncClient(proxy=proxy, follow_redirects=True, timeout=15.0) as client:
        response = await client.get(url, headers=headers)
        ok, err = validate_resolved_url(str(response.url))
        if not ok:
            raise ValueError(f"Redirect blocked: {err}")
        response.raise_for_status()

    root = ET.fromstring(response.content)
    feed_title = _first_text(root, ["./channel/title", "./{http://www.w3.org/2005/Atom}title"]) or url
    entries = _rss_items(root) or _atom_entries(root)
    lines: list[str] = [f"# {feed_title}"]
    payload_entries: list[dict[str, Any]] = []

    for entry in entries[: max(1, limit)]:
        title = _entry_text(entry, ["title", "{http://www.w3.org/2005/Atom}title"])
        link = _entry_link(entry)
        published = _entry_text(
            entry,
            ["published", "pubDate", "{http://www.w3.org/2005/Atom}published", "{http://www.w3.org/2005/Atom}updated"],
        )
        summary = _clean_html(
            _entry_text(
                entry,
                ["description", "summary", "{http://www.w3.org/2005/Atom}summary", "{http://www.w3.org/2005/Atom}content"],
            )
        )
        payload_entries.append(
            {"title": title, "url": link, "published": published, "summary": summary}
        )
        lines.append(f"\n## {title or link or 'Untitled'}")
        if link:
            lines.append(link)
        if published:
            lines.append(f"Published: {published}")
        if summary:
            lines.append(summary)

    return ContentReadResult(
        source_type="rss",
        title=feed_title,
        url=str(response.url),
        content="\n".join(lines).strip(),
        metadata={"entries": payload_entries, "entry_count": len(payload_entries)},
    )


def _rss_items(root: ET.Element) -> list[ET.Element]:
    return list(root.findall("./channel/item"))


def _atom_entries(root: ET.Element) -> list[ET.Element]:
    return list(root.findall("./{http://www.w3.org/2005/Atom}entry"))


def _first_text(root: ET.Element, paths: list[str]) -> str:
    for path in paths:
        found = root.find(path)
        if found is not None and found.text:
            return found.text.strip()
    return ""


def _entry_text(entry: ET.Element, names: list[str]) -> str:
    for name in names:
        found = entry.find(name)
        if found is not None:
            text = "".join(found.itertext()).strip()
            if text:
                return text
    return ""


def _entry_link(entry: ET.Element) -> str:
    link = entry.find("link")
    if link is not None and link.text:
        return link.text.strip()
    atom_link = entry.find("{http://www.w3.org/2005/Atom}link")
    if atom_link is not None:
        return (atom_link.attrib.get("href") or "").strip()
    return ""


def _clean_html(value: str) -> str:
    text = re.sub(r"<(script|style)[\s\S]*?</\1>", "", value, flags=re.I)
    text = re.sub(r"<[^>]+>", "", text)
    text = html.unescape(text)
    return re.sub(r"\s+", " ", text).strip()
