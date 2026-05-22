import Link from "@/components/common/link"
import { Settings, User } from "lucide-react"

import { cn } from "@/lib/utils"
import type { TFunction } from "@/lib/i18n"

export function ProfileShell({
  active,
  t,
  children,
}: {
  active: "profile" | "account"
  t: TFunction
  children: React.ReactNode
}) {
  const tabs = [
    {
      key: "profile" as const,
      href: "/user/profile",
      label: t("layout.profile.profile"),
      icon: User,
    },
    {
      key: "account" as const,
      href: "/user/profile/account",
      label: t("layout.profile.accountSettings"),
      icon: Settings,
    },
  ]

  const tabLink = (tab: (typeof tabs)[number], mobile = false) => {
    const Icon = tab.icon
    return (
      <Link
        key={tab.key}
        href={tab.href}
        className={cn(
          "flex items-center gap-1.5 rounded-md transition-all duration-300 ease-in-out",
          mobile
            ? "justify-center px-1 py-2 text-center text-xs font-medium"
            : "p-2.5",
          active === tab.key
            ? "bg-primary font-semibold text-primary-foreground shadow-md"
            : "text-foreground hover:bg-accent hover:text-primary"
        )}
      >
        <Icon className="h-4 w-4" />
        <span>{tab.label}</span>
      </Link>
    )
  }

  return (
    <section className="main">
      <div className="main-container right-main container">
        <div className="left-container">
          <div className="profile-edit-tabs-pc space-y-2 rounded-lg bg-background p-2.5">
            {tabs.map((tab) => tabLink(tab))}
          </div>
        </div>
        <div className="right-container">
          <div className="profile-edit-tabs-mobile mb-2.5 hidden rounded-lg bg-background p-1 shadow-sm">
            <ul className="flex gap-1">
              {tabs.map((tab) => (
                <li key={tab.key} className="flex-1">
                  {tabLink(tab, true)}
                </li>
              ))}
            </ul>
          </div>
          {children}
        </div>
      </div>
    </section>
  )
}
