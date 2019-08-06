package api

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services/cache"
)

type ConfigController struct {
	Ctx context.Context
}

func (this *ConfigController) GetConfigs() *simple.JsonResult {
	var (
		siteTitle       = cache.SysConfigCache.Get(model.SysConfigSiteTitle)
		siteDescription = cache.SysConfigCache.Get(model.SysConfigSiteDescription)
		siteKeywords    = cache.SysConfigCache.Get(model.SysConfigSiteKeywords)
	)

	return simple.JsonData(iris.Map{
		"siteTitle":       siteTitle,
		"siteDescription": siteDescription,
		"siteKeywords":    siteKeywords,
	})

}
