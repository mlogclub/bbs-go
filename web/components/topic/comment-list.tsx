import Link from "@/components/common/link"
import { LogIn, MessageCircle, ThumbsUp } from "lucide-react"

import { EmptyState } from "@/components/common/empty-state"
import { UserAvatar } from "@/components/common/avatar"
import type { Comment } from "@/lib/api/types"
import { prettyDate } from "@/lib/format"
import type { TFunction } from "@/lib/i18n"

function CommentItem({
  comment,
  acceptedCommentId,
  t,
}: {
  comment: Comment
  acceptedCommentId?: number
  t: TFunction
}) {
  const isAccepted = acceptedCommentId === comment.id

  return (
    <div
      className={`flex py-2.5 ${
        isAccepted
          ? "mb-2 rounded-lg border border-primary/20 bg-primary/[0.06] p-3"
          : "border-b border-border last:border-b-0"
      }`}
    >
      <div>
        <UserAvatar user={comment.user} size={30} />
      </div>
      <div className="ml-2.5 min-w-0 flex-1">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex min-w-0 items-center gap-2">
            <Link
              href={`/user/${comment.user.id}`}
              className="truncate text-sm text-foreground hover:text-primary"
            >
              {comment.user.nickname ||
                comment.user.username ||
                comment.user.id}
            </Link>
            {isAccepted ? (
              <span className="shrink-0 rounded-full bg-muted px-2 py-0.5 text-xs text-primary">
                {t("component.comment.list.acceptedAnswer")}
              </span>
            ) : null}
          </div>
          <div className="flex flex-wrap items-center gap-x-2.5">
            {comment.createTime ? (
              <time className="text-xs text-muted-foreground">
                {prettyDate(comment.createTime, t)}
              </time>
            ) : null}
            {comment.ipLocation ? (
              <span className="text-xs text-muted-foreground">
                {t("component.comment.list.ipLocation")}
                {comment.ipLocation}
              </span>
            ) : null}
          </div>
        </div>
        {comment.content ? (
          comment.contentType === "text" ? (
            <div className="mt-2.5 mb-0 whitespace-pre-wrap text-foreground">
              {comment.content}
            </div>
          ) : (
            <div
              className="bbs-content mt-2.5 mb-0 whitespace-pre-wrap text-foreground"
              dangerouslySetInnerHTML={{ __html: comment.content }}
            />
          )
        ) : null}
        {comment.imageList?.length ? (
          <div className="mt-2.5 flex flex-wrap gap-2">
            {comment.imageList.map((image, index) => (
              <img
                key={`${image.url || image.preview || index}`}
                src={image.url || image.preview}
                alt=""
                className="h-[72px] w-[72px] cursor-pointer object-cover transition-all duration-500 ease-out hover:scale-[1.04]"
              />
            ))}
          </div>
        ) : null}
        <div className="mt-2.5 flex flex-wrap items-center gap-2.5">
          <div
            className={`flex items-center gap-0.5 text-xs text-muted-foreground ${comment.liked ? "font-medium text-primary" : ""}`}
          >
            <ThumbsUp className="h-3 w-3" />
            <span>
              {comment.liked
                ? t("component.comment.list.liked")
                : t("component.comment.list.like")}
            </span>
            {comment.likeCount && comment.likeCount > 0 ? (
              <span>{comment.likeCount}</span>
            ) : null}
          </div>
          <div className="flex items-center gap-0.5 text-xs text-muted-foreground">
            <MessageCircle className="h-3 w-3" />
            <span>{t("component.comment.list.reply")}</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export function CommentList({
  t,
  comments,
  commentCount,
  title,
  acceptedCommentId,
  isLogin,
  embedded,
}: {
  t: TFunction
  comments?: Comment[]
  commentCount?: number
  title?: string
  acceptedCommentId?: number
  isLogin?: boolean
  embedded?: boolean
}) {
  const list = (
    <div className="text-sm">
      {comments?.length ? (
        comments.map((comment) => (
          <CommentItem
            key={comment.id}
            comment={comment}
            acceptedCommentId={acceptedCommentId}
            t={t}
          />
        ))
      ) : (
        <EmptyState title={t("common.noData")} className="min-h-36" />
      )}
    </div>
  )

  if (embedded) {
    return list
  }

  return (
    <section id="JComment" className="rounded-lg bg-background p-4">
      <div className="flex text-base font-medium text-foreground">
        <span>{title || t("component.comment.title")}</span>
        {commentCount && commentCount > 0 ? (
          <span>&nbsp;{commentCount}</span>
        ) : null}
      </div>
      {isLogin ? null : (
        <Link
          href="/user/signin"
          className="my-3 flex cursor-pointer items-center gap-2 rounded-lg border border-primary/40 bg-primary/[0.1] px-3 py-3 transition-colors hover:bg-primary/[0.15]"
        >
          <span className="inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <LogIn className="h-3.5 w-3.5" />
          </span>
          <div className="min-w-0 flex-1">
            <p className="text-sm font-medium text-foreground">
              {t("component.comment.loginLink")}
            </p>
          </div>
        </Link>
      )}
      {list}
    </section>
  )
}
