import { Badge } from "@/shared/ui/badge"
import { Button } from "@/shared/ui/button"
import { cn } from "@/shared/utils/cn"
import { Info, ShoppingCart, Check } from "lucide-react"
import { toast } from "sonner"

// enum ItemTag {
//   "外觀",
//   "裝飾",
//   "紀念",
// }

type ShopItemCardType = {
  id: string
  name: string
  description: string
  price: number
  tags: Array<string>
  purchased?: boolean
  pictureSrc: string
  className?: string
}

export function ShopItemCard({
  id,
  name,
  description,
  price,
  tags,
  purchased = false,
  pictureSrc,
  className,
}: ShopItemCardType) {
  const onPurchase = () => {
    // TODO: 串接 API
    toast.success(
      <div className="grid">
        <span className="font-bold text-lg">購買成功！</span>
      <span>已成功花費{" "}<strong className="font-bold ">{price} OP</strong>，購買{" "}<strong className="font-bold ">{name}</strong>。</span>
      </div>
      , {
        position: "bottom-center",
        icon: <Check />
    })
  }
  return (
    <div
      className={cn(
        "bg-accent border-foreground rounded-lg border-2 px-8 py-4",
        className,
      )}
      id={`item-${id}`}
    >
      <div className="flex gap-4">
        <div className="basis-1/3">
          <img
            src={pictureSrc}
            alt={name}
            className="aspect-square h-full rounded-lg"
          />
        </div>
        <div className="grid basis-2/3 gap-2">
          <div className="flex justify-between">
            <div className="flex gap-2">
              {tags.map((item) => {
                return <Badge key={item}>{item}</Badge>
              })}
            </div>
            <Badge>{price} OP</Badge>
          </div>
          <div className="text-2xl font-bold">{name}</div>
          <div className="text-muted-foreground">{description}</div>
          <div className="flex gap-2">
            <Button variant="secondary" className="flex-1">
              <Info />
              資訊
            </Button>
            {purchased ? (
              <Button
                variant="outline"
                disabled
                className="flex-1"
              >
                <Check />
                已擁有
              </Button>
            ) : (
              <Button className="flex-1" onClick={onPurchase}>
                <ShoppingCart />
                購買
              </Button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
