"use client"

import * as React from "react"
import { usePathname } from "@/lib/router/navigation"

export function useRouteSegment(index: number) {
  const pathname = usePathname()
  return React.useMemo(() => {
    const segments = pathname.split("/").filter(Boolean)
    const value = segments[index]
    return value ? decodeURIComponent(value) : ""
  }, [index, pathname])
}

export function useRouteData<T>(
  key: string,
  load: () => Promise<T>,
  initialData: T | null = null
) {
  const [data, setData] = React.useState<T | null>(initialData)
  const [loading, setLoading] = React.useState(!initialData)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    if (!key) return

    let active = true
    const timer = window.setTimeout(() => {
      if (!active) return
      setLoading(true)
      setError(null)
    }, 0)

    void load()
      .then((nextData) => {
        if (active) setData(nextData)
      })
      .catch((err) => {
        if (active) setError(err instanceof Error ? err.message : String(err))
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
      window.clearTimeout(timer)
    }
  }, [key, load])

  return { data, loading, error }
}
