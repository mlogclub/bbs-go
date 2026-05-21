"use client"

import { CommentSection } from "@/components/comment"
import type { Comment, PageData, SiteConfig, UserSummary } from "@/lib/api/types"

export function ArticleComments({
  entityId,
  commentCount,
  initialComments,
}: {
  entityId: string | number
  commentCount?: number
  initialComments: PageData<Comment>
  currentUser?: UserSummary | null
  config?: SiteConfig | null
}) {
  return (
    <CommentSection
      entityType="article"
      entityId={String(entityId)}
      commentCount={commentCount}
      initialData={initialComments}
    />
  )
}
