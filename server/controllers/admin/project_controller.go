package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type ProjectController struct {
	Ctx iris.Context
}

func (c *ProjectController) GetBy(id int64) *simple.JsonResult {
	t := services.ProjectService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *ProjectController) AnyList() *simple.JsonResult {
	list, paging := services.ProjectService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *ProjectController) PostCreate() *simple.JsonResult {
	t := &model.Project{}
	err := simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.ProjectService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *ProjectController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.ProjectService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.ProjectService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
