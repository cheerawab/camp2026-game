import { BattleIngamePage } from "@/pages/battle/ingame/ui/battle-ingame"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/ingame")({
  component: BattleIngamePage,
})
