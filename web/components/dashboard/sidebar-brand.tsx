"use client"

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function SidebarBrand({
  name,
  logoSrc,
  description,
}: {
  name: string
  logoSrc: string
  description: string
}) {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size="lg"
          className="h-12 cursor-default rounded-lg group-data-[collapsible=icon]:gap-0 group-data-[collapsible=icon]:justify-center"
        >
          <div className="flex aspect-square size-9 shrink-0 items-center justify-center overflow-hidden rounded-md bg-white dark:bg-sidebar-accent">
            <img
              src={logoSrc}
              alt={name}
              className="size-full object-contain p-0.5"
            />
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
            <span className="truncate font-semibold">{name}</span>
            <span className="truncate text-xs text-sidebar-foreground/60">
              {description}
            </span>
          </div>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
