import { QueryClientProvider, type QueryClient } from "@tanstack/react-query"
import type { ReactNode } from "react"

type QueryProviderProps = {
  children: ReactNode
  queryClient: QueryClient
}

export function QueryProvider({ children, queryClient }: QueryProviderProps) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
