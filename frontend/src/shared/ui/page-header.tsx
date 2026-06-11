import { Link } from "@tanstack/react-router"
import { ArrowLeft } from "lucide-react"
import type { ReactNode } from "react"

type PageHeaderType = {
  title: string
  headline: string
  backTo?: string
  rightSlot?: ReactNode
}

export function PageHeader({
  title,
  headline,
  backTo = "/",
  rightSlot,
}: PageHeaderType) {
  return (
    <div className="flex items-start gap-x-4 py-2">
      <Link
        to={backTo}
        aria-label="返回"
        className="border-ink bg-card text-ink focus-visible:outline-power grid size-11 shrink-0 place-items-center rounded-2xl border-2 transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
      >
        <ArrowLeft className="size-5" aria-hidden />
      </Link>
      <div className="flex-1">
        <p className="text-muted-foreground text-sm font-bold uppercase">
          {headline}
        </p>
        <h1 className="text-2xl font-bold">{title}</h1>
      </div>
      {rightSlot}
    </div>
  )
}
