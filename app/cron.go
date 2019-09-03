package app

import (
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/services/task"
)

func startSchedule() {
	c := cron.New()
	addCronFunc(c, "@every 10m", task.SitemapTask)
	addCronFunc(c, "@every 1h", task.CollectStudyGoLangProjectTask)
	addCronFunc(c, "@every 1h", task.CollectOschinaProjectTask)
	addCronFunc(c, "@every 6h", task.BaiduPing)
	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
