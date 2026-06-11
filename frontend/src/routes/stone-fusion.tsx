import { createFileRoute } from "@tanstack/react-router"

import { StoneFusionPage } from "@/pages/stone-fusion/ui/stone-fusion-page"

export const Route = createFileRoute("/stone-fusion")({
  component: StoneFusionPage,
})
