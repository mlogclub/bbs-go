import assert from "node:assert/strict"

import { createAdminInitialFilters } from "../lib/dashboard/default-filters.ts"

const initialFilters = createAdminInitialFilters({ status: 0 }, 20)

assert.deepEqual(initialFilters, {
  page: 1,
  limit: 20,
  status: 0,
})

initialFilters.status = 1

assert.deepEqual(createAdminInitialFilters({ status: 0 }, 50), {
  page: 1,
  limit: 50,
  status: 0,
})

console.log("admin default filters tests passed")
