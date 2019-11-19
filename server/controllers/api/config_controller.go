package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (this *ConfigController) GetConfigs() *simple.JsonResult {
	config := services.SysConfigService.GetConfigResponse()
	return simple.JsonData(config)
}

// func (this *ConfigController) GetTest() *simple.JsonResult {
// 	go func() {
// 		task.SitemapTask()
// 	}()
// 	return simple.JsonSuccess()
// }
