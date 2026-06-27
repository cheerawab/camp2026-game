import { createRouter } from "@tanstack/react-router"
import { routeTree } from "../routeTree.gen"
import { GameErrorPage } from "@/pages/error/ui/game-error-page"
import { createQueryClient } from "@/shared/lib/query-client"

export function createAppRouter() {
  const queryClient = createQueryClient()

  const router = createRouter({
    routeTree,
    context: {
      queryClient,
    },
    defaultPreload: "intent",
    defaultErrorComponent: GameErrorPage,
    scrollRestoration: true,
  })

  return router
}
