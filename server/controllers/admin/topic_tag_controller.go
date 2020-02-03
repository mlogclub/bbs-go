package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type TopicTagController struct {
	Ctx iris.Context
}

func (c *TopicTagController) GetBy(id int64) *simple.JsonResult {
	t := services.TopicTagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *TopicTagController) AnyList() *simple.JsonResult {
	list, paging := services.TopicTagService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *TopicTagController) PostCreate() *simple.JsonResult {
	t := &model.TopicTag{}
	simple.ReadForm(c.Ctx, t)

	err := services.TopicTagService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *TopicTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.TopicTagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	simple.ReadForm(c.Ctx, t)

	err = services.TopicTagService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
