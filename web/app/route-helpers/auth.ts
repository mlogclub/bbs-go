import { redirect, type RouterContextProvider } from "react-router"

import { apiFetch } from "@/lib/api/client"
import type { UserSummary } from "@/lib/api/types"
import { userCanAccessDashboard } from "@/lib/auth/roles"

import { rootDataContext } from "./context"

interface RequireUserArgs {
  request: Request
  context?: Pick<RouterContextProvider, "get">
}

export function buildSigninRedirect(requestOrUrl: Request | URL | string) {
  const url =
    typeof requestOrUrl === "string"
      ? new URL(requestOrUrl, "http://local")
      : requestOrUrl instanceof Request
        ? new URL(requestOrUrl.url)
        : requestOrUrl

  return `/user/signin?redirect=${encodeURIComponent(url.pathname + url.search)}`
}

export async function getCurrentUser(request?: Request) {
  return apiFetch<UserSummary>("/api/user/current", { request }).catch(
    () => null
  )
}

export async function requireUser({ request, context }: RequireUserArgs) {
  const getRootData = context?.get(rootDataContext)
  const user = getRootData
    ? (await getRootData()).currentUser
    : await getCurrentUser(request)

  if (!user) {
    throw redirect(buildSigninRedirect(request))
  }
  return user
}

export async function requireUserClient({ request }: { request: Request }) {
  const user = await getCurrentUser()
  if (!user) {
    throw redirect(buildSigninRedirect(request))
  }
  return user
}

function dashboardForbidden() {
  return new Response("No permission to access dashboard", {
    status: 403,
    statusText: "Forbidden",
  })
}

export async function requireDashboardAdmin(args: RequireUserArgs) {
  const user = await requireUser(args)
  if (!userCanAccessDashboard(user)) {
    throw dashboardForbidden()
  }
  return user
}

export async function requireDashboardAdminClient({
  request,
}: {
  request: Request
}) {
  const user = await requireUserClient({ request })
  if (!userCanAccessDashboard(user)) {
    throw dashboardForbidden()
  }
  return user
}
