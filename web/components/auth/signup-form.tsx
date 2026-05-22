"use client"

import Link from "@/components/common/link"
import { useRouter } from "@/lib/router/navigation"
import type { ChangeEvent, FormEvent } from "react"
import { useActionState, useEffect, useRef, useState } from "react"

import { signupAction, type AuthActionState } from "@/lib/actions/auth"
import {
  CaptchaChallenge,
  type CaptchaChallengeHandle,
} from "@/components/auth/captcha-field"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useI18n } from "@/lib/i18n/provider"
import { toast } from "@/lib/toast"

const initialState: AuthActionState = { ok: false }

function authLink(path: string, redirect?: string) {
  return redirect ? `${path}?redirect=${encodeURIComponent(redirect)}` : path
}

export function SignupForm({ redirect }: { redirect?: string }) {
  const router = useRouter()
  const { t } = useI18n()
  const [state, action, pending] = useActionState(signupAction, initialState)
  const formRef = useRef<HTMLFormElement>(null)
  const captchaRef = useRef<CaptchaChallengeHandle>(null)
  const captchaVerifiedRef = useRef(false)
  const [form, setForm] = useState({
    nickname: "",
    email: "",
    password: "",
    rePassword: "",
  })
  const wasPendingRef = useRef(false)

  useEffect(() => {
    if (state.ok && state.redirect) {
      router.replace(state.redirect)
      router.refresh()
    }
  }, [router, state])

  useEffect(() => {
    if (wasPendingRef.current && !pending && !state.ok && state.message) {
      toast.error(state.message)
      captchaVerifiedRef.current = false
      captchaRef.current?.reset()
    }

    wasPendingRef.current = pending
  }, [pending, state.message, state.ok])

  function updateField(name: keyof typeof form) {
    return (event: ChangeEvent<HTMLInputElement>) => {
      setForm((current) => ({ ...current, [name]: event.target.value }))
    }
  }

  function validateForm() {
    if (!form.nickname.trim()) {
      toast.error(t("user.signup.nicknameRequired"))
      return false
    }
    if (!form.email.trim()) {
      toast.error(t("user.signup.emailRequired"))
      return false
    }
    if (!form.password) {
      toast.error(t("user.signup.passwordRequired"))
      return false
    }
    if (form.password !== form.rePassword) {
      toast.error(t("user.signup.passwordMismatch"))
      return false
    }
    return true
  }

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    if (captchaVerifiedRef.current) {
      captchaVerifiedRef.current = false
      return
    }

    event.preventDefault()
    if (!validateForm()) {
      return
    }

    void captchaRef.current?.open()
  }

  function onCaptchaVerified() {
    captchaVerifiedRef.current = true
    formRef.current?.requestSubmit()
  }

  return (
    <div>
      <form
        ref={formRef}
        action={action}
        className="space-y-6 p-6"
        onSubmit={onSubmit}
      >
        <input type="hidden" name="redirect" value={redirect || ""} />
        <div className="space-y-2">
          <Label htmlFor="nickname">
            {t("user.signup.nickname")}
            <span className="text-red-500">*</span>
          </Label>
          <Input
            id="nickname"
            name="nickname"
            autoComplete="nickname"
            placeholder={t("user.signup.nicknamePlaceholder")}
            value={form.nickname}
            onChange={updateField("nickname")}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="email">
            {t("user.signup.email")}
            <span className="text-red-500">*</span>
          </Label>
          <Input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            placeholder={t("user.signup.emailPlaceholder")}
            value={form.email}
            onChange={updateField("email")}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="password">
            {t("user.signup.password")}
            <span className="text-red-500">*</span>
          </Label>
          <Input
            id="password"
            name="password"
            type="password"
            autoComplete="new-password"
            placeholder={t("user.signup.passwordPlaceholder")}
            value={form.password}
            onChange={updateField("password")}
          />
          <p className="text-sm text-muted-foreground">
            {t("user.signup.passwordHelp")}
          </p>
        </div>
        <div className="space-y-2">
          <Label htmlFor="rePassword">
            {t("user.signup.confirmPassword")}
            <span className="text-red-500">*</span>
          </Label>
          <Input
            id="rePassword"
            name="rePassword"
            type="password"
            autoComplete="new-password"
            placeholder={t("user.signup.confirmPasswordPlaceholder")}
            value={form.rePassword}
            onChange={updateField("rePassword")}
          />
        </div>
        <CaptchaChallenge ref={captchaRef} onVerified={onCaptchaVerified} />
        <Button type="submit" className="h-10 w-full" disabled={pending}>
          {t("user.signup.signupBtn")}
        </Button>
        <div className="text-center">
          <Link
            href={authLink("/user/signin", redirect)}
            className="text-sm text-muted-foreground hover:text-primary"
          >
            {t("user.signup.alreadyHaveAccount")}
          </Link>
        </div>
      </form>
    </div>
  )
}
