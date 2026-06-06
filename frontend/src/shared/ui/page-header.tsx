import { Button } from "./button"
import { ChevronLeft } from "lucide-react"

type PageHeaderType = {
  title: string
  headline: string
}

export function PageHeader({ title, headline }: PageHeaderType) {
  const goPrevPage = () => {
    window.history.back()
  }
  return (
    <div className="flex items-start gap-x-4 py-2">
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
  )
}
