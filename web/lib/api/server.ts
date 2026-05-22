import { cookies } from "@/lib/server/headers"

import { AUTH_COOKIE } from "@/lib/cookies"

import { apiFetch, type ApiRequestOptions } from "./client"

async function readServerToken() {
  const store = await cookies()
  return store.get(AUTH_COOKIE)?.value
}

export async function serverApiFetch<T>(
  path: string,
  options: ApiRequestOptions = {}
) {
  const token = options.token ?? (await readServerToken())
  return apiFetch<T>(path, { ...options, token })
}
