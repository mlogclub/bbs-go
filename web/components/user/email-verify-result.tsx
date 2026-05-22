"use client"

import Link from "@/components/common/link"
import * as React from "react"

import { verifyEmailAction } from "@/lib/actions/user"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { useI18n } from "@/lib/i18n/provider"

export function EmailVerifyResult({ token }: { token: string }) {
  const { t } = useI18n()
  const [state, setState] = React.useState<{
    loading: boolean
    success: boolean
    email?: string
    message?: string
  }>({
    loading: true,
    success: false,
  })

  React.useEffect(() => {
    let canceled = false
    verifyEmailAction(token).then((result) => {
      if (canceled) {
        return
      }
      setState({
        loading: false,
        success: result.ok,
        email: result.email,
        message: result.message,
      })
    })
    return () => {
      canceled = true
    }
  }, [token])

  if (state.loading) {
    return null
  }

  return (
    <Alert variant={state.success ? "default" : "destructive"}>
      <AlertTitle>{t("user.email.verify.title")}</AlertTitle>
      <AlertDescription>
        {state.success ? (
          <div>
            {t("user.email.verify.success", { email: state.email || "" })}
          </div>
        ) : (
          <div>
            {t("user.email.verify.failed")}
            {state.message ? (
              <span>
                &nbsp;{t("user.email.verify.reason", { reason: state.message })}
              </span>
            ) : null}
            {t("user.email.verify.retryInstructions")}&nbsp;
            <Link href="/user/profile/account" className="font-bold">
              {t("user.email.verify.accountSettings")}
            </Link>
            &nbsp;{t("user.email.verify.retryAction")}
          </div>
        )}
      </AlertDescription>
    </Alert>
  )
}
