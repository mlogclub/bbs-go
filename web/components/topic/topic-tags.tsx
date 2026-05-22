import Link from "@/components/common/link"

import { UserAvatar } from "@/components/common/avatar"
import type { Topic, UserSummary } from "@/lib/api/types"

export function TopicTags({
  topic,
  likeUsers,
}: {
  topic: Topic
  likeUsers?: UserSummary[] | null
}) {
  const tags = topic.tags || []

  return (
    <div className="mb-4 flex flex-col gap-2 px-4 lg:flex-row lg:items-center lg:justify-between">
      <div>
        {topic.node ? (
          <Link
            href={`/topics/node/${topic.node.id}`}
            className="mr-2.5 inline-flex items-center justify-center rounded-[12.5px] border border-border bg-muted px-2 py-0.5 text-xs text-muted-foreground hover:bg-background hover:text-primary"
          >
            {topic.node.name}
          </Link>
        ) : null}
        {tags.map((tag) => (
          <Link
            key={tag.id}
            href={`/topics/tag/${tag.id}`}
            className="mr-2.5 inline-flex items-center justify-center rounded-[12.5px] border border-border bg-muted px-2 py-0.5 text-xs text-muted-foreground hover:bg-background hover:text-primary"
          >
            #{tag.name}
          </Link>
        ))}
      </div>
      {likeUsers?.length ? (
        <div className="flex flex-wrap items-center lg:justify-end [&_.avatar-a]:-mr-0.75">
          {likeUsers.map((user) => (
            <UserAvatar
              key={user.id}
              user={user}
              size={24}
              className="border-2 border-background"
            />
          ))}
          <span className="ml-2 text-sm text-foreground">
            {topic.likeCount}
          </span>
        </div>
      ) : null}
    </div>
  )
}
