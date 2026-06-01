"""Smart-home safety gate with presence- and fact-aware policy checks."""

from __future__ import annotations

from OriginAgent.agent.action_safety import ActionDecision, ActionRequest, DefaultSafetyGate
from OriginAgent.agent.facts import FactRecord, FactStore

from .presence import PresenceStore

USER_TRIGGER = "user_initiated"
PENDING_RELEVANT_CATEGORIES = {"policy", "safety", "temporary"}


class SmartHomeActionSafetyGate:
    def __init__(self, presence_store: PresenceStore, fact_store: FactStore):
        self.presence_store = presence_store
        self.fact_store = fact_store
        self.default_gate = DefaultSafetyGate()

    def evaluate(self, request: ActionRequest) -> ActionDecision:
        default_decision = self.default_gate.evaluate(request)
        if default_decision.decision != "allow":
            return default_decision

        occupancy = self.presence_store.resolve_occupancy()
        presence_status = occupancy.status

        try:
            facts = self.fact_store.read_all()
        except Exception:
            return self._fact_read_failure(request, presence_status)

        fact_decision = self._evaluate_facts(request, facts, presence_status)
        if fact_decision is not None:
            return fact_decision

        if request.risk == "high" and request.trigger != USER_TRIGGER:
            return ActionDecision(
                decision="deny",
                reason="high-risk non-user action denied",
                presence_status=presence_status,
            )

        if request.requires_presence_empty:
            if presence_status == "occupied":
                return ActionDecision(
                    decision="deny",
                    reason="presence is occupied",
                    presence_status=presence_status,
                )
            if presence_status != "empty":
                decision = "ask_confirmation" if request.trigger == USER_TRIGGER else "deny"
                return ActionDecision(
                    decision=decision,
                    reason="presence empty is not established",
                    presence_status=presence_status,
                )

        if request.risk == "high" and request.trigger == USER_TRIGGER and presence_status == "unknown":
            return ActionDecision(
                decision="ask_confirmation",
                reason="high-risk user action needs confirmation with unknown occupancy",
                presence_status=presence_status,
            )

        if (
            request.risk == "medium"
            and request.trigger != USER_TRIGGER
            and presence_status == "unknown"
        ):
            return ActionDecision(
                decision="deny",
                reason="medium-risk non-user action denied with unknown occupancy",
                presence_status=presence_status,
            )

        return ActionDecision(
            decision="allow",
            reason="safety checks passed",
            supporting_facts=self._active_used_facts(request, facts),
            presence_status=presence_status,
        )

    def _evaluate_facts(
        self,
        request: ActionRequest,
        facts: list[FactRecord],
        presence_status: str,
    ) -> ActionDecision | None:
        if request.uses_facts:
            by_id = {fact.fact_id: fact for fact in facts}
            supporting: list[str] = []
            pending_or_missing: list[str] = []
            for fact_id in request.uses_facts:
                fact = by_id.get(fact_id)
                if fact is not None and fact.status == "active":
                    supporting.append(fact.fact_id)
                else:
                    pending_or_missing.append(fact_id)
            if pending_or_missing:
                decision = "ask_confirmation" if request.trigger == USER_TRIGGER else "deny"
                return ActionDecision(
                    decision=decision,
                    reason="requested facts are missing or not active",
                    supporting_facts=supporting,
                    pending_facts=pending_or_missing,
                    presence_status=presence_status,
                )
            return None

        pending = [
            fact.fact_id
            for fact in facts
            if fact.status == "pending_confirmation"
            and self._is_relevant_pending_fact(request, fact)
        ]
        if not pending:
            return None
        decision = "ask_confirmation" if request.trigger == USER_TRIGGER else "deny"
        return ActionDecision(
            decision=decision,
            reason="related facts are pending confirmation",
            pending_facts=pending,
            presence_status=presence_status,
        )

    def _fact_read_failure(
        self,
        request: ActionRequest,
        presence_status: str,
    ) -> ActionDecision:
        if (
            request.risk == "low"
            and request.trigger == USER_TRIGGER
            and not request.requires_presence_empty
            and not request.uses_facts
        ):
            return ActionDecision(
                decision="allow",
                reason="fact store unavailable; low-risk user action has no fact or presence constraint",
                presence_status=presence_status,
            )
        return ActionDecision(
            decision="deny",
            reason="fact store unavailable",
            presence_status=presence_status,
        )

    @staticmethod
    def _is_relevant_pending_fact(request: ActionRequest, fact: FactRecord) -> bool:
        if not _scope_related(request.scope, fact.scope):
            return False
        if fact.category in PENDING_RELEVANT_CATEGORIES:
            return True
        return request.risk in {"medium", "high"}

    @staticmethod
    def _active_used_facts(request: ActionRequest, facts: list[FactRecord]) -> list[str]:
        if not request.uses_facts:
            return []
        by_id = {fact.fact_id: fact for fact in facts if fact.status == "active"}
        return [fact_id for fact_id in request.uses_facts if fact_id in by_id]


def _scope_related(request_scope: str, fact_scope: str) -> bool:
    request_scope = (request_scope or "general").strip().lower()
    fact_scope = (fact_scope or "general").strip().lower()
    return (
        request_scope == fact_scope
        or request_scope.startswith(f"{fact_scope}.")
        or fact_scope.startswith(f"{request_scope}.")
    )
