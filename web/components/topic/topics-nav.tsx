import { TopicsNavContent } from "@/components/topic/topics-nav-content"
import { getTopicNodeNavs } from "@/lib/api/topics"

export async function TopicsNav({
  currentNodeId,
  currentRootNodeId,
}: {
  currentNodeId?: number
  currentRootNodeId?: number
}) {
  const nodes = await getTopicNodeNavs().catch(() => [])

  return (
    <TopicsNavContent
      initialNodes={nodes}
      currentNodeId={currentNodeId}
      currentRootNodeId={currentRootNodeId}
    />
  )
}
