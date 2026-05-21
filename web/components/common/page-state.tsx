import { EmptyState } from "@/components/common/empty-state"

export function PageLoading() {
  return (
    <div className="rounded-lg bg-background p-5 text-sm text-muted-foreground" />
  )
}

export function PageError({ message }: { message?: string | null }) {
  return <EmptyState title={message || "No data"} />
}
