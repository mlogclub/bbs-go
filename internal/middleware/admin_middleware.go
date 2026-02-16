package middleware

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/urls"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type pathRole struct {
	Pattern string   // path pattern
	Roles   []string // roles
}

var (
	authCfg = []pathRole{
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

// AdminMiddleware 后台权限
func AdminMiddleware(ctx iris.Context) {
	roles := getPathRoles(ctx)

	// 不需要任何角色既能访问
	if len(roles) == 0 {
		return
	}

	user := common.GetCurrentUser(ctx)
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
	for _, pathRole := range authCfg {
		if antPathMatcher.Match(pathRole.Pattern, p) {
			return pathRole.Roles
		}
	}
	return nil
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
