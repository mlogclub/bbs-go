import { useLoaderData } from "react-router"

import { UserProfileClientPage } from "@/components/user/user-profile-client-page"
import { rootDataFromMatches, userMeta } from "@/lib/seo"

import {
  loadUserProfileRouteData,
  type UserProfileRouteData,
} from "../route-helpers/user-profile"

type RouteArgs = {
  request: Request
  params: { userId?: string }
}

async function _loader({ request, params }: RouteArgs) {
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
  return userMeta(
    rootDataFromMatches(matches)?.config,
    data?.user,
    undefined,
    location.pathname
  )
}

export default function UserProfileRoute() {
  const { user } = useLoaderData<typeof loader>()
  return <UserProfileClientPage initialUser={user} />
}
