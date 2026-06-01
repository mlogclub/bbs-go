"""Read-only self-model aggregation and rendering."""

from __future__ import annotations

from collections import Counter
from contextlib import suppress
from datetime import datetime, timezone
from importlib.resources import files as pkg_files
from pathlib import Path
from typing import Any

from OriginAgent.agent.confirmation import ConfirmationRequest, PendingConfirmationStore
from OriginAgent.agent.domain_pack_governance import DomainPackGovernanceService
from OriginAgent.agent.facts import FactStore, summarize_facts
from OriginAgent.agent.memory import MemoryStore, redact_memory_text
from OriginAgent.agent.skills import SkillsLoader
from OriginAgent.agent.workflow_artifacts import (
    list_workflow_artifact_records,
    summarize_workflow_artifacts,
)
from OriginAgent.utils.helpers import truncate_text

_LIMITATION_MAX_CHARS = 220
_RENDER_LIST_MAX_ITEMS = 8


class SelfModelService:
    """Build a deterministic, read-only self-model snapshot."""

    def __init__(
        self,
        workspace: Path,
        *,
        registry: Any | None = None,
        sessions: Any | None = None,
        pending_queues: dict[str, Any] | None = None,
        cron_service: Any | None = None,
        confirmation_store: PendingConfirmationStore | None = None,
        audit_mode: str = "minimal",
        runtime_profile: str = "default",
        domain_pack_manager: Any | None = None,
        background_review_service: Any | None = None,
        curator_service: Any | None = None,
        skills_loader: SkillsLoader | None = None,
        memory_store: MemoryStore | None = None,
        facts_store: FactStore | None = None,
        review_store: Any | None = None,
        domain_governance_service: DomainPackGovernanceService | None = None,
        background_review_enabled: bool | None = None,
        curator_enabled: bool | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self._registry = registry
        self._sessions = sessions
        self._pending_queues = pending_queues or {}
        self._cron_service = cron_service
        self._confirmation_store = confirmation_store
        self._audit_mode = audit_mode
        self._runtime_profile = runtime_profile
        self._domain_pack_manager = domain_pack_manager
        self._background_review_service = background_review_service
        self._curator_service = curator_service
        self._skills_loader = skills_loader
        self._memory_store = memory_store
        self._facts_store = facts_store
        self._review_store = review_store
        self._domain_governance_service = domain_governance_service
        self._background_review_enabled = background_review_enabled
        self._curator_enabled = curator_enabled

    def build(self) -> dict[str, Any]:
        reviews, pending_reviews = self._build_reviews()
        confirmations, pending_confirmations = self._build_confirmations()
        domains = self._build_domains()
        skills = self._build_skills()
        workflows = self._build_workflows()
        facts = self._build_facts()
        memory = self._build_memory()
        limitations = self._build_limitations(
            domains=domains.get("items", []),
            skills=skills.get("items", []),
            workflows=workflows.get("items", []),
            pending_reviews=pending_reviews,
            pending_confirmations=pending_confirmations,
        )
        return {
            "schema_version": 1,
            "generated_at": datetime.now(timezone.utc).isoformat(),
            "identity": {
                "agent_name": "OriginAgent",
                "workspace_name": self.workspace.name or "workspace",
                "runtime_profile": self._runtime_profile,
                "audit_mode": self._audit_mode,
            },
            "runtime": {
                "registered_tools_count": _safe_len(getattr(self._registry, "tool_names", [])),
                "active_sessions_count": _session_count(self._sessions),
                "pending_queue_count": len(self._pending_queues),
                "cron_available": self._cron_service is not None,
                "confirmation_available": self._confirmation_store is not None or self.workspace.exists(),
                "background_review_enabled": self._review_enabled(
                    self._background_review_service,
                    self._background_review_enabled,
                ),
                "curator_enabled": self._review_enabled(
                    self._curator_service,
                    self._curator_enabled,
                ),
            },
            "domains": domains,
            "skills": skills,
            "workflows": workflows,
            "facts": facts,
            "memory": memory,
            "reviews": reviews,
            "confirmations": confirmations,
            "limitations": limitations,
        }

    def _build_domains(self) -> dict[str, Any]:
        try:
            service = self._domain_governance()
            return {
                "stats": service.stats(),
                "items": service.list_records(limit=500),
            }
        except Exception:
            return {
                "stats": {
                    "workspace_domain_pack_count": 0,
                    "builtin_domain_pack_count": 0,
                    "domain_pack_status_counts": {},
                    "active_domain_pack_count": 0,
                    "domain_pack_override_count": 0,
                    "domain_pack_eval_status_counts": {},
                    "last_domain_pack_event_at": None,
                },
                "items": [],
            }

    def _build_skills(self) -> dict[str, Any]:
        try:
            loader = self._skills()
            entries = loader.list_skills(filter_unavailable=False)
            return {
                "stats": loader.lifecycle.stats(entries),
                "items": loader.list_skill_records(filter_unavailable=False),
            }
        except Exception:
            return {
                "stats": {
                    "skills_count": 0,
                    "workspace_skills_count": 0,
                    "skill_lifecycle_status_counts": {},
                    "skill_verification_status_counts": {},
                    "unverified_skill_count": 0,
                    "deprecated_skill_count": 0,
                    "rejected_skill_count": 0,
                    "always_workspace_skill_count": 0,
                },
                "items": [],
            }

    def _build_workflows(self) -> dict[str, Any]:
        try:
            items = list_workflow_artifact_records(self.workspace)
            proposal_counts: Counter[str] = Counter()
            verification_counts: Counter[str] = Counter()
            status_counts: Counter[str] = Counter()
            managed_by_domain_pack_count = 0
            for item in items:
                status_counts[str(item.get("status") or "unknown")] += 1
                if item.get("status") == "available":
                    proposal_counts[str(item.get("proposal_status") or "unknown")] += 1
                    verification_counts[str(item.get("verification_status") or "unknown")] += 1
                if item.get("managed_by_domain_pack"):
                    managed_by_domain_pack_count += 1
            summary = summarize_workflow_artifacts(self.workspace)
            return {
                "stats": {
                    **summary,
                    "workflow_status_counts": dict(status_counts),
                    "workflow_verification_status_counts": dict(verification_counts),
                    "managed_by_domain_pack_count": managed_by_domain_pack_count,
                },
                "items": items,
            }
        except Exception:
            return {
                "stats": {
                    "workflow_artifacts_count": 0,
                    "workflow_artifact_status_counts": {},
                    "invalid_workflow_artifacts_count": 0,
                    "workflow_status_counts": {},
                    "workflow_verification_status_counts": {},
                    "managed_by_domain_pack_count": 0,
                },
                "items": [],
            }

    def _build_facts(self) -> dict[str, Any]:
        try:
            return summarize_facts(self.workspace, fact_store=self._facts())
        except Exception:
            return {
                "active_count": 0,
                "pending_confirmation_count": 0,
                "category_counts": {},
                "domain_counts": {},
            }

    def _build_memory(self) -> dict[str, Any]:
        try:
            store = self._memory()
            content = store.read_memory()
            has_memory_context = bool(content.strip()) and not _is_template_content(content, "memory/MEMORY.md")
            pending_history = store.read_unprocessed_history(since_cursor=store.get_last_dream_cursor())
            return {
                "has_memory_context": has_memory_context,
                "recent_history_pending_count": len(pending_history),
            }
        except Exception:
            return {
                "has_memory_context": False,
                "recent_history_pending_count": 0,
            }

    def _build_reviews(self) -> tuple[dict[str, Any], list[dict[str, Any]]]:
        try:
            records = self._reviews().iter_all()
        except Exception:
            return (
                {
                    "pending_count": 0,
                    "status_counts": {},
                    "type_counts": {},
                    "origin_counts": {},
                },
                [],
            )
        status_counts: Counter[str] = Counter()
        type_counts: Counter[str] = Counter()
        origin_counts: Counter[str] = Counter()
        pending: list[dict[str, Any]] = []
        for record in records:
            status = str(record.get("status") or "pending")
            proposal_type = str(record.get("proposal_type") or record.get("type") or "unknown")
            origin = str(record.get("origin") or "background_review")
            status_counts[status] += 1
            type_counts[proposal_type] += 1
            origin_counts[origin] += 1
            if status == "pending":
                pending.append(record)
        pending.sort(key=lambda item: (str(item.get("created_at") or ""), str(item.get("id") or "")))
        return (
            {
                "pending_count": len(pending),
                "status_counts": dict(status_counts),
                "type_counts": dict(type_counts),
                "origin_counts": dict(origin_counts),
            },
            pending,
        )

    def _build_confirmations(self) -> tuple[dict[str, Any], list[ConfirmationRequest]]:
        try:
            confirmations = self._confirmations().read_all()
        except Exception:
            return (
                {
                    "pending_count": 0,
                    "expired_count": 0,
                    "kind_counts": {},
                    "risk_counts": {},
                },
                [],
            )
        kind_counts: Counter[str] = Counter()
        risk_counts: Counter[str] = Counter()
        pending: list[ConfirmationRequest] = []
        expired_count = 0
        for confirmation in confirmations:
            kind_counts[confirmation.kind] += 1
            risk_counts[confirmation.risk or "unknown"] += 1
            if confirmation.status in {"pending", "notified"}:
                pending.append(confirmation)
            elif confirmation.status == "expired":
                expired_count += 1
        pending.sort(key=lambda item: (item.created_at, item.confirmation_id))
        return (
            {
                "pending_count": len(pending),
                "expired_count": expired_count,
                "kind_counts": dict(kind_counts),
                "risk_counts": dict(risk_counts),
            },
            pending,
        )

    def _build_limitations(
        self,
        *,
        domains: list[dict[str, Any]],
        skills: list[dict[str, Any]],
        workflows: list[dict[str, Any]],
        pending_reviews: list[dict[str, Any]],
        pending_confirmations: list[ConfirmationRequest],
    ) -> list[dict[str, str]]:
        limits: list[dict[str, str]] = []
        limits.extend(self._domain_limitations(domains))
        limits.extend(self._skill_limitations(skills))
        limits.extend(self._workflow_limitations(workflows))
        limits.extend(self._review_limitations(pending_reviews))
        limits.extend(self._confirmation_limitations(pending_confirmations))
        return limits

    def _domain_limitations(self, domains: list[dict[str, Any]]) -> list[dict[str, str]]:
        limits: list[dict[str, str]] = []
        status_by_id = {str(item.get("id") or ""): str(item.get("status") or "") for item in domains}
        for item in sorted(domains, key=lambda row: str(row.get("id") or "")):
            pack_id = str(item.get("id") or "")
            status = str(item.get("status") or "")
            reason = _clean_summary(
                str(item.get("unavailable_reason") or item.get("validation_summary") or "Domain pack is unavailable."),
            )
            if status == "invalid":
                limits.append(_limitation("domain_invalid", "error", "domain", pack_id, reason))
            elif status == "unavailable":
                limits.append(_limitation("domain_unavailable", "warning", "domain", pack_id, reason))
            dependencies = (
                item.get("dependencies", {}).get("packs", [])
                if isinstance(item.get("dependencies"), dict)
                else []
            )
            for dependency in sorted(str(value) for value in dependencies if str(value).strip()):
                if status_by_id.get(dependency) == "available":
                    continue
                limits.append(
                    _limitation(
                        "domain_missing_dependency",
                        "warning",
                        "domain",
                        pack_id,
                        f"Depends on domain pack `{dependency}`, which is missing or unavailable.",
                    )
                )
        return limits

    def _skill_limitations(self, skills: list[dict[str, Any]]) -> list[dict[str, str]]:
        limits: list[dict[str, str]] = []
        for item in sorted(skills, key=lambda row: (str(row.get("source") or ""), str(row.get("name") or ""))):
            name = str(item.get("name") or "")
            lifecycle = str(item.get("lifecycle_status") or "unknown")
            verification = str(item.get("verification_status") or "unknown")
            description = _clean_summary(str(item.get("description") or name or "Skill"))
            if verification == "unverified":
                limits.append(
                    _limitation(
                        "skill_unverified",
                        "warning",
                        "skill",
                        name,
                        f"{description} is still unverified.",
                    )
                )
            if lifecycle == "deprecated":
                limits.append(
                    _limitation(
                        "skill_deprecated",
                        "warning",
                        "skill",
                        name,
                        f"{description} is deprecated and not part of the default capability set.",
                    )
                )
            if lifecycle == "rejected":
                limits.append(
                    _limitation(
                        "skill_rejected",
                        "error",
                        "skill",
                        name,
                        f"{description} was rejected and cannot be loaded.",
                    )
                )
        return limits

    def _workflow_limitations(self, workflows: list[dict[str, Any]]) -> list[dict[str, str]]:
        limits: list[dict[str, str]] = []
        for item in sorted(workflows, key=lambda row: str(row.get("name") or "")):
            if str(item.get("status") or "") == "available":
                continue
            limits.append(
                _limitation(
                    "workflow_invalid",
                    "error",
                    "workflow",
                    str(item.get("name") or ""),
                    str(item.get("unavailable_reason") or "Workflow artifact is invalid."),
                )
            )
        return limits

    def _review_limitations(self, pending_reviews: list[dict[str, Any]]) -> list[dict[str, str]]:
        limits: list[dict[str, str]] = []
        for record in pending_reviews:
            payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
            subject_type = str(payload.get("subject_type") or "").strip() or "review"
            subject_id = str(payload.get("subject_id") or "").strip() or str(record.get("id") or "")
            summary = str(record.get("subject_label") or record.get("title") or record.get("content") or "Pending review.")
            limits.append(_limitation("review_pending", "info", subject_type, subject_id, summary))
        return limits

    def _confirmation_limitations(
        self,
        pending_confirmations: list[ConfirmationRequest],
    ) -> list[dict[str, str]]:
        limits: list[dict[str, str]] = []
        for confirmation in pending_confirmations:
            summary = confirmation.prompt or f"Pending {confirmation.kind} confirmation."
            limits.append(
                _limitation(
                    "confirmation_pending",
                    "warning",
                    "confirmation",
                    confirmation.confirmation_id,
                    summary,
                )
            )
        return limits

    def _skills(self) -> SkillsLoader:
        if self._skills_loader is not None:
            return self._skills_loader
        return SkillsLoader(
            self.workspace,
            domain_pack_manager=self._domain_pack_manager,
        )

    def _memory(self) -> MemoryStore:
        if self._memory_store is not None:
            return self._memory_store
        return MemoryStore(self.workspace)

    def _facts(self) -> FactStore:
        if self._facts_store is not None:
            return self._facts_store
        return self._memory().fact_store

    def _confirmations(self) -> PendingConfirmationStore:
        if self._confirmation_store is not None:
            return self._confirmation_store
        return PendingConfirmationStore(self.workspace)

    def _reviews(self):
        if self._review_store is not None:
            return self._review_store
        if self._background_review_service is not None and hasattr(self._background_review_service, "store"):
            return self._background_review_service.store
        from OriginAgent.agent.background_review import ReviewProposalStore

        return ReviewProposalStore(self.workspace)

    def _domain_governance(self) -> DomainPackGovernanceService:
        if self._domain_governance_service is not None:
            return self._domain_governance_service
        return DomainPackGovernanceService(
            self.workspace,
            domain_pack_manager=self._domain_pack_manager,
        )

    @staticmethod
    def _review_enabled(service: Any | None, explicit: bool | None) -> bool:
        if explicit is not None:
            return explicit
        if service is None:
            return False
        with suppress(Exception):
            status = service.runtime_status()
            if isinstance(status, dict):
                for key in ("background_review_enabled", "curator_enabled"):
                    if isinstance(status.get(key), bool):
                        return bool(status[key])
        return bool(getattr(service, "enabled", False))


class SelfModelRenderer:
    """Render human-readable self-model summaries."""

    def render(self, self_model: dict[str, Any]) -> str:
        identity = self_model.get("identity", {})
        runtime = self_model.get("runtime", {})
        domains = self_model.get("domains", {})
        skills = self_model.get("skills", {})
        workflows = self_model.get("workflows", {})
        facts = self_model.get("facts", {})
        memory = self_model.get("memory", {})
        reviews = self_model.get("reviews", {})
        confirmations = self_model.get("confirmations", {})
        limitations = self_model.get("limitations", [])

        active_domains = [
            str(item.get("id") or "")
            for item in domains.get("items", [])
            if item.get("active") and str(item.get("status") or "") == "available"
        ]
        verified_skills = [
            str(item.get("name") or "")
            for item in skills.get("items", [])
            if str(item.get("lifecycle_status") or "") == "active"
            and str(item.get("verification_status") or "") == "verified"
        ]
        available_workflows = [
            str(item.get("name") or "")
            for item in workflows.get("items", [])
            if str(item.get("status") or "") == "available"
        ]
        invalid_workflow_count = int(
            workflows.get("stats", {}).get("invalid_workflow_artifacts_count", 0) or 0
        )

        lines = [
            "# Self Model",
            "",
            f"- Agent: `{identity.get('agent_name') or 'OriginAgent'}`",
            f"- Workspace: `{identity.get('workspace_name') or 'workspace'}`",
            f"- Runtime profile: `{identity.get('runtime_profile') or 'default'}`",
            f"- Audit mode: `{identity.get('audit_mode') or 'minimal'}`",
            "",
            "## Runtime",
            "",
            f"- Registered tools: {int(runtime.get('registered_tools_count', 0) or 0)}",
            f"- Active sessions: {int(runtime.get('active_sessions_count', 0) or 0)}",
            f"- Pending queues: {int(runtime.get('pending_queue_count', 0) or 0)}",
            f"- Cron available: {_yes_no(bool(runtime.get('cron_available')))}",
            f"- Confirmation available: {_yes_no(bool(runtime.get('confirmation_available')))}",
            f"- Background review enabled: {_yes_no(bool(runtime.get('background_review_enabled')))}",
            f"- Curator enabled: {_yes_no(bool(runtime.get('curator_enabled')))}",
            "",
            "## Available Capabilities",
            "",
            f"- Active domains: {_render_name_list(active_domains)}",
            f"- Verified skills: {_render_name_list(verified_skills)}",
            (
                f"- Workflow knowledge: {len(available_workflows)} available"
                + (f", {invalid_workflow_count} invalid" if invalid_workflow_count else "")
            ),
            f"- Facts: {int(facts.get('active_count', 0) or 0)} active, "
            f"{int(facts.get('pending_confirmation_count', 0) or 0)} pending confirmation",
            "",
            "## Reviews & Memory",
            "",
            f"- Pending reviews: {int(reviews.get('pending_count', 0) or 0)}",
            f"- Pending confirmations: {int(confirmations.get('pending_count', 0) or 0)}",
            f"- Recent history pending: {int(memory.get('recent_history_pending_count', 0) or 0)}",
            f"- Memory context available: {_yes_no(bool(memory.get('has_memory_context')))}",
            "",
            "## Known Limitations",
            "",
        ]
        if isinstance(limitations, list) and limitations:
            for item in limitations:
                lines.append(
                    f"- [{item.get('status') or 'info'}] {item.get('subject_type') or 'item'} "
                    f"`{item.get('subject_id') or ''}`: {item.get('summary') or ''}"
                )
        else:
            lines.append("- None.")
        return "\n".join(lines)


def _limitation(
    code: str,
    status: str,
    subject_type: str,
    subject_id: str,
    summary: str,
) -> dict[str, str]:
    return {
        "code": code,
        "status": status,
        "subject_type": subject_type,
        "subject_id": subject_id,
        "summary": _clean_summary(summary),
    }


def _clean_summary(text: str) -> str:
    compact = " ".join(str(text or "").split())
    return truncate_text(redact_memory_text(compact), _LIMITATION_MAX_CHARS)


def _render_name_list(names: list[str]) -> str:
    values = [name for name in names if name]
    if not values:
        return "none"
    shown = values[:_RENDER_LIST_MAX_ITEMS]
    suffix = f", +{len(values) - len(shown)} more" if len(values) > len(shown) else ""
    return ", ".join(f"`{name}`" for name in shown) + suffix


def _yes_no(value: bool) -> str:
    return "yes" if value else "no"


def _is_template_content(content: str, template_path: str) -> bool:
    with suppress(Exception):
        template = pkg_files("OriginAgent") / "templates" / template_path
        if template.is_file():
            return content.strip() == template.read_text(encoding="utf-8").strip()
    return False


def _safe_len(value: Any) -> int:
    try:
        return len(value)
    except Exception:
        return 0


def _session_count(sessions: Any) -> int:
    for attr in ("sessions", "_sessions"):
        value = getattr(sessions, attr, None)
        if value is not None:
            return _safe_len(value)
    return 0
