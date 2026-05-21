import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { fileURLToPath } from "node:url"

import {
  getOAuthCallback,
  parseOAuthCallback,
} from "../lib/auth/oauth-callback.js"

assert.deepEqual(getOAuthCallback("github"), {
  kind: "login",
  submitPath: "/api/login/github_login_submit",
  loadingKey: "user.signin.githubLoggingIn",
  successRedirect: null,
  failureRedirect: "/user/signin",
})

assert.deepEqual(getOAuthCallback("google"), {
  kind: "login",
  submitPath: "/api/login/google_login_submit",
  loadingKey: "user.signin.googleLoggingIn",
  successRedirect: null,
  failureRedirect: "/user/signin",
})

assert.deepEqual(getOAuthCallback("weixin"), {
  kind: "login",
  submitPath: "/api/login/wx_login_submit",
  loadingKey: "user.signin.weixinLoggingIn",
  successRedirect: null,
  failureRedirect: "/user/signin",
})

assert.deepEqual(getOAuthCallback("google_bind"), {
  kind: "bind",
  submitPath: "/api/login/google_bind",
  loadingKey: "component.googleBindDialog.bindingText",
  successRedirect: "/user/profile/account",
  failureRedirect: "/user/profile/account",
})

assert.deepEqual(getOAuthCallback("weixin_bind"), {
  kind: "bind",
  submitPath: "/api/login/wx_bind",
  loadingKey: "component.wxBindDialog.bindingText",
  successRedirect: "/user/profile/account",
  failureRedirect: "/user/profile/account",
})

assert.deepEqual(
  getOAuthCallback("/user/signin/callback/weixin"),
  getOAuthCallback("weixin")
)

assert.equal(getOAuthCallback("unknown"), null)
assert.equal(getOAuthCallback(undefined), null)

assert.deepEqual(
  parseOAuthCallback("github", new URLSearchParams("code=c&state=s")),
  {
    ok: true,
    callback: getOAuthCallback("github"),
    code: "c",
    state: "s",
  }
)

assert.deepEqual(parseOAuthCallback("github", new URLSearchParams("code=c")), {
  ok: false,
  callback: getOAuthCallback("github"),
})

assert.deepEqual(
  parseOAuthCallback("unknown", new URLSearchParams("code=c&state=s")),
  {
    ok: false,
    callback: null,
  }
)

assert.deepEqual(
  parseOAuthCallback(
    "/user/signin/callback/weixin",
    new URLSearchParams("code=c&state=s")
  ),
  {
    ok: true,
    callback: getOAuthCallback("weixin"),
    code: "c",
    state: "s",
  }
)

const callbackRoute = readFileSync(
  fileURLToPath(new URL("../app/routes/user.signin_.callback.$.tsx", import.meta.url)),
  "utf8"
)
const apiClient = readFileSync(
  fileURLToPath(new URL("../lib/api/client.ts", import.meta.url)),
  "utf8"
)

assert.match(
  callbackRoute,
  /setCurrentUser\(result\.user\)/,
  "OAuth callback success must update AppProvider auth state immediately"
)

assert.match(
  callbackRoute,
  /window\.location\.pathname/,
  "OAuth callback route must fall back to the browser pathname when route splat params are unavailable"
)

assert.match(
  callbackRoute,
  /callback\.loadingKey/,
  "OAuth callback route must render provider-specific loading text while exchanging code"
)

assert.match(
  callbackRoute,
  /useDocumentTitle\(t\("user\.signin\.callbackTitle"\)\)/,
  "OAuth callback route must set a localized page title with the global site title suffix"
)

assert.match(
  callbackRoute,
  /window\.location\.replace/,
  "OAuth callback success must leave the callback URL with a hard replace after login"
)

assert.doesNotMatch(
  callbackRoute,
  /let active|active = false|if \(active\)/,
  "OAuth callback must not suppress the successful redirect when React re-runs effect cleanup"
)

assert.doesNotMatch(
  callbackRoute,
  /localStorage|persistClientAuthToken/,
  "OAuth callback must rely on the HttpOnly login cookie instead of persisting tokens in localStorage"
)

assert.match(
  apiClient,
  /credentials:\s*fetchOptions\.credentials\s*\?\?\s*"same-origin"/,
  "apiFetch must keep same-origin cookies enabled for login and current-user requests"
)

console.log("oauth callback tests passed")
