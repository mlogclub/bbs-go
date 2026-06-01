import { useAppState } from "@/components/app/app-provider"
import { RequireUser } from "@/components/auth/require-user"
import { WidgetCard } from "@/components/common/widget-card"
import { AccountSettings } from "@/components/user/account-settings"
import { ProfileBackLink } from "@/components/user/profile-back-link"
import { ProfileShell } from "@/components/user/profile-shell"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { requireUser, requireUserClient } from "../route-helpers/auth"

async function _loader(args: { request: Request }) {
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
  return noindexRouteMeta(matches, "Account settings", "账号设置")
}

export default function AccountRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.profile.account.title"))
  const { config, currentUser } = useAppState()
  return (
    <RequireUser initialUser={currentUser} redirectPath="/user/profile/account">
      <ProfileShell active="account" t={t}>
        <WidgetCard
          title={t("user.profile.account.title")}
          actions={<ProfileBackLink />}
        >
          <AccountSettings
            user={currentUser || undefined}
            config={config}
            bindInfo={{}}
          />
        </WidgetCard>
      </ProfileShell>
    </RequireUser>
  )
}
