package migrations

import (
	"bbs-go/internal/models"
	"fmt"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func migrate_heat_points_system() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 为 User 表添加 heat_points 字段
		if !tx.Migrator().HasColumn(&models.User{}, "heat_points") {
			if err := tx.Migrator().AddColumn(&models.User{}, "heat_points"); err != nil {
				return err
			}
		}

		// 2. 为 Topic 表添加 ever_viral 和 flame_locked_level 字段
		if !tx.Migrator().HasColumn(&models.Topic{}, "ever_viral") {
			if err := tx.Migrator().AddColumn(&models.Topic{}, "ever_viral"); err != nil {
				return err
			}
		}
		if !tx.Migrator().HasColumn(&models.Topic{}, "flame_locked_level") {
			if err := tx.Migrator().AddColumn(&models.Topic{}, "flame_locked_level"); err != nil {
				return err
			}
		}

		// 3. 创建热度点相关表
		autoMigrateTable := func(model interface{}, tableName string) error {
			if !tx.Migrator().HasTable(tableName) {
				return tx.Migrator().AutoMigrate(model)
			}
			return nil
		}

		// TopicStake 帖子质押记录表
		if err := autoMigrateTable(&models.TopicStake{}, "t_topic_stake"); err != nil {
			return err
		}

		// UserHeatLog 用户热度点流水表
		if err := autoMigrateTable(&models.UserHeatLog{}, "t_user_heat_log"); err != nil {
			return err
		}

		// UserHeatStats 用户热度点统计表
		if err := autoMigrateTable(&models.UserHeatStats{}, "t_user_heat_stats"); err != nil {
			return err
		}

		// HeatPublicPool 公共奖池记录表
		if err := autoMigrateTable(&models.HeatPublicPool{}, "t_heat_public_pool"); err != nil {
			return err
		}

		// SystemMintLog 系统铸币日志表
		if err := autoMigrateTable(&models.SystemMintLog{}, "t_system_mint_log"); err != nil {
			return err
		}

		// TopicInteractionSnapshot 每日互动快照表
		if err := autoMigrateTable(&models.TopicInteractionSnapshot{}, "t_topic_interaction_snapshot"); err != nil {
			return err
		}

		// HeatCirculationSnapshot 活跃流通快照表
		if err := autoMigrateTable(&models.HeatCirculationSnapshot{}, "t_heat_circulation_snapshot"); err != nil {
			return err
		}

		// DailyFlameOffset 每日火焰等级偏移量表
		if err := autoMigrateTable(&models.DailyFlameOffset{}, "t_daily_flame_offset"); err != nil {
			return err
		}

		// SettlementTaskLog 结算任务日志表
		if err := autoMigrateTable(&models.SettlementTaskLog{}, "t_settlement_task_log"); err != nil {
			return err
		}

		// 4. 初始化公共奖池余额为 0
		var poolCount int64
		if err := tx.Model(&models.HeatPublicPool{}).Where("source = ?", "InitialBalance").Count(&poolCount).Error; err != nil {
			return err
		}
		if poolCount == 0 {
			initialPool := models.HeatPublicPool{
				Source:       "InitialBalance",
				Amount:       0,
				BalanceAfter: 0,
				Remark:       "公共奖池初始余额",
				CreateTime:   dates.NowTimestamp(),
			}
			if err := tx.Create(&initialPool).Error; err != nil {
				return err
			}
		}

		// 5. 为所有现有正常用户初始化 UserHeatStats 记录（创世空投前置准备）
		var normalUsers []models.User
		if err := tx.Where("status = 0").Find(&normalUsers).Error; err != nil {
			return err
		}

		for _, user := range normalUsers {
			heatStats := models.UserHeatStats{
				UserId:             user.Id,
				TotalPoints:        0,
				StakedInWindow:     0,
				LastStakeTime:      0,
				DecayedAccumulated: 0,
				UpdateTime:         dates.NowTimestamp(),
			}
			if err := tx.Create(&heatStats).Error; err != nil {
				// 如果已存在则跳过（唯一键冲突）
				continue
			}
		}

		// 6. 创建冷存储归档表（结构与主表一致，额外增加 archive_date 字段）
		createArchiveIfNotExists := func(tx *gorm.DB, mainTable, archiveTable string) error {
			if tx.Migrator().HasTable(archiveTable) {
				return nil
			}
			// SQLite 不支持 CREATE TABLE LIKE，用 CREATE TABLE AS 替代
			return tx.Exec(fmt.Sprintf(
				"CREATE TABLE %s AS SELECT *, NULL as archive_date FROM %s WHERE 1=0",
				archiveTable, mainTable,
			)).Error
		}
		// 归档表：topic_stake → topic_stake_archive
		if err := createArchiveIfNotExists(tx, "t_topic_stake", "t_topic_stake_archive"); err != nil {
			return err
		}
		// 归档表：user_heat_log → user_heat_log_archive
		if err := createArchiveIfNotExists(tx, "t_user_heat_log", "t_user_heat_log_archive"); err != nil {
			return err
		}
		// 归档表：topic_interaction_snapshot → topic_interaction_snapshot_archive
		if err := createArchiveIfNotExists(tx, "t_topic_interaction_snapshot", "t_topic_interaction_snapshot_archive"); err != nil {
			return err
		}

		return nil
	})
}

