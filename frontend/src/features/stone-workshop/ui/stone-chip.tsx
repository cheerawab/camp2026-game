import { Gem } from "lucide-react"

import { type PebbleTone } from "@/shared/config/color-palette"
import { IconBadge } from "@/shared/ui/icon-badge"

type StoneChipProps = {
  label: string
  tone: PebbleTone
}

export function StoneChip({ label, tone }: StoneChipProps) {
  return <IconBadge label={label} tone={tone} icon={<Gem aria-hidden />} />
}
