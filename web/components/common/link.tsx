import RouterLink from "@/lib/router/link"
import type * as React from "react"

type AppLinkProps = React.ComponentProps<typeof RouterLink>

export default function Link({ prefetch = false, ...props }: AppLinkProps) {
  return <RouterLink prefetch={prefetch} {...props} />
}
