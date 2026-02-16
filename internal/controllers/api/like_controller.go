package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type LikeController struct {
	Ctx iris.Context
}

func (c *LikeController) PostLike() *web.JsonResult {
	var (
		entityType, _ = params.Get(c.Ctx, "entityType")
		entityId      = common.GetID(c.Ctx, "entityId")
		user          = common.GetCurrentUser(c.Ctx)
		err           error
	)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	switch entityType {
	case constants.EntityTopic:
		err = services.UserLikeService.TopicLike(user.Id, entityId)
	case constants.EntityArticle:
		err = services.UserLikeService.ArticleLike(user.Id, entityId)
	case constants.EntityComment:
		err = services.UserLikeService.CommentLike(user.Id, entityId)
	}
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *LikeController) PostUnlike() *web.JsonResult {
	var (
		entityType, _ = params.Get(c.Ctx, "entityType")
		entityId      = common.GetID(c.Ctx, "entityId")
		user          = common.GetCurrentUser(c.Ctx)
		err           error
	)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	switch entityType {
	case constants.EntityTopic:
		err = services.UserLikeService.TopicUnLike(user.Id, entityId)
	case constants.EntityArticle:
		err = services.UserLikeService.ArticleUnLike(user.Id, entityId)
	case constants.EntityComment:
		err = services.UserLikeService.CommentUnLike(user.Id, entityId)
	}
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *LikeController) GetLiked_ids() *web.JsonResult {
	var (
		user           = common.GetCurrentUser(c.Ctx)
		entityType     = params.FormValue(c.Ctx, "entityType")
		entityIds      = params.FormValueInt64Array(c.Ctx, "entityIds")
		likedEntityIds []int64
	)
	if user != nil {
		likedEntityIds = services.UserLikeService.IsLiked(user.Id, entityType, entityIds)
	}
	return web.JsonData(likedEntityIds)
}

func (c *LikeController) GetLiked() *web.JsonResult {
	var (
		user          = common.GetCurrentUser(c.Ctx)
		entityType, _ = params.Get(c.Ctx, "entityType")
		entityId      = common.GetID(c.Ctx, "entityId")
	)
	if user == nil || strs.IsBlank(entityType) || entityId <= 0 {
		return web.JsonData(false)
	} else {
		liked := services.UserLikeService.Exists(user.Id, entityType, entityId)
		return web.JsonData(liked)
	}
}
