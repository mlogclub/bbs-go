"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardTasksRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "tasks"),
    description: dashboardData.desc(t, "tasks"),
    listEndpoint: "/api/admin/task-config/list",
    viewPermission: PERMISSIONS.DASHBOARD_TASK_VIEW,
    defaultFilters: { status: 0 },
    detailEndpoint: (id) => `/api/admin/task-config/${id}`,
    createEndpoint: "/api/admin/task-config/create",
    createPermission: PERMISSIONS.DASHBOARD_TASK_CREATE,
    updateEndpoint: "/api/admin/task-config/update",
    updatePermission: PERMISSIONS.DASHBOARD_TASK_UPDATE,
    deleteEndpoint: "/api/admin/task-config/delete",
    deletePermission: PERMISSIONS.DASHBOARD_TASK_DELETE,
    deleteMode: "formIds",
    sortEndpoint: "/api/admin/task-config/update_sort",
    sortPermission: PERMISSIONS.DASHBOARD_TASK_UPDATE,
    dragSort: true,
    filters: [
      { name: "title", label: dashboardData.label(t, "title") },
      {
        name: "groupName",
        label: dashboardData.label(t, "groupName"),
        type: "select",
        optionsEndpoint: "/api/admin/task-config/groups",
        optionLabel: (record) => String(record.name || record.key),
        optionValue: (record) => record.key as string,
      },
      {
        name: "eventType",
        label: dashboardData.label(t, "eventType"),
        type: "select",
        optionsEndpoint: "/api/admin/common/task_event_types",
        optionLabel: (record) => String(record.title || record.value),
        optionValue: (record) => record.value as string,
      },
      {
        name: "period",
        label: dashboardData.label(t, "period"),
        type: "select",
        options: dashboardData.taskPeriodOptionsFor(t),
      },
      {
        name: "status",
        label: dashboardData.label(t, "status"),
        type: "select",
        options: dashboardData.normalDeletedOptions(t),
      },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "groupName", label: dashboardData.label(t, "groupName") },
      { key: "title", label: dashboardData.label(t, "title") },
      { key: "eventType", label: dashboardData.label(t, "eventType") },
      { key: "period", label: dashboardData.label(t, "period") },
      { key: "score", label: dashboardData.label(t, "score") },
      { key: "exp", label: dashboardData.label(t, "exp") },
      { key: "badgeId", label: dashboardData.label(t, "badgeId") },
      { key: "eventCount", label: dashboardData.label(t, "eventCount") },
      {
        key: "maxFinishCount",
        label: dashboardData.label(t, "maxFinishCount"),
      },
      {
        key: "status",
        label: dashboardData.label(t, "status"),
        render: (record) => dashboardData.disabledStatusCell(t, record.status),
      },
      {
        key: "updateTime",
        label: dashboardData.label(t, "updateTime"),
        render: (record) => dashboardData.dateCell(record.updateTime),
      },
    ],
    formFields: [
      {
        name: "groupName",
        label: dashboardData.label(t, "groupName"),
        type: "select",
        optionsEndpoint: "/api/admin/task-config/groups",
        optionLabel: (record) => String(record.name || record.key),
        optionValue: (record) => record.key as string,
      },
      {
        name: "eventType",
        label: dashboardData.label(t, "eventType"),
        type: "select",
        optionsEndpoint: "/api/admin/common/task_event_types",
        optionLabel: (record) => String(record.title || record.value),
        optionValue: (record) => record.value as string,
      },
      { name: "title", label: dashboardData.label(t, "title"), required: true },
      {
        name: "description",
        label: dashboardData.label(t, "description"),
        type: "textarea",
      },
      {
        name: "score",
        label: dashboardData.label(t, "score"),
        type: "number",
        min: 0,
        step: 1,
      },
      {
        name: "exp",
        label: dashboardData.label(t, "exp"),
        type: "number",
        min: 0,
        step: 1,
      },
      {
        name: "badgeId",
        label: dashboardData.label(t, "badgeId"),
        type: "select",
        options: [{ label: t("dashboard.task.noBadge"), value: 0 }],
        optionsEndpoint: "/api/admin/badge/list",
        optionLabel: (record) =>
          String(record.title || record.name || record.id),
        optionValue: (record) => record.id as number,
      },
      {
        name: "period",
        label: dashboardData.label(t, "period"),
        type: "select",
        required: true,
        options: dashboardData.taskPeriodOptionsFor(t),
      },
      {
        name: "eventCount",
        label: dashboardData.label(t, "eventCount"),
        type: "number",
        min: 1,
        step: 1,
      },
      {
        name: "maxFinishCount",
        label: dashboardData.label(t, "maxFinishCount"),
        type: "number",
        min: 1,
        step: 1,
      },
      { name: "btnName", label: dashboardData.label(t, "btnName") },
      { name: "actionUrl", label: dashboardData.label(t, "actionUrl") },
    ],
  }

  return <DashboardDataPage config={config} />
}
