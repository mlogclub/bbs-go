import { useParams } from "react-router"

import { useAppState } from "@/components/app/app-provider"
import { EmptyState } from "@/components/common/empty-state"
import { TopicEditForm } from "@/components/topic/topic-edit-form"
import { apiFetch } from "@/lib/api/client"
import type { TopicEditData } from "@/lib/api/topics"
import type { TopicNode } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { requireUser, requireUserClient } from "../route-helpers/auth"
import { useClientData } from "../route-helpers/client-hooks"

export async function loader(args: { request: Request }) {
  await requireUser(args)
  return null
}

export async function clientLoader(args: { request: Request }) {
  await requireUserClient(args)
  return null
}

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Edit topic", "编辑话题")
}

export default function TopicEditRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("pages.topic.edit.title"))
  const { id = "" } = useParams()
  const { config } = useAppState()
  const { data, loading, error } = useClientData(`topic-edit:${id}`, async () => {
    const [topic, nodes] = await Promise.all([
      apiFetch<TopicEditData>(`/api/topic/edit/${id}`),
      apiFetch<TopicNode[]>("/api/topic/nodes").catch(() => []),
    ])
    return { topic, nodes }
  })

  if (loading) {
    return (
      <main className="main">
        <div className="container" />
      </main>
    )
  }
  if (error || !data) return <EmptyState title={error || "No data"} />

  return (
    <main className="main">
      <div className="container">
        <TopicEditForm topic={data.topic} config={config} nodes={data.nodes} />
      </div>
    </main>
  )
}
