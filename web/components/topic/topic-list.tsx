import { EmptyState } from "@/components/common/empty-state"
import { TopicListItem } from "@/components/topic/topic-list-item"
import type { Topic } from "@/lib/api/types"
import type { TFunction } from "@/lib/i18n"

export function TopicList({ topics, showSticky, t }: { topics: Topic[]; showSticky?: boolean; t: TFunction }) {
  if (!topics.length) {
    return <EmptyState title={t("common.noData")} />
  }

  return (
    <ul className="divide-y divide-border">
      {topics.map((topic) => (
        <TopicListItem key={topic.id} topic={topic} showSticky={showSticky} t={t} />
      ))}
    </ul>
  )
}
