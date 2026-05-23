import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const componentSource = readFileSync(
  resolve(webRoot, "components/topic/category-selector.tsx"),
  "utf8"
)
const topicCss = readFileSync(resolve(webRoot, "styles/topic.css"), "utf8")

const requiredClasses = [
  "topic-subcategories",
  "topic-subcategories-body",
  "topic-subcategories-header",
  "topic-subcategories-label",
  "topic-subcategories-list",
]

for (const className of requiredClasses) {
  assert.equal(
    componentSource.includes(`"${className}"`),
    true,
    `CategoryQuickSelector should render ${className}`
  )
  assert.match(
    topicCss,
    new RegExp(`\\.publish-form\\s+\\.${className}\\b`),
    `${className} should be styled in topic.css`
  )
}

console.log("topic category selector class names are covered")
