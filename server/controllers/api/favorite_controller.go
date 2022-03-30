package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"

	"bbs-go/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

// 是否收藏了
func (c *FavoriteController) GetFavorited() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	entityType := params.FormValue(c.Ctx, "entityType")
	entityId := params.FormValueInt64Default(c.Ctx, "entityId", 0)
	if user == nil || len(entityType) == 0 || entityId <= 0 {
		return mvc.NewEmptyRspBuilder().Put("favorited", false).JsonResult()
	} else {
		tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
		return mvc.NewEmptyRspBuilder().Put("favorited", tmp != nil).JsonResult()
	}
}

// 取消收藏
func (c *FavoriteController) GetDelete() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	entityType := params.FormValue(c.Ctx, "entityType")
	entityId := params.FormValueInt64Default(c.Ctx, "entityId", 0)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
	if tmp != nil {
		services.FavoriteService.Delete(tmp.Id)
	}
	return mvc.JsonSuccess()
}
