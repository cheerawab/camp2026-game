import { LeaderboardOp } from "@/features/leaderboard/ui/leaderboard-op"
import { LeaderboardSitone } from "@/features/leaderboard/ui/leaderboard-sitone"
import { Card, CardContent } from "@/shared/ui/card"
import { PageHeader } from "@/shared/ui/page-header"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/shared/ui/tabs"

export function LeaderboardPage() {
  return (
    <main className="mx-auto grid w-full max-w-sm gap-y-2 py-4">
      <PageHeader title="小隊排行榜" headline="Leaderboard" />
      {/* 小隊資訊 */}
      <Card className="bg-accent text-accent-foreground">
        <CardContent className="flex">
          <div className="grid basis-3/4 gap-y-2">
            <span className="text-muted-foreground text-sm font-bold">
              你的隊伍
            </span>
            <span className="text-2xl font-bold">B 小隊目前排名第二！</span>
            <span className="text-wrap wrap-break-word">
              距離第一名的差距不大，僅僅 20 OP，何不嘗試拼一下呢？
            </span>
          </div>
          <div className="flex basis-1/4 items-center justify-center">
            <span className="border-foreground bg-secondary rotate-3 rounded-lg border-2 p-5 text-2xl font-bold">
              #2
            </span>
          </div>
        </CardContent>
      </Card>
      {/* 完整排行榜 */}
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
