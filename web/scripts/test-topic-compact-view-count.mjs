import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const topicListItemSource = readFileSync(
  resolve(webRoot, "components/topic/topic-list-item.tsx"),
  "utf8"
)

const compactBranchStart = topicListItemSource.indexOf(
  'if (resolvedVariant === "compact")'
)
assert.notEqual(
  compactBranchStart,
  -1,
  "TopicListItem should keep a compact variant branch"
)

const defaultBranchStart = topicListItemSource.indexOf(
  "\n  return (",
  compactBranchStart + 1
)
const compactBranchSource = topicListItemSource.slice(
  compactBranchStart,
  defaultBranchStart === -1 ? undefined : defaultBranchStart
)

assert.equal(
  topicListItemSource.includes("function formatCompactTopicViewCount("),
  true,
  "TopicListItem should format compact view count with a named helper"
)
assert.match(
  topicListItemSource,
  /viewCount\s*>\s*9999\s*\?\s*"9999\+"/,
  "Compact topic view count should display 9999+ above 9999"
)
assert.equal(
  compactBranchSource.includes("formatCompactTopicViewCount(topic.viewCount)"),
  true,
  "Compact topic list should display the topic view count"
)
assert.equal(
  compactBranchSource.includes("topic.commentCount"),
  false,
  "Compact topic list should not display comment count in the right badge"
)
assert.equal(
  compactBranchSource.includes('aria-label="views"'),
  true,
  "Compact topic list count link should be labelled as views"
)
assert.equal(
  compactBranchSource.includes("grid-cols-[minmax(0,1fr)_auto]"),
  true,
  "Compact topic list should let the right view-count column fit its content"
)
assert.equal(
  compactBranchSource.includes("grid-cols-[minmax(0,1fr)_2.5rem]"),
  false,
  "Compact topic list should not keep a fixed-width right column"
)
assert.match(
  compactBranchSource,
  /<span[\s\S]*aria-label="views"[\s\S]*formatCompactTopicViewCount\(topic\.viewCount\)/,
  "Compact topic list view count should render as a static span"
)
assert.doesNotMatch(
  compactBranchSource,
  /<Link\b(?:(?!<\/Link>).)*className="inline-flex h-7 min-w-9(?:(?!<\/Link>).)*aria-label="views"/s,
  "Compact topic list view count should not be clickable"
)

console.log("topic compact view count is covered")
