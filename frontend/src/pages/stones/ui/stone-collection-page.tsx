import {
  StoneCollectionPanel,
  WorkshopPageShell,
} from "@/features/stone-workshop"

export function StoneCollectionPage() {
  return (
    <WorkshopPageShell eyebrow="STONE FIELD KIT" title="小石收藏">
      <StoneCollectionPanel />
    </WorkshopPageShell>
  )
}
