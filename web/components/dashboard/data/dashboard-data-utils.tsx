"use client"

import {
  BarChart3Icon,
  BadgeIcon,
  BellIcon,
  BookOpenIcon,
  BriefcaseBusinessIcon,
  CircleUserRoundIcon,
  CogIcon,
  DatabaseIcon,
  FileTextIcon,
  FolderIcon,
  GlobeIcon,
  HashIcon,
  HomeIcon,
  ImageIcon,
  KeyRoundIcon,
  LayoutDashboardIcon,
  LinkIcon,
  ListIcon,
  LockIcon,
  MailIcon,
  MedalIcon,
  MenuIcon,
  MessageSquareIcon,
  NewspaperIcon,
  PackageIcon,
  PencilIcon,
  SettingsIcon,
  ShieldIcon,
  SlidersHorizontalIcon,
  StarIcon,
  TagsIcon,
  TrophyIcon,
  UploadIcon,
  UsersIcon,
  WrenchIcon,
} from "lucide-react"

import type {
  AdminFormValue,
  AdminPrimitive,
  AdminRecord,
} from "@/lib/api/admin"
import { formatDateTime } from "@/lib/format"

export const DASHBOARD_DATA_DEPTH_KEY = "__dashboardDepth"
export const DASHBOARD_DATA_KEY = "__dashboardKey"
export const DASHBOARD_DATA_HAS_CHILDREN_KEY = "__dashboardHasChildren"
export const DASHBOARD_DATA_PARENT_KEY = "__dashboardParentKey"
export const DASHBOARD_DATA_OPTION_DEPTH_KEY = "__dashboardOptionDepth"

export const DASHBOARD_DATA_ICON_OPTIONS = [
  { value: "LayoutDashboard", Icon: LayoutDashboardIcon },
  { value: "Home", Icon: HomeIcon },
  { value: "Menu", Icon: MenuIcon },
  { value: "Settings", Icon: SettingsIcon },
  { value: "Cog", Icon: CogIcon },
  { value: "SlidersHorizontal", Icon: SlidersHorizontalIcon },
  { value: "Users", Icon: UsersIcon },
  { value: "CircleUserRound", Icon: CircleUserRoundIcon },
  { value: "Shield", Icon: ShieldIcon },
  { value: "Lock", Icon: LockIcon },
  { value: "KeyRound", Icon: KeyRoundIcon },
  { value: "FileText", Icon: FileTextIcon },
  { value: "Newspaper", Icon: NewspaperIcon },
  { value: "BookOpen", Icon: BookOpenIcon },
  { value: "Mail", Icon: MailIcon },
  { value: "MessageSquare", Icon: MessageSquareIcon },
  { value: "Link", Icon: LinkIcon },
  { value: "Globe", Icon: GlobeIcon },
  { value: "Image", Icon: ImageIcon },
  { value: "Folder", Icon: FolderIcon },
  { value: "List", Icon: ListIcon },
  { value: "Tags", Icon: TagsIcon },
  { value: "Hash", Icon: HashIcon },
  { value: "Badge", Icon: BadgeIcon },
  { value: "Medal", Icon: MedalIcon },
  { value: "Trophy", Icon: TrophyIcon },
  { value: "Star", Icon: StarIcon },
  { value: "BarChart3", Icon: BarChart3Icon },
  { value: "Database", Icon: DatabaseIcon },
  { value: "Package", Icon: PackageIcon },
  { value: "BriefcaseBusiness", Icon: BriefcaseBusinessIcon },
  { value: "Upload", Icon: UploadIcon },
  { value: "Pencil", Icon: PencilIcon },
  { value: "Wrench", Icon: WrenchIcon },
]

export function getDashboardDataValue(record: AdminRecord, key: string) {
  return key.split(".").reduce<unknown>((current, part) => {
    if (current && typeof current === "object" && part in current) {
      return (current as AdminRecord)[part]
    }
    return undefined
  }, record)
}

export function toDashboardDataPrimitive(value: unknown): AdminPrimitive {
  if (
    value === undefined ||
    value === null ||
    typeof value === "string" ||
    typeof value === "number" ||
    typeof value === "boolean"
  ) {
    return value
  }
  return String(value)
}

export function dashboardDataRecordToFormValues(record: AdminRecord) {
  return Object.fromEntries(
    Object.entries(record).map(([key, value]) => [
      key,
      toDashboardDataPrimitive(value),
    ])
  ) as Record<string, AdminFormValue>
}

export function isValidDashboardDataHttpUrl(value: string) {
  try {
    const url = new URL(value)
    return url.protocol === "http:" || url.protocol === "https:"
  } catch {
    return false
  }
}

export function normalizeDashboardDataOptionRecords(data: unknown) {
  function flatten(records: AdminRecord[], depth = 0): AdminRecord[] {
    return records.flatMap((record) => {
      const children = Array.isArray(record.children)
        ? (record.children as AdminRecord[])
        : []
      return [
        { ...record, [DASHBOARD_DATA_OPTION_DEPTH_KEY]: depth },
        ...flatten(children, depth + 1),
      ]
    })
  }

  if (Array.isArray(data)) return flatten(data as AdminRecord[])
  if (
    data &&
    typeof data === "object" &&
    "results" in data &&
    Array.isArray((data as { results?: unknown }).results)
  ) {
    return flatten((data as { results: AdminRecord[] }).results)
  }
  return []
}

export function flattenDashboardDataTree(
  records: AdminRecord[],
  depth = 0,
  parentKey = "root"
): AdminRecord[] {
  return records.flatMap((record, index) => {
    const recordKey = String(record.id ?? `${parentKey}-${index}`)
    const children = Array.isArray(record.children)
      ? (record.children as AdminRecord[])
      : []
    const current = {
      ...record,
      [DASHBOARD_DATA_DEPTH_KEY]: depth,
      [DASHBOARD_DATA_KEY]: recordKey,
      [DASHBOARD_DATA_HAS_CHILDREN_KEY]: children.length > 0,
      [DASHBOARD_DATA_PARENT_KEY]: parentKey,
    }
    return [
      current,
      ...flattenDashboardDataTree(children, depth + 1, recordKey),
    ]
  })
}

export function filterVisibleDashboardDataTree(
  records: AdminRecord[],
  collapsedKeys: Set<string>
): AdminRecord[] {
  let hiddenBelowDepth: number | null = null

  return records.filter((record) => {
    const depth = Number(record[DASHBOARD_DATA_DEPTH_KEY] ?? 0)
    if (hiddenBelowDepth !== null) {
      if (depth > hiddenBelowDepth) return false
      hiddenBelowDepth = null
    }

    const recordKey = String(record[DASHBOARD_DATA_KEY] ?? "")
    if (recordKey && collapsedKeys.has(recordKey)) {
      hiddenBelowDepth = depth
    }
    return true
  })
}

export function textValue(value: unknown) {
  if (value === undefined || value === null || value === "") return "-"
  if (typeof value === "boolean") return value ? "Yes" : "No"
  return String(value)
}

export function dateCell(value: unknown) {
  const text = formatDateTime(value as string | number | null)
  return text || "-"
}

export function imageCell(value: unknown, alt = "") {
  const src = typeof value === "string" ? value : ""
  if (!src) return "-"

  return (
    <img
      src={src}
      alt={alt}
      className="size-10 rounded-md border object-cover"
      loading="lazy"
    />
  )
}
