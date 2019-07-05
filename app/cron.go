package app

import (
	"github.com/mlogclub/simple"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/mlogclub/mlog/services"
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

	// 码农日报
	addCronFunc(c, "0 0 1 * * *", func() {
		content := services.ArticleService.GetDailyContent([]int64{
			177, 79, 105, 115, 197, 88, 29, 171, 60, 53, 128, 143, 20, 190, 205, 288, 365, 326, 328, 222, 130, 68,
		})
		_, _ = services.TopicService.Publish(199, []string{"码农日报"}, "码农日报 "+simple.TimeFormat(time.Now(), simple.FMT_DATE), content)
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
