---
title: bbs-go 项目业务全景与工程认知
last_verified: 2026-08-08
---

<!-- __SYSMAP_INDEX__ -->
## 文档索引（index_updated: 2026-08-08）
| 行号 | 主题 |
|------|------|
| L19-L27 |   一、项目定位与仓库关系 |
| L28-L36 |   二、技术栈（重要认知，与旧版文档不同） |
| L37-L46 |   三、CI/CD 与部署模型（2026-08-08 起） |
| L47-L48 |   四、业务规则与已修复问题（本 fork 特有） |
| L49-L60 |     4.1 用户发帖/评论计数一致性（PR #297） |
| L61-L67 |     4.2 评论输入框 Firefox 显示 bug（PR #301） |
| L68-L73 |   五、开发注意（本仓库约定） |
<!-- __SYSMAP_INDEX_END__ -->

## 一、项目定位与仓库关系

- `bbs-go` 是一个轻量级社区/问答平台（Go 后端 + Web 前端），上游主仓库为 `mlogclub/bbs-go`。
- **本地协作模型为 fork 工作流**：
  - `origin` = 上游 `https://github.com/mlogclub/bbs-go.git`（只读，PR 目标）
  - `fork` = 自有仓库 `https://github.com/aishangwuji/bbs-go.git`（日常开发与 CI/CD 触发源）
- **生产环境只从 fork 拉取部署**，不直接使用上游镜像。
- 服务器上的源码检出：`/opt/bbs-go/repo/`（remote 指向 fork）。

## 二、技术栈（重要认知，与旧版文档不同）

- 后端：Go（`cmd/`、`internal/`，gin + gorm），SQLite 为主。
- 前端：**`web/` 目录为 React Router v7 + Vite + shadcn/ui + Tailwind v4**（`pnpm` 管理，2026 年由 Nuxt 重构而来）。
  - 组件库：`web/components/ui/`（shadcn）、`web/components/comment/`（评论）、`web/components/editor/`（富文本）等。
  - 评论输入框组件：`web/components/comment/text-editor.tsx`。
  - 样式：Tailwind v4 + `web/styles/`（含 editor.css，注意其中 `.simple-editor` 用于发帖/文章编辑器，评论框不用它）。
- 变更前端代码后必须跑：`pnpm typecheck` 与 `pnpm lint`（在 `web/` 下）。

## 三、CI/CD 与部署模型（2026-08-08 起）

- **构建在 GitHub Actions 云端完成，服务器不再本地编译**（2核4G 撑不住 `docker build`）。
- 触发：push 到 fork 的 `master` 分支 → `.github/workflows/docker-image.yml` → 多阶段 Dockerfile 构建 → 推送 GHCR。
- 镜像：`ghcr.io/aishangwuji/bbs-go:latest` + `sha-<commit>`（仓库为 PUBLIC，服务器可匿名拉取）。
- 服务器部署：`bash /opt/bbs-go/repo/deploy/remote-deploy.sh`（pull → 校验 compose → 重建 → 健康检查）。
- 回滚：`docker pull ghcr.io/aishangwuji/bbs-go:sha-<commit>` + `docker compose -f /opt/bbs-go/docker-compose.yml up -d --no-deps`。
- 服务器架构细节（nginx SNI 分流、端口、证书等）：见服务器 `/opt/server-architecture.md`（每次改动服务器服务必须同步更新）。
- 完整流程说明已写入仓库 `README.md` / `README.en-US.md`。

## 四、业务规则与已修复问题（本 fork 特有）

### 4.1 用户发帖/评论计数一致性（PR #297）

- 业务规则：`t_user.topic_count` / `comment_count` 必须与「已发布（status=0）」的帖子/评论数量一致。
- 历史问题：删除/待审核未正确扣减或计入，导致积分排行与角色框数量不符。
- 修复策略（运行时增减计数，已合入 fork master）：
  - 待审核帖子发布时**不计**入 `topic_count`，审核通过（`TopicService.Audit`）后 +1。
  - 删除已发布帖子 -1（仅当原状态为已发布）；恢复（Undelete）+1。
  - 评论删除 -1；发布 +1。
  - 计数增减用 `CASE WHEN n > 0 THEN n - 1 ELSE 0 END` 防负。
- **审阅人意见**：不要用启动 migration 全表重算（大数据量会卡死启动），历史数据手动在数据库执行。因此 `migrations/000016_*` 相关文件已从 PR 移除。
- 对应测试：`internal/services/topic_count_test.go`。

### 4.2 评论输入框 Firefox 显示 bug（PR #301）

- 现象：帖子详情页评论输入框 placeholder 遮挡「发表」与上传按钮，**仅 Firefox**。
- 根因：`text-editor.tsx` 固定高度 flex 容器中，Firefox 的 `<textarea>` `min-height:auto` 按固有高度计算，无法收缩，顶出工具条。
- 修复：textarea 加 `min-h-0`；工具条与图片区加 `shrink-0`。
- 上线验证方式：检查 CSS 产物中 `.min-h-0` 规则是否生成（`min-height:calc(var(--spacing) * 0)`）。

## 五、开发注意（本仓库约定）

- 变更后端计数逻辑时注意并发/事务（`sqls.WithTransaction`）。
- 前端尽量用既有组件库，避免引入新库造成认知负担。
- 提交 PR 到上游时，只含业务修复本身，不含本 fork 私有的 CI/CD、部署脚本等改动。
- 每次首次连接服务器先读 `/opt/server-architecture.md`；改动服务器服务后必须同步更新该文档。
