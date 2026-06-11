import { Link } from "@tanstack/react-router"
import { ArrowLeft } from "lucide-react"

import { PublicCodexPanel } from "@/features/public-codex"

export function PublicCodexPage() {
  return (
    <main className="bg-background text-foreground flex min-h-svh justify-center">
      <div className="flex min-h-svh w-full max-w-[430px] flex-col px-4 pt-[18px] pb-[30px]">
        <header className="mb-3.5 flex items-center gap-3">
          <Link
            to="/"
            aria-label="返回營隊基地"
            className="border-ink bg-card text-ink focus-visible:outline-power grid size-11 shrink-0 place-items-center rounded-2xl border-2 transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
          >
            <ArrowLeft className="size-5" aria-hidden />
          </Link>
          <div>
            <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
              PUBLIC CODEX
            </p>
            <h1 className="text-[30px] leading-[1.05] font-black tracking-tight">
              公開圖鑑
            </h1>
          </div>
        </header>
        <PublicCodexPanel />
      </div>
    </main>
  )
}
