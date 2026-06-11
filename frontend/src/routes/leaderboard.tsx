import { LeaderboardPage } from "@/pages/leaderboard/ui/leadboard-page"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/leaderboard")({
  component: LeaderboardPage,
})
