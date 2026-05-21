"use client"

import type * as React from "react"

import type { AdminFormValue } from "@/lib/api/admin"
import { Button } from "@/components/ui/button"
import { DashboardDialog } from "@/components/dashboard/dashboard-dialog"

import { DashboardDataFieldControl } from "./dashboard-data-field-control"
import type {
  DashboardDataFormField,
  DashboardDataOption,
} from "./dashboard-data-types"

function isWideFormField(field: DashboardDataFormField) {
  return (
    field.colSpan === 2 || field.type === "textarea" || field.type === "image"
  )
}

function normalizeFormLayoutFields(fields: DashboardDataFormField[]) {
  const layoutFields: DashboardDataFormField[] = []
  let pendingHalfWidthIndex: number | null = null

  fields.forEach((field) => {
    const wide = isWideFormField(field)
    if (wide && pendingHalfWidthIndex !== null) {
      layoutFields[pendingHalfWidthIndex] = {
        ...layoutFields[pendingHalfWidthIndex],
        colSpan: 2,
      }
      pendingHalfWidthIndex = null
    }

    const nextField = wide ? { ...field, colSpan: 2 as const } : field
    layoutFields.push(nextField)

    if (wide) {
      pendingHalfWidthIndex = null
    } else if (pendingHalfWidthIndex === null) {
      pendingHalfWidthIndex = layoutFields.length - 1
    } else {
      pendingHalfWidthIndex = null
    }
  })

  return layoutFields
}

export function DashboardDataFormDialog({
  open,
  formId,
  title,
  fields,
  values,
  errors,
  asyncOptions,
  submitting,
  cancelLabel,
  confirmLabel,
  onOpenChange,
  onSubmit,
  onValueChange,
}: {
  open: boolean
  formId: string
  title: string
  fields: DashboardDataFormField[]
  values: Record<string, AdminFormValue>
  errors: Record<string, string>
  asyncOptions: Record<string, DashboardDataOption[]>
  submitting: boolean
  cancelLabel: string
  confirmLabel: string
  onOpenChange: (open: boolean) => void
  onSubmit: (event: React.FormEvent) => void
  onValueChange: (name: string, value: AdminFormValue) => void
}) {
  if (!open || !fields.length) return null
  const layoutFields = normalizeFormLayoutFields(fields)

  return (
    <DashboardDialog
      open={open}
      onOpenChange={onOpenChange}
      title={title}
      size="lg"
      footer={
        <>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
          >
            {cancelLabel}
          </Button>
          <Button type="submit" form={formId} disabled={submitting}>
            {confirmLabel}
          </Button>
        </>
      }
    >
      <form
        id={formId}
        className="grid gap-4 md:grid-cols-2"
        onSubmit={onSubmit}
      >
        {layoutFields.map((field) => (
          <DashboardDataFieldControl
            key={field.name}
            field={field}
            value={values[field.name]}
            options={[
              ...(field.options ?? []),
              ...(asyncOptions[field.name] ?? []),
            ]}
            disabled={field.name === "id"}
            error={errors[field.name]}
            onChange={(value) => onValueChange(field.name, value)}
          />
        ))}
      </form>
    </DashboardDialog>
  )
}
