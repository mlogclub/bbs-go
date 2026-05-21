import { serverApiFetch as apiFetch } from "./server"

import { toFormData } from "./client"
import type {
  Article,
  Badge,
  BindInfo,
  Favorite,
  PageData,
  ScoreLog,
  SearchUser,
  Topic,
  UserMessage,
  UserSummary,
} from "./types"

type SearchUserParams = {
  keyword: string
  cursor?: string
}

export function getCurrentUser() {
  return apiFetch<UserSummary | null>("/api/user/current")
}

export function searchUsers(params: SearchUserParams) {
  return apiFetch<PageData<SearchUser>>("/api/search/user", {
    params,
  })
}

export function getScoreRank() {
  return apiFetch<UserSummary[]>("/api/user/score/rank")
}

export function getUser(userId: string) {
  return apiFetch<UserSummary>(`/api/user/${userId}`)
}

export function getUserTopics(userId: string, cursor?: string) {
  return apiFetch<PageData<Topic>>("/api/topic/user_topics", {
    params: { userId, cursor },
  })
}

export function getUserArticles(userId: string, cursor?: string) {
  return apiFetch<PageData<Article>>("/api/article/user_articles", {
    params: { userId, cursor },
  })
}

export function getUserFavorites(cursor?: string) {
  return apiFetch<PageData<Favorite>>("/api/user/favorites", {
    params: { cursor },
  })
}

export function getUserMessages(cursor?: string) {
  return apiFetch<PageData<UserMessage>>("/api/user/messages", {
    params: { cursor },
  })
}

export function getRecentUserMessages() {
  return apiFetch<{ count?: number; messages?: UserMessage[] }>(
    "/api/user/msg_recent"
  )
}

export function getUserScoreLogs(cursor?: string) {
  return apiFetch<PageData<ScoreLog>>("/api/user/score_logs", {
    params: { cursor },
  })
}

export function getUserFans(userId: string, cursor?: string) {
  return apiFetch<PageData<UserSummary>>("/api/fans/fans", {
    params: { userId, cursor },
  })
}

export function getUserFollowed(userId: string, cursor?: string) {
  return apiFetch<PageData<UserSummary>>("/api/fans/followed", {
    params: { userId, cursor },
  })
}

export function getRecentFans(userId: string) {
  return apiFetch<PageData<UserSummary>>("/api/fans/recent/fans", {
    params: { userId },
  })
}

export function getRecentFollowed(userId: string) {
  return apiFetch<PageData<UserSummary>>("/api/fans/recent/follow", {
    params: { userId },
  })
}

export function getBadges(userId: string) {
  return apiFetch<Badge[]>("/api/badge/badges", {
    params: { userId },
  })
}

export function getBindInfo(provider: "wx" | "google" | "github") {
  const path =
    provider === "wx"
      ? "/api/user/wx_bind_info"
      : provider === "google"
        ? "/api/user/google_bind_info"
        : "/api/user/github_bind_info"
  return apiFetch<BindInfo>(path)
}

export async function updateUser(
  userId: string,
  body: FormData
): Promise<void> {
  await apiFetch<null>(`/api/user/update/${userId}`, {
    method: "POST",
    body,
  })
}

export async function requestEmailVerify(): Promise<void> {
  await apiFetch<null>("/api/user/send_verify_email", { method: "POST" })
}

export function verifyEmail(token: string) {
  return apiFetch<{ email: string }>("/api/user/verify_email", {
    method: "POST",
    params: { token },
  })
}

export async function setUsername(username: string): Promise<void> {
  await apiFetch<null>("/api/user/set_username", {
    method: "POST",
    body: toFormData({ username }),
  })
}

export async function setEmail(email: string): Promise<void> {
  await apiFetch<null>("/api/user/set_email", {
    method: "POST",
    body: toFormData({ email }),
  })
}

export async function setPassword(values: {
  password: string
  rePassword: string
}): Promise<void> {
  await apiFetch<null>("/api/user/set_password", {
    method: "POST",
    body: toFormData(values),
  })
}

export async function updatePassword(values: {
  oldPassword: string
  password: string
  rePassword: string
}): Promise<void> {
  await apiFetch<null>("/api/user/update_password", {
    method: "POST",
    body: toFormData(values),
  })
}

export async function toggleFollow(
  userId: string,
  followed: boolean
): Promise<void> {
  await apiFetch<null>(followed ? "/api/fans/unfollow" : "/api/fans/follow", {
    method: "POST",
    body: toFormData({ userId }),
  })
}

export async function unbindProvider(
  provider: "wx" | "google" | "github"
): Promise<void> {
  const path =
    provider === "wx"
      ? "/api/login/wx_unbind"
      : provider === "google"
        ? "/api/login/google_unbind"
        : "/api/login/github_unbind"
  await apiFetch<null>(path, { method: "POST" })
}
