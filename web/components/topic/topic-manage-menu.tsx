"use client"

import * as React from "react"
import { useRouter } from "@/lib/router/navigation"
import { MoreVertical } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  ConfirmDialog,
  type ConfirmDialogState,
} from "@/components/common/confirm-dialog"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Textarea } from "@/components/ui/textarea"
import { apiFetch, toFormData } from "@/lib/api/client"
import type { Topic, UserSummary } from "@/lib/api/types"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"
import { msg, useToastActions } from "@/lib/toast"

function actionText(text: string, action: string) {
  return text.replace("{action}", action)
}

const reportReasonKeys = [
  "spam",
  "illegal",
  "harassment",
  "pornographic",
  "other",
] as const

type ReportReasonKey = (typeof reportReasonKeys)[number]

export function TopicManageMenu({
  topic,
  currentUser,
}: {
  topic: Topic
  currentUser?: UserSummary | null
}) {
  const router = useRouter()
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [recommend, setRecommend] = React.useState(Boolean(topic.recommend))
  const [sticky, setSticky] = React.useState(Boolean(topic.sticky))
  const [confirmState, setConfirmState] =
    React.useState<ConfirmDialogState>(null)
  const [reportOpen, setReportOpen] = React.useState(false)
  const [reportReason, setReportReason] =
    React.useState<ReportReasonKey>("spam")
  const [reportDetail, setReportDetail] = React.useState("")
  const [reporting, setReporting] = React.useState(false)
  const isTopicOwner = Boolean(currentUser && currentUser.id === topic.user.id)
  const canReport = Boolean(currentUser && !isTopicOwner)
  const canEdit = isTopicOwner && topic.type === 0
  const canDelete =
    isTopicOwner ||
    userHasPermission(currentUser, PERMISSIONS.DASHBOARD_TOPIC_DELETE)
  const canRecommend = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_TOPIC_RECOMMEND
  )
  const canSticky = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_TOPIC_STICKY
  )
  const canForbidden = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_FORBIDDEN
  )
  const canForbiddenForever = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_FORBIDDEN_FOREVER
  )
  const canManage =
    canRecommend || canSticky || canForbidden || canForbiddenForever

  if (!canEdit && !canDelete && !canManage && !canReport) {
    return null
  }

  async function forbidden(days: number) {
    try {
      await apiFetch<null>("/api/user/forbidden", {
        method: "POST",
        body: toFormData({ userId: topic.user.id, days }),
      })
      msg({ message: t("component.topicManageMenu.forbiddenSuccess") })
    } catch (error) {
      catchError(error)
    }
  }

  async function deleteTopic() {
    try {
      await apiFetch<null>(`/api/topic/delete/${topic.id}`, {
        method: "POST",
      })
      msg({
        message: t("component.topicManageMenu.deleteSuccess"),
        onClose() {
          router.push("/topics")
        },
      })
    } catch (error) {
      catchError(error)
    }
  }

  function confirmDeleteTopic() {
    setConfirmState({
      description: t("component.topicManageMenu.confirmDelete"),
      confirmText: t("component.topicManageMenu.delete"),
      onConfirm: () => {
        void deleteTopic()
      },
    })
  }

  async function switchRecommend() {
    const action = recommend
      ? t("component.topicManageMenu.cancelRecommend")
      : t("component.topicManageMenu.recommend")
    try {
      const next = !recommend
      await apiFetch<null>(`/api/topic/recommend/${topic.id}`, {
        method: "POST",
        body: toFormData({ recommend: next }),
      })
      setRecommend(next)
      msg({
        message: actionText(
          t("component.topicManageMenu.actionSuccess"),
          action
        ),
      })
    } catch (error) {
      catchError(error)
    }
  }

  function confirmSwitchRecommend() {
    const action = recommend
      ? t("component.topicManageMenu.cancelRecommend")
      : t("component.topicManageMenu.recommend")
    setConfirmState({
      description: actionText(
        t("component.topicManageMenu.confirmAction"),
        action
      ),
      confirmText: action,
      onConfirm: () => {
        void switchRecommend()
      },
    })
  }

  async function switchSticky() {
    const action = sticky
      ? t("component.topicManageMenu.cancelSticky")
      : t("component.topicManageMenu.sticky")
    try {
      const next = !sticky
      await apiFetch<null>(`/api/topic/sticky/${topic.id}`, {
        method: "POST",
        body: toFormData({ sticky: next }),
      })
      setSticky(next)
      msg({
        message: actionText(
          t("component.topicManageMenu.actionSuccess"),
          action
        ),
      })
    } catch (error) {
      catchError(error)
    }
  }

  function confirmSwitchSticky() {
    const action = sticky
      ? t("component.topicManageMenu.cancelSticky")
      : t("component.topicManageMenu.sticky")
    setConfirmState({
      description: actionText(
        t("component.topicManageMenu.confirmAction"),
        action
      ),
      confirmText: action,
      onConfirm: () => {
        void switchSticky()
      },
    })
  }

  async function submitReport(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (reporting) return

    const reasonLabel = t(
      `component.topicManageMenu.reportReason.${reportReason}`
    )
    const detail = reportDetail.trim()
    const reason = detail ? `${reasonLabel}: ${detail}` : reasonLabel

    try {
      setReporting(true)
      await apiFetch<null>("/api/user-report/submit", {
        method: "POST",
        body: toFormData({
          dataId: topic.id,
          dataType: "topic",
          reason,
        }),
      })
      setReportOpen(false)
      setReportReason("spam")
      setReportDetail("")
      msg({ message: t("component.topicManageMenu.reportSuccess") })
    } catch (error) {
      catchError(error)
    } finally {
      setReporting(false)
    }
  }

  return (
    <>
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <button
            type="button"
            className="inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
            aria-label={t("common.moreActions")}
            title={t("common.moreActions")}
          >
            <MoreVertical className="h-4 w-4" />
            <span className="sr-only">{t("common.moreActions")}</span>
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="min-w-40">
          {canEdit ? (
            <DropdownMenuItem
              onSelect={() => router.push(`/topic/edit/${topic.id}`)}
            >
              {t("component.topicManageMenu.edit")}
            </DropdownMenuItem>
          ) : null}
          {canDelete ? (
            <DropdownMenuItem onSelect={confirmDeleteTopic}>
              {t("component.topicManageMenu.delete")}
            </DropdownMenuItem>
          ) : null}
          {canManage && (canEdit || canDelete) ? (
            <DropdownMenuSeparator />
          ) : null}
          {canRecommend ? (
            <DropdownMenuItem onSelect={confirmSwitchRecommend}>
              {recommend
                ? t("component.topicManageMenu.cancelRecommend")
                : t("component.topicManageMenu.recommend")}
            </DropdownMenuItem>
          ) : null}
          {canSticky ? (
            <DropdownMenuItem onSelect={confirmSwitchSticky}>
              {sticky
                ? t("component.topicManageMenu.cancelSticky")
                : t("component.topicManageMenu.sticky")}
            </DropdownMenuItem>
          ) : null}
          {canForbidden &&
          (canEdit || canDelete || canRecommend || canSticky) ? (
            <DropdownMenuSeparator />
          ) : null}
          {canForbidden ? (
            <DropdownMenuItem onSelect={() => void forbidden(7)}>
              {t("component.topicManageMenu.forbidden7Days")}
            </DropdownMenuItem>
          ) : null}
          {canForbiddenForever ? (
            <DropdownMenuItem onSelect={() => void forbidden(-1)}>
              {t("component.topicManageMenu.forbiddenForever")}
            </DropdownMenuItem>
          ) : null}
          {canReport &&
          (canEdit ||
            canDelete ||
            canRecommend ||
            canSticky ||
            canForbidden ||
            canForbiddenForever) ? (
            <DropdownMenuSeparator />
          ) : null}
          {canReport ? (
            <DropdownMenuItem onSelect={() => setReportOpen(true)}>
              {t("component.topicManageMenu.report")}
            </DropdownMenuItem>
          ) : null}
        </DropdownMenuContent>
      </DropdownMenu>
      <Dialog open={reportOpen} onOpenChange={setReportOpen}>
        <DialogContent>
          <form className="grid gap-5" onSubmit={submitReport}>
            <DialogHeader>
              <DialogTitle>{t("component.topicManageMenu.report")}</DialogTitle>
              <DialogDescription>
                {t("component.topicManageMenu.reportDescription")}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-2">
              {reportReasonKeys.map((reason) => (
                <button
                  key={reason}
                  type="button"
                  className={`flex min-h-9 items-center rounded-md border px-3 text-left text-sm transition-colors ${
                    reportReason === reason
                      ? "border-primary bg-primary/10 text-foreground"
                      : "border-border text-muted-foreground hover:bg-muted hover:text-foreground"
                  }`}
                  onClick={() => setReportReason(reason)}
                >
                  {t(`component.topicManageMenu.reportReason.${reason}`)}
                </button>
              ))}
            </div>
            <Textarea
              value={reportDetail}
              onChange={(event) => setReportDetail(event.target.value)}
              maxLength={500}
              placeholder={t("component.topicManageMenu.reportPlaceholder")}
            />
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setReportOpen(false)}
              >
                {t("common.cancel")}
              </Button>
              <Button type="submit" disabled={reporting}>
                {reporting
                  ? t("component.topicManageMenu.reporting")
                  : t("component.topicManageMenu.reportSubmit")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
      <ConfirmDialog
        state={confirmState}
        onOpenChange={(open) => {
          if (!open) setConfirmState(null)
        }}
      />
    </>
  )
}
