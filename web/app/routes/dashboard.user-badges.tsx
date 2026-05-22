"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardUserBadgesRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "userBadges"),
    description: dashboardData.desc(t, "userBadges"),
    listEndpoint: "/api/admin/user-badge/list",
    viewPermission: PERMISSIONS.DASHBOARD_USER_BADGE_VIEW,
    filters: [
      { name: "userId", label: dashboardData.label(t, "userId") },
      {
        name: "badgeId",
        label: dashboardData.label(t, "badgeId"),
        type: "select",
        optionsEndpoint: "/api/admin/badge/list",
        optionLabel: (record) =>
          String(record.title || record.name || record.id),
        optionValue: (record) => record.id as number,
      },
      { name: "sourceType", label: dashboardData.label(t, "sourceType") },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "userId", label: dashboardData.label(t, "userId") },
      { key: "badgeId", label: dashboardData.label(t, "badgeId") },
      {
        key: "icon",
        label: dashboardData.label(t, "badgeIcon"),
        render: (record) =>
          dashboardData.imageCell(record.icon, String(record.badgeId || "")),
      },
      { key: "sourceType", label: dashboardData.label(t, "sourceType") },
      { key: "sourceId", label: dashboardData.label(t, "sourceId") },
      {
        key: "isWorn",
        label: dashboardData.label(t, "isWorn"),
        render: (record) =>
          record.isWorn
            ? t("dashboard.boolean.yes")
            : t("dashboard.boolean.no"),
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
