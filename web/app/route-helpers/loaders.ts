import { redirect } from "react-router"

import { apiFetch } from "@/lib/api/client"
import type { Article, PageData, Tag, Topic, TopicNode } from "@/lib/api/types"

import { getCurrentUser } from "./auth"

export type TopicListRouteData = {
  topics: PageData<Topic>
  nodes: TopicNode[]
  node?: TopicNode | null
  tag?: Tag | null
}

const qaStatusOptions = ["", "unsolved", "solved"]
const sortOptions = ["latestPublish", "latestReply"]

export function resolveTopicNodeId(id?: string) {
  if (id === "newest") return 0
  if (id === "recommend") return -1
  if (id === "feed") return -2
  const parsed = Number.parseInt(id || "", 10)
  return Number.isNaN(parsed) ? 0 : parsed
}

function getTopicFilters(request: Request) {
  const searchParams = new URL(request.url).searchParams
  const qaStatusValue = searchParams.get("qaStatus") || ""
  const sortValue = searchParams.get("sort") || ""

  return {
    qaStatus: qaStatusOptions.includes(qaStatusValue) ? qaStatusValue : "",
    sort: sortOptions.includes(sortValue) ? sortValue : "",
  }
}

async function loadTopicNode(
  request: Request | undefined,
  nodeId: number
): Promise<TopicNode | null> {
  if (nodeId <= 0) return null
  return apiFetch<TopicNode>("/api/topic/node", {
    request,
    params: { nodeId },
  }).catch(() => null)
}

async function loadTag(
  request: Request | undefined,
  tagId: string | number | undefined
): Promise<Tag | null> {
  if (!tagId) return null
  return apiFetch<Tag>(`/api/tag/${tagId}`, { request }).catch(() => null)
}

async function getTopicNodeFilters({
  request,
  nodeId,
}: {
  request?: Request
  nodeId: number
}) {
  if (!request || nodeId <= 0) return {}

  const filters = getTopicFilters(request)
  const node = await loadTopicNode(request, nodeId)
  if (node?.type === "qa") {
    return { qaStatus: filters.qaStatus }
  }
  if (node) {
    return { sort: filters.sort || "latestPublish" }
  }
  return {}
}

export async function loadTopicNodes(request?: Request) {
  return apiFetch<TopicNode[]>("/api/topic/node_navs", { request }).catch(
    () => []
  )
}

export async function loadTopics(params: {
  request?: Request
  cursor?: string
  nodeId?: string | number
  tagId?: string | number
  qaStatus?: string
  sort?: string
}) {
  const path = params.tagId ? "/api/topic/tag/topics" : "/api/topic/topics"

  return apiFetch<PageData<Topic>>(path, {
    request: params.request,
    params: {
      cursor: params.cursor || "",
      nodeId: params.nodeId,
      tagId: params.tagId,
      qaStatus: params.qaStatus,
      sort: params.sort,
    },
  })
}

export async function loadTopicListRouteData(
  request?: Request
): Promise<TopicListRouteData> {
  const [topics, nodes] = await Promise.all([
    loadTopics({ request }),
    loadTopicNodes(request),
  ])
  return { topics, nodes }
}

export async function loadTopicNodeRouteData({
  request,
  id,
}: {
  request?: Request
  id?: string
}): Promise<TopicListRouteData> {
  if (id === "feed") {
    const user = await getCurrentUser(request)
    if (!user) {
      throw redirect("/user/signin?redirect=/topics/node/feed")
    }
  }

  const nodeId = resolveTopicNodeId(id)
  const filters = await getTopicNodeFilters({ request, nodeId })
  const [topics, nodes, node] = await Promise.all([
    loadTopics({ request, nodeId, ...filters }),
    loadTopicNodes(request),
    loadTopicNode(request, nodeId),
  ])
  return { topics, nodes, node }
}

export async function loadTopicTagRouteData({
  request,
  id,
}: {
  request?: Request
  id?: string
}): Promise<TopicListRouteData> {
  const [topics, nodes, tag] = await Promise.all([
    loadTopics({ request, tagId: id }),
    loadTopicNodes(request),
    loadTag(request, id),
  ])
  return { topics, nodes, tag }
}

export async function loader({
  request,
  params,
}: {
  request: Request
  params?: { id?: string }
}) {
  const pathname = new URL(request.url).pathname

  if (pathname === "/" || pathname === "/topics") {
    return loadTopicListRouteData(request)
  }
  if (pathname.startsWith("/topics/node/")) {
    return loadTopicNodeRouteData({ request, id: params?.id })
  }
  if (pathname.startsWith("/topics/tag/")) {
    return loadTopicTagRouteData({ request, id: params?.id })
  }
  if (pathname === "/articles") {
    return loadArticles({ request })
  }
  if (pathname.startsWith("/articles/tag/")) {
    return loadArticleListRouteData({ request, tagId: params?.id })
  }

  return null
}

export async function loadArticles(params: {
  request?: Request
  cursor?: string
  tagId?: string | number
}) {
  const path = params.tagId
    ? "/api/article/tag/articles"
    : "/api/article/articles"

  return apiFetch<PageData<Article>>(path, {
    request: params.request,
    params: {
      cursor: params.cursor || "",
      tagId: params.tagId,
    },
  })
}

export type ArticleListRouteData = PageData<Article> & {
  tag?: Tag | null
}

export async function loadArticleListRouteData(params: {
  request?: Request
  cursor?: string
  tagId?: string | number
}): Promise<ArticleListRouteData> {
  const [articles, tag] = await Promise.all([
    loadArticles(params),
    params.tagId ? loadTag(params.request, params.tagId) : Promise.resolve(null),
  ])
  return Object.assign(articles, { tag })
}

export async function loadTopicDetail(params: {
  request?: Request
  id: string | number
}) {
  return apiFetch<Topic>(`/api/topic/${params.id}`, {
    request: params.request,
  })
}

export async function loadArticleDetail(params: {
  request?: Request
  id: string | number
}) {
  return apiFetch<Article>(`/api/article/${params.id}`, {
    request: params.request,
  })
}
