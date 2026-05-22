import { noindexRouteMeta } from "@/lib/seo"

import { requireUser, requireUserClient } from "../route-helpers/auth"
import { PrivateCenter } from "../route-helpers/private-center"

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
  return noindexRouteMeta(matches, "Messages", "消息")
}

export default function UserMessagesRoute() {
  return <PrivateCenter kind="messages" />
}
