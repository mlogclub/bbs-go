import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const commentSource = readFileSync(
  resolve(webRoot, "components/comment/index.tsx"),
  "utf8"
)
const reportDialogSource = readFileSync(
  resolve(webRoot, "components/common/user-report-dialog.tsx"),
  "utf8"
)

assert.match(
  commentSource,
  /UserReportDialog/,
  "comment component should render the shared report dialog"
)

assert.match(
  commentSource,
  /dataType="comment"/,
  "comment report dialog should submit reports as comment data"
)

assert.match(
  commentSource,
  /component\.comment\.(list|subList)\.report/,
  "comment and reply action rows should include a report label"
)

assert.match(
  reportDialogSource,
  /\/api\/user-report\/submit/,
  "shared report dialog should submit to the user report endpoint"
)

console.log("comment report entry tests passed")
