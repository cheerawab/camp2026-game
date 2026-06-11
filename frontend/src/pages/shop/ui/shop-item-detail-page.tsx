import { Badge } from "@/shared/ui/badge"
import { Button } from "@/shared/ui/button"
import { Check } from "lucide-react"
import { PageHeader } from "@/shared/ui/page-header"
import { Card, CardContent } from "@/shared/ui/card"

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

export function ShopItemDetailPage() {
  // TODO: replace with GET /api/shop/:id

  return (
    <main className="mx-auto grid w-full max-w-sm gap-y-2 px-4 py-4">
      {/* 標題 & 返回 */}
      <PageHeader title="道具資訊" headline="Item Detail" backTo="/shop" />
      {/* 物品圖片 */}
      <img
        src={itemData.pictureSrc}
        className="border-foreground mx-auto mb-2 size-36 -rotate-3 rounded-lg border-2"
      />
      {/* 詳細資訊 */}
      <Card>
        <CardContent className="grid gap-y-2">
          {/* 標籤 */}
          <div className="flex">
            {itemData.tags.map((item) => {
              return <Badge variant="outline">{item}</Badge>
            })}
          </div>
          {/* 名稱 */}
          <span className="text-2xl font-bold">{itemData.name}</span>
          {/* 簡介 */}
          <div className="grid gap-y-1">
            <span className="text-accent-foreground text-lg font-bold">
              道具簡介
            </span>
            <span>{itemData.description}</span>
          </div>
          {/* 效果 */}
          <div className="grid gap-1">
            <span className="text-accent-foreground text-lg font-bold">
              套用效果
            </span>
            <span className="">{itemData.effect}</span>
          </div>
        </CardContent>
      </Card>
      {/* 動作區 */}
      <Card>
        <CardContent className="flex items-center justify-between">
          <div className="grid">
            <span className="text-muted-foreground text-sm font-bold">
              花費開源力
            </span>
            <div className="flex items-center">
              <span className="text-2xl font-bold">{itemData.price}</span>
              <span className="whitespace-pre"> / 430 OP</span>
            </div>
          </div>
          <Button size="lg">
            <Check />
            確認購買
          </Button>
        </CardContent>
      </Card>
    </main>
  )
}
