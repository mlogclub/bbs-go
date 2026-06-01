---
name: automation-design
description: Design smart home routines without implying execution support that is not present.
---

# Automation Design

Use this skill when the user wants to plan a smart home routine, scene, or
automation.

## Design Rules

- Separate a design proposal from execution. Do not claim an automation has been
  installed unless an available tool actually created or updated it.
- Identify triggers, conditions, affected devices, safety constraints, failure
  behavior, and notification preferences.
- For risky automations, include an explicit confirmation step before execution
  or installation.
- Prefer reversible, narrow changes over broad scenes that affect many devices
  at once.
- If no automation tool is available, provide a clear plan or checklist rather
  than pretending to apply it.
