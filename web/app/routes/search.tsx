import * as React from "react"
import { Search } from "lucide-react"
import { useSearchParams } from "react-router"

import { EmptyState } from "@/components/common/empty-state"
import { LoadMore } from "@/components/common/load-more"
import { HomeAside } from "@/components/layout/home-aside"
import { MainShell } from "@/components/layout/main-shell"
import { SearchArticleList } from "@/components/search/search-article-list"
import { SearchFilters } from "@/components/search/search-filters"
import { SearchTopicList } from "@/components/search/search-topic-list"
import { SearchUserList } from "@/components/search/search-user-list"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { apiFetch } from "@/lib/api/client"
import type {
  PageData,
  SearchArticle,
  SearchUser,
  Topic,
  Category,
} from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { localizedTitle, noindexMeta, rootDataFromMatches } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { useClientData, useLoadMoreLabels } from "../route-helpers/client-hooks"

type SearchNodeOption = {
  id: number
  label: string
}

type SearchType = "topic" | "article" | "user"

type MetaLocation = {
  search?: string
}

const searchTypes: SearchType[] = ["topic", "article", "user"]

function flattenSearchNodes(
  categories: Category[] = [],
  prefix = ""
): SearchNodeOption[] {
  return categories.flatMap((category) => {
    const label = prefix ? `${prefix} / ${category.name}` : category.name || String(category.id)
    return [
      { id: Number(category.id), label },
      ...flattenSearchNodes(category.children || [], label),
    ]
  })
}

function normalizeSearchType(value: string | null): SearchType {
  return searchTypes.includes(value as SearchType)
    ? (value as SearchType)
    : "topic"
}

export function meta({
  location,
  matches,
}: {
  location?: MetaLocation
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  const searchParams = new URLSearchParams(location?.search || "")
  const keyword = searchParams.get("q") || searchParams.get("keyword") || ""
  const title = keyword
    ? `${localizedTitle(rootData?.locale, "Search", "搜索")}: ${keyword}`
    : localizedTitle(rootData?.locale, "Search", "搜索")

  return noindexMeta(rootData?.config, title)
}

export default function SearchRoute() {
  const [searchParams, setSearchParams] = useSearchParams()
  const { t } = useI18n()
  const labels = useLoadMoreLabels()
  useDocumentTitle(t("pages.search.title"))

  const keyword = searchParams.get("q") || searchParams.get("keyword") || ""
  const categoryId = Number(searchParams.get("categoryId") || 0)
  const timeRange = Number(searchParams.get("timeRange") || 0)
  const type = normalizeSearchType(searchParams.get("type"))
  const [searchKeyword, setSearchKeyword] = React.useState(keyword)
  const { data: categories } = useClientData<Category[]>("search:categories", () =>
    apiFetch<Category[]>("/api/topic/categories").catch(() => [])
  )

  React.useEffect(() => {
    setSearchKeyword(keyword)
  }, [keyword])

  function submitSearch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const nextParams = new URLSearchParams(searchParams.toString())
    const nextKeyword = searchKeyword.trim()
    if (nextKeyword) {
      nextParams.set("q", nextKeyword)
    } else {
      nextParams.delete("q")
      nextParams.delete("keyword")
    }
    setSearchParams(nextParams)
  }

  function setType(nextType: string) {
    const nextParams = new URLSearchParams(searchParams.toString())
    if (nextType === "topic") {
      nextParams.delete("type")
    } else {
      nextParams.set("type", nextType)
      if (nextType === "user") {
        nextParams.delete("categoryId")
        nextParams.delete("timeRange")
      }
      if (nextType === "article") {
        nextParams.delete("categoryId")
      }
    }
    setSearchParams(nextParams)
  }

  const emptyTitle = keyword.trim()
    ? t("pages.search.empty.title")
    : t("pages.search.empty.promptTitle")
  const emptyDescription = keyword.trim()
    ? t("pages.search.empty.description")
    : t("pages.search.empty.promptDescription")

  return (
    <MainShell aside={<HomeAside />}>
      <section
        className="mx-auto mb-4 w-full rounded-lg border border-border/70 bg-background/95 p-4 shadow-sm sm:p-5"
        style={{ maxWidth: "calc(100vw - 1rem)" }}
      >
        <div className="mb-4">
          <h1 className="text-xl font-semibold tracking-normal text-foreground">
            {t("pages.search.heading")}
          </h1>
          <p className="mt-1 text-sm leading-6 text-muted-foreground">
            {t("pages.search.description")}
          </p>
        </div>
        <form onSubmit={submitSearch}>
          <div className="flex flex-col gap-2 sm:flex-row">
            <div className="flex h-11 min-w-0 flex-1 items-center rounded-lg border border-input bg-muted/30 px-3 transition-colors focus-within:border-ring focus-within:bg-background focus-within:ring-3 focus-within:ring-ring/20">
              <Search className="mr-2 size-4 shrink-0 text-muted-foreground" />
              <Input
                value={searchKeyword}
                type="text"
                maxLength={50}
                placeholder={t("pages.search.placeholder")}
                className="h-auto border-none bg-transparent px-0 text-base shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 md:text-sm"
                autoComplete="off"
                onChange={(event) => setSearchKeyword(event.currentTarget.value)}
              />
            </div>
            <Button type="submit" className="h-11 rounded-lg px-6 sm:min-w-24">
              {t("pages.search.submit")}
            </Button>
          </div>
        </form>
      </section>

      <section
        className="mx-auto w-full overflow-hidden rounded-lg border border-border/70 bg-background/95 shadow-sm"
        style={{ maxWidth: "calc(100vw - 1rem)" }}
      >
        <div className="border-b border-border/70 px-4 py-3 sm:px-5">
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <Tabs value={type} onValueChange={setType} className="w-full md:w-auto">
              <TabsList className="grid w-full grid-cols-3 sm:inline-flex sm:w-fit">
              {searchTypes.map((item) => (
                <TabsTrigger key={item} value={item} className="min-w-0 px-3">
                  {t(`pages.search.tabs.${item}`)}
                </TabsTrigger>
              ))}
              </TabsList>
            </Tabs>
            {type === "topic" || type === "article" ? (
              <SearchFilters
                categories={flattenSearchNodes(categories || [])}
                showNode={type === "topic"}
                showTime
              />
            ) : null}
          </div>
        </div>

        {type === "topic" ? (
          <LoadMore<Topic>
            initialCursor=""
            initialHasMore
            initialLoad
            resetKey={`search-topic:${keyword}:${categoryId}:${timeRange}`}
            labels={labels}
            loadPage={({ cursor }) =>
              apiFetch<PageData<Topic>>("/api/search/topic", {
                params: { keyword, categoryId, timeRange, cursor },
              })
            }
            renderItems={(items) => <SearchTopicList results={items} />}
            renderEmpty={() => (
              <EmptyState
                title={emptyTitle}
                description={emptyDescription}
                className="py-12"
              />
            )}
          />
        ) : type === "article" ? (
          <LoadMore<SearchArticle>
            initialCursor=""
            initialHasMore
            initialLoad
            resetKey={`search-article:${keyword}:${timeRange}`}
            labels={labels}
            loadPage={({ cursor }) =>
              apiFetch<PageData<SearchArticle>>("/api/search/article", {
                params: { keyword, timeRange, cursor },
              })
            }
            renderItems={(items) => <SearchArticleList results={items} />}
            renderEmpty={() => (
              <EmptyState
                title={emptyTitle}
                description={emptyDescription}
                className="py-12"
              />
            )}
          />
        ) : (
          <LoadMore<SearchUser>
            initialCursor=""
            initialHasMore
            initialLoad
            resetKey={`search-user:${keyword}`}
            labels={labels}
            loadPage={({ cursor }) =>
              apiFetch<PageData<SearchUser>>("/api/search/user", {
                params: { keyword, cursor },
              })
            }
            renderItems={(items) => <SearchUserList results={items} />}
            renderEmpty={() => (
              <EmptyState
                title={emptyTitle}
                description={emptyDescription}
                className="py-12"
              />
            )}
          />
        )}
      </section>
    </MainShell>
  )
}
