import { queryOptions } from "@tanstack/react-query"

import { apiClient } from "@/shared/api/client"
import { HelloResponseSchema, type HelloResponse } from "../model/hello.schema"

export const helloQueryKey = ["home", "hello"] as const

export function helloQueryOptions() {
  return queryOptions({
    queryKey: helloQueryKey,
    queryFn: async (): Promise<HelloResponse> => {
      const json = await apiClient.get("/api/hello")
      return HelloResponseSchema.parse(json)
    },
    staleTime: 30_000,
  })
}
