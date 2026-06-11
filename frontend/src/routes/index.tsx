import { createFileRoute } from "@tanstack/react-router"

import { HomeBasePage } from "@/pages/home/ui/home-base-page"

export const Route = createFileRoute("/")({
  component: HomeBasePage,
})
