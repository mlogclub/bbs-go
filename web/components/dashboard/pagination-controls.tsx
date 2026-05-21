"use client"

import * as React from "react"

import { DashboardSelect } from "@/components/dashboard/dashboard-select"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useI18n } from "@/lib/i18n/provider"

const PAGE_SIZES = [10, 20, 50, 100]

export function DashboardPagination({
  page,
  pageCount,
  total,
  limit,
  loading,
  onPageChange,
  onLimitChange,
}: {
  page: number
  pageCount: number
  total: number
  limit: number
  loading?: boolean
  onPageChange: (page: number) => void
  onLimitChange: (limit: number) => void
}) {
  const { t } = useI18n()
  const [jumpPage, setJumpPage] = React.useState(String(page))

  React.useEffect(() => {
    setJumpPage(String(page))
  }, [page])

  function submitJump(event: React.FormEvent) {
    event.preventDefault()
    const nextPage = Math.max(1, Math.min(pageCount, Number(jumpPage) || 1))
    onPageChange(nextPage)
  }

  return (
    <div className="flex flex-wrap items-center justify-between gap-3 border-t px-4 py-3 text-sm">
      <span className="text-muted-foreground">
        {t("dashboard.pagination.total", { total })}
      </span>
      <div className="flex flex-wrap items-center gap-2">
        <DashboardSelect
          value={limit}
          options={PAGE_SIZES.map((size) => ({
            label: String(size),
            value: size,
          }))}
          triggerClassName="h-8 w-24"
          allowClear={false}
          onValueChange={(value) => onLimitChange(Number(value || limit))}
        />
        <Button
          variant="outline"
          size="sm"
          disabled={page <= 1 || loading}
          onClick={() => onPageChange(page - 1)}
        >
          {t("dashboard.pagination.previous")}
        </Button>
        <span className="min-w-20 text-center">
          {page} / {pageCount}
        </span>
        <Button
          variant="outline"
          size="sm"
          disabled={page >= pageCount || loading}
          onClick={() => onPageChange(page + 1)}
        >
          {t("dashboard.pagination.next")}
        </Button>
        <form className="flex items-center gap-1" onSubmit={submitJump}>
          <Input
            className="h-8 w-16"
            type="number"
            min={1}
            max={pageCount}
            value={jumpPage}
            onChange={(event) => setJumpPage(event.target.value)}
          />
          <Button variant="outline" size="sm" type="submit" disabled={loading}>
            {t("dashboard.pagination.jump")}
          </Button>
        </form>
      </div>
    </div>
  )
}
