package api

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services/cache"
)

type ConfigController struct {
	Ctx context.Context
}

func (this *ConfigController) GetConfigs() *simple.JsonResult {
	var (
		siteTitle       = cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
		siteDescription = cache.SysConfigCache.GetValue(model.SysConfigSiteDescription)
		siteKeywords    = cache.SysConfigCache.GetValue(model.SysConfigSiteKeywords)
	)

	return simple.JsonData(iris.Map{
		"siteTitle":       siteTitle,
		"siteDescription": siteDescription,
		"siteKeywords":    siteKeywords,
	})
}
