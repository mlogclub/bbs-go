"use client"

import * as React from "react"

import { apiFetch, toFormData } from "@/lib/api/client"
import type { EntityId } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { toast, useToastActions } from "@/lib/toast"

type ActionSource = "side" | "detail"

type TopicActionContextValue = {
  topicId: EntityId
  liked: boolean
  favorited: boolean
  likeCount: number
  commentCount: number
  toggleLike: (source: ActionSource) => Promise<void>
  toggleFavorite: (source: ActionSource) => Promise<void>
  scrollToComment: () => void
  scrollToTop: () => void
}

const TopicActionContext = React.createContext<TopicActionContextValue | null>(
  null
)

export function TopicActionProvider({
  topicId,
  liked: initialLiked,
  favorited: initialFavorited,
  likeCount: initialLikeCount,
  commentCount,
  children,
}: {
  topicId: EntityId
  liked?: boolean
  favorited?: boolean
  likeCount?: number
  commentCount?: number
  children: React.ReactNode
}) {
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [liked, setLiked] = React.useState(Boolean(initialLiked))
  const [favorited, setFavorited] = React.useState(Boolean(initialFavorited))
  const [likeCount, setLikeCount] = React.useState(initialLikeCount || 0)
  const [likePending, setLikePending] = React.useState(false)
  const [favoritePending, setFavoritePending] = React.useState(false)

  const toggleLike = React.useCallback(
    async (source: ActionSource) => {
      if (likePending) return
      const nextLiked = !liked
      const previousLiked = liked
      const previousCount = likeCount
      setLikePending(true)
      setLiked(nextLiked)
      setLikeCount((current) =>
        nextLiked ? current + 1 : Math.max(0, current - 1)
      )

      try {
        await apiFetch(nextLiked ? "/api/like/like" : "/api/like/unlike", {
          method: "POST",
          body: toFormData({ entityType: "topic", entityId: topicId }),
        })
        if (source === "side") {
          toast.success(
            t(
              nextLiked
                ? "component.sideActionBar.likeSuccess"
                : "component.sideActionBar.likeCancel"
            )
          )
        } else {
          toast.success(t("pages.topic.detail.likeSuccess"))
        }
      } catch (error) {
        setLiked(previousLiked)
        setLikeCount(previousCount)
        catchError(error)
      } finally {
        setLikePending(false)
      }
    },
    [catchError, likeCount, likePending, liked, t, topicId]
  )

  const toggleFavorite = React.useCallback(
    async (source: ActionSource) => {
      if (favoritePending) return
      const nextFavorited = !favorited
      const previousFavorited = favorited
      setFavoritePending(true)
      setFavorited(nextFavorited)

      try {
        await apiFetch(
          nextFavorited ? "/api/favorite/add" : "/api/favorite/delete",
          {
            method: "POST",
            body: toFormData({ entityType: "topic", entityId: topicId }),
          }
        )
        if (source === "side") {
          toast.success(
            t(
              nextFavorited
                ? "component.sideActionBar.favoriteSuccess"
                : "component.sideActionBar.favoriteCancel"
            )
          )
        } else {
          toast.success(t("pages.topic.detail.favoriteSuccess"))
        }
      } catch (error) {
        setFavorited(previousFavorited)
        catchError(error)
      } finally {
        setFavoritePending(false)
      }
    },
    [catchError, favoritePending, favorited, t, topicId]
  )

  const value = React.useMemo<TopicActionContextValue>(
    () => ({
      topicId,
      liked,
      favorited,
      likeCount,
      commentCount: commentCount || 0,
      toggleLike,
      toggleFavorite,
      scrollToComment: () => {
        const element = document.getElementById("JComment")
        if (element) {
          window.scrollTo({ top: element.offsetTop, behavior: "smooth" })
        }
      },
      scrollToTop: () => window.scrollTo({ top: 0, behavior: "smooth" }),
    }),
    [
      commentCount,
      favorited,
      likeCount,
      liked,
      toggleFavorite,
      toggleLike,
      topicId,
    ]
  )

  return (
    <TopicActionContext.Provider value={value}>
      {children}
    </TopicActionContext.Provider>
  )
}

export function useTopicActions() {
  const value = React.useContext(TopicActionContext)
  if (!value) {
    throw new Error("useTopicActions must be used within TopicActionProvider")
  }
  return value
}
