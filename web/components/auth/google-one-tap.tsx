"use client"

import * as React from "react"
import { usePathname, useRouter, useSearchParams } from "@/lib/router/navigation"

import { googleOneTapSignin } from "@/lib/actions/auth"
import { useAppConfig, useIsLogin } from "@/components/app/app-provider"
import { useToastActions } from "@/lib/toast"

declare global {
  interface Window {
    google?: {
      accounts?: {
        id?: {
          initialize: (config: {
            client_id: string
            callback: (response: { credential?: string }) => void
            auto_select?: boolean
            cancel_on_tap_outside?: boolean
            itp_support?: boolean
          }) => void
          prompt: (callback?: (notification: {
            isNotDisplayed: () => boolean
            isSkippedMoment: () => boolean
            isDismissedMoment: () => boolean
            getNotDisplayedReason: () => string
          }) => void) => void
          cancel: () => void
        }
      }
    }
  }
}

function loadGoogleIdentityScript() {
  const src = "https://accounts.google.com/gsi/client"
  const existing = document.querySelector<HTMLScriptElement>(`script[src="${src}"]`)
  if (existing) {
    return Promise.resolve()
  }

  return new Promise<void>((resolve, reject) => {
    const script = document.createElement("script")
    script.src = src
    script.async = true
    script.defer = true
    script.onload = () => resolve()
    script.onerror = () => reject(new Error("Failed to load Google Identity Services"))
    document.head.appendChild(script)
  })
}

export function GoogleOneTap() {
  const config = useAppConfig()
  const isLogin = useIsLogin()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const router = useRouter()
  const { catchError } = useToastActions()
  const initializedRef = React.useRef(false)
  const clientId = config?.loginConfig?.googleLogin?.clientId
  const enabled = Boolean(config?.loginConfig?.googleLogin?.enabled && clientId)

  React.useEffect(() => {
    if (!enabled || isLogin || initializedRef.current) {
      return
    }
    if (pathname === "/user/signin" || pathname?.startsWith("/user/signin/")) {
      return
    }

    initializedRef.current = true
    let cancelled = false

    async function initialize() {
      try {
        await loadGoogleIdentityScript()
        if (cancelled || !window.google?.accounts?.id || !clientId) {
          return
        }

        window.google.accounts.id.initialize({
          client_id: clientId,
          auto_select: false,
          cancel_on_tap_outside: true,
          itp_support: true,
          callback: async (response) => {
            if (!response.credential) {
              return
            }
            try {
              const redirect = searchParams.get("redirect") || pathname || "/"
              const target = await googleOneTapSignin(response.credential, redirect)
              router.push(target)
              router.refresh()
            } catch (error) {
              catchError(error)
            }
          },
        })
        window.google.accounts.id.prompt()
      } catch {
        initializedRef.current = false
      }
    }

    void initialize()

    return () => {
      cancelled = true
      window.google?.accounts?.id?.cancel()
    }
  }, [catchError, clientId, enabled, isLogin, pathname, router, searchParams])

  return null
}
