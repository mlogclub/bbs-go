import type { SiteConfig, UserSummary } from "@/lib/api/types"

import type { Locale } from "@/lib/i18n"

export type AppLocale = Locale

export interface RootLoaderData {
  config: SiteConfig | null
  currentUser: UserSummary | null
  locale: AppLocale
  unreadMessageCount: number
}
