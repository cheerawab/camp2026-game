import { ShopPage } from "@/pages/shop/ui/shop-page"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/shop")({
  component: ShopPage,
})
