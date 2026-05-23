"use client"

import * as React from "react"
import { ShieldCheckIcon } from "lucide-react"

import {
  DashboardDataPage,
  type DashboardDataPageConfig,
} from "@/components/dashboard/data"
import { useCurrentUser } from "@/components/app/app-provider"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer"
import { ScrollArea } from "@/components/ui/scroll-area"
import {
  adminGet,
  adminPostForm,
  type AdminFormValue,
  type AdminRecord,
} from "@/lib/api/admin"
import { userHasPermission } from "@/lib/auth/roles"
import * as dashboardData from "@/components/dashboard/data/dashboard-data-route-utils"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { msgSuccess } from "@/lib/toast"

type RolePermissionRecord = AdminRecord & {
  id?: string | number
  code?: string
  name?: string
  groupName?: string
}

function selectedValues(value: unknown) {
  return Array.isArray(value) ? value.map((item) => String(item)) : []
}

function roleSubmitValues(values: Record<string, AdminFormValue>) {
  const next: Record<string, AdminFormValue> = {
    name: values.name,
    code: values.code,
    remark: values.remark,
  }
  if (values.id !== undefined && values.id !== null && values.id !== "") {
    next.id = values.id
  }
  return next
}

function isSystemRole(record: AdminRecord) {
  return Number(record.type) === 0
}

function isOwnerRole(record: AdminRecord) {
  return String(record.code || "") === "owner"
}

function canManageRole(record: AdminRecord) {
  return !isSystemRole(record)
}

function canAssignRolePermissions(record: AdminRecord) {
  return canManageRole(record) && !isOwnerRole(record)
}

export default function DashboardRolesRoute() {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const [permissionRole, setPermissionRole] =
    React.useState<AdminRecord | null>(null)
  const canAssignPermissions = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_ROLE_PERMISSION_UPDATE
  )
  const config: DashboardDataPageConfig = {
    title: dashboardData.title(t, "roles"),
    description: dashboardData.desc(t, "roles"),
    listEndpoint: "/api/admin/role/list",
    viewPermission: PERMISSIONS.DASHBOARD_ROLE_VIEW,
    defaultFilters: { status: 0 },
    listResult: "array",
    detailEndpoint: (id) => `/api/admin/role/${id}`,
    createEndpoint: "/api/admin/role/create",
    createPermission: PERMISSIONS.DASHBOARD_ROLE_CREATE,
    updateEndpoint: "/api/admin/role/update",
    updatePermission: PERMISSIONS.DASHBOARD_ROLE_UPDATE,
    deleteEndpoint: "/api/admin/role/delete",
    deletePermission: PERMISSIONS.DASHBOARD_ROLE_DELETE,
    deleteMode: "formIds",
    sortEndpoint: "/api/admin/role/update_sort",
    sortPermission: PERMISSIONS.DASHBOARD_ROLE_SORT,
    dragSort: true,
    canEdit: canManageRole,
    canDelete: canManageRole,
    renderRowActions: (record) =>
      canAssignPermissions && canAssignRolePermissions(record) ? (
        <Button
          size="sm"
          variant="outline"
          onClick={() => setPermissionRole(record)}
        >
          <ShieldCheckIcon />
          {t("dashboard.actions.assignPermissions")}
        </Button>
      ) : null,
    transformSubmitValues: roleSubmitValues,
    filters: [
      { name: "name", label: dashboardData.label(t, "name") },
      { name: "code", label: dashboardData.label(t, "code") },
      {
        name: "status",
        label: dashboardData.label(t, "status"),
        type: "select",
        options: dashboardData.enabledDisabledOptions(t),
      },
    ],
    columns: [
      { key: "id", label: dashboardData.label(t, "id") },
      { key: "name", label: dashboardData.label(t, "name") },
      { key: "code", label: dashboardData.label(t, "code") },
      {
        key: "type",
        label: dashboardData.label(t, "type"),
        render: (record) => dashboardData.roleTypeCell(t, record.type),
      },
      {
        key: "remark",
        label: dashboardData.label(t, "remark"),
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
      { name: "name", label: dashboardData.label(t, "name"), required: true },
      { name: "code", label: dashboardData.label(t, "code"), required: true },
      {
        name: "remark",
        label: dashboardData.label(t, "remark"),
        type: "textarea",
      },
    ],
  }

  return (
    <>
      <DashboardDataPage config={config} />
      <RolePermissionDrawer
        role={permissionRole}
        onOpenChange={(open) => {
          if (!open) setPermissionRole(null)
        }}
      />
    </>
  )
}

function RolePermissionDrawer({
  role,
  onOpenChange,
}: {
  role: AdminRecord | null
  onOpenChange: (open: boolean) => void
}) {
  const { t } = useI18n()
  const [permissions, setPermissions] = React.useState<RolePermissionRecord[]>(
    []
  )
  const [checkedIds, setCheckedIds] = React.useState<Set<string>>(
    () => new Set()
  )
  const [loading, setLoading] = React.useState(false)
  const [submitting, setSubmitting] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const roleId = role?.id as AdminFormValue
  const roleName = String(role?.name || role?.code || "")

  React.useEffect(() => {
    if (!roleId) return

    let cancelled = false
    setLoading(true)
    setError(null)
    void Promise.all([
      adminGet<AdminRecord>(`/api/admin/role/${roleId}`),
      adminGet<RolePermissionRecord[]>("/api/admin/role/permissions"),
    ])
      .then(([detail, permissionList]) => {
        if (cancelled) return
        setCheckedIds(new Set(selectedValues(detail.permissionIds)))
        setPermissions(Array.isArray(permissionList) ? permissionList : [])
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : t("dashboard.errors.loadFailed")
          )
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [roleId, t])

  const groups = React.useMemo(() => {
    const grouped = new Map<string, RolePermissionRecord[]>()
    permissions.forEach((permission) => {
      const groupName = String(permission.groupName || "-")
      grouped.set(groupName, [...(grouped.get(groupName) || []), permission])
    })
    return Array.from(grouped, ([name, items]) => ({ name, items }))
  }, [permissions])

  function togglePermission(id: string, checked: boolean) {
    setCheckedIds((current) => {
      const next = new Set(current)
      if (checked) {
        next.add(id)
      } else {
        next.delete(id)
      }
      return next
    })
  }

  function toggleGroup(items: RolePermissionRecord[], checked: boolean) {
    setCheckedIds((current) => {
      const next = new Set(current)
      items.forEach((item) => {
        const id = String(item.id)
        if (checked) {
          next.add(id)
        } else {
          next.delete(id)
        }
      })
      return next
    })
  }

  async function submitPermissions(event: React.FormEvent) {
    event.preventDefault()
    if (!roleId) return

    setSubmitting(true)
    setError(null)
    try {
      await adminPostForm("/api/admin/role/update_permissions", {
        id: roleId,
        permissionIds: Array.from(checkedIds),
      })
      msgSuccess(t("dashboard.messages.saved"))
      onOpenChange(false)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("dashboard.errors.saveFailed")
      )
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Drawer open={Boolean(role)} onOpenChange={onOpenChange} direction="right">
      <DrawerContent className="min-w-3xl">
        <DrawerHeader>
          <DrawerTitle>{t("dashboard.rolePermissions.title")}</DrawerTitle>
          <DrawerDescription>
            {t("dashboard.rolePermissions.description", { role: roleName })}
          </DrawerDescription>
        </DrawerHeader>
        <form
          id="role-permission-form"
          className="flex min-h-0 flex-1 flex-col"
          onSubmit={(event) => void submitPermissions(event)}
        >
          <ScrollArea className="min-h-0 flex-1 px-4">
            {error ? (
              <p className="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">
                {error}
              </p>
            ) : null}
            {loading ? (
              <p className="py-10 text-center text-muted-foreground">
                {t("dashboard.loading")}
              </p>
            ) : groups.length ? (
              <div className="space-y-3 pb-4">
                {groups.map((group) => {
                  const ids = group.items.map((item) => String(item.id))
                  const groupLabelKey = `dashboard.permissionGroups.${group.name}`
                  const groupLabel = t(groupLabelKey)
                  const selectedCount = ids.filter((id) =>
                    checkedIds.has(id)
                  ).length
                  const groupChecked =
                    selectedCount === 0
                      ? false
                      : selectedCount === ids.length
                        ? true
                        : "indeterminate"
                  return (
                    <section
                      key={group.name}
                      className="rounded-md border bg-card"
                    >
                      <div className="flex items-center gap-3 border-b bg-muted/40 px-3 py-2">
                        <Checkbox
                          checked={groupChecked}
                          onCheckedChange={(checked) =>
                            toggleGroup(group.items, checked === true)
                          }
                        />
                        <div className="min-w-0 flex-1">
                          <h3 className="truncate text-sm font-medium">
                            {groupLabel === groupLabelKey
                              ? group.name
                              : groupLabel}
                          </h3>
                          <p className="text-xs text-muted-foreground">
                            {group.name} · {selectedCount}/{group.items.length}
                          </p>
                        </div>
                      </div>
                      <div className="grid gap-2 p-3 sm:grid-cols-2">
                        {group.items.map((permission) => {
                          const id = String(permission.id)
                          const label = String(
                            permission.name || permission.code || permission.id
                          )
                          const code = permission.code
                          return (
                            <label
                              key={id}
                              className="flex min-w-0 items-start gap-2 rounded-md px-2 py-1.5 hover:bg-muted"
                            >
                              <Checkbox
                                className="mt-0.5"
                                checked={checkedIds.has(id)}
                                onCheckedChange={(checked) =>
                                  togglePermission(id, checked === true)
                                }
                              />
                              <span className="min-w-0">
                                <span className="block truncate">{label}</span>
                                {code ? (
                                  <span className="block truncate font-mono text-xs text-muted-foreground">
                                    {code}
                                  </span>
                                ) : null}
                              </span>
                            </label>
                          )
                        })}
                      </div>
                    </section>
                  )
                })}
              </div>
            ) : (
              <p className="py-10 text-center text-muted-foreground">
                {t("dashboard.rolePermissions.empty")}
              </p>
            )}
          </ScrollArea>
          <DrawerFooter className="border-t">
            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                {t("common.cancel")}
              </Button>
              <Button type="submit" disabled={loading || submitting}>
                {submitting ? t("dashboard.actions.save") : t("common.confirm")}
              </Button>
            </div>
          </DrawerFooter>
        </form>
      </DrawerContent>
    </Drawer>
  )
}
