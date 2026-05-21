import { serverApiFetch as apiFetch } from "./server"

import type { Comment, PageData } from "./types"

export function getComments(
  entityType: string,
  entityId: string | number,
  cursor?: string
) {
  return apiFetch<PageData<Comment>>("/api/comment/comments", {
    params: { entityType, entityId, cursor },
  })
}
