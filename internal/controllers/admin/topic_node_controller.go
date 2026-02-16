package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
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
	list := services.TopicNodeService.Find(params.NewSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(list)
}

func (c *TopicNodeController) PostCreate() *web.JsonResult {
	t := &models.TopicNode{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}
	t.SortNo = services.TopicNodeService.GetNextSortNo()
	t.CreateTime = dates.NowTimestamp()
	err = services.TopicNodeService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = services.TopicNodeService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) GetNodes() *web.JsonResult {
	list := services.TopicNodeService.GetNodes()
	return web.JsonData(list)
}

func (c *TopicNodeController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := services.TopicNodeService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
