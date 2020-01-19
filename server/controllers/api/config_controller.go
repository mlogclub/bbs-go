package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
	"bbs-go/services/task"
)

type ConfigController struct {
	Ctx iris.Context
}

func (this *ConfigController) GetConfigs() *simple.JsonResult {
	config := services.SysConfigService.GetConfigResponse()
	return simple.JsonData(config)
}

func (this *ConfigController) GetTest() *simple.JsonResult {
	go func() {
		// services.SitemapService.GenerateMisc()
		// services.SitemapService.GenerateUser()
		//
		// now := time.Now()
		// dateFrom := time.Date(2002, 1, 1, 0, 0, 0, 0, now.Location())
		//
		// for {
		// 	if dateFrom.After(now) {
		// 		break
		// 	}
		// 	dateTo := dateFrom.Add(time.Hour * 24)
		// 	services.SitemapService.Generate(simple.Timestamp(dateFrom), simple.Timestamp(dateTo))
		//
		// 	dateFrom = dateFrom.Add(24 * time.Hour)
		// }

		task.SitemapTask()
	}()
	return simple.JsonSuccess()
}
