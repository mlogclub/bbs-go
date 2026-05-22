import type { AppLocale } from "./types"

export function normalizeLocale(value: unknown): AppLocale {
  return value === "zh-CN" ? "zh-CN" : "en-US"
}

export function getBrowserLocale(fallback: AppLocale): AppLocale {
  return fallback
}
