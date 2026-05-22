"use client"

import * as React from "react"
import Link from "@/components/common/link"
import { Medal } from "lucide-react"

import { UserAvatar } from "@/components/common/avatar"
import { FollowButton } from "@/components/user/follow-button"
import { BackgroundUploadButton } from "@/components/user/image-upload"
import type { Badge, UserSummary } from "@/lib/api/types"

type UserProfileSummary = UserSummary & {
  level?: number
  smallBackgroundImage?: string
}

function displayName(user: UserProfileSummary) {
  return user.nickname || user.username || `#${user.id}`
}

export function UserProfileCard({
  user,
  currentUser,
  badges = [],
}: {
  user: UserProfileSummary
  currentUser?: UserSummary | null
  badges?: Badge[]
}) {
  const [backgroundImage, setBackgroundImage] = React.useState(
    user.smallBackgroundImage || user.backgroundImage || ""
  )
  const wornBadges = badges.filter((badge) => badge.worn).slice(0, 3)
  const isOwner = currentUser?.id === user.id

  return (
    <section
      className="profile"
      style={
        backgroundImage
          ? { backgroundImage: `url(${backgroundImage})` }
          : undefined
      }
    >
      {isOwner ? (
        <BackgroundUploadButton onUploaded={setBackgroundImage} />
      ) : null}
      <div className="profile-avatar">
        <UserAvatar user={user} size={100} />
      </div>
      <div className="profile-info">
        <div className="metas">
          <div className="nickname-row">
            <span className="nickname">
              <Link
                href={`/user/${user.id}`}
                className="text-foreground hover:underline"
              >
                {displayName(user)}
              </Link>
            </span>
            {user.level !== undefined && user.level !== null ? (
              <span className="level-badge">
                <span className="level-label">Lv</span>
                <span className="tabular-nums">{user.level}</span>
              </span>
            ) : null}
            {wornBadges.map((badge) => (
              <Link
                key={badge.id}
                href={`/user/${user.id}/badges`}
                className="badge-icon shrink-0"
                title={badge.title}
              >
                {badge.icon ? (
                  <img
                    src={badge.icon}
                    alt={badge.title || ""}
                    className="h-6 w-6 object-contain"
                  />
                ) : (
                  <Medal className="h-5 w-5" />
                )}
              </Link>
            ))}
          </div>
          {user.description ? (
            <div className="description">
              <p>{user.description}</p>
            </div>
          ) : null}
        </div>
        <div className="action-btns">
          {isOwner ? null : (
            <FollowButton userId={user.id} initialFollowed={user.followed} />
          )}
        </div>
      </div>
    </section>
  )
}
