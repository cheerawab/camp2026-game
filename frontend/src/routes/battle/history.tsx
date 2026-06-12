import { createFileRoute } from "@tanstack/react-router"

import { BattleHistoryPage } from "@/pages/battle/ui/battle-history-page"

export const Route = createFileRoute("/battle/history")({
  component: BattleHistoryPage,
})
