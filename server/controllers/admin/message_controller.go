package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type MessageController struct {
	Ctx iris.Context
}

func (c *MessageController) GetBy(id int64) *simple.JsonResult {
	t := services.MessageService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *MessageController) AnyList() *simple.JsonResult {
	list, paging := services.MessageService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *MessageController) PostCreate() *simple.JsonResult {
	t := &model.Message{}
	simple.ReadForm(c.Ctx, t)

	err := services.MessageService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *MessageController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.MessageService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	simple.ReadForm(c.Ctx, t)

	err = services.MessageService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
