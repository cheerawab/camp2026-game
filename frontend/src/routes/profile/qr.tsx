import { createFileRoute } from "@tanstack/react-router"

import { ProfileQrPage } from "@/pages/profile/ui/profile-qr-page"

export const Route = createFileRoute("/profile/qr")({
  component: ProfileQrPage,
})
