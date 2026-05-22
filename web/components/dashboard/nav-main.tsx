"use client"

import * as React from "react"
import { Link, useLocation } from "react-router-dom"
import type { LucideIcon } from "lucide-react"
import { ChevronRightIcon } from "lucide-react"

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"
import { useI18n } from "@/lib/i18n/provider"

const DASHBOARD_NAV_OPEN_STORAGE_KEY = "bbsgo-dashboard-nav-open"

type NavMainItem = {
  title: string
  url: string
  icon?: LucideIcon
  isActive?: boolean
  items?: {
    title: string
    url: string
  }[]
}

function readStoredOpenState() {
  if (typeof window === "undefined") return {}
  try {
    const value = window.localStorage.getItem(DASHBOARD_NAV_OPEN_STORAGE_KEY)
    return value ? (JSON.parse(value) as Record<string, boolean>) : {}
  } catch {
    return {}
  }
}

function writeStoredOpenState(value: Record<string, boolean>) {
  if (typeof window === "undefined") return
  window.localStorage.setItem(
    DASHBOARD_NAV_OPEN_STORAGE_KEY,
    JSON.stringify(value)
  )
}

export function NavMain({ items }: { items: NavMainItem[] }) {
  const { t } = useI18n()
  const location = useLocation()
  const pathname = location.pathname
  const [openState, setOpenState] = React.useState<Record<string, boolean>>({})
  const [loadedOpenState, setLoadedOpenState] = React.useState(false)

  const isActiveUrl = (url: string) =>
    pathname === url || (url !== "/dashboard" && pathname.startsWith(`${url}/`))

  React.useEffect(() => {
    setOpenState(readStoredOpenState())
    setLoadedOpenState(true)
  }, [])

  function updateOpenState(key: string, open: boolean) {
    setOpenState((current) => {
      const next = { ...current, [key]: open }
      writeStoredOpenState(next)
      return next
    })
  }

  return (
    <SidebarGroup>
      <SidebarGroupLabel>{t("dashboard.nav.platform")}</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => {
          const Icon = item.icon
          const hasChildren = Boolean(item.items?.length)
          const active =
            isActiveUrl(item.url) ||
            Boolean(item.items?.some((child) => isActiveUrl(child.url)))

          if (!hasChildren) {
            return (
              <SidebarMenuItem key={item.url}>
                <SidebarMenuButton
                  asChild
                  tooltip={item.title}
                  isActive={active}
                  className="relative rounded-md data-[active=true]:bg-[var(--dashboard-accent-soft)] data-[active=true]:text-foreground data-[active=true]:shadow-xs data-[active=true]:before:absolute data-[active=true]:before:top-1.5 data-[active=true]:before:left-0 data-[active=true]:before:h-5 data-[active=true]:before:w-0.5 data-[active=true]:before:rounded-full data-[active=true]:before:bg-[var(--dashboard-accent)]"
                >
                  <Link to={item.url}>
                    {Icon ? <Icon /> : null}
                    <span>{item.title}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          }

          return (
            <Collapsible
              key={item.url}
              asChild
              open={loadedOpenState ? Boolean(openState[item.url]) : false}
              onOpenChange={(open) => updateOpenState(item.url, open)}
              className="group/collapsible"
            >
              <SidebarMenuItem>
                <CollapsibleTrigger asChild>
                  <SidebarMenuButton
                    tooltip={item.title}
                    isActive={active}
                    className="relative rounded-md data-[active=true]:bg-[var(--dashboard-accent-soft)] data-[active=true]:text-foreground data-[active=true]:shadow-xs data-[active=true]:before:absolute data-[active=true]:before:top-1.5 data-[active=true]:before:left-0 data-[active=true]:before:h-5 data-[active=true]:before:w-0.5 data-[active=true]:before:rounded-full data-[active=true]:before:bg-[var(--dashboard-accent)]"
                  >
                    {Icon ? <Icon /> : null}
                    <span>{item.title}</span>
                    <ChevronRightIcon className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                  </SidebarMenuButton>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    {item.items?.map((subItem) => (
                      <SidebarMenuSubItem key={subItem.title}>
                        <SidebarMenuSubButton
                          asChild
                          isActive={isActiveUrl(subItem.url)}
                          className="rounded-md data-[active=true]:bg-[var(--dashboard-accent-soft)] data-[active=true]:text-foreground"
                        >
                          <Link to={subItem.url}>
                            <span>{subItem.title}</span>
                          </Link>
                        </SidebarMenuSubButton>
                      </SidebarMenuSubItem>
                    ))}
                  </SidebarMenuSub>
                </CollapsibleContent>
              </SidebarMenuItem>
            </Collapsible>
          )
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}
