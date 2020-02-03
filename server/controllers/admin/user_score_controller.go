package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type UserScoreController struct {
	Ctx iris.Context
}

func (c *UserScoreController) GetBy(id int64) *simple.JsonResult {
	t := services.UserScoreService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *UserScoreController) AnyList() *simple.JsonResult {
	list, paging := services.UserScoreService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *UserScoreController) PostCreate() *simple.JsonResult {
	t := &model.UserScore{}
	err := simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.UserScoreService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *UserScoreController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.UserScoreService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.UserScoreService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
