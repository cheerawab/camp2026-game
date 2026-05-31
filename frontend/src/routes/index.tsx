import { createFileRoute } from "@tanstack/react-router"

import { healthQueryOptions } from "@/features/home-status/api/health.query"
import { HomeStatusError } from "@/features/home-status/ui/home-status-error"
import { HomeStatusSkeleton } from "@/features/home-status/ui/home-status-skeleton"
import { HomePage } from "@/pages/home/ui/home-page"

export const Route = createFileRoute("/")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(healthQueryOptions()),
  pendingComponent: HomeStatusSkeleton,
  errorComponent: ({ error }) => <HomeStatusError error={error} />,
  component: HomePage,
})
