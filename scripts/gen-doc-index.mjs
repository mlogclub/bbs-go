#!/usr/bin/env node
/**
 * gen-doc-index.mjs — 生成/更新 systemmap 文档的「文档索引」段（规则39 渐进式暴露协议）
 *
 * 用法：
 *   node scripts/gen-doc-index.mjs             # 处理 systemmap/*.md
 *   node scripts/gen-doc-index.mjs <file...>   # 处理指定文件（支持通配）
 *
 * 机制：
 *   - 幂等：识别 `<!-- __SYSMAP_INDEX__ -->` … `<!-- __SYSMAP_INDEX_END__ -->` 块，
 *     存在则替换，不存在则插入到 frontmatter 之后。
 *   - 索引行号 = 生成后文件的真实行号（脚本以「frontmatter + 索引段 + 正文」重写文件，
 *     正文逐字节保留，行号天然正确，无偏移计算）。
 *   - 跳过错代码块（``` 围栏）内的标题行，避免误索引。
 *
 * 规则39.3（防行号漂移）：任何对已索引文档的增删改，必须重新运行本脚本，
 * 或在同一改动中手动同步索引段行号与 index_updated。
 */

import { readFileSync, writeFileSync, readdirSync } from 'node:fs'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const INDEX_START = '<!-- __SYSMAP_INDEX__ -->'
const INDEX_END = '<!-- __SYSMAP_INDEX_END__ -->'

const scriptDir = dirname(fileURLToPath(import.meta.url))
const defaultDir = resolve(scriptDir, '../systemmap')

// ---------- 工具 ----------

/** 生成索引段的完整行数组。行号基于最终文件布局：head + 空行 + 索引段 + 空行 + body */
function buildIndexBlock(headings, headLines, bodyLength) {
  const n = headings.length
  // 块内固定行：START 标记(1) + "## 文档索引"(1) + 表头(1) + 分隔行(1) + 表体(n) + END 标记(1)
  const blockLines = 5 + n
  // body 首行在最终文件中的行号 = headLines + 1(空行) + blockLines + 1(空行) + 1
  const bodyBase = headLines + blockLines + 3

  const date = new Date()
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')

  const rows = headings.map((h, i) => {
    const start = bodyBase + h.bodyLine - 1
    const end = i < n - 1
      ? bodyBase + headings[i + 1].bodyLine - 2
      : bodyBase + bodyLength - 1
    const indent = '  '.repeat(Math.min(h.level - 1, 3))
    const title = h.title.replace(/\|/g, '\\|').trim()
    return `| L${start}-L${end} | ${indent}${title} |`
  })

  return [
    INDEX_START,
    `## 文档索引（index_updated: ${y}-${m}-${d}）`,
    '| 行号 | 主题 |',
    '|------|------|',
    ...rows,
    INDEX_END,
  ]
}

/** 移除旧的索引段（含周边多余空行），返回新 body 行数组 */
function stripIndex(body) {
  const startIdx = body.findIndex((l) => l.includes(INDEX_START))
  if (startIdx === -1) return body
  const endIdx = body.findIndex((l, i) => i > startIdx && l.includes(INDEX_END))
  const sliceEnd = endIdx === -1 ? startIdx + 1 : endIdx + 1
  let cleaned = [...body.slice(0, startIdx), ...body.slice(sliceEnd)]
  // 压缩连续空行（最多留 1 行，由拼装阶段统一补分隔空行）
  while (cleaned.length && cleaned[0].trim() === '') cleaned.shift()
  while (cleaned.length && cleaned[cleaned.length - 1].trim() === '') cleaned.pop()
  return cleaned
}

/** 扫描正文标题（跳过 fenced 代码块），返回 [{level, title, bodyLine}] */
function scanHeadings(body) {
  const headings = []
  let fenced = false
  body.forEach((line, i) => {
    if (/^\s*```/.test(line)) {
      fenced = !fenced
      return
    }
    if (fenced) return
    if (line.startsWith(INDEX_START) || line.startsWith(INDEX_END)) return
    const m = line.match(/^(#{1,4})\s+(.+)$/)
    if (m) headings.push({ level: m[1].length, title: m[2].replace(/\s+#*\s*$/, ''), bodyLine: i + 1 })
  })
  return headings
}

// ---------- 主流程 ----------

function processFile(filePath) {
  const raw = readFileSync(filePath, 'utf8')
  const eol = raw.includes('\r\n') ? '\r\n' : '\n'
  const lines = raw.split(/\r?\n/)

  // 1. 分离 frontmatter（--- 包裹）与正文
  let head = []
  let body = lines
  if (lines[0] === '---') {
    const endIdx = lines.indexOf('---', 1)
    if (endIdx > 0) {
      head = lines.slice(0, endIdx + 1)
      body = lines.slice(endIdx + 1)
    }
  }
  while (body.length && body[0].trim() === '') body.shift()
  while (body.length && body[body.length - 1].trim() === '') body.pop()

  // 2. 移除旧索引段
  body = stripIndex(body)

  // 3. 扫描标题并生成索引
  const headings = scanHeadings(body)
  const indexBlock = buildIndexBlock(headings, head.length, body.length)

  // 4. 拼装重写：head + 空行 + 索引段 + 空行 + body
  const output = [...head, '', ...indexBlock, '', ...body].join(eol) + eol
  writeFileSync(filePath, output)
  return { headings: headings.length, lines: output.split(/\r?\n/).length }
}

// ---------- 入口 ----------

const targets = process.argv.slice(2)
const files = targets.length > 0
  ? targets
  : readdirSync(defaultDir).filter((f) => f.endsWith('.md')).map((f) => resolve(defaultDir, f))

let ok = 0
for (const f of files) {
  const full = resolve(f)
  try {
    const { headings, lines } = processFile(full)
    console.log(`[OK] ${full}（${headings} 个标题，共 ${lines} 行）`)
    ok++
  } catch (e) {
    console.error(`[FAIL] ${full}: ${e.message}`)
  }
}
console.log(`完成：${ok}/${files.length} 份文档`)
