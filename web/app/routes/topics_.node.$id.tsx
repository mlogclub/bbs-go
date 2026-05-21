import { redirect, useLoaderData } from "react-router"

import { NodeTopicClientPage } from "@/components/topic/topic-dynamic-list-client-page"
import { rootDataFromMatches, topicNodeMeta } from "@/lib/seo"

import { getCurrentUser } from "../route-helpers/auth"
import {
  loadTopicNodeRouteData,
  type TopicListRouteData,
} from "../route-helpers/loaders"

export { loader } from "../route-helpers/loaders"

type RouteArgs = {
  request: Request
  params: { id?: string }
}

export async function clientLoader({ request, params }: RouteArgs) {
  if (params.id === "feed") {
    const user = await getCurrentUser()
    if (!user) {
      throw redirect("/user/signin?redirect=/topics/node/feed")
    }
  }

  return loadTopicNodeRouteData({ request, id: params.id })
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
  return topicNodeMeta(
    rootDataFromMatches(matches)?.config,
    data?.node,
    location.pathname
  )
}

export default function TopicNodeRoute() {
  const data = useLoaderData() as TopicListRouteData
  return <NodeTopicClientPage initialData={data} />
}
