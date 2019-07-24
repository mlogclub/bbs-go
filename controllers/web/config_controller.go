package web

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/simple"
	"strings"
)

type ConfigController struct {
	Ctx iris.Context
}

func (this *ConfigController) GetRecommendtags() *simple.JsonResult {
	var ret []string
	value := cache.SysConfigCache.GetValue(model.SysConfigRecommendTags)
	value = strings.TrimSpace(value)
	if len(value) > 0 {
		ss := strings.Split(value, ",")
		if len(ss) == 0 {
			return nil
		}
		for _, v := range ss {
			ret = append(ret, strings.TrimSpace(v))
		}
	}
	return simple.JsonData(ret)
}
