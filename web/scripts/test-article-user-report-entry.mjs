import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const articleMenu = readFileSync(
  resolve(webRoot, "components/article/article-manage-menu.tsx"),
  "utf8"
)
const userProfileCard = readFileSync(
  resolve(webRoot, "components/user/user-profile-card.tsx"),
  "utf8"
)
const userCenterOperations = readFileSync(
  resolve(webRoot, "components/user/user-center-operations.tsx"),
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
  articleMenu,
  /UserReportDialog/,
  "article menu should render the shared report dialog"
)

assert.match(
  articleMenu,
  /dataType="article"/,
  "article report dialog should submit reports as article data"
)

assert.match(
  articleMenu,
  /canReport/,
  "article menu should expose report actions to regular signed-in users"
)

assert.match(
  userCenterOperations,
  /UserReportDialog/,
  "user center operations should render the shared report dialog"
)

assert.match(
  userCenterOperations,
  /dataType="user"/,
  "user center operations should submit reports as user data"
)

assert.doesNotMatch(
  userProfileCard,
  /UserReportDialog|dataType="user"|component\.userProfile\.report/,
  "user profile card should not contain the report entry"
)

for (const source of [zhMessages, enMessages]) {
  assert.match(
    source,
    /articleManageMenu:[\s\S]*?report:\s*"/,
    "article menu translations should include a report label"
  )
  assert.match(
    source,
    /userCenterSidebar:[\s\S]*?report:\s*"/,
    "user center operations translations should include a report label"
  )
}

console.log("article and user report entry tests passed")
