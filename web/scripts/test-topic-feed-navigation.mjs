import assert from "node:assert/strict"
import { existsSync, readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const feedTabsPath = resolve(webRoot, "components/topic/topic-feed-tabs.tsx")
const topicsNavSource = readFileSync(
  resolve(webRoot, "components/topic/topics-nav-content.tsx"),
  "utf8"
)
assert.equal(
  existsSync(feedTabsPath),
  true,
  "TopicFeedTabs component should exist"
)
const feedTabsSource = readFileSync(feedTabsPath, "utf8")
const indexRouteSource = readFileSync(
  resolve(webRoot, "app/routes/_index.tsx"),
  "utf8"
)
const dynamicListSource = readFileSync(
  resolve(webRoot, "components/topic/topic-dynamic-list-client-page.tsx"),
  "utf8"
)

assert.match(
  topicsNavSource,
  /categories\.filter\(\(node\)\s*=>\s*node\.id\s*>\s*0\)/,
  "TopicsNavContent should hide built-in feed nodes from the left category nav"
)

for (const href of [
  "/topics/category/newest",
  "/topics/category/recommend",
  "/topics/category/feed",
]) {
  assert.equal(
    feedTabsSource.includes(`href: "${href}"`),
    true,
    `TopicFeedTabs should link to ${href}`
  )
}

assert.equal(
  indexRouteSource.includes("<TopicFeedTabs currentCategoryId={0} />"),
  true,
  "Home topic list should show the feed tabs with latest selected"
)
assert.equal(
  dynamicListSource.includes("<TopicFeedTabs currentCategoryId={categoryId} />"),
  true,
  "Dynamic topic list should show the feed tabs for built-in feeds"
)
assert.match(
  dynamicListSource,
  /categoryId\s*<=\s*0\s*\?\s*\(\s*<TopicFeedTabs/,
  "Dynamic topic list should only show feed tabs on built-in feed pages"
)

console.log("topic feed navigation is covered")
