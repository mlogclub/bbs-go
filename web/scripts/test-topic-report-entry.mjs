import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const topicMenu = readFileSync(
  resolve(webRoot, "components/topic/topic-manage-menu.tsx"),
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
  topicMenu,
  /\/api\/user-report\/submit/,
  "topic menu should submit reports to the user report endpoint"
)

assert.match(
  topicMenu,
  /dataType:\s*"topic"/,
  "topic report submission should identify reports as topic data"
)

assert.match(
  topicMenu,
  /canReport/,
  "topic menu should expose report actions to regular signed-in users"
)

for (const source of [zhMessages, enMessages]) {
  assert.match(
    source,
    /report:\s*"/,
    "topic menu translations should include a report action label"
  )
  assert.match(
    source,
    /reportSuccess:\s*"/,
    "topic menu translations should include report success feedback"
  )
  assert.match(
    source,
    /reportReason/,
    "topic menu translations should include report reasons"
  )
}

console.log("topic report entry tests passed")
