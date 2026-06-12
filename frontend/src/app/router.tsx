import { createRouter } from "@tanstack/react-router"
import { routeTree } from "../routeTree.gen"
import { createQueryClient } from "@/shared/lib/query-client"

export function createAppRouter() {
  const queryClient = createQueryClient()

  const router = createRouter({
    routeTree,
    context: {
      queryClient,
    },
    defaultPreload: "intent",
    scrollRestoration: true,
  })

  return router
}
