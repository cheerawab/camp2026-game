import { createFileRoute } from "@tanstack/react-router"

import { LeaderBoardPage } from "@/pages/leaderboard/ui/leader-board-page"

export const Route = createFileRoute("/leaderboard")({
  component: LeaderBoardPage,
})
