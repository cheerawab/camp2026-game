import { createFileRoute } from "@tanstack/react-router"

import { ShopItemDetailPage } from "@/pages/shop/ui/shop-item-detail-page"

export const Route = createFileRoute("/shop/$itemId")({
  component: ShopItemDetailPage,
})
