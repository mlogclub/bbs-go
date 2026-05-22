import { useCurrentUser } from "@/components/app/app-provider"
import { MainShell } from "@/components/layout/main-shell"
import { TasksPageContent } from "@/components/tasks/tasks-page"
import { CheckInCard, TasksUserCard } from "@/components/tasks/task-widgets"
import { useI18n } from "@/lib/i18n/provider"
import { localizedTitle, noindexMeta, rootDataFromMatches } from "@/lib/seo"
import { useDocumentTitle } from "@/lib/use-document-title"

import { requireUser, requireUserClient } from "../route-helpers/auth"
import { useMediaQuery } from "../route-helpers/client-hooks"

export async function loader(args: { request: Request }) {
  await requireUser(args)
  return null
}

export async function clientLoader(args: { request: Request }) {
  await requireUserClient(args)
  return null
}

export function meta({
  matches,
}: {
  matches: Array<{ data?: unknown; loaderData?: unknown }>
}) {
  const rootData = rootDataFromMatches(matches)
  return noindexMeta(
    rootData?.config,
    localizedTitle(rootData?.locale, "Tasks", "任务")
  )
}

export default function TasksRoute() {
  const { t } = useI18n()
  useDocumentTitle(t("user.tasks.title"))
  const user = useCurrentUser()
  const showInlineAside = useMediaQuery("(max-width: 1279px)")
  const aside = (
    <>
      <TasksUserCard user={user} badges={[]} />
      <CheckInCard initialCheckIn={null} initialRank={[]} />
    </>
  )
  return (
    <MainShell sideSize="320" aside={showInlineAside ? undefined : aside}>
      {showInlineAside ? <div className="space-y-4">{aside}</div> : null}
      <TasksPageContent groups={[]} tasks={[]} />
    </MainShell>
  )
}
