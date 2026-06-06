import { BattleIndexPage } from "@/pages/battle/index/ui/battle-index"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/")({
  component: BattleIndexPage,
})
