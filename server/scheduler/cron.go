package scheduler

import (
	"bbs-go/pkg/sitemap"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"bbs-go/services"
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

	// vip 每日积分充值
	addCronFunc(c, "0 0 1 ? * *", func() {
		services.UserService.CronUserPayScore()
	})

	// 每天4点更新ES数据
	addCronFunc(c, "0 0 4 ? * *", func() {
		services.TopicService.UpdateEsIndex()
		services.ArticleService.UpdateEsIndex()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
