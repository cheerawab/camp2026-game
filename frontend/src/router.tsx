import { createAppRouter } from "@/app/router"

export const getRouter = createAppRouter

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
