package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/model"
	"bbs-go/services"
)

type TopicNodeController struct {
	Ctx iris.Context
}

func (c *TopicNodeController) GetBy(id int64) *mvc.JsonResult {
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *TopicNodeController) AnyList() *mvc.JsonResult {
	list, paging := services.TopicNodeService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("name").PageByReq().Asc("sort_no").Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *TopicNodeController) PostCreate() *mvc.JsonResult {
	t := &model.TopicNode{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t.CreateTime = dates.NowTimestamp()
	err = services.TopicNodeService.Create(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *TopicNodeController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.TopicNodeService.Update(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *TopicNodeController) GetNodes() *mvc.JsonResult {
	list := services.TopicNodeService.GetNodes()
	return mvc.JsonData(list)
}
