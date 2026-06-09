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
    <main className="bg-paper text-ink relative min-h-svh overflow-x-clip">
      {/* decorative dashed box — behind header */}
      <div
        className="border-moss/30 pointer-events-none fixed top-24 right-[max(0px,calc((100vw-430px)/2))] z-0 size-[148px] rotate-[10deg] rounded-[42px] border-2 border-dashed"
        aria-hidden
      />
      {/* fixed header */}
      <header className="bg-paper fixed top-0 left-1/2 z-20 w-[min(100%,430px)] -translate-x-1/2 grid grid-cols-[44px_1fr_auto] items-center gap-3 px-4 pt-6 pb-3">
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

      {/* decorative circle — behind content */}
      <div
        className="border-primary/25 pointer-events-none fixed top-[368px] left-[max(0px,calc((100vw-430px)/2-36px))] z-10 size-24 rounded-full border-2"
        aria-hidden
      />
      {/* content pushed below fixed header (~68px) */}
      <section className="relative z-10 mx-auto w-full max-w-[430px] px-4 pt-[92px] pb-7">
        {children}
      </section>
    </main>
  )
}
