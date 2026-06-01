import { useLoaderData } from "react-router"

import { TopicDetailClientPage } from "@/components/topic/topic-detail-client-page"
import { rootDataFromMatches, topicMeta } from "@/lib/seo"

import { loadTopicDetail } from "../route-helpers/loaders"

type RouteArgs = {
  request: Request
  params: { id?: string }
}

async function _loader({ request, params }: RouteArgs) {
  return loadTopicDetail({ request, id: params.id || "" })
}

export async function clientLoader({ params }: RouteArgs) {
  return loadTopicDetail({ id: params.id || "" })
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
  return topicMeta(rootDataFromMatches(matches)?.config, data, location.pathname)
}

export default function TopicDetailRoute() {
  const topic = useLoaderData<typeof loader>()
  return <TopicDetailClientPage initialTopic={topic} />
}
