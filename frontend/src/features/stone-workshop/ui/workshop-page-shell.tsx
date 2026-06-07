import { Link } from "@tanstack/react-router"
import { ArrowLeft } from "lucide-react"
import type { ReactNode } from "react"

import { Button } from "@/shared/ui/button"

type WorkshopPageShellProps = {
  title: string
  eyebrow: string
  description: string
  children: ReactNode
}

export function WorkshopPageShell({
  title,
  eyebrow,
  description,
  children,
}: WorkshopPageShellProps) {
  return (
    <main className="bg-paper text-ink min-h-svh">
      <div className="mx-auto flex min-h-svh w-full max-w-[430px] flex-col px-4 py-[18px]">
        <header className="mb-3 flex items-center gap-3">
          <Button
            asChild
            variant="outline"
            size="icon"
            aria-label="返回"
            className="size-11 shrink-0 rounded-2xl"
          >
            <Link to="/">
              <ArrowLeft className="size-6" aria-hidden />
            </Link>
          </Button>
          <div className="min-w-0">
            <p className="text-muted-foreground mb-1 text-xs font-bold tracking-[0.08em] uppercase">
              {eyebrow}
            </p>
            <h1 className="text-[30px] leading-none font-extrabold tracking-normal">
              {title}
            </h1>
            {description ? (
              <p className="text-muted-foreground mt-1 text-sm leading-6 font-medium">
                {description}
              </p>
            ) : null}
          </div>
        </header>

        {children}
      </div>
    </main>
  )
}
