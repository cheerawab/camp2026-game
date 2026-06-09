import { createFileRoute } from "@tanstack/react-router"

import { BattleLobbyPage } from "@/pages/battle/ui/battle-lobby-page"

export const Route = createFileRoute("/battle/")({
  component: BattleLobbyPage,
})
