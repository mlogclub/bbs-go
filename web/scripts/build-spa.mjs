import { execSync } from "node:child_process"
import { cp, rm } from "node:fs/promises"
import path from "node:path"
import { fileURLToPath } from "node:url"

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..")

// Step 1: build with BBSGO_WEB_SPA=true (cross-platform)
process.env.BBSGO_WEB_SPA = "true"
execSync("react-router build", { cwd: root, stdio: "inherit", env: process.env })

// Step 2: copy build/client -> build/spa
const source = path.join(root, "build/client")
const target = path.join(root, "build/spa")
await rm(target, { recursive: true, force: true })
await cp(source, target, { recursive: true })
console.log("SPA build output ready at:", target)
