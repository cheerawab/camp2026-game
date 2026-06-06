import { createFileRoute } from "@tanstack/react-router"

import { ComponentsPage } from "@/pages/components/ui/components-page"

export const Route = createFileRoute("/component-gallery")({
  component: ComponentsPage,
})
