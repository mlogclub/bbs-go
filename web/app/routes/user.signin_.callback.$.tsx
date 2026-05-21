import * as React from "react"
import { useNavigate, useParams, useSearchParams } from "react-router"
import { LoaderCircle } from "lucide-react"

import { useAppState } from "@/components/app/app-provider"
import { apiFetch, toFormData } from "@/lib/api/client"
import type { LoginResult } from "@/lib/api/types"
import { parseOAuthCallback } from "@/lib/auth/oauth-callback.js"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { safeRedirect } from "@/lib/site"
import { useToastActions } from "@/lib/toast"
import { useDocumentTitle } from "@/lib/use-document-title"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Signing in", "登录中")
}

export default function SigninCallbackRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.signin.callbackTitle"))
  const navigate = useNavigate()
  const params = useParams()
  const [searchParams] = useSearchParams()
  const { setCurrentUser } = useAppState()
  const { catchError, msgError } = useToastActions()
  const submittedRef = React.useRef(false)
  const callbackPath =
    params["*"] ||
    (typeof window !== "undefined" ? window.location.pathname : "")
  const callback = parseOAuthCallback(callbackPath, searchParams).callback

  React.useEffect(() => {
    if (submittedRef.current) return
    submittedRef.current = true

    async function submitCallback() {
      const parsed = parseOAuthCallback(callbackPath, searchParams)
      if (!parsed.ok) {
        msgError(t("user.signin.missingAuthParams"))
        navigate(parsed.callback?.failureRedirect || "/user/signin", {
          replace: true,
        })
        return
      }

      try {
        if (parsed.callback.kind === "bind") {
          await apiFetch<void>(parsed.callback.submitPath, {
            method: "POST",
            body: toFormData({ code: parsed.code, state: parsed.state }),
          })
          navigate(parsed.callback.successRedirect, { replace: true })
          return
        }

        const result = await apiFetch<LoginResult>(parsed.callback.submitPath, {
          method: "POST",
          body: toFormData({ code: parsed.code, state: parsed.state }),
        })
        setCurrentUser(result.user)
        window.location.replace(
          safeRedirect(result.redirect, `/user/${result.user.id}`)
        )
      } catch (error) {
        catchError(error)
        const parsed = parseOAuthCallback(callbackPath, searchParams)
        navigate(parsed.callback?.failureRedirect || "/user/signin", {
          replace: true,
        })
      }
    }

    void submitCallback()
  }, [
    callbackPath,
    catchError,
    msgError,
    navigate,
    searchParams,
    setCurrentUser,
    t,
  ])

  return (
    <section className="main">
      <div className="container" style={{ height: 300 }}>
        <div className="fixed inset-0 z-[1000] flex items-center justify-center bg-background/70 backdrop-blur-sm">
          <div className="flex min-w-40 flex-col items-center gap-3 rounded-md border bg-background px-6 py-5 text-sm text-muted-foreground shadow-lg">
            <LoaderCircle className="h-6 w-6 animate-spin text-primary" />
            <span>
              {callback
                ? t(callback.loadingKey)
                : t("user.signin.missingAuthParams")}
            </span>
          </div>
        </div>
      </div>
    </section>
  )
}
