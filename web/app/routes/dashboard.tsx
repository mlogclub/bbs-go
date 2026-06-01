import * as React from "react"
import {
  isRouteErrorResponse,
  Outlet,
  useLocation,
  useRouteError,
} from "react-router"
import { ExternalLinkIcon, HouseIcon } from "lucide-react"

import { RequireDashboardAdmin } from "@/components/auth/require-dashboard-admin"
import Link from "@/components/common/link"
import { ErrorPage } from "@/components/common/error-page"
import { AppSidebar } from "@/components/dashboard/app-sidebar"
import {
  InstallRequiredFallback,
  isInstallRequiredRouteError,
} from "@/components/install/install-required-fallback"
import { LanguageToggle } from "@/components/language-toggle"
import { ThemeToggle } from "@/components/theme-toggle"
import { Button } from "@/components/ui/button"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { useI18n } from "@/lib/i18n/provider"
import { useDocumentTitle } from "@/lib/use-document-title"
import { resolveBreadcrumb } from "@/lib/dashboard/navigation"

import {
  requireDashboardAdmin,
  requireDashboardAdminClient,
} from "../route-helpers/auth"

const DASHBOARD_SIDEBAR_OPEN_STORAGE_KEY = "bbsgo-dashboard-sidebar-open"

type DashboardBreadcrumbItem = {
  title: string
  url?: string
}

function dashboardBreadcrumbs(
  pathname: string,
  t: ReturnType<typeof useI18n>["t"]
): DashboardBreadcrumbItem[] {
  return resolveBreadcrumb(pathname, (key) => t(key))
}

async function _loader(args: { request: Request }) {
  await requireDashboardAdmin(args)
  return null
}

export async function clientLoader(args: { request: Request }) {
  await requireDashboardAdminClient(args)
  return null
}

export function ErrorBoundary() {
  const error = useRouteError()
  if (isInstallRequiredRouteError(error)) {
    return <InstallRequiredFallback />
  }

  const statusCode = isRouteErrorResponse(error) ? error.status : 500
  const message =
    error instanceof Error
      ? error.message
      : isRouteErrorResponse(error) && statusCode !== 403
        ? error.statusText
        : undefined

  return <ErrorPage statusCode={statusCode} message={message} />
}

export default function DashboardLayout() {
  const { t } = useI18n()
  const location = useLocation()
  const breadcrumbs = dashboardBreadcrumbs(location.pathname, t)
  useDocumentTitle(t("dashboard.brand.name"), t("dashboard.brand.plan"), {
    appendSiteTitle: false,
  })
  const [sidebarOpen, setSidebarOpen] = React.useState(true)

  React.useEffect(() => {
    const value = window.localStorage.getItem(
      DASHBOARD_SIDEBAR_OPEN_STORAGE_KEY
    )
    if (value !== null) {
      setSidebarOpen(value === "true")
    }
  }, [])

  function updateSidebarOpen(open: boolean) {
    setSidebarOpen(open)
    window.localStorage.setItem(
      DASHBOARD_SIDEBAR_OPEN_STORAGE_KEY,
      String(open)
    )
  }

  return (
    <RequireDashboardAdmin>
      <SidebarProvider
        open={sidebarOpen}
        onOpenChange={updateSidebarOpen}
        data-dashboard-layout
        className="h-svh overflow-hidden bg-[var(--dashboard-surface)]"
      >
        <AppSidebar />
        <SidebarInset className="min-h-0 overflow-hidden bg-[var(--dashboard-surface)]">
          <header className="flex h-16 shrink-0 items-center gap-2 border-b bg-[var(--dashboard-panel)]/95 shadow-[0_1px_0_rgba(15,23,42,0.03)] backdrop-blur transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
            <div className="flex min-w-0 items-center gap-2 px-4">
              <SidebarTrigger className="-ml-1" />
              <Separator
                orientation="vertical"
                className="mr-2 data-vertical:h-4 data-vertical:self-auto"
              />
              <Breadcrumb>
                <BreadcrumbList>
                  {breadcrumbs.map((item, index) => {
                    const isLast = index === breadcrumbs.length - 1

                    return (
                      <React.Fragment key={`${item.url}-${item.title}`}>
                        <BreadcrumbItem
                          className={
                            index === 0 && !isLast ? "hidden md:block" : undefined
                          }
                        >
                          {isLast ? (
                            <BreadcrumbPage>{item.title}</BreadcrumbPage>
                          ) : (
                            <BreadcrumbLink asChild>
                              <Link href={item.url || "/dashboard"}>
                                {item.title}
                              </Link>
                            </BreadcrumbLink>
                          )}
                        </BreadcrumbItem>
                        {!isLast ? (
                          <BreadcrumbSeparator className="hidden md:block" />
                        ) : null}
                      </React.Fragment>
                    )
                  })}
                </BreadcrumbList>
              </Breadcrumb>
            </div>
            <div className="ml-auto flex shrink-0 items-center gap-2 px-4">
              <Button variant="outline" size="icon-sm" asChild>
                <a
                  href="/"
                  target="_blank"
                  rel="noreferrer"
                  aria-label={t("dashboard.header.siteHome")}
                  title={t("dashboard.header.siteHome")}
                >
                  <HouseIcon />
                </a>
              </Button>
              <ThemeToggle />
              <LanguageToggle />
            </div>
          </header>
          <div className="min-h-0 flex-1 overflow-y-auto">
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </RequireDashboardAdmin>
  )
}
