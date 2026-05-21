import { CheckCircle2, CircleHelp } from "lucide-react"
import Link from "@/components/common/link"

import { UserAvatar } from "@/components/common/avatar"
import { TopicManageMenu } from "@/components/topic/topic-manage-menu"
import type { Topic } from "@/lib/api/types"
import type { UserSummary } from "@/lib/api/types"
import { prettyDate } from "@/lib/format"
import type { TFunction } from "@/lib/i18n"

export function TopicMeta({
  topic,
  currentUser,
  t,
}: {
  topic: Topic
  currentUser?: UserSummary | null
  t: TFunction
}) {
  const displayName =
    topic.user.nickname || topic.user.username || topic.user.id

  return (
    <div className="mt-2.5 flex items-center justify-between gap-3">
      <div className="flex min-w-0 flex-1">
        <div className="mr-2.5 shrink-0">
          <UserAvatar user={topic.user} size={40} />
        </div>
        <div className="min-w-0">
          <div className="mb-0.5">
            <Link
              href={`/user/${topic.user.id}`}
              className="block overflow-hidden text-sm text-ellipsis whitespace-nowrap text-muted-foreground hover:text-primary"
            >
              {displayName}
            </Link>
          </div>
          <div className="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs leading-6 text-muted-foreground">
            {topic.createTime ? (
              <span>
                {t("pages.topic.detail.publishedAt")}{" "}
                <time dateTime={new Date(topic.createTime).toISOString()}>
                  {prettyDate(topic.createTime, t)}
                </time>
              </span>
            ) : null}
            {topic.type === 2 ? (
              <span
                className={`inline-flex items-center rounded-full px-2 py-0.5 text-[11px] leading-none font-medium ring-1 ${
                  topic.qaStatus === "solved"
                    ? "bg-emerald-50 text-emerald-700 ring-emerald-200"
                    : "bg-amber-50 text-amber-700 ring-amber-200"
                }`}
              >
                {topic.qaStatus === "solved" ? (
                  <CheckCircle2 className="mr-1 h-3 w-3" />
                ) : (
                  <CircleHelp className="mr-1 h-3 w-3" />
                )}
                {topic.qaStatus === "solved"
                  ? t("pages.topic.detail.qaSolved")
                  : t("pages.topic.detail.qaUnsolved")}
              </span>
            ) : null}
            {topic.type === 2 && topic.bountyScore ? (
              <span className="inline-flex items-center rounded-full bg-amber-100 px-2 py-0.5 text-[11px] leading-none font-medium text-amber-800 ring-1 ring-amber-200">
                {t("pages.topic.detail.bountyLabel", {
                  score: topic.bountyScore,
                })}
              </span>
            ) : null}
            {topic.ipLocation ? (
              <span className="text-xs">
                {t("pages.topic.detail.ipLocation")}
                {topic.ipLocation}
              </span>
            ) : null}
          </div>
        </div>
      </div>
      <div className="min-w-max">
        <TopicManageMenu topic={topic} currentUser={currentUser} />
      </div>
    </div>
  )
}
