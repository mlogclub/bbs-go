package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type ForbiddenWordController struct {
	Ctx iris.Context
}

func (c *ForbiddenWordController) GetBy(id int64) *web.JsonResult {
	t := services.ForbiddenWordService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ForbiddenWordController) AnyList() *web.JsonResult {
	list, paging := services.ForbiddenWordService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("type").LikeByReq("word").EqByReq("status").PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c ForbiddenWordController) PostDelete() *web.JsonResult {
	id, _ := params.FormValueInt64(c.Ctx, "id")
	services.ForbiddenWordService.Delete(id)
	return web.JsonSuccess()
}

func (c *ForbiddenWordController) PostCreate() *web.JsonResult {
	t := &models.ForbiddenWord{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t.CreateTime = dates.NowTimestamp()
	err = services.ForbiddenWordService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *ForbiddenWordController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.ForbiddenWordService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.ForbiddenWordService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}
