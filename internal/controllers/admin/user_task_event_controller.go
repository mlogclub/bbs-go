package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type UserTaskEventController struct {
	Ctx iris.Context
}

func (c *UserTaskEventController) GetBy(id int64) *web.JsonResult {
	t := services.UserTaskEventService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserTaskEventController) AnyList() *web.JsonResult {
	list, paging := services.UserTaskEventService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *UserTaskEventController) PostCreate() *web.JsonResult {
	t := &models.UserTaskEvent{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if err := services.UserTaskEventService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *UserTaskEventController) PostUpdate() *web.JsonResult {
	id, _ := params.GetInt64(c.Ctx, "id")
	t := services.UserTaskEventService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if err := services.UserTaskEventService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *UserTaskEventController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		services.UserTaskEventService.Delete(id)
	}
	return web.JsonSuccess()
}

