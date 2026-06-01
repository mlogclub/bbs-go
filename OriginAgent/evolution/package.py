"""Local evolution module package loading and artifact handling."""

from __future__ import annotations

import hashlib
import os
import shutil
from dataclasses import dataclass
from pathlib import Path

import yaml

from OriginAgent.agent.domain_packs import DomainPackRuntimeConfig, DomainPackValidator
from OriginAgent.evolution.manifest import EvolutionManifest, validate_manifest

EVOLUTION_MANIFEST_FILENAME = "evolution_manifest.yaml"
_EXCLUDED_DIR_NAMES = frozenset({".git", "__pycache__", ".pytest_cache", ".ruff_cache"})
_EXCLUDED_SUFFIXES = frozenset({".pyc", ".pyo"})


class EvolutionPackageError(ValueError):
    """Raised when an evolution module package cannot be staged."""


@dataclass(frozen=True)
class EvolutionPackage:
    """Validated local evolution package metadata."""

    source_path: Path
    source_name: str
    manifest: EvolutionManifest
    artifact_digest: str


def load_package(source_path: str | Path) -> EvolutionPackage:
    """Load and validate a local evolution package directory."""

    resolved_source = _resolve_source_dir(source_path)
    manifest = read_package_manifest(resolved_source)
    if manifest.module_type == "domain_pack":
        _validate_domain_pack_package(resolved_source)
    digest = compute_artifact_digest(resolved_source)
    return EvolutionPackage(
        source_path=resolved_source,
        source_name=resolved_source.name,
        manifest=manifest,
        artifact_digest=digest,
    )


def read_package_manifest(source_path: str | Path) -> EvolutionManifest:
    """Read and validate ``evolution_manifest.yaml`` from a local package directory."""

    resolved_source = _resolve_source_dir(source_path)
    manifest_path = resolved_source / EVOLUTION_MANIFEST_FILENAME
    if not manifest_path.exists():
        raise EvolutionPackageError(f"missing {EVOLUTION_MANIFEST_FILENAME}")
    try:
        raw = yaml.safe_load(manifest_path.read_text(encoding="utf-8"))
    except (OSError, yaml.YAMLError) as exc:
        raise EvolutionPackageError(f"invalid {EVOLUTION_MANIFEST_FILENAME}: {exc}") from exc
    if not isinstance(raw, dict):
        raise EvolutionPackageError(f"{EVOLUTION_MANIFEST_FILENAME} must be a mapping")
    try:
        return validate_manifest(raw)
    except ValueError as exc:
        raise EvolutionPackageError(str(exc)) from exc


def compute_artifact_digest(source_path: str | Path) -> str:
    """Compute a deterministic tree digest for a local package directory."""

    resolved_source = _resolve_source_dir(source_path)
    digest = hashlib.sha256()
    for file_path in _included_artifact_files(resolved_source):
        rel_path = file_path.relative_to(resolved_source).as_posix()
        digest.update(rel_path.encode("utf-8"))
        digest.update(b"\0")
        digest.update(file_path.read_bytes())
        digest.update(b"\0")
    return digest.hexdigest()


def copy_artifact(source_path: str | Path, target_path: str | Path) -> None:
    """Copy included package artifact files into ``target_path`` without executing them."""

    resolved_source = _resolve_source_dir(source_path)
    target = Path(target_path)
    if target.exists():
        _cleanup_path(target)
    target.mkdir(parents=True, exist_ok=True)
    for file_path in _included_artifact_files(resolved_source):
        rel_path = file_path.relative_to(resolved_source)
        destination = target / rel_path
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(file_path, destination)


def _resolve_source_dir(source_path: str | Path) -> Path:
    resolved = Path(source_path).expanduser().resolve()
    if not resolved.exists() or not resolved.is_dir():
        raise EvolutionPackageError("source path must be an existing directory")
    if resolved.is_symlink():
        raise EvolutionPackageError("source path must not be a symlink")
    return resolved


def _validate_domain_pack_package(source_path: Path) -> None:
    if not (source_path / "domain_pack.yaml").exists():
        raise EvolutionPackageError("domain_pack modules must include domain_pack.yaml")
    validated = DomainPackValidator(
        runtime_config=DomainPackRuntimeConfig(),
        strict_declarations=True,
    ).validate_pack(source_path, source="workspace")
    if validated.status == "invalid":
        reason = validated.validation_summary or validated.unavailable_reason or "invalid domain pack"
        raise EvolutionPackageError(f"invalid domain_pack.yaml: {reason}")


def _included_artifact_files(root: Path) -> list[Path]:
    files: list[Path] = []
    for dirpath, dirnames, filenames in os.walk(root):
        current = Path(dirpath)
        for dirname in list(dirnames):
            child = current / dirname
            if child.is_symlink():
                raise EvolutionPackageError("artifact must not contain symlinks")
        dirnames[:] = sorted(dirname for dirname in dirnames if dirname not in _EXCLUDED_DIR_NAMES)
        for filename in sorted(filenames):
            child = current / filename
            if child.is_symlink():
                raise EvolutionPackageError("artifact must not contain symlinks")
            if _is_excluded_file(child):
                continue
            if child.is_file():
                files.append(child)
    return sorted(files, key=lambda path: path.relative_to(root).as_posix())


def _is_excluded_file(path: Path) -> bool:
    if path.name in _EXCLUDED_DIR_NAMES:
        return True
    return path.suffix in _EXCLUDED_SUFFIXES


def _cleanup_path(path: Path) -> None:
    if path.is_dir() and not path.is_symlink():
        shutil.rmtree(path, ignore_errors=True)
    else:
        try:
            path.unlink()
        except FileNotFoundError:
            pass
