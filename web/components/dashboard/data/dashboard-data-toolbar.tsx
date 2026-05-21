"use client"

import { PlusIcon, RefreshCwIcon, SearchIcon } from "lucide-react"

import type { AdminFormValue } from "@/lib/api/admin"
import { Button } from "@/components/ui/button"

import { DashboardDataFilterControl } from "./dashboard-data-filter-control"
import type {
  DashboardDataFilter,
  DashboardDataOption,
} from "./dashboard-data-types"

export function DashboardDataToolbar({
  filters,
  values,
  asyncOptions,
  loading,
  canCreate,
  error,
  searchLabel,
  refreshLabel,
  createLabel,
  onFilterChange,
  onRefresh,
  onCreate,
}: {
  filters?: DashboardDataFilter[]
  values: Record<string, AdminFormValue>
  asyncOptions: Record<string, DashboardDataOption[]>
  loading: boolean
  canCreate: boolean
  error: string | null
  searchLabel: string
  refreshLabel: string
  createLabel: string
  onFilterChange: (name: string, value: AdminFormValue) => void
  onRefresh: () => void
  onCreate: () => void
}) {
  return (
    <div className="rounded-lg border bg-[var(--dashboard-panel)] p-3 text-card-foreground shadow-xs">
      <div className="flex flex-wrap items-end gap-2">
        {filters?.map((filter) => (
          <DashboardDataFilterControl
            key={filter.name}
            filter={filter}
            value={values[filter.name]}
            options={[
              ...(filter.options ?? []),
              ...(asyncOptions[filter.name] ?? []),
            ]}
            onChange={(value) => onFilterChange(filter.name, value)}
          />
        ))}
        <Button onClick={onRefresh} disabled={loading}>
          <SearchIcon />
          {searchLabel}
        </Button>
        <Button variant="outline" size="icon" onClick={onRefresh} disabled={loading}>
          <RefreshCwIcon />
          <span className="sr-only">{refreshLabel}</span>
        </Button>
        {canCreate ? (
          <Button className="ml-auto" onClick={onCreate}>
            <PlusIcon />
            {createLabel}
          </Button>
        ) : null}
      </div>

      {error ? (
        <div className="mt-3 rounded-md border border-destructive/25 bg-destructive/10 px-3 py-2 text-sm text-destructive">
          {error}
        </div>
      ) : null}
    </div>
  )
}
