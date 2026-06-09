import { BattleIngameTeam } from "@/features/battle-ingame/ui/battle-ingame-team"
import { Button } from "@/shared/ui/button"
import { Card, CardContent } from "@/shared/ui/card"
import { Separator } from "@/shared/ui/separator"

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

const teamDataA = [
  {
    name: "A 小石",
    type: "靈光",
    pictureSrc: "https://placehold.co/100x100/svg?text=A",
  },
  {
    name: "B 小石",
    type: "靈光",
    pictureSrc: "https://placehold.co/100x100/svg?text=B",
  },
  {
    name: "C 小石",
    type: "靈光",
    pictureSrc: "https://placehold.co/100x100/svg?text=C",
  },
  {
    name: "D 小石",
    type: "靈光",
    pictureSrc: "https://placehold.co/100x100/svg?text=D",
  },
  {
    name: "E 小石",
    type: "靈光",
    pictureSrc: "https://placehold.co/100x100/svg?text=E",
  },
]

export function BattleQuestionPage() {
  return (
    <>
      <main className="mx-auto grid w-full max-w-sm gap-y-2 px-2 py-4">
        {/* 分數 - 時間 - 分數 */}
        <div className="grid grid-cols-3 pb-2 text-center">
          <div>
            <div className="text-lg font-bold">125</div>
            <div className="text-muted-foreground text-xs">混凝土</div>
          </div>
          <div className="flex items-end justify-center">
            <div className="text-4xl font-bold">18</div>
            <div className="text-lg">s</div>
          </div>
          <div>
            <div className="text-lg font-bold">256</div>
            <div className="text-muted-foreground text-xs">義大利麵</div>
          </div>
        </div>
        {/*<BattleIngameScoreBar a={125} b={256} />*/}
        <Card>
          <CardContent className="grid">
            {/* 題目 */}
            <div className="grid gap-y-0 pb-2">
              <span className="text-muted-foreground text-sm font-bold uppercase">
                {quizData.type}
              </span>
              <span className="text-2xl font-bold">{quizData.question}</span>
            </div>
            {/**/}
            <Separator />
            {/* 選項 */}
            {quizData.answers.map((item, index) => {
              return (
                <>
                  <Button
                    variant="ghost"
                    className="h-fit justify-start rounded-none py-2 pl-0"
                  >
                    <span className="border-accent-foreground bg-accent text-muted-foreground rounded-lg border-2 px-4 py-2">
                      {["A", "B", "C", "D"][index]}
                    </span>
                    <span className="text-lg">{item}</span>
                  </Button>
                  {index < 3 && <Separator />}
                </>
              )
            })}
          </CardContent>
        </Card>
        {/* 小石資訊 */}
        <Card>
          <CardContent className="flex">
            <BattleIngameTeam team={teamDataA} highlight={2} />
            <Separator orientation="vertical" />
            <BattleIngameTeam team={teamDataA} highlight={2} reverse />
          </CardContent>
        </Card>
      </main>
    </>
  )
}
