import { createFileRoute } from "@tanstack/react-router"

import { StoneCollectionPage } from "@/pages/stone-collection/ui/stone-collection-page"

export const Route = createFileRoute("/stones")({
  component: StoneCollectionPage,
})
