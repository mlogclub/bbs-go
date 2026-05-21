"use client"

import * as React from "react"
import { useRouter } from "@/lib/router/navigation"
import { Plus, Trash2 } from "lucide-react"

import { TagInput } from "@/components/common/tag-input"
import { ContentEditor } from "@/components/editor/content-editor"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { apiFetch, toFormData } from "@/lib/api/client"
import type { Article, ArticleEditForm, ImageInfo, SiteConfig } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { msg, useToastActions } from "@/lib/toast"

export { TagInput } from "@/components/common/tag-input"

type ArticleFormMode = "create" | "edit"

type ArticleFormState = {
  title: string
  content: string
  tags: string[]
  cover: ImageInfo[]
}

function isSameTags(left: string[], right: string[]) {
  return left.length === right.length && left.every((item, index) => item === right[index])
}

async function uploadImage(file: File) {
  const body = new FormData()
  body.append("image", file, file.name)
  return apiFetch<ImageInfo>("/api/upload", { method: "POST", body })
}

function CoverUpload({
  value,
  label,
  onChange,
}: {
  value: ImageInfo[]
  label: string
  onChange: (value: ImageInfo[]) => void
}) {
  const { t } = useI18n()
  const { catchError, msgWarning } = useToastActions()
  const inputRef = React.useRef<HTMLInputElement>(null)
  const [uploading, setUploading] = React.useState(false)

  async function onFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.[0]
    if (!file) {
      return
    }

    if (value.length >= 1) {
      msgWarning(t("component.imageUpload.countLimitError", { limit: 1 }))
      event.currentTarget.value = ""
      return
    }

    setUploading(true)
    try {
      const uploaded = await uploadImage(file)
      onChange([{ ...uploaded, name: uploaded.name || file.name, size: uploaded.size || file.size }])
    } catch (error) {
      catchError(error)
    } finally {
      setUploading(false)
      event.currentTarget.value = ""
    }
  }

  return (
    <div className="flex flex-wrap gap-2.5">
      {value.map((image, index) => (
        <div key={`${image.url || image.preview || index}`} className="group relative h-[120px] w-[120px] overflow-hidden rounded-sm border">          <img src={image.url || image.preview} alt="" className="h-full w-full object-cover" />
          <button
            type="button"
            className="absolute bottom-0 left-0 hidden h-5 w-full items-center justify-center bg-black/30 text-white group-hover:flex"
            onClick={() => onChange(value.filter((_, itemIndex) => itemIndex !== index))}
          >
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      ))}
      {value.length < 1 ? (
        <button
          type="button"
          className="flex h-[120px] w-[120px] flex-col items-center justify-center rounded-sm border text-muted-foreground"
          disabled={uploading}
          onClick={() => inputRef.current?.click()}
        >
          <Plus className="h-[18px] w-[18px]" />
          <span className="text-sm font-medium">{uploading ? t("component.imageUpload.uploading") : label}</span>
        </button>
      ) : null}
      <input ref={inputRef} type="file" accept="image/*" className="hidden" onChange={onFileChange} />
    </div>
  )
}

export function ArticleForm({
  mode,
  config,
  initialArticle,
}: {
  mode: ArticleFormMode
  config: SiteConfig | null
  initialArticle?: ArticleEditForm | null
}) {
  const router = useRouter()
  const { t } = useI18n()
  const { catchError } = useToastActions()
  const [publishing, setPublishing] = React.useState(false)
  const [form, setForm] = React.useState<ArticleFormState>({
    title: initialArticle?.title || "",
    content: initialArticle?.content || "",
    tags: initialArticle?.tags || [],
    cover: initialArticle?.cover ? [initialArticle.cover] : [],
  })

  function updateForm(next: Partial<ArticleFormState>) {
    setForm((current) => ({ ...current, ...next }))
  }

  async function submit() {
    if (publishing) {
      return
    }

    setPublishing(true)
    try {
      const cover = form.cover.length ? JSON.stringify(form.cover[0]) : null
      const body = toFormData({
        title: form.title,
        content: form.content,
        tags: form.tags.length ? form.tags.join(",") : "",
        cover,
      })

      if (mode === "create") {
        const article = await apiFetch<Article>("/api/article/create", {
          method: "POST",
          body,
        })
        msg({
          message: t("pages.article.create.success"),
          onClose() {
            router.push(`/article/${article.id}`)
          },
        })
      } else if (initialArticle?.id) {
        await apiFetch<{ articleId: number }>(`/api/article/edit/${initialArticle.id}`, {
          method: "POST",
          body,
        })
        msg({
          message: t("pages.article.edit.editSuccess"),
          onClose() {
            router.push(`/article/${initialArticle.id}`)
          },
        })
      }
    } catch (error) {
      catchError(error)
      setPublishing(false)
    }
  }

  const titleKey = mode === "create" ? "pages.article.create.title" : "pages.article.edit.title"
  const titlePlaceholderKey = mode === "create" ? "pages.article.create.titlePlaceholder" : "pages.article.edit.titlePlaceholder"
  const contentPlaceholderKey = mode === "create" ? "pages.article.create.contentPlaceholder" : "pages.article.edit.contentPlaceholder"
  const submitKey = mode === "create" ? "pages.article.create.publishBtn" : "pages.article.edit.submitBtn"

  return (
    <div className="publish-form rounded-lg bg-background p-4">
      <div className="mb-4 border-b pb-2">
        <div className="text-xl font-semibold">{t(titleKey)}</div>
      </div>
      <div className="mb-4">
        <Input
          value={form.title}
          type="text"
          placeholder={t(titlePlaceholderKey)}
          onChange={(event) => updateForm({ title: event.currentTarget.value })}
        />
      </div>
      <div className="mb-4">
        <ContentEditor
          contentType="markdown"
          value={form.content}
          placeholder={t(contentPlaceholderKey)}
          height="400px"
          onChange={(content) => updateForm({ content })}
        />
      </div>
      <div className="mb-4">
        <TagInput
          value={form.tags}
          recommendTags={config?.recommendTags}
          placeholder={t("component.tagInput.placeholder")}
          onChange={(tags) => {
            if (!isSameTags(tags, form.tags)) {
              updateForm({ tags })
            }
          }}
        />
      </div>
      {mode === "create" ? (
        <div className="mb-4">
          <CoverUpload value={form.cover} label={t("pages.article.create.cover")} onChange={(cover) => updateForm({ cover })} />
        </div>
      ) : null}
      <div className="pt-2">
        <Button type="button" disabled={publishing} onClick={submit}>
          {publishing && mode === "create" ? t("pages.article.create.publishing") : t(submitKey)}
        </Button>
      </div>
    </div>
  )
}
