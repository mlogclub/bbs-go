package scheduler

import (
	"bbs-go/internal/pkg/sitemap"
	"log/slog"

	"github.com/robfig/cron/v3"

	"bbs-go/internal/services"
)

func Start() {
	c := cron.New()

	// Generate RSS
	addCronFunc(c, "@every 30m", func() {
		services.ArticleService.GenerateRss()
		services.TopicService.GenerateRss()
	})

	// Generate sitemap
	addCronFunc(c, "0 4 ? * *", func() {
		sitemap.Generate()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	_, err := c.AddFunc(sepc, cmd)
	if err != nil {
		slog.Error("add cron func error", slog.Any("err", err))
	} else {
		// slog.Info("add cron func", slog.Any("entryId", entryId))
	}
}
