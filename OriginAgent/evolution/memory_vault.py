"""Encrypted memory vault export and restore for Agent Passport migration."""

from __future__ import annotations

import base64
import binascii
import hashlib
import json
import os
import re
import tempfile
import uuid
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from cryptography.exceptions import InvalidTag
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.ledger import EvolutionLedger, canonical_dump, compute_event_hash

VAULT_SCHEMA_VERSION = "originagent.evolution.memory_vault.v1"
PAYLOAD_SCHEMA_VERSION = "originagent.evolution.memory_vault_payload.v1"
VAULT_DIGEST_DOMAIN = b"originagent.ec9.vault.v1"
ENCRYPTION_ALGORITHM = "AES-256-GCM"
AES_GCM_NONCE_BYTES = 12
VAULT_KEY_BYTES = 32

ALLOWED_MEMORY_FILES = (
    "SOUL.md",
    "USER.md",
    "memory/MEMORY.md",
    "memory/facts.jsonl",
    "memory/evolution_events.jsonl",
)

_BYTES32_RE = re.compile(r"^0x[0-9a-fA-F]{64}$")
_HEX64_RE = re.compile(r"^(?:0x)?[0-9a-fA-F]{64}$")
_SECRET_ASSIGNMENT_RE = re.compile(
    r"\b(?:api[_-]?key|secret|password|authorization|bearer)\b\s*[:=]\s*[\"']?[^\"',;\s<>]+",
    re.IGNORECASE,
)
_WINDOWS_PATH_RE = re.compile(r"(?:^|[^A-Za-z0-9])[A-Za-z]:[\\/][^\s\"'<>]+")
_UNIX_PRIVATE_PATH_RE = re.compile(r"(?:/Users/|/home/)[^\s\"'<>]+")
_URL_QUERY_RE = re.compile(r"https?://[^\s\"'<>?]+\?[^\s\"'<>]+")
_FORBIDDEN_KEYS = {
    "raw_prompt",
    "prompt",
    "file_content",
    "file_contents",
    "raw_tool_output",
    "facts",
    "facts_raw",
    "facts_text",
    "traceback",
    "hidden_reasoning",
    "private_telemetry",
    "private_key",
    "key_file",
}


class MemoryVaultError(ValueError):
    """Raised when a memory vault cannot be verified, exported, or imported."""


@dataclass(frozen=True)
class MemoryVaultImportResult:
    ok: bool
    dry_run: bool
    vault_digest: str
    restored_files: tuple[str, ...]
    source_ledger_terminal_hash: str
    import_event_hash: str = ""


def export_memory_vault(
    workspace: Path | str,
    *,
    passport_id: str,
    agent_key_hash: str,
    key_file: Path | str,
    out: Path | str,
) -> dict[str, Any]:
    workspace = Path(workspace)
    out = Path(out)
    key = _read_vault_key(Path(key_file))
    passport_id = _normalize_bytes32(passport_id, "passport_id")
    agent_key_hash = _normalize_bytes32(agent_key_hash, "agent_key_hash")

    ledger = EvolutionLedger(workspace)
    ledger_verification = ledger.verify_chain()
    if not ledger_verification.ok:
        raise MemoryVaultError(f"source evolution ledger is broken: {ledger_verification.error}")
    source_terminal_hash = ledger_verification.terminal_event_hash or ""

    files = _read_allowed_files(workspace)
    if not files:
        raise MemoryVaultError("no allowed memory files exist in source workspace")

    payload = {
        "schema_version": PAYLOAD_SCHEMA_VERSION,
        "files": files,
        "source_ledger_terminal_hash": source_terminal_hash,
    }
    payload_plaintext = _canonical_json(payload)
    payload_digest = _sha256(payload_plaintext)
    vault_nonce = f"0x{os.urandom(32).hex()}"
    aes_nonce = os.urandom(AES_GCM_NONCE_BYTES)
    metadata_aad = {
        "passport_id": passport_id,
        "agent_key_hash": agent_key_hash,
        "vault_nonce": vault_nonce,
        "created_at": datetime.now(timezone.utc).isoformat(),
        "source_ledger_terminal_hash": source_terminal_hash,
        "included_files": _included_files(files),
        "payload_digest": payload_digest,
        "encryption": {
            "algorithm": ENCRYPTION_ALGORITHM,
            "nonce_b64": _b64encode(aes_nonce),
        },
    }
    aad = _canonical_json(metadata_aad)
    encrypted_payload = AESGCM(key).encrypt(aes_nonce, payload_plaintext, aad)
    encrypted_payload_digest = _sha256(encrypted_payload)
    metadata = {
        **metadata_aad,
        "encrypted_payload_digest": encrypted_payload_digest,
    }
    vault_digest = _vault_digest(metadata, encrypted_payload_digest)
    vault = {
        "schema_version": VAULT_SCHEMA_VERSION,
        "vault_digest": vault_digest,
        "metadata": metadata,
        "encrypted_payload_b64": _b64encode(encrypted_payload),
    }
    _validate_public_vault(vault)
    _write_json_atomic(out, vault)
    ledger.append(
        EvolutionEvent.new(
            EventType.MEMORY_VAULT_EXPORTED,
            result={
                "passport_id": passport_id,
                "agent_key_hash": agent_key_hash,
                "vault_digest": vault_digest,
                "payload_digest": payload_digest,
                "encrypted_payload_digest": encrypted_payload_digest,
                "source_ledger_terminal_hash": source_terminal_hash,
            },
        )
    )
    return vault


def inspect_memory_vault(vault_path: Path | str) -> dict[str, Any]:
    vault = read_memory_vault(vault_path)
    _validate_public_vault(vault)
    return {
        "ok": True,
        "schema_version": vault["schema_version"],
        "vault_digest": vault["vault_digest"],
        "metadata": vault["metadata"],
    }


def verify_memory_vault(vault_path: Path | str, key_file: Path | str | None = None) -> dict[str, Any]:
    vault = read_memory_vault(vault_path)
    errors: list[str] = []
    try:
        _validate_public_vault(vault)
    except MemoryVaultError as exc:
        errors.append(str(exc))
    decrypted = False
    if key_file is not None and not errors:
        try:
            _decrypt_and_validate_payload(vault, _read_vault_key(Path(key_file)))
            decrypted = True
        except MemoryVaultError as exc:
            errors.append(str(exc))
    return {
        "ok": not errors,
        "errors": errors,
        "vault_digest": str(vault.get("vault_digest") or ""),
        "decrypted": decrypted,
    }


def import_memory_vault(
    vault_path: Path | str,
    *,
    key_file: Path | str,
    target_workspace: Path | str,
    apply: bool = False,
    replace: bool = False,
) -> MemoryVaultImportResult:
    vault = read_memory_vault(vault_path)
    key = _read_vault_key(Path(key_file))
    payload = _decrypt_and_validate_payload(vault, key)
    target = Path(target_workspace)
    files = list(payload["files"])
    paths = tuple(str(item["path"]) for item in files)
    conflicts = [path for path in paths if (target / path).exists()]
    if conflicts and not replace:
        raise MemoryVaultError(f"target workspace already has conflicting files: {', '.join(conflicts)}")
    if not apply:
        return MemoryVaultImportResult(
            ok=True,
            dry_run=True,
            vault_digest=str(vault["vault_digest"]),
            restored_files=paths,
            source_ledger_terminal_hash=str(payload["source_ledger_terminal_hash"]),
        )

    _write_payload_files(target, files)
    event = EvolutionLedger(target).append(
        EvolutionEvent.new(
            EventType.MEMORY_VAULT_IMPORTED,
            result={
                "passport_id": vault["metadata"]["passport_id"],
                "agent_key_hash": vault["metadata"]["agent_key_hash"],
                "vault_digest": vault["vault_digest"],
                "target_workspace_initialized": True,
            },
        )
    )
    return MemoryVaultImportResult(
        ok=True,
        dry_run=False,
        vault_digest=str(vault["vault_digest"]),
        restored_files=paths,
        source_ledger_terminal_hash=str(payload["source_ledger_terminal_hash"]),
        import_event_hash=event.event_hash,
    )


def read_memory_vault(path: Path | str) -> dict[str, Any]:
    data = json.loads(Path(path).read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise MemoryVaultError("memory vault must be a JSON object")
    return data


def _read_allowed_files(workspace: Path) -> list[dict[str, Any]]:
    files: list[dict[str, Any]] = []
    for relative in ALLOWED_MEMORY_FILES:
        path = workspace / relative
        if not path.exists():
            continue
        if not path.is_file():
            raise MemoryVaultError(f"allowed memory path is not a file: {relative}")
        content = path.read_bytes()
        files.append(
            {
                "path": relative,
                "sha256": _sha256(content),
                "size": len(content),
                "content_b64": _b64encode(content),
            }
        )
    return sorted(files, key=lambda item: str(item["path"]))


def _included_files(files: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return [
        {
            "path": item["path"],
            "sha256": item["sha256"],
            "size": item["size"],
        }
        for item in files
    ]


def _decrypt_and_validate_payload(vault: dict[str, Any], key: bytes) -> dict[str, Any]:
    _validate_public_vault(vault)
    metadata = vault["metadata"]
    encrypted_payload = _b64decode(str(vault["encrypted_payload_b64"]), "encrypted_payload_b64")
    nonce = _b64decode(str(metadata["encryption"]["nonce_b64"]), "encryption.nonce_b64")
    aad = _canonical_json(_metadata_for_aad(metadata))
    try:
        payload_bytes = AESGCM(key).decrypt(nonce, encrypted_payload, aad)
    except InvalidTag as exc:
        raise MemoryVaultError("vault decryption failed: authentication tag mismatch") from exc
    if _sha256(payload_bytes) != metadata["payload_digest"]:
        raise MemoryVaultError("payload_digest mismatch after decryption")
    try:
        payload = json.loads(payload_bytes.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError) as exc:
        raise MemoryVaultError("payload is not canonical JSON") from exc
    if not isinstance(payload, dict):
        raise MemoryVaultError("payload must be an object")
    if _canonical_json(payload) != payload_bytes:
        raise MemoryVaultError("payload is not canonical JSON")
    _validate_payload(payload, metadata)
    return payload


def _validate_public_vault(vault: dict[str, Any]) -> None:
    errors = _public_vault_errors(vault)
    if errors:
        raise MemoryVaultError("; ".join(errors))


def _public_vault_errors(vault: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    if vault.get("schema_version") != VAULT_SCHEMA_VERSION:
        errors.append(f"schema_version must be {VAULT_SCHEMA_VERSION}")
    metadata = vault.get("metadata")
    if not isinstance(metadata, dict):
        errors.append("metadata must be an object")
        return errors
    for field in (
        "passport_id",
        "agent_key_hash",
        "vault_nonce",
        "created_at",
        "source_ledger_terminal_hash",
        "included_files",
        "payload_digest",
        "encrypted_payload_digest",
        "encryption",
    ):
        if field not in metadata:
            errors.append(f"missing metadata.{field}")
    for field in ("passport_id", "agent_key_hash", "vault_nonce"):
        if field in metadata and not _BYTES32_RE.match(str(metadata[field])):
            errors.append(f"metadata.{field} must be 0x-prefixed bytes32")
    for field in ("payload_digest", "encrypted_payload_digest"):
        if field in metadata and not _is_hex64(str(metadata[field])):
            errors.append(f"metadata.{field} must be lowercase sha256 hex")
    if metadata.get("source_ledger_terminal_hash") and not _is_hex64(str(metadata["source_ledger_terminal_hash"])):
        errors.append("metadata.source_ledger_terminal_hash must be empty or lowercase sha256 hex")
    included = metadata.get("included_files")
    if not isinstance(included, list):
        errors.append("metadata.included_files must be a list")
    else:
        for item in included:
            if not _valid_file_summary(item):
                errors.append("metadata.included_files entries must contain allowed path, sha256, and size")
                break
    encryption = metadata.get("encryption")
    if not isinstance(encryption, dict):
        errors.append("metadata.encryption must be an object")
    else:
        if encryption.get("algorithm") != ENCRYPTION_ALGORITHM:
            errors.append(f"metadata.encryption.algorithm must be {ENCRYPTION_ALGORITHM}")
        try:
            nonce = _b64decode(str(encryption.get("nonce_b64") or ""), "metadata.encryption.nonce_b64")
            if len(nonce) != AES_GCM_NONCE_BYTES:
                errors.append("metadata.encryption.nonce_b64 must decode to 12 bytes")
        except MemoryVaultError as exc:
            errors.append(str(exc))
    encrypted_payload_b64 = vault.get("encrypted_payload_b64")
    try:
        encrypted_payload = _b64decode(str(encrypted_payload_b64 or ""), "encrypted_payload_b64")
    except MemoryVaultError as exc:
        encrypted_payload = b""
        errors.append(str(exc))
    if encrypted_payload and metadata.get("encrypted_payload_digest") != _sha256(encrypted_payload):
        errors.append("encrypted_payload_digest mismatch")
    if "vault_digest" not in vault or not _is_hex64(str(vault.get("vault_digest") or "")):
        errors.append("vault_digest must be lowercase sha256 hex")
    elif not errors:
        expected = _vault_digest(metadata, str(metadata["encrypted_payload_digest"]))
        if str(vault["vault_digest"]) != expected:
            errors.append("vault_digest mismatch")
    errors.extend(_privacy_errors({"schema_version": vault.get("schema_version"), "vault_digest": vault.get("vault_digest"), "metadata": metadata}))
    return errors


def _validate_payload(payload: dict[str, Any], metadata: dict[str, Any]) -> None:
    if payload.get("schema_version") != PAYLOAD_SCHEMA_VERSION:
        raise MemoryVaultError(f"payload.schema_version must be {PAYLOAD_SCHEMA_VERSION}")
    if payload.get("source_ledger_terminal_hash") != metadata.get("source_ledger_terminal_hash"):
        raise MemoryVaultError("payload source_ledger_terminal_hash does not match metadata")
    files = payload.get("files")
    if not isinstance(files, list):
        raise MemoryVaultError("payload.files must be a list")
    summaries: list[dict[str, Any]] = []
    for item in files:
        if not isinstance(item, dict):
            raise MemoryVaultError("payload file entries must be objects")
        if not _valid_allowed_path(str(item.get("path") or "")):
            raise MemoryVaultError(f"payload contains disallowed file path: {item.get('path')}")
        content = _b64decode(str(item.get("content_b64") or ""), f"payload.files[{item.get('path')}].content_b64")
        sha = _sha256(content)
        if item.get("sha256") != sha:
            raise MemoryVaultError(f"payload file sha256 mismatch: {item.get('path')}")
        if item.get("size") != len(content):
            raise MemoryVaultError(f"payload file size mismatch: {item.get('path')}")
        summaries.append({"path": item["path"], "sha256": sha, "size": len(content)})
    if sorted(summaries, key=lambda item: item["path"]) != sorted(metadata.get("included_files") or [], key=lambda item: item["path"]):
        raise MemoryVaultError("payload files do not match metadata.included_files")
    _verify_payload_ledger(files, str(metadata.get("source_ledger_terminal_hash") or ""))


def _verify_payload_ledger(files: list[dict[str, Any]], expected_terminal_hash: str) -> None:
    ledger_item = next((item for item in files if item.get("path") == "memory/evolution_events.jsonl"), None)
    if ledger_item is None:
        if expected_terminal_hash:
            raise MemoryVaultError("source_ledger_terminal_hash is set but payload has no evolution ledger")
        return
    content = _b64decode(str(ledger_item.get("content_b64") or ""), "memory/evolution_events.jsonl")
    previous_hash: str | None = None
    terminal_hash = ""
    try:
        for line_number, line in enumerate(content.decode("utf-8").splitlines(), start=1):
            if not line.strip():
                continue
            event = json.loads(line)
            if not isinstance(event, dict):
                raise MemoryVaultError(f"ledger line {line_number} must be an object")
            if event.get("previous_event_hash") != previous_hash:
                raise MemoryVaultError(f"ledger previous_event_hash mismatch at line {line_number}")
            event_hash = compute_event_hash(event)
            if event.get("event_hash") != event_hash:
                raise MemoryVaultError(f"ledger event_hash mismatch at line {line_number}")
            previous_hash = event_hash
            terminal_hash = event_hash
    except UnicodeDecodeError as exc:
        raise MemoryVaultError("evolution ledger is not UTF-8") from exc
    if terminal_hash != expected_terminal_hash:
        raise MemoryVaultError("source_ledger_terminal_hash does not match payload ledger")


def _write_payload_files(target: Path, files: list[dict[str, Any]]) -> None:
    target.mkdir(parents=True, exist_ok=True)
    temp_root = target / f".memory_vault_import_tmp_{uuid.uuid4().hex}"
    try:
        for item in files:
            relative = str(item["path"])
            content = _b64decode(str(item["content_b64"]), f"{relative}.content_b64")
            temp_path = temp_root / relative
            temp_path.parent.mkdir(parents=True, exist_ok=True)
            temp_path.write_bytes(content)
        for item in files:
            relative = str(item["path"])
            src = temp_root / relative
            dst = target / relative
            dst.parent.mkdir(parents=True, exist_ok=True)
            os.replace(src, dst)
    finally:
        _remove_tree(temp_root)


def _metadata_for_aad(metadata: dict[str, Any]) -> dict[str, Any]:
    return {key: value for key, value in metadata.items() if key != "encrypted_payload_digest"}


def _vault_digest(metadata: dict[str, Any], encrypted_payload_digest: str) -> str:
    metadata_digest = _sha256(_canonical_json(metadata))
    return hashlib.sha256(
        VAULT_DIGEST_DOMAIN + bytes.fromhex(metadata_digest) + bytes.fromhex(encrypted_payload_digest)
    ).hexdigest()


def _valid_file_summary(item: Any) -> bool:
    return (
        isinstance(item, dict)
        and _valid_allowed_path(str(item.get("path") or ""))
        and _is_hex64(str(item.get("sha256") or ""))
        and isinstance(item.get("size"), int)
        and item["size"] >= 0
    )


def _valid_allowed_path(path: str) -> bool:
    return path in ALLOWED_MEMORY_FILES and not path.startswith("/") and ".." not in Path(path).parts


def _normalize_bytes32(value: str, field: str) -> str:
    if not _BYTES32_RE.match(value):
        raise MemoryVaultError(f"{field} must be a 0x-prefixed bytes32 value")
    return f"0x{value[2:].lower()}"


def _is_hex64(value: str) -> bool:
    return bool(re.match(r"^[0-9a-f]{64}$", value))


def _read_vault_key(path: Path) -> bytes:
    data = path.read_bytes()
    stripped = data.strip()
    text = stripped.decode("ascii", errors="ignore")
    normalized = text[2:] if text.startswith("0x") else text
    if re.fullmatch(r"[0-9a-fA-F]{64}", normalized):
        return bytes.fromhex(normalized)
    try:
        decoded = base64.urlsafe_b64decode(stripped)
        if len(decoded) == VAULT_KEY_BYTES:
            return decoded
    except (binascii.Error, ValueError):
        pass
    if len(data) == VAULT_KEY_BYTES:
        return data
    raise MemoryVaultError("vault key-file must contain 32 raw bytes, base64url bytes, or 64 hex characters")


def _b64encode(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).decode("ascii")


def _b64decode(value: str, field: str) -> bytes:
    try:
        return base64.urlsafe_b64decode(value.encode("ascii"))
    except (binascii.Error, UnicodeEncodeError) as exc:
        raise MemoryVaultError(f"{field} must be base64url data") from exc


def _canonical_json(value: Any) -> bytes:
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":"), allow_nan=False).encode("utf-8")


def _sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _write_json_atomic(path: Path, data: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = path.with_name(f".{path.name}.tmp-{uuid.uuid4().hex}")
    try:
        tmp_path.write_text(json.dumps(data, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")
        os.replace(tmp_path, path)
    finally:
        tmp_path.unlink(missing_ok=True)


def _remove_tree(path: Path) -> None:
    if not path.exists():
        return
    for child in sorted(path.rglob("*"), reverse=True):
        if child.is_file() or child.is_symlink():
            child.unlink(missing_ok=True)
        elif child.is_dir():
            child.rmdir()
    path.rmdir()


def _privacy_errors(value: Any, path: str = "$") -> list[str]:
    errors: list[str] = []
    if isinstance(value, dict):
        for key, child in value.items():
            child_path = f"{path}.{key}"
            if str(key) in _FORBIDDEN_KEYS:
                errors.append(f"forbidden private field: {child_path}")
            errors.extend(_privacy_errors(child, child_path))
    elif isinstance(value, list):
        for index, child in enumerate(value):
            errors.extend(_privacy_errors(child, f"{path}[{index}]"))
    elif isinstance(value, str):
        if _WINDOWS_PATH_RE.search(value) or _UNIX_PRIVATE_PATH_RE.search(value):
            errors.append(f"local absolute path detected at {path}")
        if _URL_QUERY_RE.search(value):
            errors.append(f"URL query detected at {path}")
        if _SECRET_ASSIGNMENT_RE.search(value):
            errors.append(f"secret-like string detected at {path}")
    return errors
