# 热度点系统 - 快速参考卡

## 核心参数

```yaml
创世空投：50 点/人
每日配额：3 次质押
最低额度：1 点
最高额度：活跃流通 × 0.3%
衰减率：2%/日
热门阈值：3% 活跃流通
共识阈值：0.5% 活跃流通
```

## 火焰等级

```
🔥          > 0      (冷帖期 ±50%)
🔥🔥        ≥ 0.1%
🔥🔥🔥      ≥ 0.5%   (共识期 ±20%)
🔥🔥🔥🔥    ≥ 1%
🔥🔥🔥🔥🔥  ≥ 3%     (热门期 +2%/-30% 断崖)
```

## API 端点

```
GET  /api/stake/quota        - 查询配额
GET  /api/stake/records      - 查询记录
POST /api/stake/create       - 创建质押
POST /api/stake/redeem/:id   - 赎回质押
GET  /api/stake/heat/:id     - 火焰等级
```

## 定时任务

```
23:55 → 快照生成 (互动 + 流通 + 偏移)
00:00 → 每日结算 (利率 + 收益 + 衰减)
```

## 三阶段利率

| 阶段 | 条件 | 利率计算 | 波动范围 |
|------|------|----------|----------|
| 冷帖期 | < 0.5% | 互动×1.5 + 质押×1.0 | ±50% |
| 共识期 | 0.5%-3% | 互动×0.8 + 质押×0.6 | ±20% |
| 热门期 | ≥ 3% | 互动×0.1 | +2%/-30% |

## 关键公式

```
活跃流通 = Σ max(用户质押量，近 7 天质押总量)
利用率 = 近 7 天质押 / 持有总量
衰减量 = max(持有量 - 近 7 天质押，0) × 2%
火焰等级 = 基于 质押量/活跃流通 百分比
```

## 前端组件

```tsx
<FlameLevel level={3} />                    // 火焰图标
<StakeButton topicId={1} currentFlameLevel={2} />  // 质押按钮
<StakeRecordsDialog open={true} onOpenChange={setOpen} />  // 记录
<HeatPointsPanel />                         // 用户面板
```

## 文件位置

```
后端：/workspace/internal/services/heatpoints/
前端：/workspace/web/components/topic/
文档：/workspace/HEAT_POINTS_*.md
迁移：/workspace/migrations/000016_*.go
```

## 编译运行

```bash
# 后端
cd /workspace && go build -o bbs-go ./main.go && ./bbs-go

# 前端
cd /workspace/web && pnpm dev
```

## 测试命令

```bash
# 查询配额
curl http://localhost:8082/api/stake/quota -H "Authorization: Bearer TOKEN"

# 创建质押
curl -X POST http://localhost:8082/api/stake/create \
  -H "Authorization: Bearer TOKEN" \
  -d '{"topicId":1,"heatPoints":10}'
```

## 常见问题

**Q: 当日质押何时可赎回？**
A: 次日 00:00 结算后即可赎回

**Q: 衰减会影响已质押的吗？**
A: 不会，只针对未质押的余额

**Q: 火焰等级实时变化吗？**
A: 是的，每次加载时重新计算

**Q: 公共奖池空了怎么办？**
A: 正收益等比削减，签到照发（优先级最高）

**Q: 如何查看活跃流通量？**
A: 暂未公开，可通过 API 间接推算

---

**快速开始**: 阅读 `HEAT_POINTS_README.md`  
**完整文档**: 查看 `HEAT_POINTS_FINAL_SUMMARY.md`  
**测试指南**: 参考 `HEAT_POINTS_TEST_GUIDE.md`
