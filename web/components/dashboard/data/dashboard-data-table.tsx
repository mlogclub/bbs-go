"use client"

import * as React from "react"

import {
  ArrowDownIcon,
  ArrowUpIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  EditIcon,
  EyeIcon,
  GripVerticalIcon,
  Trash2Icon,
} from "lucide-react"

import type { AdminRecord } from "@/lib/api/admin"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { DashboardPagination } from "@/components/dashboard/pagination-controls"

import type { DashboardDataPageConfig } from "./dashboard-data-types"
import {
  DASHBOARD_DATA_DEPTH_KEY,
  DASHBOARD_DATA_HAS_CHILDREN_KEY,
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
  onReorder,
  canMove,
  canSort,
  canUpdate,
  canDelete,
  onRunAction,
  onView,
  onEdit,
  onDelete,
  isTreeRecordCollapsed,
  onToggleTreeRecord,
  batchDelete,
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
    expand: string
    collapse: string
    view: string
    edit: string
    delete: string
  }
  onPageChange: (page: number) => void
  onLimitChange: (limit: number) => void
  onMove: (index: number, direction: -1 | 1) => void
  onReorder: (fromIndex: number, toIndex: number) => void
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
  isTreeRecordCollapsed?: (record: AdminRecord) => boolean
  onToggleTreeRecord?: (record: AdminRecord) => void
  /** Optional batch delete handler. When provided, a checkbox column is shown. */
  batchDelete?: (ids: string[]) => Promise<void>
}) {
  const [draggingIndex, setDraggingIndex] = React.useState<number | null>(null)
  const [dragOverIndex, setDragOverIndex] = React.useState<number | null>(null)
  const [selectedIds, setSelectedIds] = React.useState<Set<string>>(new Set())
  const [batchDeleting, setBatchDeleting] = React.useState(false)
  const canDragSort = Boolean(
    config.dragSort && config.sortEndpoint && canSort && !config.tree
  )
  const enableBatch = Boolean(batchDelete && config.deleteEndpoint)
  const hasActions =
    enableBatch ||
    Boolean(config.formFields?.length && config.updateEndpoint && canUpdate) ||
    Boolean(config.detailFields?.length) ||
    Boolean(config.deleteEndpoint && canDelete) ||
    Boolean(config.sortEndpoint && canSort) ||
    Boolean(config.rowActions?.length) ||
    Boolean(config.renderRowActions)
  const colSpan = config.columns.length + (hasActions ? 1 : 0)

  function toggleSelect(id: string) {
    setSelectedIds((current) => {
      const next = new Set(current)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  function toggleSelectAll() {
    setSelectedIds((current) => {
      if (current.size === records.length) return new Set()
      return new Set(records.map((r) => String(r.id ?? "")))
    })
  }

  async function handleBatchDelete() {
    if (!batchDelete || selectedIds.size === 0) return
    setBatchDeleting(true)
    try {
      await batchDelete(Array.from(selectedIds))
      setSelectedIds(new Set())
    } finally {
      setBatchDeleting(false)
    }
  }

  const allSelected = records.length > 0 && selectedIds.size === records.length

  return (
    <div className="overflow-hidden rounded-lg border bg-[var(--dashboard-panel)] shadow-xs">
      {enableBatch && selectedIds.size > 0 ? (
        <div className="flex items-center gap-2 border-b bg-[var(--dashboard-accent-soft)]/30 px-4 py-2 text-sm">
          <span className="font-medium text-foreground">
            已选 {selectedIds.size} 项
          </span>
          <Button
            size="sm"
            variant="destructive"
            disabled={batchDeleting}
            onClick={() => void handleBatchDelete()}
          >
            {batchDeleting ? "删除中..." : "批量删除"}
          </Button>
        </div>
      ) : null}
      <div className="overflow-x-auto">
        <table className="w-full min-w-[980px] text-sm">
          <thead className="bg-[var(--dashboard-panel-muted)] text-muted-foreground">
            <tr>
              {enableBatch ? (
                <th className="h-10 w-10 px-3 text-left">
                  <Checkbox
                    checked={allSelected}
                    onCheckedChange={() => toggleSelectAll()}
                    aria-label="全选"
                  />
                </th>
              ) : null}
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
                  onDragOver={(event) => {
                    if (!canDragSort || draggingIndex === null) return
                    event.preventDefault()
                    event.dataTransfer.dropEffect = "move"
                    if (dragOverIndex !== index) {
                      setDragOverIndex(index)
                    }
                  }}
                  onDrop={(event) => {
                    if (!canDragSort || draggingIndex === null) return
                    event.preventDefault()
                    onReorder(draggingIndex, index)
                    setDraggingIndex(null)
                    setDragOverIndex(null)
                  }}
                  className={cn(
                    "border-t transition-colors hover:bg-[var(--dashboard-accent-soft)]/45",
                    draggingIndex === index && "opacity-50",
                    dragOverIndex === index &&
                      draggingIndex !== index &&
                      "bg-[var(--dashboard-accent-soft)]/70 outline-2 -outline-offset-2 outline-primary/45"
                  )}
                >
                  {enableBatch ? (
                    <td className="h-11 w-10 px-3 align-middle">
                      <Checkbox
                        checked={selectedIds.has(String(record.id ?? ""))}
                        onCheckedChange={() =>
                          toggleSelect(String(record.id ?? ""))
                        }
                        aria-label="选择"
                      />
                    </td>
                  ) : null}
                  {config.columns.map((column) => {
                    const isTreeIndentColumn =
                      config.treeIndentKey === column.key
                    const content = column.render
                      ? column.render(record)
                      : textValue(getDashboardDataValue(record, column.key))
                    const depth = Number(record[DASHBOARD_DATA_DEPTH_KEY] || 0)
                    const hasChildren = Boolean(
                      record[DASHBOARD_DATA_HAS_CHILDREN_KEY]
                    )
                    const collapsed = Boolean(isTreeRecordCollapsed?.(record))

                    return (
                      <td
                        key={column.key}
                        className={cn(
                          "h-11 px-3 align-middle",
                          column.className,
                          isTreeIndentColumn && "font-medium"
                        )}
                        style={
                          isTreeIndentColumn
                            ? {
                                paddingLeft: `${12 + depth * 20}px`,
                              }
                            : undefined
                        }
                      >
                        {isTreeIndentColumn ? (
                          <div className="flex min-w-0 items-center gap-1.5">
                            {hasChildren ? (
                              <Button
                                type="button"
                                size="icon-sm"
                                variant="ghost"
                                className="size-6 shrink-0"
                                onClick={() => onToggleTreeRecord?.(record)}
                              >
                                {collapsed ? (
                                  <ChevronRightIcon />
                                ) : (
                                  <ChevronDownIcon />
                                )}
                                <span className="sr-only">
                                  {collapsed ? labels.expand : labels.collapse}
                                </span>
                              </Button>
                            ) : (
                              <span className="size-6 shrink-0" />
                            )}
                            <span className="min-w-0 truncate">{content}</span>
                          </div>
                        ) : (
                          content
                        )}
                      </td>
                    )
                  })}
                  {hasActions ? (
                    <td className="px-3 py-2">
                      <div className="flex justify-end gap-1.5">
                        {canDragSort ? (
                          <Button
                            type="button"
                            size="icon-sm"
                            variant="outline"
                            draggable
                            className="cursor-grab active:cursor-grabbing"
                            onDragStart={(event) => {
                              setDraggingIndex(index)
                              setDragOverIndex(index)
                              event.dataTransfer.effectAllowed = "move"
                              event.dataTransfer.setData(
                                "text/plain",
                                String(index)
                              )
                            }}
                            onDragEnd={() => {
                              setDraggingIndex(null)
                              setDragOverIndex(null)
                            }}
                          >
                            <GripVerticalIcon />
                            <span className="sr-only">{labels.moveUp}</span>
                          </Button>
                        ) : null}
                        {config.sortEndpoint && canSort && !canDragSort ? (
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
