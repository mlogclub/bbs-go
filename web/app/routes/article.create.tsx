import { useAppState } from "@/components/app/app-provider"
import { ArticleForm } from "@/components/article/article-form"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { requireUser, requireUserClient } from "../route-helpers/auth"

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
  return noindexRouteMeta(matches, "Create article", "发布文章")
}

export default function ArticleCreateRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("pages.article.create.title"))
  const { config } = useAppState()
  return (
    <main className="main">
      <div className="container">
        <ArticleForm mode="create" config={config} />
      </div>
    </main>
  )
}
