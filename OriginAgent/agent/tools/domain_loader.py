"""Domain pack tool loading."""

from __future__ import annotations

import importlib.util
import sys
from typing import Any, Callable

from loguru import logger

from OriginAgent.agent.domain_packs import (
    DomainPack,
    DomainPackManager,
    DomainToolDeclaration,
)
from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.tools.registry import ToolRegistry
from OriginAgent.security.capabilities import CapabilitySnapshot


class DomainToolLoader:
    """Load explicitly declared tools from active domain packs."""

    def __init__(
        self,
        manager: DomainPackManager,
        *,
        evolution_capability_resolver: Callable[[str], CapabilitySnapshot | None] | None = None,
    ):
        self.manager = manager
        self._evolution_capability_resolver = evolution_capability_resolver

    def load(self, ctx: Any, registry: ToolRegistry) -> list[str]:
        registered: list[str] = []
        self.manager.clear_domain_tool_runtime()
        for pack, declaration in self.manager.active_tool_declarations():
            status = self._load_one(pack, declaration, ctx, registry)
            if status:
                registered.append(status)
        return registered

    def _load_one(
        self,
        pack: DomainPack,
        declaration: DomainToolDeclaration,
        ctx: Any,
        registry: ToolRegistry,
    ) -> str | None:
        if not declaration.available:
            self._skip(pack, declaration, declaration.unavailable_reason or "invalid declaration")
            return None
        if declaration.module_path is None:
            self._skip(pack, declaration, "missing module path")
            return None
        if registry.has(declaration.id):
            self._skip(pack, declaration, f"tool name {declaration.id} is already registered")
            return None

        tool_cls = self._load_class(pack, declaration)
        if tool_cls is None:
            return None

        try:
            if not tool_cls.enabled(ctx):
                self._skip(pack, declaration, "disabled by tool")
                return None
            tool = tool_cls.create(ctx)
        except Exception as exc:
            logger.exception(
                "Domain tool {} from pack {} failed during construction",
                declaration.id,
                pack.id,
            )
            self._skip(pack, declaration, f"construction failed: {exc}")
            return None

        if not isinstance(tool, Tool):
            self._skip(pack, declaration, "create() did not return a Tool")
            return None
        if tool.name != declaration.id:
            self._skip(
                pack,
                declaration,
                f"tool name mismatch: expected {declaration.id}, got {tool.name}",
            )
            return None
        if not tool.read_only and not declaration.permissions:
            self._skip(pack, declaration, "non-read-only domain tools require permissions")
            return None
        if registry.has(tool.name):
            self._skip(pack, declaration, f"tool name {tool.name} is already registered")
            return None

        setattr(tool, "_domain_pack_id", pack.id)
        setattr(tool, "_domain_tool_permissions", declaration.permissions)
        setattr(tool, "_domain_tool_audit", declaration.audit)
        if self._evolution_capability_resolver is not None:
            snapshot = self._evolution_capability_resolver(pack.id)
            if snapshot is not None:
                setattr(tool, "_evolution_capability_snapshot", snapshot)
        registry.register(tool)
        self.manager.record_domain_tool_runtime(pack.id, declaration.id, "registered")
        return declaration.id

    def _load_class(
        self,
        pack: DomainPack,
        declaration: DomainToolDeclaration,
    ) -> type[Tool] | None:
        module_path = declaration.module_path
        assert module_path is not None
        try:
            module_file = module_path.resolve()
            module_file.relative_to(pack.path.resolve())
        except (OSError, ValueError):
            self._skip(pack, declaration, "module path escapes domain pack")
            return None

        module_name = _module_name(pack.id, declaration)
        try:
            spec = importlib.util.spec_from_file_location(module_name, module_file)
            if spec is None or spec.loader is None:
                self._skip(pack, declaration, "could not create module spec")
                return None
            module = importlib.util.module_from_spec(spec)
            sys.modules[module_name] = module
            spec.loader.exec_module(module)
        except Exception as exc:
            logger.exception(
                "Failed to import domain tool module {} from pack {}",
                declaration.module,
                pack.id,
            )
            self._skip(pack, declaration, f"import failed: {exc}")
            return None

        attr = getattr(module, declaration.class_name, None)
        if not isinstance(attr, type):
            self._skip(pack, declaration, f"class {declaration.class_name} not found")
            return None
        if not issubclass(attr, Tool) or attr is Tool:
            self._skip(pack, declaration, f"class {declaration.class_name} is not a Tool")
            return None
        if getattr(attr, "__abstractmethods__", None):
            self._skip(pack, declaration, f"class {declaration.class_name} is abstract")
            return None
        if not getattr(attr, "_plugin_discoverable", True):
            self._skip(pack, declaration, f"class {declaration.class_name} is not discoverable")
            return None
        return attr

    def _skip(self, pack: DomainPack, declaration: DomainToolDeclaration, reason: str) -> None:
        logger.warning("Domain tool {} from pack {} skipped: {}", declaration.id, pack.id, reason)
        self.manager.record_domain_tool_runtime(pack.id, declaration.id, "skipped", reason)


def _module_name(pack_id: str, declaration: DomainToolDeclaration) -> str:
    safe_pack = pack_id.replace("-", "_")
    safe_tool = declaration.id.replace("-", "_")
    return f"_originagent_domain_pack_{safe_pack}_{safe_tool}"
