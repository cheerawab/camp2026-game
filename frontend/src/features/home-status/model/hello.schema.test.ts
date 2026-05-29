import { describe, expect, it } from "vitest"

import { HelloResponseSchema } from "./hello.schema"

describe("HelloResponseSchema", () => {
  it("accepts the mock API contract", () => {
    const parsed = HelloResponseSchema.parse({
      message: "Hello from wired mock API",
      service: "camp2026-game-frontend",
      status: "ok",
      generatedAt: "2026-05-29T00:00:00.000Z",
    })

    expect(parsed.status).toBe("ok")
  })

  it("rejects invalid status values", () => {
    expect(() =>
      HelloResponseSchema.parse({
        message: "Hello",
        service: "camp2026-game-frontend",
        status: "degraded",
        generatedAt: "2026-05-29T00:00:00.000Z",
      }),
    ).toThrow()
  })
})
