import { createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/ingame")({
  beforeLoad: () => {
    throw redirect({ to: "/battle/question" })
  },
})
