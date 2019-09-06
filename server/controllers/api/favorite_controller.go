package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services"
)

type FavoriteController struct {
	Ctx context.Context
}

// 是否收藏了
func (this *FavoriteController) GetFavorited() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	entityType := simple.FormValue(this.Ctx, "entityType")
	entityId := simple.FormValueInt64Default(this.Ctx, "entityId", 0)
	if user == nil || len(entityType) == 0 || entityId <= 0 {
		return simple.NewEmptyRspBuilder().Put("favorited", false).JsonResult()
	} else {
		tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
		return simple.NewEmptyRspBuilder().Put("favorited", tmp != nil).JsonResult()
	}
}

// 取消收藏
func (this *FavoriteController) GetDelete() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	entityType := simple.FormValue(this.Ctx, "entityType")
	entityId := simple.FormValueInt64Default(this.Ctx, "entityId", 0)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
	if tmp != nil {
		services.FavoriteService.Delete(tmp.Id)
	}
	return simple.JsonSuccess()
}
