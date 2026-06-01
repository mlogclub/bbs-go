package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/mlogclub/simple/common/strs"
)

// likeRateLimiter 点赞频率限制：同一用户对同一实体，2 秒内只能操作 1 次
var likeRateRecords sync.Map // key: "userId:entityType:entityId" → lastActionTime (int64 ms)

func checkLikeRate(userId int64, entityType string, entityId int64) bool {
	key := fmt.Sprintf("%d:%s:%d", userId, entityType, entityId)
	now := time.Now().UnixMilli()

	if val, ok := likeRateRecords.Load(key); ok {
		if now-val.(int64) < 2000 {
			return false
		}
	}

	likeRateRecords.Store(key, now)
	return true
}

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
	if !checkLikeRate(user.Id, req.EntityType, entityId) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("操作过于频繁，请稍后再试"))
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
	if !checkLikeRate(user.Id, req.EntityType, entityId) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("操作过于频繁，请稍后再试"))
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
