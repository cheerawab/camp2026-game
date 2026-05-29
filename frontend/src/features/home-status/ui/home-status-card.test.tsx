import { fireEvent, render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { HomeStatusCard } from "./home-status-card"

const apiResponse = {
  message: "Hello from wired mock API",
  service: "camp2026-game-frontend",
  status: "ok" as const,
  generatedAt: "2026-05-29T00:00:00.000Z",
}

describe("HomeStatusCard", () => {
  it("renders the parsed API payload", () => {
    render(<HomeStatusCard data={apiResponse} />)

    expect(screen.getByText("基地連線")).toBeInTheDocument()
    expect(screen.getByText("Hello from wired mock API")).toBeInTheDocument()
    expect(screen.getByText("camp2026-game-frontend")).toBeInTheDocument()
    expect(screen.getByText("ok")).toBeInTheDocument()
  })

  it("calls refresh when the action is clicked", () => {
    const onRefresh = vi.fn()

    render(<HomeStatusCard data={apiResponse} onRefresh={onRefresh} />)
    fireEvent.click(screen.getByRole("button", { name: "重新同步" }))

    expect(onRefresh).toHaveBeenCalledOnce()
  })
})
