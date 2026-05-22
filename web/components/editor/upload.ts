import { apiFetch } from "@/lib/api/client"

export async function uploadEditorImage(file: File) {
  const body = new FormData()
  body.append("image", file, file.name)
  const result = await apiFetch<{ url: string }>("/api/upload", {
    method: "POST",
    body,
  })
  return result.url
}
