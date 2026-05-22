package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/ginx"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"
)

// 取消收藏
func FavoriteAdd(ctx *gin.Context) {
	var req req.EntityActionReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var (
		user     = common.GetCurrentUser(ctx)
		entityId = req.DecodedEntityId()
	)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	var err error
	switch req.EntityType {
	case constants.EntityTopic:
		err = services.FavoriteService.AddTopicFavorite(user.Id, entityId)
	case constants.EntityArticle:
		err = services.FavoriteService.AddArticleFavorite(user.Id, entityId)
	default:
		err = errors.New("unsupported")
	}

	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func FavoriteRemove(ctx *gin.Context) {
	var req req.EntityActionReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var (
		user     = common.GetCurrentUser(ctx)
		entityId = req.DecodedEntityId()
	)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	tmp := services.FavoriteService.GetBy(user.Id, req.EntityType, entityId)
	if tmp != nil {
		services.FavoriteService.Delete(tmp.Id)
	}
	ginx.WriteJSON(ctx, nil)

}
