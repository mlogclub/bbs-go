package controllers

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type TagController struct {
	Ctx                   iris.Context
	CategoryService       *services.CategoryService
	TagService            *services.TagService
	ArticleService        *services.ArticleService
	UserArticleTagService *services.UserArticleTagService
}

// 添加标签页面
func (this *TagController) GetAdd() {
	user := utils.GetCurrentUser(this.Ctx)
	// 必须要求登录
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirect=/tag/add", iris.StatusSeeOther)
		return
	}

	render.View(this.Ctx, "/tag/add_tag.html", nil)
}

// 添加标签
func (this *TagController) PostAdd() {
	user := utils.GetCurrentUser(this.Ctx)
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

	err := this.UserArticleTagService.AddUserTag(user.Id, name)
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
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	tags := this.UserArticleTagService.GetUserTags(user.Id)
	return simple.JsonData(render.BuildTags(tags))
}

// 所有标签列表
func GetTags(ctx iris.Context) {
	page := ctx.Params().GetIntDefault("page", 1)
	activeUsers := cache.UserCache.GetActiveUsers()

	tags, paging := services.TagServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("status", model.TagStatusOk).
		Page(page, 200).Desc("id"))

	render.View(ctx, "tag/index.html", iris.Map{
		utils.GlobalFieldSiteTitle:       "标签列表",
		utils.GlobalFieldSiteDescription: "M-LOG分享",
		utils.GlobalFieldSiteKeywords:    "Go中国，Golang, Golang学习,Golang教程,Golang社区,Go基金会,Go语言中文网,Go,Go语言,主题,资源,文章,图书,开源项目,M-LOG",
		"ActiveUsers":                    render.BuildUsers(activeUsers),
		"Tags":                           render.BuildTags(tags),
		"Page":                           paging,
		"PrePageUrl":                     utils.BuildTagsUrl(page - 1),
		"NextPageUrl":                    utils.BuildTagsUrl(page + 1),
	})
}
