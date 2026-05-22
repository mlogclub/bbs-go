import { useSearchParams } from "react-router"

import { ErrorPage } from "@/components/common/error-page"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Redirect", "跳转")
}

export default function RedirectRoute() {
  const [searchParams] = useSearchParams()
  const url = searchParams.get("url") || ""
  const { t } = useI18n()
  if (
    !url.toLowerCase().startsWith("http://") &&
    !url.toLowerCase().startsWith("https://")
  ) {
    return <ErrorPage message={t("pages.redirect.error")} />
  }
  return (
    <section className="main">
      <div className="container">
        <div className="main-body redirect py-[100px] text-center">
          <a href={url} rel="nofollow">
            {t("pages.redirect.link")}
          </a>
        </div>
      </div>
    </section>
  )
}
