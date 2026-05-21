"use client"

import {
  ArrowDownIcon,
  ArrowUpIcon,
  EditIcon,
  EyeIcon,
  Trash2Icon,
} from "lucide-react"

import type { AdminRecord } from "@/lib/api/admin"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { DashboardPagination } from "@/components/dashboard/pagination-controls"

import type { DashboardDataPageConfig } from "./dashboard-data-types"
import {
  DASHBOARD_DATA_DEPTH_KEY,
  getDashboardDataValue,
  textValue,
} from "./dashboard-data-utils"

export function DashboardDataTable({
  config,
  records,
  loading,
  page,
  pageCount,
  total,
  limit,
  labels,
  onPageChange,
  onLimitChange,
  onMove,
  canMove,
  canSort,
  canUpdate,
  canDelete,
  onRunAction,
  onView,
  onEdit,
  onDelete,
}: {
  config: DashboardDataPageConfig
  records: AdminRecord[]
  loading: boolean
  page: number
  pageCount: number
  total: number
  limit: number
  labels: {
    actions: string
    loading: string
    noData: string
    moveUp: string
    moveDown: string
    view: string
    edit: string
    delete: string
  }
  onPageChange: (page: number) => void
  onLimitChange: (limit: number) => void
  onMove: (index: number, direction: -1 | 1) => void
  canMove: (index: number, direction: -1 | 1) => boolean
  canSort: boolean
  canUpdate: boolean
  canDelete: boolean
  onRunAction: (
    action: NonNullable<DashboardDataPageConfig["rowActions"]>[number],
    record: AdminRecord
  ) => void
  onView: (record: AdminRecord) => void
  onEdit: (record: AdminRecord) => void
  onDelete: (record: AdminRecord) => void
}) {
  const hasActions =
    Boolean(config.formFields?.length && config.updateEndpoint && canUpdate) ||
    Boolean(config.detailFields?.length) ||
    Boolean(config.deleteEndpoint && canDelete) ||
    Boolean(config.sortEndpoint && canSort) ||
    Boolean(config.rowActions?.length) ||
    Boolean(config.renderRowActions)
  const colSpan = config.columns.length + (hasActions ? 1 : 0)

  return (
    <div className="overflow-hidden rounded-lg border bg-[var(--dashboard-panel)] shadow-xs">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[980px] text-sm">
          <thead className="bg-[var(--dashboard-panel-muted)] text-muted-foreground">
            <tr>
              {config.columns.map((column) => (
                <th
                  key={column.key}
                  className={cn(
                    "h-10 px-3 text-left text-xs font-semibold tracking-wide uppercase",
                    column.className
                  )}
                >
                  {column.label}
                </th>
              ))}
              {hasActions ? (
                <th className="h-10 w-48 px-3 text-right text-xs font-semibold tracking-wide uppercase">
                  {labels.actions}
                </th>
              ) : null}
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td
                  colSpan={colSpan}
                  className="px-3 py-10 text-center text-muted-foreground"
                >
                  {labels.loading}
                </td>
              </tr>
            ) : records.length ? (
              records.map((record, index) => (
                <tr
                  key={String(record.id ?? index)}
                  className="border-t transition-colors hover:bg-[var(--dashboard-accent-soft)]/45"
                >
                  {config.columns.map((column) => (
                    <td
                      key={column.key}
                      className={cn(
                        "h-11 px-3 align-middle",
                        column.className,
                        config.treeIndentKey === column.key && "font-medium"
                      )}
                      style={
                        config.treeIndentKey === column.key
                          ? {
                              paddingLeft: `${12 + Number(record[DASHBOARD_DATA_DEPTH_KEY] || 0) * 20}px`,
                            }
                          : undefined
                      }
                    >
                      {column.render
                        ? column.render(record)
                        : textValue(getDashboardDataValue(record, column.key))}
                    </td>
                  ))}
                  {hasActions ? (
                    <td className="px-3 py-2">
                      <div className="flex justify-end gap-1.5">
                        {config.sortEndpoint && canSort ? (
                          <>
                            <Button
                              size="icon-sm"
                              variant="outline"
                              disabled={!canMove(index, -1)}
                              onClick={() => onMove(index, -1)}
                            >
                              <ArrowUpIcon />
                              <span className="sr-only">{labels.moveUp}</span>
                            </Button>
                            <Button
                              size="icon-sm"
                              variant="outline"
                              disabled={!canMove(index, 1)}
                              onClick={() => onMove(index, 1)}
                            >
                              <ArrowDownIcon />
                              <span className="sr-only">{labels.moveDown}</span>
                            </Button>
                          </>
                        ) : null}
                        {config.rowActions?.map((action) => (
                          <Button
                            key={action.label}
                            size="sm"
                            variant="outline"
                            onClick={() => onRunAction(action, record)}
                          >
                            {action.label}
                          </Button>
                        ))}
                        {config.renderRowActions?.(record)}
                        {config.detailFields?.length ? (
                          <Button
                            size="icon-sm"
                            variant="outline"
                            onClick={() => onView(record)}
                          >
                            <EyeIcon />
                            <span className="sr-only">{labels.view}</span>
                          </Button>
                        ) : null}
                        {config.formFields?.length &&
                        config.updateEndpoint &&
                        canUpdate &&
                        (config.canEdit?.(record) ?? true) ? (
                          <Button
                            size="icon-sm"
                            variant="outline"
                            onClick={() => onEdit(record)}
                          >
                            <EditIcon />
                            <span className="sr-only">{labels.edit}</span>
                          </Button>
                        ) : null}
                        {config.deleteEndpoint &&
                        canDelete &&
                        (config.canDelete?.(record) ?? true) ? (
                          <Button
                            size="icon-sm"
                            variant="destructive"
                            onClick={() => onDelete(record)}
                          >
                            <Trash2Icon />
                            <span className="sr-only">{labels.delete}</span>
                          </Button>
                        ) : null}
                      </div>
                    </td>
                  ) : null}
                </tr>
              ))
            ) : (
              <tr>
                <td
                  colSpan={colSpan}
                  className="px-3 py-10 text-center text-muted-foreground"
                >
                  {labels.noData}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {config.listResult !== "array" ? (
        <DashboardPagination
          page={page}
          pageCount={pageCount}
          total={total}
          limit={limit}
          loading={loading}
          onPageChange={onPageChange}
          onLimitChange={onLimitChange}
        />
      ) : null}
    </div>
  )
}
