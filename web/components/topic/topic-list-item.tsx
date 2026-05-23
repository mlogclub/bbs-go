import { CheckCircle2, CircleHelp, MessageCircle } from "lucide-react"
import Link from "@/components/common/link"

import { UserAvatar } from "@/components/common/avatar"
import { TopicLikeButton } from "@/components/topic/topic-like-button"
import { TopicVoteCard } from "@/components/topic/topic-vote-card"
import type { Topic } from "@/lib/api/types"
import { prettyDate } from "@/lib/format"
import type { TFunction } from "@/lib/i18n"

function getTopicImageSizeClass(count: number) {
  if (count <= 1) {
    return "h-[160px] w-[160px] sm:h-[210px] sm:w-[210px]"
  }
  if (count === 2) {
    return "h-[128px] w-[128px] sm:h-[180px] sm:w-[180px]"
  }
  return "h-[94px] w-[94px] sm:h-[120px] sm:w-[120px]"
}

export function TopicListItem({
  topic,
  showSticky,
  t,
}: {
  topic: Topic
  showSticky?: boolean
  t: TFunction
}) {
  const displayName =
    topic.user.nickname || topic.user.username || topic.user.id
  const topicHref = `/topic/${topic.id}`
  const imageSizeClass = getTopicImageSizeClass(topic.imageList?.length || 0)

  return (
    <li className="px-4 py-3">
      <div className="flex items-center justify-between gap-3">
        <div className="flex min-w-0 items-center gap-2">
          <UserAvatar user={topic.user} size={24} className="shrink-0" />
          <div className="flex min-w-0 items-center gap-1 text-xs md:text-sm">
            <Link
              href={`/user/${topic.user.id}`}
              target="_blank"
              className="max-w-32 truncate text-muted-foreground hover:underline"
            >
              {displayName}
            </Link>
            <span className="text-muted-foreground">·</span>
            <span className="truncate text-muted-foreground">
              {prettyDate(topic.createTime, t)}
            </span>
          </div>
        </div>
        {showSticky && topic.sticky ? (
          <span className="inline-flex items-center rounded-sm bg-orange-100 px-1.5 py-0.5 text-[11px] text-orange-700">
            {t("component.topicList.sticky")}
          </span>
        ) : null}
      </div>

      <div className="mt-2 space-y-2">
        {topic.type !== 1 ? (
          <>
            <Link
              href={topicHref}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 text-[15px] leading-6 font-semibold break-all text-foreground sm:text-base"
            >
              {topic.type === 2 ? (
                <span
                  className={`inline-flex h-5 items-center rounded-full px-2 text-[11px] leading-none font-medium ring-1 ${
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
                    ? t("component.topicList.qaSolved")
                    : t("component.topicList.qaUnsolved")}
                </span>
              ) : null}
              {topic.type === 2 && topic.bountyScore ? (
                <span className="inline-flex h-5 items-center rounded-full bg-amber-100 px-2 text-[11px] leading-none font-medium text-amber-800 ring-1 ring-amber-200">
                  {t("pages.topic.detail.bountyLabel", {
                    score: topic.bountyScore,
                  })}
                </span>
              ) : null}
              {topic.title}
            </Link>
            {topic.summary ? (
              <Link
                href={topicHref}
                target="_blank"
                rel="noopener noreferrer"
                className="line-clamp-3 block text-[15px] leading-6 break-all text-muted-foreground hover:text-foreground/80 sm:text-sm sm:leading-normal"
              >
                {topic.summary}
              </Link>
            ) : null}
          </>
        ) : (
          <>
            {topic.content ? (
              <Link
                href={topicHref}
                target="_blank"
                rel="noopener noreferrer"
                className="line-clamp-3 block text-[15px] leading-6 break-all whitespace-pre-line text-foreground sm:text-sm sm:leading-normal"
              >
                {topic.content}
              </Link>
            ) : null}
            {topic.imageList?.length ? (
              <ul className="mt-1 flex flex-wrap gap-2">
                {topic.imageList.slice(0, 9).map((image, index) => (
                  <li
                    key={`${image.preview || image.url || "image"}-${index}`}
                    className={imageSizeClass}
                  >
                    <Link
                      href={topicHref}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="block h-full w-full overflow-hidden rounded-sm bg-muted"
                    >                      <img
                        src={image.preview || image.url}
                        alt=""
                        className="h-full w-full object-cover transition-transform duration-300 hover:scale-105"
                      />
                    </Link>
                  </li>
                ))}
              </ul>
            ) : null}
          </>
        )}
      </div>

      {topic.vote ? (
        <TopicVoteCard className="mt-2 mb-2" vote={topic.vote} />
      ) : null}

      <div className="mt-2 flex flex-wrap items-center justify-between gap-2">
        <div className="min-w-0">
          {topic.category ? (
            <Link
              href={`/topics/category/${topic.category.id}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex max-w-full items-center gap-1 rounded-full bg-accent px-2.5 py-1 text-xs text-muted-foreground hover:text-foreground"
            >
              {topic.category.logo ? (
                <img
                  src={topic.category.logo}
                  alt=""
                  className="h-4 w-4 rounded-full object-cover"
                />
              ) : null}
              <span className="truncate">{topic.category.name}</span>
            </Link>
          ) : null}
        </div>
        <div className="ml-auto flex items-center gap-4 text-xs text-muted-foreground">
          <TopicLikeButton
            topicId={topic.id}
            initialLiked={topic.liked}
            initialLikeCount={topic.likeCount}
          />
          <Link
            href={topicHref}
            className="inline-flex min-h-8 items-center gap-1.5 transition-colors hover:text-primary"
          >
            <MessageCircle className="h-4 w-4" />
            <span className="min-w-[1ch] text-sm leading-none">
              {topic.commentCount || ""}
            </span>
          </Link>
        </div>
      </div>
    </li>
  )
}
