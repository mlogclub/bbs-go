"use client"

import { useRouter } from "@/lib/router/navigation"

import { cn } from "@/lib/utils"

export type NodeTopicFilterType = "qa" | "normal"

export type NodeTopicFilterItem = {
  value: string
  label: string
}

function buildFilterPath(nodeId: number, filterType: NodeTopicFilterType, value: string) {
  const params = new URLSearchParams()

  if (filterType === "qa" && value) {
    params.set("qaStatus", value)
  } else if (filterType === "normal") {
    params.set("sort", value)
  }

  const query = params.toString()
  return `/topics/node/${nodeId}${query ? `?${query}` : ""}`
}

export function NodeTopicFilters({
  nodeId,
  filterType,
  filters,
  currentValue,
  currentLabel,
}: {
  nodeId: number
  filterType: NodeTopicFilterType
  filters: NodeTopicFilterItem[]
  currentValue: string
  currentLabel: string
}) {
  const router = useRouter()

  return (
    <div className="flex justify-between px-4 py-3">
      <div className="text-base font-bold">{currentLabel}</div>
      <div className="inline-flex flex-wrap items-center gap-1 rounded-lg bg-muted p-1">
        {filters.map((item) => (
          <button
            key={item.value}
            type="button"
            className={cn(
              "inline-flex h-5 items-center rounded-md px-3 text-sm font-medium transition-colors",
              currentValue === item.value
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground",
            )}
            onClick={() => router.replace(buildFilterPath(nodeId, filterType, item.value))}
          >
            {item.label}
          </button>
        ))}
      </div>
    </div>
  )
}
