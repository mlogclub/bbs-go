package admin

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type TaskConfigController struct {
	Ctx iris.Context
}

func (c *TaskConfigController) GetBy(id int64) *web.JsonResult {
	t := services.TaskConfigService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TaskConfigController) GetGroups() *web.JsonResult {
	return web.JsonData(render.BuildTaskGroups())
}

func (c *TaskConfigController) AnyList() *web.JsonResult {
	list, paging := services.TaskConfigService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "groupName",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "eventType",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "period",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *TaskConfigController) PostCreate() *web.JsonResult {
	t := &models.TaskConfig{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	now := dates.NowTimestamp()
	t.CreateTime = now
	t.UpdateTime = now
	if err := services.TaskConfigService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *TaskConfigController) PostUpdate() *web.JsonResult {
	id, _ := params.GetInt64(c.Ctx, "id")
	t := services.TaskConfigService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.TaskConfigService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *TaskConfigController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	now := dates.NowTimestamp()
	for _, id := range ids {
		services.TaskConfigService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": now,
		})
	}
	return web.JsonSuccess()
}
