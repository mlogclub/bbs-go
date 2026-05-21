import { useLoaderData } from "react-router"

import { EmptyState } from "@/components/common/empty-state"
import { apiFetch } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"
import { localizedTitle, pageMeta, rootDataFromMatches } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

const emptyAbout: Record<string, string> = {}

export async function loader({ request }: { request: Request }) {
  return apiFetch<Record<string, string>>("/api/config/about", {
    request,
  }).catch(() => emptyAbout)
}

export async function clientLoader() {
  return apiFetch<Record<string, string>>("/api/config/about").catch(
    () => emptyAbout
  )
}

export function meta({
  location,
  matches,
}: {
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  return pageMeta(
    rootData?.config,
    localizedTitle(rootData?.locale, "About", "关于"),
    { canonicalPath: location.pathname }
  )
}

export default function AboutRoute() {
  const data = useLoaderData<typeof loader>()
  const { t } = useI18n()
  useDocumentTitle(t("pages.about.title"), { appendSiteTitle: false })
  const html = data.content || ""

  return (
    <section className="main">
      <div className="container">
        <div className="rounded-md bg-card px-8 py-3">
          {html ? (
            <div
              className="bbs-content"
              dangerouslySetInnerHTML={{ __html: html }}
            />
          ) : (
            <EmptyState title={t("common.noData")} />
          )}
        </div>
      </div>
    </section>
  )
}
