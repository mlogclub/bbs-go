import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const loadMoreSource = readFileSync(
  resolve(webRoot, "components/common/load-more.tsx"),
  "utf8"
)
const indexRouteSource = readFileSync(
  resolve(webRoot, "app/routes/_index.tsx"),
  "utf8"
)
const dynamicListSource = readFileSync(
  resolve(webRoot, "components/topic/topic-dynamic-list-client-page.tsx"),
  "utf8"
)

assert.equal(
  loadMoreSource.includes("autoLoadOnScroll = false"),
  true,
  "LoadMore should expose autoLoadOnScroll as an opt-in feature"
)
assert.equal(
  loadMoreSource.includes("sentinelRef"),
  true,
  "LoadMore should render a bottom sentinel for automatic loading"
)
assert.equal(
  loadMoreSource.includes("IntersectionObserver"),
  true,
  "LoadMore should observe the bottom sentinel"
)
assert.match(
  loadMoreSource,
  /if\s*\(\s*!autoLoadOnScroll\s*\|\|\s*!hasMore\s*\|\|\s*loading\s*\)/,
  "Automatic loading should be gated by the opt-in flag, hasMore, and loading"
)
assert.match(
  loadMoreSource,
  /entry\?\.isIntersecting[\s\S]*void loadMore\(\)/,
  "LoadMore should request the next page when the sentinel reaches the viewport"
)

const indexAutoLoadCount =
  indexRouteSource.match(/\bautoLoadOnScroll\b/g)?.length || 0
assert.equal(
  indexAutoLoadCount,
  1,
  "Home topic list should enable auto LoadMore on scroll"
)

const dynamicAutoLoadCount =
  dynamicListSource.match(/\bautoLoadOnScroll\b/g)?.length || 0
assert.equal(
  dynamicAutoLoadCount,
  2,
  "Tag and category topic lists should enable auto LoadMore on scroll"
)

console.log("topic auto load more on scroll is covered")
