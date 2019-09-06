package app

import (
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/services/task"
)

func startSchedule() {
	c := cron.New()
	addCronFunc(c, "@every 30m", task.RssTask)
	addCronFunc(c, "@every 1h", task.SitemapTask)
	addCronFunc(c, "@every 1h", task.CollectStudyGoLangProjectTask)
	addCronFunc(c, "@every 1h", task.CollectOschinaProjectTask)
	addCronFunc(c, "@every 12h", task.BaiduPing)
	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
