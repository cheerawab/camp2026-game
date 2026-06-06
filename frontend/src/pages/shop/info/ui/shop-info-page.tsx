import { Badge } from "@/shared/ui/badge"
import { Button } from "@/shared/ui/button"
import { X, Check } from "lucide-react"

const itemData = {
  id: "lantern-theme",
  name: "營燈基地佈景",
  description: "讓基地展示區換成暖黃營燈與旗繩。",
  tags: ["外觀"],
  price: 260,
  state: "可購買",
  purchased: true,
  effect: "兌換後可在基地展示區選用工程工作台風格。",
  pictureSrc: "https://placehold.co/100x100/svg",
}

type ShopInfoPageType = {
  data: string
}

export function ShopInfoPage({ data }: ShopInfoPageType) {
  console.log(data)

  // TODO: 串接 API

  return (
    <main className="mx-auto grid w-full max-w-sm gap-y-4 py-4">
      {/* TODO: place share component pageheader */}
      {/* 物品圖片 */}
      <div className="border-foreground bg-accent relative w-full rounded-lg border-2 py-8">
        <img
          src={itemData.pictureSrc}
          className="border-foreground mx-auto size-36 -rotate-2 rounded-lg border-2"
        />
        <div className="border-accent-foreground bg-muted absolute right-4 bottom-4 grid rotate-3 gap-0 rounded-lg border-2 px-4 py-2">
          <span className="text-muted-foreground text-sm leading-none">
            價格
          </span>
          <span className="text-lg leading-none font-bold">
            {itemData.price} OP
          </span>
        </div>
      </div>
      {/* 詳細資訊 */}
      <div className="border-foreground bg-accent relative grid w-full gap-y-4 rounded-lg border-2 px-8 py-4">
        {/* 標籤 */}
        <div className="flex">
          {itemData.tags.map((item) => {
            return <Badge variant="secondary">{item}</Badge>
          })}
        </div>
        {/* 名稱 */}
        <span className="text-2xl font-bold">{itemData.name}</span>
        {/* 簡介 */}
        <div className="grid gap-1">
          <span className="text-muted-foreground text-sm font-bold">
            道具簡介
          </span>
          <span>{itemData.description}</span>
        </div>
        {/* 效果 */}
        <div className="grid gap-1">
          <span className="text-muted-foreground text-sm font-bold">
            套用效果
          </span>
          <span className="">{itemData.effect}</span>
        </div>
      </div>
      {/* 動作區 */}
      <div className="border-foreground bg-accent relative flex w-full items-center justify-between gap-y-4 rounded-lg border-2 px-8 py-4">
        <div className="grid">
          <span className="text-muted-foreground text-sm font-bold">
            現在開源力
          </span>
          <span className="text-2xl font-bold">430 OP</span>
        </div>
        <Button size="lg">
          <Check />
          確認購買
        </Button>
      </div>
    </main>
  )
}
