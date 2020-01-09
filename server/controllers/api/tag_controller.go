package api

import (
	"github.com/kataras/iris/v12"
	"bbs-go/simple"

	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
	"bbs-go/services/cache"
)

type TagController struct {
	Ctx iris.Context
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
		Eq("status", model.StatusOk).
		Page(page, 200).Desc("id"))

	return simple.JsonPageData(render.BuildTags(tags), paging)
}

// 标签自动完成
func (this *TagController) PostAutocomplete() *simple.JsonResult {
	input := this.Ctx.FormValue("input")
	tags := services.TagService.Autocomplete(input)
	return simple.JsonData(tags)
}
