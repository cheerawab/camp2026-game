import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/result")({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/battle/result"!</div>
}
