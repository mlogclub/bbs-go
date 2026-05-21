"use client"

import * as React from "react"

import { ArticleList } from "@/components/article/article-list"
import { EmptyState } from "@/components/common/empty-state"
import { LoadMore } from "@/components/common/load-more"
import { HomeAside } from "@/components/layout/home-aside"
import { MainShell } from "@/components/layout/main-shell"
import { PageLoading } from "@/components/common/page-state"
import { apiFetch } from "@/lib/api/client"
import type { Article, PageData, Tag } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { useRouteData, useRouteSegment } from "@/lib/spa-route"
import { useDocumentTitle } from "@/lib/use-document-title"

export function ArticleTagClientPage({
  initialData,
}: {
  initialData?: PageData<Article>
}) {
  const tagId = useRouteSegment(2)
  const { t } = useI18n()
  const load = React.useCallback(
    () => apiFetch<Tag>(`/api/tag/${tagId}`).catch(() => ({ id: 0, name: "" })),
    [tagId]
  )
  const { data: tag, loading } = useRouteData(`article-tag:${tagId}`, load)
  const currentTag = String(tag?.id) === tagId ? tag : undefined
  useDocumentTitle(currentTag?.name, t("pages.articles.title"))
  const labels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  return (
    <MainShell aside={<HomeAside />}>
      <div className="rounded-lg bg-background">
        {loading ? <PageLoading /> : null}
        {currentTag?.name ? (
          <div className="border-b px-4 py-3 text-base font-semibold">
            {currentTag.name}
          </div>
        ) : null}
        <LoadMore<Article>
          initialItems={initialData?.results || []}
          initialCursor={initialData?.cursor || "0"}
          initialHasMore={initialData?.hasMore || false}
          initialLoad={!initialData}
          resetKey={`article-tag:${tagId}`}
          labels={labels}
          loadPage={({ cursor }) =>
            apiFetch<PageData<Article>>("/api/article/tag/articles", {
              params: { tagId, cursor },
            })
          }
          renderItems={(items) => <ArticleList articles={items} t={t} />}
          renderEmpty={() => <EmptyState title={t("common.noData")} />}
        />
      </div>
    </MainShell>
  )
}
