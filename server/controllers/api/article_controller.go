package api

import (
	"bbs-go/model/constants"
	"bbs-go/package/urls"
	"math/rand"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type ArticleController struct {
	Ctx iris.Context
}

// 文章详情
func (c *ArticleController) GetBy(articleId int64) *simple.JsonResult {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return simple.JsonErrorCode(404, "文章不存在")
	}

	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		if article.UserId != user.Id && article.Status == constants.StatusPending {
			return simple.JsonErrorCode(403, "文章审核中")
		}
	} else {
		if article.Status == constants.StatusPending {
			return simple.JsonErrorCode(403, "文章审核中")
		}
	}

	services.ArticleService.IncrViewCount(articleId) // 增加浏览量
	return simple.JsonData(render.BuildArticle(article))
}

// 发表文章
func (c *ArticleController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}
	var (
		tags    = simple.FormValueStringArray(c.Ctx, "tags")
		title   = c.Ctx.PostValue("title")
		summary = c.Ctx.PostValue("summary")
		content = c.Ctx.PostValue("content")
	)

	article, err := services.ArticleService.Publish(user.Id, title, summary, content,
		constants.ContentTypeMarkdown, tags, "")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildArticle(article))
}

// 编辑时获取详情
func (c *ArticleController) GetEditBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return simple.JsonErrorMsg("话题不存在或已被删除")
	}

	// 非作者、且非管理员
	if article.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
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
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	var (
		tags    = simple.FormValueStringArray(c.Ctx, "tags")
		title   = c.Ctx.PostValue("title")
		content = c.Ctx.PostValue("content")
	)

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return simple.JsonErrorMsg("文章不存在")
	}

	// 非作者、且非管理员
	if article.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return simple.JsonErrorMsg("无权限")
	}

	if err := services.ArticleService.Edit(articleId, tags, title, content); err != nil {
		return simple.JsonError(err)
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityArticle, articleId,
		"", c.Ctx.Request())
	return simple.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 删除文章
func (c *ArticleController) PostDeleteBy(articleId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return simple.JsonErrorMsg("文章不存在")
	}

	// 非作者、且非管理员
	if article.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return simple.JsonErrorMsg("无权限")
	}

	if err := services.ArticleService.Delete(articleId); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeDelete, constants.EntityArticle, articleId,
		"", c.Ctx.Request())
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
	if article == nil || article.Status != constants.StatusOk {
		return simple.JsonErrorMsg("文章不存在")
	}
	return simple.NewEmptyRspBuilder().Put("url", urls.ArticleUrl(articleId)).JsonResult()
}

// 最近文章
func (c *ArticleController) GetRecent() *simple.JsonResult {
	articles := services.ArticleService.Find(simple.NewSqlCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(10))
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 用户文章列表
func (c *ArticleController) GetUserArticles() *simple.JsonResult {
	userId, err := simple.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	articles, cursor := services.ArticleService.GetUserArticles(userId, cursor)
	return simple.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10))
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
	articles := services.ArticleService.GetUserNewestArticles(userId)
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 近期文章
func (c *ArticleController) GetNearlyBy(articleId int64) *simple.JsonResult {
	articles := services.ArticleService.GetNearlyArticles(articleId)
	return simple.JsonData(render.BuildSimpleArticles(articles))
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
	articles := services.ArticleService.Find(simple.NewSqlCnd().Eq("status", constants.StatusOk).Desc("id").Limit(5))
	return simple.JsonData(render.BuildSimpleArticles(articles))
}

// 热门文章
func (c *ArticleController) GetHot() *simple.JsonResult {
	articles := cache.ArticleCache.GetHotArticles()
	return simple.JsonData(render.BuildSimpleArticles(articles))
}
