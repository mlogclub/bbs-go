package admin

import (
	"github.com/mlogclub/simple/date"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type TopicNodeController struct {
	Ctx iris.Context
}

func (c *TopicNodeController) GetBy(id int64) *simple.JsonResult {
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *TopicNodeController) AnyList() *simple.JsonResult {
	list, paging := services.TopicNodeService.FindPageByParams(simple.NewQueryParams(c.Ctx).EqByReq("name").PageByReq().Asc("sort_no").Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *TopicNodeController) PostCreate() *simple.JsonResult {
	t := &model.TopicNode{}
	err := simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t.CreateTime = date.NowTimestamp()
	err = services.TopicNodeService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *TopicNodeController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TopicNodeService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *TopicNodeController) GetNodes() *simple.JsonResult {
	list := services.TopicNodeService.GetNodes()
	return simple.JsonData(list)
}
