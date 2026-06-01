"""Skills loader for agent capabilities."""

import os
import re
import shutil
from pathlib import Path
from typing import TYPE_CHECKING

import yaml

from OriginAgent.agent.metadata import read_originagent_metadata
from OriginAgent.agent.skill_lifecycle import SkillLifecycleStore

if TYPE_CHECKING:
    from OriginAgent.agent.domain_packs import DomainPackManager

# Default builtin skills directory (relative to this file)
BUILTIN_SKILLS_DIR = Path(__file__).parent.parent / "skills"

# Opening ---, YAML body (group 1), closing --- on its own line; supports CRLF.
_STRIP_SKILL_FRONTMATTER = re.compile(
    r"^---\s*\r?\n(.*?)\r?\n---\s*\r?\n?",
    re.DOTALL,
)


class SkillsLoader:
    """
    Loader for agent skills.

    Skills are markdown files (SKILL.md) that teach the agent how to use
    specific tools or perform certain tasks.
    """

    def __init__(
        self,
        workspace: Path,
        builtin_skills_dir: Path | None = None,
        disabled_skills: set[str] | None = None,
        domain_pack_manager: "DomainPackManager | None" = None,
    ):
        self.workspace = workspace
        self.workspace_skills = workspace / "skills"
        self.builtin_skills = builtin_skills_dir or BUILTIN_SKILLS_DIR
        self.disabled_skills = disabled_skills or set()
        self.domain_pack_manager = domain_pack_manager
        self.lifecycle = SkillLifecycleStore(workspace)

    def _skill_entries_from_dir(self, base: Path, source: str, *, skip_names: set[str] | None = None) -> list[dict[str, str]]:
        if not base.exists():
            return []
        entries: list[dict[str, str]] = []
        for skill_dir in base.iterdir():
            if not skill_dir.is_dir():
                continue
            skill_file = skill_dir / "SKILL.md"
            if not skill_file.exists():
                continue
            name = skill_dir.name
            if skip_names is not None and name in skip_names:
                continue
            entries.append({"name": name, "path": str(skill_file), "source": source})
        return entries

    def list_skills(self, filter_unavailable: bool = True) -> list[dict[str, str]]:
        """
        List all available skills.

        Args:
            filter_unavailable: If True, filter out skills with unmet requirements.

        Returns:
            List of skill info dicts with 'name', 'path', 'source'.
        """
        skills = self._skill_entries_from_dir(self.workspace_skills, "workspace")
        workspace_names = {entry["name"] for entry in skills}
        if self.builtin_skills and self.builtin_skills.exists():
            skills.extend(
                self._skill_entries_from_dir(self.builtin_skills, "builtin", skip_names=workspace_names)
            )
        if self.domain_pack_manager is not None:
            skills.extend(self.domain_pack_manager.active_skill_entries())

        if self.disabled_skills:
            skills = [s for s in skills if s["name"] not in self.disabled_skills]

        if filter_unavailable:
            return [skill for skill in skills if self._check_requirements(self._get_skill_meta(skill["name"]))]
        return skills

    def list_skill_records(self, filter_unavailable: bool = True) -> list[dict]:
        """List skills with lifecycle, verification, and preview metadata."""
        return self.lifecycle.list_records(self.list_skills(filter_unavailable=filter_unavailable))

    def get_skill_record(self, name: str) -> dict | None:
        """Return one enriched skill record by name."""
        entry = self._find_skill_entry(name)
        return self.lifecycle.get_record(entry)

    def load_skill(self, name: str) -> str | None:
        """
        Load a skill by name.

        Args:
            name: Skill name (directory name).

        Returns:
            Skill content or None if not found.
        """
        return self._load_skill_content(name, include_unavailable_reason=True)

    def _load_skill_content(self, name: str, *, include_unavailable_reason: bool = False) -> str | None:
        if name in self.disabled_skills:
            return None
        if name.startswith("domain:"):
            return self._load_domain_skill(name)
        if not self._is_safe_skill_name(name):
            return None

        entry = self._find_skill_entry(name)
        if entry is None:
            return None
        if entry.get("source") == "workspace":
            record = self.lifecycle.get_record(entry)
            lifecycle = str((record or {}).get("lifecycle_status") or "")
            if lifecycle == "rejected":
                if include_unavailable_reason:
                    return "Rejected skill: this workspace skill is marked rejected and cannot be loaded."
                return None
        path = Path(entry["path"])
        if path.exists():
            content = path.read_text(encoding="utf-8")
            if entry.get("source") == "workspace":
                record = self.lifecycle.get_record(entry)
                if str((record or {}).get("lifecycle_status") or "") == "deprecated":
                    return "Deprecated skill: this workspace skill is marked deprecated.\n\n" + content
            return content
        return None

    def _find_skill_entry(self, name: str) -> dict[str, str] | None:
        if name in self.disabled_skills:
            return None
        if name.startswith("domain:"):
            parsed = self._parse_domain_skill_name(name)
            if parsed is None or self.domain_pack_manager is None:
                return None
            pack_id, skill_id = parsed
            path = self.domain_pack_manager.get_active_skill_path(pack_id, skill_id)
            if path is None or not path.exists():
                return None
            return {"name": name, "path": str(path), "source": f"domain:{pack_id}"}
        if not self._is_safe_skill_name(name):
            return None
        workspace_path = self.workspace_skills / name / "SKILL.md"
        if workspace_path.exists():
            return {"name": name, "path": str(workspace_path), "source": "workspace"}
        if self.builtin_skills:
            builtin_path = self.builtin_skills / name / "SKILL.md"
            if builtin_path.exists():
                return {"name": name, "path": str(builtin_path), "source": "builtin"}
        return None

    @staticmethod
    def _is_safe_skill_name(name: str) -> bool:
        if not name or name in {".", ".."}:
            return False
        return not any(part in name for part in ("/", "\\", ".."))

    def _load_domain_skill(self, name: str) -> str | None:
        if self.domain_pack_manager is None:
            return None
        parsed = self._parse_domain_skill_name(name)
        if parsed is None:
            return None
        pack_id, skill_id = parsed
        path = self.domain_pack_manager.get_active_skill_path(pack_id, skill_id)
        if path is None or not path.exists():
            return None
        return path.read_text(encoding="utf-8")

    @staticmethod
    def _parse_domain_skill_name(name: str) -> tuple[str, str] | None:
        rest = name.removeprefix("domain:")
        if "/" not in rest:
            return None
        pack_id, skill_id = rest.split("/", 1)
        if not pack_id or not skill_id:
            return None
        if any(part in pack_id or part in skill_id for part in ("/", "\\", "..")):
            return None
        safe = re.compile(r"^[a-z0-9_-]+$")
        if not safe.fullmatch(pack_id) or not safe.fullmatch(skill_id):
            return None
        return pack_id, skill_id

    def load_skills_for_context(self, skill_names: list[str]) -> str:
        """
        Load specific skills for inclusion in agent context.

        Args:
            skill_names: List of skill names to load.

        Returns:
            Formatted skills content.
        """
        parts = [
            f"### Skill: {name}\n\n{self._strip_frontmatter(markdown)}"
            for name in skill_names
            if (markdown := self._load_skill_content(name))
        ]
        return "\n\n---\n\n".join(parts)

    def build_skills_summary(self, exclude: set[str] | None = None) -> str:
        """
        Build a summary of all skills (name, description, path, availability).

        This is used for progressive loading - the agent can read the full
        skill content using read_file when needed.

        Args:
            exclude: Set of skill names to omit from the summary.

        Returns:
            Markdown-formatted skills summary.
        """
        all_skills = self.list_skill_records(filter_unavailable=False)
        summary_records = [
            record for record in all_skills if self._include_in_skills_summary(record)
        ]
        if not summary_records:
            return ""

        lines: list[str] = []
        for entry in summary_records:
            skill_name = entry["name"]
            if exclude and skill_name in exclude:
                continue
            meta = self._get_skill_meta(skill_name)
            available = self._check_requirements(meta)
            desc = self._get_skill_description(skill_name)
            lifecycle = entry.get("lifecycle_status")
            verification = entry.get("verification_status")
            status = ""
            if entry.get("source") == "workspace" and lifecycle == "active" and verification == "verified":
                status = " (active, verified)"
            if available:
                lines.append(f"- **{skill_name}** — {desc}{status}  `{entry['path']}`")
            else:
                missing = self._get_missing_requirements(meta)
                suffix = f" (unavailable: {missing})" if missing else " (unavailable)"
                lines.append(f"- **{skill_name}** — {desc}{suffix}  `{entry['path']}`")
        return "\n".join(lines)

    def _include_in_skills_summary(self, record: dict) -> bool:
        if record.get("source") != "workspace":
            return True
        lifecycle = str(record.get("lifecycle_status") or "")
        verification = str(record.get("verification_status") or "")
        if lifecycle in {"proposed", "deprecated", "rejected"}:
            return False
        if verification == "unverified":
            return False
        return True

    def _get_missing_requirements(self, skill_meta: dict) -> str:
        """Get a description of missing requirements."""
        requires = skill_meta.get("requires", {})
        required_bins = requires.get("bins", [])
        required_env_vars = requires.get("env", [])
        return ", ".join(
            [f"CLI: {command_name}" for command_name in required_bins if not shutil.which(command_name)]
            + [f"ENV: {env_name}" for env_name in required_env_vars if not os.environ.get(env_name)]
        )

    def _get_skill_description(self, name: str) -> str:
        """Get the description of a skill from its frontmatter."""
        meta = self.get_skill_metadata(name)
        if meta and meta.get("description"):
            return meta["description"]
        return name  # Fallback to skill name

    def _strip_frontmatter(self, content: str) -> str:
        """Remove YAML frontmatter from markdown content."""
        if not content.startswith("---"):
            return content
        match = _STRIP_SKILL_FRONTMATTER.match(content)
        if match:
            return content[match.end():].strip()
        return content

    def _parse_OriginAgent_metadata(self, raw: object) -> dict:
        """Extract OriginAgent metadata from a frontmatter field.

        Legacy OpenClaw metadata is handled by the shared metadata helper.
        """
        return read_originagent_metadata(raw)

    def _check_requirements(self, skill_meta: dict) -> bool:
        """Check if skill requirements are met (bins, env vars)."""
        requires = skill_meta.get("requires", {})
        required_bins = requires.get("bins", [])
        required_env_vars = requires.get("env", [])
        return all(shutil.which(cmd) for cmd in required_bins) and all(
            os.environ.get(var) for var in required_env_vars
        )

    def _get_skill_meta(self, name: str) -> dict:
        """Get OriginAgent metadata for a skill (cached in frontmatter)."""
        raw_meta = self.get_skill_metadata(name) or {}
        return self._parse_OriginAgent_metadata(raw_meta.get("metadata"))

    def get_always_skills(self) -> list[str]:
        """Get skills marked as always=true that meet requirements."""
        return [
            record["name"]
            for record in self.list_skill_records(filter_unavailable=True)
            if record.get("effective_always")
        ]

    def get_skill_metadata(self, name: str) -> dict | None:
        """
        Get metadata from a skill's frontmatter.

        Args:
            name: Skill name.

        Returns:
            Metadata dict or None.
        """
        entry = self._find_skill_entry(name)
        if entry is None:
            return None
        try:
            content = Path(entry["path"]).read_text(encoding="utf-8")
        except OSError:
            return None
        if not content.startswith("---"):
            return None
        match = _STRIP_SKILL_FRONTMATTER.match(content)
        if not match:
            return None
        try:
            parsed = yaml.safe_load(match.group(1))
        except yaml.YAMLError:
            return None
        if not isinstance(parsed, dict):
            return None
        # yaml.safe_load returns native types (int, bool, list, etc.);
        # keep values as-is so downstream consumers get correct types.
        metadata: dict[str, object] = {}
        for key, value in parsed.items():
            metadata[str(key)] = value
        return metadata
