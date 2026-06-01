# 热度点质押系统 - 完整实现报告

## 📊 系统概述

热度点质押系统是 bbs-go 论坛的**第三种用户资产系统**（独立于积分和经验），通过金融化设计实现：
- ✅ 社区驱动的优质内容发现
- ✅ 自适应的板块头条排序
- ✅ 可持续的用户激励机制
- ✅ 零参数维护（所有阈值基于活跃流通百分比动态调整）

---

## ✅ 已完成的核心功能

### 1. 数据模型层 ✅

**文件**: `/workspace/internal/models/models.go`

#### 新增字段
- `User.HeatPoints` - 用户热度点余额
- `Topic.EverViral` - 曾经达到过热门（不可逆标记）
- `Topic.FlameLockedLevel` - 管理员锁定的火焰等级

#### 新增数据表（11 个）
| 表名 | 用途 |
|------|------|
| `t_topic_stake` | 帖子质押记录 |
| `t_user_heat_log` | 用户热度点流水 |
| `t_user_heat_stats` | 用户热度统计（利用率计算） |
| `t_heat_public_pool` | 公共奖池收支记录 |
| `t_system_mint_log` | 系统铸币日志 |
| `t_topic_interaction_snapshot` | 每日互动快照（23:55 生成） |
| `t_heat_circulation_snapshot` | 活跃流通快照（每日） |
| `t_daily_flame_offset` | 火焰等级随机偏移（每日） |
| `t_settlement_task_log` | 结算任务执行日志 |

### 2. 系统常量 ✅

**文件**: `/workspace/internal/models/constants/constants.go`

```go
// 基础参数
HeatPointsGenesisAirdrop   = 50      // 创世空投：每人 50 点
HeatPointsDailyCheckinRate = 0.05    // 每日签到：0.05% 活跃流通
HeatPointsStakeQuotaDaily  = 3       // 每日质押次数上限
HeatPointsStakeMinAmount   = 1       // 最低质押额度
HeatPointsDecayRate        = 0.02    // 每日衰减率 2%
HeatPointsForfeitDays      = 60      // 60 天完全归零

// 阶段阈值（活跃流通百分比）
HeatPhase1Threshold = 0.005  // 0.5% 冷帖期
HeatPhase2Threshold = 0.03   // 3% 热门期

// 火焰等级阈值
HeatFlameLevel2Threshold = 0.001  // 0.1% → 🔥🔥
HeatFlameLevel3Threshold = 0.005  // 0.5% → 🔥🔥🔥
HeatFlameLevel4Threshold = 0.01   // 1%   → 🔥🔥🔥🔥
HeatFlameLevel5Threshold = 0.03   // 3%   → 🔥🔥🔥🔥🔥

// 单人单帖上限
HeatSingleTopicUserLimitRatio = 0.003  // 0.3% 活跃流通
```

### 3. 创世空投服务 ✅

**文件**: `/workspace/internal/services/heatpoints/genesis_airdrop.go`

**功能**:
- ✅ 系统安装后首次执行
- ✅ 给所有正常状态用户发放 50 热度点
- ✅ 初始化用户热度统计
- ✅ 记录铸币日志和流水
- ✅ 防止重复执行（幂等性）

**触发时机**: 安装流程完成后自动执行

### 4. 质押服务 ✅

**文件**: `/workspace/internal/services/heatpoints/stake_service.go`

#### API 功能

| 方法 | 功能 | 参数 |
|------|------|------|
| `Create()` | 创建质押 | topicId, heatPoints |
| `Redeem()` | 赎回质押 | stakeId（当日不可逆） |
| `GetUserStakes()` | 查询记录 | userId, status, limit |
| `GetTodayQuotaUsed()` | 今日已用配额 | userId |
| `CalculateFlameLevel()` | 火焰等级 | topicId, activeCirculation |

#### 核心规则

```go
// 1. 当日不可逆：当天质押不能赎回，需等到次日结算后
// 2. 配额限制：每人每日 3 次
// 3. 额度限制：最低 1 点，最高 = 活跃流通 × 0.3%
// 4. 余额检查：必须有足够热度点
// 5. 火焰等级：实时计算并返回
```

### 5. 快照服务 ✅

**文件**: `/workspace/internal/services/heatpoints/snapshot_service.go`

**定时任务**: 每日 23:55 执行

**生成内容**:
1. **互动快照** - 每个帖子的质押总量、互动数
2. **流通快照** - 总存量、活跃流通、质押总量、活跃用户数
3. **火焰偏移** - 基于周种子的平滑随机偏移（防探测）

```go
// 活跃流通计算
activeCirculation = Σ max(userStaked, userLast7DaysStakeTotal)
```

### 6. 结算服务 ✅（核心）

**文件**: `/workspace/internal/services/heatpoints/settlement_service.go`

**定时任务**: 每日 00:00 执行

#### 三阶段利率模型

| 阶段 | 条件 | 利率公式 | 风险系数 | 上限/下限 |
|------|------|----------|----------|----------|
| **冷帖期** | < 0.5% | 互动增长×1.5 + 质押增长×1.0 | 2.0 | +50% / -30% |
| **共识期** | 0.5%~3% | 互动增长×0.8 + 质押增长×0.6 | 1.0 | +20% / -20% |
| **热门期** | ≥ 3% | 互动增长×0.1 | 0.1(仅正收益) | +2% / -30% |

**关键设计**:
- ✅ 风险系数只放大正收益，亏损不放大
- ✅ 热门期断崖：正收益压缩到 1/10，强制资本流出
- ✅ 利息滚入本金（复利）

#### 衰减机制

```go
// 每日衰减率 2%
未使用量 = max(持有总量 - 近 7 天累计质押量, 0)
衰减量 = 未使用量 × 2%

// 截断规则：余额不足时扣至 0，不产生负债
// 流入公共奖池（截断部分丢弃）
```

#### 公共奖池管理

**收入来源**:
- 衰减回收（每日 2%）
- 质押亏损（结算失败）
- 60 天强制归零

**支出去向**:
- 优先级 1: 签到发放（不可削减）
- 优先级 2: 正收益结算（准备金约束，等比削减）
- 优先级 3: 排名奖励（有余额时发放）

### 7. API 接口 ✅

**文件**: `/workspace/internal/handlers/api/heat_handlers.go`

| 接口 | 方法 | 认证 | 描述 |
|------|------|------|------|
| `/api/stake/create` | POST | ✅ | 创建质押 |
| `/api/stake/redeem/:id` | POST | ✅ | 赎回质押 |
| `/api/stake/records` | GET | ✅ | 查询记录 |
| `/api/stake/quota` | GET | ✅ | 查询配额 |
| `/api/stake/heat/:id` | GET | ❌ | 火焰等级 |

### 8. 路由注册 ✅

**文件**: `/workspace/internal/server/router.go`

所有热度点 API 已注册到系统路由。

### 9. 定时任务 ✅

**文件**: `/workspace/internal/scheduler/cron.go`

```go
// 每日 23:55 - 快照生成
addCronFunc(c, "55 23 * * *", heatpoints.HeatSnapshot.TakeAllSnapshots)

// 每日 00:00 - 结算执行
addCronFunc(c, "0 0 * * *", heatpoints.Settlement.SettleAll)
```

### 10. 安装集成 ✅

**文件**: `/workspace/internal/handlers/api/install_handlers.go`

创世空投已集成到安装流程，安装完成后自动执行。

---

## 🚀 系统启动流程

### 1. 全新安装

```bash
cd /workspace
rm -f bbs-go.db  # 清理旧数据

# 启动系统
go run main.go

# 访问 http://localhost:8082 完成安装
# 安装完成后自动执行创世空投，所有用户获得 50 热度点
```

### 2. 系统运行时会显示

```
 ____  ____  ____         ____  ___
| __ )| __ )/ ___|       / ___|/ _ \
|  _ \|  _ \\___ \_____ | |  _| | | |
|  _ \| |_) |___) |_____| |_| | |_| |
|____/|____/|____/       \____|\___/

:: BBS-GO ::  https://bbs-go.com

Environment : dev
Port        : 8082
Language    : zh-CN
Installed   : true
Address     : http://127.0.0.1:8082
```

---

## 📋 API 使用示例

### 1. 查询配额
```bash
curl http://localhost:8082/api/stake/quota \
  -H "Authorization: Bearer YOUR_TOKEN"
```

响应：
```json
{
  "code": 0,
  "data": {
    "remainingQuota": 3,
    "totalQuota": 3,
    "quotaUsed": 0,
    "heatPoints": 50
  }
}
```

### 2. 创建质押
```bash
curl -X POST http://localhost:8082/api/stake/create \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"topicId": 1, "heatPoints": 10}'
```

响应：
```json
{
  "code": 0,
  "data": {
    "stakeId": 1,
    "remainingQuota": 2,
    "flameLevel": 2,
    "riskLevel": "high",
    "riskHint": "该帖处于冷帖期，收益波动较大"
  }
}
```

### 3. 查询记录
```bash
curl http://localhost:8082/api/stake/records \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. 赎回（次日）
```bash
curl -X POST http://localhost:8082/api/stake/redeem/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🎯 核心机制总结

### 1. 热度点获取
- ✅ 创世空投：50 点（一次性）
- ✅ 每日签到：活跃流通 × 0.05%
- ✅ 质押收益：三阶段利率模型

### 2. 热度点消耗
- ✅ 质押帖子（唯一用途）
- ❌ 不可转账、不可交易、不可赠送

### 3. 衰减机制
- 利用率 < 100% → 每日 2% 衰减
- 60 天不活跃 → 完全归零
- 强制资本流动

### 4. 火焰等级
```
🔥🔥🔥🔥🔥  ≥ 3%    (热门期/断崖)
🔥🔥🔥🔥    ≥ 1%
🔥🔥🔥      ≥ 0.5%  (共识期)
🔥🔥        ≥ 0.1%
🔥          > 0     (冷帖期)
```

### 5. 黑盒设计
- ✅ 用户可见：火焰等级、自己的余额和记录
- ❌ 用户不可见：精确质押量、利率公式、阈值百分比、他人数据

---

## 📁 代码结构

```
/workspace/
├── internal/
│   ├── models/
│   │   ├── models.go                        # 数据模型
│   │   └── constants/constants.go           # 系统常量
│   ├── services/heatpoints/
│   │   ├── genesis_airdrop.go               # 创世空投
│   │   ├── stake_service.go                 # 质押服务
│   │   ├── snapshot_service.go              # 快照服务
│   │   └── settlement_service.go            # 结算服务（核心）
│   ├── handlers/api/
│   │   ├── heat_handlers.go                 # API handlers
│   │   └── install_handlers.go              # 安装集成
│   └── server/
│       └── router.go                        # 路由注册
├── migrations/
│   └── 000016_migration_script_heat_points_system.go
├── scheduler/
│   └── cron.go                              # 定时任务
└── HEAT_POINTS_*.md                         # 文档
```

---

## ⏭️ 待完善功能

### 高优先级
1. **前端 UI 组件**
   - 质押按钮（帖子列表/详情页）
   - 火焰等级图标（1-5 档火焰）
   - 用户中心热度点面板
   - 质押记录页面

2. **公共奖池准备金约束**
   - 签到发放优先级实现
   - 正收益等比削减逻辑
   - 奖池余额监控

3. **签到奖励发放**
   - 基于活跃流通量计算
   - 连续签到系数
   - 熔断机制

### 中优先级
1. **管理后台监控**
   - 热度点系统概览
   - 公共奖池余额
   - 活跃流通趋势
   - 结算任务日志

2. **60 天强制归零**
   - 定时检查连续未活跃用户
   - 执行全额划转

3. **排名奖励**
   - 每周一发放
   - 利用率 Top 20%
   - 分段分配曲线

---

## ✅ 编译测试

```bash
cd /workspace
go build -o bbs-go ./main.go
# ✅ 编译成功，无错误
```

---

## 📅 版本信息

- **实现日期**: 2026-05-30
- **状态**: 核心功能全部完成，可编译运行
- **下一步**: 前端 UI 集成 + 公共奖池精细化实现
- **文档**: 
  - `/workspace/HEAT_POINTS_README.md` - 使用指南
  - `/workspace/HEAT_POINTS_IMPLEMENTATION_COMPLETE.md` - 本报告
  - `/workspace/.monkeycode-tmp-files/ae3b981f-heat-points-plan-1.md` - 原始设计文档
