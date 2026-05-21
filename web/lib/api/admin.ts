import { apiFetch } from "@/lib/api/client"
import {
  normalizeAdminPageResult,
  type AdminPageResult,
} from "@/lib/api/admin-page-result"

export type AdminPrimitive = string | number | boolean | null | undefined
export type AdminFormValue = AdminPrimitive | AdminPrimitive[]
export type AdminRecord = Record<string, unknown>

export type { AdminPageResult }

export function adminFormData(values: Record<string, AdminFormValue>) {
  const form = new FormData()

  Object.entries(values).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return

    if (Array.isArray(value)) {
      const items = value.filter(
        (item) => item !== undefined && item !== null && item !== ""
      )
      if (items.length > 0) {
        form.append(key, items.map((item) => String(item)).join(","))
      }
      return
    }

    form.append(key, String(value))
  })

  return form
}

export async function adminList<T = AdminRecord>(
  path: string,
  values: Record<string, AdminFormValue>
) {
  const data = await apiFetch<AdminPageResult<T> | null>(path, {
    method: "POST",
    body: adminFormData(values),
  })
  return normalizeAdminPageResult(data)
}

export function adminPostForm<T = unknown>(
  path: string,
  values: Record<string, AdminFormValue>
) {
  return apiFetch<T>(path, {
    method: "POST",
    body: adminFormData(values),
  })
}

export function adminPostJson<T = unknown>(path: string, values: unknown) {
  return apiFetch<T>(path, {
    method: "POST",
    body: values as Record<string, unknown>,
  })
}

export function adminGet<T = AdminRecord>(path: string) {
  return apiFetch<T>(path)
}

export function adminDelete<T = unknown>(
  path: string,
  values: Record<string, AdminFormValue>
) {
  return apiFetch<T>(path, {
    method: "DELETE",
    body: adminFormData(values),
  })
}
