import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { fileURLToPath } from "node:url"

const signinForm = readFileSync(
  fileURLToPath(new URL("../components/auth/signin-form.tsx", import.meta.url)),
  "utf8"
)

const navigation = readFileSync(
  fileURLToPath(new URL("../lib/router/navigation.ts", import.meta.url)),
  "utf8"
)

assert.doesNotMatch(
  signinForm,
  /authChecked\s*&&\s*currentUser[\s\S]*router\.replace/,
  "signin page must not auto-redirect logged-in users; site keeps signin page passive"
)

assert.doesNotMatch(
  signinForm,
  /useCurrentUser|useAuthChecked/,
  "signin form must not depend on app auth state for passive page rendering"
)

assert.doesNotMatch(
  signinForm,
  /showWeixinQR/,
  "WeChat login should open a modal instead of toggling an inline QR code under the button"
)

assert.match(
  signinForm,
  /weixinDialogOpen/,
  "WeChat login should track modal open state"
)

assert.match(
  navigation,
  /React\.useMemo/,
  "useRouter must return a stable object so effects do not re-run just because the wrapper identity changed"
)

console.log("signin navigation tests passed")
