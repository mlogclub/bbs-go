"use client"

import * as React from "react"
import { Clock, Search, X } from "lucide-react"
import { useSearchParams } from "@/lib/router/navigation"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

const localStorageKey = "bbsgo.search.histories"
const maxHistoryLen = 10

export function SearchInput({
  className,
  placeholder = "",
}: {
  className?: string
  placeholder?: string
}) {
  return (
    <React.Suspense
      fallback={
        <Input
          className={cn("h-9 rounded-lg", className)}
          placeholder={placeholder}
          disabled
        />
      }
    >
      <SearchInputContent className={className} placeholder={placeholder} />
    </React.Suspense>
  )
}

function SearchInputContent({
  className,
  placeholder = "",
}: {
  className?: string
  placeholder?: string
}) {
  const searchParams = useSearchParams()
  const [keyword, setKeyword] = React.useState(searchParams.get("q") || "")
  const [inputFocus, setInputFocus] = React.useState(false)
  const [selectedIndex, setSelectedIndex] = React.useState(-1)
  const [allHistories, setAllHistories] = React.useState<string[]>(() => {
    if (typeof window === "undefined") return []
    try {
      const parsed = JSON.parse(localStorage.getItem(localStorageKey) || "[]")
      return Array.isArray(parsed) ? parsed : []
    } catch {
      return []
    }
  })

  const histories = React.useMemo(() => {
    if (!keyword) return allHistories
    return allHistories.filter((history) => history.includes(keyword))
  }, [allHistories, keyword])

  const showHistories = inputFocus && histories.length > 0

  function addHistories(query: string) {
    const next = [
      query,
      ...allHistories.filter((item) => item !== query),
    ].slice(0, maxHistoryLen)
    localStorage.setItem(localStorageKey, JSON.stringify(next))
    setAllHistories(next)
  }

  function submitSearch(query = keyword) {
    const next = query.trim()
    if (!next) return
    addHistories(next)
    window.location.assign(`/search?q=${encodeURIComponent(next)}`)
  }

  function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    submitSearch()
  }

  function changeSelect(delta: number) {
    if (!histories.length) return
    setSelectedIndex((current) => {
      const next = current + delta
      if (next < 0) return -1
      if (next >= histories.length) return 0
      return next
    })
  }

  function onKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === "ArrowDown") {
      event.preventDefault()
      changeSelect(1)
      return
    }
    if (event.key === "ArrowUp") {
      event.preventDefault()
      changeSelect(-1)
      return
    }
    if (event.key === "Enter") {
      event.stopPropagation()
      event.preventDefault()
      if (selectedIndex >= 0 && histories[selectedIndex]) {
        setKeyword(histories[selectedIndex])
        submitSearch(histories[selectedIndex])
        return
      }
      submitSearch()
    }
  }

  function deleteHistory(item: string) {
    const next = allHistories.filter((history) => history !== item)
    localStorage.setItem(localStorageKey, JSON.stringify(next))
    setAllHistories(next)
  }

  function historyItemClick(item: string) {
    setKeyword(item)
    submitSearch(item)
  }

  return (
    <div className={cn("relative", className)}>
      <form
        action="/search"
        method="GET"
        onSubmit={onSubmit}
        className={cn(
          "flex h-8 w-48 items-center gap-1 rounded-sm border border-border bg-muted px-3 py-1 text-sm transition-all duration-300 ease-in-out focus-within:w-72 focus-within:border-ring focus-within:bg-background focus-within:shadow-sm focus-within:shadow-ring/20",
          inputFocus &&
            "w-72 border-ring bg-background shadow-sm shadow-ring/20"
        )}
      >
        <Input
          type="text"
          name="q"
          value={keyword}
          maxLength={30}
          autoComplete="off"
          onFocus={() => setInputFocus(true)}
          onBlur={() => window.setTimeout(() => setInputFocus(false), 200)}
          onChange={(event) => {
            setKeyword(event.target.value)
            setSelectedIndex(-1)
          }}
          onKeyDown={onKeyDown}
          placeholder={placeholder}
          className="h-9 flex-1 border-none bg-transparent p-0 text-xs shadow-none focus-visible:ring-0 focus-visible:ring-offset-0"
        />
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className={cn(
            "h-auto p-1 text-muted-foreground hover:bg-transparent hover:text-foreground",
            inputFocus && "text-primary"
          )}
          onClick={() => submitSearch()}
        >
          <Search className="h-4 w-4" />
        </Button>
      </form>

      {showHistories ? (
        <div className="absolute top-full left-0 z-50 mt-1 w-72 animate-in duration-200 slide-in-from-top-2">
          <div className="overflow-hidden rounded-lg border border-border bg-popover shadow-lg">
            {histories.map((item, index) => (
              <div
                key={item}
                className={cn(
                  "group flex cursor-pointer items-center justify-between border-b border-border px-3 py-2 transition-colors last:border-b-0 hover:bg-accent",
                  index === selectedIndex && "bg-accent"
                )}
                onMouseOver={() => setSelectedIndex(index)}
                onMouseOut={() => setSelectedIndex(-1)}
              >
                <button
                  type="button"
                  className="flex flex-1 cursor-pointer items-center text-left text-xs text-foreground"
                  onClick={() => historyItemClick(item)}
                >
                  <Clock className="mr-2 h-3 w-3 opacity-50" />
                  {item}
                </button>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="h-auto bg-transparent p-1 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 hover:bg-transparent hover:text-destructive"
                  onClick={() => deleteHistory(item)}
                >
                  <X className="h-3 w-3" />
                </Button>
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  )
}
