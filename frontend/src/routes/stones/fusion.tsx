import { createFileRoute } from "@tanstack/react-router"

import { StoneFusionPage } from "@/pages/stones/ui/stone-fusion-page"

export const Route = createFileRoute("/stones/fusion")({
  component: StoneFusionPage,
})
