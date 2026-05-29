import { QueryClient } from "@tanstack/react-query"

import { AppError } from "@/shared/api/error"

export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        gcTime: 5 * 60_000,
        refetchOnWindowFocus: false,
        retry: (failureCount, error) =>
          error instanceof AppError && error.retryable && failureCount < 2,
        throwOnError: (error) =>
          error instanceof AppError && error.status >= 500,
      },
      mutations: {
        retry: 0,
      },
    },
  })
}
