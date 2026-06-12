import { Badge } from "@/shared/ui/badge"
import { Card, CardContent } from "@/shared/ui/card"
import { Check, Loader2 } from "lucide-react"

type BattleWaitingPlayerCardType = {
  name: string
  team: string
  ready: boolean
  loadoutCount?: number
  pictureSrc?: string
}

export function BattleWaitingPlayerCard({
  name,
  team,
  ready,
  loadoutCount = 0,
  pictureSrc,
}: BattleWaitingPlayerCardType) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between">
        {pictureSrc ? (
          <img src={pictureSrc} className="h-20" alt="" />
        ) : (
          <div className="bg-pebble-spark border-ink grid size-20 place-items-center rounded-[22px] border-2 text-2xl font-black">
            {name.trim().slice(0, 1) || "?"}
          </div>
        )}
        <div className="flex flex-col items-center justify-center">
          <span className="text-2xl font-bold">{name}</span>
          <span className="text-muted-foreground text-lg">{team}</span>
          <span className="text-muted-foreground text-sm font-bold">
            {loadoutCount > 0 ? `${loadoutCount} 顆小石` : "尚未選小石"}
          </span>
        </div>
        <Badge
          variant={ready ? "default" : "outline"}
          className={ready ? "bg-status-success" : "animate-pulse"}
        >
          {ready ? <Check /> : <Loader2 className="animate-spin" />}
          {ready ? "已準備" : "等待中"}
        </Badge>
      </CardContent>
    </Card>
  )
}
