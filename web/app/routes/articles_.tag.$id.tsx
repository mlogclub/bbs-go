import { useLoaderData } from "react-router"

import { ArticleTagClientPage } from "@/components/article/article-tag-client-page"
import { localizedTitle, rootDataFromMatches, tagPageMeta } from "@/lib/seo"

import {
  loadArticleListRouteData,
  type ArticleListRouteData,
} from "../route-helpers/loaders"



type RouteArgs = {
  request: Request
  params: { id?: string }
}

export async function clientLoader({ request, params }: RouteArgs) {
  return loadArticleListRouteData({ request, tagId: params.id })
}

export function meta({
  data,
  location,
  matches,
}: {
  data?: ArticleListRouteData
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  return tagPageMeta(
    rootData?.config,
    data?.tag,
    localizedTitle(rootData?.locale, "Articles", "文章"),
    location.pathname
  )
}

export default function ArticleTagRoute() {
  const articles = useLoaderData() as ArticleListRouteData
  return <ArticleTagClientPage initialData={articles} />
}
