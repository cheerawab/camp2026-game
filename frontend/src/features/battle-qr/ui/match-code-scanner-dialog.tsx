import { useEffect, useRef, useState } from "react"

import { normalizeMatchCode } from "@/features/battle-qr/lib/match-code"
import { Button } from "@/shared/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/dialog"
import { Input } from "@/shared/ui/input"

type BarcodeDetectorResult = {
  rawValue?: string
}

type BarcodeDetectorInstance = {
  detect(source: HTMLCanvasElement): Promise<BarcodeDetectorResult[]>
}

type BarcodeDetectorConstructor = new (options: {
  formats: string[]
}) => BarcodeDetectorInstance

type MatchCodeScannerDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCode: (code: string) => void
}

function getBarcodeDetector() {
  return (
    window as typeof window & {
      BarcodeDetector?: BarcodeDetectorConstructor
    }
  ).BarcodeDetector
}

export function MatchCodeScannerDialog({
  open,
  onOpenChange,
  onCode,
}: MatchCodeScannerDialogProps) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [manualCode, setManualCode] = useState("")
  const [message, setMessage] = useState("正在啟動相機")

  useEffect(() => {
    if (!open) return

    let cancelled = false
    let timer = 0
    let stream: MediaStream | null = null
    let activeVideo: HTMLVideoElement | null = null

    async function startScanner() {
      if (
        typeof window === "undefined" ||
        !navigator.mediaDevices?.getUserMedia
      ) {
        setMessage("這個瀏覽器無法開啟相機，請輸入房號。")
        return
      }

      const BarcodeDetector = getBarcodeDetector()
      if (!BarcodeDetector) {
        setMessage("這個瀏覽器不支援 QR 掃描，請輸入房號。")
        return
      }

      try {
        const detector = new BarcodeDetector({ formats: ["qr_code"] })
        stream = await navigator.mediaDevices.getUserMedia({
          audio: false,
          video: { facingMode: { ideal: "environment" } },
        })
        if (cancelled) return

        const video = videoRef.current
        if (!video) return

        activeVideo = video
        video.srcObject = stream
        await video.play()
        setMessage("對準房號 QR Code")

        const scan = async () => {
          if (cancelled) return

          const canvas = canvasRef.current
          const context = canvas?.getContext("2d", {
            willReadFrequently: true,
          })
          if (
            canvas &&
            context &&
            video.readyState >= HTMLMediaElement.HAVE_CURRENT_DATA &&
            video.videoWidth > 0 &&
            video.videoHeight > 0
          ) {
            canvas.width = video.videoWidth
            canvas.height = video.videoHeight
            context.drawImage(video, 0, 0, canvas.width, canvas.height)

            const detected = await detector.detect(canvas)
            const code = normalizeMatchCode(detected[0]?.rawValue ?? "")
            if (code) {
              onCode(code)
              onOpenChange(false)
              return
            }
          }

          timer = window.setTimeout(scan, 350)
        }

        await scan()
      } catch {
        if (!cancelled) {
          setMessage("相機權限未開啟，請輸入房號。")
        }
      }
    }

    void startScanner()

    return () => {
      cancelled = true
      window.clearTimeout(timer)
      stream?.getTracks().forEach((track) => track.stop())
      if (activeVideo) {
        activeVideo.srcObject = null
      }
    }
  }, [onCode, onOpenChange, open])

  function handleManualSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const code = normalizeMatchCode(manualCode)
    if (!code) {
      setMessage("請輸入房號。")
      return
    }
    onCode(code)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-4">
        <DialogHeader>
          <DialogTitle>掃描房號 QR Code</DialogTitle>
          <DialogDescription>{message}</DialogDescription>
        </DialogHeader>

        <div className="border-ink bg-ink aspect-square overflow-hidden rounded-[18px] border-2">
          <video
            ref={videoRef}
            className="h-full w-full object-cover"
            playsInline
            muted
          />
          <canvas ref={canvasRef} className="hidden" />
        </div>

        <form className="grid gap-3" onSubmit={handleManualSubmit}>
          <Input
            value={manualCode}
            onChange={(event) =>
              setManualCode(normalizeMatchCode(event.target.value))
            }
            placeholder="手動輸入房號"
            autoComplete="off"
            inputMode="text"
          />
          <DialogFooter>
            <Button type="submit" variant="secondary" className="w-full">
              加入房間
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
