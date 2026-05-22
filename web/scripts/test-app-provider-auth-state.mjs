import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { fileURLToPath } from "node:url"

const appProvider = readFileSync(
  fileURLToPath(new URL("../components/app/app-provider.tsx", import.meta.url)),
  "utf8"
)
const topicCreate = readFileSync(
  fileURLToPath(new URL("../app/routes/topic.create.tsx", import.meta.url)),
  "utf8"
)

assert.match(
  appProvider,
  /userStateTouchedRef/,
  "AppProvider must track local auth state updates"
)

assert.match(
  appProvider,
  /if \(!userStateTouchedRef\.current\)[\s\S]*setCurrentUserState\(nextState\.currentUser\)/,
  "stale client hydration must not overwrite a user set by a login callback"
)

assert.match(
  topicCreate,
  /useAuthChecked/,
  "topic create must wait for client auth hydration before redirecting to signin"
)

assert.match(
  topicCreate,
  /if \(!authChecked\)[\s\S]*return null/,
  "topic create must not redirect while auth state is still hydrating"
)

console.log("app provider auth state tests passed")
