package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *simple.JsonResult {
	config := services.SysConfigService.GetConfig()
	return simple.JsonData(config)
}
