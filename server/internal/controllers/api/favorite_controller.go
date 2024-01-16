package api

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

func (c *FavoriteController) PostAdd() *web.JsonResult {
	var (
		user       = services.UserTokenService.GetCurrent(c.Ctx)
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	var err error
	if entityType == constants.EntityTopic {
		err = services.FavoriteService.AddTopicFavorite(user.Id, entityId)
	} else if entityType == constants.EntityArticle {
		err = services.FavoriteService.AddArticleFavorite(user.Id, entityId)
	} else {
		err = errors.New("unsupproted")
	}

	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 取消收藏
func (c *FavoriteController) PostDelete() *web.JsonResult {
	var (
		user       = services.UserTokenService.GetCurrent(c.Ctx)
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	tmp := services.FavoriteService.GetBy(user.Id, entityType, entityId)
	if tmp != nil {
		services.FavoriteService.Delete(tmp.Id)
	}
	return web.JsonSuccess()
}
