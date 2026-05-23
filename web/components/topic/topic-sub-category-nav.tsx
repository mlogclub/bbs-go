"use client"

import Link from "@/components/common/link"
import * as React from "react"
import { ChevronDown, ChevronUp } from "lucide-react"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import type { Category } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { cn } from "@/lib/utils"

function categorySelected(
  rootCategoryId: number,
  children: Category[],
  currentCategoryId: number
) {
  if (currentCategoryId === rootCategoryId) {
    return rootCategoryId
  }

  return children.some((item) => Number(item.id) === currentCategoryId)
    ? currentCategoryId
    : rootCategoryId
}

function SubCategoryLink({
  id,
  selectedCategoryId,
  children,
}: {
  id: number
  selectedCategoryId: number
  children: React.ReactNode
}) {
  return (
    <Link
      href={`/topics/category/${id}`}
      data-node-id={id}
      className={cn(
        "inline-flex shrink-0 items-center rounded-md px-3 py-1 font-medium whitespace-nowrap transition-colors",
        selectedCategoryId === id
          ? "bg-primary text-primary-foreground shadow-sm"
          : "bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground"
      )}
    >
      {children}
    </Link>
  )
}

export function TopicSubCategoryNav({
  rootCategoryId,
  categories,
  currentCategoryId,
}: {
  rootCategoryId: number
  categories?: Category[]
  currentCategoryId: number
}) {
  const { t } = useI18n()
  const [open, setOpen] = React.useState(false)
  const rootRef = React.useRef<HTMLDivElement>(null)
  const listRef = React.useRef<HTMLDivElement>(null)
  const navChildren = Array.isArray(categories) ? categories : []
  const selectedCategoryId = categorySelected(rootCategoryId, navChildren, currentCategoryId)

  const getScrollViewport = React.useCallback(() => {
    return (
      rootRef.current?.querySelector<HTMLElement>(
        '[data-slot="scroll-area-viewport"]'
      ) || null
    )
  }, [])

  const scrollSelectedIntoView = React.useCallback(
    (behavior: ScrollBehavior = "smooth") => {
      window.requestAnimationFrame(() => {
        const activeNode = listRef.current?.querySelector<HTMLElement>(
          `[data-node-id="${selectedCategoryId}"]`
        )
        activeNode?.scrollIntoView({
          behavior,
          block: "nearest",
          inline: "nearest",
        })
      })
    },
    [selectedCategoryId]
  )

  React.useEffect(() => {
    scrollSelectedIntoView("auto")
  }, [navChildren.length, scrollSelectedIntoView])

  React.useEffect(() => {
    scrollSelectedIntoView()
  }, [selectedCategoryId, scrollSelectedIntoView])

  React.useEffect(() => {
    const viewport = getScrollViewport()
    if (!viewport) {
      return
    }

    function onWheel(event: WheelEvent) {
      const currentViewport = getScrollViewport()
      if (
        !currentViewport ||
        currentViewport.scrollWidth <= currentViewport.clientWidth
      ) {
        return
      }

      if (Math.abs(event.deltaY) > Math.abs(event.deltaX)) {
        event.preventDefault()
        currentViewport.scrollLeft += event.deltaY
      }
    }

    viewport.addEventListener("wheel", onWheel, { passive: false })
    return () => viewport.removeEventListener("wheel", onWheel)
  }, [getScrollViewport])

  React.useEffect(() => {
    if (!open) {
      return
    }

    window.requestAnimationFrame(() => {
      const menuContent = Array.from(
        document.querySelectorAll<HTMLElement>(
          '[data-slot="dropdown-menu-content"][data-state="open"]'
        )
      ).at(-1)
      const selectedItem = menuContent?.querySelector<HTMLElement>(
        "[data-selected-node=true]"
      )
      if (!menuContent || !selectedItem) {
        return
      }

      const targetTop =
        selectedItem.offsetTop -
        (menuContent.clientHeight - selectedItem.clientHeight) / 2
      menuContent.scrollTo({
        top: Math.max(0, targetTop),
        behavior: "smooth",
      })
    })
  }, [open, selectedCategoryId])

  if (!navChildren.length) {
    return null
  }

  return (
    <div className="border-b border-border px-4 pt-3 pb-2">
      <div className="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-2 text-xs sm:text-[13px]">
        <div ref={rootRef} className="min-w-0">
          <ScrollArea className="topic-sub-node-scroll min-w-0 whitespace-nowrap [&_[data-slot=scroll-area-viewport]]:overflow-y-hidden">
            <div
              ref={listRef}
              className="flex w-max flex-nowrap items-center gap-1 pr-1 pb-2"
            >
              <SubCategoryLink id={rootCategoryId} selectedCategoryId={selectedCategoryId}>
                {t("pages.topics.allCategories")}
              </SubCategoryLink>
              {navChildren.map((child) => (
                <SubCategoryLink
                  key={child.id}
                  id={child.id}
                  selectedCategoryId={selectedCategoryId}
                >
                  {child.name}
                </SubCategoryLink>
              ))}
            </div>
            <ScrollBar orientation="horizontal" />
          </ScrollArea>
        </div>
        <div className="pb-1.5">
          <DropdownMenu modal={false} open={open} onOpenChange={setOpen}>
            <DropdownMenuTrigger asChild>
              <button
                type="button"
                className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                aria-label={t("pages.topics.moreSubCategories")}
                title={t("pages.topics.moreSubCategories")}
              >
                {open ? (
                  <ChevronUp className="h-4 w-4" />
                ) : (
                  <ChevronDown className="h-4 w-4" />
                )}
                <span className="sr-only">
                  {t("pages.topics.moreSubCategories")}
                </span>
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              align="end"
              className="max-h-[min(60vh,420px)] w-[280px] sm:w-[320px] md:w-[360px]"
            >
              <DropdownMenuItem
                asChild
                data-selected-node={selectedCategoryId === rootCategoryId}
                className={cn(
                  selectedCategoryId === rootCategoryId &&
                    "bg-accent text-accent-foreground"
                )}
              >
                <Link href={`/topics/category/${rootCategoryId}`}>
                  {t("pages.topics.allCategories")}
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {navChildren.map((child) => (
                <DropdownMenuItem
                  key={`menu-${child.id}`}
                  asChild
                  data-selected-node={selectedCategoryId === child.id}
                  className={cn(
                    selectedCategoryId === child.id &&
                      "bg-accent text-accent-foreground"
                  )}
                >
                  <Link href={`/topics/category/${child.id}`} className="truncate">
                    {child.name}
                  </Link>
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </div>
  )
}
