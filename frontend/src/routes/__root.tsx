/// <reference types="vite/client" />

import { Outlet, createRootRouteWithContext } from "@tanstack/react-router"
import type { QueryClient } from "@tanstack/react-query"

import { QueryProvider } from "@/app/providers/query-provider"
import { AuthGate } from "@/features/auth/ui/auth-gate"
import { Toaster } from "@/shared/ui/sonner"

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient
}>()({
  component: RootComponent,
})

function RootComponent() {
  const { queryClient } = Route.useRouteContext()

  return (
    <QueryProvider queryClient={queryClient}>
      <AuthGate>
        <Outlet />
      </AuthGate>

      <Toaster position="bottom-center" />
    </QueryProvider>
  )
}
