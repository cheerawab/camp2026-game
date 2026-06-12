import { useQuery } from "@tanstack/react-query"
import { Link } from "@tanstack/react-router"

import { ShopItemCard } from "@/features/shop-index/ui/shop-item-card"
import { AppError } from "@/shared/api/error"
import { gameApi } from "@/shared/api/game"
import { Badge } from "@/shared/ui/badge"
import { Button } from "@/shared/ui/button"
import { Card, CardContent } from "@/shared/ui/card"
import { GamePageShell } from "@/shared/ui/game-page-shell"
import { PageHeader } from "@/shared/ui/page-header"

export function ShopPage() {
  const statusQuery = useQuery({
    queryKey: ["me", "status"],
    queryFn: gameApi.status,
  })
  const itemsQuery = useQuery({
    queryKey: ["shop", "items"],
    queryFn: gameApi.shopItems,
  })
  const isUnauthorized =
    (statusQuery.error instanceof AppError &&
      statusQuery.error.status === 401) ||
    (itemsQuery.error instanceof AppError && itemsQuery.error.status === 401)

  return (
    <GamePageShell contentClassName="grid content-start gap-y-2">
      <PageHeader
        title="商店"
        headline="Item Shop"
        rightSlot={
          <div className="flex flex-col items-end">
            <span className="text-muted-foreground text-sm font-bold">
              你現在持有開源力
            </span>
            <Badge className="h-fit">
              開源力 {statusQuery.data?.openPower ?? 0}
            </Badge>
          </div>
        }
      />

      {isUnauthorized ? (
        <Card>
          <CardContent className="grid gap-3">
            <h2 className="text-2xl font-bold">請先登入</h2>
            <p className="text-muted-foreground">
              登入後才能查看可兌換商品與目前開源力。
            </p>
            <Button asChild>
              <Link to="/login">前往登入</Link>
            </Button>
          </CardContent>
        </Card>
      ) : itemsQuery.isPending ? (
        <Card>
          <CardContent>
            <span className="text-muted-foreground font-bold">
              正在同步商品
            </span>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-y-2">
          {(itemsQuery.data ?? []).map((item) => (
            <ShopItemCard
              key={item.id}
              item={item}
              currentOpenPower={statusQuery.data?.openPower}
            />
          ))}
        </div>
      )}
    </GamePageShell>
  )
}
