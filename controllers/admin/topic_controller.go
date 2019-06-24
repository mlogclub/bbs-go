
package admin

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"github.com/kataras/iris"
	"strconv"
)

type TopicController struct {
	Ctx             iris.Context
	TopicService      *services.TopicService
}

func (this *TopicController) GetBy(id int64) *simple.JsonResult {
	t := this.TopicService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TopicController) AnyList() *simple.JsonResult {
	list, paging := this.TopicService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *TopicController) PostCreate() *simple.JsonResult {
	t := &model.Topic{}
	this.Ctx.ReadForm(t)

	err := this.TopicService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TopicController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := this.TopicService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = this.TopicService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

