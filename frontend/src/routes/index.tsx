import { createFileRoute } from "@tanstack/react-router"

import { helloQueryOptions } from "@/features/home-status/api/hello.query"
import { HomeStatusError } from "@/features/home-status/ui/home-status-error"
import { HomeStatusSkeleton } from "@/features/home-status/ui/home-status-skeleton"
import { HomePage } from "@/pages/home/ui/home-page"

export const Route = createFileRoute("/")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(helloQueryOptions()),
  pendingComponent: HomeStatusSkeleton,
  errorComponent: ({ error }) => <HomeStatusError error={error} />,
  component: HomePage,
})
