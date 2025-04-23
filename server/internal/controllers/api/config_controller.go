package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/pkg/config"
	"bbs-go/internal/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	cfg := config.Instance
	if cfg.Installed {
		sysConfig := services.SysConfigService.GetConfig()
		return web.NewRspBuilder(sysConfig).
			Put("installed", cfg.Installed).
			JsonResult()
	} else {
		return web.JsonData(map[string]any{
			"installed": cfg.Installed,
		})
	}
}
