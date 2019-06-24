package middleware

import (
	"github.com/mlogclub/mlog/services/oauth"
	"github.com/mlogclub/simple"
	"strings"

	"github.com/kataras/iris/context"
)

func AdminAuthHandler(ctx context.Context) {
	if !isMatchPath(ctx, "/api/admin/") {
		ctx.Next()
		return
	}
	userInfo := oauth.GetUserInfoByRequest(ctx.Request())
	if userInfo == nil {
		_, _ = ctx.JSON(simple.ErrorCode(1, "Not Login"))
		ctx.StopExecution()
		return
	}
	if !userInfo.HasRole("管理员") {
		_, _ = ctx.JSON(simple.ErrorCode(2, "无权限"))
		ctx.StopExecution()
		return
	}

	ctx.Next()
	return
}

func isMatchPath(ctx context.Context, pattern string) bool {
	path := ctx.Path()
	if strings.HasPrefix(path, pattern) {
		return true
	}
	return false
}
