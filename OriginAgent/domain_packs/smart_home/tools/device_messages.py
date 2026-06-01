"""Stable user-facing messages for device tool results."""

from __future__ import annotations

DRY_RUN_ACCEPTED = "Lighting action accepted in dry-run mode."
EXECUTED = "Lighting action completed."
POLICY_DENIED = "Lighting action was denied by policy."
VALIDATION_FAILED = "Lighting action could not be submitted."
CONFIRMATION_REQUIRED = "Confirmation required before this lighting action can run."
BACKEND_FAILED = "Lighting action could not be completed."


def device_human_message(status: str, fallback: str | None = None) -> str:
    if status == "dry_run":
        return DRY_RUN_ACCEPTED
    if status in {"executed", "success"}:
        return EXECUTED
    if status in {"pending_confirmation", "needs_confirmation"}:
        return CONFIRMATION_REQUIRED
    if status in {"denied", "deny"}:
        return POLICY_DENIED
    if status in {"backend_failed", "backend_error"}:
        return BACKEND_FAILED
    return VALIDATION_FAILED
