import { ShopInfoPage } from "@/pages/shop/info/ui/shop-info-page"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/shop/$id")({
  component: ShopInfoRoute,
})

function ShopInfoRoute() {
  const { id } = Route.useParams()
  return <ShopInfoPage data={id} />
}
