import { Badge } from "@/shared/ui/badge"
import { Button } from "@/shared/ui/button"
import { Card, CardContent } from "@/shared/ui/card"
import { useNavigate } from "@tanstack/react-router"
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
}

export function ShopItemCard({
  id,
  name,
  description,
  price,
  tags,
  purchased = false,
  pictureSrc,
}: ShopItemCardType) {
  const navigate = useNavigate()

  const onInspect = () => {
    navigate({ to: `/shop/${id}` })
  }

  const onPurchase = () => {
    // TODO: 串接 API
    toast.success(
      <div className="grid">
        <span className="text-lg font-bold">購買成功！</span>
        <span>
          已成功花費 <strong className="font-bold">{price} OP</strong>，購買{" "}
          <strong className="font-bold">{name}</strong>。
        </span>
      </div>,
      {
        position: "bottom-center",
        icon: <Check />,
      },
    )
  }
  return (
    <Card>
      <CardContent>
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
                  return (
                    <Badge variant="outline" key={item}>
                      {item}
                    </Badge>
                  )
                })}
              </div>
              <Badge>{price} OP</Badge>
            </div>
            <div className="text-2xl font-bold">{name}</div>
            <div className="text-muted-foreground">{description}</div>
            <div className="grid grid-cols-2 gap-2">
              <Button variant="secondary" onClick={onInspect}>
                <Info />
                資訊
              </Button>
              {purchased ? (
                <Button variant="outline" disabled>
                  <Check />
                  已擁有
                </Button>
              ) : (
                <Button onClick={onPurchase}>
                  <ShoppingCart />
                  購買
                </Button>
              )}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
