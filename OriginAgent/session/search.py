"""Literal search over OriginAgent session and history JSONL sources."""

from __future__ import annotations

import json
import re
import time
from dataclasses import dataclass
from datetime import datetime, time as datetime_time, timezone
from functools import lru_cache
from pathlib import Path
from typing import Any, Iterable

from loguru import logger

from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.config.loader import get_config_path
from OriginAgent.session.cold_archive import SESSION_COLD_ARCHIVE_DIR
from OriginAgent.utils.helpers import truncate_text

DEFAULT_SOURCES: tuple[str, ...] = ("sessions", "history", "webui")
SUPPORTED_SOURCES: tuple[str, ...] = ("sessions", "history", "webui", "facts", "cold")
SOURCE_PRIORITY: dict[str, int] = {"sessions": 0, "history": 1, "webui": 2, "facts": 3, "cold": 4}
SUPPORTED_MODES: tuple[str, ...] = ("literal", "hybrid", "semantic")
DEFAULT_LIMIT = 10
MAX_LIMIT = 50
DEFAULT_CACHE_TTL_S = 300.0
DEFAULT_CACHE_RECORDS_PER_SOURCE = 1000
RECENT_SCAN_WINDOW_DAYS = 30
LONG_HISTORY_TEXT_CHARS = 500
SNIPPET_MAX_CHARS = 240
SCAN_PERFORMANCE_NOTE_THRESHOLD = 50_000
_DATE_ONLY_RE = re.compile(r"^\d{4}-\d{2}-\d{2}$")


@dataclass(frozen=True)
class SearchRecord:
    source: str
    session_key: str
    role: str
    timestamp: datetime | None
    text: str
    locator: dict[str, Any]
    record_status: str = ""


@dataclass(frozen=True)
class SearchResponse:
    query: str
    mode: str
    results: list[dict[str, Any]]
    total_matches: int
    searched_sources: list[str]
    skipped_records: int
    truncated: bool
    performance_note: str | None = None
    index_stale: bool = False
    index_refresh_running: bool = False

    def to_dict(self) -> dict[str, Any]:
        return {
            "query": self.query,
            "mode": self.mode,
            "results": self.results,
            "total_matches": self.total_matches,
            "searched_sources": self.searched_sources,
            "skipped_records": self.skipped_records,
            "truncated": self.truncated,
            "performance_note": self.performance_note,
            "index_stale": self.index_stale,
            "index_refresh_running": self.index_refresh_running,
        }


@dataclass(frozen=True)
class _SourceLoad:
    records: list[SearchRecord]
    skipped_records: int
    scanned_records: int


@dataclass
class _SourceCache:
    records: list[SearchRecord]
    skipped_records: int
    scanned_records: int
    fingerprint: tuple[tuple[str, int, float], ...]
    loaded_at: float


class SessionSearchService:
    """Search workspace conversation history without mutating any state."""

    def __init__(
        self,
        workspace: Path,
        *,
        cache_ttl_s: float = DEFAULT_CACHE_TTL_S,
        cache_records_per_source: int = DEFAULT_CACHE_RECORDS_PER_SOURCE,
        webui_dir: Path | None = None,
        now: Any | None = None,
        index_service: Any | None = None,
        index_backend: str = "auto",
        semantic_enabled: bool = True,
        max_tool_refresh_ms: int = 500,
    ) -> None:
        self.workspace = Path(workspace)
        self.sessions_dir = self.workspace / "sessions"
        self.history_file = self.workspace / "memory" / "history.jsonl"
        self.cold_archive_dir = self.workspace / SESSION_COLD_ARCHIVE_DIR
        self._webui_dir = webui_dir
        self._cache_ttl_s = cache_ttl_s
        self._cache_records_per_source = cache_records_per_source
        self._now = now or time.monotonic
        self._cache: dict[str, _SourceCache] = {}
        self._index_service = index_service
        self._index_backend = index_backend
        self._semantic_enabled = semantic_enabled
        self._max_tool_refresh_ms = max_tool_refresh_ms

    def search(
        self,
        *,
        query: str,
        roles: Iterable[str] | None = None,
        sources: Iterable[str] | None = None,
        session_key: str | None = None,
        channel: str | None = None,
        chat_id: str | None = None,
        since: str | datetime | None = None,
        until: str | datetime | None = None,
        limit: int | None = DEFAULT_LIMIT,
        mode: str = "literal",
    ) -> dict[str, Any]:
        query = str(query or "").strip()
        mode = _normalize_mode(mode)
        if not query:
            return SearchResponse(
                query=query,
                mode=mode,
                results=[],
                total_matches=0,
                searched_sources=[],
                skipped_records=0,
                truncated=False,
                performance_note="Provide a non-empty literal query.",
            ).to_dict()

        requested_sources = _normalize_sources(sources)
        requested_roles = _normalize_str_set(roles)
        limit_value = _clamp_limit(limit)
        if mode != "literal":
            return self._indexed_search(
                query=query,
                mode=mode,
                requested_sources=requested_sources,
                requested_roles=requested_roles,
                session_key=session_key,
                channel=channel,
                chat_id=chat_id,
                since=since,
                until=until,
                limit_value=limit_value,
            )
        since_dt = _parse_datetime_filter(since, is_until=False)
        until_dt = _parse_datetime_filter(until, is_until=True)
        target_session_key = session_key or _session_key_from_channel_chat(channel, chat_id)
        old_range_requested = _is_old_range(since_dt) or _is_old_range(until_dt)
        use_live_scan = old_range_requested or limit_value > self._cache_records_per_source

        source_loads: dict[str, _SourceLoad] = {}
        for source in requested_sources:
            source_loads[source] = self._load_source(source, live=use_live_scan)

        query_lc = query.casefold()
        matches: list[tuple[SearchRecord, int]] = []
        skipped_records = sum(load.skipped_records for load in source_loads.values())
        scanned_records = sum(load.scanned_records for load in source_loads.values())
        for source in requested_sources:
            for record in source_loads[source].records:
                if requested_roles and record.role not in requested_roles:
                    continue
                if target_session_key and record.session_key != target_session_key:
                    continue
                if channel and _channel_from_session_key(record.session_key) != channel:
                    continue
                if chat_id and _chat_id_from_session_key(record.session_key) != chat_id:
                    continue
                if not _in_time_range(record.timestamp, since_dt, until_dt):
                    continue
                hit_count = _count_literal_matches(record.text, query_lc)
                if hit_count <= 0:
                    continue
                matches.append((record, hit_count))

        matches.sort(key=_sort_key)
        results = [
            _format_result(record, query_lc, hit_count)
            for record, hit_count in matches[:limit_value]
        ]
        truncated = len(matches) > len(results)
        performance_note = _performance_note(
            old_range_requested=old_range_requested,
            scanned_records=scanned_records,
            cache_used=not use_live_scan,
            total_matches=len(matches),
        )
        if not results and performance_note is None:
            performance_note = (
                "No literal matches found. Try alternate keywords, synonyms, "
                "or a wider since/until range."
            )

        return SearchResponse(
            query=query,
            mode=mode,
            results=results,
            total_matches=len(matches),
            searched_sources=list(requested_sources),
            skipped_records=skipped_records,
            truncated=truncated,
            performance_note=performance_note,
        ).to_dict()

    def _indexed_search(
        self,
        *,
        query: str,
        mode: str,
        requested_sources: tuple[str, ...],
        requested_roles: set[str],
        session_key: str | None,
        channel: str | None,
        chat_id: str | None,
        since: str | datetime | None,
        until: str | datetime | None,
        limit_value: int,
    ) -> dict[str, Any]:
        service = self._index()
        notes: list[str] = []
        if not self._semantic_enabled or service is None or not getattr(service, "enabled", False):
            notes.append("Indexed multilingual search is disabled; falling back to literal search.")
            literal = self.search(
                query=query,
                roles=requested_roles,
                sources=requested_sources,
                session_key=session_key,
                channel=channel,
                chat_id=chat_id,
                since=since,
                until=until,
                limit=limit_value,
                mode="literal",
            )
            literal["mode"] = mode
            literal["performance_note"] = _join_notes(notes, literal.get("performance_note"))
            return literal

        refresh_status = service.refresh_incremental(
            sources=requested_sources,
            budget_ms=self._max_tool_refresh_ms,
        )
        indexed = service.search_indexed(
            query=query,
            sources=requested_sources,
            roles=requested_roles,
            session_key=session_key,
            channel=channel,
            chat_id=chat_id,
            since=since,
            until=until,
            limit=limit_value,
            match_type="semantic" if mode == "semantic" else "fts",
        )
        if indexed.get("note"):
            notes.append(str(indexed["note"]))

        literal: dict[str, Any] | None = None
        if mode == "hybrid" or not indexed.get("results"):
            literal = self.search(
                query=query,
                roles=requested_roles,
                sources=requested_sources,
                session_key=session_key,
                channel=channel,
                chat_id=chat_id,
                since=since,
                until=until,
                limit=limit_value,
                mode="literal",
            )

        results = _merge_indexed_results(
            literal.get("results", []) if literal else [],
            indexed.get("results", []),
            limit_value,
        )
        total_matches = len(results)
        if literal is not None:
            total_matches = max(total_matches, int(literal.get("total_matches", 0) or 0))
        total_matches = max(total_matches, int(indexed.get("total_matches", 0) or 0))
        performance_note = _join_notes(
            notes,
            literal.get("performance_note") if literal else None,
        )
        return SearchResponse(
            query=query,
            mode=mode,
            results=results,
            total_matches=total_matches,
            searched_sources=list(requested_sources),
            skipped_records=int(literal.get("skipped_records", 0) if literal else 0),
            truncated=total_matches > len(results),
            performance_note=performance_note,
            index_stale=bool(refresh_status.get("session_search_index_stale")),
            index_refresh_running=bool(refresh_status.get("session_search_refresh_running")),
        ).to_dict()

    def _index(self):
        if self._index_service is not None:
            return self._index_service
        try:
            from OriginAgent.session.search_index import SessionSearchIndexService
        except Exception:
            return None
        self._index_service = SessionSearchIndexService(
            self.workspace,
            webui_dir=self._webui_dir,
            backend=self._index_backend,
            semantic_enabled=self._semantic_enabled,
        )
        return self._index_service

    def _load_source(self, source: str, *, live: bool) -> _SourceLoad:
        paths = self._source_paths(source)
        fingerprint = _fingerprint_paths(paths)
        if not live:
            cached = self._cache.get(source)
            if cached and cached.fingerprint == fingerprint:
                age = float(self._now()) - cached.loaded_at
                if age <= self._cache_ttl_s:
                    return _SourceLoad(
                        records=list(cached.records),
                        skipped_records=cached.skipped_records,
                        scanned_records=cached.scanned_records,
                    )

        loaded = self._scan_source(source, paths)
        cache_records = loaded.records[: self._cache_records_per_source]
        self._cache[source] = _SourceCache(
            records=cache_records,
            skipped_records=loaded.skipped_records,
            scanned_records=loaded.scanned_records,
            fingerprint=fingerprint,
            loaded_at=float(self._now()),
        )
        if live:
            return loaded
        return _SourceLoad(
            records=list(cache_records),
            skipped_records=loaded.skipped_records,
            scanned_records=loaded.scanned_records,
        )

    def _source_paths(self, source: str) -> list[Path]:
        if source == "sessions":
            return sorted(self.sessions_dir.glob("*.jsonl")) if self.sessions_dir.is_dir() else []
        if source == "history":
            return [self.history_file] if self.history_file.is_file() else []
        if source == "webui":
            webui_dir = self._webui_dir or (get_config_path().parent / "webui")
            return sorted(webui_dir.glob("*.jsonl")) if webui_dir.is_dir() else []
        if source == "facts":
            facts_file = self.workspace / "memory" / "facts.jsonl"
            return [facts_file] if facts_file.is_file() else []
        if source == "cold":
            if not self.cold_archive_dir.is_dir():
                return []
            return sorted(path for path in self.cold_archive_dir.glob("*.jsonl") if path.is_file())
        return []

    def _scan_source(self, source: str, paths: list[Path]) -> _SourceLoad:
        records: list[SearchRecord] = []
        skipped_records = 0
        scanned_records = 0
        for path in paths:
            parsed = self._scan_file(source, path)
            records.extend(parsed.records)
            skipped_records += parsed.skipped_records
            scanned_records += parsed.scanned_records
        records.sort(key=_recency_sort_key)
        return _SourceLoad(records=records, skipped_records=skipped_records, scanned_records=scanned_records)

    def _scan_file(self, source: str, path: Path) -> _SourceLoad:
        records: list[SearchRecord] = []
        skipped_records = 0
        scanned_records = 0
        session_key = _session_key_from_stem(path.stem, source=source)
        try:
            with open(path, encoding="utf-8") as f:
                for line_no, raw_line in enumerate(f, start=1):
                    line = raw_line.strip()
                    if not line:
                        continue
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError:
                        skipped_records += 1
                        continue
                    if not isinstance(data, dict):
                        skipped_records += 1
                        continue
                    scanned_records += 1
                    if source == "sessions" and data.get("_type") == "metadata":
                        stored_key = data.get("key")
                        if isinstance(stored_key, str) and stored_key:
                            session_key = stored_key
                        continue
                    if source == "cold":
                        records.extend(_cold_records_from_json(self.workspace, path, line_no, data))
                        continue
                    record = self._record_from_json(source, path, line_no, session_key, data)
                    if record is not None:
                        records.append(record)
        except OSError as exc:
            logger.warning("session_search failed to read {}: {}", path, exc)
            skipped_records += 1
        return _SourceLoad(records=records, skipped_records=skipped_records, scanned_records=scanned_records)

    def _record_from_json(
        self,
        source: str,
        path: Path,
        line_no: int,
        fallback_session_key: str,
        data: dict[str, Any],
    ) -> SearchRecord | None:
        if source == "sessions":
            return _session_record_from_json(self.workspace, path, line_no, fallback_session_key, data)
        if source == "history":
            return _history_record_from_json(self.workspace, path, line_no, data)
        if source == "webui":
            return _webui_record_from_json(path, line_no, fallback_session_key, data)
        if source == "facts":
            return _fact_record_from_json(self.workspace, path, line_no, data)
        return None


def _normalize_sources(sources: Iterable[str] | None) -> tuple[str, ...]:
    if sources is None:
        return DEFAULT_SOURCES
    out: list[str] = []
    for source in sources:
        if isinstance(source, str):
            value = source.strip().lower()
            if value in SUPPORTED_SOURCES and value not in out:
                out.append(value)
    return tuple(out) or DEFAULT_SOURCES


def _normalize_mode(mode: str | None) -> str:
    value = str(mode or "literal").strip().lower()
    return value if value in SUPPORTED_MODES else "literal"


def _normalize_str_set(values: Iterable[str] | None) -> set[str]:
    if values is None:
        return set()
    return {v.strip().lower() for v in values if isinstance(v, str) and v.strip()}


def _clamp_limit(limit: int | None) -> int:
    if isinstance(limit, bool):
        return DEFAULT_LIMIT
    if not isinstance(limit, int):
        return DEFAULT_LIMIT
    return max(1, min(MAX_LIMIT, limit))


def _session_key_from_channel_chat(channel: str | None, chat_id: str | None) -> str | None:
    if channel and chat_id:
        return f"{channel}:{chat_id}"
    return None


def _session_key_from_stem(stem: str, *, source: str) -> str:
    if source == "webui":
        if stem.startswith("websocket_") and len(stem) > len("websocket_"):
            return f"websocket:{stem[len('websocket_'):]}"
        return "unknown"
    if "_" in stem:
        return stem.replace("_", ":", 1)
    return stem


def _channel_from_session_key(session_key: str) -> str | None:
    if ":" not in session_key:
        return None
    return session_key.split(":", 1)[0]


def _chat_id_from_session_key(session_key: str) -> str | None:
    if ":" not in session_key:
        return None
    return session_key.split(":", 1)[1]


def _parse_datetime_filter(value: str | datetime | None, *, is_until: bool) -> datetime | None:
    if value is None or value == "":
        return None
    if isinstance(value, datetime):
        return _normalize_datetime(value)
    raw = str(value).strip()
    if not raw:
        return None
    if _DATE_ONLY_RE.match(raw):
        day = datetime.fromisoformat(raw).date()
        return datetime.combine(day, datetime_time.max if is_until else datetime_time.min)
    try:
        return _normalize_datetime(datetime.fromisoformat(raw.replace("Z", "+00:00")))
    except ValueError:
        return None


def _normalize_datetime(value: datetime) -> datetime:
    if value.tzinfo is None:
        return value
    return value.astimezone(timezone.utc).replace(tzinfo=None)


def _parse_record_timestamp(value: Any) -> datetime | None:
    if not isinstance(value, str) or not value.strip():
        return None
    try:
        return _normalize_datetime(datetime.fromisoformat(value.strip().replace("Z", "+00:00")))
    except ValueError:
        return None


def _is_old_range(since: datetime | None) -> bool:
    if since is None:
        return False
    return (datetime.now(timezone.utc).replace(tzinfo=None) - since).days > RECENT_SCAN_WINDOW_DAYS


def _in_time_range(timestamp: datetime | None, since: datetime | None, until: datetime | None) -> bool:
    if timestamp is None:
        return since is None and until is None
    if since is not None and timestamp < since:
        return False
    if until is not None and timestamp > until:
        return False
    return True


def _text_from_content(content: Any) -> str:
    if isinstance(content, str):
        return content
    if isinstance(content, list):
        parts: list[str] = []
        for block in content:
            if isinstance(block, dict):
                value = block.get("text")
                if isinstance(value, str):
                    parts.append(value)
                elif isinstance(block.get("content"), str):
                    parts.append(str(block["content"]))
            elif isinstance(block, str):
                parts.append(block)
        return "\n".join(part for part in parts if part)
    return ""


def _session_record_from_json(
    workspace: Path,
    path: Path,
    line_no: int,
    fallback_session_key: str,
    data: dict[str, Any],
) -> SearchRecord | None:
    if data.get("_type") == "metadata":
        return None
    role = str(data.get("role") or "").lower()
    if role not in {"user", "assistant", "tool", "system"}:
        return None
    text = _text_from_content(data.get("content"))
    if not text.strip():
        return None
    rel_path = _relative_path(workspace, path)
    session_key = fallback_session_key
    return SearchRecord(
        source="sessions",
        session_key=session_key,
        role=role,
        timestamp=_parse_record_timestamp(data.get("timestamp")),
        text=text,
        locator={
            "path": rel_path,
            "message_index": max(0, line_no - 2),
            "line": line_no,
            "has_full_content": False,
        },
    )


def _history_record_from_json(
    workspace: Path,
    path: Path,
    line_no: int,
    data: dict[str, Any],
) -> SearchRecord | None:
    text = _text_from_content(data.get("content"))
    if not text.strip():
        return None
    cursor = data.get("cursor")
    rel_path = _relative_path(workspace, path)
    locator: dict[str, Any] = {
        "path": rel_path,
        "line": line_no,
        "has_full_content": len(text) > LONG_HISTORY_TEXT_CHARS,
    }
    if isinstance(cursor, int) and not isinstance(cursor, bool):
        locator["cursor"] = cursor
    return SearchRecord(
        source="history",
        session_key="memory:history",
        role="archive",
        timestamp=_parse_record_timestamp(data.get("timestamp")),
        text=text,
        locator=locator,
    )


def _webui_record_from_json(
    path: Path,
    line_no: int,
    session_key: str,
    data: dict[str, Any],
) -> SearchRecord | None:
    event = data.get("event")
    kind = data.get("kind")
    role = ""
    text = ""
    if event == "user":
        role = "user"
        text = _text_from_content(data.get("text"))
    elif event == "message" and kind not in {"tool_hint", "progress", "reasoning"}:
        role = "assistant"
        text = _text_from_content(data.get("text"))
    elif event == "delta":
        role = "assistant"
        text = _text_from_content(data.get("text"))
    elif event == "message" and kind in {"tool_hint", "progress"}:
        role = "tool"
        text = _text_from_content(data.get("text"))
    if not role or not text.strip():
        return None
    locator: dict[str, Any] = {
        "path": str(path),
        "stem": path.stem,
        "line": line_no,
        "has_full_content": False,
    }
    return SearchRecord(
        source="webui",
        session_key=session_key,
        role=role,
        timestamp=_parse_record_timestamp(data.get("timestamp") or data.get("created_at")),
        text=text,
        locator=locator,
    )


def _fact_record_from_json(
    workspace: Path,
    path: Path,
    line_no: int,
    data: dict[str, Any],
) -> SearchRecord | None:
    status = str(data.get("status") or "").strip().lower()
    if status not in {"active", "pending_confirmation"}:
        return None
    content = _text_from_content(data.get("content"))
    source_excerpt = _text_from_content(data.get("source_excerpt"))
    text = "\n".join(part for part in (content, source_excerpt) if part.strip())
    if not text.strip():
        return None
    rel_path = _relative_path(workspace, path)
    locator: dict[str, Any] = {
        "path": rel_path,
        "line": line_no,
        "fact_id": str(data.get("fact_id") or ""),
        "status": status,
        "category": str(data.get("category") or ""),
        "scope": str(data.get("scope") or ""),
        "owner": str(data.get("owner") or ""),
        "has_full_content": False,
    }
    return SearchRecord(
        source="facts",
        session_key="memory:facts",
        role="archive",
        timestamp=_parse_record_timestamp(data.get("updated_at") or data.get("created_at")),
        text=text,
        locator=locator,
        record_status=status,
    )


def _cold_records_from_json(
    workspace: Path,
    path: Path,
    line_no: int,
    data: dict[str, Any],
) -> list[SearchRecord]:
    messages = data.get("messages")
    if not isinstance(messages, list):
        return []
    session_key = str(data.get("session_key") or "unknown")
    archive_id = str(data.get("archive_id") or "")
    reason = str(data.get("reason") or "")
    archived_at = data.get("archived_at")
    rel_path = _relative_path(workspace, path)
    out: list[SearchRecord] = []
    for index, message in enumerate(messages):
        if not isinstance(message, dict):
            continue
        role = str(message.get("role") or "").lower()
        if role not in {"user", "assistant", "tool", "system"}:
            continue
        text = _text_from_content(message.get("content"))
        if not text.strip():
            continue
        out.append(
            SearchRecord(
                source="cold",
                session_key=session_key,
                role=role,
                timestamp=_parse_record_timestamp(message.get("timestamp") or archived_at),
                text=text,
                locator={
                    "path": rel_path,
                    "archive_id": archive_id,
                    "batch_line": line_no,
                    "message_index": index,
                    "session_key": session_key,
                    "reason": reason,
                    "has_full_content": True,
                },
                record_status=reason,
            )
        )
    return out


@lru_cache(maxsize=2048)
def _relative_path(workspace: Path, path: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return str(path)


def _fingerprint_paths(paths: list[Path]) -> tuple[tuple[str, int, float], ...]:
    out: list[tuple[str, int, float]] = []
    for path in paths:
        try:
            stat = path.stat()
        except OSError:
            continue
        out.append((str(path), stat.st_size, stat.st_mtime))
    return tuple(out)


def _count_literal_matches(text: str, query_lc: str) -> int:
    if not query_lc:
        return 0
    return text.casefold().count(query_lc)


def _format_result(record: SearchRecord, query_lc: str, hit_count: int) -> dict[str, Any]:
    return {
        "source": record.source,
        "session_key": record.session_key,
        "role": record.role,
        "timestamp": _format_timestamp(record.timestamp),
        "snippet": _make_snippet(record.text, query_lc),
        "locator": dict(record.locator),
        "match_count": hit_count,
        "match_type": "literal",
        "score": min(1.0, 0.7 + min(hit_count, 3) * 0.1),
        "redacted": True,
        "record_status": record.record_status,
    }


def _make_snippet(text: str, query_lc: str) -> str:
    cleaned = re.sub(r"\s+", " ", text).strip()
    lower = cleaned.casefold()
    index = lower.find(query_lc)
    if index < 0:
        snippet = truncate_text(cleaned, SNIPPET_MAX_CHARS)
    else:
        half_window = SNIPPET_MAX_CHARS // 2
        start = max(0, index - half_window)
        end = min(len(cleaned), start + SNIPPET_MAX_CHARS)
        start = max(0, end - SNIPPET_MAX_CHARS)
        snippet = cleaned[start:end].strip()
        if start > 0:
            snippet = "..." + snippet
        if end < len(cleaned):
            snippet += "..."
    return redact_memory_text(snippet)


def _format_timestamp(value: datetime | None) -> str | None:
    return value.isoformat() if value is not None else None


def _timestamp_sort_value(value: datetime | None) -> float:
    if value is None:
        return float("-inf")
    return value.timestamp()


def _locator_sort_value(locator: dict[str, Any]) -> str:
    return json.dumps(locator, ensure_ascii=False, sort_keys=True, default=str)


def _recency_sort_key(record: SearchRecord) -> tuple[float, int, str]:
    line = record.locator.get("line")
    cursor = record.locator.get("cursor")
    ordinal = 0
    if isinstance(line, int) and not isinstance(line, bool):
        ordinal = line
    elif isinstance(cursor, int) and not isinstance(cursor, bool):
        ordinal = cursor
    return (
        -_timestamp_sort_value(record.timestamp),
        -ordinal,
        _locator_sort_value(record.locator),
    )


def _sort_key(item: tuple[SearchRecord, int]) -> tuple[float, int, int, str]:
    record, hit_count = item
    return (
        -_timestamp_sort_value(record.timestamp),
        -hit_count,
        SOURCE_PRIORITY.get(record.source, 99),
        _locator_sort_value(record.locator),
    )


def _performance_note(
    *,
    old_range_requested: bool,
    scanned_records: int,
    cache_used: bool,
    total_matches: int,
) -> str | None:
    notes: list[str] = []
    if old_range_requested:
        notes.append(
            "Large time range searched without a persistent index; add since/until "
            "filters when possible."
        )
    if scanned_records >= SCAN_PERFORMANCE_NOTE_THRESHOLD:
        notes.append(f"Scanned {scanned_records} history records.")
    if not cache_used and scanned_records:
        notes.append("Live JSONL scan was used because the request exceeded the recent cache.")
    if total_matches > MAX_LIMIT:
        notes.append("Results were truncated; narrow query or time range for more precision.")
    return " ".join(notes) if notes else None


def _result_identity(row: dict[str, Any]) -> tuple[str, str]:
    return (
        str(row.get("source") or ""),
        _locator_sort_value(row.get("locator") if isinstance(row.get("locator"), dict) else {}),
    )


def _merge_indexed_results(
    literal_results: list[dict[str, Any]],
    indexed_results: list[dict[str, Any]],
    limit: int,
) -> list[dict[str, Any]]:
    merged: dict[tuple[str, str], dict[str, Any]] = {}
    for row in indexed_results:
        merged[_result_identity(row)] = dict(row)
    for row in literal_results:
        identity = _result_identity(row)
        existing = merged.get(identity)
        if existing is None or existing.get("match_type") != "literal":
            merged[identity] = dict(row)
    rows = list(merged.values())

    def key(row: dict[str, Any]) -> tuple[int, float, float, int, str]:
        match_priority = 0 if row.get("match_type") == "literal" else 1
        score = float(row.get("score") or 0.0)
        timestamp = _parse_record_timestamp(row.get("timestamp"))
        source_priority = SOURCE_PRIORITY.get(str(row.get("source") or ""), 99)
        return (
            match_priority,
            -score,
            -_timestamp_sort_value(timestamp),
            source_priority,
            _locator_sort_value(row.get("locator") if isinstance(row.get("locator"), dict) else {}),
        )

    return sorted(rows, key=key)[:limit]


def _join_notes(*values: str | None) -> str | None:
    parts = [str(value).strip() for value in values if str(value or "").strip()]
    return " ".join(parts) if parts else None


__all__ = [
    "DEFAULT_CACHE_RECORDS_PER_SOURCE",
    "DEFAULT_CACHE_TTL_S",
    "DEFAULT_LIMIT",
    "DEFAULT_SOURCES",
    "MAX_LIMIT",
    "SUPPORTED_SOURCES",
    "SUPPORTED_MODES",
    "SearchRecord",
    "SearchResponse",
    "SessionSearchService",
]
