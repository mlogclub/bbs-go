import { serverApiFetch as apiFetch } from "./server"

import type { CheckInInfo, TaskGroupInfo, TaskInfo } from "./types"

export function getTaskGroups() {
  return apiFetch<TaskGroupInfo[]>("/api/task/groups")
}

export function getTasks(groupName?: string) {
  return apiFetch<TaskInfo[]>("/api/task/tasks", {
    params: { groupName },
  })
}

export function getCheckIn() {
  return apiFetch<CheckInInfo | null>("/api/checkin/checkin")
}

export function getCheckInRank() {
  return apiFetch<CheckInInfo[]>("/api/checkin/rank")
}
