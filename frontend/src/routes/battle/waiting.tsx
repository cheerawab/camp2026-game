import { BattleWaitingPage } from "@/pages/battel/waiting/ui/battle-waiting"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/waiting")({
  component: BattleWaitingPage,
})
