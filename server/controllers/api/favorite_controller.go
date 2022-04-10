package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/pkg/common"
	"bbs-go/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

// 是否收藏了
func (c *FavoriteController) GetFavorited() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	entityType := params.FormValue(c.Ctx, "entityType")
	entityId := params.FormValueInt64Default(c.Ctx, "entityId", 0)
	if user == nil || len(entityType) == 0 || entityId <= 0 {
		return web.NewEmptyRspBuilder().Put("favorited", false).JsonResult()
	} else {
		tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
		return web.NewEmptyRspBuilder().Put("favorited", tmp != nil).JsonResult()
	}
}

// 取消收藏
func (c *FavoriteController) GetDelete() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	entityType := params.FormValue(c.Ctx, "entityType")
	entityId := params.FormValueInt64Default(c.Ctx, "entityId", 0)
	if user == nil {
		return web.JsonError(common.ErrorNotLogin)
	}
	tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
	if tmp != nil {
		services.FavoriteService.Delete(tmp.Id)
	}
	return web.JsonSuccess()
}
