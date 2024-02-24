package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type RoleController struct {
	Ctx iris.Context
}

func (c *RoleController) GetBy(id int64) *web.JsonResult {
	t := services.RoleService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *RoleController) AnyList() *web.JsonResult {
	list, paging := services.RoleService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "status",
		},
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "code",
			Op:        params.Like,
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *RoleController) GetAll_roles() *web.JsonResult {
	list := services.RoleService.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
	return web.JsonData(list)
}

func (c *RoleController) PostCreate() *web.JsonResult {
	t := &models.Role{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if services.RoleService.GetByCode(t.Code) != nil {
		return web.JsonErrorMsg("角色编码已存在")
	}

	t.SortNo = services.RoleService.GetNextSortNo()
	t.Type = constants.RoleTypeCustom
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.RoleService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *RoleController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.RoleService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if t.Type == constants.RoleTypeSystem {
		return web.JsonErrorMsg("系统角色不允许编辑")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if exists := services.RoleService.GetByCode(t.Code); exists != nil && exists.Id != t.Id {
		return web.JsonErrorMsg("角色编码已存在")
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.RoleService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *RoleController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		services.RoleService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": dates.NowTimestamp(),
		})
	}
	return web.JsonSuccess()
}

func (s *RoleController) GetRoles() *web.JsonResult {
	roles := services.RoleService.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
	return web.JsonData(roles)
}

func (c *RoleController) GetRole_menu_ids() *web.JsonResult {
	roleId, _ := params.GetInt64(c.Ctx, "roleId")
	menuIds := services.RoleMenuService.GetMenuIdsByRole(roleId)
	return web.JsonData(menuIds)
}

func (c *RoleController) PostSave_role_menus() *web.JsonResult {
	roleId, _ := params.GetInt64(c.Ctx, "roleId")
	menuIds := params.GetInt64Arr(c.Ctx, "menuIds")
	if err := services.RoleMenuService.SaveRoleMenus(roleId, menuIds); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *RoleController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := services.RoleService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
