package admin

import (
	"bbs-go/model"
	"bbs-go/services"
	"github.com/mlogclub/simple"
	"github.com/kataras/iris/v12"
	"strconv"
)

type UserScoreLogController struct {
	Ctx             iris.Context
}

func (r *UserScoreLogController) GetBy(id int64) *simple.JsonResult {
	t := services.UserScoreLogService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (r *UserScoreLogController) AnyList() *simple.JsonResult {
	list, paging := services.UserScoreLogService.FindPageByParams(simple.NewQueryParams(r.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (r *UserScoreLogController) PostCreate() *simple.JsonResult {
	t := &model.UserScoreLog{}
	err := r.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.UserScoreLogService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (r *UserScoreLogController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(r.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.UserScoreLogService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = r.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.UserScoreLogService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

