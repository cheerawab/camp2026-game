import { useMemo, useState } from "react"

import { Button } from "@/shared/ui/button"
import { Card } from "@/shared/ui/card"

type PickItem = {
  id: string
  name: string
  type?: string
  count: number
  bgClassName: string
}

const stones: PickItem[] = [
  {
    id: "solder-seed",
    name: "焊點種子石",
    type: "工程",
    count: 2,
    bgClassName: "bg-pebble-engineer",
  },
  {
    id: "base-route",
    name: "基地路線石",
    type: "探索",
    count: 4,
    bgClassName: "bg-pebble-explore",
  },
  {
    id: "camp-lamp",
    name: "營燈靈光",
    type: "靈光",
    count: 1,
    bgClassName: "bg-pebble-spark",
  },
]

const materials: PickItem[] = [
  {
    id: "camp-rivet",
    name: "營釘鉚扣",
    count: 12,
    bgClassName: "bg-pebble-engineer",
  },
  {
    id: "map-cotton",
    name: "地圖棉線",
    count: 8,
    bgClassName: "bg-pebble-explore",
  },
]

// TODO: 串接 API
export function StoneFusionPanel() {
  const [selectedStoneId, setSelectedStoneId] = useState(stones[0]?.id)
  const [selectedMaterialId, setSelectedMaterialId] = useState(materials[1]?.id)
  const [status, setStatus] = useState("選擇一顆小石與一個素材，預覽可能產物。")

  const selectedStone = useMemo(
    () => stones.find((stone) => stone.id === selectedStoneId),
    [selectedStoneId],
  )
  const selectedMaterial = useMemo(
    () => materials.find((material) => material.id === selectedMaterialId),
    [selectedMaterialId],
  )

  return (
    <div className="flex flex-col pb-2">
      <Card
        className="bg-surface-raised fixed top-[84px] left-1/2 z-10 grid w-[calc(min(100vw,430px)-2rem)] -translate-x-1/2 grid-cols-[86px_1fr] items-center gap-3 rounded-[22px] px-3.5 py-3"
        aria-label="合成結果預覽"
      >
        <div
          className={[
            "border-ink grid size-[86px] rotate-[-6deg] place-items-center rounded-[26px_32px_22px_28px] border-2",
            selectedStone?.bgClassName ?? "bg-pebble-engineer",
          ].join(" ")}
          aria-hidden
        >
          <span
            className={[
              "border-card size-11 rounded-[16px] border-[3px]",
              selectedMaterial?.bgClassName ?? "bg-pebble-explore",
            ].join(" ")}
          />
        </div>
        <div className="min-w-0">
          <p className="text-muted-foreground mb-1 text-xs font-bold tracking-[0.08em] uppercase">
            預覽結果
          </p>
          <h2 className="text-2xl leading-tight font-extrabold tracking-normal">
            工程路線展示框
          </h2>
          <p className="text-muted-foreground mt-1 text-sm leading-5 font-medium">
            把工程小石與地圖素材組合成基地展示用的收藏外框。
          </p>
        </div>
      </Card>

      <div className="mt-[124px] space-y-2 pb-24">
        <PickerSection
          title="選小石"
          description="不會消耗所有同類，只消耗本次選定數量。"
        >
          {stones.map((stone, index) => (
            <PickCard
              key={stone.id}
              item={stone}
              selected={selectedStoneId === stone.id}
              onSelect={() => setSelectedStoneId(stone.id)}
              autoFocus={index === 0}
            />
          ))}
        </PickerSection>

        <PickerSection title="選素材" description="數量不足會無法合成。">
          {materials.map((material) => (
            <PickCard
              key={material.id}
              item={material}
              selected={selectedMaterialId === material.id}
              onSelect={() => setSelectedMaterialId(material.id)}
            />
          ))}
        </PickerSection>
      </div>

      <Card className="fixed bottom-4 left-1/2 z-20 mt-6 grid w-[calc(min(100vw,430px)-2rem)] -translate-x-1/2 grid-cols-[1fr_116px] items-center gap-2 rounded-[20px] px-3 py-2">
        <div className="min-w-0">
          <span className="text-muted-foreground block text-xs font-bold">
            消耗
          </span>
          <strong className="block truncate text-sm font-bold">
            1 小石 + 3 素材
          </strong>
          <span
            className="text-muted-foreground mt-0.5 block truncate text-[11px] leading-4 font-medium"
            aria-live="polite"
          >
            {status}
          </span>
        </div>
        <Button
          type="button"
          className="min-h-9 rounded-2xl px-3 text-sm font-bold"
          onClick={() => setStatus("已合成：工程路線展示框。")}
        >
          確認合成
        </Button>
      </Card>
    </div>
  )
}

function PickerSection({
  title,
  description,
  children,
}: {
  title: string
  description: string
  children: React.ReactNode
}) {
  return (
    <Card className="rounded-[22px] px-3 py-1.5">
      <div className="mb-0 flex items-start justify-between gap-3">
        <h2 className="pt-5 text-[20px] leading-none font-extrabold">
          {title}
        </h2>
        <small className="text-muted-foreground max-w-[150px] pt-5 text-right text-[11px] leading-4 font-medium">
          {description}
        </small>
      </div>
      <div className="grid gap-2">{children}</div>
    </Card>
  )
}

function PickCard({
  item,
  selected,
  onSelect,
  autoFocus,
}: {
  item: PickItem
  selected: boolean
  onSelect: () => void
  autoFocus?: boolean
}) {
  return (
    <Button
      type="button"
      variant={selected ? "secondary" : "outline"}
      className={[
        "bg-surface-raised grid h-auto min-h-[60px] grid-cols-[40px_1fr] justify-start gap-2 rounded-[18px] p-2 text-left shadow-none",
        selected ? "border-ink" : "border-border",
      ].join(" ")}
      onClick={onSelect}
      aria-pressed={selected}
      autoFocus={autoFocus}
    >
      <span
        className={[
          "border-ink size-10 rounded-[14px_18px_12px_16px] border-2",
          item.bgClassName,
        ].join(" ")}
        aria-hidden
      />
      <span className="min-w-0">
        <strong className="block truncate font-bold">{item.name}</strong>
        <small className="text-muted-foreground block text-xs font-semibold">
          {item.type || "素材"} · ×{item.count}
        </small>
      </span>
    </Button>
  )
}
