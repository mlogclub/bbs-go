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
  /name:\s*"auditStatus"[\s\S]*?type:\s*"select"/,
  "user reports dashboard should provide an auditStatus select filter"
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

assert.match(
  routeSource,
  /DASHBOARD_USER_REPORT_AUDIT/,
  "user reports dashboard actions should require the user report audit permission"
)

assert.match(
  routeSource,
  /\/api\/admin\/user-report\/audit/,
  "user reports dashboard should post audit actions to the audit endpoint"
)

assert.match(
  routeSource,
  /auditStatus:\s*1/,
  "user reports dashboard should provide a processed action"
)

assert.match(
  routeSource,
  /auditStatus:\s*2/,
  "user reports dashboard should provide an ignore action"
)

assert.match(
  routeSource,
  /detailFields:\s*\[/,
  "user reports dashboard should define detail fields"
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
  assert.match(
    source,
    /ignored:/,
    "dashboard translations should include ignored report audit status labels"
  )
  assert.match(
    source,
    /reportActions:\s*{/,
    "dashboard translations should include report action labels"
  )
}

console.log("dashboard user reports tests passed")
