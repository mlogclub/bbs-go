import path from "node:path"
import { reactRouter } from "@react-router/dev/vite"
import { defineConfig, type Plugin } from "vite"

function stripSpaRouteLoaders(): Plugin {
  const appRoutesDir = `${path.sep}app${path.sep}routes${path.sep}`

  return {
    name: "bbsgo-strip-spa-route-loaders",
    enforce: "pre",
    transform(code, id) {
      if (
        process.env.BBSGO_WEB_SPA !== "true" ||
        !id.includes(appRoutesDir)
      ) {
        return null
      }

      const nextCode = code
        .replace(
          /^\s*export\s*\{\s*loader\s*\}\s*from\s*["'][^"']+["'];?\s*$/gm,
          ""
        )
        .replace(/\bexport\s+(async\s+function\s+loader\b)/g, "$1")
        .replace(/\bexport\s+(function\s+loader\b)/g, "$1")
        .replace(/\bexport\s+(const\s+loader\s*=)/g, "$1")

      return nextCode === code ? null : { code: nextCode, map: null }
    },
  }
}

export default defineConfig({
  plugins: [stripSpaRouteLoaders(), reactRouter()],
  optimizeDeps: {
    include: ["md-editor-rt"],
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname),
      "tailwindcss": path.resolve(__dirname, "node_modules/tailwindcss/index.css"),
    },
  },
  server: {
    port: 3000,
    proxy: {
      "/api": "http://localhost:8082",
      "/res": "http://localhost:8082",
    },
  },
})
