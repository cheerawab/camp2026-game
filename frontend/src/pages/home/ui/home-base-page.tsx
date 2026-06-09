import { Link } from "@tanstack/react-router"

const PLAYER = { name: "阿洛", squad: "松鼠小隊", power: 430, stones: 8, items: 29, rank: 12 }

const ACTIONS: { label: string; desc: string; tone: string; to: string; primary?: boolean }[] = [
  { label: "知識王戰", desc: "建立或掃碼加入對戰", tone: "#E76F3C", primary: true, to: "/battle" },
  { label: "個人 QR Code", desc: "現場關卡驗證身份", tone: "#356B58", to: "/profile/qr" },
  { label: "商店", desc: "使用開源力兌換外觀", tone: "#F4C84A", to: "/shop" },
]

const COLLECTIONS = [
  { label: "小石收藏", count: "8 種", tone: "#4F8CC9", to: "/stones" },
  { label: "道具背包", count: "29 件", tone: "#31A886", to: "/inventory" },
  { label: "小石合成", count: "工作台", tone: "#9A75D6", to: "/stones/fusion" },
  { label: "排行榜", count: `#${PLAYER.rank}`, tone: "#E96F86", to: "/leaderboard" },
  { label: "公開圖鑑", count: "查詢", tone: "#6B725F", to: "/codex" },
] as const

const styles = `
  .hb-screen{min-height:100vh;display:flex;justify-content:center;background:#F5E9D2}
  .hb-canvas{width:min(100%,430px);min-height:100vh;padding:18px 16px 30px}
  .hb-btn{min-height:46px;border:2px solid #17233A;border-radius:17px;background:#E76F3C;color:#FFF8E9;box-shadow:3px 3px 0 rgba(0,0,0,.18);font-weight:950;font-family:inherit;cursor:pointer}
  .hb-btn:focus-visible{outline:3px solid #17233A;outline-offset:3px}
  .hb-btn:active{transform:translate(1px,1px)}
  .hb-eyebrow{margin:0 0 4px;color:#6B725F;font-size:12px;font-weight:950;letter-spacing:.08em}
  .hb-player-card{display:grid;grid-template-columns:64px 1fr auto;gap:12px;align-items:center;padding:14px;border:2px solid #17233A;border-radius:26px;background:#FFF8E9;box-shadow:4px 4px 0 rgba(23,35,58,.14)}
  .hb-avatar{width:64px;height:64px;display:grid;place-items:center;border:2px solid #17233A;border-radius:22px;background:#F4C84A;font-size:26px;font-weight:950}
  .hb-player-copy h1{margin:0;font-size:29px;line-height:1;letter-spacing:-.04em}
  .hb-player-copy span{color:#6B725F;font-weight:850}
  .hb-power-chip{min-width:76px;padding:8px 9px;border:2px solid #17233A;border-radius:18px;background:#17233A;color:#FFF8E9;text-align:center}
  .hb-power-chip span{display:block;color:rgba(255,248,233,.72);font-size:11px;font-weight:950}
  .hb-power-chip strong{font-size:23px}
  .hb-primary-panel{margin-top:14px;padding:18px;border:2px solid #17233A;border-radius:30px;background:#17233A;color:#FFF8E9;box-shadow:5px 5px 0 rgba(23,35,58,.16)}
  .hb-primary-panel .hb-eyebrow,.hb-primary-panel p{color:rgba(255,248,233,.75)}
  .hb-primary-panel h2{margin:0;font-size:27px;line-height:1.16;letter-spacing:-.045em}
  .hb-primary-panel p{margin:8px 0 14px;line-height:1.65}
  .hb-start-battle{width:100%}
  .hb-quick-stats{margin-top:12px;display:grid;grid-template-columns:repeat(3,1fr);gap:9px}
  .hb-quick-stats div{padding:12px 8px;border:2px solid #D8C29A;border-radius:18px;background:#FFF4D4;text-align:center}
  .hb-quick-stats span{display:block;color:#6B725F;font-size:12px;font-weight:950}
  .hb-quick-stats strong{font-size:24px}
  .hb-action-grid{margin-top:14px;display:grid;gap:10px}
  .hb-action-card{display:grid;grid-template-columns:42px 1fr 68px;gap:10px;align-items:center;padding:13px;border:2px solid #17233A;border-radius:22px;background:#FFF8E9}
  .hb-action-card.primary{background:#FFF4D4}
  .hb-stone-dot{width:42px;height:42px;border:2px solid #17233A;border-radius:16px 20px 14px 18px;transform:rotate(-7deg)}
  .hb-action-card h3{margin:0 0 3px;font-size:18px}
  .hb-action-card p{margin:0;color:#6B725F;font-size:13px;line-height:1.45}
  .hb-action-card .hb-btn{min-height:40px;background:#FFF8E9;color:#17233A;box-shadow:2px 2px 0 rgba(23,35,58,.14)}
  .hb-collection-panel,.hb-base-snapshot{margin-top:14px;padding:15px;border:2px solid #17233A;border-radius:22px;background:#FFF8E9}
  .hb-section-head h2{margin:0 0 12px;font-size:22px;letter-spacing:-.04em}
  .hb-collection-grid{display:grid;grid-template-columns:1fr 1fr;gap:9px}
  .hb-collection-tile{display:grid;grid-template-columns:24px 1fr;gap:7px;align-items:center;min-height:66px;padding:10px;background:#FFF4D4;color:#17233A;text-align:left;border:2px solid #17233A;border-radius:17px;font-family:inherit;cursor:pointer}
  .hb-collection-tile small{grid-column:2;color:#6B725F;font-size:12px}
  .hb-collection-tile strong{font-weight:900}
  .hb-tile-mark{grid-row:span 2;width:24px;height:24px;border:2px solid #17233A;border-radius:9px 12px 8px 10px}
  .hb-base-snapshot{display:grid;grid-template-columns:82px 1fr;gap:13px;align-items:center;background:#FFF4D4}
  .hb-mini-map{position:relative;height:76px;border:2px solid #17233A;border-radius:20px;background:#FFF8E9}
  .hb-mini-map span{position:absolute;width:18px;height:18px;border:2px solid #17233A;border-radius:8px}
  .hb-mini-map span:nth-child(1){left:14px;top:14px;background:#F4C84A}
  .hb-mini-map span:nth-child(2){right:15px;top:24px;background:#E76F3C}
  .hb-mini-map span:nth-child(3){left:31px;bottom:12px;background:#356B58}
  .hb-base-snapshot h2{margin:0 0 5px;font-size:18px}
  .hb-base-snapshot p{margin:0;color:#6B725F;line-height:1.55;font-size:13px}
`

export function HomeBasePage() {
  return (
    <main className="hb-screen" aria-label="營隊基地首頁">
      <style>{styles}</style>
      <section className="hb-canvas">
        <header className="hb-player-card" aria-label="玩家狀態">
          <div className="hb-avatar" aria-hidden="true">洛</div>
          <div className="hb-player-copy">
            <p className="hb-eyebrow">CAMP BASE</p>
            <h1>{PLAYER.name}</h1>
            <span>{PLAYER.squad}</span>
          </div>
          <div className="hb-power-chip" aria-label={`開源力 ${PLAYER.power}`}>
            <span>OP</span>
            <strong>{PLAYER.power}</strong>
          </div>
        </header>

        <section className="hb-primary-panel" aria-label="主要行動">
          <div>
            <p className="hb-eyebrow">現在最重要</p>
            <h2>先開始知識王戰，其他都放在下方快速入口。</h2>
            <p>首頁不放每日任務或世界事件；只放玩家現場真的會點的功能。</p>
          </div>
          <Link to="/battle" className="hb-btn hb-start-battle" style={{ display: "grid", placeItems: "center", textDecoration: "none" }}>開始 / 加入對戰</Link>
        </section>

        <section className="hb-quick-stats" aria-label="快速摘要">
          <div><span>小石</span><strong>{PLAYER.stones}</strong></div>
          <div><span>道具</span><strong>{PLAYER.items}</strong></div>
          <div><span>排行</span><strong>#{PLAYER.rank}</strong></div>
        </section>

        <section className="hb-action-grid" aria-label="核心入口">
          {ACTIONS.map((action) => (
            <article
              key={action.label}
              className={`hb-action-card${action.primary ? " primary" : ""}`}
            >
              <div
                className="hb-stone-dot"
                style={{ background: action.tone }}
                aria-hidden="true"
              />
              <div>
                <h3>{action.label}</h3>
                <p>{action.desc}</p>
              </div>
              <Link to={action.to} className="hb-btn" style={{ display: "grid", placeItems: "center", textDecoration: "none" }}>開啟</Link>
            </article>
          ))}
        </section>

        <section className="hb-collection-panel" aria-label="收藏與查詢">
          <div className="hb-section-head">
            <p className="hb-eyebrow">COLLECT &amp; CHECK</p>
            <h2>收藏、背包、合成、排行</h2>
          </div>
          <div className="hb-collection-grid">
            {COLLECTIONS.map((item) => (
              <Link key={item.label} to={item.to} className="hb-collection-tile" style={{ textDecoration: "none" }}>
                <span
                  className="hb-tile-mark"
                  style={{ background: item.tone }}
                  aria-hidden="true"
                />
                <strong>{item.label}</strong>
                <small>{item.count}</small>
              </Link>
            ))}
          </div>
        </section>

        <section className="hb-base-snapshot" aria-label="基地展示摘要">
          <div className="hb-mini-map" aria-hidden="true">
            <span /><span /><span />
          </div>
          <div>
            <h2>目前基地：營燈前哨</h2>
            <p>基地展示只做狀態摘要，不搶主要操作位置。</p>
          </div>
        </section>
      </section>
    </main>
  )
}
