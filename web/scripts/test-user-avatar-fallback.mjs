import fs from "node:fs"
import path from "node:path"
import { fileURLToPath } from "node:url"

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const source = fs.readFileSync(
  path.join(__dirname, "../components/common/avatar.tsx"),
  "utf8",
)

const expectations = [
  {
    name: "tracks avatar image load failures",
    pattern: /useState\s*\(\s*false\s*\)/,
  },
  {
    name: "switches fallback on image error",
    pattern: /onError=\{[^}]*set[A-Za-z0-9_]*Failed\s*\(\s*true\s*\)/s,
  },
  {
    name: "resets failed state when avatar source changes",
    pattern: /React\.useEffect\s*\(\s*\(\)\s*=>\s*\{\s*set[A-Za-z0-9_]*Failed\s*\(\s*false\s*\)/s,
  },
  {
    name: "only renders the image while the source has not failed",
    pattern: /src\s*&&\s*![A-Za-z0-9_]*Failed/,
  },
]

const failures = expectations.filter(({ pattern }) => !pattern.test(source))

if (failures.length > 0) {
  console.error("UserAvatar fallback behavior is incomplete:")
  for (const failure of failures) {
    console.error(`- ${failure.name}`)
  }
  process.exit(1)
}

console.log("UserAvatar fallback behavior is covered.")
