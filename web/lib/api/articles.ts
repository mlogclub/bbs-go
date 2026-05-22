import { serverApiFetch as apiFetch } from "./server"

import type { Article, ArticleEditForm, PageData, SearchArticle } from "./types"

type SearchArticleParams = {
  keyword: string
  timeRange?: number
  cursor?: string
}

export function getArticles(cursor?: string) {
  return apiFetch<PageData<Article>>("/api/article/articles", {
    params: { cursor },
  })
}

export function getArticle(id: string | number) {
  return apiFetch<Article>(`/api/article/${id}`)
}

export function getArticleEdit(id: string | number) {
  return apiFetch<ArticleEditForm>(`/api/article/edit/${id}`)
}

export function getTagArticles(tagId: string | number, cursor?: string) {
  return apiFetch<PageData<Article>>("/api/article/tag/articles", {
    params: { tagId, cursor },
  })
}

export function searchArticles(params: SearchArticleParams) {
  return apiFetch<PageData<SearchArticle>>("/api/search/article", {
    params,
  })
}
