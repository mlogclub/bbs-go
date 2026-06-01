"""Sanitized telemetry, token budgets, and local proof bundles for evolution modules."""

from __future__ import annotations

import hashlib
import json
import math
import os
import re
import uuid
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Iterable, Mapping

from filelock import FileLock

from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.evolution.capability_gate import EvolutionCapabilityGate
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger, canonical_dump
from OriginAgent.evolution.package import read_package_manifest
from OriginAgent.evolution.verifier import EvolutionModuleVerifier, EvolutionVerificationReport
from OriginAgent.utils.helpers import truncate_text

TELEMETRY_SCHEMA_VERSION = "originagent.evolution.telemetry.v1"
TOKEN_BUDGET_SCHEMA_VERSION = "originagent.evolution.token_budget.v1"
PROOF_BUNDLE_SCHEMA_VERSION = "originagent.evolution.proof_bundle.v1"
MESSAGE_MAX_CHARS = 300
BUDGET_TOKEN_TTL_SECONDS = 60

_URL_QUERY_RE = re.compile(r"(https?://[^\s?\"'<>]+)\?[^\s\"'<>]+")
_GENERIC_SECRET_RE = re.compile(
    r"(?i)\b(token|secret|password|api[_-]?key|authorization)\b"
    r"(\s*[:=]\s*)([\"']?)[^\s\"'&,;<>]+(\3)"
)
_GENERIC_BEARER_RE = re.compile(r"(?i)\bbearer\s+[A-Za-z0-9._~+/=-]{8,}")


@dataclass(frozen=True)
class EvolutionTelemetryResult:
    ok: bool
    status: str
    artifact_digest: str = ""
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    telemetry_path: str = ""
    telemetry_digest: str = ""
    event: EvolutionEvent | None = None
    error: str = ""


@dataclass(frozen=True)
class EvolutionTokenBudgetResult:
    ok: bool
    status: str
    artifact_digest: str = ""
    budget_token: str = ""
    estimated_tokens: int = 0
    consumed_tokens: int = 0
    remaining_tokens: int = 0
    frozen: bool = False
    error: str = ""


@dataclass(frozen=True)
class EvolutionProofBundleResult:
    ok: bool
    status: str
    artifact_digest: str = ""
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    proof_bundle_path: str = ""
    proof_bundle_hash: str = ""
    bundle: dict[str, Any] | None = None
    error: str = ""


class EvolutionTelemetryRecorder:
    """Record privacy-preserving local telemetry and proof snapshots.

    Proof bundle generation verifies the full ledger chain and is O(n) in
    ledger events. Phase 1-H should add checkpoints or segments before startup
    or proof paths depend on ledgers above roughly 10,000 events.
    """

    def __init__(
        self,
        workspace: Path,
        ledger: EvolutionLedger | None = None,
        lock_path: Path | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self.memory_dir = self.workspace / "memory"
        self.telemetry_root = self.memory_dir / "evolution_telemetry"
        self.proof_root = self.memory_dir / "evolution_proofs"
        self.staging_root = self.memory_dir / "evolution_staging"
        self.branch_root = self.memory_dir / "evolution_branches"
        self.ledger = ledger or EvolutionLedger(self.workspace)
        self._lock_path = Path(lock_path) if lock_path is not None else self.memory_dir / ".evolution_telemetry.lock"

    def record(
        self,
        artifact_digest: str,
        *,
        event_kind: str,
        status: str,
        actor: str = "user",
        error_code: str = "",
        exception: BaseException | None = None,
        exception_type: str = "",
        message: str = "",
        stack: str = "",
        tool_name: str = "",
        permission_denied_rule: str = "",
        token_in: int = 0,
        token_out: int = 0,
        total_tokens: int | None = None,
        duration_ms: int | None = None,
        retry_count: int = 0,
        rollback_reason: str = "",
        denial_chain: Iterable[Mapping[str, Any]] | None = None,
    ) -> EvolutionTelemetryResult:
        with self._locked():
            return self._record_unlocked(
                artifact_digest,
                event_kind=event_kind,
                status=status,
                actor=actor,
                error_code=error_code,
                exception=exception,
                exception_type=exception_type,
                message=message,
                stack=stack,
                tool_name=tool_name,
                permission_denied_rule=permission_denied_rule,
                token_in=token_in,
                token_out=token_out,
                total_tokens=total_tokens,
                duration_ms=duration_ms,
                retry_count=retry_count,
                rollback_reason=rollback_reason,
                denial_chain=denial_chain,
            )

    def preflight_token_budget(
        self,
        artifact_digest: str,
        *,
        payload_texts: Iterable[str] = (),
        estimated_tokens: int | None = None,
        actor: str = "user",
    ) -> EvolutionTokenBudgetResult:
        with self._locked():
            manifest, error = self._manifest_for_artifact(artifact_digest)
            if manifest is None:
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="manifest_failed",
                    artifact_digest=artifact_digest,
                    error=error,
                )
            token_budget = manifest.context_budget.get("token_budget")
            if token_budget is None:
                return EvolutionTokenBudgetResult(
                    ok=True,
                    status="budget_not_declared",
                    artifact_digest=artifact_digest,
                )
            if not isinstance(token_budget, int) or token_budget <= 0:
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="invalid_budget",
                    artifact_digest=artifact_digest,
                    error="token_budget must be a positive integer",
                )

            state = self._load_budget_state_unlocked(artifact_digest, token_budget)
            changed = _cleanup_pending(state)
            consumed = int(state.get("consumed_tokens") or 0)
            pending_total = _pending_total(state)
            remaining = max(token_budget - consumed - pending_total, 0)
            if state.get("frozen"):
                if changed:
                    self._write_budget_state_unlocked(artifact_digest, state)
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="blocked",
                    artifact_digest=artifact_digest,
                    consumed_tokens=consumed,
                    remaining_tokens=remaining,
                    frozen=True,
                    error=str(state.get("frozen_reason") or "token budget is frozen"),
                )

            estimate = _coerce_non_negative_int(
                estimated_tokens if estimated_tokens is not None else _estimate_tokens(payload_texts)
            )
            single_limit = max(math.floor(token_budget * 0.8), 1)
            if estimate > single_limit:
                if changed:
                    self._write_budget_state_unlocked(artifact_digest, state)
                self._record_unlocked(
                    artifact_digest,
                    event_kind="token_budget_preflight",
                    status="exceeded",
                    actor=actor,
                    error_code="single_request_token_budget_exceeded",
                    message="estimated tokens exceed per-request token budget",
                    token_in=estimate,
                    total_tokens=estimate,
                )
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="exceeded",
                    artifact_digest=artifact_digest,
                    estimated_tokens=estimate,
                    consumed_tokens=consumed,
                    remaining_tokens=remaining,
                    frozen=False,
                    error="estimated tokens exceed per-request token budget",
                )
            if consumed + pending_total + estimate > token_budget:
                if changed:
                    self._write_budget_state_unlocked(artifact_digest, state)
                self._record_unlocked(
                    artifact_digest,
                    event_kind="token_budget_preflight",
                    status="exceeded",
                    actor=actor,
                    error_code="cumulative_token_budget_exceeded",
                    message="estimated tokens exceed remaining token budget",
                    token_in=estimate,
                    total_tokens=estimate,
                )
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="exceeded",
                    artifact_digest=artifact_digest,
                    estimated_tokens=estimate,
                    consumed_tokens=consumed,
                    remaining_tokens=remaining,
                    frozen=False,
                    error="estimated tokens exceed remaining token budget",
                )

            budget_token = f"budget_{uuid.uuid4().hex}"
            now = datetime.now(timezone.utc)
            pending = dict(state.get("pending") or {})
            pending[budget_token] = {
                "created_at": now.isoformat(),
                "expires_at": datetime.fromtimestamp(
                    now.timestamp() + BUDGET_TOKEN_TTL_SECONDS,
                    tz=timezone.utc,
                ).isoformat(),
                "estimated_tokens": estimate,
            }
            state["pending"] = pending
            state["updated_at"] = now.isoformat()
            self._write_budget_state_unlocked(artifact_digest, state)
            return EvolutionTokenBudgetResult(
                ok=True,
                status="approved",
                artifact_digest=artifact_digest,
                budget_token=budget_token,
                estimated_tokens=estimate,
                consumed_tokens=consumed,
                remaining_tokens=max(token_budget - consumed - _pending_total(state), 0),
                frozen=False,
            )

    def record_postflight_usage(
        self,
        artifact_digest: str,
        *,
        budget_token: str,
        usage: Mapping[str, Any],
        actor: str = "user",
    ) -> EvolutionTokenBudgetResult:
        with self._locked():
            manifest, error = self._manifest_for_artifact(artifact_digest)
            if manifest is None:
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="manifest_failed",
                    artifact_digest=artifact_digest,
                    error=error,
                )
            token_budget = manifest.context_budget.get("token_budget")
            if token_budget is None:
                return EvolutionTokenBudgetResult(ok=True, status="budget_not_declared", artifact_digest=artifact_digest)
            if not isinstance(token_budget, int) or token_budget <= 0:
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="invalid_budget",
                    artifact_digest=artifact_digest,
                    error="token_budget must be a positive integer",
                )

            state = self._load_budget_state_unlocked(artifact_digest, token_budget)
            _cleanup_pending(state)
            consumed = int(state.get("consumed_tokens") or 0)
            if state.get("frozen"):
                self._write_budget_state_unlocked(artifact_digest, state)
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="blocked",
                    artifact_digest=artifact_digest,
                    consumed_tokens=consumed,
                    remaining_tokens=max(token_budget - consumed - _pending_total(state), 0),
                    frozen=True,
                    error=str(state.get("frozen_reason") or "token budget is frozen"),
                )

            pending = dict(state.get("pending") or {})
            if not budget_token or budget_token not in pending:
                self._write_budget_state_unlocked(artifact_digest, state)
                self._record_unlocked(
                    artifact_digest,
                    event_kind="token_budget_postflight",
                    status="invalid_budget_token",
                    actor=actor,
                    error_code="invalid_budget_token",
                    message="postflight usage did not include a valid preflight budget token",
                )
                return EvolutionTokenBudgetResult(
                    ok=False,
                    status="invalid_budget_token",
                    artifact_digest=artifact_digest,
                    consumed_tokens=consumed,
                    remaining_tokens=max(token_budget - consumed - _pending_total(state), 0),
                    error="invalid budget_token",
                )

            actual = _usage_total(usage)
            pending.pop(budget_token, None)
            consumed += actual
            now = datetime.now(timezone.utc).isoformat()
            state["pending"] = pending
            state["consumed_tokens"] = consumed
            state["updated_at"] = now
            status = "recorded"
            ok = True
            frozen = False
            error_message = ""
            if consumed > token_budget:
                state["frozen"] = True
                state["frozen_reason"] = "token_budget_exceeded"
                state["frozen_at"] = now
                status = "budget_exceeded"
                ok = False
                frozen = True
                error_message = "token budget exceeded"
            self._write_budget_state_unlocked(artifact_digest, state)
            self._record_unlocked(
                artifact_digest,
                event_kind="token_budget_postflight",
                status=status,
                actor=actor,
                error_code="token_budget_exceeded" if frozen else "",
                message=error_message or "postflight token usage recorded",
                token_in=_coerce_non_negative_int(usage.get("prompt_tokens") or usage.get("input_tokens") or 0),
                token_out=_coerce_non_negative_int(usage.get("completion_tokens") or usage.get("output_tokens") or 0),
                total_tokens=actual,
            )
            return EvolutionTokenBudgetResult(
                ok=ok,
                status=status,
                artifact_digest=artifact_digest,
                consumed_tokens=consumed,
                remaining_tokens=max(token_budget - consumed - _pending_total(state), 0),
                frozen=frozen,
                error=error_message,
            )

    def build_proof_bundle(self, artifact_digest: str, *, actor: str = "user") -> EvolutionProofBundleResult:
        with self._locked():
            verification = self.ledger.verify_chain()
            if not verification.ok:
                return EvolutionProofBundleResult(
                    ok=False,
                    status="ledger_broken",
                    artifact_digest=artifact_digest,
                    error=verification.error or "ledger chain integrity is broken",
                )
            verified_event = self._latest_ledger_event(artifact_digest, EventType.MODULE_VERIFIED)
            if verified_event is None:
                return EvolutionProofBundleResult(
                    ok=False,
                    status="unverified",
                    artifact_digest=artifact_digest,
                    error="artifact_digest has no verified module event",
                )
            report = EvolutionModuleVerifier(self.workspace, staging_root=self.staging_root).verify(artifact_digest)
            if not report.ok:
                return EvolutionProofBundleResult(
                    ok=False,
                    status="verification_failed",
                    artifact_digest=artifact_digest,
                    error=report.error,
                )

            capability_result = EvolutionCapabilityGate(self.workspace, ledger=self.ledger).snapshot_for_artifact(artifact_digest)
            activation_event_hash = ""
            capability_snapshot_digest = ""
            if capability_result.ok and capability_result.snapshot is not None:
                activation_event = self._latest_ledger_event(artifact_digest, EventType.MODULE_ACTIVATED)
                activation_event_hash = str((activation_event or {}).get("event_hash") or "")
                capability_snapshot_digest = _hash_value(capability_result.snapshot.to_dict())

            bundle: dict[str, Any] = {
                "schema_version": PROOF_BUNDLE_SCHEMA_VERSION,
                "artifact_digest": artifact_digest,
                "module_id": report.module_id,
                "module_type": report.module_type,
                "module_version": report.module_version,
                "verification_event_hash": str(verified_event.get("event_hash") or ""),
                "activation_event_hash": activation_event_hash,
                "verification_report_digest": _hash_value(_verification_report_payload(report)),
                "capability_snapshot_digest": capability_snapshot_digest,
                "telemetry_digest": self.telemetry_digest(artifact_digest),
                "state_branch_digest": self.state_branch_digest(artifact_digest),
                "ledger_tip_hash": verification.terminal_event_hash or "",
                "created_at": datetime.now(timezone.utc).isoformat(),
                "actor": actor,
                "actor_public_key": "",
                "signature_scheme": "ed25519",
                "signature": "",
                "proof_bundle_hash": "",
            }
            bundle["proof_bundle_hash"] = compute_proof_bundle_hash(bundle)
            bundle_path = self.proof_root / artifact_digest / "proof_bundle.json"
            _write_json_atomic(bundle_path, bundle)
            return EvolutionProofBundleResult(
                ok=True,
                status="created",
                artifact_digest=artifact_digest,
                module_id=report.module_id,
                module_type=report.module_type,
                module_version=report.module_version,
                proof_bundle_path=_relative_to_workspace(bundle_path, self.workspace),
                proof_bundle_hash=str(bundle["proof_bundle_hash"]),
                bundle=bundle,
            )

    def telemetry_digest(self, artifact_digest: str) -> str:
        return _hash_value(self._read_telemetry_records(artifact_digest))

    def state_branch_digest(self, artifact_digest: str) -> str:
        summaries: list[dict[str, Any]] = []
        if not self.branch_root.is_dir():
            return ""
        for branch_dir in self.branch_root.iterdir():
            branch_json = branch_dir / "branch.json"
            if not branch_json.exists():
                continue
            try:
                data = json.loads(branch_json.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError):
                continue
            if not isinstance(data, dict) or data.get("artifact_digest") != artifact_digest:
                continue
            summaries.append(
                {
                    "artifact_digest": str(data.get("artifact_digest") or ""),
                    "base_facts_hash": str(data.get("base_facts_hash") or ""),
                    "branch_id": str(data.get("branch_id") or branch_dir.name),
                    "status": str(data.get("status") or ""),
                }
            )
        if not summaries:
            return ""
        summaries.sort(key=lambda item: item["branch_id"])
        return _hash_value(summaries)

    def _record_unlocked(
        self,
        artifact_digest: str,
        *,
        event_kind: str,
        status: str,
        actor: str,
        error_code: str = "",
        exception: BaseException | None = None,
        exception_type: str = "",
        message: str = "",
        stack: str = "",
        tool_name: str = "",
        permission_denied_rule: str = "",
        token_in: int = 0,
        token_out: int = 0,
        total_tokens: int | None = None,
        duration_ms: int | None = None,
        retry_count: int = 0,
        rollback_reason: str = "",
        denial_chain: Iterable[Mapping[str, Any]] | None = None,
    ) -> EvolutionTelemetryResult:
        metadata = self._metadata_for_artifact(artifact_digest)
        message_source = message or (str(exception) if exception is not None else "")
        exception_name = exception_type or (type(exception).__name__ if exception is not None else "")
        sanitized_message = truncate_text(
            sanitize_telemetry_text(message_source, self.workspace),
            MESSAGE_MAX_CHARS,
        )[:MESSAGE_MAX_CHARS]
        total = _coerce_non_negative_int(total_tokens if total_tokens is not None else token_in + token_out)
        record = {
            "schema_version": TELEMETRY_SCHEMA_VERSION,
            "telemetry_id": f"telemetry_{uuid.uuid4().hex}",
            "created_at": datetime.now(timezone.utc).isoformat(),
            "artifact_digest": artifact_digest,
            "module_id": metadata["module_id"],
            "module_type": metadata["module_type"],
            "module_version": metadata["module_version"],
            "event_kind": _safe_small_text(event_kind, self.workspace),
            "status": _safe_small_text(status, self.workspace),
            "error_code": _safe_small_text(error_code, self.workspace),
            "exception_type": _safe_small_text(exception_name, self.workspace),
            "message_redacted": sanitized_message,
            "stack_frame_hashes": _stack_hashes(stack, self.workspace),
            "tool_name": _safe_small_text(tool_name, self.workspace),
            "permission_denied_rule": _safe_small_text(permission_denied_rule, self.workspace),
            "token_in": _coerce_non_negative_int(token_in),
            "token_out": _coerce_non_negative_int(token_out),
            "total_tokens": total,
            "duration_ms": _coerce_non_negative_int(duration_ms or 0),
            "retry_count": _coerce_non_negative_int(retry_count),
            "rollback_reason": _safe_small_text(rollback_reason, self.workspace),
            "denial_chain": _sanitize_denial_chain(denial_chain or (), self.workspace),
        }
        telemetry_path = self._telemetry_path(artifact_digest)
        telemetry_path.parent.mkdir(parents=True, exist_ok=True)
        with telemetry_path.open("a", encoding="utf-8") as handle:
            handle.write(canonical_dump(record).decode("utf-8") + "\n")
        digest = self.telemetry_digest(artifact_digest)
        event = self.ledger.append(
            EvolutionEvent.new(
                EventType.TELEMETRY_RECORDED,
                actor=actor,
                module_id=metadata["module_id"],
                module_version=metadata["module_version"],
                module_type=metadata["module_type"],
                artifact_digest=artifact_digest,
                result={
                    "status": record["status"],
                    "event_kind": record["event_kind"],
                    "telemetry_path": _relative_to_workspace(telemetry_path, self.workspace),
                    "telemetry_digest": digest,
                    "error_code": record["error_code"],
                    "message_redacted": sanitized_message,
                    "token_in": record["token_in"],
                    "token_out": record["token_out"],
                    "total_tokens": record["total_tokens"],
                },
            )
        )
        return EvolutionTelemetryResult(
            ok=True,
            status=str(record["status"]),
            artifact_digest=artifact_digest,
            module_id=metadata["module_id"],
            module_type=metadata["module_type"],
            module_version=metadata["module_version"],
            telemetry_path=_relative_to_workspace(telemetry_path, self.workspace),
            telemetry_digest=digest,
            event=event,
        )

    def _metadata_for_artifact(self, artifact_digest: str) -> dict[str, str]:
        staging_json = self.staging_root / artifact_digest / "staging.json"
        if not staging_json.exists():
            return {"module_id": "", "module_type": "", "module_version": ""}
        try:
            data = json.loads(staging_json.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            return {"module_id": "", "module_type": "", "module_version": ""}
        if not isinstance(data, dict):
            return {"module_id": "", "module_type": "", "module_version": ""}
        return {
            "module_id": str(data.get("module_id") or ""),
            "module_type": str(data.get("module_type") or ""),
            "module_version": str(data.get("version") or ""),
        }

    def _manifest_for_artifact(self, artifact_digest: str) -> tuple[Any | None, str]:
        artifact_dir = self.staging_root / artifact_digest / "artifact"
        try:
            return read_package_manifest(artifact_dir), ""
        except Exception as exc:
            return None, sanitize_telemetry_text(str(exc), self.workspace)

    def _latest_ledger_event(self, artifact_digest: str, event_type: EventType) -> dict[str, Any] | None:
        if not self.ledger.event_path.exists():
            return None
        latest: dict[str, Any] | None = None
        with self.ledger._locked():
            with self.ledger.event_path.open("r", encoding="utf-8") as handle:
                for line in handle:
                    if not line.strip():
                        continue
                    event = json.loads(line)
                    if event.get("event_type") == event_type.value and event.get("artifact_digest") == artifact_digest:
                        latest = event
        return latest

    def _read_telemetry_records(self, artifact_digest: str) -> list[dict[str, Any]]:
        telemetry_path = self._telemetry_path(artifact_digest)
        records: list[dict[str, Any]] = []
        if not telemetry_path.exists():
            return records
        with telemetry_path.open("r", encoding="utf-8") as handle:
            for line in handle:
                if line.strip():
                    record = json.loads(line)
                    if isinstance(record, dict):
                        records.append(record)
        records.sort(key=lambda item: (str(item.get("created_at") or ""), str(item.get("telemetry_id") or "")))
        return records

    def _load_budget_state_unlocked(self, artifact_digest: str, token_budget: int) -> dict[str, Any]:
        path = self._budget_path(artifact_digest)
        if path.exists():
            try:
                data = json.loads(path.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError):
                data = {}
            if isinstance(data, dict) and data.get("schema_version") == TOKEN_BUDGET_SCHEMA_VERSION:
                data.setdefault("artifact_digest", artifact_digest)
                data.setdefault("token_budget", token_budget)
                data.setdefault("consumed_tokens", 0)
                data.setdefault("pending", {})
                data.setdefault("frozen", False)
                data.setdefault("frozen_reason", "")
                data.setdefault("frozen_at", "")
                return data
        now = datetime.now(timezone.utc).isoformat()
        return {
            "schema_version": TOKEN_BUDGET_SCHEMA_VERSION,
            "artifact_digest": artifact_digest,
            "token_budget": token_budget,
            "consumed_tokens": 0,
            "pending": {},
            "frozen": False,
            "frozen_reason": "",
            "frozen_at": "",
            "created_at": now,
            "updated_at": now,
        }

    def _write_budget_state_unlocked(self, artifact_digest: str, state: dict[str, Any]) -> None:
        _write_json_atomic(self._budget_path(artifact_digest), state)

    def _telemetry_path(self, artifact_digest: str) -> Path:
        return self.telemetry_root / artifact_digest / "telemetry.jsonl"

    def _budget_path(self, artifact_digest: str) -> Path:
        return self.telemetry_root / artifact_digest / "budget.json"

    def _locked(self) -> FileLock:
        self.memory_dir.mkdir(parents=True, exist_ok=True)
        return FileLock(str(self._lock_path))


def sanitize_telemetry_text(text: str, workspace: Path) -> str:
    sanitized = redact_memory_text(str(text or ""))
    for path in _sensitive_paths(workspace):
        sanitized = sanitized.replace(path, "<path>")
        sanitized = sanitized.replace(path.replace("\\", "/"), "<path>")
        sanitized = sanitized.replace(path.replace("/", "\\"), "<path>")
    sanitized = _URL_QUERY_RE.sub(r"\1?<url>", sanitized)
    sanitized = _GENERIC_BEARER_RE.sub("Bearer <secret>", sanitized)
    sanitized = _GENERIC_SECRET_RE.sub(lambda match: f"{match.group(1)}{match.group(2)}<secret>", sanitized)
    return sanitized


def compute_proof_bundle_hash(bundle: Mapping[str, Any]) -> str:
    payload = dict(bundle)
    payload.pop("proof_bundle_hash", None)
    payload.pop("signature", None)
    return _hash_value(payload)


def _sensitive_paths(workspace: Path) -> tuple[str, ...]:
    paths: list[str] = []
    for candidate in (Path(workspace), Path.home()):
        try:
            resolved = str(candidate.expanduser().resolve(strict=False))
        except OSError:
            resolved = str(candidate.expanduser())
        if resolved and resolved not in paths:
            paths.append(resolved)
    return tuple(paths)


def _safe_small_text(value: Any, workspace: Path, max_chars: int = MESSAGE_MAX_CHARS) -> str:
    return truncate_text(sanitize_telemetry_text(str(value or ""), workspace), max_chars)[:max_chars]


def _sanitize_denial_chain(denial_chain: Iterable[Mapping[str, Any]], workspace: Path) -> list[dict[str, Any]]:
    sanitized: list[dict[str, Any]] = []
    for item in denial_chain:
        if not isinstance(item, Mapping):
            continue
        row: dict[str, Any] = {}
        for key, value in item.items():
            safe_key = _safe_small_text(key, workspace, 80)
            if isinstance(value, bool | int | float) or value is None:
                row[safe_key] = value
            else:
                row[safe_key] = _safe_small_text(value, workspace)
        sanitized.append(row)
    return sanitized


def _stack_hashes(stack: str, workspace: Path) -> list[str]:
    if not stack:
        return []
    hashes: list[str] = []
    for line in stack.splitlines():
        cleaned = sanitize_telemetry_text(line.strip(), workspace)
        if not cleaned:
            continue
        hashes.append(hashlib.sha256(cleaned.encode("utf-8")).hexdigest())
        if len(hashes) >= 16:
            break
    return hashes


def _estimate_tokens(payload_texts: Iterable[str]) -> int:
    chars = sum(len(str(item or "")) for item in payload_texts)
    return int(math.ceil(chars / 3.5)) if chars else 0


def _usage_total(usage: Mapping[str, Any]) -> int:
    total = usage.get("total_tokens")
    if total is not None and _coerce_non_negative_int(total) > 0:
        return _coerce_non_negative_int(total)
    prompt = usage.get("prompt_tokens", usage.get("input_tokens", 0))
    completion = usage.get("completion_tokens", usage.get("output_tokens", 0))
    return _coerce_non_negative_int(prompt) + _coerce_non_negative_int(completion)


def _coerce_non_negative_int(value: Any) -> int:
    try:
        return max(int(value), 0)
    except (TypeError, ValueError):
        return 0


def _cleanup_pending(state: dict[str, Any]) -> bool:
    pending = dict(state.get("pending") or {})
    now = datetime.now(timezone.utc)
    kept: dict[str, Any] = {}
    for token, entry in pending.items():
        expires_raw = str((entry or {}).get("expires_at") or "")
        try:
            expires = datetime.fromisoformat(expires_raw)
        except ValueError:
            continue
        if expires.tzinfo is None:
            expires = expires.replace(tzinfo=timezone.utc)
        if expires > now:
            kept[token] = entry
    if kept != pending:
        state["pending"] = kept
        state["updated_at"] = now.isoformat()
        return True
    return False


def _pending_total(state: dict[str, Any]) -> int:
    total = 0
    for entry in dict(state.get("pending") or {}).values():
        total += _coerce_non_negative_int((entry or {}).get("estimated_tokens"))
    return total


def _verification_report_payload(report: EvolutionVerificationReport) -> dict[str, Any]:
    return {
        "ok": report.ok,
        "status": report.status,
        "checks": list(report.checks),
        "module_id": report.module_id,
        "module_type": report.module_type,
        "module_version": report.module_version,
        "artifact_digest": report.artifact_digest,
        "staging_path": report.staging_path,
        "error": report.error,
        "permissions_evaluated": report.permissions_evaluated or {},
        "permissions_denied": list(report.permissions_denied),
        "unknown_keys_rejected": list(report.unknown_keys_rejected),
    }


def _hash_value(value: Any) -> str:
    return hashlib.sha256(canonical_dump(value).decode("utf-8").encode("utf-8")).hexdigest()


def _write_json_atomic(path: Path, data: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = path.with_name(f".{path.name}.tmp-{uuid.uuid4().hex}")
    try:
        with tmp_path.open("wb") as handle:
            handle.write(canonical_dump(data))
            handle.write(b"\n")
        os.replace(tmp_path, path)
    finally:
        if tmp_path.exists():
            tmp_path.unlink()


def _relative_to_workspace(path: Path, workspace: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return path.name
