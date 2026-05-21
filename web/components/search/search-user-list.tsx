"use client"

import Link from "@/components/common/link"
import { UserAvatar } from "@/components/common/avatar"
import { EmptyState } from "@/components/common/empty-state"
import type { SearchUser } from "@/lib/api/types"
import { formatDate } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"

export function SearchUserList({ results }: { results: SearchUser[] }) {
  const { t } = useI18n()

  if (!results.length) {
    return <EmptyState title={t("common.noData")} />
  }

  return (
    <div className="divide-y divide-border/70">
      {results.map((item) => {
        const user = item.user
        if (!user) return null
        const name = item.nickname || user.nickname || user.username || user.id
        const description = item.description || user.description
        return (
          <article
            key={user.id}
            className="group px-4 py-5 transition-colors hover:bg-muted/25 sm:px-5"
          >
            <div className="flex gap-3.5">
              <UserAvatar user={user} size={48} className="shrink-0" />
              <div className="min-w-0 flex-1">
                <Link
                  href={`/user/${user.id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-base leading-6 font-semibold text-foreground transition-colors group-hover:text-primary sm:text-[17px]"
                >
                  <span
                    className="[&_em]:rounded-sm [&_em]:bg-yellow-200/80 [&_em]:px-0.5 [&_em]:not-italic dark:[&_em]:bg-yellow-500/30"
                    dangerouslySetInnerHTML={{ __html: name }}
                  />
                </Link>
                {description ? (
                  <p
                    className="mt-1.5 line-clamp-2 text-sm leading-6 text-muted-foreground [&_em]:rounded-sm [&_em]:bg-yellow-200/70 [&_em]:px-0.5 [&_em]:not-italic dark:[&_em]:bg-yellow-500/25"
                    dangerouslySetInnerHTML={{ __html: description }}
                  />
                ) : null}
                <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1.5 text-xs text-muted-foreground/90">
                  {user.levelTitle ? <span>{user.levelTitle}</span> : null}
                  {typeof user.topicCount === "number" ? (
                    <span>
                      {t("pages.search.userTopics", {
                        count: user.topicCount,
                      })}
                    </span>
                  ) : null}
                  {typeof user.fansCount === "number" ? (
                    <span>
                      {t("pages.search.userFans", { count: user.fansCount })}
                    </span>
                  ) : null}
                  {item.createTime ? <span>{formatDate(item.createTime)}</span> : null}
                </div>
              </div>
            </div>
          </article>
        )
      })}
    </div>
  )
}
