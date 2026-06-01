import type { PermissionCode } from "@/lib/auth/permissions.generated"

export type DashboardNavItem = {
  title: string
  url: string
  icon?: React.ComponentType<{ className?: string }>
  permission?: PermissionCode
  items?: Array<{
    title: string
    url: string
    permission?: PermissionCode
  }>
}

/**
 * Central navigation tree for the dashboard.
 * Both AppSidebar (UI) and dashboard layout (breadcrumbs) read from this.
 * Add a new page here once; breadcrumbs and sidebar stay in sync.
 */
export const DASHBOARD_NAV: DashboardNavItem[] = [
  {
    title: "dashboard.nav.workspace",
    url: "/dashboard",
    permission: "DASHBOARD_VIEW",
  },
  {
    title: "dashboard.nav.content",
    url: "/dashboard/content",
    permission: "DASHBOARD_TOPIC_VIEW",
    items: [
      {
        title: "dashboard.nav.topics",
        url: "/dashboard/topics",
        permission: "DASHBOARD_TOPIC_VIEW",
      },
      {
        title: "dashboard.nav.articles",
        url: "/dashboard/articles",
        permission: "DASHBOARD_ARTICLE_VIEW",
      },
      {
        title: "dashboard.nav.categories",
        url: "/dashboard/categories",
        permission: "DASHBOARD_CATEGORY_VIEW",
      },
      {
        title: "dashboard.nav.links",
        url: "/dashboard/links",
        permission: "DASHBOARD_LINK_VIEW",
      },
      {
        title: "dashboard.nav.forbiddenWords",
        url: "/dashboard/forbidden-words",
        permission: "DASHBOARD_FORBIDDEN_WORD_VIEW",
      },
    ],
  },
  {
    title: "dashboard.nav.community",
    url: "/dashboard/users",
    permission: "DASHBOARD_USER_VIEW",
    items: [
      {
        title: "dashboard.nav.userList",
        url: "/dashboard/users",
        permission: "DASHBOARD_USER_VIEW",
      },
      {
        title: "dashboard.nav.userBadges",
        url: "/dashboard/user-badges",
        permission: "DASHBOARD_USER_BADGE_VIEW",
      },
      {
        title: "dashboard.nav.userExpLogs",
        url: "/dashboard/user-exp-logs",
        permission: "DASHBOARD_USER_EXP_LOG_VIEW",
      },
      {
        title: "dashboard.nav.userTaskLogs",
        url: "/dashboard/user-task-logs",
        permission: "DASHBOARD_USER_TASK_LOG_VIEW",
      },
      {
        title: "dashboard.nav.userReports",
        url: "/dashboard/user-reports",
        permission: "DASHBOARD_USER_REPORT_VIEW",
      },
    ],
  },
  {
    title: "dashboard.nav.growth",
    url: "/dashboard/badges",
    permission: "DASHBOARD_BADGE_VIEW",
    items: [
      {
        title: "dashboard.nav.badges",
        url: "/dashboard/badges",
        permission: "DASHBOARD_BADGE_VIEW",
      },
      {
        title: "dashboard.nav.levels",
        url: "/dashboard/levels",
        permission: "DASHBOARD_LEVEL_VIEW",
      },
      {
        title: "dashboard.nav.tasks",
        url: "/dashboard/tasks",
        permission: "DASHBOARD_TASK_VIEW",
      },
    ],
  },
  {
    title: "dashboard.nav.system",
    url: "/dashboard/settings",
    permission: "DASHBOARD_SETTING_VIEW",
    items: [
      {
        title: "dashboard.nav.siteSettings",
        url: "/dashboard/settings",
        permission: "DASHBOARD_SETTING_VIEW",
      },
      {
        title: "dashboard.nav.roles",
        url: "/dashboard/roles",
        permission: "DASHBOARD_ROLE_VIEW",
      },
      {
        title: "dashboard.nav.emailLogs",
        url: "/dashboard/email-logs",
        permission: "DASHBOARD_EMAIL_LOG_VIEW",
      },
    ],
  },
]

/**
 * Resolve breadcrumb trail from a pathname against the nav tree.
 * Returns [] when no match is found (falls back to default breadcrumb).
 */
export function resolveBreadcrumb(
  pathname: string,
  t: (key: string) => string
): Array<{ title: string; url?: string }> {
  const page = pathname.replace(/^\/dashboard\/?/, "").split("/")[0] || ""

  if (!page) {
    return [{ title: t("dashboard.nav.workspace"), url: "/dashboard" }]
  }

  if (page === "content") {
    return [{ title: t("dashboard.nav.content"), url: "/dashboard/content" }]
  }

  for (const group of DASHBOARD_NAV) {
    if (!group.items) continue
    for (const item of group.items) {
      const itemPage = item.url.replace(/^\/dashboard\//, "")
      if (itemPage === page) {
        return [
          { title: t(group.title), url: group.url },
          { title: t(item.title), url: item.url },
        ]
      }
    }
  }

  return [{ title: t("dashboard.nav.workspace"), url: "/dashboard" }]
}
