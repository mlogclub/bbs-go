"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardEmailLogsRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "emailLogs"),
    description: dashboardData.desc(t, "emailLogs"),
    listEndpoint: "/api/admin/email-log/list",
    viewPermission: PERMISSIONS.DASHBOARD_EMAIL_LOG_VIEW,
    defaultFilters: { status: 0 },
    detailEndpoint: (id) => `/api/admin/email-log/${id}`,
    filters: [
      { name: "toEmail", label: dashboardData.label(t, "toEmail") },
      { name: "bizType", label: dashboardData.label(t, "bizType") },
      {
        name: "status",
        label: dashboardData.label(t, "status"),
        type: "select",
        options: dashboardData.normalDeletedOptions(t),
      },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "toEmail", label: dashboardData.label(t, "toEmail") },
      { key: "bizType", label: dashboardData.label(t, "bizType") },
      {
        key: "subject",
        label: dashboardData.label(t, "subject"),
        className: "min-w-72",
      },
      {
        key: "status",
        label: dashboardData.label(t, "status"),
        render: (record) => dashboardData.statusCell(t, record.status),
      },
      {
        key: "errorMsg",
        label: dashboardData.label(t, "errorMsg"),
        className: "min-w-72",
      },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
    detailFields: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "toEmail", label: dashboardData.label(t, "toEmail") },
      { key: "bizType", label: dashboardData.label(t, "bizType") },
      { key: "subject", label: dashboardData.label(t, "subject") },
      {
        key: "status",
        label: dashboardData.label(t, "status"),
        render: (record) => dashboardData.statusCell(t, record.status),
      },
      { key: "errorMsg", label: dashboardData.label(t, "errorMsg") },
      {
        key: "content",
        label: dashboardData.label(t, "content"),
        render: (record) => dashboardData.codeBlock(record.content),
      },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
  }

  return <DashboardDataPage config={config} />
}
