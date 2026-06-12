import { StaffRewardsPanel } from "@/features/staff-rewards"
import { GamePageShell } from "@/shared/ui/game-page-shell"
import { PageHeader } from "@/shared/ui/page-header"

export function StaffRewardsPage() {
  return (
    <GamePageShell
      ariaLabel="工作人員發放頁"
      contentClassName="grid content-start gap-y-3"
    >
      <PageHeader title="工作人員發放" headline="Staff Rewards" />
      <StaffRewardsPanel />
    </GamePageShell>
  )
}
