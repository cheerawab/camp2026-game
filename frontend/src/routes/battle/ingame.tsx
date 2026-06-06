import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/battle/ingame')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/battle/ingame"!</div>
}
