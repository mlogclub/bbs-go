package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/model"
	"bbs-go/services"
)

type ProjectController struct {
	Ctx iris.Context
}

func (c *ProjectController) GetBy(id int64) *mvc.JsonResult {
	t := services.ProjectService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *ProjectController) AnyList() *mvc.JsonResult {
	list, paging := services.ProjectService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *ProjectController) PostCreate() *mvc.JsonResult {
	t := &model.Project{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.ProjectService.Create(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *ProjectController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t := services.ProjectService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.ProjectService.Update(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}
