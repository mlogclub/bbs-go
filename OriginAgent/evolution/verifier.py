"""Static verification for staged evolution modules."""

from __future__ import annotations

import json
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Literal

from OriginAgent.agent.domain_packs import DomainPackRuntimeConfig, DomainPackValidator
from OriginAgent.evolution.package import compute_artifact_digest, read_package_manifest

VerificationStatus = Literal["verified", "failed"]

_STAGING_SCHEMA_VERSION = "originagent.evolution.staging.v1"
_PERMISSION_KEYS = frozenset(
    {
        "read_files",
        "write_files",
        "exec",
        "send_cross_target",
        "create_cron",
        "spawn",
        "device_domains",
        "mcp_scopes",
    }
)
_DENIED_BOOL_PERMISSIONS = (
    "write_files",
    "exec",
    "send_cross_target",
    "create_cron",
    "spawn",
)
_PERMISSION_CODE_BY_KEY = {
    "read_files": "permission_read_files",
    "write_files": "permission_write_files",
    "exec": "permission_exec",
    "send_cross_target": "permission_send_cross_target",
    "create_cron": "permission_create_cron",
    "spawn": "permission_spawn",
    "device_domains": "permission_device_domains",
    "mcp_scopes": "permission_mcp_scopes",
}


@dataclass(frozen=True)
class EvolutionVerificationReport:
    ok: bool
    status: VerificationStatus
    checks: tuple[dict[str, Any], ...]
    module_id: str = ""
    module_type: str = ""
    module_version: str = ""
    artifact_digest: str = ""
    staging_path: str = ""
    error: str = ""
    permissions_evaluated: dict[str, Any] | None = None
    permissions_denied: tuple[str, ...] = ()
    unknown_keys_rejected: tuple[str, ...] = ()


class EvolutionModuleVerifier:
    """Verify staged artifact integrity and static policy without executing code."""

    def __init__(self, workspace: Path, staging_root: Path | None = None) -> None:
        self.workspace = Path(workspace)
        self.staging_root = (
            Path(staging_root)
            if staging_root is not None
            else self.workspace / "memory" / "evolution_staging"
        )

    def verify(self, artifact_digest: str) -> EvolutionVerificationReport:
        staging_dir = self.staging_root / artifact_digest
        staging_path = _relative_to_workspace(staging_dir, self.workspace)
        checks: list[dict[str, Any]] = []

        metadata = self._read_staging_metadata(staging_dir, checks)
        if metadata is None:
            return _report(
                checks=checks,
                artifact_digest=artifact_digest,
                staging_path=staging_path,
            )

        artifact_dir = staging_dir / "artifact"
        if not artifact_dir.is_dir():
            checks.append(
                _check(
                    "artifact digest 校验",
                    False,
                    "digest_match",
                    "artifact directory is missing",
                )
            )
            return _report(
                checks=checks,
                metadata=metadata,
                artifact_digest=artifact_digest,
                staging_path=staging_path,
            )

        try:
            actual_digest = compute_artifact_digest(artifact_dir)
        except Exception as exc:
            checks.append(
                _check("artifact digest 校验", False, "digest_match", str(exc))
            )
            return _report(
                checks=checks,
                metadata=metadata,
                artifact_digest=artifact_digest,
                staging_path=staging_path,
            )

        checks.append(
            _check(
                "artifact digest 校验",
                actual_digest == artifact_digest,
                "digest_match",
                "artifact digest matches staging directory"
                if actual_digest == artifact_digest
                else "artifact digest does not match staging directory",
            )
        )

        try:
            manifest = read_package_manifest(artifact_dir)
        except Exception as exc:
            checks.append(
                _check("manifest 与 staging 对齐", False, "manifest_staging_alignment", str(exc))
            )
            return _report(
                checks=checks,
                metadata=metadata,
                artifact_digest=artifact_digest,
                staging_path=staging_path,
            )

        alignment_ok = (
            metadata.get("module_id") == manifest.module_id
            and metadata.get("module_type") == manifest.module_type
            and metadata.get("version") == manifest.version
            and metadata.get("artifact_digest") == artifact_digest
        )
        checks.append(
            _check(
                "manifest 与 staging 对齐",
                alignment_ok,
                "manifest_staging_alignment",
                "manifest matches staging metadata"
                if alignment_ok
                else "manifest does not match staging metadata",
            )
        )

        if manifest.module_type == "domain_pack":
            _append_domain_pack_check(artifact_dir, checks)

        permissions_evaluated, permissions_denied, unknown_keys = _append_permission_checks(
            manifest.permissions,
            checks,
        )
        _append_external_checks(manifest.external_endpoints, manifest.external_side_effects, checks)
        _append_context_budget_checks(manifest.context_budget, checks)
        _append_python_file_check(artifact_dir, checks)

        return _report(
            checks=checks,
            metadata=metadata,
            artifact_digest=artifact_digest,
            staging_path=staging_path,
            permissions_evaluated=permissions_evaluated,
            permissions_denied=tuple(permissions_denied),
            unknown_keys_rejected=tuple(unknown_keys),
        )

    def _read_staging_metadata(
        self,
        staging_dir: Path,
        checks: list[dict[str, Any]],
    ) -> dict[str, Any] | None:
        if not staging_dir.is_dir():
            checks.append(
                _check(
                    "staging 元数据完整性",
                    False,
                    "staging_metadata_integrity",
                    "staging directory is missing",
                )
            )
            return None
        staging_json = staging_dir / "staging.json"
        if not staging_json.exists():
            checks.append(
                _check(
                    "staging 元数据完整性",
                    False,
                    "staging_metadata_integrity",
                    "staging.json is missing",
                )
            )
            return None
        try:
            data = json.loads(staging_json.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError) as exc:
            checks.append(
                _check("staging 元数据完整性", False, "staging_metadata_integrity", str(exc))
            )
            return None
        if not isinstance(data, dict):
            checks.append(
                _check(
                    "staging 元数据完整性",
                    False,
                    "staging_metadata_integrity",
                    "staging.json must be a mapping",
                )
            )
            return None
        required = ("schema_version", "module_id", "module_type", "version", "artifact_digest")
        missing = [key for key in required if not data.get(key)]
        ok = not missing and data.get("schema_version") == _STAGING_SCHEMA_VERSION
        message = (
            "staging metadata is valid"
            if ok
            else "invalid staging metadata: " + ", ".join(missing or ["schema_version"])
        )
        checks.append(_check("staging 元数据完整性", ok, "staging_metadata_integrity", message))
        return data


def _append_domain_pack_check(artifact_dir: Path, checks: list[dict[str, Any]]) -> None:
    validated = DomainPackValidator(
        runtime_config=DomainPackRuntimeConfig(),
        strict_declarations=True,
    ).validate_pack(artifact_dir, source="workspace")
    ok = validated.status != "invalid"
    if validated.status == "unavailable":
        reason = validated.unavailable_reason or validated.validation_summary or "domain pack is unavailable"
    else:
        reason = validated.validation_summary or validated.unavailable_reason or "domain pack is valid"
    checks.append(
        _check(
            "domain_pack 声明校验",
            ok,
            "manifest_staging_alignment",
            reason,
        )
    )


def _append_permission_checks(
    permissions: dict[str, Any],
    checks: list[dict[str, Any]],
) -> tuple[dict[str, Any], list[str], list[str]]:
    evaluated = {key: permissions[key] for key in permissions if key in _PERMISSION_KEYS}
    denied: list[str] = []
    unknown = sorted(key for key in permissions if key not in _PERMISSION_KEYS)

    read_files = permissions.get("read_files", False)
    read_ok = isinstance(read_files, bool)
    _append_permission_check(checks, "read_files", read_ok, "read_files is allowed")
    if not read_ok:
        denied.append("permission_read_files")

    for key in _DENIED_BOOL_PERMISSIONS:
        value = permissions.get(key, False)
        ok = isinstance(value, bool) and value is False
        _append_permission_check(
            checks,
            key,
            ok,
            f"{key} is not requested" if ok else f"{key} is denied in Phase 1-C",
        )
        if not ok:
            denied.append(_PERMISSION_CODE_BY_KEY[key])

    device_domains = permissions.get("device_domains", ())
    device_ok = _is_string_list(device_domains) and not device_domains
    _append_permission_check(
        checks,
        "device_domains",
        device_ok,
        "device domains are not requested" if device_ok else "device domains are denied in Phase 1-C",
    )
    if not device_ok:
        denied.append("permission_device_domains")

    mcp_scopes = permissions.get("mcp_scopes", ())
    mcp_ok = _is_string_list(mcp_scopes) and set(mcp_scopes) <= {"read"}
    _append_permission_check(
        checks,
        "mcp_scopes",
        mcp_ok,
        "mcp scopes are allowed" if mcp_ok else "only mcp scope 'read' is allowed",
    )
    if not mcp_ok:
        denied.append("permission_mcp_scopes")

    checks.append(
        _check(
            "未知权限键",
            not unknown,
            "permission_unknown_keys",
            "no unknown permission keys" if not unknown else "unknown permission keys: " + ", ".join(unknown),
        )
    )
    return evaluated, denied, unknown


def _append_external_checks(
    external_endpoints: tuple[str, ...],
    external_side_effects: dict[str, Any],
    checks: list[dict[str, Any]],
) -> None:
    checks.append(
        _check(
            "外部端点",
            not external_endpoints,
            "external_endpoints",
            "no external endpoints declared"
            if not external_endpoints
            else "external endpoints are denied in Phase 1-C",
        )
    )
    writes_external_state = bool(external_side_effects.get("writes_external_state", False))
    checks.append(
        _check(
            "外部写状态",
            not writes_external_state,
            "external_writes_state",
            "no external state writes declared"
            if not writes_external_state
            else "external state writes are denied in Phase 1-C",
        )
    )


def _append_context_budget_checks(context_budget: dict[str, Any], checks: list[dict[str, Any]]) -> None:
    if "token_budget" not in context_budget:
        checks.append(
            _check("Token 预算", True, "context_token_budget", "token budget is not declared")
        )
        return
    token_budget = context_budget.get("token_budget")
    ok = isinstance(token_budget, int) and token_budget > 0
    checks.append(
        _check(
            "Token 预算",
            ok,
            "context_token_budget",
            "token budget is valid" if ok else "token_budget must be a positive integer",
        )
    )


def _append_python_file_check(artifact_dir: Path, checks: list[dict[str, Any]]) -> None:
    python_count = sum(1 for path in artifact_dir.rglob("*.py") if path.is_file())
    checks.append(
        _check(
            "Python 文件存在性",
            True,
            "contains_python_files",
            f"artifact contains {python_count} Python file(s)",
        )
    )


def _append_permission_check(
    checks: list[dict[str, Any]],
    key: str,
    ok: bool,
    message: str,
) -> None:
    checks.append(_check(key, ok, _PERMISSION_CODE_BY_KEY[key], message))


def _is_string_list(value: Any) -> bool:
    return isinstance(value, (list, tuple)) and all(isinstance(item, str) for item in value)


def _report(
    *,
    checks: list[dict[str, Any]],
    artifact_digest: str,
    staging_path: str,
    metadata: dict[str, Any] | None = None,
    permissions_evaluated: dict[str, Any] | None = None,
    permissions_denied: tuple[str, ...] = (),
    unknown_keys_rejected: tuple[str, ...] = (),
) -> EvolutionVerificationReport:
    failed = [check for check in checks if not check["ok"]]
    return EvolutionVerificationReport(
        ok=not failed,
        status="verified" if not failed else "failed",
        checks=tuple(checks),
        module_id=str((metadata or {}).get("module_id") or ""),
        module_type=str((metadata or {}).get("module_type") or ""),
        module_version=str((metadata or {}).get("version") or ""),
        artifact_digest=artifact_digest,
        staging_path=staging_path,
        error=str(failed[0]["message"]) if failed else "",
        permissions_evaluated=permissions_evaluated or {},
        permissions_denied=permissions_denied,
        unknown_keys_rejected=unknown_keys_rejected,
    )


def _check(name: str, ok: bool, code: str, message: str) -> dict[str, Any]:
    return {"name": name, "ok": ok, "code": code, "message": message}


def _relative_to_workspace(path: Path, workspace: Path) -> str:
    try:
        return path.relative_to(workspace).as_posix()
    except ValueError:
        return path.name
