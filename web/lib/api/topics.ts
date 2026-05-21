import { serverApiFetch as apiFetch } from "./server"

import type {
  PageData,
  Tag,
  Topic,
  TopicHideContent,
  TopicNode,
  UserSummary,
} from "./types"

type TopicParams = Record<string, string | number | boolean | undefined>
type SearchTopicParams = {
  keyword: string
  nodeId?: number
  timeRange?: number
  cursor?: string
}

export function getTopics(params: TopicParams = {}) {
  return apiFetch<PageData<Topic>>("/api/topic/topics", { params })
}

export function getNodeTopics(
  nodeId: string | number,
  params: TopicParams = {}
) {
  return getTopics({ ...params, nodeId })
}

export function getTagTopics(tagId: string | number, cursor?: string) {
  return apiFetch<PageData<Topic>>("/api/topic/tag/topics", {
    params: { tagId, cursor },
  })
}

export function getTopic(id: string) {
  return apiFetch<Topic>(`/api/topic/${id}`)
}

export type TopicEditData = Pick<
  Topic,
  "id" | "type" | "title" | "content" | "attachments"
> & {
  nodeId: number
  contentType?: "html" | "markdown" | string
  hideContent?: string
  tags?: string[] | null
}

export function getTopicEdit(id: string) {
  return apiFetch<TopicEditData>(`/api/topic/edit/${id}`)
}

export function getTopicRecentLikes(id: string) {
  return apiFetch<UserSummary[] | null>(`/api/topic/recentlikes/${id}`)
}

export function getTopicHideContent(id: string) {
  return apiFetch<TopicHideContent>("/api/topic/hide_content", {
    params: { topicId: id },
  })
}

export function getTopicNode(nodeId: string | number) {
  return apiFetch<TopicNode>("/api/topic/node", { params: { nodeId } })
}

export function getTopicNodes() {
  return apiFetch<TopicNode[]>("/api/topic/nodes")
}

export function getTopicNodeNavs() {
  return apiFetch<TopicNode[]>("/api/topic/node_navs")
}

export function getTag(id: string | number) {
  return apiFetch<Tag>(`/api/tag/${id}`)
}

export function searchTopics(params: SearchTopicParams) {
  return apiFetch<PageData<Topic>>("/api/search/topic", {
    params,
  })
}
