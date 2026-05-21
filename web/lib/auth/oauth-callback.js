const callbacks = {
  github: {
    kind: "login",
    submitPath: "/api/login/github_login_submit",
    loadingKey: "user.signin.githubLoggingIn",
    successRedirect: null,
    failureRedirect: "/user/signin",
  },
  google: {
    kind: "login",
    submitPath: "/api/login/google_login_submit",
    loadingKey: "user.signin.googleLoggingIn",
    successRedirect: null,
    failureRedirect: "/user/signin",
  },
  weixin: {
    kind: "login",
    submitPath: "/api/login/wx_login_submit",
    loadingKey: "user.signin.weixinLoggingIn",
    successRedirect: null,
    failureRedirect: "/user/signin",
  },
  google_bind: {
    kind: "bind",
    submitPath: "/api/login/google_bind",
    loadingKey: "component.googleBindDialog.bindingText",
    successRedirect: "/user/profile/account",
    failureRedirect: "/user/profile/account",
  },
  weixin_bind: {
    kind: "bind",
    submitPath: "/api/login/wx_bind",
    loadingKey: "component.wxBindDialog.bindingText",
    successRedirect: "/user/profile/account",
    failureRedirect: "/user/profile/account",
  },
}

function callbackKey(pathname) {
  const key = pathname.replace(/^\/+|\/+$/g, "")
  const prefix = "user/signin/callback/"
  return key.startsWith(prefix) ? key.slice(prefix.length) : key
}

export function getOAuthCallback(pathname) {
  if (!pathname) return null
  const key = callbackKey(pathname)
  return callbacks[key] ?? null
}

export function parseOAuthCallback(pathname, searchParams) {
  const callback = getOAuthCallback(pathname)
  if (!callback) return { ok: false, callback: null }

  const code = searchParams.get("code") || ""
  const state = searchParams.get("state") || ""
  if (!code || !state) return { ok: false, callback }

  return { ok: true, callback, code, state }
}
