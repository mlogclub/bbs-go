"use client"

import * as React from "react"

import { DashboardImageUpload } from "@/components/dashboard/image-upload"
import {
  DashboardMultiSelect,
  DashboardSelect,
} from "@/components/dashboard/dashboard-select"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import type { AdminFormValue } from "@/lib/api/admin"
import { cn } from "@/lib/utils"

import type {
  DashboardDataFormField,
  DashboardDataOption,
} from "./dashboard-data-types"
import { DASHBOARD_DATA_ICON_OPTIONS } from "./dashboard-data-utils"

export function DashboardDataFieldControl({
  field,
  value,
  options,
  disabled,
  error,
  onChange,
}: {
  field: DashboardDataFormField
  value: AdminFormValue
  options: DashboardDataOption[]
  disabled?: boolean
  error?: string
  onChange: (value: AdminFormValue) => void
}) {
  const [optionSearch, setOptionSearch] = React.useState("")
  const isSearchableOptions = field.type === "tree-select"
  const filteredOptions =
    isSearchableOptions && optionSearch.trim()
      ? options.filter((option) =>
          option.label.toLowerCase().includes(optionSearch.trim().toLowerCase())
        )
      : options
  const filteredIconOptions = optionSearch.trim()
    ? DASHBOARD_DATA_ICON_OPTIONS.filter((option) =>
        option.value.toLowerCase().includes(optionSearch.trim().toLowerCase())
      )
    : DASHBOARD_DATA_ICON_OPTIONS

  return (
    <div
      className={cn(
        "grid gap-1.5",
        (field.colSpan === 2 || field.type === "textarea") && "md:col-span-2"
      )}
    >
      <Label>{field.label}</Label>
      {field.type === "select" ? (
        <DashboardSelect
          value={value}
          options={options}
          placeholder={field.label}
          disabled={disabled}
          onValueChange={(nextValue) => onChange(nextValue)}
        />
      ) : field.type === "tree-select" ? (
        <div className="grid gap-2 rounded-md border border-input p-2">
          <Input
            className="h-8"
            value={optionSearch}
            placeholder={field.label}
            disabled={disabled}
            onChange={(event) => setOptionSearch(event.target.value)}
          />
          <div className="max-h-56 overflow-y-auto">
            <button
              type="button"
              className={cn(
                "flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-sm hover:bg-muted",
                (value === undefined || value === "") && "bg-muted"
              )}
              disabled={disabled}
              onClick={() => onChange(undefined)}
            >
              <span className="flex size-4 items-center justify-center rounded-full border">
                {value === undefined || value === "" ? (
                  <span className="size-2 rounded-full bg-primary" />
                ) : null}
              </span>
              <span>-</span>
            </button>
            {filteredOptions.length ? (
              filteredOptions.map((option) => {
                const selected = String(value ?? "") === String(option.value)
                return (
                  <button
                    key={String(option.value)}
                    type="button"
                    className={cn(
                      "flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-sm hover:bg-muted",
                      selected && "bg-muted"
                    )}
                    disabled={disabled}
                    onClick={() => onChange(option.value)}
                  >
                    <span className="flex size-4 items-center justify-center rounded-full border">
                      {selected ? (
                        <span className="size-2 rounded-full bg-primary" />
                      ) : null}
                    </span>
                    <span className="whitespace-pre">{option.label}</span>
                  </button>
                )
              })
            ) : (
              <div className="px-2 py-3 text-sm text-muted-foreground">-</div>
            )}
          </div>
        </div>
      ) : field.type === "multiselect" ? (
        <DashboardMultiSelect
          value={value}
          options={options}
          placeholder={field.label}
          disabled={disabled}
          onValueChange={(nextValues) => onChange(nextValues)}
        />
      ) : field.type === "icon" ? (
        <div className="grid gap-2">
          <Input
            value={value === undefined || value === null ? "" : String(value)}
            disabled={disabled}
            onChange={(event) => onChange(event.target.value)}
          />
          <Input
            className="h-8"
            value={optionSearch}
            placeholder={field.label}
            disabled={disabled}
            onChange={(event) => setOptionSearch(event.target.value)}
          />
          <div className="grid max-h-56 grid-cols-6 gap-1 overflow-y-auto">
            {filteredIconOptions.length ? (
              filteredIconOptions.map(({ value: iconValue, Icon }) => (
                <button
                  key={iconValue}
                  type="button"
                  className={cn(
                    "flex size-9 items-center justify-center rounded-md border hover:bg-muted",
                    String(value || "") === iconValue &&
                      "border-primary bg-muted"
                  )}
                  disabled={disabled}
                  title={iconValue}
                  onClick={() => onChange(iconValue)}
                >
                  <Icon />
                </button>
              ))
            ) : (
              <div className="col-span-6 px-2 py-3 text-sm text-muted-foreground">
                -
              </div>
            )}
          </div>
        </div>
      ) : field.type === "image" ? (
        <DashboardImageUpload
          value={value === undefined || value === null ? "" : String(value)}
          onChange={onChange}
        />
      ) : field.type === "textarea" ? (
        <Textarea
          className="min-h-24"
          value={value === undefined || value === null ? "" : String(value)}
          disabled={disabled}
          onChange={(event) => onChange(event.target.value)}
        />
      ) : (
        <Input
          type={
            field.type === "number"
              ? "number"
              : field.type === "password"
                ? "password"
                : field.type === "url"
                  ? "url"
                  : "text"
          }
          min={field.type === "number" ? field.min : undefined}
          max={field.type === "number" ? field.max : undefined}
          step={field.type === "number" ? field.step : undefined}
          value={value === undefined || value === null ? "" : String(value)}
          disabled={disabled}
          onChange={(event) => onChange(event.target.value)}
        />
      )}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  )
}
