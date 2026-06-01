package heatpoints

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"log/slog"
	"time"

	"github.com/mlogclub/simple/sqls"
)

// StakedWindowService StakedInWindow 滑动窗口维护服务
// 每 6 小时全量重算一次，修正因质押/赎回/结算导致的累计偏差
type StakedWindowService struct{}

var StakedWindow = &StakedWindowService{}

// RecalculateAll 全量重算所有用户的 StakedInWindow
// 逻辑：对每个用户，查询 topic_stakes 表中近 7 天内（按 create_time）所有已质押记录的 heat_points 之和
// 注意：包含已赎回但仍在 7 天窗口内的记录，因为"近 7 天累计质押量"包含所有质押行为
func (s *StakedWindowService) RecalculateAll() error {
	slog.Info("开始重算 StakedInWindow")

	cutoffMs := time.Now().AddDate(0, 0, -7).UnixMilli()

	// 一条 SQL 聚合每个用户近 7 天的累计质押量
	type userWindowRow struct {
		UserId        int64
		StakedInWindow int
	}
	var rows []userWindowRow
	err := sqls.DB().Model(&models.TopicStake{}).
		Where("create_time >= ?", cutoffMs).
		Group("user_id").
		Select("user_id, COALESCE(SUM(heat_points), 0) as staked_in_window").
		Scan(&rows).Error

	if err != nil {
		slog.Error("StakedInWindow 重算查询失败", "error", err)
		return err
	}

	windowMap := make(map[int64]int)
	for _, r := range rows {
		windowMap[r.UserId] = r.StakedInWindow
	}

	// 遍历所有 UserHeatStats，逐一更新
	var stats []models.UserHeatStats
	sqls.DB().Find(&stats)

	corrected := 0
	for _, stat := range stats {
		correctValue := windowMap[stat.UserId] // 无记录则为 0
		if stat.StakedInWindow != correctValue {
			sqls.DB().Model(&models.UserHeatStats{}).
				Where("user_id = ?", stat.UserId).
				Update("staked_in_window", correctValue)
			corrected++
		}
	}

	// 同时重算 TotalPoints（余额 + 活跃质押）
	s.syncAllTotalPoints()

	slog.Info("StakedInWindow 重算完成",
		"users", len(stats),
		"corrected", corrected,
	)
	return nil
}

// syncAllTotalPoints 全量同步所有用户的 TotalPoints
func (s *StakedWindowService) syncAllTotalPoints() {
	// 批量查询所有用户的活跃质押总额
	type userStakeRow struct {
		UserId int64
		Total  int
	}
	var stakeRows []userStakeRow
	sqls.DB().Model(&models.TopicStake{}).
		Where("status = ?", constants.StakeStatusActive).
		Group("user_id").
		Select("user_id, COALESCE(SUM(heat_points), 0) as total").
		Scan(&stakeRows)

	stakeMap := make(map[int64]int)
	for _, r := range stakeRows {
		stakeMap[r.UserId] = r.Total
	}

	// 批量查询所有用户的余额
	type userBalanceRow struct {
		Id         int64
		HeatPoints int
	}
	var users []userBalanceRow
	sqls.DB().Model(&models.User{}).
		Select("id, heat_points").
		Find(&users)

	balanceMap := make(map[int64]int)
	for _, u := range users {
		balanceMap[u.Id] = u.HeatPoints
	}

	// 逐一比对并更新
	var stats []models.UserHeatStats
	sqls.DB().Find(&stats)

	for _, stat := range stats {
		correctTotal := balanceMap[stat.UserId] + stakeMap[stat.UserId]
		if stat.TotalPoints != correctTotal {
			sqls.DB().Model(&models.UserHeatStats{}).
				Where("user_id = ?", stat.UserId).
				Update("total_points", correctTotal)
		}
	}
}

// RecalculateUser 重算单个用户的 StakedInWindow（质押/赎回时调用）
func (s *StakedWindowService) RecalculateUser(userId int64) {
	cutoffMs := time.Now().AddDate(0, 0, -7).UnixMilli()

	var total int
	sqls.DB().Model(&models.TopicStake{}).
		Where("user_id = ? AND create_time >= ?", userId, cutoffMs).
		Select("COALESCE(SUM(heat_points), 0)").
		Scan(&total)

	sqls.DB().Model(&models.UserHeatStats{}).
		Where("user_id = ?", userId).
		Update("staked_in_window", total)
}
