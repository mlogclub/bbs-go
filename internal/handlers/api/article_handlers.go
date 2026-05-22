package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/spam"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/jsons"

	"bbs-go/internal/handlers/render"
	"bbs-go/internal/services"
)

func bindArticleForm(ctx *gin.Context) (req.CreateArticleReq, error) {
	var body req.ArticleReq
	if err := ginx.Bind(ctx, &body); err != nil {
		return req.CreateArticleReq{}, err
	}
	return req.CreateArticleReq{
		Title:       body.Title,
		Summary:     body.Summary,
		Content:     body.Content,
		ContentType: constants.ContentTypeMarkdown,
		Cover:       req.ParseImageDTO(body.Cover),
		Tags:        body.ParsedTags(),
	}, nil
}

// 文章详情
// PostCreate 发表文章
// 编辑时获取详情
// 编辑文章
// 删除文章
// 收藏文章
// 文章跳转链接
// 用户文章列表
// 文章列表
// 标签文章列表
func ArticleDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	articleId := id

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		ginx.WriteJSON(ctx, ginx.ErrorCode(404, locales.Get("article.not_found")))
		return
	}

	user := common.GetCurrentUser(ctx)

	// 审核中文章控制展示
	if article.Status == constants.StatusReview {
		if user != nil {
			if article.UserId != user.Id && !user.IsOwner() {
				ginx.WriteJSON(ctx, ginx.ErrorCode(403, locales.Get("article.under_review")))
				return
			}
		} else {
			ginx.WriteJSON(ctx, ginx.ErrorCode(403, locales.Get("article.under_review")))
			return
		}
	}

	// 增加浏览量
	services.ArticleService.IncrViewCount(articleId)
	ginx.WriteJSON(ctx, render.BuildArticle(article, user))

}

func ArticleCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	form, err := bindArticleForm(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Content = strings.TrimSpace(form.Content)

	if err := spam.CheckArticle(user, form); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	article, err := services.ArticleService.Publish(user.Id, form)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, render.BuildArticle(article, user))

}

func ArticleEditForm(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	articleId := id

	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.not_found")))
		return
	}

	if !services.PermissionService.CanManageOwnedResource(user, article.UserId, permissions.PermissionArticleUpdate.Code) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.no_permission")))
		return
	}

	tags := services.ArticleService.GetArticleTags(articleId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	var cover *req.ImageDTO
	if err := jsons.Parse(article.Cover, &cover); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}

	ginx.WriteJSON(ctx, map[string]any{
		"id":        article.Id,
		"articleId": article.Id,
		"title":     article.Title,
		"content":   article.Content,
		"tags":      tagNames,
		"cover":     cover,
	})

}

func ArticleEdit(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	articleId := id

	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	form, err := bindArticleForm(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Content = strings.TrimSpace(form.Content)

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.not_found")))
		return
	}

	if !services.PermissionService.CanManageOwnedResource(user, article.UserId, permissions.PermissionArticleUpdate.Code) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.no_permission")))
		return
	}

	if err := services.ArticleService.Edit(articleId, form.Tags, form.Title, form.Content, form.Cover); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityArticle, articleId,
		"", ctx.Request)
	ginx.WriteJSON(ctx, map[string]any{"articleId": article.Id})

}

func ArticleRemove(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	articleId := id

	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.not_found")))
		return
	}

	if !services.PermissionService.CanManageOwnedResource(user, article.UserId, permissions.PermissionArticleDelete.Code) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.no_permission")))
		return
	}

	if err := services.ArticleService.Delete(articleId); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeDelete, constants.EntityArticle, articleId,
		"", ctx.Request)
	ginx.WriteJSON(ctx, nil)

}

func ArticleFavorite(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	articleId := id

	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	err = services.FavoriteService.AddArticleFavorite(user.Id, articleId)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func ArticleRedirect(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	articleId := id

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != constants.StatusOk {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.not_found")))
		return
	}
	ginx.WriteJSON(ctx, map[string]any{"url": bbsurls.ArticleUrl(articleId)})

}

func ArticleUserArticles(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	if userId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("param: userId required"))
		return
	}
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	articles, cursor, hasMore := services.ArticleService.GetUserArticles(userId, cursor)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore))

}

func ArticleArticles(ctx *gin.Context) {
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	articles, cursor, hasMore := services.ArticleService.GetArticles(cursor)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore))

}

func ArticleTagArticles(ctx *gin.Context) {
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	tagId := params.FormValueInt64Default(ctx, "tagId", 0)
	articles, cursor, hasMore := services.ArticleService.GetTagArticles(tagId, cursor)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore))

}
