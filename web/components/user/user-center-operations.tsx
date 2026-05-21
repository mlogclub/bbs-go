"use client"

import * as React from "react"
import { Ban, CircleCheck, ShieldAlert } from "lucide-react"

import {
  ConfirmDialog,
  type ConfirmDialogState,
} from "@/components/common/confirm-dialog"
import { WidgetCard } from "@/components/common/widget-card"
import { apiFetch, toFormData } from "@/lib/api/client"
import type { UserSummary } from "@/lib/api/types"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"
import { msg, useToastActions } from "@/lib/toast"

export function UserCenterOperations({
  user,
  currentUser,
}: {
  user: UserSummary
  currentUser?: UserSummary | null
}) {
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [forbidden, setForbidden] = React.useState(Boolean(user.forbidden))
  const [confirmState, setConfirmState] =
    React.useState<ConfirmDialogState>(null)
  const canForbidden = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_FORBIDDEN
  )
  const canForbiddenForever = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_FORBIDDEN_FOREVER
  )

  if (!canForbidden && !canForbiddenForever) {
    return null
  }

  async function updateForbidden(days: number) {
    try {
      await apiFetch<null>("/api/user/forbidden", {
        method: "POST",
        body: toFormData({ userId: user.id, days }),
      })
      setForbidden(days !== 0)
      msg({
        message:
          days === 0
            ? t("component.userCenterSidebar.removeForbiddenSuccess")
            : t("component.userCenterSidebar.forbiddenSuccess"),
      })
    } catch (error) {
      catchError(error)
    }
  }

  function confirmUpdateForbidden(days: number) {
    if (days === 0) {
      void updateForbidden(days)
      return
    }

    setConfirmState({
      description:
        days > 0
          ? t("component.userCenterSidebar.confirmForbidden")
          : t("component.userCenterSidebar.confirmForbiddenForever"),
      confirmText:
        days > 0
          ? t("component.userCenterSidebar.forbidden7Days")
          : t("component.userCenterSidebar.forbiddenForever"),
      onConfirm: () => {
        void updateForbidden(days)
      },
    })
  }

  return (
    <>
      <WidgetCard title={t("component.userCenterSidebar.operations")}>
        <ul className="list-none space-y-2 text-sm [&_li]:hover:cursor-pointer [&_li]:hover:bg-amber-50 [&_li]:hover:text-amber-800">
          {forbidden ? (
            <li className="flex items-center gap-2">
              <CircleCheck className="shrink-0" size={14} aria-hidden="true" />
              <button
                type="button"
                className="text-primary"
                onClick={() => confirmUpdateForbidden(0)}
              >
                {t("component.userCenterSidebar.removeForbidden")}
              </button>
            </li>
          ) : (
            <>
              {canForbidden ? (
                <li className="flex items-center gap-2">
                  <Ban className="shrink-0" size={14} aria-hidden="true" />
                  <button
                    type="button"
                    className="text-primary"
                    onClick={() => confirmUpdateForbidden(7)}
                  >
                    {t("component.userCenterSidebar.forbidden7Days")}
                  </button>
                </li>
              ) : null}
              {canForbiddenForever ? (
                <li className="flex items-center gap-2">
                  <ShieldAlert
                    className="shrink-0"
                    size={14}
                    aria-hidden="true"
                  />
                  <button
                    type="button"
                    className="text-primary"
                    onClick={() => confirmUpdateForbidden(-1)}
                  >
                    {t("component.userCenterSidebar.forbiddenForever")}
                  </button>
                </li>
              ) : null}
            </>
          )}
        </ul>
      </WidgetCard>
      <ConfirmDialog
        state={confirmState}
        onOpenChange={(open) => {
          if (!open) setConfirmState(null)
        }}
      />
    </>
  )
}
