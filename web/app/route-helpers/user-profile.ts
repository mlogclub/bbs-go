import { apiFetch } from "@/lib/api/client"
import type { UserSummary } from "@/lib/api/types"

export type UserProfileRouteData = {
  user: UserSummary
}

export async function loadUserProfileRouteData({
  request,
  userId,
}: {
  request?: Request
  userId: string
}): Promise<UserProfileRouteData> {
  const user = await apiFetch<UserSummary>(`/api/user/${userId}`, {
    ...(request ? { request } : {}),
  })

  return { user }
}
