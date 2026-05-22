"use client"

import * as React from "react"
import { useRouter } from "@/lib/router/navigation"

import { useAuthChecked, useCurrentUser } from "@/components/app/app-provider"
import type { UserSummary } from "@/lib/api/types"
import { buildSigninHref } from "@/lib/toast"

const RequiredUserContext = React.createContext<UserSummary | null>(null)

export function RequireUser({
  initialUser,
  redirectPath,
  children,
}: {
  initialUser: UserSummary | null
  redirectPath: string
  children: React.ReactNode
}) {
  const router = useRouter()
  const authChecked = useAuthChecked()
  const clientUser = useCurrentUser()
  const user = clientUser || initialUser
  const checked = Boolean(initialUser) || authChecked

  React.useEffect(() => {
    if (checked && !user) {
      router.replace(buildSigninHref(redirectPath))
    }
  }, [checked, redirectPath, router, user])

  if (!checked || !user) {
    return null
  }

  return (
    <RequiredUserContext.Provider value={user}>
      {children}
    </RequiredUserContext.Provider>
  )
}

export function useRequiredUser() {
  const user = React.useContext(RequiredUserContext)

  if (!user) {
    throw new Error("useRequiredUser must be used within RequireUser")
  }

  return user
}
