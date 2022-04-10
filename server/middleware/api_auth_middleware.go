package middleware

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"bbs-go/pkg/urls"
	"bbs-go/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

var (
	config = []PathRole{
		{Pattern: "/api/admin/sys-config/**", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/user/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/user/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/topic-node/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/topic-node/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/tag/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/tag/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/**", Roles: []string{constants.RoleOwner, constants.RoleAdmin}},
	}
	antPathMatcher = urls.NewAntPathMatcher()
)

// AdminAuth 后台权限
func AdminAuth(ctx iris.Context) {
	roles := getPathRoles(ctx)

	// 不需要任何角色既能访问
	if len(roles) == 0 {
		return
	}

	user := services.UserTokenService.GetCurrent(ctx)
	if user == nil {
		notLogin(ctx)
		return
	}
	if !user.HasAnyRole(roles...) {
		noPermission(ctx)
		return
	}

	ctx.Next()
}

// getPathRoles 获取请求该路径所需的角色
func getPathRoles(ctx iris.Context) []string {
	p := ctx.Path()
	for _, pathRole := range config {
		if antPathMatcher.Match(pathRole.Pattern, p) {
			return pathRole.Roles
		}
	}
	return nil
}

// notLogin 未登录返回
func notLogin(ctx iris.Context) {
	_, _ = ctx.JSON(web.JsonError(common.ErrorNotLogin))
	ctx.StopExecution()
}

// noPermission 无权限返回
func noPermission(ctx iris.Context) {
	_, _ = ctx.JSON(web.JsonErrorCode(2, "无权限"))
	ctx.StopExecution()
}

type PathRole struct {
	Pattern string   // path pattern
	Roles   []string // roles
}
