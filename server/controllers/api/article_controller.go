package api

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/urls"
	"bbs-go/spam"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

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

// PostCreate 发表文章
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
	form := model.CreateArticleForm{
		Title:       title,
		Summary:     summary,
		Content:     content,
		ContentType: constants.ContentTypeMarkdown,
		Tags:        tags,
	}

	if err := spam.CheckArticle(user, form); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	article, err := services.ArticleService.Publish(user.Id, form)
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

// 用户文章列表
func (c *ArticleController) GetUserArticles() *simple.JsonResult {
	userId, err := simple.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	articles, cursor, hasMore := services.ArticleService.GetUserArticles(userId, cursor)
	return simple.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore)
}

// 文章列表
func (c *ArticleController) GetArticles() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	articles, cursor, hasMore := services.ArticleService.GetArticles(cursor)
	return simple.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore)
}

// 标签文章列表
func (c *ArticleController) GetTagArticles() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	tagId := simple.FormValueInt64Default(c.Ctx, "tagId", 0)
	articles, cursor, hasMore := services.ArticleService.GetTagArticles(tagId, cursor)
	return simple.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore)
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
