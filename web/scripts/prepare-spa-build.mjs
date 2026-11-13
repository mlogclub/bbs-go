import { cp, rm, readdir, writeFile } from "node:fs/promises"
import path from "node:path"
import { fileURLToPath } from "node:url"

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..")
const source = path.join(root, "build/client")
const target = path.join(root, "build/spa")

await rm(target, { recursive: true, force: true })
await cp(source, target, { recursive: true })

// Keep the Go SPA fallback aligned with the actual file routes. Resource existence
// still requires SSR; this manifest distinguishes unknown URL shapes.
const routes = (await readdir(path.join(root, "app/routes")))
  .filter((file) => file.endsWith(".tsx") && file !== "$.tsx")
  .map((file) => {
    const parts = file
      .slice(0, -4)
      .split(".")
      .filter((part) => part !== "_index")
    const segments = parts.map((part) => {
      if (part === "$") return ".*"
      if (part.startsWith("$")) return "[^/]+"
      return part.replace(/_$/, "").replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
    })
    return "^/" + segments.join("/") + "/?$"
  })
await writeFile(
  path.join(target, ".route-patterns.json"),
  JSON.stringify(routes)
)
