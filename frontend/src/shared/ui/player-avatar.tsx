import { createAvatar } from "@dicebear/core"
import * as thumbs from "@dicebear/thumbs"
import * as React from "react"

import { Avatar, AvatarImage } from "@/shared/ui/avatar"
import { cn } from "@/shared/utils"

function avatarSeed(playerId?: string, nickname?: string) {
  const id = playerId?.trim()
  if (id) return id

  const name = nickname?.trim()
  if (name) return name

  return "camp2026-player"
}

function svgDataUrl(svg: string) {
  return `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`
}

type PlayerAvatarProps = Omit<
  React.ComponentProps<typeof Avatar>,
  "children"
> & {
  playerId?: string
  nickname?: string
  svgClassName?: string
}

export function PlayerAvatar({
  playerId,
  nickname,
  className,
  svgClassName,
  "aria-label": ariaLabel,
  ...props
}: PlayerAvatarProps) {
  const seed = avatarSeed(playerId, nickname)
  const avatarSrc = React.useMemo(
    () => svgDataUrl(createAvatar(thumbs, { seed }).toString()),
    [seed],
  )

  return (
    <Avatar
      aria-hidden={ariaLabel ? undefined : true}
      aria-label={ariaLabel}
      role={ariaLabel ? "img" : undefined}
      className={cn("bg-surface-raised", className)}
      {...props}
    >
      <AvatarImage
        src={avatarSrc}
        alt=""
        aria-hidden="true"
        draggable={false}
        className={cn(
          "block size-full object-cover",
          svgClassName,
        )}
      />
    </Avatar>
  )
}
