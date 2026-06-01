"""SQLite FTS index for redacted session_search recall."""

from __future__ import annotations

import json
import hashlib
import logging
import re
import sqlite3
import threading
import time
import unicodedata
from collections import Counter
from contextlib import contextmanager
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
from typing import Any, Iterable

from loguru import logger

from OriginAgent.agent.facts import FactStore
from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.config.loader import get_config_path
from OriginAgent.session.cold_archive import SESSION_COLD_ARCHIVE_DIR
from OriginAgent.utils.helpers import ensure_dir, truncate_text

try:  # pragma: no cover - dependency presence is verified by integration tests.
    import jieba
    jieba.setLogLevel(logging.WARNING)
except Exception:  # pragma: no cover
    jieba = None

try:  # pragma: no cover - dependency presence is verified by integration tests.
    from rapidfuzz import fuzz
except Exception:  # pragma: no cover
    fuzz = None


INDEX_DB_RELATIVE_PATH = Path("memory") / "session_search.sqlite3"
INDEX_SCHEMA_VERSION = "1"
INDEX_SNIPPET_MAX_CHARS = 240
MAX_INDEX_TOKEN_LENGTH = 80
MAX_QUERY_TOKENS = 24
SUPPORTED_INDEX_SOURCES = ("sessions", "history", "webui", "facts", "cold")
DEFAULT_INDEX_SOURCES = ("sessions", "history", "webui")
_TOKEN_RE = re.compile(r"[A-Za-z0-9]+(?:[._:/\\-][A-Za-z0-9]+)*")
_CAMEL_BOUNDARY_RE = re.compile(r"(?<=[a-z0-9])(?=[A-Z])")
_CJK_RUN_RE = re.compile(r"[\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff]+")
_SAFE_TOKEN_RE = re.compile(r"^[\w\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff]+$", re.UNICODE)
_UNREDACTED_SECRET_PATTERNS = (
    re.compile(r"\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b"),
    re.compile(r"\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{20,}\b"),
    re.compile(r"\bAKIA[0-9A-Z]{16}\b"),
    re.compile(r"(?i)\bBearer\s+[A-Za-z0-9._~+/=-]{8,}"),
    re.compile(r"(?i)\b(api[_-]?key|token|secret|password)\b\s*[:=]\s*[\"']?[^\"'\s,;]{8,}"),
)


@dataclass(frozen=True)
class IndexedSearchResult:
    source: str
    session_key: str
    role: str
    timestamp: str | None
    snippet: str
    locator: dict[str, Any]
    match_type: str
    score: float
    redacted: bool
    record_status: str

    def to_dict(self) -> dict[str, Any]:
        return {
            "source": self.source,
            "session_key": self.session_key,
            "role": self.role,
            "timestamp": self.timestamp,
            "snippet": self.snippet,
            "locator": dict(self.locator),
            "match_type": self.match_type,
            "score": self.score,
            "redacted": self.redacted,
            "record_status": self.record_status,
        }


class SearchTextNormalizer:
    """Build a space-delimited token string before SQLite FTS sees text."""

    def normalize_for_index(self, text: str) -> str:
        return " ".join(self.tokens(text))

    def tokens(self, text: str) -> list[str]:
        normalized = unicodedata.normalize("NFKC", text or "")
        tokens: list[str] = []
        for match in _TOKEN_RE.finditer(normalized):
            tokens.extend(self._identifier_tokens(match.group(0)))
        for match in _CJK_RUN_RE.finditer(normalized):
            tokens.extend(self._cjk_tokens(match.group(0)))
        return _dedupe_tokens(tokens)

    def query_tokens(self, text: str) -> list[str]:
        return self.tokens(text)[:MAX_QUERY_TOKENS]

    @staticmethod
    def _identifier_tokens(value: str) -> list[str]:
        pieces = [value]
        pieces.extend(re.split(r"[._:/\\-]+", value))
        expanded: list[str] = []
        for piece in pieces:
            expanded.append(piece)
            compact = re.sub(r"[._:/\\-]+", "", piece)
            if compact and compact != piece:
                expanded.append(compact)
            expanded.extend(_CAMEL_BOUNDARY_RE.sub(" ", piece).split())
        return [token.casefold() for token in expanded if token]

    @staticmethod
    def _cjk_tokens(value: str) -> list[str]:
        out: list[str] = [value]
        if jieba is not None:
            out.extend(str(token).strip() for token in jieba.cut_for_search(value))
        else:
            out.extend(value)
        out.extend(value[index:index + 2] for index in range(max(0, len(value) - 1)))
        return out


class SessionSearchIndexService:
    """A discardable redacted FTS cache for session_search."""

    def __init__(
        self,
        workspace: Path,
        *,
        webui_dir: Path | None = None,
        backend: str = "auto",
        semantic_enabled: bool = True,
        rebuild_on_start: bool = False,
    ) -> None:
        self.workspace = Path(workspace)
        self.db_path = self.workspace / INDEX_DB_RELATIVE_PATH
        self._webui_dir = webui_dir
        self.backend = backend
        self.semantic_enabled = semantic_enabled
        self.rebuild_on_start = rebuild_on_start
        self.normalizer = SearchTextNormalizer()
        self._lock = threading.RLock()
        self._refresh_running = False
        self._last_index_error: str | None = None
        self._last_indexed_at: str | None = None
        self._last_skipped_secret_risk_count = 0
        self._fts_available_cache: bool | None = None

    @property
    def enabled(self) -> bool:
        return self.backend != "literal" and self.semantic_enabled

    def fts_available(self) -> bool:
        if self.backend == "literal":
            return False
        if self._fts_available_cache is not None:
            return self._fts_available_cache
        try:
            conn = sqlite3.connect(":memory:")
            try:
                conn.execute("CREATE VIRTUAL TABLE test_fts USING fts5(content)")
            finally:
                conn.close()
        except Exception:
            self._fts_available_cache = False
        else:
            self._fts_available_cache = True
        return self._fts_available_cache

    def refresh_incremental(
        self,
        *,
        sources: Iterable[str] | None = None,
        budget_ms: int | None = None,
        force: bool = False,
    ) -> dict[str, Any]:
        if not self.enabled or not self.fts_available():
            return self.runtime_status()
        if not self._lock.acquire(blocking=False):
            return {**self.runtime_status(), "session_search_refresh_running": True}
        self._refresh_running = True
        started = time.monotonic()
        processed_all = True
        skipped_secret_risk = 0
        try:
            self._ensure_schema()
            with self._connection() as conn:
                for source in _normalize_index_sources(sources):
                    if _budget_exhausted(started, budget_ms):
                        processed_all = False
                        break
                    source_skips, source_complete = self._refresh_source(
                        conn,
                        source,
                        force=force,
                        started=started,
                        budget_ms=budget_ms,
                    )
                    skipped_secret_risk += source_skips
                    if not source_complete:
                        processed_all = False
                        break
                conn.execute("INSERT OR REPLACE INTO meta(key, value) VALUES('schema_version', ?)", (INDEX_SCHEMA_VERSION,))
                conn.execute(
                    "INSERT OR REPLACE INTO meta(key, value) VALUES('index_stale', ?)",
                    ("0" if processed_all else "1",),
                )
                if processed_all:
                    self._last_indexed_at = datetime.now().isoformat()
                    conn.execute(
                        "INSERT OR REPLACE INTO meta(key, value) VALUES('last_indexed_at', ?)",
                        (self._last_indexed_at,),
                    )
                existing = _int_meta(conn, "skipped_secret_risk_count")
                conn.execute(
                    "INSERT OR REPLACE INTO meta(key, value) VALUES('skipped_secret_risk_count', ?)",
                    (str(existing + skipped_secret_risk),),
                )
                self._last_skipped_secret_risk_count += skipped_secret_risk
            self._last_index_error = None
        except Exception as exc:
            self._last_index_error = str(exc)
            logger.exception("session_search index refresh failed")
        finally:
            self._refresh_running = False
            self._lock.release()
        return self.runtime_status()

    def search_indexed(
        self,
        *,
        query: str,
        sources: Iterable[str] | None = None,
        roles: Iterable[str] | None = None,
        session_key: str | None = None,
        channel: str | None = None,
        chat_id: str | None = None,
        since: Any | None = None,
        until: Any | None = None,
        limit: int = 10,
        match_type: str = "fts",
    ) -> dict[str, Any]:
        if not self.enabled:
            return _empty_index_response("session_search semantic index is disabled.")
        if not self.fts_available():
            return _empty_index_response("SQLite FTS5 is not available.")
        self._ensure_schema()
        tokens = self.normalizer.query_tokens(query)
        if not tokens or _only_single_cjk_tokens(tokens):
            return _empty_index_response("Query is too short for indexed multilingual search.")
        match_query = " OR ".join(_quote_fts_token(token) for token in tokens)
        requested_sources = set(_normalize_index_sources(sources))
        requested_roles = {str(role).strip().lower() for role in roles or [] if str(role).strip()}
        target_session_key = session_key or _session_key_from_channel_chat(channel, chat_id)
        since_dt = _parse_datetime_filter(since, is_until=False)
        until_dt = _parse_datetime_filter(until, is_until=True)

        rows: list[IndexedSearchResult] = []
        try:
            with self._connection() as conn:
                sql = (
                    "SELECT d.source, d.session_key, d.role, d.timestamp, d.display_text_redacted, "
                    "d.locator_json, d.record_status, bm25(search_fts) AS rank "
                    "FROM search_fts JOIN docs d ON d.doc_id = search_fts.doc_id "
                    "WHERE search_fts MATCH ?"
                )
                params: list[Any] = [match_query]
                placeholders = ",".join("?" for _ in requested_sources)
                sql += f" AND d.source IN ({placeholders})"
                params.extend(sorted(requested_sources))
                if requested_roles:
                    role_placeholders = ",".join("?" for _ in requested_roles)
                    sql += f" AND d.role IN ({role_placeholders})"
                    params.extend(sorted(requested_roles))
                if target_session_key:
                    sql += " AND d.session_key = ?"
                    params.append(target_session_key)
                if since_dt is not None:
                    sql += " AND d.timestamp IS NOT NULL AND d.timestamp >= ?"
                    params.append(since_dt.isoformat())
                if until_dt is not None:
                    sql += " AND d.timestamp IS NOT NULL AND d.timestamp <= ?"
                    params.append(until_dt.isoformat())
                sql += " ORDER BY rank ASC, d.timestamp DESC LIMIT ?"
                params.append(max(200, min(1000, limit * 20)))
                for row in conn.execute(sql, params):
                    timestamp = row["timestamp"]
                    parsed_ts = _parse_datetime_filter(timestamp, is_until=False)
                    if not _in_time_range(parsed_ts, since_dt, until_dt):
                        continue
                    locator = _json_dict(row["locator_json"])
                    display_text = str(row["display_text_redacted"] or "")
                    rank_score = _rank_to_score(float(row["rank"] or 0.0))
                    if fuzz is not None:
                        fuzzy_score = fuzz.partial_ratio(query, display_text) / 100.0
                        score = max(rank_score, fuzzy_score * 0.95)
                    else:
                        score = rank_score
                    rows.append(
                        IndexedSearchResult(
                            source=str(row["source"]),
                            session_key=str(row["session_key"]),
                            role=str(row["role"]),
                            timestamp=timestamp,
                            snippet=_make_redacted_snippet(display_text, query),
                            locator=locator,
                            match_type=match_type,
                            score=round(max(0.0, min(1.0, score)), 4),
                            redacted=True,
                            record_status=str(row["record_status"] or ""),
                        )
                    )
        except Exception as exc:
            self._last_index_error = str(exc)
            return _empty_index_response(f"Indexed search failed: {exc}")
        rows.sort(
            key=lambda item: (
                -item.score,
                -( _timestamp_value(item.timestamp) ),
                json.dumps(item.locator, ensure_ascii=False, sort_keys=True),
            )
        )
        return {
            "results": [item.to_dict() for item in rows[:limit]],
            "total_matches": len(rows),
            "note": None,
        }

    def runtime_status(self) -> dict[str, Any]:
        stats = {
            "session_search_backend": self._effective_backend(),
            "session_search_semantic_enabled": bool(self.semantic_enabled and self.backend != "literal"),
            "session_search_index_available": False,
            "session_search_indexed_doc_count": 0,
            "session_search_indexed_source_counts": {},
            "session_search_index_stale": False,
            "session_search_refresh_running": self._refresh_running,
            "session_search_last_indexed_at": self._last_indexed_at,
            "session_search_last_index_error": self._last_index_error,
            "session_search_skipped_secret_risk_count": self._last_skipped_secret_risk_count,
        }
        if not self.db_path.exists():
            return stats
        try:
            with self._connection() as conn:
                stats["session_search_index_available"] = _table_exists(conn, "docs")
                if stats["session_search_index_available"]:
                    stats["session_search_indexed_doc_count"] = int(
                        conn.execute("SELECT COUNT(*) FROM docs").fetchone()[0]
                    )
                    counts = Counter()
                    for row in conn.execute("SELECT source, COUNT(*) AS count FROM docs GROUP BY source"):
                        counts[str(row["source"])] = int(row["count"])
                    stats["session_search_indexed_source_counts"] = dict(counts)
                    stats["session_search_index_stale"] = _str_meta(conn, "index_stale") == "1"
                    stats["session_search_last_indexed_at"] = _str_meta(conn, "last_indexed_at")
                    stats["session_search_skipped_secret_risk_count"] = _int_meta(
                        conn,
                        "skipped_secret_risk_count",
                    )
        except Exception as exc:
            stats["session_search_last_index_error"] = str(exc)
        return stats

    def _connect(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=5.0)
        conn.row_factory = sqlite3.Row
        conn.execute("PRAGMA busy_timeout=5000")
        return conn

    @contextmanager
    def _connection(self):
        conn = self._connect()
        try:
            yield conn
            conn.commit()
        except Exception:
            conn.rollback()
            raise
        finally:
            conn.close()

    def _ensure_schema(self) -> None:
        ensure_dir(self.db_path.parent)
        with self._connection() as conn:
            conn.execute("PRAGMA journal_mode=DELETE")
            conn.execute("CREATE TABLE IF NOT EXISTS meta(key TEXT PRIMARY KEY, value TEXT NOT NULL)")
            conn.execute(
                "CREATE TABLE IF NOT EXISTS source_files("
                "source TEXT NOT NULL, path TEXT NOT NULL, size INTEGER NOT NULL, "
                "mtime REAL NOT NULL, indexed_at TEXT NOT NULL, PRIMARY KEY(source, path))"
            )
            conn.execute(
                "CREATE TABLE IF NOT EXISTS docs("
                "doc_id TEXT PRIMARY KEY, source TEXT NOT NULL, source_path TEXT NOT NULL, "
                "session_key TEXT NOT NULL, role TEXT NOT NULL, timestamp TEXT, "
                "display_text_redacted TEXT NOT NULL, locator_json TEXT NOT NULL, "
                "record_status TEXT NOT NULL DEFAULT '', redacted INTEGER NOT NULL DEFAULT 1, "
                "updated_at TEXT NOT NULL)"
            )
            conn.execute(
                "CREATE VIRTUAL TABLE IF NOT EXISTS search_fts USING fts5("
                "doc_id UNINDEXED, index_text, tokenize='unicode61 remove_diacritics 0')"
            )

    def _refresh_source(
        self,
        conn: sqlite3.Connection,
        source: str,
        *,
        force: bool,
        started: float,
        budget_ms: int | None,
    ) -> tuple[int, bool]:
        skipped_secret_risk = 0
        paths = self._source_paths(source)
        current_paths = {str(path) for path in paths}
        for row in conn.execute("SELECT path FROM source_files WHERE source = ?", (source,)):
            existing_path = str(row["path"])
            if existing_path not in current_paths:
                self._delete_path_docs(conn, source, Path(existing_path))
                conn.execute(
                    "DELETE FROM source_files WHERE source = ? AND path = ?",
                    (source, existing_path),
                )
        for path in paths:
            if _budget_exhausted(started, budget_ms):
                return skipped_secret_risk, False
            fingerprint = _file_fingerprint(path)
            if fingerprint is None:
                continue
            if not force and _source_file_fingerprint(conn, source, path) == fingerprint:
                continue
            records = self._records_for_source(source, path)
            if _budget_exhausted(started, budget_ms):
                return skipped_secret_risk, False
            conn.execute("SAVEPOINT refresh_file")
            self._delete_path_docs(conn, source, path)
            for index, record in enumerate(records):
                if _budget_exhausted(started, budget_ms):
                    conn.execute("ROLLBACK TO refresh_file")
                    conn.execute("RELEASE refresh_file")
                    return skipped_secret_risk, False
                redacted_text = redact_memory_text(record["text"])
                if _has_unredacted_secret(redacted_text):
                    skipped_secret_risk += 1
                    continue
                index_text = self.normalizer.normalize_for_index(redacted_text)
                if not index_text.strip():
                    continue
                doc_id = f"{source}:{path}:{index}:{_stable_locator(record['locator'])}"
                conn.execute(
                    "INSERT OR REPLACE INTO docs("
                    "doc_id, source, source_path, session_key, role, timestamp, "
                    "display_text_redacted, locator_json, record_status, redacted, updated_at"
                    ") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?)",
                    (
                        doc_id,
                        source,
                        str(path),
                        record["session_key"],
                        record["role"],
                        record["timestamp"],
                        redacted_text,
                        json.dumps(record["locator"], ensure_ascii=False, sort_keys=True),
                        record.get("record_status", ""),
                        datetime.now().isoformat(),
                    ),
                )
                conn.execute(
                    "INSERT INTO search_fts(rowid, doc_id, index_text) "
                    "VALUES ((SELECT rowid FROM docs WHERE doc_id = ?), ?, ?)",
                    (doc_id, doc_id, index_text),
                )
            conn.execute(
                "INSERT OR REPLACE INTO source_files(source, path, size, mtime, indexed_at) "
                "VALUES (?, ?, ?, ?, ?)",
                (source, str(path), fingerprint[0], fingerprint[1], datetime.now().isoformat()),
            )
            conn.execute("RELEASE refresh_file")
        return skipped_secret_risk, True

    def _delete_path_docs(self, conn: sqlite3.Connection, source: str, path: Path) -> None:
        rowids = [
            int(row["rowid"])
            for row in conn.execute(
                "SELECT rowid FROM docs WHERE source = ? AND source_path = ?",
                (source, str(path)),
            )
        ]
        for rowid in rowids:
            conn.execute("DELETE FROM search_fts WHERE rowid = ?", (rowid,))
        conn.execute("DELETE FROM docs WHERE source = ? AND source_path = ?", (source, str(path)))

    def _source_paths(self, source: str) -> list[Path]:
        if source == "sessions":
            root = self.workspace / "sessions"
            return sorted(root.glob("*.jsonl")) if root.is_dir() else []
        if source == "history":
            path = self.workspace / "memory" / "history.jsonl"
            return [path] if path.is_file() else []
        if source == "webui":
            webui_dir = self._webui_dir or (get_config_path().parent / "webui")
            return sorted(webui_dir.glob("*.jsonl")) if webui_dir.is_dir() else []
        if source == "facts":
            path = self.workspace / "memory" / "facts.jsonl"
            return [path] if path.is_file() else []
        if source == "cold":
            root = self.workspace / SESSION_COLD_ARCHIVE_DIR
            return sorted(root.glob("*.jsonl")) if root.is_dir() else []
        return []

    def _records_for_source(self, source: str, path: Path) -> list[dict[str, Any]]:
        if source == "facts":
            return self._fact_records(path)
        from OriginAgent.session.search import SessionSearchService, _format_timestamp

        service = SessionSearchService(self.workspace, webui_dir=self._webui_dir)
        loaded = service._scan_file(source, path)
        return [
            {
                "text": record.text,
                "source": record.source,
                "session_key": record.session_key,
                "role": record.role,
                "timestamp": _format_timestamp(record.timestamp),
                "locator": record.locator,
                "record_status": record.record_status,
            }
            for record in loaded.records
        ]

    def _fact_records(self, path: Path) -> list[dict[str, Any]]:
        store = FactStore(self.workspace)
        out: list[dict[str, Any]] = []
        for fact in store.read_all():
            if fact.status not in {"active", "pending_confirmation"}:
                continue
            text = "\n".join(part for part in (fact.content, fact.source_excerpt) if part)
            out.append(
                {
                    "text": text,
                    "source": "facts",
                    "session_key": "memory:facts",
                    "role": "archive",
                    "timestamp": fact.updated_at or fact.created_at,
                    "locator": {
                        "path": _relative_path(self.workspace, path),
                        "fact_id": fact.fact_id,
                        "status": fact.status,
                        "category": fact.category,
                        "scope": fact.scope,
                        "owner": fact.owner,
                        "has_full_content": False,
                    },
                    "record_status": fact.status,
                }
            )
        return out

    def _effective_backend(self) -> str:
        if self.backend == "literal":
            return "literal"
        if not self.semantic_enabled:
            return "literal"
        if self.backend == "sqlite_fts":
            return "sqlite_fts" if self.fts_available() else "literal"
        return "sqlite_fts" if self.fts_available() else "literal"


def _normalize_index_sources(sources: Iterable[str] | None) -> tuple[str, ...]:
    if sources is None:
        return DEFAULT_INDEX_SOURCES
    out: list[str] = []
    for source in sources:
        value = str(source).strip().lower()
        if value in SUPPORTED_INDEX_SOURCES and value not in out:
            out.append(value)
    return tuple(out) or DEFAULT_INDEX_SOURCES


def _dedupe_tokens(tokens: Iterable[str]) -> list[str]:
    out: list[str] = []
    seen: set[str] = set()
    for token in tokens:
        cleaned = token.strip().casefold()
        if not cleaned or len(cleaned) > MAX_INDEX_TOKEN_LENGTH:
            continue
        if not _SAFE_TOKEN_RE.match(cleaned):
            continue
        if cleaned in seen:
            continue
        seen.add(cleaned)
        out.append(cleaned)
    return out


def _has_unredacted_secret(text: str) -> bool:
    cleaned = re.sub(r"\[REDACTED_[A-Z_]+\]", "", text or "")
    return any(pattern.search(cleaned) for pattern in _UNREDACTED_SECRET_PATTERNS)


def _only_single_cjk_tokens(tokens: list[str]) -> bool:
    return bool(tokens) and all(len(token) == 1 and _CJK_RUN_RE.fullmatch(token) for token in tokens)


def _budget_exhausted(started: float, budget_ms: int | None) -> bool:
    return budget_ms is not None and budget_ms >= 0 and (time.monotonic() - started) * 1000 >= budget_ms


def _file_fingerprint(path: Path) -> tuple[int, float] | None:
    try:
        stat = path.stat()
    except OSError:
        return None
    return int(stat.st_size), float(stat.st_mtime)


def _source_file_fingerprint(
    conn: sqlite3.Connection,
    source: str,
    path: Path,
) -> tuple[int, float] | None:
    row = conn.execute(
        "SELECT size, mtime FROM source_files WHERE source = ? AND path = ?",
        (source, str(path)),
    ).fetchone()
    if row is None:
        return None
    return int(row["size"]), float(row["mtime"])


def _table_exists(conn: sqlite3.Connection, name: str) -> bool:
    row = conn.execute(
        "SELECT 1 FROM sqlite_master WHERE type IN ('table', 'virtual table') AND name = ?",
        (name,),
    ).fetchone()
    return row is not None


def _str_meta(conn: sqlite3.Connection, key: str) -> str | None:
    try:
        row = conn.execute("SELECT value FROM meta WHERE key = ?", (key,)).fetchone()
    except sqlite3.Error:
        return None
    return str(row["value"]) if row is not None else None


def _int_meta(conn: sqlite3.Connection, key: str) -> int:
    value = _str_meta(conn, key)
    try:
        return int(value or 0)
    except (TypeError, ValueError):
        return 0


def _empty_index_response(note: str) -> dict[str, Any]:
    return {"results": [], "total_matches": 0, "note": note}


def _quote_fts_token(token: str) -> str:
    return '"' + token.replace('"', '""') + '"'


def _json_dict(value: str) -> dict[str, Any]:
    try:
        parsed = json.loads(value)
    except Exception:
        return {}
    return parsed if isinstance(parsed, dict) else {}


def _rank_to_score(rank: float) -> float:
    return 1.0 / (1.0 + abs(rank))


def _make_redacted_snippet(text: str, query: str) -> str:
    cleaned = re.sub(r"\s+", " ", text).strip()
    lower = cleaned.casefold()
    query_lc = (query or "").casefold()
    index = lower.find(query_lc) if query_lc else -1
    if index < 0:
        return redact_memory_text(truncate_text(cleaned, INDEX_SNIPPET_MAX_CHARS))
    half_window = INDEX_SNIPPET_MAX_CHARS // 2
    start = max(0, index - half_window)
    end = min(len(cleaned), start + INDEX_SNIPPET_MAX_CHARS)
    start = max(0, end - INDEX_SNIPPET_MAX_CHARS)
    snippet = cleaned[start:end].strip()
    if start > 0:
        snippet = "..." + snippet
    if end < len(cleaned):
        snippet += "..."
    return redact_memory_text(snippet)


def _stable_locator(locator: dict[str, Any]) -> str:
    payload = json.dumps(locator, ensure_ascii=False, sort_keys=True)
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()[:24]


def _timestamp_value(value: str | None) -> float:
    parsed = _parse_datetime_filter(value, is_until=False)
    if parsed is None:
        return float("-inf")
    return parsed.timestamp()


def _relative_path(workspace: Path, path: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return str(path)


def _session_key_from_channel_chat(channel: str | None, chat_id: str | None) -> str | None:
    if channel and chat_id:
        return f"{channel}:{chat_id}"
    return None


def _parse_datetime_filter(value: Any | None, *, is_until: bool) -> datetime | None:
    from datetime import time as datetime_time

    if value is None or value == "":
        return None
    if isinstance(value, datetime):
        return value.replace(tzinfo=None) if value.tzinfo else value
    raw = str(value).strip()
    if not raw:
        return None
    if re.fullmatch(r"\d{4}-\d{2}-\d{2}", raw):
        day = datetime.fromisoformat(raw).date()
        return datetime.combine(day, datetime_time.max if is_until else datetime_time.min)
    try:
        parsed = datetime.fromisoformat(raw.replace("Z", "+00:00"))
    except ValueError:
        return None
    return parsed.astimezone().replace(tzinfo=None) if parsed.tzinfo else parsed


def _in_time_range(timestamp: datetime | None, since: datetime | None, until: datetime | None) -> bool:
    if timestamp is None:
        return since is None and until is None
    if since is not None and timestamp < since:
        return False
    if until is not None and timestamp > until:
        return False
    return True


__all__ = [
    "INDEX_DB_RELATIVE_PATH",
    "SearchTextNormalizer",
    "SessionSearchIndexService",
]
