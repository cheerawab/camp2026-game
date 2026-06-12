import { useMutation } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { ArrowRight, Play, ScanQrCode } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

import {
  MatchCodeScannerDialog,
  normalizeMatchCode,
} from "@/features/battle-qr"
import { gameApi, type MatchState } from "@/shared/api/game"
import { Button } from "@/shared/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card"
import { Field } from "@/shared/ui/field"
import { GamePageShell } from "@/shared/ui/game-page-shell"
import { Input } from "@/shared/ui/input"
import { PageHeader } from "@/shared/ui/page-header"

function storeMatch(match: MatchState) {
  if (typeof window !== "undefined") {
    window.localStorage.setItem("camp2026.currentMatchId", match.matchId)
  }
}

export function BattleLobbyPage() {
  const navigate = useNavigate()
  const [code, setCode] = useState("")
  const [scannerOpen, setScannerOpen] = useState(false)
  const onMatchReady = (match: MatchState) => {
    storeMatch(match)
    navigate({ to: "/battle/room" })
  }
  const createMutation = useMutation({
    mutationFn: gameApi.createMatch,
    onSuccess: onMatchReady,
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : "建立房間失敗")
    },
  })
  const joinMutation = useMutation({
    mutationFn: gameApi.joinMatch,
    onSuccess: onMatchReady,
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : "加入房間失敗")
    },
  })

  function handleJoinCode(value: string) {
    const normalizedCode = normalizeMatchCode(value)
    if (!normalizedCode) {
      toast.error("請先輸入房號")
      return
    }
    joinMutation.mutate(normalizedCode)
  }

  function handleJoin() {
    handleJoinCode(code)
  }

  return (
    <GamePageShell contentClassName="grid content-start gap-y-2">
      <PageHeader title="知識王" headline="Battle Lobby" />
      <Card>
        <CardHeader>
          <CardTitle>快速開始</CardTitle>
          <CardDescription>建立一個雙人知識王房間</CardDescription>
        </CardHeader>
        <CardContent>
          <span>建立房間後，把房號分享給另一位學員加入對戰。</span>
        </CardContent>
        <CardFooter>
          <Button
            className="w-full"
            disabled={createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            <Play />
            {createMutation.isPending ? "建立中" : "建立房間"}
          </Button>
        </CardFooter>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>多人連線</CardTitle>
          <CardDescription>使用房號加入等待中的對戰</CardDescription>
        </CardHeader>
        <CardContent>
          <span>和其他學員連線對戰，比拼誰才是知識王。</span>
        </CardContent>
        <CardFooter className="grid gap-2">
          <Field orientation="horizontal">
            <Input
              id="input-room-id"
              type="text"
              value={code}
              onChange={(event) =>
                setCode(normalizeMatchCode(event.target.value))
              }
              placeholder="請輸入房號"
            />
            <Button
              size="icon-lg"
              type="button"
              aria-label="掃描房號 QR Code"
              disabled={joinMutation.isPending}
              onClick={() => setScannerOpen(true)}
            >
              <ScanQrCode />
            </Button>
            <Button
              className="w-full flex-1"
              size="lg"
              variant="secondary"
              disabled={joinMutation.isPending}
              onClick={handleJoin}
            >
              {joinMutation.isPending ? "加入中" : "加入房間"}
              <ArrowRight />
            </Button>
          </Field>
        </CardFooter>
      </Card>
      {scannerOpen ? (
        <MatchCodeScannerDialog
          open={scannerOpen}
          onOpenChange={setScannerOpen}
          onCode={(scannedCode) => {
            setCode(scannedCode)
            handleJoinCode(scannedCode)
          }}
        />
      ) : null}
    </GamePageShell>
  )
}
