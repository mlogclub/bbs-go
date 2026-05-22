import type { SiteConfig, Topic } from "@/lib/api/types"

type Metadata = {
  title?: string
  description?: string
  keywords?: string[]
}

export function siteTitle(config: SiteConfig | null | undefined, ...parts: Array<string | undefined>) {
  return [...parts.filter(Boolean), config?.siteTitle].filter(Boolean).join(" - ")
}

export function topicSiteTitle(config: SiteConfig | null | undefined, topic: Topic | null | undefined) {
  if (!topic) {
    return siteTitle(config)
  }

  return topic.type === 1 ? siteTitle(config, topic.content) : siteTitle(config, topic.title)
}

export function siteDescription(config: SiteConfig | null | undefined) {
  return config?.siteDescription || undefined
}

export function siteKeywords(config: SiteConfig | null | undefined) {
  return config?.siteKeywords?.length ? config.siteKeywords : undefined
}

export function siteMetadata(
  config: SiteConfig | null | undefined,
  parts: Array<string | undefined> = [],
  options: { includeSiteMeta?: boolean } = {},
): Metadata {
  return {
    title: siteTitle(config, ...parts),
    description: options.includeSiteMeta ? siteDescription(config) : undefined,
    keywords: options.includeSiteMeta ? siteKeywords(config) : undefined,
  }
}

export function safeRedirect(path: string | undefined, fallback = "/") {
  if (!path || !path.startsWith("/") || path.startsWith("//") || path.includes("\\")) {
    return fallback
  }

  try {
    const url = new URL(path, "http://local")
    if (url.origin !== "http://local" || !url.pathname.startsWith("/")) {
      return fallback
    }

    return `${url.pathname}${url.search}${url.hash}`
  } catch {
    return fallback
  }
}
