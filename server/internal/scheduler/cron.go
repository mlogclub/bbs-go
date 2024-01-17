package scheduler

import (
	"bbs-go/internal/pkg/sitemap"
	"log/slog"

	"github.com/robfig/cron"

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
	addCronFunc(c, "0 0 4 ? * *", func() {
		sitemap.Generate()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}
