"use client"

import * as React from "react"
import {
  EyeIcon,
  MessageCircleIcon,
  RefreshCwIcon,
  SearchIcon,
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
  adminList,
  adminPostForm,
  type AdminFormValue,
  type AdminRecord,
} from "@/lib/api/admin"
import type { ImageInfo } from "@/lib/api/types"
import { userHasPermission } from "@/lib/auth/roles"
import { createAdminInitialFilters } from "@/lib/dashboard/default-filters"
import { formatDateTime } from "@/lib/format"
import { useI18n } from "@/lib/i18n/provider"
import { msgSuccess } from "@/lib/toast"
import { cn } from "@/lib/utils"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

type ArticleRecord = AdminRecord & {
  id?: number
  title?: string
  summary?: string
  status?: number
  createTime?: number
  viewCount?: number
  commentCount?: number
  cover?: ImageInfo | string | null
  user?: {
    id?: number
    idEncode?: string
    nickname?: string
    username?: string
    avatar?: string
  }
  tags?: Array<{
    id?: number
    name?: string
  }>
}

type ArticleAction = "audit" | "delete"

function articleStatusLabel(
  t: ReturnType<typeof useI18n>["t"],
  status?: number
) {
  if (status === 1) return t("dashboard.topicFeed.statusDeleted")
  if (status === 2) return t("dashboard.topicFeed.statusReview")
  return t("dashboard.topicFeed.statusNormal")
}

function articleActionSuccessMessage(
  t: ReturnType<typeof useI18n>["t"],
  action: ArticleAction
) {
  if (action === "audit") return t("dashboard.messages.audited")
  return t("dashboard.messages.deleted")
}

function compactText(value: unknown) {
  if (typeof value !== "string") return ""
  return value.replace(/\s+/g, " ").trim()
}

function imageSrc(image: ImageInfo | string | null | undefined) {
  if (typeof image === "string") return image
  return image?.preview || image?.url || ""
}

export default function DashboardArticlesRoute() {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const [filters, setFilters] = React.useState<Record<string, AdminFormValue>>(
    () => createAdminInitialFilters({ status: 0 }, 20)
  )
  const [records, setRecords] = React.useState<ArticleRecord[]>([])
  const [total, setTotal] = React.useState(0)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [confirmState, setConfirmState] =
    React.useState<ConfirmDialogState>(null)

  const page = Number(filters.page || 1)
  const limit = Number(filters.limit || 20)
  const pageCount = Math.max(1, Math.ceil(total / limit))
  const canView = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_ARTICLE_VIEW)
  const canAudit = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_ARTICLE_AUDIT)
  const canDelete = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_ARTICLE_DELETE)

  const load = React.useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminList<ArticleRecord>(
        "/api/admin/article/list",
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

  function runAction(article: ArticleRecord, action: ArticleAction) {
    if (action === "audit" && !canAudit) return
    if (action === "delete" && !canDelete) return

    if (action === "delete") {
      setConfirmState({
        description: t("dashboard.confirmDelete"),
        confirmText: t("dashboard.actions.delete"),
        onConfirm: () => {
          void performAction(article, action)
        },
      })
      return
    }

    void performAction(article, action)
  }

  async function performAction(article: ArticleRecord, action: ArticleAction) {
    const id = article.id
    if (!id) return

    setError(null)
    try {
      await adminPostForm(
        action === "audit"
          ? "/api/admin/article/audit"
          : "/api/admin/article/delete",
        { id }
      )
      msgSuccess(articleActionSuccessMessage(t, action))
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
            records.map((article) => (
              <ArticleFeedItem
                key={article.id}
                article={article}
                permissions={{ audit: canAudit, delete: canDelete }}
                onAction={(action) => runAction(article, action)}
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

function ArticleFeedItem({
  article,
  permissions,
  onAction,
}: {
  article: ArticleRecord
  permissions: {
    audit: boolean
    delete: boolean
  }
  onAction: (action: ArticleAction) => void
}) {
  const { t } = useI18n()
  const userName =
    article.user?.nickname ||
    article.user?.username ||
    t("dashboard.user.anonymous")
  const articleUrl = `/article/${article.id}`
  const coverSrc = imageSrc(article.cover)
  const userUrl = article.user?.idEncode
    ? `/user/${article.user.idEncode}`
    : article.user?.id
      ? `/user/${article.user.id}`
      : undefined

  return (
    <article className="grid gap-3 p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <a
          href={articleUrl}
          target="_blank"
          rel="noreferrer"
          className="line-clamp-2 min-w-0 text-base font-semibold hover:underline"
        >
          {article.title || t("dashboard.topicFeed.untitled")}
        </a>
        <ArticleTag
          className={cn(
            article.status === 1 && "border-destructive/30 text-destructive",
            article.status === 2 && "border-amber-400/40 text-amber-700"
          )}
        >
          {articleStatusLabel(t, article.status)}
        </ArticleTag>
      </div>

      <div className="grid gap-3 md:grid-cols-[minmax(0,1fr)_150px] md:items-center">
        <div className="min-w-0">
          {compactText(article.summary) ? (
            <p className="line-clamp-3 w-full text-sm leading-6 text-muted-foreground">
              {compactText(article.summary)}
            </p>
          ) : null}

          <div className="mt-2 flex flex-wrap items-center justify-between gap-x-3 gap-y-2 text-xs text-muted-foreground">
            <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
              <a
                href={userUrl || "#"}
                target={userUrl ? "_blank" : undefined}
                rel="noreferrer"
                className="hover:text-foreground"
              >
                {userName}
              </a>
              <span>{formatDateTime(article.createTime ?? null) || "-"}</span>
              <span>ID: {article.id}</span>
            </div>
            {article.tags?.length ? (
              <div className="flex flex-wrap items-center gap-1.5">
                {article.tags.map((tag) =>
                  tag.name ? (
                    <ArticleTag key={tag.id || tag.name}>
                      #{tag.name}
                    </ArticleTag>
                  ) : null
                )}
              </div>
            ) : null}
          </div>
        </div>

        {coverSrc ? (
          <div className="aspect-[5/3] overflow-hidden rounded-md border bg-muted md:w-[150px]">
            <PreviewableImage
              src={coverSrc}
              previewSrcList={[coverSrc]}
              alt={article.title || ""}
              className="size-full cursor-zoom-in object-cover"
              loading="lazy"
            />
          </div>
        ) : null}
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
          <span className="inline-flex items-center gap-1">
            <EyeIcon className="size-3.5" />
            {article.viewCount ?? 0}
          </span>
          <span className="inline-flex items-center gap-1">
            <MessageCircleIcon className="size-3.5" />
            {article.commentCount ?? 0}
          </span>
        </div>
        <div className="flex flex-wrap justify-end gap-2">
          {permissions.audit && article.status === 2 ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction("audit")}
            >
              <Undo2Icon />
              {t("dashboard.actions.audit")}
            </Button>
          ) : null}
          {permissions.delete ? (
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

function ArticleTag({
  children,
  className,
}: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <span
      className={cn(
        "inline-flex h-6 items-center rounded-md border px-2 text-xs font-medium",
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
    <div className="grid min-w-40 gap-1.5">
      <Label className="text-xs text-muted-foreground">{label}</Label>
      <Input
        className="h-9"
        value={value === undefined || value === null ? "" : String(value)}
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
  options: Array<{ label: string; value: string | number | boolean }>
  onChange: (value: AdminFormValue) => void
}) {
  return (
    <div className="grid min-w-40 gap-1.5">
      <Label className="text-xs text-muted-foreground">{label}</Label>
      <DashboardSelect
        value={value}
        options={options}
        placeholder={label}
        triggerClassName="h-9"
        onValueChange={(nextValue) => {
          if (nextValue === undefined) {
            onChange(undefined)
            return
          }
          const option = options.find(
            (item) => String(item.value) === nextValue
          )
          onChange(option?.value ?? nextValue)
        }}
      />
    </div>
  )
}
