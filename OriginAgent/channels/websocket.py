"""WebSocket server channel: OriginAgent acts as a WebSocket server and serves connected clients."""

from __future__ import annotations

import asyncio
import base64
import binascii
import email.utils
import hashlib
import hmac
import http
import ipaddress
import json
import mimetypes
import re
import secrets
import shutil
import ssl
import socket
import time
import uuid
from pathlib import Path
from typing import TYPE_CHECKING, Any, Callable, Self
from urllib.parse import parse_qs, unquote, urlparse

from loguru import logger
from pydantic import Field, field_validator, model_validator
from websockets.asyncio.server import ServerConnection, serve
from websockets.datastructures import Headers
from websockets.exceptions import ConnectionClosed
from websockets.http11 import Request as WsRequest
from websockets.http11 import Response

from OriginAgent.bus.events import OUTBOUND_META_AGENT_UI, OutboundMessage
from OriginAgent.bus.queue import MessageBus
from OriginAgent.channels.base import BaseChannel
from OriginAgent.command.builtin import builtin_command_palette
from OriginAgent.config.paths import get_media_dir
from OriginAgent.config.schema import Base
from OriginAgent.session.goal_state import goal_state_ws_blob
from OriginAgent.utils.helpers import safe_filename
from OriginAgent.utils.media_decode import (
    FileSizeExceeded,
    save_base64_data_url,
)
from OriginAgent.utils.subagent_channel_display import scrub_subagent_messages_for_channel
from OriginAgent.utils.webui_thread_disk import delete_webui_thread
from OriginAgent.utils.webui_transcript import build_webui_thread_response

if TYPE_CHECKING:
    from OriginAgent.session.manager import SessionManager


def _strip_trailing_slash(path: str) -> str:
    if len(path) > 1 and path.endswith("/"):
        return path.rstrip("/")
    return path or "/"


def _normalize_config_path(path: str) -> str:
    return _strip_trailing_slash(path)


def _append_buttons_as_text(text: str, buttons: list[list[str]]) -> str:
    labels = [label for row in buttons for label in row if label]
    if not labels:
        return text
    fallback = "\n".join(f"{index}. {label}" for index, label in enumerate(labels, 1))
    return f"{text}\n\n{fallback}" if text else fallback


class WebSocketConfig(Base):
    """WebSocket server channel configuration.

    Clients connect with URLs like ``ws://{host}:{port}{path}?client_id=...&token=...``.
    - ``client_id``: Used for ``allow_from`` authorization; if omitted, a value is generated and logged.
    - ``token``: If non-empty, the ``token`` query param may match this static secret; short-lived tokens
      from ``token_issue_path`` are also accepted.
    - ``token_issue_path``: If non-empty, **GET** (HTTP/1.1) to this path returns JSON
      ``{"token": "...", "expires_in": <seconds>}``; use ``?token=...`` when opening the WebSocket.
      Must differ from ``path`` (the WS upgrade path). If the client runs in the **same process** as
      OriginAgent and shares the asyncio loop, use a thread or async HTTP client for GET—do not call
      blocking ``urllib`` or synchronous ``httpx`` from inside a coroutine.
    - ``token_issue_secret``: If non-empty, token requests must send ``Authorization: Bearer <secret>`` or
      ``X-OriginAgent-Auth: <secret>``.
    - ``websocket_requires_token``: If True, the handshake must include a valid token (static or issued and not expired).
    - Each connection has its own session: a unique ``chat_id`` maps to the agent session internally.
    - ``media`` field in outbound messages contains local filesystem paths; remote clients need a
      shared filesystem or an HTTP file server to access these files.
    """

    enabled: bool = False
    host: str = "127.0.0.1"
    port: int = 8765
    path: str = "/"
    token: str = ""
    token_issue_path: str = ""
    token_issue_secret: str = ""
    token_ttl_s: int = Field(default=300, ge=30, le=86_400)
    websocket_requires_token: bool = True
    allow_from: list[str] = Field(default_factory=lambda: ["*"])
    streaming: bool = True
    # Default 36 MB, upper 40 MB: supports up to 4 images at ~6 MB each after
    # client-side Worker normalization (see webui Composer). 4 × 6 MB × 1.37
    # (base64 overhead) + envelope framing stays under 36 MB; the 40 MB ceiling
    # leaves a small margin for sender slop without opening a DoS avenue.
    max_message_bytes: int = Field(default=37_748_736, ge=1024, le=41_943_040)
    ping_interval_s: float = Field(default=20.0, ge=5.0, le=300.0)
    ping_timeout_s: float = Field(default=20.0, ge=5.0, le=300.0)
    ssl_certfile: str = ""
    ssl_keyfile: str = ""

    @field_validator("path")
    @classmethod
    def path_must_start_with_slash(cls, value: str) -> str:
        if not value.startswith("/"):
            raise ValueError('path must start with "/"')
        return _normalize_config_path(value)

    @field_validator("token_issue_path")
    @classmethod
    def token_issue_path_format(cls, value: str) -> str:
        value = value.strip()
        if not value:
            return ""
        if not value.startswith("/"):
            raise ValueError('token_issue_path must start with "/"')
        return _normalize_config_path(value)

    @model_validator(mode="after")
    def token_issue_path_differs_from_ws_path(self) -> Self:
        if not self.token_issue_path:
            return self
        if _normalize_config_path(self.token_issue_path) == _normalize_config_path(self.path):
            raise ValueError("token_issue_path must differ from path (the WebSocket upgrade path)")
        return self

    @model_validator(mode="after")
    def wildcard_host_requires_auth(self) -> Self:
        if self.host not in ("0.0.0.0", "::"):
            return self
        if self.token.strip() or self.token_issue_secret.strip():
            return self
        raise ValueError(
            "host is 0.0.0.0 (all interfaces) but neither token nor "
            "token_issue_secret is set — set one to prevent unauthenticated access"
        )


def _http_json_response(data: dict[str, Any], *, status: int = 200) -> Response:
    body = json.dumps(data, ensure_ascii=False).encode("utf-8")
    headers = Headers(
        [
            ("Date", email.utils.formatdate(usegmt=True)),
            ("Connection", "close"),
            ("Content-Length", str(len(body))),
            ("Content-Type", "application/json; charset=utf-8"),
        ]
    )
    reason = http.HTTPStatus(status).phrase
    return Response(status, reason, headers, body)


def publish_runtime_model_update(
    bus: MessageBus,
    model: str,
    model_preset: str | None,
) -> None:
    """Enqueue a runtime model snapshot for all embedded WebUI subscribers."""
    bus.outbound.put_nowait(OutboundMessage(
        channel="websocket",
        chat_id="*",
        content="",
        metadata={
            "_runtime_model_updated": True,
            "model": model,
            "model_preset": model_preset,
        },
    ))


def _default_model_name_from_config() -> str | None:
    """Return the configured default model for readonly webui display."""
    try:
        from OriginAgent.config.loader import load_config

        model = load_config().resolve_preset().model.strip()
        return model or None
    except Exception as e:
        logger.debug("webui bootstrap could not load model name: {}", e)
        return None


def _resolve_bootstrap_model_name(
    runtime_name: Callable[[], str | None] | None,
) -> str | None:
    """Prefer an in-process resolver, else fall back to on-disk config."""
    if runtime_name is not None:
        try:
            raw = runtime_name()
        except Exception as e:
            logger.debug("bootstrap runtime model resolver failed: {}", e)
        else:
            if isinstance(raw, str):
                stripped = raw.strip()
                if stripped:
                    return stripped
    return _default_model_name_from_config()


def _parse_request_path(path_with_query: str) -> tuple[str, dict[str, list[str]]]:
    """Parse normalized path and query parameters in one pass."""
    parsed = urlparse("ws://x" + path_with_query)
    path = _strip_trailing_slash(parsed.path or "/")
    return path, parse_qs(parsed.query, keep_blank_values=True)


def _normalize_http_path(path_with_query: str) -> str:
    """Return the path component (no query string), with trailing slash normalized (root stays ``/``)."""
    return _parse_request_path(path_with_query)[0]


def _parse_query(path_with_query: str) -> dict[str, list[str]]:
    return _parse_request_path(path_with_query)[1]


def _query_first(query: dict[str, list[str]], key: str) -> str | None:
    """Return the first value for *key*, or None."""
    values = query.get(key)
    return values[0] if values else None


def _query_bool(query: dict[str, list[str]], key: str) -> bool | None:
    raw = _query_first(query, key)
    if raw is None:
        return None
    value = raw.strip().lower()
    if value in {"1", "true", "yes", "on"}:
        return True
    if value in {"0", "false", "no", "off"}:
        return False
    return None


def _mask_secret_hint(secret: str | None) -> str | None:
    if not secret:
        return None
    if len(secret) <= 8:
        return "••••"
    return f"{secret[:4]}••••{secret[-4:]}"


_WEB_SEARCH_PROVIDER_OPTIONS: tuple[dict[str, str], ...] = (
    {"name": "duckduckgo", "label": "DuckDuckGo", "credential": "none"},
    {"name": "brave", "label": "Brave Search", "credential": "api_key"},
    {"name": "tavily", "label": "Tavily", "credential": "api_key"},
    {"name": "searxng", "label": "SearXNG", "credential": "base_url"},
    {"name": "jina", "label": "Jina", "credential": "api_key"},
    {"name": "kagi", "label": "Kagi", "credential": "api_key"},
    {"name": "olostep", "label": "Olostep", "credential": "api_key"},
)
_WEB_SEARCH_PROVIDER_BY_NAME = {
    provider["name"]: provider for provider in _WEB_SEARCH_PROVIDER_OPTIONS
}

_RUNTIME_PROFILE_OPTIONS = {"default", "safe", "household_safe", "local_dev", "automation"}
_PROVIDER_RETRY_MODE_OPTIONS = {"standard", "persistent"}
_EVOLUTION_MODE_OPTIONS = {"conservative", "curated", "exploratory", "aggressive"}
_SESSION_SEARCH_BACKEND_OPTIONS = {"auto", "literal", "sqlite_fts"}
_EXEC_PROFILE_OPTIONS = {"secure", "local_dev", "disabled"}
_EXEC_SHELL_SYNTAX_POLICY_OPTIONS = {"restricted", "shell"}
_DEVICE_MODE_OPTIONS = {"dry_run", "real"}
_DEVICE_BACKEND_OPTIONS = {"none", "fake", "lighting_client"}
_AUDIT_MODE_OPTIONS = {"off", "minimal", "security"}

_MCP_SERVER_NAME_RE = re.compile(r"^[A-Za-z0-9_-]{1,64}$")
_MCP_SECRET_HINT = "••••"
_HA_MCP_PATH = "/api/mcp"


def _settings_runtime_controls_payload(config: Any) -> dict[str, Any]:
    defaults = config.agents.defaults
    evolution = defaults.learning.evolution
    return {
        "channels": {
            "send_progress": bool(config.channels.send_progress),
            "send_tool_hints": bool(config.channels.send_tool_hints),
            "show_reasoning": bool(config.channels.show_reasoning),
        },
        "agent": {
            "unified_session": bool(defaults.unified_session),
            "cold_archive_enabled": bool(defaults.cold_archive_enabled),
            "allow_agent_initiated_messages": bool(defaults.allow_agent_initiated_messages),
            "auxiliary_enabled": bool(defaults.auxiliary.enabled),
            "domain_packs_enabled": bool(defaults.domain_packs.enabled),
            "provider_retry_mode": defaults.provider_retry_mode,
            "dream_annotate_line_ages": bool(defaults.dream.annotate_line_ages),
        },
        "learning": {
            "background_review_enabled": bool(defaults.learning.background_review.enabled),
            "curator_enabled": bool(defaults.learning.curator.enabled),
        },
        "evolution": {
            "mode": evolution.mode,
            "allow_manual_override": bool(evolution.allow_manual_override),
            "dry_run": bool(evolution.dry_run),
            "outcome_archive_enabled": bool(evolution.outcome_archive_enabled),
            "dependency_stale_cleanup_enabled": bool(evolution.dependency_stale_cleanup_enabled),
            "auto_verify_workflows": bool(evolution.auto_verify_workflows),
            "skill_candidates_enabled": bool(evolution.skill_candidates_enabled),
            "feedback_calibration_enabled": bool(evolution.feedback_calibration_enabled),
            "sandbox_enabled": bool(evolution.sandbox.enabled),
            "trial_enabled": bool(evolution.trial.enabled),
            "trial_isolated_workspace": bool(evolution.trial.isolated_workspace),
            "trial_read_only_tools_only": bool(evolution.trial.read_only_tools_only),
        },
        "gateway": {
            "heartbeat_enabled": bool(config.gateway.heartbeat.enabled),
        },
        "security": {
            "pairing_enabled": bool(config.security.pairing.enabled),
            "pairing_allow_self_approve": bool(config.security.pairing.allow_self_approve),
        },
        "search": {
            "web_enabled": bool(config.tools.web.enable),
            "web_fetch_use_jina_reader": bool(config.tools.web.fetch.use_jina_reader),
            "session_search_enabled": bool(config.tools.session_search.enabled),
            "session_search_backend": config.tools.session_search.backend,
            "session_search_semantic_enabled": bool(config.tools.session_search.semantic_enabled),
            "session_search_rebuild_on_start": bool(config.tools.session_search.rebuild_on_start),
            "content_read_enabled": bool(config.tools.content_read.enabled),
            "content_read_use_jina_reader": bool(config.tools.content_read.use_jina_reader),
        },
        "execution": {
            "exec_enabled": bool(config.tools.exec.enable),
            "exec_profile": config.tools.exec.profile,
            "exec_allow_unsafe_exec": bool(config.tools.exec.allow_unsafe_exec),
            "exec_shell_syntax_policy": config.tools.exec.shell_syntax_policy,
            "my_enabled": bool(config.tools.my.enable),
            "my_allow_set": bool(config.tools.my.allow_set),
            "restrict_to_workspace": bool(config.tools.restrict_to_workspace),
        },
        "media": {
            "image_generation_enabled": bool(config.tools.image_generation.enabled),
        },
        "devices": {
            "device_enabled": bool(config.tools.device.enabled),
            "device_lighting_enabled": bool(config.tools.device.lighting_enabled),
            "device_mode": config.tools.device.mode,
            "device_backend": config.tools.device.backend,
        },
        "audit": {
            "audit_mode": config.tools.audit.mode,
            "audit_security_on_policy_denial": bool(config.tools.audit.security_on_policy_denial),
        },
        "runtime": {
            "profile": config.runtime.profile,
        },
    }


def _mcp_masked_mapping(values: dict[str, str]) -> dict[str, str]:
    return {key: _MCP_SECRET_HINT for key, value in values.items() if value}


def _mcp_server_payload(name: str, server: Any) -> dict[str, Any]:
    return {
        "name": name,
        "type": server.type,
        "command": server.command,
        "args": list(server.args),
        "env": _mcp_masked_mapping(server.env),
        "url": server.url,
        "headers": _mcp_masked_mapping(server.headers),
        "tool_timeout": server.tool_timeout,
        "enabled_tools": list(server.enabled_tools),
    }


def _merge_mcp_secret_fields(data: dict[str, Any], existing: Any) -> dict[str, Any]:
    merged = dict(data)
    for field in ("env", "headers"):
        incoming = merged.get(field)
        if incoming is None:
            continue
        if not isinstance(incoming, dict):
            continue
        current = dict(getattr(existing, field, {}) or {})
        cleaned: dict[str, str] = dict(current)
        for key, value in incoming.items():
            key_str = str(key).strip()
            if not key_str:
                continue
            value_str = str(value)
            if value_str == _MCP_SECRET_HINT or value_str == "":
                if key_str in current:
                    cleaned[key_str] = current[key_str]
            else:
                cleaned[key_str] = value_str
        merged[field] = cleaned
    return merged


def _home_assistant_mcp_url(address: str) -> str | None:
    raw = address.strip()
    if "://" not in raw:
        raw = f"http://{raw}"
    parsed = urlparse(raw)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        return None
    base = f"{parsed.scheme}://{parsed.netloc}"
    if _strip_trailing_slash(parsed.path) == _HA_MCP_PATH:
        return f"{base}{_HA_MCP_PATH}"
    return f"{base}{_HA_MCP_PATH}"


def _host_exact_cidrs(hostname: str) -> list[str]:
    hosts = [hostname]
    if hostname.lower() == "localhost":
        hosts = ["127.0.0.1", "::1"]
    cidrs: list[str] = []
    seen: set[str] = set()

    def add_addr(raw_addr: str) -> None:
        try:
            addr = ipaddress.ip_address(raw_addr)
        except ValueError:
            return
        cidr = f"{addr}/{addr.max_prefixlen}"
        if cidr not in seen:
            seen.add(cidr)
            cidrs.append(cidr)

    for host in hosts:
        before_count = len(cidrs)
        add_addr(host)
        if len(cidrs) > before_count:
            continue
        try:
            infos = socket.getaddrinfo(host, None, socket.AF_UNSPEC, socket.SOCK_STREAM)
        except socket.gaierror:
            continue
        for info in infos:
            add_addr(info[4][0])
    return cidrs


def _parse_inbound_payload(raw: str) -> str | None:
    """Parse a client frame into text; return None for empty or unrecognized content."""
    text = raw.strip()
    if not text:
        return None
    if text.startswith("{"):
        try:
            data = json.loads(text)
        except json.JSONDecodeError:
            return text
        if isinstance(data, dict):
            for key in ("content", "text", "message"):
                value = data.get(key)
                if isinstance(value, str) and value.strip():
                    return value
            return None
        return None
    return text


# Accept UUIDs and short scoped keys like "unified:default". Keeps the capability
# namespace small enough to rule out path traversal / quote injection tricks.
_CHAT_ID_RE = re.compile(r"^[A-Za-z0-9_:-]{1,64}$")


def _is_valid_chat_id(value: Any) -> bool:
    return isinstance(value, str) and _CHAT_ID_RE.match(value) is not None


def _parse_envelope(raw: str) -> dict[str, Any] | None:
    """Return a typed envelope dict if the frame is a new-style JSON envelope, else None.

    A frame qualifies when it parses as a JSON object with a string ``type`` field.
    Legacy frames (plain text, or ``{"content": ...}`` without ``type``) return None;
    callers should fall back to :func:`_parse_inbound_payload` for those.
    """
    text = raw.strip()
    if not text.startswith("{"):
        return None
    try:
        data = json.loads(text)
    except json.JSONDecodeError:
        return None
    if not isinstance(data, dict):
        return None
    t = data.get("type")
    if not isinstance(t, str):
        return None
    return data


# Per-message media limits. The server-side guard is a touch looser than the
# client's ``Worker`` normalization target (6 MB) — tolerate client slop, but
# still cap total ingress at ``_MAX_IMAGES_PER_MESSAGE * _MAX_IMAGE_BYTES``
# which fits comfortably inside ``max_message_bytes``.
_MAX_IMAGES_PER_MESSAGE = 4
_MAX_IMAGE_BYTES = 8 * 1024 * 1024
_MAX_VIDEOS_PER_MESSAGE = 1
_MAX_VIDEO_BYTES = 20 * 1024 * 1024

# Image MIME whitelist — matches the Composer's ``accept`` list. SVG is
# explicitly excluded to avoid the XSS surface inside embedded scripts.
_IMAGE_MIME_ALLOWED: frozenset[str] = frozenset({
    "image/png",
    "image/jpeg",
    "image/webp",
    "image/gif",
})

_VIDEO_MIME_ALLOWED: frozenset[str] = frozenset({
    "video/mp4",
    "video/webm",
    "video/quicktime",
})

_UPLOAD_MIME_ALLOWED: frozenset[str] = _IMAGE_MIME_ALLOWED | _VIDEO_MIME_ALLOWED

_DATA_URL_MIME_RE = re.compile(r"^data:([^;]+);base64,", re.DOTALL)


def _extract_data_url_mime(url: str) -> str | None:
    """Return the MIME type of a ``data:<mime>;base64,...`` URL, else ``None``."""
    if not isinstance(url, str):
        return None
    m = _DATA_URL_MIME_RE.match(url)
    if not m:
        return None
    return m.group(1).strip().lower() or None


_LOCALHOSTS = frozenset({"127.0.0.1", "::1", "localhost"})

# Matches the legacy chat-id pattern but allows file-system-safe stems too,
# so the API can address sessions whose keys came from non-WebSocket channels.
_API_KEY_RE = re.compile(r"^[A-Za-z0-9_:.-]{1,128}$")


def _decode_api_key(raw_key: str) -> str | None:
    """Decode a percent-encoded API path segment, then validate the result."""
    key = unquote(raw_key)
    if _API_KEY_RE.match(key) is None:
        return None
    return key


def _is_localhost(connection: Any) -> bool:
    """Return True if *connection* originated from the loopback interface."""
    addr = getattr(connection, "remote_address", None)
    if not addr:
        return False
    host = addr[0] if isinstance(addr, tuple) else addr
    if not isinstance(host, str):
        return False
    # ``::ffff:127.0.0.1`` is loopback in IPv6-mapped form.
    if host.startswith("::ffff:"):
        host = host[7:]
    return host in _LOCALHOSTS


def _timestamp_ms(value: Any) -> int:
    if isinstance(value, int | float):
        return int(value * 1000 if value < 10_000_000_000 else value)
    if isinstance(value, str) and value:
        try:
            from datetime import datetime

            return int(datetime.fromisoformat(value).timestamp() * 1000)
        except ValueError:
            pass
    return int(time.time() * 1000)


def _http_response(
    body: bytes,
    *,
    status: int = 200,
    content_type: str = "text/plain; charset=utf-8",
    extra_headers: list[tuple[str, str]] | None = None,
) -> Response:
    headers = [
        ("Date", email.utils.formatdate(usegmt=True)),
        ("Connection", "close"),
        ("Content-Length", str(len(body))),
        ("Content-Type", content_type),
    ]
    if extra_headers:
        headers.extend(extra_headers)
    reason = http.HTTPStatus(status).phrase
    return Response(status, reason, Headers(headers), body)


def _http_error(status: int, message: str | None = None) -> Response:
    body = (message or http.HTTPStatus(status).phrase).encode("utf-8")
    return _http_response(body, status=status)


def _bearer_token(headers: Any) -> str | None:
    """Pull a Bearer token out of standard or query-style headers."""
    auth = headers.get("Authorization") or headers.get("authorization")
    if auth and auth.lower().startswith("bearer "):
        return auth[7:].strip() or None
    return None


def _is_websocket_upgrade(request: WsRequest) -> bool:
    """Detect an actual WS upgrade; plain HTTP GETs to the same path should fall through."""
    upgrade = request.headers.get("Upgrade") or request.headers.get("upgrade")
    connection = request.headers.get("Connection") or request.headers.get("connection")
    if not upgrade or "websocket" not in upgrade.lower():
        return False
    if not connection or "upgrade" not in connection.lower():
        return False
    return True


def _b64url_encode(data: bytes) -> str:
    """URL-safe base64 without padding — compact + friendly in URL paths."""
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode("ascii")


def _b64url_decode(s: str) -> bytes:
    """Reverse of :func:`_b64url_encode`; caller handles ``ValueError``."""
    pad = "=" * (-len(s) % 4)
    return base64.urlsafe_b64decode(s + pad)


# Allowed MIME types we actually serve from the media endpoint. Anything
# outside this set is degraded to ``application/octet-stream`` so an
# attacker who somehow gets a signed URL for an unexpected file type can't
# trick the browser into sniffing executable content.
_MEDIA_ALLOWED_MIMES: frozenset[str] = frozenset({
    "image/png",
    "image/jpeg",
    "image/webp",
    "image/gif",
    "video/mp4",
    "video/webm",
    "video/quicktime",
})


def _issue_route_secret_matches(headers: Any, configured_secret: str) -> bool:
    """Return True if the token-issue HTTP request carries credentials matching ``token_issue_secret``."""
    if not configured_secret:
        return True
    authorization = headers.get("Authorization") or headers.get("authorization")
    if authorization and authorization.lower().startswith("bearer "):
        supplied = authorization[7:].strip()
        return hmac.compare_digest(supplied, configured_secret)
    header_token = headers.get("X-OriginAgent-Auth") or headers.get("x-OriginAgent-auth")
    if not header_token:
        return False
    return hmac.compare_digest(header_token.strip(), configured_secret)


class WebSocketChannel(BaseChannel):
    """Run a local WebSocket server; forward text/JSON messages to the message bus."""

    name = "websocket"
    display_name = "WebSocket"

    def __init__(
        self,
        config: Any,
        bus: MessageBus,
        *,
        session_manager: "SessionManager | None" = None,
        static_dist_path: Path | None = None,
        runtime_model_name: Callable[[], str | None] | None = None,
    ):
        if isinstance(config, dict):
            config = WebSocketConfig.model_validate(config)
        super().__init__(config, bus)
        self.config: WebSocketConfig = config
        # chat_id -> connections subscribed to it (fan-out target).
        self._subs: dict[str, set[Any]] = {}
        # connection -> chat_ids it is subscribed to (O(1) cleanup on disconnect).
        self._conn_chats: dict[Any, set[str]] = {}
        # connection -> default chat_id for legacy frames that omit routing.
        self._conn_default: dict[Any, str] = {}
        # Single-use tokens consumed at WebSocket handshake.
        self._issued_tokens: dict[str, float] = {}
        # Multi-use tokens for the embedded webui's REST surface; checked but not consumed.
        self._api_tokens: dict[str, float] = {}
        self._stop_event: asyncio.Event | None = None
        self._server_task: asyncio.Task[None] | None = None
        self._session_manager = session_manager
        self._static_dist_path: Path | None = (
            static_dist_path.resolve() if static_dist_path is not None else None
        )
        self._runtime_model_name = runtime_model_name
        # Process-local secret used to HMAC-sign media URLs. The signed URL is
        # the capability — anyone who holds a valid URL can fetch that one
        # file, nothing else. The secret regenerates on restart so links
        # become self-expiring (callers just refresh the session list).
        self._media_secret: bytes = secrets.token_bytes(32)

    # -- Subscription bookkeeping -------------------------------------------

    def _attach(self, connection: Any, chat_id: str) -> None:
        """Idempotently subscribe *connection* to *chat_id*."""
        self._subs.setdefault(chat_id, set()).add(connection)
        self._conn_chats.setdefault(connection, set()).add(chat_id)

    def _cleanup_connection(self, connection: Any) -> None:
        """Remove *connection* from every subscription set; safe to call multiple times."""
        chat_ids = self._conn_chats.pop(connection, set())
        for cid in chat_ids:
            subs = self._subs.get(cid)
            if subs is None:
                continue
            subs.discard(connection)
            if not subs:
                self._subs.pop(cid, None)
        self._conn_default.pop(connection, None)

    async def _maybe_push_active_goal_state(self, chat_id: str) -> None:
        """Replay persisted goal state after a client subscribes."""
        if self._session_manager is None:
            return
        session = self._session_manager.get_or_create(f"websocket:{chat_id}")
        blob = goal_state_ws_blob(session.metadata)
        if blob.get("active") is True:
            await self.send_goal_state(chat_id, blob)

    async def _send_event(self, connection: Any, event: str, **fields: Any) -> None:
        """Send a control event (attached, error, ...) to a single connection."""
        payload: dict[str, Any] = {"event": event}
        payload.update(fields)
        raw = json.dumps(payload, ensure_ascii=False)
        try:
            await connection.send(raw)
        except ConnectionClosed:
            self._cleanup_connection(connection)
        except Exception as e:
            self.logger.warning("failed to send {} event: {}", event, e)

    @classmethod
    def default_config(cls) -> dict[str, Any]:
        return WebSocketConfig().model_dump(by_alias=True)

    def _expected_path(self) -> str:
        return _normalize_config_path(self.config.path)

    def _build_ssl_context(self) -> ssl.SSLContext | None:
        cert = self.config.ssl_certfile.strip()
        key = self.config.ssl_keyfile.strip()
        if not cert and not key:
            return None
        if not cert or not key:
            raise ValueError(
                "ssl_certfile and ssl_keyfile must both be set for WSS, or both left empty"
            )
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        ctx.minimum_version = ssl.TLSVersion.TLSv1_2
        ctx.load_cert_chain(certfile=cert, keyfile=key)
        return ctx

    _MAX_ISSUED_TOKENS = 10_000

    def _purge_expired_issued_tokens(self) -> None:
        now = time.monotonic()
        for token_key, expiry in list(self._issued_tokens.items()):
            if now > expiry:
                self._issued_tokens.pop(token_key, None)

    def _take_issued_token_if_valid(self, token_value: str | None) -> bool:
        """Validate and consume one issued token (single use per connection attempt).

        Uses single-step pop to minimize the window between lookup and removal;
        safe under asyncio's single-threaded cooperative model.
        """
        if not token_value:
            return False
        self._purge_expired_issued_tokens()
        expiry = self._issued_tokens.pop(token_value, None)
        if expiry is None:
            return False
        if time.monotonic() > expiry:
            return False
        return True

    def _handle_token_issue_http(self, connection: Any, request: Any) -> Any:
        secret = self.config.token_issue_secret.strip()
        if secret:
            if not _issue_route_secret_matches(request.headers, secret):
                return connection.respond(401, "Unauthorized")
        else:
            self.logger.warning(
                "token_issue_path is set but token_issue_secret is empty; "
                "any client can obtain connection tokens — set token_issue_secret for production."
            )
        self._purge_expired_issued_tokens()
        if len(self._issued_tokens) >= self._MAX_ISSUED_TOKENS:
            self.logger.error(
                "too many outstanding issued tokens ({}), rejecting issuance",
                len(self._issued_tokens),
            )
            return _http_json_response({"error": "too many outstanding tokens"}, status=429)
        token_value = f"nbwt_{secrets.token_urlsafe(32)}"
        self._issued_tokens[token_value] = time.monotonic() + float(self.config.token_ttl_s)

        return _http_json_response(
            {"token": token_value, "expires_in": self.config.token_ttl_s}
        )

    # -- HTTP dispatch ------------------------------------------------------

    async def _dispatch_http(self, connection: Any, request: WsRequest) -> Any:
        """Route an inbound HTTP request to a handler or to the WS upgrade path."""
        got, query = _parse_request_path(request.path)

        # 1. Token issue endpoint (legacy, optional, gated by configured secret).
        if self.config.token_issue_path:
            issue_expected = _normalize_config_path(self.config.token_issue_path)
            if got == issue_expected:
                return self._handle_token_issue_http(connection, request)

        # 2. WebUI bootstrap: mints tokens for the embedded UI.
        if got == "/webui/bootstrap":
            return self._handle_webui_bootstrap(connection, request)

        # 3. REST surface for the embedded UI.
        if got == "/api/sessions":
            return self._handle_sessions_list(request)

        if got == "/api/settings":
            return self._handle_settings(request)

        if got == "/api/commands":
            return self._handle_commands(request)

        if got == "/api/self":
            return self._handle_self(request)

        if got == "/api/reviews":
            return self._handle_reviews_list(request)

        if got == "/api/skills":
            return self._handle_skills_list(request)

        if got == "/api/domains":
            return self._handle_domains_list(request)

        if got == "/api/settings/update":
            return self._handle_settings_update(request)

        if got == "/api/settings/provider/update":
            return self._handle_settings_provider_update(request)

        if got == "/api/settings/web-search/update":
            return self._handle_settings_web_search_update(request)

        if got == "/api/settings/learning/background-review/update":
            return self._handle_settings_learning_background_review_update(request)

        if got == "/api/settings/runtime/update":
            return self._handle_settings_runtime_update(request)

        if got == "/api/settings/mcp/upsert":
            return self._handle_settings_mcp_upsert(request)

        if got == "/api/settings/mcp/home-assistant/upsert":
            return self._handle_settings_mcp_home_assistant_upsert(request)

        if got == "/api/settings/mcp/delete":
            return self._handle_settings_mcp_delete(request)

        m = re.match(r"^/api/reviews/([^/]+)$", got)
        if m:
            return self._handle_review_detail(request, m.group(1))

        m = re.match(r"^/api/reviews/([^/]+)/(apply|approve|reject|defer)$", got)
        if m:
            return self._handle_review_action(request, m.group(1), m.group(2))

        m = re.match(r"^/api/skills/([^/]+)$", got)
        if m:
            return self._handle_skill_detail(request, m.group(1))

        m = re.match(r"^/api/skills/([^/]+)/(verify|activate|deprecate|reject|always)$", got)
        if m:
            return self._handle_skill_action(request, m.group(1), m.group(2))

        m = re.match(r"^/api/domains/([^/]+)$", got)
        if m:
            return self._handle_domain_detail(request, m.group(1))

        m = re.match(r"^/api/domains/([^/]+)/(upgrade|enable|disable|activate|deactivate|uninstall|eval)$", got)
        if m:
            return self._handle_domain_action(request, m.group(1), m.group(2))

        if got == "/api/domains/install":
            return self._handle_domains_install(request)

        m = re.match(r"^/api/sessions/([^/]+)/messages$", got)
        if m:
            return self._handle_session_messages(request, m.group(1))

        m = re.match(r"^/api/sessions/([^/]+)/webui-thread$", got)
        if m:
            return self._handle_webui_thread_get(request, m.group(1))

        # NOTE: websockets' HTTP parser only accepts GET, so we cannot expose a
        # true ``DELETE`` verb. The action is folded into the path instead.
        m = re.match(r"^/api/sessions/([^/]+)/delete$", got)
        if m:
            return self._handle_session_delete(request, m.group(1))

        # Signed media fetch: ``<sig>`` is an HMAC over ``<payload>``; the
        # payload decodes to a path inside :func:`get_media_dir`. See
        # :meth:`_sign_media_path` for the inverse direction used to build
        # these URLs when replaying a session.
        m = re.match(r"^/api/media/([A-Za-z0-9_-]+)/([A-Za-z0-9_-]+)$", got)
        if m:
            return self._handle_media_fetch(m.group(1), m.group(2))

        if got.startswith("/api/"):
            return _http_error(404, "not found")

        # 4. WebSocket upgrade (the channel's primary purpose). Only run the
        # handshake gate on requests that actually ask to upgrade; otherwise
        # a bare ``GET /`` from the browser would be rejected as an
        # unauthorized WS handshake instead of serving the SPA's index.html.
        expected_ws = self._expected_path()
        if got == expected_ws and _is_websocket_upgrade(request):
            client_id = _query_first(query, "client_id") or ""
            if len(client_id) > 128:
                client_id = client_id[:128]
            if not self.is_allowed(client_id):
                return connection.respond(403, "Forbidden")
            return self._authorize_websocket_handshake(connection, query)

        # 5. Static SPA serving (only if a build directory was wired in).
        if self._static_dist_path is not None:
            response = self._serve_static(got)
            if response is not None:
                return response

        return connection.respond(404, "Not Found")

    # -- HTTP route handlers ------------------------------------------------

    def _check_api_token(self, request: WsRequest) -> bool:
        """Validate a request against the API token pool (multi-use, TTL-bound)."""
        self._purge_expired_api_tokens()
        token = _bearer_token(request.headers) or _query_first(
            _parse_query(request.path), "token"
        )
        if not token:
            return False
        expiry = self._api_tokens.get(token)
        if expiry is None or time.monotonic() > expiry:
            self._api_tokens.pop(token, None)
            return False
        return True

    def _purge_expired_api_tokens(self) -> None:
        now = time.monotonic()
        for token_key, expiry in list(self._api_tokens.items()):
            if now > expiry:
                self._api_tokens.pop(token_key, None)

    def _handle_webui_bootstrap(self, connection: Any, request: Any) -> Response:
        # When a secret is configured (token_issue_secret or static token),
        # validate it regardless of source IP.  This secures deployments
        # behind a reverse proxy where all connections appear as localhost.
        secret = self.config.token_issue_secret.strip() or self.config.token.strip()
        if secret:
            if not _issue_route_secret_matches(request.headers, secret):
                return _http_error(401, "Unauthorized")
        elif not _is_localhost(connection):
            # No secret configured: only allow localhost (local dev mode).
            return _http_error(403, "webui bootstrap is localhost-only")
        # Cap outstanding tokens to avoid runaway growth from a misbehaving client.
        self._purge_expired_issued_tokens()
        self._purge_expired_api_tokens()
        if (
            len(self._issued_tokens) >= self._MAX_ISSUED_TOKENS
            or len(self._api_tokens) >= self._MAX_ISSUED_TOKENS
        ):
            return _http_response(
                json.dumps({"error": "too many outstanding tokens"}).encode("utf-8"),
                status=429,
                content_type="application/json; charset=utf-8",
            )
        token = f"nbwt_{secrets.token_urlsafe(32)}"
        expiry = time.monotonic() + float(self.config.token_ttl_s)
        # Same string registered in both pools: the WS handshake consumes one copy
        # while the REST surface keeps validating the other until TTL expiry.
        self._issued_tokens[token] = expiry
        self._api_tokens[token] = expiry
        return _http_json_response(
            {
                "token": token,
                "ws_path": self._expected_path(),
                "expires_in": self.config.token_ttl_s,
                "model_name": _resolve_bootstrap_model_name(self._runtime_model_name),
            }
        )

    def _handle_sessions_list(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        if self._session_manager is None:
            return _http_error(503, "session manager unavailable")
        sessions = self._session_manager.list_sessions()
        # The webui is only meaningful for websocket-channel chats — CLI /
        # Slack / Lark / Discord sessions can't be resumed from the browser,
        # so leaking them into the sidebar is just noise. Filter to the
        # ``websocket:`` prefix and strip absolute paths on the way out.
        cleaned = [
            {k: v for k, v in s.items() if k != "path"}
            for s in sessions
            if isinstance(s.get("key"), str) and s["key"].startswith("websocket:")
        ]
        return _http_json_response({"sessions": cleaned})

    def _settings_payload(self, *, requires_restart: bool = False) -> dict[str, Any]:
        from OriginAgent.config.loader import get_config_path, load_config
        from OriginAgent.providers.registry import PROVIDERS, find_by_name

        config = load_config()
        defaults = config.agents.defaults
        provider_name = config.get_provider_name(defaults.model) or defaults.provider
        provider = config.get_provider(defaults.model)
        selected_provider = provider_name
        if defaults.provider != "auto":
            spec = find_by_name(defaults.provider)
            provider_config = getattr(config.providers, spec.name, None) if spec else None
            if spec and (
                spec.is_oauth
                or spec.is_local
                or spec.is_direct
                or bool(provider_config and provider_config.api_key)
            ):
                selected_provider = spec.name
            elif spec and provider_name == defaults.provider:
                selected_provider = spec.name
        providers = []
        for spec in PROVIDERS:
            provider_config = getattr(config.providers, spec.name, None)
            if provider_config is None or spec.is_oauth or spec.is_local:
                continue
            providers.append(
                {
                    "name": spec.name,
                    "label": spec.label,
                    "configured": bool(provider_config.api_key),
                    "api_key_hint": _mask_secret_hint(provider_config.api_key),
                    "api_base": provider_config.api_base,
                    "default_api_base": spec.default_api_base or None,
                }
            )
        search_config = config.tools.web.search
        search_provider = (
            search_config.provider
            if search_config.provider in _WEB_SEARCH_PROVIDER_BY_NAME
            else "duckduckgo"
        )
        return {
            "agent": {
                "model": defaults.model,
                "provider": selected_provider,
                "resolved_provider": provider_name,
                "has_api_key": bool(provider and provider.api_key),
            },
            "providers": providers,
            "web_search": {
                "provider": search_provider,
                "api_key_hint": _mask_secret_hint(search_config.api_key),
                "base_url": search_config.base_url or None,
                "providers": list(_WEB_SEARCH_PROVIDER_OPTIONS),
            },
            "learning": {
                "background_review": {
                    "enabled": bool(defaults.learning.background_review.enabled),
                },
                "curator": {
                    "enabled": bool(defaults.learning.curator.enabled),
                },
            },
            "runtime_controls": _settings_runtime_controls_payload(config),
            "mcp": {
                "servers": [
                    _mcp_server_payload(name, server)
                    for name, server in sorted(config.tools.mcp_servers.items())
                ],
            },
            "runtime": {
                "config_path": str(get_config_path().expanduser()),
            },
            "requires_restart": requires_restart,
        }

    def _handle_settings(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        return _http_json_response(self._settings_payload())

    def _handle_commands(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        return _http_json_response({"commands": builtin_command_palette()})

    def _self_model_service(self):
        from OriginAgent.agent.confirmation import PendingConfirmationStore
        from OriginAgent.agent.domain_packs import DomainPackManager
        from OriginAgent.agent.self_model import SelfModelService
        from OriginAgent.config.loader import load_config

        config = load_config()
        manager = DomainPackManager(
            config.workspace_path,
            config=config.agents.defaults.domain_packs,
        )
        return SelfModelService(
            config.workspace_path,
            sessions=self._session_manager,
            confirmation_store=PendingConfirmationStore(config.workspace_path),
            audit_mode=config.tools.audit.mode,
            runtime_profile=config.runtime.profile,
            domain_pack_manager=manager,
            background_review_enabled=bool(config.agents.defaults.learning.background_review.enabled),
            curator_enabled=bool(config.agents.defaults.learning.curator.enabled),
        )

    def _handle_self(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        return _http_json_response({"self_model": self._self_model_service().build()})

    def _review_store(self):
        from OriginAgent.agent.background_review import ReviewProposalStore
        from OriginAgent.config.loader import load_config

        return ReviewProposalStore(load_config().workspace_path)

    def _handle_reviews_list(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        query = _parse_query(request.path)
        status = _query_first(query, "status")
        proposal_type = _query_first(query, "type")
        origin = _query_first(query, "origin")
        limit_raw = _query_first(query, "limit")
        try:
            limit = int(limit_raw) if limit_raw is not None else 50
        except ValueError:
            return _http_error(400, "limit must be an integer")
        store = self._review_store()
        return _http_json_response({
            "proposals": store.list_records(
                status=status,
                proposal_type=proposal_type,
                origin=origin,
                limit=limit,
            ),
            "stats": store.stats(
                status=status,
                proposal_type=proposal_type,
                origin=origin,
            ),
        })

    def _handle_review_detail(self, request: WsRequest, proposal_id: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        proposal_id = unquote(proposal_id)
        store = self._review_store()
        proposal = store.get(proposal_id)
        if proposal is None:
            return _http_error(404, "review proposal not found")
        return _http_json_response({"proposal": proposal, "stats": store.stats()})

    def _handle_review_action(self, request: WsRequest, proposal_id: str, action: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        proposal_id = unquote(proposal_id)
        query = _parse_query(request.path)
        reason = _query_first(query, "reason") or ""
        store = self._review_store()
        if action in {"apply", "approve"}:
            result = store.apply(proposal_id, reason=reason)
        elif action == "reject":
            result = store.reject(proposal_id, reason=reason)
        elif action == "defer":
            result = store.defer(proposal_id, reason=reason)
        else:
            return _http_error(400, "unknown review action")
        status = 200 if result.error != "not_found" else 404
        result_json = result.to_json()
        return _http_json_response({
            "result": result_json,
            "apply_result": result_json if action in {"apply", "approve"} else None,
            "proposal": result.proposal,
            "stats": store.stats(),
        }, status=status)

    def _skills_loader(self):
        from OriginAgent.agent.domain_packs import DomainPackManager
        from OriginAgent.agent.skills import SkillsLoader
        from OriginAgent.config.loader import load_config

        config = load_config()
        manager = DomainPackManager(
            config.workspace_path,
            config=config.agents.defaults.domain_packs,
        )
        return SkillsLoader(config.workspace_path, domain_pack_manager=manager)

    def _domain_governance_service(self):
        from OriginAgent.agent.domain_pack_governance import DomainPackGovernanceService
        from OriginAgent.agent.domain_packs import DomainPackManager
        from OriginAgent.config.loader import load_config

        config = load_config()
        manager = DomainPackManager(
            config.workspace_path,
            config=config.agents.defaults.domain_packs,
        )
        return DomainPackGovernanceService(config.workspace_path, domain_pack_manager=manager)

    def _handle_skills_list(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        query = _parse_query(request.path)
        source = (_query_first(query, "source") or "").strip()
        status = (_query_first(query, "status") or "").strip()
        limit_raw = _query_first(query, "limit")
        try:
            limit = max(1, min(int(limit_raw) if limit_raw is not None else 50, 200))
        except ValueError:
            return _http_error(400, "limit must be an integer")
        loader = self._skills_loader()
        records = loader.list_skill_records(filter_unavailable=False)
        if source:
            records = [record for record in records if str(record.get("source") or "") == source]
        if status:
            records = [
                record for record in records
                if str(record.get("lifecycle_status") or "") == status
            ]
        records = sorted(records, key=lambda item: (str(item.get("source") or ""), str(item.get("name") or "")))
        stats = loader.lifecycle.stats(loader.list_skills(filter_unavailable=False))
        return _http_json_response({"skills": records[:limit], "stats": stats})

    def _handle_skill_detail(self, request: WsRequest, skill_name: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        skill_name = unquote(skill_name)
        loader = self._skills_loader()
        record = loader.get_skill_record(skill_name)
        if record is None:
            return _http_error(404, "skill not found")
        return _http_json_response({
            "skill": record,
            "stats": loader.lifecycle.stats(loader.list_skills(filter_unavailable=False)),
        })

    def _handle_skill_action(self, request: WsRequest, skill_name: str, action: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.agent.skill_lifecycle import SkillLifecycleResult

        skill_name = unquote(skill_name)
        query = _parse_query(request.path)
        reason = _query_first(query, "reason") or ""
        loader = self._skills_loader()
        record = loader.get_skill_record(skill_name)
        if record is None:
            result = SkillLifecycleResult(
                skill_name=skill_name,
                status="missing",
                action=action,
                ok=False,
                message="Skill was not found.",
                error="not_found",
            )
        elif record.get("source") != "workspace":
            result = SkillLifecycleResult(
                skill_name=skill_name,
                status=str(record.get("lifecycle_status") or "unknown"),
                action=action,
                ok=False,
                message=str(record.get("disabled_reason") or "Only workspace skills can be changed in P9."),
                skill=record,
                error="read_only",
            )
        elif action == "always":
            enabled_raw = (_query_first(query, "enabled") or "").strip().lower()
            enabled = enabled_raw in {"1", "true", "yes", "on"}
            result = loader.lifecycle.transition(
                skill_name,
                action="always",
                enabled=enabled,
                reason=reason,
            )
        else:
            result = loader.lifecycle.transition(skill_name, action=action, reason=reason)
        status = 404 if result.error == "not_found" else 200
        return _http_json_response({
            "result": result.to_json(),
            "skill": result.skill,
            "stats": loader.lifecycle.stats(loader.list_skills(filter_unavailable=False)),
        }, status=status)

    def _handle_domains_list(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        query = _parse_query(request.path)
        source = (_query_first(query, "source") or "").strip()
        status = (_query_first(query, "status") or "").strip()
        limit_raw = _query_first(query, "limit")
        try:
            limit = max(1, min(int(limit_raw) if limit_raw is not None else 50, 200))
        except ValueError:
            return _http_error(400, "limit must be an integer")
        service = self._domain_governance_service()
        return _http_json_response({
            "domains": service.list_records(source=source or None, status=status or None, limit=limit),
            "stats": service.stats(),
        })

    def _handle_domain_detail(self, request: WsRequest, pack_id: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        pack_id = unquote(pack_id)
        service = self._domain_governance_service()
        record = service.get_record(pack_id)
        if record is None:
            return _http_error(404, "domain pack not found")
        return _http_json_response({"domain": record, "stats": service.stats()})

    def _handle_domains_install(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        query = _parse_query(request.path)
        source = _query_first(query, "source") or ""
        reason = _query_first(query, "reason") or ""
        service = self._domain_governance_service()
        result = service.install(source, reason=reason)
        status = 200 if result.ok else 400
        return _http_json_response({
            "result": result.to_json(),
            "domain": result.pack,
            "stats": service.stats(),
        }, status=status)

    def _handle_domain_action(self, request: WsRequest, pack_id: str, action: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        pack_id = unquote(pack_id)
        query = _parse_query(request.path)
        reason = _query_first(query, "reason") or ""
        source = _query_first(query, "source") or ""
        service = self._domain_governance_service()
        if action == "upgrade":
            result = service.upgrade(pack_id, source, reason=reason)
        elif action == "enable":
            result = service.set_enabled(pack_id, enabled=True, reason=reason)
        elif action == "disable":
            result = service.set_enabled(pack_id, enabled=False, reason=reason)
        elif action == "activate":
            result = service.set_active(pack_id, active=True, reason=reason)
        elif action == "deactivate":
            result = service.set_active(pack_id, active=False, reason=reason)
        elif action == "uninstall":
            result = service.uninstall(pack_id, reason=reason)
        elif action == "eval":
            result = service.eval_pack(pack_id)
        else:
            return _http_error(400, "unknown domain action")
        status = 200 if result.ok or result.error == "read_only" else 400
        if result.error == "not_found":
            status = 404
        return _http_json_response({
            "result": result.to_json(),
            "domain": result.pack,
            "stats": service.stats(),
        }, status=status)

    def _handle_settings_update(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config
        from OriginAgent.providers.registry import find_by_name

        query = _parse_query(request.path)
        config = load_config()
        defaults = config.agents.defaults
        changed = False

        model = _query_first(query, "model")
        if model is not None:
            model = model.strip()
            if not model:
                return _http_error(400, "model is required")
            if defaults.model != model:
                defaults.model = model
                changed = True

        provider = _query_first(query, "provider")
        if provider is not None:
            provider = provider.strip()
            if not provider:
                return _http_error(400, "provider is required")
            if find_by_name(provider) is None:
                return _http_error(400, "unknown provider")
            provider_config = getattr(config.providers, provider, None)
            if provider_config is None or not provider_config.api_key:
                return _http_error(400, "provider is not configured")
            if defaults.provider != provider:
                defaults.provider = provider
                changed = True

        if changed:
            save_config(config)
        # LLM provider/model changes are hot-reloaded by AgentLoop before each
        # new turn via the provider snapshot loader, so a restart is unnecessary.
        return _http_json_response(self._settings_payload(requires_restart=False))

    def _handle_settings_provider_update(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config
        from OriginAgent.providers.registry import find_by_name

        query = _parse_query(request.path)
        provider_name = (_query_first(query, "provider") or "").strip()
        if not provider_name:
            return _http_error(400, "provider is required")
        spec = find_by_name(provider_name)
        if spec is None or spec.is_oauth or spec.is_local:
            return _http_error(400, "unknown provider")

        config = load_config()
        provider_config = getattr(config.providers, spec.name, None)
        if provider_config is None:
            return _http_error(400, "unknown provider")

        changed = False
        if "api_key" in query or "apiKey" in query:
            api_key = _query_first(query, "api_key")
            if api_key is None:
                api_key = _query_first(query, "apiKey")
            api_key = (api_key or "").strip() or None
            if provider_config.api_key != api_key:
                provider_config.api_key = api_key
                changed = True

        if "api_base" in query or "apiBase" in query:
            api_base = _query_first(query, "api_base")
            if api_base is None:
                api_base = _query_first(query, "apiBase")
            api_base = (api_base or "").strip() or None
            if provider_config.api_base != api_base:
                provider_config.api_base = api_base
                changed = True

        if changed:
            save_config(config)
        # API key/base changes are picked up by the next provider snapshot refresh.
        return _http_json_response(self._settings_payload(requires_restart=False))

    def _handle_settings_web_search_update(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config

        query = _parse_query(request.path)
        provider_name = (_query_first(query, "provider") or "").strip().lower()
        provider_option = _WEB_SEARCH_PROVIDER_BY_NAME.get(provider_name)
        if provider_option is None:
            return _http_error(400, "unknown web search provider")

        config = load_config()
        search_config = config.tools.web.search
        previous_provider = search_config.provider
        changed = False

        def set_value(attr: str, value: str | None) -> None:
            nonlocal changed
            if getattr(search_config, attr) != value:
                setattr(search_config, attr, value)
                changed = True

        if search_config.provider != provider_name:
            search_config.provider = provider_name
            changed = True

        credential = provider_option["credential"]
        if credential == "none":
            set_value("api_key", "")
            set_value("base_url", "")
        elif credential == "base_url":
            base_url = _query_first(query, "base_url")
            if base_url is None:
                base_url = _query_first(query, "baseUrl")
            base_url = base_url.strip() if base_url is not None else None
            if not base_url and previous_provider == provider_name and search_config.base_url:
                base_url = search_config.base_url
            if not base_url:
                return _http_error(400, "base_url is required")
            set_value("base_url", base_url)
            set_value("api_key", "")
        else:
            api_key = _query_first(query, "api_key")
            if api_key is None:
                api_key = _query_first(query, "apiKey")
            api_key = api_key.strip() if api_key is not None else None
            if not api_key and previous_provider == provider_name and search_config.api_key:
                api_key = search_config.api_key
            if not api_key:
                return _http_error(400, "api_key is required")
            set_value("api_key", api_key)
            set_value("base_url", "")

        if changed:
            save_config(config)
        return _http_json_response(self._settings_payload(requires_restart=False))

    def _handle_settings_learning_background_review_update(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config

        query = _parse_query(request.path)
        enabled = _query_bool(query, "enabled")
        if enabled is None:
            return _http_error(400, "enabled must be true or false")

        config = load_config()
        target = config.agents.defaults.learning.background_review
        if target.enabled != enabled:
            target.enabled = enabled
            save_config(config)
        return _http_json_response(self._settings_payload(requires_restart=False))

    def _handle_settings_runtime_update(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config

        query = _parse_query(request.path)
        raw = _query_first(query, "config")
        if not raw:
            return _http_error(400, "config is required")
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return _http_error(400, "config must be valid JSON")
        if not isinstance(data, dict):
            return _http_error(400, "config must be an object")

        config = load_config()
        defaults = config.agents.defaults
        evolution = defaults.learning.evolution
        changed = False

        def set_bool(target: Any, attr: str, value: Any) -> None:
            nonlocal changed
            if not isinstance(value, bool):
                raise ValueError(f"{attr} must be true or false")
            if getattr(target, attr) != value:
                setattr(target, attr, value)
                changed = True

        def set_choice(target: Any, attr: str, value: Any, allowed: set[str]) -> None:
            nonlocal changed
            if not isinstance(value, str):
                raise ValueError(f"{attr} must be a string")
            candidate = value.strip()
            if candidate not in allowed:
                allowed_text = ", ".join(sorted(allowed))
                raise ValueError(f"{attr} must be one of: {allowed_text}")
            if getattr(target, attr) != candidate:
                setattr(target, attr, candidate)
                changed = True

        try:
            channels = data.get("channels")
            if isinstance(channels, dict):
                if "send_progress" in channels:
                    set_bool(config.channels, "send_progress", channels["send_progress"])
                if "send_tool_hints" in channels:
                    set_bool(config.channels, "send_tool_hints", channels["send_tool_hints"])
                if "show_reasoning" in channels:
                    set_bool(config.channels, "show_reasoning", channels["show_reasoning"])

            agent = data.get("agent")
            if isinstance(agent, dict):
                if "unified_session" in agent:
                    set_bool(defaults, "unified_session", agent["unified_session"])
                if "cold_archive_enabled" in agent:
                    set_bool(defaults, "cold_archive_enabled", agent["cold_archive_enabled"])
                if "allow_agent_initiated_messages" in agent:
                    set_bool(defaults, "allow_agent_initiated_messages", agent["allow_agent_initiated_messages"])
                if "auxiliary_enabled" in agent:
                    set_bool(defaults.auxiliary, "enabled", agent["auxiliary_enabled"])
                if "domain_packs_enabled" in agent:
                    set_bool(defaults.domain_packs, "enabled", agent["domain_packs_enabled"])
                if "provider_retry_mode" in agent:
                    set_choice(defaults, "provider_retry_mode", agent["provider_retry_mode"], _PROVIDER_RETRY_MODE_OPTIONS)
                if "dream_annotate_line_ages" in agent:
                    set_bool(defaults.dream, "annotate_line_ages", agent["dream_annotate_line_ages"])

            learning = data.get("learning")
            if isinstance(learning, dict):
                if "background_review_enabled" in learning:
                    set_bool(defaults.learning.background_review, "enabled", learning["background_review_enabled"])
                if "curator_enabled" in learning:
                    set_bool(defaults.learning.curator, "enabled", learning["curator_enabled"])

            evolution_cfg = data.get("evolution")
            if isinstance(evolution_cfg, dict):
                if "mode" in evolution_cfg:
                    set_choice(evolution, "mode", evolution_cfg["mode"], _EVOLUTION_MODE_OPTIONS)
                if "allow_manual_override" in evolution_cfg:
                    set_bool(evolution, "allow_manual_override", evolution_cfg["allow_manual_override"])
                if "dry_run" in evolution_cfg:
                    set_bool(evolution, "dry_run", evolution_cfg["dry_run"])
                if "outcome_archive_enabled" in evolution_cfg:
                    set_bool(evolution, "outcome_archive_enabled", evolution_cfg["outcome_archive_enabled"])
                if "dependency_stale_cleanup_enabled" in evolution_cfg:
                    set_bool(evolution, "dependency_stale_cleanup_enabled", evolution_cfg["dependency_stale_cleanup_enabled"])
                if "auto_verify_workflows" in evolution_cfg:
                    set_bool(evolution, "auto_verify_workflows", evolution_cfg["auto_verify_workflows"])
                if "skill_candidates_enabled" in evolution_cfg:
                    set_bool(evolution, "skill_candidates_enabled", evolution_cfg["skill_candidates_enabled"])
                if "feedback_calibration_enabled" in evolution_cfg:
                    set_bool(evolution, "feedback_calibration_enabled", evolution_cfg["feedback_calibration_enabled"])
                if "sandbox_enabled" in evolution_cfg:
                    set_bool(evolution.sandbox, "enabled", evolution_cfg["sandbox_enabled"])
                if "trial_enabled" in evolution_cfg:
                    set_bool(evolution.trial, "enabled", evolution_cfg["trial_enabled"])
                if "trial_isolated_workspace" in evolution_cfg:
                    set_bool(evolution.trial, "isolated_workspace", evolution_cfg["trial_isolated_workspace"])
                if "trial_read_only_tools_only" in evolution_cfg:
                    set_bool(evolution.trial, "read_only_tools_only", evolution_cfg["trial_read_only_tools_only"])

            gateway = data.get("gateway")
            if isinstance(gateway, dict) and "heartbeat_enabled" in gateway:
                set_bool(config.gateway.heartbeat, "enabled", gateway["heartbeat_enabled"])

            security = data.get("security")
            if isinstance(security, dict):
                if "pairing_enabled" in security:
                    set_bool(config.security.pairing, "enabled", security["pairing_enabled"])
                if "pairing_allow_self_approve" in security:
                    set_bool(config.security.pairing, "allow_self_approve", security["pairing_allow_self_approve"])

            search = data.get("search")
            if isinstance(search, dict):
                if "web_enabled" in search:
                    set_bool(config.tools.web, "enable", search["web_enabled"])
                if "web_fetch_use_jina_reader" in search:
                    set_bool(config.tools.web.fetch, "use_jina_reader", search["web_fetch_use_jina_reader"])
                if "session_search_enabled" in search:
                    set_bool(config.tools.session_search, "enabled", search["session_search_enabled"])
                if "session_search_backend" in search:
                    set_choice(config.tools.session_search, "backend", search["session_search_backend"], _SESSION_SEARCH_BACKEND_OPTIONS)
                if "session_search_semantic_enabled" in search:
                    set_bool(config.tools.session_search, "semantic_enabled", search["session_search_semantic_enabled"])
                if "session_search_rebuild_on_start" in search:
                    set_bool(config.tools.session_search, "rebuild_on_start", search["session_search_rebuild_on_start"])
                if "content_read_enabled" in search:
                    set_bool(config.tools.content_read, "enabled", search["content_read_enabled"])
                if "content_read_use_jina_reader" in search:
                    set_bool(config.tools.content_read, "use_jina_reader", search["content_read_use_jina_reader"])

            execution = data.get("execution")
            if isinstance(execution, dict):
                if "exec_enabled" in execution:
                    set_bool(config.tools.exec, "enable", execution["exec_enabled"])
                if "exec_profile" in execution:
                    set_choice(config.tools.exec, "profile", execution["exec_profile"], _EXEC_PROFILE_OPTIONS)
                if "exec_allow_unsafe_exec" in execution:
                    set_bool(config.tools.exec, "allow_unsafe_exec", execution["exec_allow_unsafe_exec"])
                if "exec_shell_syntax_policy" in execution:
                    set_choice(config.tools.exec, "shell_syntax_policy", execution["exec_shell_syntax_policy"], _EXEC_SHELL_SYNTAX_POLICY_OPTIONS)
                if "my_enabled" in execution:
                    set_bool(config.tools.my, "enable", execution["my_enabled"])
                if "my_allow_set" in execution:
                    set_bool(config.tools.my, "allow_set", execution["my_allow_set"])
                if "restrict_to_workspace" in execution:
                    set_bool(config.tools, "restrict_to_workspace", execution["restrict_to_workspace"])

            media = data.get("media")
            if isinstance(media, dict) and "image_generation_enabled" in media:
                set_bool(config.tools.image_generation, "enabled", media["image_generation_enabled"])

            devices = data.get("devices")
            if isinstance(devices, dict):
                if "device_enabled" in devices:
                    set_bool(config.tools.device, "enabled", devices["device_enabled"])
                if "device_lighting_enabled" in devices:
                    set_bool(config.tools.device, "lighting_enabled", devices["device_lighting_enabled"])
                if "device_mode" in devices:
                    set_choice(config.tools.device, "mode", devices["device_mode"], _DEVICE_MODE_OPTIONS)
                if "device_backend" in devices:
                    set_choice(config.tools.device, "backend", devices["device_backend"], _DEVICE_BACKEND_OPTIONS)

            audit = data.get("audit")
            if isinstance(audit, dict):
                if "audit_mode" in audit:
                    set_choice(config.tools.audit, "mode", audit["audit_mode"], _AUDIT_MODE_OPTIONS)
                if "audit_security_on_policy_denial" in audit:
                    set_bool(config.tools.audit, "security_on_policy_denial", audit["audit_security_on_policy_denial"])

            runtime = data.get("runtime")
            if isinstance(runtime, dict) and "profile" in runtime:
                set_choice(config.runtime, "profile", runtime["profile"], _RUNTIME_PROFILE_OPTIONS)
        except ValueError as exc:
            return _http_error(400, str(exc))

        if changed:
            save_config(config)
        return _http_json_response(self._settings_payload(requires_restart=True))

    def _handle_settings_mcp_upsert(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from pydantic import ValidationError

        from OriginAgent.config.loader import load_config, save_config
        from OriginAgent.config.schema import MCPServerConfig

        query = _parse_query(request.path)
        raw = _query_first(query, "config")
        if not raw:
            return _http_error(400, "config is required")
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return _http_error(400, "config must be valid JSON")
        if not isinstance(data, dict):
            return _http_error(400, "config must be an object")

        name = str(data.pop("name", "")).strip()
        if _MCP_SERVER_NAME_RE.fullmatch(name) is None:
            return _http_error(400, "invalid MCP server name")

        config = load_config()
        existing = config.tools.mcp_servers.get(name)
        if existing is not None:
            data = _merge_mcp_secret_fields(data, existing)

        try:
            server = MCPServerConfig.model_validate(data)
        except ValidationError as exc:
            return _http_error(400, f"invalid MCP server config: {exc.errors()[0]['msg']}")

        transport_type = server.type
        if not transport_type:
            transport_type = "stdio" if server.command else "streamableHttp" if server.url else None
        if transport_type == "stdio" and not server.command.strip():
            return _http_error(400, "command is required")
        if transport_type in {"sse", "streamableHttp"} and not server.url.strip():
            return _http_error(400, "url is required")
        if transport_type is None:
            return _http_error(400, "command or url is required")

        config.tools.mcp_servers[name] = server
        save_config(config)
        return _http_json_response(self._settings_payload(requires_restart=True))

    def _handle_settings_mcp_home_assistant_upsert(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config
        from OriginAgent.config.schema import MCPServerConfig

        query = _parse_query(request.path)
        name = (_query_first(query, "name") or "home_assistant").strip()
        if _MCP_SERVER_NAME_RE.fullmatch(name) is None:
            return _http_error(400, "invalid MCP server name")

        address = (_query_first(query, "address") or "").strip()
        url = _home_assistant_mcp_url(address)
        if url is None:
            return _http_error(400, "Home Assistant URL must start with http:// or https://")
        parsed = urlparse(url)
        if parsed.hostname is None:
            return _http_error(400, "Home Assistant URL is missing a hostname")

        token = (_query_first(query, "token") or "").strip()
        config = load_config()
        existing = config.tools.mcp_servers.get(name)
        existing_auth = (existing.headers.get("Authorization", "") if existing else "").strip()
        if token:
            authorization = token if token.lower().startswith("bearer ") else f"Bearer {token}"
        elif existing_auth:
            authorization = existing_auth
        else:
            return _http_error(400, "Home Assistant token is required")

        config.tools.mcp_servers[name] = MCPServerConfig(
            type="streamableHttp",
            url=url,
            headers={"Authorization": authorization},
            tool_timeout=30,
            enabled_tools=["*"],
        )

        for cidr in _host_exact_cidrs(parsed.hostname):
            if cidr not in config.tools.ssrf_whitelist:
                config.tools.ssrf_whitelist.append(cidr)

        save_config(config)
        return _http_json_response(self._settings_payload(requires_restart=True))

    def _handle_settings_mcp_delete(self, request: WsRequest) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        from OriginAgent.config.loader import load_config, save_config

        query = _parse_query(request.path)
        name = (_query_first(query, "name") or "").strip()
        if _MCP_SERVER_NAME_RE.fullmatch(name) is None:
            return _http_error(400, "invalid MCP server name")

        config = load_config()
        deleted = name in config.tools.mcp_servers
        if deleted:
            config.tools.mcp_servers.pop(name, None)
            save_config(config)
        payload = self._settings_payload(requires_restart=deleted)
        payload["deleted"] = deleted
        return _http_json_response(payload)

    @staticmethod
    def _is_webui_session_key(key: str) -> bool:
        """Return True when *key* belongs to the webui's websocket-only surface."""
        return key.startswith("websocket:")

    def _handle_session_messages(self, request: WsRequest, key: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        if self._session_manager is None:
            return _http_error(503, "session manager unavailable")
        decoded_key = _decode_api_key(key)
        if decoded_key is None:
            return _http_error(400, "invalid session key")
        # The embedded webui only understands websocket-channel sessions. Keep
        # its read surface aligned with ``/api/sessions`` instead of letting a
        # caller probe arbitrary CLI / Slack / Lark history by handcrafted URL.
        if not self._is_webui_session_key(decoded_key):
            return _http_error(404, "session not found")
        data = self._session_manager.read_session_file(decoded_key)
        if data is None:
            return _http_error(404, "session not found")
        messages = data.get("messages")
        if isinstance(messages, list):
            scrub_subagent_messages_for_channel(messages)
        # Decorate persisted user messages with signed media URLs so the
        # client can render previews. The raw on-disk ``media`` paths are
        # stripped on the way out — they leak server filesystem layout and
        # the client never needs them once it has the signed fetch URL.
        self._augment_media_urls(data)
        return _http_json_response(data)

    def _handle_webui_thread_get(self, request: WsRequest, key: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        decoded_key = _decode_api_key(key)
        if decoded_key is None:
            return _http_error(400, "invalid session key")
        if not self._is_webui_session_key(decoded_key):
            return _http_error(404, "session not found")
        data = build_webui_thread_response(
            decoded_key,
            augment_user_media=self._augment_transcript_user_media,
        )
        if data is None:
            data = self._build_webui_thread_from_session(decoded_key)
        elif isinstance(data.get("messages"), list):
            scrub_subagent_messages_for_channel(data["messages"])
        if data is None:
            return _http_error(404, "webui thread not found")
        return _http_json_response(data)

    def _build_webui_thread_from_session(self, key: str) -> dict[str, Any] | None:
        data = self._session_manager.read_session_file(key) if self._session_manager else None
        if data is None:
            return None
        messages = data.get("messages")
        if isinstance(messages, list):
            scrub_subagent_messages_for_channel(messages)
        self._augment_media_urls(data)
        if not isinstance(messages, list):
            return None
        ui_messages: list[dict[str, Any]] = []
        for idx, msg in enumerate(messages):
            if not isinstance(msg, dict):
                continue
            if msg.get("_command"):
                continue
            role = msg.get("role")
            if role not in {"user", "assistant", "tool"}:
                continue
            content = msg.get("content")
            if not isinstance(content, str):
                content = "" if content is None else str(content)
            row: dict[str, Any] = {
                "id": f"legacy-{idx}",
                "role": role,
                "content": content,
                "createdAt": _timestamp_ms(msg.get("timestamp")),
            }
            media = msg.get("media_urls")
            if isinstance(media, list) and media:
                row["media"] = [
                    {"kind": "image", "url": str(m["url"]), "name": str(m.get("name") or "")}
                    for m in media
                    if isinstance(m, dict) and m.get("url")
                ]
            if row["content"].strip() or row.get("media"):
                ui_messages.append(row)
        if not ui_messages:
            return None
        return {"schemaVersion": 3, "sessionKey": key, "messages": ui_messages}

    def _try_append_webui_transcript(self, chat_id: str, wire: dict[str, Any]) -> None:
        from OriginAgent.utils.webui_transcript import append_transcript_object

        if wire.get("_transcript_recorded"):
            return
        try:
            dup = json.loads(json.dumps(wire, ensure_ascii=False))
            dup.pop("_transcript_recorded", None)
            append_transcript_object(f"websocket:{chat_id}", dup)
        except (TypeError, ValueError, OSError) as e:
            self.logger.warning("webui transcript append failed: {}", e)

    def append_webui_transcript_event(self, chat_id: str, wire: dict[str, Any]) -> None:
        """Append a prebuilt WebUI transcript event.

        Agent command shortcuts persist to the session directly and may also
        call this hook when the active channel supports the WebUI transcript.
        Keeping the hook public avoids coupling the agent loop to websocket
        channel internals while still making command turns replayable after a
        refresh.
        """
        self._try_append_webui_transcript(chat_id, wire)

    def _augment_transcript_user_media(self, paths: list[str]) -> list[dict[str, Any]]:
        out: list[dict[str, Any]] = []
        for pstr in paths:
            path = Path(pstr)
            att = self._sign_or_stage_media_path(path)
            if att is None:
                continue
            mime, _ = mimetypes.guess_type(path.name)
            kind = "video" if mime and mime.startswith("video/") else "image"
            out.append({"kind": kind, "url": att["url"], "name": att.get("name", path.name)})
        return out

    def _augment_media_urls(self, payload: dict[str, Any]) -> None:
        """Mutate *payload* in place: each message's ``media`` path list is
        replaced by a parallel ``media_urls`` list of signed fetch URLs.

        Messages without media or with non-string path entries are left
        untouched. Paths that no longer live inside ``media_dir`` (e.g. the
        file was deleted, or the dir was relocated) are silently skipped;
        the client falls back to the historical-replay placeholder tile.
        """
        messages = payload.get("messages")
        if not isinstance(messages, list):
            return
        for msg in messages:
            if not isinstance(msg, dict):
                continue
            media = msg.get("media")
            if not isinstance(media, list) or not media:
                continue
            urls: list[dict[str, str]] = []
            for entry in media:
                if not isinstance(entry, str) or not entry:
                    continue
                signed = self._sign_media_path(Path(entry))
                if signed is None:
                    continue
                urls.append({"url": signed, "name": Path(entry).name})
            if urls:
                msg["media_urls"] = urls
            # Always drop the raw paths from the wire payload.
            msg.pop("media", None)

    def _sign_media_path(self, abs_path: Path) -> str | None:
        """Return a ``/api/media/<sig>/<payload>`` URL for *abs_path*, or
        ``None`` when the path does not resolve inside the media root.

        The URL is self-authenticating: the signature binds the payload to
        this process's ``_media_secret``, so only paths we chose to sign can
        be fetched. The returned path is relative to the server origin; the
        client joins it against the existing webui base.
        """
        try:
            media_root = get_media_dir().resolve()
            rel = abs_path.resolve().relative_to(media_root)
        except (OSError, ValueError):
            return None
        payload = _b64url_encode(rel.as_posix().encode("utf-8"))
        mac = hmac.new(
            self._media_secret, payload.encode("ascii"), hashlib.sha256
        ).digest()[:16]
        return f"/api/media/{_b64url_encode(mac)}/{payload}"

    def _sign_or_stage_media_path(self, path: Path) -> dict[str, str] | None:
        """Return a signed media URL payload for *path*.

        Persisted inbound media already lives under ``get_media_dir`` and can
        be signed directly. Outbound bot-generated files may live anywhere on
        disk; copy those into the websocket media bucket first so the browser
        can fetch them through the existing signed media route without
        exposing arbitrary filesystem paths.
        """
        signed = self._sign_media_path(path)
        if signed is not None:
            return {"url": signed, "name": path.name}
        try:
            if not path.is_file():
                return None
            media_dir = get_media_dir("websocket")
            safe_name = safe_filename(path.name) or "attachment"
            staged = media_dir / f"{uuid.uuid4().hex[:12]}-{safe_name}"
            shutil.copyfile(path, staged)
        except OSError as exc:
            self.logger.warning("failed to stage outbound media {}: {}", path, exc)
            return None
        signed = self._sign_media_path(staged)
        if signed is None:
            return None
        return {"url": signed, "name": path.name}

    def _handle_media_fetch(self, sig: str, payload: str) -> Response:
        """Serve a single media file previously signed via
        :meth:`_sign_media_path`. Validates the signature, decodes the
        payload to a relative path, and streams the file bytes with a
        long-lived immutable cache header (the URL already encodes the
        file identity, so caches can be aggressive)."""
        try:
            provided_mac = _b64url_decode(sig)
        except (ValueError, binascii.Error):
            return _http_error(401, "invalid signature")
        expected_mac = hmac.new(
            self._media_secret, payload.encode("ascii"), hashlib.sha256
        ).digest()[:16]
        if not hmac.compare_digest(expected_mac, provided_mac):
            return _http_error(401, "invalid signature")
        try:
            rel_bytes = _b64url_decode(payload)
            rel_str = rel_bytes.decode("utf-8")
        except (ValueError, binascii.Error, UnicodeDecodeError):
            return _http_error(400, "invalid payload")
        # An attacker who somehow bypassed the HMAC check would still need
        # the resolved path to escape the media root; guard defensively.
        try:
            media_root = get_media_dir().resolve()
            candidate = (media_root / rel_str).resolve()
            candidate.relative_to(media_root)
        except (OSError, ValueError):
            return _http_error(404, "not found")
        if not candidate.is_file():
            return _http_error(404, "not found")
        try:
            body = candidate.read_bytes()
        except OSError:
            return _http_error(500, "read error")
        mime, _ = mimetypes.guess_type(candidate.name)
        if mime not in _MEDIA_ALLOWED_MIMES:
            mime = "application/octet-stream"
        return _http_response(
            body,
            content_type=mime,
            extra_headers=[
                ("Cache-Control", "private, max-age=31536000, immutable"),
                # Paired with the MIME whitelist above: prevents browsers from
                # MIME-sniffing an octet-stream fallback into executable HTML.
                ("X-Content-Type-Options", "nosniff"),
            ],
        )

    def _handle_session_delete(self, request: WsRequest, key: str) -> Response:
        if not self._check_api_token(request):
            return _http_error(401, "Unauthorized")
        if self._session_manager is None:
            return _http_error(503, "session manager unavailable")
        decoded_key = _decode_api_key(key)
        if decoded_key is None:
            return _http_error(400, "invalid session key")
        # Same boundary as ``_handle_session_messages``: the webui may only
        # mutate websocket sessions, and deletion really does unlink the local
        # JSONL, so keep the blast radius narrow and explicit.
        if not self._is_webui_session_key(decoded_key):
            return _http_error(404, "session not found")
        deleted = self._session_manager.delete_session(decoded_key)
        delete_webui_thread(decoded_key)
        return _http_json_response({"deleted": bool(deleted)})

    def _serve_static(self, request_path: str) -> Response | None:
        """Resolve *request_path* against the built SPA directory; SPA fallback to index.html."""
        assert self._static_dist_path is not None
        rel = request_path.lstrip("/")
        if not rel:
            rel = "index.html"
        # Reject path-traversal attempts and absolute targets.
        if ".." in rel.split("/") or rel.startswith("/"):
            return _http_error(403, "Forbidden")
        candidate = (self._static_dist_path / rel).resolve()
        try:
            candidate.relative_to(self._static_dist_path)
        except ValueError:
            return _http_error(403, "Forbidden")
        if not candidate.is_file():
            # SPA history-mode fallback: unknown routes serve index.html so the
            # client-side router can render them.
            index = self._static_dist_path / "index.html"
            if index.is_file():
                candidate = index
            else:
                return None
        try:
            body = candidate.read_bytes()
        except OSError as e:
            self.logger.warning("static: failed to read {}: {}", candidate, e)
            return _http_error(500, "Internal Server Error")
        ctype, _ = mimetypes.guess_type(candidate.name)
        if ctype is None:
            ctype = "application/octet-stream"
        if ctype.startswith("text/") or ctype in {"application/javascript", "application/json"}:
            ctype = f"{ctype}; charset=utf-8"
        # Hash-named build assets are cache-friendly; index.html must stay fresh.
        if candidate.name == "index.html":
            cache = "no-cache"
        else:
            cache = "public, max-age=31536000, immutable"
        return _http_response(
            body,
            status=200,
            content_type=ctype,
            extra_headers=[("Cache-Control", cache)],
        )

    def _authorize_websocket_handshake(self, connection: Any, query: dict[str, list[str]]) -> Any:
        supplied = _query_first(query, "token")
        static_token = self.config.token.strip()

        if static_token:
            if supplied and hmac.compare_digest(supplied, static_token):
                return None
            if supplied and self._take_issued_token_if_valid(supplied):
                return None
            return connection.respond(401, "Unauthorized")

        if self.config.websocket_requires_token:
            if supplied and self._take_issued_token_if_valid(supplied):
                return None
            return connection.respond(401, "Unauthorized")

        if supplied:
            self._take_issued_token_if_valid(supplied)
        return None

    async def start(self) -> None:
        self._running = True
        self._stop_event = asyncio.Event()

        ssl_context = self._build_ssl_context()
        scheme = "wss" if ssl_context else "ws"

        async def process_request(
            connection: ServerConnection,
            request: WsRequest,
        ) -> Any:
            return await self._dispatch_http(connection, request)

        async def handler(connection: ServerConnection) -> None:
            await self._connection_loop(connection)

        self.logger.info(
            "WebSocket server listening on {}://{}:{}{}",
            scheme,
            self.config.host,
            self.config.port,
            self.config.path,
        )
        if self.config.token_issue_path:
            self.logger.info(
                "WebSocket token issue route: {}://{}:{}{}",
                scheme,
                self.config.host,
                self.config.port,
                _normalize_config_path(self.config.token_issue_path),
            )

        async def runner() -> None:
            async with serve(
                handler,
                self.config.host,
                self.config.port,
                process_request=process_request,
                max_size=self.config.max_message_bytes,
                ping_interval=self.config.ping_interval_s,
                ping_timeout=self.config.ping_timeout_s,
                ssl=ssl_context,
            ):
                assert self._stop_event is not None
                await self._stop_event.wait()

        self._server_task = asyncio.create_task(runner())
        await self._server_task

    async def _connection_loop(self, connection: Any) -> None:
        request = connection.request
        path_part = request.path if request else "/"
        _, query = _parse_request_path(path_part)
        client_id_raw = _query_first(query, "client_id")
        client_id = client_id_raw.strip() if client_id_raw else ""
        if not client_id:
            client_id = f"anon-{uuid.uuid4().hex[:12]}"
        elif len(client_id) > 128:
            self.logger.warning("client_id too long ({} chars), truncating", len(client_id))
            client_id = client_id[:128]

        default_chat_id = str(uuid.uuid4())

        try:
            await connection.send(
                json.dumps(
                    {
                        "event": "ready",
                        "chat_id": default_chat_id,
                        "client_id": client_id,
                    },
                    ensure_ascii=False,
                )
            )
            # Register only after ready is successfully sent to avoid out-of-order sends
            self._conn_default[connection] = default_chat_id
            self._attach(connection, default_chat_id)

            async for raw in connection:
                if isinstance(raw, bytes):
                    try:
                        raw = raw.decode("utf-8")
                    except UnicodeDecodeError:
                        self.logger.warning("ignoring non-utf8 binary frame")
                        continue

                envelope = _parse_envelope(raw)
                if envelope is not None:
                    await self._dispatch_envelope(connection, client_id, envelope)
                    continue

                content = _parse_inbound_payload(raw)
                if content is None:
                    continue
                await self._handle_message(
                    sender_id=client_id,
                    chat_id=default_chat_id,
                    content=content,
                    metadata={"remote": getattr(connection, "remote_address", None)},
                )
        except Exception as e:
            self.logger.debug("connection ended: {}", e)
        finally:
            self._cleanup_connection(connection)

    def _save_envelope_media(
        self,
        media: list[Any],
    ) -> tuple[list[str], str | None]:
        """Decode and persist ``media`` items from a ``message`` envelope.

        Returns ``(paths, None)`` on success or ``([], reason)`` on the first
        failure — the caller is expected to surface ``reason`` to the client
        and skip publishing so no half-formed message ever reaches the agent.
        On failure, any files already written to disk earlier in the same
        call are unlinked so partial ingress doesn't leak orphan files.
        ``reason`` is a short, stable token suitable for UI localization.

        Shape: ``list[{"data_url": str, "name"?: str | None}]``.
        """
        image_count = 0
        video_count = 0
        for item in media:
            mime = _extract_data_url_mime(item.get("data_url", "")) if isinstance(item, dict) else None
            if mime in _VIDEO_MIME_ALLOWED:
                video_count += 1
            elif mime in _IMAGE_MIME_ALLOWED:
                image_count += 1
        if image_count > _MAX_IMAGES_PER_MESSAGE:
            return [], "too_many_images"
        if video_count > _MAX_VIDEOS_PER_MESSAGE:
            return [], "too_many_videos"

        media_dir = get_media_dir("websocket")
        paths: list[str] = []

        def _abort(reason: str) -> tuple[list[str], str]:
            for p in paths:
                try:
                    Path(p).unlink(missing_ok=True)
                except OSError as exc:
                    self.logger.warning(
                        "failed to unlink partial media {}: {}", p, exc
                    )
            return [], reason

        for item in media:
            if not isinstance(item, dict):
                return _abort("malformed")
            data_url = item.get("data_url")
            if not isinstance(data_url, str) or not data_url:
                return _abort("malformed")
            mime = _extract_data_url_mime(data_url)
            if mime is None:
                return _abort("decode")
            if mime not in _UPLOAD_MIME_ALLOWED:
                return _abort("mime")
            is_video = mime in _VIDEO_MIME_ALLOWED
            max_bytes = _MAX_VIDEO_BYTES if is_video else _MAX_IMAGE_BYTES
            try:
                saved = save_base64_data_url(
                    data_url, media_dir, max_bytes=max_bytes,
                )
            except FileSizeExceeded:
                return _abort("size")
            except Exception as exc:
                self.logger.warning("media decode failed: {}", exc)
                return _abort("decode")
            if saved is None:
                return _abort("decode")
            paths.append(saved)
        return paths, None

    async def _dispatch_envelope(
        self,
        connection: Any,
        client_id: str,
        envelope: dict[str, Any],
    ) -> None:
        """Route one typed inbound envelope (``new_chat`` / ``attach`` / ``message``)."""
        t = envelope.get("type")
        if t == "new_chat":
            new_id = str(uuid.uuid4())
            self._attach(connection, new_id)
            await self._send_event(connection, "attached", chat_id=new_id)
            await self._maybe_push_active_goal_state(new_id)
            return
        if t == "attach":
            cid = envelope.get("chat_id")
            if not _is_valid_chat_id(cid):
                await self._send_event(connection, "error", detail="invalid chat_id")
                return
            self._attach(connection, cid)
            await self._send_event(connection, "attached", chat_id=cid)
            await self._maybe_push_active_goal_state(cid)
            return
        if t == "message":
            cid = envelope.get("chat_id")
            content = envelope.get("content")
            if not _is_valid_chat_id(cid):
                await self._send_event(connection, "error", detail="invalid chat_id")
                return
            if not isinstance(content, str):
                await self._send_event(connection, "error", detail="missing content")
                return

            raw_media = envelope.get("media")
            media_paths: list[str] = []
            if raw_media is not None:
                if not isinstance(raw_media, list):
                    await self._send_event(
                        connection, "error",
                        detail="image_rejected", reason="malformed",
                    )
                    return
                media_paths, reason = self._save_envelope_media(raw_media)
                if reason is not None:
                    await self._send_event(
                        connection, "error",
                        detail="image_rejected", reason=reason,
                    )
                    return

            # Allow image-only turns (content may be empty when media is attached).
            if not content.strip() and not media_paths:
                await self._send_event(connection, "error", detail="missing content")
                return

            # Auto-attach on first use so clients can one-shot without a separate attach.
            self._attach(connection, cid)
            metadata: dict[str, Any] = {"remote": getattr(connection, "remote_address", None)}
            if envelope.get("webui") is True:
                metadata["webui"] = True
            image_generation = envelope.get("image_generation")
            if isinstance(image_generation, dict) and image_generation.get("enabled") is True:
                aspect_ratio = image_generation.get("aspect_ratio")
                metadata["image_generation"] = {
                    "enabled": True,
                    "aspect_ratio": aspect_ratio if isinstance(aspect_ratio, str) else None,
                }
            await self._handle_message(
                sender_id=client_id,
                chat_id=cid,
                content=content,
                media=media_paths or None,
                metadata=metadata,
            )
            return
        await self._send_event(connection, "error", detail=f"unknown type: {t!r}")

    async def _handle_message(
        self,
        sender_id: str,
        chat_id: str,
        content: str,
        media: list[str] | None = None,
        metadata: dict[str, Any] | None = None,
        session_key: str | None = None,
    ) -> None:
        meta = metadata or {}
        if meta.get("webui"):
            user_obj: dict[str, Any] = {
                "event": "user",
                "chat_id": chat_id,
                "text": content,
            }
            if media:
                user_obj["media_paths"] = list(media)
            self._try_append_webui_transcript(chat_id, user_obj)
        await super()._handle_message(
            sender_id,
            chat_id,
            content,
            media,
            metadata,
            session_key,
        )

    async def stop(self) -> None:
        if not self._running:
            return
        self._running = False
        if self._stop_event:
            self._stop_event.set()
        if self._server_task:
            try:
                await self._server_task
            except Exception as e:
                self.logger.warning("server task error during shutdown: {}", e)
            self._server_task = None
        self._subs.clear()
        self._conn_chats.clear()
        self._conn_default.clear()
        self._issued_tokens.clear()
        self._api_tokens.clear()

    async def _safe_send_to(self, connection: Any, raw: str, *, label: str = "") -> None:
        """Send a raw frame to one connection, cleaning up on ConnectionClosed."""
        try:
            await connection.send(raw)
        except ConnectionClosed:
            self._cleanup_connection(connection)
            self.logger.warning("connection gone{}", label)
        except Exception:
            self.logger.exception("send failed{}", label)
            raise

    async def send(self, msg: OutboundMessage) -> None:
        if msg.metadata.get("_runtime_model_updated"):
            await self.send_runtime_model_updated(
                model_name=msg.metadata.get("model"),
                model_preset=msg.metadata.get("model_preset"),
            )
            return

        # Snapshot the subscriber set so ConnectionClosed cleanups mid-iteration are safe.
        conns = list(self._subs.get(msg.chat_id, ()))
        if not conns:
            if (
                msg.metadata.get("_progress")
                or msg.metadata.get("_turn_end")
                or msg.metadata.get("_session_updated")
                or msg.metadata.get("_goal_status")
                or msg.metadata.get("_goal_state_sync")
            ):
                self.logger.debug("no active subscribers for chat_id={}", msg.chat_id)
            else:
                self.logger.warning("no active subscribers for chat_id={}", msg.chat_id)
            return
        # Signal that the agent has fully finished processing the current turn.
        if msg.metadata.get("_goal_state_sync"):
            blob = msg.metadata.get("goal_state")
            await self.send_goal_state(msg.chat_id, blob if isinstance(blob, dict) else {"active": False})
            return
        if msg.metadata.get("_goal_status"):
            status = msg.metadata.get("goal_status")
            if status in {"running", "idle"}:
                started_at = msg.metadata.get("started_at", msg.metadata.get("goal_started_at"))
                await self.send_goal_status(
                    msg.chat_id,
                    status,
                    started_at=started_at if isinstance(started_at, int | float) else None,
                )
            return
        if msg.metadata.get("_turn_end"):
            lat = msg.metadata.get("latency_ms")
            gs = msg.metadata.get("goal_state")
            await self.send_turn_end(
                msg.chat_id,
                latency_ms=int(lat) if isinstance(lat, int | float) else None,
                goal_state=gs if isinstance(gs, dict) else None,
            )
            return
        if msg.metadata.get("_session_updated"):
            await self.send_session_updated(msg.chat_id)
            return
        text = msg.content
        if msg.buttons:
            text = _append_buttons_as_text(text, msg.buttons)
        payload: dict[str, Any] = {
            "event": "message",
            "chat_id": msg.chat_id,
            "text": text,
        }
        if msg.metadata.get("_webui_transcript_recorded"):
            payload["_transcript_recorded"] = True
        if msg.buttons:
            payload["buttons"] = msg.buttons
            payload["button_prompt"] = msg.content
        if msg.media:
            payload["media"] = msg.media
            urls: list[dict[str, str]] = []
            for entry in msg.media:
                signed = self._sign_or_stage_media_path(Path(entry))
                if signed is not None:
                    urls.append(signed)
            if urls:
                payload["media_urls"] = urls
        if msg.reply_to:
            payload["reply_to"] = msg.reply_to
        lat = msg.metadata.get("latency_ms")
        if isinstance(lat, int | float):
            payload["latency_ms"] = int(lat)
        if msg.metadata.get("_tool_events"):
            payload["tool_events"] = msg.metadata["_tool_events"]
        agent_ui = msg.metadata.get(OUTBOUND_META_AGENT_UI)
        if agent_ui is not None:
            payload["agent_ui"] = agent_ui
        # Mark intermediate agent breadcrumbs (tool-call hints, generic
        # progress strings) so WS clients can render them as subordinate
        # trace rows rather than conversational replies.
        if msg.metadata.get("_tool_hint"):
            payload["kind"] = "tool_hint"
        elif msg.metadata.get("_progress"):
            payload["kind"] = "progress"
        self._try_append_webui_transcript(msg.chat_id, payload)
        payload.pop("_transcript_recorded", None)
        raw = json.dumps(payload, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" ")

    async def send_reasoning_delta(
        self,
        chat_id: str,
        delta: str,
        metadata: dict[str, Any] | None = None,
    ) -> None:
        """Push one chunk of model reasoning for in-place WebUI rendering."""
        conns = list(self._subs.get(chat_id, ()))
        if not conns or not delta:
            return
        meta = metadata or {}
        body: dict[str, Any] = {
            "event": "reasoning_delta",
            "chat_id": chat_id,
            "text": delta,
        }
        if meta.get("_stream_id") is not None:
            body["stream_id"] = meta["_stream_id"]
        self._try_append_webui_transcript(chat_id, body)
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" reasoning ")

    async def send_reasoning_end(
        self,
        chat_id: str,
        metadata: dict[str, Any] | None = None,
    ) -> None:
        """Close the current reasoning stream segment."""
        conns = list(self._subs.get(chat_id, ()))
        if not conns:
            return
        meta = metadata or {}
        body: dict[str, Any] = {"event": "reasoning_end", "chat_id": chat_id}
        if meta.get("_stream_id") is not None:
            body["stream_id"] = meta["_stream_id"]
        self._try_append_webui_transcript(chat_id, body)
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" reasoning_end ")

    async def send_delta(
        self,
        chat_id: str,
        delta: str,
        metadata: dict[str, Any] | None = None,
    ) -> None:
        conns = list(self._subs.get(chat_id, ()))
        if not conns:
            return
        meta = metadata or {}
        if meta.get("_stream_end"):
            body: dict[str, Any] = {"event": "stream_end", "chat_id": chat_id}
        else:
            body = {
                "event": "delta",
                "chat_id": chat_id,
                "text": delta,
            }
        if meta.get("_stream_id") is not None:
            body["stream_id"] = meta["_stream_id"]
        self._try_append_webui_transcript(chat_id, body)
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" stream ")

    async def send_turn_end(
        self,
        chat_id: str,
        latency_ms: int | None = None,
        *,
        goal_state: dict[str, Any] | None = None,
    ) -> None:
        """Signal that the agent has fully finished processing the current turn."""
        conns = list(self._subs.get(chat_id, ()))
        if not conns:
            return
        body: dict[str, Any] = {"event": "turn_end", "chat_id": chat_id}
        if latency_ms is not None:
            body["latency_ms"] = int(latency_ms)
        if goal_state is not None:
            body["goal_state"] = goal_state
        self._try_append_webui_transcript(chat_id, body)
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" turn_end ")

    async def send_goal_state(self, chat_id: str, blob: dict[str, Any]) -> None:
        """Push persisted goal-state snapshot for one chat."""
        conns = list(self._subs.get(chat_id, ()))
        if not conns:
            return
        body = {"event": "goal_state", "chat_id": chat_id, "goal_state": blob}
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" goal_state ")

    async def send_goal_status(
        self,
        chat_id: str,
        status: str,
        *,
        started_at: float | None = None,
    ) -> None:
        """Push running/idle status for the current turn strip."""
        conns = list(self._subs.get(chat_id, ()))
        if not conns:
            return
        body: dict[str, Any] = {"event": "goal_status", "chat_id": chat_id, "status": status}
        if status == "running" and started_at is not None:
            body["started_at"] = started_at
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" goal_status ")

    async def send_session_updated(self, chat_id: str) -> None:
        """Notify clients that session metadata changed outside the main turn."""
        conns = list(self._subs.get(chat_id, ()))
        if not conns:
            return
        body: dict[str, Any] = {"event": "session_updated", "chat_id": chat_id}
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" session_updated ")

    async def send_runtime_model_updated(
        self,
        *,
        model_name: Any,
        model_preset: Any = None,
    ) -> None:
        """Broadcast runtime model changes to every open websocket connection."""
        conns = list(self._conn_chats)
        if not conns or not isinstance(model_name, str) or not model_name.strip():
            return
        body: dict[str, Any] = {
            "event": "runtime_model_updated",
            "model_name": model_name.strip(),
        }
        if isinstance(model_preset, str) and model_preset.strip():
            body["model_preset"] = model_preset.strip()
        raw = json.dumps(body, ensure_ascii=False)
        for connection in conns:
            await self._safe_send_to(connection, raw, label=" runtime_model_updated ")
