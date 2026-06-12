import { createFileRoute } from "@tanstack/react-router"

import { LoginPage } from "@/pages/login/ui/login-page"

export const Route = createFileRoute("/login")({
  validateSearch: (search): { token?: string } => {
    if (typeof search.token !== "string") return {}
    return { token: search.token }
  },
  component: LoginRoute,
})

function LoginRoute() {
  const { token } = Route.useSearch()
  return <LoginPage token={token} />
}
