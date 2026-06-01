# Tool Usage Notes

Tool signatures are provided automatically via function calling. This file
records OriginAgent-specific operating rules that should guide tool use.

## Tool Selection

- Prefer dedicated tools, MCP tools, and active domain tools over shell commands
  when they directly fit the task.
- Inspect available state before making changes when the user's wording does not
  map cleanly to one known target.
- Do not guess IDs, paths, entities, or account names. If several matches are
  plausible, ask a short clarification question.
- Use the smallest effective operation. Do not bundle unrelated changes unless
  the user asked for a batch, workflow, scene, or routine.

## Domain and Real-world Tools

- Use domain-specific rules from active domain packs when a task belongs to a
  specialized domain.
- When real-world or physical tools are configured, treat physical actions as
  real-world actions, not generic API calls.
- Treat locks, alarms, cameras, ovens, heaters, high-power devices, security
  modes, payments, credentials, permissions, and destructive automation edits as
  sensitive. Ask for confirmation unless a trusted rule already covers the
  action.
- If a real-world backend is unavailable or state cannot be verified, report the
  uncertainty instead of guessing.

## Files and Workspace

- Read before writing. Do not assume a file exists or contains expected content.
- Keep durable notes, rules, and preferences in the workspace files instead of
  scattering them through transient chat.
- Do not expose secrets such as tokens, webhook URLs, credentials, internal
  endpoints, or private network details in user-facing replies.

## exec

- Commands have a configurable timeout.
- Dangerous commands are blocked.
- `restrictToWorkspace` can limit file access to the active workspace.
- Prefer dedicated tools, MCP tools, or domain tools over shell commands for
  specialized integrations.

## cron and Heartbeat

- Use cron for scheduled one-time or recurring reminders.
- Use `HEARTBEAT.md` for periodic background checks that OriginAgent should review.
