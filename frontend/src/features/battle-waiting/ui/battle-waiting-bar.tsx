import { Progress } from "@/shared/ui/progress"

type BattleWaitingBarType = {
  max: number
  value: number
  color: string
}

export function BattleWaitingBar({ max, value, color }: BattleWaitingBarType) {
  return <Progress value={(value / max) * 100} />
}
