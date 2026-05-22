import { useSearchParams } from "react-router"

import { EmailVerifyResult } from "@/components/user/email-verify-result"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Email verification", "邮箱验证")
}

export default function EmailVerifyRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.email.verify.title"))
  const [searchParams] = useSearchParams()
  return (
    <section className="main">
      <div className="container">
        <div className="main-body no-bg">
          <EmailVerifyResult token={searchParams.get("token") || ""} />
        </div>
      </div>
    </section>
  )
}
