package api

import (
	"strings"

	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
)

type TagController struct {
	Ctx context.Context
}

func (this *TagController) PostAutocomplete() *simple.JsonResult {
	input := this.Ctx.FormValue("input")
	tags := services.TagService.Autocomplete(input)
	return simple.JsonData(tags)
}

func (this *TagController) GetRecommendtags() *simple.JsonResult {
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
