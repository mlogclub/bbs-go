"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardForbiddenWordsRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "forbiddenWords"),
    description: dashboardData.desc(t, "forbiddenWords"),
    listEndpoint: "/api/admin/forbidden-word/list",
    viewPermission: PERMISSIONS.DASHBOARD_FORBIDDEN_WORD_VIEW,
    detailEndpoint: (id) => `/api/admin/forbidden-word/${id}`,
    createEndpoint: "/api/admin/forbidden-word/create",
    createPermission: PERMISSIONS.DASHBOARD_FORBIDDEN_WORD_CREATE,
    updateEndpoint: "/api/admin/forbidden-word/update",
    updatePermission: PERMISSIONS.DASHBOARD_FORBIDDEN_WORD_UPDATE,
    deleteEndpoint: "/api/admin/forbidden-word/delete",
    deletePermission: PERMISSIONS.DASHBOARD_FORBIDDEN_WORD_DELETE,
    deleteMode: "formIds",
    filters: [
      {
        name: "type",
        label: dashboardData.label(t, "type"),
        type: "select",
        options: [
          { label: t("dashboard.forbiddenWordTypes.word"), value: "word" },
          { label: t("dashboard.forbiddenWordTypes.regex"), value: "regex" },
        ],
      },
      { name: "word", label: dashboardData.label(t, "word") },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      {
        key: "type",
        label: dashboardData.label(t, "type"),
        render: (record) =>
          record.type === "regex"
            ? t("dashboard.forbiddenWordTypes.regex")
            : t("dashboard.forbiddenWordTypes.word"),
      },
      { key: "word", label: dashboardData.label(t, "word") },
      {
        key: "remark",
        label: dashboardData.label(t, "remark"),
        className: "min-w-72",
      },
      {
        key: "createTime",
        label: dashboardData.label(t, "createTime"),
        render: (record) => dashboardData.dateCell(record.createTime),
      },
    ],
    formFields: [
      { name: "id", label: dashboardData.label(t, "id"), type: "number" },
      {
        name: "type",
        label: dashboardData.label(t, "type"),
        type: "select",
        required: true,
        options: [
          { label: t("dashboard.forbiddenWordTypes.word"), value: "word" },
          { label: t("dashboard.forbiddenWordTypes.regex"), value: "regex" },
        ],
      },
      { name: "word", label: dashboardData.label(t, "word"), required: true },
      {
        name: "remark",
        label: dashboardData.label(t, "remark"),
        type: "textarea",
      },
    ],
  }

  return <DashboardDataPage config={config} />
}
