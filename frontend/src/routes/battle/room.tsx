import { createFileRoute } from "@tanstack/react-router"

import { BattleWaitingRoomPage } from "@/pages/battle/ui/battle-waiting-room-page"

export const Route = createFileRoute("/battle/room")({
  component: BattleWaitingRoomPage,
})
