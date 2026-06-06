import { BattleIndexPage } from "@/pages/battel/index/ui/battle-index"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/")({
  component: BattleIndexPage,
})
