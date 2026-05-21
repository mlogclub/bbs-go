import type { TopicNode } from "@/lib/api/types"

export function flattenTopicNodes(nodes: TopicNode[] = []): TopicNode[] {
  const result: TopicNode[] = []

  function walk(list: TopicNode[]) {
    list.forEach((node) => {
      result.push(node)
      if (node.children?.length) {
        walk(node.children)
      }
    })
  }

  walk(Array.isArray(nodes) ? nodes : [])
  return result
}

export function filterTopicNodeTree(
  nodes: TopicNode[] = [],
  predicate: (item: TopicNode) => boolean
): TopicNode[] {
  if (!Array.isArray(nodes) || !nodes.length) {
    return []
  }

  return nodes.reduce<TopicNode[]>((result, node) => {
    const children = filterTopicNodeTree(node.children || [], predicate)
    if (!predicate(node) && !children.length) {
      return result
    }
    result.push({ ...node, children })
    return result
  }, [])
}

export function hasTopicNode(nodes: TopicNode[] = [], nodeId: number | string) {
  const targetId = Number(nodeId)
  return flattenTopicNodes(nodes).some((node) => Number(node.id) === targetId)
}

export function getFirstTopicNodeId(nodes: TopicNode[] = []) {
  const first = flattenTopicNodes(nodes)[0]
  return first ? Number(first.id) : 0
}
