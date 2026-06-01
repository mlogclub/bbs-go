"""Schema contracts for governed evolution release-candidate stores."""

from __future__ import annotations

import json
from collections import Counter
from contextlib import suppress
from dataclasses import asdict, dataclass
from pathlib import Path
from typing import Any

from OriginAgent.agent.evolution_dependencies import DEPENDENCY_SCHEMA_VERSION, DEPENDENCY_STORE_RELATIVE
from OriginAgent.agent.evolution_health_history import (
    HEALTH_HISTORY_SCHEMA_VERSION,
    HEALTH_HISTORY_STORE_RELATIVE,
)
from OriginAgent.agent.evolution_outcomes import OUTCOME_SCHEMA_VERSION, OUTCOME_STORE_RELATIVE
from OriginAgent.agent.evolution_snapshots import SNAPSHOT_ROOT_RELATIVE, SNAPSHOT_SCHEMA_VERSION
from OriginAgent.agent.evolution_trial_logs import TRIAL_LOG_SCHEMA_VERSION, TRIAL_LOG_STORE_RELATIVE

OPPORTUNITY_SIGNAL_SCHEMA_VERSION = "originagent.evolution.opportunity_signal.v1"
REVIEW_PROPOSAL_EVOLUTION_SCHEMA_VERSION = "originagent.evolution.review_payload.v1"
OPPORTUNITY_SIGNAL_STORE_RELATIVE = Path("memory") / "opportunity_signals.jsonl"
REVIEW_PROPOSAL_STORE_RELATIVE = Path("memory") / "review_proposals.jsonl"

_SEVERITIES = {"info", "warning", "pending", "reject"}
_SIGNAL_STATUSES = {"open", "converted", "suppressed", "archived"}


@dataclass(frozen=True)
class EvolutionSchemaIssue:
    """One non-mutating schema validation finding."""

    store: str
    index: int
    code: str
    severity: str
    message: str

    def to_json(self) -> dict[str, Any]:
        return asdict(self)


def validate_evolution_stores(workspace: Path) -> dict[str, Any]:
    """Validate current governed evolution JSONL stores without rewriting them."""

    workspace = Path(workspace)
    issues: list[EvolutionSchemaIssue] = []
    counts: dict[str, int] = {}

    signal_records = _read_jsonl(workspace / OPPORTUNITY_SIGNAL_STORE_RELATIVE)
    counts["opportunity_signals"] = len(signal_records)
    for index, record in enumerate(signal_records):
        issues.extend(_validate_signal(record, index=index))

    review_records = _read_jsonl(workspace / REVIEW_PROPOSAL_STORE_RELATIVE)
    auto_evolution_review_records = [
        record for record in review_records
        if _is_auto_evolution_review_record(record)
    ]
    counts["review_proposals"] = len(auto_evolution_review_records)
    for index, record in enumerate(auto_evolution_review_records):
        issues.extend(_validate_review_proposal_record(record, index=index))

    outcome_records = _read_jsonl(workspace / OUTCOME_STORE_RELATIVE)
    counts["outcomes"] = len(outcome_records)
    for index, record in enumerate(outcome_records):
        issues.extend(_validate_outcome(record, index=index))

    dependency_records = _read_jsonl(workspace / DEPENDENCY_STORE_RELATIVE)
    counts["dependencies"] = len(dependency_records)
    for index, record in enumerate(dependency_records):
        issues.extend(_validate_dependency(record, index=index))

    trial_log_records = _read_jsonl(workspace / TRIAL_LOG_STORE_RELATIVE)
    counts["trial_logs"] = len(trial_log_records)
    for index, record in enumerate(trial_log_records):
        issues.extend(_validate_trial_log(record, index=index))

    health_history_records = _read_jsonl(workspace / HEALTH_HISTORY_STORE_RELATIVE)
    counts["health_history"] = len(health_history_records)
    for index, record in enumerate(health_history_records):
        issues.extend(_validate_health_history(record, index=index))

    snapshot_records = _read_snapshot_records(workspace / SNAPSHOT_ROOT_RELATIVE)
    counts["snapshots"] = len(snapshot_records)
    for index, record in enumerate(snapshot_records):
        issues.extend(_validate_snapshot(record, index=index))

    severity_counts = Counter(issue.severity for issue in issues)
    return {
        "ok": not any(issue.severity == "reject" for issue in issues),
        "schema_versions": {
            "opportunity_signal": OPPORTUNITY_SIGNAL_SCHEMA_VERSION,
            "review_proposal_payload": REVIEW_PROPOSAL_EVOLUTION_SCHEMA_VERSION,
            "outcome": OUTCOME_SCHEMA_VERSION,
            "dependency": DEPENDENCY_SCHEMA_VERSION,
            "trial_log": TRIAL_LOG_SCHEMA_VERSION,
            "health_history": HEALTH_HISTORY_SCHEMA_VERSION,
            "snapshot": SNAPSHOT_SCHEMA_VERSION,
        },
        "record_counts": counts,
        "issue_counts": dict(sorted(severity_counts.items())),
        "issues": [issue.to_json() for issue in issues],
    }


def validate_review_proposal_payload(payload: dict[str, Any], *, index: int = 0) -> list[dict[str, Any]]:
    """Validate the auto-evolution portion of a review proposal payload."""

    issues: list[EvolutionSchemaIssue] = []
    if not isinstance(payload, dict):
        return [
            _issue("review_proposal_payload", index, "payload_not_mapping", "reject", "Payload must be a mapping.")
            .to_json()
        ]
    evolution = _mapping(payload.get("evolution"))
    if not evolution:
        issues.append(_issue(
            "review_proposal_payload",
            index,
            "missing_evolution",
            "reject",
            "Auto-evolution review payload must include payload.evolution.",
        ))
    else:
        issues.extend(_required_string(evolution, "origin", "review_proposal_payload", index))
        issues.extend(_required_string(evolution, "opportunity_id", "review_proposal_payload", index))
        issues.extend(_required_string(evolution, "kind", "review_proposal_payload", index))
        issues.extend(_required_number(evolution, "priority_score", "review_proposal_payload", index))
        evidence = evolution.get("evidence_sources")
        if not isinstance(evidence, list):
            issues.append(_issue(
                "review_proposal_payload",
                index,
                "evidence_sources_not_list",
                "reject",
                "payload.evolution.evidence_sources must be a list.",
            ))
    if "static_gate" in payload:
        issues.extend(_validate_gate_shape(_mapping(payload.get("static_gate")), "static_gate", index))
    if "promotion_gate" in payload:
        gate = _mapping(payload.get("promotion_gate"))
        issues.extend(_required_string(gate, "decision", "promotion_gate", index))
        issues.extend(_required_string(gate, "suggested_action", "promotion_gate", index))
    if "operator_insights" in payload:
        insights = _mapping(payload.get("operator_insights"))
        issues.extend(_required_mapping(insights, "trial_summary", "operator_insights", index))
        issues.extend(_required_mapping(insights, "risk_summary", "operator_insights", index))
        issues.extend(_required_string(insights, "recommended_action", "operator_insights", index))
    return [issue.to_json() for issue in issues]


def _validate_review_proposal_record(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    for key in ("id", "created_at", "proposal_type", "origin", "title"):
        issues.extend(_required_string(record, key, "review_proposals", index))
    payload = record.get("payload")
    if not isinstance(payload, dict):
        issues.append(_issue(
            "review_proposals",
            index,
            "payload_not_mapping",
            "reject",
            "Auto-evolution review proposal payload must be a mapping.",
        ))
        return issues
    issues.extend(
        EvolutionSchemaIssue(**issue)
        for issue in validate_review_proposal_payload(payload, index=index)
    )
    return issues


def _validate_signal(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    for key in ("opportunity_id", "kind", "target_key", "title", "summary", "first_seen_at", "last_seen_at"):
        issues.extend(_required_string(record, key, "opportunity_signals", index))
    issues.extend(_required_int(record, "seen_count", "opportunity_signals", index))
    issues.extend(_required_number(record, "priority_score", "opportunity_signals", index))
    if not isinstance(record.get("evidence_sources"), list):
        issues.append(_issue(
            "opportunity_signals",
            index,
            "evidence_sources_not_list",
            "reject",
            "Opportunity signal evidence_sources must be a list.",
        ))
    status = str(record.get("status") or "")
    if status and status not in _SIGNAL_STATUSES:
        issues.append(_issue(
            "opportunity_signals",
            index,
            "invalid_status",
            "reject",
            f"Opportunity signal status `{status}` is not supported.",
        ))
    return issues


def _validate_outcome(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    issues.extend(_required_schema(record, OUTCOME_SCHEMA_VERSION, "outcomes", index))
    for key in ("event_id", "timestamp", "type"):
        issues.extend(_required_string(record, key, "outcomes", index))
    if not isinstance(record.get("metadata"), dict):
        issues.append(_issue("outcomes", index, "metadata_not_mapping", "reject", "Outcome metadata must be a mapping."))
    return issues


def _validate_dependency(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    issues.extend(_required_schema(record, DEPENDENCY_SCHEMA_VERSION, "dependencies", index))
    for key in ("artifact_type", "artifact_name", "artifact_path", "updated_at", "version"):
        issues.extend(_required_string(record, key, "dependencies", index))
    for key in ("depends_on", "referenced_by"):
        if not isinstance(record.get(key), list):
            issues.append(_issue("dependencies", index, f"{key}_not_list", "reject", f"Dependency {key} must be a list."))
    return issues


def _validate_trial_log(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    issues.extend(_required_schema(record, TRIAL_LOG_SCHEMA_VERSION, "trial_logs", index))
    for key in ("trial_id", "timestamp", "status"):
        issues.extend(_required_string(record, key, "trial_logs", index))
    if not isinstance(record.get("summary"), dict):
        issues.append(_issue("trial_logs", index, "summary_not_mapping", "reject", "Trial log summary must be a mapping."))
    if not isinstance(record.get("step_logs"), list):
        issues.append(_issue("trial_logs", index, "step_logs_not_list", "reject", "Trial log step_logs must be a list."))
    return issues


def _validate_health_history(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    issues.extend(_required_schema(record, HEALTH_HISTORY_SCHEMA_VERSION, "health_history", index))
    for key in ("snapshot_id", "timestamp", "level"):
        issues.extend(_required_string(record, key, "health_history", index))
    issues.extend(_required_int(record, "score", "health_history", index))
    if not isinstance(record.get("reasons"), list):
        issues.append(_issue("health_history", index, "reasons_not_list", "reject", "Health history reasons must be a list."))
    return issues


def _validate_snapshot(record: dict[str, Any], *, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    issues.extend(_required_schema(record, SNAPSHOT_SCHEMA_VERSION, "snapshots", index))
    for key in ("snapshot_id", "created_at", "artifact_type", "artifact_name", "artifact_path", "snapshot_path", "content_hash"):
        issues.extend(_required_string(record, key, "snapshots", index))
    return issues


def _validate_gate_shape(gate: dict[str, Any], store: str, index: int) -> list[EvolutionSchemaIssue]:
    issues: list[EvolutionSchemaIssue] = []
    issues.extend(_required_string(gate, "decision", store, index))
    if not isinstance(gate.get("issues"), list):
        issues.append(_issue(store, index, "issues_not_list", "reject", f"{store}.issues must be a list."))
        return issues
    for issue_index, item in enumerate(gate.get("issues") or []):
        if not isinstance(item, dict):
            issues.append(_issue(store, index, "issue_not_mapping", "reject", f"{store}.issues[{issue_index}] must be a mapping."))
            continue
        for key in ("code", "severity", "message"):
            issues.extend(_required_string(item, key, store, index))
        severity = str(item.get("severity") or "")
        if severity and severity not in _SEVERITIES:
            issues.append(_issue(store, index, "invalid_issue_severity", "warning", f"Unsupported issue severity `{severity}`."))
    return issues


def _required_schema(
    record: dict[str, Any],
    expected: str,
    store: str,
    index: int,
) -> list[EvolutionSchemaIssue]:
    actual = str(record.get("schema_version") or "")
    if actual == expected:
        return []
    return [_issue(store, index, "schema_version_mismatch", "reject", f"Expected schema_version `{expected}`, got `{actual}`.")]


def _required_string(record: dict[str, Any], key: str, store: str, index: int) -> list[EvolutionSchemaIssue]:
    if isinstance(record.get(key), str) and record.get(key):
        return []
    return [_issue(store, index, f"missing_{key}", "reject", f"`{key}` must be a non-empty string.")]


def _required_int(record: dict[str, Any], key: str, store: str, index: int) -> list[EvolutionSchemaIssue]:
    value = record.get(key)
    if isinstance(value, int) and not isinstance(value, bool):
        return []
    return [_issue(store, index, f"invalid_{key}", "reject", f"`{key}` must be an integer.")]


def _required_number(record: dict[str, Any], key: str, store: str, index: int) -> list[EvolutionSchemaIssue]:
    value = record.get(key)
    if isinstance(value, (int, float)) and not isinstance(value, bool):
        return []
    return [_issue(store, index, f"invalid_{key}", "reject", f"`{key}` must be numeric.")]


def _required_mapping(record: dict[str, Any], key: str, store: str, index: int) -> list[EvolutionSchemaIssue]:
    if isinstance(record.get(key), dict):
        return []
    return [_issue(store, index, f"{key}_not_mapping", "reject", f"`{key}` must be a mapping.")]


def _issue(store: str, index: int, code: str, severity: str, message: str) -> EvolutionSchemaIssue:
    return EvolutionSchemaIssue(
        store=store,
        index=index,
        code=code,
        severity=severity,
        message=message,
    )


def _mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def _is_auto_evolution_review_record(record: dict[str, Any]) -> bool:
    if str(record.get("origin") or "").strip().lower() == "auto_evolution":
        return True
    payload = record.get("payload")
    evolution = payload.get("evolution") if isinstance(payload, dict) else {}
    return (
        isinstance(evolution, dict)
        and str(evolution.get("origin") or "").strip().lower() == "auto_evolution"
    )


def _read_jsonl(path: Path) -> list[dict[str, Any]]:
    records: list[dict[str, Any]] = []
    with suppress(FileNotFoundError):
        with Path(path).open("r", encoding="utf-8") as handle:
            for line in handle:
                line = line.strip()
                if not line:
                    continue
                try:
                    raw = json.loads(line)
                except json.JSONDecodeError:
                    continue
                if isinstance(raw, dict):
                    records.append(raw)
    return records


def _read_snapshot_records(root: Path) -> list[dict[str, Any]]:
    records: list[dict[str, Any]] = []
    root = Path(root)
    if not root.exists():
        return records
    for metadata_path in root.glob("*/*/*/version_metadata.json"):
        try:
            raw = json.loads(metadata_path.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            continue
        if isinstance(raw, dict):
            records.append(raw)
    return records
