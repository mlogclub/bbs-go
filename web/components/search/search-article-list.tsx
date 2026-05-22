"use client"

import Link from "@/components/common/link"
import { UserAvatar } from "@/components/common/avatar"
import { EmptyState } from "@/components/common/empty-state"
import type { SearchArticle } from "@/lib/api/types"
import { formatDate } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"

export function SearchArticleList({ results }: { results: SearchArticle[] }) {
  const { t } = useI18n()

  if (!results.length) {
    return <EmptyState title={t("common.noData")} />
  }

  return (
    <div className="divide-y divide-border/70">
      {results.map((item) => {
        const user = item.user
        const authorName = user
          ? user.nickname || user.username || user.id
          : undefined
        return (
          <article
            key={item.id}
            className="group px-4 py-5 transition-colors hover:bg-muted/25 sm:px-5"
          >
            <h3 className="text-base leading-6 font-semibold text-foreground sm:text-[17px]">
              <Link
                href={`/article/${item.id}`}
                target="_blank"
                rel="noopener noreferrer"
                className="transition-colors group-hover:text-primary"
              >
                <span
                  className="[&_em]:rounded-sm [&_em]:bg-yellow-200/80 [&_em]:px-0.5 [&_em]:not-italic dark:[&_em]:bg-yellow-500/30"
                  dangerouslySetInnerHTML={{ __html: item.title || "" }}
                />
              </Link>
            </h3>
            {item.summary ? (
              <Link
                href={`/article/${item.id}`}
                target="_blank"
                rel="noopener noreferrer"
                className="mt-2 block text-sm leading-6 text-muted-foreground transition-colors hover:text-foreground"
              >
                <span
                  className="[display:-webkit-box] overflow-hidden [-webkit-box-orient:vertical] [-webkit-line-clamp:2] [&_em]:rounded-sm [&_em]:bg-yellow-200/70 [&_em]:px-0.5 [&_em]:not-italic dark:[&_em]:bg-yellow-500/25"
                  dangerouslySetInnerHTML={{ __html: item.summary || "" }}
                />
              </Link>
            ) : null}
            <div className="mt-3 flex flex-wrap items-center gap-x-2.5 gap-y-1.5 text-xs text-muted-foreground/90">
              {user ? (
                <>
                  <UserAvatar
                    user={user}
                    size={20}
                    className="[&_a]:opacity-70 hover:[&_a]:opacity-100"
                  />
                  <Link
                    href={`/user/${user.id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-muted-foreground hover:text-foreground"
                  >
                    {authorName}
                  </Link>
                </>
              ) : null}
              {item.createTime ? <span>{formatDate(item.createTime)}</span> : null}
              {item.tags?.map((tag) => (
                <Link
                  key={tag.id}
                  href={`/articles/tag/${tag.id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="rounded-full border border-border/80 bg-background px-2 py-0.5 text-[11px] text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                >
                  {tag.name}
                </Link>
              ))}
            </div>
          </article>
        )
      })}
    </div>
  )
}
