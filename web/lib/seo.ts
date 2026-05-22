import type { MetaDescriptor } from "react-router"

import type { RootLoaderData } from "@/app/route-helpers/types"
import type {
  Article,
  SiteConfig,
  Tag,
  Topic,
  TopicNode,
  UserSummary,
} from "@/lib/api/types"

type RouteMatch = {
  data?: unknown
  loaderData?: unknown
}

type ContentMetaInput = {
  title?: string
  description?: string
  keywords?: string[]
  image?: string
  canonicalPath?: string
  structuredData?: Record<string, unknown>
  ogType?: "article" | "profile" | "website"
  noindex?: boolean
}

function trimText(value: string | null | undefined) {
  return value?.trim() || undefined
}

function siteName(config: SiteConfig | null | undefined) {
  return trimText(config?.siteTitle) || "BBS-GO"
}

function siteDescription(config: SiteConfig | null | undefined) {
  return trimText(config?.siteDescription)
}

function siteKeywords(config: SiteConfig | null | undefined) {
  const keywords = config?.siteKeywords
    ?.map((keyword) => trimText(keyword))
    .filter((keyword): keyword is string => Boolean(keyword))

  return keywords?.length ? keywords : undefined
}

function normalizeBaseURL(config: SiteConfig | null | undefined) {
  const baseURL = trimText(config?.baseURL)
  if (!baseURL || baseURL === "/") return undefined

  try {
    const url = new URL(baseURL)
    if (url.protocol !== "http:" && url.protocol !== "https:") {
      return undefined
    }
    url.pathname = url.pathname.replace(/\/+$/, "")
    url.search = ""
    url.hash = ""
    return url.toString().replace(/\/$/, "")
  } catch {
    return undefined
  }
}

function normalizePath(path: string | null | undefined) {
  if (!path) return "/"
  const parsed = new URL(path, "https://example.com")
  let pathname = parsed.pathname || "/"
  if (pathname.length > 1) pathname = pathname.replace(/\/+$/, "")
  return pathname
}

function absoluteURL(
  config: SiteConfig | null | undefined,
  pathOrURL: string | null | undefined
) {
  const value = trimText(pathOrURL)
  if (!value) return undefined
  if (/^https?:\/\//i.test(value)) return value

  const baseURL = normalizeBaseURL(config)
  if (!baseURL) return undefined

  return `${baseURL}${value.startsWith("/") ? value : `/${value}`}`
}

function canonicalURL(
  config: SiteConfig | null | undefined,
  path: string | null | undefined
) {
  return absoluteURL(config, normalizePath(path))
}

function imageURL(
  config: SiteConfig | null | undefined,
  image: string | null | undefined
) {
  return absoluteURL(config, image)
}

function compactObject<T extends Record<string, unknown>>(value: T) {
  return Object.fromEntries(
    Object.entries(value).filter(([, item]) => {
      if (Array.isArray(item)) return item.length > 0
      return item !== undefined && item !== null && item !== ""
    })
  )
}

function isRootLoaderData(value: unknown): value is RootLoaderData {
  return Boolean(
    value &&
      typeof value === "object" &&
      "config" in value &&
      "currentUser" in value &&
      "locale" in value
  )
}

export function rootDataFromMatches(matches: RouteMatch[]) {
  for (const match of matches) {
    if (isRootLoaderData(match.data)) return match.data
    if (isRootLoaderData(match.loaderData)) return match.loaderData
  }
  return undefined
}

export function localizedTitle(
  locale: RootLoaderData["locale"] | string | undefined,
  enUS: string,
  zhCN: string
) {
  return locale === "en-US" ? enUS : zhCN
}

export function tagKeywords(tags: Tag[] | null | undefined) {
  const keywords = tags
    ?.map((tag) => trimText(tag.name))
    .filter((keyword): keyword is string => Boolean(keyword))

  return keywords?.length ? keywords : undefined
}

export function siteMeta(config: SiteConfig | null | undefined) {
  const title = siteName(config)
  const description = siteDescription(config)
  const keywords = siteKeywords(config)

  return compactMeta([
    { title },
    description ? { name: "description", content: description } : undefined,
    keywords?.length
      ? { name: "keywords", content: keywords.join(",") }
      : undefined,
  ])
}

export function sitePageMeta(
  config: SiteConfig | null | undefined,
  title: string | null | undefined,
  options: {
    canonicalPath?: string
  } = {}
) {
  const description = siteDescription(config)
  const keywords = siteKeywords(config)

  return pageMeta(config, title, {
    description,
    keywords,
    canonicalPath: options.canonicalPath,
  })
}

export function contentMeta(
  config: SiteConfig | null | undefined,
  content: ContentMetaInput
) {
  const siteTitle = trimText(config?.siteTitle)
  const title = trimText(content.title)
  const description = trimText(content.description)
  const keywords = content.keywords
    ?.map((keyword) => trimText(keyword))
    .filter((keyword): keyword is string => Boolean(keyword))
  const resolvedTitle = title
    ? siteTitle
      ? `${title} - ${siteTitle}`
      : title
    : undefined
  const canonical = canonicalURL(config, content.canonicalPath)
  const image = imageURL(config, content.image)

  return compactMeta([
    resolvedTitle ? { title: resolvedTitle } : undefined,
    description ? { name: "description", content: description } : undefined,
    keywords?.length
      ? { name: "keywords", content: keywords.join(",") }
      : undefined,
    canonical
      ? { tagName: "link", rel: "canonical", href: canonical }
      : undefined,
    ...socialMeta({
      config,
      title: resolvedTitle,
      description,
      url: canonical,
      image,
      type: content.ogType || "website",
    }),
    content.structuredData
      ? { "script:ld+json": content.structuredData }
      : undefined,
    content.noindex ? noindexDescriptor() : undefined,
  ])
}

export function pageMeta(
  config: SiteConfig | null | undefined,
  title: string | null | undefined,
  options: {
    description?: string
    keywords?: string[]
    image?: string
    canonicalPath?: string
    structuredData?: Record<string, unknown>
    ogType?: "article" | "profile" | "website"
    noindex?: boolean
  } = {}
) {
  return contentMeta(config, {
    title: title ?? undefined,
    description: options.description,
    keywords: options.keywords,
    image: options.image,
    canonicalPath: options.canonicalPath,
    structuredData: options.structuredData,
    ogType: options.ogType,
    noindex: options.noindex,
  })
}

export function noindexMeta(
  config: SiteConfig | null | undefined,
  title: string | null | undefined
) {
  return pageMeta(config, title, { noindex: true })
}

export function noindexRouteMeta(
  matches: RouteMatch[],
  enUS: string,
  zhCN: string
) {
  const rootData = rootDataFromMatches(matches)
  return noindexMeta(
    rootData?.config,
    localizedTitle(rootData?.locale, enUS, zhCN)
  )
}

export function topicMeta(
  config: SiteConfig | null | undefined,
  topic: Topic | null | undefined,
  canonicalPath?: string
) {
  const title = topic?.type === 1 ? topic?.content : topic?.title
  const image = topic?.imageList?.[0]?.url || topic?.imageList?.[0]?.preview
  return contentMeta(config, {
    title,
    description: topic?.summary,
    keywords: tagKeywords(topic?.tags),
    image,
    canonicalPath,
    structuredData: topic
      ? compactObject({
          "@context": "https://schema.org",
          "@type": "DiscussionForumPosting",
          headline: title,
          text: topic.summary,
          image: imageURL(config, image),
          url: canonicalURL(config, canonicalPath),
          datePublished: timestampToIso(topic.createTime),
          dateModified: timestampToIso(topic.updateTime || topic.createTime),
          author: personStructuredData(config, topic.user),
        })
      : undefined,
    ogType: "article",
  })
}

export function topicNodeMeta(
  config: SiteConfig | null | undefined,
  node: TopicNode | null | undefined,
  canonicalPath?: string
) {
  return pageMeta(config, node?.name, {
    description: node?.description,
    image: node?.logo,
    canonicalPath,
  })
}

export function tagPageMeta(
  config: SiteConfig | null | undefined,
  tag: Tag | null | undefined,
  suffix: string,
  canonicalPath?: string
) {
  const title = trimText(tag?.name)
  return pageMeta(config, title ? `${title} - ${suffix}` : suffix, {
    description: tag?.description,
    canonicalPath,
  })
}

export function articleMeta(
  config: SiteConfig | null | undefined,
  article: Article | null | undefined,
  canonicalPath?: string
) {
  const image = article?.cover?.url || article?.cover?.preview
  return contentMeta(config, {
    title: article?.title,
    description: article?.summary,
    keywords: tagKeywords(article?.tags),
    image,
    canonicalPath,
    structuredData: article
      ? compactObject({
          "@context": "https://schema.org",
          "@type": "Article",
          headline: article.title,
          description: article.summary,
          image: imageURL(config, image),
          url: canonicalURL(config, canonicalPath),
          datePublished: timestampToIso(article.createTime),
          dateModified: timestampToIso(article.createTime),
          author: personStructuredData(config, article.user),
        })
      : undefined,
    ogType: "article",
  })
}

export function userDisplayName(user: UserSummary | null | undefined) {
  return (
    trimText(user?.nickname) ||
    trimText(user?.username) ||
    (user?.id ? `#${user.id}` : undefined)
  )
}

export function userMeta(
  config: SiteConfig | null | undefined,
  user: UserSummary | null | undefined,
  suffix?: string,
  canonicalPath?: string
) {
  const name = userDisplayName(user)
  const title = [name, suffix].filter(Boolean).join(" - ")
  return pageMeta(config, title, {
    description: user?.description,
    image: user?.avatar || user?.smallAvatar,
    canonicalPath,
    structuredData:
      user && !suffix
        ? personStructuredData(config, user, canonicalPath)
        : undefined,
    ogType: "profile",
  })
}

export function siteHomeMeta(config: SiteConfig | null | undefined) {
  const title = siteName(config)
  const description = siteDescription(config)
  const keywords = siteKeywords(config)
  const canonical = canonicalURL(config, "/")
  return compactMeta([
    { title },
    description ? { name: "description", content: description } : undefined,
    keywords?.length
      ? { name: "keywords", content: keywords.join(",") }
      : undefined,
    canonical
      ? { tagName: "link", rel: "canonical", href: canonical }
      : undefined,
    ...socialMeta({
      config,
      title,
      description,
      url: canonical,
      image: imageURL(config, config?.siteLogo),
      type: "website",
    }),
    {
      "script:ld+json": compactObject({
        "@context": "https://schema.org",
        "@type": "WebSite",
        name: title,
        description,
        url: canonical,
      }),
    },
  ])
}

function socialMeta({
  config,
  title,
  description,
  url,
  image,
  type,
}: {
  config: SiteConfig | null | undefined
  title?: string
  description?: string
  url?: string
  image?: string
  type: "article" | "profile" | "website"
}): MetaDescriptor[] {
  return compactMeta([
    title ? { property: "og:title", content: title } : undefined,
    description
      ? { property: "og:description", content: description }
      : undefined,
    { property: "og:type", content: type },
    { property: "og:site_name", content: siteName(config) },
    url ? { property: "og:url", content: url } : undefined,
    image ? { property: "og:image", content: image } : undefined,
    { name: "twitter:card", content: image ? "summary_large_image" : "summary" },
    title ? { name: "twitter:title", content: title } : undefined,
    description
      ? { name: "twitter:description", content: description }
      : undefined,
    image ? { name: "twitter:image", content: image } : undefined,
  ])
}

function personStructuredData(
  config: SiteConfig | null | undefined,
  user: UserSummary | null | undefined,
  canonicalPath?: string
) {
  if (!user) return undefined
  return compactObject({
    "@context": canonicalPath ? "https://schema.org" : undefined,
    "@type": "Person",
    name: userDisplayName(user),
    description: user.description,
    image: imageURL(config, user.avatar || user.smallAvatar),
    url: canonicalURL(config, canonicalPath),
  })
}

function timestampToIso(timestamp: number | null | undefined) {
  if (!timestamp) return undefined
  return new Date(timestamp).toISOString()
}

function noindexDescriptor(): MetaDescriptor {
  return { name: "robots", content: "noindex,nofollow" }
}

function compactMeta(items: Array<MetaDescriptor | undefined>) {
  return items.filter((item): item is MetaDescriptor => Boolean(item))
}
