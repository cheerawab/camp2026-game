import { type FormEvent, useEffect, useRef, useState } from "react"

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

type PlayerQrScannerDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onToken: (token: string) => void
}

function getBarcodeDetector() {
  return (
    window as typeof window & {
      BarcodeDetector?: BarcodeDetectorConstructor
    }
  ).BarcodeDetector
}

function normalizeToken(value: string) {
  return value.trim()
}

export function PlayerQrScannerDialog({
  open,
  onOpenChange,
  onToken,
}: PlayerQrScannerDialogProps) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [manualToken, setManualToken] = useState("")
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
        setMessage("這個瀏覽器無法開啟相機，請輸入 QR 識別碼。")
        return
      }

      const BarcodeDetector = getBarcodeDetector()
      if (!BarcodeDetector) {
        setMessage("這個瀏覽器不支援 QR 掃描，請輸入 QR 識別碼。")
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
        setMessage("對準學員的個人 QR Code")

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
            const token = normalizeToken(detected[0]?.rawValue ?? "")
            if (token) {
              onToken(token)
              onOpenChange(false)
              return
            }
          }

          timer = window.setTimeout(scan, 350)
        }

        await scan()
      } catch {
        if (!cancelled) {
          setMessage("相機權限未開啟，請輸入 QR 識別碼。")
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
  }, [onOpenChange, onToken, open])

  function handleManualSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const token = normalizeToken(manualToken)
    if (!token) {
      setMessage("請輸入 QR 識別碼。")
      return
    }
    onToken(token)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-4">
        <DialogHeader>
          <DialogTitle>掃描學員 QR Code</DialogTitle>
          <DialogDescription>{message}</DialogDescription>
        </DialogHeader>

        <div className="bg-ink border-ink aspect-square overflow-hidden rounded-[18px] border-2">
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
            value={manualToken}
            onChange={(event) => setManualToken(event.target.value)}
            placeholder="手動輸入 QR 識別碼"
            autoComplete="off"
            inputMode="text"
          />
          <DialogFooter>
            <Button type="submit" variant="secondary" className="w-full">
              確認學員
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
