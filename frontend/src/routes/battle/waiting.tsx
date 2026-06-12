import { createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/waiting")({
  beforeLoad: () => {
    throw redirect({ to: "/battle/room" })
  },
})
