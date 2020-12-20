package api

import (
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type LikeController struct {
	Ctx iris.Context
}

func (c *LikeController) GetIsLiked() *simple.JsonResult {
	var (
		user           = services.UserTokenService.GetCurrent(c.Ctx)
		entityType     = simple.FormValue(c.Ctx, "entityType")
		entityIds      = simple.FormValueInt64Array(c.Ctx, "entityIds")
		likedEntityIds []int64
	)
	if user != nil {
		likedEntityIds = services.UserLikeService.IsLiked(user.Id, entityType, entityIds)
	}
	return simple.NewEmptyRspBuilder().Put("liked", likedEntityIds).JsonResult()
}

func (c *LikeController) GetLiked() *simple.JsonResult {
	var (
		user       = services.UserTokenService.GetCurrent(c.Ctx)
		entityType = simple.FormValue(c.Ctx, "entityType")
		entityId   = simple.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil || simple.IsBlank(entityType) || entityId <= 0 {
		return simple.NewEmptyRspBuilder().Put("liked", false).JsonResult()
	} else {
		liked := services.UserLikeService.Exists(user.Id, entityType, entityId)
		return simple.NewEmptyRspBuilder().Put("liked", liked).JsonResult()
	}
}
