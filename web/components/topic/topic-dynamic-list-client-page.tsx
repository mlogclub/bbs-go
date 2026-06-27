"use client"

import * as React from "react"
import { useSearchParams } from "react-router-dom"

import { EmptyState } from "@/components/common/empty-state"
import { LoadMore } from "@/components/common/load-more"
import { HomeAside } from "@/components/layout/home-aside"
import { MainShell } from "@/components/layout/main-shell"
import { PageLoading } from "@/components/common/page-state"
import { TopicFeedTabs } from "@/components/topic/topic-feed-tabs"
import { TopicListItem } from "@/components/topic/topic-list-item"
import { TopicsNavContent } from "@/components/topic/topics-nav-content"
import { TopicSubCategoryNav } from "@/components/topic/topic-sub-category-nav"
import { apiFetch } from "@/lib/api/client"
import type { PageData, Tag, Topic, Category } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { useRouteData, useRouteSegment } from "@/lib/spa-route"
import { useDocumentTitle } from "@/lib/use-document-title"
import { cn } from "@/lib/utils"

function resolveCategoryId(id: string) {
  if (id === "newest") return 0
  if (id === "recommend") return -1
  if (id === "feed") return -2
  const parsed = Number.parseInt(id, 10)
  return Number.isNaN(parsed) ? 0 : parsed
}

const qaStatusOptions = ["", "unsolved", "solved"]
const sortOptions = ["latestPublish", "latestReply"]

type TopicListInitialData = {
  topics?: PageData<Topic>
  categories?: Category[]
}

export function TopicTagClientPage({
  initialData,
}: {
  initialData?: TopicListInitialData
}) {
  const tagId = useRouteSegment(2)
  const { t } = useI18n()
  const load = React.useCallback(
    () => apiFetch<Tag>(`/api/tag/${tagId}`).catch(() => ({ id: 0, name: "" })),
    [tagId]
  )
  const { data: tag, loading } = useRouteData(`topic-tag:${tagId}`, load)
  const currentTag = String(tag?.id) === tagId ? tag : undefined
  useDocumentTitle(currentTag?.name, t("pages.topics.title"))
  const labels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  return (
    <MainShell aside={<HomeAside />}>
      <div className="topics-wrapper">
        <TopicsNavContent initialCategories={initialData?.categories || []} />
        <div className="topics-main">
          <div className="rounded-lg bg-background">
            {loading ? <PageLoading /> : null}
            <LoadMore<Topic>
              initialItems={initialData?.topics?.results || []}
              initialCursor={initialData?.topics?.cursor || "0"}
              initialHasMore={initialData?.topics?.hasMore || false}
              initialLoad={!initialData?.topics}
              resetKey={`topic-tag:${tagId}`}
              labels={labels}
              loadPage={({ cursor }) =>
                apiFetch<PageData<Topic>>("/api/topic/tag/topics", {
                  params: { tagId, cursor },
                })
              }
              renderItems={(items) => (
                <ul className="divide-y divide-border">
                  {items.map((topic) => (
                    <TopicListItem key={topic.id} topic={topic} t={t} />
                  ))}
                </ul>
              )}
              renderEmpty={() => <EmptyState title={t("common.noData")} />}
            />
          </div>
          {currentTag?.name ? <span className="sr-only">{currentTag.name}</span> : null}
        </div>
      </div>
    </MainShell>
  )
}

export function NodeTopicClientPage({
  initialData,
}: {
  initialData?: TopicListInitialData
}) {
  const id = useRouteSegment(2)
  const categoryId = resolveCategoryId(id)
  const [searchParams, setSearchParams] = useSearchParams()
  const { t } = useI18n()
  const load = React.useCallback(
    () =>
      apiFetch<Category>("/api/topic/category", { params: { categoryId } }).catch(
        (): Category => ({ id: categoryId, name: "", children: [] })
      ),
    [categoryId]
  )
  const { data: node } = useRouteData(`category:${id}`, load)
  const currentNode =
    categoryId > 0 && String(node?.id) === String(categoryId) ? node : undefined
  const rootCategoryId = currentNode?.parentId || currentNode?.id || categoryId
  const loadRootNode = React.useCallback(() => {
    if (categoryId <= 0) return Promise.resolve<Category | null>(null)
    if (!currentNode?.parentId) return Promise.resolve(currentNode || null)
    return apiFetch<Category>("/api/topic/category", {
      params: { categoryId: currentNode.parentId },
    }).catch(() => null)
  }, [currentNode, categoryId])
  const { data: rootNode } = useRouteData<Category | null>(
    categoryId > 0 ? `topic-root-node:${rootCategoryId}` : "",
    loadRootNode,
    null
  )
  const currentRootNode =
    categoryId > 0 && String(rootNode?.id) === String(rootCategoryId)
      ? rootNode
      : currentNode?.id === rootCategoryId
        ? currentNode
        : undefined
  const subNodes =
    categoryId > 0 && currentRootNode?.children?.length
      ? currentRootNode.children
      : []
  const isQaNode = categoryId > 0 && currentNode?.type === "qa"
  const isNormalNode = categoryId > 0 && currentNode && currentNode.type !== "qa"
  const qaStatusValue = searchParams.get("qaStatus") || ""
  const qaStatus = qaStatusOptions.includes(qaStatusValue) ? qaStatusValue : ""
  const sortValue = searchParams.get("sort") || ""
  const normalSort = sortOptions.includes(sortValue) ? sortValue : "latestPublish"
  const currentFilters = isQaNode
    ? [
        { value: "", label: t("pages.qa.filterAll") },
        { value: "unsolved", label: t("pages.qa.filterUnsolved") },
        { value: "solved", label: t("pages.qa.filterSolved") },
      ]
    : isNormalNode
      ? [
          {
            value: "latestPublish",
            label: t("pages.topics.filterLatestPublish"),
          },
          {
            value: "latestReply",
            label: t("pages.topics.filterLatestReply"),
          },
        ]
      : []
  const currentFilterValue = isQaNode ? qaStatus : isNormalNode ? normalSort : ""
  const currentFilterLabel =
    currentFilters.find((item) => item.value === currentFilterValue)?.label || ""
  useDocumentTitle(currentNode?.name, t("pages.topics.title"))
  const labels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  function switchFilter(value: string) {
    const next = new URLSearchParams(searchParams)
    if (isQaNode) {
      next.delete("sort")
      if (value) {
        next.set("qaStatus", value)
      } else {
        next.delete("qaStatus")
      }
    } else if (isNormalNode) {
      next.delete("qaStatus")
      next.set("sort", value)
    }
    setSearchParams(next, { replace: true })
  }

  return (
    <MainShell aside={<HomeAside />}>
      <div className="topics-wrapper">
        <TopicsNavContent
          initialCategories={initialData?.categories || []}
          currentCategoryId={categoryId}
          currentRootCategoryId={rootCategoryId}
        />
        <div className="topics-main">
          {currentNode?.name ? (
            <div className="mb-3 rounded-lg bg-background px-4 py-4">
              <div className="flex items-start gap-3">
                {currentNode.logo ? (
                  <img
                    src={currentNode.logo}
                    alt={currentNode.name || ""}
                    className="h-10 w-10 rounded-md object-cover"
                  />
                ) : null}
                <div className="min-w-0 flex-1">
                  <div className="text-[16px] leading-snug font-semibold text-foreground">
                    {currentNode.name}
                  </div>
                  {currentNode.description ? (
                    <div className="mt-1 line-clamp-2 text-[14px] leading-snug text-muted-foreground">
                      {currentNode.description}
                    </div>
                  ) : null}
                </div>
              </div>
            </div>
          ) : null}
          <div className="rounded-lg bg-background">
            {categoryId <= 0 ? (
              <TopicFeedTabs currentCategoryId={categoryId} />
            ) : null}
            <TopicSubCategoryNav
              rootCategoryId={rootCategoryId}
              categories={subNodes}
              currentCategoryId={categoryId}
            />
            {currentFilters.length > 0 ? (
              <div className="flex justify-between border-b border-border px-4 py-3">
                <div className="text-base font-bold">
                  {currentFilterLabel}
                </div>
                <div className="inline-flex flex-wrap items-center gap-1 rounded-lg bg-muted p-1">
                  {currentFilters.map((item) => {
                    const selected = item.value === currentFilterValue
                    return (
                      <button
                        key={item.value || "all"}
                        type="button"
                        className={cn(
                          "inline-flex h-5 items-center rounded-md px-3 text-sm font-medium transition-colors",
                          selected
                            ? "bg-background text-foreground shadow-sm"
                            : "text-muted-foreground hover:text-foreground"
                        )}
                        aria-pressed={selected}
                        onClick={() => switchFilter(item.value)}
                      >
                        {item.label}
                      </button>
                    )
                  })}
                </div>
              </div>
            ) : null}
            <LoadMore<Topic>
              initialItems={initialData?.topics?.results || []}
              initialCursor={initialData?.topics?.cursor || "0"}
              initialHasMore={initialData?.topics?.hasMore || false}
              initialLoad={!initialData?.topics}
              resetKey={`category:${categoryId}:${currentNode?.type || ""}:${currentFilterValue}`}
              labels={labels}
              loadPage={({ cursor }) =>
                apiFetch<PageData<Topic>>("/api/topic/topics", {
                  params: {
                    categoryId,
                    cursor,
                    ...(isQaNode && qaStatus ? { qaStatus } : {}),
                    ...(isNormalNode ? { sort: normalSort } : {}),
                  },
                })
              }
              renderItems={(items) => (
                <ul className="divide-y divide-border">
                  {items.map((topic) => (
                    <TopicListItem
                      key={topic.id}
                      topic={topic}
                      showSticky
                      t={t}
                    />
                  ))}
                </ul>
              )}
              renderEmpty={() => <EmptyState title={t("common.noData")} />}
            />
          </div>
        </div>
      </div>
    </MainShell>
  )
}
