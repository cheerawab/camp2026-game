import { Button } from "@/shared/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card"
import { PageHeader } from "@/shared/ui/page-header"
import { Play, Plus } from "lucide-react"

export function BattleIndexPage() {
  return (
    <main className="mx-auto grid w-full max-w-sm gap-y-2 py-4">
      <PageHeader title="知識王" headline="Quiz Battle Lobby" />
      {/* 單人遊戲 */}
      <Card>
        <CardHeader>
          <CardTitle>快速開始</CardTitle>
          <CardDescription>與 CPU 對決</CardDescription>
        </CardHeader>
        <CardContent>
          <span>與 SITCON 電腦進行對決，複習上課知識！</span>
        </CardContent>
        <CardFooter>
          <Button className="w-full">開始遊戲</Button>
        </CardFooter>
      </Card>
      {/* 多人連線 */}
      <Card>
        <CardHeader>
          <CardTitle></CardTitle>
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent>
          <Input
        </CardContent>
        <CardFooter className="flex gap-2">
          <Button className="w-full flex-1" variant="secondary">加入房間</Button>
          <Button className="w-full flex-1">創建房間</Button>
        </CardFooter>
      </Card>
    </main>
  )
}
