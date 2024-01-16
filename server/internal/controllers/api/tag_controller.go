package api

import (
	"bbs-go/internal/models/constants"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/cache"
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/services"
)

type TagController struct {
	Ctx iris.Context
}

// 标签详情
func (c *TagController) GetBy(tagId int64) *web.JsonResult {
	tag := cache.TagCache.Get(tagId)
	if tag == nil {
		return web.JsonErrorMsg("标签不存在")
	}
	return web.JsonData(render.BuildTag(tag))
}

// 标签列表
func (c *TagController) GetTags() *web.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)
	tags, paging := services.TagService.FindPageByCnd(sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Page(page, 200).Desc("id"))

	return web.JsonPageData(render.BuildTags(tags), paging)
}

// 标签自动完成
func (c *TagController) PostAutocomplete() *web.JsonResult {
	input := c.Ctx.FormValue("input")
	tags := services.TagService.Autocomplete(input)
	return web.JsonData(tags)
}
