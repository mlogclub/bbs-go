package admin

import (
	"bbs-go/model"
	"bbs-go/services"
	"github.com/mlogclub/simple"
	"github.com/kataras/iris/v12"
	"strconv"
)

type UserScoreController struct {
	Ctx             iris.Context
}

func (r *UserScoreController) GetBy(id int64) *simple.JsonResult {
	t := services.UserScoreService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (r *UserScoreController) AnyList() *simple.JsonResult {
	list, paging := services.UserScoreService.FindPageByParams(simple.NewQueryParams(r.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (r *UserScoreController) PostCreate() *simple.JsonResult {
	t := &model.UserScore{}
	err := r.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.UserScoreService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (r *UserScoreController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(r.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.UserScoreService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = r.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.UserScoreService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

