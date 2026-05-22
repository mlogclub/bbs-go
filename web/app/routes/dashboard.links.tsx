"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardLinksRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "links"),
    description: dashboardData.desc(t, "links"),
    listEndpoint: "/api/admin/link/list",
    viewPermission: PERMISSIONS.DASHBOARD_LINK_VIEW,
    defaultFilters: { status: 0 },
    detailEndpoint: (id) => `/api/admin/link/${id}`,
    createEndpoint: "/api/admin/link/create",
    createPermission: PERMISSIONS.DASHBOARD_LINK_CREATE,
    updateEndpoint: "/api/admin/link/update",
    updatePermission: PERMISSIONS.DASHBOARD_LINK_UPDATE,
    filters: [
      { name: "title", label: dashboardData.label(t, "title") },
      { name: "url", label: dashboardData.label(t, "url") },
      {
        name: "status",
        label: dashboardData.label(t, "status"),
        type: "select",
        options: dashboardData.normalDeletedOptions(t),
      },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "title", label: dashboardData.label(t, "title") },
      {
        key: "url",
        label: dashboardData.label(t, "url"),
        className: "min-w-72",
      },
      {
        key: "summary",
        label: dashboardData.label(t, "summary"),
        className: "min-w-72",
      },
      {
        key: "status",
        label: dashboardData.label(t, "status"),
        render: (record) => dashboardData.statusCell(t, record.status),
      },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
    formFields: [
      { name: "id", label: dashboardData.label(t, "id"), type: "number" },
      { name: "title", label: dashboardData.label(t, "title"), required: true },
      {
        name: "url",
        label: dashboardData.label(t, "url"),
        type: "url",
        required: true,
      },
      {
        name: "summary",
        label: dashboardData.label(t, "summary"),
        type: "textarea",
      },
      {
        name: "status",
        label: dashboardData.label(t, "status"),
        type: "select",
        required: true,
        options: dashboardData.normalDeletedOptions(t),
      },
    ],
  }

  return <DashboardDataPage config={config} />
}
