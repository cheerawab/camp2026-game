import { useState } from "react"

import { Button } from "@/shared/ui/button"
import { Card } from "@/shared/ui/card"
import { PageHeader } from "@/shared/ui/page-header"

type Tab = "op" | "stones" | "battle"

type TeamRow = {
  rank: number
  name: string
  value: string
  badge: string
  colorClass: string
  mine?: boolean
}

const RANKS: Record<Tab, TeamRow[]> = {
  op: [
    { rank: 1, name: "山羌小隊", value: "1260 OP", badge: "連勝 4 場",    colorClass: "bg-pebble-spark" },
    { rank: 2, name: "松鼠小隊", value: "1188 OP", badge: "小石 42 顆",   colorClass: "bg-pebble-engineer", mine: true },
    { rank: 3, name: "雲豹小隊", value: "1110 OP", badge: "答題命中 82%", colorClass: "bg-pebble-resonate" },
    { rank: 4, name: "水鹿小隊", value: "980 OP",  badge: "今日新增 90 OP", colorClass: "bg-pebble-explore" },
    { rank: 5, name: "飛鼠小隊", value: "910 OP",  badge: "收藏進度 68%", colorClass: "bg-pebble-play" },
  ],
  stones: [
    { rank: 1, name: "雲豹小隊", value: "58 顆", badge: "本週 +12 顆", colorClass: "bg-pebble-resonate" },
    { rank: 2, name: "松鼠小隊", value: "42 顆", badge: "本週 +7 顆",  colorClass: "bg-pebble-engineer", mine: true },
    { rank: 3, name: "山羌小隊", value: "39 顆", badge: "本週 +5 顆",  colorClass: "bg-pebble-spark" },
    { rank: 4, name: "飛鼠小隊", value: "31 顆", badge: "本週 +3 顆",  colorClass: "bg-pebble-play" },
    { rank: 5, name: "水鹿小隊", value: "24 顆", badge: "本週 +2 顆",  colorClass: "bg-pebble-explore" },
  ],
  battle: [
    { rank: 1, name: "山羌小隊", value: "連勝 4", badge: "勝率 78%", colorClass: "bg-pebble-spark" },
    { rank: 2, name: "水鹿小隊", value: "連勝 2", badge: "勝率 65%", colorClass: "bg-pebble-explore" },
    { rank: 3, name: "松鼠小隊", value: "連勝 1", badge: "勝率 60%", colorClass: "bg-pebble-engineer", mine: true },
    { rank: 4, name: "飛鼠小隊", value: "0 連勝", badge: "勝率 50%", colorClass: "bg-pebble-play" },
    { rank: 5, name: "雲豹小隊", value: "0 連勝", badge: "勝率 44%", colorClass: "bg-pebble-resonate" },
  ],
}

const TABS: { key: Tab; label: string }[] = [
  { key: "op",     label: "開源力" },
  { key: "stones", label: "小石" },
  { key: "battle", label: "知識王戰" },
]

export function LeaderBoardPage() {
  const [activeTab, setActiveTab] = useState<Tab>("op")
  const ranks = RANKS[activeTab]

  return (
    <main className="bg-paper text-ink mx-auto grid w-full max-w-[430px] gap-y-2 px-4 py-4">
      {/* TODO: replace with GET /api/leaderboard */}

      <PageHeader title="小隊排行榜" headline="Leaderboard" />

      {/* 我的排名 */}
      <Card className="bg-ink text-primary-foreground grid grid-cols-[1fr_78px] items-center gap-3 rounded-[26px] px-5 py-4">
        <div>
          <p className="text-primary-foreground/70 mb-1 text-xs font-extrabold tracking-[0.08em] uppercase">
            你的隊伍
          </p>
          <h2 className="mb-1.5 text-2xl leading-tight font-extrabold">
            松鼠小隊目前第 2 名
          </h2>
          <p className="text-primary-foreground/70 text-sm leading-relaxed">
            距離第一名還差 72 OP。下一場知識王戰可以直接追近。
          </p>
        </div>
        <strong className="bg-pebble-spark text-ink border-primary-foreground grid h-[78px] place-items-center rounded-[22px] border-2 text-3xl font-extrabold">
          #2
        </strong>
      </Card>

      {/* 分類 Tabs */}
      <div className="grid grid-cols-3 gap-2" role="tablist" aria-label="排行榜分類">
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

      {/* 排名清單 */}
      <div className="grid gap-2" role="tabpanel">
        {ranks.map((team) => (
          <Card
            key={team.name}
            className={[
              "grid grid-cols-[32px_42px_1fr_auto] items-center gap-3 rounded-[18px] px-3 py-3",
              team.mine ? "border-ink bg-surface-raised shadow-sm" : "",
            ].join(" ")}
          >
            <span className="text-sm font-extrabold">#{team.rank}</span>
            <div
              className={[
                "border-ink size-[42px] rotate-[-6deg] rounded-[14px_18px_12px_16px] border-2",
                team.colorClass,
              ].join(" ")}
              aria-hidden
            />
            <div className="min-w-0">
              <p className="truncate font-bold">{team.name}</p>
              <p className="text-muted-foreground truncate text-xs font-semibold">
                {team.badge}
              </p>
            </div>
            <strong className="text-sm font-extrabold whitespace-nowrap">
              {team.value}
            </strong>
          </Card>
        ))}
      </div>
    </main>
  )
}
