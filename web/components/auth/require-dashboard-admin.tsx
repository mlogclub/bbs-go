"use client"

import * as React from "react"
import { useLocation } from "react-router"

import { useAuthChecked, useCurrentUser } from "@/components/app/app-provider"
import { ErrorPage } from "@/components/common/error-page"
import { userCanAccessDashboard } from "@/lib/auth/roles"
import { useRouter } from "@/lib/router/navigation"
import { buildSigninHref } from "@/lib/toast"

export function RequireDashboardAdmin({
  children,
}: {
  children: React.ReactNode
}) {
  const router = useRouter()
  const location = useLocation()
  const authChecked = useAuthChecked()
  const user = useCurrentUser()
  const redirectPath = `${location.pathname}${location.search}`

  React.useEffect(() => {
    if (authChecked && !user) {
      router.replace(buildSigninHref(redirectPath))
    }
  }, [authChecked, redirectPath, router, user])

  if (!authChecked || !user) {
    return null
  }

  if (!userCanAccessDashboard(user)) {
    return <ErrorPage statusCode={403} />
  }

  return children
}
