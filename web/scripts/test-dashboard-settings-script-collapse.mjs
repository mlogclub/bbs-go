import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const routeSource = readFileSync(
  resolve(webRoot, "app/routes/dashboard.settings.tsx"),
  "utf8"
)

const normalizeScriptsMatch = routeSource.match(
  /function normalizeScripts\(items: unknown\): ScriptInjection\[\] \{[\s\S]*?\n\}/
)
const scriptSettingsMatch = routeSource.match(
  /function ScriptSettings\([\s\S]*?\nfunction PageSettings\(/
)

assert.ok(
  normalizeScriptsMatch,
  "dashboard settings should normalize script injection rows"
)
assert.ok(scriptSettingsMatch, "dashboard settings should render script rows")

assert.doesNotMatch(
  normalizeScriptsMatch[0],
  /collapsed:/,
  "script injection normalization should keep persisted config separate from local collapsed state"
)

assert.match(
  scriptSettingsMatch[0],
  /useState<boolean\[\]>\(\[\]\)/,
  "script settings should store collapsed rows in separate local UI state"
)

assert.doesNotMatch(
  scriptSettingsMatch[0],
  /onChange\(scriptInjections\)/,
  "saving script settings should not overwrite local UI state with sanitized payload"
)

console.log("dashboard settings script collapse tests passed")
