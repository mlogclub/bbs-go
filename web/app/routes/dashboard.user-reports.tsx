"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardUserReportsRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "userReports"),
    description: dashboardData.desc(t, "userReports"),
    listEndpoint: "/api/admin/user-report/list",
    viewPermission: PERMISSIONS.DASHBOARD_USER_REPORT_VIEW,
    detailEndpoint: (id) => `/api/admin/user-report/${id}`,
    filters: [
      { name: "id", label: dashboardData.label(t, "id") },
      {
        name: "dataType",
        label: dashboardData.label(t, "dataType"),
        type: "select",
        options: dashboardData.reportDataTypeOptionsFor(t),
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
  }

  return <DashboardDataPage config={config} />
}
