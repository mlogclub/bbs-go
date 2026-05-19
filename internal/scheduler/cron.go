package scheduler

import (
	"log/slog"

	"bbs-go/internal/services"

	"github.com/robfig/cron/v3"
)

func Start() {
	c := cron.New()

	addCronFunc(c, "0 4 ? * *", func() {
		if err := services.SeoSitemapService.GenerateAndUpload(); err != nil {
			slog.Error("generate sitemap error", slog.Any("err", err))
		}
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, spec string, cmd func()) {
	if _, err := c.AddFunc(spec, cmd); err != nil {
		slog.Error("add cron func error", slog.Any("err", err))
	}
}
