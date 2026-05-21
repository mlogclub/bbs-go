"use client"

import * as React from "react"

import type { AdminFormValue } from "@/lib/api/admin"
import type { PermissionCode } from "@/lib/auth/permissions.generated"
import { useCurrentUser } from "@/components/app/app-provider"
import { ErrorPage } from "@/components/common/error-page"
import { ConfirmDialog } from "@/components/dashboard/confirm-dialog"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"

import { DashboardDataDetailDialog } from "./dashboard-data-detail-dialog"
import { DashboardDataFormDialog } from "./dashboard-data-form-dialog"
import { DashboardDataPasswordDialog } from "./dashboard-data-password-dialog"
import { DashboardDataTable } from "./dashboard-data-table"
import { DashboardDataToolbar } from "./dashboard-data-toolbar"
import type { DashboardDataPageConfig } from "./dashboard-data-types"
import { useDashboardDataPage } from "./use-dashboard-data-page"

export function DashboardDataPage({
  config,
}: {
  config: DashboardDataPageConfig
}) {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const formId = "dashboard-data-form"
  const canUse = (permission?: PermissionCode) =>
    !permission || userHasPermission(currentUser, permission)
  const visibleConfig = React.useMemo(
    () => ({
      ...config,
      rowActions: config.rowActions?.filter((action) =>
        canUse(action.permission)
      ),
    }),
    [config, currentUser]
  )
  const canView = canUse(config.viewPermission)
  const state = useDashboardDataPage({
    config: visibleConfig,
    messages: {
      loadFailed: t("dashboard.errors.loadFailed"),
      saveFailed: t("dashboard.errors.saveFailed"),
      deleteFailed: t("dashboard.errors.deleteFailed"),
      actionFailed: t("dashboard.errors.actionFailed"),
      sameLevelSortOnly: t("dashboard.errors.sameLevelSortOnly"),
      required: t("dashboard.errors.required"),
      invalidUrl: t("dashboard.errors.invalidUrl"),
      invalidNumber: t("dashboard.errors.invalidNumber"),
      minValue: (min) => t("dashboard.errors.minValue", { min }),
      maxValue: (max) => t("dashboard.errors.maxValue", { max }),
      saved: t("dashboard.messages.saved"),
      deleted: t("dashboard.messages.deleted"),
      actionDone: t("dashboard.messages.actionDone"),
      confirmDelete: t("dashboard.confirmDelete"),
      deleteAction: t("dashboard.actions.delete"),
    },
  })

  if (!canView) {
    return <ErrorPage statusCode={403} />
  }

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-4 md:p-6">
      <DashboardDataToolbar
        filters={visibleConfig.filters}
        values={state.filters}
        asyncOptions={state.asyncOptions}
        loading={state.loading}
        canCreate={Boolean(
          visibleConfig.createEndpoint &&
            visibleConfig.formFields?.length &&
            canUse(visibleConfig.createPermission)
        )}
        error={state.error}
        searchLabel={t("dashboard.actions.search")}
        refreshLabel={t("dashboard.actions.refresh")}
        createLabel={t("dashboard.actions.create")}
        onFilterChange={state.updateFilter}
        onRefresh={() => void state.load()}
        onCreate={state.openCreate}
      />

      <DashboardDataTable
        config={visibleConfig}
        records={state.displayRecords}
        loading={state.loading}
        page={state.page}
        pageCount={state.pageCount}
        total={state.total}
        limit={state.limit}
        labels={{
          actions: t("dashboard.actions.title"),
          loading: t("dashboard.loading"),
          noData: t("common.noData"),
          moveUp: t("dashboard.actions.moveUp"),
          moveDown: t("dashboard.actions.moveDown"),
          view: t("dashboard.actions.view"),
          edit: t("dashboard.actions.edit"),
          delete: t("dashboard.actions.delete"),
        }}
        onPageChange={(nextPage) => state.updateFilter("page", nextPage)}
        onLimitChange={(nextLimit) =>
          state.setFilters((current) => ({
            ...current,
            page: 1,
            limit: nextLimit,
          }))
        }
        onMove={(index, direction) => void state.moveRecord(index, direction)}
        canMove={state.canMoveRecord}
        canSort={canUse(visibleConfig.sortPermission)}
        canUpdate={canUse(visibleConfig.updatePermission)}
        canDelete={canUse(visibleConfig.deletePermission)}
        onRunAction={(action, record) => void state.runAction(action, record)}
        onView={(record) => void state.openView(record)}
        onEdit={(record) => void state.openEdit(record)}
        onDelete={state.requestDelete}
      />

      <DashboardDataFormDialog
        open={Boolean(state.editing)}
        formId={formId}
        title={
          state.formValues.id
            ? t("dashboard.actions.edit")
            : t("dashboard.actions.create")
        }
        fields={state.visibleFormFields}
        values={state.formValues}
        errors={state.formErrors}
        asyncOptions={state.asyncOptions}
        submitting={state.submitting}
        cancelLabel={t("common.cancel")}
        confirmLabel={t("common.confirm")}
        onOpenChange={(open) => {
          if (!open) state.setEditing(null)
        }}
        onSubmit={(event) => void state.submitForm(event)}
        onValueChange={(name, value: AdminFormValue) =>
          state.setFormValues((current) => ({
            ...current,
            [name]: value,
          }))
        }
      />

      <DashboardDataDetailDialog
        record={state.viewing}
        fields={config.detailFields}
        title={t("dashboard.actions.view")}
        cancelLabel={t("common.cancel")}
        onClose={() => state.setViewing(null)}
      />

      <DashboardDataPasswordDialog
        password={state.passwordResult}
        title={t("dashboard.resetPassword.title")}
        passwordLabel={t("dashboard.resetPassword.newPassword")}
        copyLabel={t("dashboard.resetPassword.copy")}
        copiedMessage={t("dashboard.resetPassword.copied")}
        cancelLabel={t("common.cancel")}
        onClose={() => state.setPasswordResult(null)}
      />

      <ConfirmDialog
        state={state.confirmState}
        onOpenChange={(open) => {
          if (!open) state.setConfirmState(null)
        }}
      />
    </div>
  )
}
