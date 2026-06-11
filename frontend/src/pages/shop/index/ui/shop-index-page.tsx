import { ShopItemCard } from "@/features/shop-index/ui/shop-item-card"
import { Badge } from "@/shared/ui/badge"
import { PageHeader } from "@/shared/ui/page-header"
import { Toaster } from "@/shared/ui/sonner"
import { DollarSign } from "lucide-react"

const SHOP_ITEMS = [
  {
    id: "lantern-theme",
    name: "營燈基地佈景",
    description: "讓基地展示區換成暖黃營燈與旗繩。",
    tags: ["外觀"],
    price: 260,
    state: "可購買",
    purchased: true,
    pictureSrc: "https://placehold.co/100x100/svg",
  },
  {
    id: "map-frame",
    name: "探索地圖邊框",
    tags: ["收藏"],
    description: "套用在小石卡上的路線圖外框。",
    price: 180,
    pictureSrc: "https://placehold.co/100x100/svg",
  },
  {
    id: "radio-badge",
    name: "小隊電波徽章",
    tags: ["紀念"],
    price: 520,
    description: "對戰完成後可展示在個人 QR Code 頁。",
    pictureSrc: "https://placehold.co/100x100/svg",
  },
  {
    id: "workbench-theme",
    name: "工程工作台佈景",
    tags: ["外觀"],
    description: "替基地切換成模組板、焊點與工具貼紙風格。",
    price: 360,
    pictureSrc: "https://placehold.co/100x100/svg",
  },
  {
    id: "stage-ribbon",
    name: "舞台彩帶包",
    tags: ["收藏"],
    description: "增加娛樂系小石展示用的彩帶標籤。",
    price: 120,
    pictureSrc: "https://placehold.co/100x100/svg",
  },
]

export function ShopPage() {
  return (
    <>
      <main className="mx-auto grid w-full max-w-sm gap-y-2 py-4">
        {/* 標題 & 返回 */}
        <PageHeader
          title="商店"
          headline="Item Shop"
          rightSlot={
            <div className="flex flex-col items-end">
              <span className="text-muted-foreground text-sm font-bold">
                你現在持有開源力
              </span>
              <Badge className="h-fit">
                <DollarSign className="h-4 w-4" />
                1000 OP
              </Badge>
            </div>
          }
        />
        {/* 商品列表 */}
        <div className="grid gap-y-2">
          {SHOP_ITEMS.map((item) => {
            return (
              <ShopItemCard
                id={item.id}
                name={item.name}
                description={item.description}
                price={item.price}
                tags={item.tags}
                purchased={item?.purchased}
                pictureSrc={item.pictureSrc}
              />
            )
          })}
        </div>
      </main>
      <Toaster />
    </>
  )
}
