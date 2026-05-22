import assert from "node:assert/strict"

import { normalizeAdminPageResult } from "../lib/api/admin-page-result.ts"

assert.deepEqual(normalizeAdminPageResult(null), {
  results: [],
  page: undefined,
})

assert.deepEqual(normalizeAdminPageResult({ results: null }), {
  results: [],
  page: undefined,
})

assert.deepEqual(
  normalizeAdminPageResult({
    results: [{ id: 1 }],
    page: { page: 1, limit: 20, total: 1 },
  }),
  {
    results: [{ id: 1 }],
    page: { page: 1, limit: 20, total: 1 },
  }
)

console.log("admin page result tests passed")
