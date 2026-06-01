"""Composable safety gates for action authorization decisions."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Protocol, runtime_checkable

VALID_TRIGGERS = {"user_initiated", "scheduled", "system", "subagent"}
VALID_RISKS = {"low", "medium", "high"}


@dataclass
class ActionRequest:
    action: str
    scope: str
    trigger: str
    risk: str
    requested_by: str | None = None
    requires_presence_empty: bool = False
    uses_facts: list[str] = field(default_factory=list)


@dataclass
class ActionDecision:
    decision: str
    reason: str
    supporting_facts: list[str] = field(default_factory=list)
    pending_facts: list[str] = field(default_factory=list)
    presence_status: str = "unknown"


@runtime_checkable
class SafetyGate(Protocol):
    def evaluate(self, request: ActionRequest) -> ActionDecision:
        ...


class CompositeSafetyGate:
    def __init__(self, gates: list[SafetyGate] | None = None):
        self.gates = list(gates or [])

    def evaluate(self, request: ActionRequest) -> ActionDecision:
        final = ActionDecision(decision="allow", reason="safety checks passed")
        for gate in self.gates:
            decision = gate.evaluate(request)
            if decision.decision != "allow":
                return decision
            final = decision
        return final


class DefaultSafetyGate:
    def evaluate(self, request: ActionRequest) -> ActionDecision:
        if request.trigger not in VALID_TRIGGERS:
            return ActionDecision(decision="deny", reason="invalid trigger")
        if request.risk not in VALID_RISKS:
            return ActionDecision(decision="deny", reason="invalid risk")
        return ActionDecision(decision="allow", reason="safety checks passed")


class ActionSafetyGate:
    """Compatibility gate that composes the core default gate with optional domain gates."""

    def __init__(
        self,
        presence_store: object | None = None,
        fact_store: object | None = None,
        *,
        extra_gates: list[SafetyGate] | None = None,
    ):
        gates: list[SafetyGate] = [DefaultSafetyGate()]
        gates.extend(extra_gates or [])
        if presence_store is not None or fact_store is not None:
            raise ValueError("domain-specific safety gates must be passed through extra_gates")
        self._composite = CompositeSafetyGate(gates)

    def evaluate(self, request: ActionRequest) -> ActionDecision:
        return self._composite.evaluate(request)
