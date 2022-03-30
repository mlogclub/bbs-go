package admin

import (
	"bbs-go/model"
	"bbs-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"
)

type CheckInController struct {
	Ctx iris.Context
}

func (c *CheckInController) GetBy(id int64) *mvc.JsonResult {
	t := services.CheckInService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *CheckInController) AnyList() *mvc.JsonResult {
	list, paging := services.CheckInService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *CheckInController) PostCreate() *mvc.JsonResult {
	t := &model.CheckIn{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.CheckInService.Create(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *CheckInController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t := services.CheckInService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.CheckInService.Update(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}
