"use client"

import * as React from "react"
import { CloudUpload, Upload } from "lucide-react"

import { Button } from "@/components/ui/button"
import { apiFetch, toFormData } from "@/lib/api/client"
import { useI18n } from "@/lib/i18n/provider"
import { toast } from "@/lib/toast"

async function uploadImage(file: File) {
  const body = new FormData()
  body.append("image", file, file.name)
  return apiFetch<{ url: string }>("/api/upload", { method: "POST", body })
}

export function AvatarEdit({
  value,
  onChange,
}: {
  value?: string
  onChange?: (url: string) => void
}) {
  const { t } = useI18n()
  const inputRef = React.useRef<HTMLInputElement>(null)
  const [avatar, setAvatar] = React.useState(value || "")
  const [uploading, setUploading] = React.useState(false)

  async function onFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0]
    if (!file) return

    setUploading(true)
    try {
      const result = await uploadImage(file)
      await apiFetch<null>("/api/user/update_avatar", {
        method: "POST",
        body: toFormData({ avatar: result.url }),
      })
      setAvatar(result.url)
      onChange?.(result.url)
      toast.success(t("component.avatarEdit.updateSuccess"))
    } catch {
      toast.error(t("component.avatarEdit.updateFailed"))
    } finally {
      setUploading(false)
      event.currentTarget.value = ""
    }
  }

  return (
    <div className="avatar-edit">
      <button
        type="button"
        className="avatar-view"
        style={avatar ? { backgroundImage: `url(${avatar})` } : undefined}
        disabled={uploading}
        onClick={() => inputRef.current?.click()}
      >
        <span className="upload-view">
          <Upload size="20" />
          <span>{t("component.avatarEdit.update")}</span>
        </span>
      </button>
      <input
        ref={inputRef}
        accept="image/*"
        type="file"
        onChange={onFileChange}
      />
    </div>
  )
}

export function BackgroundUploadButton({
  onUploaded,
}: {
  onUploaded?: (url: string) => void
}) {
  const { t } = useI18n()
  const inputRef = React.useRef<HTMLInputElement>(null)
  const [uploading, setUploading] = React.useState(false)

  async function onFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0]
    if (!file) return

    setUploading(true)
    try {
      const result = await uploadImage(file)
      await apiFetch<null>("/api/user/set_background_image", {
        method: "POST",
        body: toFormData({ backgroundImage: result.url }),
      })
      onUploaded?.(result.url)
      toast.success(t("component.userProfile.backgroundSuccess"))
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : t("composables.unknownError")
      )
    } finally {
      setUploading(false)
      event.currentTarget.value = ""
    }
  }

  return (
    <Button
      type="button"
      className="change-bg"
      disabled={uploading}
      onClick={() => inputRef.current?.click()}
    >
      <CloudUpload size="16" />
      <span>{t("component.userProfile.setBackground")}</span>
      <input
        ref={inputRef}
        accept="image/*"
        type="file"
        className="hidden"
        onChange={onFileChange}
      />
    </Button>
  )
}
