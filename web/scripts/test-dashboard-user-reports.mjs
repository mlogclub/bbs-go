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
  "user reports dashboard process action should require the user report audit permission"
)

assert.match(
  routeSource,
  /\/api\/admin\/user-report\/audit/,
  "user reports dashboard modal should post audit actions to the audit endpoint"
)

assert.match(
  routeSource,
  /onSubmitStatus\(1\)/,
  "user reports dashboard modal should provide a processed action"
)

assert.match(
  routeSource,
  /onSubmitStatus\(2\)/,
  "user reports dashboard modal should provide an ignore action"
)

assert.doesNotMatch(
  routeSource,
  /detailFields:\s*\[/,
  "user reports dashboard should not use the generic details button"
)

assert.doesNotMatch(
  routeSource,
  /rowActions:\s*\[/,
  "user reports dashboard should not show separate processed and ignore row buttons"
)

assert.match(
  routeSource,
  /renderRowActions:/,
  "user reports dashboard should use one custom process row button"
)

assert.match(
  routeSource,
  /ReportProcessDialog/,
  "user reports dashboard should render a dedicated report process dialog"
)

assert.match(
  routeSource,
  /record\.target/,
  "user reports dashboard dialog should display reported target content from detail data"
)

assert.match(
  routeSource,
  /HtmlImagePreview/,
  "user reports dashboard dialog should render html target content with HtmlImagePreview"
)

assert.match(
  routeSource,
  /contentType\s*===\s*"html"/,
  "user reports dashboard dialog should branch on html content type"
)

assert.match(
  routeSource,
  /dataType\s*!==\s*"comment"/,
  "user reports dashboard should hide target detail navigation for comments"
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
  assert.match(
    source,
    /openTarget:/,
    "dashboard translations should include an open target label"
  )
  assert.match(
    source,
    /targetUnavailable:/,
    "dashboard translations should include an unavailable target message"
  )
}

console.log("dashboard user reports tests passed")
