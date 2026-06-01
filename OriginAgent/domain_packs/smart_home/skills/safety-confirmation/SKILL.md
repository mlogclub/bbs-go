---
name: safety-confirmation
description: Cooperate with OriginAgent's system safety layer for physical device actions.
---

# Safety Confirmation

This skill informs the model about safety boundaries for smart home actions. It
does not replace OriginAgent's system safety layer. Final confirmation,
permission, audit, presence, and policy decisions are enforced by Core runtime
components such as the safety gate, confirmation manager, permission resolver,
and tool registry.

## How To Cooperate

- Before a risky physical action, restate the target, action, and important
  consequence in plain language.
- If a tool returns `pending`, explain that the system is waiting for
  confirmation.
- If a tool returns `denied`, `blocked`, or a permission error, report the
  denial and avoid retrying the same action without new user input.
- If the backend is unavailable or state is uncertain, report uncertainty
  instead of guessing.

## Actions Likely To Need Confirmation

- Locks, alarms, cameras, gas, and security modes.
- Ovens, heaters, high-power appliances, or HVAC extremes.
- Destructive automation edits, broad scenes, or actions affecting people.
- Any physical action with an unclear target or unclear scope.
