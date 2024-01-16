package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type UserReportController struct {
	Ctx iris.Context
}

func (c *UserReportController) GetBy(id int64) *web.JsonResult {
	t := services.UserReportService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserReportController) AnyList() *web.JsonResult {
	list, paging := services.UserReportService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *UserReportController) PostCreate() *web.JsonResult {
	t := &models.UserReport{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.UserReportService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *UserReportController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.UserReportService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.UserReportService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}
