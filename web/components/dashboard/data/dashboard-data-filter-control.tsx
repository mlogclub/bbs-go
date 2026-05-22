"use client"

import type { AdminFormValue } from "@/lib/api/admin"

import { DashboardSelect } from "@/components/dashboard/dashboard-select"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

import type {
  DashboardDataFilter,
  DashboardDataOption,
} from "./dashboard-data-types"

export function DashboardDataFilterControl({
  filter,
  value,
  options,
  onChange,
}: {
  filter: DashboardDataFilter
  value: AdminFormValue
  options: DashboardDataOption[]
  onChange: (value: AdminFormValue) => void
}) {
  return (
    <div className="grid min-w-44 gap-1.5">
      <Label className="text-xs text-muted-foreground">{filter.label}</Label>
      {filter.type === "select" ? (
        <DashboardSelect
          value={value}
          options={options}
          placeholder={filter.label}
          onValueChange={(nextValue) => onChange(nextValue)}
        />
      ) : (
        <Input
          value={value === undefined || value === null ? "" : String(value)}
          placeholder={filter.label}
          onChange={(event) => onChange(event.target.value)}
        />
      )}
    </div>
  )
}

