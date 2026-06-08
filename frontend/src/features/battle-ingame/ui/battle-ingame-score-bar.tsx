type BattleIngameScoreBarType = {
  a: number
  b: number
}

export default function BattleIngameScoreBar({
  a,
  b,
}: BattleIngameScoreBarType) {
  const total = a + b
  const aRatio = total === 0 ? 0.5 : a / total
  const bRatio = total === 0 ? 0.5 : b / total

  return (
    <div className="bg-accent flex h-2 w-full overflow-hidden rounded-full">
      <div
        className="bg-primary transition-all"
        style={{
          width: `${aRatio * 100}%`,
        }}
      />
      <div
        className="bg-secondary transition-all"
        style={{
          width: `${bRatio * 100}%`,
        }}
      />
    </div>
  )
}
