import * as React from "react"

import { useAppConfig } from "@/components/app/app-provider"

type DocumentTitleOptions = {
  appendSiteTitle?: boolean
}

type DocumentTitleArg = string | null | undefined | DocumentTitleOptions

function isDocumentTitleOptions(value: DocumentTitleArg): value is DocumentTitleOptions {
  return typeof value === "object" && value !== null && "appendSiteTitle" in value
}

export function useDocumentTitle(...args: DocumentTitleArg[]) {
  const config = useAppConfig()
  const lastArg = args[args.length - 1]
  const options = isDocumentTitleOptions(lastArg) ? lastArg : undefined
  const appendSiteTitle = options?.appendSiteTitle ?? true
  const parts = (options ? args.slice(0, -1) : args).filter(
    (part): part is string => typeof part === "string" && part.length > 0
  )
  const title = parts.join(" - ")

  React.useEffect(() => {
    if (typeof document === "undefined") {
      return
    }

    const siteTitle = config?.siteTitle || "BBS-GO"
    document.title = appendSiteTitle
      ? title
        ? `${title} - ${siteTitle}`
        : siteTitle
      : title || siteTitle
  }, [appendSiteTitle, config?.siteTitle, title])
}
