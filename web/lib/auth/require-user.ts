import { redirect } from "@/lib/router/navigation"

import { getAppState } from "@/lib/app-state/server"
import type { UserSummary } from "@/lib/api/types"

export async function requireUser(
  redirectPath: string
): Promise<UserSummary | null> {
  const { currentUser } = await getAppState()

  if (currentUser) {
    return currentUser
  }

  redirect(`/user/signin?redirect=${encodeURIComponent(redirectPath)}`)
}
