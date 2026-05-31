"use client"

import * as React from "react"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { HtmlImagePreview } from "@/components/common/image-preview"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { adminGet, adminPostForm, type AdminRecord } from "@/lib/api/admin"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { useCurrentUser } from "@/components/app/app-provider"
import { userHasPermission } from "@/lib/auth/roles"
import { msgError, msgSuccess } from "@/lib/toast"

type ReportDetailRecord = AdminRecord & {
  target?: Record<string, unknown>
}

export default function DashboardUserReportsRoute() {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const [processingReport, setProcessingReport] =
    React.useState<ReportDetailRecord | null>(null)
  const [processingLoading, setProcessingLoading] = React.useState(false)
  const [submittingStatus, setSubmittingStatus] = React.useState<number | null>(
    null
  )
  const [reloadKey, setReloadKey] = React.useState(0)
  const canAudit = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_USER_REPORT_AUDIT
  )

  async function openProcessReport(record: AdminRecord) {
    setProcessingLoading(true)
    try {
      const detail = await adminGet<ReportDetailRecord>(
        `/api/admin/user-report/${String(record.id)}`
      )
      setProcessingReport(detail)
    } catch (err) {
      msgError(
        err instanceof Error ? err.message : t("dashboard.errors.loadFailed")
      )
    } finally {
      setProcessingLoading(false)
    }
  }

  async function submitReportStatus(auditStatus: 1 | 2) {
    if (!processingReport?.id) return
    setSubmittingStatus(auditStatus)
    try {
      await adminPostForm("/api/admin/user-report/audit", {
        id: processingReport.id as number,
        auditStatus,
      })
      msgSuccess(
        auditStatus === 1
          ? t("dashboard.messages.reportProcessed")
          : t("dashboard.messages.reportIgnored")
      )
      setProcessingReport(null)
      setReloadKey((current) => current + 1)
    } catch (err) {
      msgError(
        err instanceof Error ? err.message : t("dashboard.errors.actionFailed")
      )
    } finally {
      setSubmittingStatus(null)
    }
  }

  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "userReports"),
    description: dashboardData.desc(t, "userReports"),
    listEndpoint: "/api/admin/user-report/list",
    viewPermission: PERMISSIONS.DASHBOARD_USER_REPORT_VIEW,
    refreshKey: reloadKey,
    filters: [
      { name: "dataId", label: dashboardData.label(t, "dataId") },
      {
        name: "dataType",
        label: dashboardData.label(t, "dataType"),
        type: "select",
        options: dashboardData.reportDataTypeOptionsFor(t),
      },
      {
        name: "auditStatus",
        label: dashboardData.label(t, "auditStatus"),
        type: "select",
        options: dashboardData.reportAuditStatusOptionsFor(t),
      },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      {
        key: "dataType",
        label: dashboardData.label(t, "dataType"),
        render: (record) =>
          dashboardData.reportDataTypeCell(t, record.dataType),
      },
      { key: "dataId", label: dashboardData.label(t, "dataId") },
      {
        key: "target",
        label: dashboardData.label(t, "reportTarget"),
        render: (record) => dashboardData.reportTargetCell(t, record),
      },
      { key: "userId", label: dashboardData.label(t, "userId") },
      {
        key: "reason",
        label: dashboardData.label(t, "reason"),
        className: "min-w-72",
      },
      {
        key: "auditStatus",
        label: dashboardData.label(t, "auditStatus"),
        render: (record) =>
          dashboardData.reportAuditStatusCell(t, record.auditStatus),
      },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
    renderRowActions: (record) =>
      canAudit && Number(record.auditStatus || 0) === 0 ? (
        <Button
          type="button"
          size="sm"
          variant="outline"
          disabled={processingLoading}
          onClick={() => void openProcessReport(record)}
        >
          {t("dashboard.reportActions.process")}
        </Button>
      ) : null,
  }

  return (
    <>
      <DashboardDataPage config={config} />
      <ReportProcessDialog
        record={processingReport}
        submittingStatus={submittingStatus}
        onClose={() => setProcessingReport(null)}
        onSubmitStatus={(status) => void submitReportStatus(status)}
      />
    </>
  )
}

function ReportProcessDialog({
  record,
  submittingStatus,
  onClose,
  onSubmitStatus,
}: {
  record: ReportDetailRecord | null
  submittingStatus: number | null
  onClose: () => void
  onSubmitStatus: (status: 1 | 2) => void
}) {
  const { t } = useI18n()
  const target =
    record?.target && !record.target.missing ? record.target : undefined
  const dataType = String(record?.dataType || "")
  const targetUrl =
    target?.url && dataType !== "comment" ? String(target.url) : undefined

  return (
    <Dialog open={Boolean(record)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>{t("dashboard.reportProcess.title")}</DialogTitle>
          <DialogDescription>
            {record
              ? `${dashboardData.reportDataTypeCell(t, record.dataType)} #${String(record.dataId || "-")}`
              : ""}
          </DialogDescription>
        </DialogHeader>

        {record ? (
          <div className="grid gap-4">
            <div className="grid gap-1.5">
              <div className="text-xs font-medium text-muted-foreground">
                {dashboardData.label(t, "reason")}
              </div>
              <div className="rounded-md border bg-muted/20 px-3 py-2 text-sm whitespace-pre-wrap">
                {String(record.reason || "-")}
              </div>
            </div>

            <div className="grid gap-1.5">
              <div className="text-xs font-medium text-muted-foreground">
                {dashboardData.label(t, "reportTarget")}
              </div>
              <div className="grid max-h-[52vh] gap-3 overflow-auto rounded-md border bg-muted/20 p-3 text-sm">
                {target ? (
                  <ReportTargetPreview target={target} />
                ) : (
                  <div className="text-muted-foreground">
                    {t("dashboard.reportProcess.targetUnavailable")}
                  </div>
                )}
              </div>
            </div>
          </div>
        ) : null}

        <DialogFooter>
          {targetUrl ? (
            <Button type="button" variant="outline" asChild>
              <a href={targetUrl} target="_blank" rel="noreferrer">
                {t("dashboard.reportActions.openTarget")}
              </a>
            </Button>
          ) : null}
          <Button
            type="button"
            variant="outline"
            disabled={Boolean(submittingStatus)}
            onClick={() => onSubmitStatus(2)}
          >
            {submittingStatus === 2
              ? t("dashboard.actions.save")
              : t("dashboard.reportActions.markIgnored")}
          </Button>
          <Button
            type="button"
            disabled={Boolean(submittingStatus)}
            onClick={() => onSubmitStatus(1)}
          >
            {submittingStatus === 1
              ? t("dashboard.actions.save")
              : t("dashboard.reportActions.markProcessed")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function ReportTargetPreview({ target }: { target: Record<string, unknown> }) {
  const title = target.title || target.nickname || target.username
  const content = target.content || target.description || target.summary
  const contentType = String(target.contentType || "")
  const meta = [
    target.username ? `@${String(target.username)}` : "",
    target.userId ? `userId: ${String(target.userId)}` : "",
    target.entityType ? `entityType: ${String(target.entityType)}` : "",
    target.entityId ? `entityId: ${String(target.entityId)}` : "",
  ].filter(Boolean)

  return (
    <>
      {title ? <div className="font-medium">{String(title)}</div> : null}
      {meta.length ? (
        <div className="text-xs text-muted-foreground">{meta.join(" · ")}</div>
      ) : null}
      {content ? (
        contentType === "html" ? (
          <HtmlImagePreview
            html={String(content)}
            className="bbs-content max-w-none break-words [&_img]:cursor-zoom-in"
          />
        ) : (
          <div className="whitespace-pre-wrap break-words">
            {String(content)}
          </div>
        )
      ) : (
        <div>-</div>
      )}
    </>
  )
}
