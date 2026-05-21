import type * as React from "react"

import { cn } from "@/lib/utils"

export function WidgetCard({
  title,
  actions,
  children,
  className,
  bodyClassName,
}: {
  title?: React.ReactNode
  actions?: React.ReactNode
  children: React.ReactNode
  className?: string
  bodyClassName?: string
}) {
  return (
    <section className={cn("rounded-lg bg-background px-3 py-2", className)}>
      {title || actions ? (
        <div className="flex min-h-8 items-center justify-between gap-3 border-b border-border pb-2 text-sm font-medium">
          <div className="min-w-0">{title}</div>
          {actions ? (
            <div className="shrink-0 text-xs text-primary">{actions}</div>
          ) : null}
        </div>
      ) : null}
      <div className={cn(title || actions ? "pt-2" : undefined, bodyClassName)}>
        {children}
      </div>
    </section>
  )
}
