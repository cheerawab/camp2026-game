import { Card, CardContent } from "@/shared/ui/card"

const data = [
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/svg",
    score: "A?",
  },
  {
    teamName: "B",
    pictureSrc: "https://placehold.co/100x100/svg",
    score: "B?",
  },
  {
    teamName: "C",
    pictureSrc: "https://placehold.co/100x100/svg",
    score: "C?",
  },
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/svg",
    score: "A?",
  },
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/svg",
    score: "A?",
  },
  {
    teamName: "A",
    pictureSrc: "https://placehold.co/100x100/svg",
    score: "A?",
  },
]

const highlight = 2

export function LeaderboardBattle() {
  return (
    <div className="grid gap-y-2">
      {data.map((item, index) => {
        return (
          <Card
            className={index + 1 === highlight ? "bg-secondary font-bold" : ""}
          >
            <CardContent className="flex items-center gap-x-4">
              <span>#{index + 1}</span>
              <img src={item.pictureSrc} className="h-12" />
              <div className="flex flex-1 justify-between">
                <span>{item.teamName} 隊</span>
                <span>{item.score}</span>
              </div>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
