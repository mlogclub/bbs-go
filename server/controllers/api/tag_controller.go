package api

import (
	"strings"

	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/cache"
)

type TagController struct {
	Ctx context.Context
}

// 标签详情
func (this *TagController) GetBy(tagId int64) *simple.JsonResult {
	tag := cache.TagCache.Get(tagId)
	if tag == nil {
		return simple.JsonErrorMsg("标签不存在")
	}
	return simple.JsonData(render.BuildTag(tag))
}

// 标签列表
func (this *TagController) GetTags() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)
	tags, paging := services.TagService.FindPageByCnd(simple.NewSqlCnd().
		Eq("status", model.TagStatusOk).
		Page(page, 200).Desc("id"))

	return simple.JsonPageData(render.BuildTags(tags), paging)
}

// 标签自动完成
func (this *TagController) PostAutocomplete() *simple.JsonResult {
	input := this.Ctx.FormValue("input")
	tags := services.TagService.Autocomplete(input)
	return simple.JsonData(tags)
}

// 推荐标签
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

// 活跃标签
func (this *TagController) GetActive() *simple.JsonResult {
	tags := cache.TagCache.GetActiveTags()
	return simple.JsonData(render.BuildTags(tags))
}
