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
    enabled: true,
    scriptName: "Inline full script tag",
    type: "inline",
    code: " <script>window.__bbsgoFullScriptTag = true;</script> ",
  },
  {
    enabled: true,
    scriptName: "Inline external script tag",
    type: "inline",
    code: `
      <script async src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-5683711753850351"
        crossorigin="anonymous"></script>
    `,
  },
  {
    enabled: true,
    scriptName: "Inline multiple script tags",
    type: "inline",
    code: `
      <!-- Google tag (gtag.js) -->
      <script async src="https://www.googletagmanager.com/gtag/js?id=G-5B7CC7PB9Q"></script>
      <script>
        window.dataLayer = window.dataLayer || [];
        function gtag(){dataLayer.push(arguments);}
        gtag('js', new Date());
      </script>
    `,
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
  {
    key: "script-injection-2",
    type: "inline",
    code: "window.__bbsgoFullScriptTag = true;",
  },
  {
    key: "script-injection-3",
    type: "external",
    src: "https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-5683711753850351",
    async: true,
    defer: false,
    crossOrigin: "anonymous",
  },
  {
    key: "script-injection-4-0",
    type: "external",
    src: "https://www.googletagmanager.com/gtag/js?id=G-5B7CC7PB9Q",
    async: true,
    defer: false,
    crossOrigin: undefined,
  },
  {
    key: "script-injection-4-1",
    type: "inline",
    code: "window.dataLayer = window.dataLayer || [];\n        function gtag(){dataLayer.push(arguments);}\n        gtag('js', new Date());",
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
