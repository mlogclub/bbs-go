import { serverApiFetch as apiFetch } from "./server"

import type {
  PageData,
  Tag,
  Topic,
  TopicHideContent,
  Category,
  UserSummary,
} from "./types"

type TopicParams = Record<string, string | number | boolean | undefined>
type SearchTopicParams = {
  keyword: string
  categoryId?: number
  timeRange?: number
  cursor?: string
}

export function getTopics(params: TopicParams = {}) {
  return apiFetch<PageData<Topic>>("/api/topic/topics", { params })
}

export function getCategoryTopics(
  categoryId: string | number,
  params: TopicParams = {}
) {
  return getTopics({ ...params, categoryId })
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
  categoryId: number
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

export function getCategory(categoryId: string | number) {
  return apiFetch<Category>("/api/topic/category", { params: { categoryId } })
}

export function getCategories() {
  return apiFetch<Category[]>("/api/topic/categories")
}

export function getCategoryNavs() {
  return apiFetch<Category[]>("/api/topic/category_navs")
}

export function getTag(id: string | number) {
  return apiFetch<Tag>(`/api/tag/${id}`)
}

export function searchTopics(params: SearchTopicParams) {
  return apiFetch<PageData<Topic>>("/api/search/topic", {
    params,
  })
}
