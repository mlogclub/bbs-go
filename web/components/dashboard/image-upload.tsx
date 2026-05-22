"use client"

import * as React from "react"
import { LoaderCircleIcon, PencilIcon, PlusIcon } from "lucide-react"

import { apiFetch } from "@/lib/api/client"
import { cn } from "@/lib/utils"

export function DashboardImageUpload({
  value,
  onChange,
  className,
  size = 112,
}: {
  value: string
  onChange: (value: string) => void
  className?: string
  size?: number
}) {
  const [uploading, setUploading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  async function upload(file: File) {
    const form = new FormData()
    form.append("image", file)
    setUploading(true)
    setError(null)
    try {
      const data = await apiFetch<{ url: string }>("/api/upload", {
        method: "POST",
        body: form,
      })
      onChange(data.url)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed")
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className={cn("grid w-fit gap-2", className)}>
      <label
        className={cn(
          "group relative flex cursor-pointer items-center justify-center overflow-hidden rounded-md border border-dashed bg-muted/40 text-sm text-muted-foreground transition-colors hover:border-primary hover:bg-muted",
          uploading && "pointer-events-none opacity-80"
        )}
        style={{ width: size, height: size }}
      >
        {value ? (
          <>
            <img src={value} alt="" className="size-full object-cover" />
            <span className="absolute inset-0 hidden items-center justify-center bg-black/45 text-white group-hover:flex">
              <PencilIcon className="size-5" />
            </span>
          </>
        ) : (
          <span className="flex flex-col items-center gap-2 font-medium">
            <PlusIcon className="size-5" />
            Upload
          </span>
        )}
        {uploading ? (
          <span className="absolute inset-0 flex items-center justify-center bg-background/70">
            <LoaderCircleIcon className="size-6 animate-spin text-primary" />
          </span>
        ) : null}
        <input
          type="file"
          accept="image/*"
          className="sr-only"
          disabled={uploading}
          onChange={(event) => {
            const file = event.target.files?.[0]
            if (file) void upload(file)
            event.currentTarget.value = ""
          }}
        />
      </label>
      {error ? (
        <p className="max-w-40 text-sm text-destructive">{error}</p>
      ) : null}
    </div>
  )
}
