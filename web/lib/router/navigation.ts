import {
  useLocation,
  useNavigate,
  useParams,
  useSearchParams as useReactRouterSearchParams,
} from "react-router-dom"
import * as React from "react"

export function usePathname() {
  return useLocation().pathname
}

export function useSearchParams() {
  return useReactRouterSearchParams()[0]
}

export function useRouter() {
  const navigate = useNavigate()

  return React.useMemo(
    () => ({
      push: (href: string) => void navigate(href),
      replace: (href: string) => void navigate(href, { replace: true }),
      back: () => window.history.back(),
      forward: () => window.history.forward(),
      refresh: () => window.location.reload(),
      prefetch: () => undefined,
    }),
    [navigate]
  )
}

export function redirect(href: string): never {
  if (typeof window !== "undefined") {
    window.location.replace(href)
  }
  throw new Error(`Redirect: ${href}`)
}

export function notFound(): never {
  throw new Error("Not found")
}

export function useSelectedLayoutSegment() {
  const params = useParams()
  return Object.values(params)[0] ?? null
}
