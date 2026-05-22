import { getSiteConfig } from "@/lib/api/site"

import { createT, normalizeLocale } from "."

export async function getServerLocale(fallback?: string) {
  if (fallback) {
    return normalizeLocale(fallback)
  }

  const config = await getSiteConfig().catch(() => null)
  return normalizeLocale(config?.language)
}

export async function getServerT(fallback?: string) {
  return createT(await getServerLocale(fallback))
}
