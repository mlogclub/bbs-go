import type { TFunction } from "@/lib/i18n"

export type EditorMode = "html" | "markdown"

export function getEditorModeOptions(t: TFunction) {
  return [
    { value: "html" as const, label: t("component.editorMode.visual") },
    { value: "markdown" as const, label: t("component.editorMode.markdown") },
  ]
}

export function getEditorSwitchTarget(contentType: EditorMode): EditorMode {
  return contentType === "markdown" ? "html" : "markdown"
}

export function getEditorSwitchConfirmMessage(
  contentType: EditorMode,
  t: TFunction
) {
  const target = getEditorSwitchTarget(contentType)

  return t("component.editorMode.switchConfirm", {
    mode:
      target === "markdown"
        ? t("component.editorMode.markdown")
        : t("component.editorMode.visual"),
  })
}
