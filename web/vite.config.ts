import path from "node:path"
import { reactRouter } from "@react-router/dev/vite"
import { defineConfig, loadEnv, type Plugin } from "vite"

function stripSpaRouteLoaders(): Plugin {
  const appRoutesDir = `${path.sep}app${path.sep}routes${path.sep}`
  const spaMode = process.env.BBSGO_WEB_SPA === "true"

  return {
    name: "bbsgo-strip-spa-route-loaders",
    enforce: "pre",
    transform(code, id) {
      if (!spaMode || !id.includes(appRoutesDir)) {
        return null
      }

      let changed = false

      // 1. export { loader } from "..." -> export { }
      let nextCode = code.replace(
        /^\s*export\s*\{\s*loader\s*\}\s*from\s*(["'][^"']+["']);?\s*$/gm,
        (_, src) => { changed = true; return `export {} from ${src};` }
      )

      // 2. export { loader, X } from "..." -> export { X } from "..."
      nextCode = nextCode.replace(
        /export\s*\{\s*loader\s*,\s*(.+?)\s*\}\s*from\s*(['''""][^''""]+['''""]);?/g,
        (_m: string, rest: string, src: string) => { changed = true; return `export { ${rest.trim()} } from ${src};` }
      )
      nextCode = nextCode.replace(
        /export\s*\{\s*(.+?),\s*loader\s*\}\s*from\s*(['""][^'""]+['""]);?/g,
        (_m: string, rest: string, src: string) => { changed = true; return `export { ${rest.trim()} } from ${src};` }
      )

      // 3. export async function loader -> async function _loader (rename, not delete)
      nextCode = nextCode.replace(
        /export\s+async\s+function\s+loader\b/g,
        () => { changed = true; return "async function _loader" }
      )

      // 4. export function loader -> function _loader
      nextCode = nextCode.replace(
        /export\s+function\s+loader\b/g,
        () => { changed = true; return "function _loader" }
      )

      return changed ? { code: nextCode, map: null } : null
    },
  }
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, __dirname, "")
  const serverURL = env.BBSGO_SERVER_URL

  if (!serverURL) {
    throw new Error("BBSGO_SERVER_URL is required. Set it in web/.env.")
  }
  process.env.BBSGO_SERVER_URL = serverURL

  return {
    plugins: [stripSpaRouteLoaders(), reactRouter()],
    optimizeDeps: {
      include: ["md-editor-rt"],
    },
    resolve: {
      alias: {
        "@": path.resolve(__dirname),
        "tailwindcss": path.resolve(
          __dirname,
          "node_modules/tailwindcss/index.css"
        ),
      },
    },
    server: {
      port: 3000,
      proxy: {
        "/api": serverURL,
        "/res": serverURL,
        "/sitemap.xml": serverURL,
      },
    },
  }
})
