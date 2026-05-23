"use client"

import type * as React from "react"

import type { AdminFormValue, AdminRecord } from "@/lib/api/admin"
import type { PermissionCode } from "@/lib/auth/permissions.generated"

export type DashboardDataOption = {
  label: string
  value: string | number
}

export type DashboardDataOptionSource = {
  optionsEndpoint?: string
  optionLabel?: (record: AdminRecord) => string
  optionValue?: (record: AdminRecord) => string | number
}

export type DashboardDataFilter = DashboardDataOptionSource & {
  name: string
  label: string
  type?: "text" | "select"
  options?: DashboardDataOption[]
}

export type DashboardDataFormField = DashboardDataOptionSource & {
  name: string
  label: string
  required?: boolean
  colSpan?: 1 | 2
  type?:
    | "text"
    | "textarea"
    | "number"
    | "select"
    | "tree-select"
    | "multiselect"
    | "password"
    | "url"
    | "image"
    | "icon"
  options?: DashboardDataOption[]
  min?: number
  max?: number
  step?: number
  valueFromRecord?: (record: AdminRecord) => AdminFormValue
}

export type DashboardDataColumn = {
  key: string
  label: string
  className?: string
  render?: (record: AdminRecord) => React.ReactNode
}

export type DashboardDataRowAction = {
  label: string
  endpoint: string
  permission?: PermissionCode
  method?: "POST" | "DELETE"
  payload?: (record: AdminRecord) => Record<string, AdminFormValue>
  confirm?: string
  successMessage?: string
}

export type DashboardDataDetailField = {
  key: string
  label: string
  render?: (record: AdminRecord) => React.ReactNode
}

export type DashboardDataPageConfig = {
  title: string
  description?: string
  listEndpoint: string
  viewPermission?: PermissionCode
  detailEndpoint?: (id: AdminFormValue) => string
  createEndpoint?: string
  createPermission?: PermissionCode
  updateEndpoint?: string
  updatePermission?: PermissionCode
  deleteEndpoint?: string
  deletePermission?: PermissionCode
  deleteMode?: "formId" | "formIds" | "jsonIds"
  sortEndpoint?: string
  sortPermission?: PermissionCode
  dragSort?: boolean
  tree?: boolean
  treeDefaultCollapsed?: boolean
  treeIndentKey?: string
  canEdit?: (record: AdminRecord) => boolean
  canDelete?: (record: AdminRecord) => boolean
  renderRowActions?: (record: AdminRecord) => React.ReactNode
  filters?: DashboardDataFilter[]
  defaultFilters?: Record<string, AdminFormValue>
  columns: DashboardDataColumn[]
  detailFields?: DashboardDataDetailField[]
  formFields?: DashboardDataFormField[]
  rowActions?: DashboardDataRowAction[]
  pageSize?: number
  listResult?: "page" | "array"
  formContainer?: "dialog" | "drawer"
  transformSubmitValues?: (
    values: Record<string, AdminFormValue>
  ) => Record<string, AdminFormValue>
}
