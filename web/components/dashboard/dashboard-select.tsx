"use client"

import * as React from "react"
import { CheckIcon, ChevronsUpDownIcon, SearchIcon, XIcon } from "lucide-react"
import { Popover as PopoverPrimitive } from "radix-ui"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

export type DashboardSelectOption = {
  label: string
  value: string | number | boolean
}

type DashboardSelectProps = {
  value?: unknown
  options: DashboardSelectOption[]
  placeholder?: string
  emptyLabel?: string
  emptyValue?: string
  searchPlaceholder?: string
  noResultsLabel?: string
  disabled?: boolean
  allowClear?: boolean
  className?: string
  triggerClassName?: string
  contentClassName?: string
  onValueChange: (value: string | undefined) => void
}

const DEFAULT_EMPTY_VALUE = "__empty__"

export function DashboardSelect({
  value,
  options,
  placeholder,
  emptyLabel,
  emptyValue = DEFAULT_EMPTY_VALUE,
  searchPlaceholder,
  noResultsLabel = "-",
  disabled,
  allowClear = true,
  className,
  triggerClassName,
  contentClassName,
  onValueChange,
}: DashboardSelectProps) {
  const [open, setOpen] = React.useState(false)
  const [search, setSearch] = React.useState("")
  const selectedValue =
    value === undefined || value === null || value === ""
      ? emptyValue
      : String(value)
  const selectedOption = options.find(
    (option) => String(option.value) === selectedValue
  )
  const selectedLabel =
    selectedValue === emptyValue
      ? placeholder
      : selectedOption?.label || placeholder
  const canClear = allowClear && selectedValue !== emptyValue
  const normalizedSearch = search.trim().toLowerCase()
  const filteredOptions = normalizedSearch
    ? options.filter((option) =>
        `${option.label} ${option.value}`
          .toLowerCase()
          .includes(normalizedSearch)
      )
    : options

  function selectValue(nextValue: string) {
    onValueChange(nextValue === emptyValue ? undefined : nextValue)
    setOpen(false)
    setSearch("")
  }

  function clearValue(event: React.MouseEvent<HTMLElement>) {
    event.preventDefault()
    event.stopPropagation()
    selectValue(emptyValue)
  }

  return (
    <PopoverPrimitive.Root open={open} onOpenChange={setOpen}>
      <PopoverPrimitive.Trigger asChild>
        <Button
          type="button"
          variant="outline"
          role="combobox"
          aria-expanded={open}
          disabled={disabled}
          className={cn(
            "w-full justify-between px-3 font-normal",
            (!selectedOption || selectedValue === emptyValue) &&
              "text-muted-foreground",
            triggerClassName,
            className
          )}
        >
          <span className="line-clamp-1 text-left">
            {selectedLabel || placeholder}
          </span>
          {canClear ? (
            <span
              aria-label="Clear selection"
              className="flex size-4 shrink-0 items-center justify-center rounded-sm opacity-50 transition-opacity hover:opacity-100"
              onClick={clearValue}
              onPointerDown={(event) => {
                event.preventDefault()
                event.stopPropagation()
              }}
            >
              <XIcon className="size-4" />
            </span>
          ) : (
            <ChevronsUpDownIcon className="size-4 shrink-0 opacity-50" />
          )}
        </Button>
      </PopoverPrimitive.Trigger>
      <PopoverPrimitive.Portal>
        <PopoverPrimitive.Content
          align="start"
          sideOffset={4}
          className={cn(
            "z-50 w-[var(--radix-popover-trigger-width)] rounded-md border bg-popover p-1 text-popover-foreground shadow-md outline-none data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
            contentClassName
          )}
        >
          <div className="flex items-center gap-2 border-b px-2 py-1.5">
            <SearchIcon className="size-4 shrink-0 text-muted-foreground" />
            <Input
              autoFocus
              className="h-8 border-0 px-0 shadow-none focus-visible:ring-0"
              value={search}
              placeholder={searchPlaceholder || placeholder}
              onChange={(event) => setSearch(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Escape") setOpen(false)
              }}
            />
          </div>
          <div className="max-h-64 overflow-y-auto py-1">
            {allowClear && emptyLabel ? (
              <ComboboxOptionButton
                label={emptyLabel}
                selected={selectedValue === emptyValue}
                onSelect={() => selectValue(emptyValue)}
              />
            ) : null}
            {filteredOptions.length ? (
              filteredOptions.map((option) => {
                const optionValue = String(option.value)
                return (
                  <ComboboxOptionButton
                    key={optionValue}
                    label={option.label}
                    selected={selectedValue === optionValue}
                    onSelect={() => selectValue(optionValue)}
                  />
                )
              })
            ) : (
              <div className="px-2 py-3 text-sm text-muted-foreground">
                {noResultsLabel}
              </div>
            )}
          </div>
        </PopoverPrimitive.Content>
      </PopoverPrimitive.Portal>
    </PopoverPrimitive.Root>
  )
}

export function DashboardMultiSelect({
  value,
  options,
  placeholder,
  searchPlaceholder,
  noResultsLabel = "-",
  disabled,
  className,
  triggerClassName,
  contentClassName,
  onValueChange,
}: {
  value?: unknown
  options: DashboardSelectOption[]
  placeholder?: string
  searchPlaceholder?: string
  noResultsLabel?: string
  disabled?: boolean
  className?: string
  triggerClassName?: string
  contentClassName?: string
  onValueChange: (value: string[]) => void
}) {
  const [open, setOpen] = React.useState(false)
  const [search, setSearch] = React.useState("")
  const selectedValues = React.useMemo(
    () => (Array.isArray(value) ? value.map((item) => String(item)) : []),
    [value]
  )
  const selectedOptions = selectedValues
    .map((selectedValue) =>
      options.find((option) => String(option.value) === selectedValue)
    )
    .filter((option): option is DashboardSelectOption => Boolean(option))
  const selectedLabel = selectedOptions.length
    ? selectedOptions.map((option) => option.label).join(", ")
    : placeholder
  const normalizedSearch = search.trim().toLowerCase()
  const filteredOptions = normalizedSearch
    ? options.filter((option) =>
        `${option.label} ${option.value}`
          .toLowerCase()
          .includes(normalizedSearch)
      )
    : options

  function toggleValue(nextValue: string) {
    const nextValues = selectedValues.includes(nextValue)
      ? selectedValues.filter((item) => item !== nextValue)
      : [...selectedValues, nextValue]
    onValueChange(nextValues)
  }

  function clearValue(event: React.MouseEvent<HTMLElement>) {
    event.preventDefault()
    event.stopPropagation()
    onValueChange([])
  }

  return (
    <PopoverPrimitive.Root open={open} onOpenChange={setOpen}>
      <PopoverPrimitive.Trigger asChild>
        <Button
          type="button"
          variant="outline"
          role="combobox"
          aria-expanded={open}
          disabled={disabled}
          className={cn(
            "w-full justify-between px-3 font-normal",
            !selectedOptions.length && "text-muted-foreground",
            triggerClassName,
            className
          )}
        >
          <span className="line-clamp-1 text-left">
            {selectedLabel || placeholder}
          </span>
          {selectedOptions.length ? (
            <span
              aria-label="Clear selection"
              className="flex size-4 shrink-0 items-center justify-center rounded-sm opacity-50 transition-opacity hover:opacity-100"
              onClick={clearValue}
              onPointerDown={(event) => {
                event.preventDefault()
                event.stopPropagation()
              }}
            >
              <XIcon className="size-4" />
            </span>
          ) : (
            <ChevronsUpDownIcon className="size-4 shrink-0 opacity-50" />
          )}
        </Button>
      </PopoverPrimitive.Trigger>
      <PopoverPrimitive.Portal>
        <PopoverPrimitive.Content
          align="start"
          sideOffset={4}
          className={cn(
            "z-50 w-[var(--radix-popover-trigger-width)] rounded-md border bg-popover p-1 text-popover-foreground shadow-md outline-none data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
            contentClassName
          )}
        >
          <div className="flex items-center gap-2 border-b px-2 py-1.5">
            <SearchIcon className="size-4 shrink-0 text-muted-foreground" />
            <Input
              autoFocus
              className="h-8 border-0 px-0 shadow-none focus-visible:ring-0"
              value={search}
              placeholder={searchPlaceholder || placeholder}
              onChange={(event) => setSearch(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Escape") setOpen(false)
              }}
            />
          </div>
          <div className="max-h-64 overflow-y-auto py-1">
            {filteredOptions.length ? (
              filteredOptions.map((option) => {
                const optionValue = String(option.value)
                return (
                  <ComboboxOptionButton
                    key={optionValue}
                    label={option.label}
                    selected={selectedValues.includes(optionValue)}
                    onSelect={() => toggleValue(optionValue)}
                  />
                )
              })
            ) : (
              <div className="px-2 py-3 text-sm text-muted-foreground">
                {noResultsLabel}
              </div>
            )}
          </div>
        </PopoverPrimitive.Content>
      </PopoverPrimitive.Portal>
    </PopoverPrimitive.Root>
  )
}

function ComboboxOptionButton({
  label,
  selected,
  onSelect,
}: {
  label: string
  selected: boolean
  onSelect: () => void
}) {
  return (
    <button
      type="button"
      className="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-left text-sm outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"
      onClick={onSelect}
    >
      <CheckIcon
        className={cn("size-4", selected ? "opacity-100" : "opacity-0")}
      />
      <span className="line-clamp-1">{label}</span>
    </button>
  )
}
