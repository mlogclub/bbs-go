"use client"

import * as React from "react"

import { MainShell } from "@/components/layout/main-shell"
import { PageError, PageLoading } from "@/components/common/page-state"
import { TopicActionProvider } from "@/components/topic/topic-action-context"
import { TopicAttachments } from "@/components/topic/topic-attachments"
import { TopicComments } from "@/components/topic/topic-comments"
import { TopicContent } from "@/components/topic/topic-content"
import { TopicDetailActions } from "@/components/topic/topic-detail-actions"
import { TopicHideContentLive } from "@/components/topic/topic-hide-content-live"
import { TopicMeta } from "@/components/topic/topic-meta"
import { TopicSideActionBar } from "@/components/topic/topic-side-action-bar"
import { TopicTags } from "@/components/topic/topic-tags"
import { TopicToc } from "@/components/topic/topic-toc"
import { TopicVoteCard } from "@/components/topic/topic-vote-card"
import { UserInfo } from "@/components/user/user-info"
import { useCurrentUser } from "@/components/app/app-provider"
import { apiFetch } from "@/lib/api/client"
import type {
  Comment,
  PageData,
  Topic,
  TopicHideContent,
  UserSummary,
} from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { useRouteData, useRouteSegment } from "@/lib/spa-route"
import { useDocumentTitle } from "@/lib/use-document-title"

const emptyComments: PageData<Comment> = {
  results: [],
  cursor: "0",
  hasMore: false,
}

type TopicDetailData = {
  topic: Topic
  comments: PageData<Comment>
  likeUsers: UserSummary[] | null
  hideContent: TopicHideContent | null
}

export function TopicDetailClientPage({
  initialTopic,
}: {
  initialTopic?: Topic
}) {
  const id = useRouteSegment(1)
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const initialData = React.useMemo<TopicDetailData | null>(
    () =>
      initialTopic
        ? {
            topic: initialTopic,
            comments: emptyComments,
            likeUsers: null,
            hideContent: null,
          }
        : null,
    [initialTopic]
  )
  const load = React.useCallback(async (): Promise<TopicDetailData> => {
    const comments = apiFetch<PageData<Comment>>("/api/comment/comments", {
      params: { entityType: "topic", entityId: id },
    }).catch(() => emptyComments)
    const likeUsers = apiFetch<UserSummary[] | null>(
      `/api/topic/recentlikes/${id}`
    ).catch(() => null)
    const hideContent = apiFetch<TopicHideContent>("/api/topic/hide_content", {
      params: { topicId: id },
    }).catch(() => null)

    if (initialTopic) {
      const [nextComments, nextLikeUsers, nextHideContent] = await Promise.all([
        comments,
        likeUsers,
        hideContent,
      ])
      return {
        topic: initialTopic,
        comments: nextComments,
        likeUsers: nextLikeUsers,
        hideContent: nextHideContent,
      }
    }

    const [topic, nextComments, nextLikeUsers, nextHideContent] =
      await Promise.all([
        apiFetch<Topic>(`/api/topic/${id}`),
        comments,
        likeUsers,
        hideContent,
      ])

    return {
      topic,
      comments: nextComments,
      likeUsers: nextLikeUsers,
      hideContent: nextHideContent,
    }
  }, [id, initialTopic])
  const { data, loading, error } = useRouteData(
    `topic:${id}`,
    load,
    initialData
  )
  useDocumentTitle(data?.topic.title)

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

  const { topic, comments, likeUsers, hideContent } = data
  const currentUserRoles = currentUser?.roles || []
  const canAcceptAnswer =
    topic.type === 2 &&
    Boolean(currentUser) &&
    (topic.user?.id === currentUser?.id ||
      currentUserRoles.includes("admin") ||
      currentUserRoles.includes("owner"))

  return (
    <MainShell
      sideSize="260"
      aside={
        <>
          <UserInfo user={topic.user} t={t} />
          <TopicToc topic={topic} />
        </>
      }
      containerClassName="side-size-260"
      asideClassName="!h-auto self-stretch"
    >
      <div className="main-content no-padding no-bg space-y-4">
        {topic.status === 2 ? (
          <div className="my-5 w-full rounded-md border border-amber-300 bg-amber-100 px-4 py-3 text-amber-800">
            {t("pages.topic.detail.pending")}
          </div>
        ) : null}
        <article className="mb-5 rounded-lg bg-background">
          <TopicActionProvider
            topicId={topic.id}
            liked={topic.liked}
            favorited={topic.favorited}
            likeCount={topic.likeCount}
            commentCount={topic.commentCount}
          >
            <TopicSideActionBar />
            <div className="mb-4 border-b px-4 py-3">
              {topic.title ? (
                <h1 className="text-[26px] leading-9 font-bold wrap-break-word text-foreground">
                  {topic.title}
                </h1>
              ) : null}
              <TopicMeta topic={topic} currentUser={currentUser} t={t} />
            </div>
            <TopicContent topic={topic} />
            <div className="mx-4 mb-4">
              <TopicHideContentLive
                topicId={topic.id}
                initialHideContent={hideContent}
              />
            </div>
            <div className="mb-4 px-4">
              <TopicVoteCard vote={topic.vote} />
            </div>
            <TopicAttachments attachments={topic.attachments} t={t} />
            <TopicTags topic={topic} likeUsers={likeUsers} />
            <div id="topic-actions">
              <TopicDetailActions
                topic={topic}
                labels={{
                  view: t("pages.topic.detail.view"),
                  like: t("pages.topic.detail.like"),
                  comment: t("pages.topic.detail.comment"),
                  favorite: t("pages.topic.detail.favorite"),
                }}
              />
            </div>
          </TopicActionProvider>
        </article>
        <TopicComments
          entityId={topic.id}
          commentCount={topic.commentCount}
          title={topic.type === 2 ? t("pages.topic.detail.answers") : ""}
          acceptedCommentId={topic.acceptedCommentId || 0}
          allowAcceptAnswer={canAcceptAnswer}
          initialData={comments}
        />
      </div>
    </MainShell>
  )
}
