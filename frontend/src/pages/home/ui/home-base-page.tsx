import { useQuery } from "@tanstack/react-query"
import { Link, useNavigate } from "@tanstack/react-router"
import { useEffect } from "react"
import { LogOut } from "lucide-react"
import { AppError } from "@/shared/api/error"
import { gameApi } from "@/shared/api/game"
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/avatar"
import { GamePageShell } from "@/shared/ui/game-page-shell"
import { Button } from "@/shared/ui/button"
import { apiClient } from "@/shared/api/client"

const ACTIONS: {
  label: string
  desc: string
  colorClass: string
  to: string
  primary?: boolean
}[] = [
  {
    label: "知識王戰",
    desc: "建立或掃碼加入對戰",
    colorClass: "bg-primary",
    primary: true,
    to: "/battle",
  },
  {
    label: "個人 QR Code",
    desc: "現場關卡驗證身份",
    colorClass: "bg-moss",
    to: "/profile/qr",
  },
  {
    label: "商店",
    desc: "使用開源力兌換外觀",
    colorClass: "bg-pebble-spark",
    to: "/shop",
  },
]

const STAFF_ACTION: (typeof ACTIONS)[number] = {
  label: "工作人員發放",
  desc: "掃描 QR Code 發小石或道具",
  colorClass: "bg-secondary",
  to: "/staff",
}

const COLLECTIONS: {
  label: string
  count: (input: { sitones: number; items: number; rank?: number }) => string
  colorClass: string
  to: string
}[] = [
  {
    label: "小石收藏",
    count: ({ sitones }) => `${sitones} 顆`,
    colorClass: "bg-pebble-explore",
    to: "/stones",
  },
  {
    label: "道具背包",
    count: ({ items }) => `${items} 件`,
    colorClass: "bg-pebble-engineer",
    to: "/inventory",
  },
  {
    label: "小石合成",
    count: () => "工作台",
    colorClass: "bg-pebble-play",
    to: "/stones/fusion",
  },
  {
    label: "排行榜",
    count: ({ rank }) => (rank ? `#${rank}` : "查看"),
    colorClass: "bg-pebble-resonate",
    to: "/leaderboard",
  },
  {
    label: "對戰紀錄",
    count: () => "最近一場",
    colorClass: "bg-primary",
    to: "/battle/history",
  },
  {
    label: "公開圖鑑",
    count: () => "查詢",
    colorClass: "bg-moss/60",
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
    (player?.role === "staff" ? "工作人員" : "讀取玩家資料")
  const avatarInitial = displayName.trim().slice(0, 1) || "?"
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
    <GamePageShell ariaLabel="營隊基地首頁" contentClassName="gap-3">
      <section className="grid content-start gap-3">
        {/* 玩家狀態 */}
        <header
          className="bg-card border-ink grid grid-cols-[64px_1fr_auto] items-center gap-3 rounded-[26px] border-2 p-3.5"
          style={{ boxShadow: "4px 4px 0 rgba(23,35,58,.14)" }}
          aria-label="玩家狀態"
        >
          <Avatar
            size="lg"
            className="bg-pebble-spark border-ink size-16 rounded-[22px] border-2 text-[26px]"
            aria-hidden
          >
            {player?.avatarUrl ? (
              <AvatarImage
                src={player.avatarUrl}
                alt=""
                className="object-cover"
              />
            ) : null}
            <AvatarFallback className="bg-pebble-spark text-foreground rounded-[20px] text-[26px] font-black">
              {avatarInitial}
            </AvatarFallback>
          </Avatar>
          <div>
            <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
              Camp Base
            </p>
            <h1 className="text-[29px] leading-none font-black tracking-[-0.04em]">
              {displayName}
            </h1>
            <span className="text-muted-foreground font-bold">
              {isPending ? "同步玩家資料中" : teamName}
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
              aria-label="登出"
              onClick={() => void logoutAction()}
            >
              <LogOut aria-hidden />
            </Button>
          </div>
        </header>

        {/* 主要行動 */}
        <section
          className="bg-ink text-primary-foreground rounded-[30px] border-2 border-transparent p-[18px]"
          style={{ boxShadow: "5px 5px 0 rgba(23,35,58,.16)" }}
          aria-label="主要行動"
        >
          <p className="text-primary-foreground/75 mb-1 text-xs font-black tracking-[0.08em] uppercase">
            現在最重要
          </p>
          <h2 className="mb-2 text-[27px] leading-[1.16] font-black tracking-[-0.045em]">
            先開始知識王戰，其他都放在下方快速入口。
          </h2>
          <p className="text-primary-foreground/75 mb-3.5 leading-[1.65]">
            首頁不放每日任務或世界事件；只放玩家現場真的會點的功能。
          </p>
          <Link
            to="/battle"
            className="border-ink bg-primary text-primary-foreground focus-visible:outline-power grid min-h-[46px] w-full place-items-center rounded-[17px] border-2 text-base font-black no-underline transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
            style={{ boxShadow: "3px 3px 0 rgba(0,0,0,.18)" }}
          >
            開始 / 加入對戰
          </Link>
        </section>

        {/* 快速摘要 */}
        <section className="grid grid-cols-3 gap-[9px]" aria-label="快速摘要">
          {[
            { label: "小石", value: sitoneCount },
            { label: "道具", value: itemCount },
            { label: "排行", value: rank ? `#${rank}` : "-" },
          ].map(({ label, value }) => (
            <div
              key={label}
              className="bg-surface-raised border-border rounded-[18px] border-2 px-2 py-3 text-center"
            >
              <span className="text-muted-foreground block text-xs font-black">
                {label}
              </span>
              <strong className="text-[24px] font-black">{value}</strong>
            </div>
          ))}
        </section>

        {/* 同組成員 */}
        <section
          className="bg-card border-ink rounded-[22px] border-2 p-[15px]"
          aria-label="同組成員"
        >
          <div className="mb-3 flex items-start justify-between gap-3">
            <div>
              <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
                Team
              </p>
              <h2 className="text-[22px] font-black tracking-[-0.04em]">
                同組成員
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
                const initial = member.nickname.trim().slice(0, 1) || "?"

                return (
                  <li
                    key={member.playerId}
                    className="bg-surface-raised border-border grid min-h-[56px] grid-cols-[40px_1fr_auto] items-center gap-3 rounded-[17px] border-2 px-3 py-2"
                  >
                    <Avatar
                      size="lg"
                      className="border-ink bg-pebble-resonate border-2"
                    >
                      {member.avatarUrl ? (
                        <AvatarImage
                          src={member.avatarUrl}
                          alt=""
                          className="object-cover"
                        />
                      ) : null}
                      <AvatarFallback className="bg-pebble-resonate text-primary-foreground font-black">
                        {initial}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0">
                      <strong className="block truncate text-[16px] font-black">
                        {member.nickname}
                      </strong>
                      {member.role === "staff" ? (
                        <small className="text-muted-foreground block text-xs font-bold">
                          工作人員
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
              還沒有同組成員資料
            </p>
          )}
        </section>

        {/* 核心入口 */}
        <section className="grid gap-[10px]" aria-label="核心入口">
          {actions.map((action) => (
            <article
              key={action.label}
              className={[
                "bg-card border-ink grid grid-cols-[42px_1fr_68px] items-center gap-[10px] rounded-[22px] border-2 p-[13px]",
                action.primary ? "bg-surface-raised" : "",
              ].join(" ")}
            >
              <div
                className={`border-ink size-[42px] -rotate-[7deg] rounded-[16px_20px_14px_18px] border-2 ${action.colorClass}`}
                aria-hidden
              />
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
                className="bg-card border-ink focus-visible:outline-power grid min-h-[40px] place-items-center rounded-[17px] border-2 text-sm font-black no-underline transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px"
                style={{ boxShadow: "2px 2px 0 rgba(23,35,58,.14)" }}
              >
                開啟
              </Link>
            </article>
          ))}
        </section>

        {/* 收藏與查詢 */}
        <section
          className="bg-card border-ink rounded-[22px] border-2 p-[15px]"
          aria-label="收藏與查詢"
        >
          <p className="text-muted-foreground mb-1 text-xs font-black tracking-[0.08em] uppercase">
            Collect &amp; Check
          </p>
          <h2 className="mb-3 text-[22px] font-black tracking-[-0.04em]">
            收藏、背包、合成、排行、紀錄
          </h2>
          <div className="grid grid-cols-2 gap-[9px]">
            {COLLECTIONS.map((item) => (
              <Link
                key={item.label}
                to={item.to}
                className="bg-surface-raised border-ink grid min-h-[66px] grid-cols-[24px_1fr] items-center gap-[7px] rounded-[17px] border-2 px-[10px] py-[10px] text-inherit no-underline transition-transform active:translate-y-px"
              >
                <span
                  className={`border-ink row-span-2 size-6 rounded-[9px_12px_8px_10px] border-2 ${item.colorClass}`}
                  aria-hidden
                />
                <strong className="block font-black">{item.label}</strong>
                <small className="text-muted-foreground block text-xs font-bold">
                  {item.count({
                    sitones: sitoneCount,
                    items: itemCount,
                    rank,
                  })}
                </small>
              </Link>
            ))}
          </div>
        </section>
      </section>
    </GamePageShell>
  )
}
