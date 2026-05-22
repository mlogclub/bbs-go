"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardUserTaskLogsRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "userTaskLogs"),
    description: dashboardData.desc(t, "userTaskLogs"),
    listEndpoint: "/api/admin/user-task-log/list",
    viewPermission: PERMISSIONS.DASHBOARD_USER_TASK_LOG_VIEW,
    filters: [
      { name: "userId", label: dashboardData.label(t, "userId") },
      { name: "taskId", label: dashboardData.label(t, "taskId") },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "userId", label: dashboardData.label(t, "userId") },
      { key: "taskId", label: dashboardData.label(t, "taskId") },
      { key: "periodKey", label: dashboardData.label(t, "periodKey") },
      { key: "finishNo", label: dashboardData.label(t, "finishNo") },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
  }

  return <DashboardDataPage config={config} />
}
