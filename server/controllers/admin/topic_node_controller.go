package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/model"
	"bbs-go/services"
)

type TopicNodeController struct {
	Ctx iris.Context
}

func (c *TopicNodeController) GetBy(id int64) *web.JsonResult {
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) AnyList() *web.JsonResult {
	list, paging := services.TopicNodeService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("name").PageByReq().Asc("sort_no").Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *TopicNodeController) PostCreate() *web.JsonResult {
	t := &model.TopicNode{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t.CreateTime = dates.NowTimestamp()
	err = services.TopicNodeService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.TopicNodeService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) GetNodes() *web.JsonResult {
	list := services.TopicNodeService.GetNodes()
	return web.JsonData(list)
}
