import basicStoneSkin from "@/assets/stones/basic-stone.png"
import { type PebbleTone } from "@/shared/config/color-palette"
import { IconBadge } from "@/shared/ui/icon-badge"

type StoneChipProps = {
  label: string
  tone: PebbleTone
}

export function StoneChip({ label, tone }: StoneChipProps) {
  return (
    <IconBadge
      label={label}
      tone={tone}
      icon={
        <img
          src={basicStoneSkin}
          alt=""
          aria-hidden
          className="size-4 [image-rendering:pixelated]"
        />
      }
    />
  )
}
