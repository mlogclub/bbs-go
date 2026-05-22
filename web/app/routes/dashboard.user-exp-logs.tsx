"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardUserExpLogsRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "userExpLogs"),
    description: dashboardData.desc(t, "userExpLogs"),
    listEndpoint: "/api/admin/user-exp-log/list",
    viewPermission: PERMISSIONS.DASHBOARD_USER_EXP_LOG_VIEW,
    filters: [
      { name: "userId", label: dashboardData.label(t, "userId") },
      { name: "sourceType", label: dashboardData.label(t, "sourceType") },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "userId", label: dashboardData.label(t, "userId") },
      { key: "sourceType", label: dashboardData.label(t, "sourceType") },
      { key: "sourceId", label: dashboardData.label(t, "sourceId") },
      { key: "score", label: dashboardData.label(t, "score") },
      { key: "exp", label: dashboardData.label(t, "exp") },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
  }

  return <DashboardDataPage config={config} />
}
