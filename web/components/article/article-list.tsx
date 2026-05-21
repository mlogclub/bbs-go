import Link from "@/components/common/link"
import { UserAvatar } from "@/components/common/avatar"

import { EmptyState } from "@/components/common/empty-state"
import type { Article } from "@/lib/api/types"
import { prettyDate } from "@/lib/format"
import type { TFunction } from "@/lib/i18n"
import { EyeIcon, HeartIcon, MessageCircleIcon } from "lucide-react"

export function ArticleList({
  articles,
  t,
}: {
  articles: Article[]
  t: TFunction
}) {
  if (!articles.length) {
    return <EmptyState title={t("common.noData")} />
  }

  return (
    <div className="overflow-hidden rounded-lg bg-background">
      {articles.map((article) => (
        <ArticleListItem key={article.id} article={article} t={t} />
      ))}
    </div>
  )
}

function ArticleListItem({ article, t }: { article: Article; t: TFunction }) {
  const authorName =
    article.user.nickname || article.user.username || `#${article.user.id}`
  const articleUrl = `/article/${article.id}`

  return (
    <article className="group border-b border-border/70 bg-background px-3 py-4 transition-colors last:border-b-0 hover:bg-muted/35 sm:px-4 sm:py-5">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2.5 text-sm text-muted-foreground">
            <UserAvatar user={article.user} size={28} />
            <Link
              href={`/user/${article.user.id}`}
              className="min-w-0 truncate font-medium text-foreground/80 hover:text-primary"
            >
              {authorName}
            </Link>
            {article.createTime ? (
              <>
                <span className="h-1 w-1 shrink-0 rounded-full bg-muted-foreground/40" />
                <time
                  dateTime={new Date(article.createTime).toISOString()}
                  className="shrink-0"
                >
                  <span className="sr-only">
                    {t("component.articleList.publishedAt")}{" "}
                  </span>
                  {prettyDate(article.createTime, t)}
                </time>
              </>
            ) : null}
          </div>

          <Link
            href={articleUrl}
            target="_blank"
            className="mt-3 block text-[19px] leading-7 font-semibold text-balance text-foreground transition-colors group-hover:text-primary"
          >
            {article.title}
          </Link>

          {article.summary ? (
            <p className="mt-2 line-clamp-2 text-[15px] leading-7 break-words text-muted-foreground">
              {article.summary}
            </p>
          ) : null}

          <ArticleListMeta article={article} className="mt-4 hidden sm:flex" />
        </div>

        {article.cover?.url ? (
          <Link
            href={articleUrl}
            target="_blank"
            className="block overflow-hidden rounded-md bg-muted sm:w-40 md:w-44"
            aria-label={article.title}
          >
            <img
              src={article.cover.preview || article.cover.url}
              alt={article.title}
              className="aspect-[16/10] w-full object-cover transition-transform duration-300 group-hover:scale-[1.03] sm:aspect-[4/3]"
            />
          </Link>
        ) : null}

        <ArticleListMeta article={article} className="sm:hidden" />
      </div>
    </article>
  )
}

function ArticleListMeta({
  article,
  className = "",
}: {
  article: Article
  className?: string
}) {
  return (
    <div
      className={`flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-muted-foreground ${className}`}
    >
      <ArticleMetric icon={EyeIcon} value={article.viewCount} />
      <ArticleMetric
        icon={MessageCircleIcon}
        value={article.commentCount ?? 0}
        alwaysShow
      />
      <ArticleMetric icon={HeartIcon} value={article.likeCount} />
      {article.tags?.length ? (
        <div className="flex min-w-0 flex-wrap items-center gap-2">
          {article.tags.map((tag) => (
            <Link
              key={tag.id}
              href={`/articles/tag/${tag.id}`}
              className="rounded-full bg-muted px-2.5 py-1 text-xs leading-none text-muted-foreground transition-colors hover:bg-primary/10 hover:text-primary"
            >
              {tag.name}
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  )
}

function ArticleMetric({
  icon: Icon,
  value,
  alwaysShow = false,
}: {
  icon: typeof EyeIcon
  value?: number
  alwaysShow?: boolean
}) {
  if (!alwaysShow && !value) {
    return null
  }

  const displayValue = value ?? 0

  return (
    <span className="inline-flex items-center gap-1.5">
      <Icon className="size-3.5" />
      <span>{formatCount(displayValue)}</span>
    </span>
  )
}

function formatCount(value: number) {
  if (value >= 10000) {
    return `${Number((value / 10000).toFixed(1))}w`
  }
  if (value >= 1000) {
    return `${Number((value / 1000).toFixed(1))}k`
  }
  return String(value)
}
