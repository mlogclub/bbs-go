# 热度点系统 - 前端 UI 组件

## 已创建组件

### 1. FlameLevel - 火焰等级图标

**文件**: `/workspace/web/components/topic/flame-level.tsx`

**用途**: 显示帖子的火焰等级（1-5 档）

```tsx
import { FlameLevel } from "@/components/topic/flame-level"

<FlameLevel level={3} />
// 显示：🔥🔥🔥 (橙色)

<FlameLevel level={5} className="animate-pulse" />
// 显示：🔥🔥🔥🔥🔥 (红色脉冲动画)
```

**特性**:
- 自动根据等级选择颜色（橙→红）
- Level 5 带脉冲动画
- Level 0 不显示

---

### 2. StakeButton - 质押按钮

**文件**: `/workspace/web/components/topic/stake-button.tsx`

**用途**: 主题详情页的质押操作按钮

```tsx
import { StakeButton } from "@/components/topic/stake-button"

<StakeButton 
  topicId={123} 
  currentFlameLevel={3} 
/>
```

**功能**:
- 点击弹出质押对话框
- 显示可用余额和剩余次数
- 快速选择按钮（+10, +50, 全部）
- 成功后触发 `heat-stake-updated` 事件
- 高火焰等级时显示红色警示

---

### 3. StakeRecordsDialog - 质押记录对话框

**文件**: `/workspace/web/components/topic/stake-records-dialog.tsx`

**用途**: 查看和管理质押记录

```tsx
import { StakeRecordsDialog } from "@/components/topic/stake-records-dialog"

const [open, setOpen] = useState(false)

<StakeRecordsDialog 
  open={open} 
  onOpenChange={setOpen} 
/>

// 触发打开
<Button onClick={() => setOpen(true)}>我的质押</Button>
```

**功能**:
- 统计卡片：余额、剩余次数
- 分页签：进行中 / 历史记录
- 显示火焰等级、状态
- 支持赎回操作（当日不可赎回）
- 显示历史收益

---

### 4. HeatPointsPanel - 热度点面板

**文件**: `/workspace/web/components/topic/heat-points-panel.tsx`

**用途**: 用户中心展示热度点统计

```tsx
import { HeatPointsPanel } from "@/components/topic/heat-points-panel"

<HeatPointsPanel />
```

**功能**:
- 余额大卡片（渐变背景）
- 4 个统计小卡片：
  - 已质押
  - 利用率
  - 今日剩余次数
  - 待结算收益
- 小贴士提示框

---

## 类型定义

**文件**: `/workspace/web/lib/api/types.ts`

```typescript
// 热度点配额
export interface HeatPointsQuota {
  heatPoints: number        // 余额
  stakedPoints: number      // 已质押
  pendingInterest: number   // 待结算
  remainingQuota: number    // 剩余次数
  totalQuota: number        // 总次数
  quotaUsed: number         // 已用次数
  singleLimit?: number      // 单人上限
}

// 质押记录
export interface StakeRecord {
  id: number
  topicId: number
  topicTitle?: string
  stakedPoints: number
  flameLevel: number
  status: number        // 0: 质押中，1: 已锁定，2: 已赎回
  settleStatus: number  // 0: 未结算，1: 已结算，2: 已赎回
  unsettledInterest?: number
  createTime: number
  settleTime?: number
  redeemedAt?: number
}
```

---

## API 集成

### 1. 获取配额
```typescript
const quota = await apiClient.get<HeatPointsQuota>("/api/stake/quota")
```

### 2. 创建质押
```typescript
await apiClient.post("/api/stake/create", {
  topicId: 123,
  heatPoints: 10
})
```

### 3. 查询记录
```typescript
const records = await apiClient.get<StakeRecord[]>("/api/stake/records?limit=50")
```

### 4. 赎回
```typescript
await apiClient.post(`/api/stake/redeem/${stakeId}`)
```

---

## 已集成位置

### 1. 主题详情页
**文件**: `/workspace/web/components/topic/topic-detail-actions.tsx`

```tsx
// 原有点赞、评论、收藏按钮旁
<StakeButton topicId={topic.id} currentFlameLevel={topic.flameLevel || 0} />
```

### 2. Topic 数据类型扩展
**后端**: `/workspace/internal/models/resp/response.go`
```go
type TopicResponse struct {
  // ...
  FlameLevel int `json:"flameLevel"`
}
```

**后端渲染**: `/workspace/internal/handlers/render/topic_render.go`
```go
// 计算火焰等级
activeCirculation := heatpoints.HeatSnapshot.GetActiveCirculation()
rsp.FlameLevel = heatpoints.Stake.CalculateFlameLevel(topic.Id, activeCirculation)
```

---

## 待集成位置

### 1. 主题列表项
**文件**: `/workspace/web/components/topic/topic-list-item.tsx`

建议在标题旁添加小号火焰图标：

```tsx
import { FlameLevel } from "@/components/topic/flame-level"

<Link href={`/topic/${topic.id}`}>
  {topic.title}
  {topic.flameLevel > 0 && (
    <FlameLevel level={topic.flameLevel} className="ml-2" />
  )}
</Link>
```

### 2. 用户中心页面
需要在用户中心添加"热度点管理"入口，展示 `HeatPointsPanel` 组件。

### 3. 搜索结果页
搜索结果列表中的主题也可显示火焰等级。

---

## 样式定制

### 火焰颜色
修改 `/workspace/web/components/topic/flame-level.tsx` 中的颜色映射：

```tsx
const colors = {
  1: "text-orange-400",  // 修改等级 1 颜色
  2: "text-orange-500",
  3: "text-orange-600",
  4: "text-red-500",
  5: "text-red-600 animate-pulse",
}
```

### 对话框样式
所有对话框使用 shadcn/ui 组件，可通过全局 CSS 变量定制主题色。

---

## 事件系统

组件间通过自定义事件通信：

```typescript
// 触发更新
window.dispatchEvent(new CustomEvent("heat-stake-updated"))

// 监听更新
window.addEventListener("heat-stake-updated", () => {
  // 刷新数据
})
```

---

## 编译测试

前端编译：
```bash
cd /workspace/web
pnpm build
```

开发模式：
```bash
cd /workspace/web
pnpm dev
```

---

## 下一步开发

1. **列表页集成**: 在主题列表中显示火焰等级
2. **用户中心**: 添加热度点管理页面
3. **管理后台**: 热度点系统监控面板
4. **实时通知**: 结算后推送收益通知
5. **移动端适配**: 优化移动端 UI

---

## 文件清单

```
/workspace/web/
├── components/topic/
│   ├── flame-level.tsx                  # 火焰等级图标
│   ├── stake-button.tsx                 # 质押按钮
│   ├── stake-records-dialog.tsx         # 质押记录对话框
│   └── heat-points-panel.tsx            # 热度点面板
├── lib/api/
│   └── types.ts                         # 类型定义（已扩展）
└── components/topic/
    └── topic-detail-actions.tsx         # 已集成质押按钮
```

---

## 注意事项

1. **事件同步**: 质押成功后必须触发 `heat-stake-updated` 事件
2. **错误处理**: 所有 API 调用都有 try-catch 和 toast 提示
3. **表单验证**: 质押数量必须 >= 1 且 <= 余额
4. **当日限制**: 提示用户当日质押不可赎回
5. **移动端**: 对话框使用响应式设计，支持小屏幕
