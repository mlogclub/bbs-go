package admin

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type AdvertController struct {
	Ctx iris.Context
}

func (c *AdvertController) GetBy(id int64) *web.JsonResult {
	advert := services.AdvertService.Get(id)
	if advert == nil || advert.Status == constants.StatusDeleted {
		return web.JsonErrorMsg("数据不存在")
	}
	return web.JsonData(c.buildLink(*advert))
}

func (c *AdvertController) AnyList() *web.JsonResult {
	list, paging := services.AdvertService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("status").LikeByReq("title").LikeByReq("url").PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *AdvertController) PostCreate() *web.JsonResult {
	t := &model.Advert{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}
	t.CreateTime = dates.NowTimestamp()

	err = services.AdvertService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *AdvertController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.AdvertService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = services.AdvertService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *AdvertController) buildLink(link model.Advert) map[string]interface{} {
	return map[string]interface{}{
		"id":         link.Id,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"picUrl":     link.PicUrl,
		"expireTime": link.ExpireTime,
	}
}
