import { useQuery } from "@tanstack/react-query"
import { Link, useNavigate } from "@tanstack/react-router"
import { useEffect } from "react"
import {
  Backpack,
  BookOpenText,
  ChevronRight,
  Crown,
  Hammer,
  History,
  LibraryBig,
  LogOut,
  PackageOpen,
  QrCode,
  ScanLine,
  ShieldCheck,
  ShoppingBag,
  Swords,
  Trophy,
  type LucideIcon,
} from "lucide-react"
import { AppError } from "@/shared/api/error"
import { gameApi } from "@/shared/api/game"
import { GamePageShell } from "@/shared/ui/game-page-shell"
import { Button } from "@/shared/ui/button"
import { apiClient } from "@/shared/api/client"
import { PlayerAvatar } from "@/shared/ui/player-avatar"

const ACTIONS: {
  label: string
  desc: string
  colorClass: string
  icon: LucideIcon
  to: string
  primary?: boolean
}[] = [
  {
    label: "知識王試煉",
    desc: "召開戰局或掃描戰碼進場",
    colorClass: "bg-primary",
    icon: Swords,
    primary: true,
    to: "/battle",
  },
  {
    label: "冒險者通行證",
    desc: "讓關主驗證你的玩家身份",
    colorClass: "bg-moss",
    icon: QrCode,
    to: "/profile/qr",
  },
  {
    label: "補給商店",
    desc: "消耗開源力兌換素材與外觀",
    colorClass: "bg-pebble-spark",
    icon: ShoppingBag,
    to: "/shop",
  },
]

const STAFF_ACTION: (typeof ACTIONS)[number] = {
  label: "關主發放台",
  desc: "掃描通行證發放小石與戰利品",
  colorClass: "bg-secondary",
  icon: ScanLine,
  to: "/staff",
}

const COLLECTIONS: {
  label: string
  count: (input: { sitones: number; items: number; rank?: number }) => string
  colorClass: string
  icon: LucideIcon
  to: string
}[] = [
  {
    label: "小石圖鑑",
    count: ({ sitones }) => `已收 ${sitones} 顆`,
    colorClass: "bg-pebble-explore",
    icon: BookOpenText,
    to: "/stones",
  },
  {
    label: "戰利品背包",
    count: ({ items }) => `持有 ${items} 件`,
    colorClass: "bg-pebble-engineer",
    icon: Backpack,
    to: "/inventory",
  },
  {
    label: "鍛造工坊",
    count: () => "合成開放",
    colorClass: "bg-pebble-play",
    icon: Hammer,
    to: "/stones/fusion",
  },
  {
    label: "公會排行",
    count: ({ rank }) => (rank ? `第 ${rank} 名` : "未入榜"),
    colorClass: "bg-pebble-resonate",
    icon: Trophy,
    to: "/leaderboard",
  },
  {
    label: "戰鬥回放",
    count: () => "最近戰局",
    colorClass: "bg-primary",
    icon: History,
    to: "/battle/history",
  },
  {
    label: "全服圖鑑",
    count: () => "偵查資料",
    colorClass: "bg-moss/60",
    icon: LibraryBig,
    to: "/codex",
  },
]

export function HomeBasePage() {
  const navigate = useNavigate()
  const { data, isPending, error } = useQuery({
    queryKey: ["me", "home"],
    queryFn: gameApi.home,
  })
  const unauthorized = error instanceof AppError && error.status === 401

  useEffect(() => {
    if (unauthorized) {
      navigate({ to: "/login", replace: true })
    }
  }, [navigate, unauthorized])

  if (unauthorized) return null

  const player = data?.player
  const summary = data?.summary
  const teamRank = data?.teamRank
  const displayName = player?.nickname ?? "載入中"
  const teamName =
    player?.team?.name ??
    (player?.role === "staff" ? "關主模式" : "讀取小隊資料")
  const openPower = summary?.openPower ?? player?.openPower ?? 0
  const sitoneCount = summary?.sitoneCount ?? 0
  const itemCount = summary?.itemCount ?? 0
  const rank = teamRank?.rank
  const teamMembers = player?.teamMembers ?? []
  const actions =
    player?.role === "staff" ? [STAFF_ACTION, ...ACTIONS] : ACTIONS
  const logoutAction = async () => {
    await apiClient.post("/api/auth/logout")
    navigate({ to: "/login", replace: true })
  }

  return (
    <GamePageShell ariaLabel="冒險基地首頁" contentClassName="gap-3">
      <section className="grid content-start gap-3">
        <header
          className="bg-card border-ink grid grid-cols-[64px_1fr_auto] items-center gap-3 rounded-[26px] border-2 p-3.5"
          style={{ boxShadow: "4px 4px 0 rgba(23,35,58,.14)" }}
          aria-label="冒險者狀態"
        >
          <PlayerAvatar
            playerId={player?.playerId}
            nickname={displayName}
            size="lg"
            className="bg-pebble-spark border-ink size-16 rounded-[22px] border-2 text-[26px]"
          />
          <div>
            <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
              Player Card
            </p>
            <h1 className="text-[29px] leading-none font-black tracking-normal">
              {displayName}
            </h1>
            <span className="text-muted-foreground font-bold">
              {isPending ? "召回冒險者資料中" : teamName}
            </span>
          </div>
          <div className="flex gap-x-3">
            <div
              className="bg-ink text-primary-foreground min-w-[86px] rounded-[18px] border-2 border-transparent px-[9px] py-2 text-center"
              aria-label={`開源力 ${openPower}`}
            >
              <span className="text-primary-foreground/70 block text-[11px] font-black">
                開源力
              </span>
              <strong className="text-[23px] font-black">{openPower}</strong>
            </div>
            <Button
              type="button"
              size="icon"
              variant="destructive"
              aria-label="離開基地"
              onClick={() => void logoutAction()}
            >
              <LogOut aria-hidden />
            </Button>
          </div>
        </header>

        <section
          className="bg-ink text-primary-foreground relative overflow-hidden rounded-[24px] border-2 border-transparent p-[18px]"
          style={{ boxShadow: "5px 5px 0 rgba(23,35,58,.16)" }}
          aria-label="主線任務"
        >
          <div
            className="border-primary/70 absolute top-3 right-3 grid size-12 place-items-center border-2 bg-white/5"
            aria-hidden
          >
            <Crown className="text-secondary size-7" />
          </div>
          <div
            className="border-primary/35 pointer-events-none absolute inset-x-4 bottom-4 h-3 border-x-2 border-b-2"
            aria-hidden
          />
          <p className="text-primary-foreground/75 mb-1 text-xs font-black tracking-[0.08em] uppercase">
            Main Quest
          </p>
          <h2 className="mb-2 max-w-[280px] text-[31px] leading-[1.08] font-black tracking-normal">
            知識王試煉開放
          </h2>
          <p className="text-primary-foreground/75 mb-4 max-w-[300px] leading-[1.65]">
            集結隊友、掃描戰碼，進入答題競技場奪取開源力。
          </p>
          <Link
            to="/battle"
            className="border-ink bg-primary text-primary-foreground focus-visible:outline-power relative z-10 flex min-h-[50px] w-full items-center justify-center gap-2 rounded-[14px] border-2 text-base font-black no-underline transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
            style={{ boxShadow: "3px 3px 0 rgba(0,0,0,.18)" }}
          >
            <Swords className="size-5" aria-hidden />
            進入競技場
            <ChevronRight className="size-5" aria-hidden />
          </Link>
        </section>

        <section className="grid grid-cols-3 gap-[9px]" aria-label="戰力摘要">
          {[
            {
              label: "小石隊伍",
              value: sitoneCount,
              to: "/stones",
              icon: ShieldCheck,
            },
            {
              label: "戰利品",
              value: itemCount,
              to: "/inventory",
              icon: PackageOpen,
            },
            {
              label: "名次",
              value: rank ? `#${rank}` : "-",
              to: "/leaderboard",
              icon: Crown,
            },
          ].map(({ label, value, to, icon: SummaryIcon }) => (
            <Link
              key={label}
              to={to}
              aria-label={`查看${label}`}
              className="bg-surface-raised border-border focus-visible:outline-power grid min-h-[84px] justify-items-center rounded-[16px] border-2 px-2 py-3 text-center text-inherit no-underline transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
            >
              <SummaryIcon className="text-muted-foreground mb-1 size-4" />
              <span className="text-muted-foreground block text-xs font-black">
                {label}
              </span>
              <strong className="text-[24px] font-black">{value}</strong>
            </Link>
          ))}
        </section>

        <section
          className="bg-card border-ink rounded-[22px] border-2 p-[15px]"
          aria-label="小隊成員"
        >
          <div className="mb-3 flex items-start justify-between gap-3">
            <div>
              <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
                Party
              </p>
              <h2 className="text-[22px] font-black tracking-normal">
                小隊名冊
              </h2>
            </div>
            <span className="bg-surface-raised border-border rounded-full border-2 px-2.5 py-1 text-xs font-black whitespace-nowrap">
              {isPending ? "-" : `${teamMembers.length} 人`}
            </span>
          </div>

          {isPending ? (
            <div className="grid gap-[8px]">
              {[0, 1, 2].map((item) => (
                <div
                  key={item}
                  className="bg-surface-raised border-border grid min-h-[56px] grid-cols-[40px_1fr] items-center gap-3 rounded-[17px] border-2 px-3"
                >
                  <span className="bg-muted border-border size-10 rounded-full border-2" />
                  <span className="bg-muted h-4 w-28 rounded-full" />
                </div>
              ))}
            </div>
          ) : teamMembers.length > 0 ? (
            <ul className="grid gap-[8px]">
              {teamMembers.map((member) => {
                const current = member.playerId === player?.playerId
                return (
                  <li
                    key={member.playerId}
                    className="bg-surface-raised border-border grid min-h-[56px] grid-cols-[40px_1fr_auto] items-center gap-3 rounded-[17px] border-2 px-3 py-2"
                  >
                    <PlayerAvatar
                      playerId={member.playerId}
                      nickname={member.nickname}
                      size="lg"
                      className="border-ink bg-pebble-resonate border-2"
                    />
                    <div className="min-w-0">
                      <strong className="block truncate text-[16px] font-black">
                        {member.nickname}
                      </strong>
                      {member.role === "staff" ? (
                        <small className="text-muted-foreground block text-xs font-bold">
                          關主
                        </small>
                      ) : null}
                    </div>
                    {current ? (
                      <span className="bg-secondary text-secondary-foreground border-ink rounded-full border-2 px-2 py-0.5 text-xs font-black whitespace-nowrap">
                        你
                      </span>
                    ) : null}
                  </li>
                )
              })}
            </ul>
          ) : (
            <p className="text-muted-foreground bg-surface-raised border-border rounded-[17px] border-2 px-3 py-3 text-sm font-bold">
              尚未取得小隊名冊
            </p>
          )}
        </section>

        <section className="grid gap-[10px]" aria-label="任務入口">
          {actions.map((action) => {
            const ActionIcon = action.icon

            return (
              <article
                key={action.label}
                className={[
                  "bg-card border-ink grid grid-cols-[46px_1fr_78px] items-center gap-[10px] rounded-[18px] border-2 p-[13px]",
                  action.primary ? "bg-surface-raised" : "",
                ].join(" ")}
              >
                <div
                  className={`border-ink grid size-[46px] -rotate-[4deg] place-items-center rounded-[14px] border-2 ${action.colorClass}`}
                  aria-hidden
                >
                  <ActionIcon className="size-6" />
                </div>
                <div>
                  <h3 className="mb-[3px] text-[18px] font-black">
                    {action.label}
                  </h3>
                  <p className="text-muted-foreground m-0 text-[13px] leading-[1.45]">
                    {action.desc}
                  </p>
                </div>
                <Link
                  to={action.to}
                  className="bg-card border-ink focus-visible:outline-power flex min-h-[40px] items-center justify-center gap-1 rounded-[13px] border-2 text-sm font-black no-underline transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
                  style={{ boxShadow: "2px 2px 0 rgba(23,35,58,.14)" }}
                >
                  出發
                  <ChevronRight className="size-4" aria-hidden />
                </Link>
              </article>
            )
          })}
        </section>

        <section
          className="bg-card border-ink rounded-[22px] border-2 p-[15px]"
          aria-label="冒險手冊"
        >
          <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
            Adventurer Log
          </p>
          <h2 className="mb-3 text-[22px] font-black tracking-normal">
            圖鑑、背包、鍛造與戰績
          </h2>
          <div className="grid grid-cols-2 gap-[9px]">
            {COLLECTIONS.map((item) => {
              const CollectionIcon = item.icon

              return (
                <Link
                  key={item.label}
                  to={item.to}
                  className="bg-surface-raised border-ink grid min-h-[72px] grid-cols-[30px_1fr] items-center gap-[8px] rounded-[15px] border-2 px-[10px] py-[10px] text-inherit no-underline transition-transform active:translate-y-px"
                >
                  <span
                    className={`border-ink row-span-2 grid size-7 place-items-center rounded-[10px] border-2 ${item.colorClass}`}
                    aria-hidden
                  >
                    <CollectionIcon className="size-4" />
                  </span>
                  <strong className="block font-black">{item.label}</strong>
                  <small className="text-muted-foreground block text-xs font-bold">
                    {item.count({
                      sitones: sitoneCount,
                      items: itemCount,
                      rank,
                    })}
                  </small>
                </Link>
              )
            })}
          </div>
        </section>
      </section>
    </GamePageShell>
  )
}
