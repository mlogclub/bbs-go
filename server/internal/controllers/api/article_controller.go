package api

import (
	"log/slog"
	"strconv"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/spam"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/pkg/simple/common/jsons"
	"bbs-go/internal/pkg/simple/web"
	"bbs-go/internal/pkg/simple/web/params"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type ArticleController struct {
	Ctx iris.Context
}

// 文章详情
func (c *ArticleController) GetBy(articleId int64) *web.JsonResult {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return web.JsonErrorCode(404, locales.Get("article.not_found"))
	}

	user := services.UserTokenService.GetCurrent(c.Ctx)

	// 审核中文章控制展示
	if article.Status == constants.StatusReview {
		if user != nil {
			if article.UserId != user.Id && !user.IsOwnerOrAdmin() {
				return web.JsonErrorCode(403, locales.Get("article.under_review"))
			}
		} else {
			return web.JsonErrorCode(403, locales.Get("article.under_review"))
		}
	}

	// 增加浏览量
	services.ArticleService.IncrViewCount(articleId)
	return web.JsonData(render.BuildArticle(article, user))
}

// PostCreate 发表文章
func (c *ArticleController) PostCreate() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := models.GetCreateArticleForm(c.Ctx)

	if err := spam.CheckArticle(user, form); err != nil {
		return web.JsonError(err)
	}

	article, err := services.ArticleService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(render.BuildArticle(article, user))
}

// 编辑时获取详情
func (c *ArticleController) GetEditBy(articleId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return web.JsonErrorMsg(locales.Get("article.not_found"))
	}

	// 非作者、且非管理员
	if article.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg(locales.Get("article.no_permission"))
	}

	tags := services.ArticleService.GetArticleTags(articleId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	var cover *models.ImageDTO
	if err := jsons.Parse(article.Cover, &cover); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}

	return web.NewEmptyRspBuilder().
		Put("id", article.Id).
		Put("articleId", article.Id).
		Put("title", article.Title).
		Put("content", article.Content).
		Put("tags", tagNames).
		Put("cover", cover).
		JsonResult()
}

// 编辑文章
func (c *ArticleController) PostEditBy(articleId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	var (
		tags    = params.FormValueStringArray(c.Ctx, "tags")
		title   = c.Ctx.PostValue("title")
		content = c.Ctx.PostValue("content")
		cover   = models.GetImageDTO(c.Ctx, "cover")
	)

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return web.JsonErrorMsg(locales.Get("article.not_found"))
	}

	// 非作者、且非管理员
	if article.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg(locales.Get("article.no_permission"))
	}

	if err := services.ArticleService.Edit(articleId, tags, title, content, cover); err != nil {
		return web.JsonError(err)
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityArticle, articleId,
		"", c.Ctx.Request())
	return web.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 删除文章
func (c *ArticleController) PostDeleteBy(articleId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		return web.JsonErrorMsg(locales.Get("article.not_found"))
	}

	// 非作者、且非管理员
	if article.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg(locales.Get("article.no_permission"))
	}

	if err := services.ArticleService.Delete(articleId); err != nil {
		return web.JsonError(err)
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeDelete, constants.EntityArticle, articleId,
		"", c.Ctx.Request())
	return web.JsonSuccess()
}

// 收藏文章
func (c *ArticleController) PostFavoriteBy(articleId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	err := services.FavoriteService.AddArticleFavorite(user.Id, articleId)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 文章跳转链接
func (c *ArticleController) GetRedirectBy(articleId int64) *web.JsonResult {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != constants.StatusOk {
		return web.JsonErrorMsg(locales.Get("article.not_found"))
	}
	return web.NewEmptyRspBuilder().Put("url", bbsurls.ArticleUrl(articleId)).JsonResult()
}

// 用户文章列表
func (c *ArticleController) GetUserArticles() *web.JsonResult {
	userId, err := params.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return web.JsonError(err)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	articles, cursor, hasMore := services.ArticleService.GetUserArticles(userId, cursor)
	return web.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore)
}

// 文章列表
func (c *ArticleController) GetArticles() *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	articles, cursor, hasMore := services.ArticleService.GetArticles(cursor)
	return web.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore)
}

// 标签文章列表
func (c *ArticleController) GetTagArticles() *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	tagId := params.FormValueInt64Default(c.Ctx, "tagId", 0)
	articles, cursor, hasMore := services.ArticleService.GetTagArticles(tagId, cursor)
	return web.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore)
}
