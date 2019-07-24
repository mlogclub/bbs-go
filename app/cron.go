package app

import (
	"github.com/mlogclub/mlog/services"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

func startSchedule() {
	c := cron.New()

	// 自动生成sitemap和rss
	addCronFunc(c, "@every 10m", func() {
		services.ArticleService.GenerateSitemap()
		services.ArticleService.GenerateRss()

		services.TopicService.GenerateSitemap()
		services.TopicService.GenerateRss()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
