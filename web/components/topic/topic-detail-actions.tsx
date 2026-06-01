"use client"

import { Eye, Heart, MessageCircle, Star, Flame } from "lucide-react"

import { useTopicActions } from "@/components/topic/topic-action-context"
import { StakeButton } from "@/components/topic/stake-button"
import type { Topic } from "@/lib/api/types"
import { cn } from "@/lib/utils"

function countLabel(value?: number) {
  return value && value > 0 ? `(${value})` : ""
}

export function TopicDetailActions({
  topic,
  labels,
}: {
  topic: Topic
  labels: {
    view: string
    like: string
    comment: string
    favorite: string
  }
}) {
  const {
    liked,
    favorited,
    likeCount,
    commentCount,
    toggleLike,
    toggleFavorite,
    scrollToComment,
  } = useTopicActions()

  return (
    <div className="mb-4 flex items-center justify-between border-t border-border px-4 py-2.5">
      <div className="flex flex-1 cursor-not-allowed items-center justify-center text-sm text-muted-foreground">
        <Eye className="size-[18px] stroke-2 text-muted-foreground" />
        <div className="ml-[5px] text-foreground">
          <span>{labels.view}</span>
          <span>{countLabel(topic.viewCount)}</span>
        </div>
      </div>
      <button
        type="button"
        className="group flex flex-1 items-center justify-center text-sm text-muted-foreground hover:text-primary"
        onClick={() => void toggleLike("detail")}
      >
        <Heart
          className={cn(
            "size-[18px] stroke-2 transition-all duration-200 group-hover:text-primary",
            liked
              ? "text-destructive group-hover:text-destructive"
              : "text-muted-foreground"
          )}
          fill={liked ? "currentColor" : "none"}
        />
        <div className="ml-[5px] text-foreground">
          <span>{labels.like}</span>
          <span>{countLabel(likeCount)}</span>
        </div>
      </button>
      <button
        type="button"
        className="group flex flex-1 items-center justify-center text-sm text-muted-foreground hover:text-primary"
        onClick={scrollToComment}
      >
        <MessageCircle className="size-[18px] stroke-2 text-muted-foreground transition-all duration-200 group-hover:text-primary" />
        <div className="ml-[5px] text-foreground">
          <span>{labels.comment}</span>
          <span>{countLabel(commentCount)}</span>
        </div>
      </button>
      <StakeButton topicId={topic.id} currentFlameLevel={topic.flameLevel || 0} />
    </div>
  )
}
