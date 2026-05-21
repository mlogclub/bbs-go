import type { UserSummary } from "@/lib/api/types"
import { PERMISSIONS, type PermissionCode } from "@/lib/auth/permissions.generated"

export function userHasRole(
  user: UserSummary | null | undefined,
  role: string
) {
  return Boolean(user?.roles?.includes(role))
}

export function userIsOwner(user: UserSummary | null | undefined) {
  return userHasRole(user, "owner")
}

export function userHasPermission(
  user: UserSummary | null | undefined,
  permission: PermissionCode | null | undefined
) {
  if (!permission) return true
  if (userIsOwner(user)) return true
  return Boolean(user?.permissions?.includes(permission))
}

export function userHasAnyPermission(
  user: UserSummary | null | undefined,
  permissions: PermissionCode[]
) {
  if (userIsOwner(user)) return true
  return permissions.some((permission) => userHasPermission(user, permission))
}

export function userCanAccessDashboard(user: UserSummary | null | undefined) {
  return userHasPermission(user, PERMISSIONS.DASHBOARD_VIEW)
}
