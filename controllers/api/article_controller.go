package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type ArticleController struct {
	Ctx context.Context
}

// 发表文章
func (this *ArticleController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	var (
		tags    = simple.FormValueStringArray(this.Ctx, "tags")
		title   = this.Ctx.PostValue("title")
		summary = this.Ctx.PostValue("summary")
		content = this.Ctx.PostValue("content")
	)

	article, err := services.ArticleService.Publish(user.Id, title, summary, content,
		model.ArticleContentTypeMarkdown, 0, tags, "", false)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildArticle(article))
}

// 文章详情
func (this *ArticleController) GetBy(articleId int64) *simple.JsonResult {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		return simple.JsonErrorMsg("文章不存在")
	}
	return simple.JsonData(render.BuildArticle(article))
}

// 编辑文章
func (this *ArticleController) PostEdit() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	var (
		articleId = this.Ctx.PostValueInt64Default("id", 0)
		tags      = simple.FormValueStringArray(this.Ctx, "tags")
		title     = this.Ctx.PostValue("title")
		content   = this.Ctx.PostValue("content")
	)

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		return simple.JsonErrorMsg("文章不存在")
	}

	if article.UserId != user.Id {
		return simple.JsonErrorMsg("无权限")
	}

	err := services.ArticleService.Edit(articleId, tags, title, content)
	if err != nil {
		return simple.JsonError(err)
	}
	return simple.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 删除文章
func (this *ArticleController) PostDeleteBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		return simple.JsonErrorMsg("文章不存在")
	}

	if article.UserId != user.Id {
		return simple.JsonErrorMsg("无权限")
	}

	err := services.ArticleService.Delete(articleId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 收藏文章
func (this *ArticleController) PostFavoriteBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddArticleFavorite(user.Id, articleId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 最近文章
func (this *ArticleController) GetRecent() *simple.JsonResult {
	articles, err := services.ArticleService.QueryCnd(simple.NewQueryCnd("status = ?", model.ArticleStatusPublished).Order("id desc").Size(10))
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 用户最近的文章
func (this *ArticleController) GetUserRecent() *simple.JsonResult {
	userId, err := simple.FormValueInt64(this.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	articles, err := services.ArticleService.QueryCnd(simple.NewQueryCnd("user_id = ? and status = ?", userId, model.ArticleStatusPublished).Order("id desc").Size(10))
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 用户文章列表
func (this *ArticleController) GetUserArticles() *simple.JsonResult {
	userId, err := simple.FormValueInt64(this.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	articles, paging := services.ArticleService.Query(simple.NewParamQueries(this.Ctx).
		Eq("user_id", userId).
		Eq("status", model.ArticleStatusPublished).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(render.BuildSimpleArticles(articles), paging)
}

// 文章列表
func (this *ArticleController) GetArticles() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	articles, paging := services.ArticleService.Query(simple.NewParamQueries(this.Ctx).
		Eq("status", model.ArticleStatusPublished).
		NotEq("category_id", 4).
		Page(page, 20).Desc("id"))
	return simple.JsonPageData(render.BuildSimpleArticles(articles), paging)
}

// 标签文章列表
func (this *ArticleController) GetTagArticles() *simple.JsonResult {
	tagId := simple.FormValueInt64Default(this.Ctx, "tagId", 0)
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)
	articles, paging := services.ArticleService.GetTagArticles(tagId, page)
	return simple.JsonPageData(render.BuildSimpleArticles(articles), paging)
}

// 分类文章列表
func (this *ArticleController) GetCategoryArticles() *simple.JsonResult {
	categoryId := simple.FormValueInt64Default(this.Ctx, "categoryId", 0)
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	articles, paging := services.ArticleService.Query(simple.NewParamQueries(this.Ctx).
		Eq("category_id", categoryId).
		Eq("status", model.ArticleStatusPublished).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(render.BuildSimpleArticles(articles), paging)
}
