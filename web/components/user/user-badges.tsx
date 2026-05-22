import { Award } from "lucide-react"

import type { TFunction } from "@/lib/i18n"

export function UserBadges({ t }: { t: TFunction }) {
  return (
    <section className="rounded-lg border bg-background p-4">
      <div className="flex items-center gap-2 text-sm font-medium">
        <Award className="h-4 w-4" aria-hidden="true" />
        <span>{t("user.profile.badges")}</span>
      </div>
      <p className="mt-3 text-sm text-muted-foreground">
        {t("user.profile.noBadges")}
      </p>
    </section>
  )
}
