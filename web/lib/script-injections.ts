import type { SiteConfig } from "@/lib/api/types"
import type * as React from "react"

type ScriptInjectionConfig = NonNullable<SiteConfig["scriptInjections"]>[number]

export type RenderableScriptInjection =
  | {
      key: string
      type: "external"
      src: string
      async: boolean
      defer: boolean
      crossOrigin?: React.ScriptHTMLAttributes<HTMLScriptElement>["crossOrigin"]
    }
  | {
      key: string
      type: "inline"
      code: string
    }

export function getRenderableScriptInjections(
  injections: SiteConfig["scriptInjections"] | undefined
): RenderableScriptInjection[] {
  if (!Array.isArray(injections)) return []

  return injections.flatMap((injection, index): RenderableScriptInjection[] => {
    const normalized = normalizeScriptInjection(injection)
    if (!normalized.enabled) return []

    const key = `script-injection-${index}`
    if (normalized.type === "inline") {
      if (!normalized.code) return []
      return [{ key, type: "inline", code: normalized.code }]
    }

    if (!normalized.src) return []
    return [
      {
        key,
        type: "external",
        src: normalized.src,
        async: normalized.async,
        defer: normalized.defer,
        crossOrigin: normalizeCrossOrigin(normalized.crossOrigin),
      },
    ]
  })
}

export function getScriptInjectionElementId(key: string) {
  return `bbsgo-${key}`
}

function normalizeScriptInjection(injection: ScriptInjectionConfig) {
  return {
    enabled: Boolean(injection.enabled),
    type: injection.type === "inline" ? "inline" : "external",
    src: trim(injection.src),
    code: trim(injection.code),
    async: Boolean(injection.async),
    defer: Boolean(injection.defer),
    crossOrigin: trim(injection.crossorigin),
  }
}

function trim(value: unknown) {
  return typeof value === "string" ? value.trim() : ""
}

function normalizeCrossOrigin(
  value: string
): React.ScriptHTMLAttributes<HTMLScriptElement>["crossOrigin"] | undefined {
  if (value === "anonymous" || value === "use-credentials") return value
  return undefined
}
