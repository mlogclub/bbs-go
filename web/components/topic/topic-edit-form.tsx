"use client"

import * as React from "react"
import { useRouter } from "@/lib/router/navigation"
import { Trash2 } from "lucide-react"

import { TagInput } from "@/components/common/tag-input"
import { ContentEditor } from "@/components/editor/content-editor"
import { CategoryQuickSelector } from "@/components/topic/category-selector"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { apiFetch } from "@/lib/api/client"
import type {
  SiteConfig,
  TopicAttachment,
  Category,
} from "@/lib/api/types"
import type { TopicEditData } from "@/lib/api/topics"
import { useI18n } from "@/lib/i18n/provider"
import {
  filterCategoryTree,
  getFirstCategoryId,
  hasCategory,
} from "@/lib/categories"
import { msg, useToastActions } from "@/lib/toast"

type TopicEditFormState = {
  id: string
  type: number
  categoryId: number
  title: string
  content: string
  contentType: "html" | "markdown"
  hideContent: string
  tags: string[]
}

function categoryTypeMatches(topicType: number) {
  return (node: Category) =>
    topicType === 2 ? node.type === "qa" : node.type !== "qa"
}

function normalizeEditData(topic: TopicEditData): TopicEditFormState {
  return {
    id: topic.id,
    type: Number(topic.type) || 0,
    categoryId: Number(topic.categoryId) || 0,
    title: topic.title || "",
    content: topic.content || "",
    contentType: topic.contentType === "markdown" ? "markdown" : "html",
    hideContent: topic.hideContent || "",
    tags: Array.isArray(topic.tags) ? topic.tags : [],
  }
}

function TopicAttachmentField({
  value,
  uploading,
  config,
  onUploadingChange,
  onChange,
}: {
  value: TopicAttachment[]
  uploading: boolean
  config?: SiteConfig["attachmentConfig"]
  onUploadingChange: (value: boolean) => void
  onChange: (value: TopicAttachment[]) => void
}) {
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const inputRef = React.useRef<HTMLInputElement>(null)
  const maxCount = config?.maxCount ?? 5
  const maxSizeMB = config?.maxSizeMB ?? 10
  const accept = Array.isArray(config?.allowedTypes)
    ? config.allowedTypes.join(",")
    : ".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.md,.csv,.zip,.rar,.7z,.tar,.gz"

  async function upload(file: File) {
    onUploadingChange(true)
    try {
      const body = new FormData()
      body.append("file", file, file.name)
      body.append("downloadScore", "0")
      const attachment = await apiFetch<TopicAttachment>(
        "/api/attachment/upload",
        {
          method: "POST",
          body,
        }
      )
      onChange([...value, attachment])
    } catch (error) {
      catchError(error)
    } finally {
      onUploadingChange(false)
    }
  }

  async function updateScore(
    attachment: TopicAttachment,
    downloadScore: number
  ) {
    onChange(
      value.map((item) =>
        item.id === attachment.id ? { ...item, downloadScore } : item
      )
    )
    try {
      await apiFetch<null>("/api/attachment/update_download_score", {
        method: "POST",
        body: { id: attachment.id, downloadScore },
      })
    } catch (error) {
      catchError(error)
    }
  }

  return (
    <div className="rounded-md border border-dashed bg-muted/20 p-3">
      <div className="flex flex-wrap items-center gap-2">
        <span className="text-sm text-muted-foreground">
          {t("pages.topic.create.attachment.label")}
        </span>
        <span className="text-xs text-muted-foreground">
          {t("pages.topic.create.attachment.limitHint", {
            maxCount,
            maxSizeMB,
          })}
        </span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={uploading || value.length >= maxCount}
          onClick={() => inputRef.current?.click()}
        >
          {t("pages.topic.create.attachment.add")}
        </Button>
        <input
          ref={inputRef}
          type="file"
          className="hidden"
          accept={accept}
          onChange={(event) => {
            const file = event.currentTarget.files?.[0]
            if (file) void upload(file)
            event.currentTarget.value = ""
          }}
        />
      </div>
      {uploading ? (
        <div className="mt-2 text-xs text-muted-foreground">
          {t("pages.topic.create.attachmentUploading")}
        </div>
      ) : null}
      {value.length ? (
        <ul className="mt-2 space-y-2 text-sm">
          {value.map((attachment, index) => (
            <li
              key={attachment.id || index}
              className="flex flex-col gap-2 rounded border bg-background p-2 sm:flex-row sm:items-center"
            >
              <div className="min-w-0 flex-1">
                <span className="block truncate font-medium">
                  {attachment.fileName}
                </span>
                <span className="text-xs text-muted-foreground">
                  {attachment.fileSize || 0} B
                </span>
              </div>
              <label className="flex shrink-0 items-center gap-2 text-xs text-muted-foreground">
                {t("pages.topic.create.attachment.scorePlaceholder")}
                <Input
                  type="number"
                  min="0"
                  step="1"
                  className="h-8 w-20"
                  value={attachment.downloadScore ?? 0}
                  onChange={(event) =>
                    void updateScore(
                      attachment,
                      Math.max(0, Number(event.currentTarget.value) || 0)
                    )
                  }
                />
              </label>
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                className="text-muted-foreground hover:text-destructive"
                aria-label={t("pages.topic.create.attachment.remove")}
                onClick={() =>
                  onChange(value.filter((_, itemIndex) => itemIndex !== index))
                }
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  )
}

export function TopicEditForm({
  topic,
  config,
  categories,
}: {
  topic: TopicEditData
  config: SiteConfig | null
  categories: Category[]
}) {
  const router = useRouter()
  const { t } = useI18n()
  const { catchError, msgWarning } = useToastActions()
  const lastSubmitAtRef = React.useRef(0)
  const [publishing, setPublishing] = React.useState(false)
  const [attachmentUploading, setAttachmentUploading] = React.useState(false)
  const [attachmentList, setAttachmentList] = React.useState<TopicAttachment[]>(
    () => (Array.isArray(topic.attachments) ? topic.attachments : [])
  )
  const [form, setForm] = React.useState<TopicEditFormState>(() =>
    normalizeEditData(topic)
  )
  const availableNodes = React.useMemo(
    () => filterCategoryTree(categories, categoryTypeMatches(form.type)),
    [form.type, categories]
  )
  const effectiveCategoryId = hasCategory(availableNodes, form.categoryId)
    ? form.categoryId
    : getFirstCategoryId(availableNodes)

  function updateForm(next: Partial<TopicEditFormState>) {
    setForm((current) => ({ ...current, ...next }))
  }

  async function submit() {
    const now = Date.now()
    if (now - lastSubmitAtRef.current < 500 || publishing) {
      return
    }
    lastSubmitAtRef.current = now

    if (attachmentUploading) {
      msgWarning(t("pages.topic.create.attachmentUploading"))
      return
    }

    setPublishing(true)
    try {
      await apiFetch<null>(`/api/topic/edit/${form.id}`, {
        method: "POST",
        body: {
          categoryId: effectiveCategoryId,
          title: form.title,
          content: form.content,
          hideContent: form.hideContent,
          tags: form.tags,
          attachmentIds:
            form.type === 0 ? attachmentList.map((item) => item.id) : [],
        },
      })
      msg({
        message: t("pages.topic.edit.success"),
        onClose() {
          router.push(`/topic/${form.id}`)
        },
      })
    } catch (error) {
      catchError(error)
      setPublishing(false)
    }
  }

  return (
    <div className="publish-form">
      <div className="form-title">
        <div className="form-title-name">{t("pages.topic.edit.title")}</div>
      </div>

      <div className="field">
        <CategoryQuickSelector
          value={effectiveCategoryId}
          categories={availableNodes}
          onChange={(categoryId) => updateForm({ categoryId })}
        />
      </div>

      <div className="field">
        <Input
          value={form.title}
          placeholder={t("pages.topic.edit.titlePlaceholder")}
          onChange={(event) => updateForm({ title: event.currentTarget.value })}
        />
      </div>

      <div className="field">
        <ContentEditor
          contentType={form.contentType}
          value={form.content}
          placeholder={t("pages.topic.edit.contentPlaceholder")}
          height="400px"
          onChange={(content) => updateForm({ content })}
        />
      </div>

      {form.type !== 2 && (config?.enableHideContent || form.hideContent) ? (
        <div className="field">
          <ContentEditor
            contentType="html"
            value={form.hideContent}
            height="200px"
            onChange={(hideContent) => updateForm({ hideContent })}
          />
        </div>
      ) : null}

      <div className="field">
        <TagInput
          value={form.tags}
          recommendTags={config?.recommendTags}
          placeholder={t("component.tagInput.placeholder")}
          onChange={(tags) => updateForm({ tags })}
        />
      </div>

      {form.type === 0 && config?.attachmentConfig?.enabled ? (
        <div className="field">
          <TopicAttachmentField
            value={attachmentList}
            config={config.attachmentConfig}
            uploading={attachmentUploading}
            onUploadingChange={setAttachmentUploading}
            onChange={setAttachmentList}
          />
        </div>
      ) : null}

      <div className="form-footer">
        <Button
          type="button"
          disabled={publishing || attachmentUploading}
          onClick={() => void submit()}
        >
          {t("pages.topic.edit.submitBtn")}
        </Button>
      </div>
    </div>
  )
}
