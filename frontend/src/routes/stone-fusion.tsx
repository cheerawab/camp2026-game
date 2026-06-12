import { createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/stone-fusion")({
  beforeLoad: () => {
    throw redirect({ to: "/stones/fusion" })
  },
})
