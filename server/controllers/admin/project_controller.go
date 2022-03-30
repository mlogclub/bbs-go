package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/model"
	"bbs-go/services"
)

type ProjectController struct {
	Ctx iris.Context
}

func (c *ProjectController) GetBy(id int64) *web.JsonResult {
	t := services.ProjectService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ProjectController) AnyList() *web.JsonResult {
	list, paging := services.ProjectService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *ProjectController) PostCreate() *web.JsonResult {
	t := &model.Project{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.ProjectService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *ProjectController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.ProjectService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.ProjectService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}
