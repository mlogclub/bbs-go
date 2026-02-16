package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type TopicTagController struct {
	Ctx iris.Context
}

func (c *TopicTagController) GetBy(id int64) *web.JsonResult {
	t := services.TopicTagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TopicTagController) AnyList() *web.JsonResult {
	list, paging := services.TopicTagService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *TopicTagController) PostCreate() *web.JsonResult {
	t := &models.TopicTag{}
	params.ReadForm(c.Ctx, t)

	err := services.TopicTagService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TopicTagController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.TopicTagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	params.ReadForm(c.Ctx, t)

	err = services.TopicTagService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
