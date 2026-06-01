"""Local plaintext Ed25519 identity for evolution event signing."""

from __future__ import annotations

import base64
import json
import os
from contextlib import suppress
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Mapping

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey, Ed25519PublicKey
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    NoEncryption,
    PrivateFormat,
    PublicFormat,
)

IDENTITY_SCHEMA_VERSION = "originagent.evolution.identity.v1"
IDENTITY_SCHEME = "ed25519"
PLAINTEXT_WARNING = "plaintext storage; upgrade to OS keychain in Phase 2+"


@dataclass(frozen=True)
class EvolutionIdentity:
    schema_version: str
    scheme: str
    public_key: str
    private_key: str
    created_at: str
    protection: str
    protection_warning: str


class EvolutionIdentityStore:
    """Load or create a local Ed25519 identity.

    Phase 1-H stores the private key as plaintext JSON. The file is chmod 0600
    where supported, but this is not equivalent to OS keychain protection.
    """

    def __init__(self, path: Path | None = None) -> None:
        self.path = Path(path) if path is not None else Path.home() / ".originagent" / "evolution_identity.json"
        self._identity = self._load_or_create()

    @property
    def public_key_b64(self) -> str:
        return self._identity.public_key

    def sign(self, message: bytes) -> str:
        private_key = Ed25519PrivateKey.from_private_bytes(_b64decode(self._identity.private_key))
        return _b64encode(private_key.sign(message))

    def verify(self, message: bytes, signature: str, public_key: str | None = None) -> bool:
        return verify_signature(message, signature, public_key or self._identity.public_key)

    def _load_or_create(self) -> EvolutionIdentity:
        if self.path.exists():
            data = json.loads(self.path.read_text(encoding="utf-8"))
            if not isinstance(data, dict):
                raise ValueError("evolution identity must be a mapping")
            return _identity_from_dict(data)
        private_key = Ed25519PrivateKey.generate()
        public_key = private_key.public_key()
        identity = EvolutionIdentity(
            schema_version=IDENTITY_SCHEMA_VERSION,
            scheme=IDENTITY_SCHEME,
            public_key=_b64encode(
                public_key.public_bytes(encoding=Encoding.Raw, format=PublicFormat.Raw)
            ),
            private_key=_b64encode(
                private_key.private_bytes(
                    encoding=Encoding.Raw,
                    format=PrivateFormat.Raw,
                    encryption_algorithm=NoEncryption(),
                )
            ),
            created_at=datetime.now(timezone.utc).isoformat(),
            protection="plaintext",
            protection_warning=PLAINTEXT_WARNING,
        )
        self.path.parent.mkdir(parents=True, exist_ok=True)
        self.path.write_text(json.dumps(identity.__dict__, ensure_ascii=False, sort_keys=True) + "\n", encoding="utf-8")
        with suppress(PermissionError):
            os.chmod(self.path, 0o600)
        return identity


def verify_signature(message: bytes, signature: str, public_key: str) -> bool:
    try:
        verifier = Ed25519PublicKey.from_public_bytes(_b64decode(public_key))
        verifier.verify(_b64decode(signature), message)
        return True
    except (InvalidSignature, ValueError, TypeError):
        return False


def verify_event_signature(event: Mapping[str, Any]) -> bool:
    signature = str(event.get("signature") or "")
    public_key = str(event.get("actor_public_key") or "")
    event_hash = str(event.get("event_hash") or "")
    if not signature or not public_key or not event_hash:
        return False
    return verify_signature(event_hash.encode("utf-8"), signature, public_key)


def _identity_from_dict(data: dict[str, Any]) -> EvolutionIdentity:
    identity = EvolutionIdentity(
        schema_version=str(data.get("schema_version") or ""),
        scheme=str(data.get("scheme") or ""),
        public_key=str(data.get("public_key") or ""),
        private_key=str(data.get("private_key") or ""),
        created_at=str(data.get("created_at") or ""),
        protection=str(data.get("protection") or ""),
        protection_warning=str(data.get("protection_warning") or ""),
    )
    if identity.schema_version != IDENTITY_SCHEMA_VERSION:
        raise ValueError(f"unsupported identity schema_version: {identity.schema_version}")
    if identity.scheme != IDENTITY_SCHEME:
        raise ValueError(f"unsupported identity scheme: {identity.scheme}")
    if not identity.public_key or not identity.private_key:
        raise ValueError("identity keys are required")
    return identity


def _b64encode(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).decode("ascii")


def _b64decode(data: str) -> bytes:
    return base64.urlsafe_b64decode(data.encode("ascii"))
