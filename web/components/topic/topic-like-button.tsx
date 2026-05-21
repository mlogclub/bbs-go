"use client"

import { Heart } from "lucide-react"
import * as React from "react"

import { apiFetch, toFormData } from "@/lib/api/client"
import type { EntityId } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { toast, useToastActions } from "@/lib/toast"
import { cn } from "@/lib/utils"

export function TopicLikeButton({
  topicId,
  initialLiked,
  initialLikeCount,
}: {
  topicId: EntityId
  initialLiked?: boolean
  initialLikeCount?: number
}) {
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [liked, setLiked] = React.useState(Boolean(initialLiked))
  const [likeCount, setLikeCount] = React.useState(initialLikeCount || 0)
  const [pending, setPending] = React.useState(false)

  async function toggleLike() {
    if (pending) {
      return
    }

    const nextLiked = !liked
    const previousLiked = liked
    const previousCount = likeCount
    setPending(true)
    setLiked(nextLiked)
    setLikeCount((current) => (nextLiked ? current + 1 : Math.max(0, current - 1)))

    try {
      await apiFetch<null>(nextLiked ? "/api/like/like" : "/api/like/unlike", {
        method: "POST",
        body: toFormData({ entityType: "topic", entityId: topicId }),
      })
      toast.success(t(nextLiked ? "component.topicList.likeSuccess" : "component.topicList.unlikeSuccess"))
    } catch (error) {
      setLiked(previousLiked)
      setLikeCount(previousCount)
      catchError(error)
    } finally {
      setPending(false)
    }
  }

  return (
    <button
      type="button"
      disabled={pending}
      className={cn(
        "inline-flex min-h-8 items-center gap-1.5 transition-colors hover:text-primary disabled:cursor-not-allowed disabled:opacity-70",
        liked && "text-destructive hover:text-destructive"
      )}
      onClick={toggleLike}
    >
      <Heart className="h-4 w-4" fill={liked ? "currentColor" : "none"} />
      <span className="min-w-[1ch] text-sm leading-none">{likeCount > 0 ? likeCount : ""}</span>
    </button>
  )
}
