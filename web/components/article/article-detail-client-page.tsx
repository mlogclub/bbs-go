"use client"

import * as React from "react"
import Link from "@/components/common/link"

import { ArticleComments } from "@/components/article/article-comments"
import { ArticleManageMenu } from "@/components/article/article-manage-menu"
import { HtmlImagePreview } from "@/components/common/image-preview"
import { MainShell } from "@/components/layout/main-shell"
import { PageError, PageLoading } from "@/components/common/page-state"
import { TopicToc } from "@/components/topic/topic-toc"
import { UserInfo } from "@/components/user/user-info"
import { useAppConfig, useCurrentUser } from "@/components/app/app-provider"
import { apiFetch } from "@/lib/api/client"
import type { Article, Comment, PageData } from "@/lib/api/types"
import { prettyDate } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"
import { useRouteData, useRouteSegment } from "@/lib/spa-route"
import { useDocumentTitle } from "@/lib/use-document-title"

const emptyComments: PageData<Comment> = {
  results: [],
  cursor: "0",
  hasMore: false,
}

type ArticleDetailData = {
  article: Article
  comments: PageData<Comment>
}

export function ArticleDetailClientPage({
  initialArticle,
}: {
  initialArticle?: Article
}) {
  const id = useRouteSegment(1)
  const { t } = useI18n()
  const config = useAppConfig()
  const currentUser = useCurrentUser()
  const initialData = React.useMemo<ArticleDetailData | null>(
    () =>
      initialArticle
        ? {
            article: initialArticle,
            comments: emptyComments,
          }
        : null,
    [initialArticle]
  )
  const load = React.useCallback(async (): Promise<ArticleDetailData> => {
    const comments = apiFetch<PageData<Comment>>("/api/comment/comments", {
      params: { entityType: "article", entityId: id },
    }).catch(() => emptyComments)

    if (initialArticle) {
      return { article: initialArticle, comments: await comments }
    }

    const [article, nextComments] = await Promise.all([
      apiFetch<Article>(`/api/article/${id}`),
      comments,
    ])

    return { article, comments: nextComments }
  }, [id, initialArticle])
  const { data, loading, error } = useRouteData(
    `article:${id}`,
    load,
    initialData
  )
  useDocumentTitle(data?.article.title)

  if (loading && !data)
    return (
      <MainShell>
        <PageLoading />
      </MainShell>
    )
  if (error || !data)
    return (
      <MainShell>
        <PageError message={error} />
      </MainShell>
    )

  const { article, comments } = data

  return (
    <MainShell
      sideSize="360"
      aside={
        <>
          {article.user ? <UserInfo user={article.user} t={t} /> : null}
          <TopicToc items={article.toc} />
        </>
      }
      containerClassName="side-size-360"
      asideClassName="!h-auto self-stretch"
    >
      <div className="space-y-4">
        {article.status === 2 ? (
          <div className="rounded-md border bg-muted px-3 py-2 text-sm text-muted-foreground">
            {t("pages.article.detail.pending")}
          </div>
        ) : null}
        <article className="rounded-lg bg-background p-3">
          <header className="border-b py-2.5">
            <div className="flex">
              <h1 className="w-full overflow-hidden text-lg leading-[30px] font-normal break-all text-ellipsis text-foreground">
                {article.title}
              </h1>
              <div className="min-w-max">
                <ArticleManageMenu
                  article={article}
                  currentUser={currentUser}
                />
              </div>
            </div>
            <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
              {article.user ? (
                <Link href={`/user/${article.user.id}`}>
                  {article.user.nickname}
                </Link>
              ) : null}
              {article.createTime ? (
                <span>{prettyDate(article.createTime, t)}</span>
              ) : null}
            </div>
          </header>
          <HtmlImagePreview
            html={article.content || ""}
            className="bbs-content max-w-none py-4 break-words [&_h2]:scroll-mt-20 [&_h3]:scroll-mt-20 [&_h4]:scroll-mt-20 [&_img]:cursor-zoom-in"
          />
        </article>
        <ArticleComments
          entityId={article.id}
          commentCount={article.commentCount}
          initialComments={comments}
          currentUser={currentUser}
          config={config}
        />
      </div>
    </MainShell>
  )
}
