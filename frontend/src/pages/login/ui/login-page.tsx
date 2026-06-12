import { LoginPanel, StoneGlyph, stoneTypes } from "@/features/login"
import { GamePageShell } from "@/shared/ui/game-page-shell"

type LoginPageProps = {
  token?: string
}

export function LoginPage({ token }: LoginPageProps) {
  return (
    <GamePageShell
      ariaLabel="Camp 2026 Game 登入頁"
      contentClassName="justify-between overflow-hidden"
    >
      <div
        className="border-border/70 pointer-events-none absolute inset-x-4 top-24 bottom-28 rounded-[32px] border-2 border-dashed"
        aria-hidden="true"
      />
      {stoneTypes.map((type) => (
        <span
          className="pointer-events-none absolute z-0 opacity-70"
          key={type.name}
          style={{
            left: type.x,
            top: type.y,
            transform: `rotate(${type.rotate})`,
          }}
          aria-hidden="true"
        >
          <StoneGlyph type={type} tiny />
        </span>
      ))}

      <div
        className="relative z-10 mb-6 flex items-center justify-between text-sm font-black"
        aria-hidden="true"
      >
        <span>20:26</span>
        <span className="flex items-end gap-1">
          <i className="bg-ink block h-2 w-1.5 rounded-full" />
          <i className="bg-ink block h-3 w-1.5 rounded-full" />
          <i className="bg-ink block h-4 w-1.5 rounded-full" />
        </span>
      </div>

      <div className="relative z-10 my-auto">
        <LoginPanel key={token ?? ""} initialToken={token} />
      </div>
    </GamePageShell>
  )
}
