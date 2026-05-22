import * as React from "react"

import { useI18n } from "@/lib/i18n/provider"

export function useLoadMoreLabels() {
  const { t } = useI18n()
  return {
    loadMore: t("common.loadMore.loadMore"),
    noMore: t("common.loadMore.noMore"),
  }
}

export function useClientData<T>(key: string, load: () => Promise<T>) {
  const [data, setData] = React.useState<T | null>(null)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    let mounted = true
    setLoading(true)
    setError(null)

    void load()
      .then((value) => {
        if (mounted) setData(value)
      })
      .catch((err) => {
        if (mounted) {
          setError(err instanceof Error ? err.message : String(err))
        }
      })
      .finally(() => {
        if (mounted) setLoading(false)
      })

    return () => {
      mounted = false
    }
  }, [key])

  return { data, loading, error }
}

export function useMediaQuery(query: string, defaultValue = false) {
  const subscribe = React.useCallback(
    (onStoreChange: () => void) => {
      if (typeof window === "undefined") {
        return () => {}
      }
      const mediaQueryList = window.matchMedia(query)
      mediaQueryList.addEventListener("change", onStoreChange)
      return () => mediaQueryList.removeEventListener("change", onStoreChange)
    },
    [query]
  )
  const getSnapshot = React.useCallback(() => {
    if (typeof window === "undefined") {
      return defaultValue
    }
    return window.matchMedia(query).matches
  }, [defaultValue, query])

  return React.useSyncExternalStore(subscribe, getSnapshot, () => defaultValue)
}
