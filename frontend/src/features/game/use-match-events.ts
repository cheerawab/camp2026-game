import { useQueryClient } from "@tanstack/react-query"
import { useEffect } from "react"

import { MatchStateSchema, type MatchState } from "@/shared/api/game"

const matchEventNames = [
  "match_updated",
  "player_ready",
  "player_answered",
  "round_started",
  "round_revealed",
  "match_completed",
] as const

export function useMatchEvents(
  matchID: string,
  options: { enabled?: boolean; onDeleted?: () => void } = {},
) {
  const queryClient = useQueryClient()
  const enabled = options.enabled ?? true
  const onDeleted = options.onDeleted

  useEffect(() => {
    if (!enabled || !matchID || typeof window === "undefined") return

    const source = new EventSource(
      `/api/matches/${encodeURIComponent(matchID)}/events`,
      { withCredentials: true },
    )
    const handleMessage = (event: MessageEvent<string>) => {
      try {
        const match = MatchStateSchema.parse(JSON.parse(event.data))
        queryClient.setQueryData<MatchState>(["matches", matchID], match)
      } catch {
        // Ignore malformed events and let the connection keep streaming.
      }
    }
    const handleDeleted = () => {
      queryClient.removeQueries({ queryKey: ["matches", matchID] })
      onDeleted?.()
      source.close()
    }

    for (const eventName of matchEventNames) {
      source.addEventListener(eventName, handleMessage)
    }
    source.addEventListener("match_deleted", handleDeleted)

    return () => source.close()
  }, [enabled, matchID, onDeleted, queryClient])
}

export function useMatchDeadlineRefresh(
  matchID: string,
  match: MatchState | undefined,
) {
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!matchID || match?.status !== "active") return

    const deadline =
      match.phase === "revealing" ? match.revealEndsAt : match.roundEndsAt
    if (!deadline) return

    const deadlineAt = new Date(deadline).getTime()
    if (!Number.isFinite(deadlineAt)) return

    const delay = Math.max(0, deadlineAt - Date.now() + 300)
    const timeout = window.setTimeout(() => {
      void queryClient.invalidateQueries({ queryKey: ["matches", matchID] })
    }, delay)

    return () => window.clearTimeout(timeout)
  }, [
    matchID,
    match?.phase,
    match?.revealEndsAt,
    match?.roundEndsAt,
    match?.status,
    queryClient,
  ])
}
