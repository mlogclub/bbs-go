package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type BadgeController struct {
	Ctx iris.Context
}

func (c *BadgeController) GetBy(id int64) *web.JsonResult {
	t := services.BadgeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *BadgeController) AnyList() *web.JsonResult {
	list, paging := services.BadgeService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *BadgeController) PostCreate() *web.JsonResult {
	t := &models.Badge{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	now := dates.NowTimestamp()
	t.CreateTime = now
	t.UpdateTime = now
	if err := services.BadgeService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *BadgeController) PostUpdate() *web.JsonResult {
	id, _ := params.GetInt64(c.Ctx, "id")
	t := services.BadgeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.BadgeService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *BadgeController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	now := dates.NowTimestamp()
	for _, id := range ids {
		services.BadgeService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": now,
		})
	}
	return web.JsonSuccess()
}
