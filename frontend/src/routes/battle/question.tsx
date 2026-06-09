import { createFileRoute } from "@tanstack/react-router"

import { BattleQuestionPage } from "@/pages/battle/ui/battle-question-page"

export const Route = createFileRoute("/battle/question")({
  component: BattleQuestionPage,
})
