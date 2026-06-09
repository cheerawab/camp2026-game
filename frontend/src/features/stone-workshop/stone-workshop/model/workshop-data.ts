import type { PebbleTone } from "@/shared/config/color-palette"

export type StoneType = {
  key: "all" | PebbleTone
  label: string
  shortLabel: string
}

export type Stone = {
  id: string
  name: string
  type: PebbleTone
  rarity: "普通" | "稀有" | "特殊"
  count: number
  owned: boolean
  description: string
}

export type InventoryItem = {
  id: string
  name: string
  category: "素材" | "道具" | "任務"
  rarity: "普通" | "稀有" | "特殊"
  count: number
  description: string
}

export const stoneTypes: StoneType[] = [
  { key: "all", label: "全部", shortLabel: "ALL" },
  { key: "explore", label: "探索", shortLabel: "EXP" },
  { key: "spark", label: "靈光", shortLabel: "SPK" },
  { key: "resonate", label: "共鳴", shortLabel: "ECO" },
  { key: "engineer", label: "工程", shortLabel: "BLD" },
  { key: "play", label: "娛樂", shortLabel: "PLY" },
]

export const stones: Stone[] = [
  {
    id: "trail-map-chip",
    name: "路線圖小石",
    type: "explore",
    rarity: "普通",
    count: 3,
    owned: true,
    description: "記錄探索路線與關卡線索，適合用來開啟新的事件節點。",
  },
  {
    id: "blue-telescope",
    name: "望遠鏡小石",
    type: "explore",
    rarity: "稀有",
    count: 1,
    owned: true,
    description: "看見遠方的隱藏提示，能提高探索任務的判讀效率。",
  },
  {
    id: "camp-lantern",
    name: "營燈小石",
    type: "spark",
    rarity: "稀有",
    count: 2,
    owned: true,
    description: "在討論卡關時提供靈感，常用於解謎與提案任務。",
  },
  {
    id: "idea-flint",
    name: "火花小石",
    type: "spark",
    rarity: "普通",
    count: 5,
    owned: true,
    description: "小而穩定的靈感火種，可作為合成的基礎材料。",
  },
  {
    id: "team-radio",
    name: "無線電小石",
    type: "resonate",
    rarity: "普通",
    count: 4,
    owned: true,
    description: "強化隊伍合作與訊息傳遞，適合多人任務。",
  },
  {
    id: "memory-shell",
    name: "回聲小石",
    type: "resonate",
    rarity: "特殊",
    count: 0,
    owned: false,
    description: "尚未收集。據說會在完成跨隊協作後出現。",
  },
  {
    id: "solder-seed",
    name: "焊點小石",
    type: "engineer",
    rarity: "稀有",
    count: 1,
    owned: true,
    description: "用於修復裝置與觸發技術挑戰的工程型小石。",
  },
  {
    id: "green-module",
    name: "模組小石",
    type: "engineer",
    rarity: "普通",
    count: 2,
    owned: true,
    description: "能與素材組合成穩定模組，是常見合成基底。",
  },
  {
    id: "stage-token",
    name: "舞台小石",
    type: "play",
    rarity: "稀有",
    count: 1,
    owned: true,
    description: "參與活動與舞台事件時取得，帶有娛樂能量。",
  },
  {
    id: "night-badge",
    name: "夜間徽章小石",
    type: "play",
    rarity: "特殊",
    count: 0,
    owned: false,
    description: "尚未收集。可能和晚間限定任務有關。",
  },
]

export const inventoryItems: InventoryItem[] = [
  {
    id: "map-thread",
    name: "地圖線索",
    category: "素材",
    rarity: "普通",
    count: 8,
    description: "可用於探索型合成，也能兌換路線提示。",
  },
  {
    id: "camp-rivet",
    name: "營地鉚釘",
    category: "素材",
    rarity: "普通",
    count: 12,
    description: "穩定又常見的材料，適合工程型小石合成。",
  },
  {
    id: "lantern-ticket",
    name: "燈會票券",
    category: "道具",
    rarity: "稀有",
    count: 1,
    description: "可開啟一次活動事件或作為特殊合成素材。",
  },
  {
    id: "radio-pin",
    name: "頻道別針",
    category: "任務",
    rarity: "稀有",
    count: 2,
    description: "隊伍任務留下的證明，常和共鳴型小石相關。",
  },
  {
    id: "stage-strip",
    name: "舞台布條",
    category: "道具",
    rarity: "普通",
    count: 3,
    description: "活動區取得的道具，可以提高娛樂型合成成功率。",
  },
  {
    id: "decode-card",
    name: "解碼卡",
    category: "任務",
    rarity: "特殊",
    count: 1,
    description: "完成特殊任務取得，適合保留到高階合成。",
  },
]
