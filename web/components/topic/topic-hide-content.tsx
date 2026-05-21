import { Lock } from "lucide-react"

import { HtmlImagePreview } from "@/components/common/image-preview"
import type { TopicHideContent } from "@/lib/api/types"
import type { TFunction } from "@/lib/i18n"

export function TopicHideContent({ hideContent, t }: { hideContent?: TopicHideContent | null; t: TFunction }) {
  if (!hideContent?.exists) {
    return null
  }

  return (
    <div className="bbs-content line-numbers my-5 break-words text-base leading-6 antialiased">
      {hideContent.show ? (
        <section className="rounded-md border border-border bg-background">
          <div className="flex items-center gap-2 border-b px-3 py-2 text-sm font-medium">
            <Lock size={16} />
            <span>&nbsp;{t("pages.topic.detail.hideContent")}</span>
          </div>
          <HtmlImagePreview
            html={hideContent.content || ""}
            className="p-3 [&_img]:cursor-zoom-in"
          />
        </section>
      ) : (
        <div className="flex items-center gap-2 rounded-md border border-gray-200 p-2 text-gray-600">
          <Lock size={16} />
          <span>{t("pages.topic.detail.hideContentTip")}</span>
        </div>
      )}
    </div>
  )
}
