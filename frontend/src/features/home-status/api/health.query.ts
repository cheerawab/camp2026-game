import { queryOptions } from "@tanstack/react-query"

import { apiClient } from "@/shared/api/client"
import { HealthResponseSchema, type HomeStatus } from "../model/health.schema"

export const healthQueryKey = ["home", "health"] as const

export function healthQueryOptions() {
  return queryOptions({
    queryKey: healthQueryKey,
    queryFn: async (): Promise<HomeStatus> => {
      const json = await apiClient.get("/api/healthz")
      const health = HealthResponseSchema.parse(json)

      return {
        ...health,
        service: "camp2026-game-api",
        checkedAt: new Date().toISOString(),
      }
    },
    staleTime: 30_000,
  })
}
