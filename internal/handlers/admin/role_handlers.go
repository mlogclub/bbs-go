package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	modelReq "bbs-go/internal/models/req"
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
)

func roleBuildRoleItem(role *models.Role, buildPermissions bool) map[string]interface{} {
	b := web.NewRspBuilder(role)
	if buildPermissions {
		b.Put("permissionIds", services.RolePermissionService.GetRolePermissionIds(role.Id)).
			Put("permissions", services.RolePermissionService.GetRolePermissionCodes(role.Id))
	}
	return b.Build()
}

func Roles(ctx *gin.Context) {
	roles := services.RoleService.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
	ginx.WriteJSON(ctx, roles)
}

func RoleDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t := services.RoleService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.entity_not_found")+", id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, roleBuildRoleItem(t, true))
}

func RoleList(ctx *gin.Context) {
	list := services.RoleService.Find(params.NewSqlCnd(ctx,
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
	ginx.WriteJSON(ctx, list)
}

func RolePermissions(ctx *gin.Context) {
	list := services.PermissionService.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
	for i := range list {
		if definition, ok := permissions.FindByCode(list[i].Code); ok {
			list[i].GroupName = definition.GroupName
			list[i].SortNo = definition.SortNo
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].SortNo == list[j].SortNo {
			return list[i].Id > list[j].Id
		}
		return list[i].SortNo < list[j].SortNo
	})
	ginx.WriteJSON(ctx, list)
}

func RoleCreate(ctx *gin.Context) {
	t := &models.Role{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	if services.RoleService.GetByCode(t.Code) != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.role_code_exists")))
		return
	}

	t.SortNo = services.RoleService.GetNextSortNo()
	t.Type = constants.RoleTypeCustom
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.RoleService.Create(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, roleBuildRoleItem(t, true))
}

func RoleUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	t := services.RoleService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.entity_not_found")))
		return
	}

	if t.Type == constants.RoleTypeSystem {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.system_role_edit_forbidden")))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	if exists := services.RoleService.GetByCode(t.Code); exists != nil && exists.Id != t.Id {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.role_code_exists")))
		return
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.RoleService.Update(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, roleBuildRoleItem(t, true))

}

func RoleUpdatePermissions(ctx *gin.Context) {
	var req modelReq.RolePermissionsReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	t := services.RoleService.Get(req.Id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.entity_not_found")))
		return
	}
	if t.Type == constants.RoleTypeSystem || t.Code == constants.RoleOwner {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.system_role_edit_forbidden")))
		return
	}
	if err := services.RolePermissionService.UpdateRolePermissions(t.Id, modelReq.SplitCommaInt64s(req.PermissionIds)); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, roleBuildRoleItem(t, true))

}

func RoleRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.delete_ids_required")))
		return
	}
	for _, id := range ids {
		t := services.RoleService.Get(id)
		if t == nil {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.entity_not_found")))
			return
		}
		if t.Type == constants.RoleTypeSystem {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.system_role_delete_forbidden")))
			return
		}
		if services.UserRoleService.IsRoleInUse(t.Id) {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.role_in_use_delete_forbidden")))
			return
		}
		services.RoleService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": dates.NowTimestamp(),
		})
	}
	ginx.WriteJSON(ctx, nil)

}

func RoleUpdateSort(ctx *gin.Context) {
	var ids []int64
	if err := ginx.BindJSON(ctx, &ids); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if err := services.RoleService.UpdateSort(ids); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}
