"use client"

import * as React from "react"
import { useRouter } from "@/lib/router/navigation"
import { MoreVertical } from "lucide-react"

import {
  ConfirmDialog,
  type ConfirmDialogState,
} from "@/components/common/confirm-dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { apiFetch, toFormData } from "@/lib/api/client"
import type { Article, UserSummary } from "@/lib/api/types"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"
import { msg, useToastActions } from "@/lib/toast"

export function ArticleManageMenu({
  article,
  currentUser,
}: {
  article: Article
  currentUser?: UserSummary | null
}) {
  const router = useRouter()
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [confirmState, setConfirmState] =
    React.useState<ConfirmDialogState>(null)
  const isArticleOwner = Boolean(
    currentUser && currentUser.id === article.user.id
  )
  const canEdit =
    isArticleOwner ||
    userHasPermission(currentUser, PERMISSIONS.DASHBOARD_ARTICLE_UPDATE)
  const canDelete =
    isArticleOwner ||
    userHasPermission(currentUser, PERMISSIONS.DASHBOARD_ARTICLE_DELETE)
  const canForbidden = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_FORBIDDEN
  )
  const canForbiddenForever = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_FORBIDDEN_FOREVER
  )

  if (!canEdit && !canDelete && !canForbidden && !canForbiddenForever) {
    return null
  }

  async function forbidden(days: number) {
    try {
      await apiFetch<null>("/api/user/forbidden", {
        method: "POST",
        body: toFormData({ userId: article.user.id, days }),
      })
      msg({ message: t("component.articleManageMenu.forbiddenSuccess") })
    } catch (error) {
      catchError(error)
    }
  }

  async function deleteArticle() {
    try {
      await apiFetch<null>(`/api/article/delete/${article.id}`, {
        method: "POST",
      })
      msg({
        message: t("component.articleManageMenu.deleteSuccess"),
        onClose() {
          router.push("/articles")
        },
      })
    } catch (error) {
      catchError(error)
    }
  }

  function confirmDeleteArticle() {
    setConfirmState({
      description: t("component.articleManageMenu.confirmDelete"),
      confirmText: t("component.articleManageMenu.delete"),
      onConfirm: () => {
        void deleteArticle()
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
              onSelect={() => router.push(`/article/edit/${article.id}`)}
            >
              {t("component.articleManageMenu.edit")}
            </DropdownMenuItem>
          ) : null}
          {canDelete ? (
            <DropdownMenuItem onSelect={confirmDeleteArticle}>
              {t("component.articleManageMenu.delete")}
            </DropdownMenuItem>
          ) : null}
          {canForbidden && (canEdit || canDelete) ? (
            <DropdownMenuSeparator />
          ) : null}
          {canForbidden ? (
            <DropdownMenuItem onSelect={() => void forbidden(7)}>
              {t("component.articleManageMenu.forbidden7Days")}
            </DropdownMenuItem>
          ) : null}
          {canForbiddenForever ? (
            <DropdownMenuItem onSelect={() => void forbidden(-1)}>
              {t("component.articleManageMenu.forbiddenForever")}
            </DropdownMenuItem>
          ) : null}
        </DropdownMenuContent>
      </DropdownMenu>
      <ConfirmDialog
        state={confirmState}
        onOpenChange={(open) => {
          if (!open) setConfirmState(null)
        }}
      />
    </>
  )
}
