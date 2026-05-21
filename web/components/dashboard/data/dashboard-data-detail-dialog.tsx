"use client"

import type { AdminRecord } from "@/lib/api/admin"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"

import type { DashboardDataDetailField } from "./dashboard-data-types"
import { getDashboardDataValue, textValue } from "./dashboard-data-utils"

export function DashboardDataDetailDialog({
  record,
  fields,
  title,
  cancelLabel,
  onClose,
}: {
  record: AdminRecord | null
  fields?: DashboardDataDetailField[]
  title: string
  cancelLabel: string
  onClose: () => void
}) {
  if (!record || !fields?.length) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 p-4 backdrop-blur-sm">
      <div className="max-h-[90vh] w-full max-w-3xl overflow-y-auto rounded-xl border bg-card p-4 shadow-lg">
        <div className="mb-4 flex items-center justify-between gap-3">
          <h2 className="text-lg font-semibold">{title}</h2>
          <Button type="button" variant="ghost" onClick={onClose}>
            {cancelLabel}
          </Button>
        </div>
        <div className="grid gap-4">
          {fields.map((field) => (
            <div key={field.key} className="grid gap-1.5">
              <Label className="text-xs text-muted-foreground">
                {field.label}
              </Label>
              <div className="min-h-9 rounded-md border bg-muted/20 px-3 py-2 text-sm">
                {field.render
                  ? field.render(record)
                  : textValue(getDashboardDataValue(record, field.key))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

