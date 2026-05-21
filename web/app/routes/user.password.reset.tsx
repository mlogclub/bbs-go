import { useSearchParams } from "react-router"

import { ResetPasswordForm } from "@/components/auth/password-reset-forms"
import { WidgetCard } from "@/components/common/widget-card"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Reset password", "重置密码")
}

export default function ResetRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.passwordReset.reset.title"))
  const [searchParams] = useSearchParams()
  return (
    <section className="main">
      <div className="container">
        <div className="main-body no-bg">
          <WidgetCard className="mx-auto max-w-160">
            <ResetPasswordForm token={(searchParams.get("token") || "").trim()} />
          </WidgetCard>
        </div>
      </div>
    </section>
  )
}
