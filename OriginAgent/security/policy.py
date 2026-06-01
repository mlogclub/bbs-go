"""Structured policy-denial primitives shared by tool boundaries."""

from __future__ import annotations


class PolicyDeniedError(PermissionError):
    """Raised when a non-bypassable security policy denies an operation."""

    def __init__(
        self,
        message: str,
        *,
        code: str,
        boundary: str,
        policy_rule: str,
    ) -> None:
        super().__init__(message)
        self.code = code
        self.boundary = boundary
        self.policy_rule = policy_rule

