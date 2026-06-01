"""Cold archive for persisted session messages removed from the hot window."""

from __future__ import annotations

import hashlib
import json
import os
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
from typing import Any, Iterator
from uuid import uuid4

from filelock import FileLock

from OriginAgent.utils.helpers import ensure_dir

SESSION_COLD_ARCHIVE_SCHEMA_VERSION = "originagent.session_cold_archive.v1"
SESSION_COLD_ARCHIVE_DIR = Path("memory") / "session_cold_archive"


@dataclass(frozen=True)
class SessionColdArchiveResult:
    archive_id: str
    path: Path
    line: int
    message_count: int
    content_hash: str


class SessionColdArchiveStore:
    """Append-only local archive for session messages before they are trimmed."""

    def __init__(self, workspace: Path):
        self.workspace = Path(workspace)
        self.archive_dir = ensure_dir(self.workspace / SESSION_COLD_ARCHIVE_DIR)
        self._lock_file = self.archive_dir / ".lock"

    def archive(
        self,
        session_key: str,
        messages: list[dict[str, Any]],
        *,
        reason: str,
    ) -> SessionColdArchiveResult | None:
        if not messages:
            return None

        archived_at = datetime.now().isoformat()
        stored_messages = [dict(message) for message in messages]
        content_hash = self.compute_content_hash(stored_messages)
        archive_id = self._archive_id(session_key, reason, archived_at, content_hash)
        record = {
            "schema_version": SESSION_COLD_ARCHIVE_SCHEMA_VERSION,
            "archive_id": archive_id,
            "session_key": session_key,
            "reason": reason,
            "archived_at": archived_at,
            "message_count": len(stored_messages),
            "messages": stored_messages,
            "content_hash": content_hash,
        }
        path = self._path_for_timestamp(archived_at)

        with FileLock(str(self._lock_file)):
            path.parent.mkdir(parents=True, exist_ok=True)
            line_no = self._next_line_no(path)
            with open(path, "a", encoding="utf-8") as f:
                f.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")
                f.flush()
                os.fsync(f.fileno())
            self._fsync_parent(path)

        return SessionColdArchiveResult(
            archive_id=archive_id,
            path=path,
            line=line_no,
            message_count=len(stored_messages),
            content_hash=content_hash,
        )

    def iter_archive_files(self) -> list[Path]:
        if not self.archive_dir.is_dir():
            return []
        return sorted(path for path in self.archive_dir.glob("*.jsonl") if path.is_file())

    def iter_records(self) -> Iterator[tuple[Path, int, dict[str, Any]]]:
        for path in self.iter_archive_files():
            with open(path, encoding="utf-8") as f:
                for line_no, raw_line in enumerate(f, start=1):
                    line = raw_line.strip()
                    if not line:
                        continue
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError:
                        continue
                    if isinstance(data, dict):
                        yield path, line_no, data

    @staticmethod
    def compute_content_hash(messages: list[dict[str, Any]]) -> str:
        payload = json.dumps(messages, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
        return "sha256:" + hashlib.sha256(payload.encode("utf-8")).hexdigest()

    @staticmethod
    def _archive_id(
        session_key: str,
        reason: str,
        archived_at: str,
        content_hash: str,
    ) -> str:
        seed = f"{session_key}\0{reason}\0{archived_at}\0{content_hash}\0{uuid4().hex}"
        return hashlib.sha256(seed.encode("utf-8")).hexdigest()[:32]

    def _path_for_timestamp(self, timestamp: str) -> Path:
        month = timestamp[:7] if len(timestamp) >= 7 else datetime.now().strftime("%Y-%m")
        return self.archive_dir / f"{month}.jsonl"

    @staticmethod
    def _next_line_no(path: Path) -> int:
        try:
            with open(path, encoding="utf-8") as f:
                return sum(1 for _ in f) + 1
        except FileNotFoundError:
            return 1

    @staticmethod
    def _fsync_parent(path: Path) -> None:
        try:
            fd = os.open(str(path.parent), os.O_RDONLY)
        except (PermissionError, OSError):
            return
        try:
            os.fsync(fd)
        finally:
            os.close(fd)


__all__ = [
    "SESSION_COLD_ARCHIVE_DIR",
    "SESSION_COLD_ARCHIVE_SCHEMA_VERSION",
    "SessionColdArchiveResult",
    "SessionColdArchiveStore",
]
