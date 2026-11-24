import { Navigate } from "react-router"

import { InstallWizard, type DbType } from "@/components/install/install-wizard"
import { apiFetch } from "@/lib/api/client"
import { noindexRouteMeta } from "@/lib/seo"

import { useClientData } from "../route-helpers/client-hooks"

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  return noindexRouteMeta(matches, "install_bbs_go")
}

export default function InstallRoute() {
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
