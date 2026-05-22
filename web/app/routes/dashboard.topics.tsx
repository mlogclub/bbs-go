"use client"

import * as React from "react"
import {
  CheckCircleIcon,
  EyeIcon,
  ExternalLinkIcon,
  LightbulbIcon,
  MessageCircleIcon,
  HeartIcon,
  RotateCcwIcon,
  RefreshCwIcon,
  SearchIcon,
  StarIcon,
  StarOffIcon,
  Trash2Icon,
  Undo2Icon,
} from "lucide-react"

import { DashboardSelect } from "@/components/dashboard/dashboard-select"
import {
  ConfirmDialog,
  type ConfirmDialogState,
} from "@/components/dashboard/confirm-dialog"
import { useCurrentUser } from "@/components/app/app-provider"
import { ErrorPage } from "@/components/common/error-page"
import { PreviewableImage } from "@/components/common/image-preview"
import { DashboardPagination } from "@/components/dashboard/pagination-controls"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  adminDelete,
  adminList,
  adminPostForm,
  type AdminFormValue,
  type AdminRecord,
} from "@/lib/api/admin"
import { formatDateTime } from "@/lib/format"
import { userHasPermission } from "@/lib/auth/roles"
import { createAdminInitialFilters } from "@/lib/dashboard/default-filters"
import { useI18n } from "@/lib/i18n/provider"
import { msgSuccess } from "@/lib/toast"
import { cn } from "@/lib/utils"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

type TopicRecord = AdminRecord & {
  id?: number
  idEncode?: string
  title?: string
  summary?: string
  content?: string
  type?: number
  status?: number
  qaStatus?: string
  recommend?: boolean
  createTime?: number
  viewCount?: number
  commentCount?: number
  likeCount?: number
  user?: {
    id?: number
    idEncode?: string
    nickname?: string
    username?: string
    avatar?: string
  }
  node?: {
    name?: string
  }
  tags?: Array<{
    id?: number
    name?: string
  }>
  imageList?: Array<{
    preview?: string
    url?: string
  }>
  vote?: {
    title?: string
    optionCount?: number
    voteNum?: number
    voteCount?: number
    expired?: boolean
    expiredAt?: number
    options?: Array<{
      id?: number
      content?: string
      voteCount?: number
      percent?: number
    }>
  }
}

type TopicAction =
  | "recommend"
  | "unrecommend"
  | "audit"
  | "undelete"
  | "delete"
  | "solved"
  | "unsolved"

function topicTypeLabel(t: ReturnType<typeof useI18n>["t"], type?: number) {
  if (type === 1) return t("dashboard.topicFeed.typeTweet")
  if (type === 2) return t("dashboard.topicFeed.typeQa")
  return t("dashboard.topicFeed.typeTopic")
}

function topicStatusLabel(t: ReturnType<typeof useI18n>["t"], status?: number) {
  if (status === 1) return t("dashboard.topicFeed.statusDeleted")
  if (status === 2) return t("dashboard.topicFeed.statusReview")
  return t("dashboard.topicFeed.statusNormal")
}

function topicActionSuccessMessage(
  t: ReturnType<typeof useI18n>["t"],
  action: TopicAction
) {
  const messageKeys: Record<TopicAction, string> = {
    recommend: "dashboard.messages.recommended",
    unrecommend: "dashboard.messages.unrecommended",
    audit: "dashboard.messages.audited",
    undelete: "dashboard.messages.restored",
    delete: "dashboard.messages.deleted",
    solved: "dashboard.messages.markedSolved",
    unsolved: "dashboard.messages.markedUnsolved",
  }

  return t(messageKeys[action])
}

function compactText(value: unknown) {
  if (typeof value !== "string") return ""
  return value.replace(/\s+/g, " ").trim()
}

function voteOptionPercent(
  option: NonNullable<NonNullable<TopicRecord["vote"]>["options"]>[number],
  total: number
) {
  if (typeof option.percent === "number") {
    return Math.max(0, Math.min(100, option.percent))
  }
  if (total <= 0) return 0
  return Math.max(
    0,
    Math.min(100, Math.round(((option.voteCount || 0) / total) * 100))
  )
}

export default function DashboardTopicsRoute() {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const [filters, setFilters] = React.useState<Record<string, AdminFormValue>>(
    () => createAdminInitialFilters({ status: 0 }, 20)
  )
  const [records, setRecords] = React.useState<TopicRecord[]>([])
  const [total, setTotal] = React.useState(0)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [confirmState, setConfirmState] =
    React.useState<ConfirmDialogState>(null)

  const page = Number(filters.page || 1)
  const limit = Number(filters.limit || 20)
  const pageCount = Math.max(1, Math.ceil(total / limit))
  const canView = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_TOPIC_VIEW)
  const canRecommend = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_TOPIC_RECOMMEND
  )
  const canAudit = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_TOPIC_AUDIT)
  const canDelete = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_TOPIC_DELETE)
  const canSolve = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_TOPIC_SOLVE)

  const load = React.useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminList<TopicRecord>(
        "/api/admin/topic/list",
        filters
      )
      setRecords(data.results || [])
      setTotal(data.page?.total ?? data.results?.length ?? 0)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("dashboard.errors.loadFailed")
      )
    } finally {
      setLoading(false)
    }
  }, [filters, t])

  React.useEffect(() => {
    void load()
  }, [load])

  function updateFilter(name: string, value: AdminFormValue) {
    setFilters((current) => ({
      ...current,
      [name]: value,
      page: name === "page" ? value : name === "limit" ? current.page : 1,
    }))
  }

  function runAction(topic: TopicRecord, action: TopicAction) {
    const allowed = {
      recommend: canRecommend,
      unrecommend: canRecommend,
      audit: canAudit,
      undelete: canDelete,
      delete: canDelete,
      solved: canSolve,
      unsolved: canSolve,
    }
    if (!allowed[action]) return

    if (action === "delete") {
      setConfirmState({
        description: t("dashboard.confirmDelete"),
        confirmText: t("dashboard.actions.delete"),
        onConfirm: () => {
          void performAction(topic, action)
        },
      })
      return
    }

    void performAction(topic, action)
  }

  async function performAction(topic: TopicRecord, action: TopicAction) {
    const id = topic.id
    if (!id) return

    const endpoints = {
      recommend: "/api/admin/topic/recommend",
      unrecommend: "/api/admin/topic/recommend",
      audit: "/api/admin/topic/audit",
      undelete: "/api/admin/topic/undelete",
      delete: "/api/admin/topic/delete",
      solved: "/api/admin/topic/mark_solved",
      unsolved: "/api/admin/topic/mark_unsolved",
    }

    setError(null)
    try {
      if (action === "unrecommend") {
        await adminDelete(endpoints[action], { id })
      } else {
        await adminPostForm(endpoints[action], { id })
      }
      msgSuccess(topicActionSuccessMessage(t, action))
      await load()
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("dashboard.errors.actionFailed")
      )
    }
  }

  if (!canView) {
    return <ErrorPage statusCode={403} />
  }

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-4 md:p-6">
      <section className="flex flex-col gap-3 rounded-lg border bg-[var(--dashboard-panel)] p-3 text-card-foreground shadow-xs">
        <div className="flex flex-wrap items-end gap-2">
          <FilterInput
            label={t("dashboard.fields.id")}
            value={filters.id}
            onChange={(value) => updateFilter("id", value)}
          />
          <FilterInput
            label={t("dashboard.fields.userId")}
            value={filters.userId}
            onChange={(value) => updateFilter("userId", value)}
          />
          <FilterInput
            label={t("dashboard.fields.title")}
            value={filters.title}
            onChange={(value) => updateFilter("title", value)}
          />
          <FilterSelect
            label={t("dashboard.fields.status")}
            value={filters.status}
            options={[
              { label: t("dashboard.topicFeed.statusNormal"), value: 0 },
              { label: t("dashboard.topicFeed.statusDeleted"), value: 1 },
              { label: t("dashboard.topicFeed.statusReview"), value: 2 },
            ]}
            onChange={(value) => updateFilter("status", value)}
          />
          <FilterSelect
            label={t("dashboard.fields.type")}
            value={filters.type}
            options={[
              { label: t("dashboard.topicFeed.typeTopic"), value: 0 },
              { label: t("dashboard.topicFeed.typeTweet"), value: 1 },
              { label: t("dashboard.topicFeed.typeQa"), value: 2 },
            ]}
            onChange={(value) => updateFilter("type", value)}
          />
          <FilterSelect
            label={t("dashboard.fields.recommend")}
            value={filters.recommend}
            options={[
              { label: t("dashboard.boolean.yes"), value: "true" },
              { label: t("dashboard.boolean.no"), value: "false" },
            ]}
            onChange={(value) => updateFilter("recommend", value)}
          />
          <FilterSelect
            label={t("dashboard.fields.qaStatus")}
            value={filters.qaStatus}
            options={[
              { label: t("dashboard.topicFeed.qaSolved"), value: "solved" },
              { label: t("dashboard.topicFeed.qaUnsolved"), value: "unsolved" },
            ]}
            onChange={(value) => updateFilter("qaStatus", value)}
          />
          <Button onClick={() => void load()} disabled={loading}>
            <SearchIcon />
            {t("dashboard.actions.search")}
          </Button>
          <Button
            variant="outline"
            size="icon"
            onClick={() => void load()}
            disabled={loading}
          >
            <RefreshCwIcon />
            <span className="sr-only">{t("dashboard.actions.refresh")}</span>
          </Button>
        </div>

        {error ? (
          <div className="rounded-md border border-destructive/25 bg-destructive/10 px-3 py-2 text-sm text-destructive">
            {error}
          </div>
        ) : null}
      </section>

      <section
        className="rounded-lg border bg-[var(--dashboard-panel)] shadow-xs"
        aria-busy={loading}
      >
        <div className="divide-y">
          {loading && records.length === 0 ? (
            <div className="px-4 py-16 text-center text-sm text-muted-foreground">
              {t("dashboard.loading")}
            </div>
          ) : records.length ? (
            records.map((topic) => (
              <TopicFeedItem
                key={topic.id}
                topic={topic}
                permissions={{
                  recommend: canRecommend,
                  audit: canAudit,
                  delete: canDelete,
                  solve: canSolve,
                }}
                onAction={(action) => runAction(topic, action)}
              />
            ))
          ) : (
            <div className="px-4 py-16 text-center text-sm text-muted-foreground">
              {t("common.noData")}
            </div>
          )}
        </div>

        <DashboardPagination
          page={page}
          pageCount={pageCount}
          total={total}
          limit={limit}
          loading={loading}
          onPageChange={(nextPage) => updateFilter("page", nextPage)}
          onLimitChange={(nextLimit) =>
            setFilters((current) => ({
              ...current,
              page: 1,
              limit: nextLimit,
            }))
          }
        />
      </section>
      <ConfirmDialog
        state={confirmState}
        onOpenChange={(open) => {
          if (!open) setConfirmState(null)
        }}
      />
    </div>
  )
}

function TopicFeedItem({
  topic,
  permissions,
  onAction,
}: {
  topic: TopicRecord
  permissions: {
    recommend: boolean
    audit: boolean
    delete: boolean
    solve: boolean
  }
  onAction: (action: TopicAction) => void
}) {
  const { t } = useI18n()
  const body = compactText(topic.type === 1 ? topic.content : topic.summary)
  const userName =
    topic.user?.nickname ||
    topic.user?.username ||
    t("dashboard.user.anonymous")
  const topicUrl = `/topic/${topic.idEncode || topic.id}`
  const userUrl = topic.user?.idEncode
    ? `/user/${topic.user.idEncode}`
    : topic.user?.id
      ? `/user/${topic.user.id}`
      : undefined
  const voteOptions = topic.vote?.options || []
  const previewSrcList =
    topic.imageList
      ?.map((image) => image.url || image.preview || "")
      .filter(Boolean) || []
  const voteTotal = voteOptions.reduce(
    (total, option) => total + (option.voteCount || 0),
    0
  )

  return (
    <article className="grid gap-3 p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="flex min-w-0 items-start gap-3">
          <a
            href={userUrl || "#"}
            target={userUrl ? "_blank" : undefined}
            rel={userUrl ? "noreferrer" : undefined}
            aria-disabled={!userUrl}
            className="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-full border bg-muted text-sm font-medium"
          >
            {topic.user?.avatar ? (
              <img
                src={topic.user.avatar}
                alt={userName}
                className="size-full object-cover"
                loading="lazy"
              />
            ) : (
              userName.slice(0, 1)
            )}
          </a>
          <div className="min-w-0">
            <a
              href={topicUrl}
              target="_blank"
              rel="noreferrer"
              className="line-clamp-2 text-base font-semibold hover:underline"
            >
              {topic.title || t("dashboard.topicFeed.untitled")}
            </a>
            <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
              {userUrl ? (
                <a
                  href={userUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="hover:text-foreground hover:underline"
                >
                  {userName}
                </a>
              ) : (
                <span>{userName}</span>
              )}
              <span>{formatDateTime(topic.createTime ?? null) || "-"}</span>
              <span>ID: {topic.id}</span>
            </div>
          </div>
        </div>

        <div className="flex flex-wrap justify-end gap-1">
          <TopicTag>{topicTypeLabel(t, topic.type)}</TopicTag>
          <TopicTag
            className={cn(
              topic.status === 1 && "border-destructive/30 text-destructive",
              topic.status === 2 && "border-amber-400/40 text-amber-700"
            )}
          >
            {topicStatusLabel(t, topic.status)}
          </TopicTag>
          {topic.recommend ? (
            <TopicTag className="border-primary/30 text-primary">
              {t("dashboard.fields.recommend")}
            </TopicTag>
          ) : null}
          {topic.type === 2 && topic.qaStatus ? (
            <TopicTag>
              {topic.qaStatus === "solved"
                ? t("dashboard.topicFeed.qaSolved")
                : t("dashboard.topicFeed.qaUnsolved")}
            </TopicTag>
          ) : null}
        </div>
      </div>

      {body ? (
        <p className="line-clamp-3 w-full text-sm leading-6 text-muted-foreground">
          {body}
        </p>
      ) : null}

      {topic.imageList?.length ? (
        <div className="flex flex-wrap gap-2">
          {topic.imageList.slice(0, 6).map((image, index) => {
            const src = image.url || image.preview || ""
            return src ? (
              <PreviewableImage
                key={`${src}-${index}`}
                src={src}
                previewSrcList={previewSrcList}
                initialIndex={index}
                alt=""
                className="size-24 cursor-zoom-in rounded-md border object-cover"
                loading="lazy"
              />
            ) : null
          })}
        </div>
      ) : null}

      {topic.vote ? (
        <div className="grid w-200 gap-3 rounded-lg border bg-muted/20 p-3">
          <div className="flex flex-wrap items-center justify-between gap-3 text-sm font-medium">
            <span className="inline-flex min-w-0 items-center gap-2">
              <span className="rounded bg-primary/10 px-1.5 py-0.5 text-xs font-medium text-primary">
                {t("dashboard.fields.vote")}
              </span>
              <span className="truncate">
                {topic.vote.title || t("dashboard.fields.vote")}
              </span>
              <span className="text-xs font-normal text-muted-foreground">
                {t("dashboard.topicFeed.voteMeta", {
                  optionCount: topic.vote.optionCount || voteOptions.length,
                  voteNum: topic.vote.voteNum || 1,
                })}
              </span>
            </span>
            <span className="text-xs text-muted-foreground">
              {t("dashboard.topicFeed.voteParticipants", {
                count: topic.vote.voteCount ?? voteTotal,
              })}
            </span>
          </div>
          <div className="grid gap-2">
            {voteOptions.map((option) => {
              const percent = voteOptionPercent(option, voteTotal)
              return (
                <div key={option.id ?? option.content} className="grid gap-1">
                  <div className="flex items-center justify-between gap-3 text-xs">
                    <span className="truncate">{option.content || "-"}</span>
                    <span className="text-muted-foreground">
                      {option.voteCount ?? 0} ({percent}%)
                    </span>
                  </div>
                  <div className="h-1.5 overflow-hidden rounded-full bg-muted">
                    <div
                      className="h-full rounded-full bg-primary"
                      style={{ width: `${percent}%` }}
                    />
                  </div>
                </div>
              )
            })}
          </div>
          <div className="text-xs text-muted-foreground">
            {topic.vote.expired
              ? t("dashboard.topicFeed.voteExpired")
              : topic.vote.expiredAt
                ? t("dashboard.topicFeed.voteExpiredAt", {
                    time: formatDateTime(topic.vote.expiredAt) || "-",
                  })
                : null}
          </div>
        </div>
      ) : null}

      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
          {topic.node?.name ? <span>{topic.node.name}</span> : null}
          {topic.tags?.map((tag) =>
            tag.name ? <span key={tag.id || tag.name}>#{tag.name}</span> : null
          )}
          <span className="inline-flex items-center gap-1">
            <EyeIcon className="size-3.5" />
            {topic.viewCount ?? 0}
          </span>
          <span className="inline-flex items-center gap-1">
            <MessageCircleIcon className="size-3.5" />
            {topic.commentCount ?? 0}
          </span>
          <span className="inline-flex items-center gap-1">
            <HeartIcon className="size-3.5" />
            {topic.likeCount ?? 0}
          </span>
        </div>

        <div className="flex flex-wrap justify-end gap-2">
          <Button size="sm" variant="outline" asChild>
            <a href={topicUrl} target="_blank" rel="noreferrer">
              <ExternalLinkIcon />
              {t("dashboard.actions.view")}
            </a>
          </Button>
          {permissions.recommend && topic.status === 0 && topic.recommend ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("unrecommend")}
            >
              <StarOffIcon />
              {t("dashboard.actions.unrecommend")}
            </Button>
          ) : null}
          {permissions.recommend && topic.status === 0 && !topic.recommend ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("recommend")}
            >
              <StarIcon />
              {t("dashboard.actions.recommend")}
            </Button>
          ) : null}
          {permissions.audit && topic.status === 2 ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("audit")}
            >
              <CheckCircleIcon />
              {t("dashboard.actions.audit")}
            </Button>
          ) : null}
          {permissions.delete && topic.status === 1 ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("undelete")}
            >
              <Undo2Icon />
              {t("dashboard.actions.undelete")}
            </Button>
          ) : null}
          {topic.status === 0 &&
          permissions.solve &&
          topic.type === 2 &&
          topic.qaStatus !== "solved" ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("solved")}
            >
              <LightbulbIcon />
              {t("dashboard.actions.markSolved")}
            </Button>
          ) : null}
          {topic.status === 0 &&
          permissions.solve &&
          topic.type === 2 &&
          topic.qaStatus === "solved" ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("unsolved")}
            >
              <RotateCcwIcon />
              {t("dashboard.actions.markUnsolved")}
            </Button>
          ) : null}
          {permissions.delete && (topic.status === 0 || topic.status === 2) ? (
            <Button
              size="sm"
              variant="destructive"
              onClick={() => onAction("delete")}
            >
              <Trash2Icon />
              {t("dashboard.actions.delete")}
            </Button>
          ) : null}
        </div>
      </div>
    </article>
  )
}

function TopicTag({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium",
        className
      )}
    >
      {children}
    </span>
  )
}

function FilterInput({
  label,
  value,
  onChange,
}: {
  label: string
  value: AdminFormValue
  onChange: (value: AdminFormValue) => void
}) {
  return (
    <div className="grid min-w-44 gap-1.5">
      <Label className="text-xs text-muted-foreground">{label}</Label>
      <Input
        value={value === undefined || value === null ? "" : String(value)}
        placeholder={label}
        onChange={(event) => onChange(event.target.value)}
      />
    </div>
  )
}

function FilterSelect({
  label,
  value,
  options,
  onChange,
}: {
  label: string
  value: AdminFormValue
  options: Array<{ label: string; value: string | number }>
  onChange: (value: AdminFormValue) => void
}) {
  return (
    <div className="grid min-w-44 gap-1.5">
      <Label className="text-xs text-muted-foreground">{label}</Label>
      <DashboardSelect
        value={value}
        options={options}
        placeholder={label}
        onValueChange={(nextValue) => onChange(nextValue)}
      />
    </div>
  )
}
