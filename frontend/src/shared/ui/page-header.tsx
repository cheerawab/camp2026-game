import { ReactNode } from "react"
import { Button } from "./button"
import { ChevronLeft } from "lucide-react"

type PageHeaderType = {
  title: string
  headline: string
  rightSlot?: ReactNode
}

export function PageHeader({ title, headline, rightSlot }: PageHeaderType) {
  const goPrevPage = () => {
    window.history.back()
  }
  return (
    <div className="flex items-start justify-between py-2">
      <div className="flex items-start gap-x-4">
        <Button size="icon" variant="secondary" onClick={goPrevPage}>
          <ChevronLeft />
        </Button>
        <div>
          <p className="text-muted-foreground text-sm font-bold uppercase">
            {headline}
          </p>
          <h1 className="text-2xl font-bold">{title}</h1>
        </div>
      </div>
      {rightSlot}
    </div>
  )
}
