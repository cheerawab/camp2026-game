import { useQuery } from "@tanstack/react-query"
import { useState } from "react"

import { gameApi, type Item, type Sitone } from "@/shared/api/game"
import {
  itemTypeClass,
  itemTypeLabel,
  rarityLabel,
  sitoneMeta,
} from "@/shared/lib/game-labels"
import { Button } from "@/shared/ui/button"
import { Card, CardContent } from "@/shared/ui/card"

type Tab = "stones" | "items"

type CodexEntry =
  | { kind: "stones"; data: Sitone }
  | { kind: "items"; data: Item }

function entryMeta(entry: CodexEntry) {
  if (entry.kind === "stones") {
    const meta = sitoneMeta(entry.data.type)
    return {
      name: entry.data.name,
      type: meta.label,
      rarity: rarityLabel(entry.data.rarity),
      toneClass: meta.bgClassName,
      description: entry.data.description,
      symbol: "◆",
    }
  }

  return {
    name: entry.data.name,
    type: itemTypeLabel(entry.data.type),
    rarity: rarityLabel(entry.data.rarity),
    toneClass: itemTypeClass(entry.data.type),
    description: entry.data.description,
    symbol: "▣",
  }
}

export function PublicCodexPanel() {
  const [tab, setTab] = useState<Tab>("stones")
  const sitonesQuery = useQuery({
    queryKey: ["catalog", "sitones"],
    queryFn: gameApi.catalogSitones,
  })
  const itemsQuery = useQuery({
    queryKey: ["catalog", "items"],
    queryFn: gameApi.catalogItems,
  })

  const entries: CodexEntry[] =
    tab === "stones"
      ? (sitonesQuery.data ?? []).map((data) => ({ kind: "stones", data }))
      : (itemsQuery.data ?? []).map((data) => ({ kind: "items", data }))
  const isPending =
    tab === "stones" ? sitonesQuery.isPending : itemsQuery.isPending

  return (
    <div className="flex flex-col">
      <Card className="border-ink rounded-[30px] border-2 py-0 shadow-[5px_5px_0_rgba(23,35,58,0.16)]">
        <CardContent className="p-[18px]">
          <span className="text-muted-foreground mb-1 block text-xs font-black tracking-[0.08em] uppercase">
            查詢用途
          </span>
          <h2 className="mb-2 text-[26px] leading-[1.22] font-black tracking-tight">
            查看遊戲中所有小石與道具定義。
          </h2>
          <p className="text-muted-foreground leading-[1.65]">
            這頁不顯示你的持有狀態，只整理名稱、類型、稀有度與用途描述。
          </p>
        </CardContent>
      </Card>

      <nav className="my-3.5 grid grid-cols-2 gap-2" aria-label="圖鑑分頁">
        <Button
          type="button"
          variant={tab === "stones" ? "default" : "outline"}
          className={`min-h-11 rounded-2xl border-2 font-black shadow-none ${
            tab === "stones" ? "border-ink bg-ink text-card" : "border-border"
          }`}
          onClick={() => setTab("stones")}
        >
          小石圖鑑
        </Button>
        <Button
          type="button"
          variant={tab === "items" ? "default" : "outline"}
          className={`min-h-11 rounded-2xl border-2 font-black shadow-none ${
            tab === "items" ? "border-ink bg-ink text-card" : "border-border"
          }`}
          onClick={() => setTab("items")}
        >
          道具圖鑑
        </Button>
      </nav>

      <section
        className="grid grid-cols-2 gap-2.5"
        aria-label={tab === "stones" ? "小石圖鑑" : "道具圖鑑"}
      >
        {isPending ? (
          <article className="bg-card border-ink col-span-2 min-w-0 rounded-[var(--radius)] border-2 p-4">
            <h3 className="text-[17px] leading-tight font-black tracking-tight">
              正在同步圖鑑
            </h3>
          </article>
        ) : entries.length > 0 ? (
          entries.map((entry) => {
            const meta = entryMeta(entry)
            return (
              <article
                key={`${entry.kind}-${entry.data.id}`}
                className="bg-card border-ink min-w-0 rounded-[var(--radius)] border-2 p-3"
              >
                <div
                  className={`${meta.toneClass} border-ink mb-2 flex h-[86px] items-center justify-center rounded-[20px] border-2`}
                >
                  <span
                    className="text-card/90 text-[34px] drop-shadow-[0_2px_0_rgba(23,35,58,0.3)]"
                    aria-hidden
                  >
                    {meta.symbol}
                  </span>
                </div>
                <div className="mt-[9px] mb-1.5 flex flex-wrap gap-1">
                  <span className="border-border bg-surface-raised text-muted-foreground rounded-full border px-[7px] py-[3px] text-[11px] leading-none font-black">
                    {meta.type}
                  </span>
                  <span className="border-border bg-surface-raised text-muted-foreground rounded-full border px-[7px] py-[3px] text-[11px] leading-none font-black">
                    {meta.rarity}
                  </span>
                </div>
                <h3 className="mb-1 text-[17px] leading-tight font-black tracking-tight">
                  {meta.name}
                </h3>
                <p className="text-muted-foreground text-[13px] leading-relaxed">
                  {meta.description}
                </p>
              </article>
            )
          })
        ) : (
          <article className="bg-card border-ink col-span-2 min-w-0 rounded-[var(--radius)] border-2 p-4">
            <h3 className="text-[17px] leading-tight font-black tracking-tight">
              目前沒有資料
            </h3>
          </article>
        )}
      </section>
    </div>
  )
}
