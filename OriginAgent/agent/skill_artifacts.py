"""Controlled skill artifact generation for reviewed proposals."""

from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path
from typing import Any

import yaml

from OriginAgent.agent.metadata import ORIGINAGENT_METADATA_KEY, originagent_metadata
from OriginAgent.agent.memory import redact_memory_text
from OriginAgent.agent.skills import BUILTIN_SKILLS_DIR
from OriginAgent.utils.helpers import truncate_text

SAFE_SKILL_NAME_RE = re.compile(r"^[a-z0-9][a-z0-9-]{0,63}$")
CANONICAL_SKILL_NAME_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")

_DESCRIPTION_MAX_CHARS = 512
_BODY_MAX_CHARS = 5000
_METADATA_MAX_CHARS = 512
_ALLOWED_ROOT_FILES = {"SKILL.md"}
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
    "绕过权限",
    "绕过确认",
    "不需要确认",
    "泄露密钥",
    "读取密钥",
    "伪造设备状态",
    "声称动作已完成",
)

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
class SkillArtifact:
    """Generated workspace skill artifact metadata."""

    name: str
    relative_path: str
    content: str
    validation_message: str

    def to_json(self) -> dict[str, str]:
        return {
            "skill_name": self.name,
            "path": self.relative_path,
            "validation": self.validation_message,
        }


def skill_name_from_proposal(record: dict[str, Any]) -> str:
    """Return a safe skill name from proposal payload or title."""

    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    raw = payload.get("skill_name") if isinstance(payload, dict) else None
    if isinstance(raw, str) and raw.strip():
        return raw.strip().casefold()
    title = str(record.get("title") or "")
    slug = re.sub(r"[^a-z0-9]+", "-", title.casefold()).strip("-")
    slug = re.sub(r"-{2,}", "-", slug)
    return slug[:64].strip("-")


def build_skill_artifact(record: dict[str, Any], workspace: Path) -> SkillArtifact:
    """Build and validate a SKILL.md artifact without writing it."""

    workspace = Path(workspace)
    skill_name = skill_name_from_proposal(record)
    _validate_skill_name(skill_name)
    _validate_target_paths(workspace, skill_name)

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
    body = _ensure_heading(body, skill_name)

    metadata = originagent_metadata(
        {
            "proposal_status": "proposed",
            "verification_status": "unverified",
            "lifecycle_status": "proposed",
            "review_proposal_id": _clean_metadata(record.get("id")),
            "domain_id": _clean_metadata(record.get("domain_id") or "core"),
            "created_by": _created_by_for_record(record),
        }
    )
    frontmatter = {
        "name": skill_name,
        "description": description,
        "always": False,
        "metadata": metadata,
    }
    content = "---\n" + yaml.safe_dump(frontmatter, allow_unicode=True, sort_keys=False) + "---\n\n" + body.rstrip() + "\n"
    validate_skill_artifact_content(
        content,
        expected_name=skill_name,
        expected_proposal_id=str(record.get("id") or ""),
        expected_domain_id=str(record.get("domain_id") or "core"),
    )
    return SkillArtifact(
        name=skill_name,
        relative_path=f"skills/{skill_name}/SKILL.md",
        content=content,
        validation_message="Skill artifact is valid.",
    )


def write_skill_artifact(record: dict[str, Any], workspace: Path) -> SkillArtifact:
    """Generate, validate, and write one reviewed workspace skill."""

    artifact = build_skill_artifact(record, workspace)
    workspace = Path(workspace)
    target_dir = workspace / "skills" / artifact.name
    target_file = target_dir / "SKILL.md"
    if target_dir.exists() or _builtin_skill_exists(artifact.name):
        raise ValueError(f"skill '{artifact.name}' already exists")
    target_dir.mkdir(parents=True, exist_ok=False)
    try:
        target_file.write_text(artifact.content, encoding="utf-8")
        valid, message = validate_skill_artifact_dir(
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


def validate_skill_artifact_dir(
    skill_dir: Path,
    *,
    workspace: Path,
    expected_name: str | None = None,
) -> tuple[bool, str]:
    """Validate a generated single-file workspace skill directory."""

    skill_dir = Path(skill_dir)
    workspace = Path(workspace)
    try:
        resolved_workspace = workspace.resolve()
        resolved_dir = skill_dir.resolve()
        resolved_dir.relative_to((resolved_workspace / "skills").resolve())
    except (OSError, ValueError):
        return False, "Skill directory must be inside workspace/skills"
    if not resolved_dir.is_dir():
        return False, "Skill directory not found"
    unexpected = sorted(child.name for child in resolved_dir.iterdir() if child.name not in _ALLOWED_ROOT_FILES)
    if unexpected:
        return False, "Generated review skills may only contain SKILL.md"
    skill_file = resolved_dir / "SKILL.md"
    if not skill_file.is_file():
        return False, "SKILL.md not found"
    try:
        content = skill_file.read_text(encoding="utf-8")
    except OSError as exc:
        return False, f"Could not read SKILL.md: {exc}"
    try:
        validate_skill_artifact_content(
            content,
            expected_name=expected_name or resolved_dir.name,
        )
    except ValueError as exc:
        return False, str(exc)
    return True, "Skill artifact is valid."


def validate_skill_artifact_content(
    content: str,
    *,
    expected_name: str,
    expected_proposal_id: str | None = None,
    expected_domain_id: str | None = None,
) -> None:
    frontmatter, body = _split_frontmatter(content)
    name = frontmatter.get("name")
    if name != expected_name:
        raise ValueError("skill frontmatter name must match directory name")
    _validate_skill_name(str(name))
    description = frontmatter.get("description")
    if not isinstance(description, str) or not description.strip():
        raise ValueError("skill description is required")
    if frontmatter.get("always") is not False:
        raise ValueError("review-applied skills must set always: false")
    metadata = frontmatter.get("metadata")
    if not isinstance(metadata, dict):
        raise ValueError("skill metadata is required")
    originagent_meta = metadata.get(ORIGINAGENT_METADATA_KEY)
    if not isinstance(originagent_meta, dict):
        raise ValueError("metadata.OriginAgent is required")
    required = {
        "proposal_status": "proposed",
    }
    for key, expected in required.items():
        if originagent_meta.get(key) != expected:
            raise ValueError(f"metadata.OriginAgent.{key} must be {expected}")
    verification_status = str(originagent_meta.get("verification_status") or "")
    if verification_status not in _ALLOWED_VERIFICATION_STATUSES:
        raise ValueError("metadata.OriginAgent.verification_status must be unverified or verified")
    lifecycle_status = str(originagent_meta.get("lifecycle_status") or "")
    if lifecycle_status and lifecycle_status != "proposed":
        raise ValueError("metadata.OriginAgent.lifecycle_status must be proposed")
    created_by = str(originagent_meta.get("created_by") or "")
    if created_by not in _ALLOWED_CREATED_BY:
        raise ValueError("metadata.OriginAgent.created_by must be background_review or auto_evolution")
    if expected_proposal_id and originagent_meta.get("review_proposal_id") != expected_proposal_id:
        raise ValueError("metadata.OriginAgent.review_proposal_id is incorrect")
    if expected_domain_id and originagent_meta.get("domain_id") != expected_domain_id:
        raise ValueError("metadata.OriginAgent.domain_id is incorrect")
    combined = yaml.safe_dump(frontmatter, allow_unicode=True, sort_keys=False) + "\n" + body
    _reject_unredacted_secret(combined)
    _reject_unsafe_skill_text(combined)


def _validate_skill_name(name: str) -> None:
    if not name or name in {".", ".."}:
        raise ValueError("skill name is required")
    if "/" in name or "\\" in name or ".." in name or ":" in name:
        raise ValueError("skill name must not contain path traversal or drive prefixes")
    if not SAFE_SKILL_NAME_RE.fullmatch(name):
        raise ValueError("skill name must match ^[a-z0-9][a-z0-9-]{0,63}$")
    if not CANONICAL_SKILL_NAME_RE.fullmatch(name):
        raise ValueError("skill name must use lowercase hyphen-case")


def _validate_target_paths(workspace: Path, skill_name: str) -> None:
    workspace = Path(workspace)
    target_file = workspace / "skills" / skill_name / "SKILL.md"
    try:
        expected = (workspace / "skills" / skill_name / "SKILL.md").resolve()
        if target_file.resolve() != expected:
            raise ValueError("skill target path is invalid")
        expected.relative_to((workspace / "skills").resolve())
    except (OSError, ValueError) as exc:
        raise ValueError("skill target path must stay under workspace/skills") from exc


def _builtin_skill_exists(skill_name: str) -> bool:
    return (BUILTIN_SKILLS_DIR / skill_name / "SKILL.md").exists()


def _clean_field(value: Any, max_chars: int) -> str:
    if value is None:
        return ""
    return truncate_text(redact_memory_text(str(value).strip()), max_chars)


def _clean_metadata(value: Any) -> str:
    return _clean_field(value, _METADATA_MAX_CHARS)


def _created_by_for_record(record: dict[str, Any]) -> str:
    origin = str(record.get("origin") or "").strip().lower()
    payload = record.get("payload") if isinstance(record.get("payload"), dict) else {}
    evolution = payload.get("evolution") if isinstance(payload.get("evolution"), dict) else {}
    evolution_origin = str(evolution.get("origin") or "").strip().lower()
    return "auto_evolution" if "auto_evolution" in {origin, evolution_origin} else "background_review"


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


def _ensure_heading(body: str, skill_name: str) -> str:
    body = body.strip()
    if not body:
        raise ValueError("skill body is required")
    if body.startswith("#"):
        return body
    title = skill_name.replace("-", " ").title()
    return f"# {title}\n\n{body}"


def _split_frontmatter(content: str) -> tuple[dict[str, Any], str]:
    lines = content.splitlines()
    if not lines or lines[0].strip() != "---":
        raise ValueError("Invalid frontmatter format")
    closing = None
    for index in range(1, len(lines)):
        if lines[index].strip() == "---":
            closing = index
            break
    if closing is None:
        raise ValueError("Invalid frontmatter format")
    raw = "\n".join(lines[1:closing])
    try:
        parsed = yaml.safe_load(raw)
    except yaml.YAMLError as exc:
        raise ValueError(f"Invalid YAML in frontmatter: {exc}") from exc
    if not isinstance(parsed, dict):
        raise ValueError("Frontmatter must be a YAML dictionary")
    return parsed, "\n".join(lines[closing + 1 :]).strip()


def _reject_unredacted_secret(text: str) -> None:
    redacted = redact_memory_text(text)
    redacted = redacted.replace("[REDACTED_SECRET]", "")
    redacted = redacted.replace("[REDACTED_BEARER_TOKEN]", "")
    redacted = redacted.replace("[REDACTED_PRIVATE_KEY]", "")
    for pattern in _SECRET_PATTERNS:
        if pattern.search(redacted):
            raise ValueError("skill artifact still appears to contain an unredacted secret")


def _reject_unsafe_skill_text(text: str) -> None:
    lowered = text.casefold()
    if "../" in text or "..\\" in text:
        raise ValueError("skill artifact contains path traversal text")
    for phrase in _UNSAFE_PHRASES:
        if phrase.casefold() in lowered:
            raise ValueError("skill artifact contains unsafe instructions")
