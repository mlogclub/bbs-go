import { useCurrentUser } from "@/components/app/app-provider"
import { RequireUser } from "@/components/auth/require-user"
import { PrivateUserCenterPage } from "@/components/user/private-user-center-page"
import type { Favorite, PageData, ScoreLog, UserMessage } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { useDocumentTitle } from "@/lib/use-document-title"

const emptyFavorites: PageData<Favorite> = {
  results: [],
  cursor: "",
  hasMore: false,
}
const emptyMessages: PageData<UserMessage> = {
  results: [],
  cursor: "",
  hasMore: false,
}
const emptyScores: PageData<ScoreLog> = {
  results: [],
  cursor: "",
  hasMore: false,
}

export function PrivateCenter({
  kind,
}: {
  kind: "favorites" | "messages" | "scores"
}) {
  const { t } = useI18n()
  useDocumentTitle(
    t(
      kind === "favorites"
        ? "user.favorites.title"
        : kind === "messages"
          ? "user.messages.title"
          : "user.scores.title"
    )
  )
  const data =
    kind === "favorites"
      ? emptyFavorites
      : kind === "messages"
        ? emptyMessages
        : emptyScores
  const user = useCurrentUser()
  const redirectPath =
    kind === "favorites"
      ? "/user/favorites"
      : kind === "messages"
        ? "/user/messages"
        : "/user/scores"
  return (
    <RequireUser initialUser={user} redirectPath={redirectPath}>
      <PrivateUserCenterPage
        kind={kind}
        initialData={data as never}
        initialBadges={[]}
        initialFans={[]}
        initialFollowed={[]}
        serverLoaded={false}
      />
    </RequireUser>
  )
}
