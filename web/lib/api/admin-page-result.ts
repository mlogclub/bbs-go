export interface AdminPageResult<T = Record<string, unknown>> {
  results: T[]
  page?: {
    page?: number
    limit?: number
    total?: number
  }
}

export function normalizeAdminPageResult<T = Record<string, unknown>>(
  value: { results?: T[] | null; page?: AdminPageResult<T>["page"] } | null
): AdminPageResult<T> {
  return {
    results: Array.isArray(value?.results) ? value.results : [],
    page: value?.page,
  }
}
