import { createFileRoute } from "@tanstack/react-router"

import { PublicCodexPage } from "@/pages/codex/ui/public-codex-page"

export const Route = createFileRoute("/codex")({
  component: PublicCodexPage,
})
