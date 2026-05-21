import { ErrorPage } from "@/components/common/error-page"
import { noindexRouteMeta } from "@/lib/seo"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Not found", "页面不存在")
}

export default function NotFoundRoute() {
  return <ErrorPage statusCode={404} />
}
