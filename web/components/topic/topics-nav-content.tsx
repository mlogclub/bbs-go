"use client"

import Link from "@/components/common/link"
import * as React from "react"

import { ScrollArea } from "@/components/ui/scroll-area"
import { apiFetch } from "@/lib/api/client"
import type { Category } from "@/lib/api/types"
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
  const [categories, setCategories] = React.useState(initialCategories)

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

  return (
    <div className="topics-nav">
      <nav className="dock-nav">
        <ScrollArea className="topics-scroll-area">
          <ul>
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
        </ScrollArea>
      </nav>
    </div>
  )
}
