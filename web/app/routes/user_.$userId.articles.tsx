import { UserArticlesClientPage } from "@/components/user/user-profile-client-page"
import { localizedTitle, rootDataFromMatches, userMeta } from "@/lib/seo"

import {
  loadUserProfileRouteData,
  type UserProfileRouteData,
} from "../route-helpers/user-profile"

type RouteArgs = {
  request: Request
  params: { userId?: string }
}

export async function loader({ request, params }: RouteArgs) {
  return loadUserProfileRouteData({ request, userId: params.userId || "" })
}

export async function clientLoader({ params }: RouteArgs) {
  return loadUserProfileRouteData({ userId: params.userId || "" })
}

export function meta({
  data,
  location,
  matches,
}: {
  data?: UserProfileRouteData
  location: { pathname: string }
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  return userMeta(
    rootData?.config,
    data?.user,
    localizedTitle(rootData?.locale, "Articles", "文章"),
    location.pathname
  )
}

export default UserArticlesClientPage
