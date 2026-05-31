import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query"

import { healthQueryKey, healthQueryOptions } from "../api/health.query"
import { HomeStatusCard } from "./home-status-card"

export function HomeStatusPanel() {
  const queryClient = useQueryClient()
  const { data, isFetching } = useSuspenseQuery(healthQueryOptions())

  return (
    <HomeStatusCard
      data={data}
      isFetching={isFetching}
      onRefresh={() =>
        queryClient.invalidateQueries({ queryKey: healthQueryKey })
      }
    />
  )
}
