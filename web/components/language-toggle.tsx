"use client"

import { useSyncExternalStore } from "react"
import { LanguagesIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { Locale } from "@/lib/i18n"
import { useI18n } from "@/lib/i18n/provider"

const languageOptions: Array<{
  value: Locale
  label: string
}> = [
  { value: "en-US", label: "English" },
  { value: "zh-CN", label: "中文" },
]

export function LanguageToggle() {
  const { locale, setLocale, t } = useI18n()
  const mounted = useSyncExternalStore(
    () => () => {},
    () => true,
    () => false
  )

  const activeLocale = mounted ? locale : "en-US"

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild aria-label={t("dashboard.header.language")}>
        <Button variant="outline" size="sm">
          <LanguagesIcon />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40 min-w-40">
        <DropdownMenuRadioGroup
          value={activeLocale}
          onValueChange={(value) => setLocale?.(value as Locale)}
        >
          {languageOptions.map((option) => (
            <DropdownMenuRadioItem key={option.value} value={option.value}>
              <LanguagesIcon />
              {option.label}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
