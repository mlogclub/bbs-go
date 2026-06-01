import { noindexRouteMeta } from "@/lib/seo"

import { requireUser, requireUserClient } from "../route-helpers/auth"
import { PrivateCenter } from "../route-helpers/private-center"

async function _loader(args: { request: Request }) {
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
  return noindexRouteMeta(matches, "Scores", "积分")
}

export default function UserScoresRoute() {
  return <PrivateCenter kind="scores" />
}
