import type { SiteConfig, UserSummary } from "@/lib/api/types"

export type AppLocale = "en-US" | "zh-CN"

export interface RootLoaderData {
  config: SiteConfig | null
  currentUser: UserSummary | null
  locale: AppLocale
  unreadMessageCount: number
}
