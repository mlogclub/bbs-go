import Link from "@/components/common/link"
import { ChevronRight, Medal } from "lucide-react"

import { EmptyState } from "@/components/common/empty-state"
import { WidgetCard } from "@/components/common/widget-card"
import { UserCenterOperations } from "@/components/user/user-center-operations"
import { UserFollowList } from "@/components/user/user-follow-list"
import type { Badge, UserSummary } from "@/lib/api/types"
import type { TFunction } from "@/lib/i18n"

export function UserCountsCard({
  user,
  t,
}: {
  user: UserSummary
  t: TFunction
}) {
  return (
    <WidgetCard title={t("component.myCounts.title")}>
      <ul className="extra-info">
        <li>
          <span>{t("component.myCounts.level")}</span>
          <br />
          <b>{`Lv.${user.level ?? 0}`}</b>
        </li>
        <li>
          <span>{t("component.myCounts.score")}</span>
          <br />
          <Link href="/user/scores">
            <b>{user.score ?? 0}</b>
          </Link>
        </li>
        <li>
          <span>{t("component.myCounts.topicCount")}</span>
          <br />
          <b>{user.topicCount ?? 0}</b>
        </li>
        <li>
          <span>{t("component.myCounts.commentCount")}</span>
          <br />
          <b>{user.commentCount ?? 0}</b>
        </li>
      </ul>
    </WidgetCard>
  )
}

export function UserBadgesWidget({
  user,
  badges,
  t,
}: {
  user: UserSummary
  badges: Badge[]
  t: TFunction
}) {
  const owned = badges.filter((badge) => badge.owned !== false)
  const display = owned.slice(0, 6)
  const badgesLink = `/user/${user.id}/badges`

  return (
    <WidgetCard
      title={
        <>
          <span>{t("component.userBadges.title")}</span>
          <span>&nbsp;</span>
          <span>{owned.length}</span>
        </>
      }
      actions={
        <Link href={badgesLink} className="inline-flex items-center gap-1">
          {t("component.userBadges.viewAll")}
          <ChevronRight className="h-4 w-4" />
        </Link>
      }
    >
      {!owned.length ? (
        <div className="text-sm text-slate-500 dark:text-slate-400">
          {t("component.userBadges.noBadges")}
        </div>
      ) : (
        <div className="space-y-3">
          <div className="flex flex-wrap gap-2">
            {display.map((badge) => (
              <Link
                key={badge.id}
                href={badgesLink}
                className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-amber-200/80 bg-amber-50/50 dark:border-amber-800/60 dark:bg-amber-900/20"
                title={badge.title}
              >
                {badge.icon ? (
                  <img
                    src={badge.icon}
                    alt={badge.title || ""}
                    className="h-8 w-8 object-contain"
                  />
                ) : (
                  <Medal className="h-6 w-6" />
                )}
              </Link>
            ))}
          </div>
        </div>
      )}
    </WidgetCard>
  )
}

export function MyProfileCard({
  user,
  currentUser,
  t,
}: {
  user: UserSummary
  currentUser?: UserSummary | null
  t: TFunction
}) {
  const canEdit = currentUser?.id === user.id
  return (
    <WidgetCard
      title={t("component.myProfile.title")}
      actions={
        canEdit ? (
          <Link href="/user/profile" className="inline-flex items-center gap-1">
            {t("component.myProfile.editProfile")}
            <ChevronRight className="h-4 w-4" />
          </Link>
        ) : null
      }
    >
      <div className="stable">
        <div className="str">
          <div className="slabel">{t("component.myProfile.nickname")}</div>
          <div className="svalue">{user.nickname}</div>
        </div>
        <div className="str">
          <div className="slabel">{t("component.myProfile.description")}</div>
          <div className="svalue">{user.description}</div>
        </div>
        {user.homePage ? (
          <div className="str">
            <div className="slabel">{t("component.myProfile.homePage")}</div>
            <div className="svalue">
              <a href={user.homePage} target="_blank" rel="nofollow noreferrer">
                {user.homePage}
              </a>
            </div>
          </div>
        ) : null}
      </div>
    </WidgetCard>
  )
}

export function FollowWidget({
  title,
  count,
  moreHref,
  users,
  t,
}: {
  title: string
  count?: number
  moreHref: string
  users: UserSummary[]
  t: TFunction
}) {
  return (
    <WidgetCard
      title={
        <>
          <span>{title}</span>
          <span>&nbsp;</span>
          <span>{count ?? 0}</span>
        </>
      }
      actions={
        <Link href={moreHref} className="inline-flex items-center gap-1">
          {t("component.fansWidget.more")}
          <ChevronRight className="h-4 w-4" />
        </Link>
      }
    >
      {users.length ? (
        <UserFollowList users={users} />
      ) : (
        <EmptyState title={t("common.noData")} />
      )}
    </WidgetCard>
  )
}

export function UserCenterSidebar({
  user,
  currentUser,
  badges,
  fans,
  followed,
  t,
}: {
  user: UserSummary
  currentUser?: UserSummary | null
  badges: Badge[]
  fans: UserSummary[]
  followed: UserSummary[]
  t: TFunction
}) {
  return (
    <div className="left-container space-y-4">
      <UserCountsCard user={user} t={t} />
      <UserBadgesWidget user={user} badges={badges} t={t} />
      <MyProfileCard user={user} currentUser={currentUser} t={t} />
      <FollowWidget
        title={t("component.fansWidget.title")}
        count={user.fansCount}
        moreHref={`/user/${user.id}/fans`}
        users={fans}
        t={t}
      />
      <FollowWidget
        title={t("component.followWidget.title")}
        count={user.followCount}
        moreHref={`/user/${user.id}/followed`}
        users={followed}
        t={t}
      />
      <UserCenterOperations user={user} currentUser={currentUser} />
    </div>
  )
}
