"use client"

import * as React from "react"

import { TopicHideContent } from "@/components/topic/topic-hide-content"
import { apiFetch } from "@/lib/api/client"
import type { EntityId, TopicHideContent as TopicHideContentType } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { useToastActions } from "@/lib/toast"

export const topicCommentCreatedEvent = (topicId: EntityId) =>
  `bbsgo:topic-comment-created:${topicId}`

export function TopicHideContentLive({
  topicId,
  initialHideContent,
}: {
  topicId: EntityId
  initialHideContent?: TopicHideContentType | null
}) {
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [hideContent, setHideContent] = React.useState(initialHideContent)

  React.useEffect(() => {
    setHideContent(initialHideContent)
  }, [initialHideContent])

  const refresh = React.useCallback(async () => {
    try {
      const data = await apiFetch<TopicHideContentType>("/api/topic/hide_content", {
        params: { topicId },
      })
      setHideContent(data)
    } catch (error) {
      catchError(error)
    }
  }, [catchError, topicId])

  React.useEffect(() => {
    const eventName = topicCommentCreatedEvent(topicId)
    window.addEventListener(eventName, refresh)
    return () => window.removeEventListener(eventName, refresh)
  }, [refresh, topicId])

  return <TopicHideContent hideContent={hideContent} t={t} />
}
