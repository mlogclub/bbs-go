import * as React from "react"

export default function dynamic(
  loader: () => Promise<unknown>,
  _options?: { ssr?: boolean; loading?: React.ComponentType }
): React.ComponentType<Record<string, unknown>> {
  return React.lazy(async () => {
    const mod = await loader()
    if (
      mod &&
      typeof mod === "object" &&
      "default" in mod &&
      typeof (mod as { default: unknown }).default === "function"
    ) {
      return mod as { default: React.ComponentType<unknown> }
    }
    return { default: mod as React.ComponentType<unknown> }
  })
}
