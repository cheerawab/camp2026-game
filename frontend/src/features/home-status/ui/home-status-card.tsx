import { CalendarClock, RefreshCw, Server, Wifi } from "lucide-react"

import type { HomeStatus } from "../model/health.schema"
import { Button } from "@/shared/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card"
import { InfoTile } from "@/shared/ui/info-tile"
import { StatusBadge } from "@/shared/ui/status-badge"

type HomeStatusCardProps = {
  data: HomeStatus
  isFetching?: boolean
  onRefresh?: () => void
}

export function HomeStatusCard({
  data,
  isFetching = false,
  onRefresh,
}: HomeStatusCardProps) {
  return (
    <Card>
      <CardHeader className="gap-3">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Wifi className="text-primary size-5" />
              基地連線
            </CardTitle>
            <CardDescription>後端 API 已連線</CardDescription>
          </div>
          <StatusBadge tone="success">{data.status}</StatusBadge>
        </div>
      </CardHeader>
      <CardContent className="grid gap-3 sm:grid-cols-2">
        <InfoTile label="服務" value={data.service} icon={Server} />
        <InfoTile
          label="最後同步"
          icon={CalendarClock}
          value={
            <time>
              {new Intl.DateTimeFormat("zh-TW", {
                dateStyle: "medium",
                timeStyle: "medium",
              }).format(new Date(data.checkedAt))}
            </time>
          }
        />
      </CardContent>
      <CardFooter>
        <Button type="button" onClick={onRefresh} disabled={isFetching}>
          <RefreshCw className={`size-4 ${isFetching ? "animate-spin" : ""}`} />
          {isFetching ? "同步中" : "重新同步"}
        </Button>
      </CardFooter>
    </Card>
  )
}
