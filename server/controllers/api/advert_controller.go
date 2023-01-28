package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"server/services"
)

type AdvertController struct {
	Ctx iris.Context
}

func (c *AdvertController) AnyList() *web.JsonResult {
	list, paging := services.AdvertService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("status").LikeByReq("title").LikeByReq("url").PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}
