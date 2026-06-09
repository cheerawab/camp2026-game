import { createFileRoute } from "@tanstack/react-router"

import { InventoryPage } from "@/pages/inventory/ui/inventory-page"

export const Route = createFileRoute("/inventory")({
  component: InventoryPage,
})
