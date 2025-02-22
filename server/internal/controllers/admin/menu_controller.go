package admin

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type MenuController struct {
	Ctx iris.Context
}

func (c *MenuController) GetBy(id int64) *web.JsonResult {
	t := services.MenuService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}

	apis := services.ApiService.GetByMenuId(id)

	menu := render.BuildMenu(t)
	b := web.NewRspBuilder(menu)
	b.Put("apis", apis)
	return b.JsonResult()
}

func (c *MenuController) GetTree() *web.JsonResult {
	list := services.MenuService.Find(params.NewSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "status",
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(render.BuildMenuSimpleTree(0, list))
}

func (c *MenuController) AnyList() *web.JsonResult {
	list := services.MenuService.Find(params.NewSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "status",
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(render.BuildMenuTree(0, list))
}

func (c *MenuController) PostCreate() *web.JsonResult {
	t := &models.Menu{}
	params.ReadForm(c.Ctx, t)

	if t.SortNo <= 0 {
		t.SortNo = services.MenuService.GetNextSortNo(t.ParentId)
	}
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.MenuService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	// 设置接口权限
	var apis []models.Api
	jsons.Parse(params.FormValue(c.Ctx, "apis"), &apis)
	services.MenuApiService.SetMenuApis2(t.Id, apis)

	return web.JsonData(t)
}

func (c *MenuController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.MenuService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}
	params.ReadForm(c.Ctx, t)
	t.UpdateTime = dates.NowTimestamp()
	if err := services.MenuService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	// 设置接口权限
	var apis []models.Api
	jsons.Parse(params.FormValue(c.Ctx, "apis"), &apis)
	services.MenuApiService.SetMenuApis2(id, apis)

	return web.JsonData(t)
}

func (c *MenuController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		services.MenuService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": dates.NowTimestamp(),
		})
	}
	return web.JsonSuccess()
}

func (c *MenuController) GetUser_menus() *web.JsonResult {
	user, err := services.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	list := services.MenuService.GetUserMenus(user)
	return web.JsonData(render.BuildMenuTree(0, list))
}

func (c *MenuController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := services.MenuService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
