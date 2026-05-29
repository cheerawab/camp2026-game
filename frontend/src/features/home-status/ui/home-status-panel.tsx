import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query"

import { helloQueryKey, helloQueryOptions } from "../api/hello.query"
import { HomeStatusCard } from "./home-status-card"

export function HomeStatusPanel() {
  const queryClient = useQueryClient()
  const { data, isFetching } = useSuspenseQuery(helloQueryOptions())

  return (
    <HomeStatusCard
      data={data}
      isFetching={isFetching}
      onRefresh={() =>
        queryClient.invalidateQueries({ queryKey: helloQueryKey })
      }
    />
  )
}
