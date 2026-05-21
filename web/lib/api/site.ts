import { serverApiFetch as apiFetch } from "./server"

import type { SiteConfig } from "./types"

export function getSiteConfig() {
  return apiFetch<SiteConfig>("/api/config/configs")
}
