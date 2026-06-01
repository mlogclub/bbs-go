import { useLoaderData } from "react-router"

import { ArticleList } from "@/components/article/article-list"
import { EmptyState } from "@/components/common/empty-state"
import { LoadMore } from "@/components/common/load-more"
import { HomeAside } from "@/components/layout/home-aside"
import { MainShell } from "@/components/layout/main-shell"
import { apiFetch } from "@/lib/api/client"
import type { Article, PageData } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { localizedTitle, rootDataFromMatches, sitePageMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { loadArticles } from "../route-helpers/loaders"



export async function clientLoader() {
  return loadArticles({})
}

export function meta({
  location,
  matches,
}: {
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  return sitePageMeta(
    rootData?.config,
    localizedTitle(rootData?.locale, "Articles", "文章"),
    { canonicalPath: location.pathname }
  )
}

export default function ArticlesRoute() {
  const articles = useLoaderData() as PageData<Article>
  const { t } = useI18n()
  useDocumentTitle(t("pages.articles.title"))

  return (
    <MainShell aside={<HomeAside />}>
      <div className="overflow-hidden rounded-lg bg-background">
        <LoadMore<Article>
          initialItems={articles.results}
          initialCursor={articles.cursor || ""}
          initialHasMore={articles.hasMore}
          initialLoad={false}
          resetKey="/api/article/articles"
          labels={{
            loadMore: t("common.loadMore.loadMore"),
            noMore: t("common.loadMore.noMore"),
          }}
          loadPage={({ cursor }) =>
            apiFetch<PageData<Article>>("/api/article/articles", {
              params: { cursor },
            })
          }
          renderItems={(items) => <ArticleList articles={items} t={t} />}
          renderEmpty={() => <EmptyState title={t("common.noData")} />}
        />
      </div>
    </MainShell>
  )
}
