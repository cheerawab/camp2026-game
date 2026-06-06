import { BattleResultPage } from "@/pages/battel/result/ui/battle-result"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/result")({
  component: BattleResultPage,
})
