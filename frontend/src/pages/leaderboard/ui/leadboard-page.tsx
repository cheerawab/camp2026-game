import { LeaderboardOp } from "@/features/leaderboard/ui/leaderboard-op"
import { LeaderboardSitone } from "@/features/leaderboard/ui/leaderboard-sitone"
import { PageHeader } from "@/shared/ui/page-header"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/shared/ui/tabs"

export function LeaderboardPage() {
  return (
    <main className="mx-auto w-full max-w-sm gap-y-2 py-4">
      <PageHeader title="小隊排行榜" headline="Leaderboard" />
      <Tabs defaultValue="op">
        <TabsList className="w-full">
          <TabsTrigger value="op">開源力</TabsTrigger>
          <TabsTrigger value="sitone">小石</TabsTrigger>
        </TabsList>
        <TabsContent value="op">
          <LeaderboardOp />
        </TabsContent>
        <TabsContent value="sitone">
          <LeaderboardSitone />
        </TabsContent>
      </Tabs>
    </main>
  )
}
