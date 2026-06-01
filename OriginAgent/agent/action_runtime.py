"""Safety-gated action execution boundary.

Phase 3f intentionally provides only a dry-run skeleton. Real device drivers
must plug in behind SafeActionExecutor so they cannot bypass ActionSafetyGate.
"""

from __future__ import annotations

import json
import uuid
from dataclasses import dataclass, field, replace
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable, Protocol

from filelock import FileLock
from loguru import logger

from OriginAgent.agent.action_safety import ActionDecision, ActionRequest, ActionSafetyGate
from OriginAgent.agent.audit import AuditLogger
from OriginAgent.agent.confirmation import (
    REASON_MAX_CHARS,
    ConfirmationManager,
    ConfirmationRequest,
    _sanitize_text,
)
from OriginAgent.agent.permissions import (
    PermissionDecision,
    PermissionRequest,
    PermissionResolver,
)
from OriginAgent.agent.action_privacy import FORBIDDEN_METADATA_KEYS
from OriginAgent.utils.helpers import ensure_dir

ACTION_FORBIDDEN_PAYLOAD_KEYS = {
    *FORBIDDEN_METADATA_KEYS,
    "access_token",
    "api_key",
    "apikey",
    "authorization",
    "bearer",
    "password",
    "private_key",
    "secret",
    "token",
}


@dataclass
class ActionIntent:
    action: str
    scope: str
    trigger: str
    risk: str
    requested_by: str | None = None
    requires_presence_empty: bool = False
    uses_facts: list[str] = field(default_factory=list)
    payload: dict[str, Any] = field(default_factory=dict)
    idempotency_key: str | None = None

    def to_request(self) -> ActionRequest:
        return ActionRequest(
            action=self.action,
            scope=self.scope,
            trigger=self.trigger,
            risk=self.risk,
            requested_by=self.requested_by,
            requires_presence_empty=self.requires_presence_empty,
            uses_facts=list(self.uses_facts),
        )

    def sanitized(self) -> "ActionIntent":
        return replace(
            self,
            uses_facts=list(self.uses_facts),
            payload=sanitize_action_payload(self.payload),
        )


@dataclass
class ActionExecutionResult:
    status: str
    action_id: str
    reason: str
    confirmation_id: str | None = None
    decision: ActionDecision | None = None
    backend_result: dict[str, str] = field(default_factory=dict)
    backend_called: bool = False
    permission_status: str | None = None


@dataclass
class ActionExecutionRecord:
    action_id: str
    intent: ActionIntent
    gate_decision: str
    result_status: str
    created_at: str
    executed_at: str | None = None
    confirmation_id: str | None = None


class ActionBackend(Protocol):
    def execute(self, intent: ActionIntent) -> dict[str, Any]:
        ...


class DryRunActionBackend:
    def execute(self, intent: ActionIntent) -> dict[str, Any]:
        return {
            "dry_run": True,
            "action": intent.action,
            "scope": intent.scope,
            "payload": sanitize_action_payload(intent.payload),
            "idempotency_key": intent.idempotency_key,
        }


class SuccessfulActionKeyStore:
    def __init__(self, workspace: Path, *, max_loaded: int = 10_000):
        self.action_dir = ensure_dir(workspace / "memory" / "action")
        self.keys_file = self.action_dir / "idempotency_keys.jsonl"
        self._lock_file = self.action_dir / ".idempotency.lock"
        self.max_loaded = max_loaded

    def load(self) -> set[str]:
        try:
            lines = self.keys_file.read_text(encoding="utf-8").splitlines()
        except FileNotFoundError:
            return set()
        except OSError as exc:
            logger.warning("Failed to read action idempotency keys from {}: {}", self.keys_file, exc)
            return set()

        keys: list[str] = []
        for line in lines[-self.max_loaded:]:
            if not line.strip():
                continue
            try:
                raw = json.loads(line)
            except json.JSONDecodeError:
                continue
            key = raw.get("idempotency_key") if isinstance(raw, dict) else None
            if isinstance(key, str) and key:
                keys.append(key)
        return set(keys)

    def put(self, idempotency_key: str, *, now: datetime | None = None) -> None:
        event = {
            "idempotency_key": idempotency_key,
            "created_at": _format_datetime(_normalize_datetime(now)),
        }
        line = json.dumps(event, ensure_ascii=False, sort_keys=True) + "\n"
        try:
            with FileLock(str(self._lock_file)):
                with self.keys_file.open("a", encoding="utf-8") as fh:
                    fh.write(line)
                    fh.flush()
        except OSError as exc:
            logger.warning("Failed to persist action idempotency key: {}", exc)


class SafeActionExecutor:
    def __init__(
        self,
        *,
        gate: ActionSafetyGate,
        confirmation_manager: ConfirmationManager,
        backend: ActionBackend,
        permission_resolver: PermissionResolver | None = None,
        audit_logger: AuditLogger | None = None,
        scope_redactor: Callable[[str | None], str | None] | None = None,
    ):
        self.gate = gate
        self.confirmation_manager = confirmation_manager
        self.backend = backend
        self.permission_resolver = permission_resolver or PermissionResolver()
        self.audit_logger = audit_logger
        self._scope_redactor = scope_redactor or _default_scope_redactor
        self.records: list[ActionExecutionRecord] = []
        workspace = getattr(confirmation_manager, "workspace", None)
        self._successful_key_store = (
            SuccessfulActionKeyStore(Path(workspace)) if workspace is not None else None
        )
        self._successful_idempotency_keys: set[str] = (
            self._successful_key_store.load()
            if self._successful_key_store is not None
            else set()
        )

    def submit(
        self,
        intent: ActionIntent,
        *,
        now: datetime | None = None,
    ) -> ActionExecutionResult:
        current_time = _normalize_datetime(now)
        action_id = f"action_{uuid.uuid4().hex[:12]}"
        sanitized_intent = intent.sanitized()
        if self._idempotency_already_executed(sanitized_intent.idempotency_key):
            result = ActionExecutionResult(
                status="already_executed",
                action_id=action_id,
                reason="idempotency key already executed",
            )
            self._record_without_decision(action_id, sanitized_intent, result, current_time)
            return result
        request = intent.to_request()
        decision = self.gate.evaluate(request)

        if decision.decision == "allow":
            permission = self._evaluate_permission(
                action_id,
                sanitized_intent,
                permission="execute_action",
                now=current_time,
            )
            if permission.decision != "allow":
                result = self._permission_result(action_id, decision, permission)
                self._record(action_id, sanitized_intent, decision, result, current_time)
                return result
            result = self._execute_allowed(
                action_id,
                sanitized_intent,
                decision,
                current_time,
            )
            result.permission_status = permission.decision
            self._record(action_id, sanitized_intent, decision, result, current_time)
            self._remember_successful_idempotency(sanitized_intent, result)
            return result

        if decision.decision == "ask_confirmation":
            permission = self._evaluate_permission(
                action_id,
                sanitized_intent,
                permission="confirm_action",
                now=current_time,
            )
            if permission.decision != "allow":
                result = self._permission_result(action_id, decision, permission)
                self._record(action_id, sanitized_intent, decision, result, current_time)
                return result
            try:
                confirmation = self.confirmation_manager.create_from_action_decision(
                    request,
                    decision,
                    now=current_time,
                    action_payload=sanitized_intent.payload,
                    idempotency_key=sanitized_intent.idempotency_key,
                )
            except Exception as exc:
                result = ActionExecutionResult(
                    status="failed",
                    action_id=action_id,
                    reason=_sanitize_text(str(exc), REASON_MAX_CHARS),
                    decision=decision,
                )
                self._record(action_id, sanitized_intent, decision, result, current_time)
                return result
            if confirmation is None:
                result = ActionExecutionResult(
                    status="failed",
                    action_id=action_id,
                    reason="confirmation creation failed",
                    decision=decision,
                )
            else:
                result = ActionExecutionResult(
                    status="pending_confirmation",
                    action_id=action_id,
                    reason=decision.reason,
                    confirmation_id=confirmation.confirmation_id,
                    decision=decision,
                    permission_status=permission.decision,
                )
            self._record(action_id, sanitized_intent, decision, result, current_time)
            return result

        # The production safety gate does not emit notify_only today; this path
        # only preserves compatibility for explicitly injected notification decisions.
        if decision.decision == "notify_only":
            try:
                notification = self.confirmation_manager.create_from_action_decision(
                    request,
                    decision,
                    now=current_time,
                )
            except Exception as exc:
                result = ActionExecutionResult(
                    status="failed",
                    action_id=action_id,
                    reason=_sanitize_text(str(exc), REASON_MAX_CHARS),
                    decision=decision,
                )
                self._record(action_id, sanitized_intent, decision, result, current_time)
                return result
            if notification is None:
                result = ActionExecutionResult(
                    status="failed",
                    action_id=action_id,
                    reason="notification creation failed",
                    decision=decision,
                )
            else:
                result = ActionExecutionResult(
                    status="notified",
                    action_id=action_id,
                    reason=decision.reason,
                    confirmation_id=notification.confirmation_id,
                    decision=decision,
                )
            self._record(action_id, sanitized_intent, decision, result, current_time)
            return result

        if decision.decision == "deny":
            result = ActionExecutionResult(
                status="denied",
                action_id=action_id,
                reason=decision.reason,
                decision=decision,
            )
            self._record(action_id, sanitized_intent, decision, result, current_time)
            return result

        result = ActionExecutionResult(
            status="failed",
            action_id=action_id,
            reason=f"unsupported action decision: {decision.decision}",
            decision=decision,
        )
        self._record(action_id, sanitized_intent, decision, result, current_time)
        return result

    def resume_confirmed(
        self,
        confirmation_id: str,
        *,
        reply: str,
        now: datetime | None = None,
    ) -> ActionExecutionResult:
        current_time = _normalize_datetime(now)
        action_id = f"action_{uuid.uuid4().hex[:12]}"

        confirmation = self.confirmation_manager.expire_confirmation_if_needed(
            confirmation_id,
            now=current_time,
        )
        if confirmation is None:
            return self._resume_terminal_result(
                action_id,
                "denied",
                "confirmation not found",
                current_time,
            )
        if confirmation.status == "expired":
            return self._resume_terminal_result(
                action_id,
                "denied",
                "confirmation expired",
                current_time,
                confirmation_id=confirmation_id,
            )
        if confirmation.kind == "notify_only":
            return self._resume_terminal_result(
                action_id,
                "denied",
                "notification is not executable",
                current_time,
                confirmation_id=confirmation_id,
            )
        if confirmation.kind != "action_confirmation":
            return self._resume_terminal_result(
                action_id,
                "denied",
                "confirmation is not an action confirmation",
                current_time,
                confirmation_id=confirmation_id,
            )
        if confirmation.consumed_at is not None:
            return self._resume_terminal_result(
                action_id,
                "already_executed",
                "confirmation already consumed",
                current_time,
                confirmation_id=confirmation_id,
            )

        if confirmation.status == "pending":
            reply_result = self.confirmation_manager.resolve_user_reply(
                confirmation_id,
                reply,
                now=current_time,
            )
            if reply_result.decision != "confirmed":
                status = (
                    "pending_confirmation"
                    if reply_result.decision == "unclear"
                    else "denied"
                )
                return self._resume_terminal_result(
                    action_id,
                    status,
                    reply_result.reason,
                    current_time,
                    confirmation_id=confirmation_id,
                )
            confirmation = self.confirmation_manager.store.get(confirmation_id)
            if confirmation is None:
                return self._resume_terminal_result(
                    action_id,
                    "failed",
                    "confirmation disappeared after confirmation",
                    current_time,
                    confirmation_id=confirmation_id,
                )

        if confirmation.status != "confirmed_once":
            return self._resume_terminal_result(
                action_id,
                "denied",
                f"confirmation is {confirmation.status}",
                current_time,
                confirmation_id=confirmation_id,
            )
        if confirmation.consumed_at is not None:
            return self._resume_terminal_result(
                action_id,
                "already_executed",
                "confirmation already consumed",
                current_time,
                confirmation_id=confirmation_id,
            )

        intent = _intent_from_confirmation(confirmation)
        if intent is None:
            return self._resume_terminal_result(
                action_id,
                "failed",
                "confirmation action snapshot is incomplete",
                current_time,
                confirmation_id=confirmation_id,
            )
        if self._idempotency_already_executed(intent.idempotency_key):
            return self._resume_terminal_result(
                action_id,
                "already_executed",
                "idempotency key already executed",
                current_time,
                confirmation_id=confirmation_id,
                intent=intent,
            )

        decision = self.gate.evaluate(intent.to_request())
        if decision.decision == "ask_confirmation":
            result = ActionExecutionResult(
                status="pending_confirmation_still_required",
                action_id=action_id,
                reason=decision.reason,
                confirmation_id=confirmation_id,
                decision=decision,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result
        if decision.decision == "deny":
            result = ActionExecutionResult(
                status="denied",
                action_id=action_id,
                reason=decision.reason,
                confirmation_id=confirmation_id,
                decision=decision,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result
        # Notify-only confirmations are informational and must remain non-executable.
        if decision.decision == "notify_only":
            result = ActionExecutionResult(
                status="notified",
                action_id=action_id,
                reason=decision.reason,
                confirmation_id=confirmation_id,
                decision=decision,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result
        if decision.decision != "allow":
            result = ActionExecutionResult(
                status="failed",
                action_id=action_id,
                reason=f"unsupported action decision: {decision.decision}",
                confirmation_id=confirmation_id,
                decision=decision,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result

        confirm_permission = self._evaluate_permission(
            action_id,
            intent,
            permission="confirm_action",
            confirmation_id=confirmation_id,
            now=current_time,
        )
        if confirm_permission.decision != "allow":
            result = self._permission_result(
                action_id,
                decision,
                confirm_permission,
                confirmation_id=confirmation_id,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result
        execute_permission = self._evaluate_permission(
            action_id,
            intent,
            permission="execute_action",
            confirmation_id=confirmation_id,
            now=current_time,
        )
        if execute_permission.decision != "allow":
            result = self._permission_result(
                action_id,
                decision,
                execute_permission,
                confirmation_id=confirmation_id,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result

        claimed = self.confirmation_manager.claim_consumption(
            confirmation_id,
            now=current_time,
        )
        if claimed is None:
            result = ActionExecutionResult(
                status="already_executed",
                action_id=action_id,
                reason="confirmation could not be consumed",
                confirmation_id=confirmation_id,
                decision=decision,
            )
            self._record(action_id, intent, decision, result, current_time)
            return result

        result = self._execute_allowed(action_id, intent, decision, current_time)
        result.confirmation_id = confirmation_id
        result.permission_status = execute_permission.decision
        self._record(action_id, intent, decision, result, current_time)
        self._remember_successful_idempotency(intent, result)
        return result

    def _execute_allowed(
        self,
        action_id: str,
        intent: ActionIntent,
        decision: ActionDecision,
        now: datetime,
    ) -> ActionExecutionResult:
        try:
            raw_backend_result = self.backend.execute(intent)
            is_dry_run = _is_dry_run_result(raw_backend_result)
            backend_result = sanitize_action_payload(raw_backend_result)
        except Exception as exc:
            return ActionExecutionResult(
                status="failed",
                action_id=action_id,
                reason=_sanitize_text(str(exc), REASON_MAX_CHARS),
                decision=decision,
                backend_called=True,
            )

        status = "dry_run" if is_dry_run else "executed"
        return ActionExecutionResult(
            status=status,
            action_id=action_id,
            reason=decision.reason,
            decision=decision,
            backend_result=backend_result,
            backend_called=True,
        )

    def _record(
        self,
        action_id: str,
        intent: ActionIntent,
        decision: ActionDecision,
        result: ActionExecutionResult,
        now: datetime,
    ) -> None:
        executed_at = _format_datetime(now) if result.status in {"executed", "dry_run"} else None
        self.records.append(
            ActionExecutionRecord(
                action_id=action_id,
                intent=intent,
                gate_decision=decision.decision,
                result_status=result.status,
                created_at=_format_datetime(now),
                executed_at=executed_at,
                confirmation_id=result.confirmation_id,
            )
        )
        self._audit_action_decision(
            action_id,
            intent,
            decision.decision,
            result,
            now,
            presence_status=decision.presence_status,
        )

    def _record_without_decision(
        self,
        action_id: str,
        intent: ActionIntent,
        result: ActionExecutionResult,
        now: datetime,
        *,
        gate_decision: str = "idempotent_replay",
    ) -> None:
        self.records.append(
            ActionExecutionRecord(
                action_id=action_id,
                intent=intent,
                gate_decision=gate_decision,
                result_status=result.status,
                created_at=_format_datetime(now),
                executed_at=None,
                confirmation_id=result.confirmation_id,
            )
        )
        self._audit_action_decision(
            action_id,
            intent,
            gate_decision,
            result,
            now,
            presence_status=None,
        )

    def _resume_terminal_result(
        self,
        action_id: str,
        status: str,
        reason: str,
        now: datetime,
        *,
        confirmation_id: str | None = None,
        intent: ActionIntent | None = None,
    ) -> ActionExecutionResult:
        result = ActionExecutionResult(
            status=status,
            action_id=action_id,
            reason=_sanitize_text(reason, REASON_MAX_CHARS),
            confirmation_id=confirmation_id,
        )
        self._record_without_decision(
            action_id,
            intent or _empty_action_intent(),
            result,
            now,
            gate_decision="resume_precheck",
        )
        return result

    def _idempotency_already_executed(self, idempotency_key: str | None) -> bool:
        return bool(idempotency_key and idempotency_key in self._successful_idempotency_keys)

    def _remember_successful_idempotency(
        self,
        intent: ActionIntent,
        result: ActionExecutionResult,
    ) -> None:
        if intent.idempotency_key and result.status in {"executed", "dry_run"}:
            if intent.idempotency_key in self._successful_idempotency_keys:
                return
            self._successful_idempotency_keys.add(intent.idempotency_key)
            if self._successful_key_store is not None:
                self._successful_key_store.put(intent.idempotency_key)

    def _evaluate_permission(
        self,
        action_id: str,
        intent: ActionIntent,
        *,
        permission: str,
        confirmation_id: str | None = None,
        now: datetime | None = None,
    ) -> PermissionDecision:
        request = PermissionRequest(
            actor_id=intent.requested_by,
            action=intent.action,
            scope=intent.scope,
            risk=intent.risk,
            trigger=intent.trigger,
            permission=permission,
            attributes=_permission_attributes(intent),
        )
        decision = self.permission_resolver.evaluate(request)
        self._audit_permission_decision(
            action_id,
            request,
            decision,
            confirmation_id=confirmation_id,
            now=now,
        )
        return decision

    @staticmethod
    def _permission_result(
        action_id: str,
        decision: ActionDecision,
        permission: PermissionDecision,
        *,
        confirmation_id: str | None = None,
    ) -> ActionExecutionResult:
        status = "ask_admin" if permission.decision == "ask_admin" else "denied"
        return ActionExecutionResult(
            status=status,
            action_id=action_id,
            reason=permission.reason,
            confirmation_id=confirmation_id,
            decision=decision,
            permission_status=permission.decision,
        )

    def _audit_action_decision(
        self,
        action_id: str,
        intent: ActionIntent,
        gate_decision: str,
        result: ActionExecutionResult,
        now: datetime,
        *,
        presence_status: str | None,
    ) -> None:
        if self.audit_logger is None:
            return
        try:
            metadata: dict[str, Any] = {
                "gate_decision": gate_decision,
                "result_status": result.status,
                "backend_called": result.backend_called,
                "idempotency_key_present": bool(intent.idempotency_key),
                "payload_keys": sorted(intent.payload.keys()),
                "backend_result_keys": sorted(result.backend_result.keys()),
            }
            typed_action_type = intent.payload.get("action_type")
            typed_action_domain = intent.payload.get("domain")
            if isinstance(typed_action_type, str) and isinstance(typed_action_domain, str):
                metadata["schema_validated"] = True
                metadata["typed_action_type"] = typed_action_type
                metadata["typed_action_domain"] = typed_action_domain
            if presence_status is not None:
                metadata["presence_status"] = presence_status
            if result.permission_status is not None:
                metadata["permission_status"] = result.permission_status
            self.audit_logger.log_action_decision(
                action_id=action_id,
                confirmation_id=result.confirmation_id,
                actor_id=intent.requested_by,
                action=intent.action,
                scope=_audit_scope(intent, result, self._scope_redactor),
                risk=intent.risk,
                trigger=intent.trigger,
                decision=result.status,
                reason=result.reason,
                metadata=metadata,
                created_at=now,
            )
        except Exception as exc:
            logger.warning("Audit action decision write failed: {}", exc)

    def _audit_permission_decision(
        self,
        action_id: str,
        request: PermissionRequest,
        decision: PermissionDecision,
        *,
        confirmation_id: str | None = None,
        now: datetime | None = None,
    ) -> None:
        if self.audit_logger is None:
            return
        try:
            self.audit_logger.log_permission_decision(
                action_id=action_id,
                confirmation_id=confirmation_id,
                actor_id=request.actor_id,
                action=request.action,
                scope=_audit_permission_scope(request, self._scope_redactor),
                risk=request.risk,
                trigger=request.trigger,
                permission=request.permission,
                attributes=request.attributes,
                decision=decision.decision,
                reason=decision.reason,
                actor_role=decision.actor_role,
                created_at=now,
            )
        except Exception as exc:
            logger.warning("Audit permission decision write failed: {}", exc)


def sanitize_action_payload(payload: dict[str, Any] | Any) -> dict[str, str]:
    if not isinstance(payload, dict):
        return {}
    sanitized: dict[str, str] = {}
    for key, value in payload.items():
        if not isinstance(key, str):
            continue
        normalized_key = key.strip()
        if not normalized_key:
            continue
        if normalized_key.casefold() in ACTION_FORBIDDEN_PAYLOAD_KEYS:
            continue
        if value is None:
            continue
        sanitized[normalized_key] = _sanitize_text(
            _stringify_payload_value(_sanitize_payload_value(value)),
            REASON_MAX_CHARS,
        )
    return sanitized


def _intent_from_confirmation(confirmation: ConfirmationRequest) -> ActionIntent | None:
    if not all([confirmation.action, confirmation.scope, confirmation.trigger, confirmation.risk]):
        return None
    return ActionIntent(
        action=confirmation.action,
        scope=confirmation.scope,
        trigger=confirmation.trigger,
        risk=confirmation.risk,
        requested_by=confirmation.requested_by,
        requires_presence_empty=confirmation.requires_presence_empty,
        uses_facts=list(confirmation.uses_facts),
        payload=sanitize_action_payload(confirmation.action_payload),
        idempotency_key=confirmation.idempotency_key,
    )


def _empty_action_intent() -> ActionIntent:
    return ActionIntent(
        action="unknown",
        scope="unknown",
        trigger="system",
        risk="low",
    )


def _sanitize_payload_value(value: Any) -> Any:
    if isinstance(value, dict):
        nested: dict[str, Any] = {}
        for key, item in value.items():
            if not isinstance(key, str):
                continue
            normalized_key = key.strip()
            if not normalized_key:
                continue
            if normalized_key.casefold() in ACTION_FORBIDDEN_PAYLOAD_KEYS:
                continue
            if item is None:
                continue
            nested[normalized_key] = _sanitize_payload_value(item)
        return nested
    if isinstance(value, list):
        return [
            _sanitize_payload_value(item)
            for item in value
            if item is not None
        ]
    if isinstance(value, str):
        return _sanitize_text(value, REASON_MAX_CHARS)
    return value


def _stringify_payload_value(value: Any) -> str:
    if isinstance(value, str):
        return value
    if isinstance(value, bool | int | float):
        return str(value)
    try:
        return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
    except (TypeError, ValueError):
        return str(value)


def _is_dry_run_result(result: Any) -> bool:
    if not isinstance(result, dict):
        return False
    value = result.get("dry_run")
    if isinstance(value, bool):
        return value
    if isinstance(value, str):
        return value.strip().casefold() in {"true", "1", "yes"}
    if isinstance(value, int | float):
        return value != 0
    return False


def _normalize_datetime(value: datetime | None) -> datetime:
    if value is None:
        return datetime.now(timezone.utc)
    if value.tzinfo is None:
        return value.replace(tzinfo=timezone.utc)
    return value.astimezone(timezone.utc)


def _format_datetime(value: datetime) -> str:
    return _normalize_datetime(value).isoformat()


def _audit_scope(
    intent: ActionIntent,
    result: ActionExecutionResult,
    scope_redactor: Callable[[str | None], str | None] = None,
) -> str:
    if scope_redactor is None:
        scope_redactor = _default_scope_redactor
    if (
        result.backend_result.get("device_id_present") in {"True", True}
        or intent.payload.get("device_id")
    ):
        return scope_redactor(intent.scope) or "unknown"
    return intent.scope


def _audit_permission_scope(
    request: PermissionRequest,
    scope_redactor: Callable[[str | None], str | None] = None,
) -> str:
    if scope_redactor is None:
        scope_redactor = _default_scope_redactor
    if request.attribute("device_domain") not in {None, "", "general"}:
        return scope_redactor(request.scope) or "unknown"
    return request.scope


def _permission_attributes(intent: ActionIntent) -> dict[str, str]:
    attributes: dict[str, str] = {}
    domain = intent.payload.get("domain")
    if isinstance(domain, str) and domain.strip():
        normalized_domain = domain.strip().lower()
    else:
        normalized_domain = "general"
    if normalized_domain and normalized_domain != "general":
        attributes["device_domain"] = normalized_domain
    if intent.payload.get("device_id"):
        attributes["device_id_present"] = "true"
    return attributes


def _default_scope_redactor(scope: str | None) -> str | None:
    if scope is None:
        return None
    normalized = str(scope).strip().lower()
    if not normalized:
        return None
    normalized = normalized.replace(" ", ".")
    parts = [part for part in normalized.split(".") if part]
    if len(parts) >= 3:
        return ".".join([*parts[:-1], "<target>"])
    normalized = ".".join(parts)
    return normalized or None
