package scheduler

import (
	"log/slog"

	"bbs-go/internal/services"
	"bbs-go/internal/services/heatpoints"

	"github.com/robfig/cron/v3"
)

func Start() {
	c := cron.New()

	addCronFunc(c, "0 4 ? * *", func() {
		if err := services.SeoSitemapService.GenerateAndUpload(); err != nil {
			slog.Error("generate sitemap error", slog.Any("err", err))
		}
	})

	// 热度点系统定时任务
	// 每日 23:55 生成交互快照 + 流通快照 + 火焰偏移量
	addCronFunc(c, "55 23 * * *", func() {
		if err := heatpoints.HeatSnapshot.TakeAllSnapshots(); err != nil {
			slog.Error("heat snapshot error", slog.Any("err", err))
		}
	})

	// 每日 00:00 执行结算
	addCronFunc(c, "0 0 * * *", func() {
		if err := heatpoints.Settlement.SettleAll(); err != nil {
			slog.Error("heat settlement error", slog.Any("err", err))
		}
	})

	// 每周一 04:00 检查并管理分区（自动创建新分区 + 归档旧分区）
	addCronFunc(c, "0 4 * * 1", func() {
		if err := heatpoints.Partition.CheckAndMigrate(); err != nil {
			slog.Error("partition management error", slog.Any("err", err))
		}
	})

	// 每月 1 号 03:00 执行冷数据归档（主表超阈值时迁移到归档表）
	addCronFunc(c, "0 3 1 * *", func() {
		heatpoints.ColdStorage.RunAll()
	})

	// 每 6 小时重算 StakedInWindow 滑动窗口（修正累计偏差）
	addCronFunc(c, "0 */6 * * *", func() {
		if err := heatpoints.StakedWindow.RecalculateAll(); err != nil {
			slog.Error("staked window recalc error", slog.Any("err", err))
		}
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, spec string, cmd func()) {
	if _, err := c.AddFunc(spec, cmd); err != nil {
		slog.Error("add cron func error", slog.Any("err", err))
	}
}
