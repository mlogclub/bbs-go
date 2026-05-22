"use client"

import * as React from "react"

import { CommentSection } from "@/components/comment"
import { topicCommentCreatedEvent } from "@/components/topic/topic-hide-content-live"
import type { Comment, EntityId, PageData } from "@/lib/api/types"

export function TopicComments({
  entityId,
  commentCount,
  title,
  acceptedCommentId,
  allowAcceptAnswer,
  initialData,
}: {
  entityId: EntityId
  commentCount?: number
  title?: string
  acceptedCommentId?: number
  allowAcceptAnswer?: boolean
  initialData?: PageData<Comment>
}) {
  const onCreated = React.useCallback(
    () => window.dispatchEvent(new Event(topicCommentCreatedEvent(entityId))),
    [entityId]
  )

  return (
    <CommentSection
      entityId={entityId}
      entityType="topic"
      commentCount={commentCount}
      title={title}
      acceptedCommentId={acceptedCommentId}
      allowAcceptAnswer={allowAcceptAnswer}
      initialData={initialData}
      onCreated={onCreated}
    />
  )
}
