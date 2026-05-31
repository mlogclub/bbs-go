"use client"

import * as React from "react"
import { useRouter } from "@/lib/router/navigation"
import { MoreVertical } from "lucide-react"

import {
  ConfirmDialog,
  type ConfirmDialogState,
} from "@/components/common/confirm-dialog"
import { UserReportDialog } from "@/components/common/user-report-dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { apiFetch, toFormData } from "@/lib/api/client"
import type { Topic, UserSummary } from "@/lib/api/types"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"
import { msg, useToastActions } from "@/lib/toast"

function actionText(text: string, action: string) {
  return text.replace("{action}", action)
}

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
      <UserReportDialog
        open={reportOpen}
        dataId={topic.id}
        dataType="topic"
        onOpenChange={setReportOpen}
      />
      <ConfirmDialog
        state={confirmState}
        onOpenChange={(open) => {
          if (!open) setConfirmState(null)
        }}
      />
    </>
  )
}
