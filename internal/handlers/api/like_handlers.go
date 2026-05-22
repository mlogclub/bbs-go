package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/mlogclub/simple/common/strs"
)

func LikeLike(ctx *gin.Context) {
	var req req.EntityActionReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var (
		entityId = req.DecodedEntityId()
		user     = common.GetCurrentUser(ctx)
		err      error
	)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	switch req.EntityType {
	case constants.EntityTopic:
		err = services.UserLikeService.TopicLike(user.Id, entityId)
	case constants.EntityArticle:
		err = services.UserLikeService.ArticleLike(user.Id, entityId)
	case constants.EntityComment:
		err = services.UserLikeService.CommentLike(user.Id, entityId)
	}
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func LikeUnlike(ctx *gin.Context) {
	var req req.EntityActionReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var (
		entityId = req.DecodedEntityId()
		user     = common.GetCurrentUser(ctx)
		err      error
	)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	switch req.EntityType {
	case constants.EntityTopic:
		err = services.UserLikeService.TopicUnLike(user.Id, entityId)
	case constants.EntityArticle:
		err = services.UserLikeService.ArticleUnLike(user.Id, entityId)
	case constants.EntityComment:
		err = services.UserLikeService.CommentUnLike(user.Id, entityId)
	}
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func LikeLikedIds(ctx *gin.Context) {
	var req req.LikedIdsReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var (
		user           = common.GetCurrentUser(ctx)
		entityIds      = req.ParsedEntityIds()
		likedEntityIds []int64
	)
	if user != nil {
		likedEntityIds = services.UserLikeService.IsLiked(user.Id, req.EntityType, entityIds)
	}
	ginx.WriteJSON(ctx, likedEntityIds)

}

func LikeLiked(ctx *gin.Context) {
	var req req.EntityActionReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var (
		user     = common.GetCurrentUser(ctx)
		entityId = req.DecodedEntityId()
	)
	if user == nil || strs.IsBlank(req.EntityType) || entityId <= 0 {
		ginx.WriteJSON(ctx, false)
		return
	} else {
		liked := services.UserLikeService.Exists(user.Id, req.EntityType, entityId)
		ginx.WriteJSON(ctx, liked)
		return
	}

}
