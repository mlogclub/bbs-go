package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"

	"bbs-go/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *mvc.JsonResult {
	config := services.SysConfigService.GetConfig()
	return mvc.JsonData(config)
}
