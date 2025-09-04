package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/pkg/simple/common/dates"
	"bbs-go/internal/pkg/simple/web"
	"bbs-go/internal/pkg/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *web.JsonResult {
	t := services.LinkService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *LinkController) AnyList() *web.JsonResult {
	list, paging := services.LinkService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "url",
			Op:        params.Like,
		},
	))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *LinkController) PostCreate() *web.JsonResult {
	t := &models.Link{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}
	t.CreateTime = dates.NowTimestamp()

	err = services.LinkService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *LinkController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.LinkService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = services.LinkService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
