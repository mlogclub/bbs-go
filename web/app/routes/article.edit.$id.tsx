import { useParams } from "react-router"

import { useAppState } from "@/components/app/app-provider"
import { ArticleForm } from "@/components/article/article-form"
import { EmptyState } from "@/components/common/empty-state"
import { apiFetch } from "@/lib/api/client"
import type { ArticleEditForm } from "@/lib/api/types"
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
  return noindexRouteMeta(matches, "Edit article", "编辑文章")
}

export default function ArticleEditRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("pages.article.edit.title"))
  const { id = "" } = useParams()
  const { config } = useAppState()
  const { data: article, loading, error } = useClientData(
    `article-edit:${id}`,
    () => apiFetch<ArticleEditForm>(`/api/article/edit/${id}`)
  )

  if (loading) {
    return (
      <main className="main">
        <div className="container" />
      </main>
    )
  }
  if (error) return <EmptyState title={error} />

  return (
    <main className="main">
      <div className="container">
        <ArticleForm mode="edit" config={config} initialArticle={article} />
      </div>
    </main>
  )
}
