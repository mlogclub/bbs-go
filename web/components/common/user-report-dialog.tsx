"use client"

import * as React from "react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Textarea } from "@/components/ui/textarea"
import { apiFetch, toFormData } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"
import { msg, useToastActions } from "@/lib/toast"

const reportReasonKeys = [
  "spam",
  "illegal",
  "harassment",
  "pornographic",
  "other",
] as const

type ReportReasonKey = (typeof reportReasonKeys)[number]

export function UserReportDialog({
  open,
  dataId,
  dataType,
  onOpenChange,
}: {
  open: boolean
  dataId: string | number
  dataType: string
  onOpenChange: (open: boolean) => void
}) {
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [reportReason, setReportReason] =
    React.useState<ReportReasonKey>("spam")
  const [reportDetail, setReportDetail] = React.useState("")
  const [reporting, setReporting] = React.useState(false)

  async function submitReport(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (reporting) return

    const reasonLabel = t(`component.userReport.reason.${reportReason}`)
    const detail = reportDetail.trim()
    const reason = detail ? `${reasonLabel}: ${detail}` : reasonLabel

    try {
      setReporting(true)
      await apiFetch<null>("/api/user-report/submit", {
        method: "POST",
        body: toFormData({
          dataId,
          dataType,
          reason,
        }),
      })
      onOpenChange(false)
      setReportReason("spam")
      setReportDetail("")
      msg({ message: t("component.userReport.success") })
    } catch (error) {
      catchError(error)
    } finally {
      setReporting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form className="grid gap-5" onSubmit={submitReport}>
          <DialogHeader>
            <DialogTitle>{t("component.userReport.title")}</DialogTitle>
            <DialogDescription>
              {t("component.userReport.description")}
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
                {t(`component.userReport.reason.${reason}`)}
              </button>
            ))}
          </div>
          <Textarea
            value={reportDetail}
            onChange={(event) => setReportDetail(event.target.value)}
            maxLength={500}
            placeholder={t("component.userReport.placeholder")}
          />
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              {t("common.cancel")}
            </Button>
            <Button type="submit" disabled={reporting}>
              {reporting
                ? t("component.userReport.submitting")
                : t("component.userReport.submit")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
