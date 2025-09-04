package admin

import (
	"strconv"

	"bbs-go/internal/models"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/pkg/simple/web"
	"bbs-go/internal/pkg/simple/web/params"
)

type CheckInController struct {
	Ctx iris.Context
}

func (c *CheckInController) GetBy(id int64) *web.JsonResult {
	t := services.CheckInService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *CheckInController) AnyList() *web.JsonResult {
	list, paging := services.CheckInService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *CheckInController) PostCreate() *web.JsonResult {
	t := &models.CheckIn{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = services.CheckInService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *CheckInController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.CheckInService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = services.CheckInService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
