package web

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/simple"
)

type TagController struct {
	Ctx iris.Context
}

func (this *TagController) PostAutocomplete() *simple.JsonResult {
	input := this.Ctx.FormValue("input")
	tags := services.TagService.Autocomplete(input)
	return simple.JsonData(tags)
}

// 所有标签列表
func GetTags(ctx iris.Context) {
	page := ctx.Params().GetIntDefault("page", 1)
	activeUsers := cache.UserCache.GetActiveUsers()

	tags, paging := services.TagService.Query(simple.NewParamQueries(ctx).
		Eq("status", model.TagStatusOk).
		Page(page, 200).Desc("id"))

	render.View(ctx, "tag/index.html", iris.Map{
		model.TplSiteTitle: "标签列表",
		"ActiveUsers":      render.BuildUsers(activeUsers),
		"Tags":             render.BuildTags(tags),
		"Page":             paging,
		"PrePageUrl":       utils.BuildTagsUrl(page - 1),
		"NextPageUrl":      utils.BuildTagsUrl(page + 1),
	})
}
