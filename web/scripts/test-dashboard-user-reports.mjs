import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const routeSource = readFileSync(
  resolve(webRoot, "app/routes/dashboard.user-reports.tsx"),
  "utf8"
)
const zhMessages = readFileSync(
  resolve(webRoot, "lib/i18n/messages/zh-CN.ts"),
  "utf8"
)
const enMessages = readFileSync(
  resolve(webRoot, "lib/i18n/messages/en-US.ts"),
  "utf8"
)

assert.match(
  routeSource,
  /name:\s*"dataType"[\s\S]*?type:\s*"select"/,
  "user reports dashboard should provide a dataType select filter"
)

assert.match(
  routeSource,
  /name:\s*"dataId"/,
  "user reports dashboard should provide a dataId filter"
)

assert.doesNotMatch(
  routeSource,
  /filters:\s*\[[\s\S]*?name:\s*"id"/,
  "user reports dashboard should not filter by report id"
)

assert.match(
  routeSource,
  /reportAuditStatusCell\(t,\s*record\.auditStatus\)/,
  "user reports dashboard should render auditStatus with a localized label"
)

for (const source of [zhMessages, enMessages]) {
  assert.match(
    source,
    /reportDataTypes:\s*{/,
    "dashboard translations should include report data type labels"
  )
  assert.match(
    source,
    /reportAuditStatus:\s*{/,
    "dashboard translations should include report audit status labels"
  )
}

console.log("dashboard user reports tests passed")
