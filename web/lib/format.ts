import type { TFunction } from "@/lib/i18n"

export function formatDate(timestamp?: number | string, format = "yyyy-MM-dd HH:mm") {
  if (!timestamp) {
    return ""
  }

  const date = toDate(timestamp)
  if (!date) {
    return ""
  }
  const pad = (value: number) => String(value).padStart(2, "0")

  return format
    .replace("yyyy", String(date.getFullYear()))
    .replace("MM", pad(date.getMonth() + 1))
    .replace("dd", pad(date.getDate()))
    .replace("HH", pad(date.getHours()))
    .replace("mm", pad(date.getMinutes()))
    .replace("ss", pad(date.getSeconds()))
}

export function formatDateTime(timestamp?: number | string | null) {
  return formatDate(normalizeTimestamp(timestamp), "yyyy-MM-dd HH:mm:ss")
}

export function prettyDate(timestamp: number | undefined, t: TFunction) {
  if (!timestamp) {
    return ""
  }

  const minute = 1000 * 60
  const hour = minute * 60
  const day = hour * 24
  const diffValue = Date.now() - timestamp

  if (diffValue / minute < 1) {
    return t("composables.justNow")
  }
  if (diffValue / minute < 60) {
    return t("composables.minutesAgo", { n: Math.floor(diffValue / minute) })
  }
  if (diffValue / hour <= 24) {
    return t("composables.hoursAgo", { n: Math.floor(diffValue / hour) })
  }
  if (diffValue / day <= 30) {
    return t("composables.daysAgo", { n: Math.floor(diffValue / day) })
  }
  return formatDate(timestamp)
}

function normalizeTimestamp(timestamp?: number | string | null) {
  if (timestamp === null || timestamp === undefined || timestamp === "") {
    return undefined
  }
  if (typeof timestamp === "number") {
    return timestamp
  }
  const numericTimestamp = Number(timestamp)
  return Number.isNaN(numericTimestamp) ? timestamp : numericTimestamp
}

function toDate(timestamp: number | string) {
  const date = new Date(timestamp)
  return Number.isNaN(date.getTime()) ? null : date
}
