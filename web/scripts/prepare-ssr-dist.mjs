import { cp, rm } from "node:fs/promises"
import path from "node:path"
import { fileURLToPath } from "node:url"

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..")

await rm(path.join(root, "dist/client"), { recursive: true, force: true })
await rm(path.join(root, "dist/server"), { recursive: true, force: true })
await cp(path.join(root, "build/client"), path.join(root, "dist/client"), {
  recursive: true,
})
await cp(path.join(root, "build/server"), path.join(root, "dist/server"), {
  recursive: true,
})
