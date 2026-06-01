"""Controlled workflow artifact generation for reviewed proposals."""

from __future__ import annotations

import re
from collections import Counter
from dataclasses import dataclass
from pathlib import Path
from typing import Any

import yaml

from OriginAgent.agent.metadata import ORIGINAGENT_METADATA_KEY, originagent_metadata
from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.utils.helpers import truncate_text

SAFE_WORKFLOW_NAME_RE = re.compile(r"^[a-z0-9][a-z0-9-]{0,63}$")
CANONICAL_WORKFLOW_NAME_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")

_DESCRIPTION_MAX_CHARS = 512
_BODY_MAX_CHARS = 8000
_METADATA_MAX_CHARS = 512
_STEP_TITLE_MAX_CHARS = 160
_STEP_INSTRUCTION_MAX_CHARS = 1200
_STEP_RISK_MAX_CHARS = 80
_MAX_STEPS = 25
_ALLOWED_ROOT_FILES = {"workflow.yaml"}
_ALLOWED_TOP_LEVEL_KEYS = {
    "schema_version",
    "name",
    "description",
    "kind",
    "execution",
    "body",
    "steps",
    "metadata",
}
_ALLOWED_STEP_KEYS = {"title", "instruction", "risk", "confirmation_required"}
_EXECUTION_FLAGS = {
    "auto_run": False,
    "creates_cron": False,
    "calls_tools": False,
}
_ALLOWED_VERIFICATION_STATUSES = {"unverified", "verified"}
_ALLOWED_CREATED_BY = {"background_review", "auto_evolution"}
_UNSAFE_PHRASES = (
    "bypass permission",
    "bypass permissions",
    "bypass authorization",
    "skip permission",
    "skip confirmation",
    "bypass confirmation",
    "disable confirmation",
    "ignore confirmation",
    "without confirmation",
    "without user confirmation",
    "leak secret",
    "leak secrets",
    "exfiltrate secret",
    "exfiltrate secrets",
    "read secrets",
    "steal secret",
    "forge device state",
    "fake device state",
    "pretend action completed",
    "claim action completed without",
    "create cron job",
    "auto-run",
    "autorun",
    "绕过权限",
    "绕过确认",
    "不需要确认",
    "泄露密钥",
    "读取密钥",
    "伪造设备状态",
    "声称动作已完成",
)
_EXECUTABLE_KEY_RE = re.compile(
    r"(?im)^\s*(command|script|shell|cron|schedule|webhook|tool_call|tool_calls)\s*:"
)
_UNSAFE_TRUE_FLAG_RE = re.compile(r"(?im)^\s*(auto_run|calls_tools|creates_cron)\s*:\s*true\s*$")
_SECRET_PATTERNS = (
    re.compile(r"-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----", re.DOTALL),
    re.compile(r"(?i)\bBearer\s+[A-Za-z0-9._~+/=-]{8,}"),
    re.compile(r"\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b"),
    re.compile(r"\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{20,}\b"),
    re.compile(
        r"(?i)\b(api[_-]?key|token|secret|password)\b\s*[:=]\s*[\"']?[^\"'\s,;]{8,}[\"']?"
    ),
)


@dataclass(frozen=True)
class WorkflowArtifact:
    """Generated workspace workflow artifact metadata."""

    name: str
    relative_path: str
    content: str
    validation_message: str

    def to_json(self) -> dict[str, str]:
        return {
            "artifact_type": "workflow",
            "workflow_name": self.name,
            "path": self.relative_path,
            "validation": self.validation_message,
        }


def workflow_name_from_proposal(record: dict[str, Any]) -> str:
    """Return a safe workflow name from proposal payload or title."""

    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    raw = payload.get("workflow_name") if isinstance(payload, dict) else None
    if isinstance(raw, str) and raw.strip():
        candidate = raw.strip().casefold()
        redacted = redact_memory_text(candidate)
        if redacted != candidate:
            return _slug_from_text(redacted)
        return candidate
    return _slug_from_text(str(record.get("title") or ""))


def build_workflow_artifact(
    record: dict[str, Any],
    workspace: Path,
    *,
    metadata_overrides: dict[str, Any] | None = None,
) -> WorkflowArtifact:
    """Build and validate a workflow.yaml artifact without writing it."""

    workspace = Path(workspace)
    workflow_name = workflow_name_from_proposal(record)
    _validate_workflow_name(workflow_name)
    _validate_target_paths(workspace, workflow_name)

    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    description = _clean_field(
        payload.get("description") if isinstance(payload, dict) else None,
        _DESCRIPTION_MAX_CHARS,
    )
    if not description:
        description = _clean_field(record.get("content") or record.get("title"), _DESCRIPTION_MAX_CHARS)
    body = _clean_field(payload.get("body") if isinstance(payload, dict) else None, _BODY_MAX_CHARS)
    if not body:
        body = _fallback_body(record)
    steps = _clean_steps(payload.get("steps") if isinstance(payload, dict) else None)

    originagent = {
        "proposal_status": "proposed",
        "verification_status": "unverified",
        "review_proposal_id": _clean_metadata(record.get("id")),
        "domain_id": _clean_metadata(record.get("domain_id") or "core"),
        "created_by": "background_review",
        "source_session": _clean_metadata(record.get("session_key")),
        "source_turn_id": _clean_metadata(record.get("turn_id")),
    }
    if metadata_overrides:
        for key, value in metadata_overrides.items():
            if key in {
                "verification_status",
                "created_by",
                "previous_version",
                "verified_by",
                "verified_at",
                "opportunity_id",
            }:
                originagent[key] = None if key == "previous_version" and value is None else _clean_metadata(value)

    data = {
        "schema_version": 1,
        "name": workflow_name,
        "description": description,
        "kind": "manual_guide",
        "execution": dict(_EXECUTION_FLAGS),
        "body": body,
        "steps": steps,
        "metadata": originagent_metadata(originagent),
    }
    content = yaml.safe_dump(data, allow_unicode=True, sort_keys=False)
    validate_workflow_artifact_content(
        content,
        expected_name=workflow_name,
        expected_proposal_id=str(record.get("id") or ""),
        expected_domain_id=str(record.get("domain_id") or "core"),
    )
    return WorkflowArtifact(
        name=workflow_name,
        relative_path=f"workflows/{workflow_name}/workflow.yaml",
        content=content,
        validation_message="Workflow artifact is valid.",
    )


def write_workflow_artifact(
    record: dict[str, Any],
    workspace: Path,
    *,
    metadata_overrides: dict[str, Any] | None = None,
) -> WorkflowArtifact:
    """Generate, validate, and write one reviewed workspace workflow."""

    artifact = build_workflow_artifact(record, workspace, metadata_overrides=metadata_overrides)
    workspace = Path(workspace)
    target_dir = workspace / "workflows" / artifact.name
    target_file = target_dir / "workflow.yaml"
    if target_dir.exists():
        raise ValueError(f"workflow '{artifact.name}' already exists")
    target_dir.mkdir(parents=True, exist_ok=False)
    try:
        target_file.write_text(artifact.content, encoding="utf-8")
        valid, message = validate_workflow_artifact_dir(
            target_dir,
            workspace=workspace,
            expected_name=artifact.name,
        )
        if not valid:
            raise ValueError(message)
    except Exception:
        if target_file.exists():
            target_file.unlink()
        try:
            target_dir.rmdir()
        except OSError:
            pass
        raise
    return artifact


def validate_workflow_artifact_dir(
    workflow_dir: Path,
    *,
    workspace: Path,
    expected_name: str | None = None,
) -> tuple[bool, str]:
    """Validate a generated single-file workspace workflow directory."""

    workflow_dir = Path(workflow_dir)
    workspace = Path(workspace)
    try:
        resolved_workspace = workspace.resolve()
        resolved_dir = workflow_dir.resolve()
        resolved_dir.relative_to((resolved_workspace / "workflows").resolve())
    except (OSError, ValueError):
        return False, "Workflow directory must be inside workspace/workflows"
    if not resolved_dir.is_dir():
        return False, "Workflow directory not found"
    unexpected = sorted(
        child.name for child in resolved_dir.iterdir() if child.name not in _ALLOWED_ROOT_FILES
    )
    if unexpected:
        return False, "Generated review workflows may only contain workflow.yaml"
    workflow_file = resolved_dir / "workflow.yaml"
    if not workflow_file.is_file():
        return False, "workflow.yaml not found"
    try:
        content = workflow_file.read_text(encoding="utf-8")
    except OSError as exc:
        return False, f"Could not read workflow.yaml: {exc}"
    try:
        validate_workflow_artifact_content(
            content,
            expected_name=expected_name or resolved_dir.name,
        )
    except ValueError as exc:
        return False, str(exc)
    return True, "Workflow artifact is valid."


def validate_workflow_artifact_content(
    content: str,
    *,
    expected_name: str,
    expected_proposal_id: str | None = None,
    expected_domain_id: str | None = None,
) -> None:
    data = _load_workflow_yaml(content)
    unexpected = sorted(set(data) - _ALLOWED_TOP_LEVEL_KEYS)
    if unexpected:
        raise ValueError(f"workflow.yaml contains unsupported top-level keys: {unexpected}")
    if data.get("schema_version") != 1:
        raise ValueError("workflow schema_version must be 1")
    name = data.get("name")
    if name != expected_name:
        raise ValueError("workflow name must match directory name")
    _validate_workflow_name(str(name))
    description = data.get("description")
    if not isinstance(description, str) or not description.strip():
        raise ValueError("workflow description is required")
    if data.get("kind") != "manual_guide":
        raise ValueError("workflow kind must be manual_guide")
    execution = data.get("execution")
    if execution != _EXECUTION_FLAGS:
        raise ValueError("workflow execution must disable auto_run, creates_cron, and calls_tools")
    body = data.get("body")
    if not isinstance(body, str) or not body.strip():
        raise ValueError("workflow body is required")
    steps = data.get("steps")
    if not isinstance(steps, list):
        raise ValueError("workflow steps must be a list")
    for step in steps:
        _validate_step(step)
    metadata = data.get("metadata")
    if not isinstance(metadata, dict):
        raise ValueError("workflow metadata is required")
    originagent_meta = metadata.get(ORIGINAGENT_METADATA_KEY)
    if not isinstance(originagent_meta, dict):
        raise ValueError("metadata.OriginAgent is required")
    if originagent_meta.get("proposal_status") != "proposed":
        raise ValueError("metadata.OriginAgent.proposal_status must be proposed")
    verification_status = str(originagent_meta.get("verification_status") or "")
    if verification_status not in _ALLOWED_VERIFICATION_STATUSES:
        raise ValueError("metadata.OriginAgent.verification_status must be unverified or verified")
    created_by = str(originagent_meta.get("created_by") or "")
    if created_by not in _ALLOWED_CREATED_BY:
        raise ValueError("metadata.OriginAgent.created_by must be background_review or auto_evolution")
    if expected_proposal_id and originagent_meta.get("review_proposal_id") != expected_proposal_id:
        raise ValueError("metadata.OriginAgent.review_proposal_id is incorrect")
    if expected_domain_id and originagent_meta.get("domain_id") != expected_domain_id:
        raise ValueError("metadata.OriginAgent.domain_id is incorrect")
    _reject_executable_yaml(content)
    _reject_unredacted_secret(content)
    _reject_unsafe_workflow_text(content)


def summarize_workflow_artifacts(workspace: Path) -> dict[str, Any]:
    """Return a small redacted status summary for workspace workflows."""

    status_counts: Counter[str] = Counter()
    records = list_workflow_artifact_records(workspace)
    invalid = 0
    for record in records:
        if str(record.get("status") or "") != "available":
            invalid += 1
            continue
        status = str(record.get("proposal_status") or "unknown")
        status_counts[status] += 1
    return {
        "workflow_artifacts_count": len(records),
        "workflow_artifact_status_counts": dict(status_counts),
        "invalid_workflow_artifacts_count": invalid,
    }


def list_workflow_artifact_records(workspace: Path) -> list[dict[str, Any]]:
    """List workspace workflow artifacts with validation-derived status."""

    root = Path(workspace) / "workflows"
    try:
        children = sorted(root.iterdir(), key=lambda path: path.name)
    except FileNotFoundError:
        return []
    except OSError:
        return []

    records: list[dict[str, Any]] = []
    for child in children:
        if not child.is_dir():
            continue
        records.append(_workflow_artifact_record(child, Path(workspace)))
    return records


def _slug_from_text(text: str) -> str:
    redacted = redact_memory_text(text)
    slug = re.sub(r"[^a-z0-9]+", "-", redacted.casefold()).strip("-")
    slug = re.sub(r"-{2,}", "-", slug)
    return slug[:64].strip("-")


def _workflow_artifact_record(workflow_dir: Path, workspace: Path) -> dict[str, Any]:
    name = workflow_dir.name
    relative_path = f"workflows/{name}/workflow.yaml"
    valid, message = validate_workflow_artifact_dir(
        workflow_dir,
        workspace=workspace,
        expected_name=name,
    )
    if not valid:
        return {
            "name": name,
            "path": relative_path,
            "status": "invalid",
            "proposal_status": "unknown",
            "verification_status": "unknown",
            "domain_id": "",
            "managed_by_domain_pack": False,
            "unavailable_reason": _clean_metadata(message or "Workflow artifact is invalid."),
        }

    try:
        data = _load_workflow_yaml((workflow_dir / "workflow.yaml").read_text(encoding="utf-8"))
    except (OSError, ValueError) as exc:
        return {
            "name": name,
            "path": relative_path,
            "status": "invalid",
            "proposal_status": "unknown",
            "verification_status": "unknown",
            "domain_id": "",
            "managed_by_domain_pack": False,
            "unavailable_reason": _clean_metadata(str(exc) or "Workflow artifact is invalid."),
        }

    metadata = data.get("metadata", {}) if isinstance(data, dict) else {}
    originagent = metadata.get(ORIGINAGENT_METADATA_KEY, {}) if isinstance(metadata, dict) else {}
    managed_by_domain_pack = bool(originagent.get("managed_by_domain_pack")) or bool(
        originagent.get("migrated_from_workspace")
    )
    return {
        "name": str(data.get("name") or name),
        "path": relative_path,
        "status": "available",
        "proposal_status": str(originagent.get("proposal_status") or "unknown"),
        "verification_status": str(originagent.get("verification_status") or "unknown"),
        "domain_id": str(originagent.get("domain_id") or ""),
        "managed_by_domain_pack": managed_by_domain_pack,
        "unavailable_reason": "",
    }


def _validate_workflow_name(name: str) -> None:
    if not name or name in {".", ".."}:
        raise ValueError("workflow name is required")
    if "/" in name or "\\" in name or ".." in name or ":" in name:
        raise ValueError("workflow name must not contain path traversal or drive prefixes")
    if not SAFE_WORKFLOW_NAME_RE.fullmatch(name):
        raise ValueError("workflow name must match ^[a-z0-9][a-z0-9-]{0,63}$")
    if not CANONICAL_WORKFLOW_NAME_RE.fullmatch(name):
        raise ValueError("workflow name must use lowercase hyphen-case")


def _validate_target_paths(workspace: Path, workflow_name: str) -> None:
    target_file = Path(workspace) / "workflows" / workflow_name / "workflow.yaml"
    try:
        expected = (Path(workspace) / "workflows" / workflow_name / "workflow.yaml").resolve()
        if target_file.resolve() != expected:
            raise ValueError("workflow target path is invalid")
        expected.relative_to((Path(workspace) / "workflows").resolve())
    except (OSError, ValueError) as exc:
        raise ValueError("workflow target path must stay under workspace/workflows") from exc


def _clean_field(value: Any, max_chars: int) -> str:
    if value is None:
        return ""
    return truncate_text(redact_memory_text(str(value).strip()), max_chars)


def _clean_metadata(value: Any) -> str:
    return _clean_field(value, _METADATA_MAX_CHARS)


def _clean_steps(value: Any) -> list[dict[str, Any]]:
    if value is None:
        return []
    if not isinstance(value, list):
        raise ValueError("workflow steps must be a list")
    steps: list[dict[str, Any]] = []
    for item in value[:_MAX_STEPS]:
        if not isinstance(item, dict):
            raise ValueError("workflow steps must be mappings")
        unexpected = sorted(set(str(key) for key in item) - _ALLOWED_STEP_KEYS)
        if unexpected:
            raise ValueError(f"workflow steps contain unsupported keys: {unexpected}")
        confirmation_raw = item.get("confirmation_required", False)
        if not isinstance(confirmation_raw, bool):
            raise ValueError("workflow step confirmation_required must be a boolean")
        title = _clean_field(item.get("title"), _STEP_TITLE_MAX_CHARS)
        instruction = _clean_field(item.get("instruction"), _STEP_INSTRUCTION_MAX_CHARS)
        risk = _clean_field(item.get("risk") or "unknown", _STEP_RISK_MAX_CHARS)
        if title or instruction:
            steps.append(
                {
                    "title": title or "Step",
                    "instruction": instruction,
                    "risk": risk or "unknown",
                    "confirmation_required": confirmation_raw,
                }
            )
    return steps


def _fallback_body(record: dict[str, Any]) -> str:
    lines = [
        _clean_field(record.get("content"), _BODY_MAX_CHARS),
    ]
    rationale = _clean_field(record.get("rationale"), _BODY_MAX_CHARS)
    if rationale:
        lines.extend(["", "## Rationale", rationale])
    evidence = record.get("evidence")
    if isinstance(evidence, list):
        cleaned = [_clean_field(item, 500) for item in evidence[:5]]
        cleaned = [item for item in cleaned if item]
        if cleaned:
            lines.extend(["", "## Evidence"])
            lines.extend(f"- {item}" for item in cleaned)
    return "\n".join(line for line in lines if line is not None).strip()


def _load_workflow_yaml(content: str) -> dict[str, Any]:
    try:
        parsed = yaml.safe_load(content)
    except yaml.YAMLError as exc:
        raise ValueError(f"Invalid YAML in workflow.yaml: {exc}") from exc
    if not isinstance(parsed, dict):
        raise ValueError("workflow.yaml must be a YAML dictionary")
    return parsed


def _validate_step(step: Any) -> None:
    if not isinstance(step, dict):
        raise ValueError("workflow steps must be mappings")
    unexpected = sorted(set(step) - _ALLOWED_STEP_KEYS)
    if unexpected:
        raise ValueError(f"workflow steps contain unsupported keys: {unexpected}")
    title = step.get("title")
    instruction = step.get("instruction")
    risk = step.get("risk")
    confirmation_required = step.get("confirmation_required")
    if not isinstance(title, str) or not title.strip():
        raise ValueError("workflow step title is required")
    if not isinstance(instruction, str):
        raise ValueError("workflow step instruction is required")
    if not isinstance(risk, str) or not risk.strip():
        raise ValueError("workflow step risk is required")
    if not isinstance(confirmation_required, bool):
        raise ValueError("workflow step confirmation_required must be a boolean")


def _reject_executable_yaml(text: str) -> None:
    if _EXECUTABLE_KEY_RE.search(text):
        raise ValueError("workflow artifact contains executable or schedulable keys")
    if _UNSAFE_TRUE_FLAG_RE.search(text):
        raise ValueError("workflow artifact must not enable execution flags")


def _reject_unredacted_secret(text: str) -> None:
    redacted = redact_memory_text(text)
    redacted = redacted.replace("[REDACTED_SECRET]", "")
    redacted = redacted.replace("[REDACTED_BEARER_TOKEN]", "")
    redacted = redacted.replace("[REDACTED_PRIVATE_KEY]", "")
    for pattern in _SECRET_PATTERNS:
        if pattern.search(redacted):
            raise ValueError("workflow artifact still appears to contain an unredacted secret")


def _reject_unsafe_workflow_text(text: str) -> None:
    lowered = text.casefold()
    if "../" in text or "..\\" in text:
        raise ValueError("workflow artifact contains path traversal text")
    for phrase in _UNSAFE_PHRASES:
        if phrase.casefold() in lowered:
            raise ValueError("workflow artifact contains unsafe instructions")
