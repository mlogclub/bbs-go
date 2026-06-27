import assert from "node:assert/strict"
import { existsSync, readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const feedTabsPath = resolve(webRoot, "components/topic/topic-feed-tabs.tsx")
const topicsNavSource = readFileSync(
  resolve(webRoot, "components/topic/topics-nav-content.tsx"),
  "utf8"
)
const topicCss = readFileSync(resolve(webRoot, "styles/topic.css"), "utf8")
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
assert.equal(
  topicsNavSource.includes('t("pages.topics.allCategories")'),
  true,
  "TopicsNavContent should render a local all-categories entry"
)
assert.equal(
  topicsNavSource.includes('href="/topics"'),
  true,
  "The local all-categories entry should link to the topics index"
)
assert.equal(
  topicsNavSource.includes('import { LayoutGridIcon } from "lucide-react"'),
  true,
  "The local all-categories entry should use a lucide icon"
)
assert.match(
  topicsNavSource,
  /data-node-id="all"[\s\S]*?<LayoutGridIcon[\s\S]*?className="node-logo node-logo-icon"/,
  "The local all-categories entry should render the overview icon"
)
assert.match(
  topicCss,
  /\.node-logo-icon\s*\{/,
  "lucide nav icons should be styled in topic.css"
)
assert.equal(
  topicsNavSource.includes('"categories-divider"'),
  true,
  "The all-categories entry should be separated from real categories"
)
assert.match(
  topicCss,
  /\.categories-divider\s*\{/,
  "categories divider should be styled in topic.css"
)
assert.doesNotMatch(
  topicCss,
  /\.categories-divider\s*\{[\s\S]*?margin:\s*8px\s+24px/,
  "categories divider should span the full nav width without horizontal margin"
)
assert.match(
  topicCss,
  /\.topics-wrapper\s+\.topics-nav\s*\{[\s\S]*?position:\s*sticky[\s\S]*?top:\s*calc\(52px \+ 1rem\)[\s\S]*?align-self:\s*flex-start/,
  "topics nav wrapper should own sticky positioning so it is not constrained by a short child container"
)
assert.doesNotMatch(
  topicCss,
  /\.dock-nav\s*\{[\s\S]*?position:\s*sticky/,
  "dock nav should not own sticky positioning inside the short wrapper"
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
  indexRouteSource.includes(
    "<TopicsNavContent initialCategories={categories} currentCategoryId={0} />"
  ),
  true,
  "Home topic list should mark all categories selected in the left nav"
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
