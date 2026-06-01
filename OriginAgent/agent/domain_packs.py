"""Domain pack discovery, validation, and prompt summaries."""

from __future__ import annotations

import os
import re
import shutil
import importlib.util
import sys
from dataclasses import dataclass, field, replace
from pathlib import Path
from typing import Any, Literal

import yaml

BUILTIN_DOMAIN_PACKS_DIR = Path(__file__).parent.parent / "domain_packs"
_IGNORED_PACK_DIR_NAMES = {"__pycache__"}
_DOMAIN_ID_RE = re.compile(r"^[a-z0-9_-]+$")
_TOOL_ID_RE = re.compile(r"^[a-z0-9_]{1,64}$")
_MODULE_RE = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)*$")
_CLASS_RE = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*$")
_VERIFICATION_STATUSES = {"unknown", "unverified", "verified"}
_ALLOWED_EVAL_KINDS = frozenset({"manifest", "skill", "tool", "workflow"})
_ALLOWED_TOOL_PERMISSIONS = frozenset(
    {
        "read_files",
        "write_files",
        "exec",
        "send_cross_target",
        "create_cron",
        "spawn",
        "mcp:read",
    }
)

DomainPackStatus = Literal["available", "unavailable", "invalid"]
DomainDeclarationStatus = Literal["available", "skipped"]
DomainToolRuntimeStatus = Literal["registered", "skipped"]
DomainPackSourceKind = Literal["builtin", "local_copy"]


@dataclass(frozen=True)
class DomainPackRequires:
    bins: tuple[str, ...] = ()
    env: tuple[str, ...] = ()


@dataclass(frozen=True)
class DomainPackDependencies:
    packs: tuple[str, ...] = ()


@dataclass(frozen=True)
class DomainPackSourceInfo:
    kind: DomainPackSourceKind = "local_copy"
    installed_from: str = ""
    installed_at: str = ""


@dataclass(frozen=True)
class DomainSkillDeclaration:
    id: str
    virtual_id: str
    path: Path | None = None
    status: DomainDeclarationStatus = "available"
    unavailable_reason: str = ""

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainWorkflowDeclaration:
    id: str
    path: Path | None = None
    status: DomainDeclarationStatus = "available"
    unavailable_reason: str = ""

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainFileDeclaration:
    id: str
    path: Path | None = None
    status: DomainDeclarationStatus = "available"
    unavailable_reason: str = ""

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainEvalDeclaration:
    id: str
    kind: str
    target: str = ""
    status: DomainDeclarationStatus = "available"
    unavailable_reason: str = ""

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainToolDeclaration:
    id: str
    module: str
    class_name: str
    permissions: tuple[str, ...] = ()
    audit: Literal["minimal", "security"] = "minimal"
    module_path: Path | None = None
    status: DomainDeclarationStatus = "available"
    unavailable_reason: str = ""

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainToolRuntimeRecord:
    pack_id: str
    tool_id: str
    status: DomainToolRuntimeStatus
    reason: str = ""


@dataclass(frozen=True)
class DomainRuntimeDeclaration:
    module: str
    factory: str = "build_runtime_contribution"
    module_path: Path | None = None
    status: DomainDeclarationStatus = "available"
    unavailable_reason: str = ""

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainRuntimeBuildContext:
    pack: "DomainPack"
    workspace: Path
    config: Any
    overrides: dict[str, Any] = field(default_factory=dict)


@dataclass(frozen=True)
class DomainRuntimeContribution:
    tool_context: dict[str, Any] = field(default_factory=dict)
    safety_gates: tuple[Any, ...] = ()
    permission_resolvers: tuple[Any, ...] = ()
    context_fragments: tuple[str, ...] = ()


@dataclass(frozen=True)
class DomainPack:
    id: str
    name: str
    version: str
    path: Path
    source: str
    status: DomainPackStatus
    active: bool = False
    enabled: bool = True
    active_requested: bool = False
    description: str = ""
    capabilities: tuple[str, ...] = ()
    triggers: tuple[str, ...] = ()
    requires: DomainPackRequires = field(default_factory=DomainPackRequires)
    dependencies: DomainPackDependencies = field(default_factory=DomainPackDependencies)
    unavailable_reason: str = ""
    validation_summary: str = ""
    capabilities_path: Path | None = None
    capabilities_content: str = ""
    skills: tuple[DomainSkillDeclaration, ...] = ()
    workflows: tuple[DomainWorkflowDeclaration, ...] = ()
    policies: tuple[DomainFileDeclaration, ...] = ()
    schemas: tuple[DomainFileDeclaration, ...] = ()
    tools: tuple[DomainToolDeclaration, ...] = ()
    runtime: DomainRuntimeDeclaration | None = None
    evals: tuple[DomainEvalDeclaration, ...] = ()
    manifest: dict[str, Any] = field(default_factory=dict)
    verification_status: str = "unknown"
    source_info: DomainPackSourceInfo = field(default_factory=DomainPackSourceInfo)
    overrides_builtin: bool = False

    @property
    def available(self) -> bool:
        return self.status == "available"


@dataclass(frozen=True)
class DomainPackRuntimeConfig:
    enabled: bool = True
    disabled: tuple[str, ...] = ()
    active: tuple[str, ...] = ()
    max_capability_chars: int = 4000

    @classmethod
    def from_config(cls, config: Any | None) -> "DomainPackRuntimeConfig":
        if config is None:
            return cls()
        max_capability_chars = getattr(config, "max_capability_chars", 4000)
        if max_capability_chars is None:
            max_capability_chars = 4000
        return cls(
            enabled=bool(getattr(config, "enabled", True)),
            disabled=tuple(str(item) for item in getattr(config, "disabled", []) or []),
            active=tuple(str(item) for item in getattr(config, "active", []) or []),
            max_capability_chars=int(max_capability_chars),
        )


class DomainPackValidator:
    """Validate one local domain pack directory into a deterministic record."""

    def __init__(
        self,
        *,
        runtime_config: DomainPackRuntimeConfig | None = None,
        strict_declarations: bool = False,
    ) -> None:
        self.runtime_config = runtime_config or DomainPackRuntimeConfig()
        self.strict_declarations = strict_declarations

    def validate_pack(self, pack_dir: Path, *, source: str) -> DomainPack:
        pack_dir = Path(pack_dir)
        manifest_path = pack_dir / "domain_pack.yaml"
        if not manifest_path.exists():
            return self._invalid(pack_dir.name, pack_dir, source, "missing domain_pack.yaml")

        try:
            raw = yaml.safe_load(manifest_path.read_text(encoding="utf-8"))
        except (OSError, yaml.YAMLError) as exc:
            return self._invalid(pack_dir.name, pack_dir, source, f"invalid domain_pack.yaml: {exc}")

        if not isinstance(raw, dict):
            return self._invalid(pack_dir.name, pack_dir, source, "domain_pack.yaml must be a mapping")

        pack_id = str(raw.get("id") or "").strip()
        name = str(raw.get("name") or "").strip()
        version = str(raw.get("version") or "").strip()
        if not pack_id:
            return self._invalid(pack_dir.name, pack_dir, source, "missing required field: id", raw)
        if not _DOMAIN_ID_RE.fullmatch(pack_id):
            return self._invalid(pack_id, pack_dir, source, "id must match ^[a-z0-9_-]+$", raw)
        missing = [field_name for field_name, value in (("name", name), ("version", version)) if not value]
        if missing:
            return self._invalid(pack_id, pack_dir, source, "missing required field(s): " + ", ".join(missing), raw)

        errors: list[str] = []
        capabilities_path = pack_dir / "CAPABILITIES.md"
        capabilities_content = ""
        status: DomainPackStatus = "available"
        reason = ""

        if not capabilities_path.exists():
            status = "unavailable"
            reason = "missing CAPABILITIES.md"
        else:
            try:
                capabilities_content = capabilities_path.read_text(encoding="utf-8")
            except OSError as exc:
                status = "unavailable"
                reason = f"cannot read CAPABILITIES.md: {exc}"

        source_info = self._parse_source_info(raw.get("source"), source, errors)
        verification_status = self._parse_verification_status(raw.get("verification_status"), errors)
        requires = self._parse_requires(raw.get("requires"))
        dependencies = self._parse_dependencies(raw.get("dependencies"), errors)
        skills, skill_errors = self._parse_skills(raw.get("skills"), pack_dir, pack_id)
        workflows, workflow_errors = self._parse_workflows(raw.get("workflows"), pack_dir)
        policies, policy_errors = self._parse_file_section(raw.get("policies"), pack_dir, "policies")
        schemas, schema_errors = self._parse_file_section(raw.get("schemas"), pack_dir, "schemas")
        tools, tool_errors = self._parse_tools(raw.get("tools"), pack_dir)
        runtime, runtime_errors = self._parse_runtime(raw.get("runtime"), pack_dir)
        evals, eval_errors = self._parse_evals(
            raw.get("evals"),
            skills=skills,
            workflows=workflows,
            tools=tools,
            errors=errors,
        )
        errors.extend(skill_errors)
        errors.extend(workflow_errors)
        errors.extend(policy_errors)
        errors.extend(schema_errors)
        errors.extend(tool_errors)
        errors.extend(runtime_errors)
        errors.extend(eval_errors)

        if raw.get("enabled") is False:
            status = "unavailable"
            reason = "disabled by manifest"

        enabled = pack_id not in set(self.runtime_config.disabled)
        if not enabled:
            status = "unavailable"
            reason = "disabled by config"

        if status == "available":
            missing_requirements = self._missing_requirements(requires)
            if missing_requirements:
                status = "unavailable"
                reason = "missing requirement(s): " + ", ".join(missing_requirements)

        if errors and self.strict_declarations:
            status = "invalid"
            reason = errors[0]

        active_requested = pack_id in set(self.runtime_config.active)
        active = active_requested and status == "available"
        validation_summary = "Domain pack is valid." if not errors else "; ".join(errors[:5])
        return DomainPack(
            id=pack_id,
            name=name,
            version=version,
            path=pack_dir,
            source=source,
            status=status,
            active=active,
            enabled=enabled,
            active_requested=active_requested,
            description=str(raw.get("description") or "").strip(),
            capabilities=tuple(_string_list(raw.get("capabilities"))),
            triggers=tuple(_string_list(_activation_triggers(raw.get("activation")))),
            requires=requires,
            dependencies=dependencies,
            unavailable_reason=reason,
            validation_summary=validation_summary,
            capabilities_path=capabilities_path,
            capabilities_content=capabilities_content,
            skills=tuple(skills),
            workflows=tuple(workflows),
            policies=tuple(policies),
            schemas=tuple(schemas),
            tools=tuple(tools),
            runtime=runtime,
            evals=tuple(evals),
            manifest=raw,
            verification_status=verification_status,
            source_info=source_info,
        )

    @staticmethod
    def _parse_requires(raw: Any) -> DomainPackRequires:
        if not isinstance(raw, dict):
            return DomainPackRequires()
        return DomainPackRequires(
            bins=tuple(_string_list(raw.get("bins"))),
            env=tuple(_string_list(raw.get("env"))),
        )

    @staticmethod
    def _missing_requirements(requires: DomainPackRequires) -> list[str]:
        missing: list[str] = []
        missing.extend(f"CLI: {command}" for command in requires.bins if not shutil.which(command))
        missing.extend(f"ENV: {name}" for name in requires.env if not os.environ.get(name))
        return missing

    def _parse_dependencies(self, raw: Any, errors: list[str]) -> DomainPackDependencies:
        if raw is None:
            return DomainPackDependencies()
        if not isinstance(raw, dict):
            errors.append("dependencies must be a mapping")
            return DomainPackDependencies()
        packs = []
        for value in _string_list(raw.get("packs")):
            pack_id = value.strip()
            if not _DOMAIN_ID_RE.fullmatch(pack_id):
                errors.append(f"dependencies.packs contains invalid id `{pack_id}`")
                continue
            packs.append(pack_id)
        return DomainPackDependencies(packs=tuple(packs))

    def _parse_source_info(
        self,
        raw: Any,
        source: str,
        errors: list[str],
    ) -> DomainPackSourceInfo:
        default_kind: DomainPackSourceKind = "builtin" if source == "builtin" else "local_copy"
        if raw is None:
            return DomainPackSourceInfo(kind=default_kind)
        if not isinstance(raw, dict):
            errors.append("source must be a mapping")
            return DomainPackSourceInfo(kind=default_kind)
        kind = str(raw.get("kind") or default_kind).strip()
        if kind not in {"builtin", "local_copy"}:
            errors.append("source.kind must be builtin or local_copy")
            kind = default_kind
        return DomainPackSourceInfo(
            kind=kind,
            installed_from=str(raw.get("installed_from") or "").strip(),
            installed_at=str(raw.get("installed_at") or "").strip(),
        )

    def _parse_verification_status(self, raw: Any, errors: list[str]) -> str:
        if raw is None:
            return "unknown"
        value = str(raw).strip().lower()
        if value not in _VERIFICATION_STATUSES:
            errors.append("verification_status must be unknown, unverified, or verified")
            return "unknown"
        return value

    def _parse_skills(
        self,
        raw: Any,
        pack_dir: Path,
        pack_id: str,
    ) -> tuple[list[DomainSkillDeclaration], list[str]]:
        declarations: list[DomainSkillDeclaration] = []
        errors: list[str] = []
        if raw is None:
            return declarations, errors
        if not isinstance(raw, list):
            return [
                DomainSkillDeclaration(
                    id="skills",
                    virtual_id="",
                    status="skipped",
                    unavailable_reason="skills must be a list",
                )
            ], ["skills must be a list"]
        for index, item in enumerate(raw):
            skill_id = ""
            if isinstance(item, str):
                skill_id = item.strip()
            elif isinstance(item, dict):
                skill_id = str(item.get("id") or item.get("name") or "").strip()
            label = skill_id or f"skill[{index}]"
            virtual_id = f"domain:{pack_id}/{skill_id}" if skill_id else ""
            if not skill_id or not _DOMAIN_ID_RE.fullmatch(skill_id):
                reason = "skill id must match ^[a-z0-9_-]+$"
                declarations.append(
                    DomainSkillDeclaration(
                        id=label,
                        virtual_id=virtual_id,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(reason if label.startswith("skill[") else f"{label}: {reason}")
                continue
            skill_path = pack_dir / "skills" / skill_id / "SKILL.md"
            if not skill_path.exists():
                reason = "missing SKILL.md"
                declarations.append(
                    DomainSkillDeclaration(
                        id=skill_id,
                        virtual_id=virtual_id,
                        path=skill_path,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(f"{skill_id}: {reason}")
                continue
            declarations.append(
                DomainSkillDeclaration(id=skill_id, virtual_id=virtual_id, path=skill_path)
            )
        return declarations, errors

    def _parse_workflows(
        self,
        raw: Any,
        pack_dir: Path,
    ) -> tuple[list[DomainWorkflowDeclaration], list[str]]:
        from OriginAgent.agent.workflow_artifacts import validate_workflow_artifact_dir

        declarations: list[DomainWorkflowDeclaration] = []
        errors: list[str] = []
        if raw is None:
            return declarations, errors
        if not isinstance(raw, list):
            return [
                DomainWorkflowDeclaration(
                    id="workflows",
                    status="skipped",
                    unavailable_reason="workflows must be a list",
                )
            ], ["workflows must be a list"]
        for index, item in enumerate(raw):
            workflow_id = ""
            if isinstance(item, str):
                workflow_id = item.strip()
            elif isinstance(item, dict):
                workflow_id = str(item.get("id") or item.get("name") or "").strip()
            label = workflow_id or f"workflow[{index}]"
            if not workflow_id or not _DOMAIN_ID_RE.fullmatch(workflow_id):
                reason = "workflow id must match ^[a-z0-9_-]+$"
                declarations.append(
                    DomainWorkflowDeclaration(
                        id=label,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(reason if label.startswith("workflow[") else f"{label}: {reason}")
                continue
            workflow_file = pack_dir / "workflows" / workflow_id / "workflow.yaml"
            workflow_dir = workflow_file.parent
            if not workflow_file.exists():
                reason = "missing workflow.yaml"
                declarations.append(
                    DomainWorkflowDeclaration(
                        id=workflow_id,
                        path=workflow_file,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(f"{workflow_id}: {reason}")
                continue
            valid, message = validate_workflow_artifact_dir(
                workflow_dir,
                workspace=pack_dir,
                expected_name=workflow_id,
            )
            if not valid:
                declarations.append(
                    DomainWorkflowDeclaration(
                        id=workflow_id,
                        path=workflow_file,
                        status="skipped",
                        unavailable_reason=message,
                    )
                )
                errors.append(f"{workflow_id}: {message}")
                continue
            declarations.append(DomainWorkflowDeclaration(id=workflow_id, path=workflow_file))
        return declarations, errors

    def _parse_file_section(
        self,
        raw: Any,
        pack_dir: Path,
        section: str,
    ) -> tuple[list[DomainFileDeclaration], list[str]]:
        declarations: list[DomainFileDeclaration] = []
        errors: list[str] = []
        if raw is None:
            return declarations, errors
        if not isinstance(raw, list):
            return [
                DomainFileDeclaration(
                    id=section,
                    status="skipped",
                    unavailable_reason=f"{section} must be a list",
                )
            ], [f"{section} must be a list"]
        for index, item in enumerate(raw):
            file_id = ""
            rel_path: str | None = None
            if isinstance(item, str):
                file_id = item.strip()
                rel_path = f"{section}/{file_id}"
            elif isinstance(item, dict):
                file_id = str(item.get("id") or item.get("name") or "").strip()
                rel_path = str(item.get("path") or f"{section}/{file_id}").strip()
            label = file_id or f"{section}[{index}]"
            if not file_id or not _DOMAIN_ID_RE.fullmatch(file_id):
                reason = f"{section} id must match ^[a-z0-9_-]+$"
                declarations.append(
                    DomainFileDeclaration(id=label, status="skipped", unavailable_reason=reason)
                )
                errors.append(reason if label.startswith(f"{section}[") else f"{label}: {reason}")
                continue
            path = _pack_relative_path(pack_dir, rel_path or "")
            if path is None:
                reason = f"{section} path must stay inside the pack root"
                declarations.append(
                    DomainFileDeclaration(
                        id=file_id,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(f"{file_id}: {reason}")
                continue
            if not path.exists():
                reason = f"missing declared {section[:-1] if section.endswith('s') else section} path"
                declarations.append(
                    DomainFileDeclaration(
                        id=file_id,
                        path=path,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(f"{file_id}: {reason}")
                continue
            declarations.append(DomainFileDeclaration(id=file_id, path=path))
        return declarations, errors

    def _parse_tools(
        self,
        raw: Any,
        pack_dir: Path,
    ) -> tuple[list[DomainToolDeclaration], list[str]]:
        declarations: list[DomainToolDeclaration] = []
        errors: list[str] = []
        if raw is None:
            return declarations, errors
        if not isinstance(raw, list):
            return [
                DomainToolDeclaration(
                    id="tools",
                    module="",
                    class_name="",
                    status="skipped",
                    unavailable_reason="tools must be a list",
                )
            ], ["tools must be a list"]
        for index, item in enumerate(raw):
            if not isinstance(item, dict):
                reason = "tool declaration must be a mapping"
                declarations.append(
                    DomainToolDeclaration(
                        id=f"tool[{index}]",
                        module="",
                        class_name="",
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                errors.append(reason)
                continue
            declaration = self._parse_tool(item, pack_dir, index)
            declarations.append(declaration)
            if declaration.status == "skipped":
                errors.append(f"{declaration.id}: {declaration.unavailable_reason}")
        return declarations, errors

    def _parse_tool(
        self,
        raw: dict[str, Any],
        pack_dir: Path,
        index: int,
    ) -> DomainToolDeclaration:
        tool_id = str(raw.get("id") or "").strip()
        module = str(raw.get("module") or "").strip()
        class_name = str(raw.get("class") or raw.get("class_name") or "").strip()
        label = tool_id or f"tool[{index}]"
        permissions_present = "permissions" in raw
        permissions = tuple(_string_list(raw.get("permissions")))
        audit = str(raw.get("audit") or "minimal").strip()
        module_path = _domain_tool_module_path(pack_dir, module)

        reason = ""
        if not tool_id:
            reason = "missing required field: id"
        elif not _TOOL_ID_RE.fullmatch(tool_id):
            reason = "tool id must match ^[a-z0-9_]{1,64}$"
        elif not module:
            reason = "missing required field: module"
        elif not _valid_domain_tool_module(module):
            reason = "module must be a dotted path under tools"
        elif module_path is None or not module_path.exists():
            reason = "missing tool module file"
        elif not class_name:
            reason = "missing required field: class"
        elif not _CLASS_RE.fullmatch(class_name):
            reason = "class must be a valid Python identifier"
        elif not permissions_present:
            reason = "missing permissions"
        elif any(not _is_allowed_tool_permission(permission) for permission in permissions):
            reason = "unsupported permission(s): " + ", ".join(
                permission for permission in permissions if not _is_allowed_tool_permission(permission)
            )
        elif audit not in {"minimal", "security"}:
            reason = "audit must be minimal or security"

        if reason:
            return DomainToolDeclaration(
                id=label,
                module=module,
                class_name=class_name,
                permissions=permissions,
                audit="security" if audit == "security" else "minimal",
                module_path=module_path,
                status="skipped",
                unavailable_reason=reason,
            )
        return DomainToolDeclaration(
            id=tool_id,
            module=module,
            class_name=class_name,
            permissions=permissions,
            audit="security" if audit == "security" else "minimal",
            module_path=module_path,
        )

    def _parse_runtime(
        self,
        raw: Any,
        pack_dir: Path,
    ) -> tuple[DomainRuntimeDeclaration | None, list[str]]:
        if raw is None:
            return None, []
        if not isinstance(raw, dict):
            return (
                DomainRuntimeDeclaration(
                    module="",
                    status="skipped",
                    unavailable_reason="runtime must be a mapping",
                ),
                ["runtime must be a mapping"],
            )
        module = str(raw.get("module") or "").strip()
        factory = str(raw.get("factory") or "build_runtime_contribution").strip()
        module_path = _domain_runtime_module_path(pack_dir, module)
        reason = ""
        if not module:
            reason = "missing required field: module"
        elif not _valid_domain_runtime_module(module):
            reason = "module must be a dotted path under runtime"
        elif module_path is None or not module_path.exists():
            reason = "missing runtime module file"
        elif not _CLASS_RE.fullmatch(factory):
            reason = "factory must be a valid Python identifier"
        if reason:
            return (
                DomainRuntimeDeclaration(
                    module=module,
                    factory=factory,
                    module_path=module_path,
                    status="skipped",
                    unavailable_reason=reason,
                ),
                [f"runtime: {reason}"],
            )
        return (
            DomainRuntimeDeclaration(
                module=module,
                factory=factory,
                module_path=module_path,
            ),
            [],
        )

    def _parse_evals(
        self,
        raw: Any,
        *,
        skills: list[DomainSkillDeclaration],
        workflows: list[DomainWorkflowDeclaration],
        tools: list[DomainToolDeclaration],
        errors: list[str],
    ) -> tuple[list[DomainEvalDeclaration], list[str]]:
        declarations: list[DomainEvalDeclaration] = []
        local_errors: list[str] = []
        if raw is None:
            return declarations, local_errors
        if not isinstance(raw, list):
            return [
                DomainEvalDeclaration(
                    id="evals",
                    kind="",
                    status="skipped",
                    unavailable_reason="evals must be a list",
                )
            ], ["evals must be a list"]
        known_targets = {
            "skill": {item.id for item in skills},
            "workflow": {item.id for item in workflows},
            "tool": {item.id for item in tools},
        }
        for index, item in enumerate(raw):
            if not isinstance(item, dict):
                declarations.append(
                    DomainEvalDeclaration(
                        id=f"eval[{index}]",
                        kind="",
                        status="skipped",
                        unavailable_reason="eval declaration must be a mapping",
                    )
                )
                local_errors.append("eval declaration must be a mapping")
                continue
            eval_id = str(item.get("id") or "").strip()
            kind = str(item.get("kind") or "").strip().lower()
            target = str(item.get("target") or "").strip()
            label = eval_id or f"eval[{index}]"
            reason = ""
            if not eval_id or not _DOMAIN_ID_RE.fullmatch(eval_id):
                reason = "eval id must match ^[a-z0-9_-]+$"
            elif kind not in _ALLOWED_EVAL_KINDS:
                reason = "eval kind must be one of manifest, skill, tool, workflow"
            elif target and kind in known_targets and target not in known_targets[kind]:
                reason = f"eval target `{target}` was not declared under {kind}s"
            if reason:
                declarations.append(
                    DomainEvalDeclaration(
                        id=label,
                        kind=kind,
                        target=target,
                        status="skipped",
                        unavailable_reason=reason,
                    )
                )
                local_errors.append(f"{label}: {reason}")
                continue
            declarations.append(DomainEvalDeclaration(id=eval_id, kind=kind, target=target))
        return declarations, local_errors

    @staticmethod
    def _invalid(
        pack_id: str,
        pack_dir: Path,
        source: str,
        reason: str,
        manifest: dict[str, Any] | None = None,
    ) -> DomainPack:
        safe_id = pack_id if pack_id and _DOMAIN_ID_RE.fullmatch(pack_id) else pack_dir.name
        return DomainPack(
            id=safe_id,
            name=safe_id,
            version="",
            path=pack_dir,
            source=source,
            status="invalid",
            unavailable_reason=reason,
            validation_summary=reason,
            manifest=manifest or {},
            source_info=DomainPackSourceInfo(kind="builtin" if source == "builtin" else "local_copy"),
        )


class DomainPackManager:
    """Discover local domain packs and render their agent-facing capability state."""

    def __init__(
        self,
        workspace: Path,
        *,
        config: Any | None = None,
        builtin_dir: Path | None = None,
    ) -> None:
        self.workspace = Path(workspace)
        self.workspace_dir = self.workspace / "domain_packs"
        self.builtin_dir = builtin_dir or BUILTIN_DOMAIN_PACKS_DIR
        self.config = DomainPackRuntimeConfig.from_config(config)
        self._packs: dict[str, DomainPack] | None = None
        self._domain_tool_runtime: dict[tuple[str, str], DomainToolRuntimeRecord] = {}
        self._validator = DomainPackValidator(runtime_config=self.config, strict_declarations=False)

    def refresh(self) -> None:
        self._validator = DomainPackValidator(runtime_config=self.config, strict_declarations=False)
        self._packs = None

    def list_packs(self) -> list[DomainPack]:
        if not self.config.enabled:
            return []
        if self._packs is None:
            self._packs = self._discover()
        return sorted(self._packs.values(), key=lambda pack: (pack.id, pack.source))

    def get_pack(self, pack_id: str) -> DomainPack | None:
        if not self.config.enabled:
            return None
        normalized = str(pack_id).strip()
        if not normalized:
            return None
        if self._packs is None:
            self._packs = self._discover()
        return self._packs.get(normalized)

    def build_summary(self) -> str:
        packs = self.list_packs()
        if not packs:
            return ""
        lines = ["# Domain Packs", ""]
        for pack in packs:
            marker = "active" if pack.active else pack.status
            reason = f"; {pack.unavailable_reason}" if pack.unavailable_reason else ""
            desc = f" — {pack.description}" if pack.description else ""
            lines.append(
                f"- `{pack.id}` [{pack.source}] "
                f"{pack.name} v{pack.version}: {marker}{reason}{desc}"
            )
            if pack.capabilities:
                lines.append(f"  Capabilities: {', '.join(pack.capabilities)}")
            if pack.active:
                skill_summary = _declaration_summary(pack.skills)
                tool_summary = _declaration_summary(pack.tools)
                if skill_summary:
                    lines.append(f"  Skills: {skill_summary}")
                if tool_summary:
                    lines.append(f"  Tools: {tool_summary}")
        return "\n".join(lines)

    def build_active_context(self) -> str:
        parts: list[str] = []
        limit = max(0, self.config.max_capability_chars)
        for pack in self.list_packs():
            if not pack.active:
                continue
            content = pack.capabilities_content.strip()
            if not content:
                continue
            if limit and len(content) > limit:
                content = content[:limit].rstrip() + "\n\n[Domain capabilities truncated]"
            parts.append(f"## Domain Pack: {pack.id}\n\n{content}")
        return "\n\n---\n\n".join(parts)

    def active_skill_entries(self) -> list[dict[str, str]]:
        """Return SkillsLoader-compatible entries for active domain pack skills."""
        entries: list[dict[str, str]] = []
        for pack in self.list_packs():
            if not pack.active:
                continue
            for skill in pack.skills:
                if not skill.available or skill.path is None:
                    continue
                entries.append(
                    {
                        "name": skill.virtual_id,
                        "path": str(skill.path),
                        "source": f"domain:{pack.id}",
                    }
                )
        return entries

    def get_active_skill_path(self, pack_id: str, skill_id: str) -> Path | None:
        pack = self.get_pack(pack_id)
        if pack is None or not pack.active:
            return None
        for skill in pack.skills:
            if skill.id == skill_id and skill.available:
                return skill.path
        return None

    def active_tool_declarations(self) -> list[tuple[DomainPack, DomainToolDeclaration]]:
        """Return tool declarations that active domain packs may try to register."""
        pairs: list[tuple[DomainPack, DomainToolDeclaration]] = []
        for pack in self.list_packs():
            if not pack.active:
                continue
            pairs.extend((pack, tool) for tool in pack.tools)
        return pairs

    def active_runtime_contributions(
        self,
        *,
        workspace: Path,
        config: Any,
        overrides: dict[str, Any] | None = None,
    ) -> list[DomainRuntimeContribution]:
        contributions: list[DomainRuntimeContribution] = []
        for pack in self.list_packs():
            if not pack.active or pack.runtime is None or not pack.runtime.available:
                continue
            contribution = self._load_runtime_contribution(
                pack,
                DomainRuntimeBuildContext(
                    pack=pack,
                    workspace=Path(workspace),
                    config=config,
                    overrides=dict(overrides or {}),
                ),
            )
            if contribution is not None:
                contributions.append(contribution)
        return contributions

    def clear_domain_tool_runtime(self) -> None:
        self._domain_tool_runtime.clear()

    def record_domain_tool_runtime(
        self,
        pack_id: str,
        tool_id: str,
        status: DomainToolRuntimeStatus,
        reason: str = "",
    ) -> None:
        self._domain_tool_runtime[(pack_id, tool_id)] = DomainToolRuntimeRecord(
            pack_id=pack_id,
            tool_id=tool_id,
            status=status,
            reason=reason,
        )

    def domain_tool_runtime_records(self, pack_id: str | None = None) -> list[DomainToolRuntimeRecord]:
        records = list(self._domain_tool_runtime.values())
        if pack_id is not None:
            records = [record for record in records if record.pack_id == pack_id]
        return sorted(records, key=lambda record: (record.pack_id, record.tool_id))

    def domain_tool_runtime_counts(self) -> dict[str, int]:
        counts = {"registered": 0, "skipped": 0}
        for record in self._domain_tool_runtime.values():
            counts[record.status] = counts.get(record.status, 0) + 1
        return counts

    def _discover(self) -> dict[str, DomainPack]:
        packs: dict[str, DomainPack] = {}
        builtin_packs: dict[str, DomainPack] = {}
        for pack_dir in self._pack_dirs(self.builtin_dir):
            pack = self._validator.validate_pack(pack_dir, source="builtin")
            builtin_packs[pack.id] = pack
            packs[pack.id] = pack
        for pack_dir in self._pack_dirs(self.workspace_dir):
            pack = self._validator.validate_pack(pack_dir, source="workspace")
            if pack.id in builtin_packs:
                pack = replace(pack, overrides_builtin=True)
            packs[pack.id] = pack
        return packs

    def _load_runtime_contribution(
        self,
        pack: DomainPack,
        context: DomainRuntimeBuildContext,
    ) -> DomainRuntimeContribution | None:
        declaration = pack.runtime
        if declaration is None or declaration.module_path is None:
            return None
        try:
            module_file = declaration.module_path.resolve()
            module_file.relative_to(pack.path.resolve())
            spec = importlib.util.spec_from_file_location(
                _runtime_module_name(pack.id, declaration),
                module_file,
            )
            if spec is None or spec.loader is None:
                return None
            module = importlib.util.module_from_spec(spec)
            sys.modules[spec.name] = module
            spec.loader.exec_module(module)
            factory = getattr(module, declaration.factory, None)
            if not callable(factory):
                return None
            raw = factory(context)
        except Exception:
            return None
        if raw is None:
            return None
        if isinstance(raw, DomainRuntimeContribution):
            return raw
        if isinstance(raw, dict):
            return DomainRuntimeContribution(**raw)
        return None

    @staticmethod
    def _pack_dirs(root: Path | None) -> list[Path]:
        if root is None or not root.exists() or not root.is_dir():
            return []
        return sorted(
            [
                path
                for path in root.iterdir()
                if path.is_dir() and path.name not in _IGNORED_PACK_DIR_NAMES
            ],
            key=lambda path: path.name,
        )


def _string_list(value: Any) -> list[str]:
    if value is None:
        return []
    if isinstance(value, (str, int, float, bool)):
        return [str(value)]
    if not isinstance(value, list):
        return []
    return [str(item) for item in value if isinstance(item, (str, int, float, bool))]


def _activation_triggers(value: Any) -> Any:
    if not isinstance(value, dict):
        return None
    return value.get("triggers")


def _declaration_summary(
    items: tuple[
        DomainSkillDeclaration
        | DomainWorkflowDeclaration
        | DomainFileDeclaration
        | DomainToolDeclaration
        | DomainEvalDeclaration,
        ...,
    ]
) -> str:
    if not items:
        return ""
    available = [item.id for item in items if item.status == "available"]
    skipped = [item.id for item in items if item.status == "skipped"]
    parts: list[str] = []
    if available:
        parts.append(", ".join(f"`{item}`" for item in available[:8]))
        if len(available) > 8:
            parts.append(f"+{len(available) - 8} more")
    if skipped:
        parts.append(f"{len(skipped)} skipped")
    return "; ".join(parts)


def _valid_domain_tool_module(module: str) -> bool:
    if not module.startswith("tools."):
        return False
    if any(part in module for part in ("/", "\\", "..")):
        return False
    return bool(_MODULE_RE.fullmatch(module))


def _valid_domain_runtime_module(module: str) -> bool:
    if not module.startswith("runtime."):
        return False
    if any(part in module for part in ("/", "\\", "..")):
        return False
    return bool(_MODULE_RE.fullmatch(module))


def _domain_tool_module_path(pack_dir: Path, module: str) -> Path | None:
    if not _valid_domain_tool_module(module):
        return None
    candidate = pack_dir / Path(*module.split(".")).with_suffix(".py")
    try:
        resolved_candidate = candidate.resolve()
        resolved_pack = pack_dir.resolve()
    except OSError:
        return None
    try:
        resolved_candidate.relative_to(resolved_pack)
    except ValueError:
        return None
    return candidate


def _domain_runtime_module_path(pack_dir: Path, module: str) -> Path | None:
    if not _valid_domain_runtime_module(module):
        return None
    candidate = pack_dir / Path(*module.split(".")).with_suffix(".py")
    try:
        resolved_candidate = candidate.resolve()
        resolved_pack = pack_dir.resolve()
    except OSError:
        return None
    try:
        resolved_candidate.relative_to(resolved_pack)
    except ValueError:
        return None
    return candidate


def _runtime_module_name(pack_id: str, declaration: DomainRuntimeDeclaration) -> str:
    safe_pack = pack_id.replace("-", "_")
    safe_runtime = declaration.module.replace(".", "_")
    return f"_originagent_domain_pack_{safe_pack}_{safe_runtime}"


def _pack_relative_path(pack_dir: Path, relative_path: str) -> Path | None:
    candidate = pack_dir / relative_path
    try:
        resolved_candidate = candidate.resolve()
        resolved_pack = pack_dir.resolve()
        resolved_candidate.relative_to(resolved_pack)
    except (OSError, ValueError):
        return None
    return candidate


def _is_allowed_tool_permission(permission: str) -> bool:
    if permission in _ALLOWED_TOOL_PERMISSIONS:
        return True
    return bool(re.fullmatch(r"device:[a-z0-9_-]+", permission))
