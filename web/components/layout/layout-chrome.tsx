"use client"

import type * as React from "react"
import { usePathname } from "@/lib/router/navigation"

import { SiteFooter } from "@/components/layout/site-footer"
import { SiteHeader } from "@/components/layout/site-header"

export function LayoutChrome({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()

  if (pathname === "/install" || pathname.startsWith("/dashboard")) {
    return <>{children}</>
  }

  return (
    <>
      <SiteHeader />
      {children}
      <SiteFooter />
    </>
  )
}
