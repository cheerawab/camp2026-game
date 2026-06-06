import { BattleIngamePage } from "@/pages/battel/ingame/ui/battle-ingame"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/ingame")({
  component: BattleIngamePage,
})
