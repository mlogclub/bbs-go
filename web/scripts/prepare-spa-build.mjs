import { cp, rm } from "node:fs/promises"
import path from "node:path"
import { fileURLToPath } from "node:url"

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..")
const source = path.join(root, "build/client")
const target = path.join(root, "build/spa")

await rm(target, { recursive: true, force: true })
await cp(source, target, { recursive: true })
