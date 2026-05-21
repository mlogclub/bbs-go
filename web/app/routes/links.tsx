import { useLoaderData } from "react-router"

import { EmptyState } from "@/components/common/empty-state"
import { WidgetCard } from "@/components/common/widget-card"
import { apiFetch } from "@/lib/api/client"
import type { FriendLink } from "@/lib/api/misc"
import { useI18n } from "@/lib/i18n/provider"
import { localizedTitle, pageMeta, rootDataFromMatches } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

function normalizeLinks(data: FriendLink[] | null | undefined) {
  return Array.isArray(data) ? data : []
}

export async function loader({ request }: { request: Request }) {
  return apiFetch<FriendLink[] | null>("/api/link/list", { request })
    .then(normalizeLinks)
    .catch(() => [])
}

export async function clientLoader() {
  return apiFetch<FriendLink[] | null>("/api/link/list")
    .then(normalizeLinks)
    .catch(() => [])
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
    localizedTitle(rootData?.locale, "Links", "友情链接"),
    { canonicalPath: location.pathname }
  )
}

export default function LinksRoute() {
  const links = useLoaderData<typeof loader>()
  const { t } = useI18n()
  useDocumentTitle(t("pages.links.title"), { appendSiteTitle: false })

  return (
    <section className="main">
      <div className="container">
        <WidgetCard title={t("pages.links.title")}>
          {links.length ? (
            <ul className="links px-[15px] py-2.5">
              {links.map((link) => (
                <li key={link.id} className="link">
                  <a
                    href={link.url || "#"}
                    target="_blank"
                    rel="noreferrer"
                    className="link-title"
                    title={link.title}
                  >
                    {link.title}
                  </a>
                  {link.summary ? (
                    <p className="link-summary">{link.summary}</p>
                  ) : null}
                </li>
              ))}
            </ul>
          ) : (
            <EmptyState title={t("common.noData")} />
          )}
        </WidgetCard>
      </div>
    </section>
  )
}
