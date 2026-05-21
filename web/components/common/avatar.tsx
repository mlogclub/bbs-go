import * as React from "react"

import Link from "@/components/common/link"
import { cn } from "@/lib/utils"

import type { UserSummary } from "@/lib/api/types"

function getDisplayName(user?: UserSummary | null) {
  return user?.nickname || user?.username || ""
}

function getInitial(user?: UserSummary | null) {
  const name = getDisplayName(user).trim()
  return name ? name.slice(0, 1).toUpperCase() : "U"
}

export function UserAvatar({
  user,
  size = 32,
  className,
  linkToProfile = true,
  target = "_blank",
}: {
  user?: UserSummary | null
  size?: number
  className?: string
  linkToProfile?: boolean
  target?: React.HTMLAttributeAnchorTarget
}) {
  const src = user?.smallAvatar || user?.avatar
  const name = getDisplayName(user) || "User"
  const imageRef = React.useRef<HTMLImageElement>(null)
  const [imageLoaded, setImageLoaded] = React.useState(false)
  const [imageFailed, setImageFailed] = React.useState(false)

  React.useEffect(() => {
    setImageLoaded(false)
    setImageFailed(false)
  }, [src])

  React.useEffect(() => {
    const image = imageRef.current
    if (!image?.complete) {
      return
    }
    if (image.naturalWidth > 0) {
      setImageLoaded(true)
    } else {
      setImageFailed(true)
    }
  }, [src, imageLoaded, imageFailed])

  const content =
    src && !imageFailed ? (
      <>
        {!imageLoaded ? (
          <span aria-hidden="true">{getInitial(user)}</span>
        ) : null}
        <img
          ref={imageRef}
          src={src}
          alt={name}
          className={cn("h-full w-full object-cover", !imageLoaded && "hidden")}
          referrerPolicy="no-referrer"
          onLoad={() => setImageLoaded(true)}
          onError={() => setImageFailed(true)}
        />
      </>
    ) : (
      <span aria-hidden="true">{getInitial(user)}</span>
    )
  const rootClassName = cn(
    "inline-flex shrink-0 items-center justify-center overflow-hidden rounded-full bg-muted text-sm font-medium text-muted-foreground",
    className,
  )
  const rootStyle = { width: size, height: size }

  if (linkToProfile && user?.id) {
    return (
      <Link
        href={`/user/${user.id}`}
        target={target}
        rel={target === "_blank" ? "noopener noreferrer" : undefined}
        className={rootClassName}
        style={rootStyle}
        aria-label={name}
      >
        {content}
      </Link>
    )
  }

  return (
    <span className={rootClassName} style={rootStyle} aria-label={name}>
      {content}
    </span>
  )
}
