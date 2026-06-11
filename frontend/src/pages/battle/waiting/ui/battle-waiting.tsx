import { BattleWaitingPlayerCard } from "@/features/battle-waiting/ui/battle-waiting-player-card"
import { Button } from "@/shared/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card"
import { PageHeader } from "@/shared/ui/page-header"
import { Separator } from "@/shared/ui/separator"
import { Check, DoorOpen } from "lucide-react"

const roomCode: string = "K7M2"
const qrCode = "https://placehold.co/100x100/svg"

const playerList = [
  {
    name: "混凝土",
    team: "甲隊",
    ready: true,
    pictureSrc: "https://placehold.co/100x100/svg",
  },
  {
    name: "義大利麵",
    team: "乙隊",
    ready: false,
    pictureSrc: "https://placehold.co/100x100/svg",
  },
]

export function BattleWaitingPage() {
  return (
    <main className="mx-auto grid w-full max-w-sm gap-y-2 py-4">
      <PageHeader title="等待房間" headline="Battle Room" />
      {/* 邀請碼 */}
      <Card>
        <CardHeader>
          <CardTitle>房號</CardTitle>
          <CardDescription>
            將這個房號分享給其他學員，或是透過 QRCode 並加入對戰吧！
          </CardDescription>
        </CardHeader>
        <CardContent className="flex items-center justify-center gap-4">
          <span className="block text-4xl font-bold tracking-[1rem]">
            {roomCode}
          </span>
          <Separator orientation="vertical" />
          <img src={qrCode} />
        </CardContent>
      </Card>
      {/* 房間玩家列表 */}
      {playerList.map((item) => {
        return (
          <BattleWaitingPlayerCard
            name={item.name}
            team={item.team}
            ready={item.ready}
            pictureSrc={item.pictureSrc}
          />
        )
      })}
      {/* 動作區 */}
      <Card>
        <CardContent className="grid grid-cols-2 gap-2">
          <Button variant="outline" size="lg">
            <DoorOpen />
            離開房間
          </Button>
          <Button size="lg">
            <Check />
            準備完成
          </Button>
        </CardContent>
      </Card>
    </main>
  )
}
