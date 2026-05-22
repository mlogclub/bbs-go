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
          className="h-12 cursor-default rounded-lg border border-sidebar-border/70 bg-background/60 shadow-xs hover:bg-background/60 data-[state=open]:bg-background/60 group-data-[collapsible=icon]:gap-0 group-data-[collapsible=icon]:justify-center"
        >
          <div className="flex aspect-square size-8 shrink-0 items-center justify-center overflow-hidden rounded-md bg-white ring-1 ring-sidebar-border dark:bg-sidebar-accent">
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
