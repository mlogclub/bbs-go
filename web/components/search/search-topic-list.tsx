"use client"

import Link from "@/components/common/link"

import { UserAvatar } from "@/components/common/avatar"
import { EmptyState } from "@/components/common/empty-state"
import type { Topic } from "@/lib/api/types"
import { formatDate } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"

export function SearchTopicList({ results }: { results: Topic[] }) {
  const { t } = useI18n()

  if (!results.length) {
    return <EmptyState title={t("common.noData")} />
  }

  return (
    <div className="divide-y divide-border/70">
      {results.map((item) => (
        <article
          key={item.id}
          className="group px-4 py-5 transition-colors hover:bg-muted/25 sm:px-5"
        >
          <div className="min-w-0">
            <h3 className="text-base leading-6 font-semibold text-foreground sm:text-[17px]">
              <Link
                href={`/topic/${item.id}`}
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
            <Link
              href={`/topic/${item.id}`}
              target="_blank"
              rel="noopener noreferrer"
              className="mt-2 block text-sm leading-6 text-muted-foreground transition-colors hover:text-foreground"
            >
              <span
                className="[display:-webkit-box] overflow-hidden [-webkit-box-orient:vertical] [-webkit-line-clamp:2] [&_em]:rounded-sm [&_em]:bg-yellow-200/70 [&_em]:px-0.5 [&_em]:not-italic dark:[&_em]:bg-yellow-500/25"
                dangerouslySetInnerHTML={{ __html: item.summary || "" }}
              />
            </Link>
            <div className="mt-3 flex flex-wrap items-center gap-x-2.5 gap-y-1.5 text-xs text-muted-foreground/90">
              {item.user ? (
                <UserAvatar
                  user={item.user}
                  size={20}
                  className="[&_a]:opacity-70 hover:[&_a]:opacity-100"
                />
              ) : null}
              {item.user ? (
                <Link
                  href={`/user/${item.user.id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  {item.user.nickname || item.user.username || item.user.id}
                </Link>
              ) : null}
              <span>{formatDate(item.createTime)}</span>
              {item.node ? (
                <Link
                  href={`/topics/node/${item.node.id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="rounded-full border border-border/80 bg-background px-2 py-0.5 text-[11px] text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                >
                  {item.node.name}
                </Link>
              ) : null}
            </div>
          </div>
        </article>
      ))}
    </div>
  )
}
