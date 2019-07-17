package web

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/mlog/utils/session"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type TagController struct {
	Ctx iris.Context
}

// 添加标签页面
func (this *TagController) GetAdd() {
	user := session.GetCurrentUser(this.Ctx)
	// 必须要求登录
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirect=/tag/add", iris.StatusSeeOther)
		return
	}

	render.View(this.Ctx, "/tag/add_tag.html", nil)
}

// 添加标签
func (this *TagController) PostAdd() {
	user := session.GetCurrentUser(this.Ctx)
	// 必须要求登录
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirect=/tag/add", iris.StatusSeeOther)
		return
	}

	name := strings.TrimSpace(simple.FormValue(this.Ctx, "name")) // 标签名称
	if len(name) == 0 {
		render.View(this.Ctx, "/tag/add_tag.html", iris.Map{
			"ErrMsg": "标签名称不能为空",
		})
		return
	}

	err := services.UserArticleTagService.AddUserTag(user.Id, name)
	if err != nil {
		logrus.Error(err)
		render.View(this.Ctx, "/tag/add_tag.html", iris.Map{
			"ErrMsg": err.Error(),
			"Name":   name,
		})
		return
	}
	this.Ctx.Redirect("/user/"+strconv.FormatInt(user.Id, 10)+"/tags", iris.StatusSeeOther)
}

// 用户标签列表
func (this *TagController) GetUsertags() *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	tags := services.UserArticleTagService.GetUserTags(user.Id)
	return simple.JsonData(render.BuildTags(tags))
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
