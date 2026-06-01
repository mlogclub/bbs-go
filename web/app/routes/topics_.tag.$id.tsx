import { useLoaderData } from "react-router"

import { TopicTagClientPage } from "@/components/topic/topic-dynamic-list-client-page"
import { localizedTitle, rootDataFromMatches, tagPageMeta } from "@/lib/seo"

import {
  loadTopicTagRouteData,
  type TopicListRouteData,
} from "../route-helpers/loaders"



type RouteArgs = {
  request: Request
  params: { id?: string }
}

export async function clientLoader({ request, params }: RouteArgs) {
  return loadTopicTagRouteData({ request, id: params.id })
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
  const rootData = rootDataFromMatches(matches)
  return tagPageMeta(
    rootData?.config,
    data?.tag,
    localizedTitle(rootData?.locale, "Topics", "话题"),
    location.pathname
  )
}

export default function TopicTagRoute() {
  const data = useLoaderData() as TopicListRouteData
  return <TopicTagClientPage initialData={data} />
}
