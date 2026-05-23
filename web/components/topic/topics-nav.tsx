import { TopicsNavContent } from "@/components/topic/topics-nav-content"
import { getCategoryNavs } from "@/lib/api/topics"

export async function TopicsNav({
  currentCategoryId,
  currentRootCategoryId,
}: {
  currentCategoryId?: number
  currentRootCategoryId?: number
}) {
  const nodes = await getCategoryNavs().catch(() => [])

  return (
    <TopicsNavContent
      initialCategories={nodes}
      currentCategoryId={currentCategoryId}
      currentRootCategoryId={currentRootCategoryId}
    />
  )
}
