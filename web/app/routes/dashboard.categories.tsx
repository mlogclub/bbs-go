"use client"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

export default function DashboardCategoriesRoute() {
  const { t } = useI18n()
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "categories"),
    description: dashboardData.desc(t, "categories"),
    listEndpoint: "/api/admin/category/list",
    viewPermission: PERMISSIONS.DASHBOARD_CATEGORY_VIEW,
    defaultFilters: { status: 0 },
    listResult: "array",
    detailEndpoint: (id) => `/api/admin/category/${id}`,
    createEndpoint: "/api/admin/category/create",
    createPermission: PERMISSIONS.DASHBOARD_CATEGORY_CREATE,
    updateEndpoint: "/api/admin/category/update",
    updatePermission: PERMISSIONS.DASHBOARD_CATEGORY_UPDATE,
    deleteEndpoint: "/api/admin/category/delete",
    deletePermission: PERMISSIONS.DASHBOARD_CATEGORY_DELETE,
    deleteMode: "jsonIds",
    sortEndpoint: "/api/admin/category/update_sort",
    sortPermission: PERMISSIONS.DASHBOARD_CATEGORY_SORT,
    tree: true,
    treeDefaultCollapsed: true,
    treeIndentKey: "name",
    filters: [
      { name: "name", label: dashboardData.label(t, "name") },
      {
        name: "type",
        label: dashboardData.label(t, "type"),
        type: "select",
        options: dashboardData.categoryTypeOptionsFor(t),
      },
      {
        name: "categoryId",
        label: dashboardData.label(t, "category"),
        type: "select",
        optionsEndpoint: "/api/admin/category/options",
        optionLabel: dashboardData.treeOptionLabel,
        optionValue: (record) => record.id as number,
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
      { key: "name", label: dashboardData.label(t, "name") },
      {
        key: "type",
        label: dashboardData.label(t, "type"),
        render: (record) =>
          record.type === "qa"
            ? t("dashboard.categoryTypes.qa")
            : t("dashboard.categoryTypes.normal"),
      },
      {
        key: "description",
        label: dashboardData.label(t, "description"),
        className: "min-w-72",
      },
      {
        key: "logo",
        label: dashboardData.label(t, "logo"),
        render: (record) =>
          dashboardData.imageCell(record.logo, String(record.name || "")),
      },
      { key: "sortNo", label: dashboardData.label(t, "sortNo") },
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
      {
        name: "parentId",
        label: dashboardData.label(t, "parentId"),
        type: "select",
        colSpan: 2,
        optionsEndpoint: "/api/admin/category/options",
        optionLabel: dashboardData.treeOptionLabel,
        optionValue: (record) => record.id as number,
      },
      {
        name: "type",
        label: dashboardData.label(t, "type"),
        type: "select",
        required: true,
        options: dashboardData.categoryTypeOptionsFor(t),
      },
      { name: "name", label: dashboardData.label(t, "name"), required: true },
      {
        name: "description",
        label: dashboardData.label(t, "description"),
        type: "textarea",
      },
      {
        name: "logo",
        label: dashboardData.label(t, "logo"),
        type: "image",
        colSpan: 2,
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
