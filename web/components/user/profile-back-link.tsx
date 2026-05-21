"use client"

import Link from "@/components/common/link"
import { ArrowLeft } from "lucide-react"

import { useRequiredUser } from "@/components/auth/require-user"
import { useI18n } from "@/lib/i18n/provider"

export function ProfileBackLink() {
  const user = useRequiredUser()
  const { t } = useI18n()

  return (
    <Link href={`/user/${user.id}`} className="space-x-2">
      <ArrowLeft className="inline-block h-4 w-4" />
      <span>{t("user.profile.backToProfile")}</span>
    </Link>
  )
}
