import { redirect, useLoaderData } from "react-router"

import { NodeTopicClientPage } from "@/components/topic/topic-dynamic-list-client-page"
import { rootDataFromMatches, categoryMeta } from "@/lib/seo"

import { getCurrentUser } from "../route-helpers/auth"
import {
  loadCategoryRouteData,
  type TopicListRouteData,
} from "../route-helpers/loaders"



type RouteArgs = {
  request: Request
  params: { id?: string }
}

export async function clientLoader({ request, params }: RouteArgs) {
  if (params.id === "feed") {
    const user = await getCurrentUser()
    if (!user) {
      throw redirect("/user/signin?redirect=/topics/category/feed")
    }
  }

  return loadCategoryRouteData({ request, id: params.id })
}

export function meta({
  data,
  location,
  matches,
}: {
  data?: TopicListRouteData
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return categoryMeta(
    rootDataFromMatches(matches)?.config,
    data?.category,
    location.pathname
  )
}

export default function CategoryRoute() {
  const data = useLoaderData() as TopicListRouteData
  return <NodeTopicClientPage initialData={data} />
}
