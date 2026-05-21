import { useLoaderData } from "react-router"

import { ArticleDetailClientPage } from "@/components/article/article-detail-client-page"
import { articleMeta, rootDataFromMatches } from "@/lib/seo"

import { loadArticleDetail } from "../route-helpers/loaders"

type RouteArgs = {
  request: Request
  params: { id?: string }
}

export async function loader({ request, params }: RouteArgs) {
  return loadArticleDetail({ request, id: params.id || "" })
}

export async function clientLoader({ params }: RouteArgs) {
  return loadArticleDetail({ id: params.id || "" })
}

export function meta({
  data,
  location,
  matches,
}: {
  data?: Awaited<ReturnType<typeof loader>>
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return articleMeta(
    rootDataFromMatches(matches)?.config,
    data,
    location.pathname
  )
}

export default function ArticleDetailRoute() {
  const article = useLoaderData<typeof loader>()
  return <ArticleDetailClientPage initialArticle={article} />
}
