"use client"

import * as React from "react"
import { Dialog as DialogPrimitive } from "radix-ui"
import { Check, ChevronRight, ChevronsUpDown, X } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import type { TopicNode } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { flattenTopicNodes } from "@/lib/topic-nodes"
import { cn } from "@/lib/utils"

type NodeOption = TopicNode & {
  level: number
  parentName: string
  parentId: number
}

function buildNodeOptions(nodes: TopicNode[], level = 1, parent: TopicNode | null = null): NodeOption[] {
  if (!Array.isArray(nodes) || nodes.length === 0) {
    return []
  }

  return nodes.flatMap((node) => {
    const current = {
      ...node,
      level,
      parentName: parent?.name || "",
      parentId: parent?.id ? Number(parent.id) : 0,
    }
    return [current, ...buildNodeOptions(node.children || [], level + 1, node)]
  })
}

function formatNodeLabel(node?: NodeOption | null) {
  if (!node) {
    return ""
  }
  if (node.level === 2 && node.parentName) {
    return `${node.parentName} / ${node.name}`
  }
  return node.name
}

export function TopicNodeSelector({
  value,
  nodes,
  triggerFullWidth = true,
  triggerSize = "default",
  triggerLabel,
  triggerIcon = "down",
  onChange,
}: {
  value: number
  nodes: TopicNode[]
  triggerFullWidth?: boolean
  triggerSize?: "default" | "sm"
  triggerLabel?: string
  triggerIcon?: "down" | "right"
  onChange: (value: number) => void
}) {
  const { t } = useI18n()
  const [open, setOpen] = React.useState(false)
  const [keyword, setKeyword] = React.useState("")
  const listRef = React.useRef<HTMLDivElement>(null)
  const storageKey = "bbsgo-recent-topic-node-ids"
  const maxRecentCount = 5
  const [recentNodeIds, setRecentNodeIds] = React.useState<number[]>(() => {
    if (typeof window === "undefined") {
      return []
    }
    try {
      const raw = window.localStorage.getItem(storageKey)
      if (!raw) {
        return []
      }
      const parsed = JSON.parse(raw)
      return Array.isArray(parsed) ? parsed.map((id) => Number(id)).filter((id) => Number.isFinite(id)) : []
    } catch {
      return []
    }
  })

  const nodeList = React.useMemo(() => buildNodeOptions(Array.isArray(nodes) ? nodes : []), [nodes])
  const selectedNode = nodeList.find((node) => Number(node.id) === Number(value))
  const triggerText = triggerLabel || formatNodeLabel(selectedNode) || t("pages.topic.nodeSelector.choose")
  const triggerClass = triggerFullWidth ? "h-10 w-full justify-between px-3" : "h-6 justify-between px-3"

  const filteredNodes = React.useMemo(() => {
    const query = keyword.trim().toLowerCase()
    if (!query) {
      return nodeList
    }
    return nodeList.filter((node) => {
      const name = String(node.name || "").toLowerCase()
      const parentName = String(node.parentName || "").toLowerCase()
      return name.includes(query) || parentName.includes(query)
    })
  }, [keyword, nodeList])

  const recentNodes = React.useMemo(() => {
    if (!recentNodeIds.length) {
      return []
    }
    const nodeMap = new Map(nodeList.map((node) => [Number(node.id), node]))
    return recentNodeIds.map((id) => nodeMap.get(Number(id))).filter(Boolean) as NodeOption[]
  }, [nodeList, recentNodeIds])

  React.useEffect(() => {
    if (!open) {
      return
    }
    window.requestAnimationFrame(() => {
      const selectedEl = listRef.current?.querySelector<HTMLElement>('button[data-selected="true"]')
      selectedEl?.scrollIntoView({ behavior: "smooth", block: "center" })
    })
  }, [open])

  function pushRecentNode(id: number) {
    const nextIds = recentNodeIds.filter((item) => Number(item) !== Number(id))
    nextIds.unshift(Number(id))
    const limitedIds = nextIds.slice(0, maxRecentCount)
    setRecentNodeIds(limitedIds)
    window.localStorage.setItem(storageKey, JSON.stringify(limitedIds))
  }

  function selectNode(node: NodeOption) {
    onChange(Number(node.id))
    pushRecentNode(Number(node.id))
    setOpen(false)
  }

  function onOpenChange(visible: boolean) {
    setOpen(visible)
    if (!visible) {
      setKeyword("")
    }
  }

  return (
    <div className="node-selector">
      <Button type="button" variant="outline" size={triggerSize} className={triggerClass} onClick={() => setOpen(true)}>
        <span className="truncate">{triggerText}</span>
        {triggerIcon === "right" ? <ChevronRight className="h-4 w-4 opacity-70" /> : <ChevronsUpDown className="h-4 w-4 opacity-60" />}
      </Button>

      <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}>
        <DialogPrimitive.Portal>
          <DialogPrimitive.Overlay className="fixed inset-0 z-50 bg-black/40 data-[state=closed]:animate-out data-[state=open]:animate-in data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0" />
          <DialogPrimitive.Content className="fixed top-1/2 left-1/2 z-50 grid w-[calc(100%-2rem)] max-w-lg -translate-x-1/2 -translate-y-1/2 gap-4 rounded-lg border bg-background p-6 shadow-lg duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 sm:max-w-lg">
            <div className="flex flex-col space-y-1.5 text-center sm:text-left">
              <DialogPrimitive.Title className="text-lg leading-none font-semibold tracking-tight">
                {t("pages.topic.nodeSelector.title")}
              </DialogPrimitive.Title>
            </div>

            <div className="space-y-3">
              <Input value={keyword} placeholder={t("pages.topic.nodeSelector.searchPlaceholder")} onChange={(event) => setKeyword(event.currentTarget.value)} />

              {recentNodes.length ? (
                <div className="space-y-2">
                  <div className="text-xs text-muted-foreground">{t("pages.topic.nodeSelector.recent")}</div>
                  <div className="flex flex-wrap gap-1.5">
                    {recentNodes.map((node) => (
                      <Button
                        key={`recent-${node.id}`}
                        type="button"
                        size="xs"
                        variant={Number(value) === Number(node.id) ? "default" : "secondary"}
                        className="h-6 px-2 text-xs"
                        onClick={() => selectNode(node)}
                      >
                        {formatNodeLabel(node)}
                      </Button>
                    ))}
                  </div>
                </div>
              ) : null}

              <ScrollArea className="h-72 rounded-md border">
                <div ref={listRef} className="flex flex-col gap-1 p-2">
                  {filteredNodes.map((node) => (
                    <button
                      key={node.id}
                      type="button"
                      data-node-id={node.id}
                      data-selected={Number(value) === Number(node.id)}
                      className={cn(
                        "flex w-full items-center justify-between gap-3 rounded-lg px-3 py-1.5 text-left text-sm transition-colors hover:bg-muted/60",
                        Number(value) === Number(node.id) && "bg-muted text-primary",
                        node.level === 1 && "font-medium text-foreground",
                        node.level === 2 && "pl-7 text-muted-foreground"
                      )}
                      onClick={() => selectNode(node)}
                    >
                      <span className="flex min-w-0 flex-1 items-center">
                        <span className="flex min-w-0 items-center gap-2">
                          <span className={cn("truncate leading-5", node.level === 2 && "text-foreground/90")}>{node.name}</span>
                          {node.level === 2 && node.parentName ? (
                            <span className="truncate text-[11px] leading-4 font-normal text-muted-foreground">/ {node.parentName}</span>
                          ) : null}
                        </span>
                      </span>
                      {Number(value) === Number(node.id) ? <Check className="ml-2 h-4 w-4 shrink-0" /> : null}
                    </button>
                  ))}
                  {!filteredNodes.length ? (
                    <div className="px-2 py-6 text-center text-sm text-muted-foreground">
                      {t("pages.topic.nodeSelector.empty")}
                    </div>
                  ) : null}
                </div>
              </ScrollArea>
            </div>

            <DialogPrimitive.Close className="absolute top-4 right-4 rounded-sm opacity-70 transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-none">
              <X className="h-4 w-4" />
              <span className="sr-only">Close</span>
            </DialogPrimitive.Close>
          </DialogPrimitive.Content>
        </DialogPrimitive.Portal>
      </DialogPrimitive.Root>
    </div>
  )
}

export function TopicNodeQuickSelector({
  value,
  nodes,
  onChange,
}: {
  value: number
  nodes: TopicNode[]
  onChange: (value: number) => void
}) {
  const { t } = useI18n()
  const [isMobile, setIsMobile] = React.useState(false)

  React.useEffect(() => {
    const query = window.matchMedia("(max-width: 639px)")
    const sync = () => setIsMobile(query.matches)
    sync()
    query.addEventListener("change", sync)
    return () => query.removeEventListener("change", sync)
  }, [])

  const nodeList = React.useMemo(() => flattenTopicNodes(nodes), [nodes])
  const level1NodesRaw = nodeList.filter((node) => !Number(node.parentId))
  const level2ByParent = new Map<number, TopicNode[]>()
  nodeList.forEach((node) => {
    const parentId = Number(node.parentId)
    if (parentId > 0) {
      level2ByParent.set(parentId, [...(level2ByParent.get(parentId) || []), node])
    }
  })
  const selectedNode = nodeList.find((node) => Number(node.id) === Number(value))
  const selectedLevel1Id = selectedNode ? (Number(selectedNode.parentId) > 0 ? Number(selectedNode.parentId) : Number(selectedNode.id)) : null
  const subRowNodes = selectedLevel1Id == null ? [] : level2ByParent.get(selectedLevel1Id) || []
  const visibleLimit = isMobile ? 4 : 8
  const visibleLevel1Nodes = React.useMemo(() => {
    if (level1NodesRaw.length <= visibleLimit) {
      return level1NodesRaw
    }
    const first = level1NodesRaw.slice(0, visibleLimit)
    const hasSelected = first.some((node) => Number(node.id) === selectedLevel1Id)
    if (selectedLevel1Id != null && !hasSelected && first.length > 0) {
      const selected = level1NodesRaw.find((node) => Number(node.id) === selectedLevel1Id)
      if (selected) {
        first[first.length - 1] = selected
      }
    }
    return first
  }, [level1NodesRaw, selectedLevel1Id, visibleLimit])

  function selectLevel1(node: TopicNode) {
    const children = level2ByParent.get(Number(node.id)) || []
    if (children.length) {
      const inFamily = Number(value) === Number(node.id) || children.some((child) => Number(child.id) === Number(value))
      if (!inFamily) onChange(Number(node.id))
      return
    }
    onChange(Number(node.id))
  }

  return (
    <div className="topic-tags flex flex-col gap-2">
      <div className="topic-level1-row">
        {visibleLevel1Nodes.map((node) => (
          <button
            key={node.id}
            type="button"
            className={cn("topic-tag topic-tag-level1", selectedLevel1Id === Number(node.id) ? "selected" : "topic-tag-level1-muted")}
            onClick={() => selectLevel1(node)}
          >
            <span>{node.name}</span>
          </button>
        ))}
        {level1NodesRaw.length > visibleLimit ? (
          <TopicNodeSelector
            value={value}
            nodes={nodes}
            triggerFullWidth={false}
            triggerSize="sm"
            triggerLabel={t("pages.topic.nodeSelector.more")}
            triggerIcon="right"
            onChange={onChange}
          />
        ) : null}
      </div>
      {subRowNodes.length ? (
        <div className="topic-subnodes">
          <div className="topic-subnodes-body">
            <div className="topic-subnodes-header">
              <span className="topic-subnodes-label">{t("pages.topic.nodeSelector.subnodes")}</span>
            </div>
            <div className="topic-subnodes-list">
              {subRowNodes.map((node) => (
                <button
                  key={node.id}
                  type="button"
                  className={cn("topic-tag topic-tag-sub", Number(value) === Number(node.id) && "selected")}
                  onClick={() => onChange(Number(node.id))}
                >
                  <span>{node.name}</span>
                </button>
              ))}
            </div>
          </div>
        </div>
      ) : null}
    </div>
  )
}
