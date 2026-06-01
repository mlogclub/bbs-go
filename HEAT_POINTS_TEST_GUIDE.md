# 热度点质押系统 - 测试指南

## 快速启动测试环境

### 1. 准备数据库

```bash
cd /workspace
rm -f bbs-go.db  # 清理旧数据（可选）
```

### 2. 启动系统

```bash
# 编译
go build -o bbs-go ./main.go

# 运行
./bbs-go
```

### 3. 完成安装

访问 `http://localhost:8082`，按照向导完成安装。

**安装完成后自动触发**:
- ✅ 创世空投（50 点/用户）
- ✅ 初始化系统常量

---

## 手动测试流程

### 场景 1: 创世空投验证

**步骤**:
1. 安装完成后创建/登录用户
2. 访问用户中心或直接查询 API

**API 测试**:
```bash
# 获取 auth token（从浏览器 localStorage 或 Cookie）
TOKEN="your_auth_token"

# 查询余额
curl http://localhost:8082/api/stake/quota \
  -H "Authorization: Bearer $TOKEN"
```

**预期结果**:
```json
{
  "code": 0,
  "data": {
    "heatPoints": 50,
    "stakedPoints": 0,
    "pendingInterest": 0,
    "remainingQuota": 3,
    "totalQuota": 3
  }
}
```

---

### 场景 2: 创建质押

**API 测试**:
```bash
curl -X POST http://localhost:8082/api/stake/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topicId": 1,
    "heatPoints": 10
  }'
```

**预期结果**:
```json
{
  "code": 0,
  "data": {
    "stakeId": 1,
    "remainingQuota": 2,
    "flameLevel": 1,
    "riskLevel": "normal",
    "riskHint": "该帖处于冷帖期，收益波动较大"
  }
}
```

**前端测试**:
1. 访问主题详情页 `http://localhost:8082/topic/1`
2. 点击"质押"按钮
3. 输入质押数量（如 10）
4. 确认质押
5. 检查 toast 提示和按钮状态

---

### 场景 3: 配额限制

**测试**:
```bash
# 第 1 次
curl -X POST ... -d '{"topicId": 1, "heatPoints": 5}'
# ✅ 成功

# 第 2 次
curl -X POST ... -d '{"topicId": 2, "heatPoints": 5}'
# ✅ 成功

# 第 3 次
curl -X POST ... -d '{"topicId": 3, "heatPoints": 5}'
# ✅ 成功

# 第 4 次
curl -X POST ... -d '{"topicId": 4, "heatPoints": 5}'
# ❌ 错误：今日质押次数已用完
```

---

### 场景 4: 当日不可赎回

**步骤**:
1. 创建质押
2. 立即尝试赎回

**API 测试**:
```bash
# 创建质押
curl -X POST http://localhost:8082/api/stake/create \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"topicId": 1, "heatPoints": 10}'

# 获取 stakeId（从响应中）
STAKE_ID=1

# 尝试赎回
curl -X POST http://localhost:8082/api/stake/redeem/$STAKE_ID \
  -H "Authorization: Bearer $TOKEN"
```

**预期结果**:
```json
{
  "code": 400,
  "message": "当日质押不可赎回，请等到次日结算后"
}
```

---

### 场景 5: 查询质押记录

**API 测试**:
```bash
curl http://localhost:8082/api/stake/records \
  -H "Authorization: Bearer $TOKEN"
```

**预期响应**:
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "topicId": 1,
      "topicTitle": "测试主题",
      "stakedPoints": 10,
      "flameLevel": 2,
      "status": 0,
      "settleStatus": 0,
      "createTime": 1717027200
    }
  ]
}
```

**前端测试**:
1. 打开用户中心
2. 点击"我的质押"
3. 查看记录列表
4. 切换"进行中"和"历史记录"标签

---

### 场景 6: 火焰等级显示

**步骤**:
1. 访问主题详情页
2. 查看质押按钮颜色
3. 不同火焰等级应显示不同颜色

**火焰等级验证表**:
| 质押量/流通比 | 火焰等级 | 按钮颜色 |
|--------------|---------|---------|
| > 0.1% | 🔥 | 橙色 |
| > 0.5% | 🔥🔥 | 橙色 |
| > 1% | 🔥🔥🔥 | 橙红色 |
| > 3% | 🔥🔥🔥🔥 | 红色 |
| > 10% | 🔥🔥🔥🔥🔥 | 深红脉冲 |

---

### 场景 7: 定时任务验证

**查看日志**:
```bash
# 启动时查看日志输出
./bbs-go 2>&1 | grep -E "heat|snapshot|settlement"
```

**预期日志**:
```
[23:55:00] 开始执行热度点快照
[23:55:01] 生成交互快照 100 条
[23:55:01] 生成流通快照：总存量 5000, 活跃流通 3500
[23:55:02] 生成火焰偏移量

[00:00:00] 开始执行热度点结算
[00:00:01] 结算 150 个质押记录
[00:00:02] 发放收益 235 点
[00:00:03] 执行衰减回收：12 点
```

---

## 前端 UI 测试清单

### 组件级测试

| 组件 | 测试项 | 预期行为 |
|------|-------|---------|
| `FlameLevel` | level=0 | 不显示 |
| `FlameLevel` | level=1-5 | 显示对应数量火焰 |
| `FlameLevel` | level=5 | 红色脉冲动画 |
| `StakeButton` | 未登录 | 显示登录提示 |
| `StakeButton` | 余额不足 | 输入框提示错误 |
| `StakeButton` | 成功质押 | toast 提示 + 关闭对话框 |
| `StakeRecordsDialog` | 无记录 | 显示"暂无记录" |
| `StakeRecordsDialog` | 有记录 | 正确显示列表 |
| `StakeRecordsDialog` | 赎回按钮 | 当日质押不显示赎回 |
| `HeatPointsPanel` | 加载中 | 显示骨架屏 |
| `HeatPointsPanel` | 数据加载 | 显示余额和统计 |

---

### 集成测试流程

```bash
# 1. 打开主题列表页
open http://localhost:8082/topics

# 2. 选择任一主题，检查是否显示火焰图标
# 预期：新主题显示无火焰或 1 级火焰

# 3. 进入主题详情页
open http://localhost:8082/topic/1

# 4. 点击"质押"按钮
# 预期：弹出对话框，显示余额和配额

# 5. 输入质押数量，点击"确认质押"
# 预期：toast 提示"质押成功"，对话框关闭

# 6. 再次点击"质押"（如果还有次数）
# 预期：可继续质押

# 7. 打开用户中心
open http://localhost:8082/user/current

# 8. 找到"热度点"面板
# 预期：显示余额、已质押、利用率等

# 9. 点击"查看质押记录"
# 预期：打开对话框，显示刚才的质押记录
```

---

## 三阶段利率测试

### 模拟不同阶段

**设置测试数据**（需要直接操作数据库）：

```sql
-- 冷帖期：质押量 < 0.5% 活跃流通
-- 设置活跃流通 = 10000, 帖子质押 = 10
UPDATE t_heat_circulation_snapshot SET active_circulation = 10000 WHERE id = 1;
INSERT INTO t_topic_interaction_snapshot (topic_id, staked_points) VALUES (1, 10);

-- 共识期：0.5% <= 质押量 < 3%
-- 设置帖子质押 = 100 (1%)
UPDATE t_topic_interaction_snapshot SET staked_points = 100 WHERE topic_id = 1;

-- 热门期：质押量 >= 3%
-- 设置帖子质押 = 500 (5%)
UPDATE t_topic_interaction_snapshot SET staked_points = 500 WHERE topic_id = 1;
```

### 验证利率计算

**测试步骤**:
1. 创建质押
2. 手动触发结算（修改系统时间或调用内部方法）
3. 检查利息计算

**预期利率**:
| 阶段 | 条件 | 利率范围 |
|------|------|---------|
| 冷帖期 | < 0.5% | ±50% (高波动) |
| 共识期 | 0.5%-3% | ±20% (稳定) |
| 热门期 | ≥ 3% | +2%/-30% (断崖) |

---

## 衰减机制测试

### 模拟衰减

**设置低利用率**:

```sql
-- 用户有 100 点，但近 7 天只质押了 50 点
-- 未使用量 = 50
-- 衰减量 = 50 * 2% = 1 点

UPDATE t_user_heat_stats 
SET heat_points = 100, last_7days_stake_total = 50
WHERE user_id = 1;
```

**触发结算后检查**:
```sql
-- 检查余额变化
SELECT heat_points FROM t_user_heat_stats WHERE user_id = 1;
-- 预期：99 点（衰减 1 点）

-- 检查公共奖池记录
SELECT * FROM t_heat_public_pool WHERE source_type = 'decay';
-- 预期：有一条 decay 收入记录
```

---

## 压力测试

### 批量质押测试

```bash
# 模拟 100 次质押（需要多个用户）
for i in {1..100}; do
  curl -X POST http://localhost:8082/api/stake/create \
    -H "Authorization: Bearer $TOKEN_$i" \
    -d "{\"topicId\": 1, \"heatPoints\": 1}"
done
```

### 结算性能测试

```bash
# 生成大量测试数据后，手动触发结算
# 检查日志中的处理时间和批次数

# 预期：
# - 每批 500 条记录
# - 1000 条记录应在 2 批内完成
# - 总耗时 < 5 秒
```

---

## 错误场景测试

### 1. 余额不足
```bash
curl -X POST http://localhost:8082/api/stake/create \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"topicId": 1, "heatPoints": 1000}'
# ❌ 错误：热度点余额不足
```

### 2. 不存在主题
```bash
curl -X POST http://localhost:8082/api/stake/create \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"topicId": 99999, "heatPoints": 10}'
# ❌ 错误：主题不存在
```

### 3. 赎回不存在的记录
```bash
curl -X POST http://localhost:8082/api/stake/redeem/99999 \
  -H "Authorization: Bearer $TOKEN"
# ❌ 错误：质押记录不存在
```

### 4. 重复赎回
```bash
# 第 1 次赎回
curl -X POST http://localhost:8082/api/stake/redeem/1 \
  -H "Authorization: Bearer $TOKEN"
# ✅ 成功

# 第 2 次赎回（同一记录）
curl -X POST http://localhost:8082/api/stake/redeem/1 \
  -H "Authorization: Bearer $TOKEN"
# ❌ 错误：该质押已赎回
```

---

## 验证清单

完成后检查以下项目：

- [ ] 创世空投已执行（所有用户有 50 点）
- [ ] 质押按钮正常显示
- [ ] 火焰等级图标正确渲染
- [ ] 创建质押 API 正常
- [ ] 配额限制生效
- [ ] 当日不可赎回
- [ ] 质押记录可查询
- [ ] 23:55 快照任务执行
- [ ] 00:00 结算任务执行
- [ ] 三阶段利率计算正确
- [ ] 衰减机制生效
- [ ] 公共奖池收支记录
- [ ] 前端 UI 组件正常工作
- [ ] 无编译错误
- [ ] 日志输出正常

---

## 调试技巧

### 1. 查看数据库
```bash
sqlite3 bbs-go.db
SELECT * FROM t_topic_stake ORDER BY id DESC LIMIT 10;
SELECT * FROM t_heat_circulation_snapshot ORDER BY id DESC LIMIT 5;
SELECT * FROM t_heat_public_pool ORDER BY id DESC LIMIT 10;
```

### 2. 查看日志
```bash
# 只看热度相关
./bbs-go 2>&1 | grep "heat"

# 实时查看
tail -f logs/app.log | grep -E "stake|settle|snapshot"
```

### 3. 修改系统时间（测试结算）
```bash
# Linux
date -s "2026-05-31 00:00:00"

# 注意：会影响所有定时任务
```

### 4. 手动触发结算（开发环境）
```go
// 在代码中临时添加测试端点
func TestSettle(ctx *gin.Context) {
  err := heatpoints.Settlement.SettleAll()
  ginx.WriteJSON(ctx, map[string]any{"error": err})
}
```

---

## 测试报告模板

```markdown
## 测试报告

**测试日期**: YYYY-MM-DD
**测试人员**: [姓名]
**环境**: [开发/测试/生产]

### 通过场景
- [x] 创世空投
- [x] 创建质押
- [x] 配额限制
- [ ] ...

### 发现问题
1. **问题描述**: ...
   **严重程度**: 高/中/低
   **复现步骤**: ...

### 性能数据
- 快照生成耗时: X 秒
- 结算 1000 条记录耗时: Y 秒
- API 平均响应时间: Z ms

### 结论
[通过/不通过]
```

---

## 自动化测试建议

后续可添加：
1. **API 集成测试**: Go 测试文件
2. **E2E 测试**: Playwright/Cypress
3. **性能benchmark**: 压制定量指标
