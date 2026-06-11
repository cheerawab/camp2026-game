import { Link } from "@tanstack/react-router"
import { ArrowLeft } from "lucide-react"
import type { ReactNode } from "react"

import { Button } from "@/shared/ui/button"

type WorkshopPageShellProps = {
  title: string
  eyebrow: string
  children: ReactNode
}

export function WorkshopPageShell({
  title,
  eyebrow,
  children,
}: WorkshopPageShellProps) {
  return (
    <main className="bg-paper text-ink relative flex min-h-svh justify-center overflow-hidden">
      <section className="relative min-h-svh w-full max-w-[430px] overflow-hidden px-4 py-4 pb-7">
        <div
          className="border-moss/30 pointer-events-none absolute top-24 -right-14 size-[148px] rotate-[10deg] rounded-[42px] border-2 border-dashed"
          aria-hidden
        />
        <div
          className="border-primary/25 pointer-events-none absolute top-[368px] -left-9 size-24 rounded-full border-2"
          aria-hidden
        />

        <header className="relative z-10 grid grid-cols-[44px_1fr_auto] items-center gap-3 pt-1">
          <Link
            to="/"
            aria-label="返回營隊基地"
            className="border-ink bg-card text-ink focus-visible:outline-power grid min-h-11 place-items-center rounded-2xl border-2 transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
          >
            <ArrowLeft className="size-6" aria-hidden />
          </Link>
          <div>
            <p className="text-moss text-[11px] font-extrabold tracking-[0.08em] uppercase">
              {eyebrow}
            </p>
            <h1 className="text-[27px] leading-none font-extrabold tracking-normal">
              {title}
            </h1>
          </div>
          <Button
            type="button"
            variant="secondary"
            className="border-ink bg-surface-raised focus-visible:outline-power min-h-11 rounded-2xl border-2 px-3 text-[13px] font-extrabold transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
          >
            全圖鑑
          </Button>
        </header>

        <div className="relative z-10">{children}</div>
      </section>
    </main>
  )
}
