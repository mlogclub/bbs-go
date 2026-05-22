"use client"

import * as React from "react"

import { WidgetCard } from "@/components/common/widget-card"
import type { Topic } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { cn } from "@/lib/utils"

export type TocItem = {
  id: string
  title: string
  level?: number
}

function indentClass(level?: number) {
  if (level === 3) {
    return "pl-5"
  }
  if (level && level >= 4) {
    return "pl-8"
  }
  return "pl-3"
}

function getHashId() {
  if (!window.location.hash) {
    return ""
  }

  try {
    return decodeURIComponent(window.location.hash.slice(1))
  } catch {
    return window.location.hash.slice(1)
  }
}

export function TopicToc({
  items: propItems,
  topic,
}: {
  items?: TocItem[]
  topic?: Pick<Topic, "toc">
}) {
  const { t } = useI18n()
  const items = React.useMemo(
    () =>
      (propItems || topic?.toc || []).filter(
        (item) => item.id && item.title
      ),
    [propItems, topic]
  )
  const [activeId, setActiveId] = React.useState(items[0]?.id || "")
  const tickingRef = React.useRef(false)

  const scrollToHeading = React.useCallback((id: string) => {
    const heading = document.getElementById(id)
    if (!heading) {
      return
    }

    setActiveId(id)
    heading.scrollIntoView({ behavior: "smooth", block: "start" })
    window.history.replaceState(null, "", `#${encodeURIComponent(id)}`)
  }, [])

  const updateActiveFromScroll = React.useCallback(() => {
    tickingRef.current = false
    const headings = items
      .map((item) => document.getElementById(item.id))
      .filter((item): item is HTMLElement => Boolean(item))

    if (!headings.length) {
      setActiveId("")
      return
    }

    let current = headings[0]
    for (const heading of headings) {
      if (heading.getBoundingClientRect().top <= 90) {
        current = heading
      } else {
        break
      }
    }
    setActiveId(current.id)
  }, [items])

  const scrollToHash = React.useCallback(() => {
    const id = getHashId()
    if (!id || !items.some((item) => item.id === id)) {
      return
    }

    window.setTimeout(() => scrollToHeading(id), 0)
  }, [items, scrollToHeading])

  React.useEffect(() => {
    const handleScroll = () => {
      if (tickingRef.current) {
        return
      }
      tickingRef.current = true
      window.requestAnimationFrame(updateActiveFromScroll)
    }

    window.setTimeout(() => {
      scrollToHash()
      updateActiveFromScroll()
    }, 0)
    window.addEventListener("hashchange", scrollToHash)
    window.addEventListener("scroll", handleScroll, { passive: true })

    return () => {
      window.removeEventListener("hashchange", scrollToHash)
      window.removeEventListener("scroll", handleScroll)
    }
  }, [scrollToHash, updateActiveFromScroll])

  if (!items.length) {
    return null
  }

  return (
    <div className="sticky top-18">
      <WidgetCard title={t("pages.topic.detail.tocTitle")}>
        <nav className="max-h-[calc(100vh-7rem)] overflow-y-auto py-1">
          {items.map((item) => (
            <button
              key={item.id}
              type="button"
              className={cn(
                "block w-full truncate border-l-2 py-1.5 pr-2 text-left text-sm transition-colors hover:border-primary hover:text-primary",
                activeId === item.id
                  ? "border-primary font-medium text-primary"
                  : "border-transparent text-muted-foreground",
                indentClass(item.level)
              )}
              title={item.title}
              onClick={() => scrollToHeading(item.id)}
            >
              {item.title}
            </button>
          ))}
        </nav>
      </WidgetCard>
    </div>
  )
}
