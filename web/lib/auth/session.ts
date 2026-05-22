import { ApiError } from "@/lib/api/client"
import { getCurrentUser } from "@/lib/api/users"

function isAuthError(error: unknown) {
  if (!(error instanceof ApiError)) {
    return false
  }

  if (error.errorCode === 1 || error.status === 401 || error.status === 403) {
    return true
  }

  const message = error.message.toLowerCase()
  return message.includes("unauthorized") || message.includes("unauthenticated")
}

export async function getSessionUser() {
  try {
    return await getCurrentUser()
  } catch (error) {
    if (isAuthError(error)) {
      return null
    }

    throw error
  }
}
