"""Tool discovery and registration via package scanning."""
from __future__ import annotations

import importlib
import pkgutil
from importlib.metadata import entry_points
from typing import Any

from loguru import logger

from OriginAgent.agent.tools.base import Tool
from OriginAgent.agent.tools.registry import ToolRegistry

_SKIP_MODULES = frozenset({
    "base", "schema", "registry", "context", "loader", "config",
    "file_state", "sandbox", "mcp", "__init__", "runtime_state",
})


class ToolLoader:
    def __init__(
        self,
        package: Any = None,
        *,
        test_classes: list[type[Tool]] | None = None,
        include_builtin: bool = False,
    ):
        if package is None:
            import OriginAgent.agent.tools as _pkg
            package = _pkg
        self._package = package
        self._test_classes = test_classes
        self._include_builtin = include_builtin
        self._discovered: list[type[Tool]] | None = None
        self._plugins: dict[str, type[Tool]] | None = None

    def discover(self) -> list[type[Tool]]:
        if self._test_classes is not None:
            return list(self._test_classes)
        if not self._include_builtin:
            return []
        if self._discovered is not None:
            return self._discovered
        seen: set[int] = set()
        results: list[type[Tool]] = []
        for _importer, module_name, _ispkg in pkgutil.iter_modules(self._package.__path__):
            if module_name.startswith("_") or module_name in _SKIP_MODULES:
                continue
            try:
                module = importlib.import_module(f".{module_name}", self._package.__name__)
            except Exception:
                logger.exception("Failed to import tool module: %s", module_name)
                continue
            for attr_name in dir(module):
                attr = getattr(module, attr_name)
                if (
                    isinstance(attr, type)
                    and issubclass(attr, Tool)
                    and attr is not Tool
                    and not attr_name.startswith("_")
                    and not getattr(attr, "__abstractmethods__", None)
                    and getattr(attr, "_plugin_discoverable", True)
                    and id(attr) not in seen
                ):
                    seen.add(id(attr))
                    results.append(attr)
        results.sort(key=lambda cls: cls.__name__)
        self._discovered = results
        return results

    def _discover_plugins(self) -> dict[str, type[Tool]]:
        """Discover external tool plugins registered via entry_points."""
        if self._plugins is not None:
            return self._plugins
        plugins: dict[str, type[Tool]] = {}
        try:
            eps = entry_points(group="originagent.tools")
        except Exception:
            return plugins
        for ep in eps:
            try:
                cls = ep.load()
                if (
                    isinstance(cls, type)
                    and issubclass(cls, Tool)
                    and not getattr(cls, "__abstractmethods__", None)
                    and getattr(cls, "_plugin_discoverable", True)
                ):
                    plugins[ep.name] = cls
            except Exception:
                logger.exception("Failed to load tool plugin: %s", ep.name)
        self._plugins = plugins
        return plugins

    def load(self, ctx: Any, registry: ToolRegistry, *, scope: str = "core") -> list[str]:
        registered: list[str] = []
        core_names = set(getattr(registry, "_tools", {}).keys())
        sources = [(self.discover(), False), (self._discover_plugins().values(), True)]
        for source, is_plugin_source in sources:
            for tool_cls in source:
                cls_label = tool_cls.__name__
                try:
                    if scope not in getattr(tool_cls, "_scopes", {"core"}):
                        continue
                    if not tool_cls.enabled(ctx):
                        continue
                    tool = tool_cls.create(ctx)
                    if registry.has(tool.name):
                        if is_plugin_source and tool.name in core_names:
                            logger.warning(
                                "Plugin {} skipped: conflicts with OriginAgent core tool {}",
                                cls_label, tool.name,
                            )
                            continue
                        if is_plugin_source:
                            logger.warning(
                                "Plugin {} skipped: tool name {} is already registered",
                                cls_label, tool.name,
                            )
                            continue
                        logger.warning(
                            "Tool name collision: {} from {} skipped",
                            tool.name, cls_label,
                        )
                        continue
                    registry.register(tool)
                    registered.append(tool.name)
                except Exception:
                    logger.exception("Failed to register tool: {}", cls_label)
        return registered
