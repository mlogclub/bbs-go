"use client"

import { usePathname, useRouter } from "@/lib/router/navigation"
import { toast, type ExternalToast } from "sonner"

export { toast }

import { ApiError } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"

const DEFAULT_DURATION = 800

type MsgType = "success" | "error" | "warning" | "info" | "default"

export interface MsgOptions {
  type?: MsgType
  message: string
  duration?: number
  onClose?: () => void
}

function toastOptions(options: MsgOptions): ExternalToast {
  return {
    duration: options.duration ?? DEFAULT_DURATION,
    onAutoClose: options.onClose,
  }
}

export function msg(options: MsgOptions) {
  const type = options.type ?? "success"
  const opts = toastOptions(options)

  if (type === "success") return toast.success(options.message, opts)
  if (type === "error") return toast.error(options.message, opts)
  if (type === "warning") return toast.warning(options.message, opts)
  if (type === "info") return toast.info(options.message, opts)
  return toast(options.message, opts)
}

export function msgSuccess(content: string) {
  return toast.success(content)
}

export function msgError(content: string) {
  return toast.error(content)
}

export function msgWarning(content: string) {
  return toast.warning(content)
}

function errorMessage(error: unknown, fallback: string) {
  if (error instanceof Error && error.message) {
    return error.message
  }
  if (typeof error === "string" && error) {
    return error
  }
  return fallback
}

export function isSignInError(error: unknown) {
  return (
    error instanceof ApiError &&
    (error.errorCode === 1 || error.status === 401 || error.status === 403)
  )
}

export function buildSigninHref(redirect?: string) {
  const path = redirect || "/"
  return `/user/signin?redirect=${encodeURIComponent(path)}`
}

export function useToastActions() {
  const router = useRouter()
  const pathname = usePathname()
  const { t } = useI18n()

  function msgSignIn(redirect?: string) {
    msg({
      type: "error",
      message: t("composables.pleaseSignIn"),
      onClose() {
        router.push(buildSigninHref(redirect || pathname || "/"))
      },
    })
  }

  function catchError(error: unknown) {
    if (isSignInError(error)) {
      msgSignIn()
      return
    }

    msgError(errorMessage(error, t("composables.unknownError")))
  }

  return {
    msg,
    msgSuccess,
    msgError,
    msgWarning,
    msgSignIn,
    catchError,
  }
}
