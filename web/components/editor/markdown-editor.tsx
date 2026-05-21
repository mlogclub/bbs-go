"use client"

import * as React from "react"
import dynamic from "@/lib/router/dynamic"
import type { ToolbarNames } from "md-editor-rt"

import "md-editor-rt/lib/style.css"

import { useTheme } from "@/components/theme-provider"
import { uploadEditorImage } from "@/components/editor/upload"
import { useI18n } from "@/lib/i18n/provider"

const MdEditor = dynamic(
  () => import("md-editor-rt").then((mod) => mod.MdEditor),
  { ssr: false }
)

const TOOLBARS = [
  "bold",
  "underline",
  "italic",
  "strikeThrough",
  "-",
  "title",
  "sub",
  "sup",
  "quote",
  "unorderedList",
  "orderedList",
  "task",
  "-",
  "codeRow",
  "code",
  "link",
  "image",
  "table",
  "-",
  "revoke",
  "next",
  "-",
  "preview",
  "catalog",
  "=",
  "fullscreen",
] satisfies ToolbarNames[]

export function MarkdownEditor({
  value,
  placeholder,
  height = "400px",
  onChange,
}: {
  value: string
  placeholder?: string
  height?: string
  onChange: (value: string) => void
}) {
  const { resolvedTheme } = useTheme()
  const { locale } = useI18n()

  async function uploadImg(files: File[], callback: (urls: string[]) => void) {
    const urls = await Promise.all(files.map((file) => uploadEditorImage(file)))
    callback(urls)
  }

  return (
    <MdEditor
      modelValue={value}
      theme={resolvedTheme === "dark" ? "dark" : "light"}
      toolbars={TOOLBARS}
      style={{ height }}
      placeholder={placeholder}
      preview
      language={locale}
      footers={[]}
      onChange={onChange}
      onUploadImg={uploadImg}
    />
  )
}
