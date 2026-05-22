import Link from "@/components/common/link"
import { ChevronRight, MessageSquareQuote, Trophy } from "lucide-react"

import { UserAvatar } from "@/components/common/avatar"
import { EmptyState } from "@/components/common/empty-state"
import type { Favorite, ScoreLog, UserMessage } from "@/lib/api/types"
import { formatDateTime, prettyDate } from "@/lib/format"
import type { TFunction } from "@/lib/i18n"

export function FavoriteList({
  favorites,
  t,
}: {
  favorites: Favorite[]
  t: TFunction
}) {
  if (!favorites.length) {
    return <EmptyState title={t("common.noData")} />
  }

  return (
    <ul className="favorite-list">
      {favorites.map((item) => (
        <li key={item.id} className="favorite-item">
          {item.deleted ? (
            <div className="favorite-summary">
              {t("user.favorites.contentExpired")}
            </div>
          ) : (
            <div>
              <div className="favorite-title">
                <a href={item.url || "#"} target="_blank" rel="noreferrer">
                  {item.title}
                </a>
              </div>
              {item.content ? (
                <div className="favorite-summary">{item.content}</div>
              ) : null}
              <div className="favorite-meta">
                {item.user ? (
                  <span className="favorite-meta-item">
                    <Link href={`/user/${item.user.id}`}>
                      {item.user.nickname || item.user.username}
                    </Link>
                  </span>
                ) : null}
                {item.createTime ? (
                  <span className="favorite-meta-item">
                    <time>{prettyDate(item.createTime, t)}</time>
                  </span>
                ) : null}
              </div>
            </div>
          )}
        </li>
      ))}
    </ul>
  )
}

export function MessageList({
  messages,
  t,
}: {
  messages: UserMessage[]
  t: TFunction
}) {
  if (!messages.length) {
    return (
      <div className="rounded-xl border border-dashed border-border bg-muted px-4 py-9 text-center text-sm text-muted-foreground">
        {t("common.noData")}
      </div>
    )
  }

  return (
    <ul className="flex w-full list-none flex-col gap-3">
      {messages.map((message) => (
        <li
          key={message.id}
          className="box-border w-full rounded-lg border border-border bg-background p-3 sm:p-3.5"
        >
          <div className="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1">
            <UserAvatar user={message.from} size={36} />
            {message.from?.id ? (
              <Link
                href={`/user/${message.from.id}`}
                target="_blank"
                className="text-sm font-semibold text-primary underline-offset-2 hover:underline"
              >
                {message.from.nickname}
              </Link>
            ) : (
              <span className="text-sm font-semibold text-primary">
                {message.from?.nickname}
              </span>
            )}
            <span className="text-xs leading-none text-muted-foreground">
              ·
            </span>
            {message.createTime ? (
              <span className="text-xs whitespace-nowrap text-muted-foreground">
                {prettyDate(message.createTime, t)}
              </span>
            ) : null}
            {message.title ? (
              <span className="max-w-full rounded-full border border-border bg-muted px-2 py-0.5 text-[11px] leading-4 font-medium break-all text-muted-foreground">
                {message.title}
              </span>
            ) : null}
          </div>
          <div className="mt-2 space-y-2">
            {message.content ? (
              <div className="text-sm leading-6 break-words whitespace-pre-line text-foreground">
                {message.content}
              </div>
            ) : null}
            {message.quoteContent ? (
              <div className="relative rounded-md border border-border/70 bg-muted/40 px-3 py-2 pr-8">
                <MessageSquareQuote
                  size={12}
                  className="absolute top-2 right-2 text-muted-foreground/70"
                />
                <div className="text-xs leading-5 break-words text-muted-foreground">
                  {message.quoteContent}
                </div>
              </div>
            ) : null}
            {message.detailUrl && !isSelfMessagesUrl(message.detailUrl) ? (
              <a
                href={message.detailUrl}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-0.5 text-xs font-medium text-muted-foreground transition-colors hover:text-primary"
              >
                {t("user.messages.viewDetails")}
                <ChevronRight size={12} />
              </a>
            ) : null}
          </div>
        </li>
      ))}
    </ul>
  )
}

export function ScoreLogList({
  scoreLogs,
  t,
}: {
  scoreLogs: ScoreLog[]
  t: TFunction
}) {
  if (!scoreLogs.length) {
    return (
      <div className="rounded-xl bg-muted/25 px-4 py-14 text-center">
        <Trophy
          className="mx-auto mb-3 h-10 w-10 text-muted-foreground/40"
          strokeWidth="1.5"
        />
        <p className="text-sm text-muted-foreground/80">{t("common.noData")}</p>
      </div>
    )
  }

  return (
    <ul className="flex w-full list-none flex-col gap-2">
      {scoreLogs.map((scoreLog) => {
        const gain = scoreLog.type === 0
        return (
          <li
            key={scoreLog.id}
            className="flex min-w-0 flex-col gap-1.5 rounded-xl bg-muted p-4 transition-colors hover:bg-muted/50 sm:flex-row sm:items-center sm:gap-4"
          >
            <div className="flex flex-wrap items-center gap-2 sm:flex-nowrap sm:gap-3">
              <span
                className={
                  gain
                    ? "inline-flex items-center rounded-full bg-emerald-400/20 px-2.5 py-0.5 text-xs font-medium text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300"
                    : "inline-flex items-center rounded-full bg-rose-400/20 px-2.5 py-0.5 text-xs font-medium text-rose-700 dark:bg-rose-500/15 dark:text-rose-300"
                }
              >
                {gain
                  ? t("user.scores.gainPoints")
                  : t("user.scores.losePoints")}
              </span>
              <span
                className={
                  gain
                    ? "inline-flex items-center gap-1 font-semibold text-emerald-700 tabular-nums dark:text-emerald-300"
                    : "inline-flex items-center gap-1 font-semibold text-rose-700 tabular-nums dark:text-rose-300"
                }
              >
                <Trophy className="h-3.5 w-3.5 shrink-0 opacity-70" />
                {gain ? "+" : ""}
                {scoreLog.score}
              </span>
            </div>
            {scoreLog.description ? (
              <p
                className="min-w-0 flex-1 truncate text-sm text-muted-foreground/90 sm:truncate"
                title={scoreLog.description}
              >
                {scoreLog.description}
              </p>
            ) : null}
            {scoreLog.createTime ? (
              <time className="shrink-0 text-xs text-muted-foreground/80">
                {formatDateTime(scoreLog.createTime)}
              </time>
            ) : null}
          </li>
        )
      })}
    </ul>
  )
}

function isSelfMessagesUrl(url?: string) {
  return url?.endsWith("/user/messages") || url?.endsWith("/user/messages/")
}
