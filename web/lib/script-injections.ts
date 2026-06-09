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
      const scriptTags = parseScriptTagFragment(normalized.code)
      if (scriptTags.length > 0) {
        return scriptTags.flatMap(
          (scriptTag, scriptIndex): RenderableScriptInjection[] => {
            const scriptKey =
              scriptTags.length === 1 ? key : `${key}-${scriptIndex}`
            if (scriptTag.src) {
              return [
                {
                  key: scriptKey,
                  type: "external",
                  src: scriptTag.src,
                  async: scriptTag.async,
                  defer: scriptTag.defer,
                  crossOrigin: normalizeCrossOrigin(scriptTag.crossOrigin),
                },
              ]
            }
            if (scriptTag.code) {
              return [{ key: scriptKey, type: "inline", code: scriptTag.code }]
            }
            return []
          }
        )
      }
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

function parseScriptTagFragment(code: string) {
  const scriptPattern = /<script\b([^>]*)>([\s\S]*?)<\/script>/gi
  const scripts: ReturnType<typeof createParsedScriptTag>[] = []
  let remainder = code.replace(/<!--[\s\S]*?-->/g, "")
  let match: RegExpExecArray | null

  while ((match = scriptPattern.exec(code))) {
    scripts.push(createParsedScriptTag(match[1], match[2]))
    remainder = remainder.replace(match[0], "")
  }

  return scripts.length > 0 && !remainder.trim() ? scripts : []
}

function createParsedScriptTag(attributesSource: string, code: string) {
  const attributes = parseScriptAttributes(attributesSource)
  return {
    src: trim(attributes.src),
    code: trim(code),
    async: attributes.async !== undefined,
    defer: attributes.defer !== undefined,
    crossOrigin: trim(attributes.crossorigin),
  }
}

function parseScriptAttributes(source: string) {
  const attributes: Record<string, string | true> = {}
  const attributePattern =
    /([^\s"'<>/=]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'=<>`]+)))?/g
  let match: RegExpExecArray | null

  while ((match = attributePattern.exec(source))) {
    attributes[match[1].toLowerCase()] =
      match[2] ?? match[3] ?? match[4] ?? true
  }

  return attributes
}

function normalizeCrossOrigin(
  value: string
): React.ScriptHTMLAttributes<HTMLScriptElement>["crossOrigin"] | undefined {
  if (value === "anonymous" || value === "use-credentials") return value
  return undefined
}
