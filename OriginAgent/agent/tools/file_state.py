"""Track file-read state for read-before-edit warnings and read deduplication."""

from __future__ import annotations

import hashlib
import os
from contextvars import ContextVar, Token
from dataclasses import dataclass
from pathlib import Path


@dataclass(slots=True)
class ReadState:
    mtime_ns: int
    size: int
    offset: int
    limit: int | None
    content_hash: str | None
    can_dedup: bool


def _hash_file(p: str) -> str | None:
    try:
        return hashlib.sha256(Path(p).read_bytes()).hexdigest()
    except OSError:
        return None


class FileStates:
    """Per-session read/write tracker.

    Owns its own state dict so read-dedup ("File unchanged since last read")
    and read-before-edit warnings stay scoped to one agent session and do
    not leak across sessions sharing this process.
    """

    __slots__ = ("_state",)

    def __init__(self) -> None:
        self._state: dict[str, ReadState] = {}

    def record_read(
        self,
        path: str | Path,
        offset: int = 1,
        limit: int | None = None,
        content_hash: str | None = None,
    ) -> None:
        """Record that a file was read (called after successful read)."""
        p = str(Path(path).resolve())
        try:
            stat = os.stat(p)
        except OSError:
            return
        self._state[p] = ReadState(
            mtime_ns=stat.st_mtime_ns,
            size=stat.st_size,
            offset=offset,
            limit=limit,
            content_hash=content_hash if content_hash is not None else _hash_file(p),
            can_dedup=True,
        )

    def record_write(self, path: str | Path) -> None:
        """Record that a file was written (updates mtime in state)."""
        p = str(Path(path).resolve())
        try:
            stat = os.stat(p)
        except OSError:
            self._state.pop(p, None)
            return
        self._state[p] = ReadState(
            mtime_ns=stat.st_mtime_ns,
            size=stat.st_size,
            offset=1,
            limit=None,
            content_hash=_hash_file(p),
            can_dedup=False,
        )

    def check_read(self, path: str | Path) -> str | None:
        """Check if a file has been read and is fresh.

        Returns None if OK, or a warning string.
        When mtime changed but file content is identical (e.g. touch, editor save),
        the check passes to avoid false-positive staleness warnings.
        """
        p = str(Path(path).resolve())
        entry = self._state.get(p)
        if entry is None:
            return "Warning: file has not been read yet. Read it first to verify content before editing."
        try:
            stat = os.stat(p)
        except OSError:
            return None
        if stat.st_mtime_ns != entry.mtime_ns or stat.st_size != entry.size:
            if entry.content_hash and _hash_file(p) == entry.content_hash:
                entry.mtime_ns = stat.st_mtime_ns
                entry.size = stat.st_size
                return None
            return "Warning: file has been modified since last read. Re-read to verify content before editing."
        # mtime unchanged - still check content hash to detect quick modifications
        if entry.content_hash and _hash_file(p) != entry.content_hash:
            return "Warning: file has been modified since last read. Re-read to verify content before editing."
        return None

    def is_unchanged(self, path: str | Path, offset: int = 1, limit: int | None = None) -> bool:
        """Return True if file was previously read with same params and content is unchanged."""
        p = str(Path(path).resolve())
        entry = self._state.get(p)
        if entry is None:
            return False
        if not entry.can_dedup:
            return False
        if entry.offset != offset or entry.limit != limit:
            return False
        try:
            stat = os.stat(p)
        except OSError:
            return False

        if stat.st_mtime_ns == entry.mtime_ns and stat.st_size == entry.size:
            return True

        if stat.st_size != entry.size:
            entry.can_dedup = False
            return False

        current_hash = _hash_file(p)
        if not entry.content_hash or current_hash != entry.content_hash:
            entry.can_dedup = False
            return False

        entry.mtime_ns = stat.st_mtime_ns
        entry.size = stat.st_size
        entry.content_hash = current_hash
        return True

    def get(self, path: str | Path) -> ReadState | None:
        """Return the raw ReadState entry for a path, or None."""
        return self._state.get(str(Path(path).resolve()))

    def clear(self) -> None:
        """Clear all tracked state (useful for testing)."""
        self._state.clear()


class FileStateStore:
    """Lookup table for per-session file read/write state."""

    __slots__ = ("_states_by_key",)

    def __init__(self) -> None:
        self._states_by_key: dict[str, FileStates] = {}

    def for_session(self, session_key: str | None) -> FileStates:
        key = session_key or "__default__"
        states = self._states_by_key.get(key)
        if states is None:
            states = FileStates()
            self._states_by_key[key] = states
        return states

    def clear(self) -> None:
        self._states_by_key.clear()


_current_file_states: ContextVar[FileStates | None] = ContextVar(
    "OriginAgent_file_states",
    default=None,
)


def current_file_states(default: FileStates) -> FileStates:
    """Return the FileStates bound to the current agent task, or a fallback."""
    return _current_file_states.get() or default


def bind_file_states(file_states: FileStates) -> Token[FileStates | None]:
    """Bind file read/write state for the current async task."""
    return _current_file_states.set(file_states)


def reset_file_states(token: Token[FileStates | None]) -> None:
    _current_file_states.reset(token)


# Module-level default instance, retained for backward compatibility with
# tests and callers that reach in directly. Per-session callers should hold
# their own FileStates instance instead of touching this one.
_default = FileStates()


def record_read(
    path: str | Path,
    offset: int = 1,
    limit: int | None = None,
    content_hash: str | None = None,
) -> None:
    _default.record_read(path, offset=offset, limit=limit, content_hash=content_hash)


def record_write(path: str | Path) -> None:
    _default.record_write(path)


def check_read(path: str | Path) -> str | None:
    return _default.check_read(path)


def is_unchanged(path: str | Path, offset: int = 1, limit: int | None = None) -> bool:
    return _default.is_unchanged(path, offset=offset, limit=limit)


def clear() -> None:
    _default.clear()


# Legacy attribute for callers that reached into the module-level dict
# directly (filesystem.py used to do this). Kept as a property-like accessor
# so existing imports keep working.
def __getattr__(name: str):
    if name == "_state":
        return _default._state
    raise AttributeError(name)
