import { useQuery } from "@tanstack/react-query"
import { useState } from "react"

import { gameApi, type LeaderboardType } from "@/shared/api/game"
import { Button } from "@/shared/ui/button"
import { Card } from "@/shared/ui/card"
import { GamePageShell } from "@/shared/ui/game-page-shell"
import { PageHeader } from "@/shared/ui/page-header"

const TABS: { key: LeaderboardType; label: string }[] = [
  { key: "open_power", label: "開源力" },
  { key: "sitones", label: "小石" },
]

const rowColors = [
  "bg-pebble-spark",
  "bg-pebble-engineer",
  "bg-pebble-resonate",
  "bg-pebble-explore",
  "bg-pebble-play",
]

export function LeaderBoardPage() {
  const [activeTab, setActiveTab] = useState<LeaderboardType>("open_power")
  const { data, isPending } = useQuery({
    queryKey: ["leaderboards", activeTab],
    queryFn: () => gameApi.leaderboard(activeTab),
  })
  const ranks = data?.teams ?? []
  const currentTeam = data?.currentTeam

  return (
    <GamePageShell contentClassName="grid content-start gap-y-2">
      <PageHeader title="小隊排行榜" headline="Leaderboard" />

      <Card className="bg-ink text-primary-foreground grid grid-cols-[1fr_78px] items-center gap-3 rounded-[26px] px-5 py-4">
        <div>
          <p className="text-primary-foreground/70 mb-1 text-xs font-extrabold tracking-[0.08em] uppercase">
            你的隊伍
          </p>
          <h2 className="mb-1.5 text-2xl leading-tight font-extrabold">
            {currentTeam
              ? `${currentTeam.name} 目前第 ${currentTeam.rank} 名`
              : isPending
                ? "正在同步排名"
                : "目前沒有隊伍排名"}
          </h2>
          <p className="text-primary-foreground/70 text-sm leading-relaxed">
            {currentTeam && data
              ? data.gapToPrevious > 0
                ? `距離前一名還差 ${data.gapToPrevious} ${currentTeam.metric}。`
                : "目前已經在這個分類領先或並列領先。"
              : "完成活動與知識王戰後，排行榜會自動更新。"}
          </p>
        </div>
        <strong className="bg-pebble-spark text-ink border-primary-foreground grid h-[78px] place-items-center rounded-[22px] border-2 text-3xl font-extrabold">
          {currentTeam ? `#${currentTeam.rank}` : "-"}
        </strong>
      </Card>

      <div
        className="grid grid-cols-2 gap-2"
        role="tablist"
        aria-label="排行榜分類"
      >
        {TABS.map(({ key, label }) => (
          <Button
            key={key}
            role="tab"
            aria-selected={activeTab === key}
            variant={activeTab === key ? "default" : "outline"}
            className="min-h-11 rounded-2xl font-extrabold"
            onClick={() => setActiveTab(key)}
          >
            {label}
          </Button>
        ))}
      </div>

      <div className="grid gap-2" role="tabpanel">
        {isPending ? (
          <Card className="rounded-[18px] px-3 py-4">
            <span className="text-muted-foreground text-sm font-extrabold">
              正在同步排行榜
            </span>
          </Card>
        ) : ranks.length > 0 ? (
          ranks.map((team, index) => (
            <Card
              key={team.teamId}
              className={[
                "grid grid-cols-[32px_42px_1fr_auto] items-center gap-3 rounded-[18px] px-3 py-3",
                team.current ? "border-ink bg-surface-raised shadow-sm" : "",
              ].join(" ")}
            >
              <span className="text-sm font-extrabold">#{team.rank}</span>
              <div
                className={[
                  "border-ink size-[42px] rotate-[-6deg] rounded-[14px_18px_12px_16px] border-2",
                  rowColors[index % rowColors.length],
                ].join(" ")}
                aria-hidden
              />
              <div className="min-w-0">
                <p className="truncate font-bold">{team.name}</p>
                <p className="text-muted-foreground truncate text-xs font-semibold">
                  {team.current ? "你的隊伍" : "小隊總分"}
                </p>
              </div>
              <strong className="text-sm font-extrabold whitespace-nowrap">
                {team.score} {team.metric}
              </strong>
            </Card>
          ))
        ) : (
          <Card className="rounded-[18px] px-3 py-4">
            <span className="text-muted-foreground text-sm font-extrabold">
              目前沒有排行榜資料
            </span>
          </Card>
        )}
      </div>
    </GamePageShell>
  )
}
