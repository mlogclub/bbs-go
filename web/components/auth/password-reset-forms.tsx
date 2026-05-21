"use client"

import Link from "@/components/common/link"
import { useRouter } from "@/lib/router/navigation"
import * as React from "react"

import {
  resetPasswordAction,
  sendResetPasswordEmailAction,
  type AuthActionState,
} from "@/lib/actions/auth"
import {
  CaptchaChallenge,
  type CaptchaChallengeHandle,
} from "@/components/auth/captcha-field"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { toast } from "@/lib/toast"
import { useI18n } from "@/lib/i18n/provider"

const initialState: AuthActionState = { ok: false }

export function ForgotPasswordForm() {
  const { t } = useI18n()
  const [email, setEmail] = React.useState("")
  const [sent, setSent] = React.useState(false)
  const [pending, startTransition] = React.useTransition()
  const captchaRef = React.useRef<CaptchaChallengeHandle>(null)
  const formRef = React.useRef<HTMLFormElement>(null)

  function submit(formData: FormData) {
    startTransition(async () => {
      const result = await sendResetPasswordEmailAction(initialState, formData)
      if (!result.ok) {
        toast.error(result.message || t("composables.unknownError"))
        captchaRef.current?.reset()
        return
      }
      setSent(true)
      toast.success(t("user.passwordReset.forgot.sentNotice"))
    })
  }

  return (
    <>
      <h2 className="mb-6 text-xl font-semibold">
        {t("user.passwordReset.forgot.title")}
      </h2>
      <form
        ref={formRef}
        action={submit}
        className="space-y-4"
        onSubmit={(event) => {
          if (!email) {
            event.preventDefault()
            toast.error(t("user.passwordReset.forgot.emailRequired"))
            return
          }
          if (!captchaRef.current?.hasCaptcha()) {
            event.preventDefault()
            void captchaRef.current?.open()
          }
        }}
      >
        <CaptchaChallenge
          ref={captchaRef}
          onVerified={() => formRef.current?.requestSubmit()}
        />
        <Input
          name="email"
          type="email"
          placeholder={t("user.passwordReset.forgot.emailPlaceholder")}
          value={email}
          onChange={(event) => setEmail(event.target.value)}
        />
        <Button type="submit" className="w-full" disabled={pending}>
          {pending
            ? t("user.passwordReset.forgot.sending")
            : t("user.passwordReset.forgot.sendButton")}
        </Button>
      </form>
      {sent ? (
        <Alert className="mt-4">
          <AlertDescription>
            {t("user.passwordReset.forgot.sentNotice")}
          </AlertDescription>
        </Alert>
      ) : null}
      <div className="mt-6 text-center">
        <Link href="/user/signin" className="text-sm text-primary">
          {t("user.passwordReset.backToSignin")}
        </Link>
      </div>
    </>
  )
}

export function ResetPasswordForm({ token }: { token: string }) {
  const { t } = useI18n()
  const router = useRouter()
  const [errorMsg, setErrorMsg] = React.useState("")
  const [success, setSuccess] = React.useState(false)
  const [pending, startTransition] = React.useTransition()

  function submit(formData: FormData) {
    setErrorMsg("")
    const password = String(formData.get("password") || "")
    const rePassword = String(formData.get("rePassword") || "")
    if (!token) {
      setErrorMsg(t("user.passwordReset.reset.tokenMissing"))
      return
    }
    if (!password) {
      setErrorMsg(t("user.passwordReset.reset.passwordRequired"))
      return
    }
    if (!rePassword) {
      setErrorMsg(t("user.passwordReset.reset.rePasswordRequired"))
      return
    }
    if (password !== rePassword) {
      setErrorMsg(t("user.passwordReset.reset.passwordMismatch"))
      return
    }

    startTransition(async () => {
      const result = await resetPasswordAction(initialState, formData)
      if (!result.ok) {
        setErrorMsg(result.message || t("composables.unknownError"))
        return
      }
      setSuccess(true)
      toast.success(t("user.passwordReset.reset.successNotice"))
      window.setTimeout(() => router.push("/user/signin"), 1200)
    })
  }

  return (
    <>
      <h2 className="mb-6 text-xl font-semibold">
        {t("user.passwordReset.reset.title")}
      </h2>
      {!token ? (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>
            {t("user.passwordReset.reset.tokenMissing")}
          </AlertDescription>
        </Alert>
      ) : null}
      {errorMsg ? (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{errorMsg}</AlertDescription>
        </Alert>
      ) : null}
      <form action={submit} className="space-y-4">
        <input type="hidden" name="token" value={token} />
        <Input
          name="password"
          type="password"
          placeholder={t("user.passwordReset.reset.passwordPlaceholder")}
          disabled={!token || pending || success}
        />
        <Input
          name="rePassword"
          type="password"
          placeholder={t("user.passwordReset.reset.rePasswordPlaceholder")}
          disabled={!token || pending || success}
        />
        <Button
          type="submit"
          className="w-full"
          disabled={!token || pending || success}
        >
          {pending
            ? t("user.passwordReset.reset.submitting")
            : t("user.passwordReset.reset.submitButton")}
        </Button>
      </form>
      {success ? (
        <Alert className="mt-4">
          <AlertDescription>
            {t("user.passwordReset.reset.successNotice")}
          </AlertDescription>
        </Alert>
      ) : null}
      <div className="mt-6 text-center">
        <Link href="/user/signin" className="text-sm text-primary">
          {t("user.passwordReset.backToSignin")}
        </Link>
      </div>
    </>
  )
}
