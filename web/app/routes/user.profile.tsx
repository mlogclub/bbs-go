import { useCurrentUser } from "@/components/app/app-provider"
import { RequireUser } from "@/components/auth/require-user"
import { WidgetCard } from "@/components/common/widget-card"
import { ProfileBackLink } from "@/components/user/profile-back-link"
import { ProfileForm } from "@/components/user/profile-form"
import { ProfileShell } from "@/components/user/profile-shell"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { requireUser, requireUserClient } from "../route-helpers/auth"

export async function loader(args: { request: Request }) {
  await requireUser(args)
  return null
}

export async function clientLoader(args: { request: Request }) {
  await requireUserClient(args)
  return null
}

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Profile", "个人资料")
}

export default function ProfileRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.profile.title"))
  const user = useCurrentUser()
  return (
    <RequireUser initialUser={user} redirectPath="/user/profile">
      <ProfileShell active="profile" t={t}>
        <WidgetCard title={t("user.profile.title")} actions={<ProfileBackLink />}>
          <ProfileForm user={user || undefined} />
        </WidgetCard>
      </ProfileShell>
    </RequireUser>
  )
}
