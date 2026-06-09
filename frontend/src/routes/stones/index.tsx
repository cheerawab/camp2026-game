import { createFileRoute } from "@tanstack/react-router"

import { StoneCollectionPage } from "@/pages/stones/ui/stone-collection-page"

export const Route = createFileRoute("/stones/")({
  component: StoneCollectionPage,
})
