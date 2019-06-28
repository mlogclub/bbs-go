package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/mlogclub/mlog/utils/session"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils"
)

type UserTagController struct {
	Ctx                   iris.Context
	UserArticleTagService *services.UserArticleTagService
}

// 标签列表页面
func (this *UserTagController) GetList() mvc.View {
	user := session.GetCurrentUser(this.Ctx)
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)
	list, paging := this.UserArticleTagService.Query(simple.NewParamQueries(this.Ctx).Eq("user_id", user.Id).Page(page, 20).Desc("id"))

	var tags []model.Tag
	if len(list) > 0 {
		for _, v := range list {
			tags = append(tags, *cache.TagCache.Get(v.TagId))
		}
	}
	return mvc.View{
		Name: "user/tag/list.html",
		Data: iris.Map{
			"User":                     render.BuildUser(user),
			"Tags":                     render.BuildTags(tags),
			"Page":                     paging,
			utils.GlobalFieldSiteTitle: user.Nickname + " - 标签列表",
		},
	}
}

// 添加标签页面
func (this *UserTagController) GetAdd() {

}

// 添加标签
func (this *UserTagController) PostAdd() {

}
