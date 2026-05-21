import { ForgotPasswordForm } from "@/components/auth/password-reset-forms"
import { WidgetCard } from "@/components/common/widget-card"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Forgot password", "找回密码")
}

export default function ForgotRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.passwordReset.forgot.title"))
  return (
    <section className="main">
      <div className="container">
        <div className="main-body no-bg">
          <WidgetCard className="mx-auto max-w-[520px] p-6">
            <ForgotPasswordForm />
          </WidgetCard>
        </div>
      </div>
    </section>
  )
}
