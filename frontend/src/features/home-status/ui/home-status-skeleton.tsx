import { Card, CardContent, CardFooter, CardHeader } from "@/shared/ui/card"
import { RefreshCw } from "lucide-react"
import { Skeleton } from "@/shared/ui/skeleton"

export function HomeStatusSkeleton() {
  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-5 w-28" />
        <Skeleton className="h-4 w-52" />
      </CardHeader>
      <CardContent className="grid gap-3 sm:grid-cols-2">
        <Skeleton className="h-20" />
        <Skeleton className="h-20" />
      </CardContent>
      <CardFooter>
        <div className="bg-muted text-muted-foreground flex h-9 w-28 items-center justify-center gap-2 rounded-md">
          <RefreshCw className="size-4 animate-spin" />
          <span className="text-sm">同步中</span>
        </div>
      </CardFooter>
    </Card>
  )
}
