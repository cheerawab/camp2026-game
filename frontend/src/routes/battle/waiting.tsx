import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/battle/waiting')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/battle/waiting"!</div>
}
