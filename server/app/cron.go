package app

import (
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"bbs-go/services"
	"bbs-go/sitemap"
)

func startSchedule() {
	c := cron.New()

	// Generate RSS
	addCronFunc(c, "@every 30m", func() {
		services.ArticleService.GenerateRss()
		services.TopicService.GenerateRss()
		services.ProjectService.GenerateRss()
	})

	// Generate sitemap
	addCronFunc(c, "@every 2h", func() {
		sitemap.Generate()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
