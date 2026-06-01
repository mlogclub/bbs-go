package heatpoints

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"fmt"
	"log/slog"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// ColdStorageService 冷数据归档服务
// 基于行数阈值 + 时间条件的双重触发机制，将已结束的历史数据迁移到归档表，
// 保证主表行数始终可控（topic_stakes ≤ 30 万行）。
type ColdStorageService struct{}

var ColdStorage = &ColdStorageService{}

// ArchiveConfig 单张表的归档配置
type ArchiveConfig struct {
	TableName    string // 主表名
	ArchiveTable string // 归档表名
	TimeColumn   string // 时间判断字段
	TimeDays     int    // 超过 N 天才算冷数据
	RowLimit     int    // 主表行数阈值（超过才触发）
	TargetRows   int    // 归档后目标行数
	BatchSize    int    // 每批迁移行数
	StatusColumn string // 状态字段（"status"，仅 status != 0 的记录可归档）
}

// 默认归档配置
var archiveConfigs = []ArchiveConfig{
	{
		TableName:    "t_topic_stake",
		ArchiveTable: "t_topic_stake_archive",
		TimeColumn:   "update_time",
		TimeDays:     90,
		RowLimit:     300000,
		TargetRows:   250000,
		BatchSize:    5000,
		StatusColumn: "status",
	},
	{
		TableName:    "t_user_heat_log",
		ArchiveTable: "t_user_heat_log_archive",
		TimeColumn:   "create_time",
		TimeDays:     180,
		RowLimit:     2000000,
		TargetRows:   1500000,
		BatchSize:    10000,
		StatusColumn: "", // 流水表无状态字段，全量按时间归档
	},
	{
		TableName:    "t_topic_interaction_snapshot",
		ArchiveTable: "t_topic_interaction_snapshot_archive",
		TimeColumn:   "create_time",
		TimeDays:     30,
		RowLimit:     500000,
		TargetRows:   300000,
		BatchSize:    20000,
		StatusColumn: "",
	},
}

// RunAll 对所有配置的表执行归档检查
func (s *ColdStorageService) RunAll() {
	slog.Info("开始冷数据归档检查")

	for _, cfg := range archiveConfigs {
		result := s.archiveOne(cfg)
		if result.Migrated > 0 {
			slog.Info("归档完成",
				"table", cfg.TableName,
				"migrated", result.Migrated,
				"rows_before", result.RowsBefore,
				"rows_after", result.RowsAfter,
				"trigger", result.Trigger,
				"duration_ms", result.DurationMs,
			)

			// 写入归档日志
			sqls.DB().Create(&models.ColdStorageLog{
				TableName:    cfg.TableName,
				ArchiveTable: cfg.ArchiveTable,
				TriggeredBy:  result.Trigger,
				RowsBefore:   result.RowsBefore,
				RowsMigrated: result.Migrated,
				RowsAfter:    result.RowsAfter,
				DurationMs:   result.DurationMs,
				Status:       0, // 成功
				CreateTime:   dates.NowTimestamp(),
			})
		}
	}
}

// ArchiveResult 单次归档结果
type ArchiveResult struct {
	Migrated   int64
	RowsBefore int64
	RowsAfter  int64
	Trigger    string // "time" / "row_count" / "both" / "skip"
	DurationMs int64
	Error      string
}

// archiveOne 对单张表执行归档
func (s *ColdStorageService) archiveOne(cfg ArchiveConfig) ArchiveResult {
	start := time.Now()
	result := ArchiveResult{Trigger: "skip"}

	// 1. 获取当前行数
	rowCount, err := s.getRowCount(cfg.TableName)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.RowsBefore = rowCount

	// 2. 判断是否需要归档
	needArchive := false
	cutoffMs := dates.Timestamp(time.Now().AddDate(0, 0, -cfg.TimeDays))

	// 时间条件：是否存在超过 N 天的冷数据
	coldCount, _ := s.getColdCount(cfg, cutoffMs)

	if rowCount > int64(cfg.RowLimit) && coldCount > 0 {
		needArchive = true
		result.Trigger = "both"
	} else if rowCount > int64(cfg.RowLimit) {
		needArchive = true
		result.Trigger = "row_count"
	} else if coldCount > 0 {
		needArchive = true
		result.Trigger = "time"
	}

	if !needArchive {
		return result
	}

	// 3. 计算本次需要迁移多少行
	targetRows := int64(cfg.TargetRows)
	needMigrate := rowCount - targetRows
	if needMigrate <= 0 {
		// 只按时间归档
		needMigrate = coldCount
	}
	if coldCount < needMigrate {
		needMigrate = coldCount
	}

	// 4. 分批迁移
	var totalMigrated int64
	for totalMigrated < needMigrate {
		batchSize := int64(cfg.BatchSize)
		if needMigrate-totalMigrated < batchSize {
			batchSize = needMigrate - totalMigrated
		}

		migrated, err := s.archiveBatch(cfg, cutoffMs, int(batchSize))
		if err != nil {
			result.Error = err.Error()
			break
		}
		if migrated == 0 {
			break // 无更多冷数据
		}
		totalMigrated += int64(migrated)
		time.Sleep(500 * time.Millisecond) // 避免长事务锁表
	}

	result.Migrated = totalMigrated
	result.RowsAfter = rowCount - totalMigrated
	result.DurationMs = time.Since(start).Milliseconds()
	return result
}

// archiveBatch 单批次迁移
func (s *ColdStorageService) archiveBatch(cfg ArchiveConfig, cutoffMs int64, limit int) (int, error) {
	var migrated int

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		// 构建 WHERE 条件
		where := fmt.Sprintf("%s < ?", cfg.TimeColumn)
		if cfg.StatusColumn != "" {
			where += fmt.Sprintf(" AND %s != %d", cfg.StatusColumn, constants.StakeStatusActive)
		}

		// 选出待迁移的 ID
		var ids []int64
		query := tx.Table(cfg.TableName).
			Where(where, cutoffMs).
			Order(cfg.TimeColumn + " ASC").
			Limit(limit)

		if err := query.Pluck("id", &ids).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}

		// 插入归档表
		insertSQL := fmt.Sprintf(
			"INSERT INTO %s SELECT *, CURDATE() FROM %s WHERE id IN ?",
			cfg.ArchiveTable, cfg.TableName,
		)
		if err := tx.Exec(insertSQL, ids).Error; err != nil {
			return err
		}

		// 从主表删除
		deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE id IN ?", cfg.TableName)
		if err := tx.Exec(deleteSQL, ids).Error; err != nil {
			return err
		}

		migrated = len(ids)
		return nil
	})

	return migrated, err
}

// getRowCount 获取表当前行数
func (s *ColdStorageService) getRowCount(tableName string) (int64, error) {
	var count int64
	err := sqls.DB().Table(tableName).Count(&count).Error
	return count, err
}

// getColdCount 获取满足时间条件的冷数据数量
func (s *ColdStorageService) getColdCount(cfg ArchiveConfig, cutoffMs int64) (int64, error) {
	var count int64
	query := sqls.DB().Table(cfg.TableName).Where(cfg.TimeColumn+" < ?", cutoffMs)
	if cfg.StatusColumn != "" {
		query = query.Where(cfg.StatusColumn+" != ?", constants.StakeStatusActive)
	}
	err := query.Count(&count).Error
	return count, err
}
