---
name: lighting-control
description: Use OriginAgent Core lighting tools with exact names and conservative target handling.
---

# Lighting Control

Use this skill when the active task is about controlling configured smart
lighting devices through OriginAgent Core.

## Available Tool Names

Use only the actual tool names exposed by Core:

- `originagent_device_lighting_set_power`
- `originagent_device_lighting_set_brightness`
- `originagent_device_lighting_set_color_temperature`

Do not invent shorter aliases or simplified tool names.

## Operating Rules

- If the target light, room, group, or device is ambiguous, inspect available
  state when possible or ask a concise clarification.
- Treat tool responses as authoritative. If a call is pending, denied, failed,
  or dry-run only, report that state instead of saying the light changed.
- Keep changes narrow: do not bundle unrelated lights, scenes, or routines into
  a request unless the user asked for that scope.
- For repeated or scheduled behavior, treat it as automation design unless an
  explicit automation tool is available.
