"use client"

import { useSyncExternalStore } from "react"
import { LaptopIcon, MoonIcon, SunIcon } from "lucide-react"

import { useTheme } from "@/components/theme-provider"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useI18n } from "@/lib/i18n/provider"
import { cn } from "@/lib/utils"

type ThemeMode = "light" | "dark" | "system"
type ThemeToggleButtonVariant = "outline" | "ghost"
type ThemeToggleButtonSize = "sm" | "icon-sm"

const themeOptions: Array<{
  value: ThemeMode
  icon: typeof SunIcon
}> = [
  { value: "system", icon: LaptopIcon },
  { value: "light", icon: SunIcon },
  { value: "dark", icon: MoonIcon },
]

export function ThemeToggle({
  variant = "outline",
  size = "sm",
  className,
}: {
  variant?: ThemeToggleButtonVariant
  size?: ThemeToggleButtonSize
  className?: string
}) {
  const { theme, setTheme } = useTheme()
  const { t } = useI18n()
  const mounted = useSyncExternalStore(
    () => () => {},
    () => true,
    () => false
  )

  const activeTheme = mounted
    ? ((theme as ThemeMode | undefined) ?? "system")
    : "system"
  const ActiveIcon =
    themeOptions.find((option) => option.value === activeTheme)?.icon ??
    LaptopIcon

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild aria-label={t("common.theme.toggle")}>
        <Button
          variant={variant}
          size={size}
          className={cn(
            "text-muted-foreground hover:text-foreground",
            className
          )}
        >
          <ActiveIcon />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40 min-w-40">
        <DropdownMenuRadioGroup
          value={activeTheme}
          onValueChange={(value) => setTheme(value as ThemeMode)}
        >
          {themeOptions.map((option) => {
            const Icon = option.icon
            return (
              <DropdownMenuRadioItem key={option.value} value={option.value}>
                <Icon />
                {t(`common.theme.${option.value}`)}
              </DropdownMenuRadioItem>
            )
          })}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
