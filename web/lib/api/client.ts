import type { ApiEnvelope } from "./types"

export const INSTALL_REQUIRED_STATUS = 428

export class ApiError extends Error {
  errorCode?: number
  status?: number

  constructor(
    message: string,
    options: { errorCode?: number; status?: number } = {}
  ) {
    super(message)
    this.name = "ApiError"
    this.errorCode = options.errorCode
    this.status = options.status
  }
}

type QueryValue = string | number | boolean | null | undefined

export interface ApiRequestOptions extends Omit<RequestInit, "body"> {
  params?: Record<string, QueryValue>
  body?: BodyInit | Record<string, unknown> | null
  token?: string
  request?: Request
}

function isServer() {
  return typeof window === "undefined"
}

function handleInstallRequired() {
  if (!isServer() && window.location.pathname !== "/install") {
    window.location.replace("/install")
  }

  throw new Response("Install Required", {
    status: INSTALL_REQUIRED_STATUS,
    statusText: "Install Required",
  })
}

function serverOrigin() {
  const origin = process.env.BBSGO_SERVER_URL || process.env.SERVER_URL
  if (!origin) {
    throw new Error("BBSGO_SERVER_URL is required. Set it in web/.env.")
  }
  return origin
}

function buildUrl(path: string, params?: Record<string, QueryValue>) {
  const absolute = /^https?:\/\//.test(path)
  const base = isServer() ? serverOrigin() : "http://local"
  const url = new URL(path, absolute ? undefined : base)
  Object.entries(params ?? {}).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== "") {
      url.searchParams.set(key, String(value))
    }
  })

  if (isServer() || absolute) {
    return url.toString()
  }
  return `${url.pathname}${url.search}`
}

function forwardServerHeaders(
  headers: Headers,
  request: Request | undefined,
  url: string
) {
  if (!isServer() || !request) return

  if (new URL(url).origin !== new URL(serverOrigin()).origin) return

  const cookie = request.headers.get("cookie")
  if (cookie) headers.set("cookie", cookie)

  const userAgent = request.headers.get("user-agent")
  if (userAgent) headers.set("user-agent", userAgent)

  const acceptLanguage = request.headers.get("accept-language")
  if (acceptLanguage) headers.set("accept-language", acceptLanguage)

  const token = request.headers.get("x-user-token")
  if (token) headers.set("x-user-token", token)
}

function isPlainObjectBody(
  body: ApiRequestOptions["body"]
): body is Record<string, unknown> {
  return (
    body !== null &&
    typeof body === "object" &&
    !(body instanceof FormData) &&
    !(body instanceof Blob) &&
    !(body instanceof ArrayBuffer) &&
    !(body instanceof URLSearchParams) &&
    !(body instanceof ReadableStream)
  )
}

export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const { request, ...fetchOptions } = options
  const url = buildUrl(path, fetchOptions.params)
  const headers = new Headers(fetchOptions.headers)
  forwardServerHeaders(headers, request, url)

  const token = fetchOptions.token
  if (token) {
    headers.set("X-User-Token", token)
  }

  let body = fetchOptions.body
  if (isPlainObjectBody(body)) {
    headers.set("Content-Type", "application/json")
    body = JSON.stringify(body)
  }

  const response = await fetch(url, {
    ...fetchOptions,
    body,
    headers,
    credentials: fetchOptions.credentials ?? "same-origin",
    cache: fetchOptions.cache ?? "no-store",
  })

  if (!response.ok) {
    throw new ApiError(`${response.status} ${response.statusText}`, {
      status: response.status,
    })
  }

  const envelope = (await response.json()) as ApiEnvelope<T>
  const failed =
    envelope.success === false ||
    (envelope.success === undefined &&
      envelope.errorCode !== undefined &&
      envelope.errorCode !== 0)
  if (failed) {
    if (envelope.errorCode === -1) {
      handleInstallRequired()
    }

    throw new ApiError(envelope.message || "API request failed", {
      errorCode: envelope.errorCode,
      status: response.status,
    })
  }

  return envelope.data
}

export function toFormData(values: Record<string, QueryValue>) {
  const form = new FormData()
  Object.entries(values).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      form.append(key, String(value))
    }
  })
  return form
}
