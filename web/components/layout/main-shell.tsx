import { cn } from "@/lib/utils"

export function MainShell({
  children,
  aside,
  sideSize = "320",
  className,
  containerClassName,
  asideClassName,
}: {
  children: React.ReactNode
  aside?: React.ReactNode
  sideSize?: "260" | "320" | "360"
  className?: string
  containerClassName?: string
  asideClassName?: string
}) {
  return (
    <main className={cn("main", className)}>
      <div
        className={cn(
          "container main-container",
          aside && "left-main",
          aside && `side-size-${sideSize}`,
          containerClassName,
        )}
      >
        <div className="left-container">{children}</div>
        {aside ? (
          <aside className={cn("right-container space-y-4", asideClassName)}>
            {aside}
          </aside>
        ) : null}
      </div>
    </main>
  )
}
