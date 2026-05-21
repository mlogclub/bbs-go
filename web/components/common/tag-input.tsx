"use client"

import * as React from "react"
import { X } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { apiFetch, toFormData } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"
import { cn } from "@/lib/utils"

export function TagInput({
  value,
  recommendTags,
  placeholder,
  onChange,
}: {
  value: string[]
  recommendTags?: string[]
  placeholder: string
  onChange: (value: string[]) => void
}) {
  const { t } = useI18n()
  const [input, setInput] = React.useState("")
  const [autocompleteTags, setAutocompleteTags] = React.useState<string[]>([])
  const [showRecommendTags, setShowRecommendTags] = React.useState(false)
  const [selectIndex, setSelectIndex] = React.useState(-1)
  const inputRef = React.useRef<HTMLInputElement>(null)
  const closeRecommendTimerRef = React.useRef<number | null>(null)
  const maxTagCount = 3
  const maxWordCount = 15
  const showAutocompleteTags = autocompleteTags.length > 0
  const showRecommendedOptions = !showAutocompleteTags && showRecommendTags && !!recommendTags?.length

  function addTagName(tagName: string) {
    const normalized = tagName.trim().replace(/^[,;]+|[,;]+$/g, "")
    if (!normalized || value.length >= maxTagCount || normalized.length > maxWordCount || value.includes(normalized)) {
      return
    }
    onChange([...value, normalized])
    setInput("")
    setAutocompleteTags([])
    setSelectIndex(-1)
  }

  async function autocomplete(nextInput: string) {
    setInput(nextInput)
    setShowRecommendTags(false)
    setSelectIndex(-1)

    if (!nextInput) {
      setAutocompleteTags([])
      return
    }

    try {
      const ret = await apiFetch<Array<{ name: string }>>("/api/tag/autocomplete", {
        method: "POST",
        body: toFormData({ input: nextInput }),
      })
      setAutocompleteTags((ret || []).map((item) => item.name).filter(Boolean))
    } catch {
      setAutocompleteTags([])
    }
  }

  function addTag(event: React.KeyboardEvent<HTMLInputElement>) {
    const hasSelectSuggestion = selectIndex >= 0 && autocompleteTags.length > selectIndex
    if (!hasSelectSuggestion && event.key !== "Enter" && event.key !== "," && event.key !== ";") {
      return
    }

    event.preventDefault()
    event.stopPropagation()
    addTagName(hasSelectSuggestion ? autocompleteTags[selectIndex] : event.currentTarget.value || input)
  }

  function removeLastTag(event: React.KeyboardEvent<HTMLInputElement>) {
    if ((event.key !== "Backspace" && event.key !== "Delete") || event.currentTarget.value || value.length === 0) {
      return false
    }

    event.preventDefault()
    event.stopPropagation()
    onChange(value.slice(0, -1))
    setAutocompleteTags([])
    setSelectIndex(-1)
    return true
  }

  function openRecommendTags() {
    if (closeRecommendTimerRef.current) {
      window.clearTimeout(closeRecommendTimerRef.current)
      closeRecommendTimerRef.current = null
    }
    if (recommendTags?.length) {
      setShowRecommendTags(true)
    }
  }

  function closeRecommendTags() {
    closeRecommendTimerRef.current = window.setTimeout(() => {
      setShowRecommendTags(false)
    }, 100)
  }

  React.useEffect(() => {
    return () => {
      if (closeRecommendTimerRef.current) {
        window.clearTimeout(closeRecommendTimerRef.current)
      }
    }
  }, [])

  return (
    <div className="relative">
      <input id="tags" name="tags" type="hidden" value={value.join(",")} readOnly />
      <div
        className="flex min-h-10 flex-wrap items-center gap-1.5 rounded-md border border-input bg-transparent px-2.5 py-1.5 text-sm shadow-xs transition-[color,box-shadow] focus-within:border-ring focus-within:ring-3 focus-within:ring-ring/50 dark:bg-input/30"
        onMouseDown={(event) => {
          if (event.target === event.currentTarget) {
            event.preventDefault()
            inputRef.current?.focus()
          }
        }}
      >
        {value.map((tag) => (
          <Badge key={tag} variant="secondary" className="h-6 rounded-sm pr-1">
            <span>{tag}</span>
            <Button
              type="button"
              variant="ghost"
              size="icon-xs"
              className="-mr-1 h-5 w-5 opacity-60 hover:bg-transparent hover:text-destructive hover:opacity-100"
              onClick={() => onChange(value.filter((item) => item !== tag))}
            >
              <X className="h-3 w-3" />
            </Button>
          </Badge>
        ))}
        <input
          ref={inputRef}
          value={input}
          placeholder={placeholder}
          maxLength={maxWordCount}
          className="min-w-24 flex-1 !border-0 bg-transparent py-1 text-sm !shadow-none outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
          onChange={(event) => void autocomplete(event.currentTarget.value)}
          onKeyDown={(event) => {
            if (removeLastTag(event)) {
              return
            }
            if (event.key === "ArrowUp") {
              event.preventDefault()
              event.stopPropagation()
              setSelectIndex((current) => (current < 0 ? current : current - 1))
              return
            }
            if (event.key === "ArrowDown") {
              event.preventDefault()
              event.stopPropagation()
              setSelectIndex((current) => (current < autocompleteTags.length - 1 ? current + 1 : current))
              return
            }
            if (event.key === "Escape") {
              setAutocompleteTags([])
              setSelectIndex(-1)
              closeRecommendTags()
              return
            }
            addTag(event)
          }}
          onFocus={openRecommendTags}
          onClick={openRecommendTags}
          onBlur={closeRecommendTags}
        />
      </div>
      {showAutocompleteTags ? (
        <div className="absolute top-full right-0 left-0 z-50 mt-1 overflow-hidden rounded-md bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10">
          <div className="max-h-72 overflow-y-auto p-1">
            {autocompleteTags.map((tag, index) => (
              <button
                key={tag}
                type="button"
                className={cn(
                  "flex w-full cursor-default items-center rounded-sm px-2 py-1.5 text-left text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground",
                  index === selectIndex && "bg-accent text-accent-foreground"
                )}
                onMouseDown={(event) => event.preventDefault()}
                onMouseEnter={() => setSelectIndex(index)}
                onClick={() => addTagName(tag)}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      ) : null}
      {showRecommendedOptions ? (
        <div className="absolute top-full right-0 left-0 z-50 mt-1 rounded-md bg-popover p-3 text-sm text-popover-foreground shadow-md ring-1 ring-foreground/10">
          <div className="mb-2 flex items-center justify-between border-b border-border pb-2">
            <span className="font-medium text-primary">{t("component.tagInput.recommendTags")}</span>
            <Button
              type="button"
              variant="ghost"
              size="icon-xs"
              className="text-muted-foreground hover:text-destructive"
              onMouseDown={(event) => event.preventDefault()}
              onClick={closeRecommendTags}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
          <div className="flex flex-wrap gap-1.5">
            {recommendTags.map((tag) => (
              <Button
                key={tag}
                type="button"
                variant="secondary"
                size="xs"
                className="h-6 font-normal text-primary hover:bg-primary hover:text-primary-foreground"
                onMouseDown={(event) => event.preventDefault()}
                onClick={() => {
                  addTagName(tag)
                  closeRecommendTags()
                }}
              >
                {tag}
              </Button>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  )
}
