"use client"

import * as React from "react"
import {
  AlertCircleIcon,
  ArrowUpRightIcon,
  CheckCircle2Icon,
  ClockIcon,
  FileTextIcon,
  GaugeIcon,
  MailWarningIcon,
  MessageSquareIcon,
  SettingsIcon,
  ShieldCheckIcon,
  TagsIcon,
  UsersIcon,
} from "lucide-react"

import { useCurrentUser } from "@/components/app/app-provider"
import { adminGet, type AdminRecord } from "@/lib/api/admin"
import {
  PERMISSIONS,
  type PermissionCode,
} from "@/lib/auth/permissions.generated"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"

type OverviewMetricKey =
  | "totalUsers"
  | "totalTopics"
  | "totalArticles"
  | "todayUsers"
  | "todayTopics"

type PendingKey =
  | "pendingTopics"
  | "pendingArticles"
  | "pendingReports"
  | "failedEmails"

type RecentItem = {
  id?: number | string
  title?: string
  content?: string
  nickname?: string
  createTime?: number
}

type OverviewData = {
  metrics?: Partial<Record<OverviewMetricKey, number>>
  pending?: Partial<Record<PendingKey, number>>
  recent?: {
    topics?: RecentItem[]
    users?: RecentItem[]
  }
}

const numberFormatter = new Intl.NumberFormat()
const dateFormatter = new Intl.DateTimeFormat(undefined, {
  month: "short",
  day: "numeric",
  hour: "2-digit",
  minute: "2-digit",
})

function toNumber(value: unknown) {
  return typeof value === "number" && Number.isFinite(value) ? value : 0
}

function formatNumber(value: unknown) {
  return numberFormatter.format(toNumber(value))
}

function formatDate(value: unknown) {
  if (typeof value !== "number" || !Number.isFinite(value)) return null
  return dateFormatter.format(new Date(value))
}

function toRecord(value: unknown): AdminRecord {
  return value && typeof value === "object" ? (value as AdminRecord) : {}
}

function toRecentItems(value: unknown): RecentItem[] {
  if (!Array.isArray(value)) return []
  return value
    .filter((item): item is AdminRecord => !!item && typeof item === "object")
    .map((item) => ({
      id:
        typeof item.id === "number" || typeof item.id === "string"
          ? item.id
          : undefined,
      title: typeof item.title === "string" ? item.title : undefined,
      content: typeof item.content === "string" ? item.content : undefined,
      nickname: typeof item.nickname === "string" ? item.nickname : undefined,
      createTime:
        typeof item.createTime === "number" ? item.createTime : undefined,
    }))
}

function normalizeOverview(data: AdminRecord | null): OverviewData | null {
  if (!data) return null
  const metrics = toRecord(data.metrics)
  const pending = toRecord(data.pending)
  const recent = toRecord(data.recent)

  return {
    metrics: {
      totalUsers: toNumber(metrics.totalUsers),
      totalTopics: toNumber(metrics.totalTopics),
      totalArticles: toNumber(metrics.totalArticles),
      todayUsers: toNumber(metrics.todayUsers),
      todayTopics: toNumber(metrics.todayTopics),
    },
    pending: {
      pendingTopics: toNumber(pending.pendingTopics),
      pendingArticles: toNumber(pending.pendingArticles),
      pendingReports: toNumber(pending.pendingReports),
      failedEmails: toNumber(pending.failedEmails),
    },
    recent: {
      topics: toRecentItems(recent.topics),
      users: toRecentItems(recent.users),
    },
  }
}

export function DashboardOverview() {
  const currentUser = useCurrentUser()
  const { t } = useI18n()
  const [overview, setOverview] = React.useState<OverviewData | null>(null)

  React.useEffect(() => {
    let mounted = true

    adminGet<AdminRecord>("/api/admin/common/overview")
      .then((overviewData) => {
        if (!mounted) return
        setOverview(normalizeOverview(overviewData))
      })
      .catch(() => {
        if (!mounted) return
        setOverview(normalizeOverview(null))
      })

    return () => {
      mounted = false
    }
  }, [])

  const adminName =
    currentUser?.nickname ||
    currentUser?.username ||
    t("dashboard.user.anonymous")
  const canUse = (permission?: PermissionCode) =>
    userHasPermission(currentUser, permission)

  const metricCards: Array<{
    key: OverviewMetricKey
    icon: React.ComponentType<React.SVGProps<SVGSVGElement>>
  }> = [
    { key: "totalUsers", icon: UsersIcon },
    { key: "totalTopics", icon: MessageSquareIcon },
    { key: "totalArticles", icon: FileTextIcon },
    { key: "todayUsers", icon: UsersIcon },
    { key: "todayTopics", icon: GaugeIcon },
  ]

  const pendingItems: Array<{
    key: PendingKey
    href: string
    icon: React.ComponentType<React.SVGProps<SVGSVGElement>>
    permission: PermissionCode
  }> = [
    {
      key: "pendingTopics",
      href: "/dashboard/topics",
      icon: MessageSquareIcon,
      permission: PERMISSIONS.DASHBOARD_TOPIC_VIEW,
    },
    {
      key: "pendingArticles",
      href: "/dashboard/articles",
      icon: FileTextIcon,
      permission: PERMISSIONS.DASHBOARD_ARTICLE_VIEW,
    },
    {
      key: "pendingReports",
      href: "/dashboard/user-reports",
      icon: AlertCircleIcon,
      permission: PERMISSIONS.DASHBOARD_USER_REPORT_VIEW,
    },
    {
      key: "failedEmails",
      href: "/dashboard/email-logs",
      icon: MailWarningIcon,
      permission: PERMISSIONS.DASHBOARD_EMAIL_LOG_VIEW,
    },
  ]

  const quickLinks: Array<{
    key: string
    href: string
    icon: React.ComponentType<React.SVGProps<SVGSVGElement>>
    permission: PermissionCode
  }> = [
    {
      key: "topics",
      href: "/dashboard/topics",
      icon: MessageSquareIcon,
      permission: PERMISSIONS.DASHBOARD_TOPIC_VIEW,
    },
    {
      key: "users",
      href: "/dashboard/users",
      icon: UsersIcon,
      permission: PERMISSIONS.DASHBOARD_USER_VIEW,
    },
    {
      key: "nodes",
      href: "/dashboard/nodes",
      icon: TagsIcon,
      permission: PERMISSIONS.DASHBOARD_NODE_VIEW,
    },
    {
      key: "settings",
      href: "/dashboard/settings",
      icon: SettingsIcon,
      permission: PERMISSIONS.DASHBOARD_SETTING_VIEW,
    },
    {
      key: "tasks",
      href: "/dashboard/tasks",
      icon: ShieldCheckIcon,
      permission: PERMISSIONS.DASHBOARD_TASK_VIEW,
    },
    {
      key: "levels",
      href: "/dashboard/levels",
      icon: GaugeIcon,
      permission: PERMISSIONS.DASHBOARD_LEVEL_VIEW,
    },
  ]

  const recentSections: Array<{
    key: "topics" | "users"
    href: string
  }> = [
    { key: "topics", href: "/dashboard/topics" },
    { key: "users", href: "/dashboard/users" },
  ]
  const visiblePendingItems = pendingItems
    .filter((item) => canUse(item.permission))
    .map((item) => ({
      ...item,
      count: toNumber(overview?.pending?.[item.key]),
    }))
    .sort((left, right) => right.count - left.count)
  const pendingItemsForDisplay = [
    ...visiblePendingItems.filter((item) => item.key !== "failedEmails"),
    ...visiblePendingItems.filter((item) => item.key === "failedEmails"),
  ]
  const priorityPendingItems = visiblePendingItems.filter(
    (item) => item.key !== "failedEmails"
  )
  const visibleQuickLinks = quickLinks.filter((item) => canUse(item.permission))
  const topPending = priorityPendingItems.find((item) => item.count > 0)
  const priorityPending = priorityPendingItems.reduce(
    (sum, item) => sum + item.count,
    0
  )
  const totalPending = visiblePendingItems.reduce(
    (sum, item) => sum + item.count,
    0
  )

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-4 md:gap-5 md:p-6">
      <section className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <div className="rounded-lg border bg-[var(--dashboard-panel)] p-5 text-card-foreground shadow-xs md:p-6">
          <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
            <div className="min-w-0 space-y-2">
              <p className="text-xs font-medium tracking-wide text-primary uppercase">
                {t("dashboard.overview.hero.eyebrow")}
              </p>
              <h1 className="text-2xl font-semibold tracking-tight md:text-3xl">
                {t("dashboard.overview.hero.title", { name: adminName })}
              </h1>
              <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
                {t("dashboard.overview.hero.description")}
              </p>
            </div>
            <a
              className="inline-flex h-9 shrink-0 items-center gap-2 rounded-md border bg-background px-3 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
              href="/"
              target="_blank"
              rel="noreferrer"
            >
              {t("dashboard.overview.hero.openSite")}
              <ArrowUpRightIcon className="size-4" />
            </a>
          </div>

          <div className="mt-6 grid gap-3 md:grid-cols-3">
            {metricCards.slice(0, 3).map(({ key, icon: Icon }) => (
              <div key={key} className="rounded-md border bg-background/45 p-3">
                <div className="flex items-center justify-between gap-3">
                  <p className="text-xs font-medium text-muted-foreground">
                    {t(`dashboard.overview.metrics.${key}`)}
                  </p>
                  <Icon className="size-4 text-muted-foreground" />
                </div>
                <p className="mt-2 text-2xl font-semibold tracking-tight">
                  {formatNumber(overview?.metrics?.[key])}
                </p>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-lg border bg-[var(--dashboard-panel)] p-5 text-card-foreground shadow-xs md:p-6">
          <div className="flex items-start justify-between gap-4">
            <div>
              <p className="text-sm font-medium">
                {t("dashboard.overview.priority.title")}
              </p>
              <p className="mt-1 text-sm text-muted-foreground">
                {t("dashboard.overview.priority.description")}
              </p>
            </div>
            <span className="rounded-md border bg-background px-2.5 py-1 text-sm font-semibold">
              {formatNumber(priorityPending)}
            </span>
          </div>

          {topPending ? (
            <a
              href={topPending.href}
              className="mt-5 flex items-center justify-between gap-3 rounded-lg border border-primary/25 bg-primary/8 p-4 transition-colors hover:bg-primary/12"
            >
              <span className="flex min-w-0 items-center gap-3">
                <span className="rounded-md bg-primary/10 p-2 text-primary">
                  {React.createElement(topPending.icon, {
                    className: "size-5",
                  })}
                </span>
                <span className="min-w-0">
                  <span className="block text-sm font-semibold">
                    {t(`dashboard.overview.pending.${topPending.key}`)}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {t("dashboard.overview.pending.goHandle")}
                  </span>
                </span>
              </span>
              <span className="text-2xl font-semibold">
                {formatNumber(topPending.count)}
              </span>
            </a>
          ) : (
            <div className="mt-5 flex items-center gap-3 rounded-lg border bg-background/45 p-4">
              <span className="rounded-md bg-primary/10 p-2 text-primary">
                <CheckCircle2Icon className="size-5" />
              </span>
              <div>
                <p className="text-sm font-semibold">
                  {t("dashboard.overview.priority.emptyTitle")}
                </p>
                <p className="text-xs text-muted-foreground">
                  {t("dashboard.overview.priority.emptyDescription")}
                </p>
              </div>
            </div>
          )}

          <div className="mt-4 grid grid-cols-2 gap-3">
            {metricCards.slice(3).map(({ key, icon: Icon }) => (
              <div key={key} className="rounded-md border bg-background/45 p-3">
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Icon className="size-3.5" />
                  {t(`dashboard.overview.metrics.${key}`)}
                </div>
                <p className="mt-2 text-xl font-semibold">
                  {formatNumber(overview?.metrics?.[key])}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <section className="rounded-lg border bg-[var(--dashboard-panel)] p-4 text-card-foreground shadow-xs">
          <div className="flex items-start justify-between gap-3">
            <div>
              <h2 className="text-base font-semibold">
                {t("dashboard.overview.pending.title")}
              </h2>
              <p className="mt-1 text-sm text-muted-foreground">
                {t("dashboard.overview.pending.description")}
              </p>
            </div>
            <span className="rounded-md bg-muted px-2 py-1 text-xs text-muted-foreground">
              {formatNumber(totalPending)}
            </span>
          </div>
          <div className="mt-4 grid gap-3 md:grid-cols-2">
            {pendingItemsForDisplay.map(({ key, href, icon: Icon, count }) => (
              <a
                key={key}
                className="group flex items-center justify-between gap-3 rounded-md border p-3 transition-colors hover:bg-accent hover:text-accent-foreground"
                href={href}
              >
                <span className="flex min-w-0 items-center gap-3">
                  <span className="rounded-md bg-muted p-2 text-muted-foreground group-hover:bg-background group-hover:text-foreground">
                    <Icon className="size-4" />
                  </span>
                  <span className="min-w-0">
                    <span className="block truncate text-sm font-medium">
                      {t(`dashboard.overview.pending.${key}`)}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {count > 0
                        ? t("dashboard.overview.pending.goHandle")
                        : t("dashboard.overview.priority.emptyTitle")}
                    </span>
                  </span>
                </span>
                <span className="text-lg font-semibold">
                  {formatNumber(count)}
                </span>
              </a>
            ))}
          </div>
        </section>

        <section className="rounded-lg border bg-[var(--dashboard-panel)] p-4 text-card-foreground shadow-xs">
          <h2 className="text-base font-semibold">
            {t("dashboard.overview.quick.title")}
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            {t("dashboard.overview.quick.description")}
          </p>
          <div className="mt-4 grid gap-2 sm:grid-cols-2 xl:grid-cols-1">
            {visibleQuickLinks.map(({ key, href, icon: Icon }) => (
              <a
                key={key}
                href={href}
                className="flex items-center justify-between gap-3 rounded-md border p-3 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
              >
                <span className="flex min-w-0 items-center gap-3">
                  <Icon className="size-4 shrink-0 text-muted-foreground" />
                  <span className="truncate">
                    {t(`dashboard.overview.quick.${key}`)}
                  </span>
                </span>
                <ArrowUpRightIcon className="size-3.5 shrink-0 text-muted-foreground" />
              </a>
            ))}
          </div>
        </section>
      </div>

      <section className="rounded-lg border bg-[var(--dashboard-panel)] p-4 text-card-foreground shadow-xs">
        <div className="flex items-center gap-2">
          <ClockIcon className="size-4 text-muted-foreground" />
          <h2 className="text-base font-semibold">
            {t("dashboard.overview.recent.title")}
          </h2>
        </div>
        <div className="mt-4 grid gap-4 lg:grid-cols-2">
          {recentSections.map(({ key, href }) => {
            const items = overview?.recent?.[key] ?? []
            return (
              <div key={key} className="rounded-md border p-3">
                <div className="flex items-center justify-between gap-2">
                  <h3 className="text-sm font-medium">
                    {t(`dashboard.overview.recent.${key}`)}
                  </h3>
                  <a
                    href={href}
                    className="text-xs text-primary underline-offset-4 hover:underline"
                  >
                    {t("dashboard.actions.view")}
                  </a>
                </div>
                <div className="mt-3 space-y-2">
                  {items.length > 0 ? (
                    items.map((item, index) => (
                      <div
                        key={`${item.id ?? key}-${index}`}
                        className="rounded-md bg-muted/40 px-3 py-2"
                      >
                        <p className="line-clamp-1 text-sm">
                          {item.title ||
                            item.content ||
                            item.nickname ||
                            t("common.noData")}
                        </p>
                        {item.nickname || item.createTime ? (
                          <p className="mt-1 flex flex-wrap gap-x-2 gap-y-1 text-xs text-muted-foreground">
                            {item.nickname ? (
                              <span>{item.nickname}</span>
                            ) : null}
                            {formatDate(item.createTime) ? (
                              <span>{formatDate(item.createTime)}</span>
                            ) : null}
                          </p>
                        ) : null}
                      </div>
                    ))
                  ) : (
                    <p className="rounded-md bg-muted/40 px-3 py-6 text-center text-sm text-muted-foreground">
                      {t("common.noData")}
                    </p>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      </section>
    </div>
  )
}
