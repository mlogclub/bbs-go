"use client"

import * as React from "react"
import { Bell, Heart, Trophy } from "lucide-react"

import { useRequiredUser } from "@/components/auth/require-user"
import { useSetUnreadMessageCount } from "@/components/app/app-provider"
import { LoadMore } from "@/components/common/load-more"
import { WidgetCard } from "@/components/common/widget-card"
import { UserCenterShell } from "@/components/user/user-center-shell"
import {
  FavoriteList,
  MessageList,
  ScoreLogList,
} from "@/components/user/user-lists"
import { apiFetch } from "@/lib/api/client"
import type {
  Badge,
  Favorite,
  PageData,
  ScoreLog,
  UserMessage,
  UserSummary,
} from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"

export type PrivateUserCenterKind = "favorites" | "messages" | "scores"
type PrivateUserCenterItem = Favorite | UserMessage | ScoreLog

function listPath(kind: PrivateUserCenterKind) {
  if (kind === "favorites") return "/api/user/favorites"
  if (kind === "messages") return "/api/user/messages"
  return "/api/user/score_logs"
}

function titleKey(kind: PrivateUserCenterKind) {
  if (kind === "favorites") return "user.favorites.title"
  if (kind === "messages") return "user.messages.title"
  return "user.scores.title"
}

function TitleIcon({ kind }: { kind: PrivateUserCenterKind }) {
  if (kind === "favorites") {
    return <Heart className="inline-block h-4 w-4" />
  }
  if (kind === "messages") {
    return <Bell size={18} />
  }
  return <Trophy className="h-4 w-4 shrink-0 text-emerald-500/90" />
}

function renderList(
  kind: PrivateUserCenterKind,
  items: PrivateUserCenterItem[],
  t: ReturnType<typeof useI18n>["t"]
) {
  if (kind === "favorites") {
    return <FavoriteList favorites={items as Favorite[]} t={t} />
  }
  if (kind === "messages") {
    return <MessageList messages={items as UserMessage[]} t={t} />
  }
  return <ScoreLogList scoreLogs={items as ScoreLog[]} t={t} />
}

export function PrivateUserCenterPage({
  kind,
  initialData,
  initialBadges = [],
  initialFans = [],
  initialFollowed = [],
  serverLoaded = false,
}: {
  kind: PrivateUserCenterKind
  initialData: PageData<PrivateUserCenterItem>
  initialBadges?: Badge[]
  initialFans?: UserSummary[]
  initialFollowed?: UserSummary[]
  serverLoaded?: boolean
}) {
  const { t } = useI18n()
  const user = useRequiredUser()
  const setUnreadMessageCount = useSetUnreadMessageCount()
  const [data, setData] = React.useState(initialData)
  const [badges, setBadges] = React.useState(initialBadges)
  const [fans, setFans] = React.useState(initialFans)
  const [followed, setFollowed] = React.useState(initialFollowed)
  const labels = {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }

  React.useEffect(() => {
    if (serverLoaded) return

    let mounted = true
    void Promise.all([
      apiFetch<PageData<PrivateUserCenterItem>>(listPath(kind)),
      apiFetch<Badge[]>("/api/badge/badges", {
        params: { userId: user.id },
      }).catch(() => []),
      apiFetch<PageData<UserSummary>>("/api/fans/recent/fans", {
        params: { userId: user.id },
      }).catch(() => ({ results: [], cursor: "", hasMore: false })),
      apiFetch<PageData<UserSummary>>("/api/fans/recent/follow", {
        params: { userId: user.id },
      }).catch(() => ({ results: [], cursor: "", hasMore: false })),
    ]).then(([nextData, nextBadges, nextFans, nextFollowed]) => {
      if (!mounted) return

      setData(nextData)
      if (kind === "messages") {
        setUnreadMessageCount(0)
      }
      setBadges(nextBadges)
      setFans(nextFans.results || [])
      setFollowed(nextFollowed.results || [])
    })

    return () => {
      mounted = false
    }
  }, [kind, serverLoaded, setUnreadMessageCount, user.id])

  return (
    <UserCenterShell
      user={user}
      currentUser={user}
      badges={badges}
      fans={fans}
      followed={followed}
      t={t}
    >
      <WidgetCard
        title={
          <div className="flex items-center gap-2">
            <TitleIcon kind={kind} />
            <span>{t(titleKey(kind))}</span>
          </div>
        }
      >
        {renderList(kind, data.results || [], t)}
        <LoadMore<PrivateUserCenterItem>
          key={`${kind}:${data.cursor || ""}:${data.hasMore}`}
          initialCursor={data.cursor}
          initialHasMore={data.hasMore}
          resetKey={`${kind}:${data.cursor || ""}:${data.hasMore}`}
          labels={labels}
          loadPage={({ cursor }) =>
            apiFetch<PageData<PrivateUserCenterItem>>(listPath(kind), {
              params: { cursor },
            })
          }
          renderItems={(items) => renderList(kind, items, t)}
          alwaysShowButton
        />
      </WidgetCard>
    </UserCenterShell>
  )
}
