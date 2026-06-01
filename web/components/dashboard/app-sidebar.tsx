"use client"

import * as React from "react"
import {
  BadgeIcon,
  LayoutDashboardIcon,
  MessageSquareIcon,
  Settings2Icon,
  UsersIcon,
} from "lucide-react"
import type { LucideIcon } from "lucide-react"
import { DASHBOARD_NAV } from "@/lib/dashboard/navigation"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"
import { useCurrentUser } from "@/components/app/app-provider"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"

import { NavMain } from "@/components/dashboard/nav-main"
import { NavUser } from "@/components/dashboard/nav-user"
import { SidebarBrand } from "@/components/dashboard/sidebar-brand"

/* Icon lookup for dashboard nav items (icons cannot live in shared data) */
const DASHBOARD_NAV_ICONS: Record<string, LucideIcon> = {
  "/dashboard": LayoutDashboardIcon,
  "/dashboard/content": MessageSquareIcon,
  "/dashboard/users": UsersIcon,
  "/dashboard/badges": BadgeIcon,
  "/dashboard/settings": Settings2Icon,
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const currentUser = useCurrentUser()
  const { t } = useI18n()

  const data = {
    user: {
      name:
        currentUser?.nickname ||
        currentUser?.username ||
        t("dashboard.user.anonymous"),
      email: currentUser?.email || t("dashboard.user.noEmail"),
      avatar: currentUser?.smallAvatar || currentUser?.avatar || "",
    },
    brand: {
      name: t("dashboard.brand.name"),
      logoSrc: "/logo.png",
      description: t("dashboard.brand.plan"),
    },
    navMain: DASHBOARD_NAV.map((item) => ({
      ...item,
      title: t(item.title),
      icon: DASHBOARD_NAV_ICONS[item.url],
      items: item.items?.map((sub) => ({
        ...sub,
        title: t(sub.title),
      })),
    })),
  }

  const navMain = data.navMain
    .map((item) => {
      const children = item.items?.filter((child) =>
        userHasPermission(currentUser, child.permission)
      )
      const visible =
        userHasPermission(currentUser, item.permission) ||
        Boolean(children?.length)
      if (!visible) return null
      return {
        ...item,
        items: children,
      }
    })
    .filter((item): item is NonNullable<typeof item> => Boolean(item))

  return (
    <Sidebar
      collapsible="icon"
      className="border-r bg-[var(--dashboard-panel)]"
      {...props}
    >
      <SidebarHeader>
        <SidebarBrand {...data.brand} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
