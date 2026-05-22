import { toFormData } from "./client"
import { serverApiFetch as apiFetch } from "./server"

import type { LoginResult } from "./types"

interface CaptchaFields {
  captchaId?: string
  captchaCode?: string
  captchaProtocol?: string | number
}

export function signin(
  values: {
    username: string
    password: string
    redirect?: string
  } & CaptchaFields
) {
  return apiFetch<LoginResult>("/api/login/signin", {
    method: "POST",
    body: toFormData({ ...values }),
  })
}

export function signup(
  values: {
    nickname: string
    email: string
    password: string
    rePassword: string
    redirect?: string
  } & CaptchaFields
) {
  return apiFetch<LoginResult>("/api/login/signup", {
    method: "POST",
    body: toFormData({
      ...values,
      captchaProtocol: values.captchaProtocol ?? 2,
    }),
  })
}

export function signoutRequest() {
  return apiFetch<void>("/api/login/signout")
}

export function loginSmsCode(values: { phone: string } & CaptchaFields) {
  return apiFetch<{ smsId: string }>("/api/login/login_sms_code", {
    method: "POST",
    body: toFormData({
      ...values,
      captchaProtocol: values.captchaProtocol ?? 2,
    }),
  })
}

export function loginSms(values: {
  smsId: string
  smsCode: string
  redirect?: string
  state?: string
}) {
  return apiFetch<LoginResult>("/api/login/login_sms", {
    method: "POST",
    body: toFormData(values),
  })
}

export interface OAuthLoginConfig {
  authUrl?: string
  appid?: string
  scope?: string
  redirect_uri?: string
  state?: string
}

export function githubLoginSubmit(values: { code: string; state: string }) {
  return apiFetch<LoginResult>("/api/login/github_login_submit", {
    method: "POST",
    body: toFormData(values),
  })
}

export function googleLoginSubmit(values: { code: string; state: string }) {
  return apiFetch<LoginResult>("/api/login/google_login_submit", {
    method: "POST",
    body: toFormData(values),
  })
}

export function googleOneTapLogin(values: { credential: string }) {
  return apiFetch<LoginResult>("/api/login/google_one_tap", {
    method: "POST",
    body: toFormData(values),
  })
}

export function weixinLoginSubmit(values: { code: string; state: string }) {
  return apiFetch<LoginResult>("/api/login/wx_login_submit", {
    method: "POST",
    body: toFormData(values),
  })
}

export function sendResetPasswordEmail(values: {
  email: string
  captchaId?: string
  captchaCode?: string
  captchaProtocol?: string | number
}) {
  return apiFetch<void>("/api/login/send_reset_password_email", {
    method: "POST",
    body: toFormData({
      ...values,
      captchaProtocol: values.captchaProtocol ?? 2,
    }),
  })
}

export function resetPassword(values: {
  token: string
  password: string
  rePassword: string
}) {
  return apiFetch<void>("/api/login/reset_password", {
    method: "POST",
    body: toFormData(values),
  })
}
