import assert from "node:assert/strict"
import { readFileSync } from "node:fs"
import { resolve } from "node:path"

const repoRoot = resolve(import.meta.dirname, "../..")
const topicServiceSource = readFileSync(
  resolve(repoRoot, "internal/services/topic_service.go"),
  "utf8"
)

assert.match(
  topicServiceSource,
  /const topicListPageSize = 40/,
  "topic list page size should be defined as 40"
)

function getTopicServiceFunction(functionName) {
  const start = topicServiceSource.indexOf(
    `func (s *topicService) ${functionName}(`
  )
  assert.notEqual(start, -1, `TopicService.${functionName} should exist`)

  const next = topicServiceSource.indexOf("\nfunc (s *topicService)", start + 1)
  return topicServiceSource.slice(start, next === -1 ? undefined : next)
}

for (const [functionName, message] of [
  ["GetTopics", "default topic/category/recommend list page size"],
  ["_GetFollowTopics", "follow topic list page size"],
  ["GetTagTopics", "tag topic list page size"],
]) {
  const functionSource = getTopicServiceFunction(functionName)
  assert.equal(
    functionSource.includes("topicListPageSize"),
    true,
    `TopicService.${functionName} should use topicListPageSize for ${message}`
  )
  assert.doesNotMatch(
    functionSource,
    /(var limit(?: int)? =|limit :=) 20/,
    `TopicService.${functionName} should not keep a hard-coded page size of 20`
  )
}

console.log("topic list page size is covered")
