import assert from "node:assert/strict"

import {
  getEditorModeOptions,
  getEditorSwitchTarget,
  getEditorSwitchConfirmMessage,
} from "../lib/editor-mode.ts"

const zhMessages = {
  "component.editorMode.visual": "可视化编辑",
  "component.editorMode.markdown": "Markdown",
  "component.editorMode.switchConfirm":
    "切换到{mode}会清空当前内容，是否继续？",
}

const enMessages = {
  "component.editorMode.visual": "Visual",
  "component.editorMode.markdown": "Markdown",
  "component.editorMode.switchConfirm":
    "Switching to {mode} will clear the current content. Continue?",
}

function createFakeT(messages) {
  return (key, params = {}) => {
    let value = messages[key] || key
    for (const [name, paramValue] of Object.entries(params)) {
      value = value.replace(`{${name}}`, String(paramValue))
    }
    return value
  }
}

const zhT = createFakeT(zhMessages)
const enT = createFakeT(enMessages)

assert.deepEqual(getEditorModeOptions(zhT), [
  { value: "html", label: "可视化编辑" },
  { value: "markdown", label: "Markdown" },
])

assert.equal(getEditorSwitchTarget("html"), "markdown")
assert.equal(getEditorSwitchTarget("markdown"), "html")

assert.equal(
  getEditorSwitchConfirmMessage("html", zhT),
  "切换到Markdown会清空当前内容，是否继续？"
)
assert.equal(
  getEditorSwitchConfirmMessage("markdown", zhT),
  "切换到可视化编辑会清空当前内容，是否继续？"
)
assert.equal(
  getEditorSwitchConfirmMessage("html", enT),
  "Switching to Markdown will clear the current content. Continue?"
)

console.log("editor mode tests passed")
