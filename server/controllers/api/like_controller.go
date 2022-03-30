package api

import (
	"bbs-go/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
)

type LikeController struct {
	Ctx iris.Context
}

func (c *LikeController) GetIsLiked() *mvc.JsonResult {
	var (
		user           = services.UserTokenService.GetCurrent(c.Ctx)
		entityType     = params.FormValue(c.Ctx, "entityType")
		entityIds      = params.FormValueInt64Array(c.Ctx, "entityIds")
		likedEntityIds []int64
	)
	if user != nil {
		likedEntityIds = services.UserLikeService.IsLiked(user.Id, entityType, entityIds)
	}
	return mvc.NewEmptyRspBuilder().Put("liked", likedEntityIds).JsonResult()
}

func (c *LikeController) GetLiked() *mvc.JsonResult {
	var (
		user       = services.UserTokenService.GetCurrent(c.Ctx)
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil || strs.IsBlank(entityType) || entityId <= 0 {
		return mvc.NewEmptyRspBuilder().Put("liked", false).JsonResult()
	} else {
		liked := services.UserLikeService.Exists(user.Id, entityType, entityId)
		return mvc.NewEmptyRspBuilder().Put("liked", liked).JsonResult()
	}
}
