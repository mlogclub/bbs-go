package api

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

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
		siteNavs        = cache.SysConfigCache.GetValue(model.SysConfigSiteNavs)
		recommendTags   = cache.SysConfigCache.GetValue(model.SysConfigRecommendTags)
	)

	var siteKeywordsArr []string
	if err := simple.ParseJson(siteKeywords, &siteKeywordsArr); err != nil {
		logrus.Error("站点关键词数据错误", err)
	}

	var siteNavsArr []model.SiteNav
	if err := simple.ParseJson(siteNavs, &siteNavsArr); err != nil {
		logrus.Error("站点导航数据错误", err)
	}

	var recommendTagsArr []string
	if err := simple.ParseJson(recommendTags, &recommendTagsArr); err != nil {
		logrus.Error("推荐标签数据错误", err)
	}

	return simple.JsonData(iris.Map{
		"siteTitle":       siteTitle,
		"siteDescription": siteDescription,
		"siteKeywords":    siteKeywordsArr,
		"siteNavs":        siteNavsArr,
		"recommendTags":   recommendTagsArr,
	})
}
