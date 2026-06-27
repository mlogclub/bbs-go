import { useLoaderData } from "react-router"

import { EmptyState } from "@/components/common/empty-state"
import { LoadMore } from "@/components/common/load-more"
import { HomeAside } from "@/components/layout/home-aside"
import { MainShell } from "@/components/layout/main-shell"
import { TopicFeedTabs } from "@/components/topic/topic-feed-tabs"
import { TopicListItem } from "@/components/topic/topic-list-item"
import { TopicsNavContent } from "@/components/topic/topics-nav-content"
import { apiFetch } from "@/lib/api/client"
import type { PageData, Topic } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { rootDataFromMatches, siteHomeMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import {
  loadTopicListRouteData,
  type TopicListRouteData,
} from "../route-helpers/loaders"

export { loader } from "../route-helpers/loaders"

export async function clientLoader() {
  return loadTopicListRouteData()
}

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return siteHomeMeta(rootDataFromMatches(matches)?.config)
}

export function TopicListRoute({ title }: { title?: string }) {
  const { topics, categories } = useLoaderData() as TopicListRouteData
  const { t } = useI18n()
  useDocumentTitle(title)

  return (
    <MainShell aside={<HomeAside />}>
      <div className="topics-wrapper">
        <TopicsNavContent initialCategories={categories} currentCategoryId={0} />
        <div className="topics-main">
          <div className="rounded-lg bg-background">
            <TopicFeedTabs currentCategoryId={0} />
            <LoadMore<Topic>
              initialItems={topics.results}
              initialCursor={topics.cursor || ""}
              initialHasMore={topics.hasMore}
              initialLoad={false}
              resetKey="/api/topic/topics"
              labels={{
                loadMore: t("common.loadMore.loadMore"),
                noMore: t("common.loadMore.noMore"),
              }}
              loadPage={({ cursor }) =>
                apiFetch<PageData<Topic>>("/api/topic/topics", {
                  params: { cursor },
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

export default function IndexRoute() {
  return <TopicListRoute />
}
