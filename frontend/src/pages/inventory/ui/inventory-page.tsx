import { useState } from "react"

import { Button } from "@/shared/ui/button"
import { Card } from "@/shared/ui/card"
import { PageHeader } from "@/shared/ui/page-header"

const BAG_ITEMS = [
  {
    id: "map-thread",
    name: "地圖棉線",
    type: "素材",
    rarity: "常見",
    count: 8,
    colorClass: "bg-pebble-explore",
    description: "用來標記基地佈景上的探索路線。",
  },
  {
    id: "camp-rivet",
    name: "營釘鉚扣",
    type: "素材",
    rarity: "常見",
    count: 12,
    colorClass: "bg-pebble-engineer",
    description: "可替收藏卡加上工程感邊框。",
  },
  {
    id: "lantern-ticket",
    name: "營燈佈景券",
    type: "外觀",
    rarity: "稀有",
    count: 1,
    colorClass: "bg-pebble-spark",
    description: "替小隊基地換上暖黃營燈主題。",
  },
  {
    id: "radio-pin",
    name: "小隊電波徽章",
    type: "活動紀念",
    rarity: "稀有",
    count: 2,
    colorClass: "bg-pebble-resonate",
    description: "和隊友完成同步挑戰後取得。",
  },
  {
    id: "stage-strip",
    name: "舞台彩帶",
    type: "外觀",
    rarity: "常見",
    count: 3,
    colorClass: "bg-pebble-play",
    description: "能讓展示櫃多一條娛樂系標記。",
  },
  {
    id: "decode-card",
    name: "解題提示卡",
    type: "活動紀念",
    rarity: "傳說",
    count: 1,
    colorClass: "bg-primary",
    description: "完成一場高難度知識王戰後留下的紀念卡。",
  },
]

const FILTERS = ["全部", "素材", "外觀", "活動紀念"]

type BagItem = (typeof BAG_ITEMS)[number]

function ItemIcon({
  colorClass,
  count,
}: {
  colorClass: string
  count: number
}) {
  return (
    <div
      className={[
        "border-ink relative h-[72px] overflow-hidden rounded-[20px] border-2",
        colorClass,
      ].join(" ")}
      aria-hidden
    >
      <span className="border-card/80 absolute inset-[14px_18px] -rotate-[12deg] border-x-0 border-y-[3px]" />
      <strong className="bg-card border-ink absolute right-1.5 bottom-[5px] min-w-[28px] rounded-full border-2 px-1.5 py-px text-center text-[13px]">
        {count}
      </strong>
    </div>
  )
}

function ItemCard({ item }: { item: BagItem }) {
  return (
    <Card className="border-ink grid grid-cols-[72px_1fr] items-start gap-3 rounded-[22px] border-2 p-3.5">
      <ItemIcon colorClass={item.colorClass} count={item.count} />
      <div>
        <div className="mb-1.5 flex items-start justify-between gap-2">
          <h3 className="text-[18px] leading-tight font-black tracking-[-0.02em]">
            {item.name}
          </h3>
          <strong className="text-primary shrink-0 text-[18px] font-black">
            ×{item.count}
          </strong>
        </div>
        <div className="mb-1.5 flex flex-wrap gap-1.5">
          {[item.type, item.rarity].map((tag) => (
            <span
              key={tag}
              className="bg-surface-raised border-border text-muted-foreground rounded-full border-[1.5px] px-2 py-0.5 text-xs font-black"
            >
              {tag}
            </span>
          ))}
        </div>
        <p className="text-muted-foreground text-sm leading-[1.62]">
          {item.description}
        </p>
      </div>
    </Card>
  )
}

function EmptyBag() {
  return (
    <Card
      className="border-ink rounded-[22px] border-2 border-dashed p-[22px]"
      aria-label="空背包狀態"
    >
      <div className="bg-surface-raised border-ink mb-3 grid size-[54px] place-items-center rounded-[18px] border-2 text-[28px] font-black">
        ＋
      </div>
      <h3 className="mb-1 text-[17px] font-black">這個分類目前沒有道具</h3>
      <p className="text-muted-foreground mb-3.5 leading-[1.65]">
        可以到商店兌換，或在現場活動完成挑戰後取得。
      </p>
      <Button type="button" className="w-full">
        前往商店
      </Button>
    </Card>
  )
}

export function InventoryPage() {
  // TODO: replace with GET /api/inventory
  const [filter, setFilter] = useState("全部")
  const visibleItems =
    filter === "全部"
      ? BAG_ITEMS
      : BAG_ITEMS.filter((item) => item.type === filter)
  const totalCount = BAG_ITEMS.reduce((sum, item) => sum + item.count, 0)

  return (
    <main
      className="bg-paper text-ink mx-auto grid w-full max-w-[430px] gap-y-3 px-4 py-[18px] pb-7"
      aria-label="道具背包頁"
    >
      <PageHeader
        title="道具背包"
        headline="Field Bag"
        rightSlot={
          <Button variant="secondary" size="sm" className="border-ink border-2">
            整理
          </Button>
        }
      />

      {/* 背包摘要 */}
      <Card
        className="border-ink flex justify-between gap-4 rounded-[22px] border-2 p-5"
        style={{ boxShadow: "5px 5px 0 rgba(23,35,58,.16)" }}
        aria-label="背包摘要"
      >
        <div>
          <p className="text-muted-foreground mb-[3px] text-xs font-black tracking-[0.08em] uppercase">
            目前持有
          </p>
          <strong className="block text-[48px] leading-[0.95] font-black tracking-[-0.05em]">
            {totalCount}
          </strong>
          <p className="text-muted-foreground mt-2 leading-[1.65]">
            素材、外觀與活動紀念都會先收在這裡。
          </p>
        </div>
        <div
          className="bg-surface-raised border-ink relative min-w-[92px] rounded-[22px_22px_16px_16px] border-2"
          style={{ height: 112 }}
          aria-hidden
        >
          <span className="border-ink absolute inset-[16px_24px] h-6 rounded-t-2xl border-2 border-b-0" />
          <span className="border-primary border-b-moss absolute right-5 bottom-[22px] left-5 h-[18px] border-t-[3px] border-b-[3px]" />
        </div>
      </Card>

      {/* 分類統計 */}
      <section className="grid grid-cols-3 gap-[10px]" aria-label="分類數量">
        {[
          { label: "素材", value: 20 },
          { label: "外觀", value: 4 },
          { label: "紀念", value: 3 },
        ].map(({ label, value }) => (
          <div
            key={label}
            className="bg-surface-raised border-border rounded-[18px] border-2 p-3"
          >
            <span className="text-muted-foreground block text-xs font-black">
              {label}
            </span>
            <strong className="mt-0.5 block text-[22px] font-black">
              {value}
            </strong>
          </div>
        ))}
      </section>

      {/* 篩選器 */}
      <nav className="flex gap-2 overflow-x-auto pb-1" aria-label="背包分類">
        {FILTERS.map((f) => (
          <Button
            key={f}
            variant={filter === f ? "default" : "outline"}
            size="sm"
            className="shrink-0 rounded-2xl font-black"
            onClick={() => setFilter(f)}
          >
            {f}
          </Button>
        ))}
      </nav>

      {/* 道具列表 */}
      <section className="grid gap-3" aria-label="道具列表">
        {visibleItems.length > 0 ? (
          visibleItems.map((item) => <ItemCard key={item.id} item={item} />)
        ) : (
          <EmptyBag />
        )}
      </section>
    </main>
  )
}
