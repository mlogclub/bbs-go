"use client"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { useI18n } from "@/lib/i18n/provider"

export type ConfirmDialogState = {
  title?: string
  description: string
  confirmText?: string
  onConfirm: () => void
} | null

export function ConfirmDialog({
  state,
  onOpenChange,
}: {
  state: ConfirmDialogState
  onOpenChange: (open: boolean) => void
}) {
  const { t } = useI18n()

  return (
    <AlertDialog open={Boolean(state)} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            {state?.title ?? t("dashboard.confirm.title")}
          </AlertDialogTitle>
          <AlertDialogDescription>{state?.description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{t("common.cancel")}</AlertDialogCancel>
          <AlertDialogAction
            onClick={() => {
              state?.onConfirm()
              onOpenChange(false)
            }}
          >
            {state?.confirmText ?? t("dashboard.confirm.ok")}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
