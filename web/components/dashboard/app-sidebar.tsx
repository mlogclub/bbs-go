"use client"

import * as React from "react"
import {
  BadgeIcon,
  LayoutDashboardIcon,
  MessageSquareIcon,
  Settings2Icon,
  UsersIcon,
} from "lucide-react"

import { NavMain } from "@/components/dashboard/nav-main"
import { NavUser } from "@/components/dashboard/nav-user"
import { SidebarBrand } from "@/components/dashboard/sidebar-brand"
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
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

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
    navMain: [
      {
        title: t("dashboard.nav.workspace"),
        url: "/dashboard",
        icon: LayoutDashboardIcon,
        permission: PERMISSIONS.DASHBOARD_VIEW,
      },
      {
        title: t("dashboard.nav.content"),
        url: "/dashboard/content",
        icon: MessageSquareIcon,
        permission: PERMISSIONS.DASHBOARD_TOPIC_VIEW,
        items: [
          {
            title: t("dashboard.nav.topics"),
            url: "/dashboard/topics",
            permission: PERMISSIONS.DASHBOARD_TOPIC_VIEW,
          },
          {
            title: t("dashboard.nav.articles"),
            url: "/dashboard/articles",
            permission: PERMISSIONS.DASHBOARD_ARTICLE_VIEW,
          },
          {
            title: t("dashboard.nav.categories"),
            url: "/dashboard/categories",
            permission: PERMISSIONS.DASHBOARD_CATEGORY_VIEW,
          },
          {
            title: t("dashboard.nav.links"),
            url: "/dashboard/links",
            permission: PERMISSIONS.DASHBOARD_LINK_VIEW,
          },
          {
            title: t("dashboard.nav.forbiddenWords"),
            url: "/dashboard/forbidden-words",
            permission: PERMISSIONS.DASHBOARD_FORBIDDEN_WORD_VIEW,
          },
        ],
      },
      {
        title: t("dashboard.nav.community"),
        url: "/dashboard/users",
        icon: UsersIcon,
        permission: PERMISSIONS.DASHBOARD_USER_VIEW,
        items: [
          {
            title: t("dashboard.nav.userList"),
            url: "/dashboard/users",
            permission: PERMISSIONS.DASHBOARD_USER_VIEW,
          },
          {
            title: t("dashboard.nav.userBadges"),
            url: "/dashboard/user-badges",
            permission: PERMISSIONS.DASHBOARD_USER_BADGE_VIEW,
          },
          {
            title: t("dashboard.nav.userExpLogs"),
            url: "/dashboard/user-exp-logs",
            permission: PERMISSIONS.DASHBOARD_USER_EXP_LOG_VIEW,
          },
          {
            title: t("dashboard.nav.userTaskLogs"),
            url: "/dashboard/user-task-logs",
            permission: PERMISSIONS.DASHBOARD_USER_TASK_LOG_VIEW,
          },
          {
            title: t("dashboard.nav.userReports"),
            url: "/dashboard/user-reports",
            permission: PERMISSIONS.DASHBOARD_USER_REPORT_VIEW,
          },
        ],
      },
      {
        title: t("dashboard.nav.growth"),
        url: "/dashboard/badges",
        icon: BadgeIcon,
        permission: PERMISSIONS.DASHBOARD_BADGE_VIEW,
        items: [
          {
            title: t("dashboard.nav.badges"),
            url: "/dashboard/badges",
            permission: PERMISSIONS.DASHBOARD_BADGE_VIEW,
          },
          {
            title: t("dashboard.nav.levels"),
            url: "/dashboard/levels",
            permission: PERMISSIONS.DASHBOARD_LEVEL_VIEW,
          },
          {
            title: t("dashboard.nav.tasks"),
            url: "/dashboard/tasks",
            permission: PERMISSIONS.DASHBOARD_TASK_VIEW,
          },
        ],
      },
      {
        title: t("dashboard.nav.system"),
        url: "/dashboard/settings",
        icon: Settings2Icon,
        permission: PERMISSIONS.DASHBOARD_SETTING_VIEW,
        items: [
          {
            title: t("dashboard.nav.siteSettings"),
            url: "/dashboard/settings",
            permission: PERMISSIONS.DASHBOARD_SETTING_VIEW,
          },
          {
            title: t("dashboard.nav.roles"),
            url: "/dashboard/roles",
            permission: PERMISSIONS.DASHBOARD_ROLE_VIEW,
          },
          {
            title: t("dashboard.nav.emailLogs"),
            url: "/dashboard/email-logs",
            permission: PERMISSIONS.DASHBOARD_EMAIL_LOG_VIEW,
          },
        ],
      },
    ],
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
