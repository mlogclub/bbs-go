Maintain non-MEMORY files based on the fact proposal result below.

Allowed paths:
- SOUL.md
- USER.md

Forbidden paths:
- memory/MEMORY.md
- memory/facts.jsonl
- memory/history.jsonl
- memory/.cursor
- memory/.dream_cursor

MEMORY.md is generated from memory/facts.jsonl. Do not edit it directly.
facts.jsonl is written by the validator. Do not edit it directly.

## Editing rules
- Edit SOUL.md or USER.md only for durable identity, behavior, or user-profile corrections.
- Do not create skills directly. Reusable workflow knowledge must go through background review skill proposals.
- If there is no non-MEMORY work to do, stop without calling tools.
- Do not guess paths.
- Use surgical edits only; never rewrite entire files.

## Skill creation rules
- Do not write `skills/<name>/SKILL.md` in Dream.
- If reusable workflow knowledge should become a skill, leave it for the controlled background review proposal flow.
