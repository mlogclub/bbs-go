"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardBadgesRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "badges"),
    description: dashboardData.desc(t, "badges"),
    listEndpoint: "/api/admin/badge/list",
    viewPermission: PERMISSIONS.DASHBOARD_BADGE_VIEW,
    defaultFilters: { status: 0 },
    detailEndpoint: (id) => `/api/admin/badge/${id}`,
    createEndpoint: "/api/admin/badge/create",
    createPermission: PERMISSIONS.DASHBOARD_BADGE_CREATE,
    updateEndpoint: "/api/admin/badge/update",
    updatePermission: PERMISSIONS.DASHBOARD_BADGE_UPDATE,
    deleteEndpoint: "/api/admin/badge/delete",
    deletePermission: PERMISSIONS.DASHBOARD_BADGE_DELETE,
    deleteMode: "formIds",
    filters: [
      { name: "name", label: dashboardData.label(t, "name") },
      { name: "title", label: dashboardData.label(t, "title") },
      {
        name: "status",
        label: dashboardData.label(t, "status"),
        type: "select",
        options: dashboardData.normalDeletedOptions(t),
      },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "name", label: dashboardData.label(t, "name") },
      { key: "title", label: dashboardData.label(t, "title") },
      {
        key: "description",
        label: dashboardData.label(t, "description"),
        className: "min-w-72",
      },
      {
        key: "icon",
        label: dashboardData.label(t, "icon"),
        render: (record) =>
          dashboardData.imageCell(record.icon, String(record.title || "")),
      },
      { key: "sortNo", label: dashboardData.label(t, "sortNo") },
      {
        key: "status",
        label: dashboardData.label(t, "status"),
        render: (record) => dashboardData.emailStatusCell(t, record.status),
      },
      {
        key: "updateTime",
        label: dashboardData.label(t, "updateTime"),
        render: (record) => dashboardData.dateCell(record.updateTime),
      },
    ],
    formFields: [
      {
        name: "name",
        label: dashboardData.label(t, "name"),
        required: true,
        colSpan: 2,
      },
      {
        name: "title",
        label: dashboardData.label(t, "title"),
        required: true,
        colSpan: 2,
      },
      {
        name: "description",
        label: dashboardData.label(t, "description"),
        type: "textarea",
      },
      { name: "icon", label: dashboardData.label(t, "icon"), type: "image" },
      {
        name: "sortNo",
        label: dashboardData.label(t, "sortNo"),
        type: "number",
        min: 0,
        step: 1,
      },
    ],
  }

  return <DashboardDataPage config={config} />
}
