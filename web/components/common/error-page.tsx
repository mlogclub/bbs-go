"use client"

import * as React from "react"
import {
  ArrowLeft,
  ArrowUpRight,
  Home,
  RefreshCw,
  Search,
} from "lucide-react"
import { useNavigate } from "react-router"

import Link from "@/components/common/link"
import { WidgetCard } from "@/components/common/widget-card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useI18n } from "@/lib/i18n/provider"

export function ErrorPage({
  message,
  statusCode = 500,
}: {
  message?: string
  statusCode?: number
}) {
  const navigate = useNavigate()
  const { t } = useI18n()
  const [searchKeyword, setSearchKeyword] = React.useState("")
  const isNotFound = statusCode === 404
  const isForbidden = statusCode === 403

  const statusText =
    statusCode === 404
      ? t("pages.error.notFound")
      : statusCode === 403
        ? t("pages.error.forbidden")
      : t("pages.error.unknown")
  const titleText =
    statusCode === 404
      ? t("pages.error.title404")
      : statusCode === 403
        ? t("pages.error.title403")
      : t("pages.error.title500")
  const descriptionText =
    message ||
    (statusCode === 404
      ? t("pages.error.desc404")
      : statusCode === 403
        ? t("pages.error.desc403")
      : t("pages.error.desc500"))
  const quickLinks = [
    { to: "/", label: t("pages.error.quickHome") },
    { to: "/topics", label: t("pages.error.quickTopics") },
  ]

  function submitSearch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const keyword = searchKeyword.trim()
    navigate(keyword ? `/search?q=${encodeURIComponent(keyword)}` : "/search")
  }

  return (
    <section className="main">
      <div className="container mx-auto max-w-5xl">
        <WidgetCard bodyClassName="p-0">
          <div className="grid items-center gap-8 px-3 py-8 md:grid-cols-[minmax(0,0.92fr)_minmax(0,1.08fr)] md:px-6 md:py-12">
            <div className="order-2 md:order-1">
              <div className="mb-4 inline-flex items-center gap-2 rounded-full border border-border bg-muted/50 px-3 py-1 text-xs font-medium text-muted-foreground">
                {statusText}
              </div>

              <h1 className="text-2xl leading-tight font-semibold text-foreground sm:text-3xl">
                {titleText}
              </h1>
              <p className="mt-3 max-w-xl text-sm leading-7 text-muted-foreground">
                {descriptionText}
              </p>

              {isNotFound ? (
                <form
                  className="mt-6 flex flex-col gap-2 sm:flex-row"
                  onSubmit={submitSearch}
                >
                  <div className="flex h-10 min-w-0 flex-1 items-center rounded-md border border-input bg-background px-3 shadow-xs">
                    <Search className="mr-2 size-4 shrink-0 text-muted-foreground" />
                    <Input
                      value={searchKeyword}
                      type="text"
                      maxLength={30}
                      placeholder={t("pages.error.searchPlaceholder")}
                      className="h-auto border-none bg-transparent px-0 shadow-none focus-visible:ring-0"
                      autoComplete="off"
                      onChange={(event) => setSearchKeyword(event.currentTarget.value)}
                    />
                  </div>
                  <Button type="submit" className="h-10 sm:min-w-24">
                    {t("component.searchInput.searchBtn")}
                  </Button>
                </form>
              ) : null}

              <div className="mt-6 flex flex-wrap gap-2">
                <Button asChild>
                  <Link href="/">
                    <Home className="size-4" />
                    {t("pages.error.backHome")}
                  </Link>
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => window.history.back()}
                >
                  <ArrowLeft className="size-4" />
                  {t("pages.error.goBack")}
                </Button>
                {!isNotFound && !isForbidden ? (
                  <Button
                    type="button"
                    variant="secondary"
                    onClick={() => window.location.reload()}
                  >
                    <RefreshCw className="size-4" />
                    {t("pages.error.retry")}
                  </Button>
                ) : null}
              </div>
            </div>

            <div className="order-1 md:order-2">
              <div className="relative p-2 md:p-4">
                <div className="flex items-center justify-between border-b border-border pb-4">
                  <div className="flex items-center gap-2">
                    <div className="size-2.5 rounded-full bg-destructive/70" />
                    <div className="size-2.5 rounded-full bg-amber-500/70" />
                    <div className="size-2.5 rounded-full bg-emerald-500/70" />
                  </div>
                  <span className="text-xs font-medium text-muted-foreground">
                    {t("pages.error.header")}
                  </span>
                </div>

                <div className="py-8 text-center">
                  <p className="text-[72px] leading-none font-black text-muted-foreground/15 sm:text-[112px]">
                    {statusCode}
                  </p>
                </div>

                <div className="grid gap-2 border-t border-border pt-4 text-sm">
                  {quickLinks.map((link) => (
                    <Link
                      key={link.to}
                      href={link.to}
                      className="flex items-center justify-between rounded-md px-3 py-2 text-muted-foreground transition hover:bg-accent hover:text-accent-foreground"
                    >
                      <span>{link.label}</span>
                      <ArrowUpRight className="size-4" />
                    </Link>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </WidgetCard>
      </div>
    </section>
  )
}
