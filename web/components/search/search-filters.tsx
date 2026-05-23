"use client"

import * as React from "react"
import { useRouter, useSearchParams } from "@/lib/router/navigation"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useI18n } from "@/lib/i18n/provider"

type NodeOption = {
  id: number
  label: string
}

export function SearchFilters({
  categories,
  showNode = true,
  showTime = true,
}: {
  categories: NodeOption[]
  showNode?: boolean
  showTime?: boolean
}) {
  return (
    <React.Suspense fallback={<div className="h-9 w-[268px]" />}>
      <SearchFiltersContent
        categories={categories}
        showNode={showNode}
        showTime={showTime}
      />
    </React.Suspense>
  )
}

function SearchFiltersContent({
  categories,
  showNode,
  showTime,
}: {
  categories: NodeOption[]
  showNode: boolean
  showTime: boolean
}) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { t } = useI18n()
  const categoryId = searchParams.get("categoryId") || "0"
  const timeRange = searchParams.get("timeRange") || "0"

  function setQuery(key: string, value: string) {
    const params = new URLSearchParams(searchParams.toString())
    if (Number(value) === 0) {
      params.delete(key)
    } else {
      params.set(key, value)
    }
    const query = params.toString()
    router.push(`/search${query ? `?${query}` : ""}`)
  }

  return (
    <div className="flex flex-wrap items-center gap-2">
      {showNode ? (
        <Select
          value={categoryId}
          onValueChange={(value) => setQuery("categoryId", value)}
        >
          <SelectTrigger
            className="h-9 w-full rounded-lg bg-background text-xs sm:w-[168px]"
            aria-label={t("component.search.allCategories")}
          >
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="0">{t("component.search.allCategories")}</SelectItem>
            {categories.map((category) => (
              <SelectItem key={category.id} value={String(category.id)}>
                {category.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ) : null}
      {showTime ? (
        <Select
          value={timeRange}
          onValueChange={(value) => setQuery("timeRange", value)}
        >
          <SelectTrigger
            className="h-9 w-[132px] rounded-lg bg-background text-xs"
            aria-label={t("component.search.timeRange.all")}
          >
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="0">
              {t("component.search.timeRange.all")}
            </SelectItem>
            <SelectItem value="1">
              {t("component.search.timeRange.day")}
            </SelectItem>
            <SelectItem value="2">
              {t("component.search.timeRange.week")}
            </SelectItem>
            <SelectItem value="3">
              {t("component.search.timeRange.month")}
            </SelectItem>
            <SelectItem value="4">
              {t("component.search.timeRange.year")}
            </SelectItem>
          </SelectContent>
        </Select>
      ) : null}
    </div>
  )
}
