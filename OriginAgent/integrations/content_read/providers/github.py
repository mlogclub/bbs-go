"""GitHub provider adapted from feedgrab's README-first fetcher."""

from __future__ import annotations

import base64
import re
from urllib.parse import urlparse

import httpx

from OriginAgent.integrations.content_read.types import ContentReadResult
from OriginAgent.security.network import validate_resolved_url

API_BASE = "https://api.github.com"
CHINESE_README_VARIANTS = [
    "README_CN.md",
    "README.zh-CN.md",
    "README_zh-CN.md",
    "README.zh.md",
    "README_ZH.md",
    "README.zh-Hans.md",
    "README_zh.md",
    "README.Chinese.md",
]


def is_github_url(url: str) -> bool:
    return (urlparse(url).hostname or "").lower() == "github.com"


async def fetch_github(
    url: str,
    *,
    proxy: str | None = None,
    user_agent: str | None = None,
) -> ContentReadResult:
    owner, repo, file_path = parse_github_url(url)
    headers = {
        "Accept": "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
        "User-Agent": user_agent or "OriginAgent",
    }
    async with httpx.AsyncClient(proxy=proxy, follow_redirects=True, timeout=20.0) as client:
        meta = await _get_json(client, f"{API_BASE}/repos/{owner}/{repo}", headers)
        readme_path = file_path or await _find_readme_path(client, owner, repo, headers)
        content, filename = await _fetch_readme(client, owner, repo, readme_path, headers)

    full_name = meta.get("full_name") or f"{owner}/{repo}"
    body = [
        f"# {full_name}",
        "",
        meta.get("description") or "",
        "",
        f"Repository: {meta.get('html_url') or url}",
        f"Default branch: {meta.get('default_branch') or ''}",
        "",
        f"## {filename}",
        "",
        content,
    ]
    return ContentReadResult(
        source_type="github",
        title=full_name,
        url=meta.get("html_url") or url,
        content="\n".join(part for part in body if part is not None).strip(),
        metadata={
            "description": meta.get("description") or "",
            "stars": meta.get("stargazers_count", 0),
            "forks": meta.get("forks_count", 0),
            "language": meta.get("language") or "",
            "license": (meta.get("license") or {}).get("spdx_id") or "",
            "topics": meta.get("topics") or [],
            "readme": filename,
            "default_branch": meta.get("default_branch") or "",
        },
    )


def parse_github_url(url: str) -> tuple[str, str, str | None]:
    parts = [p for p in urlparse(url).path.strip("/").split("/") if p]
    if len(parts) < 2:
        raise ValueError("GitHub URL must include owner/repo")
    owner = parts[0]
    repo = parts[1][:-4] if parts[1].endswith(".git") else parts[1]
    file_path = None
    if len(parts) >= 5 and parts[2] == "blob":
        file_path = "/".join(parts[4:])
    return owner, repo, file_path


async def _get_json(client: httpx.AsyncClient, url: str, headers: dict[str, str]) -> dict:
    response = await client.get(url, headers=headers)
    ok, err = validate_resolved_url(str(response.url))
    if not ok:
        raise ValueError(f"Redirect blocked: {err}")
    response.raise_for_status()
    return response.json()


async def _find_readme_path(
    client: httpx.AsyncClient,
    owner: str,
    repo: str,
    headers: dict[str, str],
) -> str | None:
    try:
        files = await _get_json(client, f"{API_BASE}/repos/{owner}/{repo}/contents/", headers)
    except Exception:
        return None
    if not isinstance(files, list):
        return None
    names = {item.get("name", "") for item in files if item.get("type") == "file"}
    names_lower = {name.lower(): name for name in names}
    for variant in CHINESE_README_VARIANTS:
        if variant in names:
            return variant
        actual = names_lower.get(variant.lower())
        if actual:
            return actual
    return None


async def _fetch_readme(
    client: httpx.AsyncClient,
    owner: str,
    repo: str,
    path: str | None,
    headers: dict[str, str],
) -> tuple[str, str]:
    endpoint = (
        f"{API_BASE}/repos/{owner}/{repo}/contents/{path}"
        if path
        else f"{API_BASE}/repos/{owner}/{repo}/readme"
    )
    data = await _get_json(client, endpoint, headers)
    filename = data.get("name") or path or "README.md"
    encoded = re.sub(r"\s+", "", data.get("content") or "")
    if encoded:
        return base64.b64decode(encoded).decode("utf-8", errors="replace"), filename
    download_url = data.get("download_url")
    if not download_url:
        return "", filename
    response = await client.get(download_url, headers={"User-Agent": headers["User-Agent"]})
    ok, err = validate_resolved_url(str(response.url))
    if not ok:
        raise ValueError(f"Redirect blocked: {err}")
    response.raise_for_status()
    return response.text, filename
