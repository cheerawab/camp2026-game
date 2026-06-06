import { Progress } from "@/shared/ui/progress"

type BattleWaitingBarType = {
  max: number
  value: number
}

export function BattleWaitingBar({ max, value }: BattleWaitingBarType) {
  return <Progress value={(value / max) * 100} />
}
