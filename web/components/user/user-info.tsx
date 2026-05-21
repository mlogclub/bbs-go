import Link from "@/components/common/link"

import { UserAvatar } from "@/components/common/avatar"
import type { UserSummary } from "@/lib/api/types"
import type { TFunction } from "@/lib/i18n"

type UserInfoSummary = UserSummary & {
  level?: number
  topicCount?: number
  commentCount?: number
}

function displayName(user: UserSummary) {
  return user.nickname || user.username || `#${user.id}`
}

export function UserInfo({ user, t }: { user: UserInfoSummary; t: TFunction }) {
  return (
    <section className="mb-2.5 rounded-lg bg-background">
      <div className="px-2.5 py-2.5 text-center">
        <UserAvatar user={user} size={80} className="mx-auto" />
        <h2 className="mx-auto my-2.5 truncate text-[15px] font-bold">
          <Link href={`/user/${user.id}`} className="hover:underline">
            {displayName(user)}
          </Link>
        </h2>
        {user.description ? (
          <p className="mt-1 line-clamp-2 text-center text-[13px] break-all text-muted-foreground">
            {user.description}
          </p>
        ) : null}
      </div>
      <div className="border-t bg-foreground/[0.01] py-1.5">
        <ul className="flex text-center">
          <li className="w-full">
            <span className="text-[13px] font-normal text-muted-foreground">
              {t("component.userInfo.level")}
            </span>
            <br />
            <b>{user.level ?? 0}</b>
          </li>
          <li className="w-full">
            <span className="text-[13px] font-normal text-muted-foreground">
              {t("component.userInfo.score")}
            </span>
            <br />
            <b>{user.score ?? 0}</b>
          </li>
          <li className="w-full">
            <span className="text-[13px] font-normal text-muted-foreground">
              {t("component.userInfo.topicCount")}
            </span>
            <br />
            <b>{user.topicCount ?? 0}</b>
          </li>
          <li className="w-full">
            <span className="text-[13px] font-normal text-muted-foreground">
              {t("component.userInfo.commentCount")}
            </span>
            <br />
            <b>{user.commentCount ?? 0}</b>
          </li>
        </ul>
      </div>
    </section>
  )
}
