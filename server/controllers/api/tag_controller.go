package api

import (
	"bbs-go/model/constants"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/services"
)

type TagController struct {
	Ctx iris.Context
}

// 标签详情
func (c *TagController) GetBy(tagId int64) *mvc.JsonResult {
	tag := cache.TagCache.Get(tagId)
	if tag == nil {
		return mvc.JsonErrorMsg("标签不存在")
	}
	return mvc.JsonData(render.BuildTag(tag))
}

// 标签列表
func (c *TagController) GetTags() *mvc.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)
	tags, paging := services.TagService.FindPageByCnd(sqls.NewSqlCnd().
		Eq("status", constants.StatusOk).
		Page(page, 200).Desc("id"))

	return mvc.JsonPageData(render.BuildTags(tags), paging)
}

// 标签自动完成
func (c *TagController) PostAutocomplete() *mvc.JsonResult {
	input := c.Ctx.FormValue("input")
	tags := services.TagService.Autocomplete(input)
	return mvc.JsonData(tags)
}
