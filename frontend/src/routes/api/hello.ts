import { createFileRoute } from "@tanstack/react-router"

import { HelloResponseSchema } from "@/features/home-status/model/hello.schema"

export const Route = createFileRoute("/api/hello")({
  server: {
    handlers: {
      GET: async () => {
        const payload = HelloResponseSchema.parse({
          message: "Hello from wired mock API",
          service: "camp2026-game-frontend",
          status: "ok",
          generatedAt: new Date().toISOString(),
        })

        return Response.json(payload)
      },
    },
  },
})
