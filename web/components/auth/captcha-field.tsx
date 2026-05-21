"use client"

import GoCaptcha from "go-captcha-react"
import { RefreshCw } from "lucide-react"
import { forwardRef, useCallback, useEffect, useImperativeHandle, useRef, useState } from "react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { apiFetch } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"

interface CaptchaResponse {
  captchaId: string
  captchaBase64: string
}

interface RotateCaptchaResponse {
  id: string
  imageBase64: string
  thumbBase64: string
  thumbSize: number
}

export function CaptchaField() {
  const { t } = useI18n()
  const [captcha, setCaptcha] = useState<CaptchaResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refreshCaptcha = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      const data = await apiFetch<CaptchaResponse>("/api/captcha/request")
      setCaptcha(data)
    } catch {
      setCaptcha(null)
      setError(t("captcha.loadFailed"))
    } finally {
      setLoading(false)
    }
  }, [t])

  useEffect(() => {
    let active = true

    async function loadInitialCaptcha() {
      try {
        const data = await apiFetch<CaptchaResponse>("/api/captcha/request")
        if (active) {
          setCaptcha(data)
        }
      } catch {
        if (active) {
          setCaptcha(null)
          setError(t("captcha.loadFailed"))
        }
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    void loadInitialCaptcha()

    return () => {
      active = false
    }
  }, [t])

  return (
    <div className="space-y-2">
      <input type="hidden" name="captchaId" value={captcha?.captchaId || ""} />
      <input type="hidden" name="captchaProtocol" value="0" />
      <Label htmlFor="captchaCode">{t("captcha.code")}</Label>
      <div className="flex items-center gap-2">
        <div className="flex h-10 min-w-32 flex-1 items-center justify-center rounded-md border bg-muted px-2">
          {captcha?.captchaBase64 ? (
            <img
              src={`data:image/png;base64,${captcha.captchaBase64}`}
              alt={t("captcha.code")}
              className="h-8 max-w-40 object-contain"
            />
          ) : (
            <span className="text-sm text-muted-foreground">
              {loading ? t("captcha.loading") : t("captcha.loadFailed")}
            </span>
          )}
        </div>
        <Button type="button" variant="outline" size="icon" onClick={refreshCaptcha} disabled={loading}>
          <RefreshCw className={loading ? "animate-spin" : undefined} aria-hidden="true" />
          <span className="sr-only">{t("captcha.refresh")}</span>
        </Button>
      </div>
      <Input
        id="captchaCode"
        name="captchaCode"
        autoComplete="off"
        placeholder={t("captcha.codePlaceholder")}
        required
      />
      {error ? (
        <p className="text-sm text-destructive" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  )
}

export interface CaptchaChallengeHandle {
  open: () => Promise<void>
  reset: () => void
  hasCaptcha: () => boolean
  getCaptcha: () => {
    captchaId: string
    captchaCode: string
    captchaProtocol: number
  }
}

export const CaptchaChallenge = forwardRef<CaptchaChallengeHandle, { onVerified: () => void }>(function CaptchaChallenge(
  { onVerified },
  ref,
) {
  const { t } = useI18n()
  const [captcha, setCaptcha] = useState<RotateCaptchaResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [open, setOpen] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const captchaIdRef = useRef<HTMLInputElement>(null)
  const captchaCodeRef = useRef<HTMLInputElement>(null)
  const captchaProtocolRef = useRef<HTMLInputElement>(null)

  const refreshCaptcha = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      const data = await apiFetch<RotateCaptchaResponse>("/api/captcha/request_angle")
      setCaptcha(data)
    } catch {
      setCaptcha(null)
      setError(t("captcha.loadFailed"))
    } finally {
      setLoading(false)
    }
  }, [t])

  useImperativeHandle(
    ref,
    () => ({
      open: async () => {
        setOpen(true)
        await refreshCaptcha()
      },
      reset: () => {
        if (captchaIdRef.current) captchaIdRef.current.value = ""
        if (captchaCodeRef.current) captchaCodeRef.current.value = ""
        if (captchaProtocolRef.current) captchaProtocolRef.current.value = "2"
        setCaptcha(null)
        setOpen(false)
      },
      hasCaptcha: () => Boolean(captchaIdRef.current?.value && captchaCodeRef.current?.value),
      getCaptcha: () => ({
        captchaId: captchaIdRef.current?.value || "",
        captchaCode: captchaCodeRef.current?.value || "",
        captchaProtocol: Number(captchaProtocolRef.current?.value) || 2,
      }),
    }),
    [refreshCaptcha],
  )

  return (
    <>
      <input ref={captchaIdRef} type="hidden" name="captchaId" />
      <input ref={captchaCodeRef} type="hidden" name="captchaCode" />
      <input ref={captchaProtocolRef} type="hidden" name="captchaProtocol" value="2" readOnly />
      {open ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="w-auto max-w-none overflow-hidden rounded-lg bg-white p-0 shadow-lg">
            {captcha ? (
              <GoCaptcha.Rotate
                config={{ title: t("captcha.title") }}
                data={{
                  image: captcha.imageBase64,
                  thumb: captcha.thumbBase64,
                  thumbSize: captcha.thumbSize,
                  angle: 0,
                }}
                events={{
                  refresh: () => {
                    void refreshCaptcha()
                  },
                  close: () => setOpen(false),
                  confirm: (angle) => {
                    if (captchaIdRef.current) captchaIdRef.current.value = captcha.id
                    if (captchaCodeRef.current) captchaCodeRef.current.value = String(angle)
                    setOpen(false)
                    onVerified()
                  },
                }}
              />
            ) : (
              <div className="flex min-h-60 min-w-80 items-center justify-center p-6 text-sm text-muted-foreground">
                {loading ? t("captcha.loading") : error || t("captcha.loadFailed")}
              </div>
            )}
          </div>
        </div>
      ) : null}
    </>
  )
})
