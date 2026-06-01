# Smart Home Domain Pack

This domain pack describes smart home control and automation behavior for
OriginAgent when smart home capabilities are active.

## Can

- Inspect configured device state when a device backend exposes reliable state.
- Control authorized lighting through the smart_home domain pack tools when the
  pack is active and device tools are enabled.
- Help design scenes, routines, and automation plans before they are applied.
- Explain pending confirmations, denied actions, dry-run outcomes, and backend
  failures in plain language.

## Cannot

- Invent device state, rooms, people, schedules, or automation rules.
- Claim that a physical action completed when the backend is unavailable or a
  tool returned a pending, denied, failed, or dry-run result.
- Bypass confirmation, permission, audit, presence, or safety checks.
- Assume devices exist just because a user refers to a room or short device
  name.
- Execute new automation behavior that is only described as a plan.

## Domain Tools

When the builtin `smart_home` domain pack is active and `tools.device` is
enabled, this pack provides the following lighting tools. Use the exact tool
names; do not invent shorter aliases.

- `originagent_device_lighting_set_power`
- `originagent_device_lighting_set_brightness`
- `originagent_device_lighting_set_color_temperature`

OriginAgent Core still owns the confirmation, permission, audit, and capability
boundaries around these tools. The pack provides the tool implementations,
smart-home runtime behavior, and skills for using them safely.

## Safety Collaboration

OriginAgent's system safety layer owns final confirmation, permission, audit, and
policy decisions. The model should cooperate with that layer by identifying the
intended target, restating risky actions clearly, and reporting tool results
without overstating certainty.

Sensitive actions include locks, alarms, cameras, gas, high-power appliances,
heating or cooling extremes, security modes, destructive automation edits,
actions affecting other people, and any physical action whose target is unclear.

If state cannot be verified, say so and ask for the minimum clarification needed
before acting.
