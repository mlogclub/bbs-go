import { useI18n } from "@/lib/i18n/provider"
import { localizedTitle, rootDataFromMatches, sitePageMeta } from "@/lib/seo"

import { TopicListRoute } from "./_index"

export { clientLoader } from "./_index"
export { loader } from "../route-helpers/loaders"

export function meta({
  location,
  matches,
}: {
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  return sitePageMeta(
    rootData?.config,
    localizedTitle(rootData?.locale, "Topics", "话题"),
    { canonicalPath: location.pathname }
  )
}

export default function TopicsRoute() {
  const { t } = useI18n()
  return <TopicListRoute title={t("pages.topics.title")} />
}
