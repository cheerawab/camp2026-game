import { BattleResultPage } from "@/pages/battle/result/ui/battle-result"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/result")({
  component: BattleResultPage,
})
