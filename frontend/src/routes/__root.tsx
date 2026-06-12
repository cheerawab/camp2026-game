/// <reference types="vite/client" />

import type { QueryClient } from "@tanstack/react-query"
import {
  HeadContent,
  Outlet,
  Scripts,
  createRootRouteWithContext,
} from "@tanstack/react-router"
import type { ReactNode } from "react"

import { QueryProvider } from "@/app/providers/query-provider"
import { AuthGate } from "@/features/auth/ui/auth-gate"
import appCss from "@/styles/app.css?url"
import { Toaster } from "@/shared/ui/sonner"

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient
}>()({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      {
        title: "Camp 2026 Game",
      },
      {
        name: "description",
        content: "SITCON Camp 2026 game frontend.",
      },
    ],
    links: [{ rel: "stylesheet", href: appCss }],
  }),
  component: RootComponent,
})

function RootComponent() {
  const { queryClient } = Route.useRouteContext()

  return (
    <RootDocument>
      <QueryProvider queryClient={queryClient}>
        <AuthGate>
          <Outlet />
        </AuthGate>
      </QueryProvider>
    </RootDocument>
  )
}

function RootDocument({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="zh-Hant">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Toaster position="bottom-center" />
        <Scripts />
      </body>
    </html>
  )
}
