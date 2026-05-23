import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

import { getRenderableScriptInjections } from "../lib/script-injections.ts"

const injections = getRenderableScriptInjections([
  {
    enabled: true,
    scriptName: "Analytics",
    type: "external",
    src: " https://example.com/analytics.js ",
    async: true,
    defer: false,
    crossorigin: " anonymous ",
  },
  {
    enabled: true,
    scriptName: "Inline marker",
    type: "inline",
    code: " window.__bbsgoInjected = true; ",
  },
  {
    enabled: false,
    scriptName: "Disabled",
    type: "inline",
    code: "window.disabled = true",
  },
  {
    enabled: true,
    scriptName: "Empty external",
    type: "external",
    src: " ",
  },
])

assert.deepEqual(injections, [
  {
    key: "script-injection-0",
    type: "external",
    src: "https://example.com/analytics.js",
    async: true,
    defer: false,
    crossOrigin: "anonymous",
  },
  {
    key: "script-injection-1",
    type: "inline",
    code: "window.__bbsgoInjected = true;",
  },
])

const rootSource = readFileSync(
  resolve(import.meta.dirname, "../app/root.tsx"),
  "utf8"
)

assert.match(
  rootSource,
  /getRenderableScriptInjections/,
  "root layout should read configured script injections"
)
assert.match(
  rootSource,
  /dangerouslySetInnerHTML/,
  "root layout should render inline script injections into head"
)
assert.match(
  rootSource,
  /RuntimeScriptInjections/,
  "root app should inject configured scripts after SPA client hydration"
)

console.log("script injection tests passed")
