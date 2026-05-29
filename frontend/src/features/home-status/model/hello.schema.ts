import { z } from "zod"

export const HelloResponseSchema = z.object({
  message: z.string().min(1),
  service: z.string().min(1),
  status: z.literal("ok"),
  generatedAt: z.iso.datetime(),
})

export type HelloResponse = z.infer<typeof HelloResponseSchema>
