# 热度点质押系统 - 基础框架实现说明

## 已完成的工作

### 1. 数据模型层 ✅
- **文件**: `/workspace/internal/models/models.go`
- **新增模型**:
  - `User.HeatPoints` - 用户热度点余额字段
  - `Topic.EverViral` / `Topic.FlameLockedLevel` - 帖子热门标记字段
  - `TopicStake` - 质押记录
  - `UserHeatLog` - 热度点流水
  - `UserHeatStats` - 用户热度统计
  - `HeatPublicPool` - 公共奖池
  - `SystemMintLog` - 系统铸币日志
  - `TopicInteractionSnapshot` - 每日互动快照
  - `HeatCirculationSnapshot` - 活跃流通快照
  - `DailyFlameOffset` - 火焰等级偏移量
  - `SettlementTaskLog` - 结算任务日志

### 2. 常量定义 ✅
- **文件**: `/workspace/internal/models/constants/constants.go`
- **新增常量**:
  - 热度点系统参数（创世空投 50 点、每日签到 0.05%、衰减率 2% 等）
  - 阶段阈值（冷帖期 0.5%、热门期 3%）
  - 火焰等级阈值
  - 流水类型常量
  - 公共奖池来源常量

### 3. 数据库迁移 ✅
- **文件**: `/workspace/migrations/000016_migration_script_heat_points_system.go`
- 创建所有新增表
- 为现有表添加字段
- 初始化公共奖池和用户统计

### 4. 基础服务 ✅
- **创世空投服务**: `/workspace/internal/services/heatpoints/genesis_airdrop.go`
  - 系统初始化时给所有用户发放 50 热度点
  - 记录铸币日志和流水

### 5. 定时任务集成 ✅
- **文件**: `/workspace/internal/scheduler/cron.go`
- 每日 23:55 自动生成快照

## 需要补充的功能

由于时间和复杂度限制，以下功能需要后续完善：

### 1. Repository 层
**需要创建**: `/workspace/internal/repositories/topic_stake_repository.go`
```go
// 参考 existing repository 模式实现 CRUD
```

### 2. 完整服务层
需要实现的服务（按优先级）:
1. **StakeService** - 质押核心服务
   - Create (创建质押)
   - Redeem (赎回)
   - GetUserStakes (查询记录)

2. **SettlementService** - 每日结算服务
   - 三阶段利率计算
   - 正收益/负收益结算
   - 公共奖池收支管理

3. **DecayService** - 衰减服务
   - 利用率计算
   - 每日衰减扣除
   - 60 天强制归零

### 3. API Handlers
**文件**: `/workspace/internal/handlers/api/heat_handlers.go` 已创建但需要调整

### 4. 路由注册
**文件**: `/workspace/internal/server/router.go` 已添加路由

## 系统使用流程

### 1. 系统初始化
```go
// 安装完成后、首次启动前调用
heatpoints.Airdrop.Execute()
```

### 2. 用户质押流程
```
POST /api/stake/create
Body: { "topicId": 123, "heatPoints": 10 }

# 返回值
{
  "stakeId": 456,
  "remainingQuota": 2,
  "flameLevel": 3,
  "riskLevel": "medium"
}
```

### 3. 用户赎回流程
```
POST /api/stake/redeem/{stakeId}

# 返回值
{
  "code": 0,
  "msg": ""
}
```

### 4. 查询配额
```
GET /api/stake/quota

# 返回值
{
  "remainingQuota": 2,
  "totalQuota": 3,
  "heatPoints": 150
}
```

## 后续开发建议

### 第一阶段：让系统跑起来（必须）
1. 编译并测试迁移
2. 运行创世空投
3. 测试基础质押 API

### 第二阶段：核心功能（高优先级）
1. 实现三阶段利率模型
2. 实现每日结算定时任务
3. 实现衰减机制

### 第三阶段：完善体验（中优先级）
1. 前端质押 UI 组件
2. 用户中心热度点面板
3. 帖子列表火焰等级展示

### 第四阶段：高级功能（低优先级）
1. 管理后台监控面板
2. 活跃度流通统计
3. 排名奖励系统

## 关键技术点说明

### 1. 活跃流通计算
```go
// 简化版本：直接从快照读取
activeCirculation, _ := HeatSnapshot.GetActiveCirculation()

// 完整版本需要实现：
// activeCirculation = Σ max(userStaked, userLast7DaysStakeTotal)
```

### 2. 火焰等级
```go
// 基于 ratio = 帖子质押量 / 活跃流通量
if ratio >= 3%   → 🔥🔥🔥🔥🔥 (5 档)
if ratio >= 1%   → 🔥🔥🔥🔥    (4 档)
if ratio >= 0.5% → 🔥🔥🔥      (3 档)
if ratio >= 0.1% → 🔥🔥        (2 档)
if ratio > 0     → 🔥          (1 档)
```

### 3. 公共奖池收支
- **收入**: 衰减回收、结算亏损、60 天归零
- **支出**: 签到发放（优先级 1）、正收益结算（优先级 2）、排名奖励（优先级 3）

## 文档版本
- **创建时间**: 2026-05-30
- **当前状态**: 基础框架完成，等待功能完善
- **参考文档**: `/workspace/.monkeycode-tmp-files/ae3b981f-heat-points-plan-1.md`
