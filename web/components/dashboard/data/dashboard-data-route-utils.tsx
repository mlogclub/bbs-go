"use client"

import type { TFunction } from "@/lib/i18n"

export { dateCell, imageCell } from "./dashboard-data-utils"

export function title(t: TFunction, key: string) {
  return t(`dashboard.pages.${key}.title`)
}

export function desc(t: TFunction, key: string) {
  return t(`dashboard.pages.${key}.description`)
}

export function label(t: TFunction, key: string) {
  return t(`dashboard.fields.${key}`)
}

export function userUrl(record: Record<string, unknown>) {
  const id = record.idEncode || record.id
  return id ? `/user/${String(id)}` : undefined
}

export function userLinkCell(record: Record<string, unknown>, value: unknown) {
  const text =
    value === undefined || value === null || value === "" ? "-" : String(value)
  const href = userUrl(record)
  if (!href || text === "-") return text

  return (
    <a href={href} target="_blank" rel="noreferrer" className="hover:underline">
      {text}
    </a>
  )
}

export function treeOptionLabel(record: Record<string, unknown>) {
  const depth = Number(record.__dashboardOptionDepth || 0)
  const text = String(record.title || record.name || record.id)
  return `${"  ".repeat(depth)}${depth ? "└ " : ""}${text}`
}

export function normalDeletedOptions(t: TFunction) {
  return [
    { label: t("dashboard.status.normal"), value: 0 },
    { label: t("dashboard.status.deleted"), value: 1 },
  ]
}

export function enabledDisabledOptions(t: TFunction) {
  return [
    { label: t("dashboard.status.normal"), value: 0 },
    { label: t("dashboard.status.disabled"), value: 1 },
  ]
}

export function emailStatusOptions(t: TFunction) {
  return [
    { label: t("dashboard.emailStatus.success"), value: 0 },
    { label: t("dashboard.emailStatus.failed"), value: 1 },
  ]
}

export function statusCell(t: TFunction, value: unknown) {
  return Number(value) === 0
    ? t("dashboard.status.normal")
    : t("dashboard.status.deleted")
}

export function disabledStatusCell(t: TFunction, value: unknown) {
  return Number(value) === 0
    ? t("dashboard.status.normal")
    : t("dashboard.status.disabled")
}

export function emailStatusCell(t: TFunction, value: unknown) {
  return Number(value) === 0
    ? t("dashboard.emailStatus.success")
    : t("dashboard.emailStatus.failed")
}

export function topicStatusOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.topicFeed.statusNormal"), value: 0 },
    { label: t("dashboard.topicFeed.statusDeleted"), value: 1 },
    { label: t("dashboard.topicFeed.statusReview"), value: 2 },
  ]
}

export function topicTypeOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.topicFeed.typeTopic"), value: 0 },
    { label: t("dashboard.topicFeed.typeTweet"), value: 1 },
    { label: t("dashboard.topicFeed.typeQa"), value: 2 },
  ]
}

export function booleanOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.boolean.yes"), value: "true" },
    { label: t("dashboard.boolean.no"), value: "false" },
  ]
}

export function reportDataTypeOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.reportDataTypes.topic"), value: "topic" },
    { label: t("dashboard.reportDataTypes.article"), value: "article" },
    { label: t("dashboard.reportDataTypes.comment"), value: "comment" },
    { label: t("dashboard.reportDataTypes.user"), value: "user" },
  ]
}

export function reportAuditStatusOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.reportAuditStatus.pending"), value: 0 },
    { label: t("dashboard.reportAuditStatus.processed"), value: 1 },
    { label: t("dashboard.reportAuditStatus.ignored"), value: 2 },
  ]
}

export function reportDataTypeCell(t: TFunction, value: unknown) {
  const key = String(value || "")
  if (key === "topic") return t("dashboard.reportDataTypes.topic")
  if (key === "article") return t("dashboard.reportDataTypes.article")
  if (key === "comment") return t("dashboard.reportDataTypes.comment")
  if (key === "user") return t("dashboard.reportDataTypes.user")
  return key || "-"
}

export function reportAuditStatusCell(t: TFunction, value: unknown) {
  const status = Number(value || 0)
  if (status === 0) return t("dashboard.reportAuditStatus.pending")
  if (status === 1) return t("dashboard.reportAuditStatus.processed")
  if (status === 2) return t("dashboard.reportAuditStatus.ignored")
  return t("dashboard.reportAuditStatus.unknown", { status })
}

export function reportTargetUrl(record: Record<string, unknown>) {
  const dataType = String(record.dataType || "")
  const dataId = record.dataId
  if (dataId === undefined || dataId === null || dataId === "") return undefined

  if (dataType === "topic") return `/topic/${String(dataId)}`
  if (dataType === "article") return `/article/${String(dataId)}`
  if (dataType === "user") return `/user/${String(dataId)}`
  return undefined
}

export function reportTargetCell(
  t: TFunction,
  record: Record<string, unknown>
) {
  const dataId = record.dataId
  const idText =
    dataId === undefined || dataId === null || dataId === ""
      ? "-"
      : String(dataId)
  const typeText = reportDataTypeCell(t, record.dataType)
  const text = `${typeText} #${idText}`
  const href = reportTargetUrl(record)
  if (!href || idText === "-") return text

  return (
    <a href={href} target="_blank" rel="noreferrer" className="hover:underline">
      {text}
    </a>
  )
}

export function categoryTypeOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.categoryTypes.normal"), value: "normal" },
    { label: t("dashboard.categoryTypes.qa"), value: "qa" },
  ]
}

export function codeBlock(value: unknown) {
  const text =
    typeof value === "string"
      ? value
      : value === undefined || value === null
        ? ""
        : JSON.stringify(value, null, 2)

  return (
    <pre className="max-h-96 overflow-auto text-xs leading-5 break-words whitespace-pre-wrap">
      {text || "-"}
    </pre>
  )
}

export function taskPeriodOptionsFor(t: TFunction) {
  return [
    { label: t("dashboard.period.once"), value: 0 },
    { label: t("dashboard.period.daily"), value: 1 },
  ]
}

export function roleTypeCell(t: TFunction, value: unknown) {
  return Number(value) === 0
    ? t("dashboard.roleTypes.system")
    : t("dashboard.roleTypes.normal")
}
