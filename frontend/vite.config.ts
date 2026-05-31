import tailwindcss from "@tailwindcss/vite"
import { tanstackStart } from "@tanstack/react-start/plugin/vite"
import react from "@vitejs/plugin-react"
import { defineConfig, loadEnv } from "vite"

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "")
  const apiProxyTarget = env.API_PROXY_TARGET || "http://localhost:8080"

  return {
    server: {
      port: 3000,
      strictPort: true,
      proxy: {
        "/api": {
          target: apiProxyTarget,
          changeOrigin: true,
          ws: true,
        },
      },
    },
    resolve: {
      tsconfigPaths: true,
    },
    plugins: [tanstackStart(), react(), tailwindcss()],
  }
})
