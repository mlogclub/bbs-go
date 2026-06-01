# 分区表快速参考

## 一句话总结
**MySQL 原生分区表 = 冷热数据自动分层 + 应用代码零改动 + 秒级归档历史数据**

---

## 核心参数

```yaml
分区类型：范围分区（RANGE）
分区字段：update_time / 100000000
分区粒度：每半年一个分区
保留期限：2 年（自动归档更早数据）
检查频率：每周一 04:00
```

---

## 分区边界速查

| 分区名 | 时间范围 | 边界值 |
|--------|---------|--------|
| p2026h1 | 2026 年前 | < 17672256 |
| p2026h2 | 2026 上半年 | < 17987616 |
| p2027h1 | 2026 下半年 | < 18302976 |
| p2027h2 | 2027 上半年 | < 18618336 |
| p2028h1 | 2027 下半年 | < 18933696 |
| p2028h2 | 2028 上半年 | < 19249056 |
| pmax | 2028 下半年及以后 | MAXVALUE |

---

## SQL 查询

### 查看分区状态
```sql
SELECT partition_name, table_rows, data_length/1024/1024 AS data_mb
FROM information_schema.partitions
WHERE table_name = 't_topic_stake';
```

### 手动归档分区
```sql
-- 删除 2024 年上半年数据
ALTER TABLE t_topic_stake DROP PARTITION p2024h1;
```

### 添加新分区
```sql
-- 添加 2029 年上半年分区
ALTER TABLE t_topic_stake ADD PARTITION (
    PARTITION p2029h1 VALUES LESS THAN (20459520)
);
```

---

## Go 代码调用

```go
// 自动管理（推荐）
// 定时任务会自动执行，无需手动调用

// 手动检查
err := heatpoints.Partition.CheckAndMigrate()

// 手动归档
err := heatpoints.Partition.ArchiveOldPartition(2) // 保留 2 年

// 查看统计
stats, _ := heatpoints.Partition.GetPartitionStats()
```

---

## 性能对比

| 操作 | 普通表 | 分区表 |
|------|--------|--------|
| 查询近 3 个月 | 全表扫描 | 扫描 1 个分区 |
| 删除 2 年前数据 | DELETE 数小时 | DROP 分区 <0.01 秒 |
| 统计总行数 | COUNT(*) 全表 | 各分区累加 |
| 索引维护 | 单大索引 | 多个小索引 |

---

## 监控要点

- ⚠️ 单分区行数 > 100 万
- ⚠️ 单分区大小 > 1GB
- ⚠️ 最新分区使用率 > 80%
- ✅ 每周一自动检查
- ✅ 自动创建新分区

---

## 故障处理

**Q: 分区未生效？**
```sql
SHOW CREATE TABLE t_topic_stake; -- 查看是否有 PARTITION 子句
```

**Q: 查询仍慢？**
```sql
EXPLAIN SELECT ... -- 检查是否只扫描相关分区
```

**Q: 如何回滚？**
```bash
# 从备份恢复
# 或等待下次无数据时重新迁移
```

---

## 文件位置

```
/workspace/
├── migrations/
│   └── 000017_partition_topic_stakes.go  # 迁移脚本
├── internal/models/
│   └── models.go  # TopicStake 模型
├── internal/services/heatpoints/
│   └── partition_service.go  # 分区管理
├── internal/scheduler/
│   └── cron.go  # 定时任务
└── PARTITION_TABLE_IMPLEMENTATION.md  # 完整文档
```

---

**详细说明**: 阅读 `PARTITION_TABLE_IMPLEMENTATION.md`  
**快速部署**: 系统启动时自动执行迁移
