package middleware

import (
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

// AdminMiddleware 后台权限
func AdminMiddleware(ctx iris.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		notLogin(ctx)
		return
	}
	if user.IsOwner() {
		ctx.Next()
		return
	}

	permissionCodes, ok := permissions.GetAdminPermissionCodes(ctx.Method(), ctx.Path())
	if !ok || !services.PermissionService.HasAnyPermission(user, permissionCodes...) {
		noPermission(ctx)
		return
	}

	ctx.Next()
}

// notLogin 未登录返回
func notLogin(ctx iris.Context) {
	_ = ctx.JSON(web.JsonError(errs.NotLogin()))
	ctx.StopExecution()
}

// noPermission 无权限返回
func noPermission(ctx iris.Context) {
	_ = ctx.JSON(web.JsonError(errs.NoPermission()))
	ctx.StopExecution()
}
