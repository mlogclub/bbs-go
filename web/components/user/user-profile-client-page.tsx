"use client"

import * as React from "react"
import Link from "@/components/common/link"
import { FileText, Medal, MessageSquare, UserPlus, Users } from "lucide-react"

import { ArticleList } from "@/components/article/article-list"
import { useCurrentUser } from "@/components/app/app-provider"
import { EmptyState } from "@/components/common/empty-state"
import { LoadMore } from "@/components/common/load-more"
import { PageError, PageLoading } from "@/components/common/page-state"
import { TopicListItem } from "@/components/topic/topic-list-item"
import { UserCenterShell } from "@/components/user/user-center-shell"
import { UserFollowList } from "@/components/user/user-follow-list"
import { WidgetCard } from "@/components/common/widget-card"
import { apiFetch } from "@/lib/api/client"
import type {
  Article,
  Badge,
  PageData,
  Topic,
  UserSummary,
} from "@/lib/api/types"
import { formatDate } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"
import { useRouteData, useRouteSegment } from "@/lib/spa-route"
import { useDocumentTitle } from "@/lib/use-document-title"

type UserShellData = {
  user: UserSummary
  badges: Badge[]
  fans: UserSummary[]
  followed: UserSummary[]
}

type UserProfileData = UserShellData & {
  topics: PageData<Topic>
}

type UserArticlesData = UserShellData & {
  articles: PageData<Article>
}

type UserFollowData = UserShellData & {
  pageData: PageData<UserSummary>
}

const emptyPage: PageData<Topic> = { results: [], cursor: "0", hasMore: false }
const emptyArticlePage: PageData<Article> = {
  results: [],
  cursor: "0",
  hasMore: false,
}
const emptyUserPage: PageData<UserSummary> = {
  results: [],
  cursor: "0",
  hasMore: false,
}

async function loadUserShellData(userId: string): Promise<UserShellData> {
  const [user, badges, fans, followed] = await Promise.all([
    apiFetch<UserSummary>(`/api/user/${userId}`),
    apiFetch<Badge[]>("/api/badge/badges", { params: { userId } }).catch(
      () => []
    ),
    apiFetch<PageData<UserSummary>>("/api/fans/recent/fans", {
      params: { userId },
    })
      .then((data) => data.results || [])
      .catch(() => []),
    apiFetch<PageData<UserSummary>>("/api/fans/recent/follow", {
      params: { userId },
    })
      .then((data) => data.results || [])
      .catch(() => []),
  ])

  return { user, badges, fans, followed }
}

function userDisplayName(user: UserSummary | null | undefined) {
  return (
    user?.nickname || user?.username || (user?.id ? `#${user.id}` : undefined)
  )
}

export function UserProfileClientPage({
  initialUser = null,
}: {
  initialUser?: UserSummary | null
}) {
  const userId = useRouteSegment(1)
  const currentUser = useCurrentUser()
  const { t } = useI18n()
  const load = React.useCallback(async (): Promise<UserProfileData> => {
    const [shell, topics] = await Promise.all([
      loadUserShellData(userId),
      apiFetch<PageData<Topic>>("/api/topic/user_topics", {
        params: { userId },
      }).catch(() => emptyPage),
    ])

    return { ...shell, topics }
  }, [userId])
  const { data, loading, error } = useRouteData(`user:${userId}`, load)
  useDocumentTitle(userDisplayName(data?.user ?? initialUser))

  if (loading) return <PageLoading />
  if (error || !data) return <PageError message={error} />

  const { user, topics, badges, fans, followed } = data
  const loadMoreLabels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  return (
    <UserCenterShell
      user={user}
      currentUser={currentUser}
      badges={badges}
      fans={fans}
      followed={followed}
      t={t}
    >
      <WidgetCard>
        <nav className="mb-2 inline-flex h-9 items-center justify-center rounded-lg bg-muted p-[3px] text-muted-foreground">
          <Link
            href={`/user/${user.id}`}
            className="inline-flex h-full items-center justify-center gap-1.5 rounded-md bg-background px-3 py-1 text-sm font-medium text-foreground shadow-sm"
          >
            <MessageSquare className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{t("pages.user.topics")}</span>
          </Link>
          <Link
            href={`/user/${user.id}/articles`}
            className="inline-flex h-full items-center justify-center gap-1.5 rounded-md px-3 py-1 text-sm font-medium text-foreground/60 hover:text-foreground"
          >
            <FileText className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{t("pages.user.articles")}</span>
          </Link>
        </nav>
        <LoadMore<Topic>
          initialItems={topics.results || []}
          initialCursor={topics.cursor}
          initialHasMore={topics.hasMore}
          initialLoad
          resetKey={`user-topics:${userId}:${topics.cursor}:${topics.hasMore}`}
          labels={loadMoreLabels}
          loadPage={({ cursor }) =>
            apiFetch<PageData<Topic>>("/api/topic/user_topics", {
              params: { userId, cursor },
            })
          }
          renderItems={(items) => (
            <ul className="divide-y divide-border">
              {items.map((topic) => (
                <TopicListItem key={topic.id} topic={topic} t={t} />
              ))}
            </ul>
          )}
          renderEmpty={() => <EmptyState title={t("common.noData")} />}
        />
      </WidgetCard>
    </UserCenterShell>
  )
}

export function UserArticlesClientPage() {
  const userId = useRouteSegment(1)
  const currentUser = useCurrentUser()
  const { t } = useI18n()
  const load = React.useCallback(async (): Promise<UserArticlesData> => {
    const [shell, articles] = await Promise.all([
      loadUserShellData(userId),
      apiFetch<PageData<Article>>("/api/article/user_articles", {
        params: { userId },
      }).catch(() => emptyArticlePage),
    ])

    return { ...shell, articles }
  }, [userId])
  const { data, loading, error } = useRouteData(`user-articles:${userId}`, load)

  if (loading) return <PageLoading />
  if (error || !data) return <PageError message={error} />

  const { user, articles, badges, fans, followed } = data
  const loadMoreLabels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  return (
    <UserCenterShell
      user={user}
      currentUser={currentUser}
      badges={badges}
      profileBadges={[]}
      fans={fans}
      followed={followed}
      t={t}
    >
      <WidgetCard>
        <nav className="mb-2 inline-flex h-9 items-center justify-center rounded-lg bg-muted p-[3px] text-muted-foreground">
          <Link
            href={`/user/${user.id}`}
            className="inline-flex h-full items-center justify-center gap-1.5 rounded-md px-3 py-1 text-sm font-medium text-foreground/60 hover:text-foreground"
          >
            <MessageSquare className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{t("pages.user.topics")}</span>
          </Link>
          <Link
            href={`/user/${user.id}/articles`}
            className="inline-flex h-full items-center justify-center gap-1.5 rounded-md bg-background px-3 py-1 text-sm font-medium text-foreground shadow-sm"
          >
            <FileText className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{t("pages.user.articles")}</span>
          </Link>
        </nav>
        <LoadMore<Article>
          initialItems={articles.results || []}
          initialCursor={articles.cursor}
          initialHasMore={articles.hasMore}
          initialLoad
          resetKey={`user-articles:${userId}:${articles.cursor}:${articles.hasMore}`}
          labels={loadMoreLabels}
          loadPage={({ cursor }) =>
            apiFetch<PageData<Article>>("/api/article/user_articles", {
              params: { userId, cursor },
            })
          }
          renderItems={(items) => <ArticleList articles={items} t={t} />}
          renderEmpty={() => <EmptyState title={t("common.noData")} />}
        />
      </WidgetCard>
    </UserCenterShell>
  )
}

export function UserBadgesClientPage() {
  const userId = useRouteSegment(1)
  const currentUser = useCurrentUser()
  const { t } = useI18n()
  const load = React.useCallback(() => loadUserShellData(userId), [userId])
  const { data, loading, error } = useRouteData(`user-badges:${userId}`, load)

  if (loading) return <PageLoading />
  if (error || !data) return <PageError message={error} />

  const { user, badges, fans, followed } = data

  return (
    <UserCenterShell
      user={user}
      currentUser={currentUser}
      badges={badges.filter((badge) => badge.owned)}
      fans={fans}
      followed={followed}
      t={t}
    >
      <WidgetCard>
        <div className="mb-4">
          <p className="text-lg font-semibold text-slate-900 dark:text-slate-50">
            {t("pages.user.badgesTitle")}
          </p>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            {t("pages.user.badgesSubtitle")}
          </p>
        </div>
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
          {badges.map((badge) => (
            <div
              key={badge.id}
              className={
                badge.owned
                  ? "flex flex-col items-center gap-2 rounded-xl border border-amber-200/80 bg-amber-50/50 p-4 transition dark:border-amber-800/60 dark:bg-amber-900/20"
                  : "flex flex-col items-center gap-2 rounded-xl border border-slate-200 bg-slate-50/50 p-4 opacity-70 transition dark:border-slate-700 dark:bg-slate-900/40"
              }
            >
              <div className="relative">
                {badge.icon ? (
                  <img
                    src={badge.icon}
                    alt={badge.title || ""}
                    className={
                      badge.owned
                        ? "h-14 w-14 object-contain"
                        : "h-14 w-14 object-contain opacity-40 grayscale"
                    }
                  />
                ) : (
                  <div
                    className={
                      badge.owned
                        ? "flex h-14 w-14 items-center justify-center rounded-full bg-slate-200 dark:bg-slate-700"
                        : "flex h-14 w-14 items-center justify-center rounded-full bg-slate-200 opacity-40 dark:bg-slate-700"
                    }
                  >
                    <Medal className="h-8 w-8 text-slate-600 dark:text-slate-300" />
                  </div>
                )}
                {badge.worn ? (
                  <span className="absolute -top-1 -right-4 rounded-full bg-amber-500 px-1.5 py-0.5 text-[10px] font-bold text-white">
                    {t("component.userBadges.worn")}
                  </span>
                ) : null}
              </div>
              <span
                className={
                  badge.owned
                    ? "line-clamp-2 text-center text-sm font-medium text-slate-800 dark:text-slate-100"
                    : "line-clamp-2 text-center text-sm font-medium text-slate-500 dark:text-slate-400"
                }
              >
                {badge.title}
              </span>
              {badge.owned && badge.obtainTime ? (
                <span className="text-[11px] text-slate-500 dark:text-slate-400">
                  {t("component.userBadges.obtainedAt")}{" "}
                  {formatDate(badge.obtainTime, "yyyy-MM-dd")}
                </span>
              ) : !badge.owned ? (
                <span className="text-[11px] text-slate-400 dark:text-slate-500">
                  {t("component.userBadges.notObtained")}
                </span>
              ) : null}
            </div>
          ))}
        </div>
      </WidgetCard>
    </UserCenterShell>
  )
}

export function UserFansClientPage() {
  return <UserFollowClientPage kind="fans" />
}

export function UserFollowedClientPage() {
  return <UserFollowClientPage kind="followed" />
}

function UserFollowClientPage({ kind }: { kind: "fans" | "followed" }) {
  const userId = useRouteSegment(1)
  const currentUser = useCurrentUser()
  const { t } = useI18n()
  const load = React.useCallback(async (): Promise<UserFollowData> => {
    const [shell, pageData] = await Promise.all([
      loadUserShellData(userId),
      apiFetch<PageData<UserSummary>>(
        kind === "fans" ? "/api/fans/fans" : "/api/fans/followed",
        { params: { userId } }
      ).catch(() => emptyUserPage),
    ])

    return { ...shell, pageData }
  }, [kind, userId])
  const { data, loading, error } = useRouteData(`user-${kind}:${userId}`, load)

  if (loading) return <PageLoading />
  if (error || !data) return <PageError message={error} />

  const { user, pageData, badges, fans, followed } = data
  const labels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  return (
    <UserCenterShell
      user={user}
      currentUser={currentUser}
      badges={badges}
      fans={fans}
      followed={followed}
      t={t}
    >
      <WidgetCard>
        <nav className="mb-2 inline-flex h-9 items-center justify-center rounded-lg bg-muted p-[3px] text-muted-foreground">
          <Link
            href={`/user/${user.id}/fans`}
            className={
              kind === "fans"
                ? "inline-flex h-full items-center justify-center gap-1.5 rounded-md bg-background px-3 py-1 text-sm font-medium text-foreground shadow-sm"
                : "inline-flex h-full items-center justify-center gap-1.5 rounded-md px-3 py-1 text-sm font-medium text-foreground/60 hover:text-foreground"
            }
          >
            <Users className="h-3.5 w-3.5" />
            <span>{t("pages.user.fans")}</span>
          </Link>
          <Link
            href={`/user/${user.id}/followed`}
            className={
              kind === "followed"
                ? "inline-flex h-full items-center justify-center gap-1.5 rounded-md bg-background px-3 py-1 text-sm font-medium text-foreground shadow-sm"
                : "inline-flex h-full items-center justify-center gap-1.5 rounded-md px-3 py-1 text-sm font-medium text-foreground/60 hover:text-foreground"
            }
          >
            <UserPlus className="h-3.5 w-3.5" />
            <span>{t("pages.user.followed")}</span>
          </Link>
        </nav>
        <UserFollowList users={pageData.results || []} />
        <LoadMore<UserSummary>
          initialCursor={pageData.cursor}
          initialHasMore={pageData.hasMore}
          resetKey={`user-${kind}:${userId}:${pageData.cursor}:${pageData.hasMore}`}
          labels={labels}
          loadPage={({ cursor }) =>
            apiFetch<PageData<UserSummary>>(
              kind === "fans" ? "/api/fans/fans" : "/api/fans/followed",
              { params: { userId, cursor } }
            )
          }
          renderItems={(items) => <UserFollowList users={items} />}
          alwaysShowButton
        />
      </WidgetCard>
    </UserCenterShell>
  )
}
