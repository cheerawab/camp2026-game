import { useState } from "react"
import { Link } from "@tanstack/react-router"
import { ArrowLeft } from "lucide-react"

const BAG_ITEMS = [
  { id: "map-thread", name: "地圖棉線", type: "素材", rarity: "常見", count: 8, tone: "#4F8CC9", description: "用來標記基地佈景上的探索路線。" },
  { id: "camp-rivet", name: "營釘鉚扣", type: "素材", rarity: "常見", count: 12, tone: "#31A886", description: "可替收藏卡加上工程感邊框。" },
  { id: "lantern-ticket", name: "營燈佈景券", type: "外觀", rarity: "稀有", count: 1, tone: "#F4C84A", description: "替小隊基地換上暖黃營燈主題。" },
  { id: "radio-pin", name: "小隊電波徽章", type: "活動紀念", rarity: "稀有", count: 2, tone: "#E96F86", description: "和隊友完成同步挑戰後取得。" },
  { id: "stage-strip", name: "舞台彩帶", type: "外觀", rarity: "常見", count: 3, tone: "#9A75D6", description: "能讓展示櫃多一條娛樂系標記。" },
  { id: "decode-card", name: "解題提示卡", type: "活動紀念", rarity: "傳說", count: 1, tone: "#E76F3C", description: "完成一場高難度知識王戰後留下的紀念卡。" },
]

const FILTERS = ["全部", "素材", "外觀", "活動紀念"]

type BagItem = (typeof BAG_ITEMS)[number]

function ItemIcon({ tone, count }: { tone: string; count: number }) {
  return (
    <div className="inv-item-icon" style={{ "--tone": tone } as React.CSSProperties} aria-hidden="true">
      <span className="inv-stripe" />
      <strong>{count}</strong>
    </div>
  )
}

function ItemCard({ item }: { item: BagItem }) {
  return (
    <article className="inv-item-card">
      <ItemIcon tone={item.tone} count={item.count} />
      <div className="inv-item-copy">
        <div className="inv-item-head">
          <h3>{item.name}</h3>
          <strong className="inv-quantity">×{item.count}</strong>
        </div>
        <div className="inv-tags">
          <span>{item.type}</span>
          <span>{item.rarity}</span>
        </div>
        <p>{item.description}</p>
      </div>
    </article>
  )
}

function EmptyBag() {
  return (
    <section className="inv-empty-card" aria-label="空背包狀態">
      <div className="inv-empty-icon" aria-hidden="true">＋</div>
      <h3>這個分類目前沒有道具</h3>
      <p>可以到商店兌換，或在現場活動完成挑戰後取得。</p>
      <button type="button">前往商店</button>
    </section>
  )
}

const styles = `
  .inv-screen{min-height:100vh;display:flex;justify-content:center;background:#F5E9D2}
  .inv-canvas{width:min(100%,430px);min-height:100vh;padding:18px 16px 28px}
  .inv-top-bar{display:flex;align-items:center;gap:12px;margin-bottom:16px}
  .inv-top-bar h1{margin:0;font-size:28px;line-height:1.1;letter-spacing:-.04em}
  .inv-eyebrow,.inv-label{margin:0 0 3px;color:#6B725F;font-size:12px;font-weight:900;letter-spacing:.08em}
  .inv-icon-btn,.inv-small-btn,.inv-filters button,.inv-empty-card button{min-height:44px;border:2px solid #17233A;border-radius:16px;background:#FFF8E9;box-shadow:3px 3px 0 rgba(23,35,58,.18);font-weight:900;font-family:inherit;cursor:pointer}
  .inv-icon-btn{width:44px;font-size:28px;line-height:1}
  .inv-small-btn{margin-left:auto;padding:0 14px;background:#F4C84A}
  .inv-icon-btn:focus-visible,.inv-filters button:focus-visible{outline:3px solid #17233A;outline-offset:3px}
  .inv-icon-btn:active,.inv-filters button:active{transform:translate(1px,1px);box-shadow:1px 1px 0 rgba(23,35,58,.2)}
  .inv-summary-card{display:flex;justify-content:space-between;gap:16px;padding:20px;border:2px solid #17233A;border-radius:22px;background:#FFF8E9;box-shadow:5px 5px 0 rgba(23,35,58,.16)}
  .inv-summary-card strong{display:block;font-size:48px;line-height:.95;letter-spacing:-.05em}
  .inv-summary-card p{margin:8px 0 0;color:#6B725F;line-height:1.65}
  .inv-bag-mark{width:92px;min-width:92px;height:112px;border:2px solid #17233A;border-radius:22px 22px 16px 16px;background:#FFF4D4;position:relative}
  .inv-bag-mark::before{content:"";position:absolute;inset:16px 24px auto;height:24px;border:2px solid #17233A;border-bottom:0;border-radius:16px 16px 0 0}
  .inv-bag-mark span{position:absolute;left:20px;right:20px;bottom:22px;height:18px;border-top:3px solid #E76F3C;border-bottom:3px solid #356B58}
  .inv-stats-row{display:grid;grid-template-columns:repeat(3,1fr);gap:10px;margin:14px 0}
  .inv-stats-row div{padding:12px;border:2px solid #D8C29A;border-radius:18px;background:#FFF4D4}
  .inv-stats-row span{display:block;color:#6B725F;font-size:12px;font-weight:800}
  .inv-stats-row strong{display:block;margin-top:3px;font-size:22px}
  .inv-filters{display:flex;gap:8px;overflow-x:auto;padding:2px 0 12px}
  .inv-filters button{white-space:nowrap;padding:0 14px;box-shadow:none;border-color:#D8C29A}
  .inv-filters button.active{background:#17233A;color:#FFF8E9;border-color:#17233A}
  .inv-item-list{display:grid;gap:12px}
  .inv-item-card{display:grid;grid-template-columns:72px 1fr;gap:12px;align-items:start;padding:14px;border:2px solid #17233A;border-radius:22px;background:#FFF8E9}
  .inv-item-icon{height:72px;border:2px solid #17233A;border-radius:20px;background:var(--tone);position:relative;overflow:hidden}
  .inv-stripe{position:absolute;inset:14px 18px;border-top:3px solid rgba(255,248,233,.8);border-bottom:3px solid rgba(23,35,58,.25);transform:rotate(-12deg)}
  .inv-item-icon strong{position:absolute;right:6px;bottom:5px;min-width:28px;padding:1px 6px;border:2px solid #17233A;border-radius:999px;background:#FFF8E9;text-align:center;font-size:13px}
  .inv-item-head{display:flex;justify-content:space-between;gap:8px;align-items:start}
  .inv-item-head h3{margin:0;font-size:18px;letter-spacing:-.02em}
  .inv-quantity{color:#E76F3C;font-size:18px}
  .inv-tags{display:flex;gap:6px;margin:7px 0;flex-wrap:wrap}
  .inv-tags span{padding:3px 8px;border:1.5px solid #D8C29A;border-radius:999px;background:#FFF4D4;color:#6B725F;font-size:12px;font-weight:900}
  .inv-item-copy p{margin:0;color:#6B725F;line-height:1.62;font-size:14px}
  .inv-empty-card{padding:22px;border:2px dashed #17233A;border-radius:22px;background:#FFF8E9;text-align:left}
  .inv-empty-icon{width:54px;height:54px;display:grid;place-items:center;border:2px solid #17233A;border-radius:18px;background:#FFF4D4;font-size:28px;font-weight:950}
  .inv-empty-card h3{margin:12px 0 4px}
  .inv-empty-card p{margin:0 0 14px;color:#6B725F;line-height:1.65}
  .inv-empty-card button{width:100%;background:#E76F3C;color:#FFF8E9}
`

export function InventoryPage() {
  const [filter, setFilter] = useState("全部")
  const visibleItems =
    filter === "全部" ? BAG_ITEMS : BAG_ITEMS.filter((item) => item.type === filter)
  const totalCount = BAG_ITEMS.reduce((sum, item) => sum + item.count, 0)

  return (
    <main className="inv-screen" aria-label="道具背包頁">
      <style>{styles}</style>
      <section className="inv-canvas">
        <header className="inv-top-bar">
          <Link to="/" aria-label="返回" className="border-ink bg-card text-ink focus-visible:outline-power grid size-11 shrink-0 place-items-center rounded-2xl border-2 transition-transform focus-visible:outline-3 focus-visible:outline-offset-2 active:translate-y-px" style={{ textDecoration: "none" }}><ArrowLeft className="size-5" aria-hidden /></Link>
          <div>
            <p className="inv-eyebrow">FIELD BAG</p>
            <h1>道具背包</h1>
          </div>
          <button className="inv-small-btn">整理</button>
        </header>

        <section className="inv-summary-card" aria-label="背包摘要">
          <div>
            <span className="inv-label">目前持有</span>
            <strong>{totalCount}</strong>
            <p>素材、外觀與活動紀念都會先收在這裡。</p>
          </div>
          <div className="inv-bag-mark" aria-hidden="true"><span /></div>
        </section>

        <section className="inv-stats-row" aria-label="分類數量">
          <div><span>素材</span><strong>20</strong></div>
          <div><span>外觀</span><strong>4</strong></div>
          <div><span>紀念</span><strong>3</strong></div>
        </section>

        <nav className="inv-filters" aria-label="背包分類">
          {FILTERS.map((f) => (
            <button
              key={f}
              className={filter === f ? "active" : ""}
              onClick={() => setFilter(f)}
            >
              {f}
            </button>
          ))}
        </nav>

        <section className="inv-item-list" aria-label="道具列表">
          {visibleItems.length > 0
            ? visibleItems.map((item) => <ItemCard key={item.id} item={item} />)
            : <EmptyBag />}
        </section>
      </section>
    </main>
  )
}
