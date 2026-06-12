export type SitoneToneKey = "explore" | "spark" | "echo" | "build" | "play"

type SitoneMeta = {
  key: SitoneToneKey
  label: string
  short: string
  bgClassName: string
}

const sitoneTypeMap: Record<string, SitoneMeta> = {
  exploration: {
    key: "explore",
    label: "探索",
    short: "EXP",
    bgClassName: "bg-pebble-explore",
  },
  inspiration: {
    key: "spark",
    label: "靈光",
    short: "SPK",
    bgClassName: "bg-pebble-spark",
  },
  resonance: {
    key: "echo",
    label: "共鳴",
    short: "ECO",
    bgClassName: "bg-pebble-resonate",
  },
  engineering: {
    key: "build",
    label: "工程",
    short: "BLD",
    bgClassName: "bg-pebble-engineer",
  },
  entertainment: {
    key: "play",
    label: "娛樂",
    short: "PLY",
    bgClassName: "bg-pebble-play",
  },
}

const fallbackSitoneMeta: SitoneMeta = {
  key: "explore",
  label: "小石",
  short: "STN",
  bgClassName: "bg-primary",
}

const itemTypeLabels: Record<string, string> = {
  material: "素材",
  cosmetic: "外觀",
  event: "活動紀念",
}

const itemTypeClasses: Record<string, string> = {
  material: "bg-pebble-engineer",
  cosmetic: "bg-pebble-spark",
  event: "bg-pebble-resonate",
}

const rarityLabels: Record<string, string> = {
  base: "基礎",
  common: "常見",
  rare: "稀有",
  limited: "限定",
}

export function sitoneMeta(type: string): SitoneMeta {
  return sitoneTypeMap[type] ?? fallbackSitoneMeta
}

export function itemTypeLabel(type: string) {
  return itemTypeLabels[type] ?? type
}

export function itemTypeClass(type: string) {
  return itemTypeClasses[type] ?? "bg-primary"
}

export function rarityLabel(rarity: string) {
  return rarityLabels[rarity] ?? rarity
}

export function rarityToneClass(rarity: string) {
  switch (rarity) {
    case "rare":
    case "稀有":
      return "bg-pebble-explore-muted"
    case "limited":
    case "限定":
      return "bg-pebble-play-muted"
    default:
      return "bg-pebble-engineer-muted"
  }
}
