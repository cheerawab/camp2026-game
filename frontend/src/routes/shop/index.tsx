import { ShopPage } from "@/pages/shop/index/ui/shop-index-page"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/shop/")({
  component: ShopPage,
})
