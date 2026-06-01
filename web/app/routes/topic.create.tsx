import { Navigate, useLocation, useSearchParams } from "react-router"

import { useAppState, useAuthChecked } from "@/components/app/app-provider"
import { TopicCreateForm } from "@/components/topic/topic-create-form"
import { apiFetch } from "@/lib/api/client"
import type { Category } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { requireUser, requireUserClient } from "../route-helpers/auth"
import { useClientData } from "../route-helpers/client-hooks"

async function _loader(args: { request: Request }) {
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
  return noindexRouteMeta(matches, "Create topic", "发布话题")
}

export default function TopicCreateRoute() {
  const { t } = useI18n()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const { config, currentUser } = useAppState()
  const authChecked = useAuthChecked()
  const { data: categories } = useClientData<Category[]>("topic:categories", () =>
    apiFetch<Category[]>("/api/topic/categories").catch(() => [])
  )
  const contentType =
    searchParams.get("contentType") === "markdown" ||
    searchParams.get("contentType") === "text"
      ? searchParams.get("contentType")!
      : "html"
  const categoryId = Number(searchParams.get("categoryId") || 0)
  const type = Number(searchParams.get("type") || 0)
  const title =
    type === 1
      ? t("pages.topic.create.tweet")
      : type === 2
        ? t("pages.topic.create.qa")
        : t("pages.topic.create.post")
  useDocumentTitle(title)

  if (!authChecked) {
    return null
  }

  if (!currentUser) {
    const redirectTo = `${location.pathname}${location.search}`
    return (
      <Navigate
        to={`/user/signin?redirect=${encodeURIComponent(redirectTo)}`}
        replace
      />
    )
  }

  return (
    <main className="main">
      <div className="container">
        <TopicCreateForm
          key={`${type}:${contentType}:${categoryId}`}
          contentType={contentType as "html" | "markdown" | "text"}
          currentUser={currentUser}
          config={config}
          categoryId={categoryId}
          categories={categories || []}
          type={type}
        />
      </div>
    </main>
  )
}
