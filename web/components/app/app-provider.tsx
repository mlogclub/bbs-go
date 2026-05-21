"use client"

import * as React from "react"

import { loadClientAppState } from "@/lib/app-state/client"
import type { SiteConfig, UserSummary } from "@/lib/api/types"
import type { Locale } from "@/lib/i18n"

export type ClientAppState = {
  config: SiteConfig | null
  currentUser: UserSummary | null
  locale: Locale
  unreadMessageCount: number
}

type AppContextValue = ClientAppState & {
  authChecked: boolean
  isLogin: boolean
  setConfig: React.Dispatch<React.SetStateAction<SiteConfig | null>>
  setCurrentUser: React.Dispatch<React.SetStateAction<UserSummary | null>>
}

const AppContext = React.createContext<AppContextValue | null>(null)

export function AppProvider({
  initialState,
  children,
}: {
  initialState: ClientAppState
  children: React.ReactNode
}) {
  const [config, setConfig] = React.useState(initialState.config)
  const [currentUser, setCurrentUserState] = React.useState(
    initialState.currentUser
  )
  const [authChecked, setAuthChecked] = React.useState(
    Boolean(initialState.currentUser)
  )
  const [unreadMessageCount, setUnreadMessageCount] = React.useState(
    initialState.unreadMessageCount
  )
  const userStateTouchedRef = React.useRef(false)

  const setCurrentUser = React.useCallback<
    React.Dispatch<React.SetStateAction<UserSummary | null>>
  >((nextUser) => {
    userStateTouchedRef.current = true
    setCurrentUserState(nextUser)
    setAuthChecked(true)
  }, [])

  React.useEffect(() => {
    let mounted = true

    async function hydrate() {
      const nextState = await loadClientAppState()
      if (!mounted) return

      if (nextState.config !== undefined) {
        setConfig(nextState.config)
      }
      if (!userStateTouchedRef.current) {
        setCurrentUserState(nextState.currentUser)
      }
      setUnreadMessageCount(nextState.unreadMessageCount)
      setAuthChecked(true)
    }

    void hydrate()

    return () => {
      mounted = false
    }
  }, [])

  const value = React.useMemo<AppContextValue>(
    () => ({
      config,
      currentUser,
      authChecked,
      locale: initialState.locale,
      unreadMessageCount,
      isLogin: Boolean(currentUser),
      setConfig,
      setCurrentUser,
    }),
    [authChecked, config, currentUser, initialState.locale, unreadMessageCount]
  )

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>
}

export function useAppState() {
  const value = React.useContext(AppContext)

  if (!value) {
    throw new Error("useAppState must be used within AppProvider")
  }

  return value
}

export function useAppConfig() {
  return useAppState().config
}

export function useCurrentUser() {
  return useAppState().currentUser
}

export function useUnreadMessageCount() {
  return useAppState().unreadMessageCount
}

export function useIsLogin() {
  return useAppState().isLogin
}

export function useAuthChecked() {
  return useAppState().authChecked
}

export function useAppLocale() {
  return useAppState().locale
}
