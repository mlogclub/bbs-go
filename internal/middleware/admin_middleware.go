package middleware

import (
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/services"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 后台权限
func AdminMiddleware(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		notLogin(ctx)
		return
	}
	if user.IsOwner() {
		ctx.Next()
		return
	}

	permissionCodes, ok := permissions.GetAdminPermissionCodes(ctx.Request.Method, ctx.Request.URL.Path)
	if !ok || !services.PermissionService.HasAnyPermission(user, permissionCodes...) {
		noPermission(ctx)
		return
	}

	ctx.Next()
}

// notLogin 未登录返回
func notLogin(ctx *gin.Context) {
	ginx.WriteJSON(ctx, errs.NotLogin())
	ctx.Abort()
}

// noPermission 无权限返回
func noPermission(ctx *gin.Context) {
	ginx.WriteJSON(ctx, errs.NoPermission())
	ctx.Abort()
}
