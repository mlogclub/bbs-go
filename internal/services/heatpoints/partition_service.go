package heatpoints

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// PartitionService 分区管理服务
type PartitionService struct{}

var Partition = &PartitionService{}

// CreateNextPartition 创建下个半年的分区
func (s *PartitionService) CreateNextPartition() error {
	slog.Info("检查是否需要创建新分区")
	now := time.Now()
	nextYear := now.Year()
	nextHalf := 1
	if now.Month() > 6 {
		nextHalf = 2
		nextYear++
	} else {
		nextHalf = 2
	}

	partitionName := fmt.Sprintf("p%dh%d", nextYear, nextHalf)
	exists, _ := s.partitionExists(partitionName)
	if exists {
		slog.Info("分区已存在", "partition", partitionName)
		return nil
	}

	return s.addPartition(nextYear, nextHalf)
}

func (s *PartitionService) partitionExists(name string) (bool, error) {
	var count int64
	sql := `SELECT COUNT(*) FROM information_schema.partitions WHERE table_name = 't_topic_stake' AND partition_name = ?`
	if err := sqls.DB().Raw(sql, name).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *PartitionService) addPartition(year, half int) error {
	var boundary int64
	if half == 1 {
		boundary = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli() / 100000000
	} else {
		boundary = time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC).UnixMilli() / 100000000
	}

	sql := fmt.Sprintf(`ALTER TABLE t_topic_stake ADD PARTITION (PARTITION p%dh%d VALUES LESS THAN (%d))`, year, half, boundary)
	return sqls.DB().Exec(sql).Error
}

// ArchiveOldPartition 归档旧分区
func (s *PartitionService) ArchiveOldPartition(years int) error {
	if years <= 0 {
		years = 2
	}
	slog.Info("开始归档旧分区", "years_to_keep", years)
	cutoffTime := time.Now().AddDate(-years, 0, 0)
	cutoffValue := cutoffTime.UnixMilli() / 100000000

	var partitions []string
	sql := `SELECT partition_name FROM information_schema.partitions WHERE table_name = 't_topic_stake' AND partition_name != 'pmax' AND partition_description < ? ORDER BY partition_name`
	if err := sqls.DB().Raw(sql, cutoffValue).Pluck("partition_name", &partitions).Error; err != nil {
		return err
	}

	for _, p := range partitions {
		if err := s.archiveSinglePartition(p); err != nil {
			slog.Error("归档分区失败", "partition", p, "error", err)
		}
	}
	return nil
}

func (s *PartitionService) archiveSinglePartition(name string) error {
	slog.Info("归档分区", "partition", name)
	backup := fmt.Sprintf("t_topic_stake_archive_%s", name)
	if err := sqls.DB().Exec(fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM t_topic_stake PARTITION (%s)", backup, name)).Error; err != nil {
		return err
	}
	if err := sqls.DB().Exec(fmt.Sprintf("ALTER TABLE t_topic_stake DROP PARTITION %s", name)).Error; err != nil {
		sqls.DB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", backup))
		return err
	}
	return nil
}

// PartitionStats 分区统计
type PartitionStats struct {
	PartitionName string
	TableRows     int64
	DataLength    int64
}

// GetPartitionStats 获取统计
func (s *PartitionService) GetPartitionStats() ([]PartitionStats, error) {
	var stats []PartitionStats
	sql := `SELECT partition_name, table_rows, data_length FROM information_schema.partitions WHERE table_name = 't_topic_stake' ORDER BY partition_name`
	if err := sqls.DB().Raw(sql).Scan(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

// CheckAndMigrate 自动管理
func (s *PartitionService) CheckAndMigrate() error {
	if err := s.CreateNextPartition(); err != nil {
		return err
	}
	if time.Now().Weekday() == time.Monday {
		if err := s.ArchiveOldPartition(2); err != nil {
			slog.Error("归档失败", "error", err)
		}
	}
	return nil
}

// 以下函数用于迁移脚本调用

// MigrateCreatePartitionTable 迁移入口函数
func MigrateCreatePartitionTable() error {
	slog.Info("开始执行分区表迁移")
	db := sqls.DB()

	// 检查数据库类型
	dialect := db.Dialector.Name()
	if dialect == "sqlite" {
		slog.Info("SQLite detected, skipping partition migration")
		return nil
	}

	// 检查表是否存在
	if !db.Migrator().HasTable(&TopicStakeForMigration{}) {
		slog.Info("表不存在，由 GORM 自动创建")
		return nil
	}

	var count int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 't_topic_stake' AND create_options LIKE '%partitioned%'`).Scan(&count).Error; err != nil {
		slog.Warn("检查分区状态失败", "error", err)
	}
	if count > 0 {
		slog.Info("已是分区表，跳过")
		return nil
	}

	var rows int64
	if err := db.Table("t_topic_stake").Count(&rows).Error; err != nil {
		return fmt.Errorf("检查数据量失败：%w", err)
	}

	if rows > 0 {
		return recreateWithPartition(db, rows)
	}
	return createEmptyPartitionedTable(db)
}

// TopicStakeForMigration 用于检查表结构
type TopicStakeForMigration struct {
	Id             int64
	TopicId        int64
	UserId         int64
	HeatPoints     int
	OriginalPoints int
	StakeDay       string
	Status         int
	LastSettleDay  string
	CreateTime     int64
	UpdateTime     int64
}

func createEmptyPartitionedTable(db *gorm.DB) error {
	if err := db.Exec("DROP TABLE IF EXISTS t_topic_stake").Error; err != nil {
		return err
	}

	// 动态计算分区边界值
	// 分区键 = update_time / 100000000（update_time 为毫秒时间戳）
	halfBoundary := func(year int, month time.Month) int64 {
		return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).UnixMilli() / 100000000
	}

	sql := fmt.Sprintf(`
CREATE TABLE t_topic_stake (
    id BIGINT NOT NULL AUTO_INCREMENT,
    topic_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    heat_points INT NOT NULL,
    original_points INT NOT NULL,
    stake_day VARCHAR(8) NOT NULL,
    status INT NOT NULL,
    last_settle_day VARCHAR(8),
    create_time BIGINT NOT NULL,
    update_time BIGINT NOT NULL,
    PRIMARY KEY (id, update_time),
    INDEX idx_stake_topic (topic_id),
    INDEX idx_stake_user (user_id),
    INDEX idx_stake_status (status),
    INDEX idx_stake_day (stake_day)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY RANGE (update_time / 100000000) (
    PARTITION p2026h1 VALUES LESS THAN (%d),
    PARTITION p2026h2 VALUES LESS THAN (%d),
    PARTITION p2027h1 VALUES LESS THAN (%d),
    PARTITION p2027h2 VALUES LESS THAN (%d),
    PARTITION p2028h1 VALUES LESS THAN (%d),
    PARTITION p2028h2 VALUES LESS THAN (%d),
    PARTITION pmax VALUES LESS THAN MAXVALUE
)`,
		halfBoundary(2026, time.July),
		halfBoundary(2027, time.January),
		halfBoundary(2027, time.July),
		halfBoundary(2028, time.January),
		halfBoundary(2028, time.July),
		halfBoundary(2029, time.January),
	)
	return db.Exec(sql).Error
}

func recreateWithPartition(db *gorm.DB, rows int64) error {
	slog.Info("重建表并迁移数据", "rows", rows)
	if err := db.Exec("CREATE TABLE t_topic_stake_backup AS SELECT * FROM t_topic_stake").Error; err != nil {
		return err
	}
	if err := db.Exec("DROP TABLE t_topic_stake").Error; err != nil {
		return err
	}
	if err := createEmptyPartitionedTable(db); err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO t_topic_stake SELECT * FROM t_topic_stake_backup").Error; err != nil {
		return err
	}
	if err := db.Exec("DROP TABLE t_topic_stake_backup").Error; err != nil {
		slog.Warn("清理备份表失败", "error", err)
	}
	return nil
}
