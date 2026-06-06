import { BattleWaitingBar } from "@/features/battle-waiting/ui/battle-waiting-bar"
import { Button } from "@/shared/ui/button"
import { Card, CardContent } from "@/shared/ui/card"
import { Progress } from "@/shared/ui/progress"
import { StatusDot } from "@/shared/ui/status-dot"
import { Clock, Hourglass, Timer, Watch } from "lucide-react"

const quizData = {
  progress: {
    current: 3,
    max: 5,
  },
  time: {
    current: 19,
    max: 30,
  },
  question: "下列哪一項最能描述開源授權在專案中的作用？",
  type: "Open Source",
  answers: [
    "開源專案的條款",
    "營隊隊呼的節奏",
    "午餐菜單的排序",
    "場地網路的密碼",
  ],
  state: [true, false],
}

export function BattleIngamePage() {
  return (
    <>
      <main className="mx-auto grid w-full max-w-sm gap-2 py-4">
        {/* 題目 */}
        <Card>
          <CardContent className="grid gap-y-2">
            <span className="text-muted-foreground text-sm font-bold uppercase">
              {quizData.type}
            </span>
            <span className="text-2xl font-bold">{quizData.question}</span>
            <div className="flex items-center gap-2">
              <Hourglass className="spin" />
              <BattleWaitingBar
                value={quizData.time.current}
                max={quizData.time.max}
              />
            </div>
          </CardContent>
        </Card>
        {/* 選項 */}
        {quizData.answers.map((item, index) => {
          return (
            <Button variant="outline" className="h-fit justify-start">
              <span className="border-muted-foreground bg-muted text-muted-foreground rounded-lg border-2 px-4 py-2 text-lg">
                {["A", "B", "C", "D"][index]}
              </span>
              <span className="text-lg">{item}</span>
            </Button>
          )
        })}
        {/* 作答情形 */}
        <div className="grid grid-cols-2 gap-x-2">
          <Card>
            <CardContent>
              <div className="flex items-center gap-x-2">
                <StatusDot tone="warning" />
                <span className="text-lg font-bold">你</span>
              </div>
              <span>尚未作答</span>
            </CardContent>
          </Card>
          <Card>
            <CardContent>
              <div className="flex items-center gap-x-2">
                <StatusDot tone="success" />
                <span className="text-lg font-bold">對手</span>
              </div>
              <span>已作答</span>
            </CardContent>
          </Card>
        </div>
      </main>

      <style>
        {`
        @keyframes hourGlassSpin {
          0% {
            transform: rotate(0deg);
          }
          40% {
          transform: rotate(200deg);
          }
          60% {
            transform: rotate(170deg);
          }
          70%{
            transform: rotate(180deg);
          }
          100% {
          transform: rotate(180deg);
          }
        }

        .spin {
        animation: hourGlassSpin 1s infinite linear;
        }
      `}
      </style>
    </>
  )
}
