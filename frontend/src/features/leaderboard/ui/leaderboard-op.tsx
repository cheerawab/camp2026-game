import { Card } from "@/shared/ui/card"
import { Separator } from "@/shared/ui/separator"
import { cn } from "@/shared/utils"

const data = [
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/orange/gray/?text=Group+Icon",
    backgroundSrc: "https://placehold.co/300x50/white/gray/?text=Background",
    score: "100",
  },
  {
    teamName: "B",
    pictureSrc: "https://placehold.co/100x100/orange/gray/?text=Group+Icon",
    backgroundSrc: "https://placehold.co/300x50/white/gray/?text=Background",
    score: "80",
  },
  {
    teamName: "C",
    pictureSrc: "https://placehold.co/100x100/orange/gray/?text=Group+Icon",
    backgroundSrc: "https://placehold.co/300x50/white/gray/?text=Background",
    score: "70",
  },
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/orange/gray/?text=Group+Icon",
    backgroundSrc: "https://placehold.co/300x50/white/gray/?text=Background",
    score: "60",
  },
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/orange/gray/?text=Group+Icon",
    backgroundSrc: "https://placehold.co/300x50/white/gray/?text=Background",
    score: "50",
  },
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/orange/gray/?text=Group+Icon",
    backgroundSrc: "https://placehold.co/300x50/white/gray/?text=Background",
    score: "40",
  },
]

const highlight = 2

export function LeaderboardOp() {
  return (
    <div className="grid gap-y-2">
      {data.map((item, index) => {
        const isHighlight: boolean = index + 1 === highlight
        return (
          <Card
            className={cn(
              "flex flex-row items-center gap-x-0 overflow-clip py-0 text-lg",
              isHighlight ? "bg-secondary font-bold" : "",
            )}
          >
            <span className="p-4">#{index + 1}</span>
            <Separator
              orientation="vertical"
              className={isHighlight ? "bg-accent-foreground" : ""}
            />
            <div className="relative flex h-full w-full items-center gap-x-4 p-4">
              {/* z-10 隊伍背景裝飾布局 */}
              <img
                src={item.backgroundSrc}
                className="absolute top-0 left-0 z-10 h-full w-full object-cover"
              />
              {/* z-20 半透明圖層（增加前景文字對比） */}
              <div className="absolute top-0 left-0 h-full w-full z-20 bg-background/50" />
              {/* z-30 隊伍各項資訊 */}
              <img src={item.pictureSrc} className="z-30 h-14 rounded-lg" />
              <div className="z-20 flex flex-1 justify-between">
                <span>{item.teamName} 隊</span>
                <span>{item.score} OP</span>
              </div>
            </div>
          </Card>
        )
      })}
    </div>
  )
}
