import { z } from "zod"

const EnvSchema = z.object({
  VITE_APP_NAME: z.string().min(1).default("Camp 2026 Game"),
  VITE_APP_ORIGIN: z.url().optional(),
})

export const env = EnvSchema.parse(import.meta.env)
