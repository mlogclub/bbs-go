package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/pkg/errs"
	"bbs-go/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

// 取消收藏
func (c *FavoriteController) GetDelete() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	entityType := params.FormValue(c.Ctx, "entityType")
	entityId := params.FormValueInt64Default(c.Ctx, "entityId", 0)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
	if tmp != nil {
		services.FavoriteService.Delete(tmp.Id)
	}
	return web.JsonSuccess()
}
