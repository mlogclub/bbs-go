import { serverApiFetch as apiFetch } from "./server"

export interface AboutConfig {
  content?: string
}

export interface FriendLink {
  id: string | number
  title?: string
  url?: string
  summary?: string
}

export interface InstallStatus {
  installed?: boolean
}

export function getAboutConfig() {
  return apiFetch<AboutConfig>("/api/config/about")
}

export function getFriendLinks() {
  return apiFetch<FriendLink[]>("/api/link/list")
}

export function getTopFriendLinks() {
  return apiFetch<FriendLink[]>("/api/link/top_links")
}

export function getInstallStatus() {
  return apiFetch<InstallStatus>("/api/install/status")
}
