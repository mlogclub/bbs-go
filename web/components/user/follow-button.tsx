"use client"

import * as React from "react"
import { UserCheck, UserPlus } from "lucide-react"

import { followAction } from "@/lib/actions/user"
import { Button } from "@/components/ui/button"
import { toast } from "@/lib/toast"
import { useI18n } from "@/lib/i18n/provider"

export function FollowButton({
  userId,
  initialFollowed,
  onChanged,
}: {
  userId: string
  initialFollowed?: boolean
  onChanged?: (followed: boolean) => void
}) {
  const { t } = useI18n()
  const [followed, setFollowed] = React.useState(Boolean(initialFollowed))
  const [pending, startTransition] = React.useTransition()

  function submit() {
    startTransition(async () => {
      const result = await followAction(userId, followed)
      if (!result.ok) {
        toast.error(result.message || t("composables.unknownError"))
        return
      }
      const next = Boolean(result.followed)
      setFollowed(next)
      onChanged?.(next)
    })
  }

  return (
    <Button
      type="button"
      variant={followed ? "outline" : "default"}
      size="sm"
      className="h-7 text-xs"
      disabled={pending}
      onClick={submit}
    >
      {followed ? (
        <UserCheck className="size-3" aria-hidden="true" />
      ) : (
        <UserPlus className="size-3" aria-hidden="true" />
      )}
      <span className="ml-1">
        {followed
          ? t("component.followBtn.followed")
          : t("component.followBtn.follow")}
      </span>
    </Button>
  )
}
