import type { AdminFormValue } from "@/lib/api/admin"

export function createAdminInitialFilters(
  defaultFilters: Record<string, AdminFormValue> | undefined,
  limit: number
): Record<string, AdminFormValue> {
  return {
    page: 1,
    limit,
    ...(defaultFilters ?? {}),
  }
}
