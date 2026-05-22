"use client"

import Link from "@/components/common/link"

import { UserAvatar } from "@/components/common/avatar"
import { FollowButton } from "@/components/user/follow-button"
import type { UserSummary } from "@/lib/api/types"

export function UserFollowList({ users }: { users: UserSummary[] }) {
  if (!users.length) {
    return null
  }

  return (
    <div>
      {users.map((item) => (
        <div key={item.id} className="user-follow-item">
          <UserAvatar user={item} size={40} />
          <div className="user-follow-item-info">
            <div className="nickname">
              <Link href={`/user/${item.id}`}>
                {item.nickname || item.username || `#${item.id}`}
              </Link>
            </div>
            <div className="description">{item.description}</div>
          </div>
          <FollowButton userId={item.id} initialFollowed={item.followed} />
        </div>
      ))}
    </div>
  )
}
