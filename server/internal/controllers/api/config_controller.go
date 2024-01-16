package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	config := services.SysConfigService.GetConfig()
	return web.JsonData(config)
}
