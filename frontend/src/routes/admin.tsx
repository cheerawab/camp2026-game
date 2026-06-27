import { createFileRoute } from "@tanstack/react-router"

import { AdminPanelPage } from "@/pages/admin/ui/admin-panel-page"

export const Route = createFileRoute("/admin")({
  component: AdminPanelPage,
})
