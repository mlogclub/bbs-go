import { apiFetch, toFormData } from "@/lib/api/client"
import type {
  Favorite,
  PageData,
  ScoreLog,
  UserMessage,
  UserSummary,
} from "@/lib/api/types"
import { ApiError } from "@/lib/api/client"

export interface UserActionState {
  ok: boolean
  message?: string
  profile?: {
    avatar: string
    nickname: string
    description: string
    homePage: string
  }
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback
}

function formString(formData: FormData, key: string) {
  const value = formData.get(key)
  return typeof value === "string" ? value : ""
}

export async function loadFavorites(cursor?: string) {
  return apiFetch<PageData<Favorite>>("/api/user/favorites", {
    params: { cursor },
  })
}

export async function loadMessages(cursor?: string) {
  return apiFetch<PageData<UserMessage>>("/api/user/messages", {
    params: { cursor },
  })
}

export async function loadScoreLogs(cursor?: string) {
  return apiFetch<PageData<ScoreLog>>("/api/user/score_logs", {
    params: { cursor },
  })
}

export async function loadFans(userId: string, cursor?: string) {
  return apiFetch<PageData<UserSummary>>("/api/fans/fans", {
    params: { userId, cursor },
  })
}

export async function loadFollowed(userId: string, cursor?: string) {
  return apiFetch<PageData<UserSummary>>("/api/fans/followed", {
    params: { userId, cursor },
  })
}

export async function requestEmailVerifyAction(): Promise<UserActionState> {
  try {
    await apiFetch<null>("/api/user/send_verify_email", {
      method: "POST",
    })
    return { ok: true }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function verifyEmailAction(token: string) {
  try {
    const data = await apiFetch<{ email: string }>("/api/user/verify_email", {
      method: "POST",
      params: { token },
    })
    return { ok: true, email: data.email }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "") }
  }
}

export async function saveProfileAction(
  _state: UserActionState,
  formData: FormData
): Promise<UserActionState> {
  const userId = formString(formData, "userId")
  const profile = {
    avatar: formString(formData, "avatar"),
    nickname: formString(formData, "nickname"),
    description: formString(formData, "description"),
    homePage: formString(formData, "homePage"),
  }
  const body = toFormData({
    avatar: profile.avatar,
    nickname: profile.nickname,
    description: profile.description,
    homePage: profile.homePage,
  })

  try {
    await apiFetch<null>(`/api/user/update/${userId}`, {
      method: "POST",
      body,
    })
    return { ok: true, profile }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function setUsernameAction(
  _state: UserActionState,
  formData: FormData
): Promise<UserActionState> {
  try {
    await apiFetch<null>("/api/user/set_username", {
      method: "POST",
      body: toFormData({ username: formString(formData, "username") }),
    })
    return { ok: true }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function setEmailAction(
  _state: UserActionState,
  formData: FormData
): Promise<UserActionState> {
  try {
    await apiFetch<null>("/api/user/set_email", {
      method: "POST",
      body: toFormData({ email: formString(formData, "email") }),
    })
    return { ok: true }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function setPasswordAction(
  _state: UserActionState,
  formData: FormData
): Promise<UserActionState> {
  try {
    await apiFetch<null>("/api/user/set_password", {
      method: "POST",
      body: toFormData({
        password: formString(formData, "password"),
        rePassword: formString(formData, "rePassword"),
      }),
    })
    return { ok: true }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function updatePasswordAction(
  _state: UserActionState,
  formData: FormData
): Promise<UserActionState> {
  try {
    await apiFetch<null>("/api/user/update_password", {
      method: "POST",
      body: toFormData({
        oldPassword: formString(formData, "oldPassword"),
        password: formString(formData, "password"),
        rePassword: formString(formData, "rePassword"),
      }),
    })
    return { ok: true }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function followAction(
  userId: string,
  followed: boolean
): Promise<UserActionState & { followed?: boolean }> {
  try {
    await apiFetch<null>(followed ? "/api/fans/unfollow" : "/api/fans/follow", {
      method: "POST",
      body: toFormData({ userId }),
    })
    return { ok: true, followed: !followed }
  } catch (error) {
    if (error instanceof ApiError && error.errorCode === 1) {
      return { ok: false, message: error.message }
    }
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}

export async function unbindProviderAction(
  provider: "wx" | "google" | "github"
): Promise<UserActionState> {
  const path =
    provider === "wx"
      ? "/api/login/wx_unbind"
      : provider === "google"
        ? "/api/login/google_unbind"
        : "/api/login/github_unbind"
  try {
    await apiFetch<null>(path, { method: "POST" })
    return { ok: true }
  } catch (error) {
    return { ok: false, message: errorMessage(error, "Failed") }
  }
}
