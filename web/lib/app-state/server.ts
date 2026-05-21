import { cache } from "react"

import { getSiteConfig } from "@/lib/api/site"
import type { SiteConfig, UserSummary } from "@/lib/api/types"
import { getRecentUserMessages } from "@/lib/api/users"
import { getSessionUser } from "@/lib/auth/session"
import { createT, type Locale, type TFunction } from "@/lib/i18n"
import { getServerLocale } from "@/lib/i18n/server"

export type AppState = {
  config: SiteConfig | null
  currentUser: UserSummary | null
  locale: Locale
  unreadMessageCount: number
  t: TFunction
  isLogin: boolean
}

export const getAppState = cache(async (): Promise<AppState> => {
  const [config, currentUser] = await Promise.all([
    getSiteConfig().catch(() => null),
    getSessionUser().catch(() => null),
  ])
  const locale = await getServerLocale(config?.language)
  const unreadMessageCount = currentUser
    ? await getRecentUserMessages()
        .then((data) => data.count || 0)
        .catch(() => 0)
    : 0

  return {
    config,
    currentUser,
    locale,
    unreadMessageCount,
    t: createT(locale),
    isLogin: Boolean(currentUser),
  }
})
