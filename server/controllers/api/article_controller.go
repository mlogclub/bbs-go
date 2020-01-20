package api

import (
	"math/rand"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/common/urls"
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
	"bbs-go/services/cache"
)

type ArticleController struct {
	Ctx iris.Context
}

// 文章详情
func (c *ArticleController) GetBy(articleId int64) *simple.JsonResult {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != model.StatusOk {
		return simple.JsonErrorMsg("文章不存在")
	}
	services.ArticleService.IncrViewCount(articleId) // 增加浏览量
	return simple.JsonData(render.BuildArticle(article))
}

// 发表文章
func (c *ArticleController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	var (
		tags    = simple.FormValueStringArray(c.Ctx, "tags")
		title   = c.Ctx.PostValue("title")
		summary = c.Ctx.PostValue("summary")
		content = c.Ctx.PostValue("content")
	)

	article, err := services.ArticleService.Publish(user.Id, title, summary, content,
		model.ContentTypeMarkdown, tags, "", false)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildArticle(article))
}

// 编辑时获取详情
func (c *ArticleController) GetEditBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != model.StatusOk {
		return simple.JsonErrorMsg("话题不存在或已被删除")
	}
	if article.UserId != user.Id {
		return simple.JsonErrorMsg("无权限")
	}

	tags := services.ArticleService.GetArticleTags(articleId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return simple.NewEmptyRspBuilder().
		Put("articleId", article.Id).
		Put("title", article.Title).
		Put("content", article.Content).
		Put("tags", tagNames).
		JsonResult()
}

// 编辑文章
func (c *ArticleController) PostEditBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	var (
		tags    = simple.FormValueStringArray(c.Ctx, "tags")
		title   = c.Ctx.PostValue("title")
		content = c.Ctx.PostValue("content")
	)

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.StatusDeleted {
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
func (c *ArticleController) PostDeleteBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.StatusDeleted {
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
func (c *ArticleController) PostFavoriteBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddArticleFavorite(user.Id, articleId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 文章跳转链接
func (c *ArticleController) GetRedirectBy(articleId int64) *simple.JsonResult {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != model.StatusOk {
		return simple.JsonErrorMsg("文章不存在")
	}
	if article.Share && len(article.SourceUrl) > 0 {
		return simple.NewEmptyRspBuilder().Put("url", article.SourceUrl).JsonResult()
	} else {
		return simple.NewEmptyRspBuilder().Put("url", urls.ArticleUrl(articleId)).JsonResult()
	}
}

// 最近文章
func (c *ArticleController) GetRecent() *simple.JsonResult {
	articles := services.ArticleService.Find(simple.NewSqlCnd().Where("status = ?", model.StatusOk).Desc("id").Limit(10))
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 用户最近的文章
func (c *ArticleController) GetUserRecent() *simple.JsonResult {
	userId, err := simple.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	articles := services.ArticleService.Find(simple.NewSqlCnd().Where("user_id = ? and status = ?",
		userId, model.StatusOk).Desc("id").Limit(10))
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 用户文章列表
func (c *ArticleController) GetUserArticles() *simple.JsonResult {
	userId, err := simple.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	page := simple.FormValueIntDefault(c.Ctx, "page", 1)

	articles, paging := services.ArticleService.FindPageByCnd(simple.NewSqlCnd().
		Eq("user_id", userId).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(render.BuildSimpleArticles(articles), paging)
}

// 文章列表
func (c *ArticleController) GetArticles() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	articles, cursor := services.ArticleService.GetArticles(cursor)
	return simple.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10))
}

// 标签文章列表
func (c *ArticleController) GetTagArticles() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	tagId := simple.FormValueInt64Default(c.Ctx, "tagId", 0)
	articles, cursor := services.ArticleService.GetTagArticles(tagId, cursor)
	return simple.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10))
}

// 用户最新的文章
func (c *ArticleController) GetUserNewestBy(userId int64) *simple.JsonResult {
	newestArticles := services.ArticleService.GetUserNewestArticles(userId)
	return simple.JsonData(render.BuildSimpleArticles(newestArticles))
}

// 相关文章
func (c *ArticleController) GetRelatedBy(articleId int64) *simple.JsonResult {
	relatedArticles := services.ArticleService.GetRelatedArticles(articleId)
	return simple.JsonData(render.BuildSimpleArticles(relatedArticles))
}

// 推荐
func (c *ArticleController) GetRecommend() *simple.JsonResult {
	articles := cache.ArticleCache.GetRecommendArticles()
	if articles == nil || len(articles) == 0 {
		return simple.JsonSuccess()
	} else {
		dest := make([]model.Article, len(articles))
		perm := rand.Perm(len(articles))
		for i, v := range perm {
			dest[v] = articles[i]
		}
		end := 10
		if end > len(articles) {
			end = len(articles)
		}
		ret := dest[0:end]
		return simple.JsonData(render.BuildSimpleArticles(ret))
	}
}

// 最新文章
func (c *ArticleController) GetNewest() *simple.JsonResult {
	articles := services.ArticleService.Find(simple.NewSqlCnd().Eq("status", model.StatusOk).Desc("id").Limit(5))
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 热门文章
func (c *ArticleController) GetHot() *simple.JsonResult {
	articles := cache.ArticleCache.GetHotArticles()
	return simple.JsonData(render.BuildSimpleArticles(articles))
}
