package admin

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type CommentController struct {
	Ctx            iris.Context
	CommentService *services.CommentService
}

func (this *CommentController) GetBy(id int64) *simple.JsonResult {
	t := this.CommentService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *CommentController) AnyList() *simple.JsonResult {
	list, paging := this.CommentService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *CommentController) PostCreate() *simple.JsonResult {
	t := &model.Comment{}
	this.Ctx.ReadForm(t)

	err := this.CommentService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *CommentController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := this.CommentService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = this.CommentService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
