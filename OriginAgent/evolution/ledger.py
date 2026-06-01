"""Append-only local evolution event ledger."""

from __future__ import annotations

import hashlib
import json
from dataclasses import dataclass, replace
from pathlib import Path
from typing import Any, Literal

from filelock import FileLock

from OriginAgent.evolution.events import EvolutionEvent
from OriginAgent.evolution.identity import EvolutionIdentityStore, verify_event_signature

HashChainStatus = Literal["ok", "broken"]
LEDGER_ROTATION_EVENT_THRESHOLD = 10_000


@dataclass(frozen=True)
class LedgerVerificationResult:
    chain_integrity: HashChainStatus
    event_count: int
    terminal_event_hash: str | None = None
    broken_at_index: int | None = None
    broken_line: int | None = None
    expected_hash: str | None = None
    actual_hash: str | None = None
    error: str = ""
    unsigned_event_count: int = 0
    invalid_signature_count: int = 0

    @property
    def valid(self) -> bool:
        return self.chain_integrity == "ok"

    @property
    def ok(self) -> bool:
        return self.valid


@dataclass(frozen=True)
class LedgerStatus:
    chain_integrity: HashChainStatus
    event_count: int
    terminal_event_hash: str | None = None
    event_path: str = ""
    segment_id: str = "default"
    rotation_recommended: bool = False
    unsigned_event_count: int = 0
    invalid_signature_count: int = 0
    error: str = ""


class EvolutionLedger:
    """Append-only JSONL ledger with deterministic event hashing."""

    def __init__(
        self,
        workspace: Path,
        event_path: Path | None = None,
        lock_path: Path | None = None,
        identity_store: EvolutionIdentityStore | None = None,
        sign_events: bool = False,
    ) -> None:
        self.workspace = Path(workspace)
        memory_dir = self.workspace / "memory"
        self.event_path = Path(event_path) if event_path is not None else memory_dir / "evolution_events.jsonl"
        self._lock_path = Path(lock_path) if lock_path is not None else memory_dir / ".evolution_ledger.lock"
        self._sign_events = sign_events
        self._identity_store = identity_store

    def append(self, event: EvolutionEvent) -> EvolutionEvent:
        """Append an event and return the immutable event with ledger hashes populated."""

        with self._locked():
            previous_hash = self._last_event_hash_unlocked()
            identity_store = self._identity_store
            actor_public_key = event.actor_public_key
            if self._sign_events:
                identity_store = identity_store or EvolutionIdentityStore()
                actor_public_key = identity_store.public_key_b64
            event_with_previous = replace(
                event,
                actor_public_key=actor_public_key,
                previous_event_hash=previous_hash,
                event_hash="",
                signature="",
            )
            event_hash = compute_event_hash(event_with_previous.to_dict())
            event_to_write = replace(event_with_previous, event_hash=event_hash)
            if self._sign_events and identity_store is not None:
                event_to_write = replace(
                    event_to_write,
                    signature=identity_store.sign(event_hash.encode("utf-8")),
                )
            self.event_path.parent.mkdir(parents=True, exist_ok=True)
            with self.event_path.open("a", encoding="utf-8") as handle:
                handle.write(canonical_dump(event_to_write.to_dict()).decode("utf-8") + "\n")
            return event_to_write

    def verify_chain(self, verify_signatures: bool = False) -> LedgerVerificationResult:
        """Verify the full hash chain.

        This is O(n) over all events. If event count exceeds 10,000, add
        checkpointing or segmented verification before enabling startup-wide checks.
        """

        with self._locked():
            previous_hash: str | None = None
            terminal_hash: str | None = None
            count = 0
            unsigned_count = 0
            invalid_signature_count = 0
            if not self.event_path.exists():
                return LedgerVerificationResult(chain_integrity="ok", event_count=0)
            try:
                with self.event_path.open("r", encoding="utf-8") as handle:
                    for line_number, line in enumerate(handle, start=1):
                        if not line.strip():
                            continue
                        count += 1
                        event = json.loads(line)
                        if event.get("previous_event_hash") != previous_hash:
                            return LedgerVerificationResult(
                                chain_integrity="broken",
                                event_count=count,
                                terminal_event_hash=terminal_hash,
                                broken_at_index=count - 1,
                                broken_line=line_number,
                                expected_hash=previous_hash,
                                actual_hash=event.get("previous_event_hash"),
                                error="previous_event_hash mismatch",
                                unsigned_event_count=unsigned_count,
                                invalid_signature_count=invalid_signature_count,
                            )
                        expected_hash = compute_event_hash(event)
                        if event.get("event_hash") != expected_hash:
                            return LedgerVerificationResult(
                                chain_integrity="broken",
                                event_count=count,
                                terminal_event_hash=terminal_hash,
                                broken_at_index=count - 1,
                                broken_line=line_number,
                                expected_hash=expected_hash,
                                actual_hash=event.get("event_hash"),
                                error="event_hash mismatch",
                                unsigned_event_count=unsigned_count,
                                invalid_signature_count=invalid_signature_count,
                            )
                        if verify_signatures:
                            if event.get("signature"):
                                if not verify_event_signature(event):
                                    invalid_signature_count += 1
                            else:
                                unsigned_count += 1
                        previous_hash = str(event["event_hash"])
                        terminal_hash = previous_hash
            except (OSError, json.JSONDecodeError, TypeError, ValueError) as exc:
                return LedgerVerificationResult(
                    chain_integrity="broken",
                    event_count=count,
                    terminal_event_hash=terminal_hash,
                    broken_at_index=count,
                    broken_line=count + 1,
                    error=str(exc),
                    unsigned_event_count=unsigned_count,
                    invalid_signature_count=invalid_signature_count,
                )
            return LedgerVerificationResult(
                chain_integrity="ok",
                event_count=count,
                terminal_event_hash=terminal_hash,
                unsigned_event_count=unsigned_count,
                invalid_signature_count=invalid_signature_count,
            )

    def status(self, verify_signatures: bool = False) -> LedgerStatus:
        verification = self.verify_chain(verify_signatures=verify_signatures)
        return LedgerStatus(
            chain_integrity=verification.chain_integrity,
            event_count=verification.event_count,
            terminal_event_hash=verification.terminal_event_hash,
            event_path=_relative_to_workspace(self.event_path, self.workspace),
            rotation_recommended=verification.event_count > LEDGER_ROTATION_EVENT_THRESHOLD,
            unsigned_event_count=verification.unsigned_event_count,
            invalid_signature_count=verification.invalid_signature_count,
            error=verification.error,
        )

    def _locked(self) -> FileLock:
        self.event_path.parent.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))

    def _last_event_hash_unlocked(self) -> str | None:
        if not self.event_path.exists():
            return None
        last_line = ""
        with self.event_path.open("r", encoding="utf-8") as handle:
            for line in handle:
                if line.strip():
                    last_line = line
        if not last_line:
            return None
        event = json.loads(last_line)
        event_hash = event.get("event_hash")
        if not isinstance(event_hash, str) or not event_hash:
            raise ValueError("last ledger event is missing event_hash")
        return event_hash


def canonical_dump(event: dict[str, Any]) -> bytes:
    serialized = json.dumps(
        event,
        ensure_ascii=False,
        sort_keys=True,
        separators=(",", ":"),
        allow_nan=False,
    )
    return serialized.encode("utf-8")


def compute_event_hash(event: dict[str, Any]) -> str:
    payload = dict(event)
    payload.pop("event_hash", None)
    payload.pop("signature", None)
    return hashlib.sha256(canonical_dump(payload)).hexdigest()


def _relative_to_workspace(path: Path, workspace: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return path.as_posix()
