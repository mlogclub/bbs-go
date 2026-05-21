"use client"

import Link from "@/components/common/link"
import * as React from "react"

import { ScrollArea } from "@/components/ui/scroll-area"
import { apiFetch } from "@/lib/api/client"
import type { TopicNode } from "@/lib/api/types"
import { cn } from "@/lib/utils"

function isBuiltInNode(node: TopicNode) {
  return node.id <= 0
}

function nodeHref(node: TopicNode) {
  if (node.id > 0) {
    return `/topics/node/${node.id}`
  }
  if (node.id === 0) {
    return "/topics/node/newest"
  }
  if (node.id === -1) {
    return "/topics/node/recommend"
  }
  return "/topics/node/feed"
}

function isActiveNode(
  node: TopicNode,
  currentNodeId?: number,
  currentRootNodeId?: number
) {
  if (isBuiltInNode(node)) {
    return currentNodeId === node.id
  }
  return currentRootNodeId === node.id
}

export function TopicsNavContent({
  initialNodes,
  currentNodeId,
  currentRootNodeId,
}: {
  initialNodes: TopicNode[]
  currentNodeId?: number
  currentRootNodeId?: number
}) {
  const [nodes, setNodes] = React.useState(initialNodes)

  React.useEffect(() => {
    if (initialNodes.length > 0) return

    let mounted = true
    const timer = window.setTimeout(() => {
      void apiFetch<TopicNode[]>("/api/topic/node_navs")
        .then((data) => {
          if (mounted) {
            setNodes(data)
          }
        })
        .catch(() => undefined)
    }, 0)

    return () => {
      mounted = false
      window.clearTimeout(timer)
    }
  }, [initialNodes.length])

  return (
    <div className="topics-nav">
      <nav className="dock-nav">
        <ScrollArea className="topics-scroll-area">
          <ul>
            {nodes.map((node, index) => {
              const previousNode = nodes[index - 1]
              const showDivider =
                index > 0 &&
                previousNode &&
                isBuiltInNode(previousNode) &&
                !isBuiltInNode(node)
              const active = isActiveNode(
                node,
                currentNodeId,
                currentRootNodeId
              )

              return (
                <React.Fragment key={node.id}>
                  {showDivider ? (
                    <li className="nodes-divider" aria-hidden="true" />
                  ) : null}
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
