import { createFileRoute } from "@tanstack/react-router"

import { BattleResultPage } from "@/pages/battle/ui/battle-result-page"

export const Route = createFileRoute("/battle/result")({
  component: BattleResultPage,
})
