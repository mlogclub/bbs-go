import { Navigate } from "react-router"

import { InstallWizard, type DbType } from "@/components/install/install-wizard"
import { apiFetch } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"
import { noindexRouteMeta } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { useClientData } from "../route-helpers/client-hooks"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "Install BBS-GO", "安装 BBS-GO")
}

export default function InstallRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("pages.install.title"))
  const { data, loading } = useClientData<{
    installed?: boolean
    dockerBuiltinDbType?: DbType
    dockerBuiltinMysql?: boolean
    dbType?: string
  }>("install-status", () =>
    apiFetch<{ installed?: boolean }>("/api/install/status").catch(() => ({
      installed: true,
    }))
  )

  if (loading) {
    return (
      <main className="main">
        <div className="container" />
      </main>
    )
  }
  if (data?.installed) return <Navigate to="/" replace />

  return (
    <InstallWizard
      dockerBuiltinDbType={
        data?.dockerBuiltinDbType ||
        (data?.dockerBuiltinMysql ? "mysql" : undefined)
      }
    />
  )
}
