"use client"

import Link from "@/components/common/link"
import * as React from "react"
import { LayoutGridIcon, MoreHorizontalIcon } from "lucide-react"

import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { ScrollArea } from "@/components/ui/scroll-area"
import { apiFetch } from "@/lib/api/client"
import type { Category } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { cn } from "@/lib/utils"

function nodeHref(node: Category) {
  return `/topics/category/${node.id}`
}

function isActiveNode(
  node: Category,
  currentCategoryId?: number,
  currentRootCategoryId?: number
) {
  return currentRootCategoryId === node.id
}

export function TopicsNavContent({
  initialCategories,
  currentCategoryId,
  currentRootCategoryId,
}: {
  initialCategories: Category[]
  currentCategoryId?: number
  currentRootCategoryId?: number
}) {
  const { t } = useI18n()
  const [categories, setCategories] = React.useState(initialCategories)
  const [mobileDrawerOpen, setMobileDrawerOpen] = React.useState(false)
  const mobileScrollRef = React.useRef<HTMLDivElement>(null)

  React.useEffect(() => {
    if (initialCategories.length > 0) return

    let mounted = true
    const timer = window.setTimeout(() => {
      void apiFetch<Category[]>("/api/topic/category_navs")
        .then((data) => {
          if (mounted) {
            setCategories(data)
          }
        })
        .catch(() => undefined)
    }, 0)

    return () => {
      mounted = false
      window.clearTimeout(timer)
    }
  }, [initialCategories.length])

  const visibleCategories = categories.filter((node) => node.id > 0)
  const allCategoryLabel = t("pages.topics.allCategories")
  const moreCategoriesLabel = t("pages.topics.moreCategories")
  const activeNodeId =
    currentCategoryId !== undefined && currentCategoryId <= 0
      ? "all"
      : String(currentRootCategoryId || currentCategoryId || "")

  React.useEffect(() => {
    if (!activeNodeId) return

    window.requestAnimationFrame(() => {
      const activeNode = mobileScrollRef.current?.querySelector<HTMLElement>(
        `[data-node-id="${activeNodeId}"]`
      )
      activeNode?.scrollIntoView({
        behavior: "smooth",
        block: "nearest",
        inline: "center",
      })
    })
  }, [activeNodeId, visibleCategories.length])

  function renderCategoryList() {
    return (
      <ul className="dock-nav-list">
        <li
          className={cn(
            currentCategoryId !== undefined &&
              currentCategoryId <= 0 &&
              "active"
          )}
          data-node-id="all"
        >
          <Link href="/topics">
            <LayoutGridIcon
              className="node-logo node-logo-icon"
              aria-hidden="true"
            />
            <div className="node-name">{allCategoryLabel}</div>
          </Link>
        </li>
        {/* {visibleCategories.length > 0 ? (
        <li className="categories-divider" aria-hidden="true" />
      ) : null} */}
        {visibleCategories.map((node) => {
          const active = isActiveNode(
            node,
            currentCategoryId,
            currentRootCategoryId
          )

          return (
            <React.Fragment key={node.id}>
              <li className={cn(active && "active")} data-node-id={node.id}>
                <Link href={nodeHref(node)}>
                  <i
                    className="node-logo"
                    style={
                      node.logo
                        ? { backgroundImage: `url(${node.logo})` }
                        : undefined
                    }
                  />
                  <div className="node-name">{node.name}</div>
                </Link>
              </li>
            </React.Fragment>
          )
        })}
      </ul>
    )
  }

  return (
    <div className="topics-nav">
      <nav className="dock-nav">
        <ScrollArea className="topics-scroll-area dock-nav-desktop-scroll">
          {renderCategoryList()}
        </ScrollArea>
        <div className="dock-nav-mobile-row">
          <div
            ref={mobileScrollRef}
            className="dock-nav-scroll"
            aria-label={allCategoryLabel}
          >
            {renderCategoryList()}
          </div>
          <Drawer open={mobileDrawerOpen} onOpenChange={setMobileDrawerOpen}>
            <DrawerTrigger asChild>
              <button
                type="button"
                className="topics-mobile-category-more"
                aria-label={moreCategoriesLabel}
                title={moreCategoriesLabel}
              >
                <MoreHorizontalIcon aria-hidden="true" />
                <span>{t("pages.topic.categorySelector.more")}</span>
              </button>
            </DrawerTrigger>
            <DrawerContent className="topics-mobile-category-drawer">
              <DrawerHeader>
                <DrawerTitle>{moreCategoriesLabel}</DrawerTitle>
              </DrawerHeader>
              <div className="topics-mobile-category-list">
                <DrawerClose asChild>
                  <Link
                    href="/topics"
                    className={cn(
                      "topics-mobile-category-item",
                      currentCategoryId !== undefined &&
                        currentCategoryId <= 0 &&
                        "active"
                    )}
                  >
                    <LayoutGridIcon aria-hidden="true" />
                    <span>{allCategoryLabel}</span>
                  </Link>
                </DrawerClose>
                {visibleCategories.map((node) => {
                  const active = isActiveNode(
                    node,
                    currentCategoryId,
                    currentRootCategoryId
                  )

                  return (
                    <DrawerClose key={`mobile-${node.id}`} asChild>
                      <Link
                        href={nodeHref(node)}
                        className={cn(
                          "topics-mobile-category-item",
                          active && "active"
                        )}
                      >
                        {node.logo ? (
                          <i
                            className="node-logo"
                            style={{ backgroundImage: `url(${node.logo})` }}
                          />
                        ) : (
                          <span className="node-logo node-logo-placeholder" />
                        )}
                        <span>{node.name}</span>
                      </Link>
                    </DrawerClose>
                  )
                })}
              </div>
            </DrawerContent>
          </Drawer>
        </div>
      </nav>
    </div>
  )
}
