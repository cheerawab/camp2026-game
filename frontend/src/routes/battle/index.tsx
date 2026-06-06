import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/battle/")({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/battle/"!</div>
}
