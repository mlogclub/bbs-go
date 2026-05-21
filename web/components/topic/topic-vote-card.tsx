"use client"

import * as React from "react"

import { useIsLogin } from "@/components/app/app-provider"
import { Button } from "@/components/ui/button"
import { apiFetch } from "@/lib/api/client"
import type { TopicVote } from "@/lib/api/types"
import { formatDate } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"
import { msgSuccess, msgWarning } from "@/lib/toast"
import { cn } from "@/lib/utils"

const OPTION_LIMIT = 4

function isSingleVote(vote: TopicVote) {
  return vote.type === 1 || vote.type === "single"
}

function isPkVote(vote: TopicVote) {
  return isSingleVote(vote) && (vote.options || []).length === 2
}

function getPercent(vote: TopicVote, count?: number, showInStyle?: boolean) {
  const total = (vote.options || []).reduce((sum, item) => sum + (item.voteCount || 0), 0)
  let percent = 0

  if (count || total) {
    percent = Number((((count || 0) / total) * 100).toFixed(0))
    if (showInStyle && percent < 30) {
      percent = 30
    }
    percent = Math.min(100, Math.max(0, percent))
  } else {
    percent = showInStyle ? 50 : 0
  }

  if (!vote.voted) {
    return showInStyle ? 50 : percent
  }
  return percent
}

function isChecked(checkedIds: number[], id: number) {
  return checkedIds.includes(id)
}

export function TopicVoteCard({ vote, className }: { vote?: TopicVote | null; className?: string }) {
  if (!vote?.title) {
    return null
  }

  return <TopicVoteCardContent key={vote.id} vote={vote} className={className} />
}

function TopicVoteCardContent({ vote, className }: { vote: TopicVote; className?: string }) {
  const { t } = useI18n()
  const isLogin = useIsLogin()
  const [currentVote, setCurrentVote] = React.useState<TopicVote>(vote)
  const [checkedIds, setCheckedIds] = React.useState<number[]>(vote?.optionIds || [])
  const [showMore, setShowMore] = React.useState(false)
  const [submitting, setSubmitting] = React.useState(false)

  const options = currentVote.options || []
  const canVote = Boolean(isLogin && !currentVote.expired && !currentVote.voted && !submitting)
  const maxNum = Number(currentVote.voteNum) || 1
  const single = isSingleVote(currentVote)
  const pk = isPkVote(currentVote)
  const disabledSubmit = !canVote || (single ? checkedIds.length !== 1 : checkedIds.length === 0 || checkedIds.length > maxNum)
  const visibleOptions = showMore ? options : options.slice(0, OPTION_LIMIT)

  async function submitVote(nextIds = checkedIds) {
    if (!isLogin) {
      msgWarning(t("pages.topic.detail.vote.loginToVote"))
      return
    }
    if (!canVote) {
      return
    }
    if (single && nextIds.length !== 1) {
      msgWarning(t("pages.topic.detail.vote.selectRequired"))
      return
    }
    if (!single && nextIds.length === 0) {
      msgWarning(t("pages.topic.detail.vote.selectAtLeastOne"))
      return
    }
    if (!single && nextIds.length > maxNum) {
      msgWarning(t("pages.topic.detail.vote.maxSelect", { num: maxNum }))
      return
    }

    setSubmitting(true)
    try {
      const voteId = currentVote.id
      const ret = await apiFetch<TopicVote>("/api/vote/cast", {
        method: "POST",
        body: {
          voteId,
          optionIds: nextIds,
        },
      })
      setCurrentVote(ret)
      setCheckedIds(ret.optionIds || [])
      msgSuccess(t("pages.topic.detail.vote.submitSuccess"))
    } catch (error) {
      msgWarning(error instanceof Error ? error.message : String(error))
    } finally {
      setSubmitting(false)
    }
  }

  function toggleOption(optionId: number) {
    if (!canVote || !optionId) {
      return
    }

    if (pk) {
      setCheckedIds([optionId])
      void submitVote([optionId])
      return
    }

    setCheckedIds((current) => {
      if (current.includes(optionId)) {
        return current.filter((id) => id !== optionId)
      }
      if (single) {
        return [optionId]
      }
      if (current.length >= maxNum) {
        msgWarning(t("pages.topic.detail.vote.maxSelect", { num: maxNum }))
        return current
      }
      return [...current, optionId]
    })
  }

  if (pk) {
    const left = options[0]
    const right = options[1]
    return (
      <section className={cn("rounded-lg bg-[#f6f9ff] p-5 text-[#16181f] dark:bg-card dark:text-card-foreground", className)}>
        <h2 className="mb-4 flex items-center text-base leading-none font-medium">
          <span>{currentVote.title}</span>
        </h2>
        <div className="relative flex min-h-24 overflow-hidden rounded bg-background">
          {[left, right].filter(Boolean).map((option, index) => {
            const width = getPercent(currentVote, option.voteCount, true)
            const displayPercent = getPercent(currentVote, option.voteCount, false)
            return (
              <button
                key={option.id}
                type="button"
                disabled={!canVote}
                className={cn(
                  "flex min-h-24 flex-col justify-center px-4 text-sm transition-colors disabled:cursor-default",
                  index === 0
                    ? "items-start bg-[#fceeed] text-[#f15c5b] dark:bg-red-950/40 dark:text-red-200"
                    : "items-end bg-[#deedff] text-[#1050e7] dark:bg-blue-950/40 dark:text-blue-200",
                  option.voted && !currentVote.expired && (index === 0 ? "bg-[#ff4b24] text-white" : "bg-[#3b64fc] text-white")
                )}
                style={{ width: `${width}%` }}
                title={option.content}
                onClick={() => toggleOption(option.id)}
              >
                <span className="line-clamp-2 break-all">{option.content}</span>
                <span className="mt-2 text-xs">
                  {option.voted && index !== 0 ? `${t("pages.topic.detail.vote.votedShort")}  ` : ""}
                  {option.voteCount || 0} ({displayPercent}%)
                  {option.voted && index === 0 ? `  ${t("pages.topic.detail.vote.votedShort")}` : ""}
                </span>
              </button>
            )
          })}
        </div>
        <VoteStatus vote={currentVote} t={t} />
      </section>
    )
  }

  return (
    <section className={cn("rounded-lg bg-[#f6f9ff] px-5 py-6 text-[#16181f] dark:bg-card dark:text-card-foreground", className)}>
      <div className="mb-5 flex flex-wrap items-center gap-1 text-base leading-none font-medium">
        <span className="rounded-sm bg-gradient-to-r from-[#ff603d] to-[#ff881a] px-1.5 py-0.5 text-xs text-white">
          {t("pages.topic.detail.vote.tag")}
        </span>
        <h2>{currentVote.title}</h2>
        <span>{t("pages.topic.detail.vote.titleMeta", { optionCount: options.length, voteNum: maxNum })}</span>
      </div>
      <ul className="space-y-3">
        {visibleOptions.map((option) => (
          <li key={option.id}>
            <button
              type="button"
              disabled={!canVote}
              title={option.content}
              className={cn(
                "relative min-h-9 w-full max-w-[470px] overflow-hidden rounded border border-white bg-white px-1.5 py-1.5 text-left text-sm text-[#737782] disabled:cursor-default dark:border-border dark:bg-background dark:text-muted-foreground",
                canVote && "cursor-pointer hover:text-[#ff7827] dark:hover:text-orange-300",
                canVote && isChecked(checkedIds, option.id) && "border-[#ff7827] text-[#ff7827]",
                canVote && isChecked(checkedIds, option.id) && "dark:border-orange-400 dark:text-orange-300",
                option.voted && "text-[#ff7827] dark:text-orange-300",
                !option.voted && currentVote.expired && "text-[#8b8f99] dark:text-muted-foreground/70"
              )}
              onClick={() => toggleOption(option.id)}
            >
              <span
                className={cn(
                  "absolute inset-y-[-1px] left-[-1px] bg-[#e5eeff] dark:bg-primary/20",
                  option.voted && "bg-[#ffe9dd] dark:bg-orange-500/20"
                )}
                style={{ width: `${getPercent(currentVote, option.voteCount)}%` }}
              />
              <span className="relative z-10 break-all">{option.content}</span>
              {!canVote ? (
                <span className="relative z-10 float-right ml-3 text-[#737782] dark:text-muted-foreground">{option.voteCount || 0}</span>
              ) : null}
            </button>
          </li>
        ))}
      </ul>
      {options.length > OPTION_LIMIT ? (
        <button
          type="button"
          className="mt-3 h-9 w-full max-w-[470px] rounded bg-white text-sm hover:bg-[#fff0ef] dark:bg-background dark:hover:bg-accent"
          onClick={() => setShowMore((current) => !current)}
        >
          {showMore ? t("pages.topic.detail.vote.collapseOptions") : t("pages.topic.detail.vote.expandOptions")}
        </button>
      ) : null}
      <VoteStatus vote={currentVote} t={t} />
      <Button className="mt-6 h-10 w-[188px] bg-gradient-to-r from-[#ff420e] to-[#ff7827] text-white" disabled={disabledSubmit} onClick={() => submitVote()}>
        {currentVote.voted
          ? t("pages.topic.detail.vote.voted")
          : currentVote.expired
            ? t("pages.topic.detail.vote.voteEnded")
            : t("pages.topic.detail.vote.submit")}
      </Button>
    </section>
  )
}

function VoteStatus({ vote, t }: { vote: TopicVote; t: ReturnType<typeof useI18n>["t"] }) {
  return (
    <div className="mt-3 text-sm leading-none text-[#737782] dark:text-muted-foreground">
      {t("pages.topic.detail.vote.participants", { count: vote.voteCount || 0 })}
      {vote.expired ? (
        <span className="ml-3">{t("pages.topic.detail.vote.expired")}</span>
      ) : vote.expiredAt ? (
        <span className="ml-3">
          {t("pages.topic.detail.vote.expiredAt")}: {formatDate(vote.expiredAt)}
        </span>
      ) : null}
    </div>
  )
}
