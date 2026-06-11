import { Badge } from "@/shared/ui/badge"
import { Card, CardContent } from "@/shared/ui/card"
import { Check, Loader2 } from "lucide-react"

type BattleWaitingPlayerCardType = {
  name: string
  team: string
  ready: boolean
  pictureSrc: string
}

export function BattleWaitingPlayerCard({
  name,
  team,
  ready,
  pictureSrc,
}: BattleWaitingPlayerCardType) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between">
        <img src={pictureSrc} className="h-20" />
        <div className="flex flex-col items-center justify-center">
          <span className="text-2xl font-bold">{name}</span>
          <span className="text-muted-foreground text-lg">{team}</span>
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
